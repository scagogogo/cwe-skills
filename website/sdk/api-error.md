---
title: APIError HTTP 调用错误
outline: [2, 3]
---

# ⚠️ APIError — CWE API 调用失败错误

当对 MITRE CWE REST API 的请求返回非 2xx 状态码时抛出。`APIError` 携带 `StatusCode`、`URL`、`Method`，便于按状态码分类处理（5xx 重试、4xx 告警）。

源文件：`errors.go`。由 `HTTPClient.doRequest` 与 `PostForm` 在状态码非 2xx 时构造。

## 🧱 结构体定义

```go
type APIError struct {
    *CWEError
    StatusCode int
    URL        string
    Method     string
}
```

| 字段 | 说明 |
| --- | --- |
| `CWEError.Code` | `"API_ERROR"` |
| `CWEError.Message` | `"CWE API调用失败"` |
| `CWEError.Detail` | `HTTP {statusCode} {method} {url}` |
| `StatusCode` | HTTP 状态码 |
| `URL` | 请求的完整 URL |
| `Method` | HTTP 方法（GET/POST 等） |

## 🏗️ 构造函数

```go
func NewAPIError(statusCode int, url, method string) *APIError
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `statusCode` | `int` | HTTP 状态码 |
| `url` | `string` | 请求 URL |
| `method` | `string` | HTTP 方法 |

返回 `*APIError`。

## 🔍 触发逻辑

在 `doRequest` 的重试循环中：

- 状态码 `>= 500` 且未达重试上限 → 记录 `APIError` 到 `lastErr`，继续重试。
- 状态码 `< 200` 或 `>= 300`（非 5xx）→ **立即**返回 `APIError`。
- 重试耗尽仍 5xx → 返回最后记录的 `APIError`。

::: tip 按状态码分流
5xx 是服务端临时故障，配合 [`WithRetry`](./http-retry) 自动重试；4xx（401/403/404/429）是客户端问题，应立即处理：401 检查鉴权、429 看 [`RateLimitError`](./rate-limit-error)、404 通常会被业务层的 `CWENotFoundError` 提前拦截。
:::

## 🚀 可运行示例

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewAPIClient(cweskills.WithAPIBaseURL("https://invalid.example/api/v1"))
    defer client.Close()

    _, err := client.GetVersion(context.Background())
    if err != nil {
        var apiErr *cweskills.APIError
        if errors.As(err, &apiErr) {
            fmt.Printf("HTTP %d %s %s\n", apiErr.StatusCode, apiErr.Method, apiErr.URL)
            if apiErr.StatusCode >= 500 {
                fmt.Println("服务端错误，可重试")
            }
        }
    }
}
```

::: warning URL 含完整地址
`URL` 是 `buildURL` 拼接后的完整地址（含 `baseURL`），日志里可能暴露内网镜像地址，发布前注意脱敏。
:::

## 📚 相关链接

- [错误体系概览](./errors) | [HTTP 重试](./http-retry) | [RateLimitError](./rate-limit-error) | [HTTPClient](./http-client)
