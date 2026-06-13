package main

import (
	"fmt"

	cwepkg "github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// parseCmd 解析CWE ID
var parseCmd = &cobra.Command{
	Use:   "parse [CWE-ID...]",
	Short: "解析CWE ID",
	Long: `解析CWE ID，提取数字部分并格式化为标准格式。

支持多种输入格式：
  - 纯数字: 79
  - 标准格式: CWE-79
  - 大小写不敏感: cwe-79

示例：
  cwe parse CWE-79
  cwe parse 79 cwe-89 CWE-352`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个CWE ID")
		}

		type parseResult struct {
			Input   string `json:"input"`
			ID      int    `json:"id"`
			Format  string `json:"format,omitempty"`
			Valid   bool   `json:"valid"`
			Error   string `json:"error,omitempty"`
		}

		results := make([]parseResult, 0, len(args))
		for _, input := range args {
			result := parseResult{Input: input}

			id, err := cwepkg.ParseCWEID(input)
			if err != nil {
				result.Valid = false
				result.Error = err.Error()
			} else {
				result.Valid = true
				result.ID = id
				result.Format = cwepkg.FormatCWEIDFromInt(id)
			}

			results = append(results, result)
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		for _, r := range results {
			if r.Valid {
				fmt.Fprintf(cmd.OutOrStdout(), "%s -> %s (ID: %d)\n", r.Input, r.Format, r.ID)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s -> 无效: %s\n", r.Input, r.Error)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
