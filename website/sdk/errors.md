---
title: 错误处理体系概览
outline: [2, 3]
---

# ⚠️ 错误处理体系 — errors.go

`cweskills` 包定义了一套**统一的错误类型体系**，以 `CWEError` 为根，派生出六种细分错误。所有错误都实现 `error` 接口、支持 `errors.Is`/`errors.As` 链式查找，便于上层按错误码分类处理（重试、告警、降级）。

源文件：`errors.go`。

## 🗺️ 错误类型地图

| 类型 | 错误码 | 触发场景 | 文档 |
| --- | --- | --- | --- |
| `CWEError` | （根） | 通用基础类型 | [CWEError](./cwe-error) |
| `InvalidCWEIDError` | `INVALID_CWE_ID` | CWE ID 格式非法 | [InvalidCWEIDError](./invalid-cwe-id-error) |
| `CWENotFoundError` | `CWE_NOT_FOUND` | 条目不存在 | [CWENotFoundError](./cwe-not-found-error) |
| `APIError` | `API_ERROR` | HTTP 非 2xx | [APIError](./api-error) |
| `RateLimitError` | `RATE_LIMIT` | 速率超限 | [RateLimitError](./rate-limit-error) |
| `ValidationError` | `VALIDATION_ERROR` | 字段约束失败 | [ValidationError](./validation-error) |
| `ParseError` | `PARSE_ERROR` | XML/JSON 解析失败 | [ParseError](./parse-error) |
| `RelationshipError` | `RELATIONSHIP_ERROR` | 关系操作非法 | [RelationshipError](./relationship-error) |

## 🧱 继承结构

所有细分错误都**内嵌** `*CWEError`：

```go
type InvalidCWEIDError struct {
    *CWEError
    ID string
}
```

因此每个细分错误既拥有 `CWEError` 的 `Code`/`Message`/`Detail`/`Err`，又携带自己的业务字段（`ID`、`StatusCode`、`RetryAfter` 等）。`Unwrap()` 返回内嵌的 `*CWEError`，可与 `errors.Is`/`errors.As` 配合。

## 🎯 错误处理范式

```go
var err error = client.GetWeakness(ctx, 79)

var notFound *cweskills.CWENotFoundError
if errors.As(err, &notFound) {
    // 404，降级处理
}

var apiErr *cweskills.APIError
if errors.As(err, &apiErr) && apiErr.StatusCode >= 500 {
    // 5xx，可重试
}

var invalid *cweskills.InvalidCWEIDError
if errors.As(err, &invalid) {
    // 入参问题，不应重试
}
```

::: tip 用 errors.As 而非类型断言
`errors.As` 会沿 `Unwrap` 链查找，对内嵌 `*CWEError` 的细分类型也能正确匹配。直接 `err.(*cweskills.APIError)` 可能因包装层而失败。
:::

::: warning Code 是字符串常量
错误码如 `"API_ERROR"` 是字符串，可用 `err.Code == "API_ERROR"` 判断，但更推荐 `errors.As` 类型匹配，避免拼写错误。
:::

## 📚 相关链接

- [CWEError 根类型](./cwe-error) | [APIError](./api-error) | [ParseError](./parse-error) | [HTTP 重试](./http-retry)
