---
layout: home

hero:
  name: CWE Skills
  text: AI 原生的 CWE 集成层
  tagline: 统一 MITRE REST API、离线 XML 目录与权威列表，四种方式接入你的安全工作流
  image:
    src: /architecture.png
    alt: CWE Skills 架构图
  actions:
    - theme: brand
      text: 🚀 快速开始
      link: /guide/quick-start
    - theme: alt
      text: 📖 了解项目
      link: /guide/what-is-cwe-skills
    - theme: alt
      text: 🔧 SDK API
      link: /sdk/overview

features:
  - icon: 🆔
    title: CWE ID 工具
    details: 解析、格式化、验证、从文本提取、比较 CWE ID。支持 CWE-79、cwe 79、79 等多种写法，统一规范化输出。
    link: /sdk/cwe-utils
    linkText: 查看 ID 工具 →
  - icon: 🏆
    title: 权威知名列表
    details: 内置 CWE Top 25、OWASP Top 10 (2021)、SANS Top 25 完整映射，一键判断某 CWE 是否属于高风险列表。
    link: /wellknown/cwe-top-25
    linkText: 查看知名列表 →
  - icon: 🌐
    title: MITRE REST API 客户端
    details: 内置速率限制、自动重试、结构化错误。在线获取弱点详情、父/子/祖先/后代关系与 API 版本。
    link: /sdk/api-client
    linkText: 查看 API 客户端 →
  - icon: 📥
    title: 离线 XML 目录解析
    details: 解析 MITRE 官方 XML 弱点目录，构建内存注册表与多层索引，完全离线工作，零网络依赖。
    link: /sdk/xml-parser
    linkText: 查看 XML 解析 →
  - icon: 🧭
    title: 关系导航
    details: 父/子/祖先/后代/同级/对等/链式/组合/最短路径/关系深度——比 API 更丰富的离线关系图谱查询。
    link: /sdk/navigator
    linkText: 查看导航器 →
  - icon: 🌳
    title: 层次树构建
    details: 从注册表构建 CWE 层次树/森林/视图树，支持 DFS/BFS 遍历、路径查找、叶子枚举、深度统计。
    link: /sdk/tree
    linkText: 查看树构建 →
  - icon: 🔍
    title: 搜索与过滤
    details: 按关键字、抽象层级、状态、可能性、后果范围、结构多条件过滤；按 ID/名称排序、分组与去重。
    link: /sdk/search
    linkText: 查看搜索过滤 →
  - icon: 📦
    title: 序列化与互操作
    details: JSON / XML / CSV 三格式导入导出，safeCWE 安全序列化模型，与漏洞管理平台无缝对接。
    link: /sdk/serializer
    linkText: 查看序列化 →
  - icon: 💻
    title: 40+ CLI 子命令
    details: 基于 cobra 的命令行工具，所有命令支持 text 与 JSON 双格式输出，便于脚本与流水线集成。
    link: /cli/overview
    linkText: 查看 CLI →
  - icon: 🦾
    title: AI Skills 接入
    details: 渐进式技能文档，让 Claude/GPT 等 AI 代理直接调用 cwe CLI，无需编写集成代码。
    link: /skills/
    linkText: 查看 Skills →
  - icon: ⚡
    title: 零依赖核心
    details: 核心 SDK 仅使用 Go 标准库，无第三方依赖，编译产物小、审计友好、可静态链接。
    link: /guide/performance
    linkText: 了解性能 →
  - icon: 🖥️
    title: 30+ 平台预编译
    details: Linux/macOS/Windows/BSD/AIX/Illumos/Solaris 全架构预编译二进制，覆盖 Homebrew、Scoop、deb/rpm/apk。
    link: /guide/installation
    linkText: 查看安装 →
---

## 🦾 一键复制 AI Agent 提示词

CWE Skills 是 **AI 原生**的 —— 把下面的提示词复制粘贴到你的 Codex / Claude / GPT 等 AI 客户端，它就会按指引安装并使用 `cwe` 工具，无需你手动查阅文档。

<CopyPrompt />

::: tip 💡 怎么用
1. 点击上方「📋 一键复制」按钮
2. 打开你的 AI Agent 客户端（如 OpenAI Codex、Claude Code、Cursor 等）
3. 将提示词粘贴到系统提示词 / 自定义指令 / 项目规则中
4. AI 即可根据提示词自主安装 `cwe` CLI 并调用它完成 CWE 相关任务
:::

::: info 🔗 更多接入方式
- [Skills 技能文档](/skills/) — 12 篇渐进式技能
- [快速开始指南](/guide/quick-start) — 5 分钟上手
- [四种接入方式对比](/guide/integrations)
:::

