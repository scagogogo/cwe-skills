---
title: CWE 是什么
outline: [2, 3]
---

# 🆔 CWE 是什么

要会用 CWE Skills，先得理解它操作的「CWE」到底是什么。本页讲清 CWE 的来源、ID 格式、以及四类条目（弱点 / 类别 / 视图 / 复合元素）。

---

## 🌐 MITRE CWE 背景

**CWE**（Common Weakness Enumeration，通用缺陷枚举）是由 MITRE 组织维护的**软件与硬件安全弱点分类标准**。它给每一类「可能导致安全问题的条件或行为」分配一个唯一编号，并附上描述、后果、缓解措施、层级关系。

- 官网：<https://cwe.mitre.org/>
- 维护方：MITRE Corporation（同时维护 CVE、CAPEC 等）
- 用途：SAST/DAST 工具结果标注、漏洞管理平台分类、安全培训、合规对照（PCI DSS、OWASP ASVS 等）

::: tip CWE vs CVE
**CVE** 是「具体漏洞实例」（某软件某版本的某个具体 bug），**CWE** 是「这类漏洞的根因类别」。一条 CVE 通常会引用一个或多个 CWE——例如 CVE-2024-xxxx 引用 CWE-79（跨站脚本）。CWE Skills 操作的是 **CWE**，不是 CVE。
:::

---

## 🆔 CWE ID 格式

CWE ID 的标准格式是 `CWE-NNN`，其中 `NNN` 是正整数（无前导零要求，但允许出现）。

| 写法 | 是否合法 | 说明 |
|------|---------|------|
| `CWE-79` | ✅ | 标准格式 |
| `cwe-79` | ✅ | 大小写不敏感 |
| `CWE79` | ✅ | 无连字符 |
| `CWE 79` | ✅ | 空格分隔 |
| `79` | ✅ | 纯数字 |
| `079` | ✅ | 前导零 |
| `CWE-0` | ❌ | 数字须 > 0 |
| `abc` | ❌ | 非法 |

::: tip CWE Skills 全部兼容
上表所有合法写法，`cweskills.ParseCWEID` 都能正确解析为整数 `79`，`FormatCWEID` 统一规范化为 `CWE-79`。详见 [CWE ID 工具](../sdk/cwe-utils)。
:::

```go
num, _ := cweskills.ParseCWEID("cwe-079")      // 79
formatted, _ := cweskills.FormatCWEID("79")    // "CWE-79"
cweskills.FormatCWEIDFromInt(79)               // "CWE-79"
```

---

## 📦 四类条目

CWE 目录里不止有「弱点」，共分 **四类**条目。CWE Skills 用四个类型分别建模：`CWE`（弱点）、`Category`（类别）、`View`（视图）、`CompoundElement`（复合元素）。每个条目都有自己的 ID。

### 1. 弱点（Weakness）— `CWE` 结构体

CWE 的核心：描述一类软件/硬件安全弱点，是大多数场景里你打交道最多的条目。

```go
type CWE struct {
    ID                    int                    // 数字标识符
    Name                  string                 // 名称
    Abstraction           Abstraction           // 抽象层级 Pillar/Class/Base/Variant
    Structure             Structure             // 结构 Simple/Chain/Composite
    Status                Status                // 状态 Stable/Usable/Draft/...
    Description           string                // 描述
    ExtendedDescription   string                // 扩展描述
    LikelihoodOfExploit   LikelihoodOfExploit   // 被利用可能性
    CommonConsequences    []Consequence         // 常见后果
    PotentialMitigations  []Mitigation          // 缓解措施
    DemonstrativeExamples []DemonstrativeExample// 示范示例
    ObservedExamples      []ObservedExample     // 观察示例
    References            []Reference           // 参考文献
    Relationships         []Relationship        // 关系
    ApplicablePlatforms   *ApplicablePlatforms  // 适用平台
    ModesOfIntroduction   []Introduction        // 引入方式
    AlternateTerms        []AlternateTerm       // 备用术语
    Notes                 string                // 备注
    ContentHistory        *ContentHistory       // 内容历史
    CWEType               string                // 条目类型 weakness/...
    URL                   string                // URL
}
```

::: info 便利方法
`(*CWE).CWEID()` 返回标准格式字符串 `"CWE-79"`，等价于 `FormatCWEIDFromInt(c.ID)`。
:::

典型弱点示例：

| ID | 名称 | 抽象层级 | 说明 |
|----|------|---------|------|
| CWE-79 | 跨站脚本 (XSS) | Base | 注入类经典 |
| CWE-89 | SQL 注入 | Base | 注入类经典 |
| CWE-787 | 越界写 | Base | Top 25 常客 |
| CWE-664 | 不当的资源生命周期控制 | Pillar | 最高抽象 |
| CWE-74 | 注入 | Class | 类级别 |

### 2. 类别（Category）— `Category` 结构体

**类别**是一种**非层级的分组**：它把若干 CWE 按「共同主题」聚到一起，但成员之间不一定有父子关系。详见 [类别 (Category)](./concept-category)。

典型例子：`CWE-789`（内存分配不当）这类按主题归类的集合。一个 CWE 可以同时属于多个类别。

### 3. 视图（View）— `View` 结构体

**视图**是「从某个视角看 CWE 体系」的切片，分为 Graph（层次图）、Explicit Slice（显式切片）、Implicit Slice（隐式切片）三种。视图把 CWE 重新组织成适合特定受众（研究者、开发者、硬件设计者）的结构。详见 [视图 (View)](./concept-view)。

知名视图常量：

| 视图 ID | 名称 | 类型 |
|--------|------|------|
| CWE-1000 | 研究概念 | Graph |
| CWE-699 | 软件开发 | Graph |
| CWE-1199 | 硬件设计 | Graph |
| CWE-888 | CWE 横截面 | Explicit Slice |
| CWE-1400 | 综合 CWE 字典 | Explicit Slice |

```go
cweskills.CWEViewResearchConcepts        // 1000
cweskills.CWEViewDevelopmentConcepts     // 699
cweskills.CWEViewHardwareDesign          // 1199
cweskills.CWEViewCWECrossSection         // 888
cweskills.CWEViewComprehensiveDictionary // 1400
```

### 4. 复合元素（CompoundElement）— `CompoundElement` 结构体

**复合元素**描述由多个弱点**组合**而成的复合弱点，分链式（Chain）与复合（Composite）两种结构。详见 [复合元素 (CompoundElement)](./concept-compound)。

```go
type CompoundElement struct {
    ID            int
    Name          string
    Structure     Structure  // Chain 或 Composite
    Status        Status
    Description   string
    Relationships []Relationship
}
```

---

## 🧬 条目类型字段

每个条目在 `CWEType` 字段里标记自己是哪一类：`weakness`、`category`、`view`、`compound_element`。从 XML 解析或 API 返回时，CWE Skills 会据此把条目分发到 `Registry` 对应的集合（`weaknesses` / `categories` / `views` / `compoundElements`）。

---

## 🔗 关系：CWE 是一张图

四类条目之间通过**关系**连成一张有向图。CWE 定义了 10 种关系类型，分层级、顺序、依赖、对等四类。例如 `CWE-79` 通过 `ChildOf` 关系指向更通用的 `CWE-74`（注入）。

::: tip 关系是 CWE 的灵魂
单独看一个 CWE ID 信息有限；真正有价值的是它在关系图里的位置——祖先有多通用、后代有多具体、能否前置/跟随其他弱点、是否被复合弱点依赖。CWE Skills 的 [关系导航](../sdk/navigator) 正是为此而生。关系类型详解见 [关系类型](./concept-relationship)。
:::

---

## 📖 下一步

理解了 CWE 是什么，接下来深入各类概念：

- [抽象层级 (Pillar/Class/Base/Variant)](./concept-abstraction)
- [结构类型 (Simple/Chain/Composite)](./concept-structure)
- [关系类型 (ChildOf/Requires…)](./concept-relationship)
- [后果范围与影响](./concept-consequence)
- [视图 (View)](./concept-view)
- [类别 (Category)](./concept-category)
- [复合元素 (CompoundElement)](./concept-compound)
