package cwe

import "testing"

// ==================== CWE 测试 ====================

func TestNewCWE(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		cweName  string
	}{
		{name: "创建CWE-79", id: 79, cweName: "XSS"},
		{name: "创建CWE-89", id: 89, cweName: "SQL注入"},
		{name: "创建CWE-1", id: 1, cweName: "测试"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCWE(tt.id, tt.cweName)
			if got.ID != tt.id {
				t.Errorf("ID = %v, want %v", got.ID, tt.id)
			}
			if got.Name != tt.cweName {
				t.Errorf("Name = %v, want %v", got.Name, tt.cweName)
			}
			if got.CWEType != "weakness" {
				t.Errorf("CWEType = %v, want \"weakness\"", got.CWEType)
			}
		})
	}
}

func TestCWE_CWEID(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		expected string
	}{
		{name: "CWE-79", id: 79, expected: "CWE-79"},
		{name: "CWE-1000", id: 1000, expected: "CWE-1000"},
		{name: "CWE-1", id: 1, expected: "CWE-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{ID: tt.id}
			if got := c.CWEID(); got != tt.expected {
				t.Errorf("CWEID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsWeakness(t *testing.T) {
	tests := []struct {
		name     string
		cweType  string
		expected bool
	}{
		{name: "weakness类型", cweType: "weakness", expected: true},
		{name: "category类型", cweType: "category", expected: false},
		{name: "view类型", cweType: "view", expected: false},
		{name: "compound_element类型", cweType: "compound_element", expected: false},
		{name: "空字符串", cweType: "", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{CWEType: tt.cweType}
			if got := c.IsWeakness(); got != tt.expected {
				t.Errorf("IsWeakness() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsCategory(t *testing.T) {
	tests := []struct {
		name     string
		cweType  string
		expected bool
	}{
		{name: "category类型", cweType: "category", expected: true},
		{name: "weakness类型", cweType: "weakness", expected: false},
		{name: "view类型", cweType: "view", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{CWEType: tt.cweType}
			if got := c.IsCategory(); got != tt.expected {
				t.Errorf("IsCategory() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsView(t *testing.T) {
	tests := []struct {
		name     string
		cweType  string
		expected bool
	}{
		{name: "view类型", cweType: "view", expected: true},
		{name: "weakness类型", cweType: "weakness", expected: false},
		{name: "category类型", cweType: "category", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{CWEType: tt.cweType}
			if got := c.IsView(); got != tt.expected {
				t.Errorf("IsView() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsCompoundElement(t *testing.T) {
	tests := []struct {
		name     string
		cweType  string
		expected bool
	}{
		{name: "compound_element类型", cweType: "compound_element", expected: true},
		{name: "weakness类型", cweType: "weakness", expected: false},
		{name: "category类型", cweType: "category", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{CWEType: tt.cweType}
			if got := c.IsCompoundElement(); got != tt.expected {
				t.Errorf("IsCompoundElement() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsPillar(t *testing.T) {
	tests := []struct {
		name        string
		abstraction Abstraction
		expected    bool
	}{
		{name: "Pillar级别", abstraction: AbstractionPillar, expected: true},
		{name: "Class级别", abstraction: AbstractionClass, expected: false},
		{name: "Base级别", abstraction: AbstractionBase, expected: false},
		{name: "Variant级别", abstraction: AbstractionVariant, expected: false},
		{name: "空抽象", abstraction: Abstraction(""), expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Abstraction: tt.abstraction}
			if got := c.IsPillar(); got != tt.expected {
				t.Errorf("IsPillar() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsBase(t *testing.T) {
	tests := []struct {
		name        string
		abstraction Abstraction
		expected    bool
	}{
		{name: "Base级别", abstraction: AbstractionBase, expected: true},
		{name: "Pillar级别", abstraction: AbstractionPillar, expected: false},
		{name: "Variant级别", abstraction: AbstractionVariant, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Abstraction: tt.abstraction}
			if got := c.IsBase(); got != tt.expected {
				t.Errorf("IsBase() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsVariant(t *testing.T) {
	tests := []struct {
		name        string
		abstraction Abstraction
		expected    bool
	}{
		{name: "Variant级别", abstraction: AbstractionVariant, expected: true},
		{name: "Base级别", abstraction: AbstractionBase, expected: false},
		{name: "Pillar级别", abstraction: AbstractionPillar, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Abstraction: tt.abstraction}
			if got := c.IsVariant(); got != tt.expected {
				t.Errorf("IsVariant() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsChain(t *testing.T) {
	tests := []struct {
		name      string
		structure Structure
		expected  bool
	}{
		{name: "Chain结构", structure: StructureChain, expected: true},
		{name: "Simple结构", structure: StructureSimple, expected: false},
		{name: "Composite结构", structure: StructureComposite, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Structure: tt.structure}
			if got := c.IsChain(); got != tt.expected {
				t.Errorf("IsChain() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsComposite(t *testing.T) {
	tests := []struct {
		name      string
		structure Structure
		expected  bool
	}{
		{name: "Composite结构", structure: StructureComposite, expected: true},
		{name: "Simple结构", structure: StructureSimple, expected: false},
		{name: "Chain结构", structure: StructureChain, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Structure: tt.structure}
			if got := c.IsComposite(); got != tt.expected {
				t.Errorf("IsComposite() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsStable(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{name: "Stable状态", status: StatusStable, expected: true},
		{name: "Usable状态", status: StatusUsable, expected: false},
		{name: "Draft状态", status: StatusDraft, expected: false},
		{name: "Deprecated状态", status: StatusDeprecated, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Status: tt.status}
			if got := c.IsStable(); got != tt.expected {
				t.Errorf("IsStable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_IsDeprecated(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{name: "Deprecated状态", status: StatusDeprecated, expected: true},
		{name: "Stable状态", status: StatusStable, expected: false},
		{name: "Obsolete状态", status: StatusObsolete, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Status: tt.status}
			if got := c.IsDeprecated(); got != tt.expected {
				t.Errorf("IsDeprecated() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_GetParentIDs(t *testing.T) {
	tests := []struct {
		name          string
		relationships []Relationship
		expected      []int
	}{
		{
			name:          "空关系列表",
			relationships: nil,
			expected:      nil,
		},
		{
			name:          "无ChildOf关系",
			relationships: []Relationship{{Nature: RelationshipParentOf, CWEID: 79}},
			expected:      nil,
		},
		{
			name:          "单个ChildOf关系",
			relationships: []Relationship{{Nature: RelationshipChildOf, CWEID: 74}},
			expected:      []int{74},
		},
		{
			name: "多个ChildOf关系",
			relationships: []Relationship{
				{Nature: RelationshipChildOf, CWEID: 74},
				{Nature: RelationshipChildOf, CWEID: 664},
			},
			expected: []int{74, 664},
		},
		{
			name: "混合关系类型",
			relationships: []Relationship{
				{Nature: RelationshipChildOf, CWEID: 74},
				{Nature: RelationshipParentOf, CWEID: 83},
				{Nature: RelationshipChildOf, CWEID: 664},
				{Nature: RelationshipPeerOf, CWEID: 89},
			},
			expected: []int{74, 664},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Relationships: tt.relationships}
			got := c.GetParentIDs()
			if len(got) != len(tt.expected) {
				t.Errorf("GetParentIDs() = %v, want %v", got, tt.expected)
				return
			}
			for i, id := range got {
				if id != tt.expected[i] {
					t.Errorf("GetParentIDs()[%d] = %v, want %v", i, id, tt.expected[i])
				}
			}
		})
	}
}

func TestCWE_GetChildIDs(t *testing.T) {
	tests := []struct {
		name          string
		relationships []Relationship
		expected      []int
	}{
		{
			name:          "空关系列表",
			relationships: nil,
			expected:      nil,
		},
		{
			name:          "无ParentOf关系",
			relationships: []Relationship{{Nature: RelationshipChildOf, CWEID: 74}},
			expected:      nil,
		},
		{
			name:          "单个ParentOf关系",
			relationships: []Relationship{{Nature: RelationshipParentOf, CWEID: 79}},
			expected:      []int{79},
		},
		{
			name: "多个ParentOf关系",
			relationships: []Relationship{
				{Nature: RelationshipParentOf, CWEID: 79},
				{Nature: RelationshipParentOf, CWEID: 89},
			},
			expected: []int{79, 89},
		},
		{
			name: "混合关系类型",
			relationships: []Relationship{
				{Nature: RelationshipChildOf, CWEID: 74},
				{Nature: RelationshipParentOf, CWEID: 79},
				{Nature: RelationshipParentOf, CWEID: 89},
			},
			expected: []int{79, 89},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Relationships: tt.relationships}
			got := c.GetChildIDs()
			if len(got) != len(tt.expected) {
				t.Errorf("GetChildIDs() = %v, want %v", got, tt.expected)
				return
			}
			for i, id := range got {
				if id != tt.expected[i] {
					t.Errorf("GetChildIDs()[%d] = %v, want %v", i, id, tt.expected[i])
				}
			}
		})
	}
}

func TestCWE_GetPeerIDs(t *testing.T) {
	tests := []struct {
		name          string
		relationships []Relationship
		expected      []int
	}{
		{
			name:          "空关系列表",
			relationships: nil,
			expected:      nil,
		},
		{
			name:          "无对等关系",
			relationships: []Relationship{{Nature: RelationshipChildOf, CWEID: 74}},
			expected:      nil,
		},
		{
			name:          "单个PeerOf关系",
			relationships: []Relationship{{Nature: RelationshipPeerOf, CWEID: 89}},
			expected:      []int{89},
		},
		{
			name:          "单个CanAlsoBe关系",
			relationships: []Relationship{{Nature: RelationshipCanAlsoBe, CWEID: 100}},
			expected:      []int{100},
		},
		{
			name: "混合对等关系",
			relationships: []Relationship{
				{Nature: RelationshipPeerOf, CWEID: 89},
				{Nature: RelationshipCanAlsoBe, CWEID: 100},
			},
			expected: []int{89, 100},
		},
		{
			name: "混合多种关系",
			relationships: []Relationship{
				{Nature: RelationshipChildOf, CWEID: 74},
				{Nature: RelationshipPeerOf, CWEID: 89},
				{Nature: RelationshipCanAlsoBe, CWEID: 100},
				{Nature: RelationshipCanPrecede, CWEID: 200},
			},
			expected: []int{89, 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Relationships: tt.relationships}
			got := c.GetPeerIDs()
			if len(got) != len(tt.expected) {
				t.Errorf("GetPeerIDs() = %v, want %v", got, tt.expected)
				return
			}
			for i, id := range got {
				if id != tt.expected[i] {
					t.Errorf("GetPeerIDs()[%d] = %v, want %v", i, id, tt.expected[i])
				}
			}
		})
	}
}

func TestCWE_GetChainIDs(t *testing.T) {
	tests := []struct {
		name          string
		relationships []Relationship
		expected      []int
	}{
		{
			name:          "空关系列表",
			relationships: nil,
			expected:      nil,
		},
		{
			name:          "无链式关系",
			relationships: []Relationship{{Nature: RelationshipChildOf, CWEID: 74}},
			expected:      nil,
		},
		{
			name:          "单个CanPrecede关系",
			relationships: []Relationship{{Nature: RelationshipCanPrecede, CWEID: 680}},
			expected:      []int{680},
		},
		{
			name:          "单个CanFollow关系",
			relationships: []Relationship{{Nature: RelationshipCanFollow, CWEID: 190}},
			expected:      []int{190},
		},
		{
			name: "混合链式关系",
			relationships: []Relationship{
				{Nature: RelationshipCanPrecede, CWEID: 680},
				{Nature: RelationshipCanFollow, CWEID: 190},
			},
			expected: []int{680, 190},
		},
		{
			name: "混合多种关系",
			relationships: []Relationship{
				{Nature: RelationshipChildOf, CWEID: 74},
				{Nature: RelationshipCanPrecede, CWEID: 680},
				{Nature: RelationshipCanFollow, CWEID: 190},
				{Nature: RelationshipPeerOf, CWEID: 89},
			},
			expected: []int{680, 190},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{Relationships: tt.relationships}
			got := c.GetChainIDs()
			if len(got) != len(tt.expected) {
				t.Errorf("GetChainIDs() = %v, want %v", got, tt.expected)
				return
			}
			for i, id := range got {
				if id != tt.expected[i] {
					t.Errorf("GetChainIDs()[%d] = %v, want %v", i, id, tt.expected[i])
				}
			}
		})
	}
}

func TestCWE_HasConsequenceScope(t *testing.T) {
	tests := []struct {
		name        string
		consequences []Consequence
		scope       ConsequenceScope
		expected    bool
	}{
		{
			name:        "空后果列表",
			consequences: nil,
			scope:       ScopeConfidentiality,
			expected:    false,
		},
		{
			name:        "后果列表中无匹配范围",
			consequences: []Consequence{{Scopes: []ConsequenceScope{ScopeIntegrity}}},
			scope:       ScopeConfidentiality,
			expected:    false,
		},
		{
			name:        "后果列表中有匹配范围",
			consequences: []Consequence{{Scopes: []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity}}},
			scope:       ScopeConfidentiality,
			expected:    true,
		},
		{
			name: "多个后果中第二个匹配",
			consequences: []Consequence{
				{Scopes: []ConsequenceScope{ScopeIntegrity}},
				{Scopes: []ConsequenceScope{ScopeAvailability}},
			},
			scope:    ScopeAvailability,
			expected: true,
		},
		{
			name: "多个后果都不匹配",
			consequences: []Consequence{
				{Scopes: []ConsequenceScope{ScopeIntegrity}},
				{Scopes: []ConsequenceScope{ScopeAvailability}},
			},
			scope:    ScopeConfidentiality,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CWE{CommonConsequences: tt.consequences}
			if got := c.HasConsequenceScope(tt.scope); got != tt.expected {
				t.Errorf("HasConsequenceScope() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCWE_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cwe     CWE
		wantErr bool
	}{
		{
			name:    "有效的CWE",
			cwe:     CWE{ID: 79, Name: "XSS"},
			wantErr: false,
		},
		{
			name:    "ID为零",
			cwe:     CWE{ID: 0, Name: "XSS"},
			wantErr: true,
		},
		{
			name:    "ID为负数",
			cwe:     CWE{ID: -1, Name: "XSS"},
			wantErr: true,
		},
		{
			name:    "Name为空",
			cwe:     CWE{ID: 79, Name: ""},
			wantErr: true,
		},
		{
			name:    "ID和Name都无效",
			cwe:     CWE{ID: 0, Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cwe.Validate()
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

// ==================== Category 测试 ====================

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name string
		id   int
		cat  string
	}{
		{name: "创建类别-1", id: 1000, cat: "研究概念"},
		{name: "创建类别-2", id: 699, cat: "软件开发"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCategory(tt.id, tt.cat)
			if got.ID != tt.id {
				t.Errorf("ID = %v, want %v", got.ID, tt.id)
			}
			if got.Name != tt.cat {
				t.Errorf("Name = %v, want %v", got.Name, tt.cat)
			}
		})
	}
}

// ==================== View 测试 ====================

func TestNewView(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		viewName string
		viewType ViewType
	}{
		{name: "创建图类型视图", id: 1000, viewName: "研究概念", viewType: ViewTypeGraph},
		{name: "创建显式切片视图", id: 1340, viewName: "Top 25", viewType: ViewTypeExplicitSlice},
		{name: "创建隐式切片视图", id: 999, viewName: "隐式视图", viewType: ViewTypeImplicitSlice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewView(tt.id, tt.viewName, tt.viewType)
			if got.ID != tt.id {
				t.Errorf("ID = %v, want %v", got.ID, tt.id)
			}
			if got.Name != tt.viewName {
				t.Errorf("Name = %v, want %v", got.Name, tt.viewName)
			}
			if got.Type != tt.viewType {
				t.Errorf("Type = %v, want %v", got.Type, tt.viewType)
			}
		})
	}
}

// ==================== CompoundElement 测试 ====================

func TestNewCompoundElement(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		elemName  string
		structure Structure
	}{
		{name: "创建链式复合元素", id: 680, elemName: "整数溢出到缓冲区溢出", structure: StructureChain},
		{name: "创建复合复合元素", id: 352, elemName: "CSRF", structure: StructureComposite},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCompoundElement(tt.id, tt.elemName, tt.structure)
			if got.ID != tt.id {
				t.Errorf("ID = %v, want %v", got.ID, tt.id)
			}
			if got.Name != tt.elemName {
				t.Errorf("Name = %v, want %v", got.Name, tt.elemName)
			}
			if got.Structure != tt.structure {
				t.Errorf("Structure = %v, want %v", got.Structure, tt.structure)
			}
		})
	}
}

// ==================== 支持类型字段测试 ====================

func TestMitigation_Fields(t *testing.T) {
	m := Mitigation{
		Phase:         MitigationPhaseImplementation,
		Strategy:      "输入验证",
		Description:   "对所有输入进行验证",
		Effectiveness: EffectivenessHigh,
	}
	if m.Phase != MitigationPhaseImplementation {
		t.Errorf("Phase = %v, want %v", m.Phase, MitigationPhaseImplementation)
	}
	if m.Strategy != "输入验证" {
		t.Errorf("Strategy = %v, want 输入验证", m.Strategy)
	}
	if m.Description != "对所有输入进行验证" {
		t.Errorf("Description = %v, want 对所有输入进行验证", m.Description)
	}
	if m.Effectiveness != EffectivenessHigh {
		t.Errorf("Effectiveness = %v, want %v", m.Effectiveness, EffectivenessHigh)
	}
}

func TestDemonstrativeExample_Fields(t *testing.T) {
	d := DemonstrativeExample{
		IntroText: "以下代码示例展示了...",
		BodyText:  "代码内容",
	}
	if d.IntroText != "以下代码示例展示了..." {
		t.Errorf("IntroText = %v, want 以下代码示例展示了...", d.IntroText)
	}
	if d.BodyText != "代码内容" {
		t.Errorf("BodyText = %v, want 代码内容", d.BodyText)
	}
}

func TestObservedExample_Fields(t *testing.T) {
	o := ObservedExample{
		Reference:   "CVE-2021-1234",
		Description: "某产品中观察到的示例",
		Link:        "https://example.com",
	}
	if o.Reference != "CVE-2021-1234" {
		t.Errorf("Reference = %v, want CVE-2021-1234", o.Reference)
	}
	if o.Description != "某产品中观察到的示例" {
		t.Errorf("Description = %v, want 某产品中观察到的示例", o.Description)
	}
	if o.Link != "https://example.com" {
		t.Errorf("Link = %v, want https://example.com", o.Link)
	}
}

func TestReference_Fields(t *testing.T) {
	r := Reference{
		ID:     1,
		Author: "MITRE",
		Title:  "CWE参考",
		URL:    "https://cwe.mitre.org",
	}
	if r.ID != 1 {
		t.Errorf("ID = %v, want 1", r.ID)
	}
	if r.Author != "MITRE" {
		t.Errorf("Author = %v, want MITRE", r.Author)
	}
	if r.Title != "CWE参考" {
		t.Errorf("Title = %v, want CWE参考", r.Title)
	}
	if r.URL != "https://cwe.mitre.org" {
		t.Errorf("URL = %v, want https://cwe.mitre.org", r.URL)
	}
}

func TestApplicablePlatforms_Fields(t *testing.T) {
	ap := ApplicablePlatforms{
		Languages:        []PlatformEntry{{Name: "C", Prevalence: PrevalenceOften}},
		OperatingSystems: []PlatformEntry{{Name: "Linux", Prevalence: PrevalenceSometimes}},
		Architectures:    []PlatformEntry{{Name: "x86", Prevalence: PrevalenceOften}},
		Technologies:     []PlatformEntry{{Name: "Web", Prevalence: PrevalenceOften}},
	}
	if len(ap.Languages) != 1 || ap.Languages[0].Name != "C" {
		t.Errorf("Languages = %v, want [{Name:C}]", ap.Languages)
	}
	if len(ap.OperatingSystems) != 1 || ap.OperatingSystems[0].Name != "Linux" {
		t.Errorf("OperatingSystems = %v, want [{Name:Linux}]", ap.OperatingSystems)
	}
	if len(ap.Architectures) != 1 || ap.Architectures[0].Name != "x86" {
		t.Errorf("Architectures = %v, want [{Name:x86}]", ap.Architectures)
	}
	if len(ap.Technologies) != 1 || ap.Technologies[0].Name != "Web" {
		t.Errorf("Technologies = %v, want [{Name:Web}]", ap.Technologies)
	}
}

func TestPlatformEntry_Fields(t *testing.T) {
	pe := PlatformEntry{Name: "Java", Prevalence: PrevalenceOften}
	if pe.Name != "Java" {
		t.Errorf("Name = %v, want Java", pe.Name)
	}
	if pe.Prevalence != PrevalenceOften {
		t.Errorf("Prevalence = %v, want %v", pe.Prevalence, PrevalenceOften)
	}
}

func TestIntroduction_Fields(t *testing.T) {
	i := Introduction{
		Phase:       PhaseImplementation,
		Description: "在实现阶段引入",
	}
	if i.Phase != PhaseImplementation {
		t.Errorf("Phase = %v, want %v", i.Phase, PhaseImplementation)
	}
	if i.Description != "在实现阶段引入" {
		t.Errorf("Description = %v, want 在实现阶段引入", i.Description)
	}
}

func TestAlternateTerm_Fields(t *testing.T) {
	at := AlternateTerm{
		Term:        "XSS",
		Description: "跨站脚本攻击的简称",
	}
	if at.Term != "XSS" {
		t.Errorf("Term = %v, want XSS", at.Term)
	}
	if at.Description != "跨站脚本攻击的简称" {
		t.Errorf("Description = %v, want 跨站脚本攻击的简称", at.Description)
	}
}

func TestContentHistory_Fields(t *testing.T) {
	ch := ContentHistory{
		Submission: &HistoryEntry{
			Name:         "MITRE",
			Organization: "MITRE Corporation",
			Date:         "2020-01-01",
			Comment:      "初始提交",
		},
		Modifications: []HistoryEntry{
			{Name: "John", Date: "2021-06-15", Comment: "更新描述"},
		},
	}
	if ch.Submission.Name != "MITRE" {
		t.Errorf("Submission.Name = %v, want MITRE", ch.Submission.Name)
	}
	if ch.Submission.Organization != "MITRE Corporation" {
		t.Errorf("Submission.Organization = %v, want MITRE Corporation", ch.Submission.Organization)
	}
	if ch.Submission.Date != "2020-01-01" {
		t.Errorf("Submission.Date = %v, want 2020-01-01", ch.Submission.Date)
	}
	if ch.Submission.Comment != "初始提交" {
		t.Errorf("Submission.Comment = %v, want 初始提交", ch.Submission.Comment)
	}
	if len(ch.Modifications) != 1 {
		t.Errorf("len(Modifications) = %v, want 1", len(ch.Modifications))
	}
	if ch.Modifications[0].Name != "John" {
		t.Errorf("Modifications[0].Name = %v, want John", ch.Modifications[0].Name)
	}
}

func TestHistoryEntry_Fields(t *testing.T) {
	he := HistoryEntry{
		Name:         "Alice",
		Organization: "ACME Corp",
		Date:         "2022-03-10",
		Comment:      "修复拼写错误",
	}
	if he.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", he.Name)
	}
	if he.Organization != "ACME Corp" {
		t.Errorf("Organization = %v, want ACME Corp", he.Organization)
	}
	if he.Date != "2022-03-10" {
		t.Errorf("Date = %v, want 2022-03-10", he.Date)
	}
	if he.Comment != "修复拼写错误" {
		t.Errorf("Comment = %v, want 修复拼写错误", he.Comment)
	}
}

func TestViewMember_Fields(t *testing.T) {
	vm := ViewMember{
		CWEID:     79,
		ViewID:    1000,
		Direct:    true,
		Predicate: "has",
	}
	if vm.CWEID != 79 {
		t.Errorf("CWEID = %v, want 79", vm.CWEID)
	}
	if vm.ViewID != 1000 {
		t.Errorf("ViewID = %v, want 1000", vm.ViewID)
	}
	if vm.Direct != true {
		t.Errorf("Direct = %v, want true", vm.Direct)
	}
	if vm.Predicate != "has" {
		t.Errorf("Predicate = %v, want has", vm.Predicate)
	}
}

func TestViewMember_DirectFalse(t *testing.T) {
	vm := ViewMember{
		CWEID:  89,
		ViewID: 699,
		Direct: false,
	}
	if vm.Direct != false {
		t.Errorf("Direct = %v, want false", vm.Direct)
	}
}

// ==================== 综合测试 ====================

func TestCWE_ComprehensiveExample(t *testing.T) {
	cwe := NewCWE(79, "Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')")
	cwe.Abstraction = AbstractionBase
	cwe.Structure = StructureSimple
	cwe.Status = StatusStable
	cwe.Description = "软件未对输入进行中立化处理..."
	cwe.LikelihoodOfExploit = LikelihoodHigh
	cwe.CommonConsequences = []Consequence{
		{Scopes: []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity}, Impacts: []ConsequenceImpact{ImpactHigh}},
	}
	cwe.Relationships = []Relationship{
		{Nature: RelationshipChildOf, CWEID: 74},
		{Nature: RelationshipChildOf, CWEID: 664},
		{Nature: RelationshipParentOf, CWEID: 80},
		{Nature: RelationshipPeerOf, CWEID: 89},
		{Nature: RelationshipCanPrecede, CWEID: 352},
	}

	// 测试CWEID格式化
	if got := cwe.CWEID(); got != "CWE-79" {
		t.Errorf("CWEID() = %v, want CWE-79", got)
	}

	// 测试类型判断
	if !cwe.IsWeakness() {
		t.Error("IsWeakness() = false, want true")
	}
	if cwe.IsCategory() {
		t.Error("IsCategory() = true, want false")
	}

	// 测试抽象层级
	if !cwe.IsBase() {
		t.Error("IsBase() = false, want true")
	}
	if cwe.IsPillar() {
		t.Error("IsPillar() = true, want false")
	}

	// 测试结构
	if cwe.IsChain() {
		t.Error("IsChain() = true, want false")
	}
	if cwe.IsComposite() {
		t.Error("IsComposite() = true, want false")
	}

	// 测试状态
	if !cwe.IsStable() {
		t.Error("IsStable() = false, want true")
	}
	if cwe.IsDeprecated() {
		t.Error("IsDeprecated() = true, want false")
	}

	// 测试关系查询
	parents := cwe.GetParentIDs()
	if len(parents) != 2 || parents[0] != 74 || parents[1] != 664 {
		t.Errorf("GetParentIDs() = %v, want [74 664]", parents)
	}

	children := cwe.GetChildIDs()
	if len(children) != 1 || children[0] != 80 {
		t.Errorf("GetChildIDs() = %v, want [80]", children)
	}

	peers := cwe.GetPeerIDs()
	if len(peers) != 1 || peers[0] != 89 {
		t.Errorf("GetPeerIDs() = %v, want [89]", peers)
	}

	chains := cwe.GetChainIDs()
	if len(chains) != 1 || chains[0] != 352 {
		t.Errorf("GetChainIDs() = %v, want [352]", chains)
	}

	// 测试后果范围查询
	if !cwe.HasConsequenceScope(ScopeConfidentiality) {
		t.Error("HasConsequenceScope(Confidentiality) = false, want true")
	}
	if !cwe.HasConsequenceScope(ScopeIntegrity) {
		t.Error("HasConsequenceScope(Integrity) = false, want true")
	}
	if cwe.HasConsequenceScope(ScopeAvailability) {
		t.Error("HasConsequenceScope(Availability) = true, want false")
	}

	// 测试验证
	if err := cwe.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestCWE_ValidateIDFirst(t *testing.T) {
	// 当ID和Name都无效时，应该先报ID的错误
	cwe := CWE{ID: 0, Name: ""}
	err := cwe.Validate()
	if err == nil {
		t.Error("Validate() = nil, want error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Validate() should return *ValidationError, got %T", err)
	}
	if ve.Field != "ID" {
		t.Errorf("ValidationError.Field = %v, want ID", ve.Field)
	}
}

func TestCWE_ValidateNameEmpty(t *testing.T) {
	// ID有效但Name为空
	cwe := CWE{ID: 79, Name: ""}
	err := cwe.Validate()
	if err == nil {
		t.Error("Validate() = nil, want error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("Validate() should return *ValidationError, got %T", err)
	}
	if ve.Field != "Name" {
		t.Errorf("ValidationError.Field = %v, want Name", ve.Field)
	}
}