# 🗃️ cwe registry list-categories

列出本地注册表中的所有 CWE 类别（Category）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry list-categories [flags]
```

## 参数

本命令无位置参数。

## Flags

继承自 [`registry`](./registry) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe registry list-categories --xml cwec_latest.xml
```

```text
类别 (300 项):
  CWE-1000 - ...
  CWE-1001 - ...
  ...
```

### JSON 输出

```bash
cwe registry list-categories --xml cwec_latest.xml -o json | jq 'length'
```

```json
[
  { "id": 1000, "name": "..." }
]
```

## 使用场景

- 浏览所有类别，挑选目标类别查看其成员（`registry category-members`）。
- 统计类别总数，与 [`stats`](./stats) 输出对照。
- 与 [`show category`](./show-category) 配合：先列出，再查详情。

::: tip 类别成员
源码中 `registry` 还提供 `category-members [CATEGORY-ID]` 子命令，可列出某类别下的成员 CWE，详见 [registry 总览](./registry)。
:::

## 下一步

- [registry list-views](./registry-list-views) — 列出视图。
- [show category](./show-category) — 在线查看类别详情。
- [registry](./registry) — 总览含 category-members。

## 相关文档

- [SDK Registry.GetAllCategories](../sdk/registry)
- [类别概念](../guide/concept-category)
