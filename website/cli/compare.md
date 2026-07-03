# 🆔 cwe compare

比较两个 CWE ID 的大小关系（按数字编号）。

## 语法

```bash
cwe compare <CWE-ID1> <CWE-ID2> [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID1` | 是 | 第一个 CWE ID |
| `CWE-ID2` | 是 | 第二个 CWE ID |

::: warning 参数数量
本命令固定接收 2 个参数，多于或少于都会报错。
:::

## Flags

本命令无专属 flags，仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe compare CWE-79 CWE-89
```

```text
CWE-79 is less than CWE-89
```

```bash
cwe compare 89 79
```

```text
CWE-89 is greater than CWE-79
```

```bash
cwe compare CWE-79 CWE-79
```

```text
CWE-79 is equal to CWE-79
```

### JSON 输出

```bash
cwe compare CWE-79 CWE-89 -o json
```

```json
{
  "id1": "CWE-79",
  "id2": "CWE-89",
  "comparison": "less than",
  "result": -10
}
```

`result` 为两 ID 数值之差：负数表示前者较小，正数表示较大，0 表示相等。

## 使用场景

- 对 CWE ID 集合排序时作为比较函数。
- 判断两个引用是否指向同一弱点。
- 需要按编号范围筛选时的辅助判断。

::: tip 比较的是数字
`compare` 先将输入解析为整数再比较，因此 `CWE-79` 与 `79` 等价。若你已有纯整数，可直接用 [compare-int](./compare-int)。
:::

## 下一步

- [compare-int](./compare-int) — 直接比较两个整数。
- [parse](./parse) — 了解 ID 解析规则。

## 相关文档

- [SDK CompareCWEIDs](../sdk/compare-cwe-ids)
