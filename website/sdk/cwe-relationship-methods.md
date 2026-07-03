---
title: CWE 关系获取方法
outline: [2, 3]
---

# 🧭 CWE 关系获取方法

`CWE` 结构体提供四个方法，从 `Relationships` 字段中按关系类型筛选出相关的 CWE ID 列表。它们是关系导航的**就地查询**入口——只遍历当前条目自身的关系集合，不递归、不查注册表。

## 📋 方法清单

| 方法 | 签名 | 筛选的 Nature | 返回 |
| --- | --- | --- | --- |
| `GetParentIDs` | `func (c *CWE) GetParentIDs() []int` | `ChildOf` | 父级 ID 列表 |
| `GetChildIDs` | `func (c *CWE) GetChildIDs() []int` | `ParentOf` | 子级 ID 列表 |
| `GetPeerIDs` | `func (c *CWE) GetPeerIDs() []int` | `PeerOf`、`CanAlsoBe` | 对等 ID 列表 |
| `GetChainIDs` | `func (c *CWE) GetChainIDs() []int` | `CanPrecede`、`CanFollow` | 链式 ID 列表 |

::: tip 语义解释
- `GetParentIDs`：找的是「我是谁的子项」 → 筛选 `ChildOf` 关系，返回对方（更通用的父级）的 ID
- `GetChildIDs`：找的是「我是谁的父项」 → 筛选 `ParentOf` 关系，返回对方（更具体的子级）的 ID
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.Relationships = []cweskills.Relationship{
		{CWEID: 74, Nature: cweskills.RelationshipChildOf},      // 79 是 74 的子项
		{CWEID: 732, Nature: cweskills.RelationshipChildOf},     // 79 是 732 的子项
		{CWEID: 83, Nature: cweskills.RelationshipParentOf},     // 83 是 79 的子项
		{CWEID: 80, Nature: cweskills.RelationshipPeerOf},       // 80 与 79 对等
		{CWEID: 116, Nature: cweskills.RelationshipCanPrecede},  // 116 可在 79 之前
	}

	fmt.Println(cwe.GetParentIDs()) // [74 732]
	fmt.Println(cwe.GetChildIDs())  // [83]
	fmt.Println(cwe.GetPeerIDs())   // [80]
	fmt.Println(cwe.GetChainIDs())  // [116]
}
```

## ⚠️ 边界行为

::: warning 返回 nil 而非空切片
当没有任何匹配关系时，这些方法返回 `nil`（因为内部用 `var ids []int` 声明后 append）。遍历 `nil` 切片是安全的，但若要把长度当指标，注意 `len(nil) == 0`。
:::

::: details 为什么只看本地关系？
这些方法只读 `c.Relationships`，即当前条目自带的关系列表。它们**不会**：
- 递归向上找祖先（需用 [Navigator.Ancestors](./navigator)）
- 递归向下找后代（需用 [Navigator.Descendants](./navigator)）
- 跨注册表查对端条目是否存在

如需图遍历，请配合 [Registry](./registry) 与 [Navigator](./navigator)。
:::

## 🆚 与 Navigator 的区别

| 维度 | 本组方法 | [Navigator](./navigator) |
| --- | --- | --- |
| 数据源 | 单个 `CWE.Relationships` | 整个 `Registry` 的关系索引 |
| 是否递归 | ❌ 否 | ✅ 可递归到祖先/后代 |
| 是否需注册 | ❌ 否 | ✅ 需要 `*Registry` |
| 返回 | `[]int`（ID） | `[]*CWE` 或 `[]int` |
| 适用 | 快速就地查询 | 图遍历、路径查找 |

## 🔗 相关链接

- 关系类型枚举：[RelationshipNature](./enum-relationship-nature)
- `Relationship` 结构：定义于 `relationship.go`
- 递归导航：[navigator 概览](./navigator)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
