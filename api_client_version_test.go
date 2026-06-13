package cwe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIClient_GetVersion(t *testing.T) {
	tests := []struct {
		name       string
		response   interface{}
		statusCode int
		wantErr    bool
		checkFunc  func(*VersionResponse) bool
	}{
		{
			name: "success",
			response: map[string]interface{}{
				"Data": map[string]interface{}{
					"version":     "4.10",
					"releaseDate": "2024-02-29",
					"name":        "CWE v4.10",
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(v *VersionResponse) bool {
				return v.Version == "4.10" && v.ReleaseDate == "2024-02-29" && v.Name == "CWE v4.10"
			},
		},
		{
			name:       "server error",
			response:   nil,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "invalid JSON data",
			response: map[string]interface{}{
				"Data": "not a version object",
			},
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request path
				if r.URL.Path != "/version" {
					t.Errorf("expected path /version, got %q", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer srv.Close()

			client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
			defer client.Close()

			version, err := client.GetVersion(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkFunc != nil && !tt.checkFunc(version) {
				t.Errorf("checkFunc failed for version: %+v", version)
			}
		})
	}
}
