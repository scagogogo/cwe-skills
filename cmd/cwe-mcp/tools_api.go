package main

import (
	"context"
	"fmt"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// newAPIClient 构造在线 API 客户端，默认用 MITRE base URL。
// 提取为包级变量以便测试注入指向 httptest server 的客户端。
var newAPIClient = cweskills.NewAPIClient

// registerAPITools 注册在线 MITRE API 工具。
// 这些工具调用 MITRE REST API，受速率限制。
func registerAPITools(s *server.MCPServer) {
	// get_weakness —— 在线获取弱点详情，失败时回退到离线注册表
	s.AddTool(
		mcp.NewTool("get_weakness",
			mcp.WithDescription("Fetch a CWE weakness's full details. Tries the MITRE REST API first (online, rate-limited ~0.1 req/s); if the API is unavailable, falls back to the offline XML registry (requires --xml). Returns name, description, abstraction, status, relationships, consequences, etc."),
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
			client := newAPIClient()
			defer client.Close()
			weakness, err := client.GetWeakness(ctx, id)
			if err == nil {
				return wrapJSON(map[string]any{
					"source": "online",
					"cwe":    weakness,
				})
			}
			// 在线失败，尝试离线回退
			reg, regErr := mustRegistry()
			if regErr != nil {
				return errResult(fmt.Sprintf("API call failed and no offline registry: api=%v; xml=%v", err, regErr)), nil
			}
			offline, found := reg.Get(id)
			if !found {
				return errResult(fmt.Sprintf("API call failed and CWE %d not in offline registry: api=%v", id, err)), nil
			}
			return wrapJSON(map[string]any{
				"source":     "offline",
				"cwe":        offline,
				"api_error":  err.Error(),
				"note":       "MITRE API unavailable; returned from local XML registry",
			})
		},
	)

	// get_parents —— 在线获取父级关系，失败时回退到离线注册表
	s.AddTool(
		mcp.NewTool("get_parents",
			mcp.WithDescription("Fetch the direct parent CWEs of a weakness. Tries the MITRE REST API first (online); if unavailable, falls back to the offline XML registry (requires --xml), which also covers all 10 relationship types. For the full ancestor chain, use get_ancestors."),
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
			client := newAPIClient()
			defer client.Close()
			parents, err := client.GetParents(ctx, id)
			if err == nil {
				return wrapJSON(map[string]any{
					"source":  "online",
					"cwe_id":  cweskills.FormatCWEIDFromInt(id),
					"parents": parents,
				})
			}
			// 在线失败，回退到离线注册表
			reg, regErr := mustRegistry()
			if regErr != nil {
				return errResult(fmt.Sprintf("API call failed and no offline registry: api=%v; xml=%v", err, regErr)), nil
			}
			nav := cweskills.NewNavigator(reg)
			offlineParents := nav.Parents(id)
			return wrapJSON(map[string]any{
				"source":    "offline",
				"cwe_id":    cweskills.FormatCWEIDFromInt(id),
				"parents":   offlineParents,
				"api_error": err.Error(),
				"note":      "MITRE API unavailable; returned from local XML registry",
			})
		},
	)

	// api_version —— 检查 MITRE API 版本
	s.AddTool(
		mcp.NewTool("api_version",
			mcp.WithDescription("Check the current MITRE CWE REST API version (online). Useful to verify API availability."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client := newAPIClient()
			defer client.Close()
			ver, err := client.GetVersion(ctx)
			if err != nil {
				return errResult(fmt.Sprintf("API call failed: %v", err)), nil
			}
			return wrapJSON(map[string]any{"version": ver})
		},
	)
}
