# CWE SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/cwe.svg)](https://pkg.go.dev/github.com/scagogogo/cwe)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/cwe)](https://goreportcard.com/report/github.com/scagogogo/cwe)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个完善的 [CWE（通用缺陷枚举）](https://cwe.mitre.org/) Go SDK，为构建网络安全产品提供完整的API支持。

## 功能特性

- **完整的CWE数据模型**：弱点、类别、视图、复合元素，涵盖所有字段
- **类型化枚举**：抽象层级、状态值、关系类型、后果范围等
- **CWE ID工具**：解析、格式化、验证和从文本中提取CWE ID
- **知名列表**：CWE Top 25、OWASP Top 10、SANS Top 25 及成员检查
- **MITRE REST API客户端**：完整访问CWE API，支持速率限制和重试
- **XML目录解析器**：离线解析MITRE官方XML下载文件
- **内存注册表**：存储、索引和查询CWE条目及关系索引
- **搜索与过滤**：按关键字、抽象层级、状态、可能性、后果范围等
- **关系导航**：父级、子级、祖先、后代、对等、链式、组合
- **树构建**：从CWE关系构建层次树
- **序列化**：JSON、XML和CSV的导入/导出
- **零依赖**：仅使用Go标准库

## 安装

```bash
go get github.com/scagogogo/cwe
```

## 快速开始

### 解析和验证CWE ID

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe"
)

func main() {
    // 解析CWE ID
    id, _ := cwe.ParseCWEID("CWE-79")
    fmt.Println(id) // 79

    // 格式化CWE ID
    formatted, _ := cwe.FormatCWEID("79")
    fmt.Println(formatted) // CWE-79

    // 验证
    if cwe.IsCWEID("CWE-89") {
        fmt.Println("有效的CWE ID")
    }

    // 从文本中提取
    ids := cwe.ExtractCWEIDs("参见CWE-79和CWE-89了解详情")
    fmt.Println(ids) // [CWE-79 CWE-89]
}
```

### 查询MITRE CWE API

```go
client := cwe.NewAPIClient()

// 获取弱点
weakness, err := client.GetWeakness(ctx, 79)

// 获取版本
version, err := client.GetVersion(ctx)

// 获取关系
parents, err := client.GetParents(ctx, 79)
children, err := client.GetChildren(ctx, 79, 1000) // 指定视图ID
```

### 使用注册表进行本地操作

```go
registry := cwe.NewRegistry()

// 注册CWE条目
registry.Register(&cwe.CWE{
    ID:          79,
    Name:        "跨站脚本攻击(XSS)",
    Abstraction: cwe.AbstractionBase,
    Status:      cwe.StatusStable,
})

// 构建关系索引
registry.BuildIndexes()

// 搜索和过滤
results := cwe.FindByAbstraction(registry, cwe.AbstractionBase)
filtered := cwe.Filter(results, cwe.FilterOption{Status: cwe.StatusStable})

// 导航关系
nav := cwe.NewNavigator(registry)
parents := nav.Parents(79)
ancestors := nav.Ancestors(79)
```

### 解析离线XML目录

```go
parser := cwe.NewXMLParser()
registry, err := parser.ParseFile("cwec_v4.10.xml")
```

### 检查知名列表

```go
if cwe.IsInTop25(79) {
    fmt.Println("CWE-79在Top 25中！")
}

category := cwe.GetOWASPCategory(79)
fmt.Println(category) // A03:2021-Injection
```

## API参考

### 核心类型

| 类型 | 描述 |
|------|------|
| `CWE` | 核心弱点条目，包含所有CWE字段 |
| `Category` | CWE类别，包含成员 |
| `View` | CWE视图，包含成员和类型 |
| `CompoundElement` | 链式或复合弱点 |
| `Relationship` | CWE条目之间的关系 |
| `Consequence` | 影响后果，包含范围和严重程度 |

### 枚举类型

| 类型 | 值 |
|------|-----|
| `Abstraction` | Pillar, Class, Base, Variant |
| `Structure` | Simple, Chain, Composite |
| `Status` | Stable, Usable, Draft, Incomplete, Obsolete, Deprecated |
| `LikelihoodOfExploit` | High, Medium, Low, Unknown |
| `RelationshipNature` | ChildOf, ParentOf, CanPrecede, CanFollow, Requires, RequiredBy, CanAlsoBe, PeerOf, MemberOf, HasMember |
| `ConsequenceScope` | Confidentiality, Integrity, Availability, Access Control 等 |
| `ViewType` | Graph, Explicit Slice, Implicit Slice |

### 关键函数

| 函数 | 描述 |
|------|------|
| `ParseCWEID(s)` | 从字符串解析CWE ID |
| `FormatCWEID(s)` | 格式化为 "CWE-NNN" |
| `IsCWEID(s)` | 检查是否为有效CWE ID |
| `IsInTop25(id)` | 检查是否在CWE Top 25中 |
| `IsInOWASPTop10(id)` | 检查是否在OWASP Top 10中 |
| `FindByKeyword(r, kw)` | 按关键字搜索 |
| `Filter(cwes, opts)` | 按多条件过滤 |
| `BuildTree(r, id)` | 构建层次树 |
| `MarshalJSON/UnmarshalJSON` | JSON序列化 |
| `MarshalXML/UnmarshalXML` | XML序列化 |
| `MarshalCSV/UnmarshalCSV` | CSV序列化 |

## 许可证

MIT许可证 - 详见 [LICENSE](LICENSE)
