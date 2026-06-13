package cwe

import (
	"fmt"
)

// Navigator 提供CWE条目之间的关系导航功能。
//
// Navigator基于Registry中的数据，提供丰富的关系遍历能力，
// 包括层级导航、顺序导航、依赖导航和对等导航等。
//
// 使用场景：
//   - 分析弱点的上游和下游关系
//   - 查找链式攻击路径
//   - 识别复合弱点的组成要素
//   - 计算弱点之间的关系距离
//
// 示例：
//
//	registry := cwe.NewRegistry()
//	// ... 注册CWE数据并构建索引 ...
//	nav := cwe.NewNavigator(registry)
//	parents := nav.Parents(79)
//	ancestors := nav.Ancestors(79)
type Navigator struct {
	registry *Registry
}

// NewNavigator 创建一个新的关系导航器。
//
// 参数：
//   - r: 已构建索引的CWE注册表
func NewNavigator(r *Registry) *Navigator {
	return &Navigator{registry: r}
}

// Parents 获取指定CWE条目的直接父级条目列表。
func (n *Navigator) Parents(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	parentIDs := n.registry.GetParentIDs(id)
	return n.resolveIDs(parentIDs)
}

// Children 获取指定CWE条目的直接子级条目列表。
func (n *Navigator) Children(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	childIDs := n.registry.GetChildIDs(id)
	return n.resolveIDs(childIDs)
}

// Ancestors 递归获取指定CWE条目的所有祖先条目。
func (n *Navigator) Ancestors(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	ancestorIDs := n.registry.GetAncestorIDs(id)
	return n.resolveIDs(ancestorIDs)
}

// Descendants 递归获取指定CWE条目的所有后代条目。
func (n *Navigator) Descendants(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	descendantIDs := n.registry.GetDescendantIDs(id)
	return n.resolveIDs(descendantIDs)
}

// Siblings 获取指定CWE条目的兄弟条目（具有相同父级的条目）。
func (n *Navigator) Siblings(id int) []*CWE {
	if n.registry == nil {
		return nil
	}

	parentIDs := n.registry.GetParentIDs(id)
	if len(parentIDs) == 0 {
		return nil
	}

	seen := map[int]bool{id: true}
	var result []*CWE

	for _, parentID := range parentIDs {
		for _, siblingID := range n.registry.GetChildIDs(parentID) {
			if !seen[siblingID] {
				seen[siblingID] = true
				if cwe, ok := n.registry.Get(siblingID); ok {
					result = append(result, cwe)
				}
			}
		}
	}

	return result
}

// Peers 获取指定CWE条目的对等条目。
func (n *Navigator) Peers(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	peerIDs := n.registry.GetPeerIDs(id)
	return n.resolveIDs(peerIDs)
}

// CanPrecede 获取指定CWE条目可以前置的条目（链式前驱）。
func (n *Navigator) CanPrecede(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	return n.findByRelationship(id, RelationshipCanPrecede)
}

// CanFollow 获取指定CWE条目可以后随的条目（链式后继）。
func (n *Navigator) CanFollow(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	return n.findByRelationship(id, RelationshipCanFollow)
}

// Requires 获取指定CWE条目所需的条目（复合依赖）。
func (n *Navigator) Requires(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	return n.findByRelationship(id, RelationshipRequires)
}

// RequiredBy 获取需要指定CWE条目的条目。
func (n *Navigator) RequiredBy(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	return n.findByRelationship(id, RelationshipRequiredBy)
}

// CanAlsoBe 获取指定CWE条目也可以被视为的条目。
func (n *Navigator) CanAlsoBe(id int) []*CWE {
	if n.registry == nil {
		return nil
	}
	return n.findByRelationship(id, RelationshipCanAlsoBe)
}

// ChainMembers 获取指定链式复合元素的所有链成员。
func (n *Navigator) ChainMembers(id int) []*CWE {
	if n.registry == nil {
		return nil
	}

	ce, ok := n.registry.GetCompoundElement(id)
	if !ok || ce.Structure != StructureChain {
		return nil
	}

	return n.resolveRelationshipIDs(ce.Relationships)
}

// CompositeMembers 获取指定复合元素的所有组合成员。
func (n *Navigator) CompositeMembers(id int) []*CWE {
	if n.registry == nil {
		return nil
	}

	ce, ok := n.registry.GetCompoundElement(id)
	if !ok || ce.Structure != StructureComposite {
		return nil
	}

	return n.resolveRelationshipIDs(ce.Relationships)
}

// ShortestPath 查找两个CWE条目之间的最短路径。
//
// 使用广度优先搜索在关系图中查找从from到to的最短路径。
// 返回的路径包含起点和终点。
//
// 参数：
//   - from: 起始CWE ID
//   - to: 目标CWE ID
//
// 返回值：
//   - []int: 路径上的CWE ID列表，如果不存在路径返回nil
func (n *Navigator) ShortestPath(from, to int) []int {
	if n.registry == nil {
		return nil
	}

	if from == to {
		return []int{from}
	}

	// BFS搜索
	visited := map[int]bool{from: true}
	parent := map[int]int{from: -1}
	queue := []int{from}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 获取所有相邻节点
		neighbors := n.getNeighbors(current)
		for _, next := range neighbors {
			if visited[next] {
				continue
			}
			visited[next] = true
			parent[next] = current

			if next == to {
				// 构建路径
				var path []int
				for node := to; node != -1; node = parent[node] {
					path = append([]int{node}, path...)
				}
				return path
			}

			queue = append(queue, next)
		}
	}

	return nil // 无路径
}

// IsAncestorOf 检查ancestor是否是descendant的祖先。
func (n *Navigator) IsAncestorOf(ancestor, descendant int) bool {
	if n.registry == nil {
		return false
	}
	ancestors := n.registry.GetAncestorIDs(descendant)
	for _, a := range ancestors {
		if a == ancestor {
			return true
		}
	}
	return false
}

// IsDescendantOf 检查descendant是否是ancestor的后代。
func (n *Navigator) IsDescendantOf(descendant, ancestor int) bool {
	return n.IsAncestorOf(ancestor, descendant)
}

// IsRelated 检查两个CWE条目是否存在任何关系。
func (n *Navigator) IsRelated(a, b int) bool {
	if n.registry == nil {
		return false
	}

	// 检查直接关系
	cwe, ok := n.registry.Get(a)
	if !ok {
		return false
	}
	for _, rel := range cwe.Relationships {
		if rel.CWEID == b {
			return true
		}
	}

	// 检查反向关系
	cwe2, ok := n.registry.Get(b)
	if !ok {
		return false
	}
	for _, rel := range cwe2.Relationships {
		if rel.CWEID == a {
			return true
		}
	}

	return false
}

// RelationshipDepth 计算两个CWE条目之间的层级关系深度。
//
// 如果两者不存在层级关系，返回-1。
// 直接的父子关系深度为1。
func (n *Navigator) RelationshipDepth(ancestor, descendant int) int {
	if n.registry == nil {
		return -1
	}

	if ancestor == descendant {
		return 0
	}

	depth := 0
	queue := []int{ancestor}
	visited := map[int]bool{ancestor: true}

	for len(queue) > 0 {
		depth++
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			current := queue[0]
			queue = queue[1:]

			children := n.registry.GetChildIDs(current)
			for _, childID := range children {
				if childID == descendant {
					return depth
				}
				if !visited[childID] {
					visited[childID] = true
					queue = append(queue, childID)
				}
			}
		}
	}

	return -1
}

// findByRelationship 根据关系类型查找相关条目。
func (n *Navigator) findByRelationship(id int, nature RelationshipNature) []*CWE {
	cwe, ok := n.registry.Get(id)
	if !ok {
		return nil
	}

	var ids []int
	for _, rel := range cwe.Relationships {
		if rel.Nature == nature {
			ids = append(ids, rel.CWEID)
		}
	}

	return n.resolveIDs(ids)
}

// resolveIDs 将ID列表解析为CWE条目列表。
func (n *Navigator) resolveIDs(ids []int) []*CWE {
	var result []*CWE
	for _, id := range ids {
		if cwe, ok := n.registry.Get(id); ok {
			result = append(result, cwe)
		}
	}
	return result
}

// resolveRelationshipIDs 从关系列表中提取CWE ID并解析。
func (n *Navigator) resolveRelationshipIDs(rels []Relationship) []*CWE {
	var ids []int
	for _, rel := range rels {
		ids = append(ids, rel.CWEID)
	}
	return n.resolveIDs(ids)
}

// getNeighbors 获取指定CWE条目的所有相邻节点ID。
func (n *Navigator) getNeighbors(id int) []int {
	var neighbors []int

	// 父级
	neighbors = append(neighbors, n.registry.GetParentIDs(id)...)
	// 子级
	neighbors = append(neighbors, n.registry.GetChildIDs(id)...)
	// 对等
	neighbors = append(neighbors, n.registry.GetPeerIDs(id)...)

	// 从关系列表获取其他关系
	cwe, ok := n.registry.Get(id)
	if ok {
		for _, rel := range cwe.Relationships {
			switch rel.Nature {
			case RelationshipChildOf, RelationshipParentOf,
				RelationshipPeerOf, RelationshipCanAlsoBe:
				// 已通过索引处理
			default:
				neighbors = append(neighbors, rel.CWEID)
			}
		}
	}

	return neighbors
}

// String 返回导航器的字符串表示。
func (n *Navigator) String() string {
	if n.registry == nil {
		return "Navigator<nil>"
	}
	return fmt.Sprintf("Navigator<registry size=%d>", n.registry.Size())
}
