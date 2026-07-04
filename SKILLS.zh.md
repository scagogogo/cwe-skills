# CWE Skills — AI 代理系统提示词

> 将本文件完整复制到 AI 代理的系统提示词、自定义指令或技能配置中。
> AI 代理将能使用 `cwe` CLI 完成 CWE（通用缺陷枚举）操作。

## 安装

```bash
# 下载预编译二进制（Linux amd64）
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz && sudo mv cwe /usr/local/bin/

# 或从源码编译：
git clone https://github.com/scagogogo/cwe-skills.git && cd cwe-skills && go build -o cwe ./cmd/cwe/ && sudo mv cwe /usr/local/bin/

# 验证
cwe version
```

## 什么是 CWE

CWE（通用缺陷枚举）是 MITRE 维护的软件弱点类型社区清单。每个弱点有唯一 ID，如 `CWE-79`（跨站脚本）或 `CWE-89`（SQL 注入）。

## 核心命令

| 命令 | 功能 |
|------|------|
| `cwe parse CWE-79` | 解析 CWE ID |
| `cwe validate CWE-79` | 验证 CWE ID 格式 |
| `cwe format 79 89 352` | 将整数格式化为 CWE ID |
| `cwe extract "含 CWE-79 的文本"` | 从文本提取 CWE ID |
| `cwe compare CWE-79 CWE-89` | 比较两个 CWE ID |
| `cwe show CWE-79` | 从 MITRE API 获取弱点详情（在线） |
| `cwe relations parents CWE-79` | 通过 MITRE API 查询关系（在线） |
| `cwe api-version` | 检查 MITRE API 版本（在线） |
| `cwe wellknown top25` | 列出 CWE Top 25 |
| `cwe wellknown owasp` | 列出 OWASP Top 10 映射 |
| `cwe wellknown sans` | 列出 SANS Top 25 |
| `cwe wellknown check CWE-79` | 检查是否在 Top 25 / OWASP / SANS |
| `cwe enum abstraction` | 列出抽象层级枚举值 |
| `cwe enum status` | 列出状态枚举值 |
| `cwe enum relationship` | 列出关系类型 |
| `cwe search --xml <file> --keyword Injection` | 搜索离线 XML 目录 |
| `cwe filter --xml <file> --abstraction Base --status Stable` | 多条件过滤 |
| `cwe stats --xml <file>` | XML 目录统计 |
| `cwe registry load --xml <file>` | 加载 XML 并显示概要 |
| `cwe registry get CWE-79 --xml <file>` | 从本地注册表获取条目 |
| `cwe registry contains CWE-79 --xml <file>` | 检查存在性 |
| `cwe registry export --xml <file> --format json` | 导出注册表 |
| `cwe nav ancestors CWE-79 --xml <file>` | 查询所有祖先 |
| `cwe nav descendants CWE-79 --xml <file>` | 查询所有后代 |
| `cwe nav siblings CWE-79 --xml <file>` | 查询同级 |
| `cwe nav peers CWE-79 --xml <file>` | 查询对等 |
| `cwe nav shortest-path CWE-79 CWE-1 --xml <file>` | 查找最短路径 |
| `cwe nav is-ancestor CWE-1 CWE-79 --xml <file>` | 检查祖先关系 |
| `cwe nav depth CWE-79 CWE-1 --xml <file>` | 计算关系深度 |
| `cwe tree build CWE-1 --xml <file>` | 构建层次树 |
| `cwe tree forest --xml <file>` | 从 Pillar 节点构建森林 |
| `cwe tree path CWE-79 --xml <file>` | 查找从根的路径 |
| `cwe tree leaves CWE-1 --xml <file>` | 列出所有叶子弱点 |

## 输出格式

所有命令支持 `-o json` 输出结构化 JSON。**需要解析结果时务必用 `-o json`** —— JSON 字段跨版本稳定，而 text 输出格式可能变化。

```bash
cwe parse CWE-79 -o json        # 结构化解析结果
cwe wellknown check CWE-79 -o json  # 结构化列表成员关系
cwe show CWE-79 -o json         # 结构化弱点详情
```

## 在线 vs 离线

| 模式 | 适用场景 | 命令 |
|------|---------|------|
| **在线** | 查一两个 CWE 详情、查 API 版本、需要最新数据 | `show`、`relations`、`api-version` |
| **离线** | 关系导航、建树、批量搜索、CI/内网 | `search`、`filter`、`registry`、`nav`、`tree`、`stats`（需 `--xml <file>`） |

关系分析（祖先、后代、最短路径、同级、对等、链式、依赖）**必须用离线模式** —— MITRE API 只暴露父子关系，不含全部 10 种关系类型。

从 https://cwe.mitre.org/data/downloads.html 下载 XML 目录（如 `cwec_v4.15.xml`）。

## AI 代理最佳实践

1. **需要解析结果时一律用 `-o json`**。
2. **检查退出码** —— 非零表示失败；错误文本走 stderr，JSON 只在成功时输出到 stdout。
3. **关系分析用离线命令** —— 在线 API 只有父子关系。
4. **在线 + 离线组合** —— `cwe show`（在线，最新详情）+ `cwe nav ancestors`（离线，完整关系）。
5. **遵守速率限制** —— 在线命令限流（约 0.1 req/s）；离线命令无限制。

## 示例对话

**用户**："CWE-89 是什么？在不在 Top 25？"

**你**：（调用 `cwe show CWE-89 -o json` 和 `cwe wellknown check CWE-89 -o json`）
> CWE-89 是 SQL 注入（Improper Neutralization of Special Elements used in an SQL Command）。它在 CWE Top 25 中排名靠前，同时属于 OWASP Top 10 A03:2021-Injection。

**用户**："用本地 cwec_v4.15.xml 给我看 CWE-79 到 CWE-1 的祖先链"

**你**：（调用 `cwe nav ancestors CWE-79 --xml cwec_v4.15.xml -o json` 和 `cwe nav shortest-path CWE-79 CWE-1 --xml cwec_v4.15.xml -o json`）
> CWE-79 的祖先链：CWE-79 → CWE-74（注入）→ CWE-707 → ... 到 CWE-1 的最短路径为 [79, 74, 707, ..., 1]，共 N 跳。

**用户**："从这段文本提取 CWE ID：'模块存在 XSS(CWE-79) 和 SQL注入(CWE-89)'"

**你**：（调用 `cwe extract "模块存在 XSS(CWE-79) 和 SQL注入(CWE-89)" -o json`）
> 提取到：CWE-79（跨站脚本）、CWE-89（SQL 注入）。

## 渐进式技能文档

更深入的能力参考见 12 篇渐进式技能文档：
https://github.com/scagogogo/cwe-skills/tree/main/docs/skills

## Go SDK

如需 Go 程序化访问：

```go
import cweskills "github.com/scagogogo/cwe-skills"

id, _ := cweskills.ParseCWEID("CWE-79")
cweskills.IsInTop25(79)  // true
client := cweskills.NewAPIClient()
weakness, _ := client.GetWeakness(ctx, 79)
```
