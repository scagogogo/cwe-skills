# 📚 cwe enum status

列出 CWE 条目状态（Status）的所有合法取值，如 `Stable`、`Draft`、`Deprecated` 等。

## 语法

```bash
cwe enum status [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum status
```

```text
状态 (Stable/Draft/Deprecated等) (4 项):
  - Stable
  - Draft
  - Deprecated
  - Obsolete
```

### JSON 输出

```bash
cwe enum status -o json
```

```json
[
  "Stable",
  "Draft",
  "Deprecated",
  "Obsolete"
]
```

## 取值含义

| 取值 | 含义 |
| --- | --- |
| `Stable` | 已稳定，内容成熟可广泛引用 |
| `Draft` | 草稿阶段，仍可能调整 |
| `Deprecated` | 已弃用，建议改用其他条目 |
| `Obsolete` | 已废弃，不再维护 |

::: warning 取值以实际输出为准
不同版本 XML/SDK 可能包含不同状态项，请以 `cwe enum status` 实际输出为准。
:::

## 使用场景

- 在 [`filter --status`](./filter) 前确认可选值。
- 筛选仅关注稳定弱点的条目（`--status Stable`）。
- 识别已弃用条目以便迁移引用。

## 下一步

- [enum abstraction](./enum-abstraction)
- [filter](./filter) — 按 status 过滤。

## 相关文档

- [SDK Status 枚举](../sdk/enum-status)
