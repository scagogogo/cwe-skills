---
title: GetVersion 获取版本信息
outline: [2, 3]
---

# 📥 GetVersion — 获取 CWE 数据版本信息

`GetVersion` 调用 `GET /version` 端点，返回当前 CWE 数据库的版本号、发布日期与名称。常用于启动自检、缓存失效判断或展示数据时效性。

源文件：`api_client_version.go`。

## 📐 函数签名

```go
func (c *APIClient) GetVersion(ctx context.Context) (*VersionResponse, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |

返回值：

| 返回 | 说明 |
| --- | --- |
| `*VersionResponse` | 版本信息，含 `Version`、`ReleaseDate`、`Name` |
| `error` | 请求或解析失败时非 nil |

## 🔁 内部流程

1. 路径 `/version`，底层 `HTTPClient.Get` 发请求。
2. 从 `APIResponse.Data` 反序列化为 `VersionResponse`。
3. 解析失败返回 [`ParseError`](./parse-error)。

::: tip 版本端点不计入业务配额
`/version` 是轻量端点，适合作为「连通性探针」——启动时调用一次，既验证了网络与鉴权，又拿到了数据版本。
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

    version, err := client.GetVersion(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("版本: %s\n", version.Version)
    fmt.Printf("发布日期: %s\n", version.ReleaseDate)
    fmt.Printf("名称: %s\n", version.Name)
}
```

::: details VersionResponse 结构
```go
type VersionResponse struct {
    Version     string `json:"version"`
    ReleaseDate string `json:"releaseDate"`
    Name        string `json:"name"`
}
```
完整定义见 [API 响应类型](./api-response)。
:::

## 📚 相关链接

- [API 响应类型](./api-response) | [APIClient 概览](./api-client) | [错误处理](./errors)
