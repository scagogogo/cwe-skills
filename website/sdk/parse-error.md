---
title: ParseError 解析错误
outline: [2, 3]
---

# ⚠️ ParseError — 数据解析失败错误

当 XML、JSON 或其他格式的数据解析失败时抛出。`ParseError` 携带 `Detail` 与 `Offset`，定位解析失败的位置。

源文件：`errors.go`。

## 🧱 结构体定义

```go
type ParseError struct {
    *CWEError
    Offset int64
}
```

| 字段 | 说明 |
| --- | --- |
| `CWEError.Code` | `"PARSE_ERROR"` |
| `CWEError.Message` | `"数据解析失败"` |
| `CWEError.Detail` | 解析失败的详细描述（含原始 `err`） |
| `Offset` | 解析失败的位置偏移量（`int64`，多数场景为 `0`） |

## 🏗️ 构造函数

```go
func NewParseError(detail string, offset int64) *ParseError
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `detail` | `string` | 解析失败的详细描述 |
| `offset` | `int64` | 解析失败的位置偏移量 |

返回 `*ParseError`。

## 🔍 触发点

- [`XMLParser.Parse`](./xml-parse)：`xml.Decoder.Decode` 失败、`os.Open` 失败 → `ParseError`
- [`HTTPClient`](./http-client) `Get`/`Post`/`PostForm`：`json.Unmarshal` 失败 → `ParseError`
- [`GetWeakness`](./api-get-weakness) / `GetCategory` / `GetView` / `GetCWEs` / `GetVersion`：响应 `Data` 二次反序列化失败 → `ParseError`
- 关系查询 `getRelations`：双重解析都失败 → `ParseError`

::: tip 当前 Offset 多为 0
SDK 内部构造 `ParseError` 时 `offset` 普遍传 `0`（不追踪字节偏移）。`Offset` 字段保留供未来扩展，当前不应依赖其精确值。
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
    _, err := parser.ParseBytes([]byte("<not-valid-xml>"))
    if err != nil {
        var pe *cwe.ParseError
        if errors.As(err, &pe) {
            fmt.Println("错误码:", pe.Code)
            fmt.Println("详情:", pe.Detail)
            fmt.Println("偏移:", pe.Offset)
        }
    }
}
```

::: warning 解析失败可能是数据源问题
JSON 解析失败常因 MITRE API 返回了非预期结构（如错误页 HTML），或 `baseURL` 指向了错误端点。排查时先打印响应原始字节（用 [`GetRaw`](./http-methods)）确认数据形态。
:::

## 📚 相关链接

- [错误体系概览](./errors) | [ValidationError](./validation-error) | [XMLParser](./xml-parser) | [HTTP 请求方法](./http-methods)
