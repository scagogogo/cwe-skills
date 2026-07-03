---
title: CWEError 错误根类型
outline: [2, 3]
---

# ⚠️ CWEError — 错误根类型

`CWEError` 是 `cweskills` 包所有错误的**基础类型**。它提供统一的 `Code`/`Message`/`Detail`/`Err` 四元组，实现 `error` 接口并支持 `Unwrap` 链式查找。六种细分错误都内嵌 `*CWEError`。

源文件：`errors.go`。

## 🧱 结构体定义

```go
type CWEError struct {
    Code    string
    Message string
    Detail  string
    Err     error
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Code` | `string` | 错误码，如 `"API_ERROR"` |
| `Message` | `string` | 人类可读的错误概要 |
| `Detail` | `string` | 详细描述，含上下文 |
| `Err` | `error` | 被包装的内部错误（可 `nil`） |

## 📤 方法

### Error

```go
func (e *CWEError) Error() string
```

按字段存在性分层格式化：

- 有 `Err`：`cwe: [{Code}] {Message}: {Detail}: {Err}`
- 无 `Err` 但有 `Detail`：`cwe: [{Code}] {Message}: {Detail}`
- 仅 `Code`+`Message`：`cwe: [{Code}] {Message}`

### Unwrap

```go
func (e *CWEError) Unwrap() error
```

返回 `Err` 字段，使 `errors.Is`/`errors.As` 能穿透到被包装的底层错误。

::: tip Unwrap 是链式查找的关键
`Unwrap` 让 `errors.Is(err, sql.ErrNoRows)` 这类判断生效——只要底层错误链中含目标。`CWEError` 的 `Err` 字段就是链的下一环。
:::

## 🚀 直接构造与判断

通常你不会直接 `New` 一个 `CWEError`，而是用各细分错误的构造函数。但理解它的结构有助于自定义错误：

```go
package main

import (
    "errors"
    "fmt"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    base := &cwe.CWEError{
        Code:    "CUSTOM",
        Message: "自定义错误",
        Detail:  "上下文信息",
    }
    fmt.Println(base.Error())
    // 输出: cwe: [CUSTOM] 自定义错误: 上下文信息

    var target *cwe.CWEError
    fmt.Println(errors.As(base, &target)) // true
}
```

::: warning 细分错误内嵌的是指针
`InvalidCWEIDError` 等内嵌的是 `*CWEError`（指针），构造函数里会创建该指针。`errors.As(err, &notFound)` 匹配的是外层细分类型；若想匹配 `CWEError` 本身，用 `errors.As(err, &base)` 也能命中，因为 `Unwrap` 链包含它。
:::

## 📚 相关链接

- [错误体系概览](./errors) | [APIError](./api-error) | [InvalidCWEIDError](./invalid-cwe-id-error) | [ParseError](./parse-error)
