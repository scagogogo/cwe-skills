---
title: 枚举参考总览
outline: [2, 3]
---

# 📚 枚举参考总览

CWE Skills 的 `enums.go` 为 CWE 数据模型中的全部枚举字段定义了强类型的 Go 枚举。本节是这些枚举的**参考手册**，逐个列出每个枚举的全部取值、含义、通用方法与额外方法。

所有枚举类型都遵循同一套设计约定，掌握一种即可触类旁通。

## 🧬 类型定义约定

每个枚举类型都是基于 `string` 的具名类型：

```go
type Abstraction string  // 以 Abstraction 为例
```

取值通过导出常量定义，常量名以类型名为前缀：

```go
const (
	AbstractionPillar  Abstraction = "Pillar"
	AbstractionClass   Abstraction = "Class"
	// ...
)
```

::: tip 为何不用 iota
CWE 的枚举值是字符串（与 MITRE XML / API 中的字面量一一对应），用字符串常量可直接与外部数据互相转换，无需维护额外的映射表。`iota` 适合纯数值枚举，此处不适用。
:::

## ✅ 四件套通用方法

<Badge type="tip" text="统一约定" /> 每个枚举类型 `Xxx` 都提供以下四件套：

| 方法 / 函数 | 签名 | 作用 |
| --- | --- | --- |
| `String` | `func (x Xxx) String() string` | 返回字面量字符串 |
| `IsValid` | `func (x Xxx) IsValid() bool` | 校验是否为合法取值 |
| `ParseXxx` | `func ParseXxx(s string) (Xxx, error)` | 从字符串解析，非法值返回 `ValidationError` |
| `AllXxxValues` | `func AllXxxValues() []Xxx` | 返回全部合法取值 |

```go
// 以 Structure 为例演示四件套
v, err := cweskills.ParseStructure("Chain")   // 解析
fmt.Println(v, err)                            // Chain <nil>
fmt.Println(v.String())                        // Chain
fmt.Println(v.IsValid())                       // true
fmt.Println(cweskills.AllStructureValues())    // [Simple Chain Composite]
```

::: warning ParseXxx 严格匹配
`ParseXxx` 区分大小写、不修剪空白，且不接受同义词。例如 `ParseStatus("stable")` 会返回错误，必须传入 `"Stable"`。外部输入需先做规范化。
:::

## 🔄 额外方法

部分枚举在四件套之外提供额外方法：

| 枚举 | 额外方法 | 用途 |
| --- | --- | --- |
| `Abstraction` | `AbstractionOrder() int` | 抽象层级排序权重（Pillar=4 → Variant=1） |
| `LikelihoodOfExploit` | `LikelihoodOrder() int` | 利用可能性排序权重（High=4 → Unknown=1） |
| `ConsequenceImpact` | `ImpactOrder() int` | 影响严重程度排序权重（High=4 → Unknown=1） |
| `RelationshipNature` | `IsHierarchical/IsSequential/IsDependency/IsPeer` | 关系语义分类 |

## 🏷️ 枚举清单

| 枚举 | 类型 | 含义 | 详情 |
| --- | --- | --- | --- |
| `Abstraction` | 抽象层级 | Pillar/Class/Base/Variant | [详情](./abstraction) |
| `Structure` | 结构类型 | Simple/Chain/Composite | [详情](./structure) |
| `Status` | 状态 | Stable/Usable/Draft/... | [详情](./status) |
| `LikelihoodOfExploit` | 利用可能性 | High/Medium/Low/Unknown | [详情](./likelihood) |
| `RelationshipNature` | 关系类型 | ChildOf/Requires/PeerOf/... | [详情](./relationship-nature) |
| `ConsequenceScope` | 后果范围 | Confidentiality/Integrity/... | [详情](./consequence-scope) |
| `ConsequenceImpact` | 后果影响 | High/Medium/Low/Unknown | [详情](./consequence-impact) |
| `ViewType` | 视图类型 | Graph/Explicit Slice/Implicit Slice | [详情](./view-type) |
| `PlatformType` | 平台类型 | Language/Operating System/... | [详情](./platform-type) |

::: details 关于本节未收录的枚举
`enums.go` 另有 `Prevalence`、`IntroductionPhase`、`MitigationPhase`、`Effectiveness` 等枚举，本参考手册聚焦与 CWE 条目核心字段直接对应的九类，其余可在 [SDK enums 概览](../sdk/enums) 查阅。
:::

## 🎯 何时用枚举而非裸字符串

```go
// ❌ 容易拼错、无编译期检查
if cwe.Status == "Stabel" { ... }

// ✅ 编译期检查、可校验、可遍历
if cwe.Status == cweskills.StatusStable { ... }
if !cwe.Status.IsValid() { return errors.New("非法状态") }
```

## 🔗 相关链接

- SDK 视角的枚举总览：[enums.go 概览](../sdk/enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
- CLI 枚举命令总览：[cwe enum](../cli/enum)
