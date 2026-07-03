---
title: cwe_utils 工具包概览
outline: [2, 3]
---

# 🆔 cwe_utils — CWE ID 工具包

`cweskills` 包中负责 CWE 标识符解析、格式化、校验与提取的工具集合，全部位于 `cwe_utils.go`。它提供一组**纯函数**，不依赖网络与本地数据，可在任何场景下安全调用。

## 🎯 这个包解决什么问题

CWE ID 在真实世界里的写法五花八门：

| 来源 | 写法 | 是否规范 |
| --- | --- | --- |
| MITRE 官方 XML | `CWE-79` | ✅ |
| 漏洞报告 | `cwe-79`、`CWE79` | ⚠️ 大小写/缺连字符 |
| 代码注释 | `CWE 79`、`cwe 79` | ⚠️ 带空格 |
| 数据库字段 | `79`、`079` | ⚠️ 纯数字/前导零 |

`cwe_utils` 把上述所有变体统一收敛为标准形式 `CWE-NNN`，并提供解析、判断、提取、比较能力。

::: tip 设计原则
所有函数**无副作用、无全局状态**，对同一输入永远返回同一输出，可放心用于并发场景。
:::

## 📦 函数清单

| 函数 | 签名 | 用途 |
| --- | --- | --- |
| [`FormatCWEID`](./format-cwe-id) | `func FormatCWEID(id string) (string, error)` | 字符串 → `CWE-NNN` |
| [`ParseCWEID`](./parse-cwe-id) | `func ParseCWEID(id string) (int, error)` | 字符串 → 整数 |
| [`FormatCWEIDFromInt`](./format-cwe-id-from-int) | `func FormatCWEIDFromInt(id int) string` | 整数 → `CWE-NNN` |
| [`IsCWEID`](./is-cwe-id) | `func IsCWEID(text string) bool` | 是否合法 CWE ID |
| [`ValidateCWEID`](./validate-cwe-id) | `func ValidateCWEID(text string) error` | 验证并返回详细错误 |
| [`ExtractCWEIDs`](./extract-cwe-ids) | `func ExtractCWEIDs(text string) []string` | 从文本提取全部 |
| [`ExtractFirstCWEID`](./extract-first-cwe-id) | `func ExtractFirstCWEID(text string) string` | 从文本提取首个 |
| [`CompareCWEIDs`](./compare-cwe-ids) | `func CompareCWEIDs(a, b string) (int, error)` | 按数字大小比较 |

## 🔍 核心正则

所有「提取类」函数共享同一正则：

```go
var cweIDRegex = regexp.MustCompile("(?i)CWE[-\\s]?(\\d+)")
```

- `(?i)`：大小写不敏感，匹配 `cwe`、`CWE`、`Cwe`
- `[-\\s]?`：可选的连字符或空白，兼容 `CWE-79`、`CWE 79`、`CWE79`
- `(\\d+)`：捕获数字部分

::: warning 注意
正则只匹配带 `CWE` 前缀的写法。**纯数字**（如 `"79"`）只能由 `ParseCWEID`、`IsCWEID`、`ValidateCWEID`、`FormatCWEID` 处理，`ExtractCWEIDs` 不会从一段文本里把裸数字当成 CWE ID 提取出来——否则会产生大量误报。
:::

## 🚀 快速上手

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 格式化
	id, _ := cweskills.FormatCWEID("cwe-0079")
	fmt.Println(id) // CWE-79

	// 从文本提取
	ids := cweskills.ExtractCWEIDs("See CWE-79 and cwe-89 for details")
	fmt.Println(ids) // [CWE-79 CWE-89]

	// 比较
	cmp, _ := cweskills.CompareCWEIDs("CWE-79", "CWE-89")
	fmt.Println(cmp) // -1
}
```

## 🧭 选型建议

<Badge type="tip" text="判断用" /> 需要布尔结果 → [`IsCWEID`](./is-cwe-id)
<Badge type="info" text="报错用" /> 需要错误原因 → [`ValidateCWEID`](./validate-cwe-id)
<Badge type="warning" text="批量用" /> 处理自然语言 → [`ExtractCWEIDs`](./extract-cwe-ids)

## 🔗 相关链接

- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
- 错误类型：[InvalidCWEIDError](./invalid-cwe-id-error)
- CLI 对应命令：[parse](../cli/parse)、[format](../cli/format)、[extract](../cli/extract)、[compare](../cli/compare)
