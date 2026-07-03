---
title: RelationshipNature 关系类型枚举
outline: [2, 3]
---

# 🔄 RelationshipNature — 关系类型枚举

`RelationshipNature` 表示 CWE 条目之间关系的类型。CWE 通过关系把弱点和视图组织成层次、链、复合与对等结构。

## 🧬 类型定义

```go
type RelationshipNature string
```

## 📋 全部取值

| 值 | 常量名 | 含义 |
| --- | --- | --- |
| `"ChildOf"` | `RelationshipChildOf` | 此弱点是目标弱点的子项（更具体） |
| `"ParentOf"` | `RelationshipParentOf` | 此弱点是目标弱点的父项（更通用） |
| `"CanPrecede"` | `RelationshipCanPrecede` | 此弱点可为目标弱点创造条件（链式前驱） |
| `"CanFollow"` | `RelationshipCanFollow` | 此弱点可跟随目标弱点（链式后继） |
| `"Requires"` | `RelationshipRequires` | 此复合弱点需要目标弱点存在 |
| `"RequiredBy"` | `RelationshipRequiredBy` | 此弱点被目标复合弱点所需要 |
| `"CanAlsoBe"` | `RelationshipCanAlsoBe` | 此弱点在适当上下文中也可被视为目标弱点 |
| `"PeerOf"` | `RelationshipPeerOf` | 与目标弱点有相似性，但不适合其他关系 |
| `"MemberOf"` | `RelationshipMemberOf` | 此条目是目标类别/视图的成员 |
| `"Has_Member"` | `RelationshipHasMember` | 此类别/视图包含目标条目作为成员 |

```go
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

::: warning Has_Member 含下划线
注意 `RelationshipHasMember` 的字面量是 `"Has_Member"`（带下划线），是全部取值中唯一含下划线者。解析外部数据时需保留原样，`ParseRelationshipNature("Has Member")` 会失败。
:::

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (r RelationshipNature) String() string` |
| `IsValid` | `func (r RelationshipNature) IsValid() bool` |
| `ParseRelationshipNature` | `func ParseRelationshipNature(s string) (RelationshipNature, error)` |
| `AllRelationshipNatureValues` | `func AllRelationshipNatureValues() []RelationshipNature` |

```go
r, err := cweskills.ParseRelationshipNature("ChildOf")
fmt.Println(r, err)                                      // ChildOf <nil>
fmt.Println(r.String())                                  // ChildOf
fmt.Println(cweskills.RelationshipNature("X").IsValid()) // false
fmt.Println(len(cweskills.AllRelationshipNatureValues())) // 10
```

## 🧬 额外方法：语义分类

`RelationshipNature` 提供四个语义判定方法，把十种关系归入四类：

| 方法 | 签名 | 命中取值 |
| --- | --- | --- |
| `IsHierarchical` | `func (r RelationshipNature) IsHierarchical() bool` | ChildOf、ParentOf、MemberOf、Has_Member |
| `IsSequential` | `func (r RelationshipNature) IsSequential() bool` | CanPrecede、CanFollow |
| `IsDependency` | `func (r RelationshipNature) IsDependency() bool` | Requires、RequiredBy |
| `IsPeer` | `func (r RelationshipNature) IsPeer() bool` | CanAlsoBe、PeerOf |

```go
r := cweskills.RelationshipChildOf
fmt.Println(r.IsHierarchical()) // true
fmt.Println(r.IsSequential())   // false
fmt.Println(r.IsDependency())   // false
fmt.Println(r.IsPeer())         // false
```

::: tip 用途
遍历关系时按语义分类处理：层级关系用于构建父子树，顺序关系用于重建链，依赖关系用于解析复合弱点，对等关系用于发现相似弱点。四类互斥且覆盖全部取值。
:::

## 💻 CLI 对应命令

```bash
cwe enum relationship
```

输出全部合法取值，详见 [CLI enum relationship](../cli/enum-relationship)。

## 🔗 相关链接

- SDK 视角：[RelationshipNature 关系类型枚举](../sdk/enum-relationship-nature)
- 概念背景：[关系类型 (ChildOf/Requires…)](../guide/concept-relationship)
- 关系导航方法：[CWE 关系获取方法](../sdk/cwe-relationship-methods)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
