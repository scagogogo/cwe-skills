---
title: 复合元素 (CompoundElement)
outline: [2, 3]
---

# 🧩 复合元素 (CompoundElement)

有些漏洞不是单一弱点造成的，而是**多个弱点组合**的结果。CWE 用**复合元素（CompoundElement）**来描述这类「组合型弱点」，分为**链式（Chain）**与**复合（Composite）**两种结构。

::: tip 与 [结构类型](./concept-structure) 的关系
`Structure` 枚举描述「Simple / Chain / Composite」三种结构。带 `Chain`/`Composite` 结构的条目，就是这里说的**复合元素**，在 CWE Skills 里用 `CompoundElement` 结构体建模。
:::

---

## 📦 CompoundElement 结构

```go
type CompoundElement struct {
    ID            int            // 数字标识符
    Name          string         // 名称
    Structure     Structure      // Chain 或 Composite（必填）
    Status        Status         // 状态
    Description   string         // 描述
    Relationships []Relationship // 关系（含 Requires/RequiredBy 或 CanPrecede/CanFollow）
}
```

```go
ce, ok := registry.GetCompoundElement(680)
if ok {
    fmt.Println(ce.Name, ce.Structure) // 整数溢出到缓冲区溢出 Chain
}
```

::: warning Structure 必填
`CompoundElement.Structure` 字段没有 `omitempty`，是必填的，且只能是 `Chain` 或 `Composite`。创建时用 `NewCompoundElement(id, name, structure)`。
:::

---

## 🔗 两种复合结构

### Chain（链式）— 顺序依赖

多个弱点**按顺序**可达才产生漏洞：A 成功利用后为 B 创造条件，B 成功后为 C 创造条件。

```text
[整数溢出] ──CanPrecede──► [缓冲区溢出] ──CanPrecede──► [任意代码执行]
```

- 关系类型：`CanPrecede`（可前置）、`CanFollow`（可跟随）
- 典型例子：CWE-680（整数溢出到缓冲区溢出）

### Composite（复合）— 并存依赖

多个弱点**同时存在**才产生漏洞：缺任何一个，漏洞不成立。

```text
         [复合弱点]
          Requires   Requires   Requires
            │           │          │
         [弱点A]     [弱点B]    [弱点C]   ← 缺一不可
```

- 关系类型：`Requires`（需要）、`RequiredBy`（被需要）
- 典型例子：CWE-352（CSRF，需多个弱点并存）

| 维度 | Chain | Composite |
|------|-------|-----------|
| 组合方式 | 顺序 | 并存 |
| 关系 | CanPrecede / CanFollow | Requires / RequiredBy |
| 成员关系 | 一个接一个 | 缺一不可 |
| 类比 | 多米诺骨牌 | 多把锁同时打开 |

---

## 🧭 导航复合元素成员

`Navigator` 提供专门方法访问复合元素的成员（**仅离线路径有数据**）：

```go
nav := cweskills.NewNavigator(registry)

// 链式弱点的成员
chainMembers := nav.ChainMembers(680)

// 复合弱点的成员
compositeMembers := nav.CompositeMembers(352)

// 单独查顺序/依赖关系
nav.CanPrecede(680)    // CWE-680 可前置哪些
nav.CanFollow(680)     // CWE-680 可跟随哪些
nav.Requires(352)      // CWE-352 依赖哪些
nav.RequiredBy(79)     // 哪些复合弱点依赖 CWE-79
```

::: danger 在线 API 拿不到这些
`CanPrecede`/`CanFollow`/`Requires`/`RequiredBy` 关系**只在离线 XML 里齐全**，MITRE REST API 不返回。需要导航复合元素时必须走 [离线 XML](../sdk/xml-parser)。
:::

---

## 🔍 查找复合元素

```go
// 查找所有链式弱点
chains := cweskills.FindChains(registry)

// 查找所有复合弱点
composites := cweskills.FindComposites(registry)

// 按 ID 取单个复合元素
ce, ok := registry.GetCompoundElement(680)

// 统计数量
fmt.Println("链式:", registry.CompoundElementCount())
```

---

## 💻 CLI 操作

```bash
# 查询链式弱点的成员
cwe nav chain-members CWE-680 --xml cwec_v4.15.xml

# 查询复合弱点的成员
cwe nav composite-members CWE-352 --xml cwec_v4.15.xml

# 查询顺序/依赖关系
cwe nav precede CWE-680 --xml cwec_v4.15.xml
cwe nav follow CWE-680 --xml cwec_v4.15.xml
cwe nav requires CWE-352 --xml cwec_v4.15.xml
cwe nav required-by CWE-79 --xml cwec_v4.15.xml
```

---

## 🎯 用途

### 1. 攻击路径还原

链式弱点揭示了**攻击者如何串联弱点**形成完整利用链。识别一条 Chain，等于识别了一条潜在攻击路径，有助于在早期环节阻断。

```go
// 还原 CWE-680 这条链的完整环节
fmt.Println("攻击链:")
for _, m := range nav.ChainMembers(680) {
    fmt.Println(" →", m.CWEID(), m.Name)
}
```

### 2. 防御完整性评估

复合弱点说明「**只修其中一个不够**」。识别 Composite 能避免「修了一半」的假安全感——必须同时缓解所有成员。

```go
// CWE-352 的所有必备弱点
for _, m := range nav.CompositeMembers(352) {
    fmt.Println("必须同时修复:", m.CWEID())
}
```

### 3. 告警归并与根因分析

同一条链上的多个告警可能源自同一根因。按复合元素归并，可减少重复告警、定位真正的根因环节。

---

## ⚠️ 注意事项

::: warning 复合元素用 GetCompoundElement 查
带 `Chain`/`Composite` 结构的条目是 `CompoundElement`，**不要用 `registry.Get(id)`**（那只查 weaknesses 集合）。查复合元素用 `registry.GetCompoundElement(id)`。
:::

::: info 数量不多但价值高
CWE 目录里复合元素数量远少于普通弱点，但它们往往对应真实世界的高危攻击链/复合漏洞，分析价值很高。
:::

::: tip Structure 字段决定关系语义
同一个 `CompoundElement`，`Structure=Chain` 时其 `Relationships` 是顺序关系，`Structure=Composite` 时是依赖关系。导航时 `ChainMembers` / `CompositeMembers` 会按 Structure 选择正确的关系类型。
:::

---

## 📖 相关文档

- [结构类型 (Simple/Chain/Composite)](./concept-structure)
- [关系类型](./concept-relationship)
- [关系导航 API](../sdk/navigator)
- [CWE 是什么](./concept-cwe)
