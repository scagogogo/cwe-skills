---
title: ConsequenceScope 后果范围枚举
outline: [2, 3]
---

# 🎯 ConsequenceScope — 后果范围枚举

`ConsequenceScope` 表示弱点被利用后对系统造成后果的影响范围，对应经典的安全三性及其扩展属性。

## 🧬 类型定义

```go
type ConsequenceScope string
```

## 📋 全部取值

| 值 | 常量名 | 含义 |
| --- | --- | --- |
| `"Confidentiality"` | `ScopeConfidentiality` | 机密性，信息泄露给未授权方 |
| `"Integrity"` | `ScopeIntegrity` | 完整性，数据被未授权修改 |
| `"Availability"` | `ScopeAvailability` | 可用性，服务被中断或降级 |
| `"Access Control"` | `ScopeAccessControl` | 访问控制，权限约束被绕过 |
| `"Accountability"` | `ScopeAccountability` | 可追溯性，行为无法归因到主体 |
| `"Authentication"` | `ScopeAuthentication` | 认证，身份验证机制失效 |
| `"Authorization"` | `ScopeAuthorization` | 授权，权限授予机制失效 |
| `"Non-Repudiation"` | `ScopeNonRepudiation` | 不可否认性，主体可否认其行为 |

```go
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

::: tip 三性 + 扩展
前三项（Confidentiality/Integrity/Availability）即经典 CIA 三性，是绝大多数弱点后果的核心维度；后五项是访问控制与责任追溯相关的扩展属性，常见于认证授权类弱点。
:::

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (s ConsequenceScope) String() string` |
| `IsValid` | `func (s ConsequenceScope) IsValid() bool` |
| `ParseConsequenceScope` | `func ParseConsequenceScope(s string) (ConsequenceScope, error)` |
| `AllConsequenceScopeValues` | `func AllConsequenceScopeValues() []ConsequenceScope` |

```go
s, err := cweskills.ParseConsequenceScope("Confidentiality")
fmt.Println(s, err)                                     // Confidentiality <nil>
fmt.Println(s.String())                                 // Confidentiality
fmt.Println(cweskills.ConsequenceScope("X").IsValid())  // false
fmt.Println(len(cweskills.AllConsequenceScopeValues())) // 8
```

::: warning 多词取值含空格/连字符
`"Access Control"`、`"Non-Repudiation"` 等取值含空格或连字符，`ParseConsequenceScope` 严格匹配原样。从外部 XML/JSON 读取时通常已是正确字面量，但从用户输入解析时需注意规范化。
:::

## 🔄 典型用法

```go
// 找出所有影响机密性的弱点
cwe89 := cweskills.ConsequenceScope("Confidentiality")
matches, _ := registry.FindByConsequenceScope(cwe89)
fmt.Println("影响机密性:", len(matches))

// 统计各影响范围的弱点数量
stats, _ := registry.CountByScope()
for scope, n := range stats {
	fmt.Println(scope, n)
}
```

## 💻 CLI 对应命令

```bash
cwe enum consequence-scope
```

输出全部合法取值，详见 [CLI enum consequence-scope](../cli/enum-consequence-scope)。

## 🔗 相关链接

- SDK 视角：[ConsequenceScope 后果范围枚举](../sdk/enum-consequence-scope)
- 概念背景：[后果范围与影响](../guide/concept-consequence)
- 按范围筛选：[`FindByConsequenceScope`](../sdk/find-by-consequence-scope)
- 按范围统计：[`CountByScope`](../sdk/count-by-scope)
- 后果影响（配套字段）：[ConsequenceImpact](./consequence-impact)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
