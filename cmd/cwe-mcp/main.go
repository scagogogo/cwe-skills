// Package main 是 CWE Skills 的 MCP (Model Context Protocol) 服务器。
//
// 它把 cweskills SDK 的核心能力以 MCP 工具的形式暴露给 AI 应用，
// 让 Claude Desktop、Cursor 等 MCP 兼容客户端能以标准化工具调用方式
// 访问 CWE 数据——无需 AI 自己执行 Shell 命令、无需解析 CLI 文本输出。
//
// 启动：
//
//	cwe-mcp                  # stdio 模式（默认，供 Claude Desktop 等本地客户端使用）
//	cwe-mcp --transport http --addr :8080  # SSE 模式（远程部署）
//
// 离线工具（get_ancestors / build_tree 等）需要通过 --xml 指定 XML 目录：
//
//	cwe-mcp --xml cwec_v4.15.xml
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// mcpVersion 是 MCP 服务器自身版本（独立于 SDK 版本）。
var mcpVersion = "0.1.0"

// osStderr 默认指向 os.Stderr，测试可替换为任意 io.Writer 以捕获输出。
var osStderr io.Writer = os.Stderr

// serveStdio 默认指向 server.ServeStdio，测试可替换为立即返回的 fake
// 以覆盖 stdio 分支而不阻塞。
var serveStdio = server.ServeStdio

// newSSEServer 默认指向 server.NewSSEServer，测试可替换。
var newSSEServer = server.NewSSEServer

// sseStart 可选地覆盖 *server.SSEServer.Start；nil 时回退到真实 srv.Start。
// 测试替换为立即返回 nil/err 的 fake 以覆盖 http 分支而不阻塞。
var sseStart func(addr string) error

// 全局状态：离线注册表（懒加载）
var (
	xmlPath     string
	registry    *cweskills.Registry
	registryMu  sync.Mutex
	registryErr error
)

func main() {
	os.Exit(runMain(os.Args[1:]))
}

// runMain 解析参数并按 transport 启动服务器，返回进程退出码。
// 提取自 main 以便测试覆盖 flag 解析与 --version / 未知 transport 分支，
// 而不触发 os.Exit（main 调用方负责 os.Exit）。
func runMain(args []string) int {
	fs := flag.NewFlagSet("cwe-mcp", flag.ContinueOnError)
	transport := fs.String("transport", "stdio", "传输方式: stdio 或 http")
	addr := fs.String("addr", ":8080", "HTTP 模式监听地址")
	xml := fs.String("xml", "", "CWE XML 目录文件路径（离线工具需要）")
	showVer := fs.Bool("version", false, "显示版本信息并退出")
	fs.Usage = func() {
		fmt.Fprintf(osStderr, "cwe-mcp — CWE Skills MCP 服务器（暴露 20 个 CWE 工具，供 MCP 兼容 AI 客户端调用）\n\n")
		fmt.Fprintf(osStderr, "用法:\n")
		fmt.Fprintf(osStderr, "  cwe-mcp                              stdio 模式，仅在线工具（无需 XML）\n")
		fmt.Fprintf(osStderr, "  cwe-mcp --xml cwec_v4.15.xml         stdio 模式，含离线工具\n")
		fmt.Fprintf(osStderr, "  cwe-mcp --transport http --addr :8080  SSE 模式（远程）\n")
		fmt.Fprintf(osStderr, "  cwe-mcp --version                    显示版本并退出\n\n")
		fmt.Fprintf(osStderr, "参数:\n")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *showVer {
		fmt.Printf("cwe-mcp: %s\nsdk: %s\n", mcpVersion, cweskills.Version)
		return 0
	}

	xmlPath = *xml

	if *xml != "" {
		registryMu.Lock()
		if err := loadRegistry(); err != nil {
			log.Printf("警告: 加载 XML 失败，离线工具将不可用: %v", err)
		}
		registryMu.Unlock()
	}

	s := server.NewMCPServer(
		"cwe-skills-mcp",
		mcpVersion,
		server.WithToolCapabilities(true),
	)

	registerIDTools(s)
	registerWellknownTools(s)
	registerAPITools(s)
	registerOfflineTools(s)
	registerExtraTools(s)

	switch *transport {
	case "stdio":
		if err := serveStdio(s); err != nil {
			log.Fatalf("stdio 服务器错误: %v", err)
		}
	case "http":
		srv := newSSEServer(s)
		log.Printf("MCP SSE 服务器监听 %s", *addr)
		start := sseStart
		if start == nil {
			start = srv.Start
		}
		if err := start(*addr); err != nil {
			log.Fatalf("HTTP 服务器错误: %v", err)
		}
	default:
		fmt.Fprintf(osStderr, "未知传输方式: %s\n", *transport)
		return 2
	}
	return 0
}

// loadRegistry 懒加载并构建 XML 注册表索引。
// 调用方必须持有 registryMu（内部不再加锁，避免重入死锁）。
func loadRegistry() error {
	if registry != nil || registryErr != nil {
		return registryErr
	}

	if xmlPath == "" {
		registryErr = fmt.Errorf("未指定 --xml，离线工具不可用；启动 cwe-mcp 时请加 --xml <path>")
		return registryErr
	}
	if _, err := os.Stat(xmlPath); err != nil {
		registryErr = fmt.Errorf("XML 文件不存在: %s: %w", xmlPath, err)
		return registryErr
	}

	reg, err := cweskills.NewXMLParser().ParseFile(xmlPath)
	if err != nil {
		registryErr = fmt.Errorf("解析 XML 失败: %w", err)
		return registryErr
	}
	reg.BuildIndexes()
	registry = reg
	log.Printf("已加载 XML 注册表: %d 个弱点", reg.Size())
	return nil
}

// mustRegistry 取已加载的注册表，未加载则触发懒加载。
// 加锁保证并发安全：避免多个工具调用同时看到 registry==nil 而重复加载。
func mustRegistry() (*cweskills.Registry, error) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if err := loadRegistry(); err != nil {
		return nil, err
	}
	return registry, nil
}

// jsonRaw 把任意值编码为 json.RawMessage，失败时返回带类型信息的错误 JSON。
func jsonRaw(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(fmt.Sprintf(`{"error":"json marshal failed","type":"%T","reason":"%s"}`, v, err.Error()))
	}
	return b
}

// errResult 构造一个带错误文本的工具结果。
func errResult(msg string) *mcp.CallToolResult {
	return mcp.NewToolResultError(msg)
}
