---
title: JSON 列表序列化
outline: [2, 3]
---

# 📦 JSON 列表序列化

`MarshalJSONList` 与 `UnmarshalJSONList` 处理 `[]*CWE` 与 JSON 数组之间的转换。适合批量持久化、跨进程传递弱点集合。

## 📐 函数签名

### MarshalJSONList

```go
func MarshalJSONList(cwes []*CWE) ([]byte, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `cwes` | `[]*CWE` | 待序列化的弱点列表 |
| 返回 1 | `[]byte` | JSON 数组字节流 |
| 返回 2 | `error` | 序列化错误 |

### UnmarshalJSONList

```go
func UnmarshalJSONList(data []byte) ([]*CWE, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `data` | `[]byte` | JSON 数组字节流 |
| 返回 1 | `[]*CWE` | 反序列化得到的弱点列表 |
| 返回 2 | `error` | 解析错误 |

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

	data, err := cweskills.MarshalJSONList(list)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)

	got, err := cweskills.UnmarshalJSONList(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("数量:", len(got)) // 2
	for _, c := range got {
		fmt.Println(c.CWEID(), c.Name)
	}
}
```

## ⚠️ 注意事项

::: tip 与单条目的关系
`MarshalJSONList` 本质是逐个用 `safeCWE` 转换后拼成 JSON 数组。等价于 `[{单条JSON}, {单条JSON}, ...]`。
:::

::: warning 不注册到 Registry
`UnmarshalJSONList` 只返回 `[]*CWE`，**不**把它们注册到任何 `Registry`。如需入注册表，逐个调用 `Register`，或改用 [Registry.ImportJSON](./registry-json)（它会自动注册并重建索引）。
:::

## 🆚 与 Registry.ExportJSON 的区别

| 维度 | 本组 | [`Registry.ExportJSON`](./registry-json) |
| --- | --- | --- |
| 输入 | 任意 `[]*CWE` | 注册表全部弱点 |
| 输出后 | 仅字节流 | 可直接 `ImportJSON` 重建 |
| 适用 | 部分弱点导出 | 整库快照 |

## 🔗 相关链接

- 单条版本：[JSON 序列化](./marshal-json)
- 整库版本：[注册表 JSON](./registry-json)
- CSV 列表：[CSV 序列化](./marshal-csv)
- 源文件：[`serializer.go`](https://github.com/scagogogo/cwe-skills/blob/main/serializer.go)
