# 🗃️ cwe registry get

获取本地注册表中单个 CWE 条目的详情。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe registry get [CWE-ID] [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `CWE-ID` | 是 | 待查询的 CWE ID（固定 1 个） |

## Flags

继承自 [`registry`](./registry) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe registry get CWE-79 --xml cwec_latest.xml
```

```text
=== CWE-79 ===
  名称:     Improper Neutralization of Input During Web Page Generation
  抽象层级: Base
  状态:     Stable
  结构:     Simple
  描述:     ...
  利用可能性: High
  关系:     N 项
  后果:     N 项
```

### JSON 输出

```bash
cwe registry get CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "id": 79,
  "name": "Improper Neutralization of Input During Web Page Generation",
  "abstraction": "Base",
  "status": "Stable",
  "structure": "Simple",
  "description": "...",
  "likelihood_of_exploit": "High"
}
```

::: warning 不存在的 ID
若指定 ID 不在注册表中，命令返回非零退出码并报错 `CWE-XXX 不存在于注册表中`。
:::

## 使用场景

- 离线查阅某弱点的完整详情（含描述、关系、后果等）。
- 替代在线 [`show`](./show)，零网络依赖、速度更快。
- 在脚本中按需获取单条数据。

::: tip 与 show 对比
[`show`](./show) 在线获取，[`registry get`](./registry-get) 离线获取。两者输出的核心字段一致，但 `registry get` 通常更全（依赖本地 XML 的完整内容）。已有 XML 时优先用 `registry get`。
:::

## 下一步

- [registry contains](./registry-contains) — 检查存在性。
- [registry load](./registry-load) — 加载概要。
- [show](./show) — 在线替代。

## 相关文档

- [SDK Registry.Get](../sdk/registry)
- [SDK CWE 结构体](../sdk/cwe-struct)
