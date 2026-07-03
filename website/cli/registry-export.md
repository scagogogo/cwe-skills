# 🗃️ cwe registry export

将加载的 XML 注册表导出为 JSON 或 CSV 格式。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry export [flags]
```

## 参数

本命令无位置参数。

## Flags

继承自 [`registry`](./registry) 的 `PersistentFlags`，以及 export 专属 flags：

| Flag | 默认值 | 说明 |
| --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |
| `--format` | `json` | 导出格式（`json` 或 `csv`） |
| `--output-file` | | 输出文件路径（默认输出到 stdout） |

::: danger --output-file 当前为占位实现
经核查源码 `cmd/cwe/registry.go`，当指定 `--output-file` 时调用的 `writeFile` 函数目前是占位实现——它并不真正写盘，而是返回一个错误信息（形如 `写入文件功能: <path> (<n> bytes)`）。因此当前版本只能将导出内容输出到 stdout，再用 shell 重定向落盘：

```bash
cwe registry export --xml cwec_latest.xml --format json > cwe.json
```

期待后续版本完善 `--output-file` 的真实写盘逻辑。
:::

## 示例

### 导出 JSON 到 stdout

```bash
cwe registry export --xml cwec_latest.xml --format json
```

输出为完整的 JSON 数组（缩进格式），可重定向到文件：

```bash
cwe registry export --xml cwec_latest.xml --format json > cwe.json
```

### 导出 CSV

```bash
cwe registry export --xml cwec_latest.xml --format csv > cwe.csv
```

### 配合 jq 提取子集

```bash
cwe registry export --xml cwec_latest.xml --format json | jq '.[] | select(.abstraction=="Base") | {id,name}'
```

## 使用场景

- 将 XML 数据转换为 JSON/CSV，便于其他工具消费。
- 一次性转换后供脚本、数据库导入或报表工具使用。
- 生成数据快照用于版本对比。

::: warning 格式仅支持 json/csv
`--format` 仅接受 `json` 与 `csv`，其他值会报错 `不支持的导出格式`。
:::

## 下一步

- [registry load](./registry-load) — 加载概要。
- [registry get](./registry-get) — 单条查询。
- [search](./search) / [filter](./filter) — 导出前的筛选。

## 相关文档

- [SDK Registry.ExportJSON](../sdk/registry)
- [SDK Registry.ExportCSV](../sdk/registry)
