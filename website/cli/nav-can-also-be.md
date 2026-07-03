# 🧭 cwe nav can-also-be

查询此弱点也可以是的弱点（`CanAlsoBe` 关系，本地 XML 数据）。

表示同一弱点在不同视角下也可表现为另一个弱点。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav can-also-be [CWE-ID] [flags]
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
cwe nav can-also-be CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 也可以是 (N 项):
  CWE-XXX - ...
```

### JSON 输出

```bash
cwe nav can-also-be CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "也可以是",
  "results": [
    { "id": 123, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 识别同一弱点在不同视角下的等价表示，避免遗漏修复。
- 在去重/合并弱点清单时识别可互相替代的条目。
- 关系类型含义见 [`enum relationship`](./enum-relationship)。

::: tip 与 peers 的区别
`peers` 是声明性对等关系；`can-also-be` 表示“同一弱点也可表现为另一弱点”，语义上更偏向等价性表达。
:::

## 下一步

- [nav peers](./nav-peers)
- [nav chain-members](./nav-chain-members)
- [enum relationship](./enum-relationship)

## 相关文档

- [SDK Navigator.CanAlsoBe](../sdk/navigator)
