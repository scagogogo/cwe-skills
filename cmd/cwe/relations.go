package main

import (
	"context"
	"fmt"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

var (
	relBaseURL string
	relViewID  int
)

// relCmd 关系查询命令
var relCmd = &cobra.Command{
	Use:   "relations",
	Short: "查询CWE关系",
	Long: `通过MITRE CWE REST API查询CWE条目之间的关系。

支持查询：
  - 父级弱点 (parents)
  - 子级弱点 (children)
  - 祖先弱点 (ancestors)
  - 后代弱点 (descendants)

示例：
  cwe relations parents CWE-79
  cwe relations children CWE-74
  cwe relations ancestors CWE-79`,
}

// parentsCmd 查询父级弱点
var parentsCmd = &cobra.Command{
	Use:   "parents [CWE-ID]",
	Short: "查询CWE的父级弱点",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return fmt.Errorf("无效CWE ID: %w", err)
		}

		client := cweskills.NewAPIClient(cweskills.WithAPIBaseURL(relBaseURL))
		defer client.Close()

		var rels []cweskills.Relationship
		if relViewID > 0 {
			rels, err = client.GetParents(context.Background(), id, relViewID)
		} else {
			rels, err = client.GetParents(context.Background(), id)
		}
		if err != nil {
			return fmt.Errorf("查询父级失败: %w", err)
		}

		return printRelationships(cmd, "父级弱点", id, rels)
	},
}

// childrenCmd 查询子级弱点
var childrenCmd = &cobra.Command{
	Use:   "children [CWE-ID]",
	Short: "查询CWE的子级弱点",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return fmt.Errorf("无效CWE ID: %w", err)
		}

		client := cweskills.NewAPIClient(cweskills.WithAPIBaseURL(relBaseURL))
		defer client.Close()

		var rels []cweskills.Relationship
		if relViewID > 0 {
			rels, err = client.GetChildren(context.Background(), id, relViewID)
		} else {
			rels, err = client.GetChildren(context.Background(), id)
		}
		if err != nil {
			return fmt.Errorf("查询子级失败: %w", err)
		}

		return printRelationships(cmd, "子级弱点", id, rels)
	},
}

// ancestorsCmd 查询祖先弱点
var ancestorsCmd = &cobra.Command{
	Use:   "ancestors [CWE-ID]",
	Short: "查询CWE的祖先弱点",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return fmt.Errorf("无效CWE ID: %w", err)
		}

		client := cweskills.NewAPIClient(cweskills.WithAPIBaseURL(relBaseURL))
		defer client.Close()

		rels, err := client.GetAncestors(context.Background(), id)
		if err != nil {
			return fmt.Errorf("查询祖先失败: %w", err)
		}

		return printRelationships(cmd, "祖先弱点", id, rels)
	},
}

// descendantsCmd 查询后代弱点
var descendantsCmd = &cobra.Command{
	Use:   "descendants [CWE-ID]",
	Short: "查询CWE的后代弱点",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return fmt.Errorf("无效CWE ID: %w", err)
		}

		client := cweskills.NewAPIClient(cweskills.WithAPIBaseURL(relBaseURL))
		defer client.Close()

		rels, err := client.GetDescendants(context.Background(), id)
		if err != nil {
			return fmt.Errorf("查询后代失败: %w", err)
		}

		return printRelationships(cmd, "后代弱点", id, rels)
	},
}

func printRelationships(cmd *cobra.Command, relType string, id int, rels []cweskills.Relationship) error {
	if outputFormat == "json" {
		return printJSON(cmd, map[string]interface{}{
			"cwe_id":        cweskills.FormatCWEIDFromInt(id),
			"relation_type": relType,
			"relationships": rels,
			"count":         len(rels),
		})
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s 的 %s (%d 项):\n", cweskills.FormatCWEIDFromInt(id), relType, len(rels))
	for _, rel := range rels {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s -> %s", rel.Nature, cweskills.FormatCWEIDFromInt(rel.CWEID))
		if rel.ViewID > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), " (View: %d)", rel.ViewID)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}
	return nil
}

func init() {
	relCmd.PersistentFlags().StringVar(&relBaseURL, "base-url", "https://cwe-api.mitre.org/api", "MITRE API基础URL")
	relCmd.PersistentFlags().IntVar(&relViewID, "view-id", 0, "视图ID (可选，仅对parents/children有效)")

	relCmd.AddCommand(parentsCmd)
	relCmd.AddCommand(childrenCmd)
	relCmd.AddCommand(ancestorsCmd)
	relCmd.AddCommand(descendantsCmd)
	rootCmd.AddCommand(relCmd)
}
