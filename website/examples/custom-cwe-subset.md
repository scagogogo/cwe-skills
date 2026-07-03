---
title: 实战 — 导出自定义 CWE 子集
outline: [2, 3]
---

# 📤 实战 — 导出自定义 CWE 子集

按多条件过滤出关心的 CWE 子集，去重后导出为 CSV 与 JSON，供下游工具或报告使用。

<Badge type="tip" text="SDK 实战"/>
<Badge type="info" text="离线"/>

---

## 🎬 场景

团队只关心「Base 级别 + Stable 状态」的弱点，要导出一份精简清单给漏洞管理平台，CSV 给人看、JSON 给机器用。

---

## 📋 前置准备

```bash
curl -O https://cwe.mitre.org/data/xml/cwec_latest.xml.zip
unzip cwec_latest.xml.zip

go get github.com/scagogogo/cwe-skills
```

---

## 💻 完整代码

```go
package main

import (
    "fmt"
    "os"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    // 1. 加载 XML 目录
    registry, err := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
    if err != nil {
        panic(err)
    }
    registry.BuildIndexes()

    // 2. 多条件过滤：Base + Stable
    filtered := cweskills.Filter(registry.GetAllCWEs(), cweskills.FilterOption{
        Abstraction: cweskills.AbstractionBase,
        Status:      cweskills.StatusStable,
    })
    fmt.Printf("过滤后: %d 条\n", len(filtered))

    // 3. 去重 + 排序
    filtered = cweskills.Deduplicate(filtered)
    cweskills.SortByID(filtered)
    fmt.Printf("去重后: %d 条\n", len(filtered))

    // 4. 仅保留关心的字段，导出 CSV
    csvData, err := cweskills.MarshalCSV(filtered)
    if err != nil {
        panic(err)
    }
    if err := os.WriteFile("cwe_subset.csv", csvData, 0644); err != nil {
        panic(err)
    }

    // 5. 导出 JSON（用列表序列化，含完整字段）
    jsonData, err := cweskills.MarshalJSONList(filtered)
    if err != nil {
        panic(err)
    }
    if err := os.WriteFile("cwe_subset.json", jsonData, 0644); err != nil {
        panic(err)
    }

    fmt.Println("已导出: cwe_subset.csv, cwe_subset.json")
}
```

---

## ▶️ 运行步骤

```bash
go run main.go

# 检查结果
head -5 cwe_subset.csv
cat cwe_subset.json | jq 'length'
```

---

## 📤 输出示例

`cwe_subset.csv`：

```csv
ID,Name,Abstraction,Status,Structure,Description,CWEType
79,Cross-site Scripting,Base,Stable,Simple,The product does not neutralize...,weakness
89,SQL Injection,Base,Stable,Simple,The product constructs...,weakness
352,Cross-Site Request Forgery,Base,Stable,Simple,The product does not...,weakness
...
```

终端摘要：

```text
过滤后: 312 条
去重后: 312 条
已导出: cwe_subset.csv, cwe_subset.json
```

---

## 🧩 扩展思路

- **加关键字**：`FilterOption` 加 `Keyword: "Injection"`，只导出注入类。
- **加 ID 范围**：用 `MinID`/`MaxID` 限定 ID 区间，如只导出 1–1000。
- **按后果过滤**：`FilterOption` 加 `Scope: ScopeConfidentiality`，导出影响机密性的弱点。
- **注册表级导出**：用 `registry.ExportJSON()` 一次性导出全部条目类型（弱点+类别+视图），数据更全但文件更大。
- **增量更新**：导出后记录 CWE ID 集合的哈希，下次对比发现新增条目。

::: warning CSV 有损
`MarshalCSV` 只保留基础字段。需要完整数据（关系、后果、缓解措施等）请用 `MarshalJSONList` 或 `registry.ExportJSON()`。
:::

---

## 📖 相关文档

- [技能 08 — 搜索与过滤](../skills/08-local-search-filter) · [技能 12 — 序列化](../skills/12-sdk-serialization)
- [SDK: Filter](../sdk/filter) · [FilterOption](../sdk/filter-option) · [Deduplicate](../sdk/deduplicate) · [MarshalCSV](../sdk/marshal-csv)
- [返回示例总览](./)
