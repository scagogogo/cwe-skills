---
title: JSON 序列化
outline: [2, 3]
---

# 📦 JSON 序列化

`MarshalJSON` 与 `UnmarshalJSON` 完成单个 `CWE` 与 JSON 字节流之间的转换。内部通过 `safeCWE` 安全模型中转，输出字段命名稳定。

## 📐 函数签名

### MarshalJSON

```go
func MarshalJSON(cwe *CWE) ([]byte, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cwe` | `*CWE` | 待序列化的弱点 |
| 返回 1 | `[]byte` | JSON 字节流 |
| 返回 2 | `error` | 序列化错误 |

### UnmarshalJSON

```go
func UnmarshalJSON(data []byte) (*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `data` | `[]byte` | JSON 字节流 |
| 返回 1 | `*CWE` | 反序列化得到的弱点 |
| 返回 2 | `error` | 解析错误 |

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
	cwe.Status = cweskills.StatusStable
	cwe.Description = "Cross-site Scripting"

	data, err := cweskills.MarshalJSON(cwe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)

	got, err := cweskills.UnmarshalJSON(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got.CWEID(), got.Name, got.Abstraction) // CWE-79 XSS Base
}
```

## ⚠️ 注意事项

::: tip 与 encoding/json 的关系
本组函数内部用 `safeCWE` 做字段映射后调用标准库 `json.Marshal`/`json.Unmarshal`。不应直接对 `CWE` 用 `json.Marshal`——那样会绕过安全模型，字段名可能不一致。
:::

::: warning 枚举解析
反序列化时，`Abstraction`/`Status` 等枚举字段需为合法字符串名。非法值会导致 `UnmarshalJSON` 返回错误。
:::

## 🆚 与 Registry.ExportJSON 的区别

| 维度 | 本组 | [`Registry.ExportJSON`](./registry-json) |
| --- | --- | --- |
| 粒度 | 单条 `*CWE` | 全库弱点 |
| 输出 | 单对象 JSON | JSON 数组 |
| 索引 | 不涉及 | 导入自动重建 |

## 🔗 相关链接

- 列表版本：[JSON 列表](./marshal-json-list)
- 整库版本：[注册表 JSON](./registry-json)
- XML 版本：[XML 序列化](./marshal-xml)
- 源文件：[`serializer.go`](https://github.com/scagogogo/cwe-skills/blob/main/serializer.go)
