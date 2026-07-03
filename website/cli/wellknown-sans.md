# 🏆 cwe wellknown sans

列出 SANS Top 25 Most Dangerous Software Errors 对应的 CWE ID。

## 语法

```bash
cwe wellknown sans [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe wellknown sans
```

```text
SANS Top 25 Most Dangerous Software Errors (25 项):

   1. CWE-79
   2. CWE-89
  ...
  25. CWE-...
```

### JSON 输出

```bash
cwe wellknown sans -o json
```

```json
[
  {
    "id": 79,
    "format": "CWE-79"
  }
]
```

## 使用场景

- 参考经典的 SANS Top 25 错误清单进行修复优先级评估。
- 与 [CWE Top 25](./wellknown-top25) 对比，识别在不同权威列表中均上榜的高优先级弱点。
- 历史数据/遗留系统对照分析。

::: warning 列表时效
SANS Top 25 为较早的危险错误清单，部分条目与现行 CWE Top 25 存在差异。建议以 [top25](./wellknown-top25) 作为当前基线，[sans](./wellknown-sans) 作为补充参考。
:::

## 下一步

- [wellknown top25](./wellknown-top25)
- [wellknown owasp](./wellknown-owasp)
- [wellknown check](./wellknown-check)

## 相关文档

- [SDK SANSTop25](../sdk/sans-top-25)
- [知名列表总览](../sdk/wellknown-ids)
