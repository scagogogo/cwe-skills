---
title: 性能与零依赖
outline: [2, 3]
---

# 🚀 性能与零依赖

CWE Skills 的核心 SDK 有两个工程承诺：**零第三方依赖**（仅用 Go 标准库）与**内存注册表 + 多层索引**（查询内存级速度）。本页讲清这两点的设计取舍、性能特征，以及处理大 XML 时的注意事项。

::: tip 零依赖是认真的承诺
核心 SDK 的 `go.mod` 里**没有任何第三方 require**。这意味着供应链安全风险最小、编译产物小、跨平台编译无障碍、审计简单。唯一的例外是 CLI（用了 `cobra` 做命令行解析），但 CLI 是独立二进制，不影响 SDK 用户。
:::

---

## 📦 零第三方依赖

### 范围

| 模块 | 依赖 | 说明 |
|------|------|------|
| ID 工具 | 仅标准库 | 解析/格式化/验证/提取/比较 |
| 枚举 | 仅标准库 | 类型化枚举 |
| 知名列表 | 仅标准库 | 内置静态列表 |
| API 客户端 | 仅标准库 | `net/http` |
| XML 解析器 | 仅标准库 | `encoding/xml` |
| 注册表 | 仅标准库 | `sync`（RWMutex） |
| 导航 / 树 / 搜索 | 仅标准库 | 图算法纯手写 |
| 序列化 | 仅标准库 | `encoding/json` / `encoding/xml` / `encoding/csv` |
| CLI | `cobra` | 命令行框架（仅 CLI 二进制） |

```go
// go.mod（核心 SDK 部分）
module github.com/scagogogo/cwe-skills
go 1.25
// 无 require —— 零依赖
```

### 为什么坚持零依赖

- **供应链安全**：无第三方包 = 无供应链投毒风险，适合安全敏感场景。
- **编译干净**：`go build` 产物只链标准库，体积小、启动快。
- **跨平台**：标准库全平台支持，30+ 平台预编译无障碍。
- **审计简单**：安全团队审查时只需看标准库用法，不用追第三方库的 CVE。

::: warning CLI 是例外
CLI（`cmd/cwe`）用了 `cobra`，因为它需要子命令、flag、help 的成熟方案。但 `cobra` 只在 CLI 二进制里，**SDK 用户 import `cweskills` 包时不会引入 cobra**。SDK 包本身仍是零依赖。
:::

---

## 🗄️ 内存注册表

离线路径把整个 XML 弱点目录加载到内存的 `Registry`：

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes() // 建立多层索引
```

### 存储模型

```text
Registry
├── weaknesses   map[int]*CWE          // ID → 弱点
├── categories   map[int]*Category     // ID → 类别
├── views        map[int]*View         // ID → 视图
├── compounds    map[int]*CompoundElement
└── 索引（BuildIndexes 后建立）
    ├── parentIndex    map[int][]int   // 子 → 父列表
    ├── childIndex     map[int][]int   // 父 → 子列表
    ├── peerIndex      map[int][]int   // 对等关系
    ├── memberIndex    map[int][]int   // 类/视图 → 成员
    └── memberOfIndex  map[int][]int   // 成员 → 所属
```

::: info 内存换速度
全量 XML 几百 MB，加载到内存后所有查询都是 `map` 查找，**纳秒到微秒级**。这是离线导航/搜索能秒级完成的基础。代价是内存占用（一次性加载几百 MB）。
:::

---

## ⚡ 索引加速

`BuildIndexes()` 是性能关键。调用前，查「CWE-79 的祖先」要递归遍历全部弱点的关系字段，O(n) 起步；调用后，沿 `parentIndex` 链查，O(深度)。

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")

// ❌ 没建索引：每次导航都全表扫描
nav := cweskills.NewNavigator(registry)
nav.Ancestors(79) // 慢

// ✅ 建索引后：沿索引链查，快
registry.BuildIndexes()
nav2 := cweskills.NewNavigator(registry)
nav2.Ancestors(79) // 快
```

::: danger 一定要调 BuildIndexes
`BuildIndexes()` 是**一次性 O(n) 构建、永久 O(1)/O(深度) 查询**的投入。任何导航、建树、关系查询前都必须调用它，否则性能极差。这是一条强约定。
:::

### 索引覆盖的查询

| 查询 | 用到的索引 | 复杂度 |
|------|-----------|--------|
| `Parents(79)` / `Children(79)` | parent/childIndex | O(1) |
| `Ancestors(79)` | parentIndex 链 | O(深度) |
| `Descendants(79)` | childIndex 链 | O(子树) |
| `ShortestPath(a, b)` | parent/childIndex BFS | O(图大小) |
| `GetViewMembers(1000)` | memberIndex | O(1) |
| `GetMemberOfIDs(79)` | memberOfIndex | O(1) |

---

## 🔒 并发安全

`Registry` 用 `sync.RWMutex` 保护内部 map：

- **读操作**（`Get`/`GetAll`/`Contains`/导航/搜索）：持读锁，可多 goroutine 并发读。
- **写操作**（`Register`/`RegisterCategory`/`ImportJSON`）：持写锁，互斥。

```go
// 并发读安全
var wg sync.WaitGroup
for _, id := range ids {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        _ = registry.Get(id) // 并发读 OK
    }(id)
}
wg.Wait()
```

::: warning 注册阶段不要并发读
`BuildIndexes()` 之后的**查询阶段**可安全并发读。但**注册阶段**（边 `Register` 边查）不行——写锁会阻塞读，且索引未建好时查询结果不完整。建议：先加载完所有条目 + 建好索引，再开并发查询。
:::

---

## 📈 Benchmark 意识

CWE Skills 在关键路径上注重性能，使用 Go 标准的 `testing.B` 做基准测试。设计上的性能考量：

- **map 而非 slice 查找**：注册表用 `map[int]*CWE`，O(1) 查找。
- **预建索引**：避免每次查询都递归扫描。
- **流式 XML 解析**：`encoding/xml` 的 `decoder.Token()` 流式读取，不全量 DOM。
- **零反射热路径**：枚举比较用类型化常量，不用反射。

::: tip 自己跑 Benchmark
项目里有 `*_test.go` 的 `Benchmark*` 函数。用 `go test -bench=. -benchmem` 可看每操作耗时与内存分配。关注 `ns/op` 和 `allocs/op`。
:::

---

## 🐘 大 XML 解析注意

`cwec_v4.15.xml` 是几百 MB 的大文件，解析时有几点注意：

### 1. 一次性加载的内存开销

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
// 此刻内存里持有全部 CWE + 索引，峰值内存可能达 1~2 GB
```

::: warning 容器内存要够
在 CI 容器或小内存机器上解析大 XML，要给足内存（建议 ≥ 2GB）。OOM 会导致解析中途失败。生产环境用 `cwe stats --xml <file>` 先验证可解析。
:::

### 2. 流式解析，不全量 DOM

`XMLParser` 用 `encoding/xml` 的 `decoder.Token()` 流式读取，逐条 `Register`，**不全量构建 DOM 树**。这把峰值内存从「DOM 全树 + 注册表」降到「注册表 + 当前条目」。

### 3. 只解析一次，复用注册表

```go
// ✅ 解析一次，处处复用
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()

// 长生命周期里到处用同一个 registry
nav := cweskills.NewNavigator(registry)
tree := cweskills.BuildTree(registry, 1)
results := cweskills.FindByKeyword(registry, "Injection")
```

::: danger 不要每次查询都重新解析 XML
解析 + 建索引是 O(n) 重活。把 `registry` 做成长生命周期对象（如服务启动时解析一次，全局复用），**绝对不要**在每次查询时重新 `ParseFile`。
:::

### 4. CLI 每次调用都会重新解析

CLI 是独立进程，每次 `cwe nav ... --xml <file>` 都会重新加载 XML。这是 CLI 的固有开销（约几秒）。**高频查询请用 SDK 长驻进程**，或用 [Skills](./integration-skills) 让 AI 批量调用。

::: tip SDK 长驻 vs CLI 反复启动
- SDK：解析一次，后续查询内存级速度——适合服务/长驻进程。
- CLI：每次启动都解析——适合一次性脚本，不适合高频循环。
:::

---

## 📊 性能特征速查

| 操作 | 路径 | 典型耗时 | 说明 |
|------|------|----------|------|
| `ParseCWEID` | 纯计算 | 纳秒级 | 字符串解析 |
| `IsInTop25` | 内置列表 | 纳秒级 | map 查找 |
| `ParseFile`（大 XML） | 离线 | 几秒 | 一次性 |
| `BuildIndexes` | 离线 | <1 秒 | 一次性 |
| `Get(id)`（建索引后） | 离线 | 纳秒级 | map 查找 |
| `Ancestors(79)` | 离线 | 微秒级 | 索引链 |
| `ShortestPath` | 离线 | 微秒级 | BFS |
| `FindByKeyword` | 离线 | 毫秒级 | 全表扫 |
| `GetWeakness` | 在线 | 受速率限制 | HTTP + 限流 |

---

## 📖 相关文档

- [工作原理](./how-it-works)（架构与数据流）
- [在线 vs 离线模式](./online-offline)
- [Go SDK 接入](./integration-sdk)
- [CLI 接入](./integration-cli)（CLI 每次重新解析的开销）
- [安装](./installation)
