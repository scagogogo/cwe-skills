---
title: Consequence 后果
outline: [2, 3]
---

# ⚠️ Consequence — 后果

`Consequence` 描述某个弱点被利用后可能造成的安全后果。它定义在 `consequences.go`，是 `CWE.CommonConsequences` 字段的元素类型。一个弱点可以有多种后果，每种后果可影响多个安全范围。

## 📋 结构体定义

```go
type Consequence struct {
    Scopes      []ConsequenceScope   `json:"scopes" xml:"Scopes>Scope"`
    Impacts     []ConsequenceImpact  `json:"impacts,omitempty" xml:"Impacts>Impact,omitempty"`
    Likelihood  LikelihoodOfExploit  `json:"likelihood,omitempty" xml:"Likelihood,omitempty"`
    Note        string               `json:"note,omitempty" xml:"Note,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Scopes` | `[]ConsequenceScope` | 受影响的安全范围（机密性/完整性等） |
| `Impacts` | `[]ConsequenceImpact` | 影响严重程度 |
| `Likelihood` | `LikelihoodOfExploit` | 此后果的发生可能性 |
| `Note` | `string` | 补充说明 |

::: tip 为什么 Scopes 和 Impacts 都是切片？
同一后果可能同时影响多个范围（如同时破坏机密性与完整性），每个范围又可能有不同程度。MITRE 用列表表达这种多对多关系。
:::

## 🛠️ 方法

| 方法 | 签名 | 说明 |
| --- | --- | --- |
| `HasScope` | `func (c *Consequence) HasScope(scope ConsequenceScope) bool` | 是否包含某范围 |
| `HasImpact` | `func (c *Consequence) HasImpact(impact ConsequenceImpact) bool` | 是否包含某影响 |
| `MaxImpact` | `func (c *Consequence) MaxImpact() ConsequenceImpact` | 返回最高影响；空列表返回 `ImpactUnknown` |
| `Validate` | `func (c *Consequence) Validate() error` | 校验 `Scopes` 至少一个，否则 `ValidationError` |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.CommonConsequences = []cweskills.Consequence{
		{
			Scopes:     []cweskills.ConsequenceScope{cweskills.ScopeConfidentiality, cweskills.ScopeIntegrity},
			Impacts:    []cweskills.ConsequenceImpact{cweskills.ImpactHigh, cweskills.ImpactMedium},
			Likelihood: cweskills.LikelihoodHigh,
			Note:       "可窃取会话令牌并篡改页面",
		},
	}

	cons := cwe.CommonConsequences[0]
	fmt.Println(cons.HasScope(cweskills.ScopeConfidentiality)) // true
	fmt.Println(cons.HasScope(cweskills.ScopeAvailability))    // false
	fmt.Println(cons.MaxImpact())                              // High

	// CWE 层面的便捷判断
	fmt.Println(cwe.HasConsequenceScope(cweskills.ScopeIntegrity)) // true

	// 校验
	fmt.Println(cons.Validate()) // <nil>
}
```

::: details MaxImpact 的排序依据
`MaxImpact` 用 [`ImpactOrder`](./enum-consequence-impact) 比较：`High=4 > Medium=3 > Low=2 > Unknown=1`。空 `Impacts` 列表返回 `ImpactUnknown`，不会 panic。
:::

## ⚠️ 注意事项

::: warning Scopes 不可空
`Validate` 要求 `Scopes` 至少包含一个元素。一个「没有影响范围」的后果在语义上无意义，会被判为 `ValidationError`。
:::

## 🔗 相关链接

- 宿主字段：`CWE.CommonConsequences`，见 [CWE 弱点](./cwe-struct)
- CWE 便捷方法：`CWE.HasConsequenceScope(scope)`，见 [CWE 关系获取方法](./cwe-relationship-methods)
- 枚举：[ConsequenceScope](./enum-consequence-scope)、[ConsequenceImpact](./enum-consequence-impact)、[LikelihoodOfExploit](./enum-likelihood)
- 源文件：[`consequences.go`](https://github.com/scagogogo/cwe-skills/blob/main/consequences.go)
