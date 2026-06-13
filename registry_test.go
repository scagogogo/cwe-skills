package cwe

import (
	"encoding/json"
	"testing"
)

// ========== NewRegistry ==========

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if r.Size() != 0 {
		t.Errorf("expected empty registry, got size %d", r.Size())
	}
	if r.CategoryCount() != 0 {
		t.Errorf("expected 0 categories, got %d", r.CategoryCount())
	}
	if r.ViewCount() != 0 {
		t.Errorf("expected 0 views, got %d", r.ViewCount())
	}
	if r.CompoundElementCount() != 0 {
		t.Errorf("expected 0 compound elements, got %d", r.CompoundElementCount())
	}
}

// ========== Register ==========

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	// Success
	err := r.Register(&CWE{ID: 79, Name: "XSS"})
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if r.Size() != 1 {
		t.Errorf("expected size 1, got %d", r.Size())
	}

	// Nil CWE
	err = r.Register(nil)
	if err == nil {
		t.Error("expected error for nil CWE")
	}

	// Invalid ID
	err = r.Register(&CWE{ID: 0, Name: "Invalid"})
	if err == nil {
		t.Error("expected error for ID=0")
	}

	err = r.Register(&CWE{ID: -1, Name: "Invalid"})
	if err == nil {
		t.Error("expected error for negative ID")
	}

	// Duplicate
	err = r.Register(&CWE{ID: 79, Name: "Duplicate"})
	if err == nil {
		t.Error("expected error for duplicate ID")
	}
}

// ========== RegisterCategory ==========

func TestRegistry_RegisterCategory(t *testing.T) {
	r := NewRegistry()

	err := r.RegisterCategory(&Category{ID: 1, Name: "Cat1"})
	if err != nil {
		t.Errorf("RegisterCategory() error = %v", err)
	}

	err = r.RegisterCategory(nil)
	if err == nil {
		t.Error("expected error for nil category")
	}

	err = r.RegisterCategory(&Category{ID: 0, Name: "Invalid"})
	if err == nil {
		t.Error("expected error for ID=0")
	}

	err = r.RegisterCategory(&Category{ID: 1, Name: "Duplicate"})
	if err == nil {
		t.Error("expected error for duplicate category")
	}
}

// ========== RegisterView ==========

func TestRegistry_RegisterView(t *testing.T) {
	r := NewRegistry()

	err := r.RegisterView(&View{ID: 1000, Name: "View1"})
	if err != nil {
		t.Errorf("RegisterView() error = %v", err)
	}

	err = r.RegisterView(nil)
	if err == nil {
		t.Error("expected error for nil view")
	}

	err = r.RegisterView(&View{ID: 0, Name: "Invalid"})
	if err == nil {
		t.Error("expected error for ID=0")
	}

	err = r.RegisterView(&View{ID: 1000, Name: "Duplicate"})
	if err == nil {
		t.Error("expected error for duplicate view")
	}
}

// ========== RegisterCompoundElement ==========

func TestRegistry_RegisterCompoundElement(t *testing.T) {
	r := NewRegistry()

	err := r.RegisterCompoundElement(&CompoundElement{ID: 680, Name: "Chain1", Structure: StructureChain})
	if err != nil {
		t.Errorf("RegisterCompoundElement() error = %v", err)
	}

	err = r.RegisterCompoundElement(nil)
	if err == nil {
		t.Error("expected error for nil compound element")
	}

	err = r.RegisterCompoundElement(&CompoundElement{ID: 0, Name: "Invalid", Structure: StructureChain})
	if err == nil {
		t.Error("expected error for ID=0")
	}

	err = r.RegisterCompoundElement(&CompoundElement{ID: 680, Name: "Duplicate", Structure: StructureChain})
	if err == nil {
		t.Error("expected error for duplicate compound element")
	}
}

// ========== Get ==========

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})

	cwe, ok := r.Get(79)
	if !ok {
		t.Error("expected to find CWE-79")
	}
	if cwe.Name != "XSS" {
		t.Errorf("expected Name=XSS, got %s", cwe.Name)
	}

	_, ok = r.Get(999)
	if ok {
		t.Error("expected not to find CWE-999")
	}
}

func TestRegistry_GetCategory(t *testing.T) {
	r := NewRegistry()
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1"})

	cat, ok := r.GetCategory(1)
	if !ok {
		t.Error("expected to find category 1")
	}
	if cat.Name != "Cat1" {
		t.Errorf("expected Name=Cat1, got %s", cat.Name)
	}

	_, ok = r.GetCategory(999)
	if ok {
		t.Error("expected not to find category 999")
	}
}

func TestRegistry_GetView(t *testing.T) {
	r := NewRegistry()
	r.RegisterView(&View{ID: 1000, Name: "View1"})

	view, ok := r.GetView(1000)
	if !ok {
		t.Error("expected to find view 1000")
	}
	if view.Name != "View1" {
		t.Errorf("expected Name=View1, got %s", view.Name)
	}

	_, ok = r.GetView(999)
	if ok {
		t.Error("expected not to find view 999")
	}
}

func TestRegistry_GetCompoundElement(t *testing.T) {
	r := NewRegistry()
	r.RegisterCompoundElement(&CompoundElement{ID: 680, Name: "Chain1", Structure: StructureChain})

	ce, ok := r.GetCompoundElement(680)
	if !ok {
		t.Error("expected to find compound element 680")
	}
	if ce.Name != "Chain1" {
		t.Errorf("expected Name=Chain1, got %s", ce.Name)
	}

	_, ok = r.GetCompoundElement(999)
	if ok {
		t.Error("expected not to find compound element 999")
	}
}

// ========== GetAll ==========

func TestRegistry_GetAll(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})
	r.Register(&CWE{ID: 89, Name: "SQLi"})

	all := r.GetAll()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestRegistry_GetAllCategories(t *testing.T) {
	r := NewRegistry()
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1"})

	all := r.GetAllCategories()
	if len(all) != 1 {
		t.Errorf("expected 1 category, got %d", len(all))
	}
}

func TestRegistry_GetAllViews(t *testing.T) {
	r := NewRegistry()
	r.RegisterView(&View{ID: 1000, Name: "View1"})

	all := r.GetAllViews()
	if len(all) != 1 {
		t.Errorf("expected 1 view, got %d", len(all))
	}
}

// ========== Size/Contains ==========

func TestRegistry_Size(t *testing.T) {
	r := NewRegistry()
	if r.Size() != 0 {
		t.Errorf("expected size 0, got %d", r.Size())
	}
	r.Register(&CWE{ID: 79, Name: "XSS"})
	if r.Size() != 1 {
		t.Errorf("expected size 1, got %d", r.Size())
	}
}

func TestRegistry_Contains(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})

	if !r.Contains(79) {
		t.Error("expected Contains(79) = true")
	}
	if r.Contains(999) {
		t.Error("expected Contains(999) = false")
	}
}

// ========== Remove ==========

func TestRegistry_Remove(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})

	err := r.Remove(79)
	if err != nil {
		t.Errorf("Remove() error = %v", err)
	}
	if r.Contains(79) {
		t.Error("expected CWE-79 to be removed")
	}

	err = r.Remove(999)
	if err == nil {
		t.Error("expected error removing non-existent ID")
	}
}

func TestRegistry_RemoveCategory(t *testing.T) {
	r := NewRegistry()
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1"})

	err := r.RemoveCategory(1)
	if err != nil {
		t.Errorf("RemoveCategory() error = %v", err)
	}

	err = r.RemoveCategory(999)
	if err == nil {
		t.Error("expected error removing non-existent category")
	}
}

func TestRegistry_RemoveView(t *testing.T) {
	r := NewRegistry()
	r.RegisterView(&View{ID: 1000, Name: "View1"})

	err := r.RemoveView(1000)
	if err != nil {
		t.Errorf("RemoveView() error = %v", err)
	}

	err = r.RemoveView(999)
	if err == nil {
		t.Error("expected error removing non-existent view")
	}
}

// ========== Clear ==========

func TestRegistry_Clear(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1"})
	r.RegisterView(&View{ID: 1000, Name: "View1"})

	r.Clear()
	if r.Size() != 0 {
		t.Errorf("expected size 0 after Clear, got %d", r.Size())
	}
	if r.CategoryCount() != 0 {
		t.Errorf("expected 0 categories after Clear")
	}
	if r.ViewCount() != 0 {
		t.Errorf("expected 0 views after Clear")
	}
}

// ========== BuildIndexes ==========

func TestRegistry_BuildIndexes_ParentChild(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 1, Name: "Root", Relationships: []Relationship{}})
	r.Register(&CWE{ID: 2, Name: "Child", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 1},
	}})
	r.Register(&CWE{ID: 3, Name: "Grandchild", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 2},
	}})

	r.BuildIndexes()

	parents := r.GetParentIDs(2)
	if !containsInt(parents, 1) {
		t.Errorf("expected parent 1 for CWE-2, got %v", parents)
	}

	children := r.GetChildIDs(1)
	if !containsInt(children, 2) {
		t.Errorf("expected child 2 for CWE-1, got %v", children)
	}

	ancestors := r.GetAncestorIDs(3)
	if !containsInt(ancestors, 1) {
		t.Errorf("expected ancestor 1 for CWE-3, got %v", ancestors)
	}
	if !containsInt(ancestors, 2) {
		t.Errorf("expected ancestor 2 for CWE-3, got %v", ancestors)
	}

	descendants := r.GetDescendantIDs(1)
	if !containsInt(descendants, 2) {
		t.Errorf("expected descendant 2 for CWE-1, got %v", descendants)
	}
	if !containsInt(descendants, 3) {
		t.Errorf("expected descendant 3 for CWE-1, got %v", descendants)
	}
}

func TestRegistry_BuildIndexes_PeerRelations(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", Relationships: []Relationship{
		{Nature: RelationshipPeerOf, CWEID: 89},
		{Nature: RelationshipCanAlsoBe, CWEID: 352},
	}})

	r.BuildIndexes()

	peers := r.GetPeerIDs(79)
	if !containsInt(peers, 89) {
		t.Errorf("expected peer 89, got %v", peers)
	}
	if !containsInt(peers, 352) {
		t.Errorf("expected peer 352, got %v", peers)
	}
}

func TestRegistry_BuildIndexes_CategoryMembers(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1", Relationships: []Relationship{
		{Nature: RelationshipHasMember, CWEID: 79},
	}})

	r.BuildIndexes()

	members := r.GetCategoryMembers(1)
	if !containsInt(members, 79) {
		t.Errorf("expected member 79 for category 1, got %v", members)
	}

	memberOf := r.GetMemberOfIDs(79)
	if !containsInt(memberOf, 1) {
		t.Errorf("expected memberOf 1 for CWE-79, got %v", memberOf)
	}
}

func TestRegistry_BuildIndexes_ViewMembers(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})
	r.RegisterView(&View{ID: 1000, Name: "View1", Members: []ViewMember{
		{CWEID: 79, ViewID: 1000, Direct: true},
	}})

	r.BuildIndexes()

	members := r.GetViewMembers(1000)
	if !containsInt(members, 79) {
		t.Errorf("expected member 79 for view 1000, got %v", members)
	}

	memberOf := r.GetMemberOfIDs(79)
	if !containsInt(memberOf, 1000) {
		t.Errorf("expected memberOf 1000 for CWE-79, got %v", memberOf)
	}
}

func TestRegistry_BuildIndexes_MemberOfInCategory(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 100, Name: "SomeCWE"})
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1", Relationships: []Relationship{
		{Nature: RelationshipMemberOf, CWEID: 100},
	}})

	r.BuildIndexes()

	// MemberOf means category 1 is a member of CWE-100
	memberOf := r.GetMemberOfIDs(1)
	if !containsInt(memberOf, 100) {
		t.Errorf("expected memberOf 100 for category 1, got %v", memberOf)
	}

	members := r.GetCategoryMembers(100)
	if !containsInt(members, 1) {
		t.Errorf("expected category 1 as member of CWE-100, got %v", members)
	}
}

func TestRegistry_IndexesBuilt(t *testing.T) {
	r := NewRegistry()
	if r.IndexesBuilt() {
		t.Error("expected indexes not built initially")
	}
	r.BuildIndexes()
	if !r.IndexesBuilt() {
		t.Error("expected indexes built after BuildIndexes()")
	}
	r.Register(&CWE{ID: 79, Name: "XSS"})
	if r.IndexesBuilt() {
		t.Error("expected indexes not built after modification")
	}
}

// ========== GetAncestorIDs/GetDescendantIDs ==========

func TestRegistry_GetAncestorIDs_NoAncestors(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 1, Name: "Root"})
	r.BuildIndexes()

	ancestors := r.GetAncestorIDs(1)
	if len(ancestors) != 0 {
		t.Errorf("expected no ancestors for root, got %v", ancestors)
	}
}

func TestRegistry_GetDescendantIDs_NoDescendants(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 1, Name: "Leaf"})
	r.BuildIndexes()

	descendants := r.GetDescendantIDs(1)
	if len(descendants) != 0 {
		t.Errorf("expected no descendants for leaf, got %v", descendants)
	}
}

// ========== ExportJSON/ImportJSON ==========

func TestRegistry_ExportJSON_ImportJSON(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", Abstraction: AbstractionBase})
	r.RegisterCategory(&Category{ID: 1, Name: "Cat1"})
	r.RegisterView(&View{ID: 1000, Name: "View1", Type: ViewTypeGraph})
	r.RegisterCompoundElement(&CompoundElement{ID: 680, Name: "Chain1", Structure: StructureChain})

	data, err := r.ExportJSON()
	if err != nil {
		t.Fatalf("ExportJSON() error = %v", err)
	}

	r2 := NewRegistry()
	err = r2.ImportJSON(data)
	if err != nil {
		t.Fatalf("ImportJSON() error = %v", err)
	}

	if r2.Size() != 1 {
		t.Errorf("expected 1 weakness, got %d", r2.Size())
	}
	if r2.CategoryCount() != 1 {
		t.Errorf("expected 1 category, got %d", r2.CategoryCount())
	}
	if r2.ViewCount() != 1 {
		t.Errorf("expected 1 view, got %d", r2.ViewCount())
	}
	if r2.CompoundElementCount() != 1 {
		t.Errorf("expected 1 compound element, got %d", r2.CompoundElementCount())
	}

	cwe, ok := r2.Get(79)
	if !ok || cwe.Name != "XSS" {
		t.Errorf("expected CWE-79 XSS, got %v", cwe)
	}
}

func TestRegistry_ImportJSON_InvalidData(t *testing.T) {
	r := NewRegistry()
	err := r.ImportJSON([]byte("invalid json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// ========== ExportCSV ==========

func TestRegistry_ExportCSV(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS"})

	data, err := r.ExportCSV()
	if err != nil {
		t.Fatalf("ExportCSV() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty CSV output")
	}
}

// ========== Helper ==========

func containsInt(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// Verify JSON round-trip for Registry
func TestRegistry_JSONRoundTrip(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{
		ID:          79,
		Name:        "XSS",
		Abstraction: AbstractionBase,
		Status:      StatusStable,
		Relationships: []Relationship{
			{Nature: RelationshipChildOf, CWEID: 74},
		},
		CommonConsequences: []Consequence{
			{Scopes: []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity}},
		},
	})

	data, err := r.ExportJSON()
	if err != nil {
		t.Fatalf("ExportJSON() error = %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("ExportJSON() produced invalid JSON: %v", err)
	}
}
