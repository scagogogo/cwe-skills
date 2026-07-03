---
title: SANS Top 25
outline: [2, 3]
---

# 🔍 SANS Top 25 Most Dangerous Software Errors

SANS Institute 与 MITRE 合作编制的「最危险软件错误 Top 25」，侧重于可被攻击者利用以获取系统控制权的编程错误。SDK 中以 `SANSTop25` 变量提供，配合 `IsInSANSTop25` 函数判断归属。

## 📊 背景来源

由 SANS Institute 联合 MITRE 编制：

- 关注**可被攻击者利用来获取系统控制权**的编程错误
- 按「不安全交互」「高风险资源管理」「防守性边界」三大类组织
- 是较早期清单，部分条目与现行 CWE Top 25 存在差异

相比 MITRE 的 Top 25（基于 NVD 频率+CVSS），SANS 版本更偏重**编程错误本身的可利用性**，而非年度统计数据。

## 📋 完整列表

| 排名 | CWE ID | 名称 | SDK 函数 |
| ---: | --- | --- | --- |
| 1 | CWE-89 | SQL Injection | `IsInSANSTop25(89)` |
| 2 | CWE-78 | OS Command Injection | `IsInSANSTop25(78)` |
| 3 | CWE-79 | Cross-site Scripting (XSS) | `IsInSANSTop25(79)` |
| 4 | CWE-20 | Improper Input Validation | `IsInSANSTop25(20)` |
| 5 | CWE-22 | Path Traversal | `IsInSANSTop25(22)` |
| 6 | CWE-352 | Cross-Site Request Forgery (CSRF) | `IsInSANSTop25(352)` |
| 7 | CWE-416 | Use After Free | `IsInSANSTop25(416)` |
| 8 | CWE-787 | Out-of-bounds Write | `IsInSANSTop25(787)` |
| 9 | CWE-125 | Out-of-bounds Read | `IsInSANSTop25(125)` |
| 10 | CWE-94 | Code Injection | `IsInSANSTop25(94)` |
| 11 | CWE-190 | Integer Overflow or Wraparound | `IsInSANSTop25(190)` |
| 12 | CWE-434 | Unrestricted Upload of File with Dangerous Type | `IsInSANSTop25(434)` |
| 13 | CWE-862 | Missing Authorization | `IsInSANSTop25(862)` |
| 14 | CWE-287 | Improper Authentication | `IsInSANSTop25(287)` |
| 15 | CWE-306 | Missing Authentication for Critical Function | `IsInSANSTop25(306)` |
| 16 | CWE-863 | Incorrect Authorization | `IsInSANSTop25(863)` |
| 17 | CWE-798 | Use of Hard-coded Credentials | `IsInSANSTop25(798)` |
| 18 | CWE-502 | Deserialization of Untrusted Data | `IsInSANSTop25(502)` |
| 19 | CWE-77 | Command Injection | `IsInSANSTop25(77)` |
| 20 | CWE-119 | Improper Restriction of Operations within the Bounds of a Memory Buffer | `IsInSANSTop25(119)` |
| 21 | CWE-639 | Authorization Bypass Through User-Controlled Key | `IsInSANSTop25(639)` |
| 22 | CWE-770 | Allocation of Resources Without Limits or Throttling | `IsInSANSTop25(770)` |
| 23 | CWE-918 | Server-Side Request Forgery (SSRF) | `IsInSANSTop25(918)` |
| 24 | CWE-476 | NULL Pointer Dereference | `IsInSANSTop25(476)` |
| 25 | CWE-200 | Exposure of Sensitive Information | `IsInSANSTop25(200)` |

<Badge type="info" text="SANS 独有条目" /> CWE-190（整数溢出）、CWE-287（认证不当）、CWE-798（硬编码凭证）、CWE-119（内存缓冲区限制不当）是 SANS Top 25 中较为突出、CWE Top 25 (2024) 未收录的条目，反映 SANS 对可利用编程错误的侧重。

## ✅ SDK 用法示例

```go
package main

import (
    "fmt"
    cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
    // 判断是否在 SANS Top 25
    fmt.Println(cweskills.IsInSANSTop25(798)) // true（硬编码凭证）
    fmt.Println(cweskills.IsInSANSTop25(999)) // false

    // 遍历完整 SANS Top 25
    for i, id := range cweskills.SANSTop25 {
        fmt.Printf("%2d. CWE-%d\n", i+1, id)
    }

    // 与 CWE Top 25 取交集，识别「双榜」高风险弱点
    for _, id := range cweskills.SANSTop25 {
        if cweskills.IsInTop25(id) {
            fmt.Printf("CWE-%d 同时上榜 CWE Top 25 与 SANS Top 25\n", id)
        }
    }
}
```

## 🖥️ CLI 用法

```bash
# 列出完整 SANS Top 25（text）
cwe wellknown sans
```

```text
SANS Top 25 Most Dangerous Software Errors (25 项):

   1. CWE-89
   2. CWE-78
   3. CWE-79
  ...
  25. CWE-200
```

```bash
# JSON 输出
cwe wellknown sans -o json | jq 'length'

# 检查某 CWE 是否在 SANS Top 25
cwe wellknown check CWE-798
```

```text
CWE-798: [OWASP Top 10 (A07:2021-Identification and Authentication Failures) SANS Top 25]
```

## 🎯 使用场景

<Badge type="tip" text="历史对照" /> 与 CWE Top 25 对比，识别在不同权威清单中均上榜的高优先级弱点
<Badge type="info" text="可利用性侧重" /> SANS 偏重可被攻击者直接利用的编程错误，适合红队/渗透视角
<Badge type="warning" text="遗留系统" /> 较早清单，适合评估遗留/存量系统的弱点暴露面
<Badge type="info" text="编程培训" /> 强调编程错误本身，适合作为代码安全培训的弱点清单

::: warning 时效性
SANS Top 25 为较早清单，部分排名与条目已不反映当前威胁态势。建议以 [CWE Top 25 (2024)](./cwe-top-25) 作为当前基线，本清单作为补充参考。
:::

## 🔗 相关链接

- [知名列表总览](./overview)
- [CWE Top 25 (2024)](./cwe-top-25)
- [OWASP Top 10 (2021)](./owasp-top-10)
- SDK：[SANSTop25 列表](../sdk/sans-top-25)、[IsInSANSTop25](../sdk/is-in-sans-top-25)
- CLI：[wellknown sans](../cli/wellknown-sans)、[wellknown check](../cli/wellknown-check)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
