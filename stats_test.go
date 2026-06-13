package cwe

import (
	"testing"
)

func TestComputeStatistics(t *testing.T) {
	t.Run("populated registry", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&CWE{ID: 79, Name: "XSS", Abstraction: AbstractionBase, Status: StatusStable, LikelihoodOfExploit: LikelihoodHigh, Structure: StructureSimple, CommonConsequences: []Consequence{
			{Scopes: []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity}},
		}})
		r.Register(&CWE{ID: 89, Name: "SQLi", Abstraction: AbstractionBase, Status: StatusStable, LikelihoodOfExploit: LikelihoodHigh, Structure: StructureSimple, CommonConsequences: []Consequence{
			{Scopes: []ConsequenceScope{ScopeIntegrity}},
		}})
		r.Register(&CWE{ID: 74, Name: "Injection", Abstraction: AbstractionClass, Status: StatusDraft, Structure: StructureSimple})
		r.RegisterCategory(&Category{ID: 1, Name: "Cat1"})
		r.RegisterView(&View{ID: 1000, Name: "View1"})
		r.RegisterCompoundElement(&CompoundElement{ID: 680, Name: "CE1", Structure: StructureChain})

		stats := ComputeStatistics(r)
		if stats == nil {
			t.Fatal("expected non-nil statistics")
		}
		if stats.TotalCount != 3 {
			t.Errorf("expected TotalCount 3, got %d", stats.TotalCount)
		}
		if stats.WeaknessCount != 3 {
			t.Errorf("expected WeaknessCount 3, got %d", stats.WeaknessCount)
		}
		if stats.CategoryCount != 1 {
			t.Errorf("expected CategoryCount 1, got %d", stats.CategoryCount)
		}
		if stats.ViewCount != 1 {
			t.Errorf("expected ViewCount 1, got %d", stats.ViewCount)
		}
		if stats.CompoundElementCount != 1 {
			t.Errorf("expected CompoundElementCount 1, got %d", stats.CompoundElementCount)
		}
		if stats.ByAbstraction[AbstractionBase] != 2 {
			t.Errorf("expected ByAbstraction[Base] 2, got %d", stats.ByAbstraction[AbstractionBase])
		}
		if stats.ByAbstraction[AbstractionClass] != 1 {
			t.Errorf("expected ByAbstraction[Class] 1, got %d", stats.ByAbstraction[AbstractionClass])
		}
		if stats.ByStatus[StatusStable] != 2 {
			t.Errorf("expected ByStatus[Stable] 2, got %d", stats.ByStatus[StatusStable])
		}
		if stats.ByLikelihood[LikelihoodHigh] != 2 {
			t.Errorf("expected ByLikelihood[High] 2, got %d", stats.ByLikelihood[LikelihoodHigh])
		}
		if len(stats.TopScopes) == 0 {
			t.Error("expected TopScopes to be populated")
		}
		// Integrity should appear in both CWE-79 and CWE-89
		for _, sc := range stats.TopScopes {
			if sc.Scope == ScopeIntegrity && sc.Count != 2 {
				t.Errorf("expected Integrity count 2, got %d", sc.Count)
			}
		}
	})

	t.Run("empty registry", func(t *testing.T) {
		r := NewRegistry()
		stats := ComputeStatistics(r)
		if stats == nil {
			t.Fatal("expected non-nil statistics")
		}
		if stats.TotalCount != 0 {
			t.Errorf("expected TotalCount 0, got %d", stats.TotalCount)
		}
		if len(stats.ByAbstraction) != 0 {
			t.Errorf("expected empty ByAbstraction, got %v", stats.ByAbstraction)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		stats := ComputeStatistics(nil)
		if stats == nil {
			t.Fatal("expected non-nil statistics for nil registry")
		}
		if stats.TotalCount != 0 {
			t.Errorf("expected TotalCount 0, got %d", stats.TotalCount)
		}
		if len(stats.ByAbstraction) != 0 {
			t.Errorf("expected empty ByAbstraction, got %v", stats.ByAbstraction)
		}
	})
}

func TestCountByAbstraction(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", Abstraction: AbstractionBase})
	r.Register(&CWE{ID: 89, Name: "SQLi", Abstraction: AbstractionBase})
	r.Register(&CWE{ID: 74, Name: "Injection", Abstraction: AbstractionClass})

	tests := []struct {
		name        string
		abstraction Abstraction
		want        int
	}{
		{"matches", AbstractionBase, 2},
		{"no matches", AbstractionVariant, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountByAbstraction(r, tt.abstraction); got != tt.want {
				t.Errorf("CountByAbstraction(%v) = %d, want %d", tt.abstraction, got, tt.want)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		if got := CountByAbstraction(nil, AbstractionBase); got != 0 {
			t.Errorf("expected 0 for nil registry, got %d", got)
		}
	})
}

func TestCountByStatus(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", Status: StatusStable})
	r.Register(&CWE{ID: 89, Name: "SQLi", Status: StatusStable})
	r.Register(&CWE{ID: 74, Name: "Injection", Status: StatusDraft})

	tests := []struct {
		name   string
		status Status
		want   int
	}{
		{"matches", StatusStable, 2},
		{"no matches", StatusDeprecated, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountByStatus(r, tt.status); got != tt.want {
				t.Errorf("CountByStatus(%v) = %d, want %d", tt.status, got, tt.want)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		if got := CountByStatus(nil, StatusStable); got != 0 {
			t.Errorf("expected 0 for nil registry, got %d", got)
		}
	})
}

func TestCountByLikelihood(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", LikelihoodOfExploit: LikelihoodHigh})
	r.Register(&CWE{ID: 89, Name: "SQLi", LikelihoodOfExploit: LikelihoodHigh})

	tests := []struct {
		name       string
		likelihood LikelihoodOfExploit
		want       int
	}{
		{"matches", LikelihoodHigh, 2},
		{"no matches", LikelihoodLow, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountByLikelihood(r, tt.likelihood); got != tt.want {
				t.Errorf("CountByLikelihood(%v) = %d, want %d", tt.likelihood, got, tt.want)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		if got := CountByLikelihood(nil, LikelihoodHigh); got != 0 {
			t.Errorf("expected 0 for nil registry, got %d", got)
		}
	})
}

func TestCountByScope(t *testing.T) {
	r := NewRegistry()
	r.Register(&CWE{ID: 79, Name: "XSS", CommonConsequences: []Consequence{
		{Scopes: []ConsequenceScope{ScopeConfidentiality}},
	}})
	r.Register(&CWE{ID: 89, Name: "SQLi", CommonConsequences: []Consequence{
		{Scopes: []ConsequenceScope{ScopeIntegrity, ScopeConfidentiality}},
	}})

	tests := []struct {
		name  string
		scope ConsequenceScope
		want  int
	}{
		{"matches", ScopeConfidentiality, 2},
		{"no matches", ScopeAvailability, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountByScope(r, tt.scope); got != tt.want {
				t.Errorf("CountByScope(%v) = %d, want %d", tt.scope, got, tt.want)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		if got := CountByScope(nil, ScopeConfidentiality); got != 0 {
			t.Errorf("expected 0 for nil registry, got %d", got)
		}
	})
}
