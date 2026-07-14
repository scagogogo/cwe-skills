---
title: View 视图
outline: [2, 3]
---

# 👁️ View — 视图

`View` 是 CWE 体系里的「视角」条目，提供从特定角度查看和组织弱点的方式。每个视图有一个 [`ViewType`](./enum-view-type)（Graph/Explicit Slice/Implicit Slice），并通过 `Members []ViewMember` 显式列出归属的 CWE。

## 📋 结构体定义

```go
type View struct {
    ID            int           `json:"id" xml:"ID,attr"`
    Name          string        `json:"name" xml:"Name"`
    Type          ViewType      `json:"type,omitempty" xml:"Type,omitempty"`
    Status        Status        `json:"status,omitempty" xml:"Status,omitempty"`
    Description   string        `json:"description" xml:"Description"`
    Members       []ViewMember  `json:"members,omitempty" xml:"Members>ViewMember,omitempty"`
    References    []Reference   `json:"references,omitempty" xml:"References>Reference,omitempty"`
    ContentHistory *ContentHistory `json:"content_history,omitempty" xml:"ContentHistory,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `ID` | `int` | 数字标识符 |
| `Name` | `string` | 视图名称 |
| `Type` | [`ViewType`](./enum-view-type) | 视图类型 |
| `Status` | [`Status`](./enum-status) | 状态 |
| `Description` | `string` | 视图描述 |
| `Members` | `[]ViewMember` | 视图成员列表 |
| `References` | [`[]Reference`](./reference) | 参考文献 |
| `ContentHistory` | [`*ContentHistory`](./content-history) | 内容历史 |

## 🔧 ViewMember

```go
type ViewMember struct {
    CWEID     int    `json:"cwe_id" xml:"CWE_ID"`
    ViewID    int    `json:"view_id" xml:"View_ID"`
    Direct    bool   `json:"direct" xml:"Direct"`
    Predicate string `json:"predicate,omitempty" xml:"Predicate,omitempty"`
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `CWEID` | `int` | 成员的 CWE ID |
| `ViewID` | `int` | 所属视图 ID |
| `Direct` | `bool` | 是否为直接成员 |
| `Predicate` | `string` | 谓词（可选，用于隐式切片的过滤条件） |

::: tip Direct 的含义
图类型视图中，成员分「直接」与「间接」：直接成员是显式列入视图的弱点，间接成员是通过关系图传递纳入的弱点。`Direct=true` 表示前者。
:::

## 🏗️ 构造器

```go
func NewView(id int, name string, viewType ViewType) *View
```

创建最小可用的 `View`，设置 `ID`、`Name`、`Type`。

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	view := cweskills.NewView(635, "Software Development", cweskills.ViewTypeGraph)
	view.Status = cweskills.StatusStable
	view.Description = "软件开发视角的弱点视图。"
	view.Members = []cweskills.ViewMember{
		{CWEID: 79, ViewID: 635, Direct: true},
		{CWEID: 89, ViewID: 635, Direct: true},
	}

	fmt.Printf("视图 %s (type=%s)\n", view.Name, view.Type)
	for _, m := range view.Members {
		fmt.Printf("  成员: CWE-%d (direct=%v)\n", m.CWEID, m.Direct)
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="知名列表" /> CWE Top 25、OWASP Top 10 等都是 View（多为 Explicit Slice）
<Badge type="info" text="导航" /> 通过 `Members` 找到一整组弱点
<Badge type="warning" text="过滤" /> 按 `ViewType` 区分图视图与切片视图的处理方式

## ⚠️ 注意事项

::: warning View 没有 Notes 字段
与 `Category` 不同，`View` 结构体**没有** `Notes` 字段。如需展示备注信息，请从 `Description` 或 `References` 获取。
:::

::: details 三种 ViewType 的区别
- **Graph**：层次化关系图，如 CWE-1000 研究概念视图、CWE-699 软件开发视图。成员可继承自父级。
- **Explicit Slice**：显式扁平列表，成员逐个列出，如 CWE Top 25。
- **Implicit Slice**：隐式扁平列表，通过过滤器/属性定义，如「所有 Draft 状态的条目」。`Predicate` 字段常用于描述过滤条件。
:::

## 🔗 相关链接

- 配套类别（对比）：[Category 类别](./category)
- 视图类型枚举：[ViewType](./enum-view-type)
- 知名视图实战：[知名视图](../wellknown/well-known-views)
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
