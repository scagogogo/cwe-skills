package cwe

import "testing"

func TestConsequence_HasScope(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []ConsequenceScope
		scope    ConsequenceScope
		expected bool
	}{
		{
			name:     "包含机密性范围",
			scopes:   []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity},
			scope:    ScopeConfidentiality,
			expected: true,
		},
		{
			name:     "包含完整性范围",
			scopes:   []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity},
			scope:    ScopeIntegrity,
			expected: true,
		},
		{
			name:     "不包含可用性范围",
			scopes:   []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity},
			scope:    ScopeAvailability,
			expected: false,
		},
		{
			name:     "空范围列表",
			scopes:   []ConsequenceScope{},
			scope:    ScopeConfidentiality,
			expected: false,
		},
		{
			name:     "单个范围匹配",
			scopes:   []ConsequenceScope{ScopeAvailability},
			scope:    ScopeAvailability,
			expected: true,
		},
		{
			name:     "访问控制范围",
			scopes:   []ConsequenceScope{ScopeAccessControl, ScopeAuthentication},
			scope:    ScopeAccessControl,
			expected: true,
		},
		{
			name:     "授权范围不在列表中",
			scopes:   []ConsequenceScope{ScopeAccessControl, ScopeAuthentication},
			scope:    ScopeAuthorization,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consequence{Scopes: tt.scopes}
			if got := c.HasScope(tt.scope); got != tt.expected {
				t.Errorf("HasScope() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConsequence_HasImpact(t *testing.T) {
	tests := []struct {
		name     string
		impacts  []ConsequenceImpact
		impact   ConsequenceImpact
		expected bool
	}{
		{
			name:     "包含高影响",
			impacts:  []ConsequenceImpact{ImpactHigh, ImpactMedium},
			impact:   ImpactHigh,
			expected: true,
		},
		{
			name:     "包含中等影响",
			impacts:  []ConsequenceImpact{ImpactHigh, ImpactMedium},
			impact:   ImpactMedium,
			expected: true,
		},
		{
			name:     "不包含低影响",
			impacts:  []ConsequenceImpact{ImpactHigh, ImpactMedium},
			impact:   ImpactLow,
			expected: false,
		},
		{
			name:     "空影响列表",
			impacts:  []ConsequenceImpact{},
			impact:   ImpactHigh,
			expected: false,
		},
		{
			name:     "未知影响",
			impacts:  []ConsequenceImpact{ImpactUnknown},
			impact:   ImpactUnknown,
			expected: true,
		},
		{
			name:     "nil影响列表",
			impacts:  nil,
			impact:   ImpactHigh,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consequence{Impacts: tt.impacts}
			if got := c.HasImpact(tt.impact); got != tt.expected {
				t.Errorf("HasImpact() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConsequence_MaxImpact(t *testing.T) {
	tests := []struct {
		name     string
		impacts  []ConsequenceImpact
		expected ConsequenceImpact
	}{
		{
			name:     "空影响列表返回Unknown",
			impacts:  []ConsequenceImpact{},
			expected: ImpactUnknown,
		},
		{
			name:     "nil影响列表返回Unknown",
			impacts:  nil,
			expected: ImpactUnknown,
		},
		{
			name:     "仅高影响",
			impacts:  []ConsequenceImpact{ImpactHigh},
			expected: ImpactHigh,
		},
		{
			name:     "高中低Unknown混合，High最高",
			impacts:  []ConsequenceImpact{ImpactLow, ImpactHigh, ImpactMedium, ImpactUnknown},
			expected: ImpactHigh,
		},
		{
			name:     "中低混合，Medium最高",
			impacts:  []ConsequenceImpact{ImpactLow, ImpactMedium},
			expected: ImpactMedium,
		},
		{
			name:     "仅低影响",
			impacts:  []ConsequenceImpact{ImpactLow},
			expected: ImpactLow,
		},
		{
			name:     "仅Unknown",
			impacts:  []ConsequenceImpact{ImpactUnknown},
			expected: ImpactUnknown,
		},
		{
			name:     "低和Unknown混合，Low最高",
			impacts:  []ConsequenceImpact{ImpactUnknown, ImpactLow},
			expected: ImpactLow,
		},
		{
			name:     "重复的最高值",
			impacts:  []ConsequenceImpact{ImpactHigh, ImpactHigh},
			expected: ImpactHigh,
		},
		{
			name:     "Medium和Low混合",
			impacts:  []ConsequenceImpact{ImpactLow, ImpactMedium, ImpactLow},
			expected: ImpactMedium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consequence{Impacts: tt.impacts}
			if got := c.MaxImpact(); got != tt.expected {
				t.Errorf("MaxImpact() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConsequence_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cons    Consequence
		wantErr bool
	}{
		{
			name:    "有效后果-有范围",
			cons:    Consequence{Scopes: []ConsequenceScope{ScopeConfidentiality}},
			wantErr: false,
		},
		{
			name:    "有效后果-多个范围",
			cons:    Consequence{Scopes: []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity, ScopeAvailability}},
			wantErr: false,
		},
		{
			name:    "无效后果-空范围列表",
			cons:    Consequence{Scopes: []ConsequenceScope{}},
			wantErr: true,
		},
		{
			name:    "无效后果-nil范围列表",
			cons:    Consequence{Scopes: nil},
			wantErr: true,
		},
		{
			name: "有效后果-带完整字段",
			cons: Consequence{
				Scopes:     []ConsequenceScope{ScopeConfidentiality},
				Impacts:    []ConsequenceImpact{ImpactHigh},
				Likelihood: LikelihoodHigh,
				Note:       "测试备注",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cons.Validate()
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