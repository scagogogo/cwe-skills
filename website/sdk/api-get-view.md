---
title: GetView 获取视图详情
outline: [2, 3]
---

# 📥 GetView — 获取 CWE 视图详情

`GetView` 调用 `GET /cwe/view/{id}` 端点，返回指定 ID 的 CWE 视图（View）。视图从特定视角组织弱点，例如「研究型视图」「开发型视图」。

源文件：`api_client_cwe.go`。

## 📐 函数签名

```go
func (c *APIClient) GetView(ctx context.Context, id int) (*View, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |
| `id` | `int` | 视图 ID 数字，例如 `1000`（研究型视图） |

返回值：

| 返回 | 说明 |
| --- | --- |
| `*View` | 视图详情，含 `Members` 成员列表 |
| `error` | 请求或解析失败时非 nil |

## 🔁 内部流程

1. 校验 `id > 0`，否则返回 [`InvalidCWEIDError`](./invalid-cwe-id-error)。
2. 路径 `/cwe/view/{id}`，底层 `HTTPClient.Get` 发请求。
3. 先尝试 `[]View`，再退化单条 `View`。
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

    view, err := client.GetView(context.Background(), 1000)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("View-%d: %s\n", view.ID, view.Name)
    fmt.Printf("成员数量: %d\n", len(view.Members))
}
```

::: tip 视图用于限定关系范围
[`GetParents`](./api-parents-children) 与 `GetChildren` 接受可选 `viewID` 参数，传入视图 ID 后只返回该视图内的关系，避免拉到全量图谱。
:::

## 📚 相关链接

- [GetWeakness](./api-get-weakness) | [GetCategory](./api-get-category) | [View 结构体](./view) | [父子关系](./api-parents-children)
