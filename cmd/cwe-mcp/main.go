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
	"log"
	"os"
	"sync"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// mcpVersion 是 MCP 服务器自身版本（独立于 SDK 版本）。
var mcpVersion = "0.1.0"

// 全局状态：离线注册表（懒加载）
var (
	xmlPath     string
	registry    *cweskills.Registry
	registryMu  sync.Mutex
	registryErr error
)

func main() {
	var (
		transport = flag.String("transport", "stdio", "传输方式: stdio 或 http")
		addr      = flag.String("addr", ":8080", "HTTP 模式监听地址")
		xml       = flag.String("xml", "", "CWE XML 目录文件路径（离线工具需要）")
		showVer   = flag.Bool("version", false, "显示版本信息并退出")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("cwe-mcp: %s\nsdk: %s\n", mcpVersion, cweskills.Version)
		return
	}

	xmlPath = *xml

	// 预加载 XML（若指定）—— 失败不退出，离线工具调用时报错
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
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("stdio 服务器错误: %v", err)
		}
	case "http":
		srv := server.NewSSEServer(s)
		log.Printf("MCP SSE 服务器监听 %s", *addr)
		if err := srv.Start(*addr); err != nil {
			log.Fatalf("HTTP 服务器错误: %v", err)
		}
	default:
		log.Fatalf("未知传输方式: %s", *transport)
	}
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
