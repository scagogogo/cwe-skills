# 🌳 cwe tree view

构建指定视图（View）的 CWE 层次树，并以缩进形式输出。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe tree view [VIEW-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `VIEW-ID` | 是 | 视图 CWE ID（固定 1 个） |

## Flags

继承自 [`tree`](./tree) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe tree view 1000 --xml cwec_latest.xml
```

```text
CWE-1000 - ...
  CWE-XXX - ...
  ...
```

### JSON 输出

```bash
cwe tree view 1000 --xml cwec_latest.xml -o json
```

```json
{
  "id": 1000,
  "name": "...",
  "depth": 0,
  "is_leaf": false,
  "children": [ ... ]
}
```

::: warning 视图不存在
若无法构建指定视图的树（如视图 ID 不存在），命令返回非零退出码并报错 `无法构建视图 CWE-XXX 的树`。
:::

## 使用场景

- 按视图组织展示弱点层次（如开发视图、研究视图）。
- 在特定视图下浏览弱点分类树。
- 配合 [`show view`](./show-view) 与 [`registry list-views`](./registry-list-views) 选取视图。

::: tip 视图类型
视图类型（Graph/Slice）见 [`enum viewtype`](./enum-view-type)，可用 [`registry list-views`](./registry-list-views) 浏览所有可用视图。
:::

## 下一步

- [tree build](./tree-build) — 按根节点构建。
- [tree forest](./tree-forest) — 全部 Pillar 森林。
- [show view](./show-view) — 查看视图详情。
- [registry list-views](./registry-list-views) — 列出所有视图。

## 相关文档

- [SDK BuildViewTree](../sdk/tree)
- [视图概念](../guide/concept-view)
