---
title: IsInSANSTop25 判断 CWE 是否在 SANS Top 25
outline: [2, 3]
---

# ✅ IsInSANSTop25 — 判断 CWE ID 是否在 SANS Top 25

## 📋 函数签名

```go
func IsInSANSTop25(cweID int) bool
```

## 📖 说明

`IsInSANSTop25` 遍历 [`SANSTop25`](./sans-top-25) 切片做线性查找，判断给定 CWE ID 是否上榜 SANS Top 25 最危险软件错误。

::: tip 实现与 IsInTop25 同构
两者都是 `for _, id := range 切片 { if id == cweID { return true } }`，仅数据源不同。25 个元素的线性查找成本可忽略。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cweID` | `int` | 待检查的 CWE ID 数字 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 是否上榜 | `bool` | 在 SANS Top 25 中返回 `true`，否则 `false` |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	fmt.Println(cweskills.IsInSANSTop25(190)) // true  整数溢出
	fmt.Println(cweskills.IsInSANSTop25(119)) // true  内存缓冲区操作限制不当
	fmt.Println(cweskills.IsInSANSTop25(42))  // false

	// 典型用法：编码规范培训重点筛选
	candidates := []int{89, 78, 79, 190, 42}
	for _, id := range candidates {
		if cweskills.IsInSANSTop25(id) {
			fmt.Printf("CWE-%d 纳入培训重点\n", id)
		}
	}
}
```

## ⚖️ IsInTop25 vs IsInSANSTop25

| 维度 | [`IsInTop25`](./is-in-top-25) | `IsInSANSTop25` |
| --- | --- | --- |
| 数据源 | [`CWETop25`](./cwe-top-25)（2024） | [`SANSTop25`](./sans-top-25) |
| 列表来源 | MITRE 基于 NVD 频率 + CVSS | SANS Institute + MITRE 合作 |
| 偏重 | 年度数据驱动的高危排名 | 可被利用获取控制权的编程错误 |
| 成员差异 | 不含 190、119 | 含 190、119 等 SANS 独有项 |

::: tip 两份列表互为补充
 CWE-190（整数溢出）和 CWE-119（内存缓冲区操作限制不当）在 SANS Top 25 上榜却不在 2024 版 CWE Top 25。做静态分析规则集时，建议 `IsInTop25(id) || IsInSANSTop25(id)` 联合判断以扩大覆盖。
:::

## ⚠️ 常见错误

::: warning 传字符串会编译失败
`IsInSANSTop25` 接收 `int`。若持有 `"CWE-190"`，先用 [`ParseCWEID`](./parse-cwe-id) 转 `int` 再调用。
:::

## 🔗 相关链接

- 数据源变量：[SANSTop25](./sans-top-25)
- 同族函数：[IsInTop25](./is-in-top-25) · [IsInOWASPTop10](./is-in-owasp-top-10)
- 字符串转 ID：[ParseCWEID](./parse-cwe-id)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
