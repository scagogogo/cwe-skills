# 🏆 cwe wellknown owasp

列出 OWASP Top 10 (2021) 对应的 CWE ID 及其分类。

输出按 OWASP 分类（如 `A01:2021-Broken Access Control`）组织，每个分类下列出对应的 CWE 编号。

## 语法

```bash
cwe wellknown owasp [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe wellknown owasp
```

```text
OWASP Top 10 (2021):

  A01:2021-Broken Access Control:
    CWE-22
    CWE-352
  A03:2021-Injection:
    CWE-79
    CWE-89
  ...
```

### JSON 输出

```bash
cwe wellknown owasp -o json
```

```json
[
  {
    "category": "A01:2021-Broken Access Control",
    "cwe_ids": [22, 352]
  },
  {
    "category": "A03:2021-Injection",
    "cwe_ids": [79, 89]
  }
]
```

## 使用场景

- 将自身发现的 CWE 映射到 OWASP Top 10 分类，便于合规报告。
- 生成按 OWASP 分类组织的修复优先级清单。
- 安全培训中展示各类别对应的典型弱点。

::: tip 获取某 CWE 的 OWASP 分类
若只想知道单个 CWE 属于哪个 OWASP 分类，用 [`wellknown check`](./wellknown-check)，输出会带上分类名，例如 `[OWASP Top 10 (A03:2021-Injection)]`。
:::

## 下一步

- [wellknown top25](./wellknown-top25)
- [wellknown sans](./wellknown-sans)
- [wellknown check](./wellknown-check)

## 相关文档

- [SDK OWASPTop10](../sdk/owasp-top-10)
- [知名列表总览](../sdk/wellknown-ids)
