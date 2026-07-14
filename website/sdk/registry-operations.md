---
title: 注册表基础操作
outline: [2, 3]
---

# 🗃️ 注册表基础操作

`Registry` 提供四类条目（`CWE` / `Category` / `View` / `CompoundElement`）的注册、查询、计数与删除能力。本文档汇总这些**非索引、非序列化**的基础方法。

## 📥 注册方法

| 方法 | 签名 | 说明 |
| --- | --- | --- |
| `Register` | `func (r *Registry) Register(c *CWE) error` | 注册弱点，**ID 重复时报错** |
| `RegisterCategory` | `func (r *Registry) RegisterCategory(c *Category)` | 注册分类 |
| `RegisterView` | `func (r *Registry) RegisterView(v *View)` | 注册视图 |
| `RegisterCompoundElement` | `func (r *Registry) RegisterCompoundElement(c *CompoundElement)` | 注册复合元素 |

::: warning Register 会报错
`Register` 与其它三个不同：当传入的 `CWE.ID` 已存在时返回 `error`，调用方应处理。其余三个注册方法遇重复 ID 时**覆盖**旧值且不报错。
:::

## 🔍 查询方法

| 方法 | 签名 | 返回 |
| --- | --- | --- |
| `Get` | `func (r *Registry) Get(id int) (*CWE, bool)` | 弱点指针 + 是否存在 |
| `GetCategory` | `func (r *Registry) GetCategory(id int) (*Category, bool)` | 分类 |
| `GetView` | `func (r *Registry) GetView(id int) (*View, bool)` | 视图 |
| `GetCompoundElement` | `func (r *Registry) GetCompoundElement(id int) (*CompoundElement, bool)` | 复合元素 |
| `GetAll` | `func (r *Registry) GetAll() []*CWE` | 全部弱点（无序） |
| `GetAllCategories` | `func (r *Registry) GetAllCategories() []*Category` | 全部分类 |
| `GetAllViews` | `func (r *Registry) GetAllViews() []*View` | 全部视图 |
| `Contains` | `func (r *Registry) Contains(id int) bool` | 弱点 ID 是否存在 |

## 🔢 计数方法

| 方法 | 签名 | 返回 |
| --- | --- | --- |
| `Size` | `func (r *Registry) Size() int` | 弱点数量 |
| `CategoryCount` | `func (r *Registry) CategoryCount() int` | 分类数量 |
| `ViewCount` | `func (r *Registry) ViewCount() int` | 视图数量 |
| `CompoundElementCount` | `func (r *Registry) CompoundElementCount() int` | 复合元素数量 |

## 🗑️ 删除方法

| 方法 | 签名 | 说明 |
| --- | --- | --- |
| `Remove` | `func (r *Registry) Remove(id int) error` | 删除弱点，不存在时报错 |
| `RemoveCategory` | `func (r *Registry) RemoveCategory(id int) error` | 删除分类 |
| `RemoveView` | `func (r *Registry) RemoveView(id int) error` | 删除视图 |
| `Clear` | `func (r *Registry) Clear()` | 清空全部条目 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	// 注册
	_ = r.Register(cweskills.NewCWE(79, "XSS"))
	_ = r.Register(cweskills.NewCWE(89, "SQL Injection"))

	// 重复注册报错
	if err := r.Register(cweskills.NewCWE(79, "dup")); err != nil {
		fmt.Println("重复:", err) // 重复: ...
	}

	// 查询
	fmt.Println(r.Size(), r.Contains(79))     // 2 true
	if c, ok := r.Get(89); ok {
		fmt.Println(c.Name)                   // SQL Injection
	}

	// 删除
	_ = r.Remove(79)
	fmt.Println(r.Contains(79))               // false
}
```

::: details GetAll 的顺序
`GetAll` 返回的是从内部 map 转换出的切片，**顺序不确定**。如需稳定顺序，配合 [SortByID](./sort) 排序后再使用。
:::

## 🔗 相关链接

- 索引构建：[BuildIndexes](./build-indexes)
- 整库序列化：[ExportJSON / ImportJSON](./registry-json)
- 分类/视图结构：[Category](./category)、[View](./view)
- 源文件：[`registry.go`](https://github.com/scagogogo/cwe-skills/blob/main/registry.go)
