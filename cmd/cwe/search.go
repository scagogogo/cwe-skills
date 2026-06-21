package main

import (
	"fmt"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

var (
	xmlFilePath       string
	searchKeyword     string
	searchAbstr       string
	searchStatus      string
	searchLikelihood  string
	searchStructure   string
	searchScope       string
	searchTopLevel    bool
	searchBaseOnly    bool
	searchChains      bool
	searchComposites  bool
	searchSort        string
	searchGroupBy     string
	searchDedup       bool
)

// searchCmd 本地搜索CWE条目
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "搜索CWE条目",
	Long: `从本地XML数据源搜索CWE条目。

需要先通过 --xml 参数指定MITRE CWE XML目录文件路径。
可以从 https://cwe.mitre.org/data/xml.html 下载。

示例：
  cwe search --xml cwec_latest.xml --keyword "Injection"
  cwe search --xml cwec_latest.xml --abstraction Base
  cwe search --xml cwec_latest.xml --status Stable --likelihood High
  cwe search --xml cwec_latest.xml --top-level
  cwe search --xml cwec_latest.xml --scope Confidentiality
  cwe search --xml cwec_latest.xml --sort name
  cwe search --xml cwec_latest.xml --group-by abstraction`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if xmlFilePath == "" {
			return fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
		}

		parser := cweskills.NewXMLParser()
		registry, err := parser.ParseFile(xmlFilePath)
		if err != nil {
			return fmt.Errorf("解析XML文件失败: %w", err)
		}

		var results []*cweskills.CWE

		if searchKeyword != "" {
			results = cweskills.FindByKeyword(registry, searchKeyword)
		} else if searchAbstr != "" {
			abstr, err := cweskills.ParseAbstraction(searchAbstr)
			if err != nil {
				return fmt.Errorf("无效的抽象层级: %w", err)
			}
			results = cweskills.FindByAbstraction(registry, abstr)
		} else if searchStatus != "" {
			st, err := cweskills.ParseStatus(searchStatus)
			if err != nil {
				return fmt.Errorf("无效的状态: %w", err)
			}
			results = cweskills.FindByStatus(registry, st)
		} else if searchLikelihood != "" {
			lh, err := cweskills.ParseLikelihoodOfExploit(searchLikelihood)
			if err != nil {
				return fmt.Errorf("无效的利用可能性: %w", err)
			}
			results = cweskills.FindByLikelihood(registry, lh)
		} else if searchStructure != "" {
			st, err := cweskills.ParseStructure(searchStructure)
			if err != nil {
				return fmt.Errorf("无效的结构类型: %w", err)
			}
			results = cweskills.FindByStructure(registry, st)
		} else if searchScope != "" {
			sc, err := cweskills.ParseConsequenceScope(searchScope)
			if err != nil {
				return fmt.Errorf("无效的后果范围: %w", err)
			}
			results = cweskills.FindByConsequenceScope(registry, sc)
		} else if searchTopLevel {
			results = cweskills.FindTopLevel(registry)
		} else if searchBaseOnly {
			results = cweskills.FindBaseWeaknesses(registry)
		} else if searchChains {
			results = cweskills.FindChains(registry)
		} else if searchComposites {
			results = cweskills.FindComposites(registry)
		} else {
			results = registry.GetAll()
		}

		if searchSort != "" {
			switch searchSort {
			case "id":
				cweskills.SortByID(results)
			case "name":
				cweskills.SortByName(results)
			case "abstraction":
				cweskills.SortByAbstraction(results)
			default:
				return fmt.Errorf("不支持的排序字段: %s (支持: id, name, abstraction)", searchSort)
			}
		}

		if searchDedup {
			results = cweskills.Deduplicate(results)
		}

		if searchGroupBy != "" {
			return printGroupedResults(cmd, results)
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "找到 %d 个CWE条目:\n\n", len(results))
		for _, cwe := range results {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s [%s, %s]\n",
				cweskills.FormatCWEIDFromInt(cwe.ID), cwe.Name, cwe.Abstraction, cwe.Status)
		}
		return nil
	},
}

// filterCmd 多条件过滤
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "多条件过滤CWE条目",
	Long: `使用多个条件组合过滤CWE条目。所有条件之间为AND关系。

示例：
  cwe filter --xml cwec_latest.xml --abstraction Base --status Stable --keyword Injection
  cwe filter --xml cwec_latest.xml --likelihood High --scope Confidentiality`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if xmlFilePath == "" {
			return fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
		}

		parser := cweskills.NewXMLParser()
		registry, err := parser.ParseFile(xmlFilePath)
		if err != nil {
			return fmt.Errorf("解析XML文件失败: %w", err)
		}

		all := registry.GetAll()
		opts := cweskills.FilterOption{}

		if v, _ := cmd.Flags().GetString("keyword"); v != "" {
			opts.Keyword = v
		}
		if v, _ := cmd.Flags().GetString("abstraction"); v != "" {
			abstr, err := cweskills.ParseAbstraction(v)
			if err != nil {
				return fmt.Errorf("无效的抽象层级: %w", err)
			}
			opts.Abstraction = abstr
		}
		if v, _ := cmd.Flags().GetString("status"); v != "" {
			st, err := cweskills.ParseStatus(v)
			if err != nil {
				return fmt.Errorf("无效的状态: %w", err)
			}
			opts.Status = st
		}
		if v, _ := cmd.Flags().GetString("likelihood"); v != "" {
			lh, err := cweskills.ParseLikelihoodOfExploit(v)
			if err != nil {
				return fmt.Errorf("无效的利用可能性: %w", err)
			}
			opts.Likelihood = lh
		}
		if v, _ := cmd.Flags().GetString("scope"); v != "" {
			sc, err := cweskills.ParseConsequenceScope(v)
			if err != nil {
				return fmt.Errorf("无效的后果范围: %w", err)
			}
			opts.Scope = sc
		}
		if v, _ := cmd.Flags().GetString("structure"); v != "" {
			st, err := cweskills.ParseStructure(v)
			if err != nil {
				return fmt.Errorf("无效的结构类型: %w", err)
			}
			opts.Structure = st
		}

		results := cweskills.Filter(all, opts)

		if v, _ := cmd.Flags().GetString("sort"); v != "" {
			switch v {
			case "id":
				cweskills.SortByID(results)
			case "name":
				cweskills.SortByName(results)
			case "abstraction":
				cweskills.SortByAbstraction(results)
			}
		}

		if v, _ := cmd.Flags().GetString("group-by"); v != "" {
			return printGroupedResults(cmd, results)
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "过滤结果 (%d 项):\n\n", len(results))
		for _, cwe := range results {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s [%s, %s]\n",
				cweskills.FormatCWEIDFromInt(cwe.ID), cwe.Name, cwe.Abstraction, cwe.Status)
		}
		return nil
	},
}

// statsCmd 统计CWE数据
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "统计CWE数据",
	Long: `从本地XML数据源统计CWE数据分布。

示例：
  cwe stats --xml cwec_latest.xml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if xmlFilePath == "" {
			return fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
		}

		parser := cweskills.NewXMLParser()
		registry, err := parser.ParseFile(xmlFilePath)
		if err != nil {
			return fmt.Errorf("解析XML文件失败: %w", err)
		}

		stats := cweskills.ComputeStatistics(registry)

		if outputFormat == "json" {
			return printJSON(cmd, stats)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "CWE数据统计:\n\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  总条目数:     %d\n", stats.TotalCount)
		fmt.Fprintf(cmd.OutOrStdout(), "  类别数:       %d\n", stats.CategoryCount)
		fmt.Fprintf(cmd.OutOrStdout(), "  视图数:       %d\n", stats.ViewCount)

		fmt.Fprintf(cmd.OutOrStdout(), "\n抽象层级分布:\n")
		for abstr, count := range stats.ByAbstraction {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s: %d\n", abstr, count)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "\n状态分布:\n")
		for status, count := range stats.ByStatus {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s: %d\n", status, count)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "\n利用可能性分布:\n")
		for likelihood, count := range stats.ByLikelihood {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s: %d\n", likelihood, count)
		}
		return nil
	},
}

func printGroupedResults(cmd *cobra.Command, results []*cweskills.CWE) error {
	if outputFormat == "json" {
		// Convert typed group maps to string-keyed maps for JSON
		groupStr := make(map[string]interface{})
		switch searchGroupBy {
		case "abstraction":
			for k, v := range cweskills.GroupByAbstraction(results) {
				groupStr[string(k)] = v
			}
		case "status":
			for k, v := range cweskills.GroupByStatus(results) {
				groupStr[string(k)] = v
			}
		case "likelihood":
			for k, v := range cweskills.GroupByLikelihood(results) {
				groupStr[string(k)] = v
			}
		default:
			return fmt.Errorf("不支持的分组字段: %s (支持: abstraction, status, likelihood)", searchGroupBy)
		}
		return printJSON(cmd, groupStr)
	}

	switch searchGroupBy {
	case "abstraction":
		for group, items := range cweskills.GroupByAbstraction(results) {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s (%d 项):\n", group, len(items))
			for _, cwe := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s\n", cweskills.FormatCWEIDFromInt(cwe.ID), cwe.Name)
			}
		}
	case "status":
		for group, items := range cweskills.GroupByStatus(results) {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s (%d 项):\n", group, len(items))
			for _, cwe := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s\n", cweskills.FormatCWEIDFromInt(cwe.ID), cwe.Name)
			}
		}
	case "likelihood":
		for group, items := range cweskills.GroupByLikelihood(results) {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s (%d 项):\n", group, len(items))
			for _, cwe := range items {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s\n", cweskills.FormatCWEIDFromInt(cwe.ID), cwe.Name)
			}
		}
	default:
		return fmt.Errorf("不支持的分组字段: %s (支持: abstraction, status, likelihood)", searchGroupBy)
	}
	return nil
}

func init() {
	// search flags
	searchCmd.Flags().StringVarP(&xmlFilePath, "xml", "x", "", "CWE XML目录文件路径")
	searchCmd.Flags().StringVarP(&searchKeyword, "keyword", "k", "", "按关键字搜索")
	searchCmd.Flags().StringVarP(&searchAbstr, "abstraction", "a", "", "按抽象层级搜索 (Pillar/Class/Base/Variant)")
	searchCmd.Flags().StringVarP(&searchStatus, "status", "s", "", "按状态搜索 (Stable/Draft/Deprecated)")
	searchCmd.Flags().StringVarP(&searchLikelihood, "likelihood", "l", "", "按利用可能性搜索 (High/Medium/Low)")
	searchCmd.Flags().StringVarP(&searchStructure, "structure", "t", "", "按结构类型搜索 (Simple/Chain/Composite)")
	searchCmd.Flags().StringVarP(&searchScope, "scope", "", "", "按后果范围搜索 (Confidentiality/Integrity/Availability)")
	searchCmd.Flags().BoolVar(&searchTopLevel, "top-level", false, "只显示顶层(Pillar)弱点")
	searchCmd.Flags().BoolVar(&searchBaseOnly, "base-weaknesses", false, "只显示基础(Base)弱点")
	searchCmd.Flags().BoolVar(&searchChains, "chains", false, "只显示链式弱点")
	searchCmd.Flags().BoolVar(&searchComposites, "composites", false, "只显示复合弱点")
	searchCmd.Flags().StringVar(&searchSort, "sort", "", "排序字段 (id/name/abstraction)")
	searchCmd.Flags().StringVar(&searchGroupBy, "group-by", "", "分组字段 (abstraction/status/likelihood)")
	searchCmd.Flags().BoolVar(&searchDedup, "dedup", false, "去重")

	// filter flags (use local vars)
	var filterKeyword, filterAbstr, filterStatus, filterLikelihood, filterStructure, filterScope string
	var filterSort, filterGroupBy string
	filterCmd.Flags().StringVarP(&xmlFilePath, "xml", "x", "", "CWE XML目录文件路径")
	filterCmd.Flags().StringVarP(&filterKeyword, "keyword", "k", "", "按关键字搜索")
	filterCmd.Flags().StringVarP(&filterAbstr, "abstraction", "a", "", "按抽象层级搜索")
	filterCmd.Flags().StringVarP(&filterStatus, "status", "s", "", "按状态搜索")
	filterCmd.Flags().StringVarP(&filterLikelihood, "likelihood", "l", "", "按利用可能性搜索")
	filterCmd.Flags().StringVarP(&filterStructure, "structure", "t", "", "按结构类型搜索")
	filterCmd.Flags().StringVarP(&filterScope, "scope", "", "", "按后果范围搜索")
	filterCmd.Flags().StringVar(&filterSort, "sort", "", "排序字段 (id/name/abstraction)")
	filterCmd.Flags().StringVar(&filterGroupBy, "group-by", "", "分组字段 (abstraction/status/likelihood)")

	// stats flags
	statsCmd.Flags().StringVarP(&xmlFilePath, "xml", "x", "", "CWE XML目录文件路径")

	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(filterCmd)
	rootCmd.AddCommand(statsCmd)
}
