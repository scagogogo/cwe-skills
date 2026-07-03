# 🆔 cwe format

将 CWE ID 格式化为标准形式 `CWE-XXX`（去前导零、统一大写前缀）。

## 语法

```bash
cwe format [CWE-ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID...` | 是 | 一个或多个待格式化的 CWE ID |

## Flags

本命令无专属 flags，仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe format 79 cwe-89 CWE-352
```

```text
CWE-79
CWE-89
CWE-352
```

对无效输入标注错误：

```bash
cwe format abc
```

```text
abc -> 错误: ...
```

### JSON 输出

```bash
cwe format 79 cwe-89 -o json
```

```json
[
  {
    "input": "79",
    "output": "CWE-79"
  },
  {
    "input": "cwe-89",
    "output": "CWE-89"
  }
]
```

## 使用场景

- 批量规范化 CWE ID 输出，便于对齐与去重。
- 生成报告、表格时统一展示格式。
- 作为管道中的一步，把用户输入标准化后喂给下游命令。

::: tip 与 parse 的区别
`format` 输出更简洁，每行一个标准 ID；[`parse`](./parse) 额外给出数字 `id` 与有效性，信息更全。按需选用。
:::

## 下一步

- [parse](./parse) — 解析并提取数字。
- [validate](./validate) — 仅校验格式。

## 相关文档

- [SDK FormatCWEIDFromInt](../sdk/format-cwe-id-from-int)
- [SDK FormatCWEID](../sdk/format-cwe-id)
