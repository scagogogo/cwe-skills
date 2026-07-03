---
title: AlternateTerm 备用术语
outline: [2, 3]
---

# 🔁 AlternateTerm — 备用术语

`AlternateTerm` 提供 CWE 条目的其他常用名称或术语，是 `CWE.AlternateTerms` 字段的元素类型。安全社区对同一弱点常有多种叫法，备用术语帮助建立同义词映射。

## 📋 结构体定义

```go
type AlternateTerm struct {
    Term        string `json:"term" xml:"Term"`
    Description string `json:"description,omitempty" xml:"Description,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Term` | `string` | 备用术语名称 |
| `Description` | `string` | 备用术语的描述/适用语境 |

::: tip Term 是必填字段
`Term` 的标签无 `omitempty`，序列化时一定会输出。一个没有术语文本的 `AlternateTerm` 是无意义的。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')")
	cwe.AlternateTerms = []cweskills.AlternateTerm{
		{
			Term:        "XSS",
			Description: "业界最常用的缩写。",
		},
		{
			Term:        "Cross Site Scripting",
			Description: "完整写法，历史沿用。",
		},
		{
			Term:        "CSS",
			Description: "早期缩写，因与层叠样式表冲突后改称 XSS。",
		},
	}

	fmt.Printf("%s 的备用术语：\n", cwe.Name)
	for _, at := range cwe.AlternateTerms {
		fmt.Printf("- %s：%s\n", at.Term, at.Description)
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="搜索" /> 把备用术语纳入关键词索引，提升弱点检索召回率
<Badge type="info" text="别名" /> 在 UI 上展示「又称作 ...」
<Badge type="warning" text="标准化" /> 把社区口语映射回官方 CWE 编号

## ⚠️ 注意事项

::: warning 别与 Name 混淆
`CWE.Name` 是官方正式名称（通常很长且学术化），`AlternateTerm.Term` 是社区/历史别名。两者不能互换——正式报告应使用 `Name`，搜索/展示可借助 `AlternateTerms`。
:::

## 🔗 相关链接

- 宿主字段：`CWE.AlternateTerms`，见 [CWE 弱点](./cwe-struct)
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
