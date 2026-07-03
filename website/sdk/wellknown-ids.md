---
title: wellknown_ids.go 知名列表概览
outline: [2, 3]
---

# 🏆 wellknown_ids.go — 知名视图与列表概览

<Badge type="info" text="源文件 wellknown_ids.go" />

## 📖 说明

`wellknown_ids.go` 集中维护 CWE 体系中**业界公认的高价值「知名列表」**，包括五张知名视图的 ID 常量、三份权威「最危险弱点」列表（CWE Top 25、OWASP Top 10、SANS Top 25），以及围绕它们的一组查询函数。

这些列表是安全工具排期、开发者培训、漏洞风险评估与合规检查的常用基准。本包把它们的 ID 数字固化为 Go 符号，避免在业务代码里散落魔法数。

::: tip 为何单独成文件
知名列表的成员会随年度版本更新，把它们与查询逻辑集中在 `wellknown_ids.go`，更新时只需改一处，所有依赖它的函数、CLI 子命令、文档同步生效。
:::

## 🏛️ 知名视图 ID 常量

| 常量 | 值 | 视图名称 |
| --- | --- | --- |
| `CWEViewResearchConcepts` | `1000` | 研究概念视图 |
| `CWEViewDevelopmentConcepts` | `699` | 软件开发视图 |
| `CWEViewHardwareDesign` | `1199` | 硬件设计视图 |
| `CWEViewCWECrossSection` | `888` | CWE 横截面视图 |
| `CWEViewComprehensiveDictionary` | `1400` | 综合 CWE 字典 |

详见 [IsInWellKnownView](./is-in-wellknown-view)。

## 📊 三份知名列表

| 变量 | 类型 | 说明 |
| --- | --- | --- |
| [`CWETop25`](./cwe-top-25) | `[]int` | CWE Top 25 最危险软件弱点（2024 版，按危险程度排序） |
| [`OWASPTop10`](./owasp-top-10) | `map[string][]int` | OWASP Top 10（2021 版）到 CWE ID 的映射 |
| [`SANSTop25`](./sans-top-25) | `[]int` | SANS Top 25 最危险软件错误 |

::: warning 版本固定
三份列表均内置为**快照**：CWETop25 为 2024 版、OWASPTop10 为 2021 版。若需逐年升级，请关注包版本更新。
:::

## 🛠️ 函数总览

| 函数 | 签名 | 说明 |
| --- | --- | --- |
| [`IsInTop25`](./is-in-top-25) | `func IsInTop25(cweID int) bool` | 判断 CWE ID 是否在 CWE Top 25 |
| [`IsInOWASPTop10`](./is-in-owasp-top-10) | `func IsInOWASPTop10(cweID int) bool` | 判断 CWE ID 是否在 OWASP Top 10 |
| [`IsInSANSTop25`](./is-in-sans-top-25) | `func IsInSANSTop25(cweID int) bool` | 判断 CWE ID 是否在 SANS Top 25 |
| [`GetOWASPCategory`](./get-owasp-category) | `func GetOWASPCategory(cweID int) string` | 取 CWE ID 所属的首个 OWASP 类别 |
| [`GetOWASPCategories`](./get-owasp-categories) | `func GetOWASPCategories(cweID int) []string` | 取 CWE ID 所属的全部 OWASP 类别 |
| [`IsInWellKnownView`](./is-in-wellknown-view) | `func IsInWellKnownView(viewID int) bool` | 判断视图 ID 是否为知名视图 |

## 🚀 快速上手

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 列表成员判断
	fmt.Println(cweskills.IsInTop25(89))      // true  SQL 注入
	fmt.Println(cweskills.IsInOWASPTop10(79)) // true  XSS 属于 A03 注入
	fmt.Println(cweskills.IsInSANSTop25(78))  // true  OS 命令注入

	// OWASP 类别查询
	fmt.Println(cweskills.GetOWASPCategory(79))   // A03:2021-Injection
	fmt.Println(cweskills.GetOWASPCategories(287)) // [A01:.. A04:.. A07:..] 多归属

	// 知名视图判断
	fmt.Println(cweskills.IsInWellKnownView(699)) // true  软件开发视图
}
```

## ⚠️ 实现细节

::: details 三个 IsIn* 都用线性查找
`IsInTop25`、`IsInOWASPTop10`、`IsInSANSTop25` 均遍历对应列表做 `==` 比较。列表规模很小（25 / 10 类 / 25），无需预建索引，调用成本可忽略。若你在热路径上对海量 ID 做批量判断，可考虑自行用 `map[int]struct{}` 缓存。
:::

::: warning OWASPTop10 的 map 迭代顺序不确定
`GetOWASPCategory` 返回「首个匹配」，但因 Go map 遍历顺序随机，当某个 CWE 属于多个类别时，返回的是哪一个类别**不确定**。需要稳定结果请用 [`GetOWASPCategories`](./get-owasp-categories) 自行排序。
:::

## 🔗 相关链接

- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
- CLI 入口：[`wellknown-check`](../cli/wellknown-check)
- 知名视图文档：[`../wellknown`](../wellknown/overview)
