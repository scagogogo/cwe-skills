package cweskills

// 本文件专为补全 SDK 层单元测试覆盖率而创建，
// 覆盖现有测试未触达的分支（错误路径、边界、fallback 解析等）。
//
// 注：以下分支因实现使用 bytes.Buffer/csv.Writer（永不返回错误）而不可达，
// 接受不覆盖：serializer.go:82 (xml.MarshalIndent 对 safeCWE 永不失败)、
// :166/:181/:187 (csv.Writer.Write/Flush 对 bytes.Buffer 永不失败)、
// :220 (bytes.Reader.Read 永不返回非 EOF 错误)。
//
// api_client_relations.go 的 getRelations fallback 分支（行 124-129）
// 实际可达：Relationship 比 fallback struct 多 Ordinal(string)/ChainID(int)
// 字段，传入 ordinal 为数字值时 []Relationship 解析失败（类型不匹配），
// 但 fallback struct 无 ordinal 字段会忽略该字段从而解析成功。
// 见 TestGetRelations_FallbackObjectArray。

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// ==================== api_client_cwe.go: 单 CWE/Category/View 解析失败的 fallback ====================
// 行 49-51: GetWeakness 中 []CWE 失败后，单 CWE 也失败 → NewParseError
// 行 93-95: GetCategory 中 []Category 失败后，单 Category 也失败
// 行 135-137: GetView 中 []View 失败后，单 View 也失败
// 行 185-187: GetCWEs 中 map 解析失败

// 触发单对象 fallback 成功路径（既非数组也非单对象时进入 fallback 并成功）
func TestGetWeakness_SingleObjectFallbackSuccess(t *testing.T) {
	// Data 是单个对象（非数组）→ []CWE 解析失败 → 单 CWE 解析成功
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 用 RawMessage 直接写入原始 JSON，确保 Data 字段拿到的是对象
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{
			Data: json.RawMessage(`{"id":79,"name":"XSS","cwe_type":"weakness"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	cwe, err := client.GetWeakness(context.Background(), 79)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cwe.ID != 79 || cwe.Name != "XSS" {
		t.Errorf("unexpected cwe: %+v", cwe)
	}
	if cwe.CWEType != "weakness" {
		t.Errorf("expected CWEType weakness, got %q", cwe.CWEType)
	}
}

// 触发 []CWE 解析失败 + 单 CWE 也失败 → NewParseError（行 49-51）
func TestGetWeakness_BothParseFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{
			// 既不是数组也不是有效 CWE 对象（纯数字）
			Data: json.RawMessage(`12345`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	_, err := client.GetWeakness(context.Background(), 79)
	if err == nil {
		t.Fatal("expected ParseError, got nil")
	}
}

// GetCategory 单对象 fallback 成功（行 92-96）
func TestGetCategory_SingleObjectFallbackSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{
			Data: json.RawMessage(`{"id":1,"name":"Category1"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	cat, err := client.GetCategory(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.ID != 1 || cat.Name != "Category1" {
		t.Errorf("unexpected category: %+v", cat)
	}
}

// GetCategory 两层都失败（行 93-95）
func TestGetCategory_BothParseFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{Data: json.RawMessage(`"not object or array"`)};
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	_, err := client.GetCategory(context.Background(), 1)
	if err == nil {
		t.Fatal("expected ParseError, got nil")
	}
}

// GetView 单对象 fallback 成功（行 135-138）
func TestGetView_SingleObjectFallbackSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{
			Data: json.RawMessage(`{"id":1000,"name":"View1000"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	view, err := client.GetView(context.Background(), 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.ID != 1000 || view.Name != "View1000" {
		t.Errorf("unexpected view: %+v", view)
	}
}

// GetView 两层都失败（行 135-137）
func TestGetView_BothParseFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{Data: json.RawMessage(`999`)};
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	_, err := client.GetView(context.Background(), 1000)
	if err == nil {
		t.Fatal("expected ParseError, got nil")
	}
}

// GetCWEs map 解析失败（行 185-187）
func TestGetCWEs_ParseError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{Data: json.RawMessage(`"not a map"`)};
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	_, err := client.GetCWEs(context.Background(), []int{79})
	if err == nil {
		t.Fatal("expected ParseError, got nil")
	}
}

// GetCWEs 空 ID 返回空 map（line 162-163 已覆盖，补充 nil result 填充 CWEType）
func TestGetCWEs_NilCWEInResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{Data: json.RawMessage(`{"79":null}`)};
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	result, err := client.GetCWEs(context.Background(), []int{79})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// nil CWE 的 CWEType 填充分支（line 189-193 的 if cwe != nil 跳过）
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}

// ==================== http_client.go ====================

// doRequest rateLimiter.Wait 失败（行 279-281）
func TestDoRequest_RateLimiterError(t *testing.T) {
	// burst=0 使初始 tokens=0，首次 Wait 即需等待；ctx 已取消则返回 ctx.Err()
	c := NewHTTPClient("http://example.com", WithHTTPRateLimiter(1.0, 0))
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.doRequest(ctx, http.MethodGet, "/path", nil)
	if err == nil {
		t.Fatal("expected ctx error from rateLimiter, got nil")
	}
}

// doRequest 5xx 重试耗尽后返回 lastErr（行 317-319 的 io.ReadAll 在 httptest 下难失败，
// 但行 336 lastErr 返回可达）
func TestDoRequest_5xxRetryExhausts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewHTTPClient(srv.URL, WithRetry(1, 1*time.Millisecond))
	defer c.Close()

	_, err := c.doRequest(context.Background(), http.MethodGet, "/path", nil)
	if err == nil {
		t.Fatal("expected error after retries, got nil")
	}
}

// PostForm result 解析失败（行 262-264）
func TestPostForm_ResultParseError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := NewHTTPClient(srv.URL)
	defer c.Close()

	var result struct{ X int }
	err := c.PostForm(context.Background(), "/path", url.Values{"k": {"v"}}, &result)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

// PostForm rateLimiter.Wait 失败（行 216-219）
func TestPostForm_RateLimiterError(t *testing.T) {
	// burst=0 使首次 Wait 即需等待；ctx 取消则返回 ctx.Err()
	c := NewHTTPClient("http://example.com", WithHTTPRateLimiter(1.0, 0))
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := c.PostForm(ctx, "/path", url.Values{}, nil)
	if err == nil {
		t.Fatal("expected ctx error, got nil")
	}
}

// PostForm 重试成功（行 252-255 的 5xx 重试分支）
func TestPostForm_5xxRetryThenSuccess(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := NewHTTPClient(srv.URL, WithRetry(2, 1*time.Millisecond))
	defer c.Close()

	var result map[string]bool
	err := c.PostForm(context.Background(), "/path", url.Values{"k": {"v"}}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["ok"] {
		t.Errorf("expected ok=true, got %v", result)
	}
}

// PostForm 5xx 无重试返回 APIError（行 257-259）
func TestPostForm_5xxNoRetryReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewHTTPClient(srv.URL)
	defer c.Close()

	err := c.PostForm(context.Background(), "/path", url.Values{}, nil)
	if err == nil {
		t.Fatal("expected APIError, got nil")
	}
}

// Post result 解析失败（行 194-196）
func TestPost_ResultParseError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := NewHTTPClient(srv.URL)
	defer c.Close()

	var result struct{ X int }
	err := c.Post(context.Background(), "/path", map[string]string{"k": "v"}, &result)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

// Post body marshal 失败（行 183-185）
func TestPost_MarshalError(t *testing.T) {
	c := NewHTTPClient("http://example.com")
	defer c.Close()

	var result struct{}
	err := c.Post(context.Background(), "/path", make(chan int), &result)
	if err == nil {
		t.Fatal("expected marshal error, got nil")
	}
}

// ==================== xml_parser.go: ParseFile Open 成功路径（行 109 defer 注册） ====================

func TestParseFile_OpenSuccess(t *testing.T) {
	// 创建一个有效的 XML 文件，使 os.Open 成功，触发 defer file.Close() 注册（行 109）
	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "test.xml")
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<Weakness_Catalog Name="Test" Version="7.0" Date="2026-01-01">
  <Weaknesses>
    <Weakness ID="79" Name="XSS" Abstraction="Variant" Structure="Simple" Status="Stable">
      <Description>Test desc</Description>
    </Weakness>
  </Weaknesses>
</Weakness_Catalog>`
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0644); err != nil {
		t.Fatalf("write xml: %v", err)
	}

	p := NewXMLParser()
	registry, err := p.ParseFile(xmlPath)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
	if registry.Size() != 1 {
		t.Errorf("expected 1 weakness, got %d", registry.Size())
	}
}

// byteReader EOF（行 141-143）
func TestByteReader_EOF(t *testing.T) {
	r := newByteReader([]byte{})
	n, err := r.Read(make([]byte, 10))
	if err == nil {
		// 第一次读空数据，offset>=len，应返回 io.EOF
		t.Errorf("expected error, got nil; n=%d", n)
	}
}

// ParseFile 空路径
func TestParseFile_EmptyPath(t *testing.T) {
	p := NewXMLParser()
	_, err := p.ParseFile("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

// ==================== serializer.go 可达分支 ====================

// MarshalXML nil 返回错误（行 74-76）
func TestMarshalXML_NilReturnsError(t *testing.T) {
	_, err := MarshalXML(nil)
	if err == nil {
		t.Fatal("expected error for nil CWE, got nil")
	}
}

// MarshalJSON nil 返回错误
func TestMarshalJSON_NilReturnsError(t *testing.T) {
	_, err := MarshalJSON(nil)
	if err == nil {
		t.Fatal("expected error for nil CWE, got nil")
	}
}

// MarshalCSV nil 返回空（行 158-160）
func TestMarshalCSV_NilReturnsEmpty(t *testing.T) {
	out, err := MarshalCSV(nil)
	if err != nil {
		t.Fatalf("MarshalCSV(nil) error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty for nil, got %q", out)
	}
}

// MarshalCSV 单行（覆盖循环体）
func TestMarshalCSV_SingleRow(t *testing.T) {
	out, err := MarshalCSV([]*CWE{{ID: 79, Name: "XSS", Description: "desc"}})
	if err != nil {
		t.Fatalf("MarshalCSV error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty CSV")
	}
}

// UnmarshalCSV 无效 ID 行跳过（覆盖 parseCSVInt 失败 continue）
func TestUnmarshalCSV_InvalidIntSkipped(t *testing.T) {
	data := []byte("ID,Name\nabc,BadRow\n79,XSS\n")
	cwes, err := UnmarshalCSV(data)
	if err != nil {
		t.Fatalf("UnmarshalCSV error: %v", err)
	}
	if len(cwes) != 1 || cwes[0].ID != 79 {
		t.Errorf("expected 1 valid cwe (79), got %+v", cwes)
	}
}

// UnmarshalCSV 空数据返回错误（行 196-198）
func TestUnmarshalCSV_EmptyData(t *testing.T) {
	_, err := UnmarshalCSV([]byte{})
	if err == nil {
		t.Fatal("expected error for empty data, got nil")
	}
}

// UnmarshalCSV 列数不足行跳过（行 224-226）
func TestUnmarshalCSV_RowWithOneColumnSkipped(t *testing.T) {
	data := []byte("ID,Name\n79\n80,XSS\n")
	cwes, err := UnmarshalCSV(data)
	if err != nil {
		t.Fatalf("UnmarshalCSV error: %v", err)
	}
	if len(cwes) != 1 || cwes[0].ID != 80 {
		t.Errorf("expected 1 cwe (80), got %+v", cwes)
	}
}

// ==================== tree.go: buildTreeNode visited 跳过 + Get 失败 ====================

func TestBuildTreeNode_SkipsVisited(t *testing.T) {
	// 构造一个有循环引用的树结构，使 childID 已在 visited 中
	r := NewRegistry()
	r.Register(&CWE{ID: 1, Name: "Root"})
	r.Register(&CWE{ID: 2, Name: "Child", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 1},
	}})
	r.BuildIndexes()

	root := NewTreeNode(&CWE{ID: 1, Name: "Root"})
	// 预先把 childID 2 标记为已访问
	visited := map[int]bool{1: true, 2: true}
	buildTreeNode(r, root, visited)
	// 不 panic 即成功，且不会重复添加 child 2
	if len(root.Children) != 0 {
		t.Errorf("expected 0 children (all visited), got %d", len(root.Children))
	}
}

func TestBuildTreeNode_GetFailSkips(t *testing.T) {
	// childID 指向不存在的 CWE，Get 失败跳过
	r := NewRegistry()
	r.Register(&CWE{ID: 1, Name: "Root"})
	// 注册一个 ChildOf 指向不存在的 CWE 999
	r.Register(&CWE{ID: 2, Name: "Child", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 1},
		{Nature: RelationshipParentOf, CWEID: 999}, // 999 不存在
	}})
	r.BuildIndexes()

	root := NewTreeNode(&CWE{ID: 2, Name: "Child"})
	visited := map[int]bool{2: true}
	buildTreeNode(r, root, visited)
	// CWE-2 的 child 是 999（ParentOf），Get(999) 失败跳过
	// 不 panic 即成功
}

// ==================== navigator.go ====================

// IsRelated 反向关系（行 269-277）
func TestIsRelated_ReverseRelationship(t *testing.T) {
	r := NewRegistry()
	// CWE-89 指向 79，检查 79->89 是否通过反向关系发现
	r.Register(&CWE{ID: 89, Name: "SQLi", Relationships: []Relationship{
		{Nature: RelationshipPeerOf, CWEID: 79},
	}})
	r.Register(&CWE{ID: 79, Name: "XSS", Relationships: []Relationship{}})

	nav := NewNavigator(r)
	// IsRelated(79, 89)：先查 79 的关系（无），再查 89 的关系（指向 79）→ true
	if !nav.IsRelated(79, 89) {
		t.Error("expected IsRelated(79,89)=true via reverse")
	}
}

// IsRelated nil registry（行 253-255）
func TestIsRelated_NilRegistry(t *testing.T) {
	nav := &Navigator{registry: nil}
	if nav.IsRelated(1, 2) {
		t.Error("expected false for nil registry")
	}
}

// IsRelated a 不存在（行 258-260）
func TestIsRelated_ANotFound(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})
	nav := NewNavigator(r)
	if nav.IsRelated(999, 79) {
		t.Error("expected false when a not found")
	}
}

// IsRelated b 反向不存在（行 269-272）
func TestIsRelated_BNotFoundInReverse(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 74},
	}})
	nav := NewNavigator(r)
	// b=999 不存在，反向查 999 失败
	if nav.IsRelated(79, 999) {
		t.Error("expected false when b not found")
	}
}

// findByRelationship Get 失败（行 324-327）
func TestFindByRelationship_NotFound(t *testing.T) {
	r := NewRegistry()
	nav := NewNavigator(r)
	result := nav.findByRelationship(999, RelationshipChildOf)
	if result != nil {
		t.Errorf("expected nil for non-existent id, got %v", result)
	}
}

// ==================== registry.go BuildIndexes ====================

// HasMember 在 weakness 中（行 349-351）
func TestBuildIndexes_HasMemberInWeakness(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})
	// weakness 用 RelationshipHasMember 触发 line 349-351
	r.Register(&CWE{ID: 1, Name: "Cat", Relationships: []Relationship{
		{Nature: RelationshipHasMember, CWEID: 79},
	}})
	r.BuildIndexes()
	// HasMember: memberIndex[1]=[79], memberOfIndex[79]=[1]
	memberOf := r.GetMemberOfIDs(79)
	if !containsInt(memberOf, 1) {
		t.Errorf("expected memberOf 1 for CWE-79, got %v", memberOf)
	}
}

// MemberOf 在 weakness 中（行 346-348）
func TestBuildIndexes_MemberOfInWeakness(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", Relationships: []Relationship{
		{Nature: RelationshipMemberOf, CWEID: 1},
	}})
	r.BuildIndexes()
	memberOf := r.GetMemberOfIDs(79)
	if !containsInt(memberOf, 1) {
		t.Errorf("expected memberOf 1, got %v", memberOf)
	}
}

// ParentOf 在 weakness 中（行 341-343）
func TestBuildIndexes_ParentOfRelation(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 1, Name: "Root", Relationships: []Relationship{
		{Nature: RelationshipParentOf, CWEID: 2},
	}})
	r.Register(&CWE{ID: 2, Name: "Child"})
	r.BuildIndexes()
	children := r.GetChildIDs(1)
	if !containsInt(children, 2) {
		t.Errorf("expected child 2, got %v", children)
	}
}

// Category MemberOf（行 363-365）
func TestBuildIndexes_CategoryMemberOf(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1", Relationships: []Relationship{
		{Nature: RelationshipMemberOf, CWEID: 79},
	}})
	r.BuildIndexes()
	// 触发分支即可
	memberOf := r.GetMemberOfIDs(1)
	_ = memberOf
}

// 重复 ID 去重（行 378-383 dedupIntMap）
func TestBuildIndexes_DedupDuplicates(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 2, Name: "Child", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 1},
		{Nature: RelationshipChildOf, CWEID: 1},
	}})
	r.Register(&CWE{ID: 1, Name: "Root"})
	r.BuildIndexes()

	parents := r.GetParentIDs(2)
	count := 0
	for _, p := range parents {
		if p == 1 {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected dedup to 1, got %d", count)
	}
}

// ==================== data_fetcher.go 递归错误分支 ====================

// fetchAncestorsRecursive 递归 GetWeakness 错误（行 225-227）
func TestFetchAncestorsRecursive_GetWeaknessError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)
	err := fetcher.fetchAncestorsRecursive(context.Background(), 79, 0)
	if err == nil {
		t.Fatal("expected error from GetWeakness failure, got nil")
	}
}

// fetchDescendantsRecursive 递归 GetWeakness 错误（行 242-244）
func TestFetchDescendantsRecursive_GetWeaknessError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)
	err := fetcher.fetchDescendantsRecursive(context.Background(), 79, 0)
	if err == nil {
		t.Fatal("expected error from GetWeakness failure, got nil")
	}
}

// fetchDescendantsRecursive 递归 GetChildren 错误（行 264-266）
func TestFetchDescendantsRecursive_GetChildrenError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"d","cwe_type":"weakness"}]`)}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/children":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)
	err := fetcher.fetchDescendantsRecursive(context.Background(), 79, 0)
	if err == nil {
		t.Fatal("expected error from GetChildren failure, got nil")
	}
}

// ==================== http_client.go: invalid URL（行 233, 300）====================

// PostForm invalid URL 触发 http.NewRequestWithContext 失败（行 232-235）
func TestPostForm_InvalidURLError(t *testing.T) {
	c := NewHTTPClient("http://example.com")
	defer c.Close()
	// 空字节 path 使 NewRequestWithContext 失败
	err := c.PostForm(context.Background(), "\x00bad", url.Values{"k": {"v"}}, nil)
	if err == nil {
		t.Fatal("expected NewRequest error, got nil")
	}
}

// doRequest invalid URL 触发 http.NewRequestWithContext 失败（行 299-302）
func TestDoRequest_InvalidURLError(t *testing.T) {
	c := NewHTTPClient("http://example.com")
	defer c.Close()
	_, err := c.doRequest(context.Background(), http.MethodGet, "\x00bad", nil)
	if err == nil {
		t.Fatal("expected NewRequest error, got nil")
	}
}

// ==================== http_client.go: io.ReadAll 失败（行 247, 317）====================

// errorReader 在 Read 时返回错误
type errorReader struct{}

func (errorReader) Read(p []byte) (int, error) {
	return 0, ioErrRead
}

var ioErrRead = errReadSentinel{}

type errReadSentinel struct{}

func (errReadSentinel) Error() string { return "simulated read error" }

// TestDoRequest_ReadAllError 用自定义 RoundTripper 返回会读失败的 Body
func TestDoRequest_ReadAllError(t *testing.T) {
	c := NewHTTPClient("http://example.com")
	defer c.Close()
	// 注入自定义 Transport，使 Response.Body 是 errorReader
	c.SetHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(errorReader{}),
				Header:     make(http.Header),
			}, nil
		}),
	})

	_, err := c.doRequest(context.Background(), http.MethodGet, "/path", nil)
	if err == nil {
		t.Fatal("expected ReadAll error, got nil")
	}
}

// TestPostForm_ReadAllError 用自定义 RoundTripper 触发 PostForm 的 io.ReadAll 失败（行 247-249）
func TestPostForm_ReadAllError(t *testing.T) {
	c := NewHTTPClient("http://example.com")
	defer c.Close()
	c.SetHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(errorReader{}),
				Header:     make(http.Header),
			}, nil
		}),
	})

	// 带重试，使 ReadAll 失败后 continue，最终返回 lastErr
	c.SetMaxRetries(1)
	c.SetRetryDelay(1 * time.Millisecond)

	err := c.PostForm(context.Background(), "/path", url.Values{"k": {"v"}}, nil)
	if err == nil {
		t.Fatal("expected ReadAll error after retries, got nil")
	}
}

// roundTripperFunc 适配器，将函数实现为 http.RoundTripper
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// ==================== data_fetcher.go: 递归中途失败（行 225, 264）====================

// fetchAncestorsRecursive 递归到第二层失败（行 225-227）
func TestFetchAncestorsRecursive_RecursiveError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/cwe/weakness/79":
			// 第一层 GetWeakness(79) 成功
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"d","cwe_type":"weakness"}]`)}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/parents":
			// GetParents(79) 成功，返回 parent 74
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74}]`)}
			json.NewEncoder(w).Encode(resp)
		default:
			// 递归获取 74 的 GetWeakness 失败
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)
	err := fetcher.fetchAncestorsRecursive(context.Background(), 79, 0)
	if err == nil {
		t.Fatal("expected recursive error, got nil")
	}
}

// fetchDescendantsRecursive 递归到第二层失败（行 264-266）
func TestFetchDescendantsRecursive_RecursiveError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"d","cwe_type":"weakness"}]`)}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/children":
			// GetChildren(79) 成功，返回 child 80
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{Data: json.RawMessage(`[{"nature":"ParentOf","cweId":80}]`)}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/80/children":
			// 递归获取 80 的 children 失败
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
	defer client.Close()

	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)
	err := fetcher.fetchDescendantsRecursive(context.Background(), 79, 0)
	if err == nil {
		t.Fatal("expected recursive error, got nil")
	}
}

// TestDoRequest_ClientDoError 触发 c.client.Do(req) 返回 error（行 240-242）
func TestDoRequest_ClientDoError(t *testing.T) {
	c := NewHTTPClient("http://example.com")
	defer c.Close()
	// 用一个总是返回 transport error 的 RoundTripper
	c.SetHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, ioErrRead
		}),
	})
	c.SetMaxRetries(1)
	c.SetRetryDelay(1 * time.Millisecond)

	_, err := c.doRequest(context.Background(), http.MethodGet, "/path", nil)
	if err == nil {
		t.Fatal("expected client.Do error, got nil")
	}
}

// TestPostForm_ClientDoError 触发 PostForm 的 c.client.Do(req) 失败（行 240-242 via PostForm）
func TestPostForm_ClientDoError(t *testing.T) {
	c := NewHTTPClient("http://example.com")
	defer c.Close()
	c.SetHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, ioErrRead
		}),
	})
	c.SetMaxRetries(1)
	c.SetRetryDelay(1 * time.Millisecond)

	err := c.PostForm(context.Background(), "/path", url.Values{"k": {"v"}}, nil)
	if err == nil {
		t.Fatal("expected client.Do error, got nil")
	}
}

// ==================== api_client_relations.go: getRelations 分支 ====================

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
//
// Relationship.Ordinal 为 string 类型；body 中 ordinal 传数字值 123，
// 导致 []Relationship 解析失败（类型不匹配）→ 进入 fallback。
// fallback struct 无 ordinal 字段会忽略该字段，从而解析成功。
// 第二条 nature 传无效值 "BogusNature"，使 ParseRelationshipNature 返回错误，
// 覆盖行 127 的 nature = RelationshipNature(r.Nature) fallback 赋值分支。
func TestGetRelations_FallbackObjectArray(t *testing.T) {
	body := `[{"nature":"ChildOf","cweId":100,"ordinal":123},{"nature":"BogusNature","cweId":200,"ordinal":456}]`
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
	// Go encoding/json 在 Unmarshal 报错时仍会部分填充 slice，第一段
	// []Relationship 虽失败（ordinal 数字无法填入 string 字段）却可能已填入
	// 若干条；随后 fallback 循环再 append 2 条。因此只验证 fallback 追加的
	// 尾部 2 条元素，不把部分填充的具体数量固化为契约，避免脆弱断言。
	if len(rel) < 2 {
		t.Fatalf("expected at least 2 fallback relations, got %d", len(rel))
	}
	fb0, fb1 := rel[len(rel)-2], rel[len(rel)-1]
	if fb0.CWEID != 100 || fb1.CWEID != 200 {
		t.Errorf("unexpected fallback cwe ids: %+v", rel)
	}
	// 第一条 nature 有效 → 走 ParseRelationshipNature 成功分支
	if fb0.Nature != RelationshipChildOf {
		t.Errorf("unexpected nature of fallback[0]: %v", fb0.Nature)
	}
	// 第二条 nature 无效 → ParseRelationshipNature 失败 → 行 127 原样赋值
	if fb1.Nature != RelationshipNature("BogusNature") {
		t.Errorf("unexpected nature of fallback[1]: %v", fb1.Nature)
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
