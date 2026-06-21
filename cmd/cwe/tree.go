package main

import (
	"fmt"
	"strings"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

var treeXMLPath string

// treeCmd 树构建命令组
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "CWE层次树操作",
	Long: `构建和查询CWE层次树，支持从XML目录数据构建树并遍历。

所有tree子命令都需要通过 --xml 参数指定CWE XML目录文件。

示例：
  cwe tree build CWE-1 --xml cwec_latest.xml
  cwe tree forest --xml cwec_latest.xml
  cwe tree view 1000 --xml cwec_latest.xml
  cwe tree path CWE-79 --xml cwec_latest.xml
  cwe tree leaves CWE-1 --xml cwec_latest.xml`,
}

func loadTreeRegistry() (*cweskills.Registry, error) {
	if treeXMLPath == "" {
		return nil, fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
	}
	parser := cweskills.NewXMLParser()
	registry, err := parser.ParseFile(treeXMLPath)
	if err != nil {
		return nil, fmt.Errorf("解析XML文件失败: %w", err)
	}
	registry.BuildIndexes()
	return registry, nil
}

// treeBuildCmd 构建指定根节点的树
var treeBuildCmd = &cobra.Command{
	Use:   "build [ROOT-ID]",
	Short: "构建指定根节点的CWE层次树",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadTreeRegistry()
		if err != nil {
			return err
		}

		rootID, _ := cweskills.ParseCWEID(args[0])
		tree := cweskills.BuildTree(registry, rootID)
		if tree == nil {
			return fmt.Errorf("无法构建以 %s 为根的树", cweskills.FormatCWEIDFromInt(rootID))
		}

		if outputFormat == "json" {
			return printJSON(cmd, treeNodeToMap(tree))
		}

		printTreeNode(cmd, tree, 0)
		return nil
	},
}

// treeForestCmd 构建所有pillar节点的森林
var treeForestCmd = &cobra.Command{
	Use:   "forest",
	Short: "构建所有顶层(Pillar)节点的CWE森林",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadTreeRegistry()
		if err != nil {
			return err
		}

		forest := cweskills.BuildForest(registry)
		if outputFormat == "json" {
			trees := make([]interface{}, 0, len(forest))
			for _, t := range forest {
				trees = append(trees, treeNodeToMap(t))
			}
			return printJSON(cmd, map[string]interface{}{
				"count": len(forest),
				"trees": trees,
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "CWE森林 (%d 棵树):\n\n", len(forest))
		for i, t := range forest {
			fmt.Fprintf(cmd.OutOrStdout(), "--- 树 %d ---\n", i+1)
			printTreeNode(cmd, t, 0)
			fmt.Fprintln(cmd.OutOrStdout())
		}
		return nil
	},
}

// treeViewCmd 构建指定视图的树
var treeViewCmd = &cobra.Command{
	Use:   "view [VIEW-ID]",
	Short: "构建指定视图的CWE层次树",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadTreeRegistry()
		if err != nil {
			return err
		}

		viewID, _ := cweskills.ParseCWEID(args[0])
		tree := cweskills.BuildViewTree(registry, viewID)
		if tree == nil {
			return fmt.Errorf("无法构建视图 %s 的树", cweskills.FormatCWEIDFromInt(viewID))
		}

		if outputFormat == "json" {
			return printJSON(cmd, treeNodeToMap(tree))
		}

		printTreeNode(cmd, tree, 0)
		return nil
	},
}

// treePathCmd 查找从根到指定节点的路径
var treePathCmd = &cobra.Command{
	Use:   "path [CWE-ID]",
	Short: "查找从根到指定CWE的路径",
	Long:  "查找从层次树根到指定CWE节点的路径。需要先指定一个根节点。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadTreeRegistry()
		if err != nil {
			return err
		}

		id, _ := cweskills.ParseCWEID(args[0])
		rootID, _ := cmd.Flags().GetInt("root")
		if rootID <= 0 {
			// 自动查找根（pillar祖先）
			nav := cweskills.NewNavigator(registry)
			ancestors := nav.Ancestors(id)
			if len(ancestors) > 0 {
				rootID = ancestors[len(ancestors)-1].ID
			} else {
				rootID = id
			}
		}

		tree := cweskills.BuildTree(registry, rootID)
		if tree == nil {
			return fmt.Errorf("无法构建树")
		}

		node := tree.Find(id)
		if node == nil {
			return fmt.Errorf("在树中未找到 %s", cweskills.FormatCWEIDFromInt(id))
		}

		path := node.Path()
		if path == nil {
			return fmt.Errorf("在树中未找到 %s", cweskills.FormatCWEIDFromInt(id))
		}

		if outputFormat == "json" {
			ids := make([]int, len(path))
			for i, n := range path {
				ids[i] = n.CWE.ID
			}
			return printJSON(cmd, map[string]interface{}{
				"cwe_id": cweskills.FormatCWEIDFromInt(id),
				"root":   cweskills.FormatCWEIDFromInt(rootID),
				"path":   ids,
				"depth":  len(path) - 1,
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "从 %s 到 %s 的路径 (%d 步):\n",
			cweskills.FormatCWEIDFromInt(rootID), cweskills.FormatCWEIDFromInt(id), len(path)-1)
		for i, n := range path {
			fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s - %s\n", i+1,
				cweskills.FormatCWEIDFromInt(n.CWE.ID), n.CWE.Name)
		}
		return nil
	},
}

// treeLeavesCmd 列出所有叶子节点
var treeLeavesCmd = &cobra.Command{
	Use:   "leaves [ROOT-ID]",
	Short: "列出指定根节点下的所有叶子弱点",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadTreeRegistry()
		if err != nil {
			return err
		}

		rootID, _ := cweskills.ParseCWEID(args[0])
		tree := cweskills.BuildTree(registry, rootID)
		if tree == nil {
			return fmt.Errorf("无法构建以 %s 为根的树", cweskills.FormatCWEIDFromInt(rootID))
		}

		leaves := tree.LeafNodes()
		if outputFormat == "json" {
			type leafEntry struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}
			entries := make([]leafEntry, 0, len(leaves))
			for _, l := range leaves {
				entries = append(entries, leafEntry{ID: l.CWE.ID, Name: l.CWE.Name})
			}
			return printJSON(cmd, map[string]interface{}{
				"root":   cweskills.FormatCWEIDFromInt(rootID),
				"leaves": entries,
				"count":  len(entries),
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s 的叶子节点 (%d 项):\n",
			cweskills.FormatCWEIDFromInt(rootID), len(leaves))
		for _, l := range leaves {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s\n",
				cweskills.FormatCWEIDFromInt(l.CWE.ID), l.CWE.Name)
		}
		return nil
	},
}

// treeNodeToMap 将TreeNode转为map用于JSON输出
func treeNodeToMap(node *cweskills.TreeNode) map[string]interface{} {
	result := map[string]interface{}{
		"id":      node.CWE.ID,
		"name":    node.CWE.Name,
		"depth":   node.Depth,
		"is_leaf": node.IsLeaf(),
	}
	if len(node.Children) > 0 {
		children := make([]interface{}, 0, len(node.Children))
		for _, c := range node.Children {
			children = append(children, treeNodeToMap(c))
		}
		result["children"] = children
	}
	return result
}

func printTreeNode(cmd *cobra.Command, node *cweskills.TreeNode, indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Fprintf(cmd.OutOrStdout(), "%s%s - %s\n", prefix,
		cweskills.FormatCWEIDFromInt(node.CWE.ID), node.CWE.Name)
	for _, child := range node.Children {
		printTreeNode(cmd, child, indent+1)
	}
}

func init() {
	treeCmd.PersistentFlags().StringVarP(&treeXMLPath, "xml", "x", "", "CWE XML目录文件路径")
	treePathCmd.Flags().Int("root", 0, "根节点ID（默认自动查找）")

	treeCmd.AddCommand(treeBuildCmd)
	treeCmd.AddCommand(treeForestCmd)
	treeCmd.AddCommand(treeViewCmd)
	treeCmd.AddCommand(treePathCmd)
	treeCmd.AddCommand(treeLeavesCmd)
	rootCmd.AddCommand(treeCmd)
}
