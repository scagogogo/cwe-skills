---
title: Filter 多条件过滤
outline: [2, 3]
---

# 🧹 Filter 多条件过滤

`Filter` 对一个 `[]*CWE` 切片按多个条件做 **AND 关系**过滤，返回同时满足全部条件的子集。它是组合查询的中枢——把 [FilterOption](./filter-option) 的各字段作为筛选项。

## 📐 函数签名

```go
func Filter(cwes []*CWE, opts ...FilterOption) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cwes` | `[]*CWE` | 待过滤的弱点切片（通常来自 `GetAll()`） |
| `opts` | `...FilterOption` | 一个或多个过滤选项；多选项间为 AND |
| 返回 | `[]*CWE` | 满足全部条件的弱点 |

::: tip 多选项的语义
传多个 `FilterOption` 时，条件之间是 **AND**。每个 `FilterOption` 内部各字段也按 AND 组合。空字段（零值）表示「不限制该维度」。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	c := cweskills.NewCWE(79, "XSS")
	c.Abstraction = cweskills.AbstractionBase
	c.Status = cweskills.StatusStable
	c.Description = "Cross-site Scripting"
	r.Register(c)
	r.BuildIndexes()

	// Base 且 Stable 且含 "xss"
	got := cweskills.Filter(r.GetAll(),
		cweskills.FilterOption{
			Abstraction: cweskills.AbstractionBase,
			Status:      cweskills.StatusStable,
			Keyword:     "xss",
		})
	fmt.Println(len(got)) // 1
}
```

## ⚠️ 注意事项

::: warning 零值字段不参与过滤
`FilterOption` 中零值字段（如 `""` 字符串、零枚举）被视为「不限制」。但要注意枚举零值可能与某个有效枚举冲突——若零值恰是合法选项，建议改用 `FindBy*` 单条件函数。
:::

::: details 输入不变性
`Filter` 不修改输入切片，返回新切片。但切片元素（`*CWE` 指针）是共享的，修改返回元素的字段会影响原数据。
:::

## 🔗 相关链接

- 选项结构：[FilterOption](./filter-option)
- 单维度查找：[搜索概览](./search)
- 排序结果：[排序](./sort)
- 源文件：[`filter.go`](https://github.com/scagogogo/cwe-skills/blob/main/filter.go)
