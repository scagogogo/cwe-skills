# 🧭 cwe nav is-ancestor

检查一个 CWE 是否为另一个 CWE 的祖先（基于本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav is-ancestor <ANCESTOR> <DESCENDANT> [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `ANCESTOR` | 是 | 候选祖先 CWE ID |
| `DESCENDANT` | 是 | 候选后代 CWE ID |

::: warning 参数顺序
参数顺序为“祖先 后代”，即先传候选祖先，再传候选后代。例如判断 CWE-1 是否为 CWE-79 的祖先：`cwe nav is-ancestor CWE-1 CWE-79`。
:::

## Flags

继承自 [`nav`](./nav) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe nav is-ancestor CWE-1 CWE-79 --xml cwec_latest.xml
```

```text
CWE-1 是 CWE-79 的祖先
```

否定情形：

```text
CWE-79 不是 CWE-1 的祖先
```

### JSON 输出

```bash
cwe nav is-ancestor CWE-1 CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "ancestor": "CWE-1",
  "descendant": "CWE-79",
  "is_ancestor": true
}
```

## 使用场景

- 快速判断两弱点是否在同一祖先链上，无需展开完整 ancestors 列表。
- 校验 [`tree path`](./tree-path) 的根节点选择是否合理。
- 规则引擎中作为关系判定原语。

::: tip 与 shortest-path 配合
若需进一步了解祖先到后代的完整路径，可在 `is-ancestor` 返回 true 后调用 [`nav shortest-path`](./nav-shortest-path) 获取路径详情。
:::

## 下一步

- [nav ancestors](./nav-ancestors) — 获取完整祖先链。
- [nav shortest-path](./nav-shortest-path) — 路径详情。
- [tree path](./tree-path) — 层次树路径。

## 相关文档

- [SDK Navigator.IsAncestorOf](../sdk/navigator)
