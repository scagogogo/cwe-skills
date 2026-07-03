# 🔍 cwe search

从本地 XML 目录搜索 CWE 条目。

<Badge type="tip" text="离线"/> 需通过 `--xml` 指定 MITRE CWE XML 目录文件。

## 语法

```bash
cwe search [flags]
```

## 参数

本命令无位置参数，所有筛选条件通过 flags 提供。

## Flags

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |
| `--keyword` | `-k` | | 按关键字搜索 |
| `--abstraction` | `-a` | | 按抽象层级搜索（Pillar/Class/Base/Variant） |
| `--status` | `-s` | | 按状态搜索（Stable/Draft/Deprecated） |
| `--likelihood` | `-l` | | 按利用可能性搜索（High/Medium/Low） |
| `--structure` | `-t` | | 按结构类型搜索（Simple/Chain/Composite） |
| `--scope` | | | 按后果范围搜索（Confidentiality/Integrity/Availability） |
| `--top-level` | | `false` | 只显示顶层（Pillar）弱点 |
| `--base-weaknesses` | | `false` | 只显示基础（Base）弱点 |
| `--chains` | | `false` | 只显示链式弱点 |
| `--composites` | | `false` | 只显示复合弱点 |
| `--sort` | | | 排序字段（id/name/abstraction） |
| `--group-by` | | | 分组字段（abstraction/status/likelihood） |
| `--dedup` | | `false` | 去重 |

::: warning 单一筛选
`search` 的各筛选 flag 之间为互斥分支（源码按 `if/else if` 顺序判断），一次只生效一个筛选条件。需要多条件组合请用 [`filter`](./filter)。
:::

## 示例

### 按关键字搜索并按名称排序

```bash
cwe search --xml cwec_latest.xml --keyword Injection --sort name
```

```text
找到 N 个CWE条目:

  CWE-79 - ... [Base, Stable]
  CWE-89 - ... [Base, Stable]
  ...
```

### 只看顶层 Pillar 弱点

```bash
cwe search --xml cwec_latest.xml --top-level
```

### 按抽象层级分组

```bash
cwe search --xml cwec_latest.xml --group-by abstraction
```

### JSON 输出

```bash
cwe search --xml cwec_latest.xml --keyword Injection -o json | jq 'length'
```

## 使用场景

- 在离线环境中按单一维度快速检索弱点。
- 列出所有顶层/基础/链式/复合弱点。
- 按抽象层级或状态分组浏览。

::: tip 多条件用 filter
`search` 一次只支持一个筛选条件，[`filter`](./filter) 支持多条件 AND 组合。如需“Base + Stable + High 利用可能性”同时筛选，请用 [`filter`](./filter)。
:::

## 下一步

- [filter](./filter) — 多条件组合过滤。
- [stats](./stats) — 数据分布统计。
- [enum abstraction](./enum-abstraction) — 确认可选值。

## 相关文档

- [SDK 搜索过滤](../sdk/search)
- [SDK XML 解析器](../sdk/xml-parser)
