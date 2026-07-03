---
title: 树深度与计数
outline: [2, 3]
---

# 🌳 树深度与计数

`MaxDepth` 与 `Count` 给出树的两个全局规模指标：最大层级深度与节点总数。常用于评估树规模、控制递归深度、生成报表摘要。

## 📐 方法签名

### MaxDepth

```go
func (n *TreeNode) MaxDepth() int
```

| 返回 | 说明 |
| --- | --- |
| `int` | `n` 子树的最大深度。根节点自身深度为 `0`，每下传一层 `+1` |

### Count

```go
func (n *TreeNode) Count() int
```

| 返回 | 说明 |
| --- | --- |
| `int` | `n` 子树的节点总数（含 `n` 自身） |

::: tip 深度定义
根节点深度 `0`；根的直接子节点深度 `1`。因此单节点树 `MaxDepth() == 0`、`Count() == 1`。
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
	link(791, 79, "Variant")
	r.BuildIndexes()

	root := cweskills.BuildTree(r, 703)
	fmt.Println(root.Count())    // 3
	fmt.Println(root.MaxDepth()) // 2
}
```

## ⚠️ 注意事项

::: warning 深度取最长路径
`MaxDepth` 取子树中**最长**根到叶路径的边数。若树不平衡，深度由最深分支决定，浅分支不影响结果。
:::

::: details 性能
`Count` 与 `MaxDepth` 都是 O(N) 递归遍历，N 为子树节点数。频繁调用建议缓存到结构体字段或外部变量。
:::

## 🔗 相关链接

- 节点类型：[TreeNode 节点](./tree-node)
- 叶子收集：[叶子节点](./tree-leaf-nodes)
- 树遍历：[树遍历](./tree-walk)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
