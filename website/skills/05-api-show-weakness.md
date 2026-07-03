---
title: 技能 05 — API 获取弱点详情
outline: [2, 3]
---

# 🌐 技能 05 — API 获取弱点详情

通过 MITRE CWE REST API 获取弱点的详细信息：名称、抽象层级、状态、描述、后果、缓解措施、关系等。

<Badge type="tip" text="在线 API"/>
<Badge type="warning" text="需网络"/>

---

## 🎯 技能目标

- 在线获取一个或多个 CWE 弱点的完整详情
- 获取 CWE 类别（Category）与视图（View）详情
- 配置 API 基础 URL、超时、速率限制、重试

---

## 💻 CLI 命令

### show — 获取弱点

```bash
cwe show CWE-79
cwe show 79 89 352
cwe show --base-url https://cwe-api.mitre.org/api CWE-79
cwe show --timeout 60 CWE-79
```

```text
=== CWE-79 ===
  名称:     Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')
  抽象层级: Base
  状态:     Stable
  描述:     The product does not neutralize...
  结构:     Simple
  关系:     4 项
```

JSON 输出包含完整结构化字段：`id`、`name`、`abstraction`、`status`、`description` 等。

### show category / show view

```bash
cwe show category 1     # 类别详情
cwe show view 1000      # 视图详情
```

### Flags

| Flag | 默认值 | 说明 |
|------|--------|------|
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |
| `--timeout` | `30` | 请求超时（秒） |

---

## 🔧 SDK API

### NewAPIClient

```go
// 默认客户端
client := cweskills.NewAPIClient()
defer client.Close()

// 带选项
client := cweskills.NewAPIClient(
    cweskills.WithAPIBaseURL("https://cwe-api.mitre.org/api"),
    cweskills.WithAPITimeout(30 * time.Second),
    cweskills.WithAPIRateLimit(10, time.Second),
    cweskills.WithAPIRetry(3),
)
```

### GetWeakness

```go
weakness, err := client.GetWeakness(ctx, 79)
fmt.Println(weakness.Name)         // "Improper Neutralization of Input..."
fmt.Println(weakness.Abstraction)  // "Base"
fmt.Println(weakness.Status)       // "Stable"
fmt.Println(weakness.Description)
```

### GetCategory / GetView

```go
category, _ := client.GetCategory(ctx, 1)
view, _ := client.GetView(ctx, 1000)
fmt.Println(view.Name)  // "Research Concepts"
```

::: details 客户端选项
| 选项 | 说明 |
|------|------|
| `WithAPIBaseURL(url)` | 覆盖默认 API 基础 URL |
| `WithAPITimeout(d)` | 设置 HTTP 超时 |
| `WithAPIRateLimit(n, interval)` | 速率限制：每 interval N 个请求 |
| `WithAPIRetry(max)` | 瞬时错误自动重试 |
| `WithAPIHTTPClient(c)` | 使用自定义 HTTP 客户端 |
:::

::: details 错误类型
- `CWENotFoundError`：CWE ID 在 API 中不存在
- `APIError`：HTTP 错误（5xx、网络故障）
- `RateLimitError`：超出速率限制
- `InvalidCWEIDError`：CWE ID 格式无效
:::

---

## 📝 示例

### 命令行

```bash
# 批量取详情并提取名称
cwe show 79 89 352 -o json | jq '.[].detail.name'

# 检查 API 是否可达
cwe show CWE-79 --timeout 5 -o json > /dev/null && echo "OK"
```

### Go

```go
package main

import (
    "context"
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewAPIClient()
    defer client.Close()

    w, err := client.GetWeakness(context.Background(), 79)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s (%s)\n", w.Name, w.Abstraction)
}
```

---

## 🤖 AI 代理使用提示

- 用户问「CWE-79 是什么」时，AI 用 `cwe show CWE-79 -o json` 取详情再总结。
- 批量查询时建议让 AI 加 `-o json`，结构化字段解析无歧义。
- 在线命令受速率限制，AI 连续调用可能触发自动等待，不是出错。

::: tip 记得 defer client.Close()
SDK 使用时务必 `defer client.Close()` 释放资源。批量查询用 `WithAPIRateLimit` 避免 MITRE 限流。
:::

---

## 📖 相关文档

- [技能 06 — API 关系查询](./06-api-relationships)
- [CLI: show](../cli/show) · [show category](../cli/show-category) · [show view](../cli/show-view)
- [SDK: GetWeakness](../sdk/api-get-weakness) · [NewAPIClient](../sdk/new-api-client)
- [速率限制与重试](../guide/rate-limit-retry)
- [返回 Skills 总览](./)
