package main

import (
	"fmt"

	cwepkg "github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// enumCmd 枚举类型相关命令
var enumCmd = &cobra.Command{
	Use:   "enum",
	Short: "查询CWE枚举类型",
	Long: `查询CWE规范中定义的各种枚举类型，包括：
  - Abstraction (抽象层级)
  - Structure (结构类型)
  - Status (状态)
  - LikelihoodOfExploit (利用可能性)
  - RelationshipNature (关系类型)
  - ConsequenceScope (后果范围)
  - ConsequenceImpact (后果影响)
  - ViewType (视图类型)

示例：
  cwe enum abstraction       # 列出所有抽象层级
  cwe enum relationship      # 列出所有关系类型
  cwe enum status             # 列出所有状态`,
}

func init() {
	// 为每种枚举类型创建子命令
	enumTypes := []struct {
		name   string
		short  string
		values []string
	}{
		{"abstraction", "抽象层级 (Class/Base/Variant/Pillar)", stringifySlice(cwepkg.AllAbstractionValues())},
		{"structure", "结构类型 (Simple/Chain/Composite)", stringifySlice(cwepkg.AllStructureValues())},
		{"status", "状态 (Stable/Draft/Deprecated等)", stringifySlice(cwepkg.AllStatusValues())},
		{"likelihood", "利用可能性 (High/Medium/Low)", stringifySlice(cwepkg.AllLikelihoodOfExploitValues())},
		{"relationship", "关系类型 (ChildOf/ParentOf/CanPrecede等)", stringifySlice(cwepkg.AllRelationshipNatureValues())},
		{"scope", "后果范围 (Confidentiality/Integrity/Availability)", stringifySlice(cwepkg.AllConsequenceScopeValues())},
		{"impact", "后果影响 (High/Medium/Low)", stringifySlice(cwepkg.AllConsequenceImpactValues())},
		{"viewtype", "视图类型 (Graph/Slice)", stringifySlice(cwepkg.AllViewTypeValues())},
	}

	for _, et := range enumTypes {
		name := et.name
		short := et.short
		values := et.values

		subCmd := &cobra.Command{
			Use:   name,
			Short: short,
			RunE: func(cmd *cobra.Command, args []string) error {
				if outputFormat == "json" {
					return printJSON(cmd, values)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s (%d 项):\n", short, len(values))
				for _, v := range values {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", v)
				}
				return nil
			},
		}
		enumCmd.AddCommand(subCmd)
	}

	rootCmd.AddCommand(enumCmd)
}

// stringifySlice 将任意fmt.Stringer切片转换为字符串切片
func stringifySlice[T fmt.Stringer](items []T) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.String()
	}
	return result
}
