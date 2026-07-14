---
title: Category 类别
outline: [2, 3]
---

# 📂 Category — 类别

`Category` 是 CWE 体系里的「分类」条目，用于把相关弱点按主题归组。它不是弱点本身，而是一个聚合容器。通过 `Relationships`（`MemberOf` / `Has_Member`）挂载成员弱点。

## 📋 结构体定义

```go
type Category struct {
    ID            int           `json:"id" xml:"ID,attr"`
    Name          string        `json:"name" xml:"Name"`
    Status        Status        `json:"status,omitempty" xml:"Status,omitempty"`
    Description   string        `json:"description" xml:"Description"`
    Relationships []Relationship `json:"relationships,omitempty" xml:"Relationships>Relationship,omitempty"`
    Notes         string        `json:"notes,omitempty" xml:"Notes,omitempty"`
    References    []Reference   `json:"references,omitempty" xml:"References>Reference,omitempty"`
    ContentHistory *ContentHistory `json:"content_history,omitempty" xml:"ContentHistory,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `ID` | `int` | 数字标识符 |
| `Name` | `string` | 类别名称 |
| `Status` | [`Status`](./enum-status) | 状态 |
| `Description` | `string` | 类别描述 |
| `Relationships` | `[]Relationship` | 成员关系（含 `MemberOf`/`Has_Member`） |
| `Notes` | `string` | 备注 |
| `References` | [`[]Reference`](./reference) | 参考文献 |
| `ContentHistory` | [`*ContentHistory`](./content-history) | 内容历史 |

::: tip Category 没有 CWEType
与 `CWE` 不同，`Category` 是独立结构体而非通过 `CWEType` 区分。它也没有 `Abstraction`/`Structure`/`LikelihoodOfExploit` 等弱点专属字段——类别只是归组，本身不是弱点。
:::

## 🏗️ 构造器

```go
func NewCategory(id int, name string) *Category
```

创建最小可用的 `Category`，仅设置 `ID` 与 `Name`。

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	cat := cweskills.NewCategory(1000, "Research Views")
	cat.Status = cweskills.StatusStable
	cat.Description = "研究视角下的弱点分类。"
	cat.Relationships = []cweskills.Relationship{
		{CWEID: 79, Nature: cweskills.RelationshipHasMember},
		{CWEID: 89, Nature: cweskills.RelationshipHasMember},
	}

	fmt.Printf("类别 %s (ID=%d)\n", cat.Name, cat.ID)
	// 列出成员
	for _, r := range cat.Relationships {
		fmt.Printf("  成员: CWE-%d (%s)\n", r.CWEID, r.Nature)
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="归组" /> 按主题（如「输入校验」「并发问题」）组织弱点
<Badge type="info" text="导航" /> 通过 `Has_Member` 关系找到一整类弱点
<Badge type="warning" text="报告" /> 按类别汇总扫描结果

## ⚠️ 注意事项

::: warning 与 View 的区别
- `Category` 用 `Relationships`（`Has_Member`/`MemberOf`）表达成员关系，结构较扁平。
- [`View`](./view) 用 `Members []ViewMember` 显式列出成员，并带 `ViewType`（Graph/Explicit Slice/Implicit Slice）描述组织方式。
两者都做归组，但 View 更强调「视角」与「切片」语义。
:::

## 🔗 相关链接

- 配套视图：[View 视图](./view)
- 关系类型：[RelationshipNature](./enum-relationship-nature)
- 参考文献：[Reference](./reference)、内容历史：[ContentHistory](./content-history)
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
