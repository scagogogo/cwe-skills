---
title: HTTP 重试策略 WithRetry
outline: [2, 3]
---

# 🔁 HTTP 重试策略 — WithRetry 与 doRequest

`HTTPClient` 内置对 `5xx` 服务端错误的自动重试。重试由 `maxRetries`、`retryDelay` 两个字段控制，可通过 [`WithRetry`](./http-client-option) 选项或 Setter 配置。客户端类请求错误（4xx、网络层错误）的语义各有不同，需注意区分。

源文件：`http_client.go`。

## 📐 配置入口

```go
func WithRetry(maxRetries int, delay time.Duration) HTTPClientOption

func (c *HTTPClient) GetMaxRetries() int
func (c *HTTPClient) SetMaxRetries(maxRetries int)
func (c *HTTPClient) GetRetryDelay() time.Duration
func (c *HTTPClient) SetRetryDelay(delay time.Duration)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `maxRetries` | `int` | 最大重试次数，**不含**首次请求；`0` 表示不重试 |
| `delay` | `time.Duration` | 每次重试前的等待时间 |

默认值：`maxRetries = 0`、`retryDelay = 1 * time.Second`。

## 🔁 重试触发条件

在 `doRequest` 的 `for attempt := 0; attempt <= c.maxRetries; attempt++` 循环中：

| 情况 | 行为 |
| --- | --- |
| 状态码 `>= 500` 且 `attempt < maxRetries` | 记录 `APIError` 到 `lastErr`，进入下一轮 |
| 状态码 `>= 500` 且已达上限 | 返回 `lastErr`（`APIError`） |
| 状态码 `< 200` 或 `>= 300`（非 5xx） | **立即**返回 `APIError`，不重试 |
| 网络层错误（`client.Do` 失败） | 记录到 `lastErr`，继续重试 |
| 2xx | 返回响应体，结束 |

每次重试前先 `select` 等待 `retryDelay` 或 `ctx.Done()`；若上下文取消则返回 `ctx.Err()`。

::: tip 仅重试「值得重试」的错误
4xx（如 404、401）是客户端问题，重试无意义，因此立即返回。只有 5xx（服务端临时故障）才进入重试循环——这是保守且合理的默认策略。
:::

::: warning 网络错误也会消耗重试名额
DNS 解析失败、连接超时等网络层错误同样计入 `maxRetries`。若你的网络环境不稳定，建议把 `maxRetries` 调到 `3` 以上。
:::

## 🚀 可运行示例

```go
package main

import (
    "context"
    "log"
    "time"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cwe.NewHTTPClient(
        cwe.DefaultBaseURL,
        cwe.WithRetry(3, 2*time.Second),
    )
    defer client.Close()

    // 运行时动态调整
    client.SetMaxRetries(5)
    client.SetRetryDelay(500 * time.Millisecond)

    var resp struct {
        Data json.RawMessage `json:"Data"`
    }
    if err := client.Get(context.Background(), "/version", &resp); err != nil {
        log.Fatal(err)
    }
    log.Println("成功")
}
```

::: details 与 APIClient 的关系
[`APIClient`](./api-client) 默认不开启重试（`maxRetries=0`），需用 [`WithAPIRetry(maxRetries, delay)`](./new-api-client) 显式启用，本质是写入底层 `HTTPClient` 的同名字段。
:::

## 📚 相关链接

- [HTTPClient 概览](./http-client) | [HTTPClientOption](./http-client-option) | [APIError](./api-error) | [限流器](./rate-limiter)
