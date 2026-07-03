---
title: 技能 08 — 本地搜索与过滤
outline: [2, 3]
---

# 🔎 技能 08 — 本地搜索与过滤

从离线 MITRE XML 目录中搜索与多条件过滤 CWE 条目。完全离线，无需网络。

<Badge type="tip" text="离线"/>
<Badge type="info" text="需 XML 目录"/>

XML 目录从 [MITRE](https://cwe.mitre.org/data/xml.html) 下载，如 `cwec_latest.xml`。

---

## 🎯 技能目标

- 按关键字 / 抽象层级 / 状态 / 可能性 / 结构 / 后果范围搜索
- 多条件过滤（AND 逻辑）
- 排序、分组、去重、统计

---

## 💻 CLI 命令

### search — 单条件搜索

```bash
cwe search --xml cwec_latest.xml --keyword Injection
```

| Flag | 简写 | 说明 |
|------|------|------|
| `--xml` | `-x` | **（必填）** XML 目录路径 |
| `--keyword` | `-k` | 按关键字搜（名称+描述） |
| `--abstraction` | `-a` | 按抽象层级 |
| `--status` | `-s` | 按状态 |
| `--likelihood` | `-l` | 按利用可能性 |
| `--structure` | `-t` | 按结构 |
| `--scope` | | 按后果范围 |
| `--top-level` | | 仅柱状（顶层）弱点 |
| `--base-weaknesses` | | 仅 Base 弱点 |
| `--chains` | | 仅链式弱点 |
| `--composites` | | 仅组合弱点 |
| `--sort` | | 排序：id/name/abstraction |
| `--group-by` | | 分组：abstraction/status/likelihood |
| `--dedup` | | 去重 |

### filter — 多条件过滤（AND）

```bash
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable --keyword Injection
cwe filter --xml cwec_latest.xml --likelihood High --scope Confidentiality --sort name
```

### stats — 统计

```bash
cwe stats --xml cwec_latest.xml
```

---

## 🔧 SDK API

### 单条件查找

```go
results := cweskills.FindByKeyword(registry, "Injection")
results := cweskills.FindByAbstraction(registry, cweskills.AbstractionBase)
results := cweskills.FindByConsequenceScope(registry, cweskills.ScopeConfidentiality)
results := cweskills.FindTopLevel(registry)
results := cweskills.FindBaseWeaknesses(registry)
```

### 多条件过滤

```go
filtered := cweskills.Filter(all, cweskills.FilterOption{
    Abstraction: cweskills.AbstractionBase,
    Status:      cweskills.StatusStable,
    Keyword:     "Injection",
})
```

::: details FilterOption 字段
`Abstraction` · `Status` · `Structure` · `Likelihood` · `MinID` · `MaxID` · `Keyword` · `Scope`。所有非零值条件为 AND 关系。
:::

### 排序 / 分组 / 去重

```go
cweskills.SortByID(results)
cweskills.SortByName(results)
cweskills.GroupByAbstraction(results)
cweskills.Deduplicate(results)
```

---

## 📝 示例

### 命令行

```bash
# 找出所有 Base+Stable+高可能性的注入类弱点
cwe filter --xml cwec_latest.xml \
  --abstraction Base --status Stable --likelihood High \
  --keyword Injection -o json | jq 'length'

# 按抽象层级分组统计
cwe stats --xml cwec_latest.xml -o json
```

### Go

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    registry, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
    registry.BuildIndexes()

    filtered := cweskills.Filter(registry.GetAllCWEs(), cweskills.FilterOption{
        Abstraction: cweskills.AbstractionBase,
        Keyword:     "Injection",
    })
    fmt.Printf("匹配 %d 条\n", len(filtered))
}
```

---

## 🤖 AI 代理使用提示

- 用户描述「找某种类型的弱点」时，AI 用 `cwe filter` 组合条件。
- 离线搜索不受速率限制，AI 可放心连续调用。
- 让 AI 加 `-o json`，结果便于解析与计数。

::: warning 需要 XML 目录
所有 `search`/`filter`/`stats` 命令必须 `--xml <file>` 指定本地 XML 目录。
:::

---

## 📖 相关文档

- [技能 09 — 本地注册表](./09-local-registry)
- [CLI: search](../cli/search) · [filter](../cli/filter) · [stats](../cli/stats)
- [SDK: FindByKeyword](../sdk/find-by-keyword) · [Filter](../sdk/filter) · [FilterOption](../sdk/filter-option)
- [返回 Skills 总览](./)
