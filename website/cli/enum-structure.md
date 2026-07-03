# 📚 cwe enum structure

列出 CWE 结构类型（Structure）的所有合法取值：`Simple`、`Chain`、`Composite`。

## 语法

```bash
cwe enum structure [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum structure
```

```text
结构类型 (Simple/Chain/Composite) (3 项):
  - Simple
  - Chain
  - Composite
```

### JSON 输出

```bash
cwe enum structure -o json
```

```json
[
  "Simple",
  "Chain",
  "Composite"
]
```

## 取值含义

| 取值 | 含义 |
| --- | --- |
| `Simple` | 简单弱点，单一独立描述 |
| `Chain` | 链式弱点，由多个弱点按顺序组合而成 |
| `Composite` | 复合弱点，由多个弱点同时存在而构成 |

## 使用场景

- 在 [`filter --structure`](./filter) 或 [`search --structure`](./search) 前确认可选值。
- 识别链式/复合弱点，使用 [`nav chain-members`](./nav-chain-members) 查看其成员。
- 使用 [`search --chains`](./search) / `--composites` 快速筛选链式/复合弱点。

::: tip 与 nav 配合
对 `Chain` 结构的弱点，[`nav chain-members`](./nav-chain-members) 可列出其成员；对 `Composite`，可用 `nav composite-members`。
:::

## 下一步

- [enum abstraction](./enum-abstraction)
- [nav chain-members](./nav-chain-members)
- [search](./search) — `--chains` / `--composites` 标志。

## 相关文档

- [SDK Structure 枚举](../sdk/enum-structure)
- [结构概念](../guide/concept-structure)
