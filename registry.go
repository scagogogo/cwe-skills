package cweskills

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Registry 是CWE条目的内存注册表。
//
// Registry提供了CWE条目的存储、索引和查询功能。
// 它是本地数据管理的核心组件，支持：
//   - 注册和查找CWE弱点、类别、视图、复合元素
//   - 构建关系索引（父子关系、对等关系等）
//   - 导出和导入JSON数据
//   - 并发安全的读写操作
//
// 使用场景：
//   - 从XML目录加载离线CWE数据
//   - 缓存API查询结果
//   - 构建CWE关系图
//
// 示例：
//
//	registry := cwe.NewRegistry()
//	registry.Register(&cwe.CWE{ID: 79, Name: "XSS", Abstraction: cwe.AbstractionBase})
//	cwe79, found := registry.Get(79)
type Registry struct {
	mu               sync.RWMutex
	weaknesses       map[int]*CWE
	categories       map[int]*Category
	views            map[int]*View
	compoundElements map[int]*CompoundElement
	// 关系索引
	parentIndex  map[int][]int // 子ID -> 父ID列表
	childIndex   map[int][]int // 父ID -> 子ID列表
	peerIndex    map[int][]int // ID -> 对等ID列表
	memberIndex  map[int][]int // 类别/视图ID -> 成员ID列表
	memberOfIndex map[int][]int // ID -> 所属类别/视图ID列表
	// 索引是否已构建
	indexesBuilt bool
}

// NewRegistry 创建一个新的空Registry实例。
func NewRegistry() *Registry {
	return &Registry{
		weaknesses:       make(map[int]*CWE),
		categories:       make(map[int]*Category),
		views:            make(map[int]*View),
		compoundElements: make(map[int]*CompoundElement),
		parentIndex:      make(map[int][]int),
		childIndex:       make(map[int][]int),
		peerIndex:        make(map[int][]int),
		memberIndex:      make(map[int][]int),
		memberOfIndex:    make(map[int][]int),
	}
}

// Register 向注册表中注册一个CWE弱点条目。
//
// 如果ID已存在，返回ValidationError。
//
// 参数：
//   - cwe: 要注册的CWE条目
//
// 返回值：
//   - error: 注册失败时返回错误
func (r *Registry) Register(cwe *CWE) error {
	if cwe == nil {
		return NewValidationError("CWE", "nil")
	}
	if cwe.ID <= 0 {
		return NewValidationError("ID", fmt.Sprintf("%d", cwe.ID))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.weaknesses[cwe.ID]; exists {
		return NewValidationError("ID", fmt.Sprintf("CWE-%d 已存在", cwe.ID))
	}

	r.weaknesses[cwe.ID] = cwe
	r.indexesBuilt = false
	return nil
}

// RegisterCategory 向注册表中注册一个CWE类别。
func (r *Registry) RegisterCategory(cat *Category) error {
	if cat == nil {
		return NewValidationError("Category", "nil")
	}
	if cat.ID <= 0 {
		return NewValidationError("ID", fmt.Sprintf("%d", cat.ID))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[cat.ID]; exists {
		return NewValidationError("ID", fmt.Sprintf("Category-%d 已存在", cat.ID))
	}

	r.categories[cat.ID] = cat
	r.indexesBuilt = false
	return nil
}

// RegisterView 向注册表中注册一个CWE视图。
func (r *Registry) RegisterView(view *View) error {
	if view == nil {
		return NewValidationError("View", "nil")
	}
	if view.ID <= 0 {
		return NewValidationError("ID", fmt.Sprintf("%d", view.ID))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.views[view.ID]; exists {
		return NewValidationError("ID", fmt.Sprintf("View-%d 已存在", view.ID))
	}

	r.views[view.ID] = view
	r.indexesBuilt = false
	return nil
}

// RegisterCompoundElement 向注册表中注册一个复合元素。
func (r *Registry) RegisterCompoundElement(ce *CompoundElement) error {
	if ce == nil {
		return NewValidationError("CompoundElement", "nil")
	}
	if ce.ID <= 0 {
		return NewValidationError("ID", fmt.Sprintf("%d", ce.ID))
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.compoundElements[ce.ID]; exists {
		return NewValidationError("ID", fmt.Sprintf("CompoundElement-%d 已存在", ce.ID))
	}

	r.compoundElements[ce.ID] = ce
	r.indexesBuilt = false
	return nil
}

// Get 根据ID获取CWE弱点条目。
func (r *Registry) Get(id int) (*CWE, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cwe, ok := r.weaknesses[id]
	return cwe, ok
}

// GetCategory 根据ID获取CWE类别。
func (r *Registry) GetCategory(id int) (*Category, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cat, ok := r.categories[id]
	return cat, ok
}

// GetView 根据ID获取CWE视图。
func (r *Registry) GetView(id int) (*View, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	view, ok := r.views[id]
	return view, ok
}

// GetCompoundElement 根据ID获取复合元素。
func (r *Registry) GetCompoundElement(id int) (*CompoundElement, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ce, ok := r.compoundElements[id]
	return ce, ok
}

// GetAll 获取所有CWE弱点条目。
func (r *Registry) GetAll() []*CWE {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*CWE, 0, len(r.weaknesses))
	for _, cwe := range r.weaknesses {
		result = append(result, cwe)
	}
	return result
}

// GetAllCategories 获取所有CWE类别。
func (r *Registry) GetAllCategories() []*Category {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Category, 0, len(r.categories))
	for _, cat := range r.categories {
		result = append(result, cat)
	}
	return result
}

// GetAllViews 获取所有CWE视图。
func (r *Registry) GetAllViews() []*View {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*View, 0, len(r.views))
	for _, view := range r.views {
		result = append(result, view)
	}
	return result
}

// Size 获取注册表中CWE弱点条目的数量。
func (r *Registry) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.weaknesses)
}

// CategoryCount 获取注册表中类别的数量。
func (r *Registry) CategoryCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.categories)
}

// ViewCount 获取注册表中视图的数量。
func (r *Registry) ViewCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.views)
}

// CompoundElementCount 获取注册表中复合元素的数量。
func (r *Registry) CompoundElementCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.compoundElements)
}

// Contains 检查注册表中是否包含指定ID的CWE条目。
func (r *Registry) Contains(id int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.weaknesses[id]
	return ok
}

// Remove 从注册表中移除指定ID的CWE条目。
func (r *Registry) Remove(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.weaknesses[id]; !ok {
		return NewCWENotFoundError(id)
	}

	delete(r.weaknesses, id)
	r.indexesBuilt = false
	return nil
}

// RemoveCategory 从注册表中移除指定ID的类别。
func (r *Registry) RemoveCategory(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.categories[id]; !ok {
		return NewCWENotFoundError(id)
	}

	delete(r.categories, id)
	r.indexesBuilt = false
	return nil
}

// RemoveView 从注册表中移除指定ID的视图。
func (r *Registry) RemoveView(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.views[id]; !ok {
		return NewCWENotFoundError(id)
	}

	delete(r.views, id)
	r.indexesBuilt = false
	return nil
}

// Clear 清空注册表中的所有数据。
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.weaknesses = make(map[int]*CWE)
	r.categories = make(map[int]*Category)
	r.views = make(map[int]*View)
	r.compoundElements = make(map[int]*CompoundElement)
	r.parentIndex = make(map[int][]int)
	r.childIndex = make(map[int][]int)
	r.peerIndex = make(map[int][]int)
	r.memberIndex = make(map[int][]int)
	r.memberOfIndex = make(map[int][]int)
	r.indexesBuilt = false
}

// BuildIndexes 构建关系索引。
//
// 该方法遍历所有CWE条目的关系数据，构建以下索引：
//   - 父子关系索引（基于ChildOf/ParentOf关系）
//   - 对等关系索引（基于PeerOf/CanAlsoBe关系）
//   - 成员关系索引（基于类别/视图的成员列表）
//
// 构建索引后，GetParentIDs、GetChildIDs、GetAncestorIDs等方法
// 可以快速返回结果。建议在数据加载完成后调用此方法。
func (r *Registry) BuildIndexes() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 清除现有索引
	r.parentIndex = make(map[int][]int)
	r.childIndex = make(map[int][]int)
	r.peerIndex = make(map[int][]int)
	r.memberIndex = make(map[int][]int)
	r.memberOfIndex = make(map[int][]int)

	// 从弱点关系构建索引
	for _, cwe := range r.weaknesses {
		for _, rel := range cwe.Relationships {
			switch rel.Nature {
			case RelationshipChildOf:
				r.parentIndex[cwe.ID] = append(r.parentIndex[cwe.ID], rel.CWEID)
				r.childIndex[rel.CWEID] = append(r.childIndex[rel.CWEID], cwe.ID)
			case RelationshipParentOf:
				r.childIndex[cwe.ID] = append(r.childIndex[cwe.ID], rel.CWEID)
				r.parentIndex[rel.CWEID] = append(r.parentIndex[rel.CWEID], cwe.ID)
			case RelationshipPeerOf, RelationshipCanAlsoBe:
				r.peerIndex[cwe.ID] = append(r.peerIndex[cwe.ID], rel.CWEID)
			case RelationshipMemberOf:
				r.memberOfIndex[cwe.ID] = append(r.memberOfIndex[cwe.ID], rel.CWEID)
				r.memberIndex[rel.CWEID] = append(r.memberIndex[rel.CWEID], cwe.ID)
			case RelationshipHasMember:
				r.memberIndex[cwe.ID] = append(r.memberIndex[cwe.ID], rel.CWEID)
				r.memberOfIndex[rel.CWEID] = append(r.memberOfIndex[rel.CWEID], cwe.ID)
			}
		}
	}

	// 从类别关系构建成员索引
	for _, cat := range r.categories {
		for _, rel := range cat.Relationships {
			switch rel.Nature {
			case RelationshipHasMember:
				r.memberIndex[cat.ID] = append(r.memberIndex[cat.ID], rel.CWEID)
				r.memberOfIndex[rel.CWEID] = append(r.memberOfIndex[rel.CWEID], cat.ID)
			case RelationshipMemberOf:
				r.memberOfIndex[cat.ID] = append(r.memberOfIndex[cat.ID], rel.CWEID)
				r.memberIndex[rel.CWEID] = append(r.memberIndex[rel.CWEID], cat.ID)
			}
		}
	}

	// 从视图成员构建成员索引
	for _, view := range r.views {
		for _, member := range view.Members {
			r.memberIndex[view.ID] = append(r.memberIndex[view.ID], member.CWEID)
			r.memberOfIndex[member.CWEID] = append(r.memberOfIndex[member.CWEID], view.ID)
		}
	}

	// 去重索引中的重复ID
	r.parentIndex = dedupIntMap(r.parentIndex)
	r.childIndex = dedupIntMap(r.childIndex)
	r.peerIndex = dedupIntMap(r.peerIndex)
	r.memberIndex = dedupIntMap(r.memberIndex)
	r.memberOfIndex = dedupIntMap(r.memberOfIndex)

	r.indexesBuilt = true
}

// IndexesBuilt 检查索引是否已构建。
func (r *Registry) IndexesBuilt() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.indexesBuilt
}

// GetParentIDs 获取指定CWE条目的父级ID列表。
//
// 需要先调用BuildIndexes构建索引。
func (r *Registry) GetParentIDs(id int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return copyIntSlice(r.parentIndex[id])
}

// GetChildIDs 获取指定CWE条目的子级ID列表。
//
// 需要先调用BuildIndexes构建索引。
func (r *Registry) GetChildIDs(id int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return copyIntSlice(r.childIndex[id])
}

// GetPeerIDs 获取指定CWE条目的对等ID列表。
//
// 需要先调用BuildIndexes构建索引。
func (r *Registry) GetPeerIDs(id int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return copyIntSlice(r.peerIndex[id])
}

// GetAncestorIDs 递归获取指定CWE条目的所有祖先ID。
//
// 需要先调用BuildIndexes构建索引。
// 该方法使用广度优先搜索遍历父级关系，避免循环引用导致的无限递归。
func (r *Registry) GetAncestorIDs(id int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	visited := make(map[int]bool)
	var result []int
	queue := []int{id}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		parents := r.parentIndex[current]
		for _, parentID := range parents {
			if !visited[parentID] {
				visited[parentID] = true
				result = append(result, parentID)
				queue = append(queue, parentID)
			}
		}
	}

	return result
}

// GetDescendantIDs 递归获取指定CWE条目的所有后代ID。
//
// 需要先调用BuildIndexes构建索引。
// 该方法使用广度优先搜索遍历子级关系，避免循环引用导致的无限递归。
func (r *Registry) GetDescendantIDs(id int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	visited := make(map[int]bool)
	var result []int
	queue := []int{id}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		children := r.childIndex[current]
		for _, childID := range children {
			if !visited[childID] {
				visited[childID] = true
				result = append(result, childID)
				queue = append(queue, childID)
			}
		}
	}

	return result
}

// GetViewMembers 获取指定视图的成员ID列表。
func (r *Registry) GetViewMembers(viewID int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return copyIntSlice(r.memberIndex[viewID])
}

// GetCategoryMembers 获取指定类别的成员ID列表。
func (r *Registry) GetCategoryMembers(categoryID int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return copyIntSlice(r.memberIndex[categoryID])
}

// GetMemberOfIDs 获取指定CWE条目所属的类别/视图ID列表。
func (r *Registry) GetMemberOfIDs(id int) []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return copyIntSlice(r.memberOfIndex[id])
}

// ExportJSON 将注册表数据导出为JSON格式。
func (r *Registry) ExportJSON() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := struct {
		Weaknesses       []*CWE            `json:"weaknesses"`
		Categories       []*Category       `json:"categories"`
		Views            []*View           `json:"views"`
		CompoundElements []*CompoundElement `json:"compound_elements"`
	}{
		Weaknesses:       make([]*CWE, 0, len(r.weaknesses)),
		Categories:       make([]*Category, 0, len(r.categories)),
		Views:            make([]*View, 0, len(r.views)),
		CompoundElements: make([]*CompoundElement, 0, len(r.compoundElements)),
	}

	for _, cwe := range r.weaknesses {
		data.Weaknesses = append(data.Weaknesses, cwe)
	}
	for _, cat := range r.categories {
		data.Categories = append(data.Categories, cat)
	}
	for _, view := range r.views {
		data.Views = append(data.Views, view)
	}
	for _, ce := range r.compoundElements {
		data.CompoundElements = append(data.CompoundElements, ce)
	}

	return json.MarshalIndent(data, "", "  ")
}

// ImportJSON 从JSON格式导入数据到注册表。
func (r *Registry) ImportJSON(jsonData []byte) error {
	var data struct {
		Weaknesses       []*CWE            `json:"weaknesses"`
		Categories       []*Category       `json:"categories"`
		Views            []*View           `json:"views"`
		CompoundElements []*CompoundElement `json:"compound_elements"`
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return NewParseError(fmt.Sprintf("JSON解析失败: %v", err), 0)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, cwe := range data.Weaknesses {
		if cwe != nil {
			r.weaknesses[cwe.ID] = cwe
		}
	}
	for _, cat := range data.Categories {
		if cat != nil {
			r.categories[cat.ID] = cat
		}
	}
	for _, view := range data.Views {
		if view != nil {
			r.views[view.ID] = view
		}
	}
	for _, ce := range data.CompoundElements {
		if ce != nil {
			r.compoundElements[ce.ID] = ce
		}
	}

	r.indexesBuilt = false
	return nil
}

// copyIntSlice 复制整数切片，避免外部修改影响内部数据。
func copyIntSlice(src []int) []int {
	if src == nil {
		return nil
	}
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

// dedupIntMap 对map[int][]int中的每个切片去重。
func dedupIntMap(m map[int][]int) map[int][]int {
	result := make(map[int][]int, len(m))
	for k, v := range m {
		seen := make(map[int]bool, len(v))
		deduped := make([]int, 0, len(v))
		for _, id := range v {
			if !seen[id] {
				seen[id] = true
				deduped = append(deduped, id)
			}
		}
		result[k] = deduped
	}
	return result
}
