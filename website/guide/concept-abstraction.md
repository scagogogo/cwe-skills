---
title: 抽象层级 (Pillar/Class/Base/Variant)
outline: [2, 3]
---

# 🧬 抽象层级 (Pillar / Class / Base / Variant)

CWE 按**抽象程度**从高到低分为四个层级：**Pillar（柱石）> Class（类）> Base（基础）> Variant（变体）**。抽象层级越高，描述越通用、越接近「根因大类」；越低，越具体、越接近「特定技术下的某种写法」。

::: tip Base 是首选映射层
**Base 级别是映射到漏洞根因的首选级别**——它「足够具体以推断检测/预防方法」，又不会像 Variant 那样绑死在某一具体技术上。SAST 工具结果通常映射到 Base。
:::

---

## 📊 四个层级

`Abstraction` 是一个枚举类型，定义在 `enums.go`：

```go
type Abstraction string

const (
    AbstractionPillar  Abstraction = "Pillar"
    AbstractionClass   Abstraction = "Class"
    AbstractionBase    Abstraction = "Base"
    AbstractionVariant Abstraction = "Variant"
)
```

| 层级 | 含义 | 通用度 | 典型示例 |
|------|------|--------|---------|
| **Pillar** | 柱石，最高抽象，代表一个主题 | 最通用 | CWE-664 不当的资源生命周期控制 |
| **Class** | 类，与特定语言/技术无关 | 通用 | CWE-74 注入（输出中特殊元素不当中和） |
| **Base** | 基础，足够具体以推断检测/预防方法 | 中等 | CWE-79 XSS、CWE-89 SQL 注入 |
| **Variant** | 变体，特定于某资源/技术/上下文 | 最具体 | CWE-83 网页属性中脚本的不当中和 |

抽象程度从上到下递减：

```text
Pillar  （柱石：根因大类）
  ▲ 通用
  │
Class   （类：技术无关）
  │
Base    （基础：可检测可预防）  ← 首选映射层
  │
  ▼ 具体
Variant （变体：特定技术）
```

---

## 🔢 排序权重

`AbstractionOrder()` 返回排序权重，层级越高值越大，可用于按抽象程度排序：

| 层级 | 权重 |
|------|------|
| Pillar | 4 |
| Class | 3 |
| Base | 2 |
| Variant | 1 |
| 未知 | 0 |

```go
cweskills.AbstractionPillar.AbstractionOrder()  // 4
cweskills.AbstractionBase.AbstractionOrder()    // 2
```

---

## 🛠️ 枚举 API

`Abstraction` 提供一组通用方法，所有 CWE 枚举类型都遵循同一套 API 风格：

```go
// 校验
cweskills.AbstractionBase.IsValid()              // true
cweskills.Abstraction("Foo").IsValid()          // false

// 解析（解析失败返回 ValidationError）
a, err := cweskills.ParseAbstraction("Base")    // AbstractionBase, nil
a, err := cweskills.ParseAbstraction("Foo")     // "", *ValidationError

// 穷举所有合法值
cweskills.AllAbstractionValues()
// [Pillar, Class, Base, Variant]

// 字符串表示
string(cweskills.AbstractionBase)               // "Base"
```

::: tip 用枚举而非裸字符串
代码里始终用 `cweskills.AbstractionBase` 而非 `"Base"`，编译期就能挡住拼写错误，还能用 `IsValid()` 校验外部输入。
:::

---

## 🎯 用途：为什么抽象层级重要

### 1. 优先级排序

Pillar/Class 太抽象，难以直接指导检测；Variant 太具体，覆盖面窄。**Base 级别是平衡点**，常作为安全工具的优先处理对象。

```go
// 过滤出所有 Base 级别弱点
baseWeaknesses := cweskills.FindByAbstraction(registry, cweskills.AbstractionBase)
```

### 2. 上卷/下钻分析

- **上卷**：从 Variant/Base 找到 Class/Pillar 祖先，理解根因大类。
- **下钻**：从 Pillar 下钻到 Base/Variant，找到可操作的具体弱点。

```go
nav := cweskills.NewNavigator(registry)
ancestors := nav.Ancestors(79)  // 从 Base 向上找 Class、Pillar
```

### 3. 按抽象分组统计

```go
groups := cweskills.GroupByAbstraction(registry.GetAll())
for abs, list := range groups {
    fmt.Printf("%s: %d 条\n", abs, len(list))
}
```

### 4. 过滤与排序

```go
filtered := cweskills.Filter(results, cweskills.FilterOption{
    Abstraction: cweskills.AbstractionBase,
})
sorted := cweskills.SortByAbstraction(filtered) // 按 Pillar→Variant 排序
```

::: info SortByAbstraction
`SortByAbstraction` 内部即依据 `AbstractionOrder()`，把更通用的层级排在前面。
:::

---

## 🌳 一个完整层级链示例

以跨站脚本（XSS）为例，从 Pillar 到 Variant：

```text
CWE-664  Pillar   不当的资源生命周期控制
  └─ CWE-74  Class  注入
       └─ CWE-79  Base     跨站脚本 (XSS)        ← 首选映射
            └─ CWE-83  Variant  网页属性中脚本的不当中和
```

- 报告里写 `CWE-79` 就够用（Base，可检测可预防）；
- 要做根因分析，沿祖先链上溯到 `CWE-74`（注入大类）甚至 `CWE-664`；
- 要精确描述「在 HTML 属性里」的具体变体，用 `CWE-83`。

---

## ⚠️ 注意事项

::: warning 抽象层级不等于严重程度
Pillar 不代表「最危险」，Variant 也不代表「最轻微」。抽象层级描述的是**通用度**，与漏洞的**危害程度**无关。危害程度看 [后果范围与影响](./concept-consequence) 和 `LikelihoodOfExploit`。
:::

::: info 并非所有 CWE 都有完整四层
并非每条 CWE 都能上溯到 Pillar、下钻到 Variant。有些分支只到 Class 或 Base 就停了。导航时拿到空切片是正常的。
:::

---

## 📖 相关文档

- [CWE 是什么](./concept-cwe)
- [结构类型 (Simple/Chain/Composite)](./concept-structure)
- [关系类型](./concept-relationship)
- [关系导航 API](../sdk/navigator)
- [搜索过滤 API](../sdk/search)
- [Abstraction 枚举参考](../enums/abstraction)
