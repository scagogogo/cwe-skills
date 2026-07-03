---
title: CompoundElement 复合元素
outline: [2, 3]
---

# 🔗 CompoundElement — 复合元素

`CompoundElement` 描述由多个弱点组合而成的复合弱点，包括**链式弱点**（Chain）与**复合弱点**（Composite）。它是 CWE 体系里独立于 `CWE` 的顶层条目类型，通过 `Structure` 字段区分两种组合方式。

## 📋 结构体定义

```go
type CompoundElement struct {
    ID            int           `json:"id" xml:"ID,attr"`
    Name          string        `json:"name" xml:"Name"`
    Structure     Structure     `json:"structure" xml:"Structure"`
    Status        Status        `json:"status,omitempty" xml:"Status,omitempty"`
    Description   string        `json:"description" xml:"Description"`
    Relationships []Relationship `json:"relationships,omitempty" xml:"Relationships>Relationship,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `ID` | `int` | 数字标识符 |
| `Name` | `string` | 复合元素名称 |
| `Structure` | [`Structure`](./enum-structure) | 结构类型（Chain 或 Composite） |
| `Status` | [`Status`](./enum-status) | 状态 |
| `Description` | `string` | 描述 |
| `Relationships` | `[]Relationship` | 组成关系（含 `Requires`/`RequiredBy`） |

::: tip Structure 取值
复合元素的 `Structure` 通常是：
- `StructureChain`：链式，多个弱点**按顺序**可达才产生漏洞（如 CWE-680 整数溢出 → 缓冲区溢出）
- `StructureComposite`：复合，多个弱点**同时存在**才产生漏洞（如 CWE-352 CSRF 需多个弱点并存）

`StructureSimple` 一般不用于复合元素。
:::

## 🏗️ 构造器

```go
func NewCompoundElement(id int, name string, structure Structure) *CompoundElement
```

创建最小可用的 `CompoundElement`，设置 `ID`、`Name`、`Structure`。注意第三个参数是必填的 `Structure`。

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 链式复合元素
	chain := cweskills.NewCompoundElement(680, "Integer Overflow to Buffer Overflow", cweskills.StructureChain)
	chain.Status = cweskills.StatusStable
	chain.Description = "整数溢出导致缓冲区溢出。"
	chain.Relationships = []cweskills.Relationship{
		{CWEID: 190, Nature: cweskills.RelationshipRequires}, // 需要 CWE-190 整数溢出
		{CWEID: 787, Nature: cweskills.RelationshipCanFollow}, // 后接 CWE-787 越界写
	}

	fmt.Printf("%s (structure=%s)\n", chain.Name, chain.Structure)
	for _, r := range chain.Relationships {
		fmt.Printf("  -> CWE-%d (%s)\n", r.CWEID, r.Nature)
	}

	// 复合弱点
	composite := cweskills.NewCompoundElement(352, "CSRF", cweskills.StructureComposite)
	fmt.Println(composite.Structure) // Composite
}
```

## 🎯 典型用途

<Badge type="tip" text="攻击链分析" /> 用 Chain 复合元素还原多步攻击路径
<Badge type="info" text="防御规划" /> 用 Composite 识别需同时加固的弱点组合
<Badge type="warning" text="导航" /> 通过 `Requires`/`RequiredBy` 关系找组成弱点

## ⚠️ 注意事项

::: warning CompoundElement 不是 CWE
`CompoundElement` 是独立结构体，与 `CWE` 平级。不能用 `CWE` 的 `IsCompoundElement()` 方法判断——那是判断 `CWE.CWEType == "compound_element"`。两者数据来源不同，访问字段时注意类型。
:::

::: details 与 CWE.IsChain/IsComposite 的关系
`CWE` 也有 `IsChain()`/`IsComposite()` 方法，但它们判断的是 **`CWE` 自身的 `Structure` 字段**。`CompoundElement` 是另一套结构体，专用于复合弱点条目。在 MITRE 数据里，复合弱点可能以 `CompoundElement` 形式存在，也可能以带 `Structure=Chain` 的 `CWE` 形式存在，取决于数据来源。
:::

## 🔗 相关链接

- 结构枚举：[Structure](./enum-structure)
- 关系类型：[RelationshipNature](./enum-relationship-nature)（`Requires`/`RequiredBy`/`CanPrecede`/`CanFollow`）
- CWE 对应方法：[CWE 类型判断方法](./cwe-type-methods)（`IsChain`/`IsComposite`）
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
