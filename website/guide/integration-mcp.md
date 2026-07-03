---
title: MCP 接入（规划中）
outline: [2, 3]
---

# 🌐 MCP 接入（规划中）

**MCP**（Model Context Protocol）是 AI 工具调用的标准化协议。CWE Skills 计划提供官方 MCP 服务器，让任何 MCP 兼容的 AI 工具能以「工具调用」方式访问全部 CWE 能力——无需 AI 自己跑 Shell 命令、无需解析 CLI 文本输出。

::: warning 状态：规划中
<Badge type="warning" text="规划中" /> MCP 服务器**尚未发布**。本页说明我们的愿景、动机与路线图。在 MCP 落地前，请使用 [Skills 接入](./integration-skills)（AI 跑 CLI）或 [Go SDK](./integration-sdk) 作为替代方案。
:::

---

## 🤔 什么是 MCP

MCP（Model Context Protocol）是一个开放协议，为 AI 应用与外部数据源/工具之间提供**标准化通信**。类比 USB-C：不管什么厂商的 AI 工具，只要支持 MCP，就能统一地「发现并调用」一个 MCP 服务器暴露的能力。

```text
AI 应用（Claude Desktop / IDE / 任意 MCP 客户端）
        │  MCP 协议（JSON-RPC over stdio / HTTP）
        ▼
MCP 服务器（cwe-skills-mcp）
        │  暴露工具：get_weakness / find_ancestors / search …
        ▼
cweskills 包（在线 API + 离线 XML + 注册表 + 导航）
```

::: info MCP 的核心概念
- **工具（Tool）**：服务器暴露的可调用函数，带 JSON Schema 描述参数。
- **资源（Resource）**：服务器暴露的数据源（如只读的 CWE 注册表）。
- **提示（Prompt）**：服务器提供的预设提示词模板。
- **传输**：stdio（本地）或 HTTP（远程）。
:::

---

## 🎯 为什么要做 MCP

CWE Skills 已经有 [Skills](./integration-skills)、[CLI](./integration-cli)、[Go SDK](./integration-sdk) 三种接入方式，为什么还要 MCP？

### 现状的痛点

| 方式 | 痛点 |
|------|------|
| Skills（AI 跑 CLI） | AI 必须能执行 Shell 命令；沙箱/AI 平台禁用 shell 时不可用；AI 要解析 CLI 文本输出 |
| CLI 直连 | 每个工具调用都要 spawn 进程，开销大；参数与返回都是字符串，无 schema |
| Go SDK | 只能嵌入 Go 应用，Python/Node/其他语言生态用不上 |

### MCP 的解法

- **标准化工具调用**：AI 通过 JSON-RPC 调用工具，参数与返回都是结构化 JSON，带 schema，AI 不用猜参数。
- **不依赖 Shell**：MCP 服务器是长驻进程，AI 客户端通过协议通信，无需执行 Shell——沙箱环境也能用。
- **语言无关**：MCP 客户端跨语言，Python/Node 的 AI 工具也能接。
- **能力发现**：AI 客户端能自动列出服务器提供的工具清单，无需把提示词塞进上下文。

::: tip 一句话总结
Skills 让 AI「跑命令」用 CWE，MCP 让 AI「调工具」用 CWE——后者更标准、更安全、更适配受限环境。
:::

---

## 🗺️ 愿景

MCP 服务器计划暴露的工具集（与 CLI 子命令一一对应）：

| 工具 | 对应 CLI | 能力 |
|------|----------|------|
| `parse_cwe_id` | `cwe parse` | 解析 CWE ID |
| `validate_cwe_id` | `cwe validate` | 验证 ID 格式 |
| `extract_cwe_ids` | `cwe extract` | 从文本提取 ID |
| `get_weakness` | `cwe show` | 在线取弱点详情 |
| `get_relations` | `cwe relations` | 在线关系查询 |
| `check_wellknown` | `cwe wellknown check` | Top 25/OWASP/SANS 检查 |
| `list_enum` | `cwe enum` | 枚举值列表 |
| `search_xml` | `cwe search` | 离线关键词搜索 |
| `filter_xml` | `cwe filter` | 离线多条件过滤 |
| `get_ancestors` / `get_shortest_path` | `cwe nav ...` | 离线关系导航 |
| `build_tree` | `cwe tree build` | 层次树构建 |

::: details 设计原则
1. **工具命名与 CLI 对齐**：熟悉 CLI 的用户能秒懂 MCP 工具。
2. **JSON 进、JSON 出**：所有工具参数与返回都是结构化 JSON，带 JSON Schema。
3. **在线/离线透明**：工具内部自动选择路径（`get_weakness` 走 API，`get_ancestors` 走 XML），调用方只看结果。
4. **复用 cweskills 包**：MCP 服务器是 `cweskills` 包的薄包装，不重复实现逻辑。
:::

---

## 🛣️ 路线图

| 阶段 | 内容 | 状态 |
|------|------|------|
| 1 | 基础 stdio MCP 服务器，暴露 ID 工具 + 枚举 + 知名列表 | <Badge type="info" text="待启动" /> |
| 2 | 接入在线 API 工具（`get_weakness` / `get_relations`） | <Badge type="info" text="规划中" /> |
| 3 | 接入离线 XML 工具（搜索 / 过滤 / 导航 / 树），支持 `--xml` 配置 | <Badge type="info" text="规划中" /> |
| 4 | 资源（Resource）暴露：把注册表作为可读资源 | <Badge type="info" text="规划中" /> |
| 5 | HTTP 传输模式，支持远程部署 | <Badge type="info" text="规划中" /> |

::: warning 时间未定
MCP 是明确方向，但当前优先完善 SDK / CLI / Skills 的稳定性与文档。MCP 服务器将在核心能力稳定后启动。欢迎在 [GitHub Issues](https://github.com/scagogogo/cwe-skills/issues) 提需求与用例。
:::

---

## 🔄 当前替代方案

在 MCP 落地前，按你的场景选择替代方案：

### 场景 1：AI 代理想用 CWE，且能跑 Shell

→ 用 [Skills 接入](./integration-skills)：把提示词放进 AI 系统提示词，AI 自主调用 `cwe` CLI。

```text
你: 查 CWE-79 的祖先链
AI: （调用 cwe nav ancestors CWE-79 --xml cwec_v4.15.xml -o json）
```

### 场景 2：AI 在沙箱，不能跑 Shell

→ 用 [Go SDK](./integration-sdk) 在你的应用层集成 CWE，把结果作为结构化数据喂给 AI（而非让 AI 自己调工具）。

### 场景 3：非 Go 应用（Python/Node）想集成

→ 当前只能调用 [CLI](./integration-cli)（任何语言都能 spawn 子进程解析 JSON 输出）。MCP 落地后这条路径会更优雅。

::: tip 关注进展
Star 仓库 [github.com/scagogogo/cwe-skills](https://github.com/scagogogo/cwe-skills) 关注 MCP 进展。如果你有强烈的 MCP 用例，欢迎在 Issues 里描述，会优先考虑。
:::

---

## ❓ 常见疑问

::: details MCP 和 Skills 有什么区别？
Skills 是「给 AI 一段提示词，让它跑 CLI」——依赖 AI 能执行 Shell。MCP 是「给 AI 一个标准化的工具服务器」——AI 通过协议调用，不碰 Shell。MCP 更适合受限环境与跨语言生态；Skills 更轻量、零基础设施。
:::

::: details MCP 会取代 Skills / CLI 吗？
不会。四者会长期共存：SDK 嵌入 Go 应用，CLI 给 Shell/CI，Skills 给能跑命令的 AI，MCP 给 MCP 生态的 AI 工具。按场景选。
:::

::: details 我能自己包一层 MCP 吗？
可以。MCP 协议开放，你完全可以用现有的 `cwe` CLI 或 `cweskills` SDK 包一层 MCP 服务器。官方服务器只是想提供开箱即用的标准实现。
:::

---

## 📖 相关文档

- [四种接入方式总览](./integrations)
- [Skills 接入（AI 代理）](./integration-skills)（当前 AI 接入首选）
- [CLI 接入](./integration-cli)（当前跨语言集成首选）
- [Go SDK 接入](./integration-sdk)
- [工作原理](./how-it-works)
