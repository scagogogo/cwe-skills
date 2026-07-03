---
title: 注册表导出 CSV
outline: [2, 3]
---

# 📦 注册表导出 CSV

`Registry.ExportCSV` 是整库级别的 CSV 导出方法——把注册表中全部弱点一次性写成 CSV 字节流。它是 [MarshalCSV](./marshal-csv) 的整库封装。

## 📐 方法签名

```go
func (r *Registry) ExportCSV() ([]byte, error)
```

| 项 | 说明 |
| --- | --- |
| 接收者 | `*Registry` |
| 参数 | 无 |
| 返回 1 | `[]byte` CSV 字节流（含表头 + 全部弱点行） |
| 返回 2 | `error` 序列化错误 |
| 范围 | 仅导出 `CWE`（弱点），按 [固定 7 列表头](./marshal-csv) |

## ✅ 示例

```go
package main

import (
	"fmt"
	"log"
	"os"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	r.Register(cweskills.NewCWE(79, "XSS"))
	r.Register(cweskills.NewCWE(89, "SQLi"))

	data, err := r.ExportCSV()
	if err != nil {
		log.Fatal(err)
	}

	// 写入文件
	if err := os.WriteFile("cwes.csv", data, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("导出 %d 字节，%d 行弱点\n", len(data), r.Size())
}
```

## 🆚 与 MarshalCSV 的关系

| 维度 | `Registry.ExportCSV` | [`MarshalCSV`](./marshal-csv) |
| --- | --- | --- |
| 输入 | 无（内部取全部） | 显式 `[]*CWE` |
| 范围 | 整库弱点 | 任意子集 |
| 适用 | 全库备份 | 部分导出、自定义列表 |

`ExportCSV` 等价于 `MarshalCSV(r.GetAll())`，但语义更明确。

## ⚠️ 注意事项

::: warning 仅有导出，无导入
`Registry` 提供 `ExportCSV` 但**不**提供 `ImportCSV`。如需从 CSV 导入，用包级 [UnmarshalCSV](./marshal-csv) 解析后逐个 `Register`。
:::

::: tip 不依赖索引
`ExportCSV` 遍历 `GetAll()`，无需 `BuildIndexes()`。但输出行顺序不确定，需排序时先导出再按 ID 排序。
:::

## 🔗 相关链接

- 包级版本：[CSV 序列化](./marshal-csv)
- 表头定义：[CSV 序列化](./marshal-csv)
- 整库 JSON：[注册表 JSON](./registry-json)
- 源文件：[`registry.go`](https://github.com/scagogogo/cwe-skills/blob/main/registry.go)、[`serializer.go`](https://github.com/scagogogo/cwe-skills/blob/main/serializer.go)
