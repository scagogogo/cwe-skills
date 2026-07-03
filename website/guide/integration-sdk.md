---
title: Go SDK 接入
outline: [2, 3]
---

# 🔧 Go SDK 接入

把 `cweskills` 包 import 进你的 Go 应用，以原生 Go 类型调用全部能力。这是**集成深度最高**的方式——类型化、编译期检查、结构化错误、零依赖核心。

---

## 📦 安装

```bash
go get github.com/scagogogo/cwe-skills
```

::: warning 包名是 cweskills
模块路径是 `github.com/scagogogo/cwe-skills`，但 **Go 包名是 `cweskills`**（不是 `cwe-skills`，也不是 `cwe`）。导入后用 `cweskills.` 前缀调用。
:::

要求 Go 1.25+（见 `go.mod`）。核心 SDK 仅依赖 Go 标准库，无任何第三方包。

---

## 🚀 快速上手

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    fmt.Println(cweskills.Version) // v0.0.1

    // ID 工具
    id, _ := cweskills.ParseCWEID("cwe-079")
    formatted, _ := cweskills.FormatCWEID("79")
    fmt.Println(id, formatted) // 79 CWE-79

    // 知名列表
    fmt.Println(cweskills.IsInTop25(79))       // true
    fmt.Println(cweskills.GetOWASPCategory(89)) // A03:2021-Injection

    // 枚举
    for _, a := range cweskills.AllAbstractionValues() {
        fmt.Println(a) // Pillar Class Base Variant
    }
}
```

---

## 🆔 CWE ID 工具

```go
num, _      := cweskills.ParseCWEID("CWE-79")        // 79
formatted, _ := cweskills.FormatCWEID("79")          // "CWE-79"
cweskills.FormatCWEIDFromInt(79)                     // "CWE-79"
ok          := cweskills.IsCWEID("CWE79")            // true
err         := cweskills.ValidateCWEID("CWE-79")     // nil
ids         := cweskills.ExtractCWEIDs("见 CWE-79 与 cwe89") // ["CWE-79","CWE-89"]
first       := cweskills.ExtractFirstCWEID("...CWE-79...")   // "CWE-79"
cmp, _      := cweskills.CompareCWEIDs("CWE-79", "CWE-89")   // -1
```

详见 [CWE ID 工具](../sdk/cwe-utils)。

---

## 🏆 知名列表

```go
cweskills.IsInTop25(79)         // true
cweskills.IsInOWASPTop10(89)    // true
cweskills.IsInSANSTop25(79)     // true
cweskills.GetOWASPCategory(89)  // "A03:2021-Injection"
cweskills.GetOWASPCategories(89) // []string

// 知名视图常量
cweskills.CWEViewResearchConcepts        // 1000
cweskills.CWEViewDevelopmentConcepts     // 699
cweskills.CWEViewHardwareDesign          // 1199
cweskills.CWEViewCWECrossSection         // 888
cweskills.CWEViewComprehensiveDictionary // 1400
```

详见 [知名列表](../wellknown/cwe-top-25)。

---

## 🌐 在线：MITRE REST API 客户端

```go
ctx := context.Background()

client := cweskills.NewAPIClient()
defer client.Close()

// 默认 baseURL = https://cwe-api.mitre.org/api/v1
// 默认速率限制 rate=0.1, burst=1（约每 10 秒 1 个请求）

weakness, _ := client.GetWeakness(ctx, 79)
category, _  := client.GetCategory(ctx, 789)
view, _      := client.GetView(ctx, 1000)
cwes, _      := client.GetCWEs(ctx, []int{79, 89, 352}) // 批量

parents, _   := client.GetParents(ctx, 79)
children, _  := client.GetChildren(ctx, 79)
ancestors, _ := client.GetAncestors(ctx, 79)
descendants, _ := client.GetDescendants(ctx, 79)

version, _   := client.GetVersion(ctx) // MITRE API 版本
```

自定义配置：

```go
client := cweskills.NewAPIClient(
    cweskills.WithAPIBaseURL("https://cwe-api.mitre.org/api/v1"),
    cweskills.WithAPITimeout(60*time.Second),
    cweskills.WithAPIRateLimit(0.5, 1),    // 每秒 0.5 个请求
    cweskills.WithAPIRetry(3, 2*time.Second), // 最多重试 3 次
)
```

详见 [MITRE REST API 客户端](../sdk/api-client) 与 [速率限制与重试](./rate-limit-retry)。

---

## 📥 离线：XML 解析 + 注册表

```go
// 解析 XML，构建内存注册表
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes() // 建立父子/对等/成员等多层索引

// 查询
cwe79, ok := registry.Get(79)
all       := registry.GetAll()
size      := registry.Size()
contains  := registry.Contains(79)

cat, ok := registry.GetCategory(789)
view, ok := registry.GetView(1000)
ce, ok  := registry.GetCompoundElement(680)

// 成员查询
viewMembers  := registry.GetViewMembers(1000)
catMembers   := registry.GetCategoryMembers(789)
memberOf     := registry.GetMemberOfIDs(79)
```

详见 [离线 XML 解析](../sdk/xml-parser) 与 [注册表](../sdk/registry)。

---

## 🧭 关系导航（离线）

```go
nav := cweskills.NewNavigator(registry)

// 层级
nav.Parents(79); nav.Children(79)
nav.Ancestors(79); nav.Descendants(79)
nav.Siblings(79)

// 对等
nav.Peers(79); nav.CanAlsoBe(79)

// 顺序（链式）
nav.CanPrecede(680); nav.CanFollow(680)
nav.ChainMembers(680)

// 依赖（复合）
nav.Requires(352); nav.RequiredBy(79)
nav.CompositeMembers(352)

// 图算法
path  := nav.ShortestPath(79, 1)        // []int
depth := nav.RelationshipDepth(1, 79)   // int
nav.IsAncestorOf(1, 79)                 // bool
nav.IsRelated(79, 1)                    // bool
```

::: warning 导航仅离线可用
`Navigator` 依赖注册表里的完整关系，**只在离线路径下有数据**。在线 API 不返回链式/依赖/对等关系。
:::

详见 [关系导航](../sdk/navigator)。

---

## 🌳 层次树构建（离线）

```go
tree    := cweskills.BuildTree(registry, 1)                              // 以 CWE-1 为根
forest  := cweskills.BuildForest(registry)                               // 全森林
viewTree := cweskills.BuildViewTree(registry, cweskills.CWEViewResearchConcepts)

// 遍历
tree.Walk(func(n *cweskills.TreeNode) bool { /* DFS */ return true })
tree.WalkBFS(func(n *cweskills.TreeNode) bool { /* BFS */ return true })

// 查询
node    := tree.Find(79)
path    := node.Path()        // 从根到该节点
leaves  := tree.LeafNodes()
maxDepth := tree.MaxDepth()
count   := tree.Count()
```

详见 [层次树构建](../sdk/tree)。

---

## 🔍 搜索与过滤（离线）

```go
// 单维度查找
results := cweskills.FindByKeyword(registry, "Injection")
results = cweskills.FindByAbstraction(registry, cweskills.AbstractionBase)
results = cweskills.FindByStatus(registry, cweskills.StatusStable)
results = cweskills.FindByLikelihood(registry, cweskills.LikelihoodHigh)
results = cweskills.FindByConsequenceScope(registry, cweskills.ScopeConfidentiality)
results = cweskills.FindByStructure(registry, cweskills.StructureChain)

// 多条件过滤
filtered := cweskills.Filter(results, cweskills.FilterOption{
    Keyword:            "injection",
    Abstraction:        cweskills.AbstractionBase,
    Status:             cweskills.StatusStable,
    LikelihoodOfExploit: cweskills.LikelihoodHigh,
    ConsequenceScope:   cweskills.ScopeConfidentiality,
    Structure:          cweskills.StructureSimple,
})

// 排序 / 分组 / 去重
sorted := cweskills.SortByID(filtered)
groups := cweskills.GroupByAbstraction(filtered)
unique := cweskills.Deduplicate(filtered)
```

详见 [搜索与过滤](../sdk/search)。

---

## 📦 序列化

```go
// 单条
jsonBytes, _ := cweskills.MarshalJSON(cwe79)
cweBack, _   := cweskills.UnmarshalJSON(jsonBytes)
xmlBytes, _  := cweskills.MarshalXML(cwe79)
csvRow, _    := cweskills.MarshalCSV([]*cweskills.CWE{cwe79})

// 整库
allJSON, _ := registry.ExportJSON()
allCSV, _  := registry.ExportCSV()
_ = registry.ImportJSON(allJSON) // 反向导入
```

详见 [序列化](../sdk/serializer)。

---

## ⚠️ 错误处理

所有错误都是 `*CWEError` 体系，支持 `errors.Is` / `errors.As`：

```go
weakness, err := client.GetWeakness(ctx, 999999)
if err != nil {
    var apiErr *cweskills.APIError
    if errors.As(err, &apiErr) {
        fmt.Println("HTTP 状态码:", apiErr.StatusCode) // 404
    }
    var nfErr *cweskills.CWENotFoundError
    if errors.As(err, &nfErr) {
        fmt.Println("未找到 ID:", nfErr.ID)
    }
}
```

错误类型：`InvalidCWEIDError`、`CWENotFoundError`、`APIError`、`RateLimitError`、`ValidationError`、`ParseError`、`RelationshipError`。详见 [错误处理](./error-handling)。

---

## 📖 相关文档

- [四种接入方式总览](./integrations)
- [SDK API 总览](../sdk/overview)
- [错误处理](./error-handling)
- [性能与零依赖](./performance)
- [工作原理](./how-it-works)
