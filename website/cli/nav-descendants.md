# 🧭 cwe nav descendants

查询指定 CWE 的所有后代弱点（递归向下，基于本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav descendants [CWE-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID` | 是 | 待查询的 CWE ID（固定 1 个） |

## Flags

继承自 [`nav`](./nav) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe nav descendants CWE-74 --xml cwec_latest.xml
```

```text
CWE-74 的 后代 (N 项):
  CWE-79 - ...
  CWE-89 - ...
  ...
```

### JSON 输出

```bash
cwe nav descendants CWE-74 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-74",
  "type": "后代",
  "results": [
    { "id": 79, "name": "..." },
    { "id": 89, "name": "..." }
  ],
  "count": 2
}
```

## 使用场景

- 获取某弱点下的完整子树。
- 评估某类别下的弱点规模。
- 配合 [`tree leaves`](./tree-leaves) 定位叶子弱点。

::: warning 结果规模
对上层弱点（如 Pillar）调用 `descendants` 可能返回大量条目，但本地查询无速率限制，速度优于在线 [`relations descendants`](./relations-descendants)。
:::

## 下一步

- [nav ancestors](./nav-ancestors)
- [nav children](./nav-children)
- [tree build](./tree-build) — 构建子树。
- [tree leaves](./tree-leaves) — 列出叶子。

## 相关文档

- [SDK Navigator.Descendants](../sdk/navigator)
