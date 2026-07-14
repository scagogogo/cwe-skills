---
title: 构建视图树
outline: [2, 3]
---

# 🌳 构建视图树

`BuildViewTree` 取一个 [View](./view) 的成员弱点，按其层级关系投影成一棵树。视图是 MITRE 对弱点的「视角编排」，此函数把扁平的视图成员重组成层级结构。

## 📐 函数签名

```go
func BuildViewTree(r *Registry, viewID int) *TreeNode
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 已构建索引的注册表 |
| `viewID` | `int` | 视图 ID |
| 返回 | `*TreeNode` | 视图树的根节点；视图不存在或无成员时返回 `nil` |

::: tip 数据来源
先通过 [GetViewMembers](./relationship-indexes) 取得视图成员 ID 集合，再在成员内部按 `ChildOf` 关系建树。非视图成员的弱点不会出现，即使有父子关系。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
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

	// 假设视图 1000 的成员为 703、79、89（成员关系由数据导入建立）
	root := cweskills.BuildViewTree(r, 1000)
	if root != nil {
		fmt.Println(root.CWE.Name, root.Count())
	}
}
```

## ⚠️ 注意事项

::: warning 视图成员须已建立
`BuildViewTree` 依赖 `GetViewMembers`，而后者依赖视图与弱点的成员关系索引。若数据未通过 `RegisterView` 或导入建立成员关系，返回 `nil`。
:::

::: details 视图树的根
视图成员中无父级（或父级不在视图内）的节点会成为根。若视图有多个这样的节点，`BuildViewTree` 通常以其中一个为根并挂载其余——具体策略以源码实现为准。如需多根，考虑用 [BuildForest](./build-forest) 的思路自行处理。
:::

## 🔗 相关链接

- 视图结构：[View](./view)
- 视图成员查询：[GetViewMembers](./relationship-indexes)
- 全库森林：[构建森林](./build-forest)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
