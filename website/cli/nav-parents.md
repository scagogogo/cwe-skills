# 🧭 cwe nav parents

查询指定 CWE 的父级弱点（基于本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav parents [CWE-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID` | 是 | 待查询的 CWE ID（固定 1 个） |

## Flags

继承自 [`nav`](./nav) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe nav parents CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 父级 (1 项):
  CWE-74 - Injection
```

### JSON 输出

```bash
cwe nav parents CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "父级",
  "results": [
    { "id": 74, "name": "Injection" }
  ],
  "count": 1
}
```

## 使用场景

- 离线查找某弱点的直接上层分类。
- 与 [`nav children`](./nav-children) 配合做双向遍历。
- 构建 [`tree path`](./tree-path) 前的关系预查。

::: tip 与 registry parents 区别
[`registry parents`](./registry-relations) 返回 ID+名称；`nav parents` 返回完整 CWE 对象（含抽象层级、状态等）。需要更多字段时用 `nav`。
:::

## 下一步

- [nav children](./nav-children)
- [nav ancestors](./nav-ancestors)
- [registry parents](./registry-relations)

## 相关文档

- [SDK Navigator.Parents](../sdk/navigator)
