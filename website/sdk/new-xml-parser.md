---
title: NewXMLParser 创建解析器
outline: [2, 3]
---

# 🔧 NewXMLParser — 创建 XML 解析器

`NewXMLParser` 是构造 `XMLParser` 的唯一入口。解析器无状态、无配置，零参数即可使用。

源文件：`xml_parser.go`。

## 📐 函数签名

```go
func NewXMLParser() *XMLParser
```

无参数。返回值：`*XMLParser`，零字段空结构体。

::: tip 为何没有配置选项
`XMLParser` 完全由 MITRE XML Schema 的固定结构驱动，没有可调参数。所有解析行为（字段映射、枚举转换、关系填充）都在内部 `convert*` 方法里固化。需要定制时，解析后直接操作返回的 `Registry` 即可。
:::

## 🚀 可运行示例

```go
package main

import (
    "bytes"
    "fmt"
    "log"

    cwe "github.com/scagogogo/cwe-skills"
)

func main() {
    parser := cwe.NewXMLParser()

    // 也可从字节切片解析（这里用最小占位 XML）
    xmlData := []byte(`
<Weakness_Catalog Name="Example" Version="4.10" Date="2024-01-01">
  <Weaknesses>
    <Weakness ID="79" Name="XSS" Abstraction="Variant" Structure="Simple" Status="Stable">
      <Description>Improper Neutralization</Description>
    </Weakness>
  </Weaknesses>
</Weakness_Catalog>`)

    registry, err := parser.Parse(new(bytes.Reader), nil)
    _ = err // 这里仅演示构造，实际用 ParseBytes
    registry, err = parser.ParseBytes(xmlData)
    if err != nil {
        log.Fatal(err)
    }
    if w, ok := registry.Get(79); ok {
        fmt.Println(w.Name)
    }
}
```

::: warning 上例中的 Parse 调用
为展示 `NewXMLParser` 返回值的多个方法，示例中先调用了一次 `Parse(new(bytes.Reader), nil)`——这会因 `nil` reader 触发 `ValidationError`。实际使用时请直接调 `ParseFile` 或 `ParseBytes`，详见 [Parse 方法](./xml-parse)。
:::

## 📚 相关链接

- [XMLParser 概览](./xml-parser) | [Parse / ParseFile](./xml-parse) | [Registry 模型](./model)
