---
title: ConsequenceScope 后果范围枚举
outline: [2, 3]
---

# 📚 ConsequenceScope — 后果范围枚举

`ConsequenceScope` 表示弱点被利用后受影响的安全范围（CIA 三元组及其扩展），共 8 个取值。它是 `Consequence.Scopes` 字段的元素类型。

## 📋 类型与常量

```go
type ConsequenceScope string

const (
	ScopeConfidentiality ConsequenceScope = "Confidentiality"
	ScopeIntegrity       ConsequenceScope = "Integrity"
	ScopeAvailability    ConsequenceScope = "Availability"
	ScopeAccessControl   ConsequenceScope = "Access Control"
	ScopeAccountability  ConsequenceScope = "Accountability"
	ScopeAuthentication  ConsequenceScope = "Authentication"
	ScopeAuthorization   ConsequenceScope = "Authorization"
	ScopeNonRepudiation  ConsequenceScope = "Non-Repudiation"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 |
| --- | --- | --- |
| `ScopeConfidentiality` | `"Confidentiality"` | 机密性——信息被未授权读取 |
| `ScopeIntegrity` | `"Integrity"` | 完整性——信息被未授权篡改 |
| `ScopeAvailability` | `"Availability"` | 可用性——服务被中断 |
| `ScopeAccessControl` | `"Access Control"` | 访问控制——权限被绕过 |
| `ScopeAccountability` | `"Accountability"` | 可追责性——行为无法追溯 |
| `ScopeAuthentication` | `"Authentication"` | 认证——身份验证被攻破 |
| `ScopeAuthorization` | `"Authorization"` | 授权——权限授予被滥用 |
| `ScopeNonRepudiation` | `"Non-Repudiation"` | 不可否认性——行为可被否认 |

::: tip CIA 三元组
前三个（机密性/完整性/可用性）即经典 CIA 三元组，是安全后果的核心维度。后五个是扩展范围，覆盖访问控制、认证授权等。
:::

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (s ConsequenceScope) String() string` |
| `IsValid` | `func (s ConsequenceScope) IsValid() bool` |
| `ParseConsequenceScope` | `func ParseConsequenceScope(s string) (ConsequenceScope, error)` |
| `AllConsequenceScopeValues` | `func AllConsequenceScopeValues() []ConsequenceScope` |

::: warning 两个值带空格/连字符
`ScopeAccessControl` 的值是 `"Access Control"`（带空格），`ScopeNonRepudiation` 是 `"Non-Repudiation"`（带连字符）。解析时大小写与符号必须完全匹配。
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
	sc, err := cweskills.ParseConsequenceScope("Confidentiality")
	fmt.Println(sc, err) // Confidentiality <nil>

	// 校验
	fmt.Println(cweskills.ScopeIntegrity.IsValid()) // true
	fmt.Println(cweskills.ConsequenceScope("foo").IsValid()) // false

	// 判断后果是否影响机密性
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.CommonConsequences = []cweskills.Consequence{
		{Scopes: []cweskills.ConsequenceScope{
			cweskills.ScopeConfidentiality,
			cweskills.ScopeIntegrity,
		}},
	}
	fmt.Println(cwe.HasConsequenceScope(cweskills.ScopeConfidentiality)) // true
	fmt.Println(cwe.HasConsequenceScope(cweskills.ScopeAvailability))    // false
}
```

## 🎯 典型用途

<Badge type="tip" text="风险分类" /> 按受影响范围归类弱点（机密性破坏类、可用性破坏类）
<Badge type="info" text="过滤" /> [`FindByConsequenceScope`](./find-by-consequence-scope) 找出影响某范围的弱点
<Badge type="warning" text="统计" /> [`CountByScope`](./count-by-scope) 统计各范围分布

## ⚠️ 注意事项

::: warning 一条后果可影响多个范围
`Consequence.Scopes` 是切片，一条后果可同时列入多个范围（如同时破坏机密性与完整性）。判断时用 `Consequence.HasScope(scope)` 或 `CWE.HasConsequenceScope(scope)`，不要假设只有一个。
:::

## 🔗 相关链接

- 字段归宿：`Consequence.Scopes`
- 判断方法：`Consequence.HasScope()`、`CWE.HasConsequenceScope()`，见 [Consequence](./consequence)
- 过滤 API：[FindByConsequenceScope](./find-by-consequence-scope)
- 概念背景：[后果范围与影响](../guide/concept-consequence)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
