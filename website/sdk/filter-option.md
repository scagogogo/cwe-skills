---
title: FilterOption 过滤选项
outline: [2, 3]
---

# 🧹 FilterOption 过滤选项

`FilterOption` 是 [Filter](./filter) 的配置结构，把多个过滤维度打包成一个值。零值字段表示「不限制该维度」。

## 📋 结构体定义

```go
type FilterOption struct {
    Keyword              string
    Abstraction          Abstraction
    Structure            Structure
    Status               Status
    LikelihoodOfExploit  LikelihoodOfExploit
    ConsequenceScope     ConsequenceScope
}
```

| 字段 | 类型 | 作用 |
| --- | --- | --- |
| `Keyword` | `string` | 关键词子串匹配（Name/Description/ExtendedDescription） |
| `Abstraction` | [`Abstraction`](./enum-abstraction) | 抽象层级 |
| `Structure` | [`Structure`](./enum-structure) | 结构 |
| `Status` | [`Status`](./enum-status) | 状态 |
| `LikelihoodOfExploit` | [`LikelihoodOfExploit`](./enum-likelihood) | 被利用可能性 |
| `ConsequenceScope` | [`ConsequenceScope`](./enum-consequence-scope) | 后果范围 |

::: tip 与 FindBy* 的对应
每个字段都对应一个 `FindBy*` 函数。`Filter` 相当于把多个 `FindBy*` 的条件合并成一次 AND 查询，避免多次遍历。
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
	c := cweskills.NewCWE(79, "XSS")
	c.Abstraction = cweskills.AbstractionBase
	c.Structure = cweskills.StructureSimple
	r.Register(c)
	r.BuildIndexes()

	// 只按抽象过滤，其余字段零值不限制
	opt := cweskills.FilterOption{Abstraction: cweskills.AbstractionBase}
	got := cweskills.Filter(r.GetAll(), opt)
	fmt.Println(len(got)) // 1

	// 链式追加多个选项（AND）
	got2 := cweskills.Filter(got,
		cweskills.FilterOption{Structure: cweskills.StructureSimple})
	fmt.Println(len(got2)) // 1
}
```

## ⚠️ 注意事项

::: warning 枚举零值歧义
若某枚举的零值恰是合法取值（如 `Status` 的零值可能是某个有效状态），用它作过滤条件会误把「未设置」与「明确为零值」混在一起。遇到此情况改用对应的 `FindBy*` 单条件函数。
:::

## 🔗 相关链接

- 使用入口：[Filter](./filter)
- 各枚举：[Abstraction](./enum-abstraction)、[Structure](./enum-structure)、[Status](./enum-status)
- 单条件查找：[搜索概览](./search)
- 源文件：[`filter.go`](https://github.com/scagogogo/cwe-skills/blob/main/filter.go)
