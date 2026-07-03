---
title: 树叶子节点
outline: [2, 3]
---

# 🌳 树叶子节点

`LeafNodes` 收集子树中所有叶子节点（无子节点的节点）；`IsLeaf` 判断单个节点是否为叶子。叶子通常对应最具体的 Variant/弱点，是风险清单的重点。

## 📐 方法签名

### LeafNodes

```go
func (n *TreeNode) []*TreeNode // 方法名为 LeafNodes
```

准确签名：`func (n *TreeNode) LeafNodes() []*TreeNode`

| 返回 | 说明 |
| --- | --- |
| `[]*TreeNode` | `n` 子树中所有 `Children` 为空的节点 |

### IsLeaf

```go
func (n *TreeNode) IsLeaf() bool
```

| 返回 | 说明 |
| --- | --- |
| `bool` | `Children` 为空时 `true` |

::: tip 叶子 = 最具体弱点
在 CWE 层级中，叶子节点通常是没有进一步子级的 Variant 或 Base 弱点，代表可直接出现在漏洞报告里的具体缺陷。
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
	link(79, 703, "XSS")     // 叶子
	link(89, 703, "SQLi")    // 叶子
	r.BuildIndexes()

	root := cweskills.BuildTree(r, 703)
	fmt.Println(root.IsLeaf()) // false
	for _, leaf := range root.LeafNodes() {
		fmt.Println("叶子:", leaf.CWE.CWEID()) // CWE-79 CWE-89
	}
}
```

## ⚠️ 注意事项

::: warning 单节点树也是叶子
若 `n` 自身无子节点，`LeafNodes` 返回 `[n]`（它自己既是根也是叶），`IsLeaf` 返回 `true`。
:::

::: details 顺序
`LeafNodes` 通常按 DFS 遍历顺序收集叶子。如需按 ID 排序，对结果用 [SortByID](./sort) 处理。
:::

## 🔗 相关链接

- 节点类型：[TreeNode 节点](./tree-node)
- 深度与计数：[深度与计数](./tree-depth-count)
- 树遍历：[树遍历](./tree-walk)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
