---
title: 注册表 JSON 导入导出
outline: [2, 3]
---

# 🗃️ 注册表 JSON 导入导出

`Registry` 提供整库级别的 JSON 序列化：`ExportJSON` 把全部弱点写成 JSON 字节流，`ImportJSON` 反向重建注册表。适合做本地缓存、快照备份或跨进程数据传递。

## 📤 ExportJSON

```go
func (r *Registry) ExportJSON() ([]byte, error)
```

| 项 | 说明 |
| --- | --- |
| 参数 | 无 |
| 返回 | `[]byte` 全部弱点的 JSON 数组；`error` 序列化错误 |
| 范围 | 仅导出 `CWE`（弱点），不含 Category/View/CompoundElement |

## 📥 ImportJSON

```go
func (r *Registry) ImportJSON(data []byte) error
```

| 项 | 说明 |
| --- | --- |
| 参数 | `data` 由 `ExportJSON` 产生的 JSON 字节流 |
| 返回 | `error` 解析或注册错误 |
| 行为 | 解析为 `[]*CWE` 后逐个 `Register`，**并自动重建索引** |

::: warning ImportJSON 会覆盖吗？
`ImportJSON` 内部调用 `Register`，遇到已存在的 ID 会返回错误并中止。如需「覆盖式」导入，先调用 `Clear()` 清空。
:::

::: tip 导入后索引自动就绪
`ImportJSON` 末尾会调用 `BuildIndexes()`，因此导入完成后可直接使用 [关系索引查询](./relationship-indexes)，无需手动重建。
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
	r := cweskills.NewRegistry()
	_ = r.Register(cweskills.NewCWE(79, "XSS"))
	_ = r.Register(cweskills.NewCWE(89, "SQL Injection"))

	// 导出
	data, err := r.ExportJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("导出 %d 字节\n", len(data))

	// 重建到新注册表
	r2 := cweskills.NewRegistry()
	if err := r2.ImportJSON(data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(r2.Size())         // 2
	fmt.Println(r2.IndexesBuilt()) // true（自动重建）
}
```

## 🆚 与单条目序列化的区别

| 维度 | 本组（Registry） | [Serializer](./serializer) |
| --- | --- | --- |
| 粒度 | 整库全部弱点 | 单条 `*CWE` 或 `[]*CWE` |
| 格式 | 仅 JSON | JSON / XML / CSV |
| 索引 | 导入自动重建 | 不涉及 |
| 适用 | 快照、缓存 | 单条目持久化、CSV 报表 |

## 🔗 相关链接

- 单条目 JSON：[MarshalJSON / UnmarshalJSON](./marshal-json)
- 列表 JSON：[MarshalJSONList / UnmarshalJSONList](./marshal-json-list)
- CSV 导出：[ExportCSV](./export-csv)
- 源文件：[`registry.go`](https://github.com/scagogogo/cwe-skills/blob/main/registry.go)
