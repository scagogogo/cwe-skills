---
title: 注册表与索引
outline: [2, 3]
---

# 🗃️ 注册表与索引

`Registry` 是 `cweskills` 包的**内存数据中枢**——一个线程不安全但功能完备的 CWE 仓库。它汇聚了弱点（`CWE`）、分类（`Category`）、视图（`View`）和复合元素（`CompoundElement`）四类条目，并提供注册、查询、计数、索引、序列化等全套能力。

## 🧩 在 SDK 中的位置

```text
数据源(XML/HTTP) → 解析 → CWE/Category/View/CompoundElement
                                   │
                                   ▼
                            ┌─────────────┐
                            │  Registry   │  ← 本文档主题
                            └─────────────┘
                                   │
              ┌─────────┬──────────┼──────────┬──────────┐
              ▼         ▼          ▼          ▼          ▼
          Navigator   Tree     Search      Filter     Statistics
              │         │          │          │          │
              └─────────┴──────────┴──────────┴──────────┘
                                   │
                                   ▼
                            Serializer(JSON/XML/CSV)
```

几乎所有高级功能（导航、树构建、搜索、统计）都以 `*Registry` 作为输入参数。

## 🏗️ 构造器

```go
func NewRegistry() *Registry
```

创建一个空的注册表。新实例不含任何条目，索引也未构建。

## 📚 本组文档导航

| 文档 | 主题 | 核心方法 |
| --- | --- | --- |
| [注册表基础操作](./registry-operations) | 注册、查询、计数、删除 | `Register` / `Get` / `Size` / `Remove` / `Clear` |
| [构建索引](./build-indexes) | 索引生命周期与触发 | `BuildIndexes` / `IndexesBuilt` |
| [关系索引查询](./relationship-indexes) | 基于 ID 的图查询 | `GetParentIDs` / `GetChildIDs` / `GetAncestorIDs` / `GetViewMembers` |
| [注册表 JSON 导入导出](./registry-json) | 整库序列化 | `ExportJSON` / `ImportJSON` |

## ✅ 快速上手

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	cwe := cweskills.NewCWE(79, "XSS")
	cwe.Abstraction = cweskills.AbstractionBase
	_ = r.Register(cwe)

	fmt.Println(r.Size())               // 1
	fmt.Println(r.Contains(79))         // true
	if got, ok := r.Get(79); ok {
		fmt.Println(got.Name)          // XSS
	}

	r.BuildIndexes()
	fmt.Println(r.IndexesBuilt())       // true
}
```

## ⚠️ 使用须知

::: warning 非并发安全
`Registry` 内部用普通 map 存储，**没有锁**。多协程并发读写同一个注册表需由调用方自行加锁（如 `sync.RWMutex`）。
:::

::: tip 索引需显式构建
注册/删除条目后，关系索引不会自动更新。调用 `BuildIndexes()` 重建后才能使用 `GetParentIDs` 等索引查询方法。详见 [构建索引](./build-indexes)。
:::

## 🔗 相关链接

- 关系导航（基于注册表）：[Navigator 概览](./navigator)
- 树构建（基于注册表）：[TreeNode 与构建](./tree)
- 单条目序列化：[Serializer 概览](./serializer)
- 源文件：[`registry.go`](https://github.com/scagogogo/cwe-skills/blob/main/registry.go)
