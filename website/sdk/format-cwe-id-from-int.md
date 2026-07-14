---
title: FormatCWEIDFromInt 整数转 CWE ID
outline: [2, 3]
---

# 🆔 FormatCWEIDFromInt — 将整数格式化为 CWE ID 字符串

## 📋 函数签名

```go
func FormatCWEIDFromInt(id int) string
```

## 📖 说明

`FormatCWEIDFromInt` 是最轻量的格式化函数：输入一个整数，输出 `CWE-NNN` 字符串。内部实现就是一行：

```go
func FormatCWEIDFromInt(id int) string {
	return fmt.Sprintf("CWE-%d", id)
}
```

::: tip 为什么不报错？
调用方传入 `int` 通常意味着已经过解析或来自数据库主键，**默认信任是合法整数**。本函数不做任何校验，因此：
- 不返回 error，签名简洁
- 负数或零也会原样输出（如 `CWE-0`、`CWE--1`）
- 适合在已知合法的热路径上使用

若输入来自外部不可信来源，请改用 [`FormatCWEID`](./format-cwe-id)（带校验）。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `int` | 整数形式的 CWE ID |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 字符串 | `string` | 形如 `CWE-79` |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	fmt.Println(cweskills.FormatCWEIDFromInt(79))   // CWE-79
	fmt.Println(cweskills.FormatCWEIDFromInt(1000)) // CWE-1000
	fmt.Println(cweskills.FormatCWEIDFromInt(1))    // CWE-1

	// 注意：不做校验，负数/零会原样输出
	fmt.Println(cweskills.FormatCWEIDFromInt(0))    // CWE-0
	fmt.Println(cweskills.FormatCWEIDFromInt(-1))   // CWE--1
}
```

::: details 在 CWE 结构体方法中的内部使用
`CWE` 类型的 `CWEID()` 方法正是调用本函数：
```go
func (c *CWE) CWEID() string {
	return FormatCWEIDFromInt(c.ID)
}
```
详见 [CWE 弱点](./cwe-struct)。
:::

## ⚖️ FormatCWEIDFromInt vs FormatCWEID

| 维度 | `FormatCWEIDFromInt` | [`FormatCWEID`](./format-cwe-id) |
| --- | --- | --- |
| 输入 | `int` | `string` |
| 返回 | `string` | `(string, error)` |
| 校验 | ❌ 无 | ✅ 有，失败返回 `InvalidCWEIDError` |
| 适用 | 已知合法整数（DB 主键、解析结果） | 不可信字符串输入 |

## ⚠️ 常见错误

::: warning 不要用负数或零
本函数不拦截非法值。若上游可能传入 `≤ 0` 的整数，应在调用前自行校验，或改用 `FormatCWEID` 让 SDK 帮你判断。
:::

## 🔗 相关链接

- 带校验的字符串版本：[FormatCWEID](./format-cwe-id)
- 反向操作（字符串→整数）：[ParseCWEID](./parse-cwe-id)
- 在 `CWE.CWEID()` 中的使用：[CWE 弱点](./cwe-struct)
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
