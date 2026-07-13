package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/server"
)

// withMockAPI 临时把 newAPIClient 替换为指向 srv 的客户端，返回还原函数。
// 调用方应：defer withMockAPI(srv)() 先于 defer srv.Close()。
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

// mockAPIResp 构造一个始终以 200 返回固定 JSON 体的 httptest server。
func mockAPIResp(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

// mockAPIFail 构造一个始终返回 500 的 httptest server（触发在线失败回退）。
func mockAPIFail() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"server error"}`))
	}))
}

// newAPIServer 注册全部工具到新 server。
func newAPIServer() *server.MCPServer {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)
	return s
}

// ==================== get_weakness ====================

func TestGetWeakness_OnlineSuccess(t *testing.T) {
	srv := mockAPIResp(`{"Data":{"id":79,"name":"XSS","description":"Cross-site scripting"}}`)
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "get_weakness", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "online" {
		t.Errorf("expected source=online, got %v", m["source"])
	}
	if _, ok := m["cwe"]; !ok {
		t.Errorf("expected cwe field, got %v", m)
	}
}

func TestGetWeakness_OfflineFallback(t *testing.T) {
	setupOfflineRegistry(t)
	srv := mockAPIFail()
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "get_weakness", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "offline" {
		t.Errorf("expected source=offline, got %v", m["source"])
	}
	if _, ok := m["api_error"]; !ok {
		t.Errorf("expected api_error field, got %v", m)
	}
}

func TestGetWeakness_NoOfflineRegistry(t *testing.T) {
	resetRegistryStateForTest()
	srv := mockAPIFail()
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "get_weakness", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	if !r.IsError {
		t.Errorf("expected IsError when no offline registry, got result: %v", resultText(t, r))
	}
}

func TestGetWeakness_NotInOfflineRegistry(t *testing.T) {
	setupOfflineRegistry(t)
	srv := mockAPIFail()
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "get_weakness", map[string]any{"id": "CWE-9999"})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	if !r.IsError {
		t.Errorf("expected IsError for CWE not in offline registry, got result: %v", resultText(t, r))
	}
}

// ==================== get_parents ====================

func TestGetParents_OnlineSuccess(t *testing.T) {
	// []Relationship 可直接解析：nature/cweId 字段名匹配
	srv := mockAPIResp(`{"Data":[{"nature":"ChildOf","cweId":74}]}`)
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "get_parents", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "online" {
		t.Errorf("expected source=online, got %v", m["source"])
	}
	if m["cwe_id"] != "CWE-79" {
		t.Errorf("expected cwe_id=CWE-79, got %v", m["cwe_id"])
	}
	parents, _ := m["parents"].([]any)
	if len(parents) == 0 {
		t.Errorf("expected non-empty parents, got %v", m["parents"])
	}
}

func TestGetParents_OfflineFallback(t *testing.T) {
	setupOfflineRegistry(t)
	srv := mockAPIFail()
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "get_parents", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	m := resultJSON(t, r)
	if m["source"] != "offline" {
		t.Errorf("expected source=offline, got %v", m["source"])
	}
	if m["cwe_id"] != "CWE-79" {
		t.Errorf("expected cwe_id=CWE-79, got %v", m["cwe_id"])
	}
	parents, _ := m["parents"].([]any)
	if len(parents) == 0 {
		t.Errorf("expected non-empty offline parents (CWE-79 ChildOf 74), got %v", m["parents"])
	}
}

func TestGetParents_NoOfflineRegistry(t *testing.T) {
	resetRegistryStateForTest()
	srv := mockAPIFail()
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "get_parents", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	if !r.IsError {
		t.Errorf("expected IsError when no offline registry, got result: %v", resultText(t, r))
	}
}

// ==================== api_version ====================

func TestAPIVersion_Success(t *testing.T) {
	srv := mockAPIResp(`{"Data":{"version":"4.15"}}`)
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "api_version", map[string]any{})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	m := resultJSON(t, r)
	if m["version"] == nil {
		t.Errorf("expected version field, got %v", m)
	}
}

func TestAPIVersion_Failure(t *testing.T) {
	srv := mockAPIFail()
	defer withMockAPI(srv)()
	defer srv.Close()

	s := newAPIServer()
	r, err := callTool(s, "api_version", map[string]any{})
	if err != nil {
		t.Fatalf("callTool error: %v", err)
	}
	if !r.IsError {
		t.Errorf("expected IsError on API failure, got result: %v", resultText(t, r))
	}
}
