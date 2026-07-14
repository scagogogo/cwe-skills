---
title: TreeFetcher 树获取器
outline: [2, 3]
---

# 🌳 TreeFetcher — 递归树获取器

`TreeFetcher` 递归拉取 CWE 的祖先链与后代链，构建完整血缘树并写入内部 `Registry`。它通过 `maxDepth` 控制递归层数，避免无限递归与请求爆炸。

源文件：`data_fetcher.go`。内部用 `fetchAncestorsRecursive` / `fetchDescendantsRecursive` 两个私有递归方法。

## 🧱 结构体定义

```go
type TreeFetcher struct {
    client   *APIClient
    registry *Registry
    maxDepth int
}
```

## 🏗️ 构造函数

```go
func NewTreeFetcher(client *APIClient, registry *Registry, maxDepth int) *TreeFetcher
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `client` | `*APIClient` | API 客户端；`nil` 自动新建 |
| `registry` | `*Registry` | 结果存储；`nil` 自动新建 |
| `maxDepth` | `int` | 最大递归深度；`<= 0` 自动设为 `10` |

返回 `*TreeFetcher`。

## 📤 方法

| 方法 | 签名 | 说明 |
| --- | --- | --- |
| `FetchWithAncestors` | `(ctx, id int) error` | 递归获取指定 ID 的全部祖先 |
| `FetchWithDescendants` | `(ctx, id int) error` | 递归获取指定 ID 的全部后代 |
| `FetchFullTree` | `(ctx, rootID int) error` | 先祖先后后代，构建完整树 |
| `GetRegistry` | `() *Registry` | 返回内部注册表，供查询结果 |

## 🔁 递归逻辑

### fetchAncestorsRecursive(ctx, id, depth)

1. `depth >= maxDepth` → 返回（停止）。
2. `registry.Contains(id)` → 已获取过，跳过。
3. `GetWeakness(ctx, id)` 拿本体并 `Register`。
4. `GetParents(ctx, id)` 拿父级，每条 `Nature` 改写为 `RelationshipChildOf`，追加到本体 `Relationships`。
5. 对每个父级 `rel.CWEID` 递归 `depth+1`。

### fetchDescendantsRecursive(ctx, id, depth)

1. `depth >= maxDepth` → 返回。
2. 本体未在 registry → `GetWeakness` 并 `Register`。
3. `GetChildren(ctx, id)` 拿子级，每条 `Nature` 改写为 `RelationshipParentOf`，追加到本体。
4. 对每个子级递归 `depth+1`。

::: tip 防止重复请求
两个递归都先用 `registry.Contains(id)` 判断是否已获取，避免对同一节点重复发起 `GetWeakness` 请求。这使得 `FetchFullTree` 不会因祖先与后代路径交叉而重复拉取。
:::

::: danger maxDepth 是硬上限
默认 `10` 层足以覆盖大多数 CWE 谱系，但根节点（如 CWE-1000 研究型视图根）的后代可能极多。每层每个节点都触发 `GetWeakness` + `GetChildren` 两次 API 调用，深度过大时请求数指数增长。生产环境建议 `maxDepth` 不超过 `5`，并配合限流。
:::

::: warning 不接受 viewID
`TreeFetcher` 的递归调用不传 `viewID`，拉取的是全局关系。若需限定视图，请改用 [`BasicFetcher.FetchWithRelations`](./basic-fetcher) 或直接调 [`GetParents`/`GetChildren`](./api-parents-children)。
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

    fetcher := cweskills.NewTreeFetcher(client, nil, 3)
    if err := fetcher.FetchFullTree(context.Background(), 79); err != nil {
        log.Fatal(err)
    }
    registry := fetcher.GetRegistry()
    fmt.Printf("树中共 %d 个节点\n", registry.Size())
    if w, ok := registry.Get(79); ok {
        fmt.Printf("CWE-79 关系数: %d\n", len(w.Relationships))
    }
}
```

## 📚 相关链接

- [DataFetcher 接口](./data-fetcher) | [祖先后代关系](./api-ancestors-descendants) | [父子关系](./api-parents-children) | [Registry 模型](./model)
