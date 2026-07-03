---
title: APIResponse 响应类型
outline: [2, 3]
---

# 📦 API 响应类型 — api_response.go

`api_response.go` 定义了 MITRE CWE REST API 的响应结构体。MITRE API 的统一外层是 `{ "Data": ... }`，`Data` 字段以 `json.RawMessage` 形式承载，由各业务方法按需二次反序列化。

## 🧱 结构体一览

| 结构体 | 用途 | 关键字段 |
| --- | --- | --- |
| `APIResponse` | 通用响应外壳 | `Data`、`Message` |
| `VersionResponse` | 版本端点 | `Version`、`ReleaseDate`、`Name` |
| `WeaknessesResponse` | 弱点列表 | `Data`、`Weaknesses` |
| `CategoriesResponse` | 类别列表 | `Data`、`Categories` |
| `ViewsResponse` | 视图列表 | `Data`、`Views` |
| `RelationsResponse` | 关系列表 | `Data` |
| `CWEsResponse` | 批量弱点 | `Data`、`Weaknesses map[string]*CWE` |

## 📐 APIResponse 通用外壳

```go
type APIResponse struct {
    Data    json.RawMessage `json:"Data"`
    Message string          `json:"Message,omitempty"`
}
```

`Data` 是原始 JSON 字节，调用方负责按业务形态二次解析。`Message` 为可选提示信息。

::: tip 为何用 RawMessage
MITRE 不同端点返回的 `Data` 形态各异（数组/对象/map），用 `json.RawMessage` 延迟解析，让各方法（如 `GetWeakness`）自行尝试多种结构，提高兼容性。
:::

## 📐 VersionResponse

```go
type VersionResponse struct {
    Version     string `json:"version"`
    ReleaseDate string `json:"releaseDate"`
    Name        string `json:"name"`
}
```

由 [`GetVersion`](./api-get-version) 返回。

## 📐 列表型响应

```go
type WeaknessesResponse struct {
    Data       json.RawMessage `json:"Data"`
    Weaknesses []CWE           `json:"weaknesses,omitempty"`
}
type CategoriesResponse struct {
    Data       json.RawMessage `json:"Data,omitempty"`
    Categories []Category      `json:"categories,omitempty"`
}
type ViewsResponse struct {
    Data  json.RawMessage `json:"Data,omitempty"`
    Views []View          `json:"views,omitempty"`
}
```

## 📐 关系与批量响应

```go
type RelationsResponse struct {
    Data json.RawMessage `json:"Data,omitempty"`
}
type CWEsResponse struct {
    Data       json.RawMessage      `json:"Data"`
    Weaknesses map[string]*CWE      `json:"weaknesses,omitempty"`
}
```

`CWEsResponse.Weaknesses` 以 CWE ID 字符串为键，与 [`GetCWEs`](./api-get-cwes) 的返回值形态一致。

## 🚀 直接消费 APIResponse

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cwe.NewAPIClient()
    defer client.Close()

    var resp struct {
        Data json.RawMessage `json:"Data"`
    }
    if err := client.GetHTTPClient().Get(context.Background(), "/version", &resp); err != nil {
        log.Fatal(err)
    }
    var v cwe.VersionResponse
    _ = json.Unmarshal(resp.Data, &v)
    fmt.Println(v.Version)
}
```

::: warning 业务方法已封装解析
通常无需直接操作 `APIResponse`——`GetWeakness`、`GetCWEs` 等方法内部已完成 `Data` 的二次反序列化并返回业务结构体。直接消费 `APIResponse` 仅在需要访问未封装字段时使用。
:::

## 📚 相关链接

- [APIClient 概览](./api-client) | [GetVersion](./api-get-version) | [GetCWEs](./api-get-cwes) | [HTTPClient](./http-client)
