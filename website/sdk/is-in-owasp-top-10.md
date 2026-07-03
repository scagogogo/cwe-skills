---
title: IsInOWASPTop10 判断 CWE 是否在 OWASP Top 10
outline: [2, 3]
---

# ✅ IsInOWASPTop10 — 判断 CWE ID 是否在 OWASP Top 10

## 📋 函数签名

```go
func IsInOWASPTop10(cweID int) bool
```

## 📖 说明

`IsInOWASPTop10` 遍历 [`OWASPTop10`](./owasp-top-10) map 的所有类别及其 CWE ID 列表，判断给定 ID 是否属于 OWASP Top 10（2021 版）中的任意类别。

::: warning A06 类别恒为空
`A06:2021-Vulnerable and Outdated Components` 在源码中映射为空切片，因此任何 CWE ID 都不会因 A06 而命中。该函数实际覆盖 9 个非空类别。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cweID` | `int` | 待检查的 CWE ID 数字 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 是否归属 | `bool` | 在 OWASP Top 10 任意类别中返回 `true`，否则 `false` |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	fmt.Println(cweskills.IsInOWASPTop10(79))  // true  XSS 属于 A03 注入
	fmt.Println(cweskills.IsInOWASPTop10(89))  // true  SQL 注入属于 A03
	fmt.Println(cweskills.IsInOWASPTop10(42))  // false

	// 典型用法：Web 应用扫描结果标注 OWASP 归属
	id := 918 // SSRF
	if cweskills.IsInOWASPTop10(id) {
		fmt.Printf("CWE-%d 命中 OWASP Top 10：%s\n",
			id, cweskills.GetOWASPCategory(id))
	}
}
```

## ⚖️ 与 GetOWASPCategory 的关系

<Badge type="tip" text="判断 vs 查询" />

| 维度 | `IsInOWASPTop10` | [`GetOWASPCategory`](./get-owasp-category) |
| --- | --- | --- |
| 返回值 | `bool` | `string`（类别名） |
| 适用 | 只需 yes/no | 需要类别名 |
| 成本 | 命中即返回 | 同样命中即返回 |

::: tip 先判断再查询更省心
`GetOWASPCategory` 在不命中时返回空串，因此通常无需先 `IsInOWASPTop10` 再 `GetOWASPCategory`，直接取类别、判空串即可。但若你只需布尔结果（如过滤），用 `IsInOWASPTop10` 语义更清晰。
:::

## ⚠️ 常见错误

::: warning 不要假设 CWE 只属于一个类别
一个 CWE 可能同时映射到多个 OWASP 类别（如 CWE-287 同时属于 A01、A04、A07）。`IsInOWASPTop10` 只回答「是否归属」，不区分归属几个；需要全部归属请用 [`GetOWASPCategories`](./get-owasp-categories)。
:::

::: details 为什么遍历顺序不影响布尔结果
`IsInOWASPTop10` 遍历 map 找到任意匹配即返回 `true`，因此 Go map 的随机迭代顺序不影响最终布尔值（true 就是 true）。但「命中哪一个类别」是随机的——这正是 [`GetOWASPCategory`](./get-owasp-category) 行为不确定的根因。
:::

## 🔗 相关链接

- 数据源变量：[OWASPTop10](./owasp-top-10)
- 类别查询：[GetOWASPCategory](./get-owasp-category) · [GetOWASPCategories](./get-owasp-categories)
- 同族函数：[IsInTop25](./is-in-top-25) · [IsInSANSTop25](./is-in-sans-top-25)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
