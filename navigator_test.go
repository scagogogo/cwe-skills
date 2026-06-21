package cweskills

import (
	"testing"
)

func newNavigatorTestRegistry() *Registry {
	r := NewRegistry()
	// Build a clean hierarchy: 1 -> 2,3; 2 -> 4
	// Only declare relationships from one direction to avoid duplicate index entries
	r.Register(&CWE{ID: 1, Name: "Root", Relationships: []Relationship{
		{Nature: RelationshipCanPrecede, CWEID: 10},
		{Nature: RelationshipCanAlsoBe, CWEID: 11},
	}})
	r.Register(&CWE{ID: 2, Name: "Mid", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 1},
		{Nature: RelationshipCanFollow, CWEID: 5},
		{Nature: RelationshipRequires, CWEID: 6},
		{Nature: RelationshipRequiredBy, CWEID: 7},
	}})
	r.Register(&CWE{ID: 3, Name: "Mid2", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 1},
		{Nature: RelationshipPeerOf, CWEID: 99},
	}})
	r.Register(&CWE{ID: 4, Name: "Leaf", Relationships: []Relationship{
		{Nature: RelationshipChildOf, CWEID: 2},
	}})
	r.Register(&CWE{ID: 5, Name: "Follower", Relationships: []Relationship{}})
	r.Register(&CWE{ID: 6, Name: "Required", Relationships: []Relationship{}})
	r.Register(&CWE{ID: 7, Name: "Requirer", Relationships: []Relationship{}})
	r.Register(&CWE{ID: 10, Name: "Preceded", Relationships: []Relationship{}})
	r.Register(&CWE{ID: 11, Name: "AlsoBe", Relationships: []Relationship{}})
	r.Register(&CWE{ID: 99, Name: "Peer99", Relationships: []Relationship{}})

	r.RegisterCompoundElement(&CompoundElement{ID: 100, Name: "Chain1", Structure: StructureChain, Relationships: []Relationship{
		{Nature: RelationshipCanFollow, CWEID: 1},
		{Nature: RelationshipCanPrecede, CWEID: 2},
	}})
	r.RegisterCompoundElement(&CompoundElement{ID: 101, Name: "Composite1", Structure: StructureComposite, Relationships: []Relationship{
		{Nature: RelationshipRequires, CWEID: 1},
		{Nature: RelationshipRequires, CWEID: 2},
	}})

	r.BuildIndexes()
	return r
}

func TestNewNavigator(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)
	if nav == nil {
		t.Fatal("expected non-nil navigator")
	}
}

func TestNavigator_Parents(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	parents := nav.Parents(4)
	if len(parents) != 1 || parents[0].ID != 2 {
		t.Errorf("Parents(4) = %v, expected [2]", idsFromCWEs(parents))
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.Parents(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_Children(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	children := nav.Children(1)
	if len(children) != 2 {
		t.Errorf("Children(1) = %d items, expected 2", len(children))
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.Children(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_Ancestors(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	ancestors := nav.Ancestors(4)
	if len(ancestors) != 2 {
		t.Errorf("Ancestors(4) = %d items, expected 2", len(ancestors))
	}

	t.Run("no ancestors", func(t *testing.T) {
		ancestors := nav.Ancestors(1)
		if len(ancestors) != 0 {
			t.Errorf("Ancestors(1) = %d items, expected 0", len(ancestors))
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.Ancestors(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_Descendants(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	descendants := nav.Descendants(1)
	if len(descendants) != 3 {
		t.Errorf("Descendants(1) = %d items, expected 3", len(descendants))
	}

	t.Run("no descendants", func(t *testing.T) {
		descendants := nav.Descendants(4)
		if len(descendants) != 0 {
			t.Errorf("Descendants(4) = %d items, expected 0", len(descendants))
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.Descendants(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_Siblings(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	siblings := nav.Siblings(2)
	if len(siblings) != 1 {
		t.Fatalf("Siblings(2) = %d items, expected 1", len(siblings))
	}
	if siblings[0].ID != 3 {
		t.Errorf("Siblings(2)[0].ID = %d, expected 3", siblings[0].ID)
	}

	t.Run("no parents", func(t *testing.T) {
		siblings := nav.Siblings(1)
		if siblings != nil {
			t.Errorf("Siblings(1) = %v, expected nil", siblings)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.Siblings(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_Peers(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	peers := nav.Peers(3)
	if len(peers) != 1 || peers[0].ID != 99 {
		t.Errorf("Peers(3) = %v, expected [99]", idsFromCWEs(peers))
	}

	t.Run("no peers", func(t *testing.T) {
		peers := nav.Peers(4)
		if len(peers) != 0 {
			t.Errorf("Peers(4) = %d, expected 0", len(peers))
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.Peers(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_CanPrecede(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	result := nav.CanPrecede(1)
	if len(result) != 1 || result[0].ID != 10 {
		t.Errorf("CanPrecede(1) = %v, expected [10]", idsFromCWEs(result))
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.CanPrecede(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_CanFollow(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	result := nav.CanFollow(2)
	if len(result) != 1 || result[0].ID != 5 {
		t.Errorf("CanFollow(2) = %v, expected [5]", idsFromCWEs(result))
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.CanFollow(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_Requires(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	result := nav.Requires(2)
	if len(result) != 1 || result[0].ID != 6 {
		t.Errorf("Requires(2) = %v, expected [6]", idsFromCWEs(result))
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.Requires(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_RequiredBy(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	result := nav.RequiredBy(2)
	if len(result) != 1 || result[0].ID != 7 {
		t.Errorf("RequiredBy(2) = %v, expected [7]", idsFromCWEs(result))
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.RequiredBy(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_CanAlsoBe(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	result := nav.CanAlsoBe(1)
	if len(result) != 1 || result[0].ID != 11 {
		t.Errorf("CanAlsoBe(1) = %v, expected [11]", idsFromCWEs(result))
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.CanAlsoBe(1) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_ChainMembers(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	t.Run("chain type", func(t *testing.T) {
		result := nav.ChainMembers(100)
		if len(result) != 2 {
			t.Errorf("ChainMembers(100) = %d, expected 2", len(result))
		}
	})

	t.Run("non-chain type", func(t *testing.T) {
		result := nav.ChainMembers(101)
		if result != nil {
			t.Errorf("ChainMembers(101) = %v, expected nil for composite", result)
		}
	})

	t.Run("not found", func(t *testing.T) {
		result := nav.ChainMembers(999)
		if result != nil {
			t.Errorf("ChainMembers(999) = %v, expected nil", result)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.ChainMembers(100) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_CompositeMembers(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	t.Run("composite type", func(t *testing.T) {
		result := nav.CompositeMembers(101)
		if len(result) != 2 {
			t.Errorf("CompositeMembers(101) = %d, expected 2", len(result))
		}
	})

	t.Run("non-composite type", func(t *testing.T) {
		result := nav.CompositeMembers(100)
		if result != nil {
			t.Errorf("CompositeMembers(100) = %v, expected nil for chain", result)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.CompositeMembers(101) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_ShortestPath(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	tests := []struct {
		name string
		from int
		to   int
		want []int
	}{
		{"same ID", 1, 1, []int{1}},
		{"direct parent-child", 1, 2, []int{1, 2}},
		{"multi-hop", 1, 4, []int{1, 2, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := nav.ShortestPath(tt.from, tt.to)
			if len(path) != len(tt.want) {
				t.Fatalf("ShortestPath(%d, %d) = %v, expected %v", tt.from, tt.to, path, tt.want)
			}
			for i, id := range path {
				if id != tt.want[i] {
					t.Errorf("ShortestPath(%d, %d)[%d] = %d, expected %d", tt.from, tt.to, i, id, tt.want[i])
				}
			}
		})
	}

	t.Run("no direct connection between disconnected nodes", func(t *testing.T) {
		path := nav.ShortestPath(5, 6)
		// 5 and 6 have no direct relationship
		if path != nil && len(path) < 2 {
			t.Errorf("ShortestPath(5, 6) = %v, expected nil or multi-hop path", path)
		}
	})

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.ShortestPath(1, 2) != nil {
			t.Error("expected nil for nil registry")
		}
	})
}

func TestNavigator_IsAncestorOf(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	tests := []struct {
		name     string
		ancestor int
		desc     int
		want     bool
	}{
		{"true", 1, 4, true},
		{"false", 4, 1, false},
		{"same", 1, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nav.IsAncestorOf(tt.ancestor, tt.desc); got != tt.want {
				t.Errorf("IsAncestorOf(%d, %d) = %v, want %v", tt.ancestor, tt.desc, got, tt.want)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.IsAncestorOf(1, 2) {
			t.Error("expected false for nil registry")
		}
	})
}

func TestNavigator_IsDescendantOf(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	tests := []struct {
		name       string
		descendant int
		ancestor   int
		want       bool
	}{
		{"true", 4, 1, true},
		{"false", 1, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nav.IsDescendantOf(tt.descendant, tt.ancestor); got != tt.want {
				t.Errorf("IsDescendantOf(%d, %d) = %v, want %v", tt.descendant, tt.ancestor, got, tt.want)
			}
		})
	}
}

func TestNavigator_IsRelated(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	tests := []struct {
		name string
		a    int
		b    int
		want bool
	}{
		{"true forward", 2, 1, true},
		{"true reverse", 1, 2, true},
		{"false", 5, 99, false},
		{"not found", 999, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nav.IsRelated(tt.a, tt.b); got != tt.want {
				t.Errorf("IsRelated(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.IsRelated(1, 2) {
			t.Error("expected false for nil registry")
		}
	})
}

func TestNavigator_RelationshipDepth(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	tests := []struct {
		name     string
		ancestor int
		desc     int
		want     int
	}{
		{"same ID", 1, 1, 0},
		{"direct parent", 1, 2, 1},
		{"multi-level", 1, 4, 2},
		{"no relationship", 4, 99, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nav.RelationshipDepth(tt.ancestor, tt.desc); got != tt.want {
				t.Errorf("RelationshipDepth(%d, %d) = %d, want %d", tt.ancestor, tt.desc, got, tt.want)
			}
		})
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.RelationshipDepth(1, 2) != -1 {
			t.Error("expected -1 for nil registry")
		}
	})
}

func TestNavigator_String(t *testing.T) {
	r := newNavigatorTestRegistry()
	nav := NewNavigator(r)

	s := nav.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
	if s == "Navigator<nil>" {
		t.Errorf("unexpected nil navigator string: %q", s)
	}

	t.Run("nil registry", func(t *testing.T) {
		nav := NewNavigator(nil)
		if nav.String() != "Navigator<nil>" {
			t.Errorf("expected 'Navigator<nil>', got %q", nav.String())
		}
	})
}

func idsFromCWEs(cwes []*CWE) []int {
	ids := make([]int, len(cwes))
	for i, c := range cwes {
		ids[i] = c.ID
	}
	return ids
}
