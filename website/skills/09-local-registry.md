---
title: 技能 09 — 本地注册表
outline: [2, 3]
---

# 🗃️ 技能 09 — 本地注册表

加载 MITRE XML 目录到内存注册表，并查询本地数据。所有操作完全离线。

<Badge type="tip" text="离线"/>
<Badge type="info" text="需 XML 目录"/>

---

## 🎯 技能目标

- 加载 XML 目录并查看摘要
- 获取条目详情、检查存在性
- 列出视图 / 类别、查询父子/祖先后代/对等
- 导出注册表为 JSON / CSV

---

## 💻 CLI 命令

所有 registry 命令需 `--xml <file>`。

```bash
cwe registry load --xml cwec_latest.xml              # 加载并显示摘要
cwe registry get CWE-79 --xml cwec_latest.xml        # 获取条目详情
cwe registry contains CWE-79 CWE-89 --xml <file>     # 检查存在性
cwe registry list-views --xml <file>                 # 列出所有视图
cwe registry list-categories --xml <file>            # 列出所有类别
cwe registry parents CWE-79 --xml <file>             # 父级 ID
cwe registry children CWE-74 --xml <file>            # 子级 ID
cwe registry ancestors CWE-79 --xml <file>           # 所有祖先
cwe registry descendants CWE-74 --xml <file>         # 所有后代
cwe registry peers CWE-79 --xml <file>               # 对等 ID
cwe registry view-members 1000 --xml <file>          # 视图成员
cwe registry category-members 1 --xml <file>         # 类别成员
cwe registry member-of CWE-79 --xml <file>           # 79 属于哪些
cwe registry export --xml <file> --format json       # 导出 JSON
cwe registry export --xml <file> --format csv        # 导出 CSV
```

| Flag | 简写 | 说明 |
|------|------|------|
| `--xml` | `-x` | **（必填）** XML 目录路径 |
| `--format` | | 导出格式：json 或 csv |
| `--output-file` | | 写入文件而非 stdout（用于 export） |

---

## 🔧 SDK API

### 加载与查询

```go
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
registry.BuildIndexes()

cwe, ok := registry.Get(79)
views := registry.GetAllViews()
cats := registry.GetAllCategories()
parentIDs := registry.GetParentIDs(79)
childIDs := registry.GetChildIDs(74)
ancestorIDs := registry.GetAncestorIDs(79)
descIDs := registry.GetDescendantIDs(74)
peerIDs := registry.GetPeerIDs(79)
viewMembers := registry.GetViewMembers(1000)
catMembers := registry.GetCategoryMembers(1)
memberOf := registry.GetMemberOfIDs(79)
```

### 导出

```go
json, _ := registry.ExportJSON()
csv, _ := registry.ExportCSV()
```

::: tip BuildIndexes() 不能忘
解析后调用 `registry.BuildIndexes()` 构建多层索引（父子、祖先后代、对等等），否则关系查询不生效。
:::

---

## 📝 示例

### 命令行

```bash
# 加载并看摘要
cwe registry load --xml cwec_latest.xml

# 把整个目录导出为 JSON
cwe registry export --xml cwec_latest.xml --format json --output-file cwe.json
```

### Go

```go
package main

import (
    "fmt"
    "os"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    registry, err := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
    if err != nil {
        panic(err)
    }
    registry.BuildIndexes()

    cwe, ok := registry.Get(79)
    if ok {
        fmt.Println(cwe.Name)
    }

    data, _ := registry.ExportJSON()
    _ = os.WriteFile("cwe.json", data, 0644)
}
```

---

## 🤖 AI 代理使用提示

- 用户想离线查某 CWE 时，AI 用 `cwe registry get CWE-79 --xml <file>`。
- 批量导出数据用 `cwe registry export`，AI 可引导用户保存到文件。
- 列视图/类别用 `list-views` / `list-categories`，便于 AI 了解目录全貌。

::: warning 离线命令需先下载 XML
`registry` 系列命令依赖本地 XML 目录。若环境里没有 XML，AI 应提示用户从 MITRE 下载。
:::

---

## 📖 相关文档

- [技能 08 — 本地搜索与过滤](./08-local-search-filter)
- [技能 10 — 本地关系导航](./10-local-navigation)
- [CLI: registry](../cli/registry) · [registry load](../cli/registry-load) · [registry get](../cli/registry-get) · [registry export](../cli/registry-export)
- [SDK: Registry](../sdk/registry) · [NewXMLParser](../sdk/new-xml-parser) · [ExportJSON](../sdk/registry-json)
- [返回 Skills 总览](./)
