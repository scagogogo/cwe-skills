package main

import (
	"fmt"
	"runtime"

	cwepkg "github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// versionCmd 显示版本信息
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  "显示CWE CLI工具的版本信息，包括SDK版本、Go版本和构建信息。",
	RunE: func(cmd *cobra.Command, args []string) error {
		info := struct {
			CLI       string `json:"cli"`
			SDK       string `json:"sdk"`
			Go        string `json:"go"`
			GitCommit string `json:"git_commit,omitempty"`
			BuildDate string `json:"build_date,omitempty"`
		}{
			CLI:       cliVersion,
			SDK:       cwepkg.Version,
			Go:        runtime.Version(),
			GitCommit: cliGitCommit,
			BuildDate: cliBuildDate,
		}

		if outputFormat == "json" {
			return printJSON(cmd, info)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "CWE CLI:     %s\n", info.CLI)
		fmt.Fprintf(cmd.OutOrStdout(), "CWE SDK:     %s\n", info.SDK)
		fmt.Fprintf(cmd.OutOrStdout(), "Go Version:  %s\n", info.Go)
		if info.GitCommit != "unknown" {
			fmt.Fprintf(cmd.OutOrStdout(), "Git Commit:  %s\n", info.GitCommit)
		}
		if info.BuildDate != "unknown" {
			fmt.Fprintf(cmd.OutOrStdout(), "Build Date:  %s\n", info.BuildDate)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
