# 🗃️ cwe registry ancestors / descendants

查询本地注册表中的所有祖先与所有后代关系（递归）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry ancestors [CWE-ID] [flags]
cwe registry descendants [CWE-ID] [flags]
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

### ancestors

```bash
cwe registry ancestors CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 祖先弱点 (N 项):
  CWE-74 - Injection
  CWE-285 - ...
  ...
```

### descendants

```bash
cwe registry descendants CWE-74 --xml cwec_latest.xml
```

```text
CWE-74 的 后代弱点 (N 项):
  CWE-79 - ...
  CWE-89 - ...
  ...
```

### JSON 输出

```bash
cwe registry ancestors CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "祖先弱点",
  "results": [
    { "id": 74, "name": "Injection" },
    { "id": 285, "name": "..." }
  ],
  "count": 2
}
```

## 使用场景

- 获取从某弱点到根的完整祖先链，或某弱点下的完整子树。
- 替代在线 [`relations ancestors`](./relations-ancestors)/[`relations descendants`](./relations-descendants)，无需联网且无速率限制。
- 为 [`tree path`](./tree-path)/[`tree build`](./tree-build) 提供关系预查。

::: warning 结果规模
对上层弱点调用 `descendants` 可能返回大量条目，但因为是本地查询，无 API 速率限制，速度优于在线版本。
:::

## 下一步

- [registry parents / children](./registry-relations) — 直接父子关系。
- [nav ancestors](./nav-ancestors) / [nav descendants](./nav-descendants) — 返回完整对象的版本。
- [tree build](./tree-build) — 构建子树。

## 相关文档

- [SDK Registry.GetAncestorIDs](../sdk/registry)
- [SDK Registry.GetDescendantIDs](../sdk/registry)
- [关系概念](../guide/concept-relationship)
