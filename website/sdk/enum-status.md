---
title: Status 状态枚举
outline: [2, 3]
---

# 📚 Status — 状态枚举

`Status` 表示 CWE 条目的成熟度与稳定性，共 6 个取值，从最稳定的 `Stable` 到已废弃的 `Deprecated`。

## 📋 类型与常量

```go
type Status string

const (
	StatusStable     Status = "Stable"
	StatusUsable     Status = "Usable"
	StatusDraft      Status = "Draft"
	StatusIncomplete Status = "Incomplete"
	StatusObsolete   Status = "Obsolete"
	StatusDeprecated Status = "Deprecated"
)
```

## 📝 常量说明

| 常量 | 值 | 说明 |
| --- | --- | --- |
| `StatusStable` | `"Stable"` | 所有重要元素已验证，不太可能显著变化 |
| `StatusUsable` | `"Usable"` | 已深入审查，关键元素已验证 |
| `StatusDraft` | `"Draft"` | 重要元素已填写，可能仍有问题或空缺 |
| `StatusIncomplete` | `"Incomplete"` | 并非所有重要元素都已填写，无质量保证 |
| `StatusObsolete` | `"Obsolete"` | 仍有效但不再相关，已被更新实体取代 |
| `StatusDeprecated` | `"Deprecated"` | 已从 CWE 移除，是重复或错误创建的 |

::: tip Stable vs Usable
两者都属「可用」区间。`Stable` 是最高成熟度，内容稳定；`Usable` 略低一档，已审查但可能仍有微调。生产环境通常接受这两者。
:::

## 🧩 四件套方法

| 方法/函数 | 签名 |
| --- | --- |
| `String` | `func (s Status) String() string` |
| `IsValid` | `func (s Status) IsValid() bool` |
| `ParseStatus` | `func ParseStatus(s string) (Status, error)` |
| `AllStatusValues` | `func AllStatusValues() []Status` |

::: warning Status 没有 Order 方法
`Status` 没有提供 `StatusOrder()` —— 六种状态之间没有严格的线性数值序（虽然语义上有从成熟到废弃的倾向，但 MITRE 未定义数值排序）。如需排序请自行定义。
:::

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	// 解析
	s, err := cweskills.ParseStatus("Deprecated")
	fmt.Println(s, err) // Deprecated <nil>

	// 校验
	fmt.Println(cweskills.StatusStable.IsValid())  // true
	fmt.Println(cweskills.Status("Foo").IsValid()) // false

	// 全部值
	for _, v := range cweskills.AllStatusValues() {
		fmt.Println(v)
	}

	// 业务判断：跳过废弃条目
	cwe := cweskills.NewCWE(79, "XSS")
	cwe.Status = cweskills.StatusStable
	if cwe.IsDeprecated() {
		fmt.Println("已废弃，跳过")
	} else {
		fmt.Println("可用")
	}
}
```

## 🎯 典型用途

<Badge type="tip" text="过滤" /> [`FindByStatus`](./find-by-status) 只保留稳定条目
<Badge type="info" text="守卫" /> `if cwe.IsDeprecated() { continue }` 跳过废弃
<Badge type="warning" text="统计" /> [`CountByStatus`](./count-by-status) 统计各状态分布

## ⚠️ 注意事项

::: warning Obsolete vs Deprecated
- `Obsolete`：条目本身仍有效，但已被更新的条目取代，不再相关。
- `Deprecated`：条目已从 CWE 移除，是重复或错误创建的。

两者都建议在生产中排除，但语义不同——`Obsolete` 可能有继任者，`Deprecated` 通常是重复定义。
:::

## 🔗 相关链接

- 字段归宿：`CWE.Status`、`Category.Status`、`View.Status`、`CompoundElement.Status`
- 便捷判断方法：`CWE.IsStable()`/`IsDeprecated()`，见 [CWE 类型判断方法](./cwe-type-methods)
- 概览：[enums.go](./enums)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
