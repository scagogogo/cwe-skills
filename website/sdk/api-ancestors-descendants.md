---
title: GetAncestors GetDescendants 祖先后代
outline: [2, 3]
---

# 🔗 GetAncestors / GetDescendants — 祖先与后代关系

`GetAncestors` 与 `GetDescendants` 分别调用 `GET /cwe/{id}/ancestors` 和 `GET /cwe/{id}/descendants`，返回指定 CWE 条目的**递归**祖先或后代关系（即多跳路径上的全部节点）。两者不接受 `viewID`，返回的是全局关系链。

源文件：`api_client_relations.go`。同样委托 `getRelations(ctx, path, id, relType)`。

## 📐 函数签名

```go
func (c *APIClient) GetAncestors(ctx context.Context, id int) ([]Relationship, error)
func (c *APIClient) GetDescendants(ctx context.Context, id int) ([]Relationship, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |
| `id` | `int` | 目标 CWE ID |

返回值：

| 返回 | 说明 |
| --- | --- |
| `[]Relationship` | 递归关系列表 |
| `error` | 请求或解析失败时非 nil |

## 🔁 内部流程

1. 校验 `id > 0`，否则返回 [`InvalidCWEIDError`](./invalid-cwe-id-error)。
2. 路径 `/cwe/{id}/ancestors`（或 `descendants`）。
3. 调用 `getRelations`，采用与父子关系相同的双重解析策略（严格 `[]Relationship` → 宽松结构回退）。

::: tip 祖先/后代 vs 父子
[`GetParents`](./api-parents-children) 返回直接父级（一层）；`GetAncestors` 返回递归到根的全部祖先（多层）。后代同理。前者适合「找邻居」，后者适合「画血缘树」。
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

    ancestors, err := client.GetAncestors(context.Background(), 79)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("CWE-79 的祖先共 %d 个\n", len(ancestors))
    for _, rel := range ancestors {
        fmt.Printf("  <- CWE-%d (%s)\n", rel.CWEID, rel.Nature)
    }

    descendants, err := client.GetDescendants(context.Background(), 79)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("CWE-79 的后代共 %d 个\n", len(descendants))
}
```

::: danger 返回量可能很大
某些根级弱点的后代可达数百条，且每条关系都来自一次 HTTP 调用之后的批量返回。若只需有限深度，改用 [`TreeFetcher`](./tree-fetcher) 配合 `maxDepth` 控制递归层数。
:::

::: details 与 TreeFetcher 的取舍
`GetAncestors`/`GetDescendants` 由服务端一次性返回关系列表，但不返回各节点的完整弱点详情。若需要每条祖先/后代的 `Name`、`Description` 等字段，用 [`TreeFetcher.FetchWithAncestors`](./tree-fetcher) 会逐节点拉取并注册到 `Registry`。
:::

## 📚 相关链接

- [父子关系](./api-parents-children) | [TreeFetcher](./tree-fetcher) | [RelationshipNature 枚举](./enum-relationship-nature)
