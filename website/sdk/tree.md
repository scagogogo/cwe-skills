---
title: 树构建概览
outline: [2, 3]
---

# 🌳 树构建概览

`cweskills` 提供 `TreeNode` 与三个构建函数，把 `Registry` 中扁平的弱点关系网络投影成**层级树**。树结构便于遍历、渲染、统计深度与叶子，是 CLI 展示与可视化组件的基础。

## 🧩 核心类型

```go
type TreeNode struct {
    CWE      *CWE
    Children []*TreeNode
    Parent   *TreeNode
}
```

| 字段 | 说明 |
| --- | --- |
| `CWE` | 该节点对应的弱点指针 |
| `Children` | 子节点切片（更具体的弱点） |
| `Parent` | 父节点指针（根节点为 `nil`） |

详见 [TreeNode 节点](./tree-node)。

## 🏗️ 三个构建函数

| 函数 | 签名 | 产出 | 文档 |
| --- | --- | --- | --- |
| `BuildTree` | `func BuildTree(r *Registry, rootID int) *TreeNode` | 以指定弱点为根的子树 | [构建单棵树](./build-tree) |
| `BuildForest` | `func BuildForest(r *Registry) []*TreeNode` | 全库森林（多棵根树） | [构建森林](./build-forest) |
| `BuildViewTree` | `func BuildViewTree(r *Registry, viewID int) *TreeNode` | 视图成员投影的树 | [构建视图树](./build-view-tree) |

::: tip 依赖已构建索引
三个构建函数都读取 `Registry` 的父/子索引。调用前确保 `r.BuildIndexes()` 已执行，否则树为空。
:::

## 📚 本组文档导航

| 文档 | 主题 |
| --- | --- |
| [TreeNode 节点](./tree-node) | 结构体定义、`NewTreeNode`、`AddChild` |
| [构建单棵树](./build-tree) | `BuildTree` |
| [构建森林](./build-forest) | `BuildForest` |
| [构建视图树](./build-view-tree) | `BuildViewTree` |
| [树遍历](./tree-walk) | `Walk` / `WalkBFS` |
| [路径查询](./tree-path) | `Path` / `Find` |
| [叶子节点](./tree-leaf-nodes) | `LeafNodes` / `IsLeaf` |
| [深度与计数](./tree-depth-count) | `MaxDepth` / `Count` |

## ✅ 快速上手

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	r.Register(cweskills.NewCWE(703, "Neutralization"))
	c := cweskills.NewCWE(79, "XSS")
	c.Relationships = []cweskills.Relationship{
		{CWEID: 703, Nature: cweskills.RelationshipChildOf},
	}
	r.Register(c)
	r.BuildIndexes()

	root := cweskills.BuildTree(r, 703)
	fmt.Println(root.CWE.Name)        // Neutralization
	fmt.Println(root.Count())         // 2
	fmt.Println(root.MaxDepth())      // 1
}
```

## ⚠️ 注意事项

::: warning 树是投影，非完整图
树只保留 `ChildOf`/`ParentOf` 构成的层级，`PeerOf`/`Requires` 等横向关系不会进入树。如需完整关系，配合 [Navigator](./navigator) 使用。
:::

## 🔗 相关链接

- 依赖的索引：[构建索引](./build-indexes)
- 节点方法：[TreeNode 节点](./tree-node)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
