---
title: CWE Top 25 (2024)
outline: [2, 3]
---

# 🏆 CWE Top 25 Most Dangerous Software Weaknesses (2024)

MITRE 发布的「最危险软件弱点 Top 25」2024 版，按危险程度从高到低排序。SDK 中以 `CWETop25` 变量提供，配合 `IsInTop25` 函数判断归属。

## 📊 背景来源

由 MITRE 基于 **NVD（美国国家漏洞数据库）** 中已披露漏洞的：

- **出现频率**（prevalence）—— 该弱点在真实漏洞中出现的频次
- **CVSS 评分**（severity）—— 漏洞的严重程度

综合计算得出，每年更新。2024 版代表当前对软件最严重的安全威胁集合。

## 📋 完整列表

| 排名 | CWE ID | 名称 | SDK 函数 |
| ---: | --- | --- | --- |
| 1 | CWE-79 | Cross-site Scripting (XSS) | `IsInTop25(79)` |
| 2 | CWE-89 | SQL Injection | `IsInTop25(89)` |
| 3 | CWE-352 | Cross-Site Request Forgery (CSRF) | `IsInTop25(352)` |
| 4 | CWE-862 | Missing Authorization | `IsInTop25(862)` |
| 5 | CWE-787 | Out-of-bounds Write | `IsInTop25(787)` |
| 6 | CWE-22 | Path Traversal | `IsInTop25(22)` |
| 7 | CWE-416 | Use After Free | `IsInTop25(416)` |
| 8 | CWE-125 | Out-of-bounds Read | `IsInTop25(125)` |
| 9 | CWE-78 | OS Command Injection | `IsInTop25(78)` |
| 10 | CWE-94 | Code Injection | `IsInTop25(94)` |
| 11 | CWE-120 | Buffer Copy without Checking Size of Input | `IsInTop25(120)` |
| 12 | CWE-434 | Unrestricted Upload of File with Dangerous Type | `IsInTop25(434)` |
| 13 | CWE-476 | NULL Pointer Dereference | `IsInTop25(476)` |
| 14 | CWE-121 | Stack-based Buffer Overflow | `IsInTop25(121)` |
| 15 | CWE-502 | Deserialization of Untrusted Data | `IsInTop25(502)` |
| 16 | CWE-122 | Heap-based Buffer Overflow | `IsInTop25(122)` |
| 17 | CWE-863 | Incorrect Authorization | `IsInTop25(863)` |
| 18 | CWE-20 | Improper Input Validation | `IsInTop25(20)` |
| 19 | CWE-284 | Improper Access Control | `IsInTop25(284)` |
| 20 | CWE-200 | Exposure of Sensitive Information | `IsInTop25(200)` |
| 21 | CWE-306 | Missing Authentication for Critical Function | `IsInTop25(306)` |
| 22 | CWE-918 | Server-Side Request Forgery (SSRF) | `IsInTop25(918)` |
| 23 | CWE-77 | Command Injection | `IsInTop25(77)` |
| 24 | CWE-639 | Authorization Bypass Through User-Controlled Key | `IsInTop25(639)` |
| 25 | CWE-770 | Allocation of Resources Without Limits or Throttling | `IsInTop25(770)` |

<Badge type="danger" text="排名即危险程度" /> 排名越靠前，综合威胁越大。前 5 名（XSS、SQL 注入、CSRF、缺少授权、越界写）覆盖了 Web 与内存安全两大主战场。

## ✅ SDK 用法示例

```go
package main

import (
    "fmt"
    cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
    // 判断单个 CWE 是否上榜
    fmt.Println(cweskills.IsInTop25(79))  // true（XSS，排名第 1）
    fmt.Println(cweskills.IsInTop25(999)) // false

    // 遍历完整 Top 25，按危险顺序
    for i, id := range cweskills.CWETop25 {
        fmt.Printf("%2d. CWE-%d\n", i+1, id)
    }

    // 批量过滤扫描结果，只保留上榜弱点
    findings := []int{79, 89, 352, 12345}
    for _, id := range findings {
        if cweskills.IsInTop25(id) {
            fmt.Printf("CWE-%d 上榜 Top 25，需优先修复\n", id)
        }
    }
}
```

## 🖥️ CLI 用法

```bash
# 列出完整 Top 25（text）
cwe wellknown top25
```

```text
CWE Top 25 Most Dangerous Software Weaknesses (25 项):

   1. CWE-79
   2. CWE-89
   3. CWE-352
  ...
  25. CWE-770
```

```bash
# JSON 输出，便于脚本处理
cwe wellknown top25 -o json | jq '.[0:3]'

# 检查某 CWE 是否上榜（更高效）
cwe wellknown check CWE-79 89
```

```text
CWE-79: [Top 25 OWASP Top 10 (A03:2021-Injection)]
CWE-89: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
```

::: tip check 已覆盖 Top 25
`cwe wellknown check` 一次调用即返回该 CWE 在**所有**知名列表的归属，无需单独比对 `top25`。
:::

## 🎯 使用场景

<Badge type="tip" text="修复优先级" /> 将扫描结果与 Top 25 求交集，上榜弱点优先修复
<Badge type="info" text="安全培训" /> 以 Top 25 为大纲组织开发安全培训，覆盖最致命弱点
<Badge type="warning" text="风险评估" /> 评估代码库对高危弱点的暴露面，量化安全态势
<Badge type="info" text="门禁规则" /> 在 CI 中对引入 Top 25 弱点的 PR 阻断合并

## 🔗 相关链接

- [知名列表总览](./overview)
- [OWASP Top 10 (2021)](./owasp-top-10)
- [SANS Top 25](./sans-top-25)
- SDK：[CWETop25 列表](../sdk/cwe-top-25)、[IsInTop25](../sdk/is-in-top-25)
- CLI：[wellknown top25](../cli/wellknown-top25)、[wellknown check](../cli/wellknown-check)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
