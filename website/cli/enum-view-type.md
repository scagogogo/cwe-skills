# 📚 cwe enum view-type

列出 CWE 视图类型（ViewType）的所有合法取值：`Graph`、`Slice`。

::: warning 命令名
本枚举对应的子命令名为 `viewtype`（无连字符），而非 `view-type`。即实际执行命令为 `cwe enum viewtype`。
:::

## 语法

```bash
cwe enum viewtype [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum viewtype
```

```text
视图类型 (Graph/Slice) (2 项):
  - Graph
  - Slice
```

### JSON 输出

```bash
cwe enum viewtype -o json
```

```json
[
  "Graph",
  "Slice"
]
```

## 取值含义

| 取值 | 含义 |
| --- | --- |
| `Graph` | 图视图，按关系图组织弱点 |
| `Slice` | 切片视图，按特定维度切片组织弱点 |

## 使用场景

- 在使用 [`show view`](./show-view) 或 [`tree view`](./tree-view) 前了解视图类型。
- 解读 [`registry list-views`](./registry-list-views) 输出的 `Type` 字段。
- 自定义视图分析时确认合法类型。

::: tip 与 tree view 配合
[`tree view`](./tree-view) 可基于视图构建层次树，视图类型决定其组织方式。
:::

## 下一步

- [enum abstraction](./enum-abstraction)
- [show view](./show-view)
- [tree view](./tree-view)

## 相关文档

- [SDK ViewType 枚举](../sdk/enum-view-type)
- [视图概念](../guide/concept-view)
