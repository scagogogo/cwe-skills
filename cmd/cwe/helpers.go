package main

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

// printJSON 以JSON格式输出数据
func printJSON(cmd *cobra.Command, data interface{}) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
