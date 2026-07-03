---
title: Abstraction 抽象层级枚举
outline: [2, 3]
---

# 📚 Abstraction — 抽象层级枚举

`Abstraction` 表示 CWE 条目的抽象层级。CWE 按抽象程度从高到低分为四级：**Pillar > Class > Base > Variant**。层级越高描述越通用，越低越具体。其中 **Base 是映射到漏洞根因的首选级别**。

## 🧬 类型定义

```go
type Abstraction string
```

## 🏷️ 全部取值

| 值 | 常量名 | 含义 | 示例 CWE |
| --- | --- | --- | --- |
| `"Pillar"` | `AbstractionPillar` | 柱石，最高抽象，代表一个主题 | CWE-664 资源生命周期控制不当 |
| `"Class"` | `AbstractionClass` | 类，通常与语言/技术无关 | CWE-74 注入 |
| `"Base"` | `AbstractionBase` | 基础，足够具体以推断检测/预防方法 | CWE-79 XSS、CWE-89 SQL 注入 |
| `"Variant"` | `AbstractionVariant` | 变体，特定于某资源/技术/上下文 | CWE-83 网页属性脚本不当中和 |

```go
const (
	AbstractionPillar  Abstraction = "Pillar"
	AbstractionClass   Abstraction = "Class"
	AbstractionBase    Abstraction = "Base"
	AbstractionVariant Abstraction = "Variant"
)
```

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (a Abstraction) String() string` |
| `IsValid` | `func (a Abstraction) IsValid() bool` |
| `ParseAbstraction` | `func ParseAbstraction(s string) (Abstraction, error)` |
| `AllAbstractionValues` | `func AllAbstractionValues() []Abstraction` |

```go
a, err := cweskills.ParseAbstraction("Base")
fmt.Println(a, err)                              // Base <nil>
fmt.Println(a.String())                          // Base
fmt.Println(cweskills.Abstraction("Foo").IsValid()) // false
fmt.Println(cweskills.AllAbstractionValues())    // [Pillar Class Base Variant]
```

## 📊 额外方法：AbstractionOrder

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

```go
all := cweskills.AllAbstractionValues()
sort.Slice(all, func(i, j int) bool {
	return all[i].AbstractionOrder() < all[j].AbstractionOrder()
})
fmt.Println(all) // [Variant Base Class Pillar]
```

::: tip 用途
按 `AbstractionOrder` 升序排列得到「从具体到通用」的弱点层次，降序则反之。可用于把弱点归类到最具体层级或构建层次树。
:::

::: warning 权重不是严重程度
权重方向是 Pillar=4（最高）→ Variant=1（最低），即「层级越高数字越大」。这与「重要性」无关——Variant 同样是有效弱点，只是更具体。不要把权重当作严重程度。
:::

## 💻 CLI 对应命令

```bash
cwe enum abstraction
```

输出全部合法取值，详见 [CLI enum abstraction](../cli/enum-abstraction)。

## 🔗 相关链接

- SDK 视角：[Abstraction 抽象层级枚举](../sdk/enum-abstraction)
- 概念背景：[抽象层级 (Pillar/Class/Base/Variant)](../guide/concept-abstraction)
- 按层级筛选弱点：[`FindByAbstraction`](../sdk/find-by-abstraction)
- 按层级统计：[`CountByAbstraction`](../sdk/count-by-abstraction)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
