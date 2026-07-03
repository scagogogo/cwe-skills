# 🌳 cwe tree

构建和查询 CWE 层次树，支持从 XML 目录数据构建树并遍历。

<Badge type="tip" text="离线"/> 所有 `tree` 子命令都需要通过 `--xml` 指定 CWE XML 目录文件。

## 语法

```bash
cwe tree <子命令> [CWE-ID] [flags]
```

## 子命令

| 子命令 | 说明 |
| --- | --- |
| [`build`](./tree-build) | 构建指定根节点的 CWE 层次树 |
| [`forest`](./tree-forest) | 构建所有顶层（Pillar）节点的 CWE 森林 |
| [`view`](./tree-view) | 构建指定视图的 CWE 层次树 |
| [`path`](./tree-path) | 查找从根到指定 CWE 的路径 |
| [`leaves`](./tree-leaves) | 列出指定根节点下的所有叶子弱点 |

## Flags

`tree` 的 `PersistentFlags`，对所有子命令生效：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

此外，[`tree path`](./tree-path) 还有专属 flag `--root`。

## 示例

### 构建以 CWE-1 为根的树

```bash
cwe tree build CWE-1 --xml cwec_latest.xml
```

```text
CWE-1 - ...
  CWE-20 - ...
    CWE-79 - ...
    ...
```

### 查找从根到某节点的路径

```bash
cwe tree path CWE-79 --xml cwec_latest.xml
```

### 列出某根下的叶子

```bash
cwe tree leaves CWE-1 --xml cwec_latest.xml
```

## 使用场景

- 可视化 CWE 层次结构，理解弱点的分类树。
- 查找从根到具体弱点的归属路径。
- 识别某子树下的所有叶子（最具体）弱点。

::: tip 与 nav 的关系
[`nav`](./nav) 提供关系查询（返回列表）；`tree` 提供层次树构建与遍历（返回嵌套结构）。`tree` 适合整体结构展示，`nav` 适合单点关系探索。
:::

## 下一步

- [tree build](./tree-build)
- [tree forest](./tree-forest)
- [tree path](./tree-path)

## 相关文档

- [SDK 层次树](../sdk/tree)
- [SDK Navigator](../sdk/navigator)
- [SDK 构建索引](../sdk/build-indexes)
