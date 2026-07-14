---
title: 顶层查找
outline: [2, 3]
---

# 🔍 顶层查找

`FindTopLevel` 返回注册表中所有**无父级**的顶层弱点——层级树的根。它是 [BuildForest](./build-forest) 的前置查询，也是「从最高抽象层开始浏览」的入口。

## 📐 函数签名

```go
func FindTopLevel(r *Registry) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| 返回 | `[]*CWE` | 没有 `ChildOf` 关系（无父级）的弱点列表 |

::: warning 依赖索引
`FindTopLevel` 通过 `GetParentIDs` 判断是否有父级，**必须**先 `BuildIndexes()`，否则可能把所有弱点都判为顶层。
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
	r.Register(cweskills.NewCWE(703, "Neutralization")) // 顶层

	c := cweskills.NewCWE(79, "XSS")
	c.Relationships = []cweskills.Relationship{
		{CWEID: 703, Nature: cweskills.RelationshipChildOf},
	}
	r.Register(c) // 非顶层

	r.BuildIndexes()
	top := cweskills.FindTopLevel(r)
	for _, c := range top {
		fmt.Println("顶层:", c.CWEID()) // CWE-703
	}
}
```

## 🆚 与 BuildForest 的关系

| 维度 | `FindTopLevel` | [`BuildForest`](./build-forest) |
| --- | --- | --- |
| 返回 | `[]*CWE`（仅根条目） | `[]*TreeNode`（含完整子树） |
| 后续 | 需自行 `BuildTree` 展开 | 直接可遍历 |

只需根列表用 `FindTopLevel`；需要整棵树用 `BuildForest`。

## ⚠️ 注意事项

::: details 顶层判定
「顶层」= `GetParentIDs(id)` 为空。一个弱点若只被别人 `ChildOf`（即自己是父）但自己不 `ChildOf` 任何人，仍算顶层。
:::

## 🔗 相关链接

- 构建森林：[BuildForest](./build-forest)
- 父级索引：[GetParentIDs](./relationship-indexes)
- Base 弱点查找：[链与复合/基础](./find-chains-composites)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
