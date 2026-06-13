package main

import (
	"fmt"

	cwepkg "github.com/scagogogo/cwe"
	"github.com/spf13/cobra"
)

// validateCmd 验证CWE ID格式
var validateCmd = &cobra.Command{
	Use:   "validate [CWE-ID...]",
	Short: "验证CWE ID格式",
	Long: `验证CWE ID格式是否有效。

示例：
  cwe validate CWE-79      # 有效
  cwe validate abc          # 无效
  cwe validate CWE-79 CWE-89 abc`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个CWE ID")
		}

		type validateResult struct {
			Input string `json:"input"`
			Valid bool   `json:"valid"`
			Error string `json:"error,omitempty"`
		}

		results := make([]validateResult, 0, len(args))
		allValid := true
		for _, input := range args {
			result := validateResult{Input: input}
			err := cwepkg.ValidateCWEID(input)
			if err != nil {
				result.Valid = false
				result.Error = err.Error()
				allValid = false
			} else {
				result.Valid = true
			}
			results = append(results, result)
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		for _, r := range results {
			if r.Valid {
				fmt.Fprintf(cmd.OutOrStdout(), "%s ✓ 有效\n", r.Input)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s ✗ 无效: %s\n", r.Input, r.Error)
			}
		}

		if !allValid {
			return fmt.Errorf("部分CWE ID无效")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
