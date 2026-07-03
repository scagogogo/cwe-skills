---
title: LikelihoodOfExploit 利用可能性枚举
outline: [2, 3]
---

# 📚 LikelihoodOfExploit — 利用可能性枚举

`LikelihoodOfExploit` 表示某个弱点被攻击者利用的可能性，共 4 个取值：**High / Medium / Low / Unknown**。它出现在 `CWE.LikelihoodOfExploit` 与 `Consequence.Likelihood` 两个字段。

## 📋 类型与常量

```go
type LikelihoodOfExploit string

const (
	LikelihoodHigh    LikelihoodOfExploit = "High"
	LikelihoodMedium  LikelihoodOfExploit = "Medium"
	LikelihoodLow     LikelihoodOfExploit = "Low"
	LikelihoodUnknown LikelihoodOfExploit = "Unknown"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 |
| --- | --- | --- |
| `LikelihoodHigh` | `"High"` | 高可能性被利用 |
| `LikelihoodMedium` | `"Medium"` | 中等可能性 |
| `LikelihoodLow` | `"Low"` | 低可能性 |
| `LikelihoodUnknown` | `"Unknown"` | 未知 |

::: tip 两个使用场景
1. `CWE.LikelihoodOfExploit`：弱点**整体**被利用的可能性。
2. `Consequence.Likelihood`：某个**具体后果**的发生可能性。

同一个弱点可能整体利用可能性低，但某条后果一旦发生危害极高——两字段配合刻画风险。
:::

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (l LikelihoodOfExploit) String() string` |
| `IsValid` | `func (l LikelihoodOfExploit) IsValid() bool` |
| `ParseLikelihoodOfExploit` | `func ParseLikelihoodOfExploit(s string) (LikelihoodOfExploit, error)` |
| `AllLikelihoodOfExploitValues` | `func AllLikelihoodOfExploitValues() []LikelihoodOfExploit` |

## 📊 排序权重

```go
func (l LikelihoodOfExploit) LikelihoodOrder() int
```

返回排序权重，**可能性越高值越大**：

| 取值 | 权重 |
| --- | --- |
| `LikelihoodHigh` | 4 |
| `LikelihoodMedium` | 3 |
| `LikelihoodLow` | 2 |
| `LikelihoodUnknown` | 1 |
| 未知 | 0 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"sort"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	l, err := cweskills.ParseLikelihoodOfExploit("High")
	fmt.Println(l, err) // High <nil>

	// 排序权重
	fmt.Println(cweskills.LikelihoodHigh.LikelihoodOrder())   // 4
	fmt.Println(cweskills.LikelihoodUnknown.LikelihoodOrder()) // 1

	// 按利用可能性降序排弱点
	cwes := []*cweskills.CWE{
		{LikelihoodOfExploit: cweskills.LikelihoodLow},
		{LikelihoodOfExploit: cweskills.LikelihoodHigh},
		{LikelihoodOfExploit: cweskills.LikelihoodMedium},
	}
	sort.Slice(cwes, func(i, j int) bool {
		return cwes[i].LikelihoodOfExploit.LikelihoodOrder() >
			cwes[j].LikelihoodOfExploit.LikelihoodOrder()
	})
	// 顺序: High -> Medium -> Low
}
```

## 🎯 典型用途

<Badge type="tip" text="风险评估" /> 用 `Likelihood × Impact` 估算风险
<Badge type="info" text="过滤" /> [`FindByLikelihood`](./find-by-likelihood) 只看高可能性弱点
<Badge type="warning" text="排序" /> 用 `LikelihoodOrder` 把高风险弱点排前面

## ⚠️ 注意事项

::: warning 与 ConsequenceImpact 取值相同但类型不同
`LikelihoodOfExploit` 与 [`ConsequenceImpact`](./enum-consequence-impact) 都有 High/Medium/Low/Unknown 四个值，且 `Order` 方法权重也相同，但它们是**不同的类型**，不能直接赋值或比较。语义上：前者是「发生可能性」，后者是「后果严重程度」。
:::

## 🔗 相关链接

- 字段归宿：`CWE.LikelihoodOfExploit`、`Consequence.Likelihood`
- 过滤 API：[FindByLikelihood](./find-by-likelihood)
- 统计 API：[CountByLikelihood](./count-by-likelihood)
- 配套影响枚举：[ConsequenceImpact](./enum-consequence-impact)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
