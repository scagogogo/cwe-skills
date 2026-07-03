---
title: 技能 02 — CWE ID 提取与比较
outline: [2, 3]
---

# 🔍 技能 02 — CWE ID 提取与比较

从自由文本中提取 CWE ID，并比较两个 CWE ID 的大小。适合处理漏洞报告、安全公告、合规文档。

<Badge type="tip" text="文本处理"/>
<Badge type="info" text="无需网络"/>

---

## 🎯 技能目标

- 从一段文本里抽出所有 `CWE-NNN` 引用
- 取出第一个匹配的 CWE ID
- 比较两个 CWE ID 的大小关系

---

## 💻 CLI 命令

### extract — 提取

```bash
cwe extract "This system is affected by CWE-79 and CWE-89"
```

```text
找到 2 个CWE ID:
  CWE-79
  CWE-89
```

JSON 输出：

```json
{
  "text": "This system is affected by CWE-79 and CWE-89",
  "ids": ["CWE-79", "CWE-89"],
  "count": 2
}
```

### compare — 比较

```bash
cwe compare CWE-79 CWE-89
cwe compare CWE-79 CWE-79
cwe compare CWE-89 CWE-79
```

```text
CWE-79 is less than CWE-89
CWE-79 is equal to CWE-79
CWE-89 is greater than CWE-79
```

---

## 🔧 SDK API

### ExtractCWEIDs

```go
ids := cweskills.ExtractCWEIDs("Affected by CWE-79, CWE-89, and CWE-352")
// ids = ["CWE-79", "CWE-89", "CWE-352"]
```

大小写不敏感的正则匹配 `CWE-\d+`。

### ExtractFirstCWEID

```go
id := cweskills.ExtractFirstCWEID("See CWE-79 and CWE-89")
// id = "CWE-79"，无匹配返回空串
```

### CompareCWEIDs

```go
r, _ := cweskills.CompareCWEIDs("CWE-79", "CWE-89")  // r = -1（小于）
r, _ := cweskills.CompareCWEIDs("CWE-79", "CWE-79")  // r = 0（相等）
r, _ := cweskills.CompareCWEIDs("CWE-89", "CWE-79")  // r = 1（大于）
```

比较的是数字部分：79 < 89。

---

## 📝 示例

### 命令行

```bash
# 从安全公告抽取所有 CWE
cwe extract "$(cat advisory.txt)" -o json | jq '.ids'

# 比较两个 ID 的顺序
cwe compare CWE-79 CWE-89
```

### Go

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    report := "模块存在 XSS(CWE-79) 与 SQL 注入(CWE-89) 问题"
    for _, id := range cweskills.ExtractCWEIDs(report) {
        fmt.Println("发现:", id)
    }
}
```

---

## 🤖 AI 代理使用提示

- 用户粘贴一段漏洞描述时，AI 用 `cwe extract` 自动发现涉及的 CWE。
- 需要去重时，AI 可对提取结果用 `cwe format` 规范化后再去重。
- `compare` 适合排序、分组、判断两个 CWE 是否相同。

::: tip 典型用例
1. 漏洞报告解析：从安全公告抽取所有 CWE 引用。
2. 合规检查：比较 CWE ID 决定分组或顺序。
3. 数据规范化：从自由文本字段提取并规范化 CWE。
:::

---

## 📖 相关文档

- [技能 01 — 解析与验证](./01-cwe-id-parsing-validation)
- [CLI: extract](../cli/extract) · [compare](../cli/compare)
- [SDK: ExtractCWEIDs](../sdk/extract-cwe-ids) · [CompareCWEIDs](../sdk/compare-cwe-ids)
- [返回 Skills 总览](./)
