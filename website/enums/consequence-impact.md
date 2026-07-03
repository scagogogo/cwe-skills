---
title: ConsequenceImpact 后果影响枚举
outline: [2, 3]
---

# 📊 ConsequenceImpact — 后果影响枚举

`ConsequenceImpact` 表示弱点被利用后所造成后果的严重程度，与 `ConsequenceScope`（影响范围）共同描述一个后果。

## 🧬 类型定义

```go
type ConsequenceImpact string
```

## 📋 全部取值

| 值 | 常量名 | 含义 |
| --- | --- | --- |
| `"High"` | `ImpactHigh` | 高影响，造成严重后果 |
| `"Medium"` | `ImpactMedium` | 中等影响 |
| `"Low"` | `ImpactLow` | 低影响，后果有限 |
| `"Unknown"` | `ImpactUnknown` | 未知，缺乏足够信息评估 |

```go
const (
	ImpactHigh    ConsequenceImpact = "High"
	ImpactMedium  ConsequenceImpact = "Medium"
	ImpactLow     ConsequenceImpact = "Low"
	ImpactUnknown ConsequenceImpact = "Unknown"
)
```

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (i ConsequenceImpact) String() string` |
| `IsValid` | `func (i ConsequenceImpact) IsValid() bool` |
| `ParseConsequenceImpact` | `func ParseConsequenceImpact(s string) (ConsequenceImpact, error)` |
| `AllConsequenceImpactValues` | `func AllConsequenceImpactValues() []ConsequenceImpact` |

```go
i, err := cweskills.ParseConsequenceImpact("High")
fmt.Println(i, err)                                  // High <nil>
fmt.Println(i.String())                              // High
fmt.Println(cweskills.ConsequenceImpact("X").IsValid()) // false
fmt.Println(cweskills.AllConsequenceImpactValues()) // [High Medium Low Unknown]
```

## 📊 额外方法：ImpactOrder

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

```go
all := cweskills.AllConsequenceImpactValues()
sort.Slice(all, func(i, j int) bool {
	return all[i].ImpactOrder() > all[j].ImpactOrder()
})
fmt.Println(all) // [High Medium Low Unknown]
```

::: tip 用途
按 `ImpactOrder` 降序排列得到「从最严重到最轻微」的后果序列，与 `LikelihoodOfExploit` 的排序权重配合，可计算综合风险优先级。
:::

::: warning 与 LikelihoodOfExploit 同形但不同义
`ConsequenceImpact` 与 `LikelihoodOfExploit` 的取值集合与 `Order()` 权重完全同形（High/Medium/Low/Unknown，权重 4/3/2/1），但语义不同：前者是「后果多严重」，后者是「多容易被利用」。请勿混用方法名——`ImpactOrder()` 与 `LikelihoodOrder()` 不可互换。
:::

## 🔄 典型用法

```go
// 取出某弱点所有后果中影响最高者，用于风险评估
var maxImpact cweskills.ConsequenceImpact
for _, c := range cwe.Consequences {
	if c.Impact.ImpactOrder() > maxImpact.ImpactOrder() {
		maxImpact = c.Impact
	}
}
fmt.Println("最高影响:", maxImpact)
```

## 💻 CLI 对应命令

```bash
cwe enum consequence-impact
```

输出全部合法取值，详见 [CLI enum consequence-impact](../cli/enum-consequence-impact)。

## 🔗 相关链接

- SDK 视角：[ConsequenceImpact 后果影响枚举](../sdk/enum-consequence-impact)
- 概念背景：[后果范围与影响](../guide/concept-consequence)
- 后果范围（配套字段）：[ConsequenceScope](./consequence-scope)
- 利用可能性（同形枚举）：[LikelihoodOfExploit](./likelihood)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
