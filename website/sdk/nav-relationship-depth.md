---
title: 关系深度
outline: [2, 3]
---

# 🧭 关系深度

`RelationshipDepth` 返回祖先到后代的**层级距离**（边数）。它基于祖先闭包展开，给出量化的「远近」指标，常用于排序、分组与可视化层级控制。

## 📐 方法签名

```go
func (n *Navigator) RelationshipDepth(ancestor, descendant int) int
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ancestor` | `int` | 祖先弱点 ID |
| `descendant` | `int` | 后代弱点 ID |
| 返回 | `int` | 从 `ancestor` 到 `descendant` 的边数；**无关系返回 `-1`** |

::: warning 返回 -1 而非 0
当 `ancestor` 与 `descendant` 不存在祖先-后代关系时返回 `-1`，**不是** `0`。`0` 表示同一个节点（ancestor == descendant）。判断时用 `depth >= 0` 而非 `depth != 0`。
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
	// 703 ← 79 ← 791 ← 7910
	link := func(child, parent int, name string) {
		c := cweskills.NewCWE(child, name)
		c.Relationships = []cweskills.Relationship{
			{CWEID: parent, Nature: cweskills.RelationshipChildOf},
		}
		r.Register(c)
	}
	r.Register(cweskills.NewCWE(703, "Neutralization"))
	link(79, 703, "XSS")
	link(791, 79, "Variant")
	link(7910, 791, "Reflected")
	r.BuildIndexes()

	nav := cweskills.NewNavigator(r)
	fmt.Println(nav.RelationshipDepth(703, 7910)) // 3
	fmt.Println(nav.RelationshipDepth(703, 703))  // 0
	fmt.Println(nav.RelationshipDepth(7910, 703)) // -1（方向反了）
}
```

## 🆚 与 ShortestPath 的区别

| 维度 | `RelationshipDepth` | [`ShortestPath`](./nav-shortest-path) |
| --- | --- | --- |
| 方向 | 有向（仅祖先→后代） | 无向（任意关系边） |
| 返回 | 边数 `int` | 节点序列 `[]int` |
| 无路径 | `-1` | `nil` |
| 适用 | 层级深度量化 | 任意可达路径 |

只需「隔几代」用 `RelationshipDepth` 更轻量；需要具体路径用 `ShortestPath`。

## ⚠️ 注意事项

::: danger 方向敏感
`RelationshipDepth(A, B)` 只在 A 是 B 的祖先时有意义。反过来调用 `RelationshipDepth(B, A)` 返回 `-1`（除非 B 也是 A 的祖先，即成环）。如需无向距离，用 [ShortestPath](./nav-shortest-path) 的长度减一。
:::

## 🔗 相关链接

- 路径序列：[最短路径](./nav-shortest-path)
- 祖先判定：[祖先与关联判定](./nav-ancestor-related)
- 闭包来源：[祖先与后代](./nav-ancestors-descendants)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
