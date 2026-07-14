package main

import (
	"context"
	"fmt"

	"github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerWellknownTools 注册知名列表工具（纯本地，基于内置常量）。
func registerWellknownTools(s *server.MCPServer) {
	// check_wellknown —— 检查 CWE 是否在 Top 25 / OWASP / SANS
	s.AddTool(
		mcp.NewTool("check_wellknown",
			mcp.WithDescription("Check whether a CWE ID is in the well-known lists: CWE Top 25, OWASP Top 10, and SANS Top 25. Returns the lists the CWE belongs to."),
			mcp.WithString("id", mcp.Required(), mcp.Description("CWE ID, e.g. 'CWE-79'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idStr, ok := req.Params.Arguments["id"].(string)
			if !ok {
				return errResult("missing 'id'"), nil
			}
			id, err := cweskills.ParseCWEID(idStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid CWE ID: %v", err)), nil
			}
			var inList []string
			if cweskills.IsInTop25(id) {
				inList = append(inList, "Top 25")
			}
			if cweskills.IsInOWASPTop10(id) {
				inList = append(inList, "OWASP Top 10")
			}
			if cweskills.IsInSANSTop25(id) {
				inList = append(inList, "SANS Top 25")
			}
			return wrapJSON(map[string]any{
				"cwe_id":  cweskills.FormatCWEIDFromInt(id),
				"id":      id,
				"in_list": inList,
				"in_top25": cweskills.IsInTop25(id),
			})
		},
	)

	// get_owasp_categories —— 查询 CWE 对应的 OWASP 类别
	s.AddTool(
		mcp.NewTool("get_owasp_categories",
			mcp.WithDescription("Get the OWASP Top 10 (2021) categories a CWE belongs to. Returns the list of OWASP category codes like 'A03:2021-Injection'."),
			mcp.WithString("id", mcp.Required(), mcp.Description("CWE ID, e.g. 'CWE-79'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idStr, ok := req.Params.Arguments["id"].(string)
			if !ok {
				return errResult("missing 'id'"), nil
			}
			id, err := cweskills.ParseCWEID(idStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid CWE ID: %v", err)), nil
			}
			cats := cweskills.GetOWASPCategories(id)
			return wrapJSON(map[string]any{
				"cwe_id":    cweskills.FormatCWEIDFromInt(id),
				"categories": cats,
			})
		},
	)
}
