package main

import (
	"context"
	"fmt"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerExtraTools 注册补充工具，覆盖更多 CLI 能力。
func registerExtraTools(s *server.MCPServer) {
	// format_cwe_id —— 格式化整数 ID 为 CWE-NNN
	s.AddTool(
		mcp.NewTool("format_cwe_id",
			mcp.WithDescription("Format one or more integer CWE IDs into the canonical 'CWE-NNN' form. Pass a list of integers like [79, 89, 352]."),
			mcp.WithArray("ids", mcp.Required(), mcp.Description("List of integer CWE IDs"), mcp.Items(map[string]any{"type": "integer"})),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idsRaw, ok := req.Params.Arguments["ids"].([]any)
			if !ok {
				return errResult("missing or invalid 'ids' (must be array of integers)"), nil
			}
			var formatted []string
			for _, v := range idsRaw {
				id, ok := toInt(v)
				if !ok {
					return errResult(fmt.Sprintf("invalid id value: %v", v)), nil
				}
				formatted = append(formatted, cweskills.FormatCWEIDFromInt(id))
			}
			return wrapJSON(map[string]any{"ids": formatted})
		},
	)

	// get_siblings —— 离线获取同级
	s.AddTool(
		mcp.NewTool("get_siblings",
			mcp.WithDescription("Get sibling CWEs (sharing the same parent) of a weakness from the offline XML registry (requires --xml)."),
			mcp.WithString("id", mcp.Required(), mcp.Description("CWE ID, e.g. 'CWE-79'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			idStr, ok := requireStringArg(req, "id")
			if !ok {
				return errResult("missing or invalid 'id'"), nil
			}
			id, err := cweskills.ParseCWEID(idStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid CWE ID: %v", err)), nil
			}
			nav := cweskills.NewNavigator(reg)
			siblings := nav.Siblings(id)
			return wrapJSON(map[string]any{
				"cwe_id":   cweskills.FormatCWEIDFromInt(id),
				"siblings": siblings,
			})
		},
	)

	// get_children —— 离线获取直接子级
	s.AddTool(
		mcp.NewTool("get_children",
			mcp.WithDescription("Get direct child CWEs of a weakness from the offline XML registry (requires --xml)."),
			mcp.WithString("id", mcp.Required(), mcp.Description("CWE ID, e.g. 'CWE-74'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			idStr, ok := requireStringArg(req, "id")
			if !ok {
				return errResult("missing or invalid 'id'"), nil
			}
			id, err := cweskills.ParseCWEID(idStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid CWE ID: %v", err)), nil
			}
			nav := cweskills.NewNavigator(reg)
			children := nav.Children(id)
			return wrapJSON(map[string]any{
				"cwe_id":   cweskills.FormatCWEIDFromInt(id),
				"children": children,
			})
		},
	)

	// filter_cwes —— 离线多条件过滤
	s.AddTool(
		mcp.NewTool("filter_cwes",
			mcp.WithDescription("Filter CWEs in the offline XML registry by abstraction and/or status (requires --xml). Returns matching CWE list. Use search_keyword for keyword matching; this tool is for structured attribute filtering. For a specific CWE by ID, use get_weakness."),
			mcp.WithString("abstraction", mcp.Description("Abstraction level: Pillar, Class, Base, Variant (optional)")),
			mcp.WithString("status", mcp.Description("Status: Stable, Usable, Draft, Incomplete, Obsolete, Deprecated (optional)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			opts := []cweskills.FilterOption{}
			if abs, ok := req.Params.Arguments["abstraction"].(string); ok && abs != "" {
				opts = append(opts, cweskills.FilterOption{Abstraction: cweskills.Abstraction(abs)})
			}
			if st, ok := req.Params.Arguments["status"].(string); ok && st != "" {
				opts = append(opts, cweskills.FilterOption{Status: cweskills.Status(st)})
			}
			all := reg.GetAll()
			filtered := cweskills.Filter(all, opts...)
			return wrapJSON(map[string]any{
				"count":   len(filtered),
				"results": filtered,
			})
		},
	)

	// is_ancestor —— 离线祖先判定
	s.AddTool(
		mcp.NewTool("is_ancestor",
			mcp.WithDescription("Check whether a CWE is an ancestor of another in the offline XML registry (requires --xml). Returns boolean."),
			mcp.WithString("ancestor", mcp.Required(), mcp.Description("Candidate ancestor CWE ID, e.g. 'CWE-1'")),
			mcp.WithString("descendant", mcp.Required(), mcp.Description("Candidate descendant CWE ID, e.g. 'CWE-79'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			ancStr, ok := requireStringArg(req, "ancestor")
			if !ok {
				return errResult("missing or invalid 'ancestor'"), nil
			}
			descStr, ok := requireStringArg(req, "descendant")
			if !ok {
				return errResult("missing or invalid 'descendant'"), nil
			}
			anc, err := cweskills.ParseCWEID(ancStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid 'ancestor': %v", err)), nil
			}
			desc, err := cweskills.ParseCWEID(descStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid 'descendant': %v", err)), nil
			}
			nav := cweskills.NewNavigator(reg)
			result := nav.IsAncestorOf(anc, desc)
			return wrapJSON(map[string]any{
				"ancestor":   cweskills.FormatCWEIDFromInt(anc),
				"descendant": cweskills.FormatCWEIDFromInt(desc),
				"is_ancestor": result,
			})
		},
	)
}

// toInt 把 any 安全转为 int（JSON 数字可能是 float64）。
func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	}
	return 0, false
}
