---
title: 实战 — OWASP 合规检查
outline: [2, 3]
---

# ✅ 实战 — OWASP 合规检查

输入一组 CWE ID，映射到 OWASP Top 10 (2021) 类别，生成合规覆盖报告：哪些 OWASP 类别被覆盖、哪些缺失。

<Badge type="tip" text="SDK 实战"/>
<Badge type="info" text="内置列表"/>

---

## 🎬 场景

合规审计要求证明：你的安全测试覆盖了 OWASP Top 10 的每个类别。给你一份已发现 CWE 列表，生成「覆盖了哪几项 OWASP、缺哪几项」的报告。

---

## 📋 前置准备

```bash
go get github.com/scagogogo/cwe-skills
```

无需网络（OWASP 映射是内置数据）。

---

## 💻 完整代码

```go
package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    // 1. 输入：已发现的 CWE ID 列表
    if len(os.Args) < 2 {
        fmt.Println("用法: owasp-check CWE-79 CWE-89 CWE-352 ...")
        os.Exit(1)
    }
    inputs := os.Args[1:]

    // 2. 把每个 CWE 映射到 OWASP 类别
    covered := map[string][]string{} // category -> []cweID
    for _, in := range inputs {
        id, err := cweskills.ParseCWEID(in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "跳过无效输入 %s: %v\n", in, err)
            continue
        }
        cats := cweskills.GetOWASPCategories(id)
        for _, c := range cats {
            covered[c] = append(covered[c], cweskills.FormatCWEIDFromInt(id))
        }
    }

    // 3. 全部 OWASP Top 10 类别（从内置映射取键）
    allCats := []string{}
    for k := range cweskills.OWASPTop10 {
        allCats = append(allCats, k)
    }

    // 4. 输出报告
    fmt.Println("=== OWASP Top 10 合规覆盖报告 ===")
    fmt.Println()
    for _, cat := range allCats {
        cwes, ok := covered[cat]
        if ok {
            fmt.Printf("✅ %s  覆盖(%d): %s\n", cat, len(cwes), strings.Join(cwes, ", "))
        } else {
            fmt.Printf("❌ %s  未覆盖\n", cat)
        }
    }

    // 5. 覆盖率
    coveredCount := 0
    for _, cat := range allCats {
        if _, ok := covered[cat]; ok {
            coveredCount++
        }
    }
    fmt.Printf("\n覆盖率: %d/%d (%.0f%%)\n", coveredCount, len(allCats),
        float64(coveredCount)/float64(len(allCats))*100)
}
```

---

## ▶️ 运行步骤

```bash
# 输入一组扫描发现的 CWE
go run main.go CWE-79 CWE-89 CWE-352 CWE-862 CWE-287

# 若想看完整 OWASP 映射做参照
cwe wellknown owasp
```

---

## 📤 输出示例

```text
=== OWASP Top 10 合规覆盖报告 ===

✅ A03:2021-Injection  覆盖(2): CWE-79, CWE-89
✅ A01:2021-Broken Access Control  覆盖(1): CWE-862
✅ A07:2021-Identification and Authentication Failures  覆盖(1): CWE-287
❌ A02:2021-Cryptographic Failures  未覆盖
❌ A04:2021-Insecure Design  未覆盖
❌ A05:2021-Security Misconfiguration  未覆盖
❌ A06:2021-Vulnerable and Outdated Components  未覆盖
...
覆盖率: 3/10 (30%)
```

---

## 🧩 扩展思路

- **导出报告**：把结果序列化为 CSV/HTML，作为合规审计交付物。
- **缺口补测**：对 `❌ 未覆盖` 的类别，用 `cwe wellknown owasp -o json | jq` 取该类别的 CWE 清单，指导补测方向。
- **多轮迭代**：每轮测试后重新跑报告，跟踪覆盖率提升趋势。
- **结合 Top 25**：同时用 `IsInTop25` 标注高优先级，输出「覆盖 + 优先级」双维度报告。

::: tip 内置数据无需联网
`OWASPTop10` 映射是 SDK 内置常量，合规检查完全离线，适合内网审计环境。
:::

---

## 📖 相关文档

- [技能 03 — 知名列表](../skills/03-well-known-lists) · [技能 12 — 序列化](../skills/12-sdk-serialization)
- [SDK: GetOWASPCategories](../sdk/get-owasp-categories) · [OWASPTop10](../sdk/owasp-top-10)
- [CLI: wellknown owasp](../cli/wellknown-owasp)
- [返回示例总览](./)
