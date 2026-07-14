---
title: FormatCWEID 格式化 CWE ID
outline: [2, 3]
---

# 🆔 FormatCWEID — 将 CWE ID 字符串格式化为标准形式

## 📋 函数签名

```go
func FormatCWEID(id string) (string, error)
```

## 📖 说明

`FormatCWEID` 把任意合规写法的 CWE ID 统一为 MITRE 官方标准形式 `CWE-NNN`：去空格、统一大写、数字去前导零。它内部先调用 [`ParseCWEID`](./parse-cwe-id) 取得整数，再用 `fmt.Sprintf("CWE-%d", num)` 输出。

::: tip 与 FormatCWEIDFromInt 的区别
- 本函数输入是**字符串**，需做解析与校验，可能返回错误。
- [`FormatCWEIDFromInt`](./format-cwe-id-from-int) 输入是 **`int`**，绝不报错，适合已知是合法整数的场景。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 需要格式化的 CWE ID 字符串 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 标准形式 | `string` | 形如 `CWE-79`；失败时为空串 `""` |
| `error` | `error` | 输入不合法时返回 `InvalidCWEIDError` |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	cases := []string{"79", "cwe-79", "CWE79", "CWE 79", "cwe-079"}
	for _, c := range cases {
		out, err := cweskills.FormatCWEID(c)
		fmt.Printf("%-10q => %q, %v\n", c, out, err)
	}
	// "79"       => "CWE-79", <nil>
	// "cwe-79"   => "CWE-79", <nil>
	// "CWE79"    => "CWE-79", <nil>
	// "CWE 79"   => "CWE-79", <nil>
	// "cwe-079"  => "CWE-79", <nil>

	_, err := cweskills.FormatCWEID("")
	fmt.Println(err) // InvalidCWEIDError
}
```

## ⚠️ 常见错误

::: warning 失败时返回空串
当输入无法解析时，返回的第一个值是空字符串 `""`。务必**先检查 error 再使用返回值**，避免把空串带入下游逻辑。
:::

::: details 典型误用：忽略错误直接使用
```go
// ❌ 错误示范
id, _ := cweskills.FormatCWEID(userInput)
query := "SELECT * FROM cwe WHERE id = '" + id + "'"
// userInput 非法时 id 为 ""，可能拼出非法 SQL

// ✅ 正确做法
id, err := cweskills.FormatCWEID(userInput)
if err != nil {
    return fmt.Errorf("无效的 CWE ID: %w", err)
}
```
:::

## 🔗 相关链接

- 底层依赖：[ParseCWEID](./parse-cwe-id)
- 整数版本：[FormatCWEIDFromInt](./format-cwe-id-from-int)
- 错误类型：[InvalidCWEIDError](./invalid-cwe-id-error)
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
