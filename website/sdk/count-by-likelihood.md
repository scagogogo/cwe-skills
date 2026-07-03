---
title: 按可能性计数
outline: [2, 3]
---

# 📊 按可能性计数

`CountByLikelihood` 返回指定 [LikelihoodOfExploit](./enum-likelihood) 等级的弱点数量。用于风险评估中「高利用可能性弱点数」等关键指标。

## 📐 函数签名

```go
func CountByLikelihood(r *Registry, likelihood LikelihoodOfExploit) int
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `likelihood` | [`LikelihoodOfExploit`](./enum-likelihood) | 可能性枚举 |
| 返回 | `int` | 该可能性等级的弱点数量 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	h := cweskills.NewCWE(79, "XSS")
	h.LikelihoodOfExploit = cweskills.LikelihoodHigh
	r.Register(h)
	l := cweskills.NewCWE(999, "Rare")
	l.LikelihoodOfExploit = cweskills.LikelihoodLow
	r.Register(l)

	fmt.Println("High:", cweskills.CountByLikelihood(r, cweskills.LikelihoodHigh)) // 1
	fmt.Println("Low:", cweskills.CountByLikelihood(r, cweskills.LikelihoodLow))   // 1
}
```

## 🆚 与 ComputeStatistics 的取舍

| 维度 | `CountByLikelihood` | [`ComputeStatistics`](./compute-statistics) |
| --- | --- | --- |
| 返回 | 单个 `int` | 全维度 `*Statistics` |
| 适用 | 风险指标单值 | 全可能性分布 |

只需「高风险弱点数」用本函数；需可能性分布报表用 `ComputeStatistics().ByLikelihood`。

## ⚠️ 注意事项

::: warning 零值可能性
未显式设置 `LikelihoodOfExploit` 的弱点计入零值桶。`CountByLikelihood(r, 零值)` 可能偏大，调用时区分「明确为某级」与「未设置」。
:::

::: tip 不依赖索引
本函数遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[LikelihoodOfExploit](./enum-likelihood)
- 全量统计：[计算统计](./compute-statistics)
- 按可能性查找：[按可能性查找](./find-by-likelihood)
- 源文件：[`stats.go`](https://github.com/scagogogo/cwe-skills/blob/main/stats.go)
