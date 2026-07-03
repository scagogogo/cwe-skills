---
title: IsInTop25 判断 CWE 是否在 Top 25
outline: [2, 3]
---

# ✅ IsInTop25 — 判断 CWE ID 是否在 CWE Top 25

## 📋 函数签名

```go
func IsInTop25(cweID int) bool
```

## 📖 说明

`IsInTop25` 遍历 [`CWETop25`](./cwe-top-25) 切片做线性查找，判断给定 CWE ID 是否上榜 2024 版 CWE Top 25 最危险软件弱点。

::: tip 实现极简
内部就是 `for _, id := range CWETop25 { if id == cweID { return true } }`。25 个元素的线性查找成本可忽略，无需预建索引。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cweID` | `int` | 待检查的 CWE ID 数字（如 `89`） |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 是否上榜 | `bool` | 在 CWE Top 25 中返回 `true`，否则 `false` |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	fmt.Println(cweskills.IsInTop25(89))  // true  SQL 注入（排名第 2）
	fmt.Println(cweskills.IsInTop25(79))  // true  XSS（排名第 1）
	fmt.Println(cweskills.IsInTop25(42))  // false 未上榜

	// 典型用法：扫描结果按上榜情况过滤
	type finding struct {
		cweID int
		msg   string
	}
	findings := []finding{
		{89, "SQL 注入"},
		{42, "其他问题"},
	}
	for _, f := range findings {
		if cweskills.IsInTop25(f.cweID) {
			fmt.Printf("高危（Top25）：CWE-%d %s\n", f.cweID, f.msg)
		}
	}
}
```

## ⚖️ 三个 IsIn* 对比

| 维度 | `IsInTop25` | [`IsInOWASPTop10`](./is-in-owasp-top-10) | [`IsInSANSTop25`](./is-in-sans-top-25) |
| --- | --- | --- | --- |
| 数据源 | [`CWETop25`](./cwe-top-25) | [`OWASPTop10`](./owasp-top-10) | [`SANSTop25`](./sans-top-25) |
| 列表规模 | 25 | 10 类别 | 25 |
| 适用 | 年度风险量化 | Web 应用合规 | 编码规范培训 |

::: tip 同时检查三份列表
若想判断某 CWE 是否出现在**任意一份**知名列表中，可串联调用：

```go
hit := cweskills.IsInTop25(id) ||
    cweskills.IsInOWASPTop10(id) ||
    cweskills.IsInSANSTop25(id)
```
:::

## ⚠️ 常见错误

::: warning 传错参数类型
`IsInTop25` 接收 `int`，不是字符串。若你手里是 `"CWE-89"`，先用 [`ParseCWEID`](./parse-cwe-id) 转成 `int`：

```go
id, err := cweskills.ParseCWEID("CWE-89")
if err == nil && cweskills.IsInTop25(id) {
	// ...
}
```
:::

## 🔗 相关链接

- 数据源变量：[CWETop25](./cwe-top-25)
- 同族函数：[IsInOWASPTop10](./is-in-owasp-top-10) · [IsInSANSTop25](./is-in-sans-top-25)
- 字符串转 ID：[ParseCWEID](./parse-cwe-id)
- 概览：[wellknown_ids.go](./wellknown-ids)
- 源文件：[`wellknown_ids.go`](https://github.com/scagogogo/cwe-skills/blob/main/wellknown_ids.go)
