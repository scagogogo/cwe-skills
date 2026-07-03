# 🧭 cwe nav follow

查询此弱点可以跟随的弱点（`CanFollow` 关系，本地 XML 数据）。

即：在攻击链/场景中，哪些弱点可以排在指定弱点**之前**发生。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav follow [CWE-ID] [flags]
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
cwe nav follow CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 可跟随 (N 项):
  CWE-XXX - ...
```

### JSON 输出

```bash
cwe nav follow CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "可跟随",
  "results": [
    { "id": 123, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 攻击链溯源：某弱点之前可能由哪些弱点引发。
- 根因分析，识别上游弱点以做前置防御。
- 与 [`nav precede`](./nav-precede) 组合还原完整顺序关系。

::: tip precede 与 follow 的方向
[`precede`](./nav-precede) 返回“此弱点可前置的弱点”（排在之后的）；`follow` 返回“此弱点可跟随的弱点”（排在之前的）。
:::

## 下一步

- [nav precede](./nav-precede)
- [nav required-by](./nav-required-by)
- [enum relationship](./enum-relationship)

## 相关文档

- [SDK Navigator.CanFollow](../sdk/navigator)
