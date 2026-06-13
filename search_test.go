package cwe

import (
	"testing"
)

func newTestRegistry() *Registry {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "Cross-site Scripting", Description: "XSS vulnerability", Abstraction: AbstractionBase, Status: StatusStable, LikelihoodOfExploit: LikelihoodHigh, Structure: StructureSimple})
	r.Register(&CWE{ID: 89, Name: "SQL Injection", Description: "SQL injection vulnerability", Abstraction: AbstractionBase, Status: StatusStable, LikelihoodOfExploit: LikelihoodHigh, Structure: StructureSimple})
	r.Register(&CWE{ID: 74, Name: "Injection", Description: "General injection", Abstraction: AbstractionClass, Status: StatusStable, Structure: StructureSimple})
	r.Register(&CWE{ID: 664, Name: "Improper Resource Control", Description: "Resource control issue", Abstraction: AbstractionPillar, Status: StatusStable, Structure: StructureSimple})
	r.Register(&CWE{ID: 680, Name: "Integer Overflow to Buffer Overflow", Description: "Chain overflow", Abstraction: AbstractionBase, Status: StatusDraft, LikelihoodOfExploit: LikelihoodMedium, Structure: StructureChain})
	r.Register(&CWE{ID: 352, Name: "CSRF", Description: "Cross-Site Request Forgery", Abstraction: AbstractionBase, Status: StatusUsable, LikelihoodOfExploit: LikelihoodMedium, Structure: StructureComposite})
	r.Register(&CWE{ID: 20, Name: "Improper Input Validation", Description: "Input not validated", Abstraction: AbstractionPillar, Status: StatusDeprecated, Structure: StructureSimple, CommonConsequences: []Consequence{
		{Scopes: []ConsequenceScope{ScopeIntegrity}},
	}})
	r.Register(&CWE{ID: 200, Name: "Information Exposure", Description: "Sensitive info exposed", Abstraction: AbstractionBase, Status: StatusStable, Structure: StructureSimple, CommonConsequences: []Consequence{
		{Scopes: []ConsequenceScope{ScopeConfidentiality}},
	}})
	return r
}

func TestFindByID(t *testing.T) {
	r := newTestRegistry()

	tests := []struct {
		name    string
		id      int
		want    bool
		wantNil bool
	}{
		{"found", 79, true, false},
		{"not found", 999, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwe, ok := FindByID(r, tt.id)
			if ok != tt.want {
				t.Errorf("FindByID() ok = %v, want %v", ok, tt.want)
			}
			if (cwe == nil) != tt.wantNil {
				t.Errorf("FindByID() nil = %v, wantNil %v", cwe == nil, tt.wantNil)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		cwe, ok := FindByID(nil, 79)
		if ok {
			t.Error("expected false for nil registry")
		}
		if cwe != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestFindByKeyword(t *testing.T) {
	r := newTestRegistry()

	tests := []struct {
		name        string
		keyword     string
		wantCount   int
		wantNameSub string
	}{
		{"matches in name", "Injection", 2, ""},            // 89, 74
		{"matches in description", "vulnerability", 2, ""},  // 79, 89
		{"no matches", "nonexistent", 0, ""},
		{"case insensitive", "cross-site", 2, ""}, // 79 (name), 352 (description)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindByKeyword(r, tt.keyword)
			if len(results) != tt.wantCount {
				t.Errorf("FindByKeyword(%q) = %d results, want %d", tt.keyword, len(results), tt.wantCount)
			}
		})
	}

	t.Run("empty keyword", func(t *testing.T) {
		results := FindByKeyword(r, "")
		if results != nil {
			t.Errorf("expected nil for empty keyword, got %v", results)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		results := FindByKeyword(nil, "test")
		if results != nil {
			t.Errorf("expected nil for nil registry, got %v", results)
		}
	})
}

func TestFindByAbstraction(t *testing.T) {
	r := newTestRegistry()

	tests := []struct {
		name        string
		abstraction Abstraction
		wantCount   int
	}{
		{"matches", AbstractionBase, 5},  // 79, 89, 680, 352, 200
		{"no matches", AbstractionVariant, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindByAbstraction(r, tt.abstraction)
			if len(results) != tt.wantCount {
				t.Errorf("FindByAbstraction(%v) = %d results, want %d", tt.abstraction, len(results), tt.wantCount)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		results := FindByAbstraction(nil, AbstractionBase)
		if results != nil {
			t.Errorf("expected nil for nil registry, got %v", results)
		}
	})
}

func TestFindByStatus(t *testing.T) {
	r := newTestRegistry()

	tests := []struct {
		name      string
		status    Status
		wantCount int
	}{
		{"matches", StatusStable, 5}, // 79, 89, 74, 664, 200
		{"no matches", StatusIncomplete, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindByStatus(r, tt.status)
			if len(results) != tt.wantCount {
				t.Errorf("FindByStatus(%v) = %d results, want %d", tt.status, len(results), tt.wantCount)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		results := FindByStatus(nil, StatusStable)
		if results != nil {
			t.Errorf("expected nil for nil registry, got %v", results)
		}
	})
}

func TestFindByLikelihood(t *testing.T) {
	r := newTestRegistry()

	tests := []struct {
		name       string
		likelihood LikelihoodOfExploit
		wantCount  int
	}{
		{"matches", LikelihoodHigh, 2}, // 79, 89
		{"no matches", LikelihoodLow, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindByLikelihood(r, tt.likelihood)
			if len(results) != tt.wantCount {
				t.Errorf("FindByLikelihood(%v) = %d results, want %d", tt.likelihood, len(results), tt.wantCount)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		results := FindByLikelihood(nil, LikelihoodHigh)
		if results != nil {
			t.Errorf("expected nil for nil registry, got %v", results)
		}
	})
}

func TestFindByConsequenceScope(t *testing.T) {
	r := newTestRegistry()

	tests := []struct {
		name      string
		scope     ConsequenceScope
		wantCount int
	}{
		{"matches", ScopeIntegrity, 1},    // 20
		{"matches confidentiality", ScopeConfidentiality, 1}, // 200
		{"no matches", ScopeAvailability, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindByConsequenceScope(r, tt.scope)
			if len(results) != tt.wantCount {
				t.Errorf("FindByConsequenceScope(%v) = %d results, want %d", tt.scope, len(results), tt.wantCount)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		results := FindByConsequenceScope(nil, ScopeIntegrity)
		if results != nil {
			t.Errorf("expected nil for nil registry, got %v", results)
		}
	})
}

func TestFindByStructure(t *testing.T) {
	r := newTestRegistry()

	tests := []struct {
		name      string
		structure Structure
		wantCount int
	}{
		{"chain", StructureChain, 1},     // 680
		{"composite", StructureComposite, 1}, // 352
		{"no matches variant", StructureSimple, 6}, // 79,89,74,664,20,200
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindByStructure(r, tt.structure)
			if len(results) != tt.wantCount {
				t.Errorf("FindByStructure(%v) = %d results, want %d", tt.structure, len(results), tt.wantCount)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		results := FindByStructure(nil, StructureChain)
		if results != nil {
			t.Errorf("expected nil for nil registry, got %v", results)
		}
	})
}

func TestFindTopLevel(t *testing.T) {
	r := newTestRegistry()

	results := FindTopLevel(r)
	if len(results) != 2 { // 664, 20
		t.Errorf("FindTopLevel() = %d results, want 2", len(results))
	}
}

func TestFindBaseWeaknesses(t *testing.T) {
	r := newTestRegistry()

	results := FindBaseWeaknesses(r)
	if len(results) != 5 { // 79, 89, 680, 352, 200
		t.Errorf("FindBaseWeaknesses() = %d results, want 5", len(results))
	}
}

func TestFindChains(t *testing.T) {
	r := newTestRegistry()

	results := FindChains(r)
	if len(results) != 1 {
		t.Errorf("FindChains() = %d results, want 1", len(results))
	}
	if len(results) > 0 && results[0].ID != 680 {
		t.Errorf("expected CWE-680, got CWE-%d", results[0].ID)
	}
}

func TestFindComposites(t *testing.T) {
	r := newTestRegistry()

	results := FindComposites(r)
	if len(results) != 1 {
		t.Errorf("FindComposites() = %d results, want 1", len(results))
	}
	if len(results) > 0 && results[0].ID != 352 {
		t.Errorf("expected CWE-352, got CWE-%d", results[0].ID)
	}
}
