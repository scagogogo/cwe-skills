---
title: 排序
outline: [2, 3]
---

# 🧹 排序

`SortByID`、`SortByName`、`SortByAbstraction` 三个函数对 `[]*CWE` 切片排序，返回新的有序切片。常用于在 [Filter](./filter) 之后稳定输出顺序。

## 📐 函数签名

```go
func SortByID(cwes []*CWE) []*CWE
func SortByName(cwes []*CWE) []*CWE
func SortByAbstraction(cwes []*CWE) []*CWE
```

| 函数 | 排序键 | 顺序 |
| --- | --- | --- |
| `SortByID` | `CWE.ID` | 升序（数字） |
| `SortByName` | `CWE.Name` | 升序（字典序） |
| `SortByAbstraction` | `CWE.Abstraction` | 按枚举序升序 |

::: tip 返回新切片
三个函数都返回一个新的 `[]*CWE`，**不修改**输入切片的元素顺序（内部通常复制后排序）。元素指针仍共享。
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
	r.Register(cweskills.NewCWE(89, "SQLi"))
	r.Register(cweskills.NewCWE(79, "XSS"))
	r.Register(cweskills.NewCWE(703, "Neutralization"))

	byID := cweskills.SortByID(r.GetAll())
	for _, c := range byID {
		fmt.Print(c.ID, " ") // 79 89 703 ...（升序）
	}
	fmt.Println()

	byName := cweskills.SortByName(r.GetAll())
	for _, c := range byName {
		fmt.Println(c.Name)
	}
}
```

## ⚠️ 注意事项

::: warning SortByName 大小写
字典序排序通常区分大小写，大写字母排在小写之前。若需不区分大小写，需自行排序。`SortByName` 的具体比较规则以源码实现为准。
:::

::: details 与 GetAll 配合
`GetAll()` 返回顺序不确定，输出前建议显式排序。`SortByID` 是最常用的稳定化手段。
:::

## 🔗 相关链接

- 待排序来源：[GetAll](./registry-operations)
- 过滤后排序：[Filter](./filter)
- 分组：[分组](./group-by)
- 源文件：[`filter.go`](https://github.com/scagogogo/cwe-skills/blob/main/filter.go)
