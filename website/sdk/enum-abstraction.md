---
title: Abstraction 抽象层级枚举
outline: [2, 3]
---

# 📚 Abstraction — 抽象层级枚举

`Abstraction` 表示 CWE 条目的抽象层级。CWE 按抽象程度从高到低分为四级：**Pillar > Class > Base > Variant**。层级越高描述越通用，越低越具体。其中 **Base 是映射到漏洞根因的首选级别**。

## 📋 类型与常量

```go
type Abstraction string

const (
	AbstractionPillar  Abstraction = "Pillar"
	AbstractionClass   Abstraction = "Class"
	AbstractionBase    Abstraction = "Base"
	AbstractionVariant Abstraction = "Variant"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 | 示例 |
| --- | --- | --- | --- |
| `AbstractionPillar` | `"Pillar"` | 柱石，最高抽象，代表一个主题 | CWE-664 资源生命周期控制不当 |
| `AbstractionClass` | `"Class"` | 类，与语言/技术无关 | CWE-74 注入 |
| `AbstractionBase` | `"Base"` | 基础，足够具体以推断检测/预防方法 | CWE-79 XSS、CWE-89 SQL 注入 |
| `AbstractionVariant` | `"Variant"` | 变体，特定于某资源/技术/上下文 | CWE-83 网页属性脚本不当中和 |

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (a Abstraction) String() string` |
| `IsValid` | `func (a Abstraction) IsValid() bool` |
| `ParseAbstraction` | `func ParseAbstraction(s string) (Abstraction, error)` |
| `AllAbstractionValues` | `func AllAbstractionValues() []Abstraction` |

## 📊 排序权重

```go
func (a Abstraction) AbstractionOrder() int
```

返回排序权重，**层级越高值越大**：

| 取值 | 权重 |
| --- | --- |
| `AbstractionPillar` | 4 |
| `AbstractionClass` | 3 |
| `AbstractionBase` | 2 |
| `AbstractionVariant` | 1 |
| 未知 | 0 |

::: tip 用途
按 `AbstractionOrder` 降序排列可得到「从通用到具体」的弱点层次；升序则反之。也可用于把弱点归类到最具体的层级。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"sort"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	a, err := cweskills.ParseAbstraction("Base")
	fmt.Println(a, err) // Base <nil>

	// 校验
	fmt.Println(cweskills.AbstractionPillar.IsValid())  // true
	fmt.Println(cweskills.Abstraction("Foo").IsValid()) // false

	// 排序权重
	fmt.Println(cweskills.AbstractionPillar.AbstractionOrder()) // 4
	fmt.Println(cweskills.AbstractionVariant.AbstractionOrder()) // 1

	// 全部值
	all := cweskills.AllAbstractionValues()
	sort.Slice(all, func(i, j int) bool {
		return all[i].AbstractionOrder() < all[j].AbstractionOrder()
	})
	fmt.Println(all) // [Variant Base Class Pillar]
}
```

## 🎯 典型用途

<Badge type="tip" text="映射根因" /> 优先用 Base 级弱点做漏洞归因
<Badge type="info" text="过滤" /> [`FindByAbstraction`](./find-by-abstraction) 按层级筛弱点
<Badge type="warning" text="层次展示" /> 用 `AbstractionOrder` 排序构建层次树

## ⚠️ 注意事项

::: warning AbstractionOrder 方向
权重是 **Pillar=4（最高）→ Variant=1（最低）**，即「层级越高数字越大」。这与「重要性」无关——Variant 同样是有效弱点，只是更具体。不要把权重当作严重程度。
:::

## 🔗 相关链接

- 字段归宿：`CWE.Abstraction`，见 [CWE 弱点](./cwe-struct)
- 便捷判断方法：`CWE.IsPillar()`/`IsBase()`/`IsVariant()`，见 [CWE 类型判断方法](./cwe-type-methods)
- 概念背景：[抽象层级 (Pillar/Class/Base/Variant)](../guide/concept-abstraction)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
