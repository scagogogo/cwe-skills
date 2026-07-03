# 🌐 cwe relations children

通过 MITRE CWE REST API 查询指定 CWE 的子级弱点。

<Badge type="info" text="在线"/>

## 语法

```bash
cwe relations children [CWE-ID] [flags]
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
| `--view-id` | `0` | 视图 ID（可选，指定时查询该视图下的子级） |

## 示例

### text 输出

```bash
cwe relations children CWE-74
```

```text
CWE-74 的 子级弱点 (N 项):
  ParentOf -> CWE-79
  ParentOf -> CWE-89
  ...
```

指定视图：

```bash
cwe relations children CWE-74 --view-id 1000
```

### JSON 输出

```bash
cwe relations children CWE-74 -o json
```

```json
{
  "cwe_id": "CWE-74",
  "relation_type": "子级弱点",
  "relationships": [
    { "nature": "ParentOf", "cwe_id": 79 },
    { "nature": "ParentOf", "cwe_id": 89 }
  ],
  "count": 2
}
```

## 使用场景

- 查找某弱点直接细分出的子弱点，了解具体变体。
- 在特定视图下定位子级。
- 与 [`parents`](./relations-parents) 配合做双向遍历。

::: tip 与 descendants 区别
`children` 只返回直接子级；[`descendants`](./relations-descendants) 返回递归向下的全部后代。需要完整子树用 `descendants`。
:::

## 下一步

- [relations parents](./relations-parents)
- [relations descendants](./relations-descendants)
- [nav children](./nav-children) — 离线版。

## 相关文档

- [SDK GetChildren](../sdk/api-parents-children)
