---
title: 知名列表总览
outline: [2, 3]
---

# 🏆 知名列表总览

CWE 体系庞大（上千条目），业界权威机构从中提炼出若干「最危险」清单，用于聚焦修复资源。`cweskills` 在 `wellknown_ids.go` 中内置了三份主流清单与五类知名视图，SDK 与 CLI 均可直接查询，**无需联网或加载 XML 文件**。

## 📊 内置清单一览

| 清单 | 数据来源 | 类型 | 对应变量 | 项数 |
| --- | --- | --- | --- | --- |
| CWE Top 25 (2024) | MITRE | `[]int` | `CWETop25` | 25 |
| OWASP Top 10 (2021) | OWASP | `map[string][]int` | `OWASPTop10` | 10 类别 |
| SANS Top 25 | SANS Institute + MITRE | `[]int` | `SANSTop25` | 25 |

<Badge type="tip" text="离线可用" /> 上述清单为编译期内置常量，调用零延迟、零网络依赖。

## 📋 知名视图常量

除清单外，SDK 还内置 5 个 MITRE 官方视图 ID 常量，便于按视图维度组织数据：

| 常量 | 值 | 视图名称 |
| --- | --- | --- |
| `CWEViewResearchConcepts` | 1000 | 研究概念视图 |
| `CWEViewDevelopmentConcepts` | 699 | 软件开发视图 |
| `CWEViewHardwareDesign` | 1199 | 硬件设计视图 |
| `CWEViewCWECrossSection` | 888 | CWE 横截面视图 |
| `CWEViewComprehensiveDictionary` | 1400 | 综合 CWE 字典 |

详见 [知名视图](./well-known-views)。

## 🔍 各清单背景来源

### MITRE CWE Top 25

由 MITRE 基于 NVD（美国国家漏洞数据库）中已披露漏洞的**出现频率**与**CVSS 评分**计算得出，按危险程度排序，每年更新一次。代表当前对软件**最严重的安全威胁**。

### OWASP Top 10

由 OWASP（开放式 Web 应用程序安全项目）编制，针对 **Web 应用**安全风险，是业界合规与培训的事实标准。2021 版将每个风险类别映射到一组相关 CWE ID。

### SANS Top 25

由 SANS Institute 与 MITRE 合作编制的较早清单，侧重于**可被攻击者利用以获取系统控制权**的编程错误，按错误类型分组（不安全交互、高风险资源、防守性边界）。

## ✅ SDK 判断函数

每个清单都提供 `IsIn*` 判断函数，OWASP 额外提供类别查询：

```go
package main

import (
    "fmt"
    cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
    fmt.Println(cweskills.IsInTop25(79))        // true（XSS）
    fmt.Println(cweskills.IsInOWASPTop10(89))   // true（SQL 注入）
    fmt.Println(cweskills.IsInSANSTop25(78))    // true（OS 命令注入）
    fmt.Println(cweskills.IsInWellKnownView(699)) // true（软件开发视图）

    // OWASP 类别查询
    fmt.Println(cweskills.GetOWASPCategory(89))
    // A03:2021-Injection
    fmt.Println(cweskills.GetOWASPCategories(287))
    // [A07:2021-Identification and Authentication Failures]
}
```

::: tip 函数速查
- `IsInTop25(cweID int) bool`
- `IsInOWASPTop10(cweID int) bool`
- `IsInSANSTop25(cweID int) bool`
- `GetOWASPCategory(cweID int) string` — 返回首个匹配类别，无则空串
- `GetOWASPCategories(cweID int) []string` — 返回所有匹配类别
- `IsInWellKnownView(viewID int) bool`
:::

## 🖥️ CLI 用法

`cwe wellknown` 是父命令，下设四个子命令：

```bash
# 列出三份清单
cwe wellknown top25
cwe wellknown owasp
cwe wellknown sans

# 批量检查某个 CWE 在哪些清单上榜
cwe wellknown check CWE-79 89 352
```

```text
CWE-79: [Top 25 OWASP Top 10 (A03:2021-Injection)]
CWE-89: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
CWE-352: [Top 25 OWASP Top 10 (A01:2021-Broken Access Control)]
```

详见 [CLI wellknown 总览](../cli/wellknown)。

## 🎯 使用场景

<Badge type="tip" text="修复优先级" /> 扫描器/漏洞库结果中标注「上榜」弱点，优先分配修复资源
<Badge type="info" text="合规对照" /> 将自身发现映射到 OWASP Top 10，生成合规报告
<Badge type="warning" text="安全培训" /> 以最危险清单为大纲组织开发安全培训
<Badge type="info" text="风险评估" /> 评估代码库/产品对高危弱点的暴露面

::: warning 时效性
清单反映其发布年度的威胁态势。`CWETop25` 为 2024 版，`OWASPTop10` 为 2021 版，`SANSTop25` 为较早版本。建议以较新的 `top25` 为基线，其余作为补充参考。
:::

## 🔗 相关链接

- [CWE Top 25 (2024)](./cwe-top-25)
- [OWASP Top 10 (2021)](./owasp-top-10)
- [SANS Top 25](./sans-top-25)
- [知名视图](./well-known-views)
- SDK API：[wellknown-ids 概览](../sdk/wellknown-ids)
- CLI：[wellknown 总览](../cli/wellknown)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
