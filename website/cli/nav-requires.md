# 🧭 cwe nav requires

查询此弱点所依赖的弱点（`Requires` 关系，本地 XML 数据）。

即：指定弱点的存在/可利用性，需要哪些其他弱点作为前提。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav requires [CWE-ID] [flags]
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
cwe nav requires CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 依赖 (N 项):
  CWE-XXX - ...
```

### JSON 输出

```bash
cwe nav requires CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "依赖",
  "results": [
    { "id": 123, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 识别某弱点的依赖前提，做前置修复以阻断下游。
- 风险传导分析：理解弱点间的依赖网络。
- 与 [`nav required-by`](./nav-required-by) 组合，双向理解依赖关系。

::: tip requires 与 required-by 的方向
`requires` 返回“此弱点依赖的弱点”（它的前提）；[`required-by`](./nav-required-by) 返回“依赖此弱点的弱点”（它的下游）。
:::

## 下一步

- [nav required-by](./nav-required-by)
- [nav precede](./nav-precede)
- [enum relationship](./enum-relationship)

## 相关文档

- [SDK Navigator.Requires](../sdk/navigator)
