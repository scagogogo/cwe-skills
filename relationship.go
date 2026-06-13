package cwe

import "fmt"

// Relationship 表示CWE条目之间的关系。
//
// CWE条目之间的关系描述了不同弱点之间的层级、顺序、依赖和对等联系。
// 关系类型由Nature字段决定，CWEID指定关系的目标弱点。
type Relationship struct {
	// Nature 关系类型，如ChildOf、ParentOf、CanPrecede等
	Nature RelationshipNature `json:"nature" xml:"Nature,attr"`
	// CWEID 关系目标的CWE ID
	CWEID int `json:"cwe_id" xml:"CWE_ID"`
	// ViewID 关系所属的视图ID，可选
	ViewID int `json:"view_id,omitempty" xml:"View_ID,omitempty"`
	// Ordinal 关系的序数标记，可选，如"Primary"
	Ordinal string `json:"ordinal,omitempty" xml:"Ordinal,omitempty"`
	// ChainID 关系所属的链ID，可选
	ChainID int `json:"chain_id,omitempty" xml:"Chain_ID,omitempty"`
}

// IsHierarchical 检查关系是否为层级关系。
//
// 委托给Nature字段的IsHierarchical方法。
// 层级关系包括ChildOf、ParentOf、MemberOf、HasMember。
//
// 返回值：
//   - bool: 如果是层级关系返回true，否则返回false
func (r *Relationship) IsHierarchical() bool {
	return r.Nature.IsHierarchical()
}

// IsSequential 检查关系是否为顺序关系。
//
// 委托给Nature字段的IsSequential方法。
// 顺序关系包括CanPrecede、CanFollow。
//
// 返回值：
//   - bool: 如果是顺序关系返回true，否则返回false
func (r *Relationship) IsSequential() bool {
	return r.Nature.IsSequential()
}

// IsDependency 检查关系是否为依赖关系。
//
// 委托给Nature字段的IsDependency方法。
// 依赖关系包括Requires、RequiredBy。
//
// 返回值：
//   - bool: 如果是依赖关系返回true，否则返回false
func (r *Relationship) IsDependency() bool {
	return r.Nature.IsDependency()
}

// IsPeer 检查关系是否为对等关系。
//
// 姜托给Nature字段的IsPeer方法。
// 对等关系包括PeerOf、CanAlsoBe。
//
// 返回值：
//   - bool: 如果是对等关系返回true，否则返回false
func (r *Relationship) IsPeer() bool {
	return r.Nature.IsPeer()
}

// IsPrimary 检查关系是否标记为主要关系。
//
// 当Ordinal字段等于"Primary"时返回true。
//
// 返回值：
//   - bool: 如果Ordinal为"Primary"返回true，否则返回false
func (r *Relationship) IsPrimary() bool {
	return r.Ordinal == "Primary"
}

// Validate 验证关系的有效性。
//
// 检查条件：
//   - Nature必须是有效的关系类型
//   - CWEID必须大于0
//
// 返回值：
//   - error: 如果验证失败返回ValidationError，否则返回nil
func (r *Relationship) Validate() error {
	if !r.Nature.IsValid() {
		return NewValidationError("Nature", string(r.Nature))
	}
	if r.CWEID <= 0 {
		return NewValidationError("CWEID", fmt.Sprintf("%d", r.CWEID))
	}
	return nil
}

// NewRelationship 创建一个新的关系。
//
// 参数：
//   - nature: 关系类型
//   - cweID: 目标CWE ID
//
// 返回值：
//   - *Relationship: 新创建的关系实例
func NewRelationship(nature RelationshipNature, cweID int) *Relationship {
	return &Relationship{
		Nature: nature,
		CWEID:  cweID,
	}
}

// NewRelationshipWithView 创建一个带有视图ID的新关系。
//
// 参数：
//   - nature: 关系类型
//   - cweID: 目标CWE ID
//   - viewID: 所属视图ID
//
// 返回值：
//   - *Relationship: 新创建的关系实例
func NewRelationshipWithView(nature RelationshipNature, cweID, viewID int) *Relationship {
	return &Relationship{
		Nature: nature,
		CWEID:  cweID,
		ViewID: viewID,
	}
}