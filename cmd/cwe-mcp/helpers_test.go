package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/scagogogo/cwe-skills"
	"github.com/mark3labs/mcp-go/mcp"
)

// ==================== main.go: jsonRaw / errResult ====================

func TestJsonRaw(t *testing.T) {
	raw := jsonRaw(map[string]any{"k": "v"})
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("jsonRaw output not valid JSON: %v", err)
	}
	if got["k"] != "v" {
		t.Errorf("expected k=v, got %v", got["k"])
	}
}

func TestJsonRaw_Nil(t *testing.T) {
	raw := jsonRaw(nil)
	if string(raw) != "null" {
		t.Errorf("expected null, got %q", raw)
	}
}

// jsonTypeError 是无法被 json.Marshal 编码的类型（含 chan 字段）
type jsonTypeError struct {
	Ch chan struct{}
}

func TestJsonRaw_MarshalError(t *testing.T) {
	// chan 无法被 json.Marshal，触发错误分支
	raw := jsonRaw(jsonTypeError{Ch: make(chan struct{})})
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("error fallback should be valid JSON: %v", err)
	}
	if got["error"] != "json marshal failed" {
		t.Errorf("expected error field, got %v", got["error"])
	}
	if _, ok := got["type"]; !ok {
		t.Errorf("expected type field, got %v", got)
	}
	if _, ok := got["reason"]; !ok {
		t.Errorf("expected reason field, got %v", got)
	}
}

func TestErrResult(t *testing.T) {
	r := errResult("boom")
	if r == nil {
		t.Fatal("expected non-nil result")
	}
	if !r.IsError {
		t.Error("expected IsError=true")
	}
}

// ==================== tools_id.go: wrapJSON / errMsg / requireStringArg ====================

func TestWrapJSON(t *testing.T) {
	r, err := wrapJSON(map[string]any{"k": "v"})
	if err != nil {
		t.Fatalf("wrapJSON error: %v", err)
	}
	if len(r.Content) == 0 {
		t.Fatal("expected non-empty content")
	}
	// 第一个 content item 应是 text 类型
	tc, ok := r.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", r.Content[0])
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(tc.Text), &got); err != nil {
		t.Fatalf("wrapJSON text not valid JSON: %v", err)
	}
	if got["k"] != "v" {
		t.Errorf("expected k=v, got %v", got["k"])
	}
}

func TestErrMsg_Nil(t *testing.T) {
	if got := errMsg(nil); got != "" {
		t.Errorf("expected empty for nil, got %q", got)
	}
}

func TestErrMsg_NonNil(t *testing.T) {
	e := errors.New("boom")
	if got := errMsg(e); got != "boom" {
		t.Errorf("expected boom, got %q", got)
	}
}

func TestRequireStringArg_OK(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"id": "CWE-79"}
	v, ok := requireStringArg(req, "id")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if v != "CWE-79" {
		t.Errorf("expected CWE-79, got %q", v)
	}
}

func TestRequireStringArg_Empty(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"id": ""}
	if _, ok := requireStringArg(req, "id"); ok {
		t.Error("expected ok=false for empty string")
	}
}

func TestRequireStringArg_Missing(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}
	if _, ok := requireStringArg(req, "id"); ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestRequireStringArg_WrongType(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"id": 123}
	if _, ok := requireStringArg(req, "id"); ok {
		t.Error("expected ok=false for non-string type")
	}
}

// ==================== tools_extra.go: toInt ====================

func TestToInt_Int(t *testing.T) {
	v, ok := toInt(42)
	if !ok || v != 42 {
		t.Errorf("toInt(42) = %d, %v, want 42, true", v, ok)
	}
}

func TestToInt_Int64(t *testing.T) {
	v, ok := toInt(int64(7))
	if !ok || v != 7 {
		t.Errorf("toInt(int64(7)) = %d, %v, want 7, true", v, ok)
	}
}

func TestToInt_Float64(t *testing.T) {
	// JSON 数字解码为 float64
	v, ok := toInt(float64(79))
	if !ok || v != 79 {
		t.Errorf("toInt(float64(79)) = %d, %v, want 79, true", v, ok)
	}
}

func TestToInt_Invalid(t *testing.T) {
	if _, ok := toInt("abc"); ok {
		t.Error("expected ok=false for string")
	}
	if _, ok := toInt(nil); ok {
		t.Error("expected ok=false for nil")
	}
	// 注意：toInt 对任意 float64 都返回 true（含小数，int(n) 截断）
	if v, ok := toInt(3.7); !ok || v != 3 {
		t.Errorf("toInt(3.7) = %d, %v, want 3, true (float64 截断)", v, ok)
	}
}

// ==================== tools_offline.go: serializeTreeNode ====================

func TestSerializeTreeNode_Nil(t *testing.T) {
	if got := serializeTreeNode(nil); got != nil {
		t.Errorf("expected nil for nil node, got %v", got)
	}
}

func TestSerializeTreeNode_Leaf(t *testing.T) {
	node := cweskills.NewTreeNode(&cweskills.CWE{ID: 79, Name: "XSS"})
	m := serializeTreeNode(node)
	if m["id"] != "CWE-79" {
		t.Errorf("expected CWE-79, got %v", m["id"])
	}
	if m["name"] != "XSS" {
		t.Errorf("expected XSS, got %v", m["name"])
	}
	children, ok := m["children"].([]map[string]any)
	if !ok || len(children) != 0 {
		t.Errorf("expected empty children slice, got %v", m["children"])
	}
}

func TestSerializeTreeNode_WithChildren(t *testing.T) {
	root := cweskills.NewTreeNode(&cweskills.CWE{ID: 74, Name: "Injection"})
	root.AddChild(cweskills.NewTreeNode(&cweskills.CWE{ID: 79, Name: "XSS"}))
	m := serializeTreeNode(root)
	children, ok := m["children"].([]map[string]any)
	if !ok || len(children) != 1 {
		t.Fatalf("expected 1 child, got %v", m["children"])
	}
	if children[0]["id"] != "CWE-79" {
		t.Errorf("expected child CWE-79, got %v", children[0]["id"])
	}
}

// ==================== main.go: loadRegistry / mustRegistry ====================

func resetRegistryState() {
	registry = nil
	registryErr = nil
	xmlPath = ""
}

func TestLoadRegistry_NoXMLPath(t *testing.T) {
	oldXML, oldReg, oldErr := xmlPath, registry, registryErr
	defer func() { xmlPath, registry, registryErr = oldXML, oldReg, oldErr }()
	resetRegistryState()

	err := loadRegistry()
	if err == nil {
		t.Fatal("expected error for empty xmlPath, got nil")
	}
	if !errors.Is(err, registryErr) {
		t.Errorf("expected registryErr to be set to the returned error")
	}
}

func TestLoadRegistry_NonexistentFile(t *testing.T) {
	oldXML, oldReg, oldErr := xmlPath, registry, registryErr
	defer func() { xmlPath, registry, registryErr = oldXML, oldReg, oldErr }()
	resetRegistryState()
	xmlPath = "/nonexistent/path/file.xml"

	err := loadRegistry()
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestMustRegistry_NoXMLPath(t *testing.T) {
	oldXML, oldReg, oldErr := xmlPath, registry, registryErr
	defer func() { xmlPath, registry, registryErr = oldXML, oldReg, oldErr }()
	resetRegistryState()

	reg, err := mustRegistry()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if reg != nil {
		t.Errorf("expected nil registry, got %v", reg)
	}
}

func TestMustRegistry_CachedAfterFirstError(t *testing.T) {
	// loadRegistry 失败后 registryErr 被缓存，后续调用直接返回缓存错误
	oldXML, oldReg, oldErr := xmlPath, registry, registryErr
	defer func() { xmlPath, registry, registryErr = oldXML, oldReg, oldErr }()
	resetRegistryState()

	_, err1 := mustRegistry()
	if err1 == nil {
		t.Fatal("expected first error")
	}
	cached := registryErr
	_, err2 := mustRegistry()
	if err2 != cached {
		t.Errorf("expected cached error, got %v vs %v", err2, cached)
	}
}
