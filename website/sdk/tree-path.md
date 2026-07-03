---
title: 树路径与查找
outline: [2, 3]
---

# 🌳 树路径与查找

`Path` 与 `Find` 是 `TreeNode` 上最常用的定位方法：`Find` 按 ID 在子树中定位节点，`Path` 返回从根到当前节点的完整路径。

## 📐 方法签名

### Find

```go
func (n *TreeNode) Find(id int) *TreeNode
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `int` | 目标弱点 ID |
| 返回 | `*TreeNode` | 子树中匹配的节点；未找到返回 `nil` |

在 `n` 的整棵子树中（含自身）DFS 查找 ID 等于 `id` 的节点。

### Path

```go
func (n *TreeNode) Path() []*TreeNode
```

| 返回 | 说明 |
| --- | --- |
| `[]*TreeNode` | 从根节点到 `n` 的节点序列，根在首位、`n` 在末位 |

通过 `Parent` 指针逐级回溯到根，再反转得到路径。

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
	node := root.Find(791)
	if node != nil {
		fmt.Println("找到:", node.CWE.Name) // Variant
		for _, p := range node.Path() {
			fmt.Println("  路径:", p.CWE.CWEID()) // CWE-703 CWE-79 CWE-791
		}
	}
	fmt.Println(root.Find(9999)) // <nil>
}
```

## ⚠️ 注意事项

::: warning Find 仅在子树内
`Find` 只搜索调用节点**自身及其后代**，不向上查父级链。若需全树查找，在根节点调用，或用 [Registry.Get](./registry-operations) 直接按 ID 取条目。
:::

::: details Path 的根判定
`Path` 通过 `IsRoot()`（`Parent == nil`）判定回溯终止。若节点本身是根，`Path` 返回 `[n]` 单元素切片。
:::

## 🔗 相关链接

- 节点类型：[TreeNode 节点](./tree-node)
- 树遍历：[树遍历](./tree-walk)
- 注册表按 ID 查询：[Get](./registry-operations)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
