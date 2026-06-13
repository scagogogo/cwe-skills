package cwe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIClient_GetWeakness(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func(*CWE) bool
	}{
		{
			name: "success array response",
			id:   79,
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"id": 79, "name": "XSS", "description": "Cross-site Scripting"},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(c *CWE) bool {
				return c.ID == 79 && c.Name == "XSS" && c.CWEType == "weakness"
			},
		},
		{
			name: "success single object response",
			id:   79,
			response: map[string]interface{}{
				"Data": map[string]interface{}{"id": 79, "name": "XSS", "description": "Cross-site Scripting"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(c *CWE) bool {
				return c.ID == 79 && c.Name == "XSS" && c.CWEType == "weakness"
			},
		},
		{
			name:       "not found empty array",
			id:         79,
			response:   map[string]interface{}{"Data": []interface{}{}},
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &CWENotFoundError{},
		},
		{
			name:       "invalid ID zero",
			id:         0,
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &InvalidCWEIDError{},
		},
		{
			name:       "invalid ID negative",
			id:         -1,
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &InvalidCWEIDError{},
		},
		{
			name:       "server error",
			id:         79,
			response:   nil,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer srv.Close()

			client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
			defer client.Close()

			cwe, err := client.GetWeakness(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errType != nil {
					switch tt.errType.(type) {
					case *InvalidCWEIDError:
						if _, ok := err.(*InvalidCWEIDError); !ok {
							t.Errorf("expected InvalidCWEIDError, got %T: %v", err, err)
						}
					case *CWENotFoundError:
						if _, ok := err.(*CWENotFoundError); !ok {
							t.Errorf("expected CWENotFoundError, got %T: %v", err, err)
						}
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkFunc != nil && !tt.checkFunc(cwe) {
				t.Errorf("checkFunc failed for CWE: %+v", cwe)
			}
		})
	}
}

func TestAPIClient_GetCategory(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func(*Category) bool
	}{
		{
			name: "success",
			id:   1,
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"id": 1, "name": "Category1", "description": "A category"},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(c *Category) bool {
				return c.ID == 1 && c.Name == "Category1"
			},
		},
		{
			name: "success single object",
			id:   1,
			response: map[string]interface{}{
				"Data": map[string]interface{}{"id": 1, "name": "Category1", "description": "A category"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(c *Category) bool {
				return c.ID == 1 && c.Name == "Category1"
			},
		},
		{
			name:       "not found empty array",
			id:         999,
			response:   map[string]interface{}{"Data": []interface{}{}},
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &CWENotFoundError{},
		},
		{
			name:       "invalid ID",
			id:         -1,
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &InvalidCWEIDError{},
		},
		{
			name:       "server error",
			id:         1,
			response:   nil,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer srv.Close()

			client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
			defer client.Close()

			cat, err := client.GetCategory(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errType != nil {
					switch tt.errType.(type) {
					case *InvalidCWEIDError:
						if _, ok := err.(*InvalidCWEIDError); !ok {
							t.Errorf("expected InvalidCWEIDError, got %T: %v", err, err)
						}
					case *CWENotFoundError:
						if _, ok := err.(*CWENotFoundError); !ok {
							t.Errorf("expected CWENotFoundError, got %T: %v", err, err)
						}
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkFunc != nil && !tt.checkFunc(cat) {
				t.Errorf("checkFunc failed for Category: %+v", cat)
			}
		})
	}
}

func TestAPIClient_GetView(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func(*View) bool
	}{
		{
			name: "success",
			id:   1000,
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"id": 1000, "name": "Research Concepts", "description": "Research view"},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(v *View) bool {
				return v.ID == 1000 && v.Name == "Research Concepts"
			},
		},
		{
			name: "success single object",
			id:   1000,
			response: map[string]interface{}{
				"Data": map[string]interface{}{"id": 1000, "name": "Research Concepts", "description": "Research view"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(v *View) bool {
				return v.ID == 1000 && v.Name == "Research Concepts"
			},
		},
		{
			name:       "not found empty array",
			id:         9999,
			response:   map[string]interface{}{"Data": []interface{}{}},
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &CWENotFoundError{},
		},
		{
			name:       "invalid ID",
			id:         0,
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &InvalidCWEIDError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer srv.Close()

			client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
			defer client.Close()

			view, err := client.GetView(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errType != nil {
					switch tt.errType.(type) {
					case *InvalidCWEIDError:
						if _, ok := err.(*InvalidCWEIDError); !ok {
							t.Errorf("expected InvalidCWEIDError, got %T: %v", err, err)
						}
					case *CWENotFoundError:
						if _, ok := err.(*CWENotFoundError); !ok {
							t.Errorf("expected CWENotFoundError, got %T: %v", err, err)
						}
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkFunc != nil && !tt.checkFunc(view) {
				t.Errorf("checkFunc failed for View: %+v", view)
			}
		})
	}
}

func TestAPIClient_GetCWEs(t *testing.T) {
	tests := []struct {
		name       string
		ids        []int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func(map[string]*CWE) bool
	}{
		{
			name: "success",
			ids:  []int{79, 89},
			response: map[string]interface{}{
				"Data": map[string]interface{}{
					"79": map[string]interface{}{"id": 79, "name": "XSS", "description": "XSS weakness"},
					"89": map[string]interface{}{"id": 89, "name": "SQLi", "description": "SQL Injection"},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(m map[string]*CWE) bool {
				return len(m) == 2 && m["79"] != nil && m["89"] != nil &&
					m["79"].CWEType == "weakness" && m["89"].CWEType == "weakness"
			},
		},
		{
			name:       "empty IDs",
			ids:        []int{},
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(m map[string]*CWE) bool {
				return len(m) == 0
			},
		},
		{
			name:       "invalid ID in list",
			ids:        []int{79, -1},
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &InvalidCWEIDError{},
		},
		{
			name:       "invalid ID zero in list",
			ids:        []int{79, 0},
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &InvalidCWEIDError{},
		},
		{
			name:       "server error",
			ids:        []int{79},
			response:   nil,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer srv.Close()

			client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
			defer client.Close()

			result, err := client.GetCWEs(context.Background(), tt.ids)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errType != nil {
					switch tt.errType.(type) {
					case *InvalidCWEIDError:
						if _, ok := err.(*InvalidCWEIDError); !ok {
							t.Errorf("expected InvalidCWEIDError, got %T: %v", err, err)
						}
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkFunc != nil && !tt.checkFunc(result) {
				t.Errorf("checkFunc failed for result: %+v", result)
			}
		})
	}
}
