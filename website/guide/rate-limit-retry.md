---
title: 速率限制与重试
outline: [2, 3]
---

# ⚙️ 速率限制与重试

CWE Skills 的在线路径调用 MITRE REST API（`https://cwe-api.mitre.org/api/v1`）。MITRE 对该 API 有速率限制，频繁请求会被返回 `429 Too Many Requests`。为此，SDK 内置**令牌桶速率限制**与**指数退避自动重试**，让调用方无需手写限流与重试逻辑。

::: tip 离线模式不受影响
本页讨论的速率限制与重试只针对**在线 API 路径**。离线 XML 解析、注册表、导航、树、搜索都不触网，无任何速率限制。大批量场景请走离线。见 [在线 vs 离线](./online-offline)。
:::

---

## 🪣 令牌桶速率限制

### 原理

`RateLimiter` 基于**令牌桶（token bucket）**算法：

```text
        ┌─────────────────────────────┐
        │   令牌桶（容量 = burst）     │
        │  ███████░░░  当前令牌数      │
        └────────────┬────────────────┘
                     │ 每秒补充 rate 个令牌
                     ▼
   请求来了 ──► 取 1 个令牌 ──► 有令牌？放行 : 等待
```

- **rate**：令牌补充速率（每秒补充多少令牌），即稳态请求速率上限。
- **burst**：桶容量，即允许的瞬时突发请求数。
- 每个请求消耗 1 个令牌；桶空时请求需等待令牌补充。

### 默认配置

```go
// 默认 rate=0.1, burst=1，约每 10 秒 1 个请求
client := cweskills.NewAPIClient()
```

::: warning 默认很保守
MITRE API 的速率限制较严，默认 `rate=0.1, burst=1` 是保守值，避免触发 429。如果你的场景需要更快，可调高，但要自负风险——频繁 429 可能导致 IP 被临时封禁。
:::

### API

```go
type RateLimiter struct { /* ... */ }

func NewRateLimiter(rate float64, burst int) *RateLimiter

func (r *RateLimiter) Allow() bool              // 是否有令牌（不等待）
func (r *RateLimiter) Wait(ctx) error           // 阻塞等待直到有令牌
func (r *RateLimiter) WaitForRequest(ctx) error // 等待一个请求的令牌
```

```go
rl := cweskills.NewRateLimiter(0.5, 2) // 每秒 0.5 个，突发 2
if rl.Allow() {
    // 立即放行
}
rl.WaitForRequest(ctx) // 阻塞直到拿到令牌
```

---

## 🔁 自动重试（指数退避）

### 原理

当请求返回**服务端错误（5xx）**或**速率限制（429）**时，`HTTPClient` 会自动重试，采用**指数退避（exponential backoff）**：

```text
请求失败
  │
  ├─ 第 1 次重试：等待 delay
  ├─ 第 2 次重试：等待 delay × 2
  ├─ 第 3 次重试：等待 delay × 4
  └─ ... 直到 maxRetries 次，返回最后一个错误
```

```go
func WithRetry(maxRetries int, delay time.Duration) Option
```

```go
client := cweskills.NewAPIClient(
    cweskills.WithAPIRetry(3, 2*time.Second), // 最多重试 3 次，初始退避 2s
)
```

::: info 重试触发条件
- **5xx 服务端错误**：MITRE 服务器临时故障，重试通常能恢复。
- **429 速率限制**：触发限流，重试前会读取 `Retry-After` 头（若存在）。
- **4xx 客户端错误（如 404）**：不重试，直接返回（重试也没用）。
:::

---

## ⚙️ 配置选项

### WithAPIRateLimit

调整令牌桶参数：

```go
client := cweskills.NewAPIClient(
    cweskills.WithAPIRateLimit(0.5, 1), // 每秒 0.5 个请求，突发 1
)
```

### WithAPIRetry

调整重试策略：

```go
client := cweskills.NewAPIClient(
    cweskills.WithAPIRetry(5, time.Second), // 最多重试 5 次，初始退避 1s
)
```

### WithAPITimeout

HTTP 请求超时（默认 30s）：

```go
client := cweskills.NewAPIClient(
    cweskills.WithAPITimeout(60*time.Second),
)
```

### 组合配置

```go
client := cweskills.NewAPIClient(
    cweskills.WithAPIBaseURL("https://cwe-api.mitre.org/api/v1"),
    cweskills.WithAPITimeout(60*time.Second),
    cweskills.WithAPIRateLimit(0.5, 1),
    cweskills.WithAPIRetry(3, 2*time.Second),
)
defer client.Close()
```

---

## 🛠️ 自定义 HTTPClient

如果内置的 `HTTPClient` 不能满足需求（如需要代理、自定义 Transport、OAuth、自定义日志），可以用 `WithHTTPClient` 传入自定义的 `*http.Client`：

```go
import "net/http"

customHTTP := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        Proxy: http.ProxyFromEnvironment, // 走系统代理
        // 其他 Transport 配置...
    },
}

client := cweskills.NewAPIClient(
    cweskills.WithHTTPClient(customHTTP),
)
```

::: warning 自定义 Client 的注意事项
- 默认超时 30s，自定义时记得设置 `Timeout`，避免请求挂死。
- 默认 `UserAgent` 为 `cwe-sdk-go/v0.0.1`，自定义 Client 时若覆盖了 Header，请保留合理的 UserAgent。
- 速率限制与重试仍由 SDK 的 `HTTPClient` 层处理（它包装你传入的 `*http.Client`），所以限流/重试逻辑不会丢失。
:::

::: details UserAgent 默认值
SDK 默认在请求头设置 `UserAgent: cwe-sdk-go/v0.0.1`，便于 MITRE 识别流量来源。自定义 `*http.Client` 时，SDK 仍会设置该头（除非你在 Transport 层覆盖）。
:::

---

## 🚦 触发 429 怎么办

即便有令牌桶，极端场景下仍可能触发 429（如多进程共享 IP 并发请求）。SDK 会：

1. 读取响应的 `Retry-After` 头（秒数）。
2. 封装为 `*RateLimitError`（带 `RetryAfter` 字段）。
3. 若配置了重试，按 `max(Retry-After, 退避)` 等待后重试。

```go
weakness, err := client.GetWeakness(ctx, 79)
if err != nil {
    var rlErr *cweskills.RateLimitError
    if errors.As(err, &rlErr) {
        fmt.Println("被限流，建议等待:", rlErr.RetryAfter)
        time.Sleep(rlErr.RetryAfter)
        // 重试...
    }
}
```

::: danger 多进程共享 IP 慎用高 rate
令牌桶是**进程内**的。多个进程共享同一出口 IP 时，各自独立计数，实际请求速率是各进程之和，极易触发 429。多进程场景请降低每个进程的 rate，或改用离线 XML。
:::

---

## 📈 速率与重试调优建议

| 场景 | rate | burst | maxRetries | delay |
|------|------|-------|------------|-------|
| 偶尔查一两条（默认） | 0.1 | 1 | 3 | 2s |
| 中频查询 | 0.5 | 1 | 3 | 1s |
| 大批量（不建议在线） | — | — | — | 改用离线 XML |

::: tip 大批量永远是离线
无论怎么调速率，在线 API 都不适合批量查几百个 CWE。令牌桶会把你卡在每秒零点几个请求。大批量请下载 XML 走离线，无任何限制。见 [在线 vs 离线](./online-offline)。
:::

---

## 📖 相关文档

- [在线 vs 离线模式](./online-offline)
- [错误处理](./error-handling)（`RateLimitError` 详解）
- [Go SDK 接入](./integration-sdk)
- [工作原理](./how-it-works)
- [CLI 接入](./integration-cli)
