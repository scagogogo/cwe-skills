---
title: 祖先与关联判定
outline: [2, 3]
---

# 🧭 祖先与关联判定

`IsAncestorOf`、`IsDescendantOf`、`IsRelated` 三个布尔方法回答「两个弱点在关系图上是什么关系」的问题，常用于权限校验、分类判定、可视化高亮等场景。

## 📐 方法签名

```go
func (n *Navigator) IsAncestorOf(ancestor, descendant int) bool
func (n *Navigator) IsDescendantOf(descendant, ancestor int) bool
func (n *Navigator) IsRelated(a, b int) bool
```

| 方法 | 参数顺序 | 返回 true 当且仅当 |
| --- | --- | --- |
| `IsAncestorOf` | `(ancestor, descendant)` | `ancestor` 是 `descendant` 的祖先（含传递） |
| `IsDescendantOf` | `(descendant, ancestor)` | `descendant` 是 `ancestor` 的后代（含传递） |
| `IsRelated` | `(a, b)` | `a` 与 `b` 之间存在任意关系边 |

::: tip 互为反向
`IsAncestorOf(x, y)` 与 `IsDescendantOf(y, x)` 等价。参数顺序是唯一区别，按调用方语义自然选择即可。
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
	c := cweskills.NewCWE(79, "XSS")
	c.Relationships = []cweskills.Relationship{
		{CWEID: 703, Nature: cweskills.RelationshipChildOf},
		{CWEID: 80, Nature: cweskills.RelationshipPeerOf},
	}
	r.Register(c)
	r.Register(cweskills.NewCWE(80, "Buffer Overflow"))
	r.BuildIndexes()

	nav := cweskills.NewNavigator(r)
	fmt.Println(nav.IsAncestorOf(703, 79))  // true
	fmt.Println(nav.IsDescendantOf(79, 703)) // true
	fmt.Println(nav.IsRelated(79, 80))       // true（PeerOf）
	fmt.Println(nav.IsRelated(79, 9999))     // false
}
```

## ⚠️ 注意事项

::: warning IsRelated 的范围
`IsRelated` 判断的是**直接邻接**（任意关系边），不等同于「同一连通分量」。若需判断间接可达，结合 [ShortestPath](./nav-shortest-path)：路径非 nil 即可达。
:::

::: details 含传递闭包
`IsAncestorOf`/`IsDescendantOf` 基于 [Ancestors/Descendants](./nav-ancestors-descendants) 的传递闭包，不止一跳。直接父子也算祖先关系。
:::

## 🔗 相关链接

- 闭包来源：[祖先与后代](./nav-ancestors-descendants)
- 间接可达：[最短路径](./nav-shortest-path)
- 边数量化：[关系深度](./nav-relationship-depth)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
