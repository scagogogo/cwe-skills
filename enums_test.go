package cwe

import (
	"testing"
)

// --- Abstraction ---

func TestAbstraction_String(t *testing.T) {
	tests := []struct {
		name  string
		value Abstraction
		want  string
	}{
		{"Pillar", AbstractionPillar, "Pillar"},
		{"Class", AbstractionClass, "Class"},
		{"Base", AbstractionBase, "Base"},
		{"Variant", AbstractionVariant, "Variant"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("Abstraction.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAbstraction_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value Abstraction
		want  bool
	}{
		{"Pillar is valid", AbstractionPillar, true},
		{"Class is valid", AbstractionClass, true},
		{"Base is valid", AbstractionBase, true},
		{"Variant is valid", AbstractionVariant, true},
		{"Invalid value", Abstraction("Invalid"), false},
		{"Empty value", Abstraction(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("Abstraction.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAbstraction(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Abstraction
		wantErr bool
	}{
		{"Pillar", "Pillar", AbstractionPillar, false},
		{"Class", "Class", AbstractionClass, false},
		{"Base", "Base", AbstractionBase, false},
		{"Variant", "Variant", AbstractionVariant, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAbstraction(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAbstraction(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseAbstraction(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllAbstractionValues(t *testing.T) {
	values := AllAbstractionValues()
	if len(values) != 4 {
		t.Errorf("AllAbstractionValues() returned %d values, want 4", len(values))
	}
	expected := []Abstraction{AbstractionPillar, AbstractionClass, AbstractionBase, AbstractionVariant}
	for i, v := range expected {
		if values[i] != v {
			t.Errorf("AllAbstractionValues()[%d] = %q, want %q", i, values[i], v)
		}
	}
}

func TestAbstractionOrder(t *testing.T) {
	tests := []struct {
		name  string
		value Abstraction
		want  int
	}{
		{"Pillar", AbstractionPillar, 4},
		{"Class", AbstractionClass, 3},
		{"Base", AbstractionBase, 2},
		{"Variant", AbstractionVariant, 1},
		{"Unknown", Abstraction("Unknown"), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.AbstractionOrder(); got != tt.want {
				t.Errorf("AbstractionOrder() = %d, want %d", got, tt.want)
			}
		})
	}
}

// --- Structure ---

func TestStructure_String(t *testing.T) {
	tests := []struct {
		name  string
		value Structure
		want  string
	}{
		{"Simple", StructureSimple, "Simple"},
		{"Chain", StructureChain, "Chain"},
		{"Composite", StructureComposite, "Composite"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("Structure.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStructure_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value Structure
		want  bool
	}{
		{"Simple is valid", StructureSimple, true},
		{"Chain is valid", StructureChain, true},
		{"Composite is valid", StructureComposite, true},
		{"Invalid value", Structure("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("Structure.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStructure(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Structure
		wantErr bool
	}{
		{"Simple", "Simple", StructureSimple, false},
		{"Chain", "Chain", StructureChain, false},
		{"Composite", "Composite", StructureComposite, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStructure(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStructure(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStructure(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllStructureValues(t *testing.T) {
	values := AllStructureValues()
	if len(values) != 3 {
		t.Errorf("AllStructureValues() returned %d values, want 3", len(values))
	}
}

// --- Status ---

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name  string
		value Status
		want  string
	}{
		{"Stable", StatusStable, "Stable"},
		{"Usable", StatusUsable, "Usable"},
		{"Draft", StatusDraft, "Draft"},
		{"Incomplete", StatusIncomplete, "Incomplete"},
		{"Obsolete", StatusObsolete, "Obsolete"},
		{"Deprecated", StatusDeprecated, "Deprecated"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("Status.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value Status
		want  bool
	}{
		{"Stable", StatusStable, true},
		{"Usable", StatusUsable, true},
		{"Draft", StatusDraft, true},
		{"Incomplete", StatusIncomplete, true},
		{"Obsolete", StatusObsolete, true},
		{"Deprecated", StatusDeprecated, true},
		{"Invalid", Status("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("Status.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Status
		wantErr bool
	}{
		{"Stable", "Stable", StatusStable, false},
		{"Usable", "Usable", StatusUsable, false},
		{"Draft", "Draft", StatusDraft, false},
		{"Incomplete", "Incomplete", StatusIncomplete, false},
		{"Obsolete", "Obsolete", StatusObsolete, false},
		{"Deprecated", "Deprecated", StatusDeprecated, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStatus(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllStatusValues(t *testing.T) {
	values := AllStatusValues()
	if len(values) != 6 {
		t.Errorf("AllStatusValues() returned %d values, want 6", len(values))
	}
}

// --- LikelihoodOfExploit ---

func TestLikelihoodOfExploit_String(t *testing.T) {
	tests := []struct {
		name  string
		value LikelihoodOfExploit
		want  string
	}{
		{"High", LikelihoodHigh, "High"},
		{"Medium", LikelihoodMedium, "Medium"},
		{"Low", LikelihoodLow, "Low"},
		{"Unknown", LikelihoodUnknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("LikelihoodOfExploit.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLikelihoodOfExploit_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value LikelihoodOfExploit
		want  bool
	}{
		{"High", LikelihoodHigh, true},
		{"Medium", LikelihoodMedium, true},
		{"Low", LikelihoodLow, true},
		{"Unknown", LikelihoodUnknown, true},
		{"Invalid", LikelihoodOfExploit("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("LikelihoodOfExploit.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLikelihoodOfExploit(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    LikelihoodOfExploit
		wantErr bool
	}{
		{"High", "High", LikelihoodHigh, false},
		{"Medium", "Medium", LikelihoodMedium, false},
		{"Low", "Low", LikelihoodLow, false},
		{"Unknown", "Unknown", LikelihoodUnknown, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLikelihoodOfExploit(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLikelihoodOfExploit(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLikelihoodOfExploit(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllLikelihoodOfExploitValues(t *testing.T) {
	values := AllLikelihoodOfExploitValues()
	if len(values) != 4 {
		t.Errorf("AllLikelihoodOfExploitValues() returned %d values, want 4", len(values))
	}
}

func TestLikelihoodOrder(t *testing.T) {
	tests := []struct {
		name  string
		value LikelihoodOfExploit
		want  int
	}{
		{"High", LikelihoodHigh, 4},
		{"Medium", LikelihoodMedium, 3},
		{"Low", LikelihoodLow, 2},
		{"Unknown", LikelihoodUnknown, 1},
		{"UnknownValue", LikelihoodOfExploit("Foo"), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.LikelihoodOrder(); got != tt.want {
				t.Errorf("LikelihoodOrder() = %d, want %d", got, tt.want)
			}
		})
	}
}

// --- RelationshipNature ---

func TestRelationshipNature_String(t *testing.T) {
	tests := []struct {
		name  string
		value RelationshipNature
		want  string
	}{
		{"ChildOf", RelationshipChildOf, "ChildOf"},
		{"ParentOf", RelationshipParentOf, "ParentOf"},
		{"CanPrecede", RelationshipCanPrecede, "CanPrecede"},
		{"CanFollow", RelationshipCanFollow, "CanFollow"},
		{"Requires", RelationshipRequires, "Requires"},
		{"RequiredBy", RelationshipRequiredBy, "RequiredBy"},
		{"CanAlsoBe", RelationshipCanAlsoBe, "CanAlsoBe"},
		{"PeerOf", RelationshipPeerOf, "PeerOf"},
		{"MemberOf", RelationshipMemberOf, "MemberOf"},
		{"HasMember", RelationshipHasMember, "Has_Member"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("RelationshipNature.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRelationshipNature_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value RelationshipNature
		want  bool
	}{
		{"ChildOf", RelationshipChildOf, true},
		{"ParentOf", RelationshipParentOf, true},
		{"CanPrecede", RelationshipCanPrecede, true},
		{"CanFollow", RelationshipCanFollow, true},
		{"Requires", RelationshipRequires, true},
		{"RequiredBy", RelationshipRequiredBy, true},
		{"CanAlsoBe", RelationshipCanAlsoBe, true},
		{"PeerOf", RelationshipPeerOf, true},
		{"MemberOf", RelationshipMemberOf, true},
		{"HasMember", RelationshipHasMember, true},
		{"Invalid", RelationshipNature("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("RelationshipNature.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRelationshipNature(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    RelationshipNature
		wantErr bool
	}{
		{"ChildOf", "ChildOf", RelationshipChildOf, false},
		{"ParentOf", "ParentOf", RelationshipParentOf, false},
		{"CanPrecede", "CanPrecede", RelationshipCanPrecede, false},
		{"CanFollow", "CanFollow", RelationshipCanFollow, false},
		{"Requires", "Requires", RelationshipRequires, false},
		{"RequiredBy", "RequiredBy", RelationshipRequiredBy, false},
		{"CanAlsoBe", "CanAlsoBe", RelationshipCanAlsoBe, false},
		{"PeerOf", "PeerOf", RelationshipPeerOf, false},
		{"MemberOf", "MemberOf", RelationshipMemberOf, false},
		{"HasMember", "Has_Member", RelationshipHasMember, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRelationshipNature(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRelationshipNature(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseRelationshipNature(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllRelationshipNatureValues(t *testing.T) {
	values := AllRelationshipNatureValues()
	if len(values) != 10 {
		t.Errorf("AllRelationshipNatureValues() returned %d values, want 10", len(values))
	}
}

func TestRelationshipNature_IsHierarchical(t *testing.T) {
	tests := []struct {
		name  string
		value RelationshipNature
		want  bool
	}{
		{"ChildOf is hierarchical", RelationshipChildOf, true},
		{"ParentOf is hierarchical", RelationshipParentOf, true},
		{"MemberOf is hierarchical", RelationshipMemberOf, true},
		{"HasMember is hierarchical", RelationshipHasMember, true},
		{"CanPrecede is not hierarchical", RelationshipCanPrecede, false},
		{"CanFollow is not hierarchical", RelationshipCanFollow, false},
		{"Requires is not hierarchical", RelationshipRequires, false},
		{"RequiredBy is not hierarchical", RelationshipRequiredBy, false},
		{"CanAlsoBe is not hierarchical", RelationshipCanAlsoBe, false},
		{"PeerOf is not hierarchical", RelationshipPeerOf, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsHierarchical(); got != tt.want {
				t.Errorf("IsHierarchical() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipNature_IsSequential(t *testing.T) {
	tests := []struct {
		name  string
		value RelationshipNature
		want  bool
	}{
		{"CanPrecede is sequential", RelationshipCanPrecede, true},
		{"CanFollow is sequential", RelationshipCanFollow, true},
		{"ChildOf is not sequential", RelationshipChildOf, false},
		{"ParentOf is not sequential", RelationshipParentOf, false},
		{"Requires is not sequential", RelationshipRequires, false},
		{"RequiredBy is not sequential", RelationshipRequiredBy, false},
		{"CanAlsoBe is not sequential", RelationshipCanAlsoBe, false},
		{"PeerOf is not sequential", RelationshipPeerOf, false},
		{"MemberOf is not sequential", RelationshipMemberOf, false},
		{"HasMember is not sequential", RelationshipHasMember, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsSequential(); got != tt.want {
				t.Errorf("IsSequential() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipNature_IsDependency(t *testing.T) {
	tests := []struct {
		name  string
		value RelationshipNature
		want  bool
	}{
		{"Requires is dependency", RelationshipRequires, true},
		{"RequiredBy is dependency", RelationshipRequiredBy, true},
		{"ChildOf is not dependency", RelationshipChildOf, false},
		{"ParentOf is not dependency", RelationshipParentOf, false},
		{"CanPrecede is not dependency", RelationshipCanPrecede, false},
		{"CanFollow is not dependency", RelationshipCanFollow, false},
		{"CanAlsoBe is not dependency", RelationshipCanAlsoBe, false},
		{"PeerOf is not dependency", RelationshipPeerOf, false},
		{"MemberOf is not dependency", RelationshipMemberOf, false},
		{"HasMember is not dependency", RelationshipHasMember, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsDependency(); got != tt.want {
				t.Errorf("IsDependency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipNature_IsPeer(t *testing.T) {
	tests := []struct {
		name  string
		value RelationshipNature
		want  bool
	}{
		{"PeerOf is peer", RelationshipPeerOf, true},
		{"CanAlsoBe is peer", RelationshipCanAlsoBe, true},
		{"ChildOf is not peer", RelationshipChildOf, false},
		{"ParentOf is not peer", RelationshipParentOf, false},
		{"CanPrecede is not peer", RelationshipCanPrecede, false},
		{"CanFollow is not peer", RelationshipCanFollow, false},
		{"Requires is not peer", RelationshipRequires, false},
		{"RequiredBy is not peer", RelationshipRequiredBy, false},
		{"MemberOf is not peer", RelationshipMemberOf, false},
		{"HasMember is not peer", RelationshipHasMember, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsPeer(); got != tt.want {
				t.Errorf("IsPeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- ConsequenceScope ---

func TestConsequenceScope_String(t *testing.T) {
	tests := []struct {
		name  string
		value ConsequenceScope
		want  string
	}{
		{"Confidentiality", ScopeConfidentiality, "Confidentiality"},
		{"Integrity", ScopeIntegrity, "Integrity"},
		{"Availability", ScopeAvailability, "Availability"},
		{"AccessControl", ScopeAccessControl, "Access Control"},
		{"Accountability", ScopeAccountability, "Accountability"},
		{"Authentication", ScopeAuthentication, "Authentication"},
		{"Authorization", ScopeAuthorization, "Authorization"},
		{"NonRepudiation", ScopeNonRepudiation, "Non-Repudiation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("ConsequenceScope.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConsequenceScope_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value ConsequenceScope
		want  bool
	}{
		{"Confidentiality", ScopeConfidentiality, true},
		{"Integrity", ScopeIntegrity, true},
		{"Availability", ScopeAvailability, true},
		{"AccessControl", ScopeAccessControl, true},
		{"Accountability", ScopeAccountability, true},
		{"Authentication", ScopeAuthentication, true},
		{"Authorization", ScopeAuthorization, true},
		{"NonRepudiation", ScopeNonRepudiation, true},
		{"Invalid", ConsequenceScope("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("ConsequenceScope.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConsequenceScope(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ConsequenceScope
		wantErr bool
	}{
		{"Confidentiality", "Confidentiality", ScopeConfidentiality, false},
		{"Integrity", "Integrity", ScopeIntegrity, false},
		{"Availability", "Availability", ScopeAvailability, false},
		{"AccessControl", "Access Control", ScopeAccessControl, false},
		{"Accountability", "Accountability", ScopeAccountability, false},
		{"Authentication", "Authentication", ScopeAuthentication, false},
		{"Authorization", "Authorization", ScopeAuthorization, false},
		{"NonRepudiation", "Non-Repudiation", ScopeNonRepudiation, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConsequenceScope(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConsequenceScope(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseConsequenceScope(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllConsequenceScopeValues(t *testing.T) {
	values := AllConsequenceScopeValues()
	if len(values) != 8 {
		t.Errorf("AllConsequenceScopeValues() returned %d values, want 8", len(values))
	}
}

// --- ConsequenceImpact ---

func TestConsequenceImpact_String(t *testing.T) {
	tests := []struct {
		name  string
		value ConsequenceImpact
		want  string
	}{
		{"High", ImpactHigh, "High"},
		{"Medium", ImpactMedium, "Medium"},
		{"Low", ImpactLow, "Low"},
		{"Unknown", ImpactUnknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("ConsequenceImpact.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConsequenceImpact_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value ConsequenceImpact
		want  bool
	}{
		{"High", ImpactHigh, true},
		{"Medium", ImpactMedium, true},
		{"Low", ImpactLow, true},
		{"Unknown", ImpactUnknown, true},
		{"Invalid", ConsequenceImpact("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("ConsequenceImpact.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConsequenceImpact(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ConsequenceImpact
		wantErr bool
	}{
		{"High", "High", ImpactHigh, false},
		{"Medium", "Medium", ImpactMedium, false},
		{"Low", "Low", ImpactLow, false},
		{"Unknown", "Unknown", ImpactUnknown, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConsequenceImpact(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConsequenceImpact(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseConsequenceImpact(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllConsequenceImpactValues(t *testing.T) {
	values := AllConsequenceImpactValues()
	if len(values) != 4 {
		t.Errorf("AllConsequenceImpactValues() returned %d values, want 4", len(values))
	}
}

func TestImpactOrder(t *testing.T) {
	tests := []struct {
		name  string
		value ConsequenceImpact
		want  int
	}{
		{"High", ImpactHigh, 4},
		{"Medium", ImpactMedium, 3},
		{"Low", ImpactLow, 2},
		{"Unknown", ImpactUnknown, 1},
		{"UnknownValue", ConsequenceImpact("Foo"), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.ImpactOrder(); got != tt.want {
				t.Errorf("ImpactOrder() = %d, want %d", got, tt.want)
			}
		})
	}
}

// --- ViewType ---

func TestViewType_String(t *testing.T) {
	tests := []struct {
		name  string
		value ViewType
		want  string
	}{
		{"Graph", ViewTypeGraph, "Graph"},
		{"ExplicitSlice", ViewTypeExplicitSlice, "Explicit Slice"},
		{"ImplicitSlice", ViewTypeImplicitSlice, "Implicit Slice"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("ViewType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestViewType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value ViewType
		want  bool
	}{
		{"Graph", ViewTypeGraph, true},
		{"ExplicitSlice", ViewTypeExplicitSlice, true},
		{"ImplicitSlice", ViewTypeImplicitSlice, true},
		{"Invalid", ViewType("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("ViewType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseViewType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ViewType
		wantErr bool
	}{
		{"Graph", "Graph", ViewTypeGraph, false},
		{"ExplicitSlice", "Explicit Slice", ViewTypeExplicitSlice, false},
		{"ImplicitSlice", "Implicit Slice", ViewTypeImplicitSlice, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseViewType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseViewType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseViewType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllViewTypeValues(t *testing.T) {
	values := AllViewTypeValues()
	if len(values) != 3 {
		t.Errorf("AllViewTypeValues() returned %d values, want 3", len(values))
	}
}

// --- PlatformType ---

func TestPlatformType_String(t *testing.T) {
	tests := []struct {
		name  string
		value PlatformType
		want  string
	}{
		{"Language", PlatformLanguage, "Language"},
		{"OperatingSystem", PlatformOperatingSystem, "Operating System"},
		{"Architecture", PlatformArchitecture, "Architecture"},
		{"Technology", PlatformTechnology, "Technology"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("PlatformType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPlatformType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value PlatformType
		want  bool
	}{
		{"Language", PlatformLanguage, true},
		{"OperatingSystem", PlatformOperatingSystem, true},
		{"Architecture", PlatformArchitecture, true},
		{"Technology", PlatformTechnology, true},
		{"Invalid", PlatformType("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("PlatformType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePlatformType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    PlatformType
		wantErr bool
	}{
		{"Language", "Language", PlatformLanguage, false},
		{"OperatingSystem", "Operating System", PlatformOperatingSystem, false},
		{"Architecture", "Architecture", PlatformArchitecture, false},
		{"Technology", "Technology", PlatformTechnology, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePlatformType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePlatformType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePlatformType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllPlatformTypeValues(t *testing.T) {
	values := AllPlatformTypeValues()
	if len(values) != 4 {
		t.Errorf("AllPlatformTypeValues() returned %d values, want 4", len(values))
	}
}

// --- Prevalence ---

func TestPrevalence_String(t *testing.T) {
	tests := []struct {
		name  string
		value Prevalence
		want  string
	}{
		{"Often", PrevalenceOften, "Often"},
		{"Sometimes", PrevalenceSometimes, "Sometimes"},
		{"Rarely", PrevalenceRarely, "Rarely"},
		{"Undetermined", PrevalenceUndetermined, "Undetermined"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("Prevalence.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrevalence_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value Prevalence
		want  bool
	}{
		{"Often", PrevalenceOften, true},
		{"Sometimes", PrevalenceSometimes, true},
		{"Rarely", PrevalenceRarely, true},
		{"Undetermined", PrevalenceUndetermined, true},
		{"Invalid", Prevalence("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("Prevalence.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePrevalence(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Prevalence
		wantErr bool
	}{
		{"Often", "Often", PrevalenceOften, false},
		{"Sometimes", "Sometimes", PrevalenceSometimes, false},
		{"Rarely", "Rarely", PrevalenceRarely, false},
		{"Undetermined", "Undetermined", PrevalenceUndetermined, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePrevalence(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePrevalence(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePrevalence(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllPrevalenceValues(t *testing.T) {
	values := AllPrevalenceValues()
	if len(values) != 4 {
		t.Errorf("AllPrevalenceValues() returned %d values, want 4", len(values))
	}
}

// --- IntroductionPhase ---

func TestIntroductionPhase_String(t *testing.T) {
	tests := []struct {
		name  string
		value IntroductionPhase
		want  string
	}{
		{"ArchitectureAndDesign", PhaseArchitectureAndDesign, "Architecture and Design"},
		{"Implementation", PhaseImplementation, "Implementation"},
		{"BuildAndCompilation", PhaseBuildAndCompilation, "Build and Compilation"},
		{"Operation", PhaseOperation, "Operation"},
		{"SystemConfiguration", PhaseSystemConfiguration, "System Configuration"},
		{"Installation", PhaseInstallation, "Installation"},
		{"Policy", PhasePolicy, "Policy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("IntroductionPhase.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIntroductionPhase_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value IntroductionPhase
		want  bool
	}{
		{"ArchitectureAndDesign", PhaseArchitectureAndDesign, true},
		{"Implementation", PhaseImplementation, true},
		{"BuildAndCompilation", PhaseBuildAndCompilation, true},
		{"Operation", PhaseOperation, true},
		{"SystemConfiguration", PhaseSystemConfiguration, true},
		{"Installation", PhaseInstallation, true},
		{"Policy", PhasePolicy, true},
		{"Invalid", IntroductionPhase("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("IntroductionPhase.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIntroductionPhase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    IntroductionPhase
		wantErr bool
	}{
		{"ArchitectureAndDesign", "Architecture and Design", PhaseArchitectureAndDesign, false},
		{"Implementation", "Implementation", PhaseImplementation, false},
		{"BuildAndCompilation", "Build and Compilation", PhaseBuildAndCompilation, false},
		{"Operation", "Operation", PhaseOperation, false},
		{"SystemConfiguration", "System Configuration", PhaseSystemConfiguration, false},
		{"Installation", "Installation", PhaseInstallation, false},
		{"Policy", "Policy", PhasePolicy, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIntroductionPhase(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntroductionPhase(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseIntroductionPhase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllIntroductionPhaseValues(t *testing.T) {
	values := AllIntroductionPhaseValues()
	if len(values) != 7 {
		t.Errorf("AllIntroductionPhaseValues() returned %d values, want 7", len(values))
	}
}

// --- MitigationPhase ---

func TestMitigationPhase_String(t *testing.T) {
	tests := []struct {
		name  string
		value MitigationPhase
		want  string
	}{
		{"ArchitectureAndDesign", MitigationPhaseArchitectureAndDesign, "Architecture and Design"},
		{"BuildAndCompilation", MitigationPhaseBuildAndCompilation, "Build and Compilation"},
		{"Implementation", MitigationPhaseImplementation, "Implementation"},
		{"Operation", MitigationPhaseOperation, "Operation"},
		{"SystemConfiguration", MitigationPhaseSystemConfiguration, "System Configuration"},
		{"Installation", MitigationPhaseInstallation, "Installation"},
		{"Policy", MitigationPhasePolicy, "Policy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("MitigationPhase.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMitigationPhase_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value MitigationPhase
		want  bool
	}{
		{"ArchitectureAndDesign", MitigationPhaseArchitectureAndDesign, true},
		{"BuildAndCompilation", MitigationPhaseBuildAndCompilation, true},
		{"Implementation", MitigationPhaseImplementation, true},
		{"Operation", MitigationPhaseOperation, true},
		{"SystemConfiguration", MitigationPhaseSystemConfiguration, true},
		{"Installation", MitigationPhaseInstallation, true},
		{"Policy", MitigationPhasePolicy, true},
		{"Invalid", MitigationPhase("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("MitigationPhase.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMitigationPhase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    MitigationPhase
		wantErr bool
	}{
		{"ArchitectureAndDesign", "Architecture and Design", MitigationPhaseArchitectureAndDesign, false},
		{"BuildAndCompilation", "Build and Compilation", MitigationPhaseBuildAndCompilation, false},
		{"Implementation", "Implementation", MitigationPhaseImplementation, false},
		{"Operation", "Operation", MitigationPhaseOperation, false},
		{"SystemConfiguration", "System Configuration", MitigationPhaseSystemConfiguration, false},
		{"Installation", "Installation", MitigationPhaseInstallation, false},
		{"Policy", "Policy", MitigationPhasePolicy, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMitigationPhase(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMitigationPhase(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseMitigationPhase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllMitigationPhaseValues(t *testing.T) {
	values := AllMitigationPhaseValues()
	if len(values) != 7 {
		t.Errorf("AllMitigationPhaseValues() returned %d values, want 7", len(values))
	}
}

// --- Effectiveness ---

func TestEffectiveness_String(t *testing.T) {
	tests := []struct {
		name  string
		value Effectiveness
		want  string
	}{
		{"High", EffectivenessHigh, "High"},
		{"Moderate", EffectivenessModerate, "Moderate"},
		{"Limited", EffectivenessLimited, "Limited"},
		{"DefenseInDepth", EffectivenessDefenseInDepth, "Defense in Depth"},
		{"SOARPartial", EffectivenessSOARPartial, "SOAR Partial"},
		{"Unknown", EffectivenessUnknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("Effectiveness.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEffectiveness_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value Effectiveness
		want  bool
	}{
		{"High", EffectivenessHigh, true},
		{"Moderate", EffectivenessModerate, true},
		{"Limited", EffectivenessLimited, true},
		{"DefenseInDepth", EffectivenessDefenseInDepth, true},
		{"SOARPartial", EffectivenessSOARPartial, true},
		{"Unknown", EffectivenessUnknown, true},
		{"Invalid", Effectiveness("Invalid"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.want {
				t.Errorf("Effectiveness.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseEffectiveness(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Effectiveness
		wantErr bool
	}{
		{"High", "High", EffectivenessHigh, false},
		{"Moderate", "Moderate", EffectivenessModerate, false},
		{"Limited", "Limited", EffectivenessLimited, false},
		{"DefenseInDepth", "Defense in Depth", EffectivenessDefenseInDepth, false},
		{"SOARPartial", "SOAR Partial", EffectivenessSOARPartial, false},
		{"Unknown", "Unknown", EffectivenessUnknown, false},
		{"Invalid", "Invalid", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEffectiveness(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEffectiveness(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseEffectiveness(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllEffectivenessValues(t *testing.T) {
	values := AllEffectivenessValues()
	if len(values) != 6 {
		t.Errorf("AllEffectivenessValues() returned %d values, want 6", len(values))
	}
}
