---
title: 常见问题 FAQ
outline: [2, 3]
---

# ❓ 常见问题 FAQ

本页收录 CWE Skills 的常见问题与解答。涵盖安装、数据、API、版本、贡献、商业使用等。如果你的问题不在此列，欢迎到 [GitHub Issues](https://github.com/scagogogo/cwe-skills/issues) 提问。

---

## 📥 数据与下载

### Q1：怎么下载离线 XML 弃点目录？

访问 MITRE 官方下载页 <https://cwe.mitre.org/data/downloads.html>，下载最新版（如 `cwec_v4.15.xml`，文件较大，几百 MB）。下载后放到任意路径，CLI 用 `--xml <file>` 指定，SDK 用 `ParseFile("<file>")`。

```bash
# 验证可解析
cwe stats --xml cwec_v4.15.xml
```

详见 [安装](./installation) 与 [在线 vs 离线](./online-offline)。

::: details 一定要下最新版吗？
不必。离线 XML 是快照，旧版也能用，只是数据稍旧。MITRE 每年发布新版，按需更新即可。生产环境建议锁定一个版本，定期升级。
:::

---

### Q2：在线 API 报 429（Too Many Requests）怎么办？

429 表示触发了 MITRE 的速率限制。SDK 内置令牌桶（默认约 0.1 req/s）与自动重试，正常使用很少触发。若仍触发：

1. **降低请求频率**：调小 `WithAPIRateLimit` 的 rate。
2. **等待重试**：SDK 会读取 `Retry-After` 头自动等待，详见 [速率限制与重试](./rate-limit-retry)。
3. **改用离线**：大批量场景请走 XML，无任何限制。
4. **多进程共享 IP 慎用高 rate**：令牌桶是进程内的，多进程会叠加。详见 [速率限制与重试](./rate-limit-retry)。

```go
if errors.As(err, &rlErr); rlErr != nil {
    time.Sleep(rlErr.RetryAfter) // 等待后重试
}
```

---

### Q3：`CWE-79` 和 `79` 有什么区别？

没有本质区别。`CWE-79` 是标准显示形式，`79` 是内部整数形式。SDK 用 `ParseCWEID` 解析、`FormatCWEID` 格式化，二者互通：

```go
num, _ := cweskills.ParseCWEID("CWE-79")      // 79（整数）
formatted, _ := cweskills.FormatCWEID("79")   // "CWE-79"（字符串）
cweskills.FormatCWEIDFromInt(79)               // "CWE-79"
cweskills.IsCWEID("cwe79")                     // true（宽松识别）
```

::: tip 大小写与前缀都宽松
`ParseCWEID` / `IsCWEID` 对大小写和前缀很宽松：`CWE-79`、`cwe-79`、`CWE79`、`cwe79` 都能识别。但 `FormatCWEID` 统一输出 `CWE-79` 标准形式。详见 [CWE ID 工具](../sdk/cwe-utils)。
:::

---

## 🧰 安装与使用

### Q4：`cwe` 命令提示 `command not found`？

二进制所在目录不在 `PATH` 中。

- 用 curl 安装的：`/usr/local/bin` 通常在 PATH，确认 `mv` 到了那里。
- 用 `go install` 的：二进制在 `$GOPATH/bin`（通常 `$HOME/go/bin`），把该目录加入 `PATH`：

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

详见 [安装](./installation)。

---

### Q5：macOS 提示「无法验证开发者」或「无法打开」？

macOS Gateware 对未签名二进制隔离。解除：

```bash
sudo xattr -d com.apple.quarantine /usr/local/bin/cwe
```

或右键 → 打开 → 仍要打开。

---

### Q6：在线调用报 TLS 错误？

系统缺少 CA 证书。Linux 装 `ca-certificates`：

```bash
sudo apt-get install -m ca-certificates   # Debian/Ubuntu
sudo yum install ca-certificates           # RHEL/CentOS
```

Docker 镜像里也要装（`RUN apk add --no-cache ca-certificates`）。SDK 里可用 `WithHTTPClient` 传入自定义 `*http.Client` 调整 TLS 配置。详见 [安装](./installation)。

---

## 🌐 在线 vs 离线

### Q7：为什么导航/建树/搜索都需要 `--xml`，不能在线做？

因为 MITRE REST API **只返回部分关系类型**（主要父子层级），而导航、建树、最短路径需要**全部 10 种关系类型**（含链式、依赖、对等、复合）。这些只有离线 XML 才齐全。所以这类命令必须走离线路径，带 `--xml`。详见 [在线 vs 离线](./online-offline)。

::: info 这是 MITRE API 的限制
不是 CWE Skills 的限制，而是 MITRE REST API 本身不返回完整关系。要完整关系只能用 XML。
:::

---

### Q8：能同时用在线和离线吗？

可以，且推荐。**在线查详情（新鲜），离线做导航（完整）**是最佳实践。两条路径数据互通，都是同一个 `*CWE` 类型。见 [在线 vs 离线](./online-offline) 的「混合用法」。

```go
weakness, _ := client.GetWeakness(ctx, 79)   // 在线详情
registry, _ := cweskills.NewXMLParser().ParseFile("cwec_v4.15.xml")
registry.BuildIndexes()
nav := cweskills.NewNavigator(registry)
nav.Ancestors(79)                             // 离线完整祖先链
```

---

## 🧩 概念与版本

### Q9：CWE Skills 跟 MITRE CWE 是什么关系？

CWE Skills 是 MITRE CWE 的**集成层**，不是 CWE 本身。它不定义 CWE 内容（内容由 MITRE 维护），而是提供把 MITRE 的 API、XML、权威列表统一成一套类型化 Go 对象的「胶水」，外加 CLI / Skills / MCP 等接入方式。详见 [CWE Skills 是什么](./what-is-cwe-skills)。

---

### Q10：CWE 版本兼容性如何？

- **SDK/API**：MITRE REST API 有自己的版本（`cwe version mitre`），向后兼容较好。
- **离线 XML**：每个版本（如 v4.15）是快照，CWE ID 跨版本通常稳定（已分配的 ID 不变），但新增/废弃条目会有变化。
- **SDK 自身**：遵循语义化版本（`cweskills.Version`，当前 v0.0.1）。v0.x 阶段 API 可能有变动。

::: warning v0.x 阶段
当前版本号 v0.0.1 表示还在早期开发阶段，API 可能在小版本间调整。生产环境建议锁定版本（`go get ...@v0.0.1`），升级时看 [更新日志](../changelog)。
:::

---

## 💼 许可与贡献

### Q11：可以商业使用吗？

可以。CWE Skills 基于 **MIT 许可证**发布，允许商业使用、修改、分发、私有化，只需保留版权声明。MITRE CWE 数据本身也有其许可（见 MITRE 官网），使用 XML 数据时请同时遵守 MITRE 的许可条款。

::: tip MIT 许可证最宽松
MIT 是最宽松的开源许可之一，无 copyleft 限制，商业友好。你可以在闭源商业产品里用 CWE Skills。
:::

---

### Q12：如何贡献代码或文档？

1. Fork 仓库 [github.com/scagogogo/cwe-skills](https://github.com/scagogogo/cwe-skills)。
2. 新建分支：`git checkout -b feature/my-feature`。
3. 改代码/文档，跑测试：`go test ./...`。
4. 提 PR，描述清楚动机与变更。

::: details 贡献前的准备
- 阅读 [工作原理](./how-it-works) 理解架构。
- 核心 SDK 保持零依赖——不要引入第三方包（CLI 的 cobra 例外）。
- 新增能力需同时覆盖 SDK / CLI / 文档，保持四种接入方式能力对等。
- 遵循现有代码风格与文档风格（简体中文 + emoji 标题 + VitePress 容器）。
:::

---

## 🤖 AI 与 Skills

### Q13：Skills 跟 MCP 有什么区别？该用哪个？

- **Skills**（当前可用）：给 AI 一段提示词，让 AI 跑 `cwe` CLI。依赖 AI 能执行 Shell。
- **MCP**（规划中）：给 AI 一个标准化工具服务器，AI 通过协议调用，不碰 Shell，更适合沙箱与跨语言生态。

当前 MCP 尚未发布，**AI 接入首选 Skills**。MCP 落地后按场景选。详见 [Skills 接入](./integration-skills) 与 [MCP 接入](./integration-mcp)。

---

### Q14：AI 调用 CLI 时「卡顿」是怎么回事？

大概率是触发了 MITRE 速率限制（在线命令），CLI 在自动等待令牌，不是出错。AI 连续调用 `cwe show` / `cwe relations` 时容易触发。

解决：让 AI 优先用离线命令（带 `--xml`），离线不受速率限制且关系完整。或在提示词里告知 AI「批量场景用离线」。详见 [Skills 接入](./integration-skills) 注意事项。

---

## 🔧 技术细节

### Q15：为什么包名是 `cweskills` 而模块路径是 `cwe-skills`？

Go 的模块路径（`github.com/scagogogo/cwe-skills`）允许连字符，但**包名不允许连字符**（不是合法标识符）。所以包名去掉了连字符，定为 `cweskills`。导入路径用模块全名，调用用 `cweskills.` 前缀：

```go
import "github.com/scagogogo/cwe-skills" // 导入路径
cweskills.ParseCWEID("CWE-79")           // 包名前缀
```

详见 [Go SDK 接入](./integration-sdk) 与 [安装](./installation)。

---

### Q16：SDK 真的零依赖吗？CLI 不是用了 cobra？

是的，**核心 SDK 零第三方依赖**（`go.mod` 无 require）。CLI（`cmd/cwe`）用了 `cobra` 做命令行解析，但 cobra 只在 CLI 二进制里，SDK 用户 import `cweskills` 包时不会引入 cobra。详见 [性能与零依赖](./performance)。

---

## 📖 更多

没找到答案？试试：

- [GitHub Issues](https://github.com/scagogogo/cwe-skills/issues) 提问
- [工作原理](./how-it-works) 理解内部机制
- [四种接入方式总览](./integrations) 选对入口
- [SDK API](../sdk/overview) 查 API 细节
- [CLI 命令](../cli/overview) 查命令用法
