package main

import (
	"context"
	"fmt"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerOfflineTools 注册离线 XML 工具。
// 这些工具需要启动 cwe-mcp 时通过 --xml 指定 CWE XML 目录文件。
func registerOfflineTools(s *server.MCPServer) {
	// get_ancestors —— 离线获取祖先链
	s.AddTool(
		mcp.NewTool("get_ancestors",
			mcp.WithDescription("Get all ancestor CWEs of a weakness from the offline XML registry (requires --xml). Returns the full ancestor chain. Unlike the online API, the offline registry contains all 10 relationship types."),
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
			ancestors := nav.Ancestors(id)
			return wrapJSON(map[string]any{
				"cwe_id":    cweskills.FormatCWEIDFromInt(id),
				"ancestors": ancestors,
			})
		},
	)

	// get_descendants —— 离线获取后代
	s.AddTool(
		mcp.NewTool("get_descendants",
			mcp.WithDescription("Get all descendant CWEs of a weakness from the offline XML registry (requires --xml)."),
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
			descendants := nav.Descendants(id)
			return wrapJSON(map[string]any{
				"cwe_id":     cweskills.FormatCWEIDFromInt(id),
				"descendants": descendants,
			})
		},
	)

	// get_shortest_path —— 离线最短路径
	s.AddTool(
		mcp.NewTool("get_shortest_path",
			mcp.WithDescription("Find the shortest relationship path between two CWEs in the offline XML registry (requires --xml). Returns the list of CWE IDs forming the path, or empty if no path exists."),
			mcp.WithString("from", mcp.Required(), mcp.Description("Source CWE ID, e.g. 'CWE-79'")),
			mcp.WithString("to", mcp.Required(), mcp.Description("Target CWE ID, e.g. 'CWE-1'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			fromStr, ok := requireStringArg(req, "from")
			if !ok {
				return errResult("missing or invalid 'from'"), nil
			}
			toStr, ok := requireStringArg(req, "to")
			if !ok {
				return errResult("missing or invalid 'to'"), nil
			}
			from, err := cweskills.ParseCWEID(fromStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid 'from': %v", err)), nil
			}
			to, err := cweskills.ParseCWEID(toStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid 'to': %v", err)), nil
			}
			nav := cweskills.NewNavigator(reg)
			path := nav.ShortestPath(from, to)
			// hops = 边数 = 节点数 - 1；path 为空或单点时无路径
			hops := 0
			if len(path) > 1 {
				hops = len(path) - 1
			}
			return wrapJSON(map[string]any{
				"from":     cweskills.FormatCWEIDFromInt(from),
				"to":       cweskills.FormatCWEIDFromInt(to),
				"path":     path,
				"hops":     hops,
				"found":    len(path) > 1,
			})
		},
	)

	// build_tree —— 离线构建层次树
	s.AddTool(
		mcp.NewTool("build_tree",
			mcp.WithDescription("Build a CWE hierarchy tree rooted at the given CWE from the offline XML registry (requires --xml). Returns the tree structure with nested children."),
			mcp.WithString("root", mcp.Required(), mcp.Description("Root CWE ID, e.g. 'CWE-1'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			rootStr, ok := requireStringArg(req, "root")
			if !ok {
				return errResult("missing or invalid 'root'"), nil
			}
			root, err := cweskills.ParseCWEID(rootStr)
			if err != nil {
				return errResult(fmt.Sprintf("invalid 'root': %v", err)), nil
			}
			tree := cweskills.BuildTree(reg, root)
			if tree == nil {
				return errResult(fmt.Sprintf("CWE %s not found in registry", rootStr)), nil
			}
			return wrapJSON(map[string]any{
				"root":      cweskills.FormatCWEIDFromInt(root),
				"tree":      serializeTreeNode(tree),
				"count":     tree.Count(),
				"max_depth": tree.MaxDepth(),
			})
		},
	)

	// search_keyword —— 离线关键词搜索
	s.AddTool(
		mcp.NewTool("search_keyword",
			mcp.WithDescription("Search the offline XML registry for CWEs whose name or description matches a keyword (requires --xml). Returns matching CWE list."),
			mcp.WithString("keyword", mcp.Required(), mcp.Description("Search keyword, e.g. 'Injection'")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			keyword, ok := requireStringArg(req, "keyword")
			if !ok {
				return errResult("missing or invalid 'keyword'"), nil
			}
			results := cweskills.FindByKeyword(reg, keyword)
			return wrapJSON(map[string]any{
				"keyword": keyword,
				"count":   len(results),
				"results": results,
			})
		},
	)

	// registry_stats —— 离线注册表统计
	s.AddTool(
		mcp.NewTool("registry_stats",
			mcp.WithDescription("Get statistics from the offline XML registry (requires --xml): total count, breakdown by abstraction/status."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reg, err := mustRegistry()
			if err != nil {
				return errResult(err.Error()), nil
			}
			stats := cweskills.ComputeStatistics(reg)
			return wrapJSON(map[string]any{
				"total": reg.Size(),
				"stats": stats,
			})
		},
	)
}

// serializeTreeNode 把 TreeNode 转成可 JSON 序列化的嵌套结构。
// TreeNode.Parent 是反向指针，直接序列化会触发循环引用导致栈溢出，
// 因此这里只保留 CWE 内容与 Children 子树（丢弃 Parent）。
func serializeTreeNode(node *cweskills.TreeNode) map[string]any {
	if node == nil {
		return nil
	}
	children := make([]map[string]any, 0, len(node.Children))
	for _, c := range node.Children {
		children = append(children, serializeTreeNode(c))
	}
	return map[string]any{
		"id":       cweskills.FormatCWEIDFromInt(node.CWE.ID),
		"name":     node.CWE.Name,
		"depth":    node.Depth,
		"children": children,
	}
}
