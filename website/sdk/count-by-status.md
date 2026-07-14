---
title: 按状态计数
outline: [2, 3]
---

# 📊 按状态计数

`CountByStatus` 返回指定 [Status](./enum-status) 状态的弱点数量。常用于「有效弱点数」「已弃用数」等指标。

## 📐 函数签名

```go
func CountByStatus(r *Registry, status Status) int
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `status` | [`Status`](./enum-status) | 状态枚举 |
| 返回 | `int` | 该状态的弱点数量 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	s := cweskills.NewCWE(79, "XSS")
	s.Status = cweskills.StatusStable
	r.Register(s)
	d := cweskills.NewCWE(100, "Old")
	d.Status = cweskills.StatusDeprecated
	r.Register(d)

	fmt.Println("Stable:", cweskills.CountByStatus(r, cweskills.StatusStable))         // 1
	fmt.Println("Deprecated:", cweskills.CountByStatus(r, cweskills.StatusDeprecated)) // 1
}
```

## 🆚 与 ComputeStatistics 的取舍

| 维度 | `CountByStatus` | [`ComputeStatistics`](./compute-statistics) |
| --- | --- | --- |
| 返回 | 单个 `int` | 全维度 `*Statistics` |
| 适用 | 只需某一状态数 | 需全状态分布 |

需生成「状态分布饼图」用 `ComputeStatistics().ByStatus`；只需「有效弱点数」用本函数。

## ⚠️ 注意事项

::: warning 零值状态
未设置 `Status` 的弱点会被计入零值枚举的计数。若零值不是合法状态，报表里会出现「未知」数。
:::

::: tip 不依赖索引
本函数遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[Status](./enum-status)
- 全量统计：[计算统计](./compute-statistics)
- 按状态查找：[按状态查找](./find-by-status)
- 源文件：[`stats.go`](https://github.com/scagogogo/cwe-skills/blob/main/stats.go)
