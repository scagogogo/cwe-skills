---
title: 按状态查找
outline: [2, 3]
---

# 🔍 按状态查找

`FindByStatus` 按 [Status](./enum-status) 枚举筛选弱点，返回指定状态（如 Stable、Deprecated）的全部弱点。常用于排除已弃用条目。

## 📐 函数签名

```go
func FindByStatus(r *Registry, status Status) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `status` | [`Status`](./enum-status) | 状态枚举 |
| 返回 | `[]*CWE` | 该状态的全部弱点 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	stable := cweskills.NewCWE(79, "XSS")
	stable.Status = cweskills.StatusStable
	r.Register(stable)

	deprecated := cweskills.NewCWE(100, "Old")
	deprecated.Status = cweskills.StatusDeprecated
	r.Register(deprecated)

	stables := cweskills.FindByStatus(r, cweskills.StatusStable)
	fmt.Println("Stable:", len(stables))         // 1
	fmt.Println("Deprecated:", len(cweskills.FindByStatus(r, cweskills.StatusDeprecated))) // 1
}
```

## 🆚 与 Filter 的取舍

| 维度 | `FindByStatus` | [`Filter`](./filter) + `FilterOption.Status` |
| --- | --- | --- |
| 单条件 | 简洁 | 稍冗长 |
| 多条件 | 不支持 | 支持 AND 组合 |

只按状态单维度查时用 `FindByStatus`；需组合其它维度用 `Filter`。

## ⚠️ 注意事项

::: tip 典型用途
生成「有效弱点清单」时常用 `FindByStatus(r, StatusStable)` 排除 Deprecated/Withdrawn 条目。
:::

::: warning 不依赖索引
`FindByStatus` 遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[Status](./enum-status)
- 多条件过滤：[Filter](./filter)
- 状态统计：[CountByStatus](./count-by-status)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
