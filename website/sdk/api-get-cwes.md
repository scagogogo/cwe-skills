---
title: GetCWEs 批量获取弱点
outline: [2, 3]
---

# 📥 GetCWEs — 批量获取多个 CWE 弱点

`GetCWEs` 调用 `GET /cwe/{ids}` 端点，一次性获取多个 CWE 弱点。相比循环调用 [`GetWeakness`](./api-get-weakness)，批量接口大幅减少网络往返与限流等待。

源文件：`api_client_cwe.go`。

## 📐 函数签名

```go
func (c *APIClient) GetCWEs(ctx context.Context, ids []int) (map[string]*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `ctx` | `context.Context` | 请求上下文 |
| `ids` | `[]int` | CWE ID 数字列表 |

返回值：

| 返回 | 说明 |
| --- | --- |
| `map[string]*CWE` | 以 CWE ID 字符串为键的弱点映射 |
| `error` | 请求或解析失败时非 nil |

## 🔁 内部流程

1. `ids` 为空时直接返回空 map。
2. 逐个校验 `id > 0`，发现非法值立即返回 [`InvalidCWEIDError`](./invalid-cwe-id-error)。
3. 内部函数 `joinIDs` 用逗号拼接，例如 `[79, 89]` → `"79,89"`。
4. 路径 `/cwe/{ids}`，底层 `HTTPClient.Get` 发请求。
5. 反序列化为 `map[string]*CWE`，并把每条 `CWEType` 置为 `"weakness"`。

::: tip 单次请求即拿全量
`GetCWEs` 只发起**一次** HTTP 请求，因此只消耗一个限流令牌。在默认限流下，获取 50 条弱点和获取 1 条耗时相同。
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

    ids := []int{79, 89, 119, 20, 200}
    result, err := client.GetCWEs(context.Background(), ids)
    if err != nil {
        log.Fatal(err)
    }
    for id, w := range result {
        if w != nil {
            fmt.Printf("CWE-%s: %s\n", id, w.Name)
        }
    }
}
```

::: warning map 的键是字符串
返回值键类型是 `string`（如 `"79"`），不是 `int`。遍历时注意类型。某些 ID 可能缺失（API 未返回），访问前判断 `w != nil`。
:::

::: details 批量与并存的注册表写入
配合 [`MultipleFetcher.FetchMultipleToRegistry`](./multiple-fetcher)，可把批量结果直接灌入 [`Registry`](./model)，作为离线缓存的构建方式之一。
:::

## 📚 相关链接

- [GetWeakness 单条获取](./api-get-weakness) | [MultipleFetcher](./multiple-fetcher) | [响应类型](./api-response)
