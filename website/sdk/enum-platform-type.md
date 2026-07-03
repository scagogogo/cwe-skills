---
title: PlatformType 平台类型枚举
outline: [2, 3]
---

# 📚 PlatformType — 平台类型枚举

`PlatformType` 表示适用平台的分类维度，共 4 个取值：**Language / Operating System / Architecture / Technology**。它对应 [`ApplicablePlatforms`](./applicable-platforms) 结构体的四个字段维度。

## 📋 类型与常量

```go
type PlatformType string

const (
	PlatformLanguage        PlatformType = "Language"
	PlatformOperatingSystem PlatformType = "Operating System"
	PlatformArchitecture    PlatformType = "Architecture"
	PlatformTechnology      PlatformType = "Technology"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 | 对应 ApplicablePlatforms 字段 |
| --- | --- | --- | --- |
| `PlatformLanguage` | `"Language"` | 编程语言 | `Languages` |
| `PlatformOperatingSystem` | `"Operating System"` | 操作系统 | `OperatingSystems` |
| `PlatformArchitecture` | `"Architecture"` | 架构 | `Architectures` |
| `PlatformTechnology` | `"Technology"` | 技术 | `Technologies` |

::: tip 与 ApplicablePlatforms 的对应
`PlatformType` 的四个取值恰好对应 `ApplicablePlatforms` 结构体的四个 `[]PlatformEntry` 字段。它常用于在 UI 或 CLI 里**按维度枚举**平台信息——例如让用户选择「查看语言类平台」还是「查看操作系统类平台」。
:::

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (p PlatformType) String() string` |
| `IsValid` | `func (p PlatformType) IsValid() bool` |
| `ParsePlatformType` | `func ParsePlatformType(s string) (PlatformType, error)` |
| `AllPlatformTypeValues` | `func AllPlatformTypeValues() []PlatformType` |

::: warning Operating System 带空格
`PlatformOperatingSystem` 的值是 `"Operating System"`（带空格）。解析时必须完全匹配。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	pt, err := cweskills.ParsePlatformType("Language")
	fmt.Println(pt, err) // Language <nil>

	// 校验
	fmt.Println(cweskills.PlatformArchitecture.IsValid()) // true
	fmt.Println(cweskills.PlatformType("foo").IsValid())   // false

	// 按维度展示平台
	cwe := cweskills.NewCWE(89, "SQL Injection")
	cwe.ApplicablePlatforms = &cweskills.ApplicablePlatforms{
		Languages: []cweskills.PlatformEntry{
			{Name: "Java", Prevalence: cweskills.PrevalenceOften},
		},
		Technologies: []cweskills.PlatformEntry{
			{Name: "Web Server", Prevalence: cweskills.PrevalenceOften},
		},
	}

	for _, t := range cweskills.AllPlatformTypeValues() {
		var entries []cweskills.PlatformEntry
		switch t {
		case cweskills.PlatformLanguage:
			entries = cwe.ApplicablePlatforms.Languages
		case cweskills.PlatformOperatingSystem:
			entries = cwe.ApplicablePlatforms.OperatingSystems
		case cweskills.PlatformArchitecture:
			entries = cwe.ApplicablePlatforms.Architectures
		case cweskills.PlatformTechnology:
			entries = cwe.ApplicablePlatforms.Technologies
		}
		fmt.Printf("[%s] %d 个平台\n", t, len(entries))
	}
	// [Language] 1 个平台
	// [Operating System] 0 个平台
	// [Architecture] 0 个平台
	// [Technology] 1 个平台
}
```

## 🎯 典型用途

<Badge type="tip" text="维度枚举" /> 在 UI/CLI 按平台维度展示弱点适用范围
<Badge type="info" text="过滤" /> 只关注影响「Language」维度的弱点
<Badge type="warning" text="报告" /> 按维度分块输出适用平台清单

## ⚠️ 注意事项

::: warning PlatformType 本身不存于数据中
`PlatformType` 主要用于**程序内部按维度操作**，MITRE 原始数据里 `ApplicablePlatforms` 直接是四个具名字段，并不存储一个 `PlatformType` 字段。它是工具型枚举，便于遍历与选择。
:::

## 🔗 相关链接

- 字段归宿：[`ApplicablePlatforms`](./applicable-platforms) 的四个维度字段
- 平台条目：`PlatformEntry`（含 `Prevalence`）
- 普遍程度枚举：`Prevalence`（见 `enums.go`）
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
