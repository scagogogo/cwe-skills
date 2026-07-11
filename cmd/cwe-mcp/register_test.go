package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	cweskills "github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// getHandler 通过反射从 MCPServer 取出已注册工具的 handler。
// mcp-go v0.20.0 的 MCPServer.tools 是未导出 map[string]ServerTool，
// 反射读出的 map 值不可直接 Interface()（CanInterface()==false），
// 因此用 unsafe.Pointer 把字段地址转成 *map 取出 ServerTool（其字段均导出）。
func getHandler(s *server.MCPServer, name string) server.ToolHandlerFunc {
	v := reflect.ValueOf(s).Elem()
	toolsField := v.FieldByName("tools")
	if !toolsField.IsValid() {
		panic("MCPServer.tools field not found via reflection")
	}
	// 用 unsafe 取未导出 map 字段的指针
	toolsPtr := unsafe.Pointer(toolsField.UnsafeAddr())
	tools := *(*map[string]server.ServerTool)(toolsPtr)
	entry, ok := tools[name]
	if !ok {
		panic("tool " + name + " not registered")
	}
	return entry.Handler
}

// callTool 调用指定工具，args 为参数 map。
func callTool(s *server.MCPServer, name string, args map[string]any) (*mcp.CallToolResult, error) {
	h := getHandler(s, name)
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return h(context.Background(), req)
}

// newServerWith 注册全部工具到新 server
func newServerWith() *server.MCPServer {
	s := server.NewMCPServer("test", "0", server.WithToolCapabilities(true))
	registerIDTools(s)
	registerWellknownTools(s)
	registerAPITools(s)
	registerOfflineTools(s)
	registerExtraTools(s)
	return s
}

// resultText 提取工具结果的文本内容
func resultText(t *testing.T, r *mcp.CallToolResult) string {
	t.Helper()
	if len(r.Content) == 0 {
		return ""
	}
	if tc, ok := r.Content[0].(mcp.TextContent); ok {
		return tc.Text
	}
	return ""
}

// resultJSON 把工具结果文本解析为 map
func resultJSON(t *testing.T, r *mcp.CallToolResult) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(resultText(t, r)), &m); err != nil {
		t.Fatalf("result not JSON: %v\nraw: %s", err, resultText(t, r))
	}
	return m
}

// setupOfflineRegistry 写一个最小 CWE XML 并加载到全局 registry。
func setupOfflineRegistry(t *testing.T) {
	t.Helper()
	oldXML, oldReg, oldErr := xmlPath, registry, registryErr
	t.Cleanup(func() { xmlPath, registry, registryErr = oldXML, oldReg, oldErr })

	xml := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="Test" Version="7.0" Date="2026-01-01">
  <Weaknesses>
    <Weakness ID="79" Name="XSS" Abstraction="Variant" Structure="Simple" Status="Stable">
      <Description>Cross-site scripting</Description>
      <Relationships>
        <Relationship Nature="ChildOf" CWE_ID="74" View_ID="1000"/>
      </Relationships>
    </Weakness>
    <Weakness ID="74" Name="Injection" Abstraction="Base" Structure="Simple" Status="Stable">
      <Description>Injection</Description>
      <Relationships>
        <Relationship Nature="ChildOf" CWE_ID="707" View_ID="1000"/>
      </Relationships>
    </Weakness>
    <Weakness ID="707" Name="InjectionGeneral" Abstraction="Class" Structure="Simple" Status="Stable">
      <Description>General injection</Description>
    </Weakness>
  </Weaknesses>
</Weakness_Catalog>`

	path := filepath.Join(t.TempDir(), "test.xml")
	if err := os.WriteFile(path, []byte(xml), 0644); err != nil {
		t.Fatalf("write xml: %v", err)
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = nil
	registryErr = nil
	xmlPath = path
	if err := loadRegistry(); err != nil {
		t.Fatalf("loadRegistry: %v", err)
	}
}

// resetRegistryStateForTest 清空全局 registry 状态（用于测无 XML 分支）
func resetRegistryStateForTest() {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = nil
	registryErr = nil
	xmlPath = ""
}

// ==================== registerIDTools ====================

func TestRegisterIDTools_Parse(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "parse_cwe_id", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("parse_cwe_id error: %v", err)
	}
	m := resultJSON(t, r)
	if m["valid"] != true {
		t.Errorf("expected valid=true, got %v", m["valid"])
	}
	if m["format"] != "CWE-79" {
		t.Errorf("expected CWE-79, got %v", m["format"])
	}
}

func TestRegisterIDTools_ParseMissingID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "parse_cwe_id", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterIDTools_ParseInvalidID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "parse_cwe_id", map[string]any{"id": "abc"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

func TestRegisterIDTools_Validate(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "validate_cwe_id", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["valid"] != true {
		t.Errorf("expected valid=true, got %v", m["valid"])
	}
	if m["error"] != "" {
		t.Errorf("expected empty error, got %v", m["error"])
	}
}

func TestRegisterIDTools_ValidateInvalid(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "validate_cwe_id", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["valid"] != false {
		t.Errorf("expected valid=false, got %v", m["valid"])
	}
	if m["error"] == "" {
		t.Error("expected non-empty error")
	}
}

func TestRegisterIDTools_ValidateMissingID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "validate_cwe_id", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterIDTools_Extract(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "extract_cwe_ids", map[string]any{"text": "CWE-79 and CWE-89"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["count"].(float64) != 2 {
		t.Errorf("expected count 2, got %v", m["count"])
	}
}

func TestRegisterIDTools_ExtractMissingText(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "extract_cwe_ids", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing text")
	}
}

func TestRegisterIDTools_Compare(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "compare_cwe_ids", map[string]any{"a": "CWE-79", "b": "CWE-89"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["comparison"].(float64) != -1 {
		t.Errorf("expected -1, got %v", m["comparison"])
	}
}

func TestRegisterIDTools_CompareMissingA(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "compare_cwe_ids", map[string]any{"b": "CWE-89"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing a")
	}
}

func TestRegisterIDTools_CompareMissingB(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "compare_cwe_ids", map[string]any{"a": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing b")
	}
}

func TestRegisterIDTools_CompareInvalid(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerIDTools(s)

	r, err := callTool(s, "compare_cwe_ids", map[string]any{"a": "bad", "b": "CWE-89"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid a")
	}
}

// ==================== registerWellknownTools ====================

func TestRegisterWellknown_Check(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerWellknownTools(s)

	// CWE-79 在多个知名列表
	r, err := callTool(s, "check_wellknown", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	inList, _ := m["in_list"].([]any)
	if len(inList) == 0 {
		t.Error("expected CWE-79 to be in some well-known list")
	}
}

func TestRegisterWellknown_CheckMissingID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerWellknownTools(s)

	r, err := callTool(s, "check_wellknown", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterWellknown_CheckInvalidID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerWellknownTools(s)

	r, err := callTool(s, "check_wellknown", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

func TestRegisterWellknown_CheckNotInList(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerWellknownTools(s)

	// CWE-99999 不在任何列表
	r, err := callTool(s, "check_wellknown", map[string]any{"id": "99999"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	inList, _ := m["in_list"].([]any)
	if len(inList) != 0 {
		t.Errorf("expected empty in_list, got %v", inList)
	}
}

func TestRegisterWellknown_OWASPCategories(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerWellknownTools(s)

	r, err := callTool(s, "get_owasp_categories", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if _, ok := m["categories"]; !ok {
		t.Errorf("expected categories field, got %v", m)
	}
}

func TestRegisterWellknown_OWASPMissingID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerWellknownTools(s)

	r, err := callTool(s, "get_owasp_categories", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterWellknown_OWASPInvalidID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerWellknownTools(s)

	r, err := callTool(s, "get_owasp_categories", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

// ==================== registerExtraTools ====================

func TestRegisterExtra_FormatInt(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "format_cwe_id", map[string]any{"ids": []any{79, 89}})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	ids, _ := m["ids"].([]any)
	if len(ids) != 2 || ids[0] != "CWE-79" || ids[1] != "CWE-89" {
		t.Errorf("unexpected ids: %v", ids)
	}
}

func TestRegisterExtra_FormatMissingIDs(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "format_cwe_id", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing ids")
	}
}

func TestRegisterExtra_FormatInvalidID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "format_cwe_id", map[string]any{"ids": []any{"abc"}})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id value")
	}
}

func TestRegisterExtra_FormatEmptyList(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "format_cwe_id", map[string]any{"ids": []any{}})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	ids, _ := m["ids"].([]any)
	if len(ids) != 0 {
		t.Errorf("expected empty ids, got %v", ids)
	}
}

func TestRegisterExtra_Siblings(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "get_siblings", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	resultJSON(t, r) // 仅验证不报错、输出合法 JSON
}

func TestRegisterExtra_SiblingsMissingID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "get_siblings", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterExtra_SiblingsInvalidID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "get_siblings", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

func TestRegisterExtra_Children(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "get_children", map[string]any{"id": "CWE-74"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	children, _ := m["children"].([]any)
	if len(children) == 0 {
		t.Errorf("expected children for CWE-74, got %v", children)
	}
}

func TestRegisterExtra_ChildrenMissingID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "get_children", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterExtra_ChildrenInvalidID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "get_children", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

func TestRegisterExtra_Filter(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "filter_cwes", map[string]any{"abstraction": "Variant", "status": "Stable"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["count"].(float64) == 0 {
		t.Error("expected non-zero filtered count")
	}
}

func TestRegisterExtra_FilterNoMatch(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "filter_cwes", map[string]any{"abstraction": "Pillar"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["count"].(float64) != 0 {
		t.Errorf("expected 0, got %v", m["count"])
	}
}

func TestRegisterExtra_FilterEmptyArgs(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "filter_cwes", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["count"].(float64) == 0 {
		t.Error("expected non-zero count for empty filter (returns all)")
	}
}

func TestRegisterExtra_IsAncestor(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	// CWE-74 是 CWE-79 的祖先（79 ChildOf 74 ChildOf 707）
	r, err := callTool(s, "is_ancestor", map[string]any{"ancestor": "CWE-74", "descendant": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["is_ancestor"] != true {
		t.Errorf("expected is_ancestor=true, got %v", m["is_ancestor"])
	}
}

func TestRegisterExtra_IsAncestorNotAncestor(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "is_ancestor", map[string]any{"ancestor": "CWE-79", "descendant": "CWE-74"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["is_ancestor"] != false {
		t.Errorf("expected is_ancestor=false, got %v", m["is_ancestor"])
	}
}

func TestRegisterExtra_IsAncestorMissingAnc(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "is_ancestor", map[string]any{"descendant": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing ancestor")
	}
}

func TestRegisterExtra_IsAncestorMissingDesc(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "is_ancestor", map[string]any{"ancestor": "CWE-74"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing descendant")
	}
}

func TestRegisterExtra_IsAncestorInvalidAnc(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "is_ancestor", map[string]any{"ancestor": "bad", "descendant": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid ancestor")
	}
}

func TestRegisterExtra_IsAncestorInvalidDesc(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "is_ancestor", map[string]any{"ancestor": "CWE-74", "descendant": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid descendant")
	}
}

func TestRegisterExtra_NoXML_Siblings(t *testing.T) {
	// 无 XML 时 mustRegistry 返回错误
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)

	r, err := callTool(s, "get_siblings", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError when no XML loaded")
	}
}

// ==================== registerOfflineTools ====================

func TestRegisterOffline_Ancestors(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_ancestors", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	anc, _ := m["ancestors"].([]any)
	if len(anc) == 0 {
		t.Error("expected ancestors for CWE-79")
	}
}

func TestRegisterOffline_AncestorsMissingID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_ancestors", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterOffline_AncestorsInvalidID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_ancestors", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

func TestRegisterOffline_Descendants(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_descendants", map[string]any{"id": "CWE-74"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	desc, _ := m["descendants"].([]any)
	if len(desc) == 0 {
		t.Error("expected descendants for CWE-74")
	}
}

func TestRegisterOffline_DescendantsMissingID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_descendants", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterOffline_DescendantsInvalidID(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_descendants", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

func TestRegisterOffline_ShortestPath(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_shortest_path", map[string]any{"from": "CWE-79", "to": "CWE-707"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["found"] != true {
		t.Errorf("expected found=true, got %v", m["found"])
	}
}

func TestRegisterOffline_ShortestPathNoPath(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	// CWE-89 不在 registry，与 CWE-79 无路径
	r, err := callTool(s, "get_shortest_path", map[string]any{"from": "CWE-79", "to": "CWE-89"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["found"] != false {
		t.Errorf("expected found=false, got %v", m["found"])
	}
}

func TestRegisterOffline_ShortestPathMissingFrom(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_shortest_path", map[string]any{"to": "CWE-707"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing from")
	}
}

func TestRegisterOffline_ShortestPathMissingTo(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_shortest_path", map[string]any{"from": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing to")
	}
}

func TestRegisterOffline_ShortestPathInvalidFrom(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_shortest_path", map[string]any{"from": "bad", "to": "CWE-707"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid from")
	}
}

func TestRegisterOffline_ShortestPathInvalidTo(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_shortest_path", map[string]any{"from": "CWE-79", "to": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid to")
	}
}

func TestRegisterOffline_BuildTree(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "build_tree", map[string]any{"root": "CWE-74"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["root"] != "CWE-74" {
		t.Errorf("expected root CWE-74, got %v", m["root"])
	}
}

func TestRegisterOffline_BuildTreeNotFound(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "build_tree", map[string]any{"root": "99999"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for not-found root")
	}
}

func TestRegisterOffline_BuildTreeMissingRoot(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "build_tree", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing root")
	}
}

func TestRegisterOffline_BuildTreeInvalidRoot(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "build_tree", map[string]any{"root": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid root")
	}
}

func TestRegisterOffline_SearchKeyword(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "search_keyword", map[string]any{"keyword": "Injection"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["count"].(float64) == 0 {
		t.Error("expected non-zero count for 'Injection'")
	}
}

func TestRegisterOffline_SearchKeywordNoMatch(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "search_keyword", map[string]any{"keyword": "NoSuchThingHere"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["count"].(float64) != 0 {
		t.Errorf("expected 0, got %v", m["count"])
	}
}

func TestRegisterOffline_SearchKeywordMissing(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "search_keyword", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing keyword")
	}
}

func TestRegisterOffline_RegistryStats(t *testing.T) {
	setupOfflineRegistry(t)
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "registry_stats", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	m := resultJSON(t, r)
	if m["total"].(float64) != 3 {
		t.Errorf("expected total 3, got %v", m["total"])
	}
}

func TestRegisterOffline_NoXML_Ancestors(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)

	r, err := callTool(s, "get_ancestors", map[string]any{"id": "CWE-79"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError when no XML loaded")
	}
}

// ==================== registerAPITools（仅测参数验证分支，不触网） ====================

func TestRegisterAPI_GetWeakness_MissingID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_weakness", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterAPI_GetWeakness_InvalidID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_weakness", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

func TestRegisterAPI_GetParents_MissingID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_parents", map[string]any{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for missing id")
	}
}

func TestRegisterAPI_GetParents_InvalidID(t *testing.T) {
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerAPITools(s)

	r, err := callTool(s, "get_parents", map[string]any{"id": "bad"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !r.IsError {
		t.Error("expected IsError for invalid id")
	}
}

// ==================== 全量注册 smoke test ====================

func TestNewServerWith_AllToolsRegistered(t *testing.T) {
	s := newServerWith()
	// 确保所有工具都注册成功（反射能取到 handler 即说明已注册）
	for _, name := range []string{
		"parse_cwe_id", "validate_cwe_id", "extract_cwe_ids", "compare_cwe_ids",
		"check_wellknown", "get_owasp_categories",
		"get_weakness", "get_parents", "api_version",
		"get_ancestors", "get_descendants", "get_shortest_path", "build_tree",
		"search_keyword", "registry_stats",
		"format_cwe_id", "get_siblings", "get_children", "filter_cwes", "is_ancestor",
	} {
		// 不 panic 即说明已注册
		_ = getHandler(s, name)
	}
}

// 确保 cweskills 包被引用（某些测试可能不直接使用其类型）
var _ = cweskills.NewRegistry

// ==================== 无 XML 时各离线工具的错误分支 ====================

func TestRegisterOffline_NoXML_Descendants(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)
	r, err := callTool(s, "get_descendants", map[string]any{"id": "CWE-79"})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

func TestRegisterOffline_NoXML_ShortestPath(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)
	r, err := callTool(s, "get_shortest_path", map[string]any{"from": "CWE-79", "to": "CWE-89"})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

func TestRegisterOffline_NoXML_BuildTree(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)
	r, err := callTool(s, "build_tree", map[string]any{"root": "CWE-79"})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

func TestRegisterOffline_NoXML_SearchKeyword(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)
	r, err := callTool(s, "search_keyword", map[string]any{"keyword": "x"})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

func TestRegisterOffline_NoXML_RegistryStats(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerOfflineTools(s)
	r, err := callTool(s, "registry_stats", map[string]any{})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

func TestRegisterExtra_NoXML_Children(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)
	r, err := callTool(s, "get_children", map[string]any{"id": "CWE-79"})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

func TestRegisterExtra_NoXML_Filter(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)
	r, err := callTool(s, "filter_cwes", map[string]any{})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

func TestRegisterExtra_NoXML_IsAncestor(t *testing.T) {
	resetRegistryStateForTest()
	s := server.NewMCPServer("t", "0", server.WithToolCapabilities(true))
	registerExtraTools(s)
	r, err := callTool(s, "is_ancestor", map[string]any{"ancestor": "CWE-74", "descendant": "CWE-79"})
	if err != nil { t.Fatalf("error: %v", err) }
	if !r.IsError { t.Error("expected IsError when no XML loaded") }
}

// ==================== loadRegistry 无效 XML（触发 ParseFile 错误） ====================

func TestLoadRegistry_InvalidXML(t *testing.T) {
	oldXML, oldReg, oldErr := xmlPath, registry, registryErr
	defer func() { xmlPath, registry, registryErr = oldXML, oldReg, oldErr }()

	path := filepath.Join(t.TempDir(), "bad.xml")
	os.WriteFile(path, []byte("<not valid xml"), 0644)

	registryMu.Lock()
	defer registryMu.Unlock()
	registry = nil
	registryErr = nil
	xmlPath = path

	err := loadRegistry()
	if err == nil {
		t.Fatal("expected error for invalid XML, got nil")
	}
}
