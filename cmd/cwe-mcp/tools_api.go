package main

import (
	"context"
	"fmt"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerAPITools 注册在线 MITRE API 工具。
// 这些工具调用 MITRE REST API，受速率限制。
func registerAPITools(s *server.MCPServer) {
	// get_weakness —— 在线获取弱点详情
	s.AddTool(
		mcp.NewTool("get_weakness",
			mcp.WithDescription("Fetch a CWE weakness's full details from the MITRE REST API (online). Includes name, description, abstraction, status, relationships, consequences, etc. Rate-limited (~0.1 req/s)."),
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
			client := cweskills.NewAPIClient()
			defer client.Close()
			weakness, err := client.GetWeakness(ctx, id)
			if err != nil {
				return errResult(fmt.Sprintf("API call failed: %v", err)), nil
			}
			return wrapJSON(weakness)
		},
	)

	// get_parents —— 在线获取父级关系
	s.AddTool(
		mcp.NewTool("get_parents",
			mcp.WithDescription("Fetch the parent CWEs of a weakness via the MITRE REST API (online). Note: the API only exposes parent/child relationships — for full relationship types (siblings, peers, chains, dependencies), use the offline get_ancestors tool with --xml."),
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
			client := cweskills.NewAPIClient()
			defer client.Close()
			parents, err := client.GetParents(ctx, id)
			if err != nil {
				return errResult(fmt.Sprintf("API call failed: %v", err)), nil
			}
			return wrapJSON(map[string]any{
				"cwe_id":  cweskills.FormatCWEIDFromInt(id),
				"parents": parents,
			})
		},
	)

	// api_version —— 检查 MITRE API 版本
	s.AddTool(
		mcp.NewTool("api_version",
			mcp.WithDescription("Check the current MITRE CWE REST API version (online). Useful to verify API availability."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client := cweskills.NewAPIClient()
			defer client.Close()
			ver, err := client.GetVersion(ctx)
			if err != nil {
				return errResult(fmt.Sprintf("API call failed: %v", err)), nil
			}
			return wrapJSON(map[string]any{"version": ver})
		},
	)
}
