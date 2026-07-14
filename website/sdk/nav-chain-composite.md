---
title: 链与复合成员导航
outline: [2, 3]
---

# 🧭 链与复合成员导航

`ChainMembers` 与 `CompositeMembers` 针对 MITRE 的**聚合型弱点**：`Structure=Chain`（链式）和 `Structure=Composite`（复合）的弱点本身由多个成员弱点组合而成。这两个方法返回这些成员。

## 📐 方法签名

```go
func (n *Navigator) ChainMembers(id int) []*CWE
func (n *Navigator) CompositeMembers(id int) []*CWE
```

| 方法 | 适用对象 | 返回 |
| --- | --- | --- |
| `ChainMembers` | `Structure=Chain` 的弱点 | 链中成员弱点 `[]*CWE` |
| `CompositeMembers` | `Structure=Composite` 的弱点 | 复合体成员弱点 `[]*CWE` |

::: tip 与 CanPrecede 的区别
- `CanPrecede`：任意两个弱点间的时序边。
- `ChainMembers`：某个**链类型弱点**（如 CWE-680，Structure=Chain）所**封装**的成员列表。

后者是「聚合关系」，前者是「直接关系边」。
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
	r.Register(cweskills.NewCWE(79, "XSS"))
	r.Register(cweskills.NewCWE(89, "SQLi"))

	// 链式弱点 680 由 89、79 组成
	chain := cweskills.NewCWE(680, "Chain Example")
	chain.Structure = cweskills.StructureChain
	chain.Relationships = []cweskills.Relationship{
		{CWEID: 89, Nature: cweskills.RelationshipCanPrecede},
		{CWEID: 79, Nature: cweskills.RelationshipCanPrecede},
	}
	r.Register(chain)

	r.BuildIndexes()
	nav := cweskills.NewNavigator(r)

	for _, m := range nav.ChainMembers(680) {
		fmt.Println("链成员:", m.CWEID()) // CWE-89 CWE-79
	}
}
```

## ⚠️ 注意事项

::: warning 对象须为聚合型
对非 Chain/Composite 类型的弱点调用这两个方法，返回空列表。建议先用 `cwe.IsChain()` / `cwe.IsComposite()` 判断，参见 [CWE 类型判断方法](./cwe-type-methods)。
:::

::: details 成员来源
成员通常通过 `CanPrecede`（链）或组合关系（复合）边聚合。具体边类型以源码实现为准，调用方应把结果视为「该聚合弱点的组成成员集合」。
:::

## 🔗 相关链接

- 时序边查询：[前置与后继](./nav-precede-follow)
- 类型判断：[CWE 类型判断方法](./cwe-type-methods)
- Structure 枚举：[Structure](./enum-structure)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
