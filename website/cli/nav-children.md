# 🧭 cwe nav children

查询指定 CWE 的子级弱点（基于本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav children [CWE-ID] [flags]
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
cwe nav children CWE-74 --xml cwec_latest.xml
```

```text
CWE-74 的 子级 (2 项):
  CWE-79 - ...
  CWE-89 - ...
```

### JSON 输出

```bash
cwe nav children CWE-74 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-74",
  "type": "子级",
  "results": [
    { "id": 79, "name": "..." },
    { "id": 89, "name": "..." }
  ],
  "count": 2
}
```

## 使用场景

- 离线查找某弱点的直接细分变体。
- 与 [`nav parents`](./nav-parents) 配合做双向遍历。
- 在 [`tree build`](./tree-build) 前预览子级规模。

## 下一步

- [nav parents](./nav-parents)
- [nav descendants](./nav-descendants)
- [registry children](./registry-relations)

## 相关文档

- [SDK Navigator.Children](../sdk/navigator)
