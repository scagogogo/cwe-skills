package main

import (
	"fmt"
	"os"

	cwepkg "github.com/scagogogo/cwe-skills"
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

SDK版本: %s`, cwepkg.Version),
	SilenceUsage:  true,
	SilenceErrors: true,
}

// outputFormat 输出格式：text 或 json
var outputFormat string

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "输出格式 (text|json)")
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
