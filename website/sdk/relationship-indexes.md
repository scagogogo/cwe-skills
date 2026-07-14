---
title: 关系索引查询
outline: [2, 3]
---

# 🗃️ 关系索引查询

`Registry` 在 [BuildIndexes](./build-indexes) 之后，暴露一组**基于 ID 的图查询方法**。它们返回 `[]int`（ID 列表），不返回 `*CWE`，调用方需自行 `Get` 取条目。若需直接拿 `[]*CWE`，改用 [Navigator](./navigator)。

::: warning 前置条件
本组方法依赖已构建的索引。调用前请确认 `IndexesBuilt()` 为 `true`，否则结果不可靠。
:::

## 🧬 传递关系查询

| 方法 | 签名 | 含义 |
| --- | --- | --- |
| `GetParentIDs` | `func (r *Registry) GetParentIDs(id int) []int` | 直接父级（更抽象的弱点） |
| `GetChildIDs` | `func (r *Registry) GetChildIDs(id int) []int` | 直接子级（更具体的弱点） |
| `GetPeerIDs` | `func (r *Registry) GetPeerIDs(id int) []int` | 对等弱点（`PeerOf`/`CanAlsoBe`） |
| `GetAncestorIDs` | `func (r *Registry) GetAncestorIDs(id int) []int` | 全部祖先（递归向上，传递闭包） |
| `GetDescendantIDs` | `func (r *Registry) GetDescendantIDs(id int) []int` | 全部后代（递归向下，传递闭包） |

::: tip Ancestor/Descendant 是闭包
`GetAncestorIDs(79)` 返回的是从 79 出发沿 `ChildOf` 边一路向上的**所有**节点，不止直接父级。同理 `GetDescendantIDs` 沿 `ParentOf` 边向下展开全部层级。
:::

## 🗂️ 成员关系查询

| 方法 | 签名 | 含义 |
| --- | --- | --- |
| `GetViewMembers` | `func (r *Registry) GetViewMembers(viewID int) []int` | 视图包含的弱点 ID |
| `GetCategoryMembers` | `func (r *Registry) GetCategoryMembers(categoryID int) []int` | 分类包含的弱点 ID |
| `GetMemberOfIDs` | `func (r *Registry) GetMemberOfIDs(id int) []int` | 该弱点所属的视图/分类 ID |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	// pillar → base → variant
	r.Register(cweskills.NewCWE(703, "Neutralization"))
	b := cweskills.NewCWE(79, "XSS")
	b.Relationships = []cweskills.Relationship{
		{CWEID: 703, Nature: cweskills.RelationshipChildOf},
	}
	r.Register(b)
	v := cweskills.NewCWE(791, "XSS Variant")
	v.Relationships = []cweskills.Relationship{
		{CWEID: 79, Nature: cweskills.RelationshipChildOf},
	}
	r.Register(v)

	r.BuildIndexes()

	fmt.Println(r.GetAncestorIDs(791))    // [79 703]  递归向上
	fmt.Println(r.GetDescendantIDs(703))  // [79 791]  递归向下
	fmt.Println(r.GetParentIDs(79))       // [703]
}
```

## 🆚 与 Navigator 的取舍

| 维度 | 本组方法（Registry） | [Navigator](./navigator) |
| --- | --- | --- |
| 返回类型 | `[]int`（ID） | `[]*CWE`（指针） |
| 额外对象 | 无 | 需 `NewNavigator(r)` |
| 路径/深度查询 | 不支持 | `ShortestPath` / `RelationshipDepth` |
| 适用 | 只需 ID 的轻量查询 | 需要条目本体或图算法 |

## 🔗 相关链接

- 索引构建：[BuildIndexes](./build-indexes)
- 返回 `[]*CWE` 的导航：[Navigator 概览](./navigator)
- 单条目的本地关系：[CWE 关系获取方法](./cwe-relationship-methods)
- 源文件：[`registry.go`](https://github.com/scagogogo/cwe-skills/blob/main/registry.go)
