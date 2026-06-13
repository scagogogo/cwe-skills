package main

import (
	"fmt"

	cwepkg "github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

var (
	xmlFilePath    string
	searchKeyword  string
	searchAbstr    string
	searchStatus   string
	searchLikelihood string
	searchStructure  string
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
  cwe search --xml cwec_latest.xml --status Stable`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if xmlFilePath == "" {
			return fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
		}

		parser := cwepkg.NewXMLParser()
		registry, err := parser.ParseFile(xmlFilePath)
		if err != nil {
			return fmt.Errorf("解析XML文件失败: %w", err)
		}

		var results []*cwepkg.CWE

		if searchKeyword != "" {
			results = cwepkg.FindByKeyword(registry, searchKeyword)
		} else if searchAbstr != "" {
			abstr, err := cwepkg.ParseAbstraction(searchAbstr)
			if err != nil {
				return fmt.Errorf("无效的抽象层级: %w", err)
			}
			results = cwepkg.FindByAbstraction(registry, abstr)
		} else if searchStatus != "" {
			st, err := cwepkg.ParseStatus(searchStatus)
			if err != nil {
				return fmt.Errorf("无效的状态: %w", err)
			}
			results = cwepkg.FindByStatus(registry, st)
		} else if searchLikelihood != "" {
			lh, err := cwepkg.ParseLikelihoodOfExploit(searchLikelihood)
			if err != nil {
				return fmt.Errorf("无效的利用可能性: %w", err)
			}
			results = cwepkg.FindByLikelihood(registry, lh)
		} else if searchStructure != "" {
			st, err := cwepkg.ParseStructure(searchStructure)
			if err != nil {
				return fmt.Errorf("无效的结构类型: %w", err)
			}
			results = cwepkg.FindByStructure(registry, st)
		} else {
			// 无过滤条件，返回所有条目
			results = registry.GetAll()
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "找到 %d 个CWE条目:\n\n", len(results))
		for _, cwe := range results {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s - %s [%s, %s]\n",
				cwepkg.FormatCWEIDFromInt(cwe.ID), cwe.Name, cwe.Abstraction, cwe.Status)
		}
		return nil
	},
}

// statsCmd 统计CWE数据
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "统计CWE数据",
	Long: `从本地XML数据源统计CWE数据分布。

需要先通过 --xml 参数指定MITRE CWE XML目录文件路径。

示例：
  cwe stats --xml cwec_latest.xml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if xmlFilePath == "" {
			return fmt.Errorf("请通过 --xml 参数指定CWE XML目录文件路径")
		}

		parser := cwepkg.NewXMLParser()
		registry, err := parser.ParseFile(xmlFilePath)
		if err != nil {
			return fmt.Errorf("解析XML文件失败: %w", err)
		}

		stats := cwepkg.ComputeStatistics(registry)

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

func init() {
	searchCmd.Flags().StringVarP(&xmlFilePath, "xml", "x", "", "CWE XML目录文件路径")
	searchCmd.Flags().StringVarP(&searchKeyword, "keyword", "k", "", "按关键字搜索")
	searchCmd.Flags().StringVarP(&searchAbstr, "abstraction", "a", "", "按抽象层级搜索 (Class/Base/Variant/Pillar)")
	searchCmd.Flags().StringVarP(&searchStatus, "status", "s", "", "按状态搜索 (Stable/Draft/Deprecated)")
	searchCmd.Flags().StringVarP(&searchLikelihood, "likelihood", "l", "", "按利用可能性搜索 (High/Medium/Low)")
	searchCmd.Flags().StringVarP(&searchStructure, "structure", "t", "", "按结构类型搜索 (Simple/Chain/Composite)")

	statsCmd.Flags().StringVarP(&xmlFilePath, "xml", "x", "", "CWE XML目录文件路径")

	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)
}