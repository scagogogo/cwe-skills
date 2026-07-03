---
title: Skills 接入总览
outline: [2, 3]
---

# 🦾 Skills 接入总览

**Skills** 是 CWE Skills 四种接入方式中最面向 **AI 代理**的一种：把一段提示词复制进 AI 代理的系统提示词，AI 就能自主调用 `cwe` CLI 完成全部 CWE 操作——解析 ID、查弱点、查关系、判断是否在 Top 25、构建层次树……

<Badge type="tip" text="零代码"/>
<Badge type="info" text="面向 AI 代理"/>
<Badge type="warning" text="需安装 cwe CLI"/>

::: tip 为什么 Skills 面向 AI 代理
AI 代理（Claude、GPT 等）擅长「决定调用哪个命令、解读 JSON 结果、给出结论」，但不擅长「处理 HTTP 边角、记速率限制、解析半结构化文本」。Skills 把这些脏活儿交给 `cwe` CLI，AI 只管推理与决策。
:::

---

## 🎯 什么是 Skills 接入

Skills = **一段 Markdown 提示词 + 一个已安装的 `cwe` CLI**。

- 提示词告诉 AI：「你可以用 `cwe` 命令做 CWE 操作，这是命令清单」。
- AI 收到用户自然语言提问后，自主选择并调用合适的 `cwe` 命令，读取输出（推荐 `-o json`），再回答用户。

整个过程不需要写一行 Go 代码，也不需要 AI 代理内置 CWE 知识。

---

## 📋 前置条件

1. **AI 代理**：支持执行外部命令或具备工具调用能力的代理（Claude、GPT 等）。
2. **`cwe` CLI**：已安装在 AI 代理可执行的环境里。

```bash
# 验证 CLI 可用
cwe version
```

3. **（可选）离线 XML 目录**：若要让 AI 做导航/建树/搜索，需在环境里放一份 `cwec_v4.15.xml`，从 [MITRE](https://cwe.mitre.org/data/xml.html) 下载。

::: warning CLI 必须装在 AI 能跑的环境
若 AI 运行在沙箱、无法执行外部命令，则 Skills 不可用——此时改用 [Go SDK](../sdk/overview) 或等待 [MCP](../guide/integration-mcp)。
:::

---

## 📝 安装提示词

把下面这段 Markdown 复制到 AI 代理的系统提示词、或技能配置文件中：

```markdown
## CWE Skills

你可以使用 `cwe` CLI 工具进行 CWE（通用缺陷枚举）操作。

### 安装
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz && sudo mv cwe /usr/local/bin/

### 核心命令
| 命令 | 功能 |
|------|------|
| `cwe parse CWE-79` | 解析 CWE ID |
| `cwe validate CWE-79` | 验证 CWE ID 格式 |
| `cwe show CWE-79` | 从 MITRE API 获取弱点详情 |
| `cwe wellknown check CWE-79` | 检查是否在 Top 25/OWASP/SANS |
| `cwe enum abstraction` | 列出枚举值 |
| `cwe search --xml <file> --keyword Injection` | 搜索离线 XML 目录 |
| `cwe filter --xml <file> --abstraction Base --status Stable` | 多条件过滤 |
| `cwe registry get CWE-79 --xml <file>` | 从本地注册表获取条目 |
| `cwe nav ancestors CWE-79 --xml <file>` | 离线导航关系 |
| `cwe nav shortest-path CWE-79 CWE-1 --xml <file>` | 查找两个 CWE 间最短路径 |
| `cwe tree build CWE-1 --xml <file>` | 构建层次树 |
| `cwe stats --xml <file>` | XML 目录统计 |

### 输出格式
所有命令支持 `-o json` 输出结构化 JSON。示例: `cwe parse CWE-79 -o json`
```

::: details 更完整的提示词
仓库根目录的 `README.zh.md` 含更完整版提示词（含 Go SDK 摘要）。渐进式技能文档见 `docs/skills/` 目录，也可整体粘贴进 AI 上下文以增强能力。
:::

---

## 🧠 12 个渐进式技能索引

12 篇技能文档从简到深，覆盖 CLI 命令、SDK API 与示例。可直接作为 AI 能力参考，也可逐篇喂给 AI。

| # | 技能 | CLI 命令 | 文档 |
|---|------|----------|------|
| 1 | 🆔 CWE ID 解析与验证 | `parse` `validate` `format` | [01 →](./01-cwe-id-parsing-validation) |
| 2 | 🔍 CWE ID 提取与比较 | `extract` `compare` | [02 →](./02-cwe-id-extraction-comparison) |
| 3 | 🏆 知名列表 | `wellknown top25/owasp/sans/check` | [03 →](./03-well-known-lists) |
| 4 | 📚 枚举类型 | `enum` | [04 →](./04-enumeration-types) |
| 5 | 🌐 API 获取弱点详情 | `show` | [05 →](./05-api-show-weakness) |
| 6 | 🧭 API 关系查询 | `relations` | [06 →](./06-api-relationships) |
| 7 | 📦 API 版本检查 | `api-version` | [07 →](./07-api-version) |
| 8 | 🔎 本地搜索与过滤 | `search` `filter` `stats` | [08 →](./08-local-search-filter) |
| 9 | 🗃️ 本地注册表 | `registry` | [09 →](./09-local-registry) |
| 10 | 🧭 本地关系导航 | `nav` | [10 →](./10-local-navigation) |
| 11 | 🌳 本地树构建 | `tree` | [11 →](./11-local-tree) |
| 12 | 📦 SDK 序列化 | — | [12 →](./12-sdk-serialization) |

::: tip 喂给 AI 增强能力
把这 12 篇（或其中几篇）粘进 AI 上下文，能让 AI 更精确地知道「什么场景该用哪个命令」。文档越全，AI 调用越准。
:::

---

## 💬 使用示例

配置好提示词后，直接用自然语言向 AI 提问：

```text
你: CWE-89 是什么？在不在 Top 25？

AI: （调用 cwe show CWE-89、cwe wellknown check CWE-89）
    CWE-89 是 SQL 注入。它在 CWE Top 25 中排名靠前，
    同时属于 OWASP Top 10 的 A03:2021-Injection 类别。
```

```text
你: 给我看 CWE-79 到 CWE-1 的祖先链，用本地的 cwec_v4.15.xml

AI: （调用 cwe nav ancestors CWE-79 --xml cwec_v4.15.xml）
    CWE-79 的祖先链为：CWE-79 → CWE-74（注入）→ CWE-707 → ...
```

::: tip 让 AI 用 JSON 输出更稳
AI 解析 CLI 的 text 输出可能有歧义，建议让 AI 加 `-o json`，结果结构化、解析无歧义。提示词里已写明「所有命令支持 `-o json`」。
:::

---

## 📖 相关文档

- [Skills 接入指南](../guide/integration-skills)
- [CLI 命令参考](../cli/overview)
- [Go SDK API](../sdk/overview)
- [在线 vs 离线模式](../guide/online-offline)
- [示例与教程](../examples/)
