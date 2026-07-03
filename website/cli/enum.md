# 📚 cwe enum

查询 CWE 规范中定义的各种枚举类型。

`enum` 是父命令，本身不直接执行查询，需指定子命令列出对应枚举的全部合法取值。

## 语法

```bash
cwe enum <子命令> [flags]
```

## 子命令

`enum` 为以下枚举类型动态生成子命令：

| 子命令 | 说明 |
| --- | --- |
| [`abstraction`](./enum-abstraction) | 抽象层级（Pillar/Class/Base/Variant） |
| [`status`](./enum-status) | 状态（Stable/Draft/Deprecated 等） |
| [`relationship`](./enum-relationship) | 关系类型（ChildOf/ParentOf/CanPrecede 等） |
| [`structure`](./enum-structure) | 结构类型（Simple/Chain/Composite） |
| [`likelihood`](./enum-likelihood) | 利用可能性（High/Medium/Low） |
| [`consequence-scope`](./enum-consequence-scope) | 后果范围（Confidentiality/Integrity/Availability） |
| [`consequence-impact`](./enum-consequence-impact) | 后果影响（High/Medium/Low） |
| [`enum-view-type`](./enum-view-type) | 视图类型（Graph/Slice），命令名为 `viewtype` |
| `enum-platform` | 平台类型 |

::: warning 命令名注意
部分子命令的实际命令名与文件名风格不同：
- 后果范围/影响的命令名是 `scope` / `impact`（见 [consequence-scope](./enum-consequence-scope)、[consequence-impact](./enum-consequence-impact)）。
- 视图类型命令名是 `viewtype`（无连字符，见 [enum-view-type](./enum-view-type)）。
:::

## 示例

### 列出抽象层级

```bash
cwe enum abstraction
```

```text
抽象层级 (Class/Base/Variant/Pillar) (4 项):
  - Pillar
  - Class
  - Base
  - Variant
```

### 列出关系类型

```bash
cwe enum relationship
```

```text
关系类型 (ChildOf/ParentOf/CanPrecede等) (N 项):
  - ChildOf
  - ParentOf
  - CanPrecede
  ...
```

### JSON 输出

```bash
cwe enum status -o json
```

```json
[
  "Stable",
  "Draft",
  "Deprecated",
  "Obsolete"
]
```

## 使用场景

- 查阅某枚举字段（如 `Abstraction`、`Status`）的所有合法取值。
- 在调用 `search`/`filter` 的 `--abstraction`、`--status` 等过滤参数前确认可选值。
- 文档生成、数据校验时获取枚举全集。

::: tip 与过滤命令配合
`enum` 列出的取值可直接用于 [`filter`](./filter) 的 `--abstraction`、`--status`、`--likelihood`、`--scope`、`--structure` 等参数。
:::

## 下一步

- [enum abstraction](./enum-abstraction) — 最常用的抽象层级枚举。
- [filter](./filter) — 多条件过滤。

## 相关文档

- [SDK 枚举总览](../sdk/enums)
- [SDK Abstraction 枚举](../sdk/enum-abstraction)
