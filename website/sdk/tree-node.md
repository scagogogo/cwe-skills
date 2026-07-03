---
title: TreeNode 节点
outline: [2, 3]
---

# 🌳 TreeNode 节点

`TreeNode` 是树构建的基本单元，封装一个 `CWE` 及其父子链接。它既是 `BuildTree` 等函数的返回类型，也可由调用方手动拼装。

## 📋 结构体定义

```go
type TreeNode struct {
    CWE      *CWE
    Children []*TreeNode
    Parent   *TreeNode
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `CWE` | `*CWE` | 该节点对应的弱点 |
| `Children` | `[]*TreeNode` | 子节点列表 |
| `Parent` | `*TreeNode` | 父节点；根节点为 `nil` |

## 🏗️ 构造与拼接

### NewTreeNode

```go
func NewTreeNode(cwe *CWE) *TreeNode
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cwe` | `*CWE` | 节点承载的弱点 |
| 返回 | `*TreeNode` | 无子节点、无父节点的新节点 |

### AddChild

```go
func (n *TreeNode) AddChild(child *TreeNode)
```

把 `child` 挂到 `n` 的 `Children` 末尾，并设置 `child.Parent = n`。

## 🔎 状态判断

| 方法 | 签名 | 返回 true 当 |
| --- | --- | --- |
| `IsLeaf` | `func (n *TreeNode) IsLeaf() bool` | `Children` 为空 |
| `IsRoot` | `func (n *TreeNode) IsRoot() bool` | `Parent` 为 `nil` |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	root := cweskills.NewTreeNode(cweskills.NewCWE(703, "Neutralization"))
	child := cweskills.NewTreeNode(cweskills.NewCWE(79, "XSS"))
	root.AddChild(child)

	fmt.Println(root.IsRoot())        // true
	fmt.Println(child.IsRoot())       // false
	fmt.Println(child.IsLeaf())       // true
	fmt.Println(root.Children[0].CWE.Name) // XSS
	fmt.Println(child.Parent.CWE.Name)     // Neutralization
}
```

## 📝 字符串表示

```go
func (n *TreeNode) String() string
```

返回该节点的可读形式（通常含 CWE ID 与名称），便于调试输出。

## 🔗 相关链接

- 遍历方法：[树遍历](./tree-walk)
- 路径与查找：[路径查询](./tree-path)
- 自动构建：[构建单棵树](./build-tree)
- 源文件：[`tree.go`](https://github.com/scagogogo/cwe-skills/blob/main/tree.go)
