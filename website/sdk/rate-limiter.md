---
title: RateLimiter 令牌桶限流器
outline: [2, 3]
---

# 🚦 RateLimiter — 令牌桶限流器

`RateLimiter` 是 `cweskills` 包基于**令牌桶算法**实现的速率限制器，仅依赖 Go 标准库。它被 `HTTPClient` 用于控制对 MITRE API 的请求频率，也可独立用于任何需要节流的场景。

源文件：`http_rate_limiter.go`。

## 🧱 结构体定义

```go
type RateLimiter struct {
    mu         sync.Mutex
    rate       float64       // 每秒允许的令牌数
    burst      int           // 桶容量（最大突发数）
    tokens     float64       // 当前令牌数
    lastRefill time.Time     // 上次填充令牌的时间
    interval   time.Duration // 请求间隔（兼容旧接口）
    lastReq    time.Time     // 上次请求时间（兼容旧接口）
}
```

字段全部私有，通过方法访问。`tokens` 初始为满桶（`float64(burst)`）。

## 🏗️ 构造函数

```go
func NewRateLimiter(rate float64, burst int) *RateLimiter
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `rate` | `float64` | 每秒令牌数；`1.0` = 每秒 1 个，`0.1` = 每 10 秒 1 个 |
| `burst` | `int` | 桶容量，允许的瞬时突发请求数 |

返回值：`*RateLimiter`，初始满桶，`lastRefill` 设为当前时间。

## 🔄 核心方法

| 方法 | 签名 | 行为 |
| --- | --- | --- |
| `Allow` | `() bool` | 非阻塞；有令牌则消耗并返回 `true`，否则 `false` |
| `Wait` | `(ctx context.Context) error` | 阻塞等待至有令牌，或 `ctx` 取消 |
| `WaitForRequest` | `()` | 阻塞确保两次请求间隔不小于 `interval`（兼容旧接口） |
| `refill` | （私有） | 按经过时间补充令牌，上限为 `burst` |

::: tip Allow vs Wait
`Allow` 是「试一下就走」，适合可丢弃的请求；`Wait` 是「等到能发为止」，适合必须执行的请求。`HTTPClient.doRequest` 用的是 `Wait`，保证请求一定会发出。
:::

## 🚀 可运行示例

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    limiter := cweskills.NewRateLimiter(1.0, 5) // 每秒1个，突发5个

    // 突发5个立即通过
    for i := 0; i < 5; i++ {
        fmt.Println("Allow:", limiter.Allow())
    }
    // 第6个被拒
    fmt.Println("第6次:", limiter.Allow())

    // 阻塞等待
    if err := limiter.Wait(context.Background()); err != nil {
        fmt.Println("等待被取消:", err)
    }
    fmt.Println("当前令牌:", limiter.Tokens())
}
```

::: warning Wait 的最小等待粒度
`Wait` 计算的等待时间若小于 `1ms` 会被提升到 `1ms`，避免忙等。在极高 `rate` 下实际吞吐会略低于设定值。
:::

## 📚 相关链接

- [限流器 API 速查](./rate-limiter-api) | [HTTPClient](./http-client) | [HTTPClientOption](./http-client-option)
