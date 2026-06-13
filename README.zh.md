# CWE Skills — AI原生的CWE Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/cwe-skills.svg)](https://pkg.go.dev/github.com/scagogogo/cwe-skills)
[![CI](https://github.com/scagogogo/cwe-skills/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/cwe-skills/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AI原生的 [CWE（通用缺陷枚举）](https://cwe.mitre.org/) SDK和CLI工具** — 为构建网络安全产品、SAST/DAST工具、漏洞管理平台和AI安全代理提供完整的API支持。

> 🇬🇧 [English](README.md)

## 为什么是"AI原生"？

CWE Skills 从设计之初就面向AI代理集成：

- **结构化JSON输出** — 每个CLI命令都支持`-o json`，机器可直接解析
- **零交互提示** — 所有操作非交互、可脚本化
- **渐进式技能文档** — AI代理可以渐进式发现能力
- **完整SDK API** — 100+公开函数，97.4%测试覆盖率
- **离线优先** — 加载MITRE XML目录支持离线环境
- **多格式序列化** — JSON/XML/CSV支持数据管道集成

## 功能特性

- **完整的CWE数据模型**：弱点、类别、视图、复合元素
- **类型化枚举**：抽象层级、状态值、关系类型、后果范围
- **CWE ID工具**：解析、格式化、验证和提取
- **知名列表**：CWE Top 25、OWASP Top 10、SANS Top 25
- **MITRE REST API客户端**：支持速率限制和重试
- **XML目录解析器**：离线解析MITRE官方XML
- **内存注册表**：存储、索引和查询CWE条目
- **搜索与过滤**：关键字、抽象层级、状态、可能性、后果范围
- **关系导航**：父/子/祖先/后代/同级/对等/链式/组合
- **树构建**：构建和遍历CWE层次树
- **序列化**：JSON、XML和CSV导入导出
- **Cobra CLI**：40+子命令，text/JSON双格式输出
- **零依赖**：核心SDK仅使用Go标准库

## 安装

### SDK

```bash
go get github.com/scagogogo/cwe-skills
```

### CLI — 从GitHub Release下载（推荐）

从 [Releases](https://github.com/scagogogo/cwe-skills/releases/latest) 下载：

```bash
# Linux (amd64)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_darwin_aarch64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_windows_x86_64.zip -OutFile cwe.zip
Expand-Archive cwe.zip
```

### CLI — 从源码编译

```bash
git clone https://github.com/scagogogo/cwe-skills.git
cd cwe-skills
go build -o cwe ./cmd/cwe/
```

### CLI — 包管理器

```bash
# Homebrew
brew install scagogogo/tap/cwe-skills

# Scoop (Windows)
scoop bucket add scagogogo https://github.com/scagogogo/scoop-bucket
scoop install cwe-skills

# Go Install
go install github.com/scagogogo/cwe-skills/cmd/cwe@latest
```

### 验证

```bash
cwe version
```

## 快速开始

### 使用CLI

```bash
# 解析和验证
cwe parse CWE-79 89 cwe-352
cwe validate CWE-79 CWE-89

# 格式化和提取
cwe format 79 89 352
cwe extract "受CWE-79和CWE-89影响"

# 检查知名列表
cwe wellknown top25
cwe wellknown check CWE-79

# 查询MITRE API（在线）
cwe show CWE-79
cwe relations parents CWE-79

# 本地搜索和过滤（离线，需XML目录）
cwe search --xml cwec_latest.xml --keyword Injection
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable

# 本地注册表操作（离线）
cwe registry get CWE-79 --xml cwec_latest.xml
cwe registry parents CWE-79 --xml cwec_latest.xml
cwe registry export --xml cwec_latest.xml --format json

# 本地关系导航（离线）
cwe nav siblings CWE-79 --xml cwec_latest.xml
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_latest.xml

# 树操作（离线）
cwe tree build CWE-1 --xml cwec_latest.xml
cwe tree forest --xml cwec_latest.xml

# 所有命令支持JSON输出
cwe parse CWE-79 -o json
```

### 使用Go SDK

```go
package main

import (
    "fmt"
    "context"
    cwepkg "github.com/scagogogo/cwe-skills"
)

func main() {
    // 解析CWE ID
    id, _ := cwepkg.ParseCWEID("CWE-79")
    fmt.Println(id) // 79

    // 查询MITRE API
    client := cwepkg.NewAPIClient()
    defer client.Close()
    weakness, _ := client.GetWeakness(context.Background(), 79)

    // 本地注册表
    registry := cwepkg.NewRegistry()
    registry.Register(&cwepkg.CWE{ID: 79, Name: "XSS", Abstraction: cwepkg.AbstractionBase})
    registry.BuildIndexes()

    // 导航关系
    nav := cwepkg.NewNavigator(registry)
    parents := nav.Parents(79)
    ancestors := nav.Ancestors(79)

    // 构建树
    tree := cwepkg.BuildTree(registry, 1)
    leaves := tree.LeafNodes()

    // 检查知名列表
    if cwepkg.IsInTop25(79) {
        fmt.Println("CWE-79在Top 25中！")
    }
}
```

### AI代理接入

CWE Skills专为AI代理集成设计。复制以下提示词到你的AI代理配置中：

```markdown
## CWE Skills 集成

你可以使用 `cwe` CLI工具进行CWE（通用缺陷枚举）操作。

### 安装
下载: https://github.com/scagogogo/cwe-skills/releases/latest
或从源码编译: `git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/`

### 常用命令
- `cwe parse CWE-79` — 解析CWE ID
- `cwe validate CWE-79` — 验证CWE ID格式
- `cwe show CWE-79` — 从MITRE API获取弱点详情
- `cwe wellknown check CWE-79` — 检查是否在Top 25/OWASP/SANS列表
- `cwe search --xml <file> --keyword <term>` — 搜索离线XML目录
- `cwe nav ancestors CWE-79 --xml <file>` — 离线导航关系
- `cwe tree build CWE-1 --xml <file>` — 构建层次树

### 输出格式
所有命令支持 `-o json` 输出结构化JSON。

### SDK (Go)
```go
import cwepkg "github.com/scagogogo/cwe-skills"
id, _ := cwepkg.ParseCWEID("CWE-79")
cwepkg.IsInTop25(79) // true
```

### 文档
完整技能文档: https://github.com/scagogogo/cwe-skills/tree/main/docs/skills
```

## CLI命令参考

| 命令 | 描述 |
|------|------|
| `cwe version` | 显示版本信息 |
| `cwe parse [IDs...]` | 解析CWE ID |
| `cwe validate [IDs...]` | 验证CWE ID格式 |
| `cwe format [IDs...]` | 格式化为CWE-NNN |
| `cwe extract [text...]` | 从文本提取CWE ID |
| `cwe compare <ID1> <ID2>` | 比较两个CWE ID |
| `cwe enum <type>` | 列出枚举值 |
| `cwe wellknown top25` | CWE Top 25列表 |
| `cwe wellknown owasp` | OWASP Top 10列表 |
| `cwe wellknown sans` | SANS Top 25列表 |
| `cwe wellknown check [IDs...]` | 检查列表成员 |
| `cwe show [IDs...]` | 从MITRE API获取 |
| `cwe relations parents/children/ancestors/descendants [ID]` | API关系查询 |
| `cwe api-version` | 检查MITRE API版本 |
| `cwe search --xml <file> [flags]` | 搜索离线XML |
| `cwe filter --xml <file> [flags]` | 多条件过滤 |
| `cwe stats --xml <file>` | 统计XML数据 |
| `cwe registry load/get/contains/... --xml <file>` | 注册表操作 |
| `cwe nav parents/children/siblings/peers/... --xml <file>` | 本地导航 |
| `cwe tree build/forest/view/path/leaves --xml <file>` | 树操作 |

## Skills技能文档

渐进式技能文档，面向AI代理和开发者：

| # | 技能 | 描述 |
|---|------|------|
| 1 | [CWE ID 解析与验证](docs/skills/01-cwe-id-parsing-validation.md) | 解析、验证、格式化CWE ID |
| 2 | [CWE ID 提取与比较](docs/skills/02-cwe-id-extraction-comparison.md) | 从文本提取、比较ID |
| 3 | [知名列表](docs/skills/03-well-known-lists.md) | CWE Top 25、OWASP Top 10、SANS Top 25 |
| 4 | [枚举类型](docs/skills/04-enumeration-types.md) | 抽象层级、状态、关系类型 |
| 5 | [API: 获取弱点详情](docs/skills/05-api-show-weakness.md) | 从MITRE API获取 |
| 6 | [API: 关系查询](docs/skills/06-api-relationships.md) | 通过API查询关系 |
| 7 | [API: 版本检查](docs/skills/07-api-version.md) | 检查MITRE API版本 |
| 8 | [本地: 搜索与过滤](docs/skills/08-local-search-filter.md) | 搜索和多条件过滤 |
| 9 | [本地: 注册表操作](docs/skills/09-local-registry.md) | 加载、查询、导出本地数据 |
| 10 | [本地: 关系导航](docs/skills/10-local-navigation.md) | 离线导航关系 |
| 11 | [本地: 树构建](docs/skills/11-local-tree.md) | 构建和遍历层次树 |
| 12 | [SDK: 序列化](docs/skills/12-sdk-serialization.md) | JSON、XML、CSV导入导出 |

→ **[完整技能索引](docs/skills/README.md)**

## 许可证

MIT许可证 - 详见 [LICENSE](LICENSE)
