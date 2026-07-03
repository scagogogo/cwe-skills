---
title: ParseCWEID 解析 CWE ID
outline: [2, 3]
---

# 🆔 ParseCWEID — 将 CWE ID 字符串解析为整数

## 📋 函数签名

```go
func ParseCWEID(id string) (int, error)
```

## 📖 说明

`ParseCWEID` 是 cwe_utils 工具链的**底层解析函数**。它接受各种常见写法的 CWE ID，返回其数字部分。该函数被 [`FormatCWEID`](./format-cwe-id)、[`IsCWEID`](./is-cwe-id)、[`ValidateCWEID`](./validate-cwe-id)、[`CompareCWEIDs`](./compare-cwe-ids) 共同复用，是整个工具集的基石。

支持下列输入格式：

| 格式 | 示例 | 说明 |
| --- | --- | --- |
| 标准格式 | `CWE-79` | MITRE 官方写法 |
| 无连字符 | `CWE79` | 紧凑写法 |
| 带空格 | `CWE 79`、`cwe 79` | 文本中常见 |
| 大小写不敏感 | `cwe-79`、`Cwe-79` | 任一大小写 |
| 前导零 | `CWE-079`、`079` | 数字去前导零 |
| 纯数字 | `79` | 仅本函数支持，提取类不支持 |

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 需要解析的 CWE ID 字符串 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 数字 ID | `int` | 解析出的整数，如 `79` |
| `error` | `error` | 失败时返回 `InvalidCWEIDError`，见 [InvalidCWEIDError](./invalid-cwe-id-error) |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cases := []string{"CWE-79", "79", "cwe-079", "CWE79", "cwe 79"}
	for _, c := range cases {
		num, err := cweskills.ParseCWEID(c)
		fmt.Printf("%-10q => %d, %v\n", c, num, err)
	}
	// "CWE-79"   => 79, <nil>
	// "79"       => 79, <nil>
	// "cwe-079"  => 79, <nil>
	// "CWE79"    => 79, <nil>
	// "cwe 79"   => 79, <nil>

	_, err := cweskills.ParseCWEID("")
	fmt.Println(err) // InvalidCWEIDError

	_, err = cweskills.ParseCWEID("abc")
	fmt.Println(err) // InvalidCWEIDError
}
```

## ⚠️ 常见错误

::: warning 解析失败的几种情况
1. **空字符串** `""` → 立即返回 `InvalidCWEIDError`
2. **非数字文本** `"abc"`、`"CWE-abc"` → 正则与纯数字解析均失败
3. **数字 ≤ 0** `"CWE-0"`、`"-1"` → 校验失败，CWE ID 必须 > 0
:::

::: details 为什么 0 也不合法？
CWE 编号从 1 开始（实际从个位数起），`0` 不是任何真实弱点的编号。SDK 主动拒绝 `≤ 0` 的值，避免下游产生空对象。
:::

## 🔗 相关链接

- 上层封装：[FormatCWEID](./format-cwe-id)、[IsCWEID](./is-cwe-id)、[ValidateCWEID](./validate-cwe-id)
- 整数→字符串的反向操作：[FormatCWEIDFromInt](./format-cwe-id-from-int)
- 错误类型：[InvalidCWEIDError](./invalid-cwe-id-error)
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
