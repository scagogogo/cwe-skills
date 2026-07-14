---
title: 按后果范围查找
outline: [2, 3]
---

# 🔍 按后果范围查找

`FindByConsequenceScope` 按 [ConsequenceScope](./enum-consequence-scope) 枚举筛选弱点，返回具有指定后果范围（如 Confidentiality、Integrity、Availability）的弱点。用于按安全三要素（CIA）分类检索。

## 📐 函数签名

```go
func FindByConsequenceScope(r *Registry, scope ConsequenceScope) []*CWE
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `scope` | [`ConsequenceScope`](./enum-consequence-scope) | 后果范围枚举 |
| 返回 | `[]*CWE` | `CommonConsequences` 中含该范围的弱点 |

::: tip 数据来源
该函数遍历每个弱点的 `CommonConsequences` 字段，检查其中任一 `Consequence` 的 `Scope` 是否匹配。一个弱点的后果可能跨多个范围，因此会出现在多个不同 `scope` 的查询结果中。
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

	c := cweskills.NewCWE(79, "XSS")
	c.CommonConsequences = []cweskills.Consequence{
		{Scope: cweskills.ConsequenceScopeConfidentiality},
		{Scope: cweskills.ConsequenceScopeIntegrity},
	}
	r.Register(c)

	conf := cweskills.FindByConsequenceScope(r, cweskills.ConsequenceScopeConfidentiality)
	fmt.Println("Confidentiality:", len(conf)) // 1
	avail := cweskills.FindByConsequenceScope(r, cweskills.ConsequenceScopeAvailability)
	fmt.Println("Availability:", len(avail))   // 0
}
```

## ⚠️ 注意事项

::: warning 依赖后果字段
仅当弱点显式声明 `CommonConsequences` 且其中含匹配 `Scope` 时才返回。无后果字段的弱点不会出现。
:::

::: details 不依赖关系索引
该函数遍历 `GetAll()` 与各弱点的 `CommonConsequences`，无需 `BuildIndexes()`。
:::

## 🔗 相关链接

- 枚举定义：[ConsequenceScope](./enum-consequence-scope)
- 后果结构：[Consequence](./consequence)
- 范围统计：[CountByScope](./count-by-scope)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
