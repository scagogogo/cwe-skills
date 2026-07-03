---
title: BasicFetcher 基础获取器
outline: [2, 3]
---

# 📥 BasicFetcher — 基础获取器

`BasicFetcher` 是最简单的 `DataFetcher` 实现：内嵌一个 `*APIClient`，单条获取弱点、类别、视图，或一次性拿弱点 + 直接父子关系。它实现了 [`DataFetcher`](./data-fetcher) 接口的 `Fetch` 方法。

源文件：`data_fetcher.go`。

## 🧱 结构体定义

```go
type BasicFetcher struct {
    client *APIClient
}
```

## 🏗️ 构造函数

```go
func NewBasicFetcher(client *APIClient) *BasicFetcher
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `client` | `*APIClient` | API 客户端；传 `nil` 会自动 `NewAPIClient()` |

返回 `*BasicFetcher`。

## 📤 方法

| 方法 | 签名 | 委托 | 文档 |
| --- | --- | --- | --- |
| `Fetch` | `(ctx, id int) (*CWE, error)` | `client.GetWeakness` | [GetWeakness](./api-get-weakness) |
| `FetchCategory` | `(ctx, id int) (*Category, error)` | `client.GetCategory` | [GetCategory](./api-get-category) |
| `FetchView` | `(ctx, id int) (*View, error)` | `client.GetView` | [GetView](./api-get-view) |
| `FetchWithRelations` | `(ctx, id int, viewID ...int) (*CWE, error)` | 组合 `GetWeakness`+`GetParents`+`GetChildren` | 见下 |

## FetchWithRelations

```go
func (f *BasicFetcher) FetchWithRelations(ctx context.Context, id int, viewID ...int) (*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `int` | 目标 CWE ID |
| `viewID` | `...int` | 可选视图过滤，透传给 `GetParents`/`GetChildren` |

流程：

1. `GetWeakness(ctx, id)` 拿到本体。
2. `GetParents(ctx, id, viewID...)` 拿父级，每条 `Nature` 改写为 `RelationshipChildOf`，追加到 `cwe.Relationships`。
3. `GetChildren(ctx, id, viewID...)` 拿子级，每条 `Nature` 改写为 `RelationshipParentOf`，追加到 `cwe.Relationships`。
4. 父/子获取出错时**静默跳过**（`err == nil` 才追加），本体仍返回。

::: tip 一步拿到邻居
`FetchWithRelations` 等价于「本体 + 一层父 + 一层子」三次 API 调用的组合，是浏览单条弱点时最常用的便捷方法。需要多层血缘用 [`TreeFetcher`](./tree-fetcher)。
:::

::: warning 关系 Nature 被重写
`GetParents`/`GetChildren` 返回的 `Nature` 是服务端原始值；`FetchWithRelations` 会按方向**覆盖**为 `RelationshipChildOf`/`RelationshipParentOf`，便于在本体上统一表达「我 -> 父」「我 -> 子」。若需要原始 Nature，请直接调 `GetParents`。
:::

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

    fetcher := cwe.NewBasicFetcher(client)

    w, err := fetcher.FetchWithRelations(context.Background(), 79, 1000)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("CWE-79: %s\n", w.Name)
    fmt.Printf("关系数: %d\n", len(w.Relationships))
    for _, rel := range w.Relationships {
        fmt.Printf("  %s CWE-%d\n", rel.Nature, rel.CWEID)
    }
}
```

## 📚 相关链接

- [DataFetcher 接口](./data-fetcher) | [GetWeakness](./api-get-weakness) | [父子关系](./api-parents-children) | [TreeFetcher](./tree-fetcher)
