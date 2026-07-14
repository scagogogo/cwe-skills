---
title: GetOWASPCategory 取首个 OWASP 类别
outline: [2, 3]
---

# 🔍 GetOWASPCategory — 取 CWE ID 所属的首个 OWASP 类别

## 📋 函数签名

```go
func GetOWASPCategory(cweID int) string
```

## 📖 说明

`GetOWASPCategory` 遍历 [`OWASPTop10`](./owasp-top-10) map，返回给定 CWE ID 所属的**第一个匹配** OWASP 类别名（如 `"A03:2021-Injection"`）。若不属于任何类别，返回空字符串 `""`。

::: warning 多归属时返回结果不确定
当某个 CWE 同时属于多个 OWASP 类别（如 CWE-287 属于 A01、A04、A07），本函数只返回「第一个匹配」。但由于 Go map 遍历顺序随机，**这里的「第一个」是不确定的**——每次运行可能返回不同类别。需要稳定、完整的归属请用 [`GetOWASPCategories`](./get-owasp-categories)。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cweID` | `int` | 待查询的 CWE ID 数字 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 类别名 | `string` | 第一个匹配的 OWASP 类别（如 `"A01:2021-Broken Access Control"`）；无匹配返回 `""` |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	// 单一归属
	fmt.Println(cweskills.GetOWASPCategory(79))  // A03:2021-Injection
	fmt.Println(cweskills.GetOWASPCategory(918)) // A10:2021-Server-Side Request Forgery (SSRF)

	// 无归属
	fmt.Println(cweskills.GetOWASPCategory(42))  // (空串)

	// 典型用法：判空串即可，无需先 IsInOWASPTop10
	if cat := cweskills.GetOWASPCategory(89); cat != "" {
		fmt.Printf("CWE-89 归属 OWASP 类别：%s\n", cat)
	}
}
```

## ⚖️ GetOWASPCategory vs GetOWASPCategories

<Badge type="tip" text="单值 vs 全部" />

| 维度 | `GetOWASPCategory` | [`GetOWASPCategories`](./get-owasp-categories) |
| --- | --- | --- |
| 返回值 | `string` | `[]string` |
| 多归属时 | 只返回一个（**不确定**哪个） | 返回全部 |
| 适用 | 只需展示一个标签 | 需要完整归属（合规报告） |
| 稳定性 | ⚠️ 多归属时不稳定 | ✅ 返回全部，但顺序仍随 map 迭代 |

::: tip 何时选哪个
- UI 上只显示一个 OWASP 标签 → `GetOWASPCategory`，简单。
- 生成合规报告、需要穷举所有风险类别 → `GetOWASPCategories`，再自行 `sort.Strings`。
:::

## ⚠️ 常见错误

::: warning 不要依赖返回值判断「主归属」
因 map 迭代随机，CWE-287 多次调用 `GetOWASPCategory` 可能交替返回 A01、A04、A07。若业务逻辑依赖「该 CWE 的主 OWASP 类别」，请在应用层固化映射（如自建 `map[int]string` 缓存），不要依赖本函数的随机结果。
:::

::: details 为什么不用 map 顺序固定方案
Go 规范明确 map 遍历顺序随机化，是为了避免开发者依赖顺序。本包保持与标准库一致，未对 `OWASPTop10` 做 key 排序。若需固定单值结果，最稳妥的做法是 `GetOWASPCategories` + 排序后取 `[0]`。
:::

## 🔗 相关链接

- 数据源变量：[OWASPTop10](./owasp-top-10)
- 全部归属版本：[GetOWASPCategories](./get-owasp-categories)
- 布尔判断：[IsInOWASPTop10](./is-in-owasp-top-10)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
