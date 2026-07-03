---
title: ExtractFirstCWEID 提取首个 CWE ID
outline: [2, 3]
---

# 🆔 ExtractFirstCWEID — 从文本中提取第一个 CWE ID

## 📋 函数签名

```go
func ExtractFirstCWEID(text string) string
```

## 📖 说明

`ExtractFirstCWEID` 是 [`ExtractCWEIDs`](./extract-cwe-ids) 的「只取首个」快捷版本：用 `cweIDRegex.FindStringSubmatch` 找到第一处匹配，格式化为 `CWE-NNN` 返回。相比先取全部再取首元素，它**只做一次匹配即返回**，在长文本上更高效。

::: tip 无匹配的行为
找不到任何 CWE ID 时返回**空字符串 `""`**（不是 `nil`，因为返回类型是 `string`）。调用方可用 `if id == ""` 判断。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `text` | `string` | 待搜索的文本 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 首个 CWE ID | `string` | 标准形式 `CWE-NNN`；无匹配时为 `""` |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	fmt.Println(cweskills.ExtractFirstCWEID("See CWE-79 and CWE-89"))  // CWE-79
	fmt.Println(cweskills.ExtractFirstCWEID("issue: cwe 89 found"))   // CWE-89
	fmt.Println(cweskills.ExtractFirstCWEID("No CWE IDs here"))       // (空串)
	fmt.Println(cweskills.ExtractFirstCWEID(""))                      // (空串)

	// 典型用法：从漏洞描述中取主 CWE
	desc := "This flaw is a CWE-79 cross-site scripting issue."
	if id := cweskills.ExtractFirstCWEID(desc); id != "" {
		fmt.Println("主弱点:", id) // 主弱点: CWE-79
	}
}
```

## ⚖️ ExtractFirstCWEID vs ExtractCWEIDs

| 维度 | `ExtractFirstCWEID` | [`ExtractCWEIDs`](./extract-cwe-ids) |
| --- | --- | --- |
| 返回 | `string`（单个） | `[]string`（全部） |
| 内部调用 | `FindStringSubmatch`（一次） | `FindAllStringSubmatch`（全部） |
| 无匹配 | `""` | `[]string{}` |
| 性能 | 长文本更优 | 需扫描全文 |
| 适用 | 只关心首个/主弱点 | 需要完整清单 |

## ⚠️ 常见错误

::: warning 用空串判断而非 nil
返回类型是 `string`，无匹配返回 `""`。不要写成 `if id == nil`（编译错误）或 `if id == "nil"`，正确写法是 `if id == ""`。
:::

## 🔗 相关链接

- 取全部：[ExtractCWEIDs](./extract-cwe-ids)
- 单个解析：[ParseCWEID](./parse-cwe-id)
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
