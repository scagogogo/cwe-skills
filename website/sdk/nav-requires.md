---
title: 依赖关系导航
outline: [2, 3]
---

# 🧭 依赖关系导航

`Requires` 与 `RequiredBy` 描述弱点之间的**前置依赖**：一个弱点的存在或成立，依赖于另一个弱点先行出现。二者互为反向查询。

## 📐 方法签名

```go
func (n *Navigator) Requires(id int) []*CWE
func (n *Navigator) RequiredBy(id int) []*CWE
```

| 方法 | 含义 | 关系 Nature |
| --- | --- | --- |
| `Requires` | 当前弱点**依赖**的弱点（前置条件） | `Requires`（我依赖对方） |
| `RequiredBy` | **依赖当前弱点**的弱点（被谁依赖） | `RequiredBy`（对方依赖我） |

::: tip 与 CanPrecede 的区别
- `CanPrecede`：弱时序/因果，「可」在前，非强制。
- `Requires`：强依赖，「必须」在前，缺一不可。

语义上 `Requires` 比 `CanPrecede` 更严格。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	// 89 依赖 20（输入校验缺失 → SQLi）
	c89 := cweskills.NewCWE(89, "SQL Injection")
	c89.Relationships = []cweskills.Relationship{
		{CWEID: 20, Nature: cweskills.RelationshipRequires},
	}
	r.Register(c89)
	r.Register(cweskills.NewCWE(20, "Improper Input Validation"))

	r.BuildIndexes()
	nav := cweskills.NewNavigator(r)

	for _, c := range nav.Requires(89) {
		fmt.Println("89 依赖:", c.CWEID()) // CWE-20
	}
	for _, c := range nav.RequiredBy(20) {
		fmt.Println("20 被依赖:", c.CWEID()) // CWE-89
	}
}
```

## ⚠️ 边界行为

::: warning 方向不要搞反
`Requires(A)` 返回的是 A 所依赖的节点（A 的前提），不是依赖 A 的节点。若想要「谁依赖我」，用 `RequiredBy(A)`。
:::

::: details 不递归
`Requires`/`RequiredBy` 只返回直接声明的依赖边，不递归展开传递依赖。如需完整依赖闭包，可基于结果自行迭代调用。
:::

## 🔗 相关链接

- 时序关系：[前置与后继](./nav-precede-follow)
- 替代关系：[可作为](./nav-can-also-be)
- 关系类型枚举：[RelationshipNature](./enum-relationship-nature)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
