# 🗃️ cwe registry parents / children

查询本地注册表中的直接父级与直接子级关系。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry parents [CWE-ID] [flags]
cwe registry children [CWE-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID` | 是 | 待查询的 CWE ID（固定 1 个） |

## Flags

继承自 [`registry`](./registry) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### parents

```bash
cwe registry parents CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 父级弱点 (1 项):
  CWE-74 - Injection
```

### children

```bash
cwe registry children CWE-74 --xml cwec_latest.xml
```

```text
CWE-74 的 子级弱点 (2 项):
  CWE-79 - ...
  CWE-89 - ...
```

### JSON 输出

```bash
cwe registry children CWE-74 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-74",
  "type": "子级弱点",
  "results": [
    { "id": 79, "name": "..." },
    { "id": 89, "name": "..." }
  ],
  "count": 2
}
```

## 使用场景

- 离线快速查看某弱点的直接上层/下层，无需联网。
- 替代在线 [`relations parents`](./relations-parents)/[`relations children`](./relations-children)。
- 与 [`nav`](./nav) 配合做更深层的关系导航。

::: tip 与 nav 的区别
`registry parents/children` 只返回直接父子 ID 与名称；[`nav parents`](./nav-parents)/[`nav children`](./nav-children) 返回完整 CWE 对象，且 `nav` 还支持 siblings/peers/shortest-path 等更丰富查询。
:::

## 下一步

- [registry ancestors / descendants](./registry-anc-desc) — 递归查询。
- [nav parents](./nav-parents) / [nav children](./nav-children) — 更丰富的本地导航。
- [relations](./relations) — 在线版本。

## 相关文档

- [SDK Registry.GetParentIDs](../sdk/registry)
- [SDK Registry.GetChildIDs](../sdk/registry)
- [关系概念](../guide/concept-relationship)
