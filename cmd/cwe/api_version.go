package main

import (
	"context"
	"fmt"

	cwepkg "github.com/scagogogo/cwe"
	"github.com/spf13/cobra"
)

var apiVersionBaseURL string

// apiVersionCmd 查询MITRE API版本
var apiVersionCmd = &cobra.Command{
	Use:   "api-version",
	Short: "查询MITRE CWE API版本",
	Long:  "查询MITRE CWE REST API的版本信息。",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cwepkg.NewAPIClient(cwepkg.WithAPIBaseURL(apiVersionBaseURL))
		defer client.Close()

		version, err := client.GetVersion(context.Background())
		if err != nil {
			return fmt.Errorf("查询API版本失败: %w", err)
		}

		if outputFormat == "json" {
			return printJSON(cmd, version)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "MITRE CWE API版本: %s\n", version.Version)
		if version.ReleaseDate != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "发布日期: %s\n", version.ReleaseDate)
		}
		if version.Name != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "版本名称: %s\n", version.Name)
		}
		return nil
	},
}

func init() {
	apiVersionCmd.Flags().StringVar(&apiVersionBaseURL, "base-url", "https://cwe-api.mitre.org/api", "MITRE API基础URL")
	rootCmd.AddCommand(apiVersionCmd)
}
