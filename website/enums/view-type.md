---
title: ViewType 视图类型枚举
outline: [2, 3]
---

# 🗺️ ViewType — 视图类型枚举

`ViewType` 表示 CWE 视图（View）的组织类型。视图是对 CWE 条目的特定切分，用于不同分析视角。

## 🧬 类型定义

```go
type ViewType string
```

## 📋 全部取值

| 值 | 常量名 | 含义 | 示例视图 |
| --- | --- | --- | --- |
| `"Graph"` | `ViewTypeGraph` | 图类型，层次化的关系表示 | CWE-1000 研究概念、CWE-699 软件开发 |
| `"Explicit Slice"` | `ViewTypeExplicitSlice` | 显式切片，外部因素相关的扁平列表 | CWE Top 25、OWASP Top Ten |
| `"Implicit Slice"` | `ViewTypeImplicitSlice` | 隐式切片，过滤器/属性定义的扁平列表 | 所有草稿状态条目 |

```go
const (
	ViewTypeGraph         ViewType = "Graph"
	ViewTypeExplicitSlice ViewType = "Explicit Slice"
	ViewTypeImplicitSlice ViewType = "Implicit Slice"
)
```

::: tip 三种视图的差异
**Graph** 用父子关系构建层次树，适合漫游整个 CWE 体系；**Explicit Slice** 是人工挑选的扁平列表（如年度 Top 25），成员固定；**Implicit Slice** 由过滤器动态生成（如「所有 Draft 状态」），成员随数据变化。前两者是静态的，后者是动态的。
:::

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (v ViewType) String() string` |
| `IsValid` | `func (v ViewType) IsValid() bool` |
| `ParseViewType` | `func ParseViewType(s string) (ViewType, error)` |
| `AllViewTypeValues` | `func AllViewTypeValues() []ViewType` |

```go
v, err := cweskills.ParseViewType("Graph")
fmt.Println(v, err)                              // Graph <nil>
fmt.Println(v.String())                          // Graph
fmt.Println(cweskills.ViewType("X").IsValid())   // false
fmt.Println(cweskills.AllViewTypeValues())       // [Graph Explicit Slice Implicit Slice]
```

::: warning 取值含空格
`"Explicit Slice"` 与 `"Implicit Slice"` 含空格，`ParseViewType` 严格匹配原样。从 MITRE XML 读取时字面量已正确，但从用户输入解析时需注意不要误写成 `ExplicitSlice`。
:::

## 🔄 典型用法

```go
// 区分图视图与切片视图，采用不同遍历策略
view, _ := registry.GetView("CWE-1000")
switch view.Type {
case cweskills.ViewTypeGraph:
	// 用 tree 构建层次树
	buildTreeView(view)
case cweskills.ViewTypeExplicitSlice, cweskills.ViewTypeImplicitSlice:
	// 切片视图直接列举成员
	listMembers(view)
}
```

## 💻 CLI 对应命令

```bash
cwe enum view-type
```

输出全部合法取值，详见 [CLI enum view-type](../cli/enum-view-type)。

## 🔗 相关链接

- SDK 视角：[ViewType 视图类型枚举](../sdk/enum-view-type)
- 概念背景：[视图 (View)](../guide/concept-view)
- 视图数据模型：[View 视图](../sdk/view)
- 按视图构建树：[`BuildViewTree`](../sdk/build-view-tree)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
