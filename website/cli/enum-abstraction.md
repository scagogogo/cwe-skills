# 📚 cwe enum abstraction

列出 CWE 抽象层级（Abstraction）的所有合法取值。

CWE 通过抽象层级描述弱点的泛化程度，从最抽象到最具体通常为 Pillar → Class → Base → Variant。

## 语法

```bash
cwe enum abstraction [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum abstraction
```

```text
抽象层级 (Class/Base/Variant/Pillar) (4 项):
  - Pillar
  - Class
  - Base
  - Variant
```

### JSON 输出

```bash
cwe enum abstraction -o json
```

```json
[
  "Pillar",
  "Class",
  "Base",
  "Variant"
]
```

## 取值含义

| 取值 | 含义 |
| --- | --- |
| `Pillar` | 最顶层抽象，代表一类弱点的根本主题 |
| `Class` | 比 Pillar 更具体的类别 |
| `Base` | 可独立描述、可被利用的具体弱点，最为常用 |
| `Variant` | Base 之下的特定变体 |

## 使用场景

- 在 [`filter --abstraction`](./filter) 前确认可选值。
- 生成报告时按抽象层级分组统计。
- 教学/文档场景介绍 CWE 层次模型。

::: tip 配合 stats
[`stats`](./stats) 会输出按抽象层级的分布统计，可与本命令配合理解各取值含义。
:::

## 下一步

- [enum status](./enum-status)
- [filter](./filter) — 按 abstraction 过滤。
- [stats](./stats) — 查看抽象层级分布。

## 相关文档

- [SDK Abstraction 枚举](../sdk/enum-abstraction)
- [抽象层级概念](../guide/concept-abstraction)
