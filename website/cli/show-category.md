# 🌐 cwe show category

通过 MITRE CWE REST API 获取 CWE 类别（Category）的详细信息。

<Badge type="info" text="在线"/>

## 语法

```bash
cwe show category [ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `ID...` | 是 | 一个或多个类别 ID |

## Flags

继承自父命令 [`show`](./show) 的 `PersistentFlags`：

| Flag | 默认值 | 说明 |
| --- | --- | --- |
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |
| `--timeout` | `30` | API 请求超时时间（秒） |

## 示例

### text 输出

```bash
cwe show category 1000
```

```text
=== CWE-1000 (类别) ===
  名称: ...
  描述: ...
```

批量查询多个类别：

```bash
cwe show category 1000 1001 1002
```

### JSON 输出

```bash
cwe show category 1000 -o json
```

```json
{
  "id": 1000,
  "name": "...",
  "description": "..."
}
```

::: warning JSON 模式逐条输出
源码中 `show category` 在 JSON 模式下对每个输入 ID 逐条调用 `printJSON`，因此多个 ID 时会输出多个 JSON 对象（非数组）。建议一次只查一个 ID，或在脚本中按需解析。
:::

## 使用场景

- 查阅某类别的官方定义与描述。
- 了解类别下涵盖的弱点主题。
- 与 [`show`](./show)（弱点）、[`show view`](./show-view) 配合，覆盖 CWE 的三类条目。

## 下一步

- [show](./show) — 获取弱点详情。
- [show view](./show-view) — 获取视图详情。
- [registry list-categories](./registry-list-categories) — 离线列出所有类别。

## 相关文档

- [SDK GetCategory](../sdk/api-get-category)
- [类别概念](../guide/concept-category)
