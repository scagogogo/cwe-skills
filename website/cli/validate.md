# 🆔 cwe validate

验证 CWE ID 格式是否有效。

## 语法

```bash
cwe validate [CWE-ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID...` | 是 | 一个或多个待验证的 CWE ID |

## Flags

本命令无专属 flags，仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe validate CWE-79 CWE-89 abc
```

```text
CWE-79 ✓ 有效
CWE-89 ✓ 有效
abc ✗ 无效: ...
```

::: warning 退出码
当存在任意无效 ID 时，命令返回非零退出码并输出错误到 stderr，便于在脚本中据此中断流程。
:::

### JSON 输出

```bash
cwe validate CWE-79 abc -o json
```

```json
[
  {
    "input": "CWE-79",
    "valid": true
  },
  {
    "input": "abc",
    "valid": false,
    "error": "..."
  }
]
```

## 使用场景

- 在 CI/流水线中校验漏洞报告引用的 CWE ID 是否合法。
- 表单或配置文件录入前的格式校验。
- 与 `parse` 配合：先用 `validate` 过滤，再用 `parse` 提取数字。

::: tip 仅校验格式
`validate` 只检查字符串是否符合 CWE ID 格式，不验证该编号是否真实存在于 MITRE 数据库。若需确认存在性，请用 [registry contains](./registry-contains)（离线）或 [show](./show)（在线）。
:::

## 下一步

- [parse](./parse) — 解析并提取数字。
- [format](./format) — 格式化为标准形式。
- [registry contains](./registry-contains) — 检查是否真实存在。

## 相关文档

- [SDK ValidateCWEID](../sdk/validate-cwe-id)
- [SDK IsCWEID](../sdk/is-cwe-id)
