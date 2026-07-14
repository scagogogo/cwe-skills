---
title: 按 ID 查找
outline: [2, 3]
---

# 🔍 按 ID 查找

`FindByID` 是最直接的查找方式：给定弱点 ID，返回对应的 `*CWE`。它本质是 `Registry.Get` 的函数式封装。

## 📐 函数签名

```go
func FindByID(r *Registry, id int) (*CWE, bool)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `r` | `*Registry` | 注册表 |
| `id` | `int` | 弱点 ID |
| 返回 1 | `*CWE` | 找到的弱点；未找到为 `nil` |
| 返回 2 | `bool` | 是否找到 |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	r := cweskills.NewRegistry()
	_ = r.Register(cweskills.NewCWE(79, "XSS"))
	_ = r.Register(cweskills.NewCWE(89, "SQLi"))

	if c, ok := cweskills.FindByID(r, 79); ok {
		fmt.Println(c.CWEID(), c.Name) // CWE-79 XSS
	}
	if _, ok := cweskills.FindByID(r, 9999); !ok {
		fmt.Println("不存在") // 不存在
	}
}
```

## 🆚 与 Registry.Get 的关系

| 维度 | `FindByID` | [`Registry.Get`](./registry-operations) |
| --- | --- | --- |
| 形式 | 包级函数 | 方法 |
| 签名 | `(*CWE, bool)` | `(*CWE, bool)` |
| 行为 | 完全等价 | 完全等价 |

`FindByID` 提供统一的函数式 API，便于与其它 `FindBy*` 函数在管道中混用。

## ⚠️ 注意事项

::: tip 不依赖索引
`FindByID` 直接查注册表内部 map，无需 `BuildIndexes()`，O(1) 复杂度。
:::

## 🔗 相关链接

- 方法版本：[Registry.Get](./registry-operations)
- 关键词查找：[按关键词查找](./find-by-keyword)
- 源文件：[`search.go`](https://github.com/scagogogo/cwe-skills/blob/main/search.go)
