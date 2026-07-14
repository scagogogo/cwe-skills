---
title: 祖先与后代
outline: [2, 3]
---

# 🧭 祖先与后代导航

`Ancestors` 与 `Descendants` 在层级关系图上做**传递闭包**展开——从当前节点出发，沿 `ChildOf`/`ParentOf` 边一路递归，返回全部可达的祖先或后代。

## 📐 方法签名

```go
func (n *Navigator) Ancestors(id int) []*CWE
func (n *Navigator) Descendants(id int) []*CWE
```

| 方法 | 方向 | 边类型 | 返回 |
| --- | --- | --- | --- |
| `Ancestors` | 向上（更抽象） | `ChildOf` | 全部祖先 `[]*CWE` |
| `Descendants` | 向下（更具体） | `ParentOf` | 全部后代 `[]*CWE` |

::: tip 与 Parents/Children 的关系
`Ancestors(79)` = `Parents(79)` ∪ `Parents` 的 `Parents` ∪ …（递归并集）。
`Descendants` 同理对应 `Children` 的递归并集。
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
	// 703 ← 79 ← 791  ←  7910
	addLink := func(child, parent int, name string) {
		c := cweskills.NewCWE(child, name)
		c.Relationships = []cweskills.Relationship{
			{CWEID: parent, Nature: cweskills.RelationshipChildOf},
		}
		r.Register(c)
	}
	r.Register(cweskills.NewCWE(703, "Neutralization"))
	addLink(79, 703, "XSS")
	addLink(791, 79, "XSS Variant")
	addLink(7910, 791, "Reflected XSS")
	r.BuildIndexes()

	nav := cweskills.NewNavigator(r)
	for _, a := range nav.Ancestors(7910) {
		fmt.Println("祖先:", a.CWEID()) // CWE-791 CWE-79 CWE-703
	}
	for _, d := range nav.Descendants(703) {
		fmt.Println("后代:", d.CWEID()) // CWE-79 CWE-791 CWE-7910
	}
}
```

## ⚠️ 注意事项

::: danger 环路会导致无限递归
MITRE 数据理论上不应出现环路，但若数据有误（A ChildOf B 且 B ChildOf A），传递闭包可能无限展开。SDK 内部用访问集合去重，返回的祖先/后代列表中每个节点只出现一次。
:::

::: details 返回顺序
实现通常为 BFS/DFS 遍历顺序，**不保证**按深度排序。如需按深度排序，结合 [RelationshipDepth](./nav-relationship-depth) 自行排序。
:::

## 🔗 相关链接

- 仅一跳：[父子导航](./nav-parents-children)
- 深度计算：[关系深度](./nav-relationship-depth)
- 祖先判定：[祖先与关联判定](./nav-ancestor-related)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
