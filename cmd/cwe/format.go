package main

import (
	"fmt"
	"strconv"
	"strings"

	cwepkg "github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// formatCmd 格式化CWE ID
var formatCmd = &cobra.Command{
	Use:   "format [CWE-ID...]",
	Short: "格式化CWE ID",
	Long: `将CWE ID格式化为标准格式 (CWE-XXX)。

示例：
  cwe format 79          # CWE-79
  cwe format cwe-89      # CWE-89
  cwe format CWE-352     # CWE-352`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个CWE ID")
		}

		type formatResult struct {
			Input  string `json:"input"`
			Output string `json:"output,omitempty"`
			Error  string `json:"error,omitempty"`
		}

		results := make([]formatResult, 0, len(args))
		for _, input := range args {
			result := formatResult{Input: input}
			id, err := cwepkg.ParseCWEID(input)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Output = cwepkg.FormatCWEIDFromInt(id)
			}
			results = append(results, result)
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		for _, r := range results {
			if r.Error != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%s -> 错误: %s\n", r.Input, r.Error)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", r.Output)
			}
		}
		return nil
	},
}

// extractCmd 从文本中提取CWE ID
var extractCmd = &cobra.Command{
	Use:   "extract [text...]",
	Short: "从文本中提取CWE ID",
	Long: `从输入文本中提取所有CWE ID。

支持以下格式：
  - CWE-79
  - CWE-89
  - CWE-352

示例：
  cwe extract "受CWE-79和CWE-89影响"
  cwe extract "漏洞: CWE-79, CWE-89, CWE-352"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一段文本")
		}

		text := strings.Join(args, " ")
		ids := cwepkg.ExtractCWEIDs(text)

		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"text":  text,
				"ids":   ids,
				"count": len(ids),
			})
		}

		if len(ids) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "未找到CWE ID\n")
			return nil
		}

		fmt.Fprintf(cmd.OutOrStdout(), "找到 %d 个CWE ID:\n", len(ids))
		for _, id := range ids {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", id)
		}
		return nil
	},
}

// compareCmd 比较两个CWE ID
var compareCmd = &cobra.Command{
	Use:   "compare <CWE-ID1> <CWE-ID2>",
	Short: "比较两个CWE ID",
	Long: `比较两个CWE ID的大小关系。

示例：
  cwe compare CWE-79 CWE-89    # CWE-79 < CWE-89
  cwe compare 89 79             # CWE-89 > CWE-79
  cwe compare CWE-79 CWE-79    # CWE-79 == CWE-79`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := cwepkg.CompareCWEIDs(args[0], args[1])
		if err != nil {
			return fmt.Errorf("比较失败: %w", err)
		}

		var comparison string
		switch {
		case result < 0:
			comparison = "less than"
		case result > 0:
			comparison = "greater than"
		default:
			comparison = "equal to"
		}

		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"id1":        args[0],
				"id2":        args[1],
				"comparison": comparison,
				"result":     result,
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s is %s %s\n", args[0], comparison, args[1])
		return nil
	},
}

// compareIntCmd 比较两个整数CWE ID
var compareIntCmd = &cobra.Command{
	Use:   "compare-int <int1> <int2>",
	Short: "比较两个整数CWE ID",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id1, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("第一个参数不是有效整数: %w", err)
		}
		id2, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("第二个参数不是有效整数: %w", err)
		}

		result := id1 - id2
		var comparison string
		switch {
		case result < 0:
			comparison = "less than"
		case result > 0:
			comparison = "greater than"
		default:
			comparison = "equal to"
		}

		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"id1":        cwepkg.FormatCWEIDFromInt(id1),
				"id2":        cwepkg.FormatCWEIDFromInt(id2),
				"comparison": comparison,
				"result":     result,
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s is %s %s\n",
			cwepkg.FormatCWEIDFromInt(id1), comparison, cwepkg.FormatCWEIDFromInt(id2))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(compareCmd)
	rootCmd.AddCommand(compareIntCmd)
}
