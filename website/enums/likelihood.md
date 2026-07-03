---
title: LikelihoodOfExploit 利用可能性枚举
outline: [2, 3]
---

# ⚡ LikelihoodOfExploit — 利用可能性枚举

`LikelihoodOfExploit` 表示漏洞被利用的可能性高低，是风险评估与优先级排序的关键字段。

## 🧬 类型定义

```go
type LikelihoodOfExploit string
```

## 📋 全部取值

| 值 | 常量名 | 含义 |
| --- | --- | --- |
| `"High"` | `LikelihoodHigh` | 高可能性，攻击者容易利用 |
| `"Medium"` | `LikelihoodMedium` | 中等可能性 |
| `"Low"` | `LikelihoodLow` | 低可能性，利用条件苛刻 |
| `"Unknown"` | `LikelihoodUnknown` | 未知，缺乏足够信息评估 |

```go
const (
	LikelihoodHigh    LikelihoodOfExploit = "High"
	LikelihoodMedium  LikelihoodOfExploit = "Medium"
	LikelihoodLow     LikelihoodOfExploit = "Low"
	LikelihoodUnknown LikelihoodOfExploit = "Unknown"
)
```

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (l LikelihoodOfExploit) String() string` |
| `IsValid` | `func (l LikelihoodOfExploit) IsValid() bool` |
| `ParseLikelihoodOfExploit` | `func ParseLikelihoodOfExploit(s string) (LikelihoodOfExploit, error)` |
| `AllLikelihoodOfExploitValues` | `func AllLikelihoodOfExploitValues() []LikelihoodOfExploit` |

```go
l, err := cweskills.ParseLikelihoodOfExploit("High")
fmt.Println(l, err)                                   // High <nil>
fmt.Println(l.String())                               // High
fmt.Println(cweskills.LikelihoodOfExploit("X").IsValid()) // false
fmt.Println(cweskills.AllLikelihoodOfExploitValues()) // [High Medium Low Unknown]
```

## 📊 额外方法：LikelihoodOrder

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

::: tip 用途
按 `LikelihoodOrder` 降序排列即可得到「从最可能被利用到最不可能」的优先级序列，适合用于风险排序与修复优先级决策。
:::

```go
all := cweskills.AllLikelihoodOfExploitValues()
sort.Slice(all, func(i, j int) bool {
	return all[i].LikelihoodOrder() > all[j].LikelihoodOrder()
})
fmt.Println(all) // [High Medium Low Unknown]
```

::: warning 与 ConsequenceImpact 区分
`LikelihoodOfExploit` 描述「被利用的可能性」，`ConsequenceImpact` 描述「后果的严重程度」。两者独立——一个低可能性但高影响的弱点同样值得关注。综合风险应同时参考两者。
:::

## 💻 CLI 对应命令

```bash
cwe enum likelihood
```

输出全部合法取值，详见 [CLI enum likelihood](../cli/enum-likelihood)。

## 🔗 相关链接

- SDK 视角：[LikelihoodOfExploit 利用可能性枚举](../sdk/enum-likelihood)
- 按可能性筛选：[`FindByLikelihood`](../sdk/find-by-likelihood)
- 按可能性统计：[`CountByLikelihood`](../sdk/count-by-likelihood)
- 后果影响（配套字段）：[ConsequenceImpact](./consequence-impact)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
