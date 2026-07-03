---
title: CLI 接入（命令行与脚本）
outline: [2, 3]
---

# 💻 CLI 接入（命令行与脚本）

`cwe` 是 CWE Skills 的命令行入口：一个**单文件静态二进制**，40+ 子命令覆盖全部能力，所有命令支持 `-o text|json` 双格式输出。它是 Shell 脚本、CI/CD 流水线、快速手动查询的最佳选择，也是 [Skills 接入](./integration-skills) 的后端。

::: tip CLI = SDK 的薄包装
`cwe` 二进制内部调用的就是 `cweskills` Go 包，能力与 [Go SDK](./integration-sdk) 完全对等。区别只在「调用语法」——SDK 是 Go 函数，CLI 是子命令。两者数据互通、结果一致。
:::

---

## 📦 安装

### 1. curl 一键安装（推荐）

**Linux (amd64)**：

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/
```

**macOS (Apple Silicon)**：

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_darwin_aarch64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/
```

**Windows (PowerShell)**：

```powershell
Invoke-WebRequest -Uri https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_windows_x86_64.zip -OutFile cwe.zip
Expand-Archive cwe.zip
Move-Item cwe.exe C:\Windows\cwe.exe
```

### 2. go install（需本地 Go 工具链）

```bash
go install github.com/scagogogo/cwe-skills/cmd/cwe@latest
```

::: info 装到哪
二进制装到 `$GOBIN`（通常 `$GOPATH/bin` 或 `$HOME/go/bin`），请确保该目录在 `PATH` 中。
:::

### 3. 包管理器

```bash
# Homebrew (macOS/Linux)
brew install scagogogo/tap/cwe-skills

# Scoop (Windows)
scoop install cwe-skills
```

### 4. 验证

```bash
cwe version
```

```text
cwe CLI 版本: dev
SDK版本: v0.0.1
```

更多安装方式（deb/rpm/apk、Docker、源码编译）见 [安装](./installation)。

---

## 🗂️ 命令总览

`cwe` 的子命令按能力域分组。下表列出主要命令，完整列表见 [CLI 命令参考](../cli/overview)。

| 分组 | 命令 | 说明 |
|------|------|------|
| 🆔 ID 工具 | `cwe parse CWE-79` | 解析 CWE ID 为整数 |
| | `cwe validate CWE-79` | 验证 ID 格式 |
| | `cwe format 79` | 格式化为标准形式 |
| | `cwe extract "见 CWE-79 与 CWE-89"` | 从文本提取 ID |
| | `cwe compare CWE-79 CWE-89` | 比较两个 ID |
| 📚 枚举 | `cwe enum abstraction` | 列出抽象层级枚举 |
| | `cwe enum structure` / `status` / `relationship` … | 其他枚举 |
| 🏆 知名列表 | `cwe wellknown check CWE-79` | 检查 Top 25/OWASP/SANS |
| | `cwe wellknown top25` | 列出 Top 25 |
| 🌐 在线 API | `cwe show CWE-79` | 从 MITRE API 取弱点详情 |
| | `cwe relations CWE-79` | 在线关系查询 |
| | `cwe version mitre` | MITRE API 版本 |
| 📥 离线 XML | `cwe stats --xml <file>` | XML 目录统计 |
| | `cwe registry get CWE-79 --xml <file>` | 本地注册表查询 |
| 🧭 关系导航 | `cwe nav ancestors CWE-79 --xml <file>` | 祖先链 |
| | `cwe nav shortest-path CWE-79 CWE-1 --xml <file>` | 最短路径 |
| 🌳 树构建 | `cwe tree build CWE-1 --xml <file>` | 构建层次树 |
| 🔍 搜索过滤 | `cwe search --xml <file> --keyword Injection` | 关键词搜索 |
| | `cwe filter --xml <file> --abstraction Base --status Stable` | 多条件过滤 |

::: warning 离线命令需要 --xml
导航、建树、搜索、注册表查询等离线命令都需要 `--xml <file>` 指向本地 `cwec_v4.15.xml`。在线命令（`show`/`relations`/`version mitre`）无需此参数，但受速率限制。详见 [在线 vs 离线](./online-offline)。
:::

---

## 🌐 全局参数

所有命令共享以下全局参数：

| 参数 | 说明 | 示例 |
|------|------|------|
| `-o, --output` | 输出格式 `text` 或 `json`（默认 `text`） | `cwe show CWE-79 -o json` |
| `--xml` | 离线 XML 文件路径 | `--xml cwec_v4.15.xml` |
| `-h, --help` | 命令帮助 | `cwe parse --help` |
| `-v, --version` | 版本信息 | `cwe version` |

::: tip -o json 是脚本友好的关键
`-o json` 输出结构化 JSON，便于 `jq` 管道与 AI 解析。text 格式面向人类阅读。详见 [输出格式](./output-format)。
:::

---

## 📜 脚本集成示例

### 示例 1：JSON 管道 + jq 提取字段

```bash
# 取 CWE-79 的名称，只输出名字一行
cwe show CWE-79 -o json | jq -r '.name'
# => Cross-site Scripting (XSS)

# 批量查多个 CWE 的名称
for id in 79 89 352; do
  echo "$id: $(cwe show CWE-$id -o json | jq -r '.name')"
done
```

### 示例 2：从扫描报告提取 CWE 并查 Top 25

```bash
# 从文本报告提取所有 CWE ID，去重，逐个检查是否在 Top 25
cwe extract "$(cat scan-report.txt)" -o json | jq -r '.[]' | sort -u | while read cweid; do
  if cwe wellknown check "$cweid" -o json | jq -e '.top25' >/dev/null; then
    echo "⚠️ $cweid 在 Top 25，需优先修复"
  fi
done
```

### 示例 3：CI 流水线门禁

```bash
# .github/workflows/security.yml 片段
- name: 检查漏洞报告中的 CWE 是否高优
  run: |
    high_risk=$(cwe extract "$(cat findings.txt)" -o json \
      | jq -r '.[]' | sort -u \
      | while read id; do cwe wellknown check "$id" -o json | jq -r 'select(.top25) | .id'; done)
    if [ -n "$high_risk" ]; then
      echo "::error::发现 Top 25 高危 CWE: $high_risk"
      exit 1
    fi
```

### 示例 4：离线批量导航

```bash
# 离线计算一组 CWE 到 CWE-1 的最短路径，输出 JSON 数组
for id in 79 89 352 787; do
  cwe nav shortest-path "CWE-$id" CWE-1 --xml cwec_v4.15.xml -o json
done
```

::: details 为什么脚本里用 -o json 而不是 text
text 格式带人类可读的标签、缩进、说明文本，正则解析脆弱。JSON 字段稳定、结构化，`jq` 一行就能取值，CI 里更不容易因输出格式微调而崩。这也是 [Skills 接入](./integration-skills) 里反复提示 AI「加 `-o json`」的原因。
:::

---

## 🔗 与 SDK 的关系

CLI 和 SDK 背后是同一套 `cweskills` 包：

| 维度 | CLI | Go SDK |
|------|-----|--------|
| 调用方 | Shell / 脚本 / AI | Go 应用 |
| 语法 | 子命令 + 参数 | Go 函数 + 类型 |
| 输出 | text / JSON 字符串 | Go 对象（`*CWE` 等） |
| 错误 | 退出码 + stderr 文本 | `*CWEError` 结构化错误 |
| 类型检查 | 运行时 | 编译期 |

```bash
# CLI
cwe show CWE-79 -o json
```

```go
// 等价的 SDK 调用
weakness, _ := client.GetWeakness(ctx, 79)
```

::: tip 何时从 CLI 升级到 SDK
当你需要：① 在 Go 程序里深度集成；② 类型化对象而非字符串解析；③ 复杂数据流（注册表 + 导航 + 树组合）；④ 精细控制 HTTP/速率/重试。否则 CLI 足够。见 [Go SDK 接入](./integration-sdk)。
:::

---

## ⚠️ 注意事项

::: warning 在线命令受速率限制
`cwe show`、`cwe relations` 等在线命令受 MITRE API 速率限制（默认约 0.1 req/s）。脚本里批量循环调用会触发限流，CLI 会自动等待。大批量场景请改用离线 `--xml`。见 [速率限制与重试](./rate-limit-retry)。
:::

::: info CLI 是 Skills 的后端
[Skills 接入](./integration-skills) 让 AI 调用的就是这个 `cwe` 二进制。AI 所在环境必须装好 CLI 且在 PATH 中。若 AI 跑在无 CLI 的沙箱，需改用 SDK 或等待 [MCP](./integration-mcp)。
:::

---

## 📖 相关文档

- [四种接入方式总览](./integrations)
- [Go SDK 接入](./integration-sdk)
- [Skills 接入（AI 代理）](./integration-skills)
- [输出格式 (text/JSON)](./output-format)
- [在线 vs 离线模式](./online-offline)
- [安装](./installation)
