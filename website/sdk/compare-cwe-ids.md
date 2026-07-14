---
title: CompareCWEIDs 比较两个 CWE ID
outline: [2, 3]
---

# 🆔 CompareCWEIDs — 按数字大小比较两个 CWE ID

## 📋 函数签名

```go
func CompareCWEIDs(a, b string) (int, error)
```

## 📖 说明

`CompareCWEIDs` 把两个 CWE ID 字符串各自解析为整数，再按数值大小比较，返回 `-1 / 0 / 1`。它遵循 `sort.Compare` 风格约定，可直接用于排序比较函数。

::: tip 为什么不直接比较字符串？
字符串比较 `"CWE-100" < "CWE-79"` 为 `true`（按字典序 `'1' < '7'`），但数值上 `100 > 79`。CWE ID 必须按数字比较，本函数正是为此而存在。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `a` | `string` | 第一个 CWE ID |
| `b` | `string` | 第二个 CWE ID |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 比较结果 | `int` | `a < b` → `-1`；`a == b` → `0`；`a > b` → `1` |
| `error` | `error` | 任一参数无法解析时返回包装后的错误 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	cmp, _ := cweskills.CompareCWEIDs("CWE-79", "CWE-89")
	fmt.Println(cmp) // -1

	cmp, _ = cweskills.CompareCWEIDs("CWE-79", "CWE-79")
	fmt.Println(cmp) // 0

	cmp, _ = cweskills.CompareCWEIDs("CWE-89", "CWE-79")
	fmt.Println(cmp) // 1

	// 字典序陷阱：字符串比较会出错，本函数正确
	cmp, _ = cweskills.CompareCWEIDs("CWE-100", "CWE-79")
	fmt.Println(cmp) // 1 （100 > 79）

	// 错误情况
	_, err := cweskills.CompareCWEIDs("abc", "CWE-79")
	fmt.Println(err) // 比较失败，第一个CWE ID无效: ...
}
```

::: details 用于排序
```go
ids := []string{"CWE-89", "CWE-79", "CWE-100"}
sort.Slice(ids, func(i, j int) bool {
	cmp, _ := cweskills.CompareCWEIDs(ids[i], ids[j])
	return cmp < 0
})
// ids = [CWE-79 CWE-89 CWE-100]
```
:::

## ⚠️ 常见错误

::: warning 两个参数都会被校验
若 `a` 或 `b` 任一无法解析，函数返回 `0` 和错误。**此时不要使用返回的 `0`**——它不代表「相等」，而是错误占位。务必先判 error。
:::

::: details 错误信息区分先后
错误信息会指明是第几个参数无效：
- 第一个无效：`"比较失败，第一个CWE ID无效: ..."`
- 第二个无效：`"比较失败，第二个CWE ID无效: ..."`
便于定位是哪个输入出了问题。
:::

## 🔗 相关链接

- 底层依赖：[ParseCWEID](./parse-cwe-id)
- 整数比较也可先解析再比：`ParseCWEID` 后直接 `int` 比较
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
