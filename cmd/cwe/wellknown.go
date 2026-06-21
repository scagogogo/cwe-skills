package main

import (
	"fmt"
	"sort"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// wellknownCmd 知名CWE列表相关命令
var wellknownCmd = &cobra.Command{
	Use:   "wellknown",
	Short: "知名CWE列表查询",
	Long: `查询CWE的各种知名列表，包括：
  - CWE Top 25 Most Dangerous Software Weaknesses
  - OWASP Top 10
  - SANS Top 25

示例：
  cwe wellknown top25          # 列出CWE Top 25
  cwe wellknown owasp          # 列出OWASP Top 10
  cwe wellknown sans           # 列出SANS Top 25
  cwe wellknown check CWE-79   # 检查CWE-79是否在知名列表中`,
}

// top25Cmd 列出CWE Top 25
var top25Cmd = &cobra.Command{
	Use:   "top25",
	Short: "列出CWE Top 25 Most Dangerous Software Weaknesses",
	Long:  "列出MITRE CWE Top 25 Most Dangerous Software Weaknesses.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := cweskills.CWETop25
		return printIDList(cmd, "CWE Top 25 Most Dangerous Software Weaknesses", ids)
	},
}

// owaspCmd 列出OWASP Top 10
var owaspCmd = &cobra.Command{
	Use:   "owasp",
	Short: "列出OWASP Top 10 (2021) 对应的CWE ID",
	Long:  "列出OWASP Top 10 (2021) 对应的CWE ID及其分类。",
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputFormat == "json" {
			type owaspEntry struct {
				Category string `json:"category"`
				CWEIDs   []int  `json:"cwe_ids"`
			}
			var entries []owaspEntry
			for _, cat := range sortedOWASPKeys() {
				entries = append(entries, owaspEntry{
					Category: cat,
					CWEIDs:   cweskills.OWASPTop10[cat],
				})
			}
			return printJSON(cmd, entries)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "OWASP Top 10 (2021):\n\n")
		for _, cat := range sortedOWASPKeys() {
			ids := cweskills.OWASPTop10[cat]
			fmt.Fprintf(cmd.OutOrStdout(), "  %s:\n", cat)
			for _, id := range ids {
				fmt.Fprintf(cmd.OutOrStdout(), "    %s\n", cweskills.FormatCWEIDFromInt(id))
			}
		}
		return nil
	},
}

// sansCmd 列出SANS Top 25
var sansCmd = &cobra.Command{
	Use:   "sans",
	Short: "列出SANS Top 25 Most Dangerous Software Errors",
	Long:  "列出SANS Top 25 Most Dangerous Software Errors对应的CWE ID.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := cweskills.SANSTop25
		return printIDList(cmd, "SANS Top 25 Most Dangerous Software Errors", ids)
	},
}

// checkListCmd 检查CWE ID是否在知名列表中
var checkListCmd = &cobra.Command{
	Use:   "check [CWE-ID...]",
	Short: "检查CWE ID是否在知名列表中",
	Long: `检查给定的CWE ID是否属于CWE Top 25、OWASP Top 10或SANS Top 25.

示例：
  cwe wellknown check CWE-79
  cwe wellknown check 79 89 352`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("请提供至少一个CWE ID")
		}

		type checkResult struct {
			CWEID  string   `json:"cwe_id"`
			InList []string `json:"in_list,omitempty"`
		}

		results := make([]checkResult, 0, len(args))
		for _, input := range args {
			id, err := cweskills.ParseCWEID(input)
			if err != nil {
				results = append(results, checkResult{CWEID: input, InList: []string{}})
				continue
			}

			var inList []string
			if cweskills.IsInTop25(id) {
				inList = append(inList, "Top 25")
			}
			if cweskills.IsInOWASPTop10(id) {
				cat := cweskills.GetOWASPCategory(id)
				if cat != "" {
					inList = append(inList, "OWASP Top 10 ("+cat+")")
				} else {
					inList = append(inList, "OWASP Top 10")
				}
			}
			if cweskills.IsInSANSTop25(id) {
				inList = append(inList, "SANS Top 25")
			}

			results = append(results, checkResult{
				CWEID:  cweskills.FormatCWEIDFromInt(id),
				InList: inList,
			})
		}

		if outputFormat == "json" {
			return printJSON(cmd, results)
		}

		for _, r := range results {
			if len(r.InList) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: 不在任何知名列表中\n", r.CWEID)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %v\n", r.CWEID, r.InList)
			}
		}
		return nil
	},
}

func printIDList(cmd *cobra.Command, title string, ids []int) error {
	if outputFormat == "json" {
		type idEntry struct {
			ID     int    `json:"id"`
			Format string `json:"format"`
		}
		entries := make([]idEntry, len(ids))
		for i, id := range ids {
			entries[i] = idEntry{ID: id, Format: cweskills.FormatCWEIDFromInt(id)}
		}
		return printJSON(cmd, entries)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s (%d 项):\n\n", title, len(ids))
	for i, id := range ids {
		fmt.Fprintf(cmd.OutOrStdout(), "  %2d. %s\n", i+1, cweskills.FormatCWEIDFromInt(id))
	}
	return nil
}

func sortedOWASPKeys() []string {
	keys := make([]string, 0, len(cweskills.OWASPTop10))
	for k := range cweskills.OWASPTop10 {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func init() {
	wellknownCmd.AddCommand(top25Cmd)
	wellknownCmd.AddCommand(owaspCmd)
	wellknownCmd.AddCommand(sansCmd)
	wellknownCmd.AddCommand(checkListCmd)
	rootCmd.AddCommand(wellknownCmd)
}
