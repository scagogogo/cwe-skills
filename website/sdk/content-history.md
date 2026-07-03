---
title: ContentHistory 内容历史
outline: [2, 3]
---

# 📜 ContentHistory — 内容历史

`ContentHistory` 记录 CWE 条目的提交与修改轨迹，是 `CWE.ContentHistory` 字段的类型（指针，可空）。它包含首次提交信息与历次修改记录，便于追溯条目的演变。

## 📋 结构体定义

```go
type ContentHistory struct {
    Submission     *HistoryEntry  `json:"submission,omitempty" xml:"Submission,omitempty"`
    Modifications  []HistoryEntry `json:"modifications,omitempty" xml:"Modifications>Modification,omitempty"`
}
```

## 📝 字段说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Submission` | `*HistoryEntry` | 首次提交信息 |
| `Modifications` | `[]HistoryEntry` | 修改记录列表 |

## 🔧 HistoryEntry

```go
type HistoryEntry struct {
    Name         string `json:"name,omitempty" xml:"Name,omitempty"`
    Organization string `json:"organization,omitempty" xml:"Organization,omitempty"`
    Date         string `json:"date,omitempty" xml:"Date,omitempty"`
    Comment      string `json:"comment,omitempty" xml:"Comment,omitempty"`
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Name` | `string` | 提交者/修改者姓名 |
| `Organization` | `string` | 所属组织 |
| `Date` | `string` | 日期（字符串形式） |
| `Comment` | `string` | 本次提交/修改的注释 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.ContentHistory = &cweskills.ContentHistory{
		Submission: &cweskills.HistoryEntry{
			Name:         "MITRE",
			Organization: "MITRE Corporation",
			Date:         "2006-07-19",
			Comment:      "首次提交。",
		},
		Modifications: []cweskills.HistoryEntry{
			{
				Name:    "CWE Content Team",
				Date:    "2023-06-29",
				Comment: "更新了演示示例与适用平台。",
			},
		},
	}

	if cwe.ContentHistory != nil && cwe.ContentHistory.Submission != nil {
		s := cwe.ContentHistory.Submission
		fmt.Printf("提交: %s @ %s (%s)\n", s.Name, s.Date, s.Organization)
	}
	for _, m := range cwe.ContentHistory.Modifications {
		fmt.Printf("修改: %s @ %s — %s\n", m.Name, m.Date, m.Comment)
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="溯源" /> 判断条目是否近期被维护
<Badge type="info" text="审计" /> 在合规场景展示弱点定义的演变
<Badge type="warning" text="过滤" /> 按修改日期筛出最近更新的条目

## ⚠️ 注意事项

::: warning 指针类型，逐层判空
`CWE.ContentHistory` 是指针，且其内部 `Submission` 也是指针。访问历史前需**两重判空**：先 `cwe.ContentHistory != nil`，再 `cwe.ContentHistory.Submission != nil`。`Modifications` 切片则可能为 `nil` 或空，遍历是安全的。
:::

::: details Date 是字符串而非时间类型
`Date` 以字符串存储（如 `"2006-07-19"`），SDK 不强制解析为 `time.Time`，因为 MITRE 原始日期格式可能不一致。需要时间运算时请自行 `time.Parse`。
:::

## 🔗 相关链接

- 宿主字段：`CWE.ContentHistory`，也用于 `Category.ContentHistory`、`View.ContentHistory`
- 模型概览：[model.go](./model)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)
