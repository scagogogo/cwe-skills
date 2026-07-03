# 🆔 cwe parse

解析 CWE ID，提取数字部分并格式化为标准形式。

支持多种输入写法：纯数字 `79`、标准格式 `CWE-79`、大小写不敏感 `cwe-79`。

## 语法

```bash
cwe parse [CWE-ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID...` | 是 | 一个或多个待解析的 CWE ID，写法不限 |

## Flags

本命令无专属 flags，仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe parse CWE-79 89 cwe-352
```

```text
CWE-79 -> CWE-79 (ID: 79)
89 -> CWE-89 (ID: 89)
cwe-352 -> CWE-352 (ID: 352)
```

对无效输入会标注错误：

```bash
cwe parse abc
```

```text
abc -> 无效: ...
```

### JSON 输出

```bash
cwe parse CWE-79 89 -o json
```

```json
[
  {
    "input": "CWE-79",
    "id": 79,
    "format": "CWE-79",
    "valid": true
  },
  {
    "input": "89",
    "id": 89,
    "format": "CWE-89",
    "valid": true
  }
]
```

## 使用场景

- 将用户或报告中混杂写法的 CWE ID 统一规范化。
- 批量提取 CWE 的整数编号供后续处理。
- 在导入漏洞数据前做清洗。

::: tip 与 format 的区别
`parse` 同时返回数字 `id` 和标准 `format`，信息更全；[`format`](./format) 只输出标准格式字符串，更简洁。
:::

## 下一步

- [validate](./validate) — 仅校验格式是否合法。
- [format](./format) — 仅输出标准格式。

## 相关文档

- [SDK ParseCWEID](../sdk/parse-cwe-id)
- [SDK FormatCWEIDFromInt](../sdk/format-cwe-id-from-int)
