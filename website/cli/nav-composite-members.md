# 🧭 `cwe nav composite-members`

<Badge type="info" text="离线" /> <Badge type="tip" text="需 --xml" /> <Badge type="info" text="nav 子命令" />

查询一个**复合弱点（Composite）**的所有成员弱点。复合弱点要求多个弱点**同时存在**才会成立（参见 [结构类型：Composite](../enums/structure)）。

::: tip 与 chain-members 的区别
- `chain-members`：链式（Chain）弱点的成员，弱点须**按顺序**可达
- `composite-members`：复合（Composite）弱点的成员，弱点须**同时存在**
:::

## 📖 语法

```bash
cwe nav composite-members [CWE-ID] [flags]
```

| 参数 | 说明 |
|------|------|
| `CWE-ID` | 复合弱点的 CWE ID（必填，1 个） |

| Flag | 说明 |
|------|------|
| `--xml <file>` | 本地 MITRE CWE XML 目录文件（必填） |
| `-o, --output` | 输出格式 `text`（默认）或 `json` |

## 🚀 示例

### 文本输出

```bash
cwe nav composite-members CWE-352 --xml cwec_latest.xml
```

```text
CWE-352 (Cross-Site Request Forgery) 的复合成员:
  - CWE-352 → ...
```

### JSON 输出

```bash
cwe nav composite-members CWE-352 --xml cwec_latest.xml -o json
```

```json
[
  { "id": 352, "name": "Cross-Site Request Forgery (CSRF)" }
]
```

## 🧠 对应 SDK API

`Navigator.CompositeMembers(id int) []*CWE` —— 详见 [SDK nav-chain-composite](../sdk/nav-chain-composite)。

```go
nav := cweskills.NewNavigator(reg)
members := nav.CompositeMembers(352)
```

## 🎯 使用场景

- 分析 CSRF（CWE-352）等多弱点耦合场景的组成
- 评估某复合弱点的修复复杂度（需同时消除所有成员弱点）

## 🔗 相关

- [nav chain-members](./nav-chain-members) — 链式成员
- [nav 总览](./nav) · [SDK 复合成员](../sdk/nav-chain-composite) · [结构类型 Composite](../enums/structure)
