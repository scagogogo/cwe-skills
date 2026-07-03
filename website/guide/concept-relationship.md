---
title: 关系类型 (ChildOf/Requires…)
outline: [2, 3]
---

# 🔗 关系类型 (ChildOf / Requires …)

CWE 不只是一份编号列表，而是一张**多语义有向图**。条目之间通过 10 种**关系类型**相连，分为层级、顺序、依赖、对等四类。理解关系类型，是读懂 CWE 体系、用好 [关系导航](../sdk/navigator) 的前提。

---

## 📚 全部 10 种关系

`RelationshipNature` 是枚举类型，定义在 `enums.go`：

```go
type RelationshipNature string

const (
    RelationshipChildOf    RelationshipNature = "ChildOf"
    RelationshipParentOf   RelationshipNature = "ParentOf"
    RelationshipCanPrecede RelationshipNature = "CanPrecede"
    RelationshipCanFollow  RelationshipNature = "CanFollow"
    RelationshipRequires   RelationshipNature = "Requires"
    RelationshipRequiredBy RelationshipNature = "RequiredBy"
    RelationshipCanAlsoBe  RelationshipNature = "CanAlsoBe"
    RelationshipPeerOf     RelationshipNature = "PeerOf"
    RelationshipMemberOf   RelationshipNature = "MemberOf"
    RelationshipHasMember  RelationshipNature = "Has_Member"
)
```

::: warning Has_Member 带下划线
注意 `RelationshipHasMember` 的字符串值是 `"Has_Member"`（带下划线），是 MITRE 原始定义如此，其余 9 种都是驼峰无分隔。`ParseRelationshipNature` 只认这 10 个精确字符串。
:::

---

## 🧩 四类关系

CWE Skills 提供 `IsHierarchical()` / `IsSequential()` / `IsDependency()` / `IsPeer()` 四个方法判定关系类别：

| 类别 | 关系 | 判定方法 | 含义 |
|------|------|---------|------|
| **层级** (Hierarchical) | `ChildOf`、`ParentOf`、`MemberOf`、`Has_Member` | `IsHierarchical()` | 父子 / 成员归属，可建树 |
| **顺序** (Sequential) | `CanPrecede`、`CanFollow` | `IsSequential()` | 链式前后，攻击路径 |
| **依赖** (Dependency) | `Requires`、`RequiredBy` | `IsDependency()` | 复合弱点的并存依赖 |
| **对等** (Peer) | `PeerOf`、`CanAlsoBe` | `IsPeer()` | 相似性，非层级 |

```go
cweskills.RelationshipChildOf.IsHierarchical()   // true
cweskills.RelationshipCanPrecede.IsSequential()  // true
cweskills.RelationshipRequires.IsDependency()    // true
cweskills.RelationshipPeerOf.IsPeer()            // true
```

---

## 1️⃣ 层级关系（Hierarchical）

层级关系构成 CWE 的「族谱」，是建树与上下钻的基础。

| 关系 | 方向 | 含义 |
|------|------|------|
| `ChildOf` | A ChildOf B | A 是 B 的子项（更具体） |
| `ParentOf` | A ParentOf B | A 是 B 的父项（更通用） |
| `MemberOf` | A MemberOf B | A 是类别/视图 B 的成员 |
| `Has_Member` | A Has_Member B | 类别/视图 A 包含成员 B |

```text
        CWE-74 (注入, Class)
         ▲ ChildOf
         │
        CWE-79 (XSS, Base)
```

::: tip ChildOf 是主力
CWE 层级树主要靠 `ChildOf` 构建。`ParentOf` 是 `ChildOf` 的反向；`MemberOf`/`Has_Member` 用于条目归入类别/视图（非父子层级）。见 [类别](./concept-category)、[视图](./concept-view)。
:::

导航方法：

```go
nav := cweskills.NewNavigator(registry)
nav.Parents(79)      // 直接父级
nav.Children(79)     // 直接子级
nav.Ancestors(79)    // 所有祖先（递归）
nav.Descendants(79)  // 所有后代（递归）
nav.Siblings(79)     // 同级（同父级的其他子级）
```

---

## 2️⃣ 顺序关系（Sequential）

顺序关系描述**链式弱点**中环节的前后：A 成功利用后，为 B 创造了可达条件。

| 关系 | 方向 | 含义 |
|------|------|------|
| `CanPrecede` | A CanPrecede B | A 可前置 B（A 发生后 B 才可能） |
| `CanFollow` | A CanFollow B | A 可跟随 B（B 发生后 A 才可能） |

```text
[整数溢出] --CanPrecede--> [缓冲区溢出] --CanPrecede--> [任意代码执行]
```

```go
nav.CanPrecede(680)  // CWE-680 可以前置哪些弱点
nav.CanFollow(680)   // CWE-680 可以跟随哪些弱点
nav.ChainMembers(680) // 整条链的成员
```

::: info CanPrecede 与 CanFollow 互为反向
`A CanPrecede B` 等价于 `B CanFollow A`。导航时用哪个取决于你想从哪个方向看。
:::

---

## 3️⃣ 依赖关系（Dependency）

依赖关系描述**复合弱点**中成员的并存要求：复合弱点成立，要求所有成员**同时存在**。

| 关系 | 方向 | 含义 |
|------|------|------|
| `Requires` | A Requires B | A（复合）需要 B 存在 |
| `RequiredBy` | A RequiredBy B | A 被 B（复合）所需要 |

```text
         [复合弱点 CSRF]
          │ Requires   │ Requires   │ Requires
          ▼            ▼            ▼
       [弱点A]      [弱点B]      [弱点C]   ← 缺一不可
```

```go
nav.Requires(352)       // CWE-352 依赖哪些弱点
nav.RequiredBy(79)      // 哪些复合弱点依赖 CWE-79
nav.CompositeMembers(352) // 复合弱点的全部成员
```

---

## 4️⃣ 对等关系（Peer）

对等关系表示两个弱点有**相似性**，但不适合用层级/顺序/依赖描述。

| 关系 | 方向 | 含义 |
|------|------|------|
| `PeerOf` | A PeerOf B | A 与 B 相似 |
| `CanAlsoBe` | A CanAlsoBe B | 在适当上下文中 A 也可视为 B |

```go
nav.Peers(79)        // CWE-79 的对等弱点
nav.CanAlsoBe(79)    // CWE-79 也可以是哪些
```

::: tip 对等是双向语义
`PeerOf` 通常是对称的（A PeerOf B 蕴含 B PeerOf A），但 `CanAlsoBe` 更偏向「在某些上下文下可等同」。导航结果按关系方向返回。
:::

---

## 🧭 关系图总览

```text
            ┌──────────── 层级 (Hierarchical) ────────────┐
            │  ChildOf  ParentOf  MemberOf  Has_Member    │  → 建树、上下钻
            └─────────────────────────────────────────────┘
            ┌──────────── 顺序 (Sequential) ──────────────┐
            │  CanPrecede  CanFollow                       │  → 链式攻击路径
            └─────────────────────────────────────────────┘
            ┌──────────── 依赖 (Dependency) ──────────────┐
            │  Requires  RequiredBy                       │  → 复合弱点并存
            └─────────────────────────────────────────────┘
            ┌──────────── 对等 (Peer) ────────────────────┐
            │  PeerOf  CanAlsoBe                          │  → 相似性
            └─────────────────────────────────────────────┘
```

---

## 🛠️ 枚举 API

```go
// 校验
cweskills.RelationshipChildOf.IsValid()                 // true
cweskills.RelationshipNature("Foo").IsValid()           // false

// 解析
r, err := cweskills.ParseRelationshipNature("ChildOf")  // RelationshipChildOf
r, err := cweskills.ParseRelationshipNature("Has_Member") // RelationshipHasMember

// 穷举
cweskills.AllRelationshipNatureValues() // 10 个全部

// 类别判定
r.IsHierarchical(); r.IsSequential(); r.IsDependency(); r.IsPeer()
```

---

## 🚦 在线 vs 离线的关系覆盖

::: danger 在线 API 关系不全
MITRE REST API **只返回部分关系类型**（主要是 `ChildOf` 等层级关系）。`CanPrecede`/`CanFollow`/`Requires`/`RequiredBy`/`PeerOf`/`CanAlsoBe` 等**只有在离线 XML 里才齐全**。

| 关系 | 在线 API | 离线 XML |
|------|---------|---------|
| ChildOf / ParentOf | ✅ 部分 | ✅ 完整 |
| MemberOf / Has_Member | ✅ 部分 | ✅ 完整 |
| CanPrecede / CanFollow | ❌ | ✅ |
| Requires / RequiredBy | ❌ | ✅ |
| PeerOf / CanAlsoBe | ❌ | ✅ |

需要链式/依赖/对等导航时，必须走 [离线 XML](../sdk/xml-parser)。详见 [在线 vs 离线](./online-offline)。
:::

---

## 🎯 高级导航：最短路径与深度

层级关系之上，`Navigator` 还提供图算法级别的查询（仅离线）：

```go
// 最短路径：从 79 到 1 经过的 ID 序列
path := nav.ShortestPath(79, 1) // 例如 [79, 74, 707, ..., 1]

// 关系深度：1 是 79 的第几代祖先
depth := nav.RelationshipDepth(1, 79) // 例如 5

// 祖先/后代判定
nav.IsAncestorOf(1, 79)   // true：1 是 79 的祖先
nav.IsDescendantOf(79, 1) // true：79 是 1 的后代
nav.IsRelated(79, 1)      // true：两者有关系
```

::: details 最短路径用什么关系？
`ShortestPath` / `IsRelated` 综合考虑注册表里所有已建立的关系边（父子、对等、成员等），在混合关系图上做 BFS。适合回答「这两个 CWE 之间到底隔了几层」。
:::

---

## 📖 相关文档

- [关系导航 API](../sdk/navigator)
- [抽象层级](./concept-abstraction)
- [结构类型](./concept-structure)
- [复合元素](./concept-compound)
- [在线 vs 离线](./online-offline)
- [RelationshipNature 枚举参考](../enums/relationship-nature)
