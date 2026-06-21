package cweskills

import (
	"testing"
)

func newFilterTestCWEs() []*CWE {
	return []*CWE{
		{ID: 79, Name: "Cross-site Scripting", Description: "XSS vulnerability", Abstraction: AbstractionBase, Status: StatusStable, LikelihoodOfExploit: LikelihoodHigh, Structure: StructureSimple, CommonConsequences: []Consequence{{Scopes: []ConsequenceScope{ScopeConfidentiality}}}},
		{ID: 89, Name: "SQL Injection", Description: "SQL injection vulnerability", Abstraction: AbstractionBase, Status: StatusStable, LikelihoodOfExploit: LikelihoodHigh, Structure: StructureSimple, CommonConsequences: []Consequence{{Scopes: []ConsequenceScope{ScopeIntegrity}}}},
		{ID: 74, Name: "Injection", Description: "General injection", Abstraction: AbstractionClass, Status: StatusStable, LikelihoodOfExploit: LikelihoodMedium, Structure: StructureSimple},
		{ID: 680, Name: "Integer Overflow", Description: "Integer overflow to buffer overflow", Abstraction: AbstractionBase, Status: StatusDraft, LikelihoodOfExploit: LikelihoodMedium, Structure: StructureChain, CommonConsequences: []Consequence{{Scopes: []ConsequenceScope{ScopeAvailability}}}},
		{ID: 352, Name: "CSRF", Description: "Cross-Site Request Forgery", Abstraction: AbstractionBase, Status: StatusUsable, LikelihoodOfExploit: LikelihoodMedium, Structure: StructureComposite},
	}
}

func TestFilter(t *testing.T) {
	cwes := newFilterTestCWEs()

	tests := []struct {
		name      string
		cwes      []*CWE
		opts      FilterOption
		wantCount int
	}{
		{"empty list", []*CWE{}, FilterOption{}, 0},
		{"no options", cwes, FilterOption{}, 5},
		{"by abstraction", cwes, FilterOption{Abstraction: AbstractionBase}, 4},
		{"by status", cwes, FilterOption{Status: StatusStable}, 3},
		{"by structure", cwes, FilterOption{Structure: StructureChain}, 1},
		{"by likelihood", cwes, FilterOption{Likelihood: LikelihoodHigh}, 2},
		{"by MinID", cwes, FilterOption{MinID: 80}, 3},
		{"by MaxID", cwes, FilterOption{MaxID: 80}, 2},
		{"by keyword", cwes, FilterOption{Keyword: "injection"}, 2},
		{"by scope", cwes, FilterOption{Scope: ScopeConfidentiality}, 1},
		{"combined filters", cwes, FilterOption{Abstraction: AbstractionBase, Status: StatusStable}, 2},
		{"MinID and MaxID", cwes, FilterOption{MinID: 74, MaxID: 89}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.cwes, tt.opts)
			if len(result) != tt.wantCount {
				t.Errorf("Filter() = %d results, want %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestFilter_NoOptions(t *testing.T) {
	cwes := newFilterTestCWEs()
	result := Filter(cwes)
	if len(result) != len(cwes) {
		t.Errorf("Filter with no opts = %d results, want %d", len(result), len(cwes))
	}
}

func TestSortByID(t *testing.T) {
	tests := []struct {
		name  string
		cwes  []*CWE
		want  []int
	}{
		{"empty list", []*CWE{}, nil},
		{"sorted result", []*CWE{
			{ID: 89, Name: "SQLi"},
			{ID: 79, Name: "XSS"},
			{ID: 74, Name: "Injection"},
		}, []int{74, 79, 89}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortByID(tt.cwes)
			if len(result) != len(tt.want) {
				t.Fatalf("SortByID() = %d items, want %d", len(result), len(tt.want))
			}
			for i, cwe := range result {
				if cwe.ID != tt.want[i] {
					t.Errorf("result[%d].ID = %d, want %d", i, cwe.ID, tt.want[i])
				}
			}
		})
	}
}

func TestSortByName(t *testing.T) {
	tests := []struct {
		name string
		cwes []*CWE
		want []string
	}{
		{"empty list", []*CWE{}, nil},
		{"sorted result", []*CWE{
			{ID: 89, Name: "SQL Injection"},
			{ID: 79, Name: "Cross-site Scripting"},
			{ID: 74, Name: "Injection"},
		}, []string{"Cross-site Scripting", "Injection", "SQL Injection"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortByName(tt.cwes)
			if len(result) != len(tt.want) {
				t.Fatalf("SortByName() = %d items, want %d", len(result), len(tt.want))
			}
			for i, cwe := range result {
				if cwe.Name != tt.want[i] {
					t.Errorf("result[%d].Name = %q, want %q", i, cwe.Name, tt.want[i])
				}
			}
		})
	}
}

func TestSortByAbstraction(t *testing.T) {
	tests := []struct {
		name string
		cwes []*CWE
		want []Abstraction
	}{
		{"empty list", []*CWE{}, nil},
		{"sorted result", []*CWE{
			{ID: 79, Name: "XSS", Abstraction: AbstractionBase},
			{ID: 664, Name: "Resource", Abstraction: AbstractionPillar},
			{ID: 74, Name: "Injection", Abstraction: AbstractionClass},
			{ID: 83, Name: "Variant XSS", Abstraction: AbstractionVariant},
		}, []Abstraction{AbstractionPillar, AbstractionClass, AbstractionBase, AbstractionVariant}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortByAbstraction(tt.cwes)
			if len(result) != len(tt.want) {
				t.Fatalf("SortByAbstraction() = %d items, want %d", len(result), len(tt.want))
			}
			for i, cwe := range result {
				if cwe.Abstraction != tt.want[i] {
					t.Errorf("result[%d].Abstraction = %v, want %v", i, cwe.Abstraction, tt.want[i])
				}
			}
		})
	}
}

func TestGroupByAbstraction(t *testing.T) {
	cwes := newFilterTestCWEs()
	groups := GroupByAbstraction(cwes)

	if len(groups[AbstractionBase]) != 4 {
		t.Errorf("GroupByAbstraction()[Base] = %d, want 4", len(groups[AbstractionBase]))
	}
	if len(groups[AbstractionClass]) != 1 {
		t.Errorf("GroupByAbstraction()[Class] = %d, want 1", len(groups[AbstractionClass]))
	}
}

func TestGroupByStatus(t *testing.T) {
	cwes := newFilterTestCWEs()
	groups := GroupByStatus(cwes)

	if len(groups[StatusStable]) != 3 {
		t.Errorf("GroupByStatus()[Stable] = %d, want 3", len(groups[StatusStable]))
	}
	if len(groups[StatusDraft]) != 1 {
		t.Errorf("GroupByStatus()[Draft] = %d, want 1", len(groups[StatusDraft]))
	}
}

func TestGroupByLikelihood(t *testing.T) {
	cwes := newFilterTestCWEs()
	groups := GroupByLikelihood(cwes)

	if len(groups[LikelihoodHigh]) != 2 {
		t.Errorf("GroupByLikelihood()[High] = %d, want 2", len(groups[LikelihoodHigh]))
	}
	if len(groups[LikelihoodMedium]) != 3 {
		t.Errorf("GroupByLikelihood()[Medium] = %d, want 3", len(groups[LikelihoodMedium]))
	}
}

func TestDeduplicate(t *testing.T) {
	tests := []struct {
		name      string
		cwes      []*CWE
		wantCount int
	}{
		{"with duplicates", []*CWE{
			{ID: 79, Name: "XSS"},
			{ID: 89, Name: "SQLi"},
			{ID: 79, Name: "XSS Dup"},
		}, 2},
		{"without duplicates", []*CWE{
			{ID: 79, Name: "XSS"},
			{ID: 89, Name: "SQLi"},
		}, 2},
		{"empty list", []*CWE{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Deduplicate(tt.cwes)
			if len(result) != tt.wantCount {
				t.Errorf("Deduplicate() = %d, want %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestFilter_KeywordMatching(t *testing.T) {
	cwes := newFilterTestCWEs()

	// Test keyword matching in name
	result := Filter(cwes, FilterOption{Keyword: "CSRF"})
	if len(result) != 1 || result[0].ID != 352 {
		t.Errorf("Filter keyword=CSRF: expected CWE-352, got %v", result)
	}

	// Test keyword matching in description (case insensitive)
	result = Filter(cwes, FilterOption{Keyword: "VULNERABILITY"})
	if len(result) != 2 {
		t.Errorf("Filter keyword=VULNERABILITY: expected 2, got %d", len(result))
	}
}

func TestFilter_ScopeFilter(t *testing.T) {
	cwes := newFilterTestCWEs()

	result := Filter(cwes, FilterOption{Scope: ScopeAvailability})
	if len(result) != 1 || result[0].ID != 680 {
		t.Errorf("Filter scope=Availability: expected CWE-680, got %v", result)
	}
}

func TestFilter_CombinedMinIDMaxID(t *testing.T) {
	cwes := newFilterTestCWEs()

	result := Filter(cwes, FilterOption{MinID: 74, MaxID: 89})
	if len(result) != 3 {
		t.Errorf("Filter MinID=74 MaxID=89: expected 3, got %d", len(result))
	}
}
