---
title: 去重
outline: [2, 3]
---

# 🧹 去重

`Deduplicate` 对 `[]*CWE` 切片按 ID 去重，保留首次出现的元素，返回新切片。用于合并多个查询结果时消除重复条目。

## 📐 函数签名

```go
func Deduplicate(cwes []*CWE) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cwes` | `[]*CWE` | 可能含重复的切片 |
| 返回 | `[]*CWE` | 去重后的新切片，按首次出现顺序排列 |

::: tip 去重依据
按 `CWE.ID` 去重。ID 相同即视为同一条目，仅保留第一次出现的指针实例。
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
	r.Register(cweskills.NewCWE(79, "XSS"))
	r.Register(cweskills.NewCWE(89, "SQLi"))

	// 模拟两个查询结果合并
	resultA := cweskills.FindByKeyword(r, "xss")
	resultB := cweskills.FindByKeyword(r, "sql")
	merged := append(resultA, resultB...)

	deduped := cweskills.Deduplicate(merged)
	fmt.Println("合并前:", len(merged), "去重后:", len(deduped))
}
```

## ⚠️ 注意事项

::: warning 保留首次出现
重复 ID 出现多次时，只保留**第一次**的 `*CWE` 指针。若同一 ID 对应不同指针实例（理论上不应发生），以首次为准。
:::

::: details 顺序保持
`Deduplicate` 按输入切片的顺序输出，不排序。如需排序，先去重再调 [SortByID](./sort)。
:::

## 🔗 相关链接

- 合并来源：[搜索概览](./search)
- 排序：[排序](./sort)
- 注册表查重：[Register](./registry-operations)（ID 重复报错）
- 源文件：[`filter.go`](https://github.com/scagogogo/cwe-skills/blob/main/filter.go)
