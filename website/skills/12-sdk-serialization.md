---
title: 技能 12 — SDK 序列化
outline: [2, 3]
---

# 📦 技能 12 — SDK 序列化

以 JSON、XML、CSV 三种格式导入导出 CWE 数据。支持单条目与通过注册表的批量操作。

<Badge type="tip" text="互操作"/>
<Badge type="info" text="三格式"/>

---

## 🎯 技能目标

- 单个 CWE 与 CWE 列表的 JSON / XML / CSV 序列化
- 通过注册表批量导入导出
- 解析 MITRE XML 目录为内存注册表

---

## 🔧 JSON 序列化

### 单个 CWE

```go
cwe := &cweskills.CWE{ID: 79, Name: "XSS", Abstraction: cweskills.AbstractionBase}

data, err := cweskills.MarshalJSON(cwe)        // 序列化
parsed, err := cweskills.UnmarshalJSON(data)   // 反序列化
```

### CWE 列表

```go
cwes := []*cweskills.CWE{cwe1, cwe2}
data, err := cweskills.MarshalJSONList(cwes)
parsed, err := cweskills.UnmarshalJSONList(data)
```

### 注册表 JSON 往返

```go
data, _ := registry.ExportJSON()             // 导出整个注册表
newRegistry := cweskills.NewRegistry()
err := newRegistry.ImportJSON(data)          // 导入到新注册表
```

::: tip ExportJSON 保留全部类型
注册表 JSON 导出保留全部条目类型：弱点、类别、视图、复合元素。
:::

---

## 🔧 XML 序列化

### 单个 CWE

```go
data, err := cweskills.MarshalXML(cwe)        // 含 <?xml ...?> 头
parsed, err := cweskills.UnmarshalXML(data)
```

### MITRE XML 目录

```go
parser := cweskills.NewXMLParser()
registry, err := parser.ParseFile("cwec_v4.15.xml")   // 从文件
registry, err := parser.Parse(reader)                  // 从 reader
registry, err := parser.ParseBytes(xmlData)            // 从字节
```

XML 解析器处理官方 MITRE CWE 目录格式，把所有条目类型转换为注册表条目。

---

## 🔧 CSV 序列化

### CWE 列表

```go
cwes := []*cweskills.CWE{cwe1, cwe2, cwe3}
data, err := cweskills.MarshalCSV(cwes)
parsed, err := cweskills.UnmarshalCSV(data)
```

### 注册表 CSV 导出

```go
csvData, err := registry.ExportCSV()
```

CSV 格式：

```csv
ID,Name,Abstraction,Status,Structure,Description,CWEType
79,Cross-site Scripting,Base,Stable,Simple,The product does not...,weakness
89,SQL Injection,Base,Stable,Simple,The product constructs...,weakness
```

::: warning CSV 有损
CSV 只保留基础字段（ID、Name、Abstraction、Status、Structure、Description、CWEType）。含逗号或引号的字段会正确转义；列数不足的行仍可解析（缺失字段默认为空）。
:::

---

## 📝 示例

### Go — 导出过滤结果为 CSV

```go
package main

import (
    "os"
    "github.com/scagogogo/cwe-skills"
)

func main() {
    registry, _ := cweskills.NewXMLParser().ParseFile("cwec_latest.xml")
    registry.BuildIndexes()

    filtered := cweskills.Filter(registry.GetAllCWEs(), cweskills.FilterOption{
        Abstraction: cweskills.AbstractionBase,
    })

    data, _ := cweskills.MarshalCSV(filtered)
    _ = os.WriteFile("base_weaknesses.csv", data, 0644)
}
```

### Go — 注册表 JSON 往返

```go
data, _ := registry.ExportJSON()
_ = os.WriteFile("cwe.json", data, 0644)

// 另一处再加载
newReg := cweskills.NewRegistry()
_ = newReg.ImportJSON(data)
newReg.BuildIndexes()
```

---

## 🤖 AI 代理使用提示

- 用户要「导出数据」时，AI 用 `cwe registry export`（CLI）或上述 SDK 函数。
- JSON 适合完整保留数据；CSV 适合表格工具（Excel）查看。
- 把序列化结果写入文件用 `os.WriteFile`，CLI 用 `--output-file`。

::: details 错误类型
- `ParseError`：反序列化时格式无效
- `ValidationError`：字段值无效
- `ImportJSON` 对畸形 JSON 返回错误
:::

::: tip JSON 处理 interface{} 字段
JSON 序列化能优雅处理 `interface{}` 字段（如 show 结果里的），不会因类型不确定而失败。
:::

---

## 📖 相关文档

- [技能 09 — 本地注册表](./09-local-registry)
- [SDK: MarshalJSON](../sdk/marshal-json) · [MarshalCSV](../sdk/marshal-csv) · [MarshalXML](../sdk/marshal-xml) · [ExportCSV](../sdk/export-csv) · [Registry JSON](../sdk/registry-json)
- [CLI: registry export](../cli/registry-export)
- [返回 Skills 总览](./)
