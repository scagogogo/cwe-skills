// Package main 是 CWE Skills 的 MCP (Model Context Protocol) 服务器。
//
// 它把 cweskills SDK 的核心能力以 MCP 工具的形式暴露给 AI 应用，
// 让 Claude Desktop、Cursor 等 MCP 兼容客户端能以标准化工具调用方式
// 访问 CWE 数据——无需 AI 自己执行 Shell 命令、无需解析 CLI 文本输出。
//
// 启动：
//
//	cwe-mcp                  # stdio 模式（默认，供 Claude Desktop 等本地客户端使用）
//	cwe-mcp --transport http --addr :8080  # HTTP 模式（远程部署）
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
	)
	flag.Parse()

	xmlPath = *xml

	// 预加载 XML（若指定）—— 失败不退出，离线工具调用时报错
	if *xml != "" {
		if err := loadRegistry(); err != nil {
			log.Printf("警告: 加载 XML 失败，离线工具将不可用: %v", err)
		}
	}

	s := server.NewMCPServer(
		"cwe-skills-mcp",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	registerIDTools(s)
	registerWellknownTools(s)
	registerAPITools(s)
	registerOfflineTools(s)

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

// loadRegistry 懒加载并构建 XML 注册表索引（线程安全）。
func loadRegistry() error {
	registryMu.Lock()
	defer registryMu.Unlock()

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

// mustRegistry 取已加载的注册表，未加载则返回错误给 MCP 客户端。
func mustRegistry() (*cweskills.Registry, error) {
	if registry != nil {
		return registry, nil
	}
	if err := loadRegistry(); err != nil {
		return nil, err
	}
	return registry, nil
}

// jsonRaw 把任意值编码为 json.RawMessage，失败时返回错误 JSON。
func jsonRaw(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{"error":"json marshal failed"}`)
	}
	return b
}

// errResult 构造一个带错误文本的工具结果。
func errResult(msg string) *mcp.CallToolResult {
	return mcp.NewToolResultError(msg)
}
