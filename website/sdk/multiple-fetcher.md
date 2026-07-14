---
title: MultipleFetcher 批量获取器
outline: [2, 3]
---

# 📥 MultipleFetcher — 批量获取器

`MultipleFetcher` 专攻批量场景：一次请求拿多条弱点，或直接把批量结果灌入 [`Registry`](./model)。它内嵌一个 `*APIClient`，核心方法委托 [`GetCWEs`](./api-get-cwes)。

源文件：`data_fetcher.go`。

## 🧱 结构体定义

```go
type MultipleFetcher struct {
    client *APIClient
}
```

## 🏗️ 构造函数

```go
func NewMultipleFetcher(client *APIClient) *MultipleFetcher
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `client` | `*APIClient` | API 客户端；传 `nil` 自动 `NewAPIClient()` |

返回 `*MultipleFetcher`。

## 📤 方法

### FetchMultiple

```go
func (f *MultipleFetcher) FetchMultiple(ctx context.Context, ids []int) (map[string]*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ids` | `[]int` | CWE ID 列表 |

直接委托 `client.GetCWEs`，返回以 CWE ID 字符串为键的映射。只发**一次** HTTP 请求。

### FetchMultipleToRegistry

```go
func (f *MultipleFetcher) FetchMultipleToRegistry(ctx context.Context, ids []int, registry *Registry) error
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ids` | `[]int` | 要获取的 ID 列表 |
| `registry` | `*Registry` | 目标注册表 |

流程：

1. `registry` 为 `nil` → 返回 `ValidationError(field="registry", value="nil")`。
2. `client.GetCWEs(ctx, ids)` 拿到映射。
3. 遍历映射，对每个非空 `*CWE` 调 `registry.Register`，重复注册错误被静默忽略。

::: tip 构建离线缓存的最佳拍档
`FetchMultipleToRegistry` 把批量结果直接写入 `Registry`，配合 [`XMLParser`](./xml-parser) 预加载，可实现「XML 兜底 + API 补缺」的混合缓存：先解析全量 XML，再用 `MultipleFetcher` 把关注的新版本弱点覆盖进同一个 `Registry`。
:::

::: warning nil 条目会被跳过
`GetCWEs` 返回的 map 中，某些 ID 对应的值可能为 `nil`（API 未返回该条）。`FetchMultipleToRegistry` 会跳过这些空值，不会写入注册表，调用方需自行核对缺失项。
:::

## 🚀 可运行示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewAPIClient()
    defer client.Close()

    fetcher := cweskills.NewMultipleFetcher(client)

    // 1. 批量获取到 map
    ids := []int{79, 89, 119, 20}
    result, err := fetcher.FetchMultiple(context.Background(), ids)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("获取到:", len(result), "条")

    // 2. 批量灌入 Registry
    registry := cweskills.NewRegistry()
    if err := fetcher.FetchMultipleToRegistry(context.Background(), ids, registry); err != nil {
        log.Fatal(err)
    }
    if w, ok := registry.Get(79); ok {
        fmt.Println("Registry 中:", w.Name)
    }
}
```

## 📚 相关链接

- [DataFetcher 接口](./data-fetcher) | [GetCWEs](./api-get-cwes) | [Registry 模型](./model) | [BasicFetcher](./basic-fetcher)
