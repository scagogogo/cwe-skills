---
title: 错误处理
outline: [2, 3]
---

# ⚠️ 错误处理

CWE Skills 把所有错误统一为 **`CWEError` 体系**：一个基础错误类型 + 一组场景化子类型。每个子类型携带与该场景相关的字段（如 `APIError` 带 `StatusCode`、`RateLimitError` 带 `RetryAfter`），且都实现 `Unwrap()`，可与标准库 `errors.Is` / `errors.As` 无缝协作。

::: tip 为什么不用 errorf 字符串
字符串错误无法被程序化判别（只能字符串匹配，脆弱）。`CWEError` 体系让调用方能**按类型分支处理**——404 走 fallback、429 走退避、解析错误走跳过——而不靠字符串猜。
:::

---

## 🌳 错误体系图

```text
                    CWEError (基础)
                  ┌────┴─────────────────────────┐
                  │ Code / Message                │ Unwrap()
                  └────┬──────────────────────────┘
            ┌──────────┼──────────────────────────┐
            │          │                          │
   InvalidCWEIDError  CWENotFoundError          APIError
   (输入的 ID 非法)   (ID 合法但不存在)          (statusCode/url/method)
            │          │                          │
            │          │                  ┌───────┴───────┐
            │          │                  │               │
            │          │          RateLimitError     ValidationError
            │          │          (retryAfter)        (field/value)
            │          │
            │          │          ParseError      RelationshipError
            │          │          (detail/offset) (from/to/nature)
```

::: info 都实现 Unwrap()
每个错误类型都实现 `Unwrap() error`，因此 `errors.Is(err, someSentinel)` 和 `errors.As(err, &typed)` 都能正常工作，包括包装嵌套场景。
:::

---

## 📋 错误类型一览

| 类型 | 触发场景 | 关键字段 |
|------|----------|----------|
| `CWEError` | 基础类型，其他错误的包装根 | `Code`, `Message` |
| `InvalidCWEIDError` | 输入的 CWE ID 格式非法（如 `"CWE-abc"`） | ID 字符串 |
| `CWENotFoundError` | ID 格式合法但 MITRE/注册表里没有 | `ID` |
| `APIError` | HTTP 调用失败（非 429） | `StatusCode`, `URL`, `Method` |
| `RateLimitError` | 触发 429 速率限制 | `RetryAfter` |
| `ValidationError` | 字段值校验失败 | `Field`, `Value` |
| `ParseError` | XML/JSON 解析失败 | `Detail`, `Offset` |
| `RelationshipError` | 关系操作非法（如不存在的 from/to） | `From`, `To`, `Nature` |

---

## 🎯 errors.As：按类型分支

`errors.As` 把错误「拆」成具体类型，读取其字段：

```go
import "errors"

weakness, err := client.GetWeakness(ctx, 999999)
if err != nil {
    var apiErr *cweskills.APIError
    if errors.As(err, &apiErr) {
        fmt.Println("HTTP 状态码:", apiErr.StatusCode) // 404
        fmt.Println("请求 URL:", apiErr.URL)
        fmt.Println("方法:", apiErr.Method)
    }

    var nfErr *cweskills.CWENotFoundError
    if errors.As(err, &nfErr) {
        fmt.Println("未找到 CWE ID:", nfErr.ID)
    }
}
```

::: tip errors.As 链式匹配
错误可能被多层包装。`errors.As` 会沿 `Unwrap()` 链向下查找第一个匹配的类型，所以即使错误被 `fmt.Errorf("...: %w", err)` 包装过，仍能匹配到内部的具体类型。
:::

---

## 🔍 errors.Is：按语义判定

`errors.Is` 判定错误是否「是」某个哨兵错误（或可等价的错误）：

```go
// 判定是否某种已知错误
if errors.Is(err, cweskills.ErrCWENotFound) {
    // 走 fallback 逻辑
}

// 判定是否速率限制（无论包装层级）
var rlErr *cweskills.RateLimitError
if errors.As(err, &rlErr) {
    time.Sleep(rlErr.RetryAfter)
}
```

::: details errors.Is vs errors.As
- `errors.Is(err, target)`：判定 err 是否等于 target（或其包装链中含 target）。用于**哨兵错误**（sentinel）。
- `errors.As(err, &typed)`：把 err 拆成 typed 类型实例，读字段。用于**结构化错误**。
CWE Skills 推荐用 `errors.As` 读字段（因为错误都带数据），`errors.Is` 用于快速语义判定。
:::

---

## 🎬 每种错误的触发场景

### 1. InvalidCWEIDError

```go
_, err := cweskills.ParseCWEID("CWE-abc") // 非法
// → *InvalidCWEIDError
```

**恢复策略**：提示用户重新输入，或用 `ExtractCWEIDs` 从文本里清洗。

### 2. CWENotFoundError

```go
_, err := client.GetWeakness(ctx, 999999) // ID 合法但不存在
// → *CWENotFoundError
```

**恢复策略**：fallback 到已知列表，或提示「该 CWE 不在当前版本中」。

### 3. APIError

```go
_, err := client.GetWeakness(ctx, 79) // 网络故障 / 5xx
// → *APIError（StatusCode=500 等）
```

**恢复策略**：检查 `StatusCode`——5xx 走重试（SDK 已内置），4xx（非 404/429）检查请求参数。

### 4. RateLimitError

```go
_, err := client.GetWeakness(ctx, 79) // 触发 429
// → *RateLimitError（RetryAfter）
```

**恢复策略**：等待 `RetryAfter` 后重试，或切换到离线 XML。见 [速率限制与重试](./rate-limit-retry)。

### 5. ValidationError

```go
// 传入非法枚举值或字段
// → *ValidationError（Field="abstraction", Value="Invalid"）
```

**恢复策略**：用 `cwe enum <type>` 列出合法值，修正后重试。

### 6. ParseError

```go
_, err := cweskills.NewXMLParser().ParseFile("corrupted.xml")
// → *ParseError（Detail, Offset）
```

**恢复策略**：检查 `Offset` 定位 XML 损坏位置，重新下载官方 XML。

### 7. RelationshipError

```go
nav := cweskills.NewNavigator(registry)
nav.ShortestPath(999999, 1) // from 不存在
// → *RelationshipError（From=999999, To=1, Nature="shortest-path"）
```

**恢复策略**：先用 `registry.Contains(id)` 确认节点存在，再导航。

---

## 🧩 统一错误处理模式

一个健壮的调用通常这样处理错误：

```go
weakness, err := client.GetWeakness(ctx, 79)
if err != nil {
    var nfErr *cweskills.CWENotFoundError
    if errors.As(err, &nfErr) {
        return fmt.Errorf("CWE-%d 不存在: %w", nfErr.ID, err)
    }

    var rlErr *cweskills.RateLimitError
    if errors.As(err, &rlErr) {
        time.Sleep(rlErr.RetryAfter)
        return retry(ctx, client, 79) // 自定义重试
    }

    var apiErr *cweskills.APIError
    if errors.As(err, &apiErr) && apiErr.StatusCode >= 500 {
        return fmt.Errorf("MITRE 服务端故障，稍后重试: %w", err)
    }

    // 兜底
    return fmt.Errorf("查询失败: %w", err)
}
```

::: tip 用 %w 包装保留链
`fmt.Errorf("...: %w", err)` 用 `%w` 动词包装错误，保留 `Unwrap()` 链，上层仍能用 `errors.As` 匹配内部类型。不要用 `%v`，那会断链。
:::

---

## 💻 CLI 的错误呈现

CLI 把 `CWEError` 转成人类可读的 stderr 文本 + 非零退出码，便于脚本判断：

```bash
$ cwe show CWE-999999
Error: CWE-999999 未找到 (CWE_NOT_FOUND)
exit code: 1

$ cwe show CWE-79 -o json | jq '.name'   # JSON 输出不含错误，错误走 stderr
```

::: info 脚本里判断错误
脚本里用退出码判断成败：`0` 成功，非 `0` 失败。错误详情在 stderr。`-o json` 输出只在成功时产生；失败时 JSON 不输出，避免污染管道。见 [输出格式](./output-format)。
:::

---

## 📖 相关文档

- [速率限制与重试](./rate-limit-retry)（`RateLimitError` 与 429 处理）
- [在线 vs 离线模式](./online-offline)
- [Go SDK 接入](./integration-sdk)
- [CLI 接入](./integration-cli)
- [输出格式 (text/JSON)](./output-format)
