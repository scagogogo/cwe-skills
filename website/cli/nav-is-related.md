# 🧭 `cwe nav is-related`

<Badge type="info" text="离线" /> <Badge type="tip" text="需 --xml" /> <Badge type="info" text="nav 子命令" />

检查两个 CWE 之间是否存在**任意关系**（父子、对等、依赖、顺序等任一）。

::: tip 与 is-ancestor 的区别
- `is-ancestor`：仅判断**层级祖先**关系（A 是否是 B 的祖先）
- `is-related`：判断**任意类型**关系（包括非层级的 PeerOf/Requires/CanPrecede 等）
:::

## 📖 语法

```bash
cwe nav is-related <CWE-ID1> <CWE-ID2> [flags]
```

| 参数 | 说明 |
|------|------|
| `CWE-ID1` | 第一个 CWE ID（必填） |
| `CWE-ID2` | 第二个 CWE ID（必填） |

| Flag | 说明 |
|------|------|
| `--xml <file>` | 本地 XML 目录文件（必填） |
| `-o, --output` | 输出格式 `text`（默认）或 `json` |

## 🚀 示例

### 文本输出

```bash
cwe nav is-related CWE-79 CWE-74 --xml cwec_latest.xml
```

```text
CWE-79 与 CWE-74 存在关系: true
```

### JSON 输出

```bash
cwe nav is-related CWE-79 CWE-74 --xml cwec_latest.xml -o json
```

```json
{ "related": true }
```

## 🧠 对应 SDK API

`Navigator.IsRelated(a, b int) bool` —— 详见 [SDK nav-ancestor-related](../sdk/nav-ancestor-related)。

```go
nav := cweskills.NewNavigator(reg)
if nav.IsRelated(79, 74) {
    fmt.Println("存在关系")
}
```

## 🎯 使用场景

- 漏洞关联分析：判断两个报告的弱点是否属于同一关系簇
- 合并去重：识别对等（PeerOf）或可互换的 CWE

## 🔗 相关

- [nav is-ancestor](./nav-is-ancestor) — 祖先判断
- [nav shortest-path](./nav-shortest-path) · [SDK IsRelated](../sdk/nav-ancestor-related)
