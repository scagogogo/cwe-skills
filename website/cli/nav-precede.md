# 🧭 cwe nav precede

查询此弱点可以前置的弱点（`CanPrecede` 关系，本地 XML 数据）。

即：在攻击链/场景中，哪些弱点可以排在指定弱点**之后**发生。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav precede [CWE-ID] [flags]
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
cwe nav precede CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 可前置 (N 项):
  CWE-XXX - ...
```

### JSON 输出

```bash
cwe nav precede CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "可前置",
  "results": [
    { "id": 123, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 攻击链分析：某弱点发生后，可能进一步引发哪些弱点。
- 顺序关系建模，识别弱点间的因果链。
- 配合 [`nav follow`](./nav-follow) 构建完整顺序关系。

::: tip precede 与 follow 的方向
`precede` 返回“此弱点可前置的弱点”（排在之后的）；[`follow`](./nav-follow) 返回“此弱点可跟随的弱点”（排在之前的）。两者方向相反，组合使用可还原完整顺序。
:::

## 下一步

- [nav follow](./nav-follow)
- [nav requires](./nav-requires)
- [enum relationship](./enum-relationship)

## 相关文档

- [SDK Navigator.CanPrecede](../sdk/navigator)
