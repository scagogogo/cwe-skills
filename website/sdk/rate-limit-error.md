---
title: RateLimitError 速率超限错误
outline: [2, 3]
---

# ⚠️ RateLimitError — 请求速率超限错误

当请求频率超过 API 速率限制时抛出。`RateLimitError` 携带 `RetryAfter`，提示调用方应等待多久后再重试。

源文件：`errors.go`。

## 🧱 结构体定义

```go
type RateLimitError struct {
    *CWEError
    RetryAfter time.Duration
}
```

| 字段 | 说明 |
| --- | --- |
| `CWEError.Code` | `"RATE_LIMIT"` |
| `CWEError.Message` | `"请求速率超限"` |
| `CWEError.Detail` | `建议等待: {retryAfter}` |
| `RetryAfter` | 建议等待时长（`time.Duration`） |

## 🏗️ 构造函数

```go
func NewRateLimitError(retryAfter time.Duration) *RateLimitError
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `retryAfter` | `time.Duration` | 建议等待时间 |

返回 `*RateLimitError`。

::: tip 与客户端限流器的关系
SDK 内置的 [`RateLimiter`](./rate-limiter) 是**客户端主动节流**，通过 `Wait` 阻塞避免超频。`RateLimitError` 是**服务端拒绝**后的错误反馈，二者互补：客户端节流正常工作时通常不会触发服务端 429。若仍收到此错误，说明限流配置过宽，应收紧 `WithAPIRateLimit`。
:::

## 🚀 可运行示例

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewAPIClient()
    defer client.Close()

    _, err := client.GetWeakness(context.Background(), 79)
    if err != nil {
        var rl *cweskills.RateLimitError
        if errors.As(err, &rl) {
            fmt.Printf("被限流，建议等待 %v\n", rl.RetryAfter)
            time.Sleep(rl.RetryAfter)
            // 重试一次
        }
    }
}
```

::: warning SDK 默认不会返回此错误
当前 `HTTPClient` 在传输层用 `RateLimiter.Wait` 阻塞，不会主动构造 `RateLimitError`。它主要为上层应用或自定义 Transport 解析 HTTP 429 后抛出而保留。若你实现了 429 检测，应 `return cweskills.NewRateLimitError(retryAfter)`。
:::

## 📚 相关链接

- [错误体系概览](./errors) | [RateLimiter 限流器](./rate-limiter) | [APIError](./api-error) | [NewAPIClient](./new-api-client)
