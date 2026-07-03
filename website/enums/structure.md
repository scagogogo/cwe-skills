---
title: Structure 结构类型枚举
outline: [2, 3]
---

# 🏷️ Structure — 结构类型枚举

`Structure` 表示 CWE 条目的结构类型，描述弱点之间的组成关系：单一弱点、链式弱点或复合弱点。

## 🧬 类型定义

```go
type Structure string
```

## 📋 全部取值

| 值 | 常量名 | 含义 | 示例 CWE |
| --- | --- | --- | --- |
| `"Simple"` | `StructureSimple` | 简单弱点，不依赖其他弱点的存在 | 大多数独立 CWE |
| `"Chain"` | `StructureChain` | 链式弱点，弱点必须按顺序可达才产生漏洞 | CWE-680 整数溢出到缓冲区溢出 |
| `"Composite"` | `StructureComposite` | 复合弱点，多个弱点必须同时存在 | CWE-352 CSRF |

```go
const (
	StructureSimple    Structure = "Simple"
	StructureChain     Structure = "Chain"
	StructureComposite Structure = "Composite"
)
```

::: tip Chain 与 Composite 的区别
**Chain**（链）要求弱点**按顺序**依次触发——前者为后者创造条件；**Composite**（复合）要求多个弱点**同时存在**——任一缺失则漏洞不成立。两者结构不同，分析方法也不同。
:::

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (s Structure) String() string` |
| `IsValid` | `func (s Structure) IsValid() bool` |
| `ParseStructure` | `func ParseStructure(s string) (Structure, error)` |
| `AllStructureValues` | `func AllStructureValues() []Structure` |

```go
s, err := cweskills.ParseStructure("Chain")
fmt.Println(s, err)                                // Chain <nil>
fmt.Println(s.String())                            // Chain
fmt.Println(cweskills.Structure("Foo").IsValid())  // false
fmt.Println(cweskills.AllStructureValues())        // [Simple Chain Composite]
```

::: warning 无额外方法
`Structure` 不提供 `Order()` 之类的排序权重——Simple/Chain/Composite 之间没有天然的严重程度或层级顺序，不应假定排序语义。
:::

## 🔄 典型用法

```go
// 找出所有链式与复合弱点，做依赖分析
chains, _ := registry.FindByStructure(cweskills.StructureChain)
composites, _ := registry.FindByStructure(cweskills.StructureComposite)
fmt.Println("链式:", len(chains), "复合:", len(composites))
```

## 💻 CLI 对应命令

```bash
cwe enum structure
```

输出全部合法取值，详见 [CLI enum structure](../cli/enum-structure)。

## 🔗 相关链接

- SDK 视角：[Structure 结构类型枚举](../sdk/enum-structure)
- 概念背景：[结构类型 (Simple/Chain/Composite)](../guide/concept-structure)
- 按结构筛选：[`FindByStructure`](../sdk/find-by-structure)
- 查找链/复合：[`FindChainsComposites`](../sdk/find-chains-composites)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
