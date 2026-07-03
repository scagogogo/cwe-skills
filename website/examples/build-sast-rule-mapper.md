---
title: 实战 — SAST 规则映射到 CWE
outline: [2, 3]
---

# 🛠️ 实战 — SAST 规则映射到 CWE

把 SAST（静态应用安全测试）工具的规则输出批量映射到 CWE，生成带弱点详情的对照报告。

<Badge type="tip" text="SDK 实战"/>
<Badge type="info" text="在线 API"/>

---

## 🎬 场景

你有一个 SAST 工具，每条规则带名称、描述、可能引用的 CWE。你想生成一份报告：每条规则 → 对应 CWE → CWE 详情（名称、抽象层级、是否在 Top 25）。

---

## 📋 前置准备

```bash
go get github.com/scagogogo/cwe-skills
```

需要网络访问 MITRE API。

---

## 💻 完整代码

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "github.com/scagogogo/cwe-skills"
)

// SAST 规则（示例结构）
type SASTRule struct {
    ID       string `json:"rule_id"`
    Name     string `json:"name"`
    Severity string `json:"severity"`
    Desc     string `json:"description"` // 可能包含 CWE-NNN 引用
}

// 映射结果
type Mapping struct {
    RuleID    string `json:"rule_id"`
    CWEIDs    []string `json:"cwe_ids"`
    InTop25   []string `json:"in_top25"`
    Weaknesses []WeakSummary `json:"weaknesses"`
}

type WeakSummary struct {
    ID          string `json:"cwe_id"`
    Name        string `json:"name"`
    Abstraction string `json:"abstraction"`
}

func main() {
    // 1. 读取 SAST 规则输出
    data, _ := os.ReadFile("sast_rules.json")
    var rules []SASTRule
    _ = json.Unmarshal(data, &rules)

    client := cweskills.NewAPIClient(cweskills.WithAPIRateLimit(5, 0))
    defer client.Close()
    ctx := context.Background()

    var results []Mapping
    for _, r := range rules {
        m := Mapping{RuleID: r.ID}

        // 2. 从规则描述提取 CWE ID
        text := r.Name + " " + r.Desc
        m.CWEIDs = cweskills.ExtractCWEIDs(text)

        // 3. 逐个取详情 + 判断 Top 25
        for _, cweStr := range m.CWEIDs {
            id, err := cweskills.ParseCWEID(cweStr)
            if err != nil {
                continue
            }
            if cweskills.IsInTop25(id) {
                m.InTop25 = append(m.InTop25, cweStr)
            }
            if w, err := client.GetWeakness(ctx, id); err == nil {
                m.Weaknesses = append(m.Weaknesses, WeakSummary{
                    ID:          cweStr,
                    Name:        w.Name,
                    Abstraction: w.Abstraction.String(),
                })
            }
        }
        results = append(results, m)
    }

    // 4. 输出报告
    out, _ := json.MarshalIndent(results, "", "  ")
    _ = os.WriteFile("sast_cwe_mapping.json", out, 0644)

    // 打印摘要
    for _, m := range results {
        top := strings.Join(m.InTop25, ",")
        if top == "" {
            top = "—"
        }
        fmt.Printf("%s → %d 个CWE，Top25: %s\n", m.RuleID, len(m.CWEIDs), top)
    }
}
```

---

## ▶️ 运行步骤

```bash
# 1. 准备 SAST 规则 JSON（含 description 字段引用 CWE）
cat > sast_rules.json <<'EOF'
[
  {"rule_id":"R001","name":"SQL Injection Sink","severity":"high","description":"Detects CWE-89 patterns in SQL queries"},
  {"rule_id":"R002","name":"XSS Reflected","severity":"high","description":"Reflected XSS, see CWE-79"},
  {"rule_id":"R003","name":"Hardcoded Password","severity":"medium","description":"No CWE reference"}
]
EOF

# 2. 运行
go run main.go

# 3. 查看结果
cat sast_cwe_mapping.json | jq '.[] | {rule_id, cwe_ids, in_top25}'
```

---

## 📤 输出示例

```text
R001 → 1 个CWE，Top25: CWE-89
R002 → 1 个CWE，Top25: CWE-79
R003 → 0 个CWE，Top25: —
```

```json
[
  {
    "rule_id": "R001",
    "cwe_ids": ["CWE-89"],
    "in_top25": ["CWE-89"],
    "weaknesses": [{
      "cwe_id": "CWE-89",
      "name": "Improper Neutralization of Special Elements used in an SQL Command...",
      "abstraction": "Base"
    }]
  }
]
```

---

## 🧩 扩展思路

- **缺失 CWE 的规则**：对没有 CWE 引用的规则，用 `cwe search --keyword <规则名>` 离线匹配候选 CWE。
- **批量加速**：用 `WithAPIRateLimit` 控制速率，避免被 MITRE 限流；规则量大时分批处理。
- **映射到 OWASP**：把 `IsInTop25` 换成 `GetOWASPCategories`，生成 OWASP 合规维度的报告。
- **持久化**：把结果写入数据库，作为后续漏洞管理的基础。

---

## 📖 相关文档

- [技能 02 — 提取](../skills/02-cwe-id-extraction-comparison) · [技能 03 — Top 25](../skills/03-well-known-lists) · [技能 05 — GetWeakness](../skills/05-api-show-weakness)
- [SDK: ExtractCWEIDs](../sdk/extract-cwe-ids) · [GetWeakness](../sdk/api-get-weakness) · [IsInTop25](../sdk/is-in-top-25)
- [返回示例总览](./)
