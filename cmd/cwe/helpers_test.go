package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/scagogogo/cwe-skills"
	"github.com/spf13/cobra"
)

// newCmdWithBuffer 构造一个输出到 buffer 的 cobra.Command
func newCmdWithBuffer() (*cobra.Command, *bytes.Buffer) {
	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	return cmd, &buf
}

// ==================== helpers.go: printJSON ====================

func TestPrintJSON(t *testing.T) {
	cmd, buf := newCmdWithBuffer()
	data := map[string]string{"key": "value"}
	if err := printJSON(cmd, data); err != nil {
		t.Fatalf("printJSON error: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if got["key"] != "value" {
		t.Errorf("expected key=value, got %v", got)
	}
}

func TestPrintJSON_Nil(t *testing.T) {
	cmd, buf := newCmdWithBuffer()
	if err := printJSON(cmd, nil); err != nil {
		t.Fatalf("printJSON(nil) error: %v", err)
	}
	// json 编码 nil 为 "null"
	if !strings.Contains(buf.String(), "null") {
		t.Errorf("expected null for nil, got %q", buf.String())
	}
}

func TestPrintJSON_Slice(t *testing.T) {
	cmd, buf := newCmdWithBuffer()
	if err := printJSON(cmd, []int{1, 2, 3}); err != nil {
		t.Fatalf("printJSON error: %v", err)
	}
	var got []int
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if len(got) != 3 || got[2] != 3 {
		t.Errorf("expected [1,2,3], got %v", got)
	}
}

// ==================== enums.go: stringifySlice ====================

type testStringer int

func (s testStringer) String() string { return fmt.Sprintf("item-%d", int(s)) }

func TestStringifySlice(t *testing.T) {
	items := []testStringer{1, 2, 3}
	got := stringifySlice(items)
	if len(got) != 3 {
		t.Fatalf("expected 3 items, got %d", len(got))
	}
	if got[0] != "item-1" || got[1] != "item-2" || got[2] != "item-3" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestStringifySlice_Empty(t *testing.T) {
	got := stringifySlice([]testStringer{})
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

// ==================== navigate.go: parseIDArg & boolStr ====================

func TestParseIDArg(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"plain number", "79", 79, false},
		{"CWE-dash format", "CWE-79", 79, false},
		{"CWE space format", "CWE 79", 79, false},
		{"CWE no separator", "CWE79", 79, false},
		{"with spaces", "  79  ", 79, false},
		{"empty", "", 0, true},
		{"invalid string", "abc", 0, true},
		{"zero", "0", 0, true},
		{"negative", "-1", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIDArg(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIDArg(%q) err=%v wantErr=%v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseIDArg(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestBoolStr(t *testing.T) {
	if got := boolStr(true, "yes", "no"); got != "yes" {
		t.Errorf("boolStr(true)=%q want yes", got)
	}
	if got := boolStr(false, "yes", "no"); got != "no" {
		t.Errorf("boolStr(false)=%q want no", got)
	}
}

// ==================== tree.go: treeNodeToMap ====================

func TestTreeNodeToMap(t *testing.T) {
	child := cweskills.NewTreeNode(&cweskills.CWE{ID: 80, Name: "Child"})
	root := cweskills.NewTreeNode(&cweskills.CWE{ID: 79, Name: "XSS"})
	root.AddChild(child)

	m := treeNodeToMap(root)
	if m == nil {
		t.Fatal("expected non-nil map")
	}
	if m["id"] != 79 {
		t.Errorf("expected id=79, got %v", m["id"])
	}
	if m["name"] != "XSS" {
		t.Errorf("expected name=XSS, got %v", m["name"])
	}
	children, ok := m["children"].([]interface{})
	if !ok || len(children) != 1 {
		t.Fatalf("expected 1 child, got %v", m["children"])
	}
	childMap, ok := children[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected child to be map, got %T", children[0])
	}
	if childMap["id"] != 80 {
		t.Errorf("expected child id=80, got %v", childMap["id"])
	}
}

func TestTreeNodeToMap_LeafNode(t *testing.T) {
	// 叶子节点无 children
	root := cweskills.NewTreeNode(&cweskills.CWE{ID: 79, Name: "XSS"})
	m := treeNodeToMap(root)
	if m["is_leaf"] != true {
		t.Errorf("expected is_leaf=true, got %v", m["is_leaf"])
	}
	if _, ok := m["children"]; ok {
		t.Error("leaf node should not have children key")
	}
}

// ==================== tree.go: printTreeNode ====================

func TestPrintTreeNode(t *testing.T) {
	cmd, buf := newCmdWithBuffer()
	root := cweskills.NewTreeNode(&cweskills.CWE{ID: 79, Name: "XSS"})
	root.AddChild(cweskills.NewTreeNode(&cweskills.CWE{ID: 80, Name: "Child"}))

	printTreeNode(cmd, root, 0)
	out := buf.String()
	if !strings.Contains(out, "CWE-79") {
		t.Errorf("expected CWE-79 in output, got %q", out)
	}
	if !strings.Contains(out, "CWE-80") {
		t.Errorf("expected CWE-80 in output, got %q", out)
	}
	if !strings.Contains(out, "XSS") {
		t.Errorf("expected XSS in output, got %q", out)
	}
}

// ==================== wellknown.go: printIDList & sortedOWASPKeys ====================

func TestPrintIDList_Text(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	if err := printIDList(cmd, "Test List", []int{79, 89}); err != nil {
		t.Fatalf("printIDList error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Test List") {
		t.Errorf("expected title, got %q", out)
	}
	if !strings.Contains(out, "CWE-79") || !strings.Contains(out, "CWE-89") {
		t.Errorf("expected CWE-79 and CWE-89, got %q", out)
	}
}

func TestPrintIDList_JSON(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "json"

	cmd, buf := newCmdWithBuffer()
	if err := printIDList(cmd, "Test List", []int{79}); err != nil {
		t.Fatalf("printIDList error: %v", err)
	}
	var entries []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("expected JSON array, got %q: %v", buf.String(), err)
	}
	if len(entries) != 1 || entries[0]["format"] != "CWE-79" {
		t.Errorf("unexpected entries: %v", entries)
	}
}

func TestPrintIDList_Empty(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	if err := printIDList(cmd, "Empty", []int{}); err != nil {
		t.Fatalf("printIDList error: %v", err)
	}
	if !strings.Contains(buf.String(), "Empty") {
		t.Errorf("expected title, got %q", buf.String())
	}
}

func TestSortedOWASPKeys(t *testing.T) {
	keys := sortedOWASPKeys()
	if len(keys) == 0 {
		t.Fatal("expected non-empty OWASP keys")
	}
	// 验证已排序
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("keys not sorted: %q > %q", keys[i-1], keys[i])
		}
	}
}

// ==================== relations.go: printRelationships ====================

func TestPrintRelationships_Text(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	rels := []cweskills.Relationship{
		{Nature: cweskills.RelationshipChildOf, CWEID: 74, ViewID: 1000},
		{Nature: cweskills.RelationshipPeerOf, CWEID: 89},
	}
	if err := printRelationships(cmd, "父级", 79, rels); err != nil {
		t.Fatalf("printRelationships error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CWE-79") {
		t.Errorf("expected CWE-79, got %q", out)
	}
	if !strings.Contains(out, "CWE-74") {
		t.Errorf("expected CWE-74, got %q", out)
	}
	if !strings.Contains(out, "View: 1000") {
		t.Errorf("expected View: 1000, got %q", out)
	}
}

func TestPrintRelationships_JSON(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "json"

	cmd, buf := newCmdWithBuffer()
	rels := []cweskills.Relationship{
		{Nature: cweskills.RelationshipChildOf, CWEID: 74},
	}
	if err := printRelationships(cmd, "父级", 79, rels); err != nil {
		t.Fatalf("printRelationships error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("expected JSON object: %v\noutput: %s", err, buf.String())
	}
	if got["cwe_id"] != "CWE-79" {
		t.Errorf("expected cwe_id CWE-79, got %v", got["cwe_id"])
	}
	if got["count"].(float64) != 1 {
		t.Errorf("expected count 1, got %v", got["count"])
	}
}

func TestPrintRelationships_Empty(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	if err := printRelationships(cmd, "父级", 79, nil); err != nil {
		t.Fatalf("printRelationships error: %v", err)
	}
	if !strings.Contains(buf.String(), "0 项") {
		t.Errorf("expected 0 items, got %q", buf.String())
	}
}

// ==================== search.go: printGroupedResults ====================

func TestPrintGroupedResults_Text(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "text"
	searchGroupBy = "abstraction"

	cmd, buf := newCmdWithBuffer()
	results := []*cweskills.CWE{
		{ID: 79, Name: "XSS", Abstraction: cweskills.AbstractionVariant},
		{ID: 89, Name: "SQLi", Abstraction: cweskills.AbstractionBase},
	}
	if err := printGroupedResults(cmd, results); err != nil {
		t.Fatalf("printGroupedResults error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CWE-79") || !strings.Contains(out, "XSS") {
		t.Errorf("expected CWE-79/XSS, got %q", out)
	}
}

func TestPrintGroupedResults_JSON(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "json"
	searchGroupBy = "status"

	cmd, buf := newCmdWithBuffer()
	results := []*cweskills.CWE{{ID: 79, Name: "XSS", Status: cweskills.StatusStable}}
	if err := printGroupedResults(cmd, results); err != nil {
		t.Fatalf("printGroupedResults error: %v", err)
	}
	var got interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("expected valid JSON: %v\noutput: %s", err, buf.String())
	}
}

func TestPrintGroupedResults_JSONAbstraction(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "json"
	searchGroupBy = "abstraction"

	cmd, buf := newCmdWithBuffer()
	results := []*cweskills.CWE{{ID: 79, Name: "XSS", Abstraction: cweskills.AbstractionVariant}}
	if err := printGroupedResults(cmd, results); err != nil {
		t.Fatalf("printGroupedResults error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("expected JSON: %v", err)
	}
	if _, ok := got["Variant"]; !ok {
		t.Errorf("expected Variant group, got %v", got)
	}
}

func TestPrintGroupedResults_JSONLikelihood(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "json"
	searchGroupBy = "likelihood"

	cmd, buf := newCmdWithBuffer()
	results := []*cweskills.CWE{{ID: 79, Name: "XSS", LikelihoodOfExploit: cweskills.LikelihoodHigh}}
	if err := printGroupedResults(cmd, results); err != nil {
		t.Fatalf("printGroupedResults error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("expected JSON: %v", err)
	}
}

func TestPrintGroupedResults_JSONInvalidGroup(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "json"
	searchGroupBy = "invalid"

	cmd, _ := newCmdWithBuffer()
	err := printGroupedResults(cmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid group in JSON, got nil")
	}
}

func TestPrintGroupedResults_Likelihood(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "text"
	searchGroupBy = "likelihood"

	cmd, buf := newCmdWithBuffer()
	results := []*cweskills.CWE{
		{ID: 79, Name: "XSS", LikelihoodOfExploit: cweskills.LikelihoodHigh},
	}
	if err := printGroupedResults(cmd, results); err != nil {
		t.Fatalf("printGroupedResults error: %v", err)
	}
	if !strings.Contains(buf.String(), "CWE-79") {
		t.Errorf("expected CWE-79, got %q", buf.String())
	}
}

func TestPrintGroupedResults_Status(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "text"
	searchGroupBy = "status"

	cmd, buf := newCmdWithBuffer()
	results := []*cweskills.CWE{
		{ID: 79, Name: "XSS", Status: cweskills.StatusStable},
	}
	if err := printGroupedResults(cmd, results); err != nil {
		t.Fatalf("printGroupedResults error: %v", err)
	}
	if !strings.Contains(buf.String(), "CWE-79") {
		t.Errorf("expected CWE-79, got %q", buf.String())
	}
}

func TestPrintGroupedResults_InvalidGroup(t *testing.T) {
	oldFmt, oldGroup := outputFormat, searchGroupBy
	defer func() { outputFormat, searchGroupBy = oldFmt, oldGroup }()
	outputFormat = "text"
	searchGroupBy = "invalid"

	cmd, _ := newCmdWithBuffer()
	err := printGroupedResults(cmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid group, got nil")
	}
}

// ==================== registry.go: printCWEDetail, printIDResults, writeFile ====================

func TestPrintCWEDetail(t *testing.T) {
	cmd, buf := newCmdWithBuffer()
	cwe := &cweskills.CWE{
		ID:                  79,
		Name:                "XSS",
		Abstraction:         cweskills.AbstractionVariant,
		Status:              cweskills.StatusStable,
		Structure:           cweskills.StructureSimple,
		Description:         "desc",
		LikelihoodOfExploit: cweskills.LikelihoodHigh,
		Relationships: []cweskills.Relationship{
			{Nature: cweskills.RelationshipChildOf, CWEID: 74},
		},
		CommonConsequences: []cweskills.Consequence{
			{Note: "test consequence"},
		},
	}
	printCWEDetail(cmd, cwe)
	out := buf.String()
	if !strings.Contains(out, "CWE-79") {
		t.Errorf("expected CWE-79, got %q", out)
	}
	if !strings.Contains(out, "XSS") {
		t.Errorf("expected XSS, got %q", out)
	}
	if !strings.Contains(out, "desc") {
		t.Errorf("expected description, got %q", out)
	}
	if !strings.Contains(out, "关系") {
		t.Errorf("expected 关系, got %q", out)
	}
}

func TestPrintCWEDetail_Minimal(t *testing.T) {
	// 只 ID+Name，触发各 if 分支跳过
	cmd, buf := newCmdWithBuffer()
	cwe := &cweskills.CWE{ID: 79, Name: "XSS"}
	printCWEDetail(cmd, cwe)
	out := buf.String()
	if !strings.Contains(out, "CWE-79") {
		t.Errorf("expected CWE-79, got %q", out)
	}
	if strings.Contains(out, "结构") {
		t.Errorf("minimal cwe should not print 结构, got %q", out)
	}
}

func TestPrintIDResults_Text(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	r := cweskills.NewRegistry()
	r.Register(&cweskills.CWE{ID: 79, Name: "XSS"})
	r.Register(&cweskills.CWE{ID: 80, Name: "Child"})

	if err := printIDResults(cmd, "父级", 79, []int{79, 80, 999}, r); err != nil {
		t.Fatalf("printIDResults error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CWE-79") {
		t.Errorf("expected CWE-79, got %q", out)
	}
	if !strings.Contains(out, "XSS") {
		t.Errorf("expected XSS name, got %q", out)
	}
	// 999 不在 registry，应只打印 CWE-999 无名称
	if !strings.Contains(out, "CWE-999") {
		t.Errorf("expected CWE-999, got %q", out)
	}
}

func TestPrintIDResults_JSON(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "json"

	cmd, buf := newCmdWithBuffer()
	r := cweskills.NewRegistry()
	r.Register(&cweskills.CWE{ID: 79, Name: "XSS"})

	if err := printIDResults(cmd, "父级", 79, []int{79}, r); err != nil {
		t.Fatalf("printIDResults error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("expected JSON: %v\noutput: %s", err, buf.String())
	}
	if got["cwe_id"] != "CWE-79" {
		t.Errorf("expected cwe_id CWE-79, got %v", got["cwe_id"])
	}
}

func TestPrintIDResults_Empty(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	r := cweskills.NewRegistry()
	if err := printIDResults(cmd, "父级", 79, nil, r); err != nil {
		t.Fatalf("printIDResults error: %v", err)
	}
	if !strings.Contains(buf.String(), "0 项") {
		t.Errorf("expected 0 items, got %q", buf.String())
	}
}

func TestWriteFile_AlwaysReturnsError(t *testing.T) {
	// writeFile 是 stub 实现，总是返回错误
	err := writeFile("/tmp/out.txt", []byte("hello"))
	if err == nil {
		t.Fatal("expected error from writeFile stub, got nil")
	}
	if !strings.Contains(err.Error(), "/tmp/out.txt") {
		t.Errorf("expected path in error, got %v", err)
	}
	if !strings.Contains(err.Error(), "5 bytes") {
		t.Errorf("expected byte count in error, got %v", err)
	}
}

// ==================== navigate.go: printNavResults ====================

func TestPrintNavResults_Text(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	cwes := []*cweskills.CWE{
		{ID: 74, Name: "Parent"},
		{ID: 707, Name: "Ancestor"},
	}
	if err := printNavResults(cmd, "父级", 79, cwes, nil); err != nil {
		t.Fatalf("printNavResults error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CWE-79") {
		t.Errorf("expected CWE-79, got %q", out)
	}
	if !strings.Contains(out, "Parent") {
		t.Errorf("expected Parent, got %q", out)
	}
}

func TestPrintNavResults_JSON(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "json"

	cmd, buf := newCmdWithBuffer()
	cwes := []*cweskills.CWE{{ID: 74, Name: "Parent"}}
	if err := printNavResults(cmd, "父级", 79, cwes, nil); err != nil {
		t.Fatalf("printNavResults error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("expected JSON: %v\noutput: %s", err, buf.String())
	}
	if got["cwe_id"] != "CWE-79" {
		t.Errorf("expected cwe_id CWE-79, got %v", got["cwe_id"])
	}
	if got["count"].(float64) != 1 {
		t.Errorf("expected count 1, got %v", got["count"])
	}
}

func TestPrintNavResults_Empty(t *testing.T) {
	old := outputFormat
	defer func() { outputFormat = old }()
	outputFormat = "text"

	cmd, buf := newCmdWithBuffer()
	if err := printNavResults(cmd, "父级", 79, nil, nil); err != nil {
		t.Fatalf("printNavResults error: %v", err)
	}
	if !strings.Contains(buf.String(), "0 项") {
		t.Errorf("expected 0 items, got %q", buf.String())
	}
}
