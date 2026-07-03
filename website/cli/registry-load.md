# 🗃️ cwe registry load

加载 MITRE CWE XML 目录文件，构建索引并显示概要信息。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry load [flags]
```

## 参数

本命令无位置参数。

## Flags

继承自 [`registry`](./registry) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe registry load --xml cwec_latest.xml
```

```text
已加载CWE注册表:
  弱点:     1400+
  类别:     300+
  视图:     40+
  复合元素: N
  索引:     true
```

### JSON 输出

```bash
cwe registry load --xml cwec_latest.xml -o json
```

```json
{
  "weaknesses": 1400,
  "categories": 300,
  "views": 40,
  "compounds": 10,
  "indexed": true
}
```

## 使用场景

- 验证 XML 文件可被正确解析、索引构建成功。
- 在执行其他 registry/nav/tree 命令前确认数据可用。
- 对比不同版本 XML 的条目规模。

::: tip 索引构建
`load` 显示的 `索引: true` 表示注册表的查找索引已构建，这是后续 `get`/`contains`/关系查询快速响应的前提。
:::

## 下一步

- [registry get](./registry-get) — 获取条目详情。
- [registry contains](./registry-contains) — 检查存在性。
- [stats](./stats) — 更详细的分布统计。

## 相关文档

- [SDK Registry](../sdk/registry)
- [SDK 构建索引](../sdk/build-indexes)
