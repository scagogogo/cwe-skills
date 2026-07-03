# 🗃️ cwe registry contains

检查给定 CWE ID 是否存在于本地注册表。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry contains [CWE-ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID...` | 是 | 一个或多个待检查的 CWE ID |

## Flags

继承自 [`registry`](./registry) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe registry contains CWE-79 CWE-999 --xml cwec_latest.xml
```

```text
CWE-79 ✓ 存在 (Weakness)
CWE-999 ✗ 不存在
```

### JSON 输出

```bash
cwe registry contains CWE-79 CWE-999 --xml cwec_latest.xml -o json
```

```json
[
  {
    "cwe_id": "CWE-79",
    "exists": true,
    "type": "Weakness"
  },
  {
    "cwe_id": "CWE-999",
    "exists": false
  }
]
```

## 使用场景

- 批量验证一组 CWE ID 是否真实存在于 MITRE 数据库（区别于 [`validate`](./validate) 仅校验格式）。
- 导入漏洞数据前做存在性清洗。
- 配合 [`wellknown check`](./wellknown-check) 区分“格式合法但不存在”与“存在但不在知名列表”。

::: tip type 字段
对存在的条目，输出会附带 `type`（如 `Weakness`、`Category`、`View`），可据此判断条目类型。
:::

## 下一步

- [registry get](./registry-get) — 获取详情。
- [validate](./validate) — 仅校验格式。
- [wellknown check](./wellknown-check) — 检查是否在知名列表。

## 相关文档

- [SDK Registry.Contains](../sdk/registry)
- [SDK CWE 类型方法](../sdk/cwe-type-methods)
