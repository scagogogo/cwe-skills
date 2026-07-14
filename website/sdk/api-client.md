---
title: APIClient MITRE API 客户端概览
outline: [2, 3]
---

# 🌐 APIClient — MITRE CWE REST API 客户端

`APIClient` 是 `cweskills` 包访问 MITRE CWE REST API 的**统一入口**。它在底层封装了 [`HTTPClient`](./http-client) 与 [`RateLimiter`](./rate-limiter)，对外暴露面向领域语义的方法：获取弱点、类别、视图，查询父子/祖先后代关系，读取数据版本。所有 API 调用都自动遵守 MITRE 默认限流（每 10 秒 1 个请求）并支持重试。

源文件：`api_client.go`、`api_client_cwe.go`、`api_client_relations.go`、`api_client_version.go`。

## 🧩 结构体定义

```go
type APIClient struct {
    httpClient *HTTPClient
    baseURL    string
}
```

`APIClient` 仅持有两个私有字段：底层 `HTTPClient`（负责实际网络请求、重试、限流）和 `baseURL`（默认 `https://cwe-api.mitre.org/api/v1`）。两个值始终保持同步——`SetBaseURL` 会同时更新二者。

## 🗺️ 能力地图

| 能力分组 | 方法 | 文档 |
| --- | --- | --- |
| 构造与配置 | `NewAPIClient`、`WithAPI*` 选项 | [创建客户端](./new-api-client) |
| 弱点获取 | `GetWeakness`、`GetCWEs` | [GetWeakness](./api-get-weakness)、[GetCWEs](./api-get-cwes) |
| 类别 / 视图 | `GetCategory`、`GetView` | [GetCategory](./api-get-category)、[GetView](./api-get-view) |
| 版本信息 | `GetVersion` | [GetVersion](./api-get-version) |
| 父子关系 | `GetParents`、`GetChildren` | [Parents/Children](./api-parents-children) |
| 祖先后代 | `GetAncestors`、`GetDescendants` | [Ancestors/Descendants](./api-ancestors-descendants) |
| 响应类型 | `APIResponse` 等 | [API 响应类型](./api-response) |

## 🔧 Getter / Setter

```go
func (c *APIClient) GetBaseURL() string
func (c *APIClient) SetBaseURL(url string)
func (c *APIClient) GetHTTPClient() *HTTPClient
func (c *APIClient) SetHTTPClient(client *HTTPClient)
func (c *APIClient) GetRateLimiter() *RateLimiter
func (c *APIClient) SetRateLimiter(limiter *RateLimiter)
func (c *APIClient) Close()
```

`GetHTTPClient` / `SetHTTPClient` 暴露底层传输层，便于复用连接池或注入测试桩。`Close` 释放底层 `http.Client` 的空闲连接，应在客户端生命周期结束时调用。

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

    version, err := client.GetVersion(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("CWE 版本: %s (%s)\n", version.Version, version.ReleaseDate)

    weakness, err := client.GetWeakness(context.Background(), 79)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("CWE-79: %s\n", weakness.Name)
}
```

::: tip 默认限流为何这么慢
MITRE 官方 API 建议每 10 秒最多 1 个请求，`NewAPIClient` 默认配置 `WithHTTPRateLimiter(0.1, 1)`。若你接入的是自建镜像或代理，可用 [`WithAPIRateLimit`](./new-api-client) 放宽限制。
:::

::: warning 不要遗忘 Close
`APIClient` 内部持有一个 `http.Client`，长驻进程中若反复 `NewAPIClient` 而不 `Close`，会泄漏底层 TCP 连接。推荐用 `defer client.Close()`。
:::

## 📚 相关链接

- [HTTPClient 传输层](./http-client) | [RateLimiter 限流器](./rate-limiter) | [错误处理](./errors)
