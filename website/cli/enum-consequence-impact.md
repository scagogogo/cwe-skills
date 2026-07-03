# 📚 cwe enum consequence-impact

列出 CWE 后果影响（ConsequenceImpact）的所有合法取值。

::: warning 命令名
本枚举对应的子命令名为 `impact`，而非 `consequence-impact`。即实际执行命令为 `cwe enum impact`。
:::

## 语法

```bash
cwe enum impact [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum impact
```

```text
后果影响 (High/Medium/Low) (3 项):
  - High
  - Medium
  - Low
```

### JSON 输出

```bash
cwe enum impact -o json
```

```json
[
  "High",
  "Medium",
  "Low"
]
```

## 取值含义

| 取值 | 含义 |
| --- | --- |
| `High` | 后果严重，损失重大 |
| `Medium` | 中等影响 |
| `Low` | 影响较轻 |

## 使用场景

- 评估弱点被利用后的后果严重程度。
- 与 [likelihood](./enum-likelihood) 结合做风险分级（可能性 × 影响）。
- 优先级排序时参考。

::: tip 与 likelihood 区别
`likelihood` 描述“多容易被利用”，`impact` 描述“利用后有多严重”，两者共同决定风险等级。两者是独立维度，可同时出现在过滤条件中。
:::

## 下一步

- [enum likelihood](./enum-likelihood)
- [enum consequence-scope](./enum-consequence-scope)

## 相关文档

- [SDK ConsequenceImpact 枚举](../sdk/enum-consequence-impact)
- [后果概念](../guide/concept-consequence)
