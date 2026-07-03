# 🌳 cwe tree path

查找从层次树根到指定 CWE 节点的路径。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe tree path [CWE-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID` | 是 | 目标 CWE ID（固定 1 个） |

## Flags

继承自 [`tree`](./tree) 的 `PersistentFlags`，以及 `path` 专属 flag：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |
| `--root` | | `0` | 根节点 ID（默认 `0` 表示自动查找：取目标节点的最后一位祖先作为根） |

## 示例

### 自动查找根

```bash
cwe tree path CWE-79 --xml cwec_latest.xml
```

```text
从 CWE-1 到 CWE-79 的路径 (3 步):
  1. CWE-1 - Location
  2. CWE-20 - ...
  3. CWE-79 - ...
```

未指定 `--root` 时，命令会自动用目标的最后一位祖先（Pillar）作为根。

### 指定根节点

```bash
cwe tree path CWE-79 --root 1 --xml cwec_latest.xml
```

### JSON 输出

```bash
cwe tree path CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "root": "CWE-1",
  "path": [1, 20, 79],
  "depth": 2
}
```

::: warning 节点不在树中
若指定的 `--root` 构建的树中不包含目标节点，命令报错 `在树中未找到 CWE-XXX`。
:::

## 使用场景

- 定位某弱点在层次树中的完整归属路径。
- 生成“根 → ... → 目标”的可视化路径。
- 确认某弱点归属于哪个顶层 Pillar。

::: tip 与 nav ancestors 配合
[`nav ancestors`](./nav-ancestors) 返回祖先列表；`tree path` 返回从根到目标的有序路径，且支持自定义根。`path` 的自动查找根正是基于 ancestors 的最后一位。
:::

## 下一步

- [tree build](./tree-build) — 构建完整树。
- [tree leaves](./tree-leaves) — 列出叶子。
- [nav ancestors](./nav-ancestors) — 祖先链。

## 相关文档

- [SDK BuildTree](../sdk/tree)
- [SDK TreeNode.Path](../sdk/tree-node)
