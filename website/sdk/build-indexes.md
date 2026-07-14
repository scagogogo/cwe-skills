---
title: 构建索引
outline: [2, 3]
---

# 🗃️ 构建索引

`Registry` 维护一套**关系索引**（父→子、子→父、对等、祖先、后代、视图成员等）。索引不会在每次注册时自动维护，而是由调用方在数据准备就绪后显式触发 `BuildIndexes()`。

## 🛠️ 核心方法

### BuildIndexes

```go
func (r *Registry) BuildIndexes()
```

| 项 | 说明 |
| --- | --- |
| 参数 | 无 |
| 返回 | 无 |
| 作用 | 扫描全部 `CWE` 的 `Relationships` 字段，构建父/子/对等/祖先/后代及视图/分类成员索引 |
| 复杂度 | O(N + E)，N 为弱点数、E 为关系总数 |

### IndexesBuilt

```go
func (r *Registry) IndexesBuilt() bool
```

| 项 | 说明 |
| --- | --- |
| 返回 | 索引是否已构建（`true` 表示可安全调用索引查询方法） |

::: warning 调用索引查询前必须先构建
若未调用 `BuildIndexes()` 而直接调用 `GetParentIDs` 等[关系索引查询](./relationship-indexes)方法，可能得到空结果或 panic。务必先 `BuildIndexes()` 或用 `IndexesBuilt()` 守卫。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	parent := cweskills.NewCWE(703, "Improper Neutralization")
	child := cweskills.NewCWE(79, "XSS")
	child.Relationships = []cweskills.Relationship{
		{CWEID: 703, Nature: cweskills.RelationshipChildOf},
	}
	_ = r.Register(parent)
	_ = r.Register(child)

	fmt.Println(r.IndexesBuilt()) // false
	r.BuildIndexes()
	fmt.Println(r.IndexesBuilt()) // true

	fmt.Println(r.GetChildIDs(703)) // [79]
}
```

## 🔄 何时重建索引

::: tip 数据变更后重建
下列操作会改变关系结构，执行后应再次调用 `BuildIndexes()`：
- `Register` / `Remove` 增删弱点
- 直接修改已注册 `CWE` 的 `Relationships` 字段
- `ImportJSON` 导入新数据（导入后会自动重建，无需手动调用）
:::

::: details 为什么不自动维护？
关系索引是**全局视图**，单条注册时无法高效更新祖先/后代这类传递闭包。批量构建一次的均摊成本远低于每次注册都维护，因此设计为显式触发。这也是 `Registry` 选择「不并发加锁」的代价之一。
:::

## 🔗 相关链接

- 索引查询方法：[关系索引查询](./relationship-indexes)
- 注册/删除：[注册表基础操作](./registry-operations)
- 基于索引的图导航：[Navigator 概览](./navigator)
- 源文件：[`registry.go`](https://github.com/scagogogo/cwe-skills/blob/main/registry.go)
