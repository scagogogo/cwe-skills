---
title: HTTPClientOption 配置选项
outline: [2, 3]
---

# 🔧 HTTPClientOption — HTTP 客户端配置选项

`HTTPClientOption` 是 `HTTPClient` 的函数式选项类型，配合 [`NewHTTPClient`](./http-client) 实现按需配置。每个选项都是一个接收 `*HTTPClient` 的闭包。

源文件：`http_client.go`。

## 📐 类型定义

```go
type HTTPClientOption func(*HTTPClient)
```

## 🧰 选项清单

| 选项 | 签名 | 作用 |
| --- | --- | --- |
| `WithRetry` | `(maxRetries int, delay time.Duration)` | 设置 5xx 重试次数与间隔 |
| `WithHTTPRateLimiter` | `(rate float64, burst int)` | 安装令牌桶限流器 |
| `WithHTTPTimeout` | `(timeout time.Duration)` | 设置请求超时 |
| `WithUserAgent` | `(ua string)` | 自定义 User-Agent 头 |
| `WithHTTPClient` | `(client *http.Client)` | 注入自定义底层 `http.Client` |

## WithRetry

```go
func WithRetry(maxRetries int, delay time.Duration) HTTPClientOption
```

设置 `maxRetries`（不含首次请求，`0` 表示不重试）与 `retryDelay`。仅在响应状态码 `>= 500` 且未达上限时重试。详见 [HTTP 重试](./http-retry)。

## WithHTTPRateLimiter

```go
func WithHTTPRateLimiter(rate float64, burst int) HTTPClientOption
```

调用 [`NewRateLimiter(rate, burst)`](./rate-limiter) 创建限流器并赋给 `rateLimiter` 字段。`rate` 为每秒令牌数，`burst` 为桶容量。

## WithHTTPTimeout

```go
func WithHTTPTimeout(timeout time.Duration) HTTPClientOption
```

直接写入 `client.Timeout`，覆盖 `DefaultTimeout`（30s）。等价于 `APIClient` 的 [`WithAPITimeout`](./new-api-client)。

## WithUserAgent

```go
func WithUserAgent(ua string) HTTPClientOption
```

覆盖 `DefaultUserAgent`（`cwe-sdk-go/{Version}`）。每个请求的 `User-Agent` 头都会被设为此值，便于服务端识别调用方。

## WithHTTPClient

```go
func WithHTTPClient(client *http.Client) HTTPClientOption
```

整体替换底层 `*http.Client`，适合注入自定义 Transport（如代理、mTLS、连接池调优）或测试桩。

## 🚀 可运行示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    customHTTP := &http.Client{
        Timeout: 45 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns: 50,
        },
    }

    client := cwe.NewHTTPClient(
        cwe.DefaultBaseURL,
        cwe.WithHTTPClient(customHTTP),
        cwe.WithUserAgent("my-scanner/2.1"),
        cwe.WithRetry(3, time.Second),
        cwe.WithHTTPRateLimiter(0.5, 2),
    )
    defer client.Close()

    raw, err := client.GetRaw(context.Background(), "/version")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(raw))
}
```

::: warning WithHTTPClient 会覆盖 timeout
`WithHTTPClient` 用你传入的 `*http.Client` 整体替换默认实例，之后再调用 `WithHTTPTimeout` 才会生效。注意选项顺序：先 `WithHTTPClient`，再 `WithHTTPTimeout`。
:::

::: tip 与 APIClient 选项的关系
`APIClient` 提供 `WithAPIHTTPClient(opts ...HTTPClientOption)` 作为透传通道，把这里的选项转发给底层 `HTTPClient`，无需直接构造 `HTTPClient`。
:::

## 📚 相关链接

- [HTTPClient 概览](./http-client) | [HTTP 重试](./http-retry) | [限流器](./rate-limiter) | [NewAPIClient](./new-api-client)
