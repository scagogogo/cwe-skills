# 📦 包与版本

<Badge type="info" text="cwe.go" /> <Badge type="warning" text="v0.0.1" />

`cwe.go` 是 SDK 的包级文档文件，定义了版本常量与包的整体说明。

## 📋 版本常量

```go
// Package cweskills 提供了对CWE（Common Weakness Enumeration，通用缺陷枚举）的完整支持
package cweskills

// Version 表示本SDK的版本号
const Version = "v0.0.1"
```

### `Version`

<Badge type="warning" text="常量" />

```go
const Version = "v0.0.1"
```

SDK 的语义化版本号字符串。

| 属性 | 值 |
|------|-----|
| 类型 | `string` |
| 当前值 | `"v0.0.1"` |
| 用途 | User-Agent 标识、API 客户端标识、CLI 版本显示 |

::: tip 用法
`Version` 被 `HTTPClient` 用作默认 UserAgent 的一部分（`cwe-sdk-go/` + Version），也用于 CLI `cwe version` 命令的 SDK 版本输出。
:::

## 📖 包文档

`cwe.go` 顶部的包注释概述了 SDK 的完整能力范围：

- CWE ID 的格式化、解析、验证和提取
- 完整的枚举类型定义（抽象层级、状态、关系类型等）
- 结构化错误类型
- 知名 CWE 列表（Top 25、OWASP Top 10 等）
- 核心数据模型（Weakness、Category、View、CompoundElement）
- 关系导航（父/子/祖先/后代/链/组合）
- 内存注册表与索引
- 搜索与过滤
- MITRE CWE REST API 客户端
- XML 目录解析（离线模式）
- JSON/XML/CSV 序列化
- 速率限制的 HTTP 客户端

## 🎯 设计目标

::: tip 零依赖核心
核心 SDK **仅使用 Go 标准库**，无任何第三方依赖。这使得：
- 编译产物小，可静态链接
- 供应链审计友好，无第三方漏洞面
- 构建快速，无依赖解析开销
:::

::: warning CLI 例外
CLI 工具（`cmd/cwe`）使用 `github.com/spf13/cobra`，但这是 CLI 层的依赖，不影响 SDK 本身的零依赖特性。
:::

## 🔗 相关

- [SDK 总览](./overview)
- [HTTP 客户端](./http-client)（使用 `Version` 作 UserAgent）
- [CLI version 命令](../cli/version)（显示 `Version`）
- [更新日志](../changelog)
