package cweskills

import "context"

// DataFetcher 定义CWE数据获取的接口。
//
// 该接口抽象了CWE数据的获取方式，允许上层应用
// 根据需要选择API获取或本地注册表查找等不同实现。
type DataFetcher interface {
	// Fetch 获取指定ID的CWE条目
	Fetch(ctx context.Context, id int) (*CWE, error)
}

// BasicFetcher 通过API获取单个CWE条目的获取器。
//
// 示例：
//
//	client := cwe.NewAPIClient()
//	fetcher := cwe.NewBasicFetcher(client)
//	cwe79, err := fetcher.Fetch(context.Background(), 79)
type BasicFetcher struct {
	client *APIClient
}

// NewBasicFetcher 创建一个新的基础获取器。
//
// 参数：
//   - client: API客户端实例
func NewBasicFetcher(client *APIClient) *BasicFetcher {
	if client == nil {
		client = NewAPIClient()
	}
	return &BasicFetcher{client: client}
}

// Fetch 通过API获取指定ID的CWE弱点条目。
func (f *BasicFetcher) Fetch(ctx context.Context, id int) (*CWE, error) {
	return f.client.GetWeakness(ctx, id)
}

// FetchCategory 通过API获取指定ID的CWE类别。
func (f *BasicFetcher) FetchCategory(ctx context.Context, id int) (*Category, error) {
	return f.client.GetCategory(ctx, id)
}

// FetchView 通过API获取指定ID的CWE视图。
func (f *BasicFetcher) FetchView(ctx context.Context, id int) (*View, error) {
	return f.client.GetView(ctx, id)
}

// FetchWithRelations 通过API获取指定ID的CWE条目及其关系。
//
// 该方法不仅获取CWE条目本身，还获取其父级和子级关系信息，
// 并将关系信息填充到CWE条目的Relationships字段中。
func (f *BasicFetcher) FetchWithRelations(ctx context.Context, id int, viewID ...int) (*CWE, error) {
	cwe, err := f.client.GetWeakness(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取父级关系
	parents, err := f.client.GetParents(ctx, id, viewID...)
	if err == nil {
		for _, rel := range parents {
			rel.Nature = RelationshipChildOf
			cwe.Relationships = append(cwe.Relationships, rel)
		}
	}

	// 获取子级关系
	children, err := f.client.GetChildren(ctx, id, viewID...)
	if err == nil {
		for _, rel := range children {
			rel.Nature = RelationshipParentOf
			cwe.Relationships = append(cwe.Relationships, rel)
		}
	}

	return cwe, nil
}

// MultipleFetcher 批量获取CWE条目的获取器。
type MultipleFetcher struct {
	client *APIClient
}

// NewMultipleFetcher 创建一个新的批量获取器。
func NewMultipleFetcher(client *APIClient) *MultipleFetcher {
	if client == nil {
		client = NewAPIClient()
	}
	return &MultipleFetcher{client: client}
}

// FetchMultiple 批量获取多个CWE弱点条目。
//
// 参数：
//   - ctx: 请求上下文
//   - ids: 要获取的CWE ID列表
//
// 返回值：
//   - map[string]*CWE: 以CWE ID字符串为键的弱点映射
//   - error: 获取失败时返回错误
func (f *MultipleFetcher) FetchMultiple(ctx context.Context, ids []int) (map[string]*CWE, error) {
	return f.client.GetCWEs(ctx, ids)
}

// FetchMultipleToRegistry 批量获取CWE弱点条目并注册到Registry中。
//
// 参数：
//   - ctx: 请求上下文
//   - ids: 要获取的CWE ID列表
//   - registry: 目标注册表
//
// 返回值：
//   - error: 获取或注册失败时返回错误
func (f *MultipleFetcher) FetchMultipleToRegistry(ctx context.Context, ids []int, registry *Registry) error {
	if registry == nil {
		return NewValidationError("registry", "nil")
	}

	result, err := f.client.GetCWEs(ctx, ids)
	if err != nil {
		return err
	}

	for _, cwe := range result {
		if cwe != nil {
			_ = registry.Register(cwe) // 忽略重复注册错误
		}
	}

	return nil
}

// TreeFetcher 递归获取CWE树结构的获取器。
type TreeFetcher struct {
	client   *APIClient
	registry *Registry
	maxDepth int
}

// NewTreeFetcher 创建一个新的树获取器。
//
// 参数：
//   - client: API客户端实例
//   - registry: 用于存储获取结果的注册表
//   - maxDepth: 最大递归深度，0表示无限制
func NewTreeFetcher(client *APIClient, registry *Registry, maxDepth int) *TreeFetcher {
	if client == nil {
		client = NewAPIClient()
	}
	if registry == nil {
		registry = NewRegistry()
	}
	if maxDepth <= 0 {
		maxDepth = 10 // 默认最大深度
	}
	return &TreeFetcher{
		client:   client,
		registry: registry,
		maxDepth: maxDepth,
	}
}

// FetchWithAncestors 获取指定CWE条目及其所有祖先。
//
// 该方法递归地获取CWE条目的父级，直到达到根节点或最大深度。
func (f *TreeFetcher) FetchWithAncestors(ctx context.Context, id int) error {
	return f.fetchAncestorsRecursive(ctx, id, 0)
}

// FetchWithDescendants 获取指定CWE条目及其所有后代。
//
// 该方法递归地获取CWE条目的子级，直到没有子级或达到最大深度。
func (f *TreeFetcher) FetchWithDescendants(ctx context.Context, id int) error {
	return f.fetchDescendantsRecursive(ctx, id, 0)
}

// FetchFullTree 获取以指定ID为根的完整CWE树。
//
// 该方法同时获取祖先和后代，构建完整的CWE关系树。
func (f *TreeFetcher) FetchFullTree(ctx context.Context, rootID int) error {
	if err := f.FetchWithAncestors(ctx, rootID); err != nil {
		return err
	}
	return f.FetchWithDescendants(ctx, rootID)
}

// GetRegistry 获取用于存储结果的注册表。
func (f *TreeFetcher) GetRegistry() *Registry {
	return f.registry
}

// fetchAncestorsRecursive 递归获取祖先。
func (f *TreeFetcher) fetchAncestorsRecursive(ctx context.Context, id int, depth int) error {
	if depth >= f.maxDepth {
		return nil
	}

	// 如果已经获取过，跳过
	if f.registry.Contains(id) {
		return nil
	}

	// 获取当前CWE条目
	cwe, err := f.client.GetWeakness(ctx, id)
	if err != nil {
		return err
	}

	_ = f.registry.Register(cwe)

	// 获取父级关系
	parents, err := f.client.GetParents(ctx, id)
	if err != nil {
		return err
	}

	for _, rel := range parents {
		rel.Nature = RelationshipChildOf
		cwe.Relationships = append(cwe.Relationships, rel)

		// 递归获取父级
		if err := f.fetchAncestorsRecursive(ctx, rel.CWEID, depth+1); err != nil {
			return err
		}
	}

	return nil
}

// fetchDescendantsRecursive 递归获取后代。
func (f *TreeFetcher) fetchDescendantsRecursive(ctx context.Context, id int, depth int) error {
	if depth >= f.maxDepth {
		return nil
	}

	// 获取当前CWE条目（如果还没有的话）
	if !f.registry.Contains(id) {
		cwe, err := f.client.GetWeakness(ctx, id)
		if err != nil {
			return err
		}
		_ = f.registry.Register(cwe)
	}

	// 获取子级关系
	children, err := f.client.GetChildren(ctx, id)
	if err != nil {
		return err
	}

	// 获取当前CWE并添加关系
	if cwe, ok := f.registry.Get(id); ok {
		for _, rel := range children {
			rel.Nature = RelationshipParentOf
			cwe.Relationships = append(cwe.Relationships, rel)
		}
	}

	// 递归获取子级
	for _, rel := range children {
		if err := f.fetchDescendantsRecursive(ctx, rel.CWEID, depth+1); err != nil {
			return err
		}
	}

	return nil
}
