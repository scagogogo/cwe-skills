# 🏆 cwe wellknown top25

列出 MITRE CWE Top 25 Most Dangerous Software Weaknesses。

## 语法

```bash
cwe wellknown top25 [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe wellknown top25
```

```text
CWE Top 25 Most Dangerous Software Weaknesses (25 项):

   1. CWE-79
   2. CWE-89
   3. CWE-352
  ...
  25. CWE-94
```

### JSON 输出

```bash
cwe wellknown top25 -o json
```

```json
[
  {
    "id": 79,
    "format": "CWE-79"
  },
  {
    "id": 89,
    "format": "CWE-89"
  }
]
```

## 使用场景

- 获取当前年度最危险的 25 个软件弱点编号清单。
- 作为漏洞修复优先级排序的基线。
- 与 [`wellknown check`](./wellknown-check) 配合判断单点弱点是否上榜。

::: tip 配合 check
若只需判断个别 CWE 是否上榜，直接用 [`cwe wellknown check CWE-79`](./wellknown-check) 更高效，无需自行比对列表。
:::

## 下一步

- [wellknown owasp](./wellknown-owasp) — OWASP Top 10 映射。
- [wellknown sans](./wellknown-sans) — SANS Top 25。
- [wellknown check](./wellknown-check) — 检查归属。

## 相关文档

- [SDK CWETop25](../sdk/cwe-top-25)
- [知名列表总览](../sdk/wellknown-ids)
