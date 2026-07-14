---
title: ValidateCWEID 验证 CWE ID
outline: [2, 3]
---

# 🆔 ValidateCWEID — 验证 CWE ID 并返回详细错误

## 📋 函数签名

```go
func ValidateCWEID(text string) error
```

## 📖 说明

`ValidateCWEID` 对 CWE ID 做完整校验，校验通过返回 `nil`，失败返回带原因的 [`InvalidCWEIDError`](./invalid-cwe-id-error)。与只返回布尔的 [`IsCWEID`](./is-cwe-id) 相比，它**保留错误对象**，便于调用方记录日志、向用户展示或包装后向上抛出。

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `text` | `string` | 待验证的 CWE ID 字符串 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| `error` | `error` | 通过返回 `nil`；失败返回 `InvalidCWEIDError` |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	err := cweskills.ValidateCWEID("CWE-79")
	fmt.Println(err) // <nil>

	err = cweskills.ValidateCWEID("")
	fmt.Println(err) // InvalidCWEIDError: ...

	err = cweskills.ValidateCWEID("abc")
	fmt.Println(err) // InvalidCWEIDError: ...

	// 典型用法：把原因带向上层
	if err := cweskills.ValidateCWEID(input); err != nil {
		return fmt.Errorf("用户输入的 CWE ID 无效: %w", err)
	}
}
```

## ⚖️ IsCWEID vs ValidateCWEID

| 维度 | [`IsCWEID`](./is-cwe-id) | `ValidateCWEID` |
| --- | --- | --- |
| 返回值 | `bool` | `error` |
| 错误原因 | 无 | 有，`InvalidCWEIDError` |
| 适用 | 简单守卫 | 需要错误传播/展示 |

::: details 校验失败的具体原因有哪些？
`ValidateCWEID` 转发 `ParseCWEID` 的错误，失败情形包括：
1. 空字符串
2. 既不匹配 `CWE[-\s]?\d+` 也不是纯数字
3. 数字部分 `≤ 0`

错误对象是 `InvalidCWEIDError`，可通过 `errors.As` 取出原始输入做进一步处理。
:::

## ⚠️ 常见错误

::: warning 不要忽略返回的 error
```go
// ❌ 丢弃错误
_ = cweskills.ValidateCWEID(input)

// ✅ 检查并传播
if err := cweskills.ValidateCWEID(input); err != nil {
    return err
}
```
:::

## 🔗 相关链接

- 布尔版本：[IsCWEID](./is-cwe-id)
- 底层依赖：[ParseCWEID](./parse-cwe-id)
- 错误类型：[InvalidCWEIDError](./invalid-cwe-id-error)
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
