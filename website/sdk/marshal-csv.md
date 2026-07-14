---
title: CSV 序列化
outline: [2, 3]
---

# 📦 CSV 序列化

`MarshalCSV` 与 `UnmarshalCSV` 处理 `[]*CWE` 与 CSV 字节流之间的转换。CSV 格式扁平、列固定，适合导入电子表格或做简报。

## 📐 函数签名

### MarshalCSV

```go
func MarshalCSV(cwes []*CWE) ([]byte, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cwes` | `[]*CWE` | 待序列化的弱点列表 |
| 返回 1 | `[]byte` | CSV 字节流（含表头行） |
| 返回 2 | `error` | 序列化错误 |

### UnmarshalCSV

```go
func UnmarshalCSV(data []byte) ([]*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `data` | `[]byte` | CSV 字节流 |
| 返回 1 | `[]*CWE` | 反序列化得到的弱点列表 |
| 返回 2 | `error` | 解析错误 |

## 📋 固定表头

CSV 使用固定 7 列表头，顺序如下：

| 列序 | 字段 |
| --- | --- |
| 1 | `ID` |
| 2 | `Name` |
| 3 | `Abstraction` |
| 4 | `Structure` |
| 5 | `Status` |
| 6 | `Description` |
| 7 | `LikelihoodOfExploit` |

```go
var csvHeader = []string{"ID", "Name", "Abstraction", "Structure", "Status", "Description", "LikelihoodOfExploit"}
```

::: warning 字段子集
CSV 只导出上述 7 个字段，`CommonConsequences`/`References` 等嵌套字段**不**导出。需要完整数据请用 [JSON](./marshal-json) 或 [XML](./marshal-xml)。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"log"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	list := []*cweskills.CWE{
		cweskills.NewCWE(79, "XSS"),
		cweskills.NewCWE(89, "SQLi"),
	}

	data, err := cweskills.MarshalCSV(list)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)

	got, err := cweskills.UnmarshalCSV(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("数量:", len(got)) // 2
}
```

## ⚠️ 注意事项

::: tip Description 含逗号
CSV 写入时对含逗号/引号/换行的字段做标准 CSV 转义（加引号、双写引号）。读取时按标准规则还原，调用方无需手动处理。
:::

::: details 枚举的 CSV 表示
`Abstraction`/`Structure`/`Status`/`LikelihoodOfExploit` 在 CSV 中以枚举的字符串名输出。反序列化时按名解析回枚举，非法名报错。
:::

## 🔗 相关链接

- 整库 CSV：[ExportCSV](./export-csv)
- JSON 列表：[JSON 列表](./marshal-json-list)
- 表头定义：[`serializer.go`](https://github.com/scagogogo/cwe-skills/blob/main/serializer.go)
- 源文件：[`serializer.go`](https://github.com/scagogogo/cwe-skills/blob/main/serializer.go)
