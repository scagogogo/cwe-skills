---
title: CWENotFoundError 未找到错误
outline: [2, 3]
---

# ⚠️ CWENotFoundError — CWE 条目未找到错误

当 API 成功响应但返回的弱点/类别/视图列表为空时抛出。语义上是「404 等价物」，但源于业务层判空而非 HTTP 状态码。

源文件：`errors.go`。

## 🧱 结构体定义

```go
type CWENotFoundError struct {
    *CWEError
    ID int
}
```

| 字段 | 说明 |
| --- | --- |
| `CWEError.Code` | `"CWE_NOT_FOUND"` |
| `CWEError.Message` | `"CWE条目未找到"` |
| `CWEError.Detail` | `CWE ID: {id}` |
| `ID` | 未找到的 CWE ID（`int`） |

## 🏗️ 构造函数

```go
func NewCWENotFoundError(id int) *CWENotFoundError
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `int` | 未找到的 CWE ID |

返回 `*CWENotFoundError`。

## 🔍 触发点

在 [`GetWeakness`](./api-get-weakness) / [`GetCategory`](./api-get-category) / [`GetView`](./api-get-view) 中，当响应 `Data` 解析出的列表 `len == 0` 时返回此错误。

::: tip 与 APIError 的区别
`APIError` 由 HTTP 非 2xx 触发（如真的 404）；`CWENotFoundError` 是 2xx 但业务空数据。两者都表示「没有这条」，但来源不同——前者是传输层，后者是业务层。
:::

## 🚀 可运行示例

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewAPIClient()
    defer client.Close()

    _, err := client.GetWeakness(context.Background(), 999999)
    if err != nil {
        var notFound *cweskills.CWENotFoundError
        if errors.As(err, &notFound) {
            fmt.Printf("CWE-%d 不存在\n", notFound.ID)
        }
    }
}
```

::: warning 不应重试
条目不存在是稳定状态，重试无意义。捕获后应返回 404 或在本地注册表中标记为「已知缺失」。
:::

## 📚 相关链接

- [错误体系概览](./errors) | [InvalidCWEIDError](./invalid-cwe-id-error) | [APIError](./api-error) | [GetWeakness](./api-get-weakness)
