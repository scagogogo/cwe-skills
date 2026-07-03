# 📚 cwe enum relationship

列出 CWE 关系类型（RelationshipNature）的所有合法取值，如 `ChildOf`、`ParentOf`、`CanPrecede` 等。

## 语法

```bash
cwe enum relationship [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe enum relationship
```

```text
关系类型 (ChildOf/ParentOf/CanPrecede等) (N 项):
  - ChildOf
  - ParentOf
  - CanPrecede
  - CanFollow
  - Requires
  - RequiredBy
  - CanAlsoBe
  - PeerOf
  ...
```

### JSON 输出

```bash
cwe enum relationship -o json
```

```json
[
  "ChildOf",
  "ParentOf",
  "CanPrecede",
  "CanFollow",
  "Requires",
  "RequiredBy",
  "CanAlsoBe",
  "PeerOf"
]
```

## 常见取值含义

| 取值 | 含义 |
| --- | --- |
| `ChildOf` / `ParentOf` | 层级父子关系 |
| `CanPrecede` / `CanFollow` | 顺序关系，前者可先于后者发生 |
| `Requires` / `RequiredBy` | 依赖关系 |
| `CanAlsoBe` | 同一弱点也可表现为另一弱点 |
| `PeerOf` | 对等关系 |

## 使用场景

- 理解 [`nav`](./nav) 各子命令对应的关系语义。
- 解读 [`show`](./show) 或 [`registry get`](./registry-get) 输出的 `Relationships` 字段。
- 自定义关系分析时获取合法取值集合。

::: tip 与导航命令对应
[`nav`](./nav) 的 `precede`/`follow`/`requires`/`required-by`/`can-also-be`/`peers` 等子命令正是基于这些关系类型实现。
:::

## 下一步

- [nav](./nav) — 本地关系导航。
- [relations](./relations) — 在线 API 关系查询。

## 相关文档

- [SDK RelationshipNature 枚举](../sdk/enum-relationship-nature)
- [关系概念](../guide/concept-relationship)
