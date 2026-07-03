# 🗃️ `cwe registry member-of`

<Badge type="info" text="离线" /> <Badge type="tip" text="需 --xml" /> <Badge type="info" text="registry 子命令" />

反向查询：列出指定 CWE **所属的**所有视图和类别 ID。与 `view-members` / `category-members` 互为逆操作。

## 📖 语法

```bash
cwe registry member-of [CWE-ID] [flags]
```

| 参数 | 说明 |
|------|------|
| `CWE-ID` | 要查询归属的 CWE ID（必填，1 个） |

| Flag | 说明 |
|------|------|
| `--xml <file>` | 本地 XML 目录文件（必填） |
| `-o, --output` | 输出格式 `text`（默认）或 `json` |

## 🚀 示例

### 文本输出

```bash
# 查询 CWE-79 属于哪些视图/类别
cwe registry member-of 79 --xml cwec_latest.xml
```

```text
CWE-79 所属的视图/类别:
  - CWE-699 (Software Development)
  - CWE-1000 (Research Concepts)
  - CWE-XXX (...)
```

### JSON 输出

```bash
cwe registry member-of 79 --xml cwec_latest.xml -o json
```

```json
{ "id": 79, "member_of": [699, 1000, ...] }
```

::: tip 使用场景
- 检查某弱点是否在知名视图（如研究概念 CWE-1000）中
- 合规审计：确认弱点是否落在某行业分类下
- 与 [IsInWellKnownView](../sdk/is-in-wellknown-view) 配合判断视图归属
:::

## 🧠 对应 SDK API

`Registry.GetMemberOfIDs(id int) []int` —— 详见 [SDK 关系索引](../sdk/relationship-indexes)。

```go
reg.BuildIndexes()
memberOf := reg.GetMemberOfIDs(79)
```

## 🔗 相关

- [registry view-members](./registry-view-members) · [registry category-members](./registry-category-members) · [SDK 关系索引](../sdk/relationship-indexes) · [知名视图](../wellknown/well-known-views)
