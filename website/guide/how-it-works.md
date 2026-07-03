---
title: 工作原理
outline: [2, 3]
---

# 🧠 工作原理

要会用一个库，最好先看清它内部怎么转。本页讲清 CWE Skills 的整体架构、数据流，以及「在线」与「离线」两条数据路径如何各自又如何合流。

![CWE Skills 架构图](/architecture.png)

::: tip 图怎么看
架构图自上而下分三层：上层是你的应用 / AI 代理；中间是 CWE Skills 集成层（ID 工具、枚举、知名列表、注册表、导航、树、搜索、序列化）；下层是两条数据来源——MITRE REST API（在线）与 XML 弱点目录（离线）。集成层把两条路统一成同一套类型化对象。
:::

---

## 🧩 核心模块

CWE Skills 由若干个相互解耦的模块组成。每个模块只管一件事，模块之间通过明确的接口协作。

| 模块 | 职责 | 关键类型 / 函数 |
|------|------|----------------|
| 🆔 ID 工具 | CWE ID 的解析/格式化/验证/提取/比较 | `ParseCWEID` `FormatCWEID` `IsCWEID` `ExtractCWEIDs` `CompareCWEIDs` |
| 📚 枚举 | 类型化枚举与校验/解析/排序 | `Abstraction` `Structure` `Status` `RelationshipNature` `ConsequenceScope` `ViewType` … |
| 🏆 知名列表 | Top 25 / OWASP / SANS 与知名视图 | `CWETop25` `OWASPTop10` `SANSTop25` `IsInTop25` `CWEViewResearchConcepts` |
| 🌐 API 客户端 | 调用 MITRE REST API | `APIClient` `GetWeakness` `GetParents` `GetVersion` |
| ⚙️ HTTP 客户端 | 速率限制 + 重试 + 结构化错误 | `HTTPClient` `RateLimiter` `CWEError` |
| 📥 XML 解析器 | 解析官方 XML 目录 | `XMLParser` `ParseFile` |
| 🗄️ 注册表 | 内存存储 + 多层索引 | `Registry` `Register` `Get` `BuildIndexes` |
| 🧭 导航器 | 关系遍历与图算法 | `Navigator` `Ancestors` `ShortestPath` `RelationshipDepth` |
| 🌳 树 | 层次树构建与遍历 | `BuildTree` `BuildForest` `BuildViewTree` `TreeNode` |
| 🔍 搜索过滤 | 查找与多条件过滤 | `FindByKeyword` `Filter` `SortByID` `GroupByAbstraction` |
| 📦 序列化 | JSON/XML/CSV 互转 | `MarshalJSON` `MarshalXML` `MarshalCSV` `ExportJSON` |
| 💻 CLI | 40+ 子命令 | `cwe parse` `cwe show` `cwe nav` `cwe tree` … |

---

## 🌊 数据流：两条路径

CWE Skills 同时支持两条数据来源。无论走哪条，最终都汇入同一套类型化对象（`*CWE`、`*Category`、`*View`、`*CompoundElement`），上层代码无需感知数据来自哪里。

![数据流图](/data-flow.png)

### 路径 A：在线（MITRE REST API）

适合「只需要查一两个 CWE」「需要最新版本」「无法下载大 XML」的场景。

```text
调用方
  │  client.GetWeakness(ctx, 79)
  ▼
APIClient ──► HTTPClient ──► RateLimiter (令牌桶, 默认 0.1/s, burst=1)
  │              │              │  Allow() / Wait(ctx)
  │              │              ▼
  │              │         HTTP 请求 https://cwe-api.mitre.org/api/v1/cwe/79
  │              │              │
  │              │              ▼  状态码 >=500 ? 重试(maxRetries) : 返回
  │              │         解析 JSON → *CWE
  │              ▼
  └── 返回 *CWE / *Category / *View / []Relationship / *VersionResponse
```

::: warning 在线路径的关系类型不全
MITRE REST API 只返回**部分**关系类型（主要是父子层级）。链式（CanPrecede/CanFollow）、依赖（Requires/RequiredBy）、对等（PeerOf/CanAlsoBe）等关系，**只有 XML 路径才齐全**。需要这些关系时请走离线路径。
:::

### 路径 B：离线（XML 弱点目录）

适合「需要全量数据」「需要全部关系类型」「内网/CI 禁止外网」「要批量导航与建树」的场景。

```text
调用方
  │  cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
  ▼
XMLParser ──► encoding/xml 流式解析
  │              │
  │              ▼
  │         构建空 Registry，逐条 Register/RegisterCategory/RegisterView/RegisterCompoundElement
  │
  │  registry.BuildIndexes()
  ▼
Registry 建立多层索引：
  parentIndex(子→父)  childIndex(父→子)
  peerIndex(对等)      memberIndex(类/视图→成员)
  memberOfIndex(成员→所属)
  │
  ├──► Navigator  (Ancestors / ShortestPath / RelationshipDepth …)
  ├──► BuildTree  (层次树 / 森林 / 视图树)
  ├──► FindBy* / Filter (搜索过滤)
  └──► ExportJSON/CSV (序列化导出)
```

::: tip 离线路径才完整
离线 XML 包含 MITRE 维护的**全部 10 种关系类型**，因此 `Navigator` 的链式/依赖/对等导航只有在离线路径下才有数据。这也是为什么「关系导航」类 CLI 子命令都需要 `--xml <file>` 参数。
:::

---

## 🔗 模块协作：一个完整例子

下面这个例子同时用到了 ID 工具、API 客户端、XML 解析、注册表、导航、树、知名列表——串起几乎所有模块：

```go
package main

import (
    "context"
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    // 1) ID 工具：规范化输入
    id, _ := cweskills.ParseCWEID("cwe-79") // 79

    // 2) 知名列表：快速判定优先级
    if cweskills.IsInTop25(id) {
        fmt.Println("CWE-79 属于 Top 25，高优先级")
    }

    // 3) 在线路径：取最新详情（受速率限制）
    client := cweskills.NewAPIClient()
    defer client.Close()
    weakness, err := client.GetWeakness(context.Background(), id)
    if err == nil {
        fmt.Println("在线详情:", weakness.Name)
    }

    // 4) 离线路径：加载完整目录，做 API 做不到的事
    registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
    registry.BuildIndexes()

    // 5) 导航：祖先链 + 最短路径（仅离线有完整关系）
    nav := cweskills.NewNavigator(registry)
    fmt.Println("祖先数:", len(nav.Ancestors(id)))
    fmt.Println("到 CWE-1 的最短路径:", nav.ShortestPath(id, 1))

    // 6) 树：以 CWE-1 为根建树，统计叶子
    tree := cweskills.BuildTree(registry, 1)
    fmt.Println("CWE-1 子树叶子数:", len(tree.LeafNodes()))

    // 7) 序列化：导出整库供下游
    jsonBytes, _ := registry.ExportJSON()
    _ = jsonBytes
}
```

数据在模块间的流向：`ID 工具` 规范输入 → `知名列表` / `API 客户端`（在线）或 `XMLParser` → `Registry`（离线）→ `Navigator` / `BuildTree` / `Filter` 消费注册表 → `序列化` 导出。

---

## 🔄 在线 vs 离线：如何选择

| 维度 | 在线 API | 离线 XML |
|------|---------|---------|
| 数据新鲜度 | 实时 | 取决于 XML 版本 |
| 网络依赖 | 必须 | 完全离线 |
| 速率限制 | 受限（默认 ~0.1 req/s） | 无 |
| 关系类型 | 部分（主要是父子） | **全部 10 种** |
| 数据量 | 按需单条 | 全量（几百 MB） |
| 适合场景 | 查一两条、查版本、查最新 | 导航、建树、批量过滤、内网 |

详细对比与选型建议见 [在线 vs 离线模式](./online-offline)。

---

## ⚙️ 关键工程特性

- **零依赖核心**：上述除 CLI（依赖 `cobra`）外的全部模块，仅用 Go 标准库实现。见 [性能与零依赖](./performance)。
- **并发安全注册表**：`Registry` 用 `sync.RWMutex` 保护，可被多个 goroutine 并发读。
- **令牌桶速率限制**：`RateLimiter` 基于令牌桶，支持突发与阻塞等待。见 [速率限制与重试](./rate-limit-retry)。
- **结构化错误**：所有错误统一为 `CWEError` 体系，支持 `errors.Is/As`。见 [错误处理](./error-handling)。
- **双格式输出**：CLI 全部命令支持 `-o text|json`，便于人类与脚本/AI 共用。见 [输出格式](./output-format)。

---

## 📖 接下来

- 两种数据路径的取舍 → [在线 vs 离线模式](./online-offline)
- 马上跑一遍 → [快速开始](./quick-start)
- 各模块 API 细节 → [SDK API](../sdk/overview)
