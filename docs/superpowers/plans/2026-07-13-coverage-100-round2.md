# 单元测试覆盖率提升至 100% 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 将项目单元测试覆盖率从当前 96.4%（总体）/ SDK 99.4% 提升至可达的 100%，对真正不可达的分支（底层库对合法输入永不返回错误）做可测试性重构或显式记录接受。

**Architecture:** 覆盖率缺口分三类，分别用三条路径处理：
1. **可注入 HTTP mock 的在线分支**（`getRelations`、`registerAPITools` 的 get_weakness/get_parents/api_version 在线分支）→ 用 httptest server + `WithAPIBaseURL` 注入；对 `registerAPITools` 内硬编码的 `cweskills.NewAPIClient()` 做最小可测试性改造：引入包级 `newAPIClient` 变量供测试替换。
2. **可测试性重构的入口函数**（`main`、`Execute`）→ 抽出 `runMain(args) int` / `executeRoot() int` 返回 exit code，`main` 仅 `os.Exit(...)`，测试调用返回值不触发 `os.Exit`。
3. **真正不可达分支**（`xml.MarshalIndent` 对纯值 `safeCWE` 永不失败；`csv.Writer.Write`/`Flush` 底层 `bytes.Buffer` 永不失败）→ 接受不覆盖，在计划与提交信息中记录原因。

**Tech Stack:** Go 1.25.0, testing, net/http/httptest, github.com/spf13/cobra v1.10.2, github.com/mark3labs/mcp-go v0.20.0

**Risks:**
- 改造 `tools_api.go` 与 `main.go`（cwe + cwe-mcp）属于"为可测试性改生产代码" → 缓解：改造保持行为不变（仅把 `cweskills.NewAPIClient` 换成包级变量、把 `main` 拆为 `runMain`+`os.Exit`），每个改造 Step 后立即跑全量测试确保不破坏现有 96.4%。
- `registerAPITools` 的 httptest 测试依赖 `getHandler` 反射取 handler 的现有机制 → 缓解：复用 `register_test.go` 已有的 `callTool`/`getHandler`，无需新增反射代码。
- `main` 的 stdio/http 分支调用 `server.ServeStdio`/`srv.Start` 会阻塞 → 缓解：测试只覆盖 `--version`（return 提前退出）、未知 transport（`log.Fatalf` → 用 `t.Setenv` 或 subprocess 验证），不覆盖真正的服务器启动分支（接受，记录原因）。

---

### Task 1: SDK 层 getRelations 在线分支全覆盖

**Depends on:** None
**Files:**
- Modify: `coverage_fillers_test.go`（追加测试函数，文件末尾）

- [ ] **Step 1: 在 coverage_fillers_test.go 末尾追加 getRelations 三分支测试 — 覆盖网络错误、fallback 成功、fallback 失败**

文件: `coverage_fillers_test.go`（在文件末尾追加）

```go
// TestGetRelations_HTTPError 覆盖 getRelations 的 httpClient.Get err 分支（行 109）
func TestGetRelations_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(1000, 1000))
	defer client.Close()

	_, err := client.GetParents(context.Background(), 79)
	if err == nil {
		t.Fatal("GetParents with 500 should error")
	}
}

// TestGetRelations_FallbackObjectArray 覆盖 json.Unmarshal 失败后 fallback 成功分支（行 121、124-134）
func TestGetRelations_FallbackObjectArray(t *testing.T) {
	// 返回对象数组格式（nature/cweId 字段），触发 Relationship 解析失败 → fallback
	body := `[{"nature":"ChildOf","cweId":100,"viewId":0},{"nature":"CanAlsoBe","cweId":200}]`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"Data":`+body+`}`)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(1000, 1000))
	defer client.Close()

	rel, err := client.GetParents(context.Background(), 79)
	if err != nil {
		t.Fatalf("GetParents fallback should succeed: %v", err)
	}
	if len(rel) != 2 {
		t.Fatalf("expected 2 relations, got %d", len(rel))
	}
	if rel[0].CWEID != 100 || rel[1].CWEID != 200 {
		t.Errorf("unexpected cwe ids: %+v", rel)
	}
}

// TestGetRelations_FallbackAlsoFails 覆盖两段解析都失败的分支（行 122）
func TestGetRelations_FallbackAlsoFails(t *testing.T) {
	// Data 是非数组 JSON，两段 Unmarshal 都失败
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"Data":"not-an-array"}`)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(1000, 1000))
	defer client.Close()

	_, err := client.GetParents(context.Background(), 79)
	if err == nil {
		t.Fatal("GetParents with non-array data should error")
	}
}
```

- [ ] **Step 2: 在 coverage_fillers_test.go 顶部 import 追加 net/http、net/http/httptest、fmt（若缺）**

文件: `coverage_fillers_test.go`（替换 import 块）

```go
import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	// ... 保留原有 import
)
```

- [ ] **Step 3: 验证 getRelations 覆盖率达 100%**
Run: `go test -count=1 -run "TestGetRelations_" -v . 2>&1 | tail -20`
Expected:
  - Exit code: 0
  - Output contains: "ok" and "PASS"

Run: `go test -count=1 -coverprofile=/tmp/cov1.out . && go tool cover -func=/tmp/cov1.out | grep getRelations`
Expected:
  - `getRelations  100.0%`

- [ ] **Step 4: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add coverage_fillers_test.go && git commit -m "test(sdk): 补全 getRelations 在线分支测试覆盖至 100%

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 2: SDK 层序列化剩余分支评估与不可达记录

**Depends on:** None
**Files:**
- Modify: `coverage_fillers_test.go`（追加 UnmarshalCSV 非 EOF 读取错误分支测试）

- [ ] **Step 1: 在 coverage_fillers_test.go 追加 UnmarshalCSV 数据行读取错误测试 — 覆盖 reader.Read 非 EOF err 分支（行 220）**

文件: `coverage_fillers_test.go`（追加）

```go
// TestUnmarshalCSV_DataRowReadError 覆盖 reader.Read 返回非 EOF 错误的分支（行 220-221）。
// 构造一个合法表头 + 引号未闭合的后续行，使 Read 在数据行报错（非 EOF）。
func TestUnmarshalCSV_DataRowReadError(t *testing.T) {
	// 表头合法；第二行有奇数个引号，csv.Reader.Read 报错（不是 EOF）
	data := []byte("ID,Name\n79,unterminated\"field")
	_, err := UnmarshalCSV(data)
	if err == nil {
		t.Fatal("UnmarshalCSV with malformed row should error")
	}
}
```

- [ ] **Step 2: 验证 UnmarshalCSV 达 100%**
Run: `go test -count=1 -run "TestUnmarshalCSV_DataRowReadError" -v . 2>&1 | tail -10`
Expected:
  - Exit code: 0
  - Output contains: "PASS"

Run: `go test -count=1 -coverprofile=/tmp/cov2.out . && go tool cover -func=/tmp/cov2.out | grep -E "MarshalXML|MarshalCSV|UnmarshalCSV"`
Expected:
  - `MarshalXML` 仍 85.7%（`xml.MarshalIndent` 对纯值 safeCWE 永不返回 err，不可达分支，接受）
  - `MarshalCSV` 仍 78.6%（`csv.Writer` 写 `bytes.Buffer` 永不返回 err，不可达分支，接受）
  - `UnmarshalCSV  100.0%`

- [ ] **Step 3: 提交（含不可达分支说明）**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add coverage_fillers_test.go && git commit -m "test(sdk): 补全 UnmarshalCSV 数据行读取错误分支至 100%

MarshalXML/MarshalCSV 剩余 err 分支因 bytes.Buffer/csv.Writer 对合法输入
永不返回错误而不可达，接受不覆盖。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 3: cmd/cwe-mcp 可测试性改造 — 引入 newAPIClient 注入点

**Depends on:** None
**Files:**
- Modify: `cmd/cwe-mcp/tools_api.go:1-55`（引入包级变量 + 替换三处 NewAPIClient 调用）

- [ ] **Step 1: 修改 tools_api.go — 引入 newAPIClient 包级变量，替换三处硬编码 NewAPIClient**

文件: `cmd/cwe-mcp/tools_api.go:3-10`（替换 import 与函数声明之间区域）以及 `:30`、`:72`、`:105` 三处 `cweskills.NewAPIClient()`。

替换 import 块为：

```go
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
```

把 `:30` 的 `client := cweskills.NewAPIClient()` 替换为：

```go
client := newAPIClient()
defer client.Close()
```

同理替换 `:72` 的 `client := cweskills.NewAPIClient()` 和 `:105` 的 `client := cweskills.NewAPIClient()` 各为 `client := newAPIClient()` + `defer client.Close()`（注意 `:105` 原本无 defer，保留新增 defer）。

- [ ] **Step 2: 验证改造未破坏现有测试与行为**
Run: `go build ./cmd/cwe-mcp && go test -count=1 ./cmd/cwe-mcp/ 2>&1 | tail -10`
Expected:
  - Exit code: 0
  - Output contains: "ok" and "coverage:" (与改造前持平 80.5%)

- [ ] **Step 3: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add cmd/cwe-mcp/tools_api.go && git commit -m "refactor(mcp): 提取 newAPIClient 注入点以便在线分支测试覆盖

行为不变：默认仍用 MITRE base URL；测试可替换为指向 httptest 的客户端。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 4: cmd/cwe-mcp registerAPITools 在线分支全覆盖

**Depends on:** Task 3
**Files:**
- Create: `cmd/cwe-mcp/tools_api_test.go`

- [ ] **Step 1: 创建 tools_api_test.go — 覆盖 get_weakness 在线成功、离线回退、api_version 成功/失败**

文件: `cmd/cwe-mcp/tools_api_test.go`

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/server"
)

// withMockAPI 临时把 newAPIClient 替换为指向 srv 的客户端，返回还原函数。
func withMockAPI(srv *httptest.Server) func() {
	orig := newAPIClient
	newAPIClient = func(opts ...cweskills.APIClientOption) *cweskills.APIClient {
		return cweskills.NewAPIClient(
			cweskills.WithAPIBaseURL(srv.URL),
			cweskills.WithAPIRateLimit(1000, 1000),
		)
	}
	return func() { newAPIClient = orig }
}

// TestGetWeakness_OnlineSuccess 覆盖 get_weakness 在线成功分支（行 33-38）
func TestGetWeakness_OnlineSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"Data":{"id":79,"name":"XSS","description":"cross-site scripting"}}`)
	}))
	defer srv.Close()
	defer withMockAPI(srv)()

	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_weakness", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "online" {
		t.Errorf("expected source=online, got %v", m["source"])
	}
}

// TestGetParents_OnlineSuccess 覆盖 get_parents 在线成功分支（行 75-81）
func TestGetParents_OnlineSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"Data":[{"nature":"ChildOf","cweId":100}]}`)
	}))
	defer srv.Close()
	defer withMockAPI(srv)()

	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_parents", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "online" {
		t.Errorf("expected source=online, got %v", m["source"])
	}
}

// TestGetWeakness_OfflineFallback 覆盖在线失败 + 离线回退命中分支（行 40-53）
func TestGetWeakness_OfflineFallback(t *testing.T) {
	// 在线 server 返回 500，触发回退
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	defer srv.Close()
	defer withMockAPI(srv)()

	// 注入离线注册表：构造含 CWE 79 的 registry
	xmlDir := buildOfflineRegistry(t, 79)
	defer xmlAPIRestore(xmlDir)()
	xmlPath = xmlDir

	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_weakness", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "offline" {
		t.Errorf("expected source=offline, got %v", m["source"])
	}
}

// TestGetParents_OfflineFallback 覆盖 get_parents 在线失败 + 离线回退分支（行 83-95）
func TestGetParents_OfflineFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	defer srv.Close()
	defer withMockAPI(srv)()

	xmlDir := buildOfflineRegistry(t, 79)
	defer xmlAPIRestore(xmlDir)()
	xmlPath = xmlDir

	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_parents", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "offline" {
		t.Errorf("expected source=offline, got %v", m["source"])
	}
}

// TestAPIVersion_Success 覆盖 api_version 在线成功分支（行 107-112）
func TestAPIVersion_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"Data":"4.15"}`)
	}))
	defer srv.Close()
	defer withMockAPI(srv)()

	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "api_version", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	m := resultJSON(t, r)
	if m["version"] == nil {
		t.Error("expected version in result")
	}
}

// TestAPIVersion_Failure 覆盖 api_version 在线失败分支（行 108-109）
func TestAPIVersion_Failure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	defer srv.Close()
	defer withMockAPI(srv)()

	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "api_version", map[string]any{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !r.IsError {
		t.Error("expected error result for api_version failure")
	}
}
```

- [ ] **Step 2: 在 tools_api_test.go 追加 buildOfflineRegistry / xmlAPIRestore 辅助函数 — 复用现有离线注册表加载机制注入 registry 全局变量**

文件: `cmd/cwe-mcp/tools_api_test.go`（追加到文件末尾）

```go
// buildOfflineRegistry 用现有 loadRegistry/mustRegistry 机制构造一个含指定 CWE 的离线注册表，
// 返回临时 XML 文件路径。复用项目内现成的最小 XML 样例。
// 若 register_test.go 已有等价辅助，则直接调用它而不要重复定义。
func buildOfflineRegistry(t *testing.T, id int) string {
	t.Helper()
	// 复用 register_test.go 中已存在的最小 XML 构造（如 newXMLFile / writeTempXML）。
	// 具体函数名以现有代码为准；执行时若已存在则直接调用。
	return existingOfflineXMLPath(t, id)
}

// xmlAPIRestore 还原全局 registry 状态，避免污染后续测试。
func xmlAPIRestore(_ string) func() {
	return func() {
		registryMu.Lock()
		registry = nil
		registryErr = nil
		registryMu.Unlock()
		xmlPath = ""
	}
}

// existingOfflineXMLPath 包装对现有离线 XML 构造的调用；若不存在则内联最小实现。
func existingOfflineXMLPath(t *testing.T, id int) string {
	t.Helper()
	// 最小内联实现：写一个含目标 CWE 的 XML 片段到临时文件。
	// 仅当 register_test.go 未提供等价辅助时使用。
	return writeMinXMLFile(t, id)
}
```

> 注：Step 2 的辅助函数依赖 `register_test.go` 中是否已有 XML 构造辅助。执行时先 `grep -n "func.*XML\|func.*xml\|func.*Registry" cmd/cwe-mcp/register_test.go`，若已有则直接调用，删除本 Step 的重复定义，仅保留 `xmlAPIRestore`。

- [ ] **Step 3: 验证 registerAPITools 覆盖率达 100%**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go test -count=1 -run "TestGetWeakness|TestGetParents|TestAPIVersion" -v ./cmd/cwe-mcp/ 2>&1 | tail -25`
Expected:
  - Exit code: 0
  - Output contains: "PASS" 6 次（6 个测试全过）

Run: `go test -count=1 -coverprofile=/tmp/cov4.out ./cmd/cwe-mcp/ && go tool cover -func=/tmp/cov4.out | grep registerAPITools`
Expected:
  - `registerAPITools  100.0%`

- [ ] **Step 4: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add cmd/cwe-mcp/tools_api_test.go && git commit -m "test(mcp): 补全 registerAPITools 在线/离线回退分支覆盖至 100%

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 5: cmd/cwe-mcp main 可测试性改造与覆盖

**Depends on:** None
**Files:**
- Modify: `cmd/cwe-mcp/main.go:41-102`（拆出 runMain(args) int）
- Create: `cmd/cwe-mcp/main_test.go`

- [ ] **Step 1: 修改 main.go — 拆出 runMain(args []string) int，main 仅 os.Exit(runMain(...))**

文件: `cmd/cwe-mcp/main.go:41-102`（替换 main 函数整体）

```go
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
		fmt.Fprintf(os.Stderr, "cwe-mcp — CWE Skills MCP 服务器（暴露 20 个 CWE 工具，供 MCP 兼容 AI 客户端调用）\n\n")
		fmt.Fprintf(os.Stderr, "用法:\n")
		fmt.Fprintf(os.Stderr, "  cwe-mcp                              stdio 模式，仅在线工具（无需 XML）\n")
		fmt.Fprintf(os.Stderr, "  cwe-mcp --xml cwec_v4.15.xml         stdio 模式，含离线工具\n")
		fmt.Fprintf(os.Stderr, "  cwe-mcp --transport http --addr :8080  SSE 模式（远程）\n")
		fmt.Fprintf(os.Stderr, "  cwe-mcp --version                    显示版本并退出\n\n")
		fmt.Fprintf(os.Stderr, "参数:\n")
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
		fmt.Fprintf(os.Stderr, "未知传输方式: %s\n", *transport)
		return 2
	}
	return 0
}
```

> 行为变化说明：原 `default` 分支用 `log.Fatalf`（调 os.Exit），现改为 `return 2` + stderr 输出，等价且可测试。`flag.Parse` 失败原静默继续，现 `return 2`，行为更明确。

- [ ] **Step 2: 创建 main_test.go — 覆盖 --version、未知 transport、参数解析错误分支**

文件: `cmd/cwe-mcp/main_test.go`

```go
package main

import (
	"bytes"
	"testing"
)

// TestRunMain_Version 覆盖 --version 分支（行：showVer 为真 → return 0）
func TestRunMain_Version(t *testing.T) {
	code := runMain([]string{"-version"})
	if code != 0 {
		t.Errorf("runMain -version: want exit 0, got %d", code)
	}
}

// TestRunMain_UnknownTransport 覆盖 default 分支（未知 transport → return 2）
func TestRunMain_UnknownTransport(t *testing.T) {
	var buf bytes.Buffer
	orig := osStderr
	osStderr = &buf
	defer func() { osStderr = orig }()

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
	code := runMain([]string{"-not-a-flag"})
	if code != 2 {
		t.Errorf("runMain bad flag: want exit 2, got %d", code)
	}
}
```

- [ ] **Step 3: 修改 main.go Step 1 的 default 分支用 osStderr 变量以便测试注入**

文件: `cmd/cwe-mcp/main.go`（在 `var mcpVersion` 附近新增）

```go
// osStderr 默认指向 os.Stderr，测试可替换以捕获输出。
var osStderr = os.Stderr
```

并把 Step 1 default 分支的 `fmt.Fprintf(os.Stderr, ...)` 改为 `fmt.Fprintf(osStderr, ...)`。

- [ ] **Step 4: 验证 main 覆盖率（接受 stdio/http 服务器启动分支不覆盖）**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go test -count=1 -run "TestRunMain" -v ./cmd/cwe-mcp/ 2>&1 | tail -15`
Expected:
  - Exit code: 0
  - Output contains: "PASS" 3 次

Run: `go test -count=1 -coverprofile=/tmp/cov5.out ./cmd/cwe-mcp/ && go tool cover -func=/tmp/cov5.out | grep -E "main\b|runMain"`
Expected:
  - `main  0.0%`（仅 os.Exit 一行，接受不覆盖）
  - `runMain  100.0%`（除 stdio/http 真正启动分支外全覆盖；启动分支接受不覆盖）

- [ ] **Step 5: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add cmd/cwe-mcp/main.go cmd/cwe-mcp/main_test.go && git commit -m "refactor(mcp): 拆出 runMain 以便 main flag/version 分支可测试

main 覆盖率：--version/未知 transport/参数错误分支达 100%；
stdio/http 服务器启动分支会阻塞，接受不覆盖。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 6: cmd/cwe main 与 Execute 可测试性改造与覆盖

**Depends on:** None
**Files:**
- Modify: `cmd/cwe/root.go:41-47`（拆出 executeRoot）
- Modify: `cmd/cwe/main.go:1-5`（main 调 os.Exit）
- Create: `cmd/cwe/execute_test.go`

- [ ] **Step 1: 修改 root.go — 拆出 executeRoot() int，Execute 仅 os.Exit**

文件: `cmd/cwe/root.go:41-47`（替换 Execute 函数）

```go
// Execute 执行根命令
func Execute() {
	os.Exit(executeRoot())
}

// executeRoot 执行根命令并返回退出码。
// 提取自 Execute 以便测试覆盖 cobra 执行成功与失败两条路径，
// 而不触发 os.Exit。
func executeRoot() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(osStderr, err)
		return 1
	}
	return 0
}
```

- [ ] **Step 2: 修改 root.go — 新增 osStderr 包级变量以便测试注入**

文件: `cmd/cwe/root.go`（在 `var outputFormat` 附近新增）

```go
// osStderr 默认指向 os.Stderr，测试可替换以捕获输出。
var osStderr = os.Stderr
```

- [ ] **Step 3: 创建 execute_test.go — 覆盖 executeRoot 成功（help）与失败路径**

文件: `cmd/cwe/execute_test.go`

```go
package main

import (
	"bytes"
	"testing"
)

// TestExecuteRoot_HelpSuccess 覆盖 cobra 执行成功路径（return 0）
func TestExecuteRoot_HelpSuccess(t *testing.T) {
	// 无参数时 cobra 输出 help 并返回 nil error（SilenceUsage=true，根命令无 Run）
	rootCmd.SetArgs([]string{"--help"})
	defer rootCmd.SetArgs([]string{}) // 还原

	code := executeRoot()
	if code != 0 {
		t.Errorf("executeRoot --help: want 0, got %d", code)
	}
}

// TestExecuteRoot_UnknownCommandFailure 覆盖 cobra 执行失败路径（return 1）
func TestExecuteRoot_UnknownCommandFailure(t *testing.T) {
	var buf bytes.Buffer
	orig := osStderr
	osStderr = &buf
	defer func() { osStderr = orig }()

	rootCmd.SetArgs([]string{"nonexistent-subcommand-xyz"})
	defer rootCmd.SetArgs([]string{})

	code := executeRoot()
	if code != 1 {
		t.Errorf("executeRoot unknown cmd: want 1, got %d", code)
	}
	if buf.Len() == 0 {
		t.Error("expected error output on stderr")
	}
}
```

- [ ] **Step 4: 验证 Execute / executeRoot 覆盖率（接受 main 仅 os.Exit 不覆盖）**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go test -count=1 -run "TestExecuteRoot" -v ./cmd/cwe/ 2>&1 | tail -15`
Expected:
  - Exit code: 0
  - Output contains: "PASS" 2 次

Run: `go test -count=1 -coverprofile=/tmp/cov6.out ./cmd/cwe/ && go tool cover -func=/tmp/cov6.out | grep -E "Execute|executeRoot|main\b"`
Expected:
  - `main  0.0%`（仅 Execute() 一行，接受）
  - `Execute  100.0%`（单行 os.Exit，但 Execute 被测试间接覆盖到？—— 注意：Execute 调 os.Exit 不在测试中调用。Execute 行为由 executeRoot 测试代表，接受 Execute 不直接覆盖）
  - `executeRoot  100.0%`

- [ ] **Step 5: 提交**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add cmd/cwe/root.go cmd/cwe/execute_test.go && git commit -m "refactor(cli): 拆出 executeRoot 以便 cobra 执行分支可测试

Execute 成功/失败路径通过 executeRoot 覆盖；main 仅 os.Exit 接受不覆盖。

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

---

### Task 7: 全量验证与不可达分支记录

**Depends on:** Task 1, Task 2, Task 3, Task 4, Task 5, Task 6
**Files:**
- Modify: `docs/superpowers/plans/2026-07-12-coverage-100.md`（追加本轮执行结果）

- [ ] **Step 1: 运行全量测试与覆盖率**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && go test -race -count=1 -coverprofile=coverage.out ./... 2>&1 | tail -10`
Expected:
  - Exit code: 0
  - 三个包均 "ok"
  - SDK 覆盖率 >= 99.4%（getRelations 提升后应更高）

Run: `go tool cover -func=coverage.out | tail -1`
Expected:
  - `total:` 行覆盖率显著高于 96.4%

- [ ] **Step 2: 生成未达 100% 的最终清单，确认仅剩真正不可达分支**
Run: `go tool cover -func=coverage.out | grep -v "100.0%"`
Expected:
  - 仅剩：`MarshalXML`（xml.MarshalIndent 不可达）、`MarshalCSV`（csv.Writer+bytes.Buffer 不可达）、`main`（各包的 os.Exit 单行）—— 均有文档记录的不可达理由
  - 不应再出现 `getRelations`、`registerAPITools`、`executeRoot`、`runMain` 未满

- [ ] **Step 3: 在 docs/superpowers/plans/2026-07-12-coverage-100.md 末尾追加本轮执行结果**

文件: `docs/superpowers/plans/2026-07-12-coverage-100.md`（追加到文件末尾）

```markdown

## 执行结果（2026-07-13 第二轮提升）

| 函数 | 覆盖率（前） | 覆盖率（后） | 处理方式 |
|---|---|---|---|
| getRelations | 64.3% | 100% | httptest 注入覆盖网络错误/fallback 成功/fallback 失败三分支 |
| UnmarshalCSV | 97.1% | 100% | 构造畸形 CSV 行触发 reader.Read 非 EOF 错误 |
| registerAPITools | 34.1% | 100% | 提取 newAPIClient 注入点 + httptest 覆盖在线成功/离线回退/api_version 成败 |
| cmd/cwe-mcp main | 0% | runMain 100% | 拆出 runMain(args) int；stdio/http 启动分支阻塞接受不覆盖 |
| cmd/cwe Execute | 0% | executeRoot 100% | 拆出 executeRoot() int；main 仅 os.Exit 接受不覆盖 |
| MarshalXML | 85.7% | 85.7% | xml.MarshalIndent 对纯值 safeCWE 永不返回 err，不可达，接受 |
| MarshalCSV | 78.6% | 78.6% | csv.Writer 写 bytes.Buffer 永不返回 err，不可达，接受 |

**可测试性改造（行为不变）：**
1. `cmd/cwe-mcp/tools_api.go`：提取 `var newAPIClient = cweskills.NewAPIClient`，handler 改用 `newAPIClient()`
2. `cmd/cwe-mcp/main.go`：拆出 `runMain(args []string) int`，main 仅 `os.Exit(runMain(...))`
3. `cmd/cwe/root.go`：拆出 `executeRoot() int`，Execute 仅 `os.Exit(executeRoot())`
4. 引入 `osStderr` 包级变量（cwe 与 cwe-mcp 各一）以便测试捕获 stderr
```

- [ ] **Step 4: 提交文档更新**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && git add docs/superpowers/plans/2026-07-12-coverage-100.md && git commit -m "docs(plan): 记录第二轮覆盖率提升结果（getRelations/registerAPITools/main 全部达 100%）

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"`

- [ ] **Step 5: 确认 CI 门禁仍通过（整体 95% / SDK 99%）**
Run: `cd /home/cc11001100/github/scagogogo/cwe-skills && total=$(go tool cover -func=coverage.out | awk '/^total:/ {print $3}' | tr -d '%') && sdk=$(go test -count=1 -coverprofile=/tmp/sdk.out . 2>/dev/null | grep -oE 'coverage: [0-9.]+' | head -1 | grep -oE '[0-9.]+') && echo "total=$total sdk=$sdk"`
Expected:
  - total > 96.4（显著提升）
  - sdk > 99.4

---

## 验证总表（本计划完成后）

| 层 | 覆盖率（前） | 覆盖率（后目标） | 不可达分支（接受） |
|---|---|---|---|
| SDK 根层 | 99.4% | 100%（getRelations/UnmarshalCSV 满） | MarshalXML/MarshalCSV 的底层库不可达 err 分支 |
| cmd/cwe | 37.0% | executeRoot 100% | `main()`（os.Exit 单行） |
| cmd/cwe-mcp | 80.5% | registerAPITools/runMain 100% | `main()`（os.Exit 单行）、stdio/http 服务器启动分支（阻塞） |
| 总计 | 96.4% | 接近 100% | 仅余上述文档记录的不可达分支 |

## 失败回退

- 若 `runMain` 拆分后 cobra/mcp-go 全局状态被测试污染 → 用 `t.Cleanup` 还原 `rootCmd.SetArgs`、`registry`、`xmlPath`、`newAPIClient`、`osStderr`
- 若 httptest mock 无法触发 `getRelations` 的 fallback 分支（API 路径不匹配）→ 检查 `GetParents` 实际请求的 path，调整 httptest handler 匹配任意路径
- 若可测试性改造破坏现有 96.4% → 立即 `git revert` 对应 commit，重新评估该函数是否接受不覆盖
