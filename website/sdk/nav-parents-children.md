---
title: 父子导航
outline: [2, 3]
---

# 🧭 父子导航

`Parents` 与 `Children` 返回与当前弱点**直接相邻**的父级（更抽象）与子级（更具体）弱点。它们是层级树遍历的最基础一跳。

## 📐 方法签名

```go
func (n *Navigator) Parents(id int) []*CWE
func (n *Navigator) Children(id int) []*CWE
```

| 方法 | 参数 | 返回 | 关系 Nature |
| --- | --- | --- | --- |
| `Parents` | `id int` 弱点 ID | 父级 `[]*CWE` | `ChildOf`（我是对方的子项） |
| `Children` | `id int` 弱点 ID | 子级 `[]*CWE` | `ParentOf`（我是对方的父项） |

::: tip 语义记忆
`Parents(79)` = 「79 是谁的子项？」→ 返回更抽象的父级。
`Children(79)` = 「谁是 79 的子项？」→ 返回更具体的子级。
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
	r.Register(cweskills.NewCWE(703, "Neutralization")) // 父

	xss := cweskills.NewCWE(79, "XSS")                  // 子
	xss.Relationships = []cweskills.Relationship{
		{CWEID: 703, Nature: cweskills.RelationshipChildOf},
	}
	r.Register(xss)

	// 79 的子项
	variant := cweskills.NewCWE(791, "XSS Variant")
	variant.Relationships = []cweskills.Relationship{
		{CWEID: 79, Nature: cweskills.RelationshipChildOf},
	}
	r.Register(variant)

	r.BuildIndexes()
	nav := cweskills.NewNavigator(r)

	for _, p := range nav.Parents(79) {
		fmt.Println("父:", p.Name) // 父: Neutralization
	}
	for _, c := range nav.Children(79) {
		fmt.Println("子:", c.Name) // 子: XSS Variant
	}
}
```

## ⚠️ 边界行为

::: warning 一个弱点可有多个父级
MITRE 允许一个弱点 `ChildOf` 多个父级，因此 `Parents` 可能返回多个元素。同理 `Children` 可返回多个子级。
:::

::: details 只走一跳，不递归
`Parents`/`Children` 只返回**直接**邻接节点。如需递归到全部祖先/后代，改用 [Ancestors / Descendants](./nav-ancestors-descendants)。
:::

## 🔗 相关链接

- 递归祖先/后代：[祖先与后代](./nav-ancestors-descendants)
- 兄弟节点：[兄弟与对等](./nav-siblings-peers)
- 返回 ID 的索引版本：[GetParentIDs / GetChildIDs](./relationship-indexes)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
