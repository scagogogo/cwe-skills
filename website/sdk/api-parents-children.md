---
title: GetParents GetChildren 父子关系
outline: [2, 3]
---

# 🔗 GetParents / GetChildren — CWE 父子关系

`GetParents` 与 `GetChildren` 分别调用 `GET /cwe/{id}/parents` 和 `GET /cwe/{id}/children`，返回指定 CWE 条目在某个视图下的直接父级或子级关系。两者都接受可选 `viewID` 限定范围。

源文件：`api_client_relations.go`。内部都委托给私有方法 `getRelations(ctx, path, id, relType)`。

## 📐 函数签名

```go
func (c *APIClient) GetParents(ctx context.Context, id int, viewID ...int) ([]Relationship, error)
func (c *APIClient) GetChildren(ctx context.Context, id int, viewID ...int) ([]Relationship, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |
| `id` | `int` | 目标 CWE ID |
| `viewID` | `...int` | 可变参数，最多取第一个非零值作为视图过滤 |

返回值：

| 返回 | 说明 |
| --- | --- |
| `[]Relationship` | 关系列表，每条含 `Nature`、`CWEID`、`ViewID` |
| `error` | 请求或解析失败时非 nil |

## 🔁 内部流程

1. 校验 `id > 0`，否则返回 [`InvalidCWEIDError`](./invalid-cwe-id-error)。
2. 拼接路径 `/cwe/{id}/parents`（或 `children`），当 `viewID[0] > 0` 时追加 `?view={viewID}`。
3. 调用 `getRelations`：先尝试反序列化为 `[]Relationship`；失败则尝试宽松结构 `{nature, cweId, viewId}` 再用 [`ParseRelationshipNature`](./enum-relationship-nature) 转换。
4. 两种解析都失败返回 [`ParseError`](./parse-error)。

::: tip 何时传 viewID
CWE 关系图是跨视图的全集。若只关心某个视图（如研究型视图 1000）下的父子，传 `viewID` 可显著减少返回量并保证语义一致。
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

    // 不限定视图
    parents, err := client.GetParents(context.Background(), 79)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("CWE-79 的父级：")
    for _, rel := range parents {
        fmt.Printf("  -> CWE-%d (%s)\n", rel.CWEID, rel.Nature)
    }

    // 限定研究型视图
    children, err := client.GetChildren(context.Background(), 79, 1000)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("视图1000下CWE-79的子级: %d 条\n", len(children))
}
```

::: warning 只返回直接关系
`GetParents`/`GetChildren` 只返回**一层**邻居。需要全链路祖先或后代请用 [`GetAncestors`/`GetDescendants`](./api-ancestors-descendants)。
:::

## 📚 相关链接

- [祖先后代关系](./api-ancestors-descendants) | [RelationshipNature 枚举](./enum-relationship-nature) | [FetchWithRelations](./basic-fetcher)
