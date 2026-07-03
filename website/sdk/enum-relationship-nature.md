---
title: RelationshipNature 关系类型枚举
outline: [2, 3]
---

# 📚 RelationshipNature — 关系类型枚举

`RelationshipNature` 表示 CWE 条目之间的关系类型，共 10 种。它是关系导航与图谱构建的基础，[`CWE.GetParentIDs`](./cwe-relationship-methods) 等方法都基于它筛选关系。

## 📋 类型与常量

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

## 📝 常量说明

| 常量 | 值 | 说明 |
| --- | --- | --- |
| `RelationshipChildOf` | `"ChildOf"` | 此弱点是目标弱点的子项（更具体） |
| `RelationshipParentOf` | `"ParentOf"` | 此弱点是目标弱点的父项（更通用） |
| `RelationshipCanPrecede` | `"CanPrecede"` | 此弱点可创建条件使目标弱点成为可能（链式前驱） |
| `RelationshipCanFollow` | `"CanFollow"` | 此弱点可跟随目标弱点（链式后继） |
| `RelationshipRequires` | `"Requires"` | 此复合弱点需要目标弱点存在 |
| `RelationshipRequiredBy` | `"RequiredBy"` | 此弱点被目标复合弱点所需要 |
| `RelationshipCanAlsoBe` | `"CanAlsoBe"` | 在适当上下文中也可视为目标弱点 |
| `RelationshipPeerOf` | `"PeerOf"` | 与目标弱点相似，不适合其他关系类型 |
| `RelationshipMemberOf` | `"MemberOf"` | 此条目是目标类别/视图的成员 |
| `RelationshipHasMember` | `"Has_Member"` | 此类别/视图包含目标条目作为成员 |

::: warning Has_Member 带下划线
`RelationshipHasMember` 的字符串值是 `"Has_Member"`（中间有下划线），其余 9 个均无下划线。这是 MITRE 规范的原始写法，解析/比较时务必注意。
:::

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (r RelationshipNature) String() string` |
| `IsValid` | `func (r RelationshipNature) IsValid() bool` |
| `ParseRelationshipNature` | `func ParseRelationshipNature(s string) (RelationshipNature, error)` |
| `AllRelationshipNatureValues` | `func AllRelationshipNatureValues() []RelationshipNature` |

## 🔀 分类判断方法

`RelationshipNature` 额外提供 4 个**分类判断**方法，把 10 种关系归为 4 类：

| 方法 | 包含的关系 | 语义 |
| --- | --- | --- |
| `IsHierarchical() bool` | ChildOf, ParentOf, MemberOf, Has_Member | 层级关系（父子/成员归属） |
| `IsSequential() bool` | CanPrecede, CanFollow | 顺序关系（链式先后） |
| `IsDependency() bool` | Requires, RequiredBy | 依赖关系（复合元素所需） |
| `IsPeer() bool` | PeerOf, CanAlsoBe | 对等关系（相似/可互换） |

::: tip 10 种关系恰好四类全覆盖
四种分类方法合起来恰好覆盖全部 10 个常量，互不重叠。可用它们做关系归桶处理。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	r, err := cweskills.ParseRelationshipNature("ChildOf")
	fmt.Println(r, err) // ChildOf <nil>

	// 分类判断
	fmt.Println(cweskills.RelationshipChildOf.IsHierarchical()) // true
	fmt.Println(cweskills.RelationshipCanPrecede.IsSequential()) // true
	fmt.Println(cweskills.RelationshipRequires.IsDependency())   // true
	fmt.Println(cweskills.RelationshipPeerOf.IsPeer())           // true

	// Has_Member 的下划线陷阱
	fmt.Println(cweskills.RelationshipHasMember.String()) // Has_Member
	fmt.Println(cweskills.RelationshipHasMember.IsHierarchical()) // true

	// 全部值
	fmt.Println(len(cweskills.AllRelationshipNatureValues())) // 10
}
```

## 🎯 典型用途

<Badge type="tip" text="导航" /> [`CWE.GetParentIDs`](./cwe-relationship-methods) 用 ChildOf 找父级
<Badge type="info" text="归组" /> [`Category`](./category) 用 Has_Member 挂成员
<Badge type="warning" text="链分析" /> 用 CanPrecede/CanFollow 还原攻击链

## ⚠️ 注意事项

::: warning 方向性
关系是**有向**的。「A ChildOf B」表示 A 是 B 的子项，方向从 A 指向 B。`CWE.GetParentIDs` 遍历自身 `Relationships` 里 `Nature==ChildOf` 的条目，返回的是**对方（父级）**的 ID，不是自身。
:::

## 🔗 相关链接

- 字段归宿：`CWE.Relationships`、`Category.Relationships`、`CompoundElement.Relationships` 中每个 `Relationship.Nature`
- 关系获取方法：[CWE 关系获取方法](./cwe-relationship-methods)
- 递归导航：[navigator 概览](./navigator)
- 概念背景：[关系类型 (ChildOf/Requires…)](../guide/concept-relationship)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
