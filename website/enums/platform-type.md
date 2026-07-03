---
title: PlatformType 平台类型枚举
outline: [2, 3]
---

# 🖥️ PlatformType — 平台类型枚举

`PlatformType` 表示 CWE 适用平台的类别。CWE 会标注弱点适用的编程语言、操作系统、架构或技术，便于按平台筛选相关弱点。

## 🧬 类型定义

```go
type PlatformType string
```

## 📋 全部取值

| 值 | 常量名 | 含义 |
| --- | --- | --- |
| `"Language"` | `PlatformLanguage` | 编程语言（如 C、Java、Python） |
| `"Operating System"` | `PlatformOperatingSystem` | 操作系统（如 Linux、Windows） |
| `"Architecture"` | `PlatformArchitecture` | 处理器架构（如 x86、ARM） |
| `"Technology"` | `PlatformTechnology` | 技术栈/运行时（如 Android、iOS） |

```go
const (
	PlatformLanguage        PlatformType = "Language"
	PlatformOperatingSystem PlatformType = "Operating System"
	PlatformArchitecture    PlatformType = "Architecture"
	PlatformTechnology      PlatformType = "Technology"
)
```

::: tip 用 AllPlatformTypeValues 取完整列表
本枚举取值固定为上述四类「平台类别」。具体平台实例（C、Java、Linux 等）不在此枚举内，而是作为字符串保存在 `ApplicablePlatforms` 中。如需获取这四个类别常量，调用 `AllPlatformTypeValues()` 即可。
:::

## ✅ 通用方法（四件套）

| 方法 / 函数 | 签名 |
| --- | --- |
| `String` | `func (p PlatformType) String() string` |
| `IsValid` | `func (p PlatformType) IsValid() bool` |
| `ParsePlatformType` | `func ParsePlatformType(s string) (PlatformType, error)` |
| `AllPlatformTypeValues` | `func AllPlatformTypeValues() []PlatformType` |

```go
p, err := cweskills.ParsePlatformType("Language")
fmt.Println(p, err)                               // Language <nil>
fmt.Println(p.String())                           // Language
fmt.Println(cweskills.PlatformType("X").IsValid()) // false
fmt.Println(cweskills.AllPlatformTypeValues())    // [Language Operating System Architecture Technology]
```

::: warning 取值含空格
`"Operating System"` 含空格，`ParsePlatformType` 严格匹配原样。不要写成 `"OperatingSystem"` 或 `"OS"`。
:::

## 🔄 典型用法

```go
// 按平台类别归类某弱点的适用平台
for _, ap := range cwe.ApplicablePlatforms {
	switch ap.Type {
	case cweskills.PlatformLanguage:
		fmt.Println("语言:", ap.Name)
	case cweskills.PlatformOperatingSystem:
		fmt.Println("操作系统:", ap.Name)
	case cweskills.PlatformArchitecture:
		fmt.Println("架构:", ap.Name)
	case cweskills.PlatformTechnology:
		fmt.Println("技术:", ap.Name)
	}
}
```

## 💻 CLI 对应命令

```bash
cwe enum platform
```

输出全部合法取值，详见 [CLI enum platform](../cli/enum-platform)。

## 🔗 相关链接

- SDK 视角：[PlatformType 平台类型枚举](../sdk/enum-platform-type)
- 适用平台数据模型：[ApplicablePlatforms](../sdk/applicable-platforms)
- 源文件：[`enums.go`](https://github.com/scagogogo/cwe-skills/blob/main/enums.go)
