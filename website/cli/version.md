# 💻 cwe version

显示 CWE CLI 工具的版本信息，包括 SDK 版本、Go 版本和构建信息。

## 语法

```bash
cwe version [flags]
```

## 参数

本命令无位置参数。

## Flags

仅支持 [全局参数](./global-flags) `-o/--output`。

## 示例

### text 输出

```bash
cwe version
```

```text
CWE CLI:     dev
CWE SDK:     v0.x.x
Go Version:  go1.21.x
```

当通过 `-ldflags` 注入构建信息时，会额外显示 Git 提交与构建日期：

```text
CWE CLI:     v1.0.0
CWE SDK:     v0.x.x
Go Version:  go1.21.x
Git Commit:  abc1234
Build Date:  2026-07-03T...
```

### JSON 输出

```bash
cwe version -o json
```

```json
{
  "cli": "dev",
  "sdk": "v0.x.x",
  "go": "go1.21.x",
  "git_commit": "unknown",
  "build_date": "unknown"
}
```

## 字段说明

| 字段 | JSON 键 | 说明 |
| --- | --- | --- |
| CLI 版本 | `cli` | CLI 工具自身版本，构建时通过 `-ldflags` 注入 `main.cliVersion` |
| SDK 版本 | `sdk` | 依赖的 `github.com/scagogogo/cwe-skills` SDK 版本（`cweskills.Version`） |
| Go 版本 | `go` | 编译用的 Go 版本（`runtime.Version()`） |
| Git 提交 | `git_commit` | 构建时的 Git 提交哈希，未注入时为 `unknown` |
| 构建日期 | `build_date` | 构建日期，未注入时为 `unknown` |

## 使用场景

- 排查环境问题，确认 CLI 与 SDK 版本。
- 反馈 bug 时附带版本信息以便复现。
- 验证自定义构建是否正确注入了版本元数据。

::: tip 注入构建信息
构建时通过 `-ldflags` 注入 `main.cliVersion`、`main.cliGitCommit`、`main.cliBuildDate` 即可在 `version` 输出中看到完整信息，详见 [安装](./install)。
:::

## 下一步

- [CLI 总览](./overview)
- [安装](./install)
- [全局参数](./global-flags)

## 相关文档

- [SDK Version](../sdk/package)
- [SDK 总览](../sdk/overview)
