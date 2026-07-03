---
title: 技能 10 — 本地关系导航
outline: [2, 3]
---

# 🧭 技能 10 — 本地关系导航

用 Navigator 离线导航 CWE 关系。比在线 API 的 `cwe relations` 更丰富：含同级、对等、最短路径、关系深度等。

<Badge type="tip" text="离线"/>
<Badge type="info" text="需 XML 目录"/>
<Badge type="warning" text="比 API 更全"/>

---

## 🎯 技能目标

- 离线查询父/子/祖先/后代/同级/对等
- 查询顺序关系（CanPrecede/CanFollow）、依赖关系（Requires/RequiredBy）
- 查询最短路径、是否祖先、是否相关、关系深度

---

## 💻 CLI 命令

所有 nav 命令需 `--xml <file>`。

```bash
# 层级关系
cwe nav parents CWE-79 --xml <file>
cwe nav children CWE-74 --xml <file>
cwe nav ancestors CWE-79 --xml <file>
cwe nav descendants CWE-74 --xml <file>
cwe nav siblings CWE-79 --xml <file>

# 对等与顺序
cwe nav peers CWE-79 --xml <file>
cwe nav precede CWE-89 --xml <file>
cwe nav follow CWE-79 --xml <file>

# 依赖
cwe nav requires CWE-79 --xml <file>
cwe nav required-by CWE-79 --xml <file>
cwe nav can-also-be CWE-79 --xml <file>

# 复合
cwe nav chain-members 680 --xml <file>
cwe nav composite-members 680 --xml <file>

# 路径查询
cwe nav shortest-path CWE-79 CWE-1 --xml <file>
cwe nav is-ancestor CWE-1 CWE-79 --xml <file>
cwe nav is-related CWE-79 CWE-89 --xml <file>
cwe nav depth CWE-79 CWE-1 --xml <file>
```

| Flag | 简写 | 说明 |
|------|------|------|
| `--xml` | `-x` | **（必填）** XML 目录路径 |

---

## 🔧 SDK API

### 创建 Navigator

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
registry.BuildIndexes()
nav := cweskills.NewNavigator(registry)
```

### 关系查询

```go
parents := nav.Parents(79)
children := nav.Children(74)
ancestors := nav.Ancestors(79)
descendants := nav.Descendants(74)
siblings := nav.Siblings(79)
peers := nav.Peers(79)
```

### 顺序与依赖

```go
nav.CanPrecede(89)
nav.CanFollow(79)
nav.Requires(79)
nav.RequiredBy(79)
nav.CanAlsoBe(79)
nav.ChainMembers(680)
nav.CompositeMembers(680)
```

### 路径与判定

```go
path := nav.ShortestPath(79, 1)            // []int，无路径返回 nil
isAncestor := nav.IsAncestorOf(1, 79)      // bool
isRelated := nav.IsRelated(79, 89)         // bool
depth := nav.RelationshipDepth(79, 1)      // int，无关返回 -1
```

---

## 📝 示例

### 命令行

```bash
# 找 CWE-79 到根 CWE-1 的最短路径
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_latest.xml -o json

# 判断 CWE-1 是否是 CWE-79 的祖先
cwe nav is-ancestor CWE-1 CWE-79 --xml cwec_latest.xml
```

### Go

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    registry, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
    registry.BuildIndexes()
    nav := cweskills.NewNavigator(registry)

    path := nav.ShortestPath(79, 1)
    fmt.Printf("最短路径: %v\n", path)
    fmt.Printf("关系深度: %d\n", nav.RelationshipDepth(79, 1))
}
```

---

## 🤖 AI 代理使用提示

- 用户问「CWE-79 和 CWE-1 什么关系」时，AI 用 `cwe nav shortest-path` 和 `cwe nav is-related`。
- 查同级弱点用 `cwe nav siblings`，查替代方案用 `cwe nav can-also-be`。
- 离线导航不受速率限制，关系类型完整，优先于在线 `relations`。

::: tip 比在线 API 更全
Navigator 提供同级（siblings）、对等（peers）、最短路径、关系深度等 API 不支持的查询，且无速率限制。关系分析优先用离线导航。
:::

---

## 📖 相关文档

- [技能 06 — API 关系查询](./06-api-relationships)（在线版）
- [技能 11 — 本地树构建](./11-local-tree)
- [CLI: nav](../cli/nav) · [nav shortest-path](../cli/nav-shortest-path) · [nav is-ancestor](../cli/nav-is-ancestor)
- [SDK: Navigator](../sdk/navigator) · [ShortestPath](../sdk/nav-shortest-path)
- [返回 Skills 总览](./)
