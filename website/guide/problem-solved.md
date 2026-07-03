---
title: 解决了什么问题
outline: [2, 3]
---

# 🎯 解决了什么问题

本页把工程实践中遇到的真实痛点，与 CWE Skills 提供的解决方案**逐条对应**。每一条都说明「以前怎么做」→「现在怎么做」，让你清楚看到 CWE Skills 替你省掉了什么。

---

## 🆔 痛点 1：CWE ID 写法不统一

**以前**：不同系统对同一个 CWE 的写法各异——`CWE-79`、`cwe-79`、`CWE79`、`CWE 79`、`79`、`079`。字符串直接比较恒为 False，正则要自己写、前导零要自己剥、大小写要自己归一。

```go
// 自己写的脆弱代码
id := strings.TrimSpace(input)
id = strings.TrimPrefix(id, "CWE-")
id = strings.TrimPrefix(id, "cwe-")
num, err := strconv.Atoi(id)
// 还得处理 "CWE79"、"CWE 79"、空字符串、非法字符……
```

**现在**：`cweskills` 包提供一套覆盖所有合法写法的 ID 工具。

```go
num, _       := cweskills.ParseCWEID("cwe-079")            // 79
formatted, _ := cweskills.FormatCWEID("79")                // "CWE-79"
ok           := cweskills.IsCWEID("CWE79")                 // true
ids          := cweskills.ExtractCWEIDs("受CWE-79与cwe89影响") // ["CWE-79","CWE-89"]
cmp, _       := cweskills.CompareCWEIDs("CWE-79", "CWE-89") // -1
```

::: tip 为什么 CompareCWEIDs 必须按数值比较
字符串比较会得到 `"CWE-100" < "CWE-79"`（字典序 `'1'<'7'`），但数值上 `100 > 79`。CWE ID 必须按数字比较，否则排序就错了。
:::

**省掉了**：正则编写、边界测试、各项目重复实现。→ [CWE ID 工具](../sdk/cwe-utils)

---

## 📚 痛点 2：枚举值靠裸字符串传来传去

**以前**：抽象层级、状态、关系类型等概念在不同模块里就是字符串 `"Base"`、`"Stable"`、`"ChildOf"`，拼写错一个字母要到运行时才暴露，且无法穷举、无法校验、无法排序。

**现在**：每个概念都是一个**带方法的枚举类型**，提供 `IsValid()`、`ParseXxx()`、`AllXxxValues()`、`XxxOrder()`（排序权重）。

```go
a, _ := cweskills.ParseAbstraction("Base")   // AbstractionBase
a.IsValid()                                   // true
cweskills.AllAbstractionValues()              // [Pillar Class Base Variant]
a.AbstractionOrder()                          // 2，可用于排序

r := cweskills.RelationshipChildOf
r.IsHierarchical()  // true
r.IsSequential()    // false
r.IsDependency()    // false
r.IsPeer()          // false
```

覆盖的枚举：`Abstraction`、`Structure`、`Status`、`LikelihoodOfExploit`、`RelationshipNature`、`ConsequenceScope`、`ConsequenceImpact`、`ViewType`、`PlatformType`。

**省掉了**：字符串拼写 bug、手写 switch 校验、手写排序权重。→ [枚举参考](../enums/abstraction)

---

## 🏆 痛点 3：权威列表各自维护、年久失修

**以前**：CWE Top 25、OWASP Top 10、SANS Top 25 散落在不同文档，格式各异（Top 25 是数字列表，OWASP 是类别→多 ID 映射）。每个项目复制一份，年份一更新就全部失效，还得自己写「某 CWE 属于哪个 OWASP 类别」的查询。

**现在**：内置 `CWETop25`（2024）、`OWASPTop10`（2021）、`SANSTop25`，并提供查询函数与知名视图常量。

```go
cweskills.IsInTop25(79)            // true
cweskills.IsInOWASPTop10(89)       // true
cweskills.IsInSANSTop25(79)        // true
cweskills.GetOWASPCategory(89)     // "A03:2021-Injection"

// 知名视图常量
cweskills.CWEViewResearchConcepts         // 1000
cweskills.CWEViewDevelopmentConcepts      // 699
cweskills.CWEViewHardwareDesign           // 1199
cweskills.CWEViewCWECrossSection          // 888
cweskills.CWEViewComprehensiveDictionary  // 1400
```

**省掉了**：列表复制粘贴、年份同步、类别反查逻辑。→ [知名列表](../wellknown/cwe-top-25)

---

## 🌐 痛点 4：MITRE API 难用、易撞限

**以前**：直接 `http.Get` 调 MITRE API，要自己拼 URL、处理 429 速率限制、写重试、解析 JSON、把 HTTP 错误码翻译成业务错误。MITRE 默认大约每 10 秒 1 个请求，批量查询一不小心就被限流。

**现在**：`APIClient` 默认即配置好 MITRE 的速率限制（`rate=0.1, burst=1`）与 30 秒超时，并可通过 Option 灵活调整重试与速率。

```go
client := cweskills.NewAPIClient()
defer client.Close()

// 可选：自定义速率与重试
client = cweskills.NewAPIClient(
    cweskills.WithAPITimeout(60*time.Second),
    cweskills.WithAPIRateLimit(0.5, 1),
    cweskills.WithAPIRetry(3, 2*time.Second),
)

weakness, _ := client.GetWeakness(ctx, 79)
parents,   _ := client.GetParents(ctx, 79)
children,  _ := client.GetChildren(ctx, 79)
ancestors, _ := client.GetAncestors(ctx, 79)
version,   _ := client.GetVersion(ctx)
```

**省掉了**：速率限制器、重试循环、URL 拼接、错误翻译。→ [MITRE REST API 客户端](../sdk/api-client)、[速率限制与重试](./rate-limit-retry)

---

## 📥 痛点 5：离线场景只能啃巨大 XML

**以前**：内网、CI、受控环境禁止外网，只能用 MITRE 的 XML 弱点目录（几百 MB）。自己写 `encoding/xml` 解析，结构层层嵌套、关系散落各处，字段容易漏，更别提建索引了。

**现在**：`XMLParser` 解析官方 XML，构建内存 `Registry`，一次 `BuildIndexes()` 即建立父子、子父、对等、成员、所属等多层索引。

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()

cwe79, ok := registry.Get(79)
all   := registry.GetAll()
n     := registry.Size()
members := registry.GetViewMembers(cweskills.CWEViewResearchConcepts)
```

**省掉了**：XML 结构研究、字段映射、索引设计。→ [离线 XML 目录解析](../sdk/xml-parser)

---

## 🧭 痛点 6：关系导航与最短路径要自造

**以前**：CWE 是一张多语义的有向图。要回答「79 的所有祖先」「79 到 1 的最短路径」「79 是不是 1 的后代」，得自己写 BFS/DFS、自己管 visited、自己处理环。

**现在**：`Navigator` 在 `Registry` 之上提供一整套关系导航，且**关系类型比在线 API 更全**（含 CanPrecede/CanFollow/Requires/RequiredBy/PeerOf/CanAlsoBe 等）。

```go
nav := cweskills.NewNavigator(registry)
parents    := nav.Parents(79)
ancestors  := nav.Ancestors(79)
siblings   := nav.Siblings(79)
peers      := nav.Peers(79)
precede    := nav.CanPrecede(79)
requires   := nav.Requires(79)
chainM     := nav.ChainMembers(79)
path       := nav.ShortestPath(79, 1)   // []int
depth      := nav.RelationshipDepth(1, 79) // int
ok         := nav.IsAncestorOf(1, 79)  // bool
```

**省掉了**：图遍历算法、防环逻辑、关系类型补全。→ [关系导航](../sdk/navigator)

---

## 🌳 痛点 7：层次树构建与遍历重复造轮子

**以前**：要把 CWE 以树形展示或做层级统计，得自己递归建树、自己实现 DFS/BFS、自己算深度和叶子。

**现在**：`BuildTree` / `BuildForest` / `BuildViewTree` 一行建树，`TreeNode` 自带 `Walk`（DFS）、`WalkBFS`、`Find`、`Path`、`LeafNodes`、`MaxDepth`、`Count`。

```go
tree    := cweskills.BuildTree(registry, 1)
forest  := cweskills.BuildForest(registry)
viewTree := cweskills.BuildViewTree(registry, cweskills.CWEViewResearchConcepts)

leaves   := tree.LeafNodes()
maxDepth := tree.MaxDepth()
tree.Walk(func(n *cweskills.TreeNode) bool {
    fmt.Println(n.CWE.CWEID(), "depth=", n.Depth)
    return true
})
node := tree.Find(79)
chain := node.Path() // 从根到该节点的路径
```

**省掉了**：递归建树、遍历器、路径收集。→ [层次树构建](../sdk/tree)

---

## 🔍 痛点 8：多条件过滤与搜索各自实现

**以前**：要按「关键字 + 抽象=Base + 状态=Stable + 可能性=High」过滤，得自己写一堆 `if`，再自己写排序、分组、去重。

**现在**：`FindByKeyword` / `FindByAbstraction` 等单维度查找 + `Filter` 多条件过滤 + `SortByID/Name/Abstraction`、`GroupByAbstraction/Status/Likelihood`、`Deduplicate`。

```go
results := cweskills.FindByKeyword(registry, "Injection")
filtered := cweskills.Filter(results, cweskills.FilterOption{
    Abstraction: cweskills.AbstractionBase,
    Status:      cweskills.StatusStable,
})
sorted := cweskills.SortByID(filtered)
groups := cweskills.GroupByAbstraction(filtered)
unique := cweskills.Deduplicate(filtered)
```

**省掉了**：条件组合、排序/分组/去重样板代码。→ [搜索与过滤](../sdk/search)

---

## 📦 痛点 9：序列化格式不统一、字段裸露

**以前**：导出 CWE 给下游系统时，要么手写 JSON（字段名各处不一致），要么直接吐 XML struct tag，内部字段（如 `CWEType`、`URL`）不该外暴露却跟着导出。

**现在**：`MarshalJSON/XML/CSV` + `safeCWE` 安全模型，统一对外字段名；`Registry.ExportJSON/CSV`、`ImportJSON` 支持整库导入导出。

```go
data, _   := cweskills.MarshalJSON(cwe79)
cweBack, _ := cweskills.UnmarshalJSON(data)
jsonBytes, _ := registry.ExportJSON()
csvBytes, _  := registry.ExportCSV()
_ = registry.ImportJSON(jsonBytes)
```

**省掉了**：字段命名协商、安全字段过滤、CSV 列设计。→ [序列化](../sdk/serializer)

---

## 🤖 痛点 10：AI 代理用 CWE 没有顺口工具

**以前**：让 AI 代理分析 CWE，要么喂它一大段静态文档（知识会过时），要么让它自己 curl MITRE API（处理速率限制和 JSON 解析对 AI 是负担）。

**现在**：**Skills 接入**——把提示词放进 AI 系统提示词，AI 直接调用 `cwe` CLI，所有命令支持 `-o json` 输出结构化结果，AI 解析无障碍。

```text
[系统提示词里包含 Skills 提示词]
→ AI 自主决定调用：cwe wellknown check CWE-79 -o json
→ 拿到 JSON，继续推理
```

**省掉了**：为 AI 写专用集成、喂静态文档、教 AI 处理 HTTP。→ [Skills 接入](./integration-skills)

---

## 📊 一表速览

| # | 痛点 | CWE Skills 方案 | 相关文档 |
|---|------|----------------|---------|
| 1 | CWE ID 写法混乱 | `ParseCWEID`/`FormatCWEID`/`ExtractCWEIDs`/`CompareCWEIDs` | [ID 工具](../sdk/cwe-utils) |
| 2 | 枚举靠裸字符串 | 类型化枚举 + `IsValid`/`Parse`/`AllValues`/`Order` | [枚举参考](../enums/abstraction) |
| 3 | 权威列表分散 | 内置 Top 25/OWASP/SANS + 查询函数 | [知名列表](../wellknown/cwe-top-25) |
| 4 | MITRE API 难用 | `APIClient` + 速率限制 + 重试 + 结构化错误 | [API 客户端](../sdk/api-client) |
| 5 | 离线啃 XML | `XMLParser` → `Registry` + 多层索引 | [XML 解析](../sdk/xml-parser) |
| 6 | 关系导航自造 | `Navigator` 全套关系 + 最短路径 + 深度 | [导航器](../sdk/navigator) |
| 7 | 层次树自造 | `BuildTree`/`BuildForest`/`BuildViewTree` + 遍历 | [树构建](../sdk/tree) |
| 8 | 过滤搜索自造 | `FindBy*` + `Filter` + 排序/分组/去重 | [搜索过滤](../sdk/search) |
| 9 | 序列化不统一 | `MarshalJSON/XML/CSV` + `safeCWE` 安全模型 | [序列化](../sdk/serializer) |
| 10 | AI 无顺口工具 | Skills 提示词 + `cwe` CLI + JSON 输出 | [Skills 接入](./integration-skills) |

接下来推荐看 [工作原理](./how-it-works) 了解这些模块如何协作。
