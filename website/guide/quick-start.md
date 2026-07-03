---
title: 快速开始
outline: [2, 3]
---

# 🚀 快速开始

本页带你用 **5 分钟**跑通 CWE Skills：安装 CLI → 跑第一条命令 → 写第一个 SDK 程序。三段任选其一，按你需要的方式接入。

---

## 🅰️ 方式一：CLI（最快上手）

### 1. 安装

从 [Releases](https://github.com/scagogogo/cwe-skills/releases/latest) 下载预编译二进制（Linux amd64 为例）：

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/
```

macOS (Apple Silicon) 把 `linux_x86_64` 换成 `darwin_aarch64`；Windows 下载 `.zip` 解压。更多方式见 [安装](./installation)。

### 2. 验证

```bash
cwe version
```

```text
cwe CLI 版本: dev
SDK版本: v0.0.1
```

### 3. 第一条命令：解析与知名列表

```bash
# 解析多种写法的 CWE ID
cwe parse CWE-79 89 cwe-352

# 验证格式
cwe validate CWE-79 CWE-89

# 检查是否属于 Top 25 / OWASP / SANS
cwe wellknown check CWE-79
```

输出类似：

```text
CWE-79 (CWE-79): 跨站脚本(XSS)
  ✓ CWE Top 25
  ✓ OWASP Top 10
  ✓ SANS Top 25
```

### 4. 在线查详情（MITRE API）

```bash
cwe show CWE-79
cwe relations parents CWE-79
cwe api-version
```

::: info 首次在线查询可能稍慢
MITRE API 有速率限制（默认约每 10 秒 1 个请求），首次调用若触发限流会自动等待。详见 [速率限制与重试](./rate-limit-retry)。
:::

### 5. 离线搜索与导航（需先下载 XML）

从 [MITRE 下载页](https://cwe.mitre.org/data/downloads.html) 下载 `cwec_v4.15.xml`，然后：

```bash
# 关键字搜索
cwe search --xml cwec_v4.15.xml --keyword Injection

# 多条件过滤
cwe filter --xml cwec_v4.15.xml --abstraction Base --status Stable

# 离线导航：CWE-79 的祖先链
cwe nav ancestors CWE-79 --xml cwec_v4.15.xml

# 最短路径
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_v4.15.xml

# 建树
cwe tree build CWE-1 --xml cwec_v4.15.xml
```

### 6. JSON 输出（脚本/AI 友好）

所有命令加 `-o json` 即输出结构化 JSON：

```bash
cwe parse CWE-79 -o json
cwe wellknown check CWE-79 -o json | jq .
```

详见 [输出格式](./output-format) 与 [CLI 命令](../cli/overview)。

---

## 🅱️ 方式二：Go SDK

### 1. 安装

```bash
go get github.com/scagogogo/cwe-skills
```

### 2. 第一个程序

新建 `main.go`：

```go
package main

import (
    "context"
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    // 1) ID 工具
    id, _ := cweskills.ParseCWEID("cwe-079")
    fmt.Printf("解析得到 ID = %d\n", id)

    // 2) 知名列表
    fmt.Printf("CWE-79 在 Top 25? %v\n", cweskills.IsInTop25(id))
    fmt.Printf("CWE-89 的 OWASP 类别: %s\n", cweskills.GetOWASPCategory(89))

    // 3) 在线 API
    client := cweskills.NewAPIClient()
    defer client.Close()
    w, err := client.GetWeakness(context.Background(), id)
    if err != nil {
        fmt.Println("API 调用失败:", err)
        return
    }
    fmt.Printf("在线详情: %s\n", w.Name)

    // 4) 枚举
    for _, a := range cweskills.AllAbstractionValues() {
        fmt.Println("抽象层级:", a)
    }
}
```

运行：

```bash
go run main.go
```

### 3. 离线版（注册表 + 导航 + 树）

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()

nav := cweskills.NewNavigator(registry)
fmt.Println("CWE-79 祖先数:", len(nav.Ancestors(79)))
fmt.Println("最短路径 79→1:", nav.ShortestPath(79, 1))

tree := cweskills.BuildTree(registry, 1)
fmt.Println("CWE-1 子树叶子数:", len(tree.LeafNodes()))

jsonData, _ := registry.ExportJSON()
fmt.Printf("导出 JSON %d 字节\n", len(jsonData))
```

---

## 🅲️ 方式三：Skills（AI 代理）

把 [Skills 提示词](./integration-skills) 复制到你的 AI 代理（Claude 等）系统提示词或技能配置中，AI 即可自主调用 `cwe` CLI。

示例对话：

```text
你: 帮我查一下 CWE-79 是什么，是否在 Top 25，它的祖先链是什么。
AI: （调用 cwe show CWE-79、cwe wellknown check CWE-79、cwe nav ancestors CWE-79 --xml ...）
    CWE-79 是跨站脚本(XSS)，在 Top 25 第 1 位，祖先链为 CWE-79 → CWE-74 → ...
```

::: tip Skills 不需要写代码
Skills 接入零代码——AI 通过 `cwe` CLI 完成一切，输出用 `-o json` 便于 AI 解析。详见 [Skills 接入](./integration-skills)。
:::

---

## ✅ 自检清单

跑完上面任意一段，你应当能做到：

- [ ] `cwe version` 正常输出版本
- [ ] `cwe parse` / `cwe validate` 能处理多种写法
- [ ] `cwe wellknown check` 能判断 Top 25
- [ ] SDK 能 `ParseCWEID` / `IsInTop25` / `NewAPIClient().GetWeakness`
- [ ] 知道在线 vs 离线两条路（[在线 vs 离线](./online-offline)）

---

## 📖 下一步

- 完整安装（30+ 平台 / Homebrew / Scoop / deb/rpm）→ [安装](./installation)
- 理解在线 vs 离线取舍 → [在线 vs 离线模式](./online-offline)
- CWE 概念入门 → [CWE 是什么](./concept-cwe)
- 全部 CLI 命令 → [CLI 命令](../cli/overview)
- 全部 SDK API → [SDK API](../sdk/overview)
