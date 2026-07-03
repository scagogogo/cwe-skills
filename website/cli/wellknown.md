# 🏆 cwe wellknown

查询 CWE 的各种知名列表，包括 CWE Top 25、OWASP Top 10、SANS Top 25。

`wellknown` 是父命令，需指定子命令执行具体查询。

## 语法

```bash
cwe wellknown <子命令> [flags]
```

## 子命令

| 子命令 | 说明 |
| --- | --- |
| [`top25`](./wellknown-top25) | 列出 CWE Top 25 Most Dangerous Software Weaknesses |
| [`owasp`](./wellknown-owasp) | 列出 OWASP Top 10 (2021) 对应的 CWE 及分类 |
| [`sans`](./wellknown-sans) | 列出 SANS Top 25 Most Dangerous Software Errors |
| [`check`](./wellknown-check) | 检查给定 CWE 是否在上述知名列表中 |

## 示例

### 列出 Top 25

```bash
cwe wellknown top25
```

```text
CWE Top 25 Most Dangerous Software Weaknesses (25 项):

   1. CWE-79
   2. CWE-89
   ...
```

### 检查某 CWE 是否在列表中

```bash
cwe wellknown check CWE-79 89 352
```

```text
CWE-79: [Top 25 OWASP Top 10 (A03:2021-Injection)]
...
```

### JSON 输出

```bash
cwe wellknown owasp -o json | jq '.[].category'
```

## 使用场景

- 评估漏洞库/扫描器结果中哪些弱点属于“最危险”集合，优先修复。
- 合规对照：将自身发现映射到 OWASP Top 10。
- 安全培训/报告素材：展示当前最具威胁的弱点。

::: tip 数据来源
列表为 SDK 内置的权威映射（`CWETop25`、`OWASPTop10`、`SANSTop25`），无需联网或 XML 文件即可使用。
:::

## 下一步

- [wellknown top25](./wellknown-top25) — 最常用的危险弱点列表。
- [wellknown check](./wellknown-check) — 批量检查。

## 相关文档

- [SDK 知名列表总览](../sdk/wellknown-ids)
- [OWASP Top 10 映射](../sdk/owasp-top-10)
