# 🧭 cwe nav peers

查询指定 CWE 的对等弱点（基于显式 `PeerOf` 关系，本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav peers [CWE-ID] [flags]
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
cwe nav peers CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 对等 (N 项):
  CWE-XXX - ...
```

### JSON 输出

```bash
cwe nav peers CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "对等",
  "results": [
    { "id": 123, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 查找 MITRE 显式声明与某弱点“对等”的弱点。
- 扩展弱点集合，发现语义相近的条目。
- 关系类型含义见 [`enum relationship`](./enum-relationship)。

::: tip 与 siblings 的区别
[`siblings`](./nav-siblings) 基于“同一父级”推断；`peers` 基于 CWE 显式声明的 `PeerOf` 关系。前者是结构性推导，后者是声明性关系。
:::

## 下一步

- [nav siblings](./nav-siblings)
- [nav can-also-be](./nav-can-also-be)
- [enum relationship](./enum-relationship)

## 相关文档

- [SDK Navigator.Peers](../sdk/navigator)
