---
title: 兄弟与对等导航
outline: [2, 3]
---

# 🧭 兄弟与对等导航

`Siblings` 与 `Peers` 返回与当前弱点「平级」相关的节点。两者语义相近但取数路径不同：`Siblings` 通过共同父级推导，`Peers` 直接读 `PeerOf`/`CanAlsoBe` 关系边。

## 📐 方法签名

```go
func (n *Navigator) Siblings(id int) []*CWE
func (n *Navigator) Peers(id int) []*CWE
```

| 方法 | 推导方式 | 返回 |
| --- | --- | --- |
| `Siblings` | 取所有父级 → 取父级的所有子级 → 排除自身 | 共享父级的兄弟 `[]*CWE` |
| `Peers` | 直接筛选 `PeerOf` 与 `CanAlsoBe` 关系 | 显式声明的对等 `[]*CWE` |

::: tip 二者区别
- `Siblings`：A 和 B 有同一个父级，**未必**显式声明彼此关系。
- `Peers`：A 的 `Relationships` 中显式写了 `PeerOf B` 或 `CanAlsoBe B`。

一个节点可以是 `Siblings` 但不是 `Peers`，反之亦然。
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

	// 79、89 都是 703 的子项 → 互为兄弟
	for _, id := range []int{79, 89} {
		c := cweskills.NewCWE(id, fmt.Sprintf("CWE-%d", id))
		c.Relationships = []cweskills.Relationship{
			{CWEID: 703, Nature: cweskills.RelationshipChildOf},
		}
		r.Register(c)
	}

	// 80 显式声明与 79 对等
	c80 := cweskills.NewCWE(80, "Buffer Overflow")
	c80.Relationships = []cweskills.Relationship{
		{CWEID: 79, Nature: cweskills.RelationshipPeerOf},
	}
	r.Register(c80)

	r.BuildIndexes()
	nav := cweskills.NewNavigator(r)

	for _, s := range nav.Siblings(79) {
		fmt.Println("兄弟:", s.CWEID()) // CWE-89
	}
	for _, p := range nav.Peers(79) {
		fmt.Println("对等:", p.CWEID()) // CWE-80
	}
}
```

## ⚠️ 边界行为

::: warning Siblings 排除自身
`Siblings(79)` 的结果中不会包含 79 自身，即使数据有自环。
:::

::: details Peers 合并两种 Nature
`Peers` 同时纳入 `PeerOf` 与 `CanAlsoBe`。若只想取其中一种，需自行对 `CWE.Relationships` 过滤，参见 [CWE 关系获取方法](./cwe-relationship-methods)。
:::

## 🔗 相关链接

- 父级查询：[父子导航](./nav-parents-children)
- `CanAlsoBe` 单独导航：[可作为](./nav-can-also-be)
- 本地关系读取：[CWE 关系获取方法](./cwe-relationship-methods)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
