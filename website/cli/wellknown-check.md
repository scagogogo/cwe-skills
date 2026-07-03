# 🏆 cwe wellknown check

检查给定的 CWE ID 是否属于 CWE Top 25、OWASP Top 10 或 SANS Top 25。

## 语法

```bash
cwe wellknown check [CWE-ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID...` | 是 | 一个或多个待检查的 CWE ID |

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe wellknown check CWE-79 89 352
```

```text
CWE-79: [Top 25 OWASP Top 10 (A03:2021-Injection)]
CWE-89: [Top 25 OWASP Top 10 (A03:2021-Injection) SANS Top 25]
CWE-352: [Top 25 OWASP Top 10 (A01:2021-Broken Access Control)]
```

不在任何列表中的 CWE：

```bash
cwe wellknown check CWE-999
```

```text
CWE-999: 不在任何知名列表中
```

### JSON 输出

```bash
cwe wellknown check CWE-79 89 -o json
```

```json
[
  {
    "cwe_id": "CWE-79",
    "in_list": ["Top 25", "OWASP Top 10 (A03:2021-Injection)"]
  },
  {
    "cwe_id": "CWE-89",
    "in_list": ["Top 25", "OWASP Top 10 (A03:2021-Injection)", "SANS Top 25"]
  }
]
```

## 使用场景

- 批量判断扫描结果中的弱点是否属于高风险列表，决定修复优先级。
- 漏洞报告自动标注“上榜/未上榜”。
- 筛选需要重点关注的子集。

::: tip 输入容错
无效的 CWE ID 会被原样保留并返回空 `in_list`，不会中断整体检查流程。
:::

## 下一步

- [wellknown top25](./wellknown-top25) — 查看完整 Top 25。
- [wellknown owasp](./wellknown-owasp) — 查看完整 OWASP 映射。

## 相关文档

- [SDK IsInTop25](../sdk/is-in-top-25)
- [SDK IsInOWASPTop10](../sdk/is-in-owasp-top-10)
- [SDK IsInSANSTop25](../sdk/is-in-sans-top-25)
