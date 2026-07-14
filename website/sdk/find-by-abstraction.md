---
title: 按抽象查找
outline: [2, 3]
---

# 🔍 按抽象查找

`FindByAbstraction` 按 [Abstraction](./enum-abstraction) 枚举筛选弱点，返回指定抽象层级的全部弱点。MITRE 的四级抽象：Pillar > Class > Base > Variant。

## 📐 函数签名

```go
func FindByAbstraction(r *Registry, abstraction Abstraction) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `abstraction` | [`Abstraction`](./enum-abstraction) | 抽象层级枚举 |
| 返回 | `[]*CWE` | 该抽象层级的全部弱点 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	pillar := cweskills.NewCWE(703, "Neutralization")
	pillar.Abstraction = cweskills.AbstractionPillar
	r.Register(pillar)

	base := cweskills.NewCWE(79, "XSS")
	base.Abstraction = cweskills.AbstractionBase
	r.Register(base)

	variant := cweskills.NewCWE(791, "Variant")
	variant.Abstraction = cweskills.AbstractionVariant
	r.Register(variant)

	fmt.Println("Pillar:", len(cweskills.FindByAbstraction(r, cweskills.AbstractionPillar))) // 1
	fmt.Println("Base:", len(cweskills.FindByAbstraction(r, cweskills.AbstractionBase)))     // 1
	fmt.Println("Variant:", len(cweskills.FindByAbstraction(r, cweskills.AbstractionVariant))) // 1
}
```

## 🆚 与 Filter 的取舍

| 维度 | `FindByAbstraction` | [`Filter`](./filter) + `FilterOption.Abstraction` |
| --- | --- | --- |
| 单条件 | 简洁 | 稍冗长 |
| 多条件 | 不支持 | 支持 AND 组合 |
| 输入 | `*Registry` | `[]*CWE` |

只按抽象单维度查时用 `FindByAbstraction`；需同时按状态、关键词等多维度筛选用 `Filter`。

## ⚠️ 注意事项

::: tip 不依赖索引
`FindByAbstraction` 遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[Abstraction](./enum-abstraction)
- 多条件过滤：[Filter](./filter)
- 抽象统计：[CountByAbstraction](./count-by-abstraction)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
