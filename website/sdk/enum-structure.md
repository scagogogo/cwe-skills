---
title: Structure 结构类型枚举
outline: [2, 3]
---

# 📚 Structure — 结构类型枚举

`Structure` 表示 CWE 条目的结构类型，描述弱点之间的组成关系。共有三种：**Simple**（单一）、**Chain**（链式）、**Composite**（复合）。

## 📋 类型与常量

```go
type Structure string

const (
	StructureSimple    Structure = "Simple"
	StructureChain     Structure = "Chain"
	StructureComposite Structure = "Composite"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 | 示例 |
| --- | --- | --- | --- |
| `StructureSimple` | `"Simple"` | 单一弱点，不依赖其他弱点的存在 | 大多数普通弱点 |
| `StructureChain` | `"Chain"` | 链式，多个弱点**按顺序**可达才产生漏洞 | CWE-680 整数溢出→缓冲区溢出 |
| `StructureComposite` | `"Composite"` | 复合，多个弱点**同时存在**才产生漏洞 | CWE-352 CSRF 需多弱点并存 |

::: tip Chain vs Composite 的区别
- **Chain**：弱点 A 发生后，为弱点 B 创造条件，A→B 有先后顺序。
- **Composite**：弱点 A 与 B 必须同时存在，漏洞才会显现，无先后。

这是 CWE 区分复合弱点的核心维度。
:::

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (s Structure) String() string` |
| `IsValid` | `func (s Structure) IsValid() bool` |
| `ParseStructure` | `func ParseStructure(s string) (Structure, error)` |
| `AllStructureValues` | `func AllStructureValues() []Structure` |

::: warning Structure 没有 Order 方法
与 `Abstraction`/`Likelihood`/`Impact` 不同，`Structure` **没有** `StructureOrder()` 方法——三种结构类型之间没有自然的数值序关系。
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
	s, err := cweskills.ParseStructure("Chain")
	fmt.Println(s, err) // Chain <nil>

	// 校验
	fmt.Println(cweskills.StructureSimple.IsValid())    // true
	fmt.Println(cweskills.Structure("Foo").IsValid())   // false

	// 全部值
	for _, v := range cweskills.AllStructureValues() {
		fmt.Println(v)
	}
	// Simple
	// Chain
	// Composite
}
```

## 🎯 典型用途

<Badge type="tip" text="分类" /> 区分单一弱点与复合弱点
<Badge type="info" text="过滤" /> [`FindByStructure`](./find-by-structure) / [`FindChains`](./find-chains-composites)
<Badge type="warning" text="复合元素" /> [`NewCompoundElement`](./compound-element) 用它指定是 Chain 还是 Composite

## ⚠️ 注意事项

::: warning StructureSimple 用于复合元素无意义
`CompoundElement.Structure` 通常取 `Chain` 或 `Composite`。`Simple` 一般不用于复合元素——「简单」意味着不依赖其他弱点，本就不是复合。但类型上不强制禁止。
:::

## 🔗 相关链接

- 字段归宿：`CWE.Structure`、`CompoundElement.Structure`
- 便捷判断方法：`CWE.IsChain()`/`IsComposite()`，见 [CWE 类型判断方法](./cwe-type-methods)
- 复合元素结构体：[CompoundElement](./compound-element)
- 概念背景：[结构类型 (Simple/Chain/Composite)](../guide/concept-structure)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
