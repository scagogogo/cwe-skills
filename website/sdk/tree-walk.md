---
title: 树遍历
outline: [2, 3]
---

# 🌳 树遍历

`Walk` 与 `WalkBFS` 提供两种遍历 `TreeNode` 的方式：深度优先（DFS）与广度优先（BFS）。两者都通过回调函数访问节点，回调返回 `false` 时**提前终止**遍历。

## 📐 方法签名

```go
func (n *TreeNode) Walk(fn func(*TreeNode) bool)
func (n *TreeNode) WalkBFS(fn func(*TreeNode) bool)
```

| 方法 | 参数 | 返回 | 遍历顺序 |
| --- | --- | --- | --- |
| `Walk` | `fn func(*TreeNode) bool` | 无 | 深度优先（先访问当前节点，再递归子节点） |
| `WalkBFS` | `fn func(*TreeNode) bool` | 无 | 广度优先（按层访问） |

::: tip 回调返回值
回调返回 `true` 继续遍历下一个节点；返回 `false` 立即停止整个遍历。这是「找到即停」的标准用法。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	r.Register(cweskills.NewCWE(703, "Neutralization"))
	link := func(child, parent int, name string) {
		c := cweskills.NewCWE(child, name)
		c.Relationships = []cweskills.Relationship{
			{CWEID: parent, Nature: cweskills.RelationshipChildOf},
		}
		r.Register(c)
	}
	link(79, 703, "XSS")
	link(89, 703, "SQLi")
	r.BuildIndexes()

	root := cweskills.BuildTree(r, 703)

	fmt.Println("DFS:")
	root.Walk(func(n *cweskills.TreeNode) bool {
		fmt.Println("  ", n.CWE.CWEID())
		return true // 继续遍历
	})

	fmt.Println("BFS 找到 XSS 即停:")
	root.WalkBFS(func(n *cweskills.TreeNode) bool {
		fmt.Println("  访问:", n.CWE.CWEID())
		if n.CWE.ID == 79 {
			fmt.Println("  找到，停止")
			return false
		}
		return true
	})
}
```

## ⚠️ 注意事项

::: warning Walk 不含根节点之后访问
DFS 实现通常先访问当前节点再访问子节点（前序）。若需后序（先子后根），需自行递归实现。
:::

::: details 遍历起点
`Walk`/`WalkBFS` 从调用它们的节点开始遍历其**整个子树**，包含该节点本身。要从根遍历，在根节点上调用。
:::

## 🔗 相关链接

- 节点类型：[TreeNode 节点](./tree-node)
- 路径与查找：[路径查询](./tree-path)
- 构建树：[构建单棵树](./build-tree)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
