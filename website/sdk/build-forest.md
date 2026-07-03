---
title: 构建森林
outline: [2, 3]
---

# 🌳 构建森林

`BuildForest` 扫描整个 `Registry`，找出所有**无父级**的顶层弱点作为根，对每个根调用 `BuildTree`，返回多棵树组成的森林。

## 📐 函数签名

```go
func BuildForest(r *Registry) []*TreeNode
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 已构建索引的注册表 |
| 返回 | `[]*TreeNode` | 顶层根节点列表（每个根自带完整子树） |

::: tip 顶层根的判定
「顶层」指在 `ChildOf` 关系中没有父级的弱点，等价于 [FindTopLevel](./find-top-level) 的结果。每棵树的根 `Parent` 为 `nil`。
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
	// 两个独立的根
	r.Register(cweskills.NewCWE(703, "Neutralization"))
	r.Register(cweskills.NewCWE(1000, "Another Pillar"))

	c := cweskills.NewCWE(79, "XSS")
	c.Relationships = []cweskills.Relationship{
		{CWEID: 703, Nature: cweskills.RelationshipChildOf},
	}
	r.Register(c)
	r.BuildIndexes()

	forest := cweskills.BuildForest(r)
	fmt.Println("树的数量:", len(forest)) // 2
	for _, root := range forest {
		fmt.Println("根:", root.CWE.CWEID(), "节点数:", root.Count())
	}
}
```

## ⚠️ 注意事项

::: warning 孤立节点也是一棵树
没有父级也没有子级的弱点，自身构成一棵单节点树，仍会出现在森林中。森林大小 = 顶层弱点数量。
:::

::: details 与 BuildTree 的关系
`BuildForest` 本质是「找全部顶层根 → 对每个根 `BuildTree`」。若只关心某个根的子树，直接用 [`BuildTree`](./build-tree) 更高效。
:::

## 🔗 相关链接

- 单棵树：[构建单棵树](./build-tree)
- 顶层弱点查询：[FindTopLevel](./find-top-level)
- 节点方法：[TreeNode 节点](./tree-node)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
