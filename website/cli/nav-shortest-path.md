# 🧭 cwe nav shortest-path

查找两个 CWE 之间的最短路径（基于本地 XML 数据）。

<Badge type="tip" text="离线"/>

## 语法

```bash
cwe nav shortest-path <FROM> <TO> [flags]
```

## 参数

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `FROM` | 是 | 起点 CWE ID |
| `TO` | 是 | 终点 CWE ID |

::: warning 参数数量
本命令固定接收 2 个参数。
:::

## Flags

继承自 [`nav`](./nav) 的 `PersistentFlags`：

| Flag | 简写 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `--xml` | `-x` | （必填） | CWE XML 目录文件路径 |

## 示例

### text 输出

```bash
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_latest.xml
```

```text
最短路径 (3 步):
  1. CWE-79
  2. CWE-74
  3. CWE-1
```

无路径时：

```text
CWE-79 和 CWE-XXX 之间没有路径
```

### JSON 输出

```bash
cwe nav shortest-path CWE-79 CWE-1 --xml cwec_latest.xml -o json
```

```json
{
  "from": "CWE-79",
  "to": "CWE-1",
  "path": [79, 74, 1],
  "depth": 2
}
```

## 使用场景

- 分析两个弱点在关系网络中的距离与连接路径。
- 评估某弱点是否会经由中间弱点传导到另一弱点。
- 关系网络可视化时的路径提取。

::: tip depth 字段
JSON 输出中 `depth` 为路径步数（即 `path` 长度减 1）。`path` 为 CWE 数字编号数组。
:::

::: warning 路径方向
最短路径基于本地注册表构建的关系图，路径方向取决于底层关系建模。结果可能为 0 步（起点==终点）或无路径（返回空）。
:::

## 下一步

- [nav is-ancestor](./nav-is-ancestor) — 祖先关系判定。
- [nav ancestors](./nav-ancestors) — 完整祖先链。
- [tree path](./tree-path) — 层次树中的路径。

## 相关文档

- [SDK Navigator.ShortestPath](../sdk/navigator)
