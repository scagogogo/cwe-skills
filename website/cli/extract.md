# 🆔 cwe extract

从输入文本中提取所有 CWE ID。

支持识别 `CWE-79`、`cwe-89` 等格式的 ID，多个参数会被拼接为一段文本后扫描。

## 语法

```bash
cwe extract [text...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `text...` | 是 | 一段或多段待扫描文本，多参数以空格拼接 |

## Flags

本命令无专属 flags，仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe extract "受CWE-79和CWE-89影响"
```

```text
找到 2 个CWE ID:
  CWE-79
  CWE-89
```

未匹配时：

```bash
cwe extract "无相关漏洞"
```

```text
未找到CWE ID
```

### JSON 输出

```bash
cwe extract "漏洞: CWE-79, CWE-89, CWE-352" -o json
```

```json
{
  "text": "漏洞: CWE-79, CWE-89, CWE-352",
  "ids": [
    "CWE-79",
    "CWE-89",
    "CWE-352"
  ],
  "count": 3
}
```

## 使用场景

- 从漏洞通告、安全报告、工单描述中自动提取引用的 CWE。
- 批量处理扫描器输出文本，汇总涉及的 CWE 集合。
- 与 [`jq`](https://stedolan.github.io/jq/) 配合去重统计：

```bash
cwe extract "CWE-79 CWE-89 CWE-79" -o json | jq '.ids | unique'
```

::: tip 多段拼接
传入多个参数时它们以空格拼接成一整段文本再扫描，因此可分别传入标题与正文：`cwe extract "标题CWE-79" "正文CWE-89"`。
:::

## 下一步

- [parse](./parse) — 对提取出的 ID 做解析。
- [format](./format) — 对提取出的 ID 做规范化。

## 相关文档

- [SDK ExtractCWEIDs](../sdk/extract-cwe-ids)
- [SDK ExtractFirstCWEID](../sdk/extract-first-cwe-id)
