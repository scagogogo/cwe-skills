---
title: model.go 数据模型概览
outline: [2, 3]
---

# 🧱 model.go — 数据模型概览

`model.go` 定义了 `cweskills` 包的全部核心数据结构。这些结构体是 SDK 与 MITRE CWE 数据之间的**契约层**：XML 解析器把官方目录反序列化成它们，API 客户端把 HTTP 响应映射到它们，序列化器再把它们导出为 JSON/CSV。理解这些结构，就理解了 SDK 的数据骨架。

## 🗺️ 模型地图

### 主要条目类型

CWE 体系里有四类「顶层条目」，每类对应一个结构体：

| 结构体 | 对应 `CWEType` | 说明 |
| --- | --- | --- |
| [`CWE`](./cwe-struct) | `weakness` | 弱点，最核心、字段最丰富 |
| [`Category`](./category) | `category` | 类别，把相关弱点分组 |
| [`View`](./view) | `view` | 视图，特定视角的组织 |
| [`CompoundElement`](./compound-element) | `compound_element` | 复合元素（链/复合弱点） |

::: tip 四者的统一标识
所有条目都用 `int` 类型 `ID` 作为主键，对外通过 `CWE-NNN` 形式表示。`CWE` 结构体额外有 `CWEType` 字段区分自己是哪一类。
:::

### CWE 内嵌的子模型

`CWE` 是个「聚合根」，挂载了大量子结构：

| 子模型 | 字段 | 文档 |
| --- | --- | --- |
| `Consequence` | `CommonConsequences` | [Consequence 后果](./consequence) |
| `Mitigation` | `PotentialMitigations` | [Mitigation 缓解措施](./mitigation) |
| `DemonstrativeExample` | `DemonstrativeExamples` | [DemonstrativeExample](./demonstrative-example) |
| `ObservedExample` | `ObservedExamples` | [ObservedExample](./observed-example) |
| `Reference` | `References` | [Reference 参考文献](./reference) |
| `ApplicablePlatforms` | `ApplicablePlatforms` | [ApplicablePlatforms](./applicable-platforms) |
| `Introduction` | `ModesOfIntroduction` | [Introduction 引入方式](./introduction) |
| `AlternateTerm` | `AlternateTerms` | [AlternateTerm 备用术语](./alternate-term) |
| `ContentHistory` | `ContentHistory` | [ContentHistory 内容历史](./content-history) |

## 🏷️ 序列化标签约定

每个字段同时携带 JSON 与 XML 标签，适配两种数据源：

```go
type CWE struct {
    ID          int    `json:"id" xml:"ID,attr"`
    Name        string `json:"name" xml:"Name"`
    Description string `json:"description" xml:"Description"`
    // ...
}
```

- **JSON tag**：下划线命名，如 `extended_description`、`common_consequences`
- **XML tag**：PascalCase，与 MITRE 官方 XML 一致，如 `Extended_Description`、`CommonConsequences>Consequence`
- 可选字段普遍带 `omitempty`，避免空值污染输出

::: details 嵌套集合的 XML tag 写法
`CommonConsequences>Consequence` 表示 XML 中 `<CommonConsequences><Consequence>...</Consequence></CommonConsequences>` 的嵌套关系，Go 的 `encoding/xml` 会自动展开。
:::

## 🛠️ 构造器一览

SDK 为四个顶层条目都提供了构造函数，方便快速创建实例：

| 构造器 | 签名 | 文档 |
| --- | --- | --- |
| `NewCWE` | `func NewCWE(id int, name string) *CWE` | [CWE 弱点](./cwe-struct) |
| `NewCategory` | `func NewCategory(id int, name string) *Category` | [Category 类别](./category) |
| `NewView` | `func NewView(id int, name string, viewType ViewType) *View` | [View 视图](./view) |
| `NewCompoundElement` | `func NewCompoundElement(id int, name string, structure Structure) *CompoundElement` | [CompoundElement](./compound-element) |

## 🔗 相关链接

- 枚举类型（结构体字段大量引用）：[enums.go 概览](./enums)
- 序列化：[序列化概览](./serializer)
- XML 解析：[xml_parser 概览](./xml-parser)
- 源文件：[`model.go`](https://github.com/scagogogo/cwe-skills/blob/main/model.go)、[`consequences.go`](https://github.com/scagogogo/cwe-skills/blob/main/consequences.go)
