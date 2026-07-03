---
title: 技能 06 — API 关系查询
outline: [2, 3]
---

# 🧭 技能 06 — API 关系查询

通过 MITRE CWE REST API 查询 CWE 关系：父级、子级、祖先（传递父级至根）、后代（传递子级至叶）。理解这些关系是导航 CWE 层次结构的关键。

<Badge type="tip" text="在线 API"/>
<Badge type="warning" text="需网络"/>

---

## 🎯 技能目标

- 在线查询某 CWE 的父级 / 子级弱点
- 查询全部祖先 / 全部后代（传递关系）
- 用视图 ID（如 1000 Research Concepts）过滤关系

---

## 💻 CLI 命令

### relations parents / children

```bash
cwe relations parents CWE-79
cwe relations parents CWE-79 --view-id 1000
cwe relations children CWE-74
```

```text
CWE-79 的 父级弱点 (1 项):
  ChildOf -> CWE-74 (View: 1000)
```

### relations ancestors / descendants

```bash
cwe relations ancestors CWE-79     # 所有祖先（传递至根）
cwe relations descendants CWE-74   # 所有后代（传递至叶）
```

### Flags

| Flag | 默认值 | 说明 |
|------|--------|------|
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |
| `--view-id` | `0` | 按视图过滤（仅 parents/children） |

---

## 🔧 SDK API

### GetParents / GetChildren

```go
parents, _ := client.GetParents(ctx, 79)
parents, _ = client.GetParents(ctx, 79, 1000)   // 带视图过滤

children, _ := client.GetChildren(ctx, 74)
children, _ = client.GetChildren(ctx, 74, 1000) // 带视图过滤
```

`GetParents`/`GetChildren` 用可变参数支持可选 `viewID`。

### GetAncestors / GetDescendants

```go
ancestors, _ := client.GetAncestors(ctx, 79)     // 不支持 viewID
descendants, _ := client.GetDescendants(ctx, 74) // 不支持 viewID
```

::: warning 祖先/后代不支持视图过滤
`GetAncestors` 与 `GetDescendants` 不接受 `viewID` 参数。
:::

### Relationship 结构

```go
type Relationship struct {
    Nature RelationshipNature  // ChildOf, ParentOf 等
    CWEID  int                 // 关联的 CWE ID
    ViewID int                 // 视图上下文（无则 0）
}
```

### RelationshipNature 分类

| 类别 | 值 | 说明 |
|------|----|------|
| 层级 | `ChildOf` · `ParentOf` | CWE 树中的父子 |
| 顺序 | `CanPrecede` · `CanFollow` | 时序先后 |
| 依赖 | `Requires` · `RequiredBy` | 功能依赖 |
| 对等 | `PeerOf` · `CanAlsoBe` | 相关替代 |
| 成员 | `MemberOf` · `HasMember` | 类别/视图成员 |

---

## 📝 示例

### 命令行

```bash
# 看 CWE-79 的完整祖先链
cwe relations ancestors CWE-79 -o json | jq '.[].cwe_id'

# 用 Research Concepts 视图过滤父子
cwe relations children CWE-74 --view-id 1000 -o json
```

### Go

```go
package main

import (
    "context"
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    client := cweskills.NewAPIClient()
    defer client.Close()

    anc, _ := client.GetAncestors(context.Background(), 79)
    for _, r := range anc {
        fmt.Printf("祖先: CWE-%d (%s)\n", r.CWEID, r.Nature)
    }
}
```

---

## 🤖 AI 代理使用提示

- 用户问「CWE-79 的根因是什么」时，AI 用 `cwe relations ancestors CWE-79` 往上追溯。
- 影响分析用 `cwe relations descendants` 看一个基类弱点影响哪些子类。
- 标准层次用 `--view-id 1000`（Research Concepts）过滤。

::: tip 典型用例
1. 影响分析：查某基类弱点的所有后代，理解影响范围。
2. 根因分析：沿祖先链走到根类别。
3. 视图过滤：用 `--view-id 1000` 取标准层次。
4. 合规映射：追溯特定弱点与类别间的关系。
:::

::: tip 离线更全
关系导航的离线版（技能 10）提供同级、对等、最短路径、关系深度等更丰富的查询，不受速率限制。优先用离线命令做关系分析。
:::

---

## 📖 相关文档

- [技能 10 — 本地关系导航](./10-local-navigation)（离线版，更丰富）
- [CLI: relations](../cli/relations) · [relations parents](../cli/relations-parents) · [relations ancestors](../cli/relations-ancestors)
- [SDK: GetParents](../sdk/api-parents-children) · [GetAncestors](../sdk/api-ancestors-descendants)
- [返回 Skills 总览](./)
