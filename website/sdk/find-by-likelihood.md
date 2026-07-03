---
title: 按可能性查找
outline: [2, 3]
---

# 🔍 按被利用可能性查找

`FindByLikelihood` 按 [LikelihoodOfExploit](./enum-likelihood) 枚举筛选弱点，返回被利用可能性等级（如 High、Medium、Low）匹配的弱点。常用于风险评估优先级排序。

## 📐 函数签名

```go
func FindByLikelihood(r *Registry, likelihood LikelihoodOfExploit) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `likelihood` | [`LikelihoodOfExploit`](./enum-likelihood) | 可能性枚举 |
| 返回 | `[]*CWE` | 该可能性等级的全部弱点 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	high := cweskills.NewCWE(79, "XSS")
	high.LikelihoodOfExploit = cweskills.LikelihoodHigh
	r.Register(high)

	low := cweskills.NewCWE(999, "Rare")
	low.LikelihoodOfExploit = cweskills.LikelihoodLow
	r.Register(low)

	highs := cweskills.FindByLikelihood(r, cweskills.LikelihoodHigh)
	fmt.Println("High:", len(highs)) // 1
	fmt.Println("Low:", len(cweskills.FindByLikelihood(r, cweskills.LikelihoodLow))) // 1
}
```

## ⚠️ 注意事项

::: warning 未设置的可能性
未显式赋值的弱点 `LikelihoodOfExploit` 为零值。`FindByLikelihood(r, 零值)` 可能匹配到大量未设置的弱点，调用时注意区分「明确为 Low」与「未设置」。
:::

::: tip 不依赖索引
`FindByLikelihood` 遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[LikelihoodOfExploit](./enum-likelihood)
- 多条件过滤：[Filter](./filter)
- 可能性统计：[CountByLikelihood](./count-by-likelihood)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
