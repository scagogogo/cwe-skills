---
title: ExtractCWEIDs 从文本提取 CWE ID
outline: [2, 3]
---

# 🆔 ExtractCWEIDs — 从文本中提取所有 CWE ID

## 📋 函数签名

```go
func ExtractCWEIDs(text string) []string
```

## 📖 说明

`ExtractCWEIDs` 在一段自然语言文本中扫描所有匹配 CWE ID 模式的子串，返回标准化的 `CWE-NNN` 列表。它使用包级正则 `cweIDRegex`：

```go
var cweIDRegex = regexp.MustCompile("(?i)CWE[-\\s]?(\\d+)")
```

匹配到的每个结果都会被解析为整数并重新格式化为 `CWE-NNN`，因此**输出始终是标准形式**（大写、无空格、无前导零）。

::: tip 输出特性
- 按文本中出现顺序返回
- 每个匹配项都会出现（**不去重**，如需去重请自行处理）
- 无匹配时返回**空切片 `[]string{}`**，非 `nil`
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `text` | `string` | 待搜索的文本 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| CWE ID 列表 | `[]string` | 标准形式；无匹配时为空切片 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	text := "See CWE-79 and cwe-89 for details. CWE79 is also mentioned."
	ids := cweskills.ExtractCWEIDs(text)
	fmt.Println(ids) // [CWE-79 CWE-89 CWE-79]

	// 去重保序
	seen := map[string]bool{}
	uniq := make([]string, 0, len(ids))
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			uniq = append(uniq, id)
		}
	}
	fmt.Println(uniq) // [CWE-79 CWE-89]

	// 空文本与无匹配
	fmt.Println(cweskills.ExtractCWEIDs(""))            // []
	fmt.Println(cweskills.ExtractCWEIDs("no cwe here")) // []
}
```

## ⚠️ 常见错误

::: warning 不会提取裸数字
正则要求带 `CWE` 前缀。文本中的 `79` 不会被识别为 CWE ID——这是有意为之，避免把版本号、行号等误判为弱点编号。若需把纯数字也视为 CWE ID，请用 [`ParseCWEID`](./parse-cwe-id) 单独处理。
:::

::: details 为什么不去重？
真实文本里同一 CWE ID 可能被多次提及，是否去重取决于业务语义（统计频次 vs 唯一列表）。SDK 保持原始顺序与重复，把去重决策留给调用方。
:::

## 🔗 相关链接

- 只要首个：[ExtractFirstCWEID](./extract-first-cwe-id)
- 单个解析：[ParseCWEID](./parse-cwe-id)
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
