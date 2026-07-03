# 📚 cwe enum platform

列出 CWE 适用平台类型（PlatformType）的合法取值。

::: danger 当前 CLI 未注册此子命令
经核查 `cmd/cwe/enums.go`，当前版本的 `cwe enum` 父命令动态生成的子命令为：`abstraction`、`structure`、`status`、`likelihood`、`relationship`、`scope`、`impact`、`viewtype`，**未包含 `platform`**。直接执行 `cwe enum platform` 会报“unknown command”错误。
:::

平台类型枚举在 SDK 层面存在（`AllPlatformTypeValues()`），值包括：

| 取值 | 含义 |
| --- | --- |
| `Language` | 编程语言相关平台 |
| `Operating System` | 操作系统相关平台 |
| `Architecture` | 硬件架构相关平台 |
| `Technology` | 技术栈相关平台 |

## 语法（如未来启用）

```bash
cwe enum platform [flags]
```

## 替代获取方式

由于 CLI 暂未注册该子命令，可通过以下方式获取平台类型取值：

### 方式一：调用 SDK

```go
package main

import (
	"fmt"
	cweskills "github.com/scagogogo/cwe-skills"
)

func main() {
	for _, p := range cweskills.AllPlatformTypeValues() {
		fmt.Println(p)
	}
}
```

### 方式二：参考本文表格

直接使用上表中的四个取值：`Language`、`Operating System`、`Architecture`、`Technology`。

## 使用场景

- 解读弱点 `ApplicablePlatforms` 字段中的平台分类。
- 按平台维度筛选/归类弱点。
- 自定义平台相关分析时确认合法取值。

::: tip 期待支持
若希望 CLI 提供 `cwe enum platform` 子命令，可参考 [SDK ApplicablePlatform](../sdk/applicable-platforms) 与 [PlatformType 枚举](../sdk/enum-platform-type) 自行扩展，或向项目提 issue。
:::

## 下一步

- [enum](./enum) — 查看当前已注册的全部枚举子命令。
- [enum view-type](./enum-view-type) — 另一个枚举示例。

## 相关文档

- [SDK PlatformType 枚举](../sdk/enum-platform-type)
- [SDK ApplicablePlatform](../sdk/applicable-platforms)
