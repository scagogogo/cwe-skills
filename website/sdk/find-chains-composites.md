---
title: 链、复合与基础查找
outline: [2, 3]
---

# 🔍 链、复合与基础弱点查找

`FindBaseWeaknesses`、`FindChains`、`FindComposites` 三个便捷函数分别返回 Base 抽象、Chain 结构、Composite 结构的弱点。它们是 `FindByAbstraction`/`FindByStructure` 的常用快捷方式。

## 📐 函数签名

```go
func FindBaseWeaknesses(r *Registry) []*CWE
func FindChains(r *Registry) []*CWE
func FindComposites(r *Registry) []*CWE
```

| 函数 | 等价于 | 含义 |
| --- | --- | --- |
| `FindBaseWeaknesses` | `FindByAbstraction(r, AbstractionBase)` | Base 抽象层级的弱点 |
| `FindChains` | `FindByStructure(r, StructureChain)` | 链式结构的弱点 |
| `FindComposites` | `FindByStructure(r, StructureComposite)` | 复合结构的弱点 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	base := cweskills.NewCWE(79, "XSS")
	base.Abstraction = cweskills.AbstractionBase
	r.Register(base)

	chain := cweskills.NewCWE(680, "Chain")
	chain.Structure = cweskills.StructureChain
	r.Register(chain)

	composite := cweskills.NewCWE(1000, "Composite")
	composite.Structure = cweskills.StructureComposite
	r.Register(composite)

	fmt.Println("Base:", len(cweskills.FindBaseWeaknesses(r))) // 1
	fmt.Println("Chains:", len(cweskills.FindChains(r)))       // 1
	fmt.Println("Composites:", len(cweskills.FindComposites(r))) // 1
}
```

## ⚠️ 注意事项

::: tip 为什么单独提供？
Base/Chain/Composite 是 CWE 分析中最常关注的三类：Base 是可直接检测的具体弱点，Chain/Composite 揭示攻击组合。提供专用函数让调用方语义更清晰。
:::

::: warning 不依赖索引
三个函数都遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 通用版本：[按抽象查找](./find-by-abstraction)、[按结构查找](./find-by-structure)
- 枚举：[Abstraction](./enum-abstraction)、[Structure](./enum-structure)
- 链/复合成员导航：[链与复合成员](./nav-chain-composite)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
