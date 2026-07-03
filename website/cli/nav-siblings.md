# 🧭 cwe nav siblings

查询指定 CWE 的同级弱点（同一父级的其他子级，基于本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav siblings [CWE-ID] [flags]
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
cwe nav siblings CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的 同级 (N 项):
  CWE-89 - ...
  ...
```

### JSON 输出

```bash
cwe nav siblings CWE-79 --xml cwec_latest.xml -o json
```

```json
{
  "cwe_id": "CWE-79",
  "type": "同级",
  "results": [
    { "id": 89, "name": "..." }
  ],
  "count": 1
}
```

## 使用场景

- 查找与某弱点共享同一父级的兄弟弱点，常用于发现同类的可替代/相关弱点。
- 安全分析中扩展弱点集合（“既然有 CWE-79，同级还有哪些？”）。
- 培训/文档场景展示同族弱点。

::: tip 与 peers 的区别
`siblings` 基于“同一父级”的层级关系；[`peers`](./nav-peers) 基于 CWE 显式声明的 `PeerOf` 对等关系，两者语义不同。
:::

## 下一步

- [nav peers](./nav-peers)
- [nav parents](./nav-parents)
- [nav children](./nav-children)

## 相关文档

- [SDK Navigator.Siblings](../sdk/navigator)
