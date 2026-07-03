---
title: 安装
outline: [2, 3]
---

# ⬇️ 安装

CWE Skills 的安装分两类：**CLI 工具**（Shell 脚本与开发工作流用）和 **Go SDK**（嵌入 Go 应用用）。两者数据互通、能力对等，按需选择。

::: tip 先确定要装哪个
- 只想在命令行跑 `cwe` 命令、写 Shell 脚本、或给 AI 代理用 → 装 **CLI**
- 想在 Go 程序里 `import` 调用 → 装 **Go SDK**
- 两者都要 → 先装 SDK（`go get`），再装 CLI 二进制
:::

---

## 🅰️ 安装 CLI

### 1. 从 Release 下载（推荐）

预编译二进制覆盖 30+ 平台。从 [Releases 最新版](https://github.com/scagogogo/cwe-skills/releases/latest) 下载对应平台的压缩包解压即可。

**Linux (amd64)**：

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_linux_x86_64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/
```

**macOS (Apple Silicon)**：

```bash
curl -sL https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_darwin_aarch64.tar.gz | tar xz
sudo mv cwe /usr/local/bin/
```

**macOS (Intel)**：把 `darwin_aarch64` 换成 `darwin_x86_64`。

**Windows (PowerShell)**：

```powershell
Invoke-WebRequest -Uri https://github.com/scagogogo/cwe-skills/releases/latest/download/cwe-skills_latest_windows_x86_64.zip -OutFile cwe.zip
Expand-Archive cwe.zip
Move-Item cwe.exe C:\Windows\cwe.exe
```

### 2. 包管理器

```bash
# Homebrew (macOS/Linux)
brew install scagogogo/tap/cwe-skills

# Scoop (Windows)
scoop install cwe-skills
```

### 3. Go 安装（需本地有 Go 工具链）

```bash
go install github.com/scagogogo/cwe-skills/cmd/cwe@latest
```

::: info go install 装到哪
二进制会装到 `$GOBIN`（通常在 `$GOPATH/bin` 或 `$HOME/go/bin`），请确保该目录在 `PATH` 中。
:::

### 4. 系统包（deb / rpm / apk）

Release 附带 `.deb`、`.rpm`、`.apk` 安装包，适合容器与系统级部署：

```bash
# Debian/Ubuntu
sudo dpkg -i cwe-skills_*_linux_amd64.deb

# RHEL/CentOS/Fedora
sudo rpm -i cwe-skills-*_linux_amd64.rpm

# Alpine
sudo apk add --allow-untrusted cwe-skills-*-linux-amd64.apk
```

### 5. 从源码编译

```bash
git clone https://github.com/scagogogo/cwe-skills.git
cd cwe-skills
go build -o cwe ./cmd/cwe/
sudo mv cwe /usr/local/bin/
```

::: details 注入版本号（可选）
CLI 的版本号通过 `-ldflags` 注入 `cmd/cwe` 包的 `cliVersion` 变量，默认为 `dev`：

```bash
go build -ldflags "-X main.cliVersion=v0.0.1" -o cwe ./cmd/cwe/
```
:::

### 6. 验证

```bash
cwe version
```

```text
cwe CLI 版本: dev
SDK版本: v0.0.1
```

看到 `SDK版本: v0.0.1` 即安装成功。`cwe --help` 可查看全部子命令。

---

## 🅱️ 安装 Go SDK

```bash
go get github.com/scagogogo/cwe-skills
```

在代码中导入（注意包名是 `cweskills`，与模块路径最后一段 `cwe-skills` 不同）：

```go
import "github.com/scagogogo/cwe-skills"

fmt.Println(cweskills.Version) // v0.0.1
```

::: warning 包名注意
模块路径是 `github.com/scagogogo/cwe-skills`，但 Go 包名是 `cweskills`。导入后用 `cweskills.` 前缀调用，不要写成 `cweskills.` 也不要写成 `cwe.`。
:::

---

## 🖥️ 预编译平台一览

CLI 预编译二进制覆盖 30+ 平台，下表列出主要组合（完整列表见 [Releases](https://github.com/scagogogo/cwe-skills/releases)）：

| 操作系统 | 架构 | 备注 |
|---------|------|------|
| Linux | amd64, 386, arm64, arm, mips, mipsle, mips64, mips64le, ppc64, ppc64le, s390x, riscv64 | 覆盖最广 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) | |
| Windows | amd64, 386, arm64 | |
| FreeBSD | amd64, 386, arm64, arm | |
| NetBSD | amd64, 386, arm64, arm | |
| OpenBSD | amd64, 386, arm64, arm | |
| AIX | ppc64 | IBM Power |
| Illumos | amd64 | illumos/OmniOS |
| Solaris | amd64 | Oracle Solaris |

::: tip 找不到你的平台？
`goreleaser` 配置覆盖很广；若仍缺，可用「从源码编译」方式，只要你的平台有 Go 工具链即可。
:::

---

## 📥 下载离线 XML（可选，离线模式需要）

离线模式（搜索、过滤、导航、建树）需要 MITRE 官方 XML 弱点目录：

1. 访问 <https://cwe.mitre.org/data/downloads.html>
2. 下载最新版（如 `cwec_v4.15.xml`，文件较大）
3. 放到任意路径，使用时通过 `--xml <path>` 指定

```bash
# 验证可解析
cwe stats --xml cwec_v4.15.xml
```

详见 [在线 vs 离线模式](./online-offline)。

---

## 🐳 Docker 中使用

CLI 是单文件静态二进制，直接复制进镜像即可：

```dockerfile
FROM alpine:latest
COPY cwe /usr/local/bin/cwe
RUN apk add --no-cache ca-certificates
ENTRYPOINT ["cwe"]
```

::: warning 别忘了 ca-certificates
在线模式调用 HTTPS 的 MITRE API，镜像里需要有 CA 证书（`ca-certificates`），否则 TLS 握手会失败。
:::

---

## 🔄 升级

- CLI：重新执行 Release 下载命令，或 `brew upgrade` / `scoop update cwe-skills` / `go install ...@latest`。
- SDK：`go get github.com/scagogogo/cwe-skills@latest`。

---

## 🧰 环境要求

- **CLI 二进制**：无任何运行时依赖，下载即用。
- **从源码编译 / SDK**：Go 1.25+（见 `go.mod`）。
- **在线模式**：能访问 `https://cwe-api.mitre.org`。
- **离线模式**：仅需本地 XML 文件，无网络。

---

## ❓ 常见安装问题

::: details `command not found: cwe`
二进制所在目录不在 `PATH` 中。`/usr/local/bin` 通常在 PATH 内；若用 `go install`，确保 `$GOPATH/bin` 在 PATH。
:::

::: details macOS 提示「无法验证开发者」
`sudo xattr -d com.apple.quarantine /usr/local/bin/cwe` 解除隔离属性。
:::

::: details 在线调用报 TLS 错误
系统缺少 CA 证书。Linux 装 `ca-certificates`，或在 SDK 里用 `WithHTTPClient` 传入自定义 `*http.Client`。
:::

更多问题见 [FAQ](./faq)。

---

## 📖 下一步

- 跑通第一个命令 → [快速开始](./quick-start)
- 理解两条数据路径 → [在线 vs 离线模式](./online-offline)
- 全部 CLI 命令 → [CLI 命令](../cli/overview)
