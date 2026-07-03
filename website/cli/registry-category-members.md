# 🗃️ `cwe registry category-members`

<Badge type="info" text="离线" /> <Badge type="tip" text="需 --xml" /> <Badge type="info" text="registry 子命令" />

列出指定**类别（Category）**包含的所有 CWE 成员 ID。类别是对弱点的非层级分组（如「PHP 相关弱点」「缓冲区溢出相关」），详见 [类别概念](../guide/concept-category)。

## 📖 语法

```bash
cwe registry category-members [CATEGORY-ID] [flags]
```

| 参数 | 说明 |
|------|------|
| `CATEGORY-ID` | 类别的 CWE ID（必填，1 个） |

| Flag | 说明 |
|------|------|
| `--xml <file>` | 本地 XML 目录文件（必填） |
| `-o, --output` | 输出格式 `text`（默认）或 `json` |

## 🚀 示例

### 文本输出

```bash
cwe registry category-members 25 --xml cwec_latest.xml
```

```text
类别 CWE-25 (Path Traversal) 成员 (N 个):
  - CWE-XXX
  ...
```

### JSON 输出

```bash
cwe registry category-members 25 --xml cwec_latest.xml -o json
```

```json
{ "category": 25, "members": [22, 23, ...] }
```

::: warning 类别 vs 视图
- **类别**：扁平的、按主题的弱点集合（MemberOf 关系）
- **视图**：可层级化的、按视角的组织（Graph/Explicit Slice/Implicit Slice）
两者都用 `MemberOf`/`Has_Member` 关系，但 `GetCategoryMembers` 与 `GetViewMembers` 是不同的索引查询。详见 [类别](../guide/concept-category) vs [视图](../guide/concept-view)。
:::

## 🧠 对应 SDK API

`Registry.GetCategoryMembers(categoryID int) []int` —— 详见 [SDK 关系索引](../sdk/relationship-indexes)。

```go
reg.BuildIndexes()
members := reg.GetCategoryMembers(25)
```

## 🔗 相关

- [registry list-categories](./registry-list-categories) · [registry view-members](./registry-view-members) · [SDK 关系索引](../sdk/relationship-indexes) · [类别概念](../guide/concept-category)
