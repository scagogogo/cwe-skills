---
title: enums.go 枚举类型概览
outline: [2, 3]
---

# 📚 enums.go — 枚举类型概览

`enums.go` 集中定义了 `cweskills` 包的全部枚举类型。这些枚举贯穿数据模型的各个字段，把 MITRE CWE 规范里的离散取值约束为强类型，避免字符串拼写错误。配套的 `consequences.go` 定义了与后果相关的结构（但其字段类型仍来自这里）。

## 🗺️ 枚举清单

| 枚举类型 | 用途 | 文档 |
| --- | --- | --- |
| `Abstraction` | 抽象层级 Pillar/Class/Base/Variant | [Abstraction](./enum-abstraction) |
| `Structure` | 结构 Simple/Chain/Composite | [Structure](./enum-structure) |
| `Status` | 状态 Stable/Deprecated 等 6 种 | [Status](./enum-status) |
| `LikelihoodOfExploit` | 被利用可能性 High/Medium/Low/Unknown | [LikelihoodOfExploit](./enum-likelihood) |
| `RelationshipNature` | 关系类型 10 种 | [RelationshipNature](./enum-relationship-nature) |
| `ConsequenceScope` | 后果范围 8 种 | [ConsequenceScope](./enum-consequence-scope) |
| `ConsequenceImpact` | 后果影响 High/Medium/Low/Unknown | [ConsequenceImpact](./enum-consequence-impact) |
| `ViewType` | 视图类型 Graph/Explicit Slice/Implicit Slice | [ViewType](./enum-view-type) |
| `PlatformType` | 平台类型 4 种 | [PlatformType](./enum-platform-type) |

::: tip 额外枚举
`enums.go` 还定义了几个不在任务清单但模型用到的枚举：`Prevalence`（平台普遍程度）、`IntroductionPhase`（引入阶段）、`MitigationPhase`（缓解阶段）、`Effectiveness`（缓解有效性）。它们同样遵循下文的统一接口约定。
:::

## 🧩 统一接口约定

每个枚举类型都提供**四件套**，命名规律一致：

| 能力 | 方法/函数 | 示例（以 Abstraction 为例） |
| --- | --- | --- |
| 转字符串 | `(T) String() string` | `AbstractionPillar.String()` → `"Pillar"` |
| 校验 | `(T) IsValid() bool` | `AbstractionPillar.IsValid()` → `true` |
| 解析 | `ParseXxx(s string) (T, error)` | `ParseAbstraction("Pillar")` |
| 全部值 | `AllXxxValues() []T` | `AllAbstractionValues()` |

::: warning 解析失败返回 ValidationError
所有 `ParseXxx` 在字符串不匹配任何合法值时返回 `ValidationError`（见 [ValidationError](./validation-error)），而非 `InvalidCWEIDError`。可用 `errors.As` 取出字段名与非法值。
:::

## 🏷️ 常量命名规则

枚举常量统一采用「**类型缩写 + 值**」的前缀命名，便于 IDE 自动补全与代码检索：

| 类型 | 前缀 | 示例常量 |
| --- | --- | --- |
| `Abstraction` | `Abstraction` | `AbstractionPillar`、`AbstractionBase` |
| `Status` | `Status` | `StatusStable`、`StatusDeprecated` |
| `LikelihoodOfExploit` | `Likelihood` | `LikelihoodHigh`、`LikelihoodMedium` |
| `RelationshipNature` | `Relationship` | `RelationshipChildOf`、`RelationshipPeerOf` |
| `ConsequenceScope` | `Scope` | `ScopeConfidentiality`、`ScopeIntegrity` |
| `ConsequenceImpact` | `Impact` | `ImpactHigh`、`ImpactLow` |
| `ViewType` | `ViewType` | `ViewTypeGraph` |
| `PlatformType` | `Platform` | `PlatformLanguage` |

## 📊 排序权重方法

部分枚举额外提供 `XxxOrder() int`，返回可比较的权重，用于排序：

| 方法 | 权重（高→低） |
| --- | --- |
| `Abstraction.AbstractionOrder()` | Pillar=4, Class=3, Base=2, Variant=1 |
| `LikelihoodOfExploit.LikelihoodOrder()` | High=4, Medium=3, Low=2, Unknown=1 |
| `ConsequenceImpact.ImpactOrder()` | High=4, Medium=3, Low=2, Unknown=1 |

::: tip 用途
这些 `Order` 方法被 [`Consequence.MaxImpact()`](./consequence) 等方法用于求最高影响，也可直接传入 `sort.Slice` 做自定义排序。
:::

## 🚀 快速上手

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	a, err := cweskills.ParseAbstraction("Base")
	fmt.Println(a, err) // Base <nil>

	// 遍历全部值
	for _, s := range cweskills.AllStatusValues() {
		fmt.Println(s)
	}

	// 排序权重
	fmt.Println(cweskills.AbstractionBase.AbstractionOrder()) // 2
	fmt.Println(cweskills.LikelihoodHigh.LikelihoodOrder())   // 4
}
```

## 🔗 相关链接

- 数据模型（枚举的字段归宿）：[model.go](./model)
- 错误类型：[ValidationError](./validation-error)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
- CLI 枚举命令：[enum 总览](../cli/enum)
