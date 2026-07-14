---
title: CWE 弱点结构体
outline: [2, 3]
---

# 🧱 CWE — 弱点结构体

`CWE` 是 `cweskills` 包的**核心类型**，描述一条 MITRE CWE 弱点条目的全部信息。它是 SDK 数据模型的聚合根，几乎所有查询、序列化、导航都以它为中心。

## 📋 结构体定义

```go
type CWE struct {
    ID                    int                   `json:"id" xml:"ID,attr"`
    Name                  string                `json:"name" xml:"Name"`
    Abstraction           Abstraction           `json:"abstraction,omitempty" xml:"Abstraction,omitempty"`
    Structure             Structure             `json:"structure,omitempty" xml:"Structure,omitempty"`
    Status                Status                `json:"status,omitempty" xml:"Status,omitempty"`
    Description           string                `json:"description" xml:"Description"`
    ExtendedDescription   string                `json:"extended_description,omitempty" xml:"Extended_Description,omitempty"`
    LikelihoodOfExploit   LikelihoodOfExploit   `json:"likelihood_of_exploit,omitempty" xml:"LikelihoodOfExploit,omitempty"`
    CommonConsequences    []Consequence         `json:"common_consequences,omitempty" xml:"CommonConsequences>Consequence,omitempty"`
    PotentialMitigations  []Mitigation          `json:"potential_mitigations,omitempty" xml:"PotentialMitigations>Mitigation,omitempty"`
    DemonstrativeExamples []DemonstrativeExample `json:"demonstrative_examples,omitempty" xml:"DemonstrativeExamples>DemonstrativeExample,omitempty"`
    ObservedExamples      []ObservedExample     `json:"observed_examples,omitempty" xml:"ObservedExamples>ObservedExample,omitempty"`
    References            []Reference           `json:"references,omitempty" xml:"References>Reference,omitempty"`
    Relationships         []Relationship        `json:"relationships,omitempty" xml:"Relationships>Relationship,omitempty"`
    ApplicablePlatforms   *ApplicablePlatforms  `json:"applicable_platforms,omitempty" xml:"ApplicablePlatforms,omitempty"`
    ModesOfIntroduction   []Introduction        `json:"modes_of_introduction,omitempty" xml:"ModesOfIntroduction>Introduction,omitempty"`
    AlternateTerms        []AlternateTerm       `json:"alternate_terms,omitempty" xml:"AlternateTerms>AlternateTerm,omitempty"`
    Notes                 string                `json:"notes,omitempty" xml:"Notes,omitempty"`
    ContentHistory        *ContentHistory       `json:"content_history,omitempty" xml:"ContentHistory,omitempty"`
    CWEType               string                `json:"cwe_type" xml:"-"`
    URL                   string                `json:"url,omitempty" xml:"-"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `ID` | `int` | 数字标识符，主键 |
| `Name` | `string` | 弱点名称 |
| `Abstraction` | [`Abstraction`](./enum-abstraction) | 抽象层级 Pillar/Class/Base/Variant |
| `Structure` | [`Structure`](./enum-structure) | 结构 Simple/Chain/Composite |
| `Status` | [`Status`](./enum-status) | 状态 Stable/Deprecated 等 |
| `Description` | `string` | 简要描述 |
| `ExtendedDescription` | `string` | 扩展描述 |
| `LikelihoodOfExploit` | [`LikelihoodOfExploit`](./enum-likelihood) | 被利用可能性 |
| `CommonConsequences` | [`[]Consequence`](./consequence) | 常见后果 |
| `PotentialMitigations` | [`[]Mitigation`](./mitigation) | 缓解措施 |
| `DemonstrativeExamples` | [`[]DemonstrativeExample`](./demonstrative-example) | 示范示例 |
| `ObservedExamples` | [`[]ObservedExample`](./observed-example) | 真实观察示例 |
| `References` | [`[]Reference`](./reference) | 参考文献 |
| `Relationships` | `[]Relationship` | 与其他 CWE 的关系 |
| `ApplicablePlatforms` | [`*ApplicablePlatforms`](./applicable-platforms) | 适用平台 |
| `ModesOfIntroduction` | [`[]Introduction`](./introduction) | 引入方式 |
| `AlternateTerms` | [`[]AlternateTerm`](./alternate-term) | 备用术语 |
| `Notes` | `string` | 备注 |
| `ContentHistory` | [`*ContentHistory`](./content-history) | 内容历史 |
| `CWEType` | `string` | 条目类型 `weakness`/`category`/`view`/`compound_element` |
| `URL` | `string` | MITRE 官方页面地址 |

::: tip CWEType 的取值
`CWEType` 由数据来源填充，`NewCWE` 默认设为 `"weakness"`。可用 [类型判断方法](./cwe-type-methods) 区分四种条目。
:::

## 🏗️ 构造器

```go
func NewCWE(id int, name string) *CWE
```

创建一个最小可用的 `CWE`，`CWEType` 默认为 `"weakness"`。

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')")
	cwe.Abstraction = cweskills.AbstractionBase
	cwe.Status = cweskills.StatusStable
	cwe.LikelihoodOfExploit = cweskills.LikelihoodHigh
	cwe.Description = "XSS 弱点"

	fmt.Println(cwe.CWEID())      // CWE-79
	fmt.Println(cwe.IsWeakness()) // true
	fmt.Println(cwe.IsBase())     // true
	fmt.Println(cwe.IsStable())   // true

	// 校验
	if err := cwe.Validate(); err != nil {
		fmt.Println(err)
	}
}
```

## 🔗 方法导航

- 类型判断（IsWeakness/IsCategory/IsPillar...）：[CWE 类型判断方法](./cwe-type-methods)
- 关系获取（GetParentIDs/GetChildIDs...）：[CWE 关系获取方法](./cwe-relationship-methods)
- `CWEID() string`：内部调用 [`FormatCWEIDFromInt`](./format-cwe-id-from-int)
- `Validate() error`：要求 `ID > 0` 且 `Name` 非空

## ⚠️ 注意事项

::: warning Relationships 字段
`Relationships` 存储的是该弱点与其他 CWE 的关联，类型为 `[]Relationship`（定义于 `relationship.go`）。关系导航方法（[GetParentIDs](./cwe-relationship-methods) 等）都基于此字段遍历。
:::

## 🔗 相关链接

- 模型概览：[model.go](./model)
- 子模型逐一文档：[Mitigation](./mitigation)、[Consequence](./consequence) 等
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
