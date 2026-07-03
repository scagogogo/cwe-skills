---
title: GetCategory 获取类别详情
outline: [2, 3]
---

# 📥 GetCategory — 获取 CWE 类别详情

`GetCategory` 调用 `GET /cwe/category/{id}` 端点，返回指定 ID 的 CWE 类别（Category）。类别把相关弱点归组，本身不是弱点。

源文件：`api_client_cwe.go`。

## 📐 函数签名

```go
func (c *APIClient) GetCategory(ctx context.Context, id int) (*Category, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |
| `id` | `int` | 类别 ID 数字 |

返回值：

| 返回 | 说明 |
| --- | --- |
| `*Category` | 类别详情 |
| `error` | 请求或解析失败时非 nil |

## 🔁 内部流程

1. 校验 `id > 0`，否则返回 [`InvalidCWEIDError`](./invalid-cwe-id-error)。
2. 路径 `/cwe/category/{id}`，底层 `HTTPClient.Get` 发请求。
3. 先尝试 `[]Category`，再退化单条 `Category`。
4. 列表为空返回 [`CWENotFoundError`](./cwe-not-found-error)。

## 🚀 可运行示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cwe.NewAPIClient()
    defer client.Close()

    cat, err := client.GetCategory(context.Background(), 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Category-%d: %s\n", cat.ID, cat.Name)
    fmt.Println(cat.Description)
}
```

::: tip Category 与 CWE 的区别
`Category` 是「分组容器」，通过 `Relationships` 引用其下属弱点；它没有 `Abstraction`、`Structure` 等弱点专属字段。详见 [Category 结构体](./category)。
:::

## 📚 相关链接

- [GetWeakness](./api-get-weakness) | [GetView](./api-get-view) | [Category 结构体](./category)
