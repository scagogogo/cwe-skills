package cwe

import (
	"strings"
	"testing"
)

func TestParseCWEID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"empty string", "", 0, true},
		{"pure number", "79", 79, false},
		{"standard format", "CWE-79", 79, false},
		{"lowercase format", "cwe-79", 79, false},
		{"mixed case format", "Cwe-79", 79, false},
		{"no hyphen", "CWE79", 79, false},
		{"leading zeros", "CWE-079", 79, false},
		{"leading zeros pure number", "079", 79, false},
		{"negative number", "-1", 0, true},
		{"zero", "0", 0, true},
		{"CWE-0", "CWE-0", 0, true},
		{"non-numeric", "abc", 0, true},
		{"whitespace padding", " CWE-79 ", 79, false},
		{"whitespace only", "   ", 0, true},
		{"large number", "999999", 999999, false},
		{"CWE with space", "CWE 79", 79, false},
		{"number one", "1", 1, false},
		{"CWE-1", "CWE-1", 1, false},
		{"overflow number in CWE format", "CWE-99999999999999999999", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCWEID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCWEID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseCWEID(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatCWEID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid pure number", "79", "CWE-79", false},
		{"valid standard format", "CWE-79", "CWE-79", false},
		{"valid lowercase", "cwe-79", "CWE-79", false},
		{"valid with leading zeros", "CWE-079", "CWE-79", false},
		{"empty string", "", "", true},
		{"invalid non-numeric", "abc", "", true},
		{"invalid zero", "0", "", true},
		{"valid large number", "1000", "CWE-1000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatCWEID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatCWEID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatCWEID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatCWEIDFromInt(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{"positive number", 79, "CWE-79"},
		{"large number", 1000, "CWE-1000"},
		{"one", 1, "CWE-1"},
		{"zero", 0, "CWE-0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatCWEIDFromInt(tt.input)
			if got != tt.want {
				t.Errorf("FormatCWEIDFromInt(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsCWEID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid standard", "CWE-79", true},
		{"valid pure number", "79", true},
		{"valid lowercase", "cwe-79", true},
		{"invalid non-numeric", "abc", false},
		{"empty string", "", false},
		{"invalid zero", "0", false},
		{"invalid negative", "-1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCWEID(tt.input)
			if got != tt.want {
				t.Errorf("IsCWEID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateCWEID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid CWE ID", "CWE-79", false},
		{"valid pure number", "89", false},
		{"empty string", "", true},
		{"invalid non-numeric", "abc", true},
		{"invalid zero", "0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCWEID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCWEID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestExtractCWEIDs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"no matches", "No CWE IDs here", []string{}},
		{"one match", "See CWE-79 for details", []string{"CWE-79"}},
		{"multiple matches", "See CWE-79 and CWE-89 for details", []string{"CWE-79", "CWE-89"}},
		{"empty string", "", []string{}},
		{"lowercase match", "see cwe-79 here", []string{"CWE-79"}},
		{"no hyphen match", "CWE79 is a thing", []string{"CWE-79"}},
		{"mixed formats", "CWE-79, cwe-89 and CWE125", []string{"CWE-79", "CWE-89", "CWE-125"}},
		{"zero ID excluded", "CWE-0 is invalid but CWE-1 is valid", []string{"CWE-1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractCWEIDs(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractCWEIDs(%q) = %v, want %v", tt.input, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ExtractCWEIDs(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestExtractFirstCWEID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"found", "See CWE-79 and CWE-89", "CWE-79"},
		{"not found", "No CWE IDs here", ""},
		{"empty string", "", ""},
		{"lowercase", "see cwe-89 here", "CWE-89"},
		{"zero ID returns empty", "CWE-0 is invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractFirstCWEID(tt.input)
			if got != tt.want {
				t.Errorf("ExtractFirstCWEID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCompareCWEIDs(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		want    int
		wantErr bool
	}{
		{"less than", "CWE-79", "CWE-89", -1, false},
		{"equal", "CWE-79", "CWE-79", 0, false},
		{"greater than", "CWE-89", "CWE-79", 1, false},
		{"first invalid", "abc", "CWE-79", 0, true},
		{"second invalid", "CWE-79", "abc", 0, true},
		{"both invalid", "abc", "xyz", 0, true},
		{"equal with pure number", "79", "CWE-79", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareCWEIDs(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareCWEIDs(%q, %q) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareCWEIDs(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCompareCWEIDsErrorMessage(t *testing.T) {
	_, err := CompareCWEIDs("abc", "CWE-79")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "第一个CWE ID无效") {
		t.Errorf("error message should mention first CWE ID invalid, got: %v", err)
	}

	_, err = CompareCWEIDs("CWE-79", "abc")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "第二个CWE ID无效") {
		t.Errorf("error message should mention second CWE ID invalid, got: %v", err)
	}
}
