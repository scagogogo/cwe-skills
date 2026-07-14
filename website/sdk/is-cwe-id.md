---
title: IsCWEID 判断是否合法 CWE ID
outline: [2, 3]
---

# 🆔 IsCWEID — 判断字符串是否为合法 CWE ID

## 📋 函数签名

```go
func IsCWEID(text string) bool
```

## 📖 说明

`IsCWEID` 是一个**轻量判断函数**：它只回答「格式上是不是合法的 CWE ID」，返回布尔值。内部实现等价于：

```go
func IsCWEID(text string) bool {
    _, err := ParseCWEID(text)
    return err == nil
}
```

::: warning 仅校验格式，不查存在性
本函数**不会**查询 MITRE 数据库确认该编号是否真实存在。`IsCWEID("CWE-99999")` 返回 `true`，但 99999 并非真实弱点。需要确认存在性请配合 [Registry](./registry) 或 [API 客户端](./api-client)。
:::

## 📥 参数

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `text` | `string` | 待检查的字符串 |

## 📤 返回值

| 返回值 | 类型 | 说明 |
| --- | --- | --- |
| 合法性 | `bool` | 合法返回 `true`，否则 `false` |

## ✅ 示例

```go
package main

import (
	"fmt"
	"github.com/scagogogo/cwe-skills"
)

func main() {
	fmt.Println(cweskills.IsCWEID("CWE-79"))  // true
	fmt.Println(cweskills.IsCWEID("79"))      // true
	fmt.Println(cweskills.IsCWEID("cwe-079")) // true
	fmt.Println(cweskills.IsCWEID("abc"))     // false
	fmt.Println(cweskills.IsCWEID(""))        // false
	fmt.Println(cweskills.IsCWEID("CWE-0"))   // false （必须 > 0）

	// 典型用法：守卫用户输入
 userInput := "CWE-79"
	if !cweskills.IsCWEID(userInput) {
		fmt.Println("请输入合法的 CWE ID")
		return
	}
}
```

## ⚖️ IsCWEID vs ValidateCWEID

<Badge type="tip" text="二选一" />

| 维度 | `IsCWEID` | [`ValidateCWEID`](./validate-cwe-id) |
| --- | --- | --- |
| 返回值 | `bool` | `error` |
| 是否给出原因 | ❌ 否 | ✅ 是，带 `InvalidCWEIDError` |
| 适用场景 | 守卫判断、过滤 | 需向用户/日志输出错误原因 |
| 实现差异 | 仅判空 + `ParseCWEID` | 同样调用 `ParseCWEID`，但保留 error |

::: tip 选择建议
- 只需要 yes/no → `IsCWEID`，代码更简洁。
- 需要把错误透传给上层或展示给用户 → `ValidateCWEID`。
:::

## ⚠️ 常见错误

::: warning 不要用它做存在性校验
`IsCWEID("CWE-42")` 为 `true` 不代表 CWE-42 真实存在。在需要真实性的业务里，请再用 `Registry.Contains(42)` 或 `APIClient.GetWeakness(42)` 二次确认。
:::

## 🔗 相关链接

- 底层依赖：[ParseCWEID](./parse-cwe-id)
- 需要错误原因时：[ValidateCWEID](./validate-cwe-id)
- 源文件：[`cwe_utils.go`](https://github.com/scagogogo/cwe-skills/blob/main/cwe_utils.go)
