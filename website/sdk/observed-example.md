---
title: ObservedExample 观察到的示例
outline: [2, 3]
---

# 🔍 ObservedExample — 观察到的示例

`ObservedExample` 记录来自真实漏洞报告或安全事件的示例，是 `CWE.ObservedExamples` 字段的元素类型。与教学性质的 [`DemonstrativeExample`](./demonstrative-example) 不同，它对应**真实发生过的**案例。

## 📋 结构体定义

```go
type ObservedExample struct {
    Reference   string `json:"reference,omitempty" xml:"Reference,omitempty"`
    Description string `json:"description" xml:"Description"`
    Link        string `json:"link,omitempty" xml:"Link,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Reference` | `string` | 参考编号（如 CVE 编号或厂商公告 ID） |
| `Description` | `string` | 示例描述 |
| `Link` | `string` | 相关链接（URL） |

::: tip 与 Reference 类型的区别
这里的 `Reference` 是**字符串字段**（参考编号文本），不是 [`Reference`](./reference) 结构体。后者是独立的参考文献对象，挂在 `CWE.References` 上。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.ObservedExamples = []cweskills.ObservedExample{
		{
			Reference:   "CVE-2020-1234",
			Description: "某 CMS 的评论功能存在存储型 XSS，攻击者可注入恶意脚本。",
			Link:        "https://nvd.nist.gov/vuln/detail/CVE-2020-1234",
		},
		{
			Reference:   "CVE-2021-5678",
			Description: "某论坛个人签名栏未做输出编码，导致反射型 XSS。",
			Link:        "https://nvd.nist.gov/vuln/detail/CVE-2021-5678",
		},
	}

	for _, ex := range cwe.ObservedExamples {
		fmt.Printf("[%s] %s\n  %s\n", ex.Reference, ex.Description, ex.Link)
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="举证" /> 在风险评估中引用真实漏洞佐证弱点严重性
<Badge type="info" text="关联" /> 通过 `Reference` 字段与 CVE 库做交叉引用

## ⚠️ 注意事项

::: warning Reference 字段格式不固定
`Reference` 是自由文本，可能是 `CVE-XXXX-XXXX`、厂商公告编号或自定义标识。解析时不要假设固定模式，如需提取 CVE 编号可用正则单独处理。
:::

## 🔗 相关链接

- 宿主字段：`CWE.ObservedExamples`，见 [CWE 弱点](./cwe-struct)
- 教学示例（对比）：[DemonstrativeExample](./demonstrative-example)
- 参考文献（结构体）：[Reference](./reference)
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
