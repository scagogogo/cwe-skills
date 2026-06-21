package main

import (
	"fmt"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

var registryXMLPath string

// registryCmd 注册表操作命令组
var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "本地注册表操作",
	Long: `对本地CWE注册表进行操作，包括加载XML、查询、导出等。

所有registry子命令都需要通过 --xml 参数指定CWE XML目录文件，
或者通过 --json 参数指定之前导出的JSON文件。

示例：
  cwe registry load --xml cwec_latest.xml
  cwe registry get CWE-79 --xml cwec_latest.xml
  cwe registry parents CWE-79 --xml cwec_latest.xml
  cwe registry export --xml cwec_latest.xml --format json
  cwe registry import data.json`,
}

// registryLoadCmd 加载XML并显示概要
var registryLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "加载XML目录并显示概要",
	Long:  "加载MITRE CWE XML目录文件，构建索引并显示概要信息。",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(cmd, map[string]interface{}{
				"weaknesses":  registry.Size(),
				"categories":  registry.CategoryCount(),
				"views":       registry.ViewCount(),
				"compounds":   registry.CompoundElementCount(),
				"indexed":     registry.IndexesBuilt(),
			})
		}

		fmt.Fprintf(cmd.OutOrStdout(), "已加载CWE注册表:\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  弱点:     %d\n", registry.Size())
		fmt.Fprintf(cmd.OutOrStdout(), "  类别:     %d\n", registry.CategoryCount())
		fmt.Fprintf(cmd.OutOrStdout(), "  视图:     %d\n", registry.ViewCount())
		fmt.Fprintf(cmd.OutOrStdout(), "  复合元素: %d\n", registry.CompoundElementCount())
		fmt.Fprintf(cmd.OutOrStdout(), "  索引:     %v\n", registry.IndexesBuilt())
		return nil
	},
}

// registryGetCmd 获取单个CWE条目
var registryGetCmd = &cobra.Command{
	Use:   "get [CWE-ID]",
	Short: "获取CWE条目详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		cwe, ok := registry.Get(id)
		if !ok {
			return fmt.Errorf("CWE-%d 不存在于注册表中", id)
		}

		if outputFormat == "json" {
			return printJSON(cmd, cwe)
		}

		printCWEDetail(cmd, cwe)
		return nil
	},
}

// registryContainsCmd 检查CWE ID是否在注册表中
var registryContainsCmd = &cobra.Command{
	Use:   "contains [CWE-ID...]",
	Short: "检查CWE ID是否存在于注册表",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个CWE ID")
		}

		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		type containsResult struct {
			CWEID   string `json:"cwe_id"`
			Exists  bool   `json:"exists"`
			Type    string `json:"type,omitempty"`
		}

		results := make([]containsResult, 0, len(args))
		for _, input := range args {
			id, err := cweskills.ParseCWEID(input)
			if err != nil {
				results = append(results, containsResult{CWEID: input, Exists: false})
				continue
			}

			exists := registry.Contains(id)
			var entryType string
			if exists {
				if cwe, ok := registry.Get(id); ok {
					entryType = cwe.CWEType
				}
			}
			results = append(results, containsResult{
				CWEID:  cweskills.FormatCWEIDFromInt(id),
				Exists: exists,
				Type:   entryType,
			})
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		for _, r := range results {
			if r.Exists {
				fmt.Fprintf(cmd.OutOrStdout(), "%s ✓ 存在 (%s)\n", r.CWEID, r.Type)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s ✗ 不存在\n", r.CWEID)
			}
		}
		return nil
	},
}

// registryListViewsCmd 列出所有视图
var registryListViewsCmd = &cobra.Command{
	Use:   "list-views",
	Short: "列出注册表中的所有视图",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		views := registry.GetAllViews()
		if outputFormat == "json" {
			return printJSON(cmd, views)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "视图 (%d 项):\n", len(views))
		for _, v := range views {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s [%s]\n",
				cweskills.FormatCWEIDFromInt(v.ID), v.Name, v.Type)
		}
		return nil
	},
}

// registryListCategoriesCmd 列出所有类别
var registryListCategoriesCmd = &cobra.Command{
	Use:   "list-categories",
	Short: "列出注册表中的所有类别",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		cats := registry.GetAllCategories()
		if outputFormat == "json" {
			return printJSON(cmd, cats)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "类别 (%d 项):\n", len(cats))
		for _, c := range cats {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s\n",
				cweskills.FormatCWEIDFromInt(c.ID), c.Name)
		}
		return nil
	},
}

// registryParentsCmd 查询本地父级关系
var registryParentsCmd = &cobra.Command{
	Use:   "parents [CWE-ID]",
	Short: "查询本地注册表中的父级关系",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		parentIDs := registry.GetParentIDs(id)
		return printIDResults(cmd, "父级弱点", id, parentIDs, registry)
	},
}

// registryChildrenCmd 查询本地子级关系
var registryChildrenCmd = &cobra.Command{
	Use:   "children [CWE-ID]",
	Short: "查询本地注册表中的子级关系",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		childIDs := registry.GetChildIDs(id)
		return printIDResults(cmd, "子级弱点", id, childIDs, registry)
	},
}

// registryAncestorsCmd 查询本地祖先关系
var registryAncestorsCmd = &cobra.Command{
	Use:   "ancestors [CWE-ID]",
	Short: "查询本地注册表中的所有祖先关系",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		ancestorIDs := registry.GetAncestorIDs(id)
		return printIDResults(cmd, "祖先弱点", id, ancestorIDs, registry)
	},
}

// registryDescendantsCmd 查询本地后代关系
var registryDescendantsCmd = &cobra.Command{
	Use:   "descendants [CWE-ID]",
	Short: "查询本地注册表中的所有后代关系",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		descIDs := registry.GetDescendantIDs(id)
		return printIDResults(cmd, "后代弱点", id, descIDs, registry)
	},
}

// registryPeersCmd 查询本地对等关系
var registryPeersCmd = &cobra.Command{
	Use:   "peers [CWE-ID]",
	Short: "查询本地注册表中的对等关系",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		peerIDs := registry.GetPeerIDs(id)
		return printIDResults(cmd, "对等弱点", id, peerIDs, registry)
	},
}

// registryViewMembersCmd 查询视图成员
var registryViewMembersCmd = &cobra.Command{
	Use:   "view-members [VIEW-ID]",
	Short: "查询视图的成员列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		viewID, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		memberIDs := registry.GetViewMembers(viewID)
		return printIDResults(cmd, "视图成员", viewID, memberIDs, registry)
	},
}

// registryCategoryMembersCmd 查询类别成员
var registryCategoryMembersCmd = &cobra.Command{
	Use:   "category-members [CATEGORY-ID]",
	Short: "查询类别的成员列表",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		catID, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		memberIDs := registry.GetCategoryMembers(catID)
		return printIDResults(cmd, "类别成员", catID, memberIDs, registry)
	},
}

// registryMemberOfCmd 查询CWE所属的类别/视图
var registryMemberOfCmd = &cobra.Command{
	Use:   "member-of [CWE-ID]",
	Short: "查询CWE所属的类别和视图",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		id, err := cweskills.ParseCWEID(args[0])
		if err != nil {
			return err
		}

		memberOfIDs := registry.GetMemberOfIDs(id)
		return printIDResults(cmd, "所属类别/视图", id, memberOfIDs, registry)
	},
}

// registryExportCmd 导出注册表
var registryExportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出注册表为JSON或CSV",
	Long:  "将加载的XML注册表导出为JSON或CSV格式。",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := loadRegistry()
		if err != nil {
			return err
		}

		exportFormat, _ := cmd.Flags().GetString("format")
		exportOutput, _ := cmd.Flags().GetString("output-file")

		var data []byte
		switch exportFormat {
		case "json":
			data, err = registry.ExportJSON()
		case "csv":
			data, err = registry.ExportCSV()
		default:
			return fmt.Errorf("不支持的导出格式: %s (支持: json, csv)", exportFormat)
		}
		if err != nil {
			return fmt.Errorf("导出失败: %w", err)
		}

		if exportOutput != "" {
			return writeFile(exportOutput, data)
		}

		fmt.Fprint(cmd.OutOrStdout(), string(data))
		return nil
	},
}

// 辅助函数

func loadRegistry() (*cweskills.Registry, error) {
	if registryXMLPath == "" {
		return nil, fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
	}
	parser := cweskills.NewXMLParser()
	registry, err := parser.ParseFile(registryXMLPath)
	if err != nil {
		return nil, fmt.Errorf("解析XML文件失败: %w", err)
	}
	return registry, nil
}

func printCWEDetail(cmd *cobra.Command, cwe *cweskills.CWE) {
	fmt.Fprintf(cmd.OutOrStdout(), "=== %s ===\n", cweskills.FormatCWEIDFromInt(cwe.ID))
	fmt.Fprintf(cmd.OutOrStdout(), "  名称:     %s\n", cwe.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "  抽象层级: %s\n", cwe.Abstraction)
	fmt.Fprintf(cmd.OutOrStdout(), "  状态:     %s\n", cwe.Status)
	if cwe.Structure != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  结构:     %s\n", cwe.Structure)
	}
	if cwe.Description != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  描述:     %s\n", cwe.Description)
	}
	if cwe.LikelihoodOfExploit != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  利用可能性: %s\n", cwe.LikelihoodOfExploit)
	}
	if len(cwe.Relationships) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "  关系:     %d 项\n", len(cwe.Relationships))
	}
	if len(cwe.CommonConsequences) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "  后果:     %d 项\n", len(cwe.CommonConsequences))
	}
}

func printIDResults(cmd *cobra.Command, relType string, id int, ids []int, registry *cweskills.Registry) error {
	type idEntry struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
	}

	entries := make([]idEntry, 0, len(ids))
	for _, cid := range ids {
		entry := idEntry{ID: cid}
		if cwe, ok := registry.Get(cid); ok {
			entry.Name = cwe.Name
		}
		entries = append(entries, entry)
	}

	if outputFormat == "json" {
		return printJSON(cmd, map[string]interface{}{
			"cwe_id":  cweskills.FormatCWEIDFromInt(id),
			"type":    relType,
			"results": entries,
			"count":   len(entries),
		})
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s 的 %s (%d 项):\n", cweskills.FormatCWEIDFromInt(id), relType, len(entries))
	for _, e := range entries {
		if e.Name != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s\n", cweskills.FormatCWEIDFromInt(e.ID), e.Name)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", cweskills.FormatCWEIDFromInt(e.ID))
		}
	}
	return nil
}

func writeFile(path string, data []byte) error {
	return fmt.Errorf("写入文件功能: %s (%d bytes)", path, len(data))
}

func init() {
	registryCmd.PersistentFlags().StringVarP(&registryXMLPath, "xml", "x", "", "CWE XML目录文件路径")

	registryExportCmd.Flags().String("format", "json", "导出格式 (json|csv)")
	registryExportCmd.Flags().String("output-file", "", "输出文件路径（默认输出到stdout）")

	registryCmd.AddCommand(registryLoadCmd)
	registryCmd.AddCommand(registryGetCmd)
	registryCmd.AddCommand(registryContainsCmd)
	registryCmd.AddCommand(registryListViewsCmd)
	registryCmd.AddCommand(registryListCategoriesCmd)
	registryCmd.AddCommand(registryParentsCmd)
	registryCmd.AddCommand(registryChildrenCmd)
	registryCmd.AddCommand(registryAncestorsCmd)
	registryCmd.AddCommand(registryDescendantsCmd)
	registryCmd.AddCommand(registryPeersCmd)
	registryCmd.AddCommand(registryViewMembersCmd)
	registryCmd.AddCommand(registryCategoryMembersCmd)
	registryCmd.AddCommand(registryMemberOfCmd)
	registryCmd.AddCommand(registryExportCmd)
	rootCmd.AddCommand(registryCmd)
}
