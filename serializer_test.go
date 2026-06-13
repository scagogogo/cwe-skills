package cwe

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

// ==================== MarshalJSON 测试 ====================

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		cwe     *CWE
		wantErr bool
	}{
		{
			name: "success",
			cwe: &CWE{
				ID:                  79,
				Name:                "XSS",
				Abstraction:         AbstractionBase,
				Structure:           StructureSimple,
				Status:              StatusStable,
				Description:         "Improper Neutralization of Input",
				ExtendedDescription: "Extended description",
				LikelihoodOfExploit: LikelihoodHigh,
				URL:                 "https://cwe.mitre.org/data/definitions/79.html",
				CWEType:             "weakness",
			},
			wantErr: false,
		},
		{
			name:    "nil CWE",
			cwe:     nil,
			wantErr: true,
		},
		{
			name: "minimal CWE",
			cwe: &CWE{
				ID:          1,
				Name:        "Test",
				Description: "desc",
				CWEType:     "weakness",
			},
			wantErr: false,
		},
		{
			name: "CWE with relationships",
			cwe: &CWE{
				ID:          79,
				Name:        "XSS",
				Description: "desc",
				CWEType:     "weakness",
				Relationships: []Relationship{
					{Nature: RelationshipChildOf, CWEID: 74, ViewID: 1000},
					{Nature: RelationshipParentOf, CWEID: 80},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalJSON(tt.cwe)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(data) == 0 {
					t.Error("MarshalJSON() returned empty data")
				}
				// Verify it's valid JSON
				var m map[string]interface{}
				if err := json.Unmarshal(data, &m); err != nil {
					t.Errorf("MarshalJSON() produced invalid JSON: %v", err)
				}
			}
		})
	}
}

// ==================== UnmarshalJSON 测试 ====================

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		wantID  int
		wantName string
	}{
		{
			name: "success",
			data: func() []byte {
				d, _ := MarshalJSON(&CWE{ID: 79, Name: "XSS", Description: "desc", CWEType: "weakness"})
				return d
			}(),
			wantErr:  false,
			wantID:   79,
			wantName: "XSS",
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			data:    []byte("{invalid json}"),
			wantErr: true,
		},
		{
			name: "valid JSON with all fields",
			data: func() []byte {
				d, _ := MarshalJSON(&CWE{
					ID:                  89,
					Name:                "SQL Injection",
					Abstraction:         AbstractionBase,
					Structure:           StructureSimple,
					Status:              StatusStable,
					Description:         "SQL injection description",
					ExtendedDescription: "Extended info",
					LikelihoodOfExploit: LikelihoodHigh,
					CWEType:             "weakness",
				})
				return d
			}(),
			wantErr:  false,
			wantID:   89,
			wantName: "SQL Injection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwe, err := UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cwe.ID != tt.wantID {
					t.Errorf("UnmarshalJSON() ID = %v, want %v", cwe.ID, tt.wantID)
				}
				if cwe.Name != tt.wantName {
					t.Errorf("UnmarshalJSON() Name = %v, want %v", cwe.Name, tt.wantName)
				}
			}
		})
	}
}

// ==================== MarshalJSONList 测试 ====================

func TestMarshalJSONList(t *testing.T) {
	tests := []struct {
		name    string
		cwes    []*CWE
		wantErr bool
	}{
		{
			name: "success",
			cwes: []*CWE{
				{ID: 79, Name: "XSS", Description: "desc", CWEType: "weakness"},
				{ID: 89, Name: "SQLi", Description: "desc", CWEType: "weakness"},
			},
			wantErr: false,
		},
		{
			name:    "nil list",
			cwes:    nil,
			wantErr: false,
		},
		{
			name:    "empty list",
			cwes:    []*CWE{},
			wantErr: false,
		},
		{
			name: "single element",
			cwes: []*CWE{
				{ID: 79, Name: "XSS", Description: "desc", CWEType: "weakness"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalJSONList(tt.cwes)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSONList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.cwes == nil {
					if string(data) != "[]" {
						t.Errorf("MarshalJSONList(nil) = %s, want []", string(data))
					}
				} else {
					if len(data) == 0 {
						t.Error("MarshalJSONList() returned empty data")
					}
					var result []*CWE
					if err := json.Unmarshal(data, &result); err != nil {
						t.Errorf("MarshalJSONList() produced invalid JSON: %v", err)
					}
				}
			}
		})
	}
}

// ==================== UnmarshalJSONList 测试 ====================

func TestUnmarshalJSONList(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantErr  bool
		wantLen  int
	}{
		{
			name: "success",
			data: func() []byte {
				d, _ := MarshalJSONList([]*CWE{
					{ID: 79, Name: "XSS", Description: "desc", CWEType: "weakness"},
					{ID: 89, Name: "SQLi", Description: "desc", CWEType: "weakness"},
				})
				return d
			}(),
			wantErr: false,
			wantLen: 2,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			data:    []byte("[invalid]"),
			wantErr: true,
		},
		{
			name:    "empty JSON array",
			data:    []byte("[]"),
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwes, err := UnmarshalJSONList(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSONList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(cwes) != tt.wantLen {
					t.Errorf("UnmarshalJSONList() len = %v, want %v", len(cwes), tt.wantLen)
				}
			}
		})
	}
}

// ==================== MarshalXML 测试 ====================

func TestMarshalXML(t *testing.T) {
	tests := []struct {
		name    string
		cwe     *CWE
		wantErr bool
	}{
		{
			name: "success",
			cwe: &CWE{
				ID:                  79,
				Name:                "XSS",
				Abstraction:         AbstractionBase,
				Structure:           StructureSimple,
				Status:              StatusStable,
				Description:         "Improper Neutralization of Input",
				LikelihoodOfExploit: LikelihoodHigh,
				URL:                 "https://cwe.mitre.org/data/definitions/79.html",
				CWEType:             "weakness",
			},
			wantErr: false,
		},
		{
			name:    "nil CWE",
			cwe:     nil,
			wantErr: true,
		},
		{
			name: "CWE with relationships",
			cwe: &CWE{
				ID:          79,
				Name:        "XSS",
				Description: "desc",
				CWEType:     "weakness",
				Relationships: []Relationship{
					{Nature: RelationshipChildOf, CWEID: 74},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalXML(tt.cwe)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(data) == 0 {
					t.Error("MarshalXML() returned empty data")
				}
				// Verify it has XML header
				if string(data[:5]) != "<?xml" {
					t.Errorf("MarshalXML() should start with XML header, got: %s", string(data[:5]))
				}
			}
		})
	}
}

// ==================== UnmarshalXML 测试 ====================

func TestUnmarshalXML(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantErr  bool
		wantID   int
		wantName string
	}{
		{
			name: "success",
			data: func() []byte {
				d, _ := MarshalXML(&CWE{ID: 79, Name: "XSS", Description: "desc", CWEType: "weakness"})
				return d
			}(),
			wantErr:  false,
			wantID:   79,
			wantName: "XSS",
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "invalid XML",
			data:    []byte("<invalid><xml>"),
			wantErr: true,
		},
		{
			name: "valid XML with all fields",
			data: func() []byte {
				d, _ := MarshalXML(&CWE{
					ID:                  89,
					Name:                "SQLi",
					Abstraction:         AbstractionBase,
					Structure:           StructureSimple,
					Status:              StatusStable,
					Description:         "SQL injection",
					ExtendedDescription: "Extended",
					LikelihoodOfExploit: LikelihoodHigh,
					CWEType:             "weakness",
				})
				return d
			}(),
			wantErr:  false,
			wantID:   89,
			wantName: "SQLi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwe, err := UnmarshalXML(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cwe.ID != tt.wantID {
					t.Errorf("UnmarshalXML() ID = %v, want %v", cwe.ID, tt.wantID)
				}
				if cwe.Name != tt.wantName {
					t.Errorf("UnmarshalXML() Name = %v, want %v", cwe.Name, tt.wantName)
				}
				// UnmarshalXML always sets CWEType to "weakness" (via fromSafeCWE)
				if cwe.CWEType != "weakness" {
					t.Errorf("UnmarshalXML() CWEType = %v, want weakness", cwe.CWEType)
				}
			}
		})
	}
}

// ==================== MarshalCSV 测试 ====================

func TestMarshalCSV(t *testing.T) {
	tests := []struct {
		name    string
		cwes    []*CWE
		wantErr bool
		wantLen int // minimum expected length
	}{
		{
			name: "success",
			cwes: []*CWE{
				{ID: 79, Name: "XSS", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "desc", LikelihoodOfExploit: LikelihoodHigh},
				{ID: 89, Name: "SQLi", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "desc2", LikelihoodOfExploit: LikelihoodMedium},
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "nil list",
			cwes:    nil,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "empty list",
			cwes:    []*CWE{},
			wantErr: false,
			wantLen: 1, // still has header
		},
		{
			name: "single element",
			cwes: []*CWE{
				{ID: 79, Name: "XSS", Description: "desc"},
			},
			wantErr: false,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalCSV(tt.cwes)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.cwes == nil {
					if len(data) != 0 {
						t.Errorf("MarshalCSV(nil) should return empty bytes, got %d bytes", len(data))
					}
				} else {
					if len(data) < tt.wantLen {
						t.Errorf("MarshalCSV() returned too few bytes: %d", len(data))
					}
				}
			}
		})
	}
}

// ==================== UnmarshalCSV 测试 ====================

func TestUnmarshalCSV(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		wantLen int
	}{
		{
			name: "success",
			data: func() []byte {
				d, _ := MarshalCSV([]*CWE{
					{ID: 79, Name: "XSS", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "desc", LikelihoodOfExploit: LikelihoodHigh},
				})
				return d
			}(),
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "invalid CSV",
			data:    []byte("\"unclosed quote"),
			wantErr: true,
		},
		{
			name:    "rows with too few columns",
			data:    []byte("ID,Name\n79\n"),
			wantErr: false,
			wantLen: 0, // "79" is single-field row, skipped since len(record) < 2
		},
		{
			name:    "header with only one column",
			data:    []byte("ID\n79\n"),
			wantErr: true,
		},
		{
			name:    "CSV with non-numeric ID",
			data:    []byte("ID,Name\nabc,Test\n"),
			wantErr: false,
			wantLen: 0, // rows with invalid IDs are skipped
		},
		{
			name:    "CSV with multiple rows some invalid",
			data:    []byte("ID,Name\nabc,Invalid\n89,SQLi\n"),
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "round-trip with all fields",
			data: func() []byte {
				d, _ := MarshalCSV([]*CWE{
					{ID: 79, Name: "XSS", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "Cross-site Scripting", LikelihoodOfExploit: LikelihoodHigh},
					{ID: 89, Name: "SQLi", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "SQL Injection", LikelihoodOfExploit: LikelihoodMedium},
				})
				return d
			}(),
			wantErr: false,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwes, err := UnmarshalCSV(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(cwes) != tt.wantLen {
					t.Errorf("UnmarshalCSV() len = %v, want %v", len(cwes), tt.wantLen)
				}
			}
		})
	}
}

// ==================== ExportCSV 测试 ====================

func TestSerializerRegistry_ExportCSV(t *testing.T) {
	tests := []struct {
		name       string
		register   func(*Registry)
		wantErr    bool
		wantLenMin int
	}{
		{
			name: "success",
			register: func(r *Registry) {
				_ = r.Register(&CWE{ID: 79, Name: "XSS", Abstraction: AbstractionBase, Description: "desc", CWEType: "weakness"})
				_ = r.Register(&CWE{ID: 89, Name: "SQLi", Abstraction: AbstractionBase, Description: "desc2", CWEType: "weakness"})
			},
			wantErr:    false,
			wantLenMin: 1,
		},
		{
			name:       "empty registry",
			register:   func(r *Registry) {},
			wantErr:    false,
			wantLenMin: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistry()
			tt.register(r)
			data, err := r.ExportCSV()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(data) < tt.wantLenMin {
				t.Errorf("ExportCSV() returned too few bytes: %d", len(data))
			}
		})
	}
}

// ==================== JSON Round-trip 测试 ====================

func TestJSON_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		cwe  *CWE
	}{
		{
			name: "basic CWE",
			cwe: &CWE{
				ID:          79,
				Name:        "XSS",
				Description: "Improper Neutralization of Input",
				CWEType:     "weakness",
			},
		},
		{
			name: "full CWE",
			cwe: &CWE{
				ID:                  89,
				Name:                "SQL Injection",
				Abstraction:         AbstractionBase,
				Structure:           StructureSimple,
				Status:              StatusStable,
				Description:         "SQL injection description",
				ExtendedDescription: "Extended info",
				LikelihoodOfExploit: LikelihoodHigh,
				URL:                 "https://cwe.mitre.org/data/definitions/89.html",
				CWEType:             "weakness",
				Relationships: []Relationship{
					{Nature: RelationshipChildOf, CWEID: 74, ViewID: 1000},
					{Nature: RelationshipParentOf, CWEID: 560},
				},
			},
		},
		{
			name: "CWE with consequences",
			cwe: &CWE{
				ID:          79,
				Name:        "XSS",
				Description: "desc",
				CWEType:     "weakness",
				CommonConsequences: []Consequence{
					{Scopes: []ConsequenceScope{ScopeConfidentiality, ScopeIntegrity}, Impacts: []ConsequenceImpact{ImpactHigh}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalJSON(tt.cwe)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			got, err := UnmarshalJSON(data)
			if err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}

			// Verify key fields match
			if got.ID != tt.cwe.ID {
				t.Errorf("Round-trip ID mismatch: got %d, want %d", got.ID, tt.cwe.ID)
			}
			if got.Name != tt.cwe.Name {
				t.Errorf("Round-trip Name mismatch: got %q, want %q", got.Name, tt.cwe.Name)
			}
			if got.Description != tt.cwe.Description {
				t.Errorf("Round-trip Description mismatch: got %q, want %q", got.Description, tt.cwe.Description)
			}
			if got.Abstraction != tt.cwe.Abstraction {
				t.Errorf("Round-trip Abstraction mismatch: got %q, want %q", got.Abstraction, tt.cwe.Abstraction)
			}
			if got.Structure != tt.cwe.Structure {
				t.Errorf("Round-trip Structure mismatch: got %q, want %q", got.Structure, tt.cwe.Structure)
			}
			if got.Status != tt.cwe.Status {
				t.Errorf("Round-trip Status mismatch: got %q, want %q", got.Status, tt.cwe.Status)
			}
			if got.LikelihoodOfExploit != tt.cwe.LikelihoodOfExploit {
				t.Errorf("Round-trip LikelihoodOfExploit mismatch: got %q, want %q", got.LikelihoodOfExploit, tt.cwe.LikelihoodOfExploit)
			}
			if len(got.Relationships) != len(tt.cwe.Relationships) {
				t.Errorf("Round-trip Relationships len mismatch: got %d, want %d", len(got.Relationships), len(tt.cwe.Relationships))
			}
		})
	}
}

// ==================== JSON List Round-trip 测试 ====================

func TestJSONList_RoundTrip(t *testing.T) {
	cwes := []*CWE{
		{ID: 79, Name: "XSS", Description: "desc1", CWEType: "weakness"},
		{ID: 89, Name: "SQLi", Description: "desc2", CWEType: "weakness"},
		{ID: 119, Name: "Buffer Overflow", Description: "desc3", CWEType: "weakness"},
	}

	data, err := MarshalJSONList(cwes)
	if err != nil {
		t.Fatalf("MarshalJSONList() error = %v", err)
	}

	got, err := UnmarshalJSONList(data)
	if err != nil {
		t.Fatalf("UnmarshalJSONList() error = %v", err)
	}

	if len(got) != len(cwes) {
		t.Fatalf("Round-trip len mismatch: got %d, want %d", len(got), len(cwes))
	}

	for i, c := range got {
		if c.ID != cwes[i].ID {
			t.Errorf("Round-trip cwes[%d].ID = %d, want %d", i, c.ID, cwes[i].ID)
		}
		if c.Name != cwes[i].Name {
			t.Errorf("Round-trip cwes[%d].Name = %q, want %q", i, c.Name, cwes[i].Name)
		}
	}
}

// ==================== XML Round-trip 测试 ====================

func TestXML_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		cwe  *CWE
	}{
		{
			name: "basic CWE",
			cwe: &CWE{
				ID:          79,
				Name:        "XSS",
				Description: "Improper Neutralization of Input",
				CWEType:     "weakness",
			},
		},
		{
			name: "full CWE",
			cwe: &CWE{
				ID:                  89,
				Name:                "SQL Injection",
				Abstraction:         AbstractionBase,
				Structure:           StructureSimple,
				Status:              StatusStable,
				Description:         "SQL injection description",
				ExtendedDescription: "Extended info",
				LikelihoodOfExploit: LikelihoodHigh,
				URL:                 "https://cwe.mitre.org/data/definitions/89.html",
				CWEType:             "weakness",
				Relationships: []Relationship{
					{Nature: RelationshipChildOf, CWEID: 74, ViewID: 1000},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalXML(tt.cwe)
			if err != nil {
				t.Fatalf("MarshalXML() error = %v", err)
			}

			got, err := UnmarshalXML(data)
			if err != nil {
				t.Fatalf("UnmarshalXML() error = %v", err)
			}

			if got.ID != tt.cwe.ID {
				t.Errorf("Round-trip ID mismatch: got %d, want %d", got.ID, tt.cwe.ID)
			}
			if got.Name != tt.cwe.Name {
				t.Errorf("Round-trip Name mismatch: got %q, want %q", got.Name, tt.cwe.Name)
			}
			if got.Description != tt.cwe.Description {
				t.Errorf("Round-trip Description mismatch: got %q, want %q", got.Description, tt.cwe.Description)
			}
			if got.Abstraction != tt.cwe.Abstraction {
				t.Errorf("Round-trip Abstraction mismatch: got %q, want %q", got.Abstraction, tt.cwe.Abstraction)
			}
			if got.Structure != tt.cwe.Structure {
				t.Errorf("Round-trip Structure mismatch: got %q, want %q", got.Structure, tt.cwe.Structure)
			}
			if got.Status != tt.cwe.Status {
				t.Errorf("Round-trip Status mismatch: got %q, want %q", got.Status, tt.cwe.Status)
			}
			if got.LikelihoodOfExploit != tt.cwe.LikelihoodOfExploit {
				t.Errorf("Round-trip LikelihoodOfExploit mismatch: got %q, want %q", got.LikelihoodOfExploit, tt.cwe.LikelihoodOfExploit)
			}
			if got.URL != tt.cwe.URL {
				t.Errorf("Round-trip URL mismatch: got %q, want %q", got.URL, tt.cwe.URL)
			}
			// Note: CWEType is always set to "weakness" by UnmarshalXML (fromSafeCWE)
			if got.CWEType != "weakness" {
				t.Errorf("Round-trip CWEType = %q, want %q", got.CWEType, "weakness")
			}
			if len(got.Relationships) != len(tt.cwe.Relationships) {
				t.Errorf("Round-trip Relationships len mismatch: got %d, want %d", len(got.Relationships), len(tt.cwe.Relationships))
			}
		})
	}
}

// ==================== CSV Round-trip 测试 ====================

func TestCSV_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		cwes []*CWE
	}{
		{
			name: "single CWE with all CSV fields",
			cwes: []*CWE{
				{ID: 79, Name: "XSS", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "Cross-site Scripting", LikelihoodOfExploit: LikelihoodHigh},
			},
		},
		{
			name: "multiple CWEs",
			cwes: []*CWE{
				{ID: 79, Name: "XSS", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "Cross-site Scripting", LikelihoodOfExploit: LikelihoodHigh},
				{ID: 89, Name: "SQLi", Abstraction: AbstractionBase, Structure: StructureSimple, Status: StatusStable, Description: "SQL Injection", LikelihoodOfExploit: LikelihoodMedium},
			},
		},
		{
			name: "CWE with minimal fields",
			cwes: []*CWE{
				{ID: 100, Name: "Test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalCSV(tt.cwes)
			if err != nil {
				t.Fatalf("MarshalCSV() error = %v", err)
			}

			got, err := UnmarshalCSV(data)
			if err != nil {
				t.Fatalf("UnmarshalCSV() error = %v", err)
			}

			if len(got) != len(tt.cwes) {
				t.Fatalf("Round-trip len mismatch: got %d, want %d", len(got), len(tt.cwes))
			}

			for i, c := range got {
				if c.ID != tt.cwes[i].ID {
					t.Errorf("Round-trip cwes[%d].ID = %d, want %d", i, c.ID, tt.cwes[i].ID)
				}
				if c.Name != tt.cwes[i].Name {
					t.Errorf("Round-trip cwes[%d].Name = %q, want %q", i, c.Name, tt.cwes[i].Name)
				}
				// CSV only preserves a subset of fields
				if c.Abstraction != tt.cwes[i].Abstraction {
					t.Errorf("Round-trip cwes[%d].Abstraction = %q, want %q", i, c.Abstraction, tt.cwes[i].Abstraction)
				}
				if c.Structure != tt.cwes[i].Structure {
					t.Errorf("Round-trip cwes[%d].Structure = %q, want %q", i, c.Structure, tt.cwes[i].Structure)
				}
				if c.Status != tt.cwes[i].Status {
					t.Errorf("Round-trip cwes[%d].Status = %q, want %q", i, c.Status, tt.cwes[i].Status)
				}
				if c.Description != tt.cwes[i].Description {
					t.Errorf("Round-trip cwes[%d].Description = %q, want %q", i, c.Description, tt.cwes[i].Description)
				}
				if c.LikelihoodOfExploit != tt.cwes[i].LikelihoodOfExploit {
					t.Errorf("Round-trip cwes[%d].LikelihoodOfExploit = %q, want %q", i, c.LikelihoodOfExploit, tt.cwes[i].LikelihoodOfExploit)
				}
				// CWEType is always set to "weakness" by UnmarshalCSV
				if c.CWEType != "weakness" {
					t.Errorf("Round-trip cwes[%d].CWEType = %q, want %q", i, c.CWEType, "weakness")
				}
			}
		})
	}
}

// ==================== parseCSVInt 测试 ====================

func TestParseCSVInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "valid positive integer",
			input:   "79",
			want:    79,
			wantErr: false,
		},
		{
			name:    "valid large integer",
			input:   "1000",
			want:    1000,
			wantErr: false,
		},
		{
			name:    "valid single digit",
			input:   "1",
			want:    1,
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid chars - letters",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid chars - mixed",
			input:   "7a9",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid chars - special",
			input:   "79!",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid chars - negative sign",
			input:   "-1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid chars - space",
			input:   "7 9",
			want:    0,
			wantErr: true,
		},
		{
			name:    "zero",
			input:   "0",
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCSVInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCSVInt(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseCSVInt(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ==================== toSafeCWE / fromSafeCWE 测试 ====================

func TestSafeCWE_Conversion(t *testing.T) {
	original := &CWE{
		ID:                  79,
		Name:                "XSS",
		Abstraction:         AbstractionBase,
		Structure:           StructureSimple,
		Status:              StatusStable,
		Description:         "Cross-site Scripting",
		ExtendedDescription: "Extended",
		LikelihoodOfExploit: LikelihoodHigh,
		URL:                 "https://cwe.mitre.org/data/definitions/79.html",
		Relationships: []Relationship{
			{Nature: RelationshipChildOf, CWEID: 74},
		},
	}

	safe := toSafeCWE(original)
	if safe.ID != original.ID {
		t.Errorf("toSafeCWE ID = %d, want %d", safe.ID, original.ID)
	}
	if safe.Name != original.Name {
		t.Errorf("toSafeCWE Name = %q, want %q", safe.Name, original.Name)
	}
	if safe.Description != original.Description {
		t.Errorf("toSafeCWE Description = %q, want %q", safe.Description, original.Description)
	}
	if len(safe.Relationships) != len(original.Relationships) {
		t.Errorf("toSafeCWE Relationships len = %d, want %d", len(safe.Relationships), len(original.Relationships))
	}

	recovered := fromSafeCWE(safe)
	if recovered.ID != original.ID {
		t.Errorf("fromSafeCWE ID = %d, want %d", recovered.ID, original.ID)
	}
	if recovered.Name != original.Name {
		t.Errorf("fromSafeCWE Name = %q, want %q", recovered.Name, original.Name)
	}
	if recovered.CWEType != "weakness" {
		t.Errorf("fromSafeCWE CWEType = %q, want %q", recovered.CWEType, "weakness")
	}
}

// ==================== UnmarshalCSV with partial columns ====================

func TestUnmarshalCSV_PartialColumns(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		wantLen      int
		wantID       int
		wantName     string
		wantAbstr    Abstraction
		wantStruct   Structure
		wantStatus   Status
		wantDesc     string
		wantLikeli   LikelihoodOfExploit
	}{
		{
			name: "only ID and Name columns",
			data: []byte("ID,Name\n79,XSS\n"),
			wantLen:    1,
			wantID:     79,
			wantName:   "XSS",
			wantAbstr:  "",
			wantStruct: "",
			wantStatus: "",
			wantDesc:   "",
			wantLikeli: "",
		},
		{
			name: "ID, Name, Abstraction columns",
			data: []byte("ID,Name,Abstraction\n79,XSS,Base\n"),
			wantLen:    1,
			wantID:     79,
			wantName:   "XSS",
			wantAbstr:  AbstractionBase,
			wantStruct: "",
			wantStatus: "",
			wantDesc:   "",
			wantLikeli: "",
		},
		{
			name: "all seven columns",
			data: []byte("ID,Name,Abstraction,Structure,Status,Description,LikelihoodOfExploit\n89,SQLi,Base,Simple,Stable,SQL Injection,High\n"),
			wantLen:    1,
			wantID:     89,
			wantName:   "SQLi",
			wantAbstr:  AbstractionBase,
			wantStruct: StructureSimple,
			wantStatus: StatusStable,
			wantDesc:   "SQL Injection",
			wantLikeli: LikelihoodHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwes, err := UnmarshalCSV(tt.data)
			if err != nil {
				t.Fatalf("UnmarshalCSV() error = %v", err)
			}
			if len(cwes) != tt.wantLen {
				t.Fatalf("UnmarshalCSV() len = %d, want %d", len(cwes), tt.wantLen)
			}
			c := cwes[0]
			if c.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", c.ID, tt.wantID)
			}
			if c.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", c.Name, tt.wantName)
			}
			if c.Abstraction != tt.wantAbstr {
				t.Errorf("Abstraction = %q, want %q", c.Abstraction, tt.wantAbstr)
			}
			if c.Structure != tt.wantStruct {
				t.Errorf("Structure = %q, want %q", c.Structure, tt.wantStruct)
			}
			if c.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", c.Status, tt.wantStatus)
			}
			if c.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", c.Description, tt.wantDesc)
			}
			if c.LikelihoodOfExploit != tt.wantLikeli {
				t.Errorf("LikelihoodOfExploit = %q, want %q", c.LikelihoodOfExploit, tt.wantLikeli)
			}
		})
	}
}

// ==================== XML specific structure tests ====================

func TestSafeCWE_XMLTags(t *testing.T) {
	cwe := &CWE{
		ID:          79,
		Name:        "XSS",
		Description: "desc",
		CWEType:     "weakness",
		Relationships: []Relationship{
			{Nature: RelationshipChildOf, CWEID: 74},
			{Nature: RelationshipParentOf, CWEID: 80},
		},
	}

	data, err := MarshalXML(cwe)
	if err != nil {
		t.Fatalf("MarshalXML() error = %v", err)
	}

	// Verify the XML contains expected elements
	xmlStr := string(data)
	if !contains(xmlStr, "<CWE") {
		t.Error("XML should contain <CWE element")
	}
	if !contains(xmlStr, "<Name>XSS</Name>") {
		t.Error("XML should contain <Name>XSS</Name>")
	}
	if !contains(xmlStr, "<Relationships>") {
		t.Error("XML should contain <Relationships>")
	}
}

// contains checks if a string contains a substring (simple helper)
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || containsSubstr(s, sub))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// ==================== UnmarshalCSV with row having fewer than 2 columns ====================

func TestUnmarshalCSV_RowWithOneColumn(t *testing.T) {
	// A row with only 1 column should be skipped (len < 2)
	data := []byte("ID,Name\n79\n89,SQLi\n")
	cwes, err := UnmarshalCSV(data)
	if err != nil {
		t.Fatalf("UnmarshalCSV() error = %v", err)
	}
	// Only the row with 2 columns should be parsed
	// But actually, CSV reader might split "79" as a single field row
	// len(record) < 2 means it gets skipped
	if len(cwes) != 1 {
		t.Errorf("UnmarshalCSV() len = %d, want 1", len(cwes))
	}
	if len(cwes) > 0 && cwes[0].ID != 89 {
		t.Errorf("UnmarshalCSV()[0].ID = %d, want 89", cwes[0].ID)
	}
}

// ==================== XML header check ====================

func TestMarshalXML_IncludesHeader(t *testing.T) {
	cwe := &CWE{ID: 79, Name: "XSS", Description: "desc", CWEType: "weakness"}
	data, err := MarshalXML(cwe)
	if err != nil {
		t.Fatalf("MarshalXML() error = %v", err)
	}
	if string(data)[:len(xml.Header)] != xml.Header {
		t.Errorf("MarshalXML() should include XML header, got: %s", string(data[:len(xml.Header)]))
	}
}

// ==================== XML unmarshal with missing optional fields ====================

func TestUnmarshalXML_MinimalXML(t *testing.T) {
	// Only required fields in XML
	xmlData := []byte(`<CWE ID="79"><Name>XSS</Name><Description>desc</Description></CWE>`)
	cwe, err := UnmarshalXML(xmlData)
	if err != nil {
		t.Fatalf("UnmarshalXML() error = %v", err)
	}
	if cwe.ID != 79 {
		t.Errorf("ID = %d, want 79", cwe.ID)
	}
	if cwe.Name != "XSS" {
		t.Errorf("Name = %q, want %q", cwe.Name, "XSS")
	}
	if cwe.Description != "desc" {
		t.Errorf("Description = %q, want %q", cwe.Description, "desc")
	}
}
