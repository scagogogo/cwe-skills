---
title: ValidationError 验证错误
outline: [2, 3]
---

# ⚠️ ValidationError — 模型验证失败错误

当字段值不符合约束条件（如 `nil`、空字符串）时抛出。`ValidationError` 携带 `Field` 与 `Value`，定位到具体违规字段。

源文件：`errors.go`。

## 🧱 结构体定义

```go
type ValidationError struct {
    *CWEError
    Field string
    Value string
}
```

| 字段 | 说明 |
| --- | --- |
| `CWEError.Code` | `"VALIDATION_ERROR"` |
| `CWEError.Message` | `"模型验证失败"` |
| `CWEError.Detail` | `字段 {field} 的值 {value} 无效` |
| `Field` | 验证失败的字段名 |
| `Value` | 验证失败的值 |

## 🏗️ 构造函数

```go
func NewValidationError(field, value string) *ValidationError
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `field` | `string` | 失败字段名 |
| `value` | `string` | 失败的值（字符串描述） |

返回 `*ValidationError`。

## 🔍 触发点

- [`XMLParser.Parse`](./xml-parse)：`reader` 为 `nil` → `ValidationError(field="reader", value="nil")`
- `XMLParser.ParseFile`：`path` 为空 → `ValidationError(field="path", value="empty")`
- `XMLParser.ParseBytes`：`data` 为空 → `ValidationError(field="data", value="empty")`
- [`MultipleFetcher.FetchMultipleToRegistry`](./multiple-fetcher)：`registry` 为 `nil` → `ValidationError(field="registry", value="nil")`

::: tip 区分入参与解析错误
`ValidationError` 表征「调用方传入了不符合前置条件的参数」（nil/空），属于编程错误，应在代码中规避而非运行时重试。XML 内容本身的语法问题走 [`ParseError`](./parse-error)。
:::

## 🚀 可运行示例

```go
package main

import (
    "errors"
    "fmt"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    parser := cwe.NewXMLParser()
    _, err := parser.ParseFile("")
    if err != nil {
        var ve *cwe.ValidationError
        if errors.As(err, &ve) {
            fmt.Printf("字段 %q 值 %q 无效\n", ve.Field, ve.Value)
        }
    }
}
```

::: warning 不应重试
`ValidationError` 是稳定的输入错误，重试无意义。捕获后修正调用方代码即可。
:::

## 📚 相关链接

- [错误体系概览](./errors) | [ParseError](./parse-error) | [XMLParser](./xml-parser) | [MultipleFetcher](./multiple-fetcher)
