# 🧭 cwe nav chain-members

查询链式弱点（Chain 结构）的成员弱点（本地 XML 数据）。

链式弱点由多个弱点按顺序组合而成，本命令列出其组成成员。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav chain-members [CWE-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID` | 是 | 链式弱点的 CWE ID（固定 1 个） |

## Flags

继承自 [`nav`](./nav) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe nav chain-members CWE-680 --xml cwec_latest.xml
```

```text
CWE-680 的 链式成员 (N 项):
  CWE-XXX - ...
  ...
```

### JSON 输出

```bash
cwe nav chain-members CWE-680 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-680",
  "type": "链式成员",
  "results": [
    { "id": 123, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 拆解链式弱点，识别其成员以便逐项修复。
- 攻击链分析：理解链式弱点如何串联多个独立弱点。
- 配合 [`enum structure`](./enum-structure) 与 [`search --chains`](./search) 使用。

::: tip 链式 vs 复合
链式弱点（Chain）成员按顺序组合；复合弱点（Composite）成员需同时存在。源码中 `nav` 还提供 `composite-members` 子命令查询复合成员，详见 [nav 总览](./nav)。
:::

## 下一步

- [nav can-also-be](./nav-can-also-be)
- [nav](./nav) — 总览含 composite-members。
- [enum structure](./enum-structure)
- [search](./search) — `--chains` 标志。

## 相关文档

- [SDK Navigator.ChainMembers](../sdk/navigator)
- [结构概念](../guide/concept-structure)
