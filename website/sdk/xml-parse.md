---
title: Parse ParseFile ParseBytes 解析方法
outline: [2, 3]
---

# 📥 Parse / ParseFile / ParseBytes — XML 解析方法

`XMLParser` 提供三个解析入口，分别面向文件路径、`io.Reader` 与字节切片，底层都汇聚到 `Parse(io.Reader)`。三者返回的 `*Registry` 已填入四类条目。

源文件：`xml_parser.go`。

## 📐 方法签名

```go
func (p *XMLParser) ParseFile(path string) (*Registry, error)
func (p *XMLParser) Parse(reader io.Reader) (*Registry, error)
func (p *XMLParser) ParseBytes(data []byte) (*Registry, error)
```

## ParseFile

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `path` | `string` | XML 文件路径 |

`path` 为空返回 `ValidationError(field="path")`；`os.Open` 失败返回 `ParseError`；随后委托 `Parse(file)`。

## Parse

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `reader` | `io.Reader` | XML 数据源 |

`reader` 为 `nil` 返回 `ValidationError(field="reader", value="nil")`。用 `xml.NewDecoder` 解码到 `xmlWeaknessCatalog`，失败返回 `ParseError`。成功后遍历四类节点填入新 `Registry`。

## ParseBytes

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `data` | `[]byte` | XML 字节切片 |

`len(data)==0` 返回 `ValidationError(field="data", value="empty")`，否则用内部 `byteReader` 包装后委托 `Parse`。

## 🔁 转换细节

`convert*` 方法会把 XML 字符串字段经对应枚举解析器转换：

- `Abstraction` → [`ParseAbstraction`](./enum-abstraction)
- `Structure` → [`ParseStructure`](./enum-structure)
- `Status` → [`ParseStatus`](./enum-status)
- `LikelihoodOfExploit` → [`ParseLikelihoodOfExploit`](./enum-likelihood)
- `Relationship.Nature` → [`ParseRelationshipNature`](./enum-relationship-nature)
- `View.Type` → `ParseViewType`

枚举解析失败时以零值兜底，不中断整体流程。每条 `CWE` 的 `URL` 字段自动生成为 `https://cwe.mitre.org/data/definitions/{ID}.html`。

::: tip 重复注册被静默忽略
`registry.Register` 对重复 ID 返回错误，但 `Parse` 用 `_ =` 丢弃，确保即使 XML 内有重复条目也能整体解析完成。
:::

## 🚀 可运行示例

```go
package main

import (
    "fmt"
    "log"
    "strings"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    parser := cwe.NewXMLParser()

    // 1. 从文件解析
    registry, err := parser.ParseFile("cwec_v4.10.xml")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("文件解析条目数:", registry.Size())

    // 2. 从 io.Reader 解析
    xmlStream := `<Weakness_Catalog Version="4.10">
  <Weaknesses><Weakness ID="89" Name="SQLi" Abstraction="Variant" Structure="Simple" Status="Stable">
    <Description>SQL Injection</Description>
  </Weakness></Weaknesses>
</Weakness_Catalog>`
    r2, err := parser.Parse(strings.NewReader(xmlStream))
    if err != nil {
        log.Fatal(err)
    }
    if w, ok := r2.Get(89); ok {
        fmt.Println("Reader 解析:", w.Name)
    }

    // 3. 从字节解析
    r3, err := parser.ParseBytes([]byte(xmlStream))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("字节解析条目数:", r3.Size())
}
```

::: warning Registry.Size() 取决于实现
上例假设 `Registry` 暴露 `Size()` 与 `Get(int)`；若当前版本方法名不同，请以 [Registry 模型](./model) 文档为准调整。
:::

## 📚 相关链接

- [XMLParser 概览](./xml-parser) | [NewXMLParser](./new-xml-parser) | [枚举总览](./enums) | [ValidationError](./validation-error)
