package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

// TestRunMain_Version 覆盖 --version 分支（showVer 为真 → return 0）
func TestRunMain_Version(t *testing.T) {
	code := runMain([]string{"-version"})
	if code != 0 {
		t.Errorf("runMain -version: want exit 0, got %d", code)
	}
}

// TestRunMain_UnknownTransport 覆盖 default 分支（未知 transport → return 2 + stderr）
func TestRunMain_UnknownTransport(t *testing.T) {
	var buf bytes.Buffer
	orig := osStderr
	osStderr = &buf
	t.Cleanup(func() { osStderr = orig })

	code := runMain([]string{"-transport", "weird"})
	if code != 2 {
		t.Errorf("runMain unknown transport: want exit 2, got %d", code)
	}
	if !bytes.Contains(buf.Bytes(), []byte("未知传输方式")) {
		t.Errorf("expected stderr to mention unknown transport, got: %s", buf.String())
	}
}

// TestRunMain_FlagParseError 覆盖 fs.Parse 失败分支（return 2）
func TestRunMain_FlagParseError(t *testing.T) {
	// -version=notabool 让 flag 包报 "invalid boolean value"
	code := runMain([]string{"-version=notabool"})
	if code != 2 {
		t.Errorf("runMain bad bool flag: want exit 2, got %d", code)
	}
}

// TestRunMain_StdioSuccess 覆盖 transport=stdio 成功分支：
// 注入 serveStdio 立即返回 nil，runMain 走完 stdio case 返回 0。
func TestRunMain_StdioSuccess(t *testing.T) {
	orig := serveStdio
	serveStdio = func(s *server.MCPServer, opts ...server.StdioOption) error {
		return nil
	}
	t.Cleanup(func() { serveStdio = orig })

	code := runMain([]string{"-transport", "stdio"})
	if code != 0 {
		t.Errorf("runMain stdio success: want exit 0, got %d", code)
	}
}

// TestRunMain_HTTPSuccess 覆盖 transport=http 成功分支：
// 注入 newSSEServer 返回 stub、sseStart 立即返回 nil，不阻塞。
func TestRunMain_HTTPSuccess(t *testing.T) {
	origNew := newSSEServer
	origStart := sseStart
	newSSEServer = func(s *server.MCPServer, opts ...server.SSEOption) *server.SSEServer {
		return nil // sseStart 已注入，runMain 不再访问返回的 srv
	}
	sseStart = func(addr string) error { return nil }
	t.Cleanup(func() {
		newSSEServer = origNew
		sseStart = origStart
	})

	code := runMain([]string{"-transport", "http", "-addr", ":0"})
	if code != 0 {
		t.Errorf("runMain http success: want exit 0, got %d", code)
	}
}

// TestRunMain_XMLLoadWarning 覆盖 if *xml != "" 块（main.go）：
// 传 -xml 指向不存在路径，loadRegistry 失败走 log.Printf 警告分支，
// runMain 不中断继续注册离线工具并走 stdio 成功路径返回 0。
func TestRunMain_XMLLoadWarning(t *testing.T) {
	origStdio := serveStdio
	serveStdio = func(s *server.MCPServer, opts ...server.StdioOption) error { return nil }
	t.Cleanup(func() { serveStdio = origStdio })

	code := runMain([]string{"-transport", "stdio", "-xml", "/nonexistent/cwe.xml"})
	if code != 0 {
		t.Errorf("runMain with bad xml: want exit 0 (warning only), got %d", code)
	}
}

// TestRunMain_HTTPStartFallbackAndFatalf 覆盖 sseStart=nil 回退赋值与 HTTP Fatalf 调用点：
// sseStart 保持 nil（默认）走 start=srv.Start 回退；newSSEServer 注入返回真实 SSEServer；
// -addr 绑不上的地址使 srv.Start 返回 error，fatalf 注入不退出覆盖调用点。
func TestRunMain_HTTPStartFallbackAndFatalf(t *testing.T) {
	origNew := newSSEServer
	origStart := sseStart
	origFatalf := fatalf
	fatalfCalled := false
	sseStart = nil // 走 start = srv.Start 回退
	newSSEServer = func(s *server.MCPServer, opts ...server.SSEOption) *server.SSEServer {
		return server.NewSSEServer(s) // 真实 SSEServer，Start 会尝试监听
	}
	fatalf = func(format string, a ...any) { fatalfCalled = true }
	t.Cleanup(func() {
		newSSEServer = origNew
		sseStart = origStart
		fatalf = origFatalf
	})

	// bad:addr 形式非法，net.Listen 立即返回 error，srv.Start 返回 error
	code := runMain([]string{"-transport", "http", "-addr", "bad:addr:99999"})
	if !fatalfCalled {
		t.Fatal("runMain http bad addr: expected fatalf to be called")
	}
	// fatalf 注入后不 os.Exit，runMain 继续落到 return 0
	if code != 0 {
		t.Errorf("runMain http bad addr with fatalf injected: want exit 0, got %d", code)
	}
}

// TestRunMain_StdioFatalf 覆盖 stdio Fatalf 调用点：
// serveStdio 注入返回 error，fatalf 注入不退出，覆盖调用点。
func TestRunMain_StdioFatalf(t *testing.T) {
	origStdio := serveStdio
	origFatalf := fatalf
	fatalfCalled := false
	serveStdio = func(s *server.MCPServer, opts ...server.StdioOption) error {
		return fmt.Errorf("injected stdio error")
	}
	fatalf = func(format string, a ...any) { fatalfCalled = true }
	t.Cleanup(func() { serveStdio = origStdio; fatalf = origFatalf })

	code := runMain([]string{"-transport", "stdio"})
	if !fatalfCalled {
		t.Fatal("runMain stdio error: expected fatalf to be called")
	}
	// fatalf 注入后不 os.Exit，runMain 继续落到 return 0
	if code != 0 {
		t.Errorf("runMain stdio error with fatalf injected: want exit 0, got %d", code)
	}
}

// TestMain_InjectsOSExit 覆盖 main 的 os.Exit 调用点（main.go）：
// 注入 osExit 不退出，断言 main 调用了 osExit。
func TestMain_InjectsOSExit(t *testing.T) {
	origExit := osExit
	origStdio := serveStdio
	called := false
	osExit = func(code int) { called = true }
	serveStdio = func(s *server.MCPServer, opts ...server.StdioOption) error { return nil }
	t.Cleanup(func() { osExit = origExit; serveStdio = origStdio })

	main()
	if !called {
		t.Fatal("main: expected osExit to be called")
	}
}
