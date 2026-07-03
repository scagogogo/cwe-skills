---
title: XML 序列化
outline: [2, 3]
---

# 📦 XML 序列化

`MarshalXML` 与 `UnmarshalXML` 完成单个 `CWE` 与 XML 字节流之间的转换。XML 格式与 MITRE 官方 CWE XML 的元素命名对齐，便于与官方数据互换。

## 📐 函数签名

### MarshalXML

```go
func MarshalXML(cwe *CWE) ([]byte, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cwe` | `*CWE` | 待序列化的弱点 |
| 返回 1 | `[]byte` | XML 字节流 |
| 返回 2 | `error` | 序列化错误 |

### UnmarshalXML

```go
func UnmarshalXML(data []byte) (*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `data` | `[]byte` | XML 字节流 |
| 返回 1 | `*CWE` | 反序列化得到的弱点 |
| 返回 2 | `error` | 解析错误 |

::: tip 与官方 XML 的关系
`CWE` 结构体的 `xml` tag（如 `xml:"Name"`、`xml:"CommonConsequences>Consequence"`）与 MITRE 官方 CWE XML schema 的元素命名一致，因此 `MarshalXML` 产出的 XML 可被官方工具识别。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"log"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.Abstraction = cweskills.AbstractionBase
	cwe.Description = "Cross-site Scripting"

	data, err := cweskills.MarshalXML(cwe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)

	got, err := cweskills.UnmarshalXML(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got.CWEID(), got.Name) // CWE-79 XSS
}
```

## ⚠️ 注意事项

::: warning 仅单条目
本组函数只处理**单个** `CWE`，不提供列表/整库 XML 序列化。如需批量 XML，需自行循环调用或用官方 XML parser，参见 [XML Parser](./xml-parser)。
:::

::: details 枚举的 XML 表示
枚举字段通过其 `MarshalText`/`UnmarshalText` 在 XML 中以文本节点形式出现。反序列化时非法文本会导致错误。
:::

## 🔗 相关链接

- JSON 版本：[JSON 序列化](./marshal-json)
- 官方 XML 解析：[XML Parser](./xml-parser)
- 数据模型 tag：[CWE 结构体](./cwe-struct)
- 源文件：[`serializer.go`](https://github.com/scagogogo/cwe-skills/blob/main/serializer.go)
