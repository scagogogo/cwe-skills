---
title: 输出格式 (text/JSON)
outline: [2, 3]
---

# 📤 输出格式 (text/JSON)

CWE Skills 的 CLI 全部命令支持 **`-o text|json`** 双格式输出：`text` 面向人类阅读（带标签、缩进、说明），`json` 面向脚本与 AI 解析（结构化、字段稳定）。默认 `text`。这是 CLI 能同时服务开发者和自动化管道的关键设计。

::: tip 一个参数统天下
`-o` / `--output` 是 CLI 的**全局参数**，所有子命令都支持。无需为每个命令记不同的「格式参数」。
:::

---

## 🎛️ 统一 -o 参数

| 参数 | 值 | 默认 | 说明 |
|------|----|----|------|
| `-o, --output` | `text` / `json` | `text` | 输出格式 |

```bash
cwe show CWE-79            # text 格式（默认）
cwe show CWE-79 -o json    # JSON 格式
cwe show CWE-79 --output json # 长写法
```

::: info 全局参数
`-o` 是全局参数，写在子命令前或后都行：`cwe -o json show CWE-79` 与 `cwe show CWE-79 -o json` 等价。
:::

---

## 👤 text 格式（人类可读）

```bash
$ cwe show CWE-79
```

```text
CWE-79: Cross-site Scripting (XSS)
  抽象层级: Base
  状态:     Stable
  结构:     Simple
  描述:     The software does not neutralize ...
  常见后果: 机密性 High / 完整性 High
```

`text` 格式带对齐的标签、中文说明、缩进，适合终端里直接阅读。但不适合程序解析——标签和格式可能随版本微调。

---

## 🧮 JSON 格式（脚本/AI 友好）

```bash
$ cwe show CWE-79 -o json
```

```json
{
  "id": 79,
  "name": "Cross-site Scripting (XSS)",
  "abstraction": "Base",
  "status": "Stable",
  "structure": "Simple",
  "description": "The software does not neutralize ...",
  "consequences": [
    { "scope": "Confidentiality", "impact": "High" }
  ]
}
```

::: tip 字段稳定
JSON 字段名与结构是 SDK 类型（`*CWE` 等）的直接序列化，**跨版本稳定**。这是 `json` 输出适合 `jq` 管道与 AI 解析的根本原因。
:::

---

## 🔧 用 jq 管道处理 JSON

`-o json` 配合 `jq` 是 CLI 脚本集成的标准姿势：

### 提取单字段

```bash
# 只取名称
cwe show CWE-79 -o json | jq -r '.name'
# => Cross-site Scripting (XSS)

# 只取抽象层级
cwe show CWE-79 -o json | jq -r '.abstraction'
# => Base
```

### 批量处理

```bash
# 批量查名称
for id in 79 89 352 787; do
  echo "$id: $(cwe show CWE-$id -o json | jq -r '.name')"
done
```

### 过滤

```bash
# 列出所有 Base 级别 + Stable 的 CWE（离线）
cwe search --xml cwec_v4.15.xml --keyword injection -o json \
  | jq '.[] | select(.abstraction=="Base" and .status=="Stable") | .id'
```

### 合并结果

```bash
# 在线详情 + 离线祖先链，合并成一个 JSON
jq -s '{weakness: .[0], ancestors: .[1]}' \
  <(cwe show CWE-79 -o json) \
  <(cwe nav ancestors CWE-79 --xml cwec_v4.15.xml -o json)
```

::: details 为什么不用 text + grep
`text` 输出的标签、对齐、多行描述用 `grep`/`awk` 解析脆弱——格式微调就崩。JSON + `jq` 用字段名取值，稳定且可读。见 [CLI 接入](./integration-cli) 的脚本示例。
:::

---

## 📦 SDK 序列化

Go SDK 侧的序列化能力与 CLI 的 `-o json` 同源——都是把 `cweskills` 类型序列化：

```go
// 单条
jsonBytes, _ := cweskills.MarshalJSON(cwe79)
cweBack, _   := cweskills.UnmarshalJSON(jsonBytes)

// 整库
allJSON, _ := registry.ExportJSON()
_ = registry.ImportJSON(allJSON) // 反向导入

// 也支持 XML / CSV
xmlBytes, _ := cweskills.MarshalXML(cwe79)
csvRow, _   := cweskills.MarshalCSV([]*cweskills.CWE{cwe79})
```

::: tip CLI 与 SDK 输出一致
CLI 的 `-o json` 内部就是调用 `cweskills.MarshalJSON`，所以 CLI 的 JSON 输出与 SDK 的 `MarshalJSON` 字段完全一致——你可以用 CLI 探索字段结构，再在 SDK 里用同样的字段名。
:::

---

## 🤖 AI 用 JSON 更稳

AI 代理调用 CLI 时，`text` 输出可能有歧义（多行、标签、缩进），AI 解析易错。**建议让 AI 加 `-o json`**：

```text
你: 查 CWE-79 的详情，用 JSON 拿结果再总结。

AI: （调用 cwe show CWE-79 -o json，解析 JSON 字段）
    名称：Cross-site Scripting (XSS)
    抽象层级：Base
    状态：Stable
    ...
```

::: warning 提示 AI 用 -o json
在 [Skills 提示词](./integration-skills) 里已写明「所有命令支持 `-o json`」。若 AI 忘了加，追加一句「请用 -o json 重新调用」即可。JSON 输出让 AI 解析无歧义。
:::

---

## ⚠️ 错误时的输出

| 场景 | stdout | stderr | 退出码 |
|------|--------|--------|--------|
| 成功 | text/JSON 结果 | 空 | 0 |
| 失败 | 空 | 错误文本 | 非 0 |

::: danger JSON 输出只在成功时产生
失败时 CLI **不会**输出 JSON（避免污染管道），错误信息走 stderr，退出码非 0。脚本里先判退出码，再解析 stdout 的 JSON：

```bash
if ! cwe show CWE-79 -o json > /tmp/cwe79.json 2>/tmp/err.txt; then
  echo "失败: $(cat /tmp/err.txt)"
  exit 1
fi
jq '.name' /tmp/cwe79.json
```
:::

---

## 📖 相关文档

- [CLI 接入](./integration-cli)（脚本集成示例）
- [Skills 接入（AI 代理）](./integration-skills)
- [错误处理](./error-handling)（错误时的输出与退出码）
- [Go SDK 接入](./integration-sdk)（SDK 序列化）
- [在线 vs 离线模式](./online-offline)
