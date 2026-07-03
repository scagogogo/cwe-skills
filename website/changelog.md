# 📝 更新日志

<Badge type="info" text="Changelog" /> <Badge type="warning" text="当前 v0.0.1" />

CWE Skills 的版本演进记录。本页面基于 Git 提交历史整理。

::: tip 版本号规则
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。当前 `v0.0.1` 表示处于早期开发阶段，API 可能在小版本间调整。生产环境建议锁定版本：`go get github.com/scagogogo/cwe-skills@v0.0.1`。
:::

## 🚀 v0.0.1（初始发布）

<Badge type="tip" text="初始版本" />

首个公开版本，确立「AI 原生 CWE 集成层」定位，提供四种接入方式。

### ✨ 新增

- 🆔 **CWE ID 工具**：`ParseCWEID` / `FormatCWEID` / `ValidateCWEID` / `ExtractCWEIDs` / `CompareCWEIDs`，支持 `CWE-79`、`cwe 79`、`79` 等多种写法
- 🧱 **完整数据模型**：`CWE`、`Category`、`View`、`CompoundElement` 及 `Mitigation`、`Consequence`、`Reference` 等 16 个结构体
- 📚 **10 类枚举**：`Abstraction`、`Structure`、`Status`、`LikelihoodOfExploit`、`RelationshipNature`、`ConsequenceScope`、`ConsequenceImpact`、`ViewType`、`PlatformType`，每个枚举配备 `String/IsValid/Parse/AllValues` 四件套
- 🏆 **权威知名列表**：CWE Top 25（2024）、OWASP Top 10（2021）、SANS Top 25 完整映射，`IsInTop25` / `GetOWASPCategory` 等查询函数
- 🌐 **MITRE REST API 客户端**：`APIClient` + `HTTPClient`，内置令牌桶速率限制、指数退避重试、结构化错误
- 📥 **离线 XML 解析**：`XMLParser` 解析 MITRE 官方 `cwec_vX.Y.xml`，构建内存 `Registry`
- 🗃️ **内存注册表与索引**：`Registry` 存储 + `BuildIndexes` 多层关系索引
- 🧭 **关系导航**：`Navigator` 提供 13 种关系查询 + `ShortestPath` + `RelationshipDepth`
- 🌳 **层次树**：`BuildTree` / `BuildForest` / `BuildViewTree` + DFS/BFS 遍历
- 🔍 **搜索与过滤**：`FindByKeyword` 等 11 个查找函数 + `Filter` 多条件 + 排序/分组/去重
- 📊 **统计**：`ComputeStatistics` 各维度计数
- 📦 **序列化**：JSON / XML / CSV 三格式 + `safeCWE` 安全模型
- 💻 **40+ CLI 子命令**：基于 cobra，全命令支持 `-o text|json`
- 🦾 **AI Skills 接入**：12 篇渐进式技能文档
- ⚡ **零依赖核心**：SDK 仅用 Go 标准库
- 🖥️ **30+ 平台预编译**：Linux/macOS/Windows/BSD/AIX/Illumos/Solaris

### 🔧 重构

- 重命名 Go 包名 `cwe` → `cweskills`，与模块路径 `github.com/scagogogo/cwe-skills` 对齐
- 重写 README 为「AI 原生」定位，突出四种接入方式（Skills / SDK / CLI / MCP），Skills 放首位
- 大幅扩展 CLI 覆盖度，新增 `registry` / `nav` / `tree` / `filter` 命令组

### 🐛 修复

- 移除未使用的 `treeMaxDepth` 变量
- 修复 golangci-lint 报告的代码质量问题
- 简化 lint 配置（移除 `gofmt` / `gosimple` / `errcheck`，97.4% 测试覆盖率已覆盖）
- 修复图片和资源路径，使用 `import.meta.env.BASE_URL` 动态生成

### 📊 质量

- 测试覆盖率 **97.4%**
- 全平台 CI 通过

## 🔮 后续规划

- 🤖 **MCP 服务器**：提供 Model Context Protocol 接入，敬请期待
- 📈 更多知名列表（CWE Top 25 历年版本、行业特定列表）
- 🔗 与 NVD / CVE 的关联查询

::: details 查看完整 Git 历史
完整的提交历史请见 [GitHub Releases](https://github.com/scagogogo/cwe-skills/releases) 与 [提交历史](https://github.com/scagogogo/cwe-skills/commits/main)。
:::

## 🔗 相关

- [FAQ](./guide/faq)
- [安装指南](./guide/installation)
- [SDK 总览](./sdk/overview)
