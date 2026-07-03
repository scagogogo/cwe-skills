# 🗃️ `cwe registry view-members`

<Badge type="info" text="离线" /> <Badge type="tip" text="需 --xml" /> <Badge type="info" text="registry 子命令" />

列出指定**视图（View）**包含的所有 CWE 成员 ID。视图是 CWE 的组织方式（如 CWE-1000 研究概念、CWE-699 软件开发），详见 [视图概念](../guide/concept-view)。

## 📖 语法

```bash
cwe registry view-members [VIEW-ID] [flags]
```

| 参数 | 说明 |
|------|------|
| `VIEW-ID` | 视图的 CWE ID（必填，1 个） |

| Flag | 说明 |
|------|------|
| `--xml <file>` | 本地 XML 目录文件（必填） |
| `-o, --output` | 输出格式 `text`（默认）或 `json` |

## 🚀 示例

### 文本输出

```bash
# 列出软件开发视图(CWE-699)的所有成员
cwe registry view-members 699 --xml cwec_latest.xml
```

```text
视图 CWE-699 (Software Development) 成员 (N 个):
  - CWE-XXX
  - CWE-YYY
  ...
```

### JSON 输出

```bash
cwe registry view-members 699 --xml cwec_latest.xml -o json
```

```json
{ "view": 699, "members": [20, 22, 74, 79, ...] }
```

::: tip 知名视图速查
- `699` 软件开发视图
- `1000` 研究概念视图
- `1199` 硬件设计视图
- `888` CWE 横截面视图
- `1400` 综合字典
详见 [知名视图](../wellknown/well-known-views)。
:::

## 🧠 对应 SDK API

`Registry.GetViewMembers(viewID int) []int` —— 详见 [SDK 关系索引](../sdk/relationship-indexes)。

```go
reg.BuildIndexes()
members := reg.GetViewMembers(699)
```

## 🔗 相关

- [registry list-views](./registry-list-views) · [tree view](./tree-view) · [SDK 关系索引](../sdk/relationship-indexes) · [知名视图](../wellknown/well-known-views)
