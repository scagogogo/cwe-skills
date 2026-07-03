---
title: 技能 11 — 本地树构建
outline: [2, 3]
---

# 🌳 技能 11 — 本地树构建

从离线 XML 数据构建并遍历 CWE 层次树。可视化完整的弱点分类体系。

<Badge type="tip" text="离线"/>
<Badge type="info" text="需 XML 目录"/>

---

## 🎯 技能目标

- 从根节点构建层次树
- 构建森林（所有柱状根节点）
- 按视图构建树
- 查找从根到某 CWE 的路径、列出叶子节点

---

## 💻 CLI 命令

所有 tree 命令需 `--xml <file>`。

```bash
# 从根节点构建树
cwe tree build CWE-1 --xml <file>

# 构建森林（所有柱状节点）
cwe tree forest --xml <file>

# 按视图构建树
cwe tree view 1000 --xml <file>

# 从根到某 CWE 的路径
cwe tree path CWE-79 --xml <file>
cwe tree path CWE-79 --xml <file> --root 1

# 列出某根下所有叶子
cwe tree leaves CWE-1 --xml <file>
```

| Flag | 简写 | 说明 |
|------|------|------|
| `--xml` | `-x` | **（必填）** XML 目录路径 |
| `--root` | | `path` 命令的根节点 ID（省略则自动检测） |

---

## 🔧 SDK API

### 构建树

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
registry.BuildIndexes()

tree := cweskills.BuildTree(registry, 1)          // *TreeNode
forest := cweskills.BuildForest(registry)           // []*TreeNode
viewTree := cweskills.BuildViewTree(registry, 1000) // *TreeNode
```

### 遍历

```go
tree.Walk(func(node *cweskills.TreeNode) bool {
    fmt.Printf("%s%s\n", strings.Repeat("  ", node.Depth), node.CWE.Name)
    return true
})

tree.WalkBFS(func(node *cweskills.TreeNode) bool {
    return true
})
```

`Walk` 是深度优先，`WalkBFS` 是广度优先。

### 查询

```go
path := tree.Find(79).Path()    // []*TreeNode，从根到 79
leaves := tree.LeafNodes()       // []*TreeNode
maxDepth := tree.MaxDepth()      // int
count := tree.Count()            // int
isLeaf := tree.IsLeaf()          // bool
```

::: details TreeNode 方法
`Walk` · `WalkBFS` · `Find(id)` · `Path()` · `LeafNodes()` · `MaxDepth()` · `Count()` · `IsLeaf()`。每个节点有 `CWE`、`Children`、`Depth`、`Parent` 字段。
:::

---

## 📝 示例

### 命令行

```bash
# 看 CWE-1 下的层次树
cwe tree build CWE-1 --xml cwec_latest.xml

# 从根到 CWE-79 的路径
cwe tree path CWE-79 --xml cwec_latest.xml -o json
```

### Go

```go
package main

import (
    "fmt"
    "strings"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    registry, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
    registry.BuildIndexes()

    tree := cweskills.BuildTree(registry, 1)
    tree.Walk(func(n *cweskills.TreeNode) bool {
        if n.Depth <= 2 {
            fmt.Printf("%sCWE-%d %s\n", strings.Repeat("  ", n.Depth), n.CWE.ID, n.CWE.Name)
        }
        return true
    })
}
```

---

## 🤖 AI 代理使用提示

- 用户想「看 CWE 分类树」时，AI 用 `cwe tree build` 或 `cwe tree forest`。
- 查某 CWE 在树里的位置用 `cwe tree path`。
- 列出某根下所有具体弱点（叶子）用 `cwe tree leaves`。

::: tip 视图树
`cwe tree view 1000` 用 Research Concepts 视图构建树，是最常引用的标准层次视图。
:::

---

## 📖 相关文档

- [技能 10 — 本地关系导航](./10-local-navigation)
- [CLI: tree](../cli/tree) · [tree build](../cli/tree-build) · [tree forest](../cli/tree-forest) · [tree path](../cli/tree-path) · [tree leaves](../cli/tree-leaves)
- [SDK: BuildTree](../sdk/build-tree) · [BuildForest](../sdk/build-forest) · [TreeNode](../sdk/tree-node) · [Walk](../sdk/tree-walk)
- [返回 Skills 总览](./)
