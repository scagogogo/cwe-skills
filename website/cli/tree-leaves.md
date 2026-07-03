# 🌳 cwe tree leaves

列出指定根节点下的所有叶子弱点（无子级的节点）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe tree leaves [ROOT-ID] [flags]
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
cwe tree leaves CWE-1 --xml cwec_latest.xml
```

```text
CWE-1 的叶子节点 (N 项):
  CWE-79 - ...
  CWE-89 - ...
  ...
```

### JSON 输出

```bash
cwe tree leaves CWE-1 --xml cwec_latest.xml -o json | jq '.count'
```

```json
{
  "root": "CWE-1",
  "leaves": [
    { "id": 79, "name": "..." },
    { "id": 89, "name": "..." }
  ],
  "count": 2
}
```

::: warning 根不存在
若无法构建以指定根的树，命令返回非零退出码并报错 `无法构建以 CWE-XXX 为根的树`。
:::

## 使用场景

- 识别某子树下最具体的弱点（叶子），通常是需要直接修复的对象。
- 评估某分支的细化程度（叶子越多，分类越细）。
- 配合 [`tree build`](./tree-build) 先看结构，再用 `leaves` 提取末端集合。

::: tip 叶子 = 最具体
叶子节点在层次树中没有子级，对应 CWE 中最具体的弱点描述（通常是 Variant 或具体 Base）。
:::

## 下一步

- [tree build](./tree-build) — 构建完整树。
- [tree path](./tree-path) — 查找路径。
- [nav descendants](./nav-descendants) — 全部后代。

## 相关文档

- [SDK BuildTree](../sdk/tree)
- [SDK TreeNode.LeafNodes](../sdk/tree-node)
