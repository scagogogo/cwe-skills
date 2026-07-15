package main

import (
	"fmt"
	"io"
	"os"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// cliVersion CLI工具的版本号，构建时通过 -ldflags 注入
var cliVersion = "dev"

// cliGitCommit Git提交哈希，构建时通过 -ldflags 注入
var cliGitCommit = "unknown"

// cliBuildDate 构建日期，构建时通过 -ldflags 注入
var cliBuildDate = "unknown"

// rootCmd 是CWE CLI的根命令
var rootCmd = &cobra.Command{
	Use:   "cwe",
	Short: "CWE (Common Weakness Enumeration) 命令行工具",
	Long: fmt.Sprintf(`CWE (Common Weakness Enumeration) 命令行工具

基于 github.com/scagogogo/cwe-skills SDK 构建的CWE通用缺陷枚举命令行工具。
支持CWE ID的解析、验证、搜索，以及MITRE CWE REST API的查询。

SDK版本: %s`, cweskills.Version),
	SilenceUsage:  true,
	SilenceErrors: true,
}

// outputFormat 输出格式：text 或 json
var outputFormat string

// osStderr 默认指向 os.Stderr，测试可替换以捕获输出。
var osStderr io.Writer = os.Stderr

// osExit 默认指向 os.Exit，测试可替换为不退出进程的 fake
// 以覆盖 main/Execute 的 os.Exit 调用点。
var osExit = os.Exit

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "输出格式 (text|json)")
}

// Execute 执行根命令
func Execute() {
	osExit(executeRoot())
}

// executeRoot 执行根命令并返回退出码。
// 提取自 Execute 以便测试覆盖 cobra 执行成功与失败两条路径，
// 而不触发 os.Exit。
func executeRoot() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(osStderr, err)
		return 1
	}
	return 0
}
