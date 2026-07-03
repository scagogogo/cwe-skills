# 🛠️ 全局参数

`cwe` 的全局参数通过 cobra 的 `PersistentFlags` 注册，对所有子命令生效。

## 语法

```bash
cwe [全局参数] <子命令> [子命令参数/flags]
```

## 全局参数

| 参数 | 简写 | 类型 | 默认值 | 说明 |
| --- | --- | --- | --- | --- |
| `--output` | `-o` | string | `text` | 输出格式，可选 `text` 或 `json` |

::: tip 简写
`-o json` 与 `--output json` 等价。`text` 为默认值，可省略。
:::

## `--output` 详解

控制所有命令的输出形式：

- `text`（默认）：人类可读的中文文本，含对齐、序号、勾叉符号等，适合终端查看。
- `json`：缩进 2 空格的结构化 JSON，适合管道处理与脚本消费。

### text 输出

```bash
cwe parse CWE-79
```

```text
CWE-79 -> CWE-79 (ID: 79)
```

### json 输出

```bash
cwe parse CWE-79 -o json
```

```json
[
  {
    "input": "CWE-79",
    "id": 79,
    "format": "CWE-79",
    "valid": true
  }
]
```

## 与 jq 配合

JSON 模式便于与 [`jq`](https://stedolan.github.io/jq/) 配合做二次处理：

```bash
# 提取所有有效 CWE ID 的数字部分
cwe parse CWE-79 abc 89 -o json | jq '.[] | select(.valid) | .id'

# 统计搜索结果数量
cwe search --xml cwec_latest.xml --keyword Injection -o json | jq 'length'
```

## 子命令专属参数

除全局参数外，各子命令还有自己的 flags，例如：

- 在线命令的 `--base-url`、`--timeout`（见 [show](./show)、[relations](./relations)）。
- 离线命令的 `--xml/-x`（见 [search](./search)、[registry](./registry)、[nav](./nav)、[tree](./tree)）。
- `tree path` 的 `--root`（见 [tree path](./tree-path)）。

这些参数仅对相应子命令生效，不影响其他命令。

::: warning 参数位置
`-o` 是全局参数，可置于子命令前或后；子命令专属 flags 必须跟在子命令之后。例如 `cwe -o json parse CWE-79` 与 `cwe parse CWE-79 -o json` 均合法。
:::

## 退出码

- `0`：命令成功执行。
- 非 `0`：发生错误（如参数缺失、CWE ID 无效、XML 解析失败、API 请求失败等），错误信息输出到 stderr。

部分命令在“部分成功”时也会返回非零退出码，例如 [`validate`](./validate) 在存在无效 ID 时返回错误。

## 下一步

- [CLI 总览](./overview)
- 体验各命令的 text/json 双输出：[parse](./parse)、[wellknown top25](./wellknown-top25)

## 相关文档

- [SDK API 客户端](../sdk/api-client)
- [错误处理](../guide/error-handling)
