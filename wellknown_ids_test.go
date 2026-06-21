package cweskills

import (
	"testing"
)

func TestIsInTop25(t *testing.T) {
	tests := []struct {
		name  string
		cweID int
		want  bool
	}{
		{"CWE-79 XSS is in Top25", 79, true},
		{"CWE-89 SQLi is in Top25", 89, true},
		{"CWE-787 OOB Write is in Top25", 787, true},
		{"Non-member", 999999, false},
		{"Zero is not in Top25", 0, false},
		{"CWE-20 is in Top25", 20, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInTop25(tt.cweID); got != tt.want {
				t.Errorf("IsInTop25(%d) = %v, want %v", tt.cweID, got, tt.want)
			}
		})
	}
}

func TestIsInOWASPTop10(t *testing.T) {
	tests := []struct {
		name  string
		cweID int
		want  bool
	}{
		{"CWE-79 is in OWASP Top10", 79, true},
		{"CWE-918 SSRF is in OWASP Top10", 918, true},
		{"Non-member", 999999, false},
		{"CWE-22 is in OWASP Top10", 22, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInOWASPTop10(tt.cweID); got != tt.want {
				t.Errorf("IsInOWASPTop10(%d) = %v, want %v", tt.cweID, got, tt.want)
			}
		})
	}
}

func TestIsInSANSTop25(t *testing.T) {
	tests := []struct {
		name  string
		cweID int
		want  bool
	}{
		{"CWE-89 is in SANS Top25", 89, true},
		{"CWE-79 is in SANS Top25", 79, true},
		{"Non-member", 999999, false},
		{"CWE-190 is in SANS Top25", 190, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInSANSTop25(tt.cweID); got != tt.want {
				t.Errorf("IsInSANSTop25(%d) = %v, want %v", tt.cweID, got, tt.want)
			}
		})
	}
}

func TestGetOWASPCategory(t *testing.T) {
	tests := []struct {
		name  string
		cweID int
		want  string
	}{
		{"CWE-79 belongs to Injection", 79, "A03:2021-Injection"},
		{"CWE-918 belongs to SSRF", 918, "A10:2021-Server-Side Request Forgery (SSRF)"},
		{"Non-member returns empty", 999999, ""},
		{"CWE-22 belongs to Broken Access Control", 22, "A01:2021-Broken Access Control"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetOWASPCategory(tt.cweID)
			if got != tt.want {
				t.Errorf("GetOWASPCategory(%d) = %q, want %q", tt.cweID, got, tt.want)
			}
		})
	}
}

func TestGetOWASPCategories(t *testing.T) {
	tests := []struct {
		name       string
		cweID      int
		wantMinLen int
		wantEmpty  bool
	}{
		{"CWE-306 belongs to multiple categories", 306, 2, false},
		{"Non-member returns empty", 999999, 0, true},
		{"CWE-79 belongs to at least one", 79, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetOWASPCategories(tt.cweID)
			if tt.wantEmpty && len(got) != 0 {
				t.Errorf("GetOWASPCategories(%d) = %v, want empty", tt.cweID, got)
			}
			if !tt.wantEmpty && len(got) < tt.wantMinLen {
				t.Errorf("GetOWASPCategories(%d) = %v, want at least %d categories", tt.cweID, got, tt.wantMinLen)
			}
		})
	}
}

func TestIsInWellKnownView(t *testing.T) {
	tests := []struct {
		name    string
		viewID  int
		want    bool
	}{
		{"Research Concepts (1000)", CWEViewResearchConcepts, true},
		{"Development Concepts (699)", CWEViewDevelopmentConcepts, true},
		{"Hardware Design (1199)", CWEViewHardwareDesign, true},
		{"CWE Cross Section (888)", CWEViewCWECrossSection, true},
		{"Comprehensive Dictionary (1400)", CWEViewComprehensiveDictionary, true},
		{"Unknown view (999)", 999, false},
		{"Zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInWellKnownView(tt.viewID); got != tt.want {
				t.Errorf("IsInWellKnownView(%d) = %v, want %v", tt.viewID, got, tt.want)
			}
		})
	}
}

func TestCWETop25Has25Entries(t *testing.T) {
	if len(CWETop25) != 25 {
		t.Errorf("CWETop25 has %d entries, want 25", len(CWETop25))
	}
}

func TestOWASPTop10Has10Categories(t *testing.T) {
	if len(OWASPTop10) != 10 {
		t.Errorf("OWASPTop10 has %d categories, want 10", len(OWASPTop10))
	}
}

func TestSANSTop25Has25Entries(t *testing.T) {
	if len(SANSTop25) != 25 {
		t.Errorf("SANSTop25 has %d entries, want 25", len(SANSTop25))
	}
}
