# 🌐 cwe show view

通过 MITRE CWE REST API 获取 CWE 视图（View）的详细信息。

<Badge type="info" text="在线"/>

## 语法

```bash
cwe show view [ID...] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `ID...` | 是 | 一个或多个视图 ID |

## Flags

继承自父命令 [`show`](./show) 的 `PersistentFlags`：

| Flag | 默认值 | 说明 |
| --- | --- | --- |
| `--base-url` | `https://cwe-api.mitre.org/api` | MITRE API 基础 URL |
| `--timeout` | `30` | API 请求超时时间（秒） |

## 示例

### text 输出

```bash
cwe show view 1000
```

```text
=== CWE-1000 (视图) ===
  名称: ...
  类型: Graph
  描述: ...
```

批量查询：

```bash
cwe show view 1000 699 1003
```

### JSON 输出

```bash
cwe show view 1000 -o json
```

```json
{
  "id": 1000,
  "name": "...",
  "type": "Graph",
  "description": "..."
}
```

::: warning JSON 模式逐条输出
与 [`show category`](./show-category) 一致，`show view` 在 JSON 模式下对每个 ID 逐条输出 JSON 对象，建议单次查询单个视图。
:::

## 使用场景

- 查阅视图（如开发视图、研究视图）的定义与类型。
- 在使用 [`tree view`](./tree-view) 基于视图构建层次树前，先确认视图类型。
- 合规分析时定位特定视图下的弱点集合。

::: tip 视图类型
`type` 字段取值见 [`enum viewtype`](./enum-view-type)：`Graph` 或 `Slice`。
:::

## 下一步

- [show](./show) — 获取弱点详情。
- [show category](./show-category) — 获取类别详情。
- [tree view](./tree-view) — 基于视图构建树。
- [registry list-views](./registry-list-views) — 离线列出所有视图。

## 相关文档

- [SDK GetView](../sdk/api-get-view)
- [视图概念](../guide/concept-view)
