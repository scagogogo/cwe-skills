---
title: IsInWellKnownView 判断知名视图
outline: [2, 3]
---

# ✅ IsInWellKnownView — 判断视图 ID 是否为知名视图

## 📋 函数签名

```go
func IsInWellKnownView(viewID int) bool
```

## 📖 说明

`IsInWellKnownView` 用 `switch` 判断给定视图 ID 是否属于 CWE 体系中的五个**知名视图**。这五个视图是 MITRE 官方维护、被工具链广泛引用的高价值视图，本包把它们固化为常量。

::: tip 用常量而非魔法数
判断时优先使用包级常量（如 `cweskills.CWEViewResearchConcepts`），而非裸数字 `1000`，可读性与可维护性更好。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `viewID` | `int` | 待检查的视图 ID 数字 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 是否知名 | `bool` | 属于五个知名视图之一返回 `true`，否则 `false` |

## 🏛️ 五个知名视图常量

| 常量 | 值 | 视图名称 | 说明 |
| --- | --- | --- | --- |
| `CWEViewResearchConcepts` | `1000` | 研究概念视图 | 按抽象概念组织 CWE 条目的层次结构 |
| `CWEViewDevelopmentConcepts` | `699` | 软件开发视图 | 按软件开发活动组织 CWE 条目 |
| `CWEViewHardwareDesign` | `1199` | 硬件设计视图 | 按硬件设计活动组织 CWE 条目 |
| `CWEViewCWECrossSection` | `888` | CWE 横截面视图 | 提供 CWE 条目的横截面视图 |
| `CWEViewComprehensiveDictionary` | `1400` | 综合 CWE 字典 | 包含所有 CWE 条目的综合视图 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 用常量更清晰
	fmt.Println(cweskills.IsInWellKnownView(cweskills.CWEViewResearchConcepts)) // true
	fmt.Println(cweskills.IsInWellKnownView(cweskills.CWEViewDevelopmentConcepts)) // true
	fmt.Println(cweskills.IsInWellKnownView(699)) // true  软件开发视图

	// 非知名视图
	fmt.Println(cweskills.IsInWellKnownView(100)) // false
	fmt.Println(cweskills.IsInWellKnownView(0))   // false

	// 典型用法：从 API 拿到视图后做知名性标注
	viewID := 1000
	if cweskills.IsInWellKnownView(viewID) {
		fmt.Printf("视图 %d 为知名视图，建议优先展示\n", viewID)
	}
}
```

## ⚖️ 与列表类 IsIn* 的区别

| 维度 | `IsInWellKnownView` | [`IsInTop25`](./is-in-top-25) 等 |
| --- | --- | --- |
| 判断对象 | 视图 ID（view） | 弱点 ID（weakness） |
| 数据源 | 5 个 `const` | `CWETop25` 等切片/map |
| 实现 | `switch` | 线性遍历 |
| 用途 | 视图导航/筛选 | 弱点风险标注 |

::: warning 视图 ID 与弱点 ID 是不同空间
视图 ID（如 1000）和弱点 ID（如 CWE-1000）数值可能重叠但语义不同。`IsInWellKnownView` 判断的是**视图 ID**，不要用它检查弱点是否上榜 Top 25——那是 [`IsInTop25`](./is-in-top-25) 的职责。
:::

## ⚠️ 常见错误

::: warning 别把弱点 ID 当视图 ID 传
CWE-1000（研究概念视图）本身是一个视图条目，其 ID `1000` 恰好也是知名视图常量值。但 CWE-699 是弱点 ID 还是视图 ID 取决于上下文：作为视图时它指「软件开发视图」。调用本函数前请确认你手里的是视图 ID。
:::

::: details 为什么只固定这五个
这五个视图覆盖了研究（1000）、开发（699）、硬件（1199）、横截面（888）、综合字典（1400）五个维度，是 MITRE 官方推荐工具优先支持的视图。其他视图（如某些领域专用视图）未纳入知名集合，调用本函数会返回 `false`，属预期行为。
:::

## 🔗 相关链接

- 概览：[wellknown_ids.go](./wellknown-ids)
- 视图文档：[`../wellknown`](../wellknown/overview)
- 视图相关 API：[api-get-view](./api-get-view) · [build-view-tree](./build-view-tree)
- 同族列表判断：[IsInTop25](./is-in-top-25) · [IsInOWASPTop10](./is-in-owasp-top-10) · [IsInSANSTop25](./is-in-sans-top-25)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
