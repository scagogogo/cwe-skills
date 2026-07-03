# 🌳 cwe tree build

构建指定根节点的 CWE 层次树，并以缩进形式输出。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe tree build [ROOT-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `ROOT-ID` | 是 | 根节点 CWE ID（固定 1 个） |

## Flags

继承自 [`tree`](./tree) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe tree build CWE-1 --xml cwec_latest.xml
```

```text
CWE-1 - Location
  CWE-20 - ...
    CWE-79 - ...
    CWE-89 - ...
  ...
```

缩进反映层次深度。

### JSON 输出

```bash
cwe tree build CWE-1 --xml cwec_latest.xml -o json
```

```json
{
  "id": 1,
  "name": "Location",
  "depth": 0,
  "is_leaf": false,
  "children": [
    {
      "id": 20,
      "name": "...",
      "depth": 1,
      "is_leaf": false,
      "children": [ ... ]
    }
  ]
}
```

::: warning 根不存在
若无法构建以指定根的树（如根 ID 不存在），命令返回非零退出码并报错 `无法构建以 CWE-XXX 为根的树`。
:::

## 使用场景

- 可视化某弱点及其子树的整体结构。
- 整体浏览某分支下的层次组织。
- 为 [`tree path`](./tree-path)/[`tree leaves`](./tree-leaves) 选定根节点。

::: tip 配合 leaves
构建树后，可用 [`tree leaves`](./tree-leaves) 列出该树下的所有叶子弱点。
:::

## 下一步

- [tree forest](./tree-forest) — 所有 Pillar 根的森林。
- [tree view](./tree-view) — 基于视图构建树。
- [tree leaves](./tree-leaves) — 列出叶子。

## 相关文档

- [SDK BuildTree](../sdk/tree)
