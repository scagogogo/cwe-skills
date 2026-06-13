package cwe

import (
	"fmt"
)

// TreeNode 表示CWE树中的一个节点。
//
// TreeNode包装了CWE条目，添加了树结构信息（父节点、子节点、深度），
// 用于可视化展示和树形遍历。
type TreeNode struct {
	// CWE 当前节点对应的CWE条目
	CWE *CWE
	// Children 子节点列表
	Children []*TreeNode
	// Parent 父节点
	Parent *TreeNode
	// Depth 节点深度（根节点为0）
	Depth int
}

// NewTreeNode 创建一个新的TreeNode。
//
// 参数：
//   - cwe: 节点对应的CWE条目
//
// 返回值：
//   - *TreeNode: 新创建的树节点
func NewTreeNode(cwe *CWE) *TreeNode {
	return &TreeNode{
		CWE:      cwe,
		Children: make([]*TreeNode, 0),
		Depth:    0,
	}
}

// AddChild 向当前节点添加一个子节点。
//
// 该方法会自动设置子节点的Parent和Depth。
//
// 参数：
//   - child: 要添加的子节点
func (n *TreeNode) AddChild(child *TreeNode) {
	child.Parent = n
	child.Depth = n.Depth + 1
	n.Children = append(n.Children, child)
}

// Walk 使用深度优先搜索遍历树。
//
// 遍历过程中对每个节点调用fn函数。
// 如果fn返回false，则停止遍历该分支。
//
// 参数：
//   - fn: 对每个节点调用的函数，返回false停止遍历
func (n *TreeNode) Walk(fn func(*TreeNode) bool) {
	if n == nil {
		return
	}
	if !fn(n) {
		return
	}
	for _, child := range n.Children {
		child.Walk(fn)
	}
}

// WalkBFS 使用广度优先搜索遍历树。
//
// 参数：
//   - fn: 对每个节点调用的函数，返回false停止遍历
func (n *TreeNode) WalkBFS(fn func(*TreeNode) bool) {
	if n == nil {
		return
	}

	queue := []*TreeNode{n}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if !fn(current) {
			return
		}
		queue = append(queue, current.Children...)
	}
}

// Find 在树中查找指定ID的节点。
//
// 参数：
//   - id: 要查找的CWE ID
//
// 返回值：
//   - *TreeNode: 找到的节点，如果未找到返回nil
func (n *TreeNode) Find(id int) *TreeNode {
	if n == nil {
		return nil
	}

	if n.CWE != nil && n.CWE.ID == id {
		return n
	}

	for _, child := range n.Children {
		if found := child.Find(id); found != nil {
			return found
		}
	}

	return nil
}

// Path 返回从根节点到当前节点的路径。
//
// 返回值：
//   - []*TreeNode: 路径上的节点列表，从根节点到当前节点
func (n *TreeNode) Path() []*TreeNode {
	if n == nil {
		return nil
	}

	var path []*TreeNode
	current := n
	for current != nil {
		path = append([]*TreeNode{current}, path...)
		current = current.Parent
	}

	return path
}

// LeafNodes 获取当前节点下所有的叶子节点。
//
// 返回值：
//   - []*TreeNode: 所有叶子节点列表
func (n *TreeNode) LeafNodes() []*TreeNode {
	if n == nil {
		return nil
	}

	var leaves []*TreeNode
	n.Walk(func(node *TreeNode) bool {
		if len(node.Children) == 0 {
			leaves = append(leaves, node)
		}
		return true
	})

	return leaves
}

// MaxDepth 获取当前节点下子树的最大深度。
//
// 返回值：
//   - int: 最大深度
func (n *TreeNode) MaxDepth() int {
	if n == nil {
		return 0
	}

	maxDepth := n.Depth
	n.Walk(func(node *TreeNode) bool {
		if node.Depth > maxDepth {
			maxDepth = node.Depth
		}
		return true
	})

	return maxDepth
}

// Count 获取当前节点下子树的节点总数（包括自身）。
//
// 返回值：
//   - int: 节点总数
func (n *TreeNode) Count() int {
	if n == nil {
		return 0
	}

	count := 0
	n.Walk(func(node *TreeNode) bool {
		count++
		return true
	})

	return count
}

// IsLeaf 检查当前节点是否为叶子节点。
func (n *TreeNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// IsRoot 检查当前节点是否为根节点。
func (n *TreeNode) IsRoot() bool {
	return n.Parent == nil
}

// String 返回节点的字符串表示。
func (n *TreeNode) String() string {
	if n == nil || n.CWE == nil {
		return "TreeNode<nil>"
	}
	return fmt.Sprintf("TreeNode<CWE-%d %s depth=%d children=%d>",
		n.CWE.ID, n.CWE.Name, n.Depth, len(n.Children))
}

// BuildTree 从注册表构建以指定ID为根的树。
//
// 该方法递归地构建CWE条目的树结构。
// 需要先调用Registry.BuildIndexes构建索引。
//
// 参数：
//   - r: CWE注册表
//   - rootID: 根节点的CWE ID
//
// 返回值：
//   - *TreeNode: 构建的树根节点，如果根ID不存在返回nil
func BuildTree(r *Registry, rootID int) *TreeNode {
	if r == nil {
		return nil
	}

	root, ok := r.Get(rootID)
	if !ok {
		return nil
	}

	rootNode := NewTreeNode(root)
	buildTreeNode(r, rootNode, map[int]bool{rootID: true})

	return rootNode
}

// BuildForest 从注册表构建森林（所有顶层Pillar节点的树）。
//
// 参数：
//   - r: CWE注册表
//
// 返回值：
//   - []*TreeNode: 所有顶层节点的树列表
func BuildForest(r *Registry) []*TreeNode {
	if r == nil {
		return nil
	}

	pillars := FindTopLevel(r)
	var forest []*TreeNode

	for _, pillar := range pillars {
		tree := BuildTree(r, pillar.ID)
		if tree != nil {
			forest = append(forest, tree)
		}
	}

	return forest
}

// BuildViewTree 从注册表构建指定视图的树。
//
// 参数：
//   - r: CWE注册表
//   - viewID: 视图ID
//
// 返回值：
//   - *TreeNode: 视图的根节点树
func BuildViewTree(r *Registry, viewID int) *TreeNode {
	if r == nil {
		return nil
	}

	view, ok := r.GetView(viewID)
	if !ok {
		return nil
	}

	// 创建虚拟根节点代表视图
	viewCWE := &CWE{
		ID:          view.ID,
		Name:        view.Name,
		Description: view.Description,
		CWEType:     "view",
	}
	rootNode := NewTreeNode(viewCWE)

	// 添加视图成员作为子节点
	visited := map[int]bool{viewID: true}
	for _, member := range view.Members {
		if visited[member.CWEID] {
			continue
		}
		if cwe, ok := r.Get(member.CWEID); ok {
			childNode := NewTreeNode(cwe)
			buildTreeNode(r, childNode, visited)
			rootNode.AddChild(childNode)
		}
	}

	return rootNode
}

// buildTreeNode 递归构建树节点。
func buildTreeNode(r *Registry, node *TreeNode, visited map[int]bool) {
	childIDs := r.GetChildIDs(node.CWE.ID)

	for _, childID := range childIDs {
		if visited[childID] {
			continue
		}
		visited[childID] = true

		child, ok := r.Get(childID)
		if !ok {
			continue
		}

		childNode := NewTreeNode(child)
		node.AddChild(childNode)
		buildTreeNode(r, childNode, visited)
	}
}
