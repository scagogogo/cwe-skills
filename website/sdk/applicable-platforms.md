---
title: ApplicablePlatforms 适用平台
outline: [2, 3]
---

# 🖥️ ApplicablePlatforms — 适用平台

`ApplicablePlatforms` 描述 CWE 条目适用的平台范围（语言、操作系统、架构、技术），是 `CWE.ApplicablePlatforms` 字段的类型（指针，可空）。每个维度都是 `[]PlatformEntry` 列表。

## 📋 结构体定义

```go
type ApplicablePlatforms struct {
    Languages        []PlatformEntry `json:"languages,omitempty" xml:"Languages>Language,omitempty"`
    OperatingSystems []PlatformEntry `json:"operating_systems,omitempty" xml:"Operating_Systems>Operating_System,omitempty"`
    Architectures    []PlatformEntry `json:"architectures,omitempty" xml:"Architectures>Architecture,omitempty"`
    Technologies     []PlatformEntry `json:"technologies,omitempty" xml:"Technologies>Technology,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Languages` | `[]PlatformEntry` | 适用的编程语言 |
| `OperatingSystems` | `[]PlatformEntry` | 适用的操作系统 |
| `Architectures` | `[]PlatformEntry` | 适用的架构 |
| `Technologies` | `[]PlatformEntry` | 适用的技术 |

::: tip PlatformEntry 携带普遍程度
每个 `PlatformEntry` 除了名称还带 `Prevalence`（使用普遍程度：Often/Sometimes/Rarely/Undetermined），可用于评估弱点在该平台的常见度。
:::

## 🔧 PlatformEntry

```go
type PlatformEntry struct {
    Name       string     `json:"name" xml:"Name,attr"`
    Prevalence Prevalence `json:"prevalence,omitempty" xml:"Prevalence,attr,omitempty"`
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Name` | `string` | 平台名称（如 "Java"、"Linux"） |
| `Prevalence` | `Prevalence` | 使用普遍程度 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(89, "SQL Injection")
	cwe.ApplicablePlatforms = &cweskills.ApplicablePlatforms{
		Languages: []cweskills.PlatformEntry{
			{Name: "Java", Prevalence: cweskills.PrevalenceOften},
			{Name: "PHP", Prevalence: cweskills.PrevalenceOften},
			{Name: "C#", Prevalence: cweskills.PrevalenceSometimes},
		},
		Technologies: []cweskills.PlatformEntry{
			{Name: "Web Server", Prevalence: cweskills.PrevalenceOften},
		},
	}

	// 列出所有适用语言
	if cwe.ApplicablePlatforms != nil {
		for _, lang := range cwe.ApplicablePlatforms.Languages {
			fmt.Printf("语言: %s (普遍程度: %s)\n", lang.Name, lang.Prevalence)
		}
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="过滤" /> 只关注影响本团队技术栈的弱点
<Badge type="info" text="评估" /> 用 `Prevalence` 估算暴露面
<Badge type="warning" text="报告" /> 在弱点详情里展示「影响哪些平台」

## ⚠️ 注意事项

::: warning 指针类型，判空再访问
`CWE.ApplicablePlatforms` 是 `*ApplicablePlatforms`，很多弱点（尤其是抽象层级较高的）没有平台信息，序列化时为 `nil`。访问前务必判空，否则解引用空指针会 panic。
:::

::: details Prevalence 枚举配套方法
`Prevalence` 同样提供 `String()`、`IsValid()`、`ParsePrevalence()`、`AllPrevalenceValues()`，定义在 `enums.go`。取值为 `Often`、`Sometimes`、`Rarely`、`Undetermined`。
:::

## 🔗 相关链接

- 宿主字段：`CWE.ApplicablePlatforms`，见 [CWE 弱点](./cwe-struct)
- 平台类型枚举：[PlatformType](./enum-platform-type)
- 普遍程度枚举：`Prevalence`（见 `enums.go`）
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
