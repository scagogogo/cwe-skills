# 🗃️ cwe registry list-views

列出本地注册表中的所有 CWE 视图（View）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry list-views [flags]
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
cwe registry list-views --xml cwec_latest.xml
```

```text
视图 (40 项):
  CWE-1000 - ... [Graph]
  CWE-699 - ... [Slice]
  ...
```

### JSON 输出

```bash
cwe registry list-views --xml cwec_latest.xml -o json | jq '.[].id'
```

```json
[
  { "id": 1000, "name": "...", "type": "Graph" },
  { "id": 699, "name": "...", "type": "Slice" }
]
```

## 使用场景

- 浏览可用的视图清单，挑选目标视图用于 [`tree view`](./tree-view)。
- 查找特定视图 ID 用于 [`relations --view-id`](./relations)。
- 与 [`show view`](./show-view) 配合：先列出视图，再查看详情。

::: tip 视图类型
输出中 `type` 字段取值见 [`enum viewtype`](./enum-view-type)：`Graph` 或 `Slice`。
:::

## 下一步

- [registry list-categories](./registry-list-categories) — 列出类别。
- [tree view](./tree-view) — 基于视图构建树。
- [show view](./show-view) — 在线查看视图详情。

## 相关文档

- [SDK Registry.GetAllViews](../sdk/registry)
- [视图概念](../guide/concept-view)
