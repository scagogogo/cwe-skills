package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// printJSON 以JSON格式输出数据
func printJSON(cmd *cobra.Command, data interface{}) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printError 输出错误信息
func printError(cmd *cobra.Command, err error) {
	if outputFormat == "json" {
		_ = printJSON(cmd, map[string]string{"error": err.Error()})
	} else {
		fmt.Fprintf(cmd.ErrOrStderr(), "错误: %v\n", err)
	}
}
