---
title: RateLimiter API 速查
outline: [2, 3]
---

# 📋 RateLimiter API 速查

`RateLimiter` 的查询与配置方法集合。所有公开方法都线程安全（内部用 `sync.Mutex` 保护）。本页是 [RateLimiter 概览](./rate-limiter) 的方法索引。

源文件：`http_rate_limiter.go`。

## 🔧 Getter / Setter

### GetRate / GetBurst

```go
func (r *RateLimiter) GetRate() float64
func (r *RateLimiter) GetBurst() int
```

返回构造时设定的 `rate`（每秒令牌数）与 `burst`（桶容量）。只读，不会触发 `refill`。

### GetInterval / SetInterval

```go
func (r *RateLimiter) GetInterval() time.Duration
func (r *RateLimiter) SetInterval(interval time.Duration)
```

`interval` 是兼容旧接口的「最小请求间隔」。`SetInterval` 在 `interval > 0` 时会同步反算 `rate = float64(time.Second) / float64(interval)`，保持两套语义一致。

::: warning SetInterval 会改写 rate
`SetInterval` 同时修改 `interval` 与 `rate`，调用后再 `GetRate` 会得到反算后的新值。若你同时依赖两个字段，注意顺序。
:::

### Tokens

```go
func (r *RateLimiter) Tokens() float64
```

返回当前可用令牌数。**会触发 `refill`**，因此每次调用都会按经过时间补充令牌（但不消耗）。适合监控面板。

### ResetLastRequest

```go
func (r *RateLimiter) ResetLastRequest()
```

重置 `lastReq` 为零值、`lastRefill` 为当前时间、`tokens` 重新置为满桶。调用后下一次 `WaitForRequest` 立即返回，`Wait` 也能立即取到一个令牌。适合长时闲置后「重新激活」。

## 🔄 行为方法

| 方法 | 签名 | 阻塞 | 说明 |
| --- | --- | --- | --- |
| `Allow` | `() bool` | 否 | 有令牌消耗返回 `true`，否则 `false` |
| `Wait` | `(ctx context.Context) error` | 是 | 等到有令牌，或 `ctx` 取消返回 `ctx.Err()` |
| `WaitForRequest` | `()` | 是 | 按 `interval` 模式保证最小间隔，首次立即返回 |

## 🚀 配合 HTTPClient 使用

```go
package main

import (
    "context"
    "fmt"
    "time"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    limiter := cwe.NewRateLimiter(0.5, 3)
    fmt.Printf("rate=%.2f burst=%d\n", limiter.GetRate(), limiter.GetBurst())

    limiter.SetInterval(2 * time.Second)
    fmt.Printf("调整后 rate=%.2f\n", limiter.GetRate())

    client := cwe.NewHTTPClient(
        cwe.DefaultBaseURL,
        cwe.WithHTTPRateLimiter(1.0, 5),
    )
    defer client.Close()

    // 替换限流器
    client.SetRateLimiter(limiter)
    fmt.Println("当前令牌:", client.GetRateLimiter().Tokens())

    // 用完一轮后重置
    limiter.ResetLastRequest()
    _ = client.GetRaw(context.Background(), "/version")
}
```

::: tip 替换 HTTPClient 的限流器
`HTTPClient.SetRateLimiter` 可在运行时整体替换限流器实例。`APIClient` 也暴露 `GetRateLimiter`/`SetRateLimiter` 转发到底层。
:::

## 📚 相关链接

- [RateLimiter 概览](./rate-limiter) | [HTTPClient](./http-client) | [HTTPClientOption](./http-client-option)
