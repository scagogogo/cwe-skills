# 🧭 `cwe nav depth`

<Badge type="info" text="离线" /> <Badge type="tip" text="需 --xml" /> <Badge type="info" text="nav 子命令" />

计算两个 CWE 之间的**关系深度**（最短路径的边数）。用于衡量两个弱点在层次结构中的「距离」。

::: warning 返回值
- 深度 `0`：两个 ID 相同
- 深度 `N`：从 A 到 B 的最短路径经过 N 条边
- 深度 `-1`：两 CWE 无任何关系路径
:::

## 📖 语法

```bash
cwe nav depth <FROM> <TO> [flags]
```

| 参数 | 说明 |
|------|------|
| `FROM` | 起点 CWE ID（必填） |
| `TO` | 终点 CWE ID（必填） |

| Flag | 说明 |
|------|------|
| `--xml <file>` | 本地 XML 目录文件（必填） |
| `-o, --output` | 输出格式 `text`（默认）或 `json` |

## 🚀 示例

### 文本输出

```bash
cwe nav depth CWE-79 CWE-1 --xml cwec_latest.xml
```

```text
CWE-79 到 CWE-1 的关系深度: 4
```

### JSON 输出

```bash
cwe nav depth CWE-79 CWE-1 --xml cwec_latest.xml -o json
```

```json
{ "depth": 4 }
```

## 🧠 对应 SDK API

`Navigator.RelationshipDepth(ancestor, descendant int) int` —— 详见 [SDK nav-relationship-depth](../sdk/nav-relationship-depth)。

```go
nav := cweskills.NewNavigator(reg)
d := nav.RelationshipDepth(79, 1)
// d == 4 表示 79 是 1 的后代，深度为 4
```

::: tip 与 shortest-path 的关系
`shortest-path` 返回**完整路径**（节点序列），`depth` 仅返回**路径长度**（边数 = 节点数 - 1）。当你只关心距离而非具体路径时，用 `depth` 更轻量。
:::

## 🎯 使用场景

- 量化弱点的层级距离，用于优先级/影响范围评估
- 构建弱点关系图谱时的边权重计算

## 🔗 相关

- [nav shortest-path](./nav-shortest-path) — 完整路径
- [SDK RelationshipDepth](../sdk/nav-relationship-depth)
