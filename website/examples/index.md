---
title: 示例与教程
outline: [2, 3]
---

# 🗺️ 示例与教程

端到端实战教程，把 CWE Skills 的 CLI 与 SDK 串起来解决真实安全问题。每个示例包含完整代码/脚本、运行步骤与输出示例。

<Badge type="tip" text="端到端实战"/>

---

## 📚 教程列表

| # | 教程 | 场景 | 主要能力 |
|---|------|------|----------|
| 1 | [SAST 规则映射到 CWE](./build-sast-rule-mapper) | 把 SAST 工具规则批量映射到 CWE | `ExtractCWEIDs` · `GetWeakness` |
| 2 | [离线 CWE 浏览器](./offline-cwe-explorer) | 用 XML 目录构建离线浏览器 | `NewXMLParser` · `BuildTree` · `nav` |
| 3 | [漏洞分诊](./vulnerability-triage) | 输入漏洞描述自动定优先级 | `extract` · `wellknown` · `show` |
| 4 | [CI 流水线集成](./ci-pipeline-integration) | GitHub Actions 检查 Top 25 引入 | `cwe` + JSON + jq |
| 5 | [OWASP 合规检查](./owasp-compliance-check) | CWE 列表生成 OWASP 合规报告 | `GetOWASPCategories` |
| 6 | [导出自定义 CWE 子集](./custom-cwe-subset) | 过滤去重导出 CSV/JSON | `filter` · `Deduplicate` · `ExportCSV` |
| 7 | [树可视化](./tree-visualization) | 生成 mermaid/ASCII 树图 | `BuildForest` · `Walk` |

::: tip 先读 Skills
这些教程基于 [Skills](../skills/) 与 [SDK](../sdk/overview) 能力。若不熟悉基础命令，建议先过一遍 12 篇技能文档。
:::

---

## 🚀 快速选择

- **想做合规报告** → [OWASP 合规检查](./owasp-compliance-check)
- **想给扫描结果排优先级** → [漏洞分诊](./vulnerability-triage)
- **想卡 CI 防止引入 Top 25** → [CI 流水线集成](./ci-pipeline-integration)
- **想离线浏览 CWE** → [离线 CWE 浏览器](./offline-cwe-explorer)
- **想把 SAST 规则对齐 CWE** → [SAST 规则映射](./build-sast-rule-mapper)
- **想导出筛选后的 CWE** → [自定义 CWE 子集](./custom-cwe-subset)
- **想画 CWE 层次图** → [树可视化](./tree-visualization)

---

## 📦 前置准备

多数教程需要 `cwe` CLI 已安装：

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/
cwe version
```

离线教程还需 MITRE XML 目录：

```bash
# 从 https://cwe.mitre.org/data/xml.html 下载
curl -O https://cwe.mitre.org/data/xml/cwec_latest.xml.zip
unzip cwec_latest.xml.zip
```

SDK 教程需安装 Go 包：

```bash
go get github.com/scagogogo/cwe-skills
```

---

## 📖 相关文档

- [Skills 总览](../skills/)
- [CLI 命令参考](../cli/overview)
- [Go SDK API](../sdk/overview)
- [安装指南](../guide/installation)
