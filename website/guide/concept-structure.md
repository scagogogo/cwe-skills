---
title: 结构类型 (Simple/Chain/Composite)
outline: [2, 3]
---

# 🧱 结构类型 (Simple / Chain / Composite)

**结构类型**（Structure）描述一个 CWE 条目是「单一弱点」还是「由多个弱点组合而成」，以及组合的方式。CWE 定义了三种结构：**Simple（简单）**、**Chain（链式）**、**Composite（复合）**。

::: tip 结构 vs 抽象
[抽象层级](./concept-abstraction)回答「多通用」，**结构类型**回答「是单一弱点还是组合弱点」。两者正交——一个 Base 级弱点既可能是 Simple，也可能属于某个 Chain/Composite 的成员。
:::

---

## 📊 三种结构

`Structure` 是枚举类型，定义在 `enums.go`：

```go
type Structure string

const (
    StructureSimple    Structure = "Simple"
    StructureChain     Structure = "Chain"
    StructureComposite Structure = "Composite"
)
```

| 结构 | 含义 | 组合方式 | 典型示例 |
|------|------|---------|---------|
| **Simple** | 单一弱点，不依赖其他弱点的存在 | 无 | CWE-79 XSS、CWE-89 SQL 注入 |
| **Chain** | 链式弱点，多个弱点**按顺序**可达才产生漏洞 | 顺序依赖 | CWE-680 整数溢出 → 缓冲区溢出 |
| **Composite** | 复合弱点，多个弱点**同时存在**才产生漏洞 | 并存依赖 | CWE-352 CSRF（需多个弱点并存） |

```text
Simple     :  [A]              单点

Chain      :  [A] → [B] → [C]  必须依次发生
              （A 的成功为 B 创造条件，B 的成功为 C 创造条件）

Composite  :  [A] + [B] + [C]  必须同时存在
              （任一缺失，漏洞不成立）
```

---

## 🔗 链式 vs 复合：关键区别

| 维度 | Chain（链式） | Composite（复合） |
|------|--------------|------------------|
| 组合方式 | 顺序（前驱为后继创造条件） | 并存（同时存在） |
| 关系类型 | `CanPrecede` / `CanFollow` | `Requires` / `RequiredBy` |
| 成员关系 | 一个接一个 | 缺一不可 |
| 类比 | 推倒多米诺骨牌 | 多把锁同时打开 |

::: details 用关系类型区分
- **Chain** 的成员之间是顺序关系：`CanPrecede`（A 可前置 B）、`CanFollow`（B 可跟随 A）。
- **Composite** 的成员之间是依赖关系：`Requires`（A 需要 B 存在）、`RequiredBy`（A 被 B 所需要）。

详见 [关系类型](./concept-relationship)。
:::

---

## 🛠️ 枚举 API

```go
// 校验
cweskills.StructureSimple.IsValid()         // true
cweskills.Structure("Foo").IsValid()        // false

// 解析
s, err := cweskills.ParseStructure("Chain") // StructureChain, nil

// 穷举
cweskills.AllStructureValues()              // [Simple, Chain, Composite]
```

::: warning Composite 只用于复合元素
`StructureComposite` 主要出现在 `CompoundElement`（复合元素）上。普通弱点（`CWE` 结构体）的 `Structure` 字段通常是 `Simple`；链式/复合弱点的 `Structure` 是 `Chain`/`Composite`，且它们作为 `CompoundElement` 注册。见 [复合元素](./concept-compound)。
:::

---

## 🧭 导航链式与复合成员

`Navigator` 提供专门的方法访问链式/复合弱点的成员（**仅离线路径有数据**，因在线 API 不返回这些关系）：

```go
nav := cweskills.NewNavigator(registry)

// 链式弱点的成员（按 CanPrecede/CanFollow 关系）
chainMembers := nav.ChainMembers(680)

// 复合弱点的成员（按 Requires/RequiredBy 关系）
compositeMembers := nav.CompositeMembers(352)

// 单独查顺序/依赖关系
precede   := nav.CanPrecede(680)   // 此弱点可以前置哪些
follow    := nav.CanFollow(680)    // 此弱点可以跟随哪些
requires  := nav.Requires(352)     // 此弱点依赖哪些
requiredBy := nav.RequiredBy(352)  // 哪些弱点依赖此弱点
```

::: info 在线 API 拿不到这些
`CanPrecede`/`CanFollow`/`Requires`/`RequiredBy` 等关系**只在离线 XML 里齐全**。若用 `APIClient.GetParents` 之类在线方法，链式/复合成员通常查不到。需要时走 [离线 XML](../sdk/xml-parser)。
:::

---

## 🔍 按结构查找与过滤

```go
// 查找所有链式弱点
chains := cweskills.FindChains(registry)

// 查找所有复合弱点
composites := cweskills.FindComposites(registry)

// 按结构过滤
filtered := cweskills.Filter(results, cweskills.FilterOption{
    Structure: cweskills.StructureChain,
})

// 也可用专门的查找函数
simpleOnes := cweskills.FindByStructure(registry, cweskills.StructureSimple)
```

---

## 🎯 用途：为什么结构类型重要

### 1. 攻击路径分析

链式弱点揭示了**攻击者如何串联多个弱点**形成完整利用链。识别一条 Chain，等于识别了一条潜在攻击路径。

```go
// 找出 CWE-680 这条链的所有环节
for _, member := range nav.ChainMembers(680) {
    fmt.Println(member.CWEID(), member.Name)
}
```

### 2. 防御覆盖评估

复合弱点说明「**只修其中一个不够**」——必须同时缓解所有成员。安全评估时，识别 Composite 能避免「修了一半」的假安全感。

### 3. 工具结果去重与归并

同一条链上的多个弱点，可能是同一个根因的不同表现。按结构归并可减少重复告警。

---

## ⚠️ 注意事项

::: warning Chain/Composite 是 CompoundElement
带 `Chain`/`Composite` 结构的条目通常是 `CompoundElement`，而非普通 `CWE` 弱点。用 `registry.GetCompoundElement(id)` 获取，而非 `registry.Get(id)`（后者只查 weaknesses）。
:::

::: info Simple 是绝大多数
CWE 目录里绝大多数条目是 `Simple`。链式和复合条目数量不多，但分析价值高。
:::

---

## 📖 相关文档

- [复合元素 (CompoundElement)](./concept-compound)
- [关系类型](./concept-relationship)
- [关系导航 API](../sdk/navigator)
- [搜索过滤 API](../sdk/search)
- [Structure 枚举参考](../enums/structure)
