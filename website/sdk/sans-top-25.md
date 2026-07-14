---
title: SANSTop25 SANS 最危险软件错误
outline: [2, 3]
---

# 🏆 SANSTop25 — SANS Top 25 最危险软件错误

<Badge type="info" text="var []int" />

## 📋 变量声明

```go
var SANSTop25 = []int{
	89, 78, 79, 20, 22, 352, 416, 787, 125, 94,
	190, 434, 862, 287, 306, 863, 798, 502, 77, 119,
	639, 770, 918, 476, 200,
}
```

## 📖 说明

`SANSTop25` 是 SANS Institute 与 MITRE 合作编制的**25 个最危险软件错误列表**。它侧重于可被攻击者利用来获取系统控制权或进行数据窃取的编程错误，常用于编码规范培训与代码审计重点筛选。

::: tip 与 CWETop25 的排序含义不同
`CWETop25` 按危险程度排名；`SANSTop25` 的切片顺序是 SANS 发布时的呈现顺序，并非严格的「第 1 到第 25 危险」。不要把下标当作权威排名使用。
:::

## 📊 完整 25 项列表

| 序号 | CWE ID | 弱点名称 |
| --- | --- | --- |
| 1 | 89 | SQL Injection |
| 2 | 78 | OS Command Injection |
| 3 | 79 | Cross-site Scripting (XSS) |
| 4 | 20 | Improper Input Validation |
| 5 | 22 | Path Traversal |
| 6 | 352 | Cross-Site Request Forgery (CSRF) |
| 7 | 416 | Use After Free |
| 8 | 787 | Out-of-bounds Write |
| 9 | 125 | Out-of-bounds Read |
| 10 | 94 | Code Injection |
| 11 | 190 | Integer Overflow or Wraparound |
| 12 | 434 | Unrestricted Upload of File with Dangerous Type |
| 13 | 862 | Missing Authorization |
| 14 | 287 | Improper Authentication |
| 15 | 306 | Missing Authentication for Critical Function |
| 16 | 863 | Incorrect Authorization |
| 17 | 798 | Use of Hard-coded Credentials |
| 18 | 502 | Deserialization of Untrusted Data |
| 19 | 77 | Command Injection |
| 20 | 119 | Improper Restriction of Operations within the Bounds of a Memory Buffer |
| 21 | 639 | Authorization Bypass Through User-Controlled Key |
| 22 | 770 | Allocation of Resources Without Limits or Throttling |
| 23 | 918 | Server-Side Request Forgery (SSRF) |
| 24 | 476 | NULL Pointer Dereference |
| 25 | 200 | Exposure of Sensitive Information |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	// 遍历列表
	for _, id := range cweskills.SANSTop25 {
		fmt.Printf("CWE-%d\n", id)
	}

	// 判断成员
	fmt.Println(cweskills.IsInSANSTop25(190)) // true  整数溢出
	fmt.Println(cweskills.IsInSANSTop25(42))  // false
}
```

## ⚠️ 常见错误

::: warning 不要修改切片
`SANSTop25` 是包级 `var`，`append` 或改写会污染全局状态。需要定制列表请先 `copy` 一份再操作。
:::

::: details SANSTop25 与 CWETop25 的成员差异
两份列表高度重叠（SQL 注入、XSS、越界写等都同时上榜），但各有独有项：例如 CWE-190（整数溢出）与 CWE-119（内存缓冲区操作限制不当）出现在 SANSTop25，却不在 2024 版 CWETop25 中。交叉比对两份列表可识别「公认高危但未被年度数据突出」的弱点。
:::

## 🔗 相关链接

- 成员判断函数：[IsInSANSTop25](./is-in-sans-top-25)
- 另两份列表：[CWETop25](./cwe-top-25) · [OWASPTop10](./owasp-top-10)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
