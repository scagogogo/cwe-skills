---
title: 按抽象计数
outline: [2, 3]
---

# 📊 按抽象计数

`CountByAbstraction` 返回指定 [Abstraction](./enum-abstraction) 层级的弱点数量。当只需单个数字时，比 `ComputeStatistics` 更轻量。

## 📐 函数签名

```go
func CountByAbstraction(r *Registry, abstraction Abstraction) int
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `abstraction` | [`Abstraction`](./enum-abstraction) | 抽象层级枚举 |
| 返回 | `int` | 该抽象层级的弱点数量 |

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
	p := cweskills.NewCWE(703, "Neutralization")
	p.Abstraction = cweskills.AbstractionPillar
	r.Register(p)

	fmt.Println("Base:", cweskills.CountByAbstraction(r, cweskills.AbstractionBase))     // 1
	fmt.Println("Pillar:", cweskills.CountByAbstraction(r, cweskills.AbstractionPillar)) // 1
	fmt.Println("Variant:", cweskills.CountByAbstraction(r, cweskills.AbstractionVariant)) // 0
}
```

## 🆚 与 ComputeStatistics 的取舍

| 维度 | `CountByAbstraction` | [`ComputeStatistics`](./compute-statistics) |
| --- | --- | --- |
| 返回 | 单个 `int` | 全维度 `*Statistics` |
| 遍历 | 每次调用一次 O(N) | 一次 O(N) 得全部 |
| 适用 | 只需一个数 | 需多维度报表 |

只需 Base 数量用本函数；需同时看 Pillar/Class/Variant 分布用 `ComputeStatistics` 一次算全。

## ⚠️ 注意事项

::: tip 不依赖索引
本函数遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[Abstraction](./enum-abstraction)
- 全量统计：[计算统计](./compute-statistics)
- 按抽象查找：[按抽象查找](./find-by-abstraction)
- 源文件：[`stats.go`](https://github.com/scagogogo/cwe-skills/blob/main/stats.go)
