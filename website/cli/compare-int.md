# 🆔 cwe compare-int

比较两个整数形式的 CWE ID 大小。

与 [compare](./compare) 的区别：`compare-int` 直接接收纯整数，不解析 `CWE-` 前缀，性能略好且输入更直观。

## 语法

```bash
cwe compare-int <int1> <int2> [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `int1` | 是 | 第一个整数 |
| `int2` | 是 | 第二个整数 |

::: danger 仅接受整数
参数必须为纯数字，`CWE-79` 这类带前缀的写法会被当作无效整数报错。需要解析前缀请用 [compare](./compare)。
:::

## Flags

本命令无专属 flags，仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe compare-int 79 89
```

```text
CWE-79 is less than CWE-89
```

```bash
cwe compare-int 89 79
```

```text
CWE-89 is greater than CWE-79
```

```bash
cwe compare-int 79 79
```

```text
CWE-79 is equal to CWE-79
```

### JSON 输出

```bash
cwe compare-int 79 89 -o json
```

```json
{
  "id1": "CWE-79",
  "id2": "CWE-89",
  "comparison": "less than",
  "result": -10
}
```

输出中的 `id1`/`id2` 已被格式化为标准 `CWE-XXX` 形式，便于直接引用。

## 使用场景

- 已知纯编号时的高效比较，避免前缀解析开销。
- 在脚本中按数字范围筛选 CWE（如“编号大于 1000 的弱点”）。
- 与 `compare` 互补：前者用于原始数字，后者用于带前缀写法。

## 下一步

- [compare](./compare) — 比较带 `CWE-` 前缀的 ID。
- [format](./format) — 将整数格式化为标准形式。

## 相关文档

- [SDK CompareCWEIDs](../sdk/compare-cwe-ids)
- [SDK FormatCWEIDFromInt](../sdk/format-cwe-id-from-int)
