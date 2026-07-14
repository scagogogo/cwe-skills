---
title: InvalidCWEIDError 无效 ID 错误
outline: [2, 3]
---

# ⚠️ InvalidCWEIDError — 无效 CWE ID 错误

当 CWE ID 不符合预期（数字部分 `<= 0` 或格式非法）时抛出。这是**入参错误**，不应重试，应修正调用方逻辑。

源文件：`errors.go`。

## 🧱 结构体定义

```go
type InvalidCWEIDError struct {
    *CWEError
    ID string
}
```

| 字段 | 说明 |
| --- | --- |
| `CWEError.Code` | `"INVALID_CWE_ID"` |
| `CWEError.Message` | `"CWE ID格式无效"` |
| `CWEError.Detail` | `输入值: %q`（含原始输入） |
| `ID` | 原始无效输入字符串 |

## 🏗️ 构造函数

```go
func NewInvalidCWEIDError(id string) *InvalidCWEIDError
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 无效的 CWE ID 字符串 |

返回 `*InvalidCWEIDError`，`Detail` 形如 `输入值: "0"` 或 `输入值: "abc"`。

## 🔍 触发点

- [`GetWeakness`](./api-get-weakness) / [`GetCategory`](./api-get-category) / [`GetView`](./api-get-view)：`id <= 0`
- [`GetCWEs`](./api-get-cwes)：`ids` 中任一 `id <= 0`
- [`GetParents`](./api-parents-children) / `GetChildren` / `GetAncestors` / `GetDescendants`：`id <= 0`

::: tip 错误码固定
`Code` 恒为 `"INVALID_CWE_ID"`，可作为 `errors.As` 之外的二次判别依据。
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

    _, err := client.GetWeakness(context.Background(), 0)
    if err != nil {
        var invalid *cweskills.InvalidCWEIDError
        if errors.As(err, &invalid) {
            fmt.Println("错误码:", invalid.Code)
            fmt.Println("无效输入:", invalid.ID)
        }
    }
}
```

::: warning 不要对 InvalidCWEIDError 重试
这是调用方传错参数，重试只会再次失败。捕获后应记录、修正入参或返回 4xx 给上游。
:::

## 📚 相关链接

- [错误体系概览](./errors) | [CWEError 根类型](./cwe-error) | [CWENotFoundError](./cwe-not-found-error) | [GetWeakness](./api-get-weakness)
