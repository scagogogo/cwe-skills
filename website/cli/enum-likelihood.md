# 📚 cwe enum likelihood

列出 CWE 利用可能性（LikelihoodOfExploit）的所有合法取值：`High`、`Medium`、`Low`。

## 语法

```bash
cwe enum likelihood [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum likelihood
```

```text
利用可能性 (High/Medium/Low) (3 项):
  - High
  - Medium
  - Low
```

### JSON 输出

```bash
cwe enum likelihood -o json
```

```json
[
  "High",
  "Medium",
  "Low"
]
```

## 取值含义

| 取值 | 含义 |
| --- | --- |
| `High` | 被利用可能性高，需优先缓解 |
| `Medium` | 中等可能性 |
| `Low` | 较低可能性 |

## 使用场景

- 在 [`filter --likelihood`](./filter) 或 [`search --likelihood`](./search) 前确认可选值。
- 按“高利用可能性 + 高影响”组合筛选需优先处理的弱点。
- 风险评估时查阅取值集合。

::: tip 配合 filter
筛选高利用可能性的稳定基础弱点：`cwe filter --xml file.xml --abstraction Base --status Stable --likelihood High`
:::

## 下一步

- [enum consequence-impact](./enum-consequence-impact) — 配合使用评估风险。
- [filter](./filter) — 按 likelihood 过滤。

## 相关文档

- [SDK LikelihoodOfExploit 枚举](../sdk/enum-likelihood)
