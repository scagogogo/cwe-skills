---
title: HTTPClient 请求方法 Get Post
outline: [2, 3]
---

# 📤 HTTP 请求方法 — Get / GetRaw / Post / PostForm

`HTTPClient` 提供四个公开方法覆盖常见 HTTP 场景。`Get`/`Post` 自动 JSON 解码；`GetRaw` 返回原始字节；`PostForm` 处理表单提交。所有方法都经私有 `doRequest` 统一处理限流、重试与状态码校验。

源文件：`http_client.go`。

## 📐 方法签名

```go
func (c *HTTPClient) Get(ctx context.Context, path string, result interface{}) error
func (c *HTTPClient) GetRaw(ctx context.Context, path string) ([]byte, error)
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error
func (c *HTTPClient) PostForm(ctx context.Context, path string, data url.Values, result interface{}) error
```

## Get

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |
| `path` | `string` | 相对 `baseURL` 的路径 |
| `result` | `interface{}` | 解析目标指针，`nil` 则不解析 |

返回 `error`，解析失败返回 [`ParseError`](./parse-error)，非 2xx 返回 [`APIError`](./api-error)。

## GetRaw

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |
| `path` | `string` | 请求路径 |

返回 `([]byte, error)`。适合响应非 JSON 或需保留原始字节的场景。

## Post

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `body` | `interface{}` | 请求体，会被 `json.Marshal` 序列化 |
| `result` | `interface{}` | 解析目标指针 |

请求头自动设 `Content-Type: application/json`。

## PostForm

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `data` | `url.Values` | 表单键值对 |
| `result` | `interface{}` | 解析目标指针 |

请求头自动设 `Content-Type: application/x-www-form-urlencoded`。`PostForm` 自行实现限流与重试循环，未走 `doRequest`，但行为一致。

## 🔁 doRequest 内部流程

1. `buildURL` 拼接 `baseURL + path`。
2. 若安装了 `rateLimiter`，调用 `Wait(ctx)` 阻塞等待令牌。
3. 循环 `attempt = 0..maxRetries`：每轮创建请求、设置 `User-Agent`（有 body 时设 `Content-Type: application/json`）、发送。
4. 状态码 `>= 500` 且未达重试上限 → 记录 `APIError` 继续；非 2xx → 立即返回 `APIError`；2xx → 返回响应体。
5. 重试前等待 `retryDelay`，期间若 `ctx` 取消则返回 `ctx.Err()`。

## 🚀 可运行示例

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/url"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewHTTPClient(cweskills.DefaultBaseURL)
    defer client.Close()

    // Get + 自动解码
    var resp struct {
        Data json.RawMessage `json:"Data"`
    }
    if err := client.Get(context.Background(), "/version", &resp); err != nil {
        log.Fatal(err)
    }

    // GetRaw 原始字节
    raw, err := client.GetRaw(context.Background(), "/cwe/weakness/79")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("原始响应长度:", len(raw))

    // PostForm 示例（占位端点）
    form := url.Values{}
    form.Set("q", "xss")
    var result map[string]interface{}
    _ = client.PostForm(context.Background(), "/search", form, &result)
}
```

::: tip 路径以 / 开头
`path` 应以 `/` 开头，`buildURL` 直接做字符串拼接 `baseURL + path`。若 `baseURL` 末尾带 `/` 且 `path` 也带 `/`，会产生双斜杠，请保持 `baseURL` 不带尾斜杠。
:::

::: warning PostForm 的重试与 Get 略有差异
`PostForm` 没有复用 `doRequest`，而是内联了限流与重试循环，但语义一致：仅 5xx 重试，非 2xx 立即返回 `APIError`。若你定制了 `rateLimiter`，表单请求同样受其约束。
:::

## 📚 相关链接

- [HTTPClient 概览](./http-client) | [重试策略](./http-retry) | [APIError](./api-error) | [ParseError](./parse-error)
