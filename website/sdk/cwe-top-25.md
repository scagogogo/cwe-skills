---
title: CWETop25 CWE 最危险弱点列表
outline: [2, 3]
---

# 🏆 CWETop25 — CWE Top 25 最危险软件弱点（2024）

<Badge type="tip" text="2024 版" /> <Badge type="info" text="var []int" />

## 📋 变量声明

```go
var CWETop25 = []int{
	79, 89, 352, 862, 787, 22, 416, 125, 78, 94,
	120, 434, 476, 121, 502, 122, 863, 20, 284, 200,
	306, 918, 77, 639, 770,
}
```

## 📖 说明

`CWETop25` 是 MITRE 发布的**年度最危险软件弱点列表（2024 版）**，按危险程度由高到低排序。排名基于 NVD 漏洞数据的出现频率分析与 CVSS 评分计算，代表了对软件最严重的安全威胁。

::: tip 排序即危险程度
切片顺序即排名：`CWETop25[0]`（79，XSS）最危险，`CWETop25[24]`（770，资源分配无限制）第 25 位。如需展示「第几名」可取切片下标 +1。
:::

适用场景：

- 安全工具优先级排序
- 开发者安全培训重点
- 漏洞管理风险评估

## 📊 完整 25 项列表

| 排名 | CWE ID | 弱点名称 |
| --- | --- | --- |
| 1 | 79 | Cross-site Scripting (XSS) |
| 2 | 89 | SQL Injection |
| 3 | 352 | Cross-Site Request Forgery (CSRF) |
| 4 | 862 | Missing Authorization |
| 5 | 787 | Out-of-bounds Write |
| 6 | 22 | Path Traversal |
| 7 | 416 | Use After Free |
| 8 | 125 | Out-of-bounds Read |
| 9 | 78 | OS Command Injection |
| 10 | 94 | Code Injection |
| 11 | 120 | Buffer Copy without Checking Size of Input |
| 12 | 434 | Unrestricted Upload of File with Dangerous Type |
| 13 | 476 | NULL Pointer Dereference |
| 14 | 121 | Stack-based Buffer Overflow |
| 15 | 502 | Deserialization of Untrusted Data |
| 16 | 122 | Heap-based Buffer Overflow |
| 17 | 863 | Incorrect Authorization |
| 18 | 20 | Improper Input Validation |
| 19 | 284 | Improper Access Control |
| 20 | 200 | Exposure of Sensitive Information |
| 21 | 306 | Missing Authentication for Critical Function |
| 22 | 918 | Server-Side Request Forgery (SSRF) |
| 23 | 77 | Command Injection |
| 24 | 639 | Authorization Bypass Through User-Controlled Key |
| 25 | 770 | Allocation of Resources Without Limits or Throttling |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 直接遍历排名
	for rank, id := range cweskills.CWETop25 {
		fmt.Printf("#%2d  CWE-%d\n", rank+1, id)
	}

	// 判断某弱点是否上榜
	fmt.Println(cweskills.IsInTop25(89)) // true
	fmt.Println(cweskills.IsInTop25(42)) // false
}
```

## ⚠️ 常见错误

::: warning 不要修改切片
`CWETop25` 是包级 `var`，直接 `append` 或改写元素会污染全局状态，影响所有调用方。若需定制列表，请复制一份：

```go
mine := make([]int, len(cweskills.CWETop25))
copy(mine, cweskills.CWETop25)
```
:::

::: details 何时该用 CWETop25 而非 SANSTop25
两者都是「25 个最危险弱点」，但口径不同：CWETop25 基于 NVD 频率 + CVSS 评分逐年更新，偏数据驱动；SANSTop25 由 SANS Institute 与 MITRE 合作编制，偏「可被攻击者利用获取控制权」的编程错误。做年度风险量化优先 CWETop25，做编码规范培训可参考 SANSTop25。
:::

## 🔗 相关链接

- 成员判断函数：[IsInTop25](./is-in-top-25)
- 另两份列表：[OWASPTop10](./owasp-top-10) · [SANSTop25](./sans-top-25)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
