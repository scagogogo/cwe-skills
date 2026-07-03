---
title: 为什么需要 CWE Skills
outline: [2, 3]
---

# 💡 为什么需要 CWE Skills

CWE 是安全领域事实上的「弱点字典」，几乎所有 SAST 工具、漏洞管理平台、合规报告都会引用 CWE ID。但「用 CWE」这件事本身，处处是坑。本页把真实工程中最常见的几类痛点摆出来——CWE Skills 正是为消除这些痛点而生。

---

## 😖 痛点 1：CWE ID 的写法五花八门

同一个 `CWE-79`，在不同来源里可能是 `CWE-79`、`cwe-79`、`CWE79`、`CWE 79`、甚至纯数字 `79`、带前导零的 `079`。扫描器报告里写 `CWE-79`，工单系统里填 `cwe79`，合规表格里又是 `79`——**字符串直接比较永远是 False**。

```go
// 自己手写解析？正则、大小写、前导零、空格……每个边界都要处理
if strings.EqualFold(id, "CWE-79") { ... }   // ❌ 写死，且不通用
```

**CWE Skills 的答案**：一套覆盖所有合法写法的 `ParseCWEID` / `FormatCWEID` / `IsCWEID` / `ExtractCWEIDs` / `CompareCWEIDs`，统一规范化为 `CWE-NNN`。

```go
id, _ := cweskills.ParseCWEID("cwe-079")     // 79
formatted, _ := cweskills.FormatCWEID("79")  // "CWE-79"
ids := cweskills.ExtractCWEIDs("受 CWE-79 和 cwe89 影响") // ["CWE-79", "CWE-89"]
```

详见 [CWE ID 工具](../sdk/cwe-utils)。

---

## 😖 痛点 2：MITRE REST API 不好直接用

MITRE 提供了 REST API（`https://cwe-api.mitre.org/api/v1`），但它有几个让人头疼的特性：

1. **严格速率限制** —— 默认大约每 10 秒 1 个请求，批量查询极易撞限。
2. **关系类型不全** —— 在线 API 只返回部分关系类型，链式/组合等关系往往拿不到。
3. **裸 HTTP 调用繁琐** —— 要自己拼 URL、处理状态码、解析 JSON、做重试。
4. **错误不结构化** —— 出错时拿到的是 HTTP 状态码 + 一段文本，上层很难分类处理。

自己实现一套「带速率限制 + 自动重试 + 结构化错误」的客户端，工作量不小，且容易写错。

**CWE Skills 的答案**：`APIClient` 开箱即用，默认配置已是 MITRE 的速率限制（`rate=0.1, burst=1`），支持可配置重试，错误统一为 `CWEError` 体系。

```go
client := cweskills.NewAPIClient()
defer client.Close()
weakness, err := client.GetWeakness(ctx, 79)
parents, err := client.GetParents(ctx, 79)
```

详见 [MITRE REST API 客户端](../sdk/api-client) 与 [速率限制与重试](./rate-limit-retry)。

---

## 😖 痛点 3：权威列表散落各处，格式不一

「CWE Top 25 最危险软件弱点」「OWASP Top 10」「SANS Top 25」是安全优先级排序的事实标准。但它们的来源、年份、格式各不相同：

- CWE Top 25 是一串数字 ID 列表；
- OWASP Top 10 是 `A01:2021-...` 类别名到多个 CWE ID 的映射；
- SANS Top 25 又是另一份列表。

每个项目都自己复制粘贴一份，年份一更新就全部失效。

**CWE Skills 的答案**：内置 `CWETop25`（2024）、`OWASPTop10`（2021）、`SANSTop25`，并提供 `IsInTop25` / `IsInOWASPTop10` / `IsInSANSTop25` / `GetOWASPCategory` 等查询函数，一处维护、处处可用。

```go
cweskills.IsInTop25(79)               // true
cweskills.IsInOWASPTop10(79)          // true
cweskills.GetOWASPCategory(89)        // "A03:2021-Injection"
```

详见 [知名列表](../wellknown/cwe-top-25)。

---

## 😖 痛点 4：没有像样的离线方案

很多安全场景是**离线**的：内网扫描器、受控环境、CI 流水线里禁止外网访问。MITRE 提供的 XML 弱点目录（如 `cwec_v4.15.xml`）是离线场景唯一的完整数据源，但它有几百 MB、结构复杂、关系层层嵌套，自己写解析器既慢又容易漏字段。

而且，**在线 API 返回的关系类型不完整**——链式（CanPrecede/CanFollow）、依赖（Requires/RequiredBy）、对等（PeerOf/CanAlsoBe）等关系，只有在 XML 里才齐全。

**CWE Skills 的答案**：`XMLParser` 解析官方 XML，构建内存 `Registry` 并建立父子、对等、成员等多层索引；`Navigator` 在此之上提供比 API 更丰富的关系导航；所有查询完全离线、零网络依赖。

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()
nav := cweskills.NewNavigator(registry)
ancestors := nav.Ancestors(79)
path := nav.ShortestPath(79, 1)
```

详见 [离线 XML 目录解析](../sdk/xml-parser) 与 [关系导航](../sdk/navigator)。

---

## 😖 痛点 5：关系导航与层级树要自己造轮子

CWE 是一张**有向图**（甚至带层级、顺序、依赖、对等多种语义）。要回答「CWE-79 的所有祖先是谁」「CWE-79 到 CWE-1 的最短路径」「以 CWE-1 为根的整棵树有多少叶子」这类问题，自己写图遍历既繁琐又容易出 bug。

**CWE Skills 的答案**：`Navigator` 提供父/子/祖先/后代/同级/对等/前置/跟随/依赖/被依赖/链成员/组合成员/最短路径/关系深度等一整套导航；`BuildTree` / `BuildForest` / `BuildViewTree` 构建层次树，支持 DFS/BFS 遍历、路径查找、叶子枚举、深度统计。

```go
tree := cweskills.BuildTree(registry, 1)
leaves := tree.LeafNodes()
tree.Walk(func(n *cweskills.TreeNode) bool { /* DFS */ return true })
```

详见 [关系导航](../sdk/navigator) 与 [层次树构建](../sdk/tree)。

---

## 😖 痛点 6：AI 代理「想用 CWE」却没有顺手的工具

越来越多的安全分析由 AI 代理（Claude、GPT 等）完成。让 AI 自己去拼 curl、解析 JSON、记速率限制，既低效又不可靠——AI 擅长调用工具，不擅长处理 HTTP 边角。

**CWE Skills 的答案**：**Skills 接入**——把一段提示词放进 AI 代理的系统提示词或技能配置，AI 就能直接调用 `cwe` CLI 完成解析、查询、导航、树构建等全部操作，输出还可指定为 JSON 便于 AI 解析。

详见 [Skills 接入](./integration-skills)。

---

## 🎯 小结：CWE Skills 的价值主张

| 痛点 | CWE Skills 提供的能力 |
|------|----------------------|
| CWE ID 写法混乱 | 类型化 ID 工具，统一规范化 |
| MITRE API 难用 | 带速率限制/重试/结构化错误的客户端 |
| 权威列表分散 | 内置 Top 25 / OWASP / SANS，一键查询 |
| 无离线方案 | XML 解析 + 内存注册表 + 多层索引 |
| 关系导航/树要自造 | Navigator + BuildTree 全套图与树能力 |
| AI 代理无顺口工具 | Skills 提示词 + JSON 输出 |

> 一句话：**CWE Skills 让「用对 CWE」这件事从「每个团队重新实现一遍」变成「一行 go get / 一段提示词」。**

想看更详细的「痛点 → 方案」逐条对照，继续读 [解决了什么问题](./problem-solved)。
