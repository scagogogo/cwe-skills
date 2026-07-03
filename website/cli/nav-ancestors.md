# 🧭 cwe nav ancestors

查询指定 CWE 的所有祖先弱点（递归向上，基于本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav ancestors [CWE-ID] [flags]
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
cwe nav ancestors CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 祖先 (N 项):
  CWE-74 - Injection
  CWE-285 - ...
  ...
```

### JSON 输出

```bash
cwe nav ancestors CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "祖先",
  "results": [
    { "id": 74, "name": "Injection" },
    { "id": 285, "name": "..." }
  ],
  "count": 2
}
```

## 使用场景

- 获取从某弱点到根的完整祖先链。
- 判断两个弱点是否在同一祖先路径上（配合 [`is-ancestor`](./nav-is-ancestor)）。
- 为 [`tree path`](./tree-path) 定位根节点（`tree path` 在未指定 `--root` 时会自动用最后一位祖先作为根）。

::: tip 与 nav parents 区别
`ancestors` 递归返回全部上层；[`parents`](./nav-parents) 只返回直接父级。
:::

## 下一步

- [nav descendants](./nav-descendants)
- [nav parents](./nav-parents)
- [nav is-ancestor](./nav-is-ancestor)
- [tree path](./tree-path)

## 相关文档

- [SDK Navigator.Ancestors](../sdk/navigator)
