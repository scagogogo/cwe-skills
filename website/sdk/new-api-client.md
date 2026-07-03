---
title: NewAPIClient 创建客户端与配置选项
outline: [2, 3]
---

# 🚀 NewAPIClient — 创建客户端与配置选项

`NewAPIClient` 是构造 `APIClient` 的唯一入口，配合一组函数式选项（`APIClientOption`）实现按需配置。所有选项都可选，零参数调用即得到一个开箱即用、遵守 MITRE 默认限流的客户端。

源文件：`api_client.go`。

## 📐 构造函数签名

```go
func NewAPIClient(opts ...APIClientOption) *APIClient

type APIClientOption func(*APIClient)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `opts` | `...APIClientOption` | 零个或多个配置选项 |

返回值：`*APIClient`，已初始化底层 `HTTPClient`，`baseURL` 默认 `DefaultBaseURL`。

::: tip 默认值清单
零参数调用时：`baseURL = "https://cwe-api.mitre.org/api/v1"`、超时 `30s`、限流 `0.1` req/s burst `1`、重试 `0` 次。
:::

## 🔧 配置选项

| 选项 | 签名 | 作用 |
| --- | --- | --- |
| `WithAPIBaseURL` | `(url string)` | 替换基础 URL，指向镜像或代理 |
| `WithAPITimeout` | `(timeout time.Duration)` | 设置请求超时 |
| `WithAPIRateLimit` | `(rate float64, burst int)` | 自定义令牌桶限流 |
| `WithAPIRetry` | `(maxRetries int, delay time.Duration)` | 5xx 自动重试 |
| `WithAPIHTTPClient` | `(opts ...HTTPClientOption)` | 透传到底层 `HTTPClient` 的选项 |

### WithAPIBaseURL

```go
func WithAPIBaseURL(url string) APIClientOption
```

设置 `baseURL`，并同步到底层 `HTTPClient`。适用于接入企业内网镜像。

### WithAPITimeout

```go
func WithAPITimeout(timeout time.Duration) APIClientOption
```

直接写入底层 `http.Client.Timeout`，覆盖 `DefaultTimeout`（30s）。

### WithAPIRateLimit

```go
func WithAPIRateLimit(rate float64, burst int) APIClientOption
```

用 [`NewRateLimiter`](./rate-limiter) 重建底层限流器。`rate` 为每秒令牌数（`0.1` = 每 10 秒 1 个），`burst` 为桶容量。

### WithAPIRetry

```go
func WithAPIRetry(maxRetries int, delay time.Duration) APIClientOption
```

同时设置 `maxRetries` 与 `retryDelay`，对 5xx 响应自动重试。详见 [HTTP 重试](./http-retry)。

### WithAPIHTTPClient

```go
func WithAPIHTTPClient(opts ...HTTPClientOption) APIClientOption
```

把 [`HTTPClientOption`](./http-client-option) 逐个应用到底层 `HTTPClient`，是「穿透」配置通道。

## 🚀 可运行示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    cwe "github.com/scaggogo/cwe-skills"
)

func main() {
    client := cwe.NewAPIClient(
        cwe.WithAPIBaseURL("https://cwe-mirror.internal/api/v1"),
        cwe.WithAPITimeout(60*time.Second),
        cwe.WithAPIRateLimit(1.0, 5),       // 每秒1个，突发5个
        cwe.WithAPIRetry(3, 2*time.Second), // 5xx最多重试3次
        cwe.WithAPIHTTPClient(
            cwe.WithUserAgent("my-app/1.0"),
        ),
    )
    defer client.Close()

    w, err := client.GetWeakness(context.Background(), 89)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(w.Name)
}
```

::: warning WithAPIHTTPClient 与 WithAPI* 的顺序
`WithAPIHTTPClient` 把选项透传到底层 `HTTPClient`。若在它之后调用 `WithAPIRateLimit` / `WithAPIRetry`，后者会覆盖前者对限流器/重试的设置。建议把 `WithAPIHTTPClient` 放在最后。
:::

::: details 零参数等价写法
`cwe.NewAPIClient()` 与下面这段完全等价：

```go
cwe.NewAPIClient(
    cwe.WithAPIBaseURL(cwe.DefaultBaseURL),
    cwe.WithAPITimeout(cwe.DefaultTimeout),
    cwe.WithAPIRateLimit(0.1, 1),
)
```
:::

## 📚 相关链接

- [APIClient 概览](./api-client) | [HTTPClient 选项](./http-client-option) | [限流器](./rate-limiter)
