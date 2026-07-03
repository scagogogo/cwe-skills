# 📚 cwe enum consequence-scope

列出 CWE 后果范围（ConsequenceScope）的所有合法取值。

::: warning 命令名
本枚举对应的子命令名为 `scope`，而非 `consequence-scope`。即实际执行命令为 `cwe enum scope`。
:::

## 语法

```bash
cwe enum scope [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum scope
```

```text
后果范围 (Confidentiality/Integrity/Availability) (3 项):
  - Confidentiality
  - Integrity
  - Availability
```

### JSON 输出

```bash
cwe enum scope -o json
```

```json
[
  "Confidentiality",
  "Integrity",
  "Availability"
]
```

## 取值含义

| 取值 | 含义 |
| --- | --- |
| `Confidentiality` | 机密性，信息泄露 |
| `Integrity` | 完整性，数据被篡改 |
| `Availability` | 可用性，服务被中断 |

即经典的 CIA 三元组。

## 使用场景

- 在 [`filter --scope`](./filter) 或 [`search --scope`](./search) 前确认可选值。
- 按安全属性筛选弱点，例如仅关注机密性影响。
- 合规分析中映射弱点到 CIA 维度。

::: tip 配合 filter
筛选影响机密性的弱点：`cwe filter --xml file.xml --scope Confidentiality`
:::

## 下一步

- [enum consequence-impact](./enum-consequence-impact)
- [filter](./filter) — 按 scope 过滤。

## 相关文档

- [SDK ConsequenceScope 枚举](../sdk/enum-consequence-scope)
- [后果概念](../guide/concept-consequence)
