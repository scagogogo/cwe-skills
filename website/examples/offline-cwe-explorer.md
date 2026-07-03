---
title: 实战 — 离线 CWE 浏览器
outline: [2, 3]
---

# 🔌 实战 — 离线 CWE 浏览器

用 MITRE XML 目录构建一个命令行离线浏览器：加载 XML、构建层次树、按 ID/关键字/导航查询，完全离线。

<Badge type="tip" text="SDK 实战"/>
<Badge type="info" text="离线"/>

---

## 🎬 场景

内网/断网环境需要一个 CWE 查询工具：输入 CWE ID 看详情、输入关键字搜索、输入两个 ID 看最短路径。所有数据来自本地 XML 目录。

---

## 📋 前置准备

```bash
# 下载 MITRE XML 目录
curl -O https://cwe.mitre.org/data/xml/cwec_latest.xml.zip
unzip cwec_latest.xml.zip

go get github.com/scagogogo/cwe-skills
```

---

## 💻 完整代码

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/scagogogo/cwe-skills"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("用法: explorer <cwec_latest.xml>")
        os.Exit(1)
    }
    xmlPath := os.Args[1]

    // 1. 加载 XML 目录
    registry, err := cweskills.NewXMLParser().ParseFile(xmlPath)
    if err != nil {
        panic(err)
    }
    registry.BuildIndexes()
    nav := cweskills.NewNavigator(registry)

    fmt.Printf("已加载 %d 个 CWE 条目\n", len(registry.GetAllCWEs()))
    fmt.Println("命令: show <id> | search <keyword> | path <from> <to> | ancestors <id> | quit")

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        line := strings.TrimSpace(scanner.Text())
        parts := strings.Fields(line)
        if len(parts) == 0 {
            continue
        }
        switch parts[0] {
        case "quit", "exit":
            return
        case "show":
            if len(parts) < 2 {
                continue
            }
            id, _ := cweskills.ParseCWEID(parts[1])
            if cwe, ok := registry.Get(id); ok {
                fmt.Printf("CWE-%d %s\n  抽象: %s  状态: %s\n  %s\n",
                    cwe.ID, cwe.Name, cwe.Abstraction, cwe.Status, cwe.Description)
            }
        case "search":
            if len(parts) < 2 {
                continue
            }
            kw := strings.Join(parts[1:], " ")
            results := cweskills.FindByKeyword(registry, kw)
            fmt.Printf("找到 %d 条:\n", len(results))
            for i, c := range results {
                if i >= 10 {
                    fmt.Println("  ...（仅显示前 10 条）")
                    break
                }
                fmt.Printf("  CWE-%d %s\n", c.ID, c.Name)
            }
        case "path":
            if len(parts) < 3 {
                continue
            }
            from, _ := cweskills.ParseCWEID(parts[1])
            to, _ := cweskills.ParseCWEID(parts[2])
            p := nav.ShortestPath(from, to)
            if p == nil {
                fmt.Println("无路径")
            } else {
                ids := make([]string, len(p))
                for i, x := range p {
                    ids[i] = fmt.Sprintf("CWE-%d", x)
                }
                fmt.Printf("最短路径(%d跳): %s\n", len(p)-1, strings.Join(ids, " → "))
            }
        case "ancestors":
            if len(parts) < 2 {
                continue
            }
            id, _ := cweskills.ParseCWEID(parts[1])
            for _, a := range nav.Ancestors(id) {
                fmt.Printf("  CWE-%d %s\n", a.ID, a.Name)
            }
        default:
            fmt.Println("未知命令")
        }
    }
}
```

---

## ▶️ 运行步骤

```bash
# 1. 编译运行
go run main.go cwec_latest.xml

# 2. 交互查询
# > show CWE-79
# > search Injection
# > path CWE-79 CWE-1
# > ancestors CWE-79
# > quit
```

---

## 📤 输出示例

```text
已加载 1298 个 CWE 条目
命令: show <id> | search <keyword> | path <from> <to> | ancestors <id> | quit
> show CWE-79
CWE-79 Improper Neutralization of Input During Web Page Generation...
  抽象: Base  状态: Stable
  The product does not neutralize...
> path CWE-79 CWE-1
最短路径(4跳): CWE-79 → CWE-74 → CWE-707 → CWE-664 → CWE-1
> search Injection
找到 37 条:
  CWE-74 Injection
  CWE-79 Cross-site Scripting
  ...
```

---

## 🧩 扩展思路

- **加树视图**：用 `BuildTree(registry, id)` + `Walk` 在浏览器里展示子树。
- **加分页**：`FindByKeyword` 结果多时做分页显示。
- **持久化索引**：把 `registry.ExportJSON()` 存盘，下次直接 `ImportJSON` 加载，跳过 XML 解析。
- **做成 TUI**：用 [bubbletea](https://github.com/charmbracelet/bubbletea) 把命令行交互升级为终端 UI。

---

## 📖 相关文档

- [技能 09 — 注册表](../skills/09-local-registry) · [技能 10 — 导航](../skills/10-local-navigation) · [技能 11 — 树](../skills/11-local-tree)
- [SDK: NewXMLParser](../sdk/new-xml-parser) · [NewNavigator](../sdk/navigator) · [ShortestPath](../sdk/nav-shortest-path)
- [返回示例总览](./)
