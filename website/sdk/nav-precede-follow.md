---
title: 前置与后继导航
outline: [2, 3]
---

# 🧭 前置与后继导航

`CanPrecede` 与 `CanFollow` 描述 CWE 之间的**时序/因果链**：一个弱点「可以发生在前」或「可以发生在后」。这类关系用于刻画攻击链中弱点的先后顺序。

## 📐 方法签名

```go
func (n *Navigator) CanPrecede(id int) []*CWE
func (n *Navigator) CanFollow(id int) []*CWE
```

| 方法 | 含义 | 关系 Nature |
| --- | --- | --- |
| `CanPrecede` | 当前弱点之后可能出现的弱点（当前 → 后继） | `CanPrecede`（我能在前） |
| `CanFollow` | 当前弱点之前可能出现的弱点（当前 ← 前置） | `CanFollow`（我能在后） |

::: tip 语义方向
设 A 的关系里写了 `CanPrecede B`，意为「A 可在 B 之前」：
- `CanPrecede(A)` → 返回 B（A 之后可接谁）
- `CanFollow(B)` → 返回 A（B 之前可有谁）
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
	// 89 可在 79 之前（SQLi → XSS 的攻击链）
	c89 := cweskills.NewCWE(89, "SQL Injection")
	c89.Relationships = []cweskills.Relationship{
		{CWEID: 79, Nature: cweskills.RelationshipCanPrecede},
	}
	r.Register(c89)
	r.Register(cweskills.NewCWE(79, "XSS"))

	r.BuildIndexes()
	nav := cweskills.NewNavigator(r)

	for _, c := range nav.CanPrecede(89) {
		fmt.Println("89 之后可接:", c.CWEID()) // CWE-79
	}
	for _, c := range nav.CanFollow(79) {
		fmt.Println("79 之前可有:", c.CWEID()) // CWE-89
	}
}
```

## ⚠️ 注意事项

::: warning 不递归
`CanPrecede`/`CanFollow` 只返回**直接**声明的邻接节点，不沿攻击链递归展开。如需完整链条，结合 [ChainMembers](./nav-chain-composite) 或自行递归。
:::

::: details 与 Chain 的关系
`CanPrecede`/`CanFollow` 是**关系边**查询；而 [ChainMembers](./nav-chain-composite) 返回的是某条 Chain 类型弱点（`Structure=Chain`）的组成成员。两者数据来源不同。
:::

## 🔗 相关链接

- 链成员：[链与复合成员](./nav-chain-composite)
- 依赖关系（`Requires`）：[依赖关系](./nav-requires)
- 关系类型枚举：[RelationshipNature](./enum-relationship-nature)
- 源文件：[`navigator.go`](https://github.com/scagogogo/cwe-skills/blob/main/navigator.go)
