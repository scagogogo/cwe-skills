---
title: HTTPClient HTTP 客户端
outline: [2, 3]
---

# 🚦 HTTPClient — 传输层 HTTP 客户端

`HTTPClient` 是 `cweskills` 包的**传输层**，封装标准库 `http.Client`，叠加了自动重试、令牌桶限流、超时控制与 User-Agent 标识。[`APIClient`](./api-client) 内部即持有一个 `HTTPClient` 实例；你也可以独立使用它访问任意遵循 `{ "Data": ... }` 约定的 JSON API。

源文件：`http_client.go`。

## 🔢 常量

```go
const DefaultBaseURL    = "https://cwe-api.mitre.org/api/v1"
const DefaultTimeout    = 30 * time.Second
const DefaultUserAgent  = "cwe-sdk-go/" + Version
```

## 🧱 结构体定义

```go
type HTTPClient struct {
    client      *http.Client
    baseURL     string
    userAgent   string
    maxRetries  int
    retryDelay  time.Duration
    rateLimiter *RateLimiter
}
```

字段全部私有，通过 Getter/Setter 访问。`maxRetries` 默认 `0`（不重试），`retryDelay` 默认 `1s`。

## 🏗️ 构造函数

```go
func NewHTTPClient(baseURL string, opts ...HTTPClientOption) *HTTPClient
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `baseURL` | `string` | API 基础 URL |
| `opts` | `...HTTPClientOption` | 配置选项，详见 [HTTPClientOption](./http-client-option) |

返回值：`*HTTPClient`，已设置默认超时与 `DefaultUserAgent`。

## 📤 请求方法

| 方法 | 签名 | 说明 |
| --- | --- | --- |
| `Get` | `(ctx, path, result interface{}) error` | GET 并自动 JSON 解码到 `result` |
| `GetRaw` | `(ctx, path) ([]byte, error)` | GET 返回原始响应体 |
| `Post` | `(ctx, path, body, result interface{}) error` | POST JSON body 并解码 |
| `PostForm` | `(ctx, path, data url.Values, result interface{}) error` | POST 表单并解码 |

所有方法都经私有 `doRequest(ctx, method, path, body)` 统一处理限流、重试与状态码校验，路径由 `buildURL` 拼接到 `baseURL` 之后。

::: tip Get 的 result 可为 nil
`Get` 的 `result` 传 `nil` 时只发请求不解析，适合只关心状态码的场景。需要原始字节用 `GetRaw`。
:::

## 🔧 Getter / Setter

```go
func (c *HTTPClient) GetBaseURL() string
func (c *HTTPClient) SetBaseURL(url string)
func (c *HTTPClient) GetMaxRetries() int
func (c *HTTPClient) SetMaxRetries(maxRetries int)
func (c *HTTPClient) GetRetryDelay() time.Duration
func (c *HTTPClient) SetRetryDelay(delay time.Duration)
func (c *HTTPClient) GetRateLimiter() *RateLimiter
func (c *HTTPClient) SetRateLimiter(limiter *RateLimiter)
func (c *HTTPClient) GetHTTPClient() *http.Client
func (c *HTTPClient) SetHTTPClient(client *http.Client)
func (c *HTTPClient) Close()
```

`Close` 调用 `client.CloseIdleConnections()` 释放空闲连接，长驻进程应在退出前调用。

## 🚀 可运行示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewHTTPClient(
        cweskills.DefaultBaseURL,
        cweskills.WithRetry(3, 2*time.Second),
        cweskills.WithHTTPRateLimiter(1.0, 5),
        cweskills.WithHTTPTimeout(30*time.Second),
    )
    defer client.Close()

    var resp struct {
        Data json.RawMessage `json:"Data"`
    }
    if err := client.Get(context.Background(), "/version", &resp); err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(resp.Data))
}
```

::: warning 需手动导入 json
上面示例中 `json.RawMessage` 来自标准库 `encoding/json`，使用时需自行 `import`。
:::

## 📚 相关链接

- [HTTPClientOption 配置](./http-client-option) | [请求方法](./http-methods) | [重试策略](./http-retry) | [限流器](./rate-limiter)
