# 🔍 cwe stats

从本地 XML 目录统计 CWE 数据分布。

<Badge type="tip" text="离线"/> 需通过 `--xml` 指定 MITRE CWE XML 目录文件。

## 语法

```bash
cwe stats [flags]
```

## 参数

本命令无位置参数。

## Flags

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe stats --xml cwec_latest.xml
```

```text
CWE数据统计:

  总条目数:     1400+
  类别数:       300+
  视图数:       40+

抽象层级分布:
  Pillar: 10+
  Class: 70+
  Base: 800+
  Variant: 500+

状态分布:
  Stable: ...
  Draft: ...

利用可能性分布:
  High: ...
  Medium: ...
  Low: ...
```

### JSON 输出

```bash
cwe stats --xml cwec_latest.xml -o json | jq '.total_count'
```

```json
{
  "total_count": 1400,
  "category_count": 300,
  "view_count": 40,
  "by_abstraction": { "Base": 800, "Variant": 500 },
  "by_status": { "Stable": 1300 },
  "by_likelihood": { "High": 200 }
}
```

## 使用场景

- 了解本地 XML 数据集的整体规模与分布。
- 决定后续 [`filter`](./filter)/[`search`](./search) 的筛选维度。
- 数据质量评估：观察 Draft/Deprecated 占比。

::: tip 与 enum 配合
`stats` 输出的分布维度（抽象层级、状态、利用可能性）的取值含义可对照 [enum abstraction](./enum-abstraction)、[enum status](./enum-status)、[enum likelihood](./enum-likelihood)。
:::

## 下一步

- [search](./search) — 按维度检索。
- [filter](./filter) — 多条件过滤。
- [registry load](./registry-load) — 加载并查看注册表概要。

## 相关文档

- [SDK ComputeStatistics](../sdk/compute-statistics)
- [SDK XML 解析器](../sdk/xml-parser)
