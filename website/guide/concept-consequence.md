---
title: 后果范围与影响
outline: [2, 3]
---

# ⚠️ 后果范围与影响

一个弱点被利用后，会造成什么安全后果？CWE 用 **后果（Consequence）** 来描述，每条后果包含**影响范围（Scope）**和**影响严重程度（Impact）**两个维度。理解这两个维度，是做风险评估与优先级排序的关键。

---

## 📦 Consequence 结构

`Consequence` 定义在 `consequences.go`，是 `CWE.CommonConsequences` 字段的元素类型。一个弱点可以有多种后果，每种后果可影响多个范围：

```go
type Consequence struct {
    Scopes     []ConsequenceScope   // 影响范围（可多个）
    Impacts    []ConsequenceImpact  // 影响严重程度（可多个）
    Likelihood LikelihoodOfExploit  // 利用可能性
    Note       string               // 说明
}
```

```go
// 取 CWE-79 的常见后果
for _, c := range weakness.CommonConsequences {
    fmt.Println("范围:", c.Scopes, "影响:", c.Impacts, "可能性:", c.Likelihood)
}
```

---

## 🛡️ 影响范围：8 种 Scope

`ConsequenceScope` 描述弱点影响的安全属性。前三个即经典的 **CIA 三要素**：

```go
const (
    ScopeConfidentiality ConsequenceScope = "Confidentiality" // 机密性
    ScopeIntegrity       ConsequenceScope = "Integrity"       // 完整性
    ScopeAvailability    ConsequenceScope = "Availability"    // 可用性
    ScopeAccessControl   ConsequenceScope = "Access Control"  // 访问控制
    ScopeAccountability  ConsequenceScope = "Accountability"  // 可追溯性
    ScopeAuthentication  ConsequenceScope = "Authentication"  // 身份认证
    ScopeAuthorization   ConsequenceScope = "Authorization"   // 授权
    ScopeNonRepudiation  ConsequenceScope = "Non-Repudiation" // 不可否认性
)
```

| 范围 | 含义 | 典型弱点 |
|------|------|---------|
| **Confidentiality** 机密性 | 信息被未授权读取 | CWE-200 敏感信息暴露 |
| **Integrity** 完整性 | 信息被未授权修改 | CWE-89 SQL 注入改数据 |
| **Availability** 可用性 | 服务被中断或降级 | CWE-400 资源耗尽 |
| **Access Control** 访问控制 | 绕过访问限制 | CWE-862 缺失授权 |
| **Accountability** 可追溯性 | 无法追溯操作来源 | 日志被篡改 |
| **Authentication** 身份认证 | 认证被绕过/伪造 | CWE-287 认证不当 |
| **Authorization** 授权 | 越权操作 | CWE-863 授权不当 |
| **Non-Repudiation** 不可否认性 | 行为可被否认 | 审计链断裂 |

### CIA 三要素

**C**onfidentiality（机密性）、**I**ntegrity（完整性）、**A**vailability（可用性）是信息安全的基石，称为 CIA 三要素。多数弱点的后果至少影响其中之一：

- SQL 注入 → 完整性（改数据）+ 机密性（读数据）
- XSS → 完整性（注入内容）+ 机密性（窃取 cookie）
- DoS → 可用性

::: tip 一个后果可影响多个范围
`Consequence.Scopes` 是切片。例如某注入类弱点的后果可能同时影响 Confidentiality、Integrity、Availability 三项。
:::

---

## 📊 影响严重程度：4 级 Impact

`ConsequenceImpact` 描述后果的严重程度：

```go
const (
    ImpactHigh    ConsequenceImpact = "High"
    ImpactMedium  ConsequenceImpact = "Medium"
    ImpactLow     ConsequenceImpact = "Low"
    ImpactUnknown ConsequenceImpact = "Unknown"
)
```

| 等级 | 权重 (`ImpactOrder()`) | 含义 |
|------|----------------------|------|
| High | 4 | 严重 |
| Medium | 3 | 中等 |
| Low | 2 | 轻微 |
| Unknown | 1 | 未知 |

```go
cweskills.ImpactHigh.ImpactOrder()   // 4
cweskills.ImpactUnknown.ImpactOrder() // 1
```

::: info Impact 与范围是两个维度
「影响 Confidentiality」说的是**影响哪个安全属性**，「High」说的是**影响多严重**。一个后果 = 范围（可多个）× 影响等级（可多个）。例如 `{Scopes:[Confidentiality,Integrity], Impacts:[High]}`。
:::

---

## 🎲 利用可能性：LikelihoodOfExploit

`LikelihoodOfExploit` 描述弱点**被实际利用的可能性**，与后果严重程度共同决定风险：

```go
const (
    LikelihoodHigh    LikelihoodOfExploit = "High"
    LikelihoodMedium  LikelihoodOfExploit = "Medium"
    LikelihoodLow     LikelihoodOfExploit = "Low"
    LikelihoodUnknown LikelihoodOfExploit = "Unknown"
)
```

| 等级 | 权重 (`LikelihoodOrder()`) |
|------|---------------------------|
| High | 4 |
| Medium | 3 |
| Low | 2 |
| Unknown | 1 |

::: tip 风险 = 可能性 × 影响
经典风险公式：**风险 = LikelihoodOfExploit × Impact**。CWE Skills 把两者都做成带排序权重的枚举，方便你做风险评分与排序。`Likelihood` 既出现在 `Consequence.Likelihood`，也出现在 `CWE.LikelihoodOfExploit` 字段。
:::

---

## 🛠️ 枚举 API

三个枚举遵循统一 API：

```go
// 校验
cweskills.ScopeConfidentiality.IsValid()    // true
cweskills.ImpactHigh.IsValid()              // true
cweskills.LikelihoodHigh.IsValid()          // true

// 解析
s, _ := cweskills.ParseConsequenceScope("Integrity")
i, _ := cweskills.ParseConsequenceImpact("High")
l, _ := cweskills.ParseLikelihoodOfExploit("Medium")

// 穷举
cweskills.AllConsequenceScopeValues()   // 8 个
cweskills.AllConsequenceImpactValues()  // 4 个
cweskills.AllLikelihoodOfExploitValues()// 4 个
```

---

## 🔍 按后果查找与过滤

```go
// 查找所有影响机密性的弱点
confidentialityOnes := cweskills.FindByConsequenceScope(registry, cweskills.ScopeConfidentiality)

// 按利用可能性查找
highLikelihood := cweskills.FindByLikelihood(registry, cweskills.LikelihoodHigh)

// 过滤：影响机密性 + 高可能性
filtered := cweskills.Filter(results, cweskills.FilterOption{
    ConsequenceScope:     cweskills.ScopeConfidentiality,
    LikelihoodOfExploit:  cweskills.LikelihoodHigh,
})

// 按可能性分组
groups := cweskills.GroupByLikelihood(registry.GetAll())
```

::: details FilterOption 里的后果字段
`FilterOption` 支持 `ConsequenceScope` 字段，按「该弱点的任一后果涉及此范围」过滤。结合 `LikelihoodOfExploit` 可精准筛出「高可能 + 影响机密性」的高风险集合。
:::

---

## 🎯 用途：风险评估

### 1. 风险评分

```go
func riskScore(c *cweskills.CWE) int {
    score := 0
    for _, con := range c.CommonConsequences {
        for _, imp := range con.Impacts {
            score += imp.ImpactOrder() * c.LikelihoodOfExploit.LikelihoodOrder()
        }
    }
    return score
}
```

### 2. 优先级排序

按「影响 CIA 三要素 + High 影响 + High 可能性」排序，先修最危险的：

```go
critical := cweskills.Filter(registry.GetAll(), cweskills.FilterOption{
    ConsequenceScope:    cweskills.ScopeConfidentiality,
    LikelihoodOfExploit: cweskills.LikelihoodHigh,
})
```

### 3. 合规映射

许多合规框架（PCI DSS、SOC 2）要求保护 CIA。按范围过滤可快速定位影响合规要求的弱点集。

---

## ⚠️ 注意事项

::: warning Impact 是后果维度，不是抽象层级
不要把 `ConsequenceImpact`（后果严重程度）和 [抽象层级](./concept-abstraction) 混淆。前者描述「后果多严重」，后者描述「弱点多通用」。
:::

::: info Not all weaknesses have Consequences
并非每条 CWE 都填了 `CommonConsequences`。Pillar/Class 等高层级条目往往没有具体后果描述；Variant/Base 通常有。过滤时拿到空切片是正常的。
:::

---

## 📖 相关文档

- [CWE 是什么](./concept-cwe)
- [抽象层级](./concept-abstraction)
- [搜索过滤 API](../sdk/search)
- [Consequence 结构参考](../sdk/consequence)
