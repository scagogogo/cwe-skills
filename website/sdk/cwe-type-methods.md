---
title: CWE 类型判断方法
outline: [2, 3]
---

# 🧱 CWE 类型判断方法

`CWE` 结构体提供一组零参数的布尔方法，用于快速判断条目的**类型**与**抽象/结构/状态**层级。它们都是简单的字段比较，无副作用，可放心用于过滤与条件分支。

## 📋 方法清单

### 条目类型判断（按 `CWEType` 字段）

| 方法 | 签名 | 为真条件 |
| --- | --- | --- |
| `IsWeakness` | `func (c *CWE) IsWeakness() bool` | `CWEType == "weakness"` |
| `IsCategory` | `func (c *CWE) IsCategory() bool` | `CWEType == "category"` |
| `IsView` | `func (c *CWE) IsView() bool` | `CWEType == "view"` |
| `IsCompoundElement` | `func (c *CWE) IsCompoundElement() bool` | `CWEType == "compound_element"` |

### 抽象层级判断（按 `Abstraction` 字段）

| 方法 | 为真条件 |
| --- | --- |
| `IsPillar() bool` | `Abstraction == AbstractionPillar` |
| `IsBase() bool` | `Abstraction == AbstractionBase` |
| `IsVariant() bool` | `Abstraction == AbstractionVariant` |

::: tip 缺少的 IsClass
SDK 没有提供 `IsClass()` 便捷方法。若需判断 Class 层级，直接比较：`cwe.Abstraction == cweskills.AbstractionClass`。
:::

### 结构判断（按 `Structure` 字段）

| 方法 | 为真条件 |
| --- | --- |
| `IsChain() bool` | `Structure == StructureChain` |
| `IsComposite() bool` | `Structure == StructureComposite` |

### 状态判断（按 `Status` 字段）

| 方法 | 为真条件 |
| --- | --- |
| `IsStable() bool` | `Status == StatusStable` |
| `IsDeprecated() bool` | `Status == StatusDeprecated` |

::: warning 状态判断不完整
仅提供 `IsStable` 与 `IsDeprecated` 两个便捷方法。其余状态（Usable/Draft/Incomplete/Obsolete）需直接比较 `cwe.Status` 字段，例如：`cwe.Status == cweskills.StatusDraft`。
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
	cwe.Abstraction = cweskills.AbstractionBase
	cwe.Structure = cweskills.StructureSimple
	cwe.Status = cweskills.StatusStable

	fmt.Println(cwe.IsWeakness())      // true
	fmt.Println(cwe.IsCategory())      // false
	fmt.Println(cwe.IsBase())          // true
	fmt.Println(cwe.IsVariant())       // false
	fmt.Println(cwe.IsChain())         // false
	fmt.Println(cwe.IsStable())        // true
	fmt.Println(cwe.IsDeprecated())    // false

	// 过滤：只保留稳定的基础弱点
	cwes := []*cweskills.CWE{cwe}
	stableBase := make([]*cweskills.CWE, 0)
	for _, c := range cwes {
		if c.IsStable() && c.IsBase() {
			stableBase = append(stableBase, c)
		}
	}
	fmt.Println(len(stableBase)) // 1
}
```

## 🎯 典型用途

<Badge type="tip" text="过滤" /> 在 [搜索/过滤](./filter) 中按抽象层级筛选弱点
<Badge type="info" text="分类" /> 渲染报告时按类型分桶展示
<Badge type="warning" text="守卫" /> 跳过已废弃条目：`if cwe.IsDeprecated() { continue }`

## 🔗 相关链接

- 字段定义来源：[CWE 弱点结构体](./cwe-struct)
- 关系获取方法：[CWE 关系获取方法](./cwe-relationship-methods)
- 枚举取值：[Abstraction](./enum-abstraction)、[Structure](./enum-structure)、[Status](./enum-status)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
