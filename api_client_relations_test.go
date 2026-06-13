package cwe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIClient_GetParents(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		viewID     []int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func([]Relationship) bool
	}{
		{
			name:   "success",
			id:     79,
			viewID: []int{},
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "ChildOf", "cweId": 74, "viewId": 1000},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(rels []Relationship) bool {
				return len(rels) == 1 && rels[0].CWEID == 74
			},
		},
		{
			name:   "with viewID",
			id:     79,
			viewID: []int{1000},
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "ChildOf", "cweId": 74, "viewId": 1000},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(rels []Relationship) bool {
				return len(rels) == 1
			},
		},
		{
			name:       "invalid ID",
			id:         -1,
			viewID:     []int{},
			response:   nil,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    &InvalidCWEIDError{},
		},
		{
			name:   "raw relation format",
			id:     79,
			viewID: []int{},
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "ChildOf", "cweId": 74},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(rels []Relationship) bool {
				return len(rels) == 1 && rels[0].CWEID == 74
			},
		},
		{
			name:       "server error",
			id:         79,
			viewID:     []int{},
			response:   nil,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify viewID query parameter when provided
				if len(tt.viewID) > 0 && tt.viewID[0] > 0 && tt.statusCode == http.StatusOK {
					viewParam := r.URL.Query().Get("view")
					if viewParam != "1000" {
						t.Errorf("expected view=1000 query param, got %q", viewParam)
					}
				}
				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer srv.Close()

			client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
			defer client.Close()

			var rels []Relationship
			var err error
			if len(tt.viewID) > 0 {
				rels, err = client.GetParents(context.Background(), tt.id, tt.viewID...)
			} else {
				rels, err = client.GetParents(context.Background(), tt.id)
			}

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
			if tt.checkFunc != nil && !tt.checkFunc(rels) {
				t.Errorf("checkFunc failed for relationships: %+v", rels)
			}
		})
	}
}

func TestAPIClient_GetChildren(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		viewID     []int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func([]Relationship) bool
	}{
		{
			name:   "success",
			id:     74,
			viewID: []int{},
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "ParentOf", "cweId": 79, "viewId": 1000},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(rels []Relationship) bool {
				return len(rels) == 1 && rels[0].CWEID == 79
			},
		},
		{
			name:   "with viewID",
			id:     74,
			viewID: []int{1000},
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "ParentOf", "cweId": 79, "viewId": 1000},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(rels []Relationship) bool {
				return len(rels) == 1
			},
		},
		{
			name:       "invalid ID",
			id:         0,
			viewID:     []int{},
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

			var rels []Relationship
			var err error
			if len(tt.viewID) > 0 {
				rels, err = client.GetChildren(context.Background(), tt.id, tt.viewID...)
			} else {
				rels, err = client.GetChildren(context.Background(), tt.id)
			}

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
			if tt.checkFunc != nil && !tt.checkFunc(rels) {
				t.Errorf("checkFunc failed for relationships: %+v", rels)
			}
		})
	}
}

func TestAPIClient_GetAncestors(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func([]Relationship) bool
	}{
		{
			name: "success",
			id:   79,
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "ChildOf", "cweId": 74},
					{"nature": "ChildOf", "cweId": 707},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(rels []Relationship) bool {
				return len(rels) == 2
			},
		},
		{
			name:       "invalid ID",
			id:         -5,
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

			rels, err := client.GetAncestors(context.Background(), tt.id)

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
			if tt.checkFunc != nil && !tt.checkFunc(rels) {
				t.Errorf("checkFunc failed for relationships: %+v", rels)
			}
		})
	}
}

func TestAPIClient_GetDescendants(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		response   interface{}
		statusCode int
		wantErr    bool
		errType    interface{}
		checkFunc  func([]Relationship) bool
	}{
		{
			name: "success",
			id:   74,
			response: map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "ParentOf", "cweId": 79},
					{"nature": "ParentOf", "cweId": 89},
				},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
			checkFunc: func(rels []Relationship) bool {
				return len(rels) == 2
			},
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

			rels, err := client.GetDescendants(context.Background(), tt.id)

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
			if tt.checkFunc != nil && !tt.checkFunc(rels) {
				t.Errorf("checkFunc failed for relationships: %+v", rels)
			}
		})
	}
}

func TestAPIClient_GetRelations_ParseFallback(t *testing.T) {
	// Test the fallback parsing path in getRelations when JSON doesn't
	// directly unmarshal to []Relationship
	t.Run("invalid nature uses raw string", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// Send data that won't parse as []Relationship but will as rawRelations
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Data": []map[string]interface{}{
					{"nature": "UnknownNature", "cweId": 100},
				},
			})
		}))
		defer srv.Close()

		client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
		defer client.Close()

		rels, err := client.GetParents(context.Background(), 79)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(rels) != 1 {
			t.Fatalf("expected 1 relationship, got %d", len(rels))
		}
		// The nature should be the raw string since ParseRelationshipNature fails
		if string(rels[0].Nature) != "UnknownNature" {
			t.Errorf("expected nature 'UnknownNature', got %q", rels[0].Nature)
		}
	})

	t.Run("completely unparseable data", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Data": "not an array or object",
			})
		}))
		defer srv.Close()

		client := NewAPIClient(WithAPIBaseURL(srv.URL), WithAPIRateLimit(100, 10))
		defer client.Close()

		_, err := client.GetParents(context.Background(), 79)
		if err == nil {
			t.Fatal("expected error for unparseable data")
		}
	})
}
