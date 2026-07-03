---
title: 最短路径
outline: [2, 3]
---

# 🧭 最短路径

`ShortestPath` 用 **BFS（广度优先搜索）** 在 CWE 关系图上寻找从 `from` 到 `to` 的最短跳数路径，返回路径上的节点 ID 序列。

## 📐 方法签名

```go
func (n *Navigator) ShortestPath(from, to int) []int
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `from` | `int` | 起点弱点 ID |
| `to` | `int` | 终点弱点 ID |
| 返回 | `[]int` | 从 `from` 到 `to` 的节点 ID 序列（含两端）；**无路径返回 `nil`** |

::: tip 算法特性
- 采用 BFS，保证跳数最少。
- 边为无向遍历（同时考虑 `ChildOf`/`ParentOf`/`PeerOf` 等各类关系边）。
- 起点等于终点时返回单元素 `[from]`。
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
	path := nav.ShortestPath(7910, 703)
	fmt.Println(path) // [7910 791 79 703]

	fmt.Println(nav.ShortestPath(7910, 9999)) // [] (nil，无路径)
}
```

## ⚠️ 边界行为

::: warning nil 表示无路径
终点不可达时返回 `nil`，不是空切片。判断前用 `if path != nil` 而非 `if len(path) > 0`（后者对 nil 也为 true 的零值场景需注意）。
:::

::: details BFS 复杂度
BFS 时间复杂度 O(V + E)，V 为可达节点数、E 为边数。在大规模注册表上反复调用建议缓存结果。路径不保证唯一——存在多条等长最短路径时返回其中一条。
:::

## 🔗 相关链接

- 关系深度（边数）：[关系深度](./nav-relationship-depth)
- 是否可达判定：[祖先与关联判定](./nav-ancestor-related)
- 祖先/后代闭包：[祖先与后代](./nav-ancestors-descendants)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
