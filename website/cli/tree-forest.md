# 🌳 cwe tree forest

构建所有顶层（Pillar）节点的 CWE 森林，逐棵输出。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe tree forest [flags]
```

## 参数

本命令无位置参数。

## Flags

继承自 [`tree`](./tree) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe tree forest --xml cwec_latest.xml
```

```text
CWE森林 (N 棵树):

--- 树 1 ---
CWE-XXX - ...
  CWE-YYY - ...
  ...

--- 树 2 ---
CWE-ZZZ - ...
...
```

### JSON 输出

```bash
cwe tree forest --xml cwec_latest.xml -o json | jq '.count'
```

```json
{
  "count": 9,
  "trees": [
    { "id": ..., "name": "...", "children": [ ... ] }
  ]
}
```

## 使用场景

- 浏览 CWE 顶层 Pillar 弱点及其各自子树，掌握全局结构。
- 选取感兴趣的 Pillar 作为 [`tree build`](./tree-build) 的根。
- 数据规模评估：查看森林包含多少棵顶层树。

::: tip 与 search --top-level 配合
[`search --top-level`](./search) 仅列出顶层 Pillar 的 ID 与名称；`tree forest` 则展开每棵子树，信息更完整但输出更大。
:::

## 下一步

- [tree build](./tree-build) — 构建单棵树。
- [tree view](./tree-view) — 视图树。
- [search](./search) — `--top-level` 标志。

## 相关文档

- [SDK BuildForest](../sdk/tree)
