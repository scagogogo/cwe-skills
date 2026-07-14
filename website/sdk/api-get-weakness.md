---
title: GetWeakness 获取弱点详情
outline: [2, 3]
---

# 📥 GetWeakness — 获取 CWE 弱点详情

`GetWeakness` 调用 MITRE CWE REST API 的 `GET /cwe/weakness/{id}` 端点，返回指定 ID 的完整弱点信息。这是最常用的单条查询入口。

源文件：`api_client_cwe.go`。

## 📐 函数签名

```go
func (c *APIClient) GetWeakness(ctx context.Context, id int) (*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文，控制超时与取消 |
| `id` | `int` | CWE ID 数字，例如 `79` |

返回值：

| 返回 | 说明 |
| --- | --- |
| `*CWE` | 弱点详情，`CWEType` 字段会被置为 `"weakness"` |
| `error` | 请求或解析失败时非 nil |

## 🔁 内部流程

1. 校验 `id > 0`，否则返回 [`InvalidCWEIDError`](./invalid-cwe-id-error)。
2. 拼接路径 `/cwe/weakness/{id}`，经底层 `HTTPClient.Get` 发起请求。
3. 从 `APIResponse.Data` 中先尝试反序列化为 `[]CWE`；若失败再退化为单条 `CWE`。
4. 列表为空时返回 [`CWENotFoundError`](./cwe-not-found-error)。

::: tip 兼容两种响应形态
MITRE API 偶尔以数组、偶尔以对象返回 `Data`。`GetWeakness` 做了双重尝试，无论哪种形态都能正确解析。
:::

## 🚀 可运行示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewAPIClient()
    defer client.Close()

    weakness, err := client.GetWeakness(context.Background(), 79)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("CWE-%d: %s\n", weakness.ID, weakness.Name)
    // 输出: CWE-79: Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')
    fmt.Println("详情页:", weakness.URL)
}
```

::: warning 限流会阻塞
底层默认 `0.1` req/s，连续调用 `GetWeakness` 时第二次会阻塞约 10 秒。需要更高吞吐请用 [`WithAPIRateLimit`](./new-api-client) 调整。
:::

## 📚 相关链接

- [GetCWEs 批量获取](./api-get-cwes) | [GetCategory](./api-get-category) | [错误处理](./errors)
