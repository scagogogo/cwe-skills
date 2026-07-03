---
title: 按关键词查找
outline: [2, 3]
---

# 🔍 按关键词查找

`FindByKeyword` 在全库弱点中做**文本搜索**，匹配 `Name`、`Description`、`ExtendedDescription` 三个字段。常用于「按漏洞名/描述片段」检索弱点。

## 📐 函数签名

```go
func FindByKeyword(r *Registry, keyword string) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `keyword` | `string` | 搜索关键词 |
| 返回 | `[]*CWE` | 任意字段包含关键词的弱点列表 |

::: tip 匹配规则
匹配通常**不区分大小写**，对三个字段做子串包含判断。空关键词的返回行为以源码实现为准（一般返回空或全部）。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	c1 := cweskills.NewCWE(79, "Improper Neutralization XSS")
	c1.Description = "Cross-site Scripting"
	r.Register(c1)
	c2 := cweskills.NewCWE(89, "SQL Injection")
	c2.Description = "SQL injection weakness"
	r.Register(c2)

	for _, c := range cweskills.FindByKeyword(r, "sql") {
		fmt.Println(c.CWEID()) // CWE-89
	}
	for _, c := range cweskills.FindByKeyword(r, "xss") {
		fmt.Println(c.CWEID()) // CWE-79
	}
}
```

## ⚠️ 注意事项

::: warning 仅子串匹配
`FindByKeyword` 是朴素的子串包含搜索，不支持正则、通配符或分词。复杂查询建议先 `GetAll()` 再用 [Filter](./filter) 组合条件。
:::

::: details 字段范围
搜索覆盖 `Name`、`Description`、`ExtendedDescription`，**不**搜索 `References`、`Notes` 等字段。如需全文检索，需自行遍历。
:::

## 🔗 相关链接

- 多条件过滤：[Filter](./filter)
- 按 ID 精确查找：[按 ID 查找](./find-by-id)
- 注册表遍历：[GetAll](./registry-operations)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
