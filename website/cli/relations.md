# 🌐 cwe relations

通过 MITRE CWE REST API 查询 CWE 条目之间的关系。

<Badge type="info" text="在线"/> `relations` 是父命令，需指定子命令。

## 语法

```bash
cwe relations <子命令> [CWE-ID] [flags]
```

## 子命令

| 子命令 | 说明 | 是否支持 `--view-id` |
| --- | --- | --- |
| [`parents`](./relations-parents) | 父级弱点 | 是 |
| [`children`](./relations-children) | 子级弱点 | 是 |
| [`ancestors`](./relations-ancestors) | 所有祖先弱点 | 否 |
| [`descendants`](./relations-descendants) | 所有后代弱点 | 否 |

## Flags

`relations` 的 `PersistentFlags`，对所有子命令生效：

| Flag | 默认值 | 说明 |
| --- | --- | --- |
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |
| `--view-id` | `0` | 视图 ID（可选，仅对 `parents`/`children` 生效；为 0 时查询全局关系） |

## 示例

### 查询父级

```bash
cwe relations parents CWE-79
```

```text
CWE-79 的 父级弱点 (N 项):
  ChildOf -> CWE-74
  ...
```

### 在指定视图下查询子级

```bash
cwe relations children CWE-74 --view-id 1000
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
    { "nature": "ChildOf", "cwe_id": 74 }
  ],
  "count": 1
}
```

## 使用场景

- 在线探索弱点的层级关系，理解其归类路径。
- 在特定视图（如开发视图）下查询父子关系。
- 构建关系图谱用于分析。

::: tip 离线替代
[`nav`](./nav) 与 [`registry`](./registry) 提供更丰富的本地关系查询（含 siblings、peers、shortest-path 等），且无需联网。在线 `relations` 适合临时查询或无本地 XML 时使用。
:::

## 下一步

- [relations parents](./relations-parents)
- [relations children](./relations-children)
- [relations ancestors](./relations-ancestors)
- [relations descendants](./relations-descendants)
- [nav](./nav) — 更强的本地导航。

## 相关文档

- [SDK 父子关系](../sdk/api-parents-children)
- [SDK 祖先后代](../sdk/api-ancestors-descendants)
- [关系概念](../guide/concept-relationship)
