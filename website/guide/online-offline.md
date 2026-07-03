---
title: 在线 vs 离线模式
outline: [2, 3]
---

# 🌐 在线 vs 离线模式

CWE Skills 同时支持**两条数据路径**：在线（MITRE REST API）与离线（本地 XML 弃点目录）。两者最终都汇入同一套类型化对象（`*CWE`、`*Category`、`*View`），但数据新鲜度、关系完整度、速率限制、网络依赖截然不同。理解差异、按场景选对路径，是用好 CWE Skills 的关键。

::: tip 两条路殊途同归
无论在线还是离线，最终拿到的都是同一个 `*cweskills.CWE` 结构体。上层代码（导航、树、搜索）无需感知数据来源——但**关系导航类能力只有离线路径有完整数据**。
:::

---

## 📊 对比总表

| 维度 | 🌐 在线 API | 📥 离线 XML |
|------|------------|------------|
| **数据来源** | `https://cwe-api.mitre.org/api/v1` | 本地 `cwec_v4.15.xml` 文件 |
| **数据新鲜度** | 实时（MITRE 最新版） | 取决于下载的 XML 版本 |
| **网络依赖** | 必须能访问 MITRE API | 完全离线，零网络 |
| **速率限制** | 受限（默认 ~0.1 req/s，令牌桶） | 无限制 |
| **关系类型** | 部分（主要父子层级） | **全部 10 种** |
| **数据量** | 按需单条/批量 | 全量（几百 MB） |
| **首次准备** | 无需下载 | 需下载大 XML |
| **适合场景** | 查一两条、查版本、查最新 | 导航、建树、批量过滤、内网/CI |
| **CLI 参数** | 无需 `--xml` | 需要 `--xml <file>` |

::: warning 关系完整性是最大差异
MITRE REST API **只返回部分关系类型**（主要是 `ChildOf`/`ParentOf` 父子层级）。链式（`CanPrecede`/`CanFollow`）、依赖（`Requires`/`RequiredBy`）、对等（`PeerOf`/`CanAlsoBe`）等关系，**只有离线 XML 才齐全**（共 10 种）。需要这些关系时必须走离线。
:::

---

## 🌐 在线模式（MITRE REST API）

### 何时用

- 只需要查一两个 CWE 的详情
- 需要确认 MITRE API 当前版本
- 无法/不愿下载几百 MB 的 XML
- 想要最新数据（XML 可能滞后）

### 用法（CLI）

```bash
cwe show CWE-79              # 取弱点详情
cwe relations CWE-79         # 在线关系（仅父子）
cwe version mitre            # MITRE API 版本
```

### 用法（SDK）

```go
client := cweskills.NewAPIClient()
defer client.Close()

weakness, _ := client.GetWeakness(ctx, 79)
parents, _  := client.GetParents(ctx, 79)
version, _  := client.GetVersion(ctx)
```

::: danger 受速率限制
在线模式默认速率约 0.1 req/s（每 10 秒 1 个请求）。批量循环调用会触发限流，CLI/SDK 会自动等待或重试。详见 [速率限制与重试](./rate-limit-retry)。
:::

---

## 📥 离线模式（XML 弃点目录）

### 何时用

- 需要完整关系（链式/依赖/对等/复合）
- 需要导航、建树、最短路径等图算法
- 需要批量搜索与多条件过滤
- 内网/CI 环境禁止外网
- 需要全量数据做离线分析

### 准备 XML

1. 访问 <https://cwe.mitre.org/data/downloads.html>
2. 下载最新版（如 `cwec_v4.15.xml`，文件较大）
3. 放到任意路径

```bash
# 验证可解析
cwe stats --xml cwec_v4.15.xml
```

### 用法（CLI）

```bash
cwe stats --xml cwec_v4.15.xml
cwe registry get CWE-79 --xml cwec_v4.15.xml
cwe nav ancestors CWE-79 --xml cwec_v4.15.xml
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_v4.15.xml
cwe tree build CWE-1 --xml cwec_v4.15.xml
cwe search --xml cwec_v4.15.xml --keyword Injection
cwe filter --xml cwec_v4.15.xml --abstraction Base --status Stable
```

### 用法（SDK）

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes() // 建立多层索引

nav := cweskills.NewNavigator(registry)
nav.Ancestors(79)
nav.ShortestPath(79, 1)

tree := cweskills.BuildTree(registry, 1)
```

::: tip 离线更快更全
离线命令不受速率限制，关系类型完整。鼓励优先用离线做关系分析。XML 加载到内存后，所有查询都是内存级速度。
:::

---

## 🧭 何时用哪个

```text
你的需求是？
│
├─ 查一两个 CWE 详情 / 查 API 版本 ─────► 在线（无需准备 XML）
│
├─ 只需要父子层级关系 ───────────────► 在线（够用）
│
├─ 需要链式/依赖/对等/复合关系 ───────► 离线（在线没有）
│
├─ 导航 / 建树 / 最短路径 ───────────► 离线（需要完整图）
│
├─ 批量搜索 / 多条件过滤 ───────────► 离线（不受限流）
│
├─ 内网 / CI 禁止外网 ────────────────► 离线（唯一选择）
│
└─ 想要最新数据 / XML 滞后 ──────────► 在线（实时）
```

---

## 🔀 混合用法：API 查详情 + XML 导航

实际项目里，**两条路径常组合使用**：在线 API 取最新详情，离线 XML 做完整导航。这是 CWE Skills 推荐的「最佳实践」之一。

### 典型场景

> 想要 CWE-79 的最新描述 + 它到 CWE-1 的完整祖先链。

- 最新描述 → 在线 API（实时）
- 完整祖先链 → 离线 XML（含全部关系）

### SDK 示例

```go
// 在线：取最新详情
client := cweskills.NewAPIClient()
defer client.Close()
weakness, _ := client.GetWeakness(ctx, 79) // 最新描述

// 离线：完整祖先链
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()
nav := cweskills.NewNavigator(registry)
ancestors := nav.Ancestors(79) // 完整祖先链

fmt.Println("名称:", weakness.Name)
fmt.Println("祖先链:", ancestors)
```

### CLI 示例

```bash
# 在线取详情
cwe show CWE-79 -o json > cwe79.json

# 离线取祖先链
cwe nav ancestors CWE-79 --xml cwec_v4.15.xml -o json > ancestors.json

# 用 jq 合并
jq -s '.[0] as $w | {weakness: $w, ancestors: .[1]}' cwe79.json ancestors.json
```

::: tip 详情在线、关系离线
这是一个实用的经验法则：**单条详情查在线（新鲜），关系导航查离线（完整）**。两者各取所长。
:::

---

## ⚠️ 常见误区

::: danger 在线做不了完整导航
`Navigator` 依赖注册表里的完整关系，**只在离线路径下有数据**。用在线 API 的 `GetParents`/`GetChildren` 只能拿到父子层级，链式/依赖/对等都拿不到。需要这些就必须离线。
:::

::: warning XML 版本会过期
离线 XML 是某个时间点的快照。MITRE 发布新版后，你的本地 XML 会滞后。定期从 <https://cwe.mitre.org/data/downloads.html> 更新。
:::

::: info 在线 vs 离线不互斥
两条路径可以同时用、数据互通。在线拿到的 `*CWE` 和离线拿到的 `*CWE` 是同一个类型，可以混用。
:::

---

## 📖 相关文档

- [工作原理](./how-it-works)（两条路径的数据流详解）
- [安装](./installation)（下载离线 XML）
- [速率限制与重试](./rate-limit-retry)
- [Go SDK 接入](./integration-sdk)
- [CLI 接入](./integration-cli)
