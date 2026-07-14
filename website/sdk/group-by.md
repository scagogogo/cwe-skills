---
title: 分组
outline: [2, 3]
---

# 🧹 分组

`GroupByAbstraction`、`GroupByStatus`、`GroupByLikelihood` 把 `[]*CWE` 按枚举维度分桶，返回 `map[枚举][]*CWE`。常用于报表分类展示。

## 📐 函数签名

```go
func GroupByAbstraction(cwes []*CWE) map[Abstraction][]*CWE
func GroupByStatus(cwes []*CWE) map[Status][]*CWE
func GroupByLikelihood(cwes []*CWE) map[LikelihoodOfExploit][]*CWE
```

| 函数 | 分组键 | 返回 |
| --- | --- | --- |
| `GroupByAbstraction` | [`Abstraction`](./enum-abstraction) | `map[Abstraction][]*CWE` |
| `GroupByStatus` | [`Status`](./enum-status) | `map[Status][]*CWE` |
| `GroupByLikelihood` | [`LikelihoodOfExploit`](./enum-likelihood) | `map[LikelihoodOfExploit][]*CWE` |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	b := cweskills.NewCWE(79, "XSS")
	b.Abstraction = cweskills.AbstractionBase
	r.Register(b)
	v := cweskills.NewCWE(791, "Variant")
	v.Abstraction = cweskills.AbstractionVariant
	r.Register(v)

	groups := cweskills.GroupByAbstraction(r.GetAll())
	for abs, list := range groups {
		fmt.Printf("%s: %d 个\n", abs, len(list))
	}
}
```

## ⚠️ 注意事项

::: warning 零值分组
未设置该字段的弱点会被归入枚举零值那一桶。若零值不是合法取值，报表里会出现「未知」桶；调用方应决定是过滤还是展示。
:::

::: details map 顺序
返回的 `map` 遍历顺序不确定。报表展示建议按枚举常量定义顺序输出，而非依赖 map 迭代。
:::

## 🔗 相关链接

- 枚举：[Abstraction](./enum-abstraction)、[Status](./enum-status)、[LikelihoodOfExploit](./enum-likelihood)
- 计数版本：[ComputeStatistics](./compute-statistics)
- 去重：[去重](./deduplicate)
- 源文件：[`filter.go`](https://github.com/scagogogo/cwe-skills/blob/main/filter.go)
