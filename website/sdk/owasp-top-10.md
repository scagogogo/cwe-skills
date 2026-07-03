---
title: OWASPTop10 OWASP 十大风险映射
outline: [2, 3]
---

# 🏆 OWASPTop10 — OWASP Top 10（2021）到 CWE 的映射

<Badge type="tip" text="2021 版" /> <Badge type="info" text="var map[string][]int" />

## 📋 变量声明

```go
var OWASPTop10 = map[string][]int{
	"A01:2021-Broken Access Control":                {22, 23, 35, ...},
	"A02:2021-Cryptographic Failures":               {260, 261, ...},
	// ... 共 10 个类别
}
```

## 📖 说明

`OWASPTop10` 是 OWASP Top 10（2021 版）Web 应用安全风险类别到相关 CWE ID 的映射。OWASP Top 10 是 Web 应用安全风险的行业标准列表，每个类别下挂若干 CWE。

::: warning A06 类别为空
`A06:2021-Vulnerable and Outdated Components`（易受攻击和过时的组件）在源码中映射为空切片 `{}`。该类别强调的是依赖管理与漏洞数据库联动，而非映射到具体 CWE 编号，因此 `IsInOWASPTop10` 永远不会因 A06 返回 `true`。
:::

适用场景：

- Web 应用安全评估
- 合规性检查
- 安全测试用例设计

## 📊 10 类别映射表

| 类别 | CWE ID 集合 | 数量 |
| --- | --- | --- |
| A01:2021-Broken Access Control | 22, 23, 35, 59, 78, 94, 200, 201, 219, 255, 269, 276, 284, 285, 287, 306, 346, 639, 651, 668, 862, 863, 922 | 23 |
| A02:2021-Cryptographic Failures | 260, 261, 295, 310, 311, 312, 319, 325, 326, 327, 328, 329, 330, 337, 338, 340, 347, 522, 757, 759, 760, 780 | 22 |
| A03:2021-Injection | 20, 74, 75, 77, 78, 79, 80, 83, 87, 88, 89, 90, 91, 94, 95, 96, 97, 98, 99, 100, 113, 116, 138, 141, 147, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171 | 47 |
| A04:2021-Insecure Design | 209, 235, 256, 267, 284, 285, 287, 311, 326, 384, 393, 664, 863 | 13 |
| A05:2021-Security Misconfiguration | 2, 5, 11, 13, 15, 16, 260, 315, 520, 526, 537, 540, 544, 546, 547, 548, 611, 613, 614, 759, 760, 1021 | 22 |
| A06:2021-Vulnerable and Outdated Components | *(空)* | 0 |
| A07:2021-Identification and Authentication Failures | 255, 256, 258, 259, 260, 287, 288, 290, 294, 295, 297, 306, 307, 346, 384, 521, 522, 523, 613, 620, 640, 798 | 22 |
| A08:2021-Software and Data Integrity Failures | 311, 345, 353, 426, 494, 502, 565, 610, 653, 754, 829, 912 | 12 |
| A09:2021-Security Logging and Monitoring Failures | 117, 223, 532, 778 | 4 |
| A10:2021-Server-Side Request Forgery (SSRF) | 918, 1021 | 2 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 遍历类别与其 CWE
	for category, ids := range cweskills.OWASPTop10 {
		fmt.Printf("%s -> %d 个 CWE\n", category, len(ids))
	}

	// 判断 + 查类别
	fmt.Println(cweskills.IsInOWASPTop10(79))        // true
	fmt.Println(cweskills.GetOWASPCategory(79))      // A03:2021-Injection
	fmt.Println(cweskills.GetOWASPCategories(287))   // 多归属示例
}
```

## ⚠️ 常见错误

::: warning map 迭代顺序不确定
`OWASPTop10` 是 `map[string][]int`，Go 的 map 遍历顺序随机。直接 `for k, v := range` 打印时类别顺序不固定；如需稳定输出，请先排序 key：

```go
keys := make([]string, 0, len(cweskills.OWASPTop10))
for k := range cweskills.OWASPTop10 {
	keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
	fmt.Println(k, cweskills.OWASPTop10[k])
}
```
:::

::: details 一个 CWE 可能属于多个类别
例如 CWE-287（Improper Authentication）同时出现在 A01、A04、A07 三类。因此：
- [`GetOWASPCategory`](./get-owasp-category) 返回的「首个匹配」因 map 遍历随机而**不确定**。
- 需要全部归属请用 [`GetOWASPCategories`](./get-owasp-categories)。
:::

## 🔗 相关链接

- 成员判断：[IsInOWASPTop10](./is-in-owasp-top-10)
- 类别查询：[GetOWASPCategory](./get-owasp-category) · [GetOWASPCategories](./get-owasp-categories)
- 另两份列表：[CWETop25](./cwe-top-25) · [SANSTop25](./sans-top-25)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
