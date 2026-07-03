# 🌐 cwe relations parents

通过 MITRE CWE REST API 查询指定 CWE 的父级弱点。

<Badge type="info" text="在线"/>

## 语法

```bash
cwe relations parents [CWE-ID] [flags]
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
| `--view-id` | `0` | 视图 ID（可选，指定时查询该视图下的父子关系） |

## 示例

### text 输出

```bash
cwe relations parents CWE-79
```

```text
CWE-79 的 父级弱点 (1 项):
  ChildOf -> CWE-74
```

指定视图：

```bash
cwe relations parents CWE-79 --view-id 1000
```

### JSON 输出

```bash
cwe relations parents CWE-79 -o json
```

```json
{
  "cwe_id": "CWE-79",
  "relation_type": "父级弱点",
  "relationships": [
    { "nature": "ChildOf", "cwe_id": 74, "view_id": 0 }
  ],
  "count": 1
}
```

## 使用场景

- 查找某弱点的直接上层分类，理解其归类。
- 在特定视图下定位父级。
- 与 [`children`](./relations-children) 配合做双向遍历。

::: tip 与 ancestors 区别
`parents` 只返回直接父级；[`ancestors`](./relations-ancestors) 返回递归向上的全部祖先。需要完整路径用 `ancestors`。
:::

## 下一步

- [relations children](./relations-children)
- [relations ancestors](./relations-ancestors)
- [nav parents](./nav-parents) — 离线版。

## 相关文档

- [SDK GetParents](../sdk/api-parents-children)
