# 🛠️ 安装

本文介绍如何获取并安装 `cwe` 命令行工具。

## 前置条件

- 操作系统：Linux / macOS / Windows（及 WSL）
- 网络访问（仅在使用 <Badge type="info" text="在线"/> 命令时需要，访问 MITRE API）
- 本地 XML 目录文件（仅在使用 <Badge type="tip" text="离线"/> 命令时需要，见下文）

## 方式一：从源码构建（推荐）

需要本地安装 [Go](https://go.dev/) 1.21+。

```bash
# 克隆仓库
git clone https://github.com/scagogogo/cwe-skills.git
cd cwe-skills

# 直接构建 CLI 入口
go build -o cwe ./cmd/cwe

# 安装到 $GOPATH/bin（需确保该目录在 PATH 中）
go install ./cmd/cwe
```

构建时可注入版本信息：

```bash
VERSION="v1.0.0"
COMMIT=$(git rev-parse HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -ldflags "-X main.cliVersion=$VERSION \
  -X main.cliGitCommit=$COMMIT \
  -X main.cliBuildDate=$DATE" -o cwe ./cmd/cwe
```

注入后 [`cwe version`](./version) 会显示完整的构建元信息。

## 方式二：go install

若仓库已发布且模块路径可解析，可一行安装：

```bash
go install github.com/scagogogo/cwe-skills/cmd/cwe@latest
```

::: warning 模块路径
包名已于近期从 `cwe` 重命名为 `cweskills`，与模块路径 `github.com/scagogogo/cwe-skills` 对齐。若遇导入错误，请确保使用最新代码。
:::

## 验证安装

```bash
cwe version
```

预期输出（text 模式）：

```text
CWE CLI:     dev
CWE SDK:     v0.x.x
Go Version:  go1.21.x
```

## 准备离线数据

使用 `search`、`registry`、`nav`、`tree` 等命令前，需下载 MITRE 官方 XML 目录：

```bash
# 从 https://cwe.mitre.org/data/xml.html 下载，例如：
wget https://cwe.mitre.org/data/downloads/cwec_latest.xml.zip
unzip cwec_latest.xml.zip
# 得到 cwec_latest.xml
```

随后在命令中通过 `--xml cwec_latest.xml`（或简写 `-x`）指定路径。

## 下一步

- [CLI 总览](./overview)
- [全局参数](./global-flags)
- 离线命令示例：[search](./search)、[registry load](./registry-load)

## 相关文档

- [SDK 安装与引入](../sdk/introduction)
- [XML 解析器](../sdk/xml-parser)
