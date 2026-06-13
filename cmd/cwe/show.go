package main

import (
	"context"
	"fmt"
	"time"

	cwepkg "github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

var (
	showBaseURL string
	showTimeout int
)

// showCmd 通过MITRE API获取CWE详细信息
var showCmd = &cobra.Command{
	Use:   "show [CWE-ID...]",
	Short: "获取CWE详细信息",
	Long: `通过MITRE CWE REST API获取CWE条目的详细信息。

需要网络连接到MITRE API (https://cwe-api.mitre.org/api)。

示例：
  cwe show CWE-79
  cwe show 79 89 352
  cwe show --base-url https://cwe-api.mitre.org/api CWE-79`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个CWE ID")
		}

		opts := []cwepkg.APIClientOption{
			cwepkg.WithAPIBaseURL(showBaseURL),
		}
		if showTimeout > 0 {
			opts = append(opts, cwepkg.WithAPITimeout(time.Duration(showTimeout)*time.Second))
		}

		client := cwepkg.NewAPIClient(opts...)
		defer client.Close()

		type showResult struct {
			CWEID  string      `json:"cwe_id"`
			Detail interface{} `json:"detail,omitempty"`
			Error  string      `json:"error,omitempty"`
		}

		results := make([]showResult, 0, len(args))
		for _, input := range args {
			id, err := cwepkg.ParseCWEID(input)
			if err != nil {
				results = append(results, showResult{CWEID: input, Error: err.Error()})
				continue
			}

			cwe, err := client.GetWeakness(context.Background(), id)
			if err != nil {
				results = append(results, showResult{
					CWEID: cwepkg.FormatCWEIDFromInt(id),
					Error: err.Error(),
				})
				continue
			}

			results = append(results, showResult{
				CWEID:  cwepkg.FormatCWEIDFromInt(id),
				Detail: cwe,
			})
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		for _, r := range results {
			if r.Error != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: 错误 - %s\n", r.CWEID, r.Error)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "=== %s ===\n", r.CWEID)
				if cwe, ok := r.Detail.(*cwepkg.CWE); ok {
					fmt.Fprintf(cmd.OutOrStdout(), "  名称:     %s\n", cwe.Name)
					fmt.Fprintf(cmd.OutOrStdout(), "  抽象层级: %s\n", cwe.Abstraction)
					fmt.Fprintf(cmd.OutOrStdout(), "  状态:     %s\n", cwe.Status)
					if cwe.Description != "" {
						fmt.Fprintf(cmd.OutOrStdout(), "  描述:     %s\n", cwe.Description)
					}
					if cwe.Structure != "" {
						fmt.Fprintf(cmd.OutOrStdout(), "  结构:     %s\n", cwe.Structure)
					}
					if len(cwe.Relationships) > 0 {
						fmt.Fprintf(cmd.OutOrStdout(), "  关系:     %d 项\n", len(cwe.Relationships))
					}
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}
		}
		return nil
	},
}

// showCategoryCmd 获取CWE类别详细信息
var showCategoryCmd = &cobra.Command{
	Use:   "category [ID...]",
	Short: "获取CWE类别详细信息",
	Long:  "通过MITRE CWE REST API获取CWE类别的详细信息。",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个类别ID")
		}

		client := cwepkg.NewAPIClient(cwepkg.WithAPIBaseURL(showBaseURL))
		defer client.Close()

		for _, input := range args {
			id, err := cwepkg.ParseCWEID(input)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: 无效ID - %v\n", input, err)
				continue
			}

			cat, err := client.GetCategory(context.Background(), id)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: 错误 - %v\n", cwepkg.FormatCWEIDFromInt(id), err)
				continue
			}

			if outputFormat == "json" {
				_ = printJSON(cmd, cat)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "=== %s (类别) ===\n", cwepkg.FormatCWEIDFromInt(id))
				fmt.Fprintf(cmd.OutOrStdout(), "  名称: %s\n", cat.Name)
				if cat.Description != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  描述: %s\n", cat.Description)
				}
			}
		}
		return nil
	},
}

// showViewCmd 获取CWE视图详细信息
var showViewCmd = &cobra.Command{
	Use:   "view [ID...]",
	Short: "获取CWE视图详细信息",
	Long:  "通过MITRE CWE REST API获取CWE视图的详细信息。",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个视图ID")
		}

		client := cwepkg.NewAPIClient(cwepkg.WithAPIBaseURL(showBaseURL))
		defer client.Close()

		for _, input := range args {
			id, err := cwepkg.ParseCWEID(input)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: 无效ID - %v\n", input, err)
				continue
			}

			view, err := client.GetView(context.Background(), id)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: 错误 - %v\n", cwepkg.FormatCWEIDFromInt(id), err)
				continue
			}

			if outputFormat == "json" {
				_ = printJSON(cmd, view)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "=== %s (视图) ===\n", cwepkg.FormatCWEIDFromInt(id))
				fmt.Fprintf(cmd.OutOrStdout(), "  名称: %s\n", view.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "  类型: %s\n", view.Type)
				if view.Description != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  描述: %s\n", view.Description)
				}
			}
		}
		return nil
	},
}

func init() {
	showCmd.PersistentFlags().StringVar(&showBaseURL, "base-url", "https://cwe-api.mitre.org/api", "MITRE API基础URL")
	showCmd.PersistentFlags().IntVar(&showTimeout, "timeout", 30, "API请求超时时间(秒)")

	showCmd.AddCommand(showCategoryCmd)
	showCmd.AddCommand(showViewCmd)
	rootCmd.AddCommand(showCmd)
}
