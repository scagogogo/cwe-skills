---
title: 技能 04 — 枚举类型
outline: [2, 3]
---

# 📚 技能 04 — 枚举类型

列出并理解 CWE 规范定义的类型化枚举：抽象层级、状态值、关系类型、后果范围与影响等。这些枚举贯穿所有 SDK API。

<Badge type="tip" text="类型参考"/>
<Badge type="info" text="内置数据"/>

---

## 🎯 技能目标

- 列出每种枚举类型的全部合法值
- 把字符串解析为类型化枚举值
- 校验、排序、分类枚举值

---

## 💻 CLI 命令

### enum — 列出枚举值

```bash
cwe enum abstraction      # 抽象层级
cwe enum structure         # 结构类型
cwe enum status            # 状态值
cwe enum likelihood        # 利用可能性
cwe enum relationship      # 关系类型
cwe enum scope             # 后果范围
cwe enum impact            # 后果影响
cwe enum viewtype          # 视图类型
```

```text
抽象层级 (Class/Base/Variant/Pillar) (4 项):
  - Pillar
  - Class
  - Base
  - Variant
```

JSON：`["Pillar", "Class", "Base", "Variant"]`

---

## 🔧 SDK API

### Parse 函数

每种枚举有 `Parse*` 函数把字符串转为类型化值：

```go
abstr, _ := cweskills.ParseAbstraction("Base")            // AbstractionBase
status, _ := cweskills.ParseStatus("Stable")              // StatusStable
nature, _ := cweskills.ParseRelationshipNature("ChildOf") // RelationshipChildOf
structure, _ := cweskills.ParseStructure("Chain")         // StructureChain
likelihood, _ := cweskills.ParseLikelihoodOfExploit("High") // LikelihoodHigh
scope, _ := cweskills.ParseConsequenceScope("Confidentiality") // ScopeConfidentiality
impact, _ := cweskills.ParseConsequenceImpact("High")     // ImpactHigh
viewType, _ := cweskills.ParseViewType("Graph")           // ViewTypeGraph
```

### 校验与字符串化

```go
abstr.IsValid()   // bool
abstr.String()    // 规范字符串
```

### 全部值

```go
cweskills.AllAbstractionValues()         // []Abstraction
cweskills.AllStatusValues()              // []Status
cweskills.AllRelationshipNatureValues()  // []RelationshipNature
cweskills.AllConsequenceScopeValues()    // []ConsequenceScope
cweskills.AllConsequenceImpactValues()   // []ConsequenceImpact
cweskills.AllViewTypeValues()            // []ViewType
```

### 排序与分类

```go
cweskills.AbstractionOrder(cweskills.AbstractionPillar)  // 0
cweskills.AbstractionOrder(cweskills.AbstractionBase)    // 2

cweskills.LikelihoodOrder(cweskills.LikelihoodHigh)      // 2
cweskills.ImpactOrder(cweskills.ImpactHigh)              // 2

nature.IsHierarchical()  // ChildOf/ParentOf → true
nature.IsSequential()    // CanPrecede/CanFollow → true
nature.IsDependency()    // Requires/RequiredBy → true
nature.IsPeer()          // PeerOf/CanAlsoBe → true
```

---

## 📝 示例

### 命令行

```bash
# 查看抽象层级有哪些值
cwe enum abstraction -o json

# 查关系类型分类
cwe enum relationship
```

### Go

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    for _, a := range cweskills.AllAbstractionValues() {
        fmt.Printf("%s (order=%d)\n", a, cweskills.AbstractionOrder(a))
    }
}
```

---

## 📊 枚举速查

### Abstraction 抽象层级

| 值 | 说明 |
|----|------|
| `Pillar` | 顶层抽象（如 CWE-664） |
| `Class` | 中层类别（如 CWE-74 注入） |
| `Base` | 具体弱点类型（如 CWE-79 XSS） |
| `Variant` | 平台/语言特定变体 |

### Status 状态

`Stable` · `Usable` · `Draft` · `Incomplete` · `Obsolete` · `Deprecated`

### RelationshipNature 关系类型

| 类别 | 值 |
|------|----|
| 层级 | `ChildOf` · `ParentOf` |
| 顺序 | `CanPrecede` · `CanFollow` |
| 依赖 | `Requires` · `RequiredBy` |
| 对等 | `PeerOf` · `CanAlsoBe` |
| 成员 | `MemberOf` · `HasMember` |

::: details 结构类型 Structure
`Simple`（单一弱点） · `Chain`（顺序链） · `Composite`（同时组合）
:::

---

## 🤖 AI 代理使用提示

- 用户问「CWE 有哪些抽象层级」时，AI 直接 `cwe enum abstraction`。
- 解析用户输入的枚举字符串前，可用 `cwe enum <type>` 确认合法值。
- 解释关系类型时，用 `IsHierarchical/IsSequential` 等分类方法说明语义。

---

## 📖 相关文档

- [枚举参考总览](../enums/overview)
- [Abstraction](../enums/abstraction) · [Status](../enums/status) · [RelationshipNature](../enums/relationship-nature)
- [CLI: enum](../cli/enum)
- [SDK: enums](../sdk/enums)
- [返回 Skills 总览](./)
