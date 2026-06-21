package main

import (
	"fmt"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

var navXMLPath string

// navCmd 导航命令组
var navCmd = &cobra.Command{
	Use:   "nav",
	Short: "本地关系导航",
	Long: `基于本地XML数据的关系导航，提供比API更丰富的关系查询。

所有nav子命令都需要通过 --xml 参数指定CWE XML目录文件。

支持的关系类型：
  - 层级关系: parents, children, ancestors, descendants, siblings
  - 对等关系: peers, can-also-be
  - 顺序关系: precede, follow
  - 依赖关系: requires, required-by
  - 复合关系: chain-members, composite-members
  - 路径查询: shortest-path, is-ancestor, is-related, depth

示例：
  cwe nav parents CWE-79 --xml cwec_latest.xml
  cwe nav siblings CWE-79 --xml cwec_latest.xml
  cwe nav shortest-path CWE-79 CWE-1 --xml cwec_latest.xml
  cwe nav is-ancestor CWE-1 CWE-79 --xml cwec_latest.xml`,
}

func loadNavRegistry() (*cweskills.Registry, *cweskills.Navigator, error) {
	if navXMLPath == "" {
		return nil, nil, fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
	}
	parser := cweskills.NewXMLParser()
	registry, err := parser.ParseFile(navXMLPath)
	if err != nil {
		return nil, nil, fmt.Errorf("解析XML文件失败: %w", err)
	}
	registry.BuildIndexes()
	nav := cweskills.NewNavigator(registry)
	return registry, nav, nil
}

func parseIDArg(arg string) (int, error) {
	return cweskills.ParseCWEID(arg)
}

// navParentsCmd
var navParentsCmd = &cobra.Command{
	Use:   "parents [CWE-ID]", Short: "查询父级弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "父级", id, nav.Parents(id), r)
	},
}

// navChildrenCmd
var navChildrenCmd = &cobra.Command{
	Use:   "children [CWE-ID]", Short: "查询子级弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "子级", id, nav.Children(id), r)
	},
}

// navAncestorsCmd
var navAncestorsCmd = &cobra.Command{
	Use:   "ancestors [CWE-ID]", Short: "查询所有祖先弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "祖先", id, nav.Ancestors(id), r)
	},
}

// navDescendantsCmd
var navDescendantsCmd = &cobra.Command{
	Use:   "descendants [CWE-ID]", Short: "查询所有后代弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "后代", id, nav.Descendants(id), r)
	},
}

// navSiblingsCmd
var navSiblingsCmd = &cobra.Command{
	Use:   "siblings [CWE-ID]", Short: "查询同级弱点（同一父级的其他子级）", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "同级", id, nav.Siblings(id), r)
	},
}

// navPeersCmd
var navPeersCmd = &cobra.Command{
	Use:   "peers [CWE-ID]", Short: "查询对等弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "对等", id, nav.Peers(id), r)
	},
}

// navPrecedeCmd
var navPrecedeCmd = &cobra.Command{
	Use:   "precede [CWE-ID]", Short: "查询此弱点可以前置的弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "可前置", id, nav.CanPrecede(id), r)
	},
}

// navFollowCmd
var navFollowCmd = &cobra.Command{
	Use:   "follow [CWE-ID]", Short: "查询此弱点可以跟随的弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "可跟随", id, nav.CanFollow(id), r)
	},
}

// navRequiresCmd
var navRequiresCmd = &cobra.Command{
	Use:   "requires [CWE-ID]", Short: "查询此弱点所依赖的弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "依赖", id, nav.Requires(id), r)
	},
}

// navRequiredByCmd
var navRequiredByCmd = &cobra.Command{
	Use:   "required-by [CWE-ID]", Short: "查询依赖此弱点的弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "被依赖", id, nav.RequiredBy(id), r)
	},
}

// navCanAlsoBeCmd
var navCanAlsoBeCmd = &cobra.Command{
	Use:   "can-also-be [CWE-ID]", Short: "查询此弱点也可以是的弱点", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "也可以是", id, nav.CanAlsoBe(id), r)
	},
}

// navChainMembersCmd
var navChainMembersCmd = &cobra.Command{
	Use:   "chain-members [CWE-ID]", Short: "查询链式弱点的成员", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "链式成员", id, nav.ChainMembers(id), r)
	},
}

// navCompositeMembersCmd
var navCompositeMembersCmd = &cobra.Command{
	Use:   "composite-members [CWE-ID]", Short: "查询复合弱点的成员", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, nav, err := loadNavRegistry()
		if err != nil { return err }
		id, _ := parseIDArg(args[0])
		return printNavResults(cmd, "复合成员", id, nav.CompositeMembers(id), r)
	},
}

// navShortestPathCmd
var navShortestPathCmd = &cobra.Command{
	Use:   "shortest-path <FROM> <TO>", Short: "查找两个CWE之间的最短路径", Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, nav, err := loadNavRegistry()
		if err != nil { return err }
		from, _ := parseIDArg(args[0])
		to, _ := parseIDArg(args[1])

		path := nav.ShortestPath(from, to)
		if path == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "%s 和 %s 之间没有路径\n",
				cweskills.FormatCWEIDFromInt(from), cweskills.FormatCWEIDFromInt(to))
			return nil
		}

		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"from": cweskills.FormatCWEIDFromInt(from),
				"to":   cweskills.FormatCWEIDFromInt(to),
				"path": path,
				"depth": len(path) - 1,
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "最短路径 (%d 步):\n", len(path)-1)
		for i, id := range path {
			fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s\n", i+1, cweskills.FormatCWEIDFromInt(id))
		}
		return nil
	},
}

// navIsAncestorCmd
var navIsAncestorCmd = &cobra.Command{
	Use:   "is-ancestor <ANCESTOR> <DESCENDANT>", Short: "检查是否为祖先关系", Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, nav, err := loadNavRegistry()
		if err != nil { return err }
		ancestor, _ := parseIDArg(args[0])
		descendant, _ := parseIDArg(args[1])

		result := nav.IsAncestorOf(ancestor, descendant)
		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"ancestor":   cweskills.FormatCWEIDFromInt(ancestor),
				"descendant": cweskills.FormatCWEIDFromInt(descendant),
				"is_ancestor": result,
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s 的祖先\n",
			cweskills.FormatCWEIDFromInt(ancestor), boolStr(result, "是", "不是"), cweskills.FormatCWEIDFromInt(descendant))
		return nil
	},
}

// navIsRelatedCmd
var navIsRelatedCmd = &cobra.Command{
	Use:   "is-related <CWE-ID1> <CWE-ID2>", Short: "检查两个CWE是否有关系", Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, nav, err := loadNavRegistry()
		if err != nil { return err }
		a, _ := parseIDArg(args[0])
		b, _ := parseIDArg(args[1])

		result := nav.IsRelated(a, b)
		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"cwe_id_a":   cweskills.FormatCWEIDFromInt(a),
				"cwe_id_b":   cweskills.FormatCWEIDFromInt(b),
				"is_related": result,
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s 和 %s %s\n",
			cweskills.FormatCWEIDFromInt(a), cweskills.FormatCWEIDFromInt(b), boolStr(result, "有关联", "无关联"))
		return nil
	},
}

// navDepthCmd
var navDepthCmd = &cobra.Command{
	Use:   "depth <FROM> <TO>", Short: "计算两个CWE之间的关系深度", Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, nav, err := loadNavRegistry()
		if err != nil { return err }
		from, _ := parseIDArg(args[0])
		to, _ := parseIDArg(args[1])

		depth := nav.RelationshipDepth(from, to)
		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"from":  cweskills.FormatCWEIDFromInt(from),
				"to":    cweskills.FormatCWEIDFromInt(to),
				"depth": depth,
			})
		}

		if depth < 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "%s 和 %s 之间没有关系\n",
				cweskills.FormatCWEIDFromInt(from), cweskills.FormatCWEIDFromInt(to))
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s 到 %s 的关系深度: %d\n",
				cweskills.FormatCWEIDFromInt(from), cweskills.FormatCWEIDFromInt(to), depth)
		}
		return nil
	},
}

func printNavResults(cmd *cobra.Command, relType string, id int, cwes []*cweskills.CWE, r *cweskills.Registry) error {
	if outputFormat == "json" {
		type entry struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		entries := make([]entry, 0, len(cwes))
		for _, c := range cwes {
			entries = append(entries, entry{ID: c.ID, Name: c.Name})
		}
		return printJSON(cmd, map[string]interface{}{
			"cwe_id":  cweskills.FormatCWEIDFromInt(id),
			"type":    relType,
			"results": entries,
			"count":   len(entries),
		})
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s 的 %s (%d 项):\n", cweskills.FormatCWEIDFromInt(id), relType, len(cwes))
	for _, c := range cwes {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s\n", cweskills.FormatCWEIDFromInt(c.ID), c.Name)
	}
	return nil
}

func boolStr(b bool, trueStr, falseStr string) string {
	if b { return trueStr }
	return falseStr
}

func init() {
	navCmd.PersistentFlags().StringVarP(&navXMLPath, "xml", "x", "", "CWE XML目录文件路径")

	navCmd.AddCommand(navParentsCmd)
	navCmd.AddCommand(navChildrenCmd)
	navCmd.AddCommand(navAncestorsCmd)
	navCmd.AddCommand(navDescendantsCmd)
	navCmd.AddCommand(navSiblingsCmd)
	navCmd.AddCommand(navPeersCmd)
	navCmd.AddCommand(navPrecedeCmd)
	navCmd.AddCommand(navFollowCmd)
	navCmd.AddCommand(navRequiresCmd)
	navCmd.AddCommand(navRequiredByCmd)
	navCmd.AddCommand(navCanAlsoBeCmd)
	navCmd.AddCommand(navChainMembersCmd)
	navCmd.AddCommand(navCompositeMembersCmd)
	navCmd.AddCommand(navShortestPathCmd)
	navCmd.AddCommand(navIsAncestorCmd)
	navCmd.AddCommand(navIsRelatedCmd)
	navCmd.AddCommand(navDepthCmd)

	rootCmd.AddCommand(navCmd)
}
