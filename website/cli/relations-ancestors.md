# 🌐 cwe relations ancestors

通过 MITRE CWE REST API 查询指定 CWE 的所有祖先弱点（递归向上）。

<Badge type="info" text="在线"/>

## 语法

```bash
cwe relations ancestors [CWE-ID] [flags]
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
`--view-id` 仅对 `parents`/`children` 生效，对 `ancestors` 无意义，源码中未使用。
:::

## 示例

### text 输出

```bash
cwe relations ancestors CWE-79
```

```text
CWE-79 的 祖先弱点 (N 项):
  ChildOf -> CWE-74
  ChildOf -> CWE-...
  ...
```

### JSON 输出

```bash
cwe relations ancestors CWE-79 -o json
```

```json
{
  "cwe_id": "CWE-79",
  "relation_type": "祖先弱点",
  "relationships": [
    { "nature": "ChildOf", "cwe_id": 74 },
    { "nature": "ChildOf", "cwe_id": 285 }
  ],
  "count": 2
}
```

## 使用场景

- 获取从某弱点到根的完整祖先链，理解其在层次树中的位置。
- 判断两个弱点是否在同一祖先路径上。
- 构建 [`tree path`](./tree-path) 前的关系预查。

::: tip 与 parents 区别
`ancestors` 递归返回全部上层弱点；[`parents`](./relations-parents) 只返回直接父级。
:::

## 下一步

- [relations descendants](./relations-descendants)
- [relations parents](./relations-parents)
- [nav ancestors](./nav-ancestors) — 离线版。
- [tree path](./tree-path) — 查找从根到节点的路径。

## 相关文档

- [SDK GetAncestors](../sdk/api-ancestors-descendants)
