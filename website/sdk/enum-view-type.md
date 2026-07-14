---
title: ViewType 视图类型枚举
outline: [2, 3]
---

# 📚 ViewType — 视图类型枚举

`ViewType` 表示 CWE 视图的组织方式，共 3 个取值：**Graph / Explicit Slice / Implicit Slice**。它是 `View.Type` 字段的类型。

## 📋 类型与常量

```go
type ViewType string

const (
	ViewTypeGraph         ViewType = "Graph"
	ViewTypeExplicitSlice ViewType = "Explicit Slice"
	ViewTypeImplicitSlice ViewType = "Implicit Slice"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 | 示例 |
| --- | --- | --- | --- |
| `ViewTypeGraph` | `"Graph"` | 图类型视图，层次化关系表示 | CWE-1000 研究概念、CWE-699 软件开发 |
| `ViewTypeExplicitSlice` | `"Explicit Slice"` | 显式切片，外部因素相关的扁平列表 | CWE Top 25、OWASP Top Ten |
| `ViewTypeImplicitSlice` | `"Implicit Slice"` | 隐式切片，过滤器/属性定义的扁平列表 | 所有 Draft 状态的条目 |

::: tip 三种视图的本质区别
- **Graph**：成员通过关系图组织，可继承；`ViewMember.Direct` 区分直接/间接成员。
- **Explicit Slice**：成员逐个显式列出，是「清单」式视图。
- **Implicit Slice**：成员由谓词（`ViewMember.Predicate`）动态定义，是「查询」式视图。
:::

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (v ViewType) String() string` |
| `IsValid` | `func (v ViewType) IsValid() bool` |
| `ParseViewType` | `func ParseViewType(s string) (ViewType, error)` |
| `AllViewTypeValues` | `func AllViewTypeValues() []ViewType` |

::: warning 两个值带空格
`ViewTypeExplicitSlice` 是 `"Explicit Slice"`，`ViewTypeImplicitSlice` 是 `"Implicit Slice"`，中间都有空格。解析时必须完全匹配。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	vt, err := cweskills.ParseViewType("Graph")
	fmt.Println(vt, err) // Graph <nil>

	// 校验
	fmt.Println(cweskills.ViewTypeExplicitSlice.IsValid()) // true
	fmt.Println(cweskills.ViewType("foo").IsValid())        // false

	// 创建视图时指定类型
	view := cweskills.NewView(1000, "Research Concepts", cweskills.ViewTypeGraph)
	fmt.Println(view.Type) // Graph

	// 全部值
	for _, v := range cweskills.AllViewTypeValues() {
		fmt.Println(v)
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="知名列表" /> CWE Top 25、OWASP Top 10 通常是 Explicit Slice
<Badge type="info" text="导航" /> Graph 视图配合 [BuildViewTree](./build-view-tree) 构建层次树
<Badge type="warning" text="过滤" /> 按 ViewType 选择不同的成员解析策略

## ⚠️ 注意事项

::: warning Graph 视图需配合关系数据
`ViewTypeGraph` 视图的完整成员可能不在 `View.Members` 里全部列出——部分通过关系图传递纳入。处理 Graph 视图时建议配合 [Registry](./registry) 与 [Navigator](./navigator) 做递归展开。
:::

## 🔗 相关链接

- 字段归宿：`View.Type`，见 [View 视图](./view)
- 构造器：`NewView(id, name, viewType)`
- 视图树构建：[BuildViewTree](./build-view-tree)
- 概念背景：[视图 (View)](../guide/concept-view)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
