---
title: 技能 01 — CWE ID 解析与验证
outline: [2, 3]
---

# 🆔 技能 01 — CWE ID 解析与验证

解析、验证、格式化 CWE ID。这是最基础的技能——**所有 CWE 操作都从拿到一个合法的 CWE ID 开始**。

<Badge type="tip" text="基础技能"/>
<Badge type="info" text="无需网络"/>

CWE ID 遵循 `CWE-NNN` 格式（NNN 为正整数）。SDK 接受多种写法（`79`、`CWE-79`、`cwe-79`），统一规范化为 `CWE-NNN`。

---

## 🎯 技能目标

- 从任意输入解析出 CWE ID 的数字部分
- 验证输入是否为合法 CWE ID
- 把多种写法格式化为标准 `CWE-NNN`

---

## 💻 CLI 命令

### parse — 解析

```bash
cwe parse CWE-79 89 cwe-352
```

```text
CWE-79 -> CWE-79 (ID: 79)
89 -> CWE-89 (ID: 89)
cwe-352 -> CWE-352 (ID: 352)
```

JSON 输出，无效输入会被标注但不致命：

```bash
cwe parse CWE-79 abc -o json
```

```json
[
  {"input": "CWE-79", "id": 79, "format": "CWE-79", "valid": true},
  {"input": "abc", "id": 0, "valid": false, "error": "cwe: [INVALID_CWE_ID] CWE ID格式无效: 输入值: \"abc\""}
]
```

### validate — 验证

```bash
cwe validate CWE-79 CWE-89 abc
```

```text
CWE-79 ✓ 有效
CWE-89 ✓ 有效
abc ✗ 无效
```

全部有效退出码 `0`，有无效项退出码 `1`。

### format — 格式化

```bash
cwe format 79 cwe-89 CWE-352
```

```text
CWE-79
CWE-89
CWE-352
```

---

## 🔧 SDK API

### ParseCWEID

```go
id, err := cweskills.ParseCWEID("CWE-79")  // id = 79
id, err := cweskills.ParseCWEID("abc")     // id = 0, err = InvalidCWEIDError
```

大小写不敏感，可省略前缀；`CWE-0`、`CWE--1` 被拒绝。

### FormatCWEID / FormatCWEIDFromInt

```go
s, err := cweskills.FormatCWEID("79")        // "CWE-79"
s := cweskills.FormatCWEIDFromInt(79)        // "CWE-79"
```

### IsCWEID / ValidateCWEID

```go
cweskills.IsCWEID("CWE-89")        // true — 快速布尔判断
err := cweskills.ValidateCWEID("abc")     // InvalidCWEIDError
err := cweskills.ValidateCWEID("CWE-79")  // nil
```

`IsCWEID` 是布尔快查；`ValidateCWEID` 返回带详情的结构化错误。

::: details 错误类型
`InvalidCWEIDError`：输入不匹配 CWE ID 模式（空、无数字、负数）。
:::

---

## 📝 示例

### 命令行

```bash
# 从漏洞报告字段批量解析
cwe parse CWE-79 CWE-89 CWE-352 -o json | jq '.[] | .format'

# 把用户输入规范化
cwe format 79 89 352
```

### Go

```go
package main

import (
    "fmt"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    inputs := []string{"CWE-79", "89", "cwe-352", "abc"}
    for _, in := range inputs {
        if id, err := cweskills.ParseCWEID(in); err == nil {
            fmt.Printf("%s -> %s\n", in, cweskills.FormatCWEIDFromInt(id))
        } else {
            fmt.Printf("%s 无效: %v\n", in, err)
        }
    }
}
```

---

## 🤖 AI 代理使用提示

- 用户提到 CWE ID 时，先用 `cwe parse` 规范化，再去查详情。
- 批量解析用 `-o json`，AI 解析无歧义。
- 用户输入不确定是否合法时，用 `cwe validate` 判断退出码。

---

## 📖 相关文档

- [技能 02 — 提取与比较](./02-cwe-id-extraction-comparison)
- [CLI: parse](../cli/parse) · [validate](../cli/validate) · [format](../cli/format)
- [SDK: ParseCWEID](../sdk/parse-cwe-id) · [FormatCWEID](../sdk/format-cwe-id)
- [返回 Skills 总览](./)
