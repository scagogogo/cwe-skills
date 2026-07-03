---
title: 技能 07 — API 版本检查
outline: [2, 3]
---

# 📦 技能 07 — API 版本检查

检查 MITRE CWE REST API 的版本。用于验证 API 连通性、确认当前 CWE 数据版本。

<Badge type="tip" text="在线 API"/>
<Badge type="info" text="轻量调用"/>

---

## 🎯 技能目标

- 查询 MITRE CWE API 的版本号、发布日期、版本名称
- 用作连通性健康检查

---

## 💻 CLI 命令

### api-version

```bash
cwe api-version
cwe api-version --base-url https://cwe-api.mitre.org/api
```

```text
MITRE CWE API版本: 4.15
发布日期: 2024-11-19
版本名称: CWE v4.15
```

JSON 输出：

```json
{
  "version": "4.15",
  "releaseDate": "2024-11-19",
  "name": "CWE v4.15"
}
```

---

## 🔧 SDK API

### GetVersion

```go
version, err := client.GetVersion(ctx)
fmt.Println(version.Version)     // "4.15"
fmt.Println(version.ReleaseDate) // "2024-11-19"
fmt.Println(version.Name)        // "CWE v4.15"
```

### VersionResponse 结构

```go
type VersionResponse struct {
    Version     string `json:"version"`
    ReleaseDate string `json:"releaseDate"`
    Name        string `json:"name"`
}
```

---

## 📝 示例

### 命令行

```bash
# 健康检查
cwe api-version -o json | jq -r .version

# 验证自建镜像
cwe api-version --base-url https://my-cwe-mirror.local/api
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

    v, err := client.GetVersion(context.Background())
    if err != nil {
        fmt.Println("API 不可达:", err)
        return
    }
    fmt.Printf("已连接，当前版本 %s（%s）\n", v.Version, v.ReleaseDate)
}
```

---

## 🤖 AI 代理使用提示

- AI 在执行任何在线操作前，可用 `cwe api-version` 快速确认网络与 API 可达。
- 这是一个轻量调用，没有速率限制顾虑，适合做健康检查。

::: tip 用作健康检查
版本检查是最轻量的 API 调用，适合在脚本启动或 CI 流水线里验证 MITRE API 连通性。
:::

---

## 📖 相关文档

- [技能 05 — API 获取弱点详情](./05-api-show-weakness)
- [CLI: api-version](../cli/api-version)
- [SDK: GetVersion](../sdk/api-get-version)
- [返回 Skills 总览](./)
