---
title: OWASP Top 10 (2021)
outline: [2, 3]
---

# 📊 OWASP Top 10 (2021) CWE 映射

OWASP 发布的 Web 应用安全风险 Top 10（2021 版），每个风险类别映射到一组相关 CWE ID。SDK 中以 `OWASPTop10`（`map[string][]int`）提供，配套 `IsInOWASPTop10`、`GetOWASPCategory`、`GetOWASPCategories` 三组函数。

## 📊 背景来源

由 OWASP（开放式 Web 应用程序安全项目）编制，针对 **Web 应用**安全风险：

- 基于社区投票与行业数据调研
- 每 3–4 年更新一次，2021 版是当前最新
- 每个风险类别关联多个 CWE ID，体现「一个风险类别对应多种具体弱点」

是业界合规检查、渗透测试、安全培训的事实标准。

## 📋 完整映射

| 类别 | CWE 数量 | CWE ID 列表 |
| --- | ---: | --- |
| A01:2021-Broken Access Control | 23 | 22, 23, 35, 59, 78, 94, 200, 201, 219, 255, 269, 276, 284, 285, 287, 306, 346, 639, 651, 668, 862, 863, 922 |
| A02:2021-Cryptographic Failures | 22 | 260, 261, 295, 310, 311, 312, 319, 325, 326, 327, 328, 329, 330, 337, 338, 340, 347, 522, 757, 759, 760, 780 |
| A03:2021-Injection | 47 | 20, 74, 75, 77, 78, 79, 80, 83, 87, 88, 89, 90, 91, 94, 95, 96, 97, 98, 99, 100, 113, 116, 138, 141, 147, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171 |
| A04:2021-Insecure Design | 13 | 209, 235, 256, 267, 284, 285, 287, 311, 326, 384, 393, 664, 863 |
| A05:2021-Security Misconfiguration | 22 | 2, 5, 11, 13, 15, 16, 260, 315, 520, 526, 537, 540, 544, 546, 547, 548, 611, 613, 614, 759, 760, 1021 |
| A06:2021-Vulnerable and Outdated Components | 0 | _（空，无 CWE 映射）_ |
| A07:2021-Identification and Authentication Failures | 22 | 255, 256, 258, 259, 260, 287, 288, 290, 294, 295, 297, 306, 307, 346, 384, 521, 522, 523, 613, 620, 640, 798 |
| A08:2021-Software and Data Integrity Failures | 12 | 311, 345, 353, 426, 494, 502, 565, 610, 653, 754, 829, 912 |
| A09:2021-Security Logging and Monitoring Failures | 4 | 117, 223, 532, 778 |
| A10:2021-Server-Side Request Forgery (SSRF) | 2 | 918, 1021 |

::: warning A06 无 CWE 映射
`A06:2021-Vulnerable and Outdated Components`（易受攻击和过时的组件）关注的是**依赖管理**问题，OWASP 未将其映射到具体 CWE，因此 `OWASPTop10` 中该类别对应空切片。`IsInOWASPTop10` 永远不会因 A06 返回 true。
:::

::: details 一个 CWE 可属于多个类别
例如 CWE-287（认证不当）同时出现在 A01（访问控制失效）与 A07（身份认证失效）。此时：
- `GetOWASPCategory(287)` 只返回首个匹配（A01）
- `GetOWASPCategories(287)` 返回全部 `[A01:2021-Broken Access Control, A07:2021-Identification and Authentication Failures]`
:::

## ✅ SDK 用法示例

```go
package main

import (
    "fmt"
    cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
    // 1. 判断是否属于某 OWASP 类别
    fmt.Println(cweskills.IsInOWASPTop10(89)) // true

    // 2. 取首个匹配类别
    fmt.Println(cweskills.GetOWASPCategory(89))
    // A03:2021-Injection

    // 3. 取全部匹配类别（CWE-287 跨两个类别）
    for _, c := range cweskills.GetOWASPCategories(287) {
        fmt.Println(c)
    }
    // A01:2021-Broken Access Control
    // A07:2021-Identification and Authentication Failures

    // 4. 遍历整个映射
    for category, ids := range cweskills.OWASPTop10 {
        fmt.Printf("%s -> %d 个 CWE\n", category, len(ids))
    }
}
```

## 🖥️ CLI 用法

```bash
# 按类别列出（text）
cwe wellknown owasp
```

```text
OWASP Top 10 (2021):

  A01:2021-Broken Access Control:
    CWE-22
    CWE-35
  ...
```

```bash
# JSON 输出
cwe wellknown owasp -o json | jq '.[] | {category, count: (.cwe_ids|length)}'

# 查某 CWE 属于哪个 OWASP 类别
cwe wellknown check CWE-89
```

```text
CWE-89: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
```

## 🎯 使用场景

<Badge type="tip" text="合规对照" /> 将扫描/审计发现映射到 OWASP Top 10，生成合规报告
<Badge type="info" text="修复优先级" /> 按类别聚合弱点，从风险类别维度规划修复批次
<Badge type="warning" text="测试用例" /> 以每个类别为单元设计安全测试用例套件
<Badge type="info" text="培训大纲" /> OWASP Top 10 是 Web 安全培训的经典框架

## 🔗 相关链接

- [知名列表总览](./overview)
- [CWE Top 25 (2024)](./cwe-top-25)
- [SANS Top 25](./sans-top-25)
- SDK：[OWASPTop10 映射](../sdk/owasp-top-10)、[IsInOWASPTop10](../sdk/is-in-owasp-top-10)、[GetOWASPCategory](../sdk/get-owasp-category)、[GetOWASPCategories](../sdk/get-owasp-categories)
- CLI：[wellknown owasp](../cli/wellknown-owasp)、[wellknown check](../cli/wellknown-check)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
