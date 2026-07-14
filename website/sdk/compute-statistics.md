---
title: 计算统计
outline: [2, 3]
---

# 📊 计算统计

`ComputeStatistics` 是统计模块的入口——一次遍历 `Registry`，聚合出 [Statistics](./stats) 结构体的全部字段。适合生成全库报表。

## 📐 函数签名

```go
func ComputeStatistics(r *Registry) *Statistics
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| 返回 | `*Statistics` | 包含总量与各维度分布的统计结构 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	c1 := cweskills.NewCWE(79, "XSS")
	c1.Abstraction = cweskills.AbstractionBase
	c1.Status = cweskills.StatusStable
	r.Register(c1)
	c2 := cweskills.NewCWE(703, "Neutralization")
	c2.Abstraction = cweskills.AbstractionPillar
	r.Register(c2)
	r.BuildIndexes()

	s := cweskills.ComputeStatistics(r)
	fmt.Println("总数:", s.TotalCWEs)                                   // 2
	fmt.Println("Base:", s.ByAbstraction[cweskills.AbstractionBase])     // 1
	fmt.Println("Pillar:", s.ByAbstraction[cweskills.AbstractionPillar]) // 1
	fmt.Println("Stable:", s.ByStatus[cweskills.StatusStable])           // 1
	fmt.Println("顶层:", s.TopLevelCount)                                // 1
}
```

## ⚠️ 注意事项

::: warning TopLevelCount 依赖索引
`Statistics.TopLevelCount` 通过判断弱点是否有父级得到，需先 `BuildIndexes()`。未构建索引时该字段可能不准确。
:::

::: details 各 map 含零值桶
`ByAbstraction`/`ByStatus` 等 map 会包含零值键（对应未设置该字段的弱点计数）。报表展示时注意是否过滤零值桶。
:::

## 🆚 与 CountBy* 的取舍

| 维度 | `ComputeStatistics` | [`CountBy*`](./count-by-abstraction) |
| --- | --- | --- |
| 范围 | 全维度一次性 | 单维度单次 |
| 适用 | 全库快照 | 只需某一维度的计数 |

需多维度报表用 `ComputeStatistics`；只需一个数用对应 `CountBy*` 更轻量。

## 🔗 相关链接

- 结构体：[统计概览](./stats)
- 单维度计数：[按抽象计数](./count-by-abstraction)
- 源文件：[`stats.go`](https://github.com/scagogogo/cwe-skills/blob/main/stats.go)
