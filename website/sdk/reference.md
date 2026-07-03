---
title: Reference 参考文献
outline: [2, 3]
---

# 📚 Reference — 参考文献

`Reference` 提供关于某个弱点的更多信息来源，是 `CWE.References` 字段的元素类型（`Category`、`View` 也各自带 `References`）。每条参考文献包含作者、标题与可选链接。

## 📋 结构体定义

```go
type Reference struct {
    ID     int    `json:"id,omitempty" xml:"Reference_ID,omitempty"`
    Author string `json:"author,omitempty" xml:"Author,omitempty"`
    Title  string `json:"title" xml:"Title"`
    URL    string `json:"url,omitempty" xml:"URL,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `ID` | `int` | 参考文献的内部标识符 |
| `Author` | `string` | 作者 |
| `Title` | `string` | 标题（必填，无 `omitempty`） |
| `URL` | `string` | 链接地址 |

::: tip Title 是必填字段
`Title` 的 JSON/XML 标签都没有 `omitempty`，意味着即使为空也会序列化出来。这是 MITRE 数据规范的体现——参考文献至少要有标题。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.References = []cweskills.Reference{
		{
			ID:     1,
			Author: "OWASP",
			Title:  "Cross Site Scripting (XSS) Attack",
			URL:    "https://owasp.org/www-community/attacks/xss/",
		},
		{
			ID:     2,
			Author: "MITRE",
			Title:  "CWE-79 Definition",
			URL:    "https://cwe.mitre.org/data/definitions/79.html",
		},
	}

	for _, r := range cwe.References {
		fmt.Printf("[%d] %s - %s\n  %s\n", r.ID, r.Author, r.Title, r.URL)
	}
}
```

::: details Category 与 View 也带 References
不止 `CWE`，`Category.References` 与 `View.References` 同样使用此类型。详见 [Category 类别](./category)、[View 视图](./view)。
:::

## 🎯 典型用途

<Badge type="tip" text="溯源" /> 在报告中列出弱点的权威资料链接
<Badge type="info" text="扩展阅读" /> 自动生成「了解更多」区块

## ⚠️ 注意事项

::: warning ID 不等于 CWE ID
`Reference.ID` 是参考文献自身的内部编号（XML 中为 `Reference_ID`），**不是** CWE 编号。不要把它与 `CWE.ID` 混淆。
:::

## 🔗 相关链接

- 宿主字段：`CWE.References`、`Category.References`、`View.References`
- 配套的观察示例（含参考编号文本）：[ObservedExample](./observed-example)
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
