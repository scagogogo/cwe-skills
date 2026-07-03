# 🧭 cwe nav required-by

查询依赖此弱点的弱点（`RequiredBy` 关系，本地 XML 数据）。

即：哪些其他弱点的存在/可利用性，以指定弱点为前提。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav required-by [CWE-ID] [flags]
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
cwe nav required-by CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 被依赖 (N 项):
  CWE-XXX - ...
```

### JSON 输出

```bash
cwe nav required-by CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "被依赖",
  "results": [
    { "id": 123, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 评估修复某弱点的影响范围：哪些下游弱点会因此失去前提。
- 风险传导分析：识别某弱点作为根因会引发哪些下游问题。
- 与 [`nav requires`](./nav-requires) 组合，双向理解依赖关系。

::: tip requires 与 required-by 的方向
[`requires`](./nav-requires) 返回“此弱点依赖的弱点”（它的前提）；`required-by` 返回“依赖此弱点的弱点”（它的下游）。
:::

## 下一步

- [nav requires](./nav-requires)
- [nav follow](./nav-follow)
- [enum relationship](./enum-relationship)

## 相关文档

- [SDK Navigator.RequiredBy](../sdk/navigator)
