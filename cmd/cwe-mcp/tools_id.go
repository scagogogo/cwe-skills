package main

import (
	"context"
	"fmt"

	"github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerIDTools 注册 CWE ID 工具类（纯本地，无需 XML 或网络）。
func registerIDTools(s *server.MCPServer) {
	// parse_cwe_id —— 解析 CWE ID
	s.AddTool(
		mcp.NewTool("parse_cwe_id",
			mcp.WithDescription("Parse a CWE ID string (e.g. 'CWE-79', 'cwe-79', '79') into a normalized form and integer ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("CWE ID string to parse")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := req.Params.Arguments["id"].(string)
			if !ok {
				return errResult("missing or invalid 'id' parameter"), nil
			}
			parsed, err := cweskills.ParseCWEID(id)
			if err != nil {
				return errResult(fmt.Sprintf("parse failed: %v", err)), nil
			}
			return wrapJSON(map[string]any{
				"input":  id,
				"id":     parsed,
				"format": cweskills.FormatCWEIDFromInt(parsed),
				"valid":  true,
			})
		},
	)

	// validate_cwe_id —— 验证 CWE ID 格式
	s.AddTool(
		mcp.NewTool("validate_cwe_id",
			mcp.WithDescription("Validate whether a string is a well-formed CWE ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("CWE ID string to validate")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := req.Params.Arguments["id"].(string)
			if !ok {
				return errResult("missing 'id'"), nil
			}
			err := cweskills.ValidateCWEID(id)
			return wrapJSON(map[string]any{
				"input": id,
				"valid": err == nil,
				"error": errMsg(err),
			})
		},
	)

	// extract_cwe_ids —— 从文本提取 CWE ID
	s.AddTool(
		mcp.NewTool("extract_cwe_ids",
			mcp.WithDescription("Extract all CWE IDs mentioned in a piece of text (e.g. a vulnerability report). Returns the list of CWE ID strings found."),
			mcp.WithString("text", mcp.Required(), mcp.Description("Text to scan for CWE IDs")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			text, ok := req.Params.Arguments["text"].(string)
			if !ok {
				return errResult("missing 'text'"), nil
			}
			ids := cweskills.ExtractCWEIDs(text)
			return wrapJSON(map[string]any{
				"count": len(ids),
				"ids":   ids,
			})
		},
	)

	// compare_cwe_ids —— 比较两个 CWE ID
	s.AddTool(
		mcp.NewTool("compare_cwe_ids",
			mcp.WithDescription("Compare two CWE IDs numerically. Returns -1, 0, or 1 (a < b, a == b, a > b)."),
			mcp.WithString("a", mcp.Required(), mcp.Description("First CWE ID")),
			mcp.WithString("b", mcp.Required(), mcp.Description("Second CWE ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a, ok := requireStringArg(req, "a")
			if !ok {
				return errResult("missing or invalid 'a'"), nil
			}
			b, ok := requireStringArg(req, "b")
			if !ok {
				return errResult("missing or invalid 'b'"), nil
			}
			cmp, err := cweskills.CompareCWEIDs(a, b)
			if err != nil {
				return errResult(fmt.Sprintf("compare failed: %v", err)), nil
			}
			return wrapJSON(map[string]any{"a": a, "b": b, "comparison": cmp})
		},
	)
}

// wrapJSON 返回一个以 JSON 文本为内容的工具结果。
func wrapJSON(v any) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(string(jsonRaw(v))), nil
}

// errMsg 把 error 转为字符串（nil 时返回空串）。
func errMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// requireStringArg 从工具参数中提取必填字符串字段。
// 返回 (value, ok)——ok 为 false 时 caller 应返回 errResult 提示参数缺失/类型错误。
// 相比裸 .(string)，它能区分"缺失"与"类型错误"，给 AI 更精准的报错。
func requireStringArg(req mcp.CallToolRequest, key string) (string, bool) {
	v, ok := req.Params.Arguments[key].(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}
