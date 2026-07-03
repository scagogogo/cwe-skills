---
title: 统计概览
outline: [2, 3]
---

# 📊 统计概览

`stats.go` 提供 `Statistics` 聚合结构与一组计数函数，对 `Registry` 做多维度汇总。它是「全库快照」——一次计算得到总量、按抽象/状态/可能性/后果范围的分布与顶层计数。

## 📋 核心结构

```go
type Statistics struct {
    TotalCWEs         int
    TotalCategories   int
    TotalViews        int
    ByAbstraction     map[Abstraction]int
    ByStatus          map[Status]int
    ByLikelihood      map[LikelihoodOfExploit]int
    ByConsequenceScope map[ConsequenceScope]int
    TopLevelCount     int
}

type ConsequenceScopeCount struct {
    Scope ConsequenceScope
    Count int
}
```

| 字段 | 说明 |
| --- | --- |
| `TotalCWEs` | 弱点总数 |
| `TotalCategories` | 分类总数 |
| `TotalViews` | 视图总数 |
| `ByAbstraction` | 按抽象层级的计数 |
| `ByStatus` | 按状态的计数 |
| `ByLikelihood` | 按被利用可能性的计数 |
| `ByConsequenceScope` | 按后果范围的计数 |
| `TopLevelCount` | 顶层（无父级）弱点数 |

## 📚 本组文档导航

| 文档 | 主题 | 函数 |
| --- | --- | --- |
| [计算统计](./compute-statistics) | 一次性聚合 | `ComputeStatistics` |
| [按抽象计数](./count-by-abstraction) | 单维度计数 | `CountByAbstraction` |
| [按状态计数](./count-by-status) | 单维度计数 | `CountByStatus` |
| [按可能性计数](./count-by-likelihood) | 单维度计数 | `CountByLikelihood` |
| [按范围计数](./count-by-scope) | 单维度计数 | `CountByScope` |

## ✅ 快速上手

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
	c.Status = cweskills.StatusStable
	r.Register(c)
	r.BuildIndexes()

	s := cweskills.ComputeStatistics(r)
	fmt.Println(s.TotalCWEs)                       // 1
	fmt.Println(s.ByAbstraction[cweskills.AbstractionBase]) // 1
	fmt.Println(cweskills.CountByStatus(r, cweskills.StatusStable)) // 1
}
```

## ⚠️ 注意事项

::: tip 计数函数不依赖索引
`CountBy*` 系列遍历 `GetAll()`，无需 `BuildIndexes()`。但 `ComputeStatistics` 中的 `TopLevelCount` 需要索引判断父级，调用前应 `BuildIndexes()`。
:::

## 🔗 相关链接

- 数据来源：[Registry 基础操作](./registry-operations)
- 查找过滤：[搜索概览](./search)
- 源文件：[`stats.go`](https://github.com/scagogogo/cwe-skills/blob/main/stats.go)
