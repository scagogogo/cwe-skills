---
title: 技能 03 — 知名列表
outline: [2, 3]
---

# 🏆 技能 03 — 知名列表

查询并判断 CWE 是否属于权威知名列表：**CWE Top 25 最危险软件弱点**、**OWASP Top 10 (2021)**、**SANS Top 25 最危险软件错误**。这些列表是安全优先级排序与合规的关键依据。

<Badge type="tip" text="优先级排序"/>
<Badge type="info" text="内置数据"/>
<Badge type="warning" text="合规关键"/>

---

## 🎯 技能目标

- 列出 CWE Top 25 / OWASP Top 10 / SANS Top 25
- 检查某个 CWE 是否属于任意知名列表
- 查询 CWE 对应的 OWASP 分类

---

## 💻 CLI 命令

### wellknown top25 / owasp / sans

```bash
cwe wellknown top25      # CWE Top 25 最危险软件弱点
cwe wellknown owasp      # OWASP Top 10 (2021) 及其关联 CWE
cwe wellknown sans       # SANS Top 25 最危险软件错误
```

### wellknown check — 成员检查

```bash
cwe wellknown check CWE-79 CWE-89 CWE-999
```

```text
CWE-79: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
CWE-89: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
CWE-999: 不在任何知名列表中
```

JSON 输出：

```json
[
  {"cwe_id": "CWE-79", "in_list": ["Top 25", "OWASP Top 10 (A03:2021-Injection)", "SANS Top 25"]},
  {"cwe_id": "CWE-999", "in_list": []}
]
```

---

## 🔧 SDK API

### 成员判断

```go
cweskills.IsInTop25(79)       // true
cweskills.IsInOWASPTop10(79)  // true
cweskills.IsInSANSTop25(79)   // true
```

### OWASP 分类

```go
cat := cweskills.GetOWASPCategory(79)      // "A03:2021-Injection"（首个匹配）
cats := cweskills.GetOWASPCategories(79)   // []string{...}（全部匹配）
```

### 知名视图与预定义列表

```go
cweskills.IsInWellKnownView(1000)  // true — View 1000 是知名研究视图

top25 := cweskills.CWETop25    // []int{79, 89, 352, 862, ...}
sans25 := cweskills.SANSTop25   // []int{119, 20, 79, ...}
owasp := cweskills.OWASPTop10   // map[string][]int
```

::: details 一个 CWE 可属于多个列表
CWE-79 同时出现在 Top 25、OWASP Top 10、SANS Top 25 三个列表里。`GetOWASPCategory` 返回首个匹配，`GetOWASPCategories` 返回全部。
:::

---

## 📝 示例

### 命令行

```bash
# 检查一批扫描结果是否命中 Top 25
cwe wellknown check CWE-79 CWE-89 CWE-352 CWE-999 -o json | jq '.[] | select(.in_list|length>0)'

# 列出 OWASP A03 注入类别的所有 CWE
cwe wellknown owasp -o json | jq '.[] | select(.category=="A03:2021-Injection")'
```

### Go

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    findings := []int{79, 89, 352, 999}
    for _, id := range findings {
        if cweskills.IsInTop25(id) {
            fmt.Printf("CWE-%d 命中 Top 25，优先修复！\n", id)
        }
    }
}
```

---

## 🤖 AI 代理使用提示

- 用户给出扫描结果或 CWE 列表时，AI 用 `cwe wellknown check` 标注优先级。
- 合规场景下，AI 用 `cwe wellknown owasp` 把 CWE 映射到 OWASP 类别。
- 风险评分：命中知名列表的发现应加权更高。

::: tip 典型用例
1. 安全优先级排序：聚焦 Top 25 弱点的修复。
2. 合规报告：把漏洞映射到 OWASP Top 10。
3. 风险评分：命中知名列表的发现权重更高。
4. 过滤：自动标记高优先级 CWE 条目。
:::

---

## 📖 相关文档

- [知名列表总览](../wellknown/overview)
- [CWE Top 25](../wellknown/cwe-top-25) · [OWASP Top 10](../wellknown/owasp-top-10) · [SANS Top 25](../wellknown/sans-top-25)
- [CLI: wellknown](../cli/wellknown) · [wellknown check](../cli/wellknown-check)
- [SDK: IsInTop25](../sdk/is-in-top-25) · [GetOWASPCategories](../sdk/get-owasp-categories)
- [返回 Skills 总览](./)
