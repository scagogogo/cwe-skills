# 🌐 cwe relations descendants

通过 MITRE CWE REST API 查询指定 CWE 的所有后代弱点（递归向下）。

<Badge type="info" text="在线"/>

## 语法

```bash
cwe relations descendants [CWE-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID` | 是 | 待查询的 CWE ID（固定 1 个） |

## Flags

继承自 [`relations`](./relations) 的 `PersistentFlags`：

| Flag | 默认值 | 说明 |
| --- | --- | --- |
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |

::: warning 不支持 --view-id
`--view-id` 仅对 `parents`/`children` 生效，对 `descendants` 无意义。
:::

## 示例

### text 输出

```bash
cwe relations descendants CWE-74
```

```text
CWE-74 的 后代弱点 (N 项):
  ParentOf -> CWE-79
  ParentOf -> CWE-89
  ...
```

### JSON 输出

```bash
cwe relations descendants CWE-74 -o json
```

```json
{
  "cwe_id": "CWE-74",
  "relation_type": "后代弱点",
  "relationships": [
    { "nature": "ParentOf", "cwe_id": 79 },
    { "nature": "ParentOf", "cwe_id": 89 }
  ],
  "count": 2
}
```

## 使用场景

- 获取某弱点下的完整子树，了解其所有细化变体。
- 评估某类别下的弱点规模。
- 与 [`tree leaves`](./tree-leaves) 配合定位叶子弱点。

::: warning 结果规模
对上层弱点（如 Pillar）调用 `descendants` 可能返回大量条目，注意 MITRE API 速率限制。如需完整子树分析，建议改用离线 [`nav descendants`](./nav-descendants) 或 [`tree build`](./tree-build)。
:::

## 下一步

- [relations ancestors](./relations-ancestors)
- [relations children](./relations-children)
- [nav descendants](./nav-descendants) — 离线版。
- [tree build](./tree-build) — 构建子树。

## 相关文档

- [SDK GetDescendants](../sdk/api-ancestors-descendants)
