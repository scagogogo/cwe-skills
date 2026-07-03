# 🔍 cwe filter

使用多个条件组合过滤 CWE 条目，所有条件之间为 **AND** 关系。

<Badge type="tip" text="离线"/> 需通过 `--xml` 指定 MITRE CWE XML 目录文件。

## 语法

```bash
cwe filter [flags]
```

## 参数

本命令无位置参数，所有条件通过 flags 提供。

## Flags

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |
| `--keyword` | `-k` | | 按关键字搜索 |
| `--abstraction` | `-a` | | 按抽象层级搜索 |
| `--status` | `-s` | | 按状态搜索 |
| `--likelihood` | `-l` | | 按利用可能性搜索 |
| `--structure` | `-t` | | 按结构类型搜索 |
| `--scope` | | | 按后果范围搜索 |
| `--sort` | | | 排序字段（id/name/abstraction） |
| `--group-by` | | | 分组字段（abstraction/status/likelihood） |

::: tip AND 关系
`filter` 同时指定多个条件时，结果必须同时满足全部条件（AND）。这与 [`search`](./search) 的互斥分支不同。
:::

## 示例

### 多条件组合

```bash
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable --likelihood High
```

```text
过滤结果 (N 项):

  CWE-79 - ... [Base, Stable]
  ...
```

### 关键字 + 范围

```bash
cwe filter --xml cwec_latest.xml --keyword Injection --scope Confidentiality
```

### 按名称排序

```bash
cwe filter --xml cwec_latest.xml --abstraction Base --sort name
```

### JSON 输出

```bash
cwe filter --xml cwec_latest.xml --abstraction Base --status Stable -o json | jq '.[].id'
```

## 使用场景

- 精细化筛选，例如“稳定的、高利用可能性的基础弱点”。
- 风险评估：按 `--likelihood High --scope Confidentiality` 定位高风险弱点。
- 报告生成：组合多维度筛选输出特定子集。

::: warning 枚举值需合法
`--abstraction`/`--status`/`--likelihood`/`--structure`/`--scope` 的取值必须为合法枚举值，可用 [`enum`](./enum) 系列命令查阅。非法值会报错。
:::

## 下一步

- [search](./search) — 单一条件快速检索。
- [stats](./stats) — 查看数据分布。
- [enum abstraction](./enum-abstraction) — 确认可选值。

## 相关文档

- [SDK Filter](../sdk/search)
- [SDK FilterOption](../sdk/filter-option)
