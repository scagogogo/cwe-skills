---
title: RelationshipError 关系错误
outline: [2, 3]
---

# ⚠️ RelationshipError — 关系操作失败错误

当尝试建立无效的 CWE 关系时抛出。`RelationshipError` 携带 `From`、`To`、`Nature`，描述被拒绝的关系三元组。

源文件：`errors.go`。

## 🧱 结构体定义

```go
type RelationshipError struct {
    *CWEError
    From   string
    To     string
    Nature RelationshipNature
}
```

| 字段 | 说明 |
| --- | --- |
| `CWEError.Code` | `"RELATIONSHIP_ERROR"` |
| `CWEError.Message` | `"关系操作失败"` |
| `CWEError.Detail` | `无法建立 {from} -> {to} (类型: {nature}) 的关系` |
| `From` | 源 CWE ID（`string`） |
| `To` | 目标 CWE ID（`string`） |
| `Nature` | 关系类型，见 [`RelationshipNature`](./enum-relationship-nature) |

## 🏗️ 构造函数

```go
func NewRelationshipError(from, to string, nature RelationshipNature) *RelationshipError
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `from` | `string` | 源 CWE ID |
| `to` | `string` | 目标 CWE ID |
| `nature` | `RelationshipNature` | 关系类型 |

返回 `*RelationshipError`。

::: tip From/To 是字符串
注意 `From`/`To` 为 `string` 而非 `int`，便于承载 `"CWE-79"` 这样的完整标识。构造时按需格式化。
:::

## 🎯 关系类型

`Nature` 字段类型是 [`RelationshipNature`](./enum-relationship-nature)，常用值包括：

| 常量 | 语义 |
| --- | --- |
| `RelationshipChildOf` | 子级关系（A 是 B 的子） |
| `RelationshipParentOf` | 父级关系（A 是 B 的父） |
| `RelationshipCanAlsoBe` | 也可以是 |
| `RelationshipCanPrecede` | 可以先于 |
| `RelationshipRequires` | 依赖于 |

::: warning 当前 API 客户端不主动抛出
`GetParents`/`GetChildren`/`GetAncestors`/`GetDescendants` 不构造 `RelationshipError`——它们直接返回服务端关系。`RelationshipError` 主要供上层应用在**本地构建关系图**时校验用，例如尝试在 `Registry` 上建立自环或重复边时抛出。
:::

## 🚀 可运行示例

```go
package main

import (
    "errors"
    "fmt"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    err := cweskills.NewRelationshipError("CWE-79", "CWE-79", cweskills.RelationshipChildOf)
    var relErr *cweskills.RelationshipError
    if errors.As(err, &relErr) {
        fmt.Printf("拒绝 %s -> %s (%s)\n", relErr.From, relErr.To, relErr.Nature)
    }
    fmt.Println(err.Error())
}
```

::: details 典型校验场景
本地构建 CWE 关系图时，可这样使用：

```go
func addEdge(graph *Graph, from, to string, nature cweskills.RelationshipNature) error {
    if from == to {
        return cweskills.NewRelationshipError(from, to, nature) // 禁止自环
    }
    if graph.hasEdge(from, to, nature) {
        return cweskills.NewRelationshipError(from, to, nature) // 禁止重复
    }
    graph.addEdge(from, to, nature)
    return nil
}
```
:::

## 📚 相关链接

- [错误体系概览](./errors) | [RelationshipNature 枚举](./enum-relationship-nature) | [父子关系](./api-parents-children) | [CWEError 根类型](./cwe-error)
