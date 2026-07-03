---
title: 按后果范围计数
outline: [2, 3]
---

# 📊 按后果范围计数

`CountByScope` 返回具有指定 [ConsequenceScope](./enum-consequence-scope) 后果范围的弱点数量。用于按安全三要素（CIA：Confidentiality/Integrity/Availability）统计风险分布。

## 📐 函数签名

```go
func CountByScope(r *Registry, scope ConsequenceScope) int
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `scope` | [`ConsequenceScope`](./enum-consequence-scope) | 后果范围枚举 |
| 返回 | `int` | `CommonConsequences` 中含该范围的弱点数量 |

::: tip 数据来源
遍历每个弱点的 `CommonConsequences`，只要其中任一 `Consequence.Scope` 匹配即计入。一个弱点可能跨多个范围，因此各范围计数之和可能大于弱点总数。
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
	c.CommonConsequences = []cweskills.Consequence{
		{Scope: cweskills.ConsequenceScopeConfidentiality},
		{Scope: cweskills.ConsequenceScopeIntegrity},
	}
	r.Register(c)

	fmt.Println("Confidentiality:", cweskills.CountByScope(r, cweskills.ConsequenceScopeConfidentiality)) // 1
	fmt.Println("Integrity:", cweskills.CountByScope(r, cweskills.ConsequenceScopeIntegrity))           // 1
	fmt.Println("Availability:", cweskills.CountByScope(r, cweskills.ConsequenceScopeAvailability))     // 0
}
```

## 🆚 与 ComputeStatistics 的取舍

| 维度 | `CountByScope` | [`ComputeStatistics`](./compute-statistics) |
| --- | --- | --- |
| 返回 | 单个 `int` | 全维度 `*Statistics` |
| 适用 | 单范围指标 | 全范围分布 |

只需「影响机密性的弱点数」用本函数；需 CIA 全分布用 `ComputeStatistics().ByConsequenceScope`。

## ⚠️ 注意事项

::: warning 依赖后果字段
仅当弱点显式声明 `CommonConsequences` 且含匹配 `Scope` 时才计数。无后果字段的弱点不计入任何范围。
:::

::: tip 不依赖关系索引
本函数遍历 `GetAll()` 与各弱点的 `CommonConsequences`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[ConsequenceScope](./enum-consequence-scope)
- 后果结构：[Consequence](./consequence)
- 全量统计：[计算统计](./compute-statistics)
- 按范围查找：[按后果范围查找](./find-by-consequence-scope)
- 源文件：[`stats.go`](https://github.com/scagogogo/cwe-skills/blob/main/stats.go)
