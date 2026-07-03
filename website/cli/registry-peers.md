# 🗃️ `cwe registry peers`

<Badge type="info" text="离线" /> <Badge type="tip" text="需 --xml" /> <Badge type="info" text="registry 子命令" />

查询本地注册表中指定 CWE 的**对等（PeerOf）关系**弱点。对等弱点在语义上等价或可互换。

::: tip 与 nav peers 的区别
`registry peers` 直接基于注册表的关系索引查询；`nav peers` 通过 Navigator 封装，底层一致但支持更多关系类型组合。两者数据源相同。
:::

## 📖 语法

```bash
cwe registry peers [CWE-ID] [flags]
```

| 参数 | 说明 |
|------|------|
| `CWE-ID` | 要查询的 CWE ID（必填，1 个） |

| Flag | 说明 |
|------|------|
| `--xml <file>` | 本地 MITRE CWE XML 目录文件（必填） |
| `-o, --output` | 输出格式 `text`（默认）或 `json` |

## 🚀 示例

### 文本输出

```bash
cwe registry peers CWE-79 --xml cwec_latest.xml
```

```text
CWE-79 的对等弱点:
  - CWE-XXX (Name)
```

### JSON 输出

```bash
cwe registry peers CWE-79 --xml cwec_latest.xml -o json
```

```json
[{ "id": 79, "peers": [123] }]
```

## 🧠 对应 SDK API

`Registry.GetPeerIDs(id int) []int` —— 详见 [SDK 关系索引](../sdk/relationship-indexes)。

```go
reg, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
reg.BuildIndexes()
peers := reg.GetPeerIDs(79)
```

## 🔗 相关

- [registry 总览](./registry) · [nav peers](./nav-peers) · [SDK 关系索引](../sdk/relationship-indexes)
