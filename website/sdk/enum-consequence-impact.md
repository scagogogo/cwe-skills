---
title: ConsequenceImpact 后果影响枚举
outline: [2, 3]
---

# 📚 ConsequenceImpact — 后果影响枚举

`ConsequenceImpact` 表示后果的严重程度，共 4 个取值：**High / Medium / Low / Unknown**。它是 `Consequence.Impacts` 字段的元素类型。

## 📋 类型与常量

```go
type ConsequenceImpact string

const (
	ImpactHigh    ConsequenceImpact = "High"
	ImpactMedium  ConsequenceImpact = "Medium"
	ImpactLow     ConsequenceImpact = "Low"
	ImpactUnknown ConsequenceImpact = "Unknown"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 |
| --- | --- | --- |
| `ImpactHigh` | `"High"` | 严重影响 |
| `ImpactMedium` | `"Medium"` | 中等影响 |
| `ImpactLow` | `"Low"` | 较小影响 |
| `ImpactUnknown` | `"Unknown"` | 未知影响 |

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (i ConsequenceImpact) String() string` |
| `IsValid` | `func (i ConsequenceImpact) IsValid() bool` |
| `ParseConsequenceImpact` | `func ParseConsequenceImpact(s string) (ConsequenceImpact, error)` |
| `AllConsequenceImpactValues` | `func AllConsequenceImpactValues() []ConsequenceImpact` |

## 📊 排序权重

```go
func (i ConsequenceImpact) ImpactOrder() int
```

返回排序权重，**影响越严重值越大**：

| 取值 | 权重 |
| --- | --- |
| `ImpactHigh` | 4 |
| `ImpactMedium` | 3 |
| `ImpactLow` | 2 |
| `ImpactUnknown` | 1 |
| 未知 | 0 |

::: tip 用途
[`Consequence.MaxImpact()`](./consequence) 正是用 `ImpactOrder` 找出一条后果里最严重的影响。也可用于把多条后果按严重度排序。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	imp, err := cweskills.ParseConsequenceImpact("High")
	fmt.Println(imp, err) // High <nil>

	// 排序权重
	fmt.Println(cweskills.ImpactHigh.ImpactOrder())   // 4
	fmt.Println(cweskills.ImpactUnknown.ImpactOrder()) // 1

	// 取一条后果的最高影响
	cons := cweskills.Consequence{
		Impacts: []cweskills.ConsequenceImpact{
			cweskills.ImpactLow,
			cweskills.ImpactHigh,
			cweskills.ImpactMedium,
		},
	}
	fmt.Println(cons.MaxImpact()) // High
}
```

## 🎯 典型用途

<Badge type="tip" text="风险评级" /> 取 `MaxImpact` 作为后果严重度
<Badge type="info" text="排序" /> 用 `ImpactOrder` 把高影响后果排前面
<Badge type="warning" text="过滤" /> 筛出影响达到 High 的弱点

## ⚠️ 注意事项

::: warning 与 LikelihoodOfExploit 取值相同但类型不同
`ConsequenceImpact` 与 [`LikelihoodOfExploit`](./enum-likelihood) 都有 High/Medium/Low/Unknown，`Order` 权重也一致，但类型不同，不可互换。语义：本枚举是「后果严重程度」，后者是「发生可能性」。风险 = 可能性 × 严重度，两枚举常配合使用。
:::

::: details MaxImpact 对空列表的处理
`Consequence.MaxImpact()` 在 `Impacts` 为空时返回 `ImpactUnknown`（而非 panic 或零值）。这意味着没有影响信息的后果会被视为「未知严重度」，下游可据此决定是否降级处理。
:::

## 🔗 相关链接

- 字段归宿：`Consequence.Impacts`
- 取最高影响：`Consequence.MaxImpact()`，见 [Consequence](./consequence)
- 配套可能性枚举：[LikelihoodOfExploit](./enum-likelihood)
- 概念背景：[后果范围与影响](../guide/concept-consequence)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
