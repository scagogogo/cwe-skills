---
title: 按结构查找
outline: [2, 3]
---

# 🔍 按结构查找

`FindByStructure` 按 [Structure](./enum-structure) 枚举筛选弱点，返回指定结构（Simple、Chain、Composite）的弱点。用于区分简单弱点与聚合型弱点。

## 📐 函数签名

```go
func FindByStructure(r *Registry, structure Structure) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `structure` | [`Structure`](./enum-structure) | 结构枚举 |
| 返回 | `[]*CWE` | 该结构的全部弱点 |

## ✅ 示例

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()

	simple := cweskills.NewCWE(79, "XSS")
	simple.Structure = cweskills.StructureSimple
	r.Register(simple)

	chain := cweskills.NewCWE(680, "Chain Example")
	chain.Structure = cweskills.StructureChain
	r.Register(chain)

	composite := cweskills.NewCWE(1000, "Composite Example")
	composite.Structure = cweskills.StructureComposite
	r.Register(composite)

	fmt.Println("Simple:", len(cweskills.FindByStructure(r, cweskills.StructureSimple)))         // 1
	fmt.Println("Chain:", len(cweskills.FindByStructure(r, cweskills.StructureChain)))           // 1
	fmt.Println("Composite:", len(cweskills.FindByStructure(r, cweskills.StructureComposite)))   // 1
}
```

## 🆚 与 FindChains/FindComposites 的关系

| 维度 | `FindByStructure` | [`FindChains`](./find-chains-composites) / `FindComposites` |
| --- | --- | --- |
| 通用性 | 通用，按任意 `Structure` 查 | 专用，固定查 Chain/Composite |
| 等价性 | `FindByStructure(r, StructureChain)` ≡ `FindChains(r)` | 同理 |

二者结果等价，`FindChains`/`FindComposites` 是常用场景的便捷别名。

## ⚠️ 注意事项

::: tip 不依赖索引
`FindByStructure` 遍历 `GetAll()`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[Structure](./enum-structure)
- 便捷函数：[链与复合查找](./find-chains-composites)
- 结构统计：[ComputeStatistics](./compute-statistics)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
