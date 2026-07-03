---
title: DemonstrativeExample 示范性示例
outline: [2, 3]
---

# 💡 DemonstrativeExample — 示范性示例

`DemonstrativeExample` 展示某个弱点在实际代码或场景中可能出现的情况，是 `CWE.DemonstrativeExamples` 字段的元素类型。它通常包含一段介绍文本与主体内容，帮助理解弱点的具体形态。

## 📋 结构体定义

```go
type DemonstrativeExample struct {
    IntroText string `json:"intro_text,omitempty" xml:"IntroText,omitempty"`
    BodyText  string `json:"body_text,omitempty" xml:"BodyText,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `IntroText` | `string` | 示例的介绍/背景文本 |
| `BodyText` | `string` | 示例的主体内容（常含代码片段） |

::: tip 内容形态
`BodyText` 在 MITRE 原始数据中通常是 HTML 片段（含 `<code>`、`<p>` 等标签），因为是直接从 XML 转储而来。渲染时需要做 HTML 清洗或转换。
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
	cwe.DemonstrativeExamples = []cweskills.DemonstrativeExample{
		{
			IntroText: "以下代码直接把用户输入写回页面：",
			BodyText:  "<code>print("&lt;p&gt;" + input + "&lt;/p&gt;")</code>",
		},
	}

	for i, ex := range cwe.DemonstrativeExamples {
		fmt.Printf("=== 示例 %d ===\n", i+1)
		fmt.Println(ex.IntroText)
		fmt.Println(ex.BodyText)
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="教学" /> 在文档/培训材料里展示弱点的真实样貌
<Badge type="info" text="检测" /> 提取 `BodyText` 中的代码模式作为静态扫描规则参考

## ⚠️ 注意事项

::: warning BodyText 可能含 HTML
MITRE 原始示例的 `BodyText` 常带 HTML 标签。直接拼接进网页前应做转义/清洗，避免引入存储型 XSS——这正是 CWE-79 所描述的弱点，颇具讽刺。
:::

## 🔗 相关链接

- 宿主字段：`CWE.DemonstrativeExamples`，见 [CWE 弱点](./cwe-struct)
- 配套的真实示例：[ObservedExample](./observed-example)
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
