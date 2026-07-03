---
title: GetOWASPCategories 取全部 OWASP 类别
outline: [2, 3]
---

# 🔍 GetOWASPCategories — 取 CWE ID 所属的全部 OWASP 类别

## 📋 函数签名

```go
func GetOWASPCategories(cweID int) []string
```

## 📖 说明

`GetOWASPCategories` 遍历 [`OWASPTop10`](./owasp-top-10) map 的所有类别，收集给定 CWE ID 所属的**全部** OWASP 类别名。与 [`GetOWASPCategory`](./get-owasp-category) 只返回首个匹配不同，它能完整列出多归属 CWE 的所有风险类别。

::: tip 解决多归属不确定性
CWE-287 同时属于 A01、A04、A07 三个类别。本函数返回全部三个，避免了 [`GetOWASPCategory`](./get-owasp-category) 因 map 迭代随机导致「返回哪个不确定」的问题。需要稳定顺序时，对返回切片 `sort.Strings` 即可。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cweID` | `int` | 待查询的 CWE ID 数字 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 类别列表 | `[]string` | 所有匹配的 OWASP 类别名；无匹配返回 `nil`（长度 0） |

## ✅ 示例

```go
package main

import (
	"fmt"
	"sort"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 单一归属
	fmt.Println(cweskills.GetOWASPCategories(79))  // [A03:2021-Injection]

	// 多归属（CWE-287 同时属于 A01、A04、A07）
	cats := cweskills.GetOWASPCategories(287)
	sort.Strings(cats) // 排序保证稳定输出
	fmt.Println(cats)  // [A01:2021-Broken Access Control A04:2021-Insecure Design A07:2021-Identification and Authentication Failures]

	// 无归属
	fmt.Println(cweskills.GetOWASPCategories(42)) // []

	// 典型用法：合规报告穷举所有风险类别
	for _, c := range cweskills.GetOWASPCategories(287) {
		fmt.Printf("CWE-287 触发 OWASP 风险：%s\n", c)
	}
}
```

## ⚖️ GetOWASPCategories vs GetOWASPCategory

| 维度 | `GetOWASPCategories` | [`GetOWASPCategory`](./get-owasp-category) |
| --- | --- | --- |
| 返回值 | `[]string` | `string` |
| 多归属 | 返回全部 | 仅首个（不确定） |
| 无匹配 | `nil`（空切片） | `""`（空串） |
| 适用 | 合规报告、完整风险画像 | 单标签展示 |
| 稳定性 | 内容稳定（顺序需排序） | 多归属时不稳定 |

::: details 返回 nil 还是空切片？
源码中 `var categories []string` 初始为 `nil`，未匹配时不会 append，因此返回 `nil`。在 Go 中 `nil` 切片长度为 0，可直接 `for range` 安全遍历，无需特判。但若要序列化为 JSON，`nil` 会输出 `null` 而非 `[]`，需要时可用 `make([]string, 0)` 替代。
:::

## ⚠️ 常见错误

::: warning 切片顺序仍随 map 迭代
虽然返回了全部类别，但其**顺序**仍取决于 map 迭代顺序，多次调用顺序可能不同。展示给用户前请 `sort.Strings(cats)` 固定顺序。
:::

::: warning 不要假设长度
一个 CWE 可能属于 0、1 或多个类别。不要写 `if len(cats) == 1` 之类的硬编码假设，应基于「是否为空」或「是否包含某类别」判断。
:::

## 🔗 相关链接

- 数据源变量：[OWASPTop10](./owasp-top-10)
- 单值版本：[GetOWASPCategory](./get-owasp-category)
- 布尔判断：[IsInOWASPTop10](./is-in-owasp-top-10)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
