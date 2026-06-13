package cwe

import "testing"

func TestNewRelationship(t *testing.T) {
	tests := []struct {
		name     string
		nature   RelationshipNature
		cweID    int
		expected *Relationship
	}{
		{
			name:     "创建ChildOf关系",
			nature:   RelationshipChildOf,
			cweID:    79,
			expected: &Relationship{Nature: RelationshipChildOf, CWEID: 79},
		},
		{
			name:     "创建ParentOf关系",
			nature:   RelationshipParentOf,
			cweID:    1000,
			expected: &Relationship{Nature: RelationshipParentOf, CWEID: 1000},
		},
		{
			name:     "创建CanPrecede关系",
			nature:   RelationshipCanPrecede,
			cweID:    89,
			expected: &Relationship{Nature: RelationshipCanPrecede, CWEID: 89},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRelationship(tt.nature, tt.cweID)
			if got.Nature != tt.expected.Nature {
				t.Errorf("NewRelationship().Nature = %v, want %v", got.Nature, tt.expected.Nature)
			}
			if got.CWEID != tt.expected.CWEID {
				t.Errorf("NewRelationship().CWEID = %v, want %v", got.CWEID, tt.expected.CWEID)
			}
			if got.ViewID != 0 {
				t.Errorf("NewRelationship().ViewID = %v, want 0", got.ViewID)
			}
			if got.Ordinal != "" {
				t.Errorf("NewRelationship().Ordinal = %q, want empty", got.Ordinal)
			}
			if got.ChainID != 0 {
				t.Errorf("NewRelationship().ChainID = %v, want 0", got.ChainID)
			}
		})
	}
}

func TestNewRelationshipWithView(t *testing.T) {
	tests := []struct {
		name   string
		nature RelationshipNature
		cweID  int
		viewID int
	}{
		{
			name:   "创建带视图ID的ChildOf关系",
			nature: RelationshipChildOf,
			cweID:  79,
			viewID: 1000,
		},
		{
			name:   "创建带视图ID的MemberOf关系",
			nature: RelationshipMemberOf,
			cweID:  89,
			viewID: 699,
		},
		{
			name:   "视图ID为零",
			nature: RelationshipPeerOf,
			cweID:  100,
			viewID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRelationshipWithView(tt.nature, tt.cweID, tt.viewID)
			if got.Nature != tt.nature {
				t.Errorf("Nature = %v, want %v", got.Nature, tt.nature)
			}
			if got.CWEID != tt.cweID {
				t.Errorf("CWEID = %v, want %v", got.CWEID, tt.cweID)
			}
			if got.ViewID != tt.viewID {
				t.Errorf("ViewID = %v, want %v", got.ViewID, tt.viewID)
			}
		})
	}
}

func TestRelationship_IsHierarchical(t *testing.T) {
	tests := []struct {
		name     string
		nature   RelationshipNature
		expected bool
	}{
		{name: "ChildOf是层级关系", nature: RelationshipChildOf, expected: true},
		{name: "ParentOf是层级关系", nature: RelationshipParentOf, expected: true},
		{name: "MemberOf是层级关系", nature: RelationshipMemberOf, expected: true},
		{name: "HasMember是层级关系", nature: RelationshipHasMember, expected: true},
		{name: "CanPrecede不是层级关系", nature: RelationshipCanPrecede, expected: false},
		{name: "CanFollow不是层级关系", nature: RelationshipCanFollow, expected: false},
		{name: "Requires不是层级关系", nature: RelationshipRequires, expected: false},
		{name: "RequiredBy不是层级关系", nature: RelationshipRequiredBy, expected: false},
		{name: "PeerOf不是层级关系", nature: RelationshipPeerOf, expected: false},
		{name: "CanAlsoBe不是层级关系", nature: RelationshipCanAlsoBe, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relationship{Nature: tt.nature, CWEID: 1}
			if got := r.IsHierarchical(); got != tt.expected {
				t.Errorf("IsHierarchical() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelationship_IsSequential(t *testing.T) {
	tests := []struct {
		name     string
		nature   RelationshipNature
		expected bool
	}{
		{name: "CanPrecede是顺序关系", nature: RelationshipCanPrecede, expected: true},
		{name: "CanFollow是顺序关系", nature: RelationshipCanFollow, expected: true},
		{name: "ChildOf不是顺序关系", nature: RelationshipChildOf, expected: false},
		{name: "ParentOf不是顺序关系", nature: RelationshipParentOf, expected: false},
		{name: "Requires不是顺序关系", nature: RelationshipRequires, expected: false},
		{name: "PeerOf不是顺序关系", nature: RelationshipPeerOf, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relationship{Nature: tt.nature, CWEID: 1}
			if got := r.IsSequential(); got != tt.expected {
				t.Errorf("IsSequential() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelationship_IsDependency(t *testing.T) {
	tests := []struct {
		name     string
		nature   RelationshipNature
		expected bool
	}{
		{name: "Requires是依赖关系", nature: RelationshipRequires, expected: true},
		{name: "RequiredBy是依赖关系", nature: RelationshipRequiredBy, expected: true},
		{name: "ChildOf不是依赖关系", nature: RelationshipChildOf, expected: false},
		{name: "CanPrecede不是依赖关系", nature: RelationshipCanPrecede, expected: false},
		{name: "PeerOf不是依赖关系", nature: RelationshipPeerOf, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relationship{Nature: tt.nature, CWEID: 1}
			if got := r.IsDependency(); got != tt.expected {
				t.Errorf("IsDependency() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelationship_IsPeer(t *testing.T) {
	tests := []struct {
		name     string
		nature   RelationshipNature
		expected bool
	}{
		{name: "PeerOf是对等关系", nature: RelationshipPeerOf, expected: true},
		{name: "CanAlsoBe是对等关系", nature: RelationshipCanAlsoBe, expected: true},
		{name: "ChildOf不是对等关系", nature: RelationshipChildOf, expected: false},
		{name: "CanPrecede不是对等关系", nature: RelationshipCanPrecede, expected: false},
		{name: "Requires不是对等关系", nature: RelationshipRequires, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relationship{Nature: tt.nature, CWEID: 1}
			if got := r.IsPeer(); got != tt.expected {
				t.Errorf("IsPeer() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelationship_IsPrimary(t *testing.T) {
	tests := []struct {
		name     string
		ordinal  string
		expected bool
	}{
		{name: "Ordinal为Primary返回true", ordinal: "Primary", expected: true},
		{name: "Ordinal为空返回false", ordinal: "", expected: false},
		{name: "Ordinal为Secondary返回false", ordinal: "Secondary", expected: false},
		{name: "Ordinal为小写primary返回false", ordinal: "primary", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Relationship{Nature: RelationshipChildOf, CWEID: 1, Ordinal: tt.ordinal}
			if got := r.IsPrimary(); got != tt.expected {
				t.Errorf("IsPrimary() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRelationship_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rel     Relationship
		wantErr bool
	}{
		{
			name:    "有效关系",
			rel:     Relationship{Nature: RelationshipChildOf, CWEID: 79},
			wantErr: false,
		},
		{
			name:    "无效的关系类型",
			rel:     Relationship{Nature: RelationshipNature("Invalid"), CWEID: 79},
			wantErr: true,
		},
		{
			name:    "CWEID为零",
			rel:     Relationship{Nature: RelationshipChildOf, CWEID: 0},
			wantErr: true,
		},
		{
			name:    "CWEID为负数",
			rel:     Relationship{Nature: RelationshipChildOf, CWEID: -1},
			wantErr: true,
		},
		{
			name:    "空的关系类型",
			rel:     Relationship{Nature: RelationshipNature(""), CWEID: 79},
			wantErr: true,
		},
		{
			name:    "CanAlsoBe有效关系",
			rel:     Relationship{Nature: RelationshipCanAlsoBe, CWEID: 100},
			wantErr: false,
		},
		{
			name:    "HasMember有效关系",
			rel:     Relationship{Nature: RelationshipHasMember, CWEID: 200},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rel.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if _, ok := err.(*ValidationError); !ok {
					t.Errorf("Validate() should return ValidationError, got %T", err)
				}
			}
		})
	}
}