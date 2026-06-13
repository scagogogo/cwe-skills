package cwe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ==================== Helper: create test API client backed by httptest ====================

// newTestAPIClient creates an APIClient that points to the given httptest server.
func newTestAPIClient(server *httptest.Server) *APIClient {
	client := NewAPIClient(WithAPIBaseURL(server.URL))
	// Remove the rate limiter for tests to avoid delays
	client.SetRateLimiter(nil)
	client.GetHTTPClient().SetMaxRetries(0)
	return client
}

// ==================== BasicFetcher 测试 ====================

func TestNewBasicFetcher(t *testing.T) {
	tests := []struct {
		name       string
		client     *APIClient
		wantNilCli bool
	}{
		{
			name:       "with nil client uses default",
			client:     nil,
			wantNilCli: false, // should create a default client
		},
		{
			name:       "with provided client",
			client:     NewAPIClient(),
			wantNilCli: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewBasicFetcher(tt.client)
			if f == nil {
				t.Error("NewBasicFetcher() returned nil")
			}
			if f.client == nil {
				t.Error("BasicFetcher.client should not be nil")
			}
		})
	}
}

func TestBasicFetcher_Fetch(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		handler    http.HandlerFunc
		wantErr    bool
		wantID     int
		wantName   string
	}{
		{
			name: "success",
			id:   79,
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := struct {
					Data json.RawMessage `json:"Data"`
				}{
					Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"Cross-site Scripting","cwe_type":"weakness"}]`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			wantErr:  false,
			wantID:   79,
			wantName: "XSS",
		},
		{
			name: "error from API",
			id:   79,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "not found - empty array",
			id:   999,
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := struct {
					Data json.RawMessage `json:"Data"`
				}{
					Data: json.RawMessage(`[]`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			wantErr: true,
		},
		{
			name: "invalid ID",
			id:   -1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Should not be called
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			fetcher := NewBasicFetcher(client)

			cwe, err := fetcher.Fetch(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cwe.ID != tt.wantID {
					t.Errorf("Fetch() ID = %v, want %v", cwe.ID, tt.wantID)
				}
				if cwe.Name != tt.wantName {
					t.Errorf("Fetch() Name = %v, want %v", cwe.Name, tt.wantName)
				}
			}
		})
	}
}

func TestBasicFetcher_FetchCategory(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		handler  http.HandlerFunc
		wantErr  bool
		wantID   int
		wantName string
	}{
		{
			name: "success",
			id:   1000,
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := struct {
					Data json.RawMessage `json:"Data"`
				}{
					Data: json.RawMessage(`[{"id":1000,"name":"Research Concepts","description":"desc","status":"Stable"}]`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			wantErr:  false,
			wantID:   1000,
			wantName: "Research Concepts",
		},
		{
			name: "error from API",
			id:   1000,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "invalid ID",
			id:   -1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			fetcher := NewBasicFetcher(client)

			cat, err := fetcher.FetchCategory(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cat.ID != tt.wantID {
					t.Errorf("FetchCategory() ID = %v, want %v", cat.ID, tt.wantID)
				}
				if cat.Name != tt.wantName {
					t.Errorf("FetchCategory() Name = %v, want %v", cat.Name, tt.wantName)
				}
			}
		})
	}
}

func TestBasicFetcher_FetchView(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		handler  http.HandlerFunc
		wantErr  bool
		wantID   int
		wantName string
	}{
		{
			name: "success",
			id:   1000,
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := struct {
					Data json.RawMessage `json:"Data"`
				}{
					Data: json.RawMessage(`[{"id":1000,"name":"Research Concepts","description":"desc","type":"Graph","status":"Stable"}]`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			wantErr:  false,
			wantID:   1000,
			wantName: "Research Concepts",
		},
		{
			name: "error from API",
			id:   1000,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "invalid ID",
			id:   -1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			fetcher := NewBasicFetcher(client)

			view, err := fetcher.FetchView(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchView() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if view.ID != tt.wantID {
					t.Errorf("FetchView() ID = %v, want %v", view.ID, tt.wantID)
				}
				if view.Name != tt.wantName {
					t.Errorf("FetchView() Name = %v, want %v", view.Name, tt.wantName)
				}
			}
		})
	}
}

func TestBasicFetcher_FetchWithRelations(t *testing.T) {
	tests := []struct {
		name              string
		id                int
		handler           http.HandlerFunc
		wantErr           bool
		wantID            int
		wantRelCount      int
	}{
		{
			name: "success with parents and children",
			id:   79,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/79":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/parents":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ParentOf","cweId":80,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			wantErr:      false,
			wantID:       79,
			wantRelCount: 2,
		},
		{
			name: "error on main fetch",
			id:   79,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "error on parents still succeeds",
			id:   79,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/79":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/parents":
					w.WriteHeader(http.StatusInternalServerError)
				case "/cwe/79/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ParentOf","cweId":80,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			wantErr:      false,
			wantID:       79,
			wantRelCount: 1, // only child relation since parents errored
		},
		{
			name: "error on children still succeeds",
			id:   79,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/79":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/parents":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/children":
					w.WriteHeader(http.StatusInternalServerError)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			wantErr:      false,
			wantID:       79,
			wantRelCount: 1, // only parent relation since children errored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			fetcher := NewBasicFetcher(client)

			cwe, err := fetcher.FetchWithRelations(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWithRelations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cwe.ID != tt.wantID {
					t.Errorf("FetchWithRelations() ID = %v, want %v", cwe.ID, tt.wantID)
				}
				if len(cwe.Relationships) != tt.wantRelCount {
					t.Errorf("FetchWithRelations() Relationships len = %v, want %v", len(cwe.Relationships), tt.wantRelCount)
				}
				// Verify relationship natures are set correctly
				for _, rel := range cwe.Relationships {
					if rel.Nature != RelationshipChildOf && rel.Nature != RelationshipParentOf {
						t.Errorf("Unexpected relationship nature: %v", rel.Nature)
					}
				}
			}
		})
	}
}

func TestBasicFetcher_FetchWithRelations_NatureOverridden(t *testing.T) {
	// Verify that FetchWithRelations overrides nature: parents become ChildOf, children become ParentOf
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/children":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ParentOf","cweId":80}]`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	fetcher := NewBasicFetcher(client)

	cwe, err := fetcher.FetchWithRelations(context.Background(), 79)
	if err != nil {
		t.Fatalf("FetchWithRelations() error = %v", err)
	}

	if len(cwe.Relationships) != 2 {
		t.Fatalf("Expected 2 relationships, got %d", len(cwe.Relationships))
	}

	// Parent relation should have Nature = ChildOf
	foundChildOf := false
	foundParentOf := false
	for _, rel := range cwe.Relationships {
		if rel.Nature == RelationshipChildOf && rel.CWEID == 74 {
			foundChildOf = true
		}
		if rel.Nature == RelationshipParentOf && rel.CWEID == 80 {
			foundParentOf = true
		}
	}
	if !foundChildOf {
		t.Error("Expected ChildOf relationship with CWEID=74")
	}
	if !foundParentOf {
		t.Error("Expected ParentOf relationship with CWEID=80")
	}
}

// ==================== MultipleFetcher 测试 ====================

func TestNewMultipleFetcher(t *testing.T) {
	tests := []struct {
		name   string
		client *APIClient
	}{
		{
			name:   "with nil client uses default",
			client: nil,
		},
		{
			name:   "with provided client",
			client: NewAPIClient(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewMultipleFetcher(tt.client)
			if f == nil {
				t.Error("NewMultipleFetcher() returned nil")
			}
			if f.client == nil {
				t.Error("MultipleFetcher.client should not be nil")
			}
		})
	}
}

func TestMultipleFetcher_FetchMultiple(t *testing.T) {
	tests := []struct {
		name     string
		ids      []int
		handler  http.HandlerFunc
		wantErr  bool
		wantLen  int
	}{
		{
			name: "success",
			ids:  []int{79, 89},
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := struct {
					Data json.RawMessage `json:"Data"`
				}{
					Data: json.RawMessage(`{"79":{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"},"89":{"id":89,"name":"SQLi","description":"desc","cwe_type":"weakness"}}`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "error from API",
			ids:  []int{79},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "empty ids",
			ids:  []int{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "invalid id in list",
			ids:  []int{-1},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			fetcher := NewMultipleFetcher(client)

			result, err := fetcher.FetchMultiple(context.Background(), tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchMultiple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(result) != tt.wantLen {
					t.Errorf("FetchMultiple() len = %v, want %v", len(result), tt.wantLen)
				}
			}
		})
	}
}

func TestMultipleFetcher_FetchMultipleToRegistry(t *testing.T) {
	tests := []struct {
		name     string
		ids      []int
		registry *Registry
		handler  http.HandlerFunc
		wantErr  bool
		wantSize int
	}{
		{
			name: "success",
			ids:  []int{79, 89},
			registry: NewRegistry(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := struct {
					Data json.RawMessage `json:"Data"`
				}{
					Data: json.RawMessage(`{"79":{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"},"89":{"id":89,"name":"SQLi","description":"desc","cwe_type":"weakness"}}`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			},
			wantErr:  false,
			wantSize: 2,
		},
		{
			name:     "nil registry error",
			ids:      []int{79},
			registry: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true,
		},
		{
			name: "API error",
			ids:  []int{79},
			registry: NewRegistry(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:  true,
			wantSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			fetcher := NewMultipleFetcher(client)

			err := fetcher.FetchMultipleToRegistry(context.Background(), tt.ids, tt.registry)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchMultipleToRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.registry != nil {
				if tt.registry.Size() != tt.wantSize {
					t.Errorf("FetchMultipleToRegistry() registry size = %v, want %v", tt.registry.Size(), tt.wantSize)
				}
			}
		})
	}
}

// ==================== TreeFetcher 测试 ====================

func TestNewTreeFetcher(t *testing.T) {
	tests := []struct {
		name     string
		client   *APIClient
		registry *Registry
		maxDepth int
		wantDepth int
	}{
		{
			name:     "with nil client and registry uses defaults",
			client:   nil,
			registry: nil,
			maxDepth: 0, // should default to 10
			wantDepth: 10,
		},
		{
			name:     "with provided client and registry",
			client:   NewAPIClient(),
			registry: NewRegistry(),
			maxDepth: 5,
			wantDepth: 5,
		},
		{
			name:     "negative maxDepth defaults to 10",
			client:   nil,
			registry: nil,
			maxDepth: -1,
			wantDepth: 10,
		},
		{
			name:     "maxDepth 0 defaults to 10",
			client:   nil,
			registry: nil,
			maxDepth: 0,
			wantDepth: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewTreeFetcher(tt.client, tt.registry, tt.maxDepth)
			if f == nil {
				t.Error("NewTreeFetcher() returned nil")
			}
			if f.client == nil {
				t.Error("TreeFetcher.client should not be nil")
			}
			if f.registry == nil {
				t.Error("TreeFetcher.registry should not be nil")
			}
			if f.maxDepth != tt.wantDepth {
				t.Errorf("TreeFetcher.maxDepth = %v, want %v", f.maxDepth, tt.wantDepth)
			}
		})
	}
}

func TestTreeFetcher_GetRegistry(t *testing.T) {
	reg := NewRegistry()
	f := NewTreeFetcher(nil, reg, 5)
	if f.GetRegistry() != reg {
		t.Error("GetRegistry() should return the same registry instance")
	}
}

func TestTreeFetcher_FetchWithAncestors(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		maxDepth int
		handler  http.HandlerFunc
		wantErr  bool
		wantSize int // expected registry size
	}{
		{
			name:     "success - fetches self and parent",
			id:       79,
			maxDepth: 5,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/79":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/parents":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/weakness/74":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/74/parents":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					// Return empty for any unexpected path
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				}
			},
			wantErr:  false,
			wantSize: 2, // CWE-79 and CWE-74
		},
		{
			name:     "max depth limit",
			id:       79,
			maxDepth: 1, // only depth 0, can fetch CWE-79 but not recurse further
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/79":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/parents":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					// Should not be called due to depth limit
					t.Errorf("Unexpected request to %s with maxDepth=1", r.URL.Path)
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				}
			},
			wantErr:  false,
			wantSize: 1, // only CWE-79 because depth limit prevents recursion
		},
		{
			name:     "error on initial fetch",
			id:       79,
			maxDepth: 5,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			registry := NewRegistry()
			fetcher := NewTreeFetcher(client, registry, tt.maxDepth)

			err := fetcher.FetchWithAncestors(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWithAncestors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if registry.Size() != tt.wantSize {
					t.Errorf("FetchWithAncestors() registry size = %v, want %v", registry.Size(), tt.wantSize)
				}
			}
		})
	}
}

func TestTreeFetcher_FetchWithDescendants(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		maxDepth int
		handler  http.HandlerFunc
		wantErr  bool
		wantSize int
	}{
		{
			name:     "success - fetches self and children",
			id:       74,
			maxDepth: 5,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/74":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/74/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ParentOf","cweId":79,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/weakness/79":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				}
			},
			wantErr:  false,
			wantSize: 2, // CWE-74 and CWE-79
		},
		{
			name:     "max depth limit",
			id:       74,
			maxDepth: 1,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/74":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/74/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ParentOf","cweId":79,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				}
			},
			wantErr:  false,
			wantSize: 1, // only CWE-74 because depth limit prevents recursion
		},
		{
			name:     "already fetched - skips duplicate fetch",
			id:       74,
			maxDepth: 5,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/74":
					// This should not be called since CWE-74 is already in registry
					t.Error("Should not fetch CWE-74 again since it's already in registry")
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/74/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				}
			},
			wantErr:  false,
			wantSize: 1, // already registered
		},
		{
			name:     "error on children fetch",
			id:       74,
			maxDepth: 5,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/74":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/74/children":
					w.WriteHeader(http.StatusInternalServerError)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			registry := NewRegistry()

			// For the "already fetched" test case, pre-register the CWE
			if tt.name == "already fetched - skips duplicate fetch" {
				_ = registry.Register(&CWE{ID: 74, Name: "Injection", Description: "desc", CWEType: "weakness"})
			}

			fetcher := NewTreeFetcher(client, registry, tt.maxDepth)

			err := fetcher.FetchWithDescendants(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWithDescendants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if registry.Size() != tt.wantSize {
					t.Errorf("FetchWithDescendants() registry size = %v, want %v", registry.Size(), tt.wantSize)
				}
			}
		})
	}
}

func TestTreeFetcher_FetchFullTree(t *testing.T) {
	tests := []struct {
		name     string
		rootID   int
		maxDepth int
		handler  http.HandlerFunc
		wantErr  bool
	}{
		{
			name:     "success - fetches ancestors and descendants",
			rootID:   79,
			maxDepth: 5,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch r.URL.Path {
				case "/cwe/weakness/79":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/parents":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/weakness/74":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/74/parents":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/79/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"nature":"ParentOf","cweId":80,"viewId":1000}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/weakness/80":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[{"id":80,"name":"Browser XSS","description":"desc","cwe_type":"weakness"}]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/80/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				case "/cwe/74/children":
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				default:
					resp := struct {
						Data json.RawMessage `json:"Data"`
					}{
						Data: json.RawMessage(`[]`),
					}
					json.NewEncoder(w).Encode(resp)
				}
			},
			wantErr: false,
		},
		{
			name:     "error on ancestors",
			rootID:   79,
			maxDepth: 5,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := newTestAPIClient(server)
			registry := NewRegistry()
			fetcher := NewTreeFetcher(client, registry, tt.maxDepth)

			err := fetcher.FetchFullTree(context.Background(), tt.rootID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchFullTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Should have fetched the root, its parent, and its child
				size := registry.Size()
				if size < 1 {
					t.Errorf("FetchFullTree() registry should have at least 1 entry, got %d", size)
				}
			}
		})
	}
}

// ==================== TreeFetcher edge cases ====================

func TestTreeFetcher_FetchWithAncestors_SkipsAlreadyFetched(t *testing.T) {
	// When a CWE is already in the registry, fetchAncestorsRecursive should skip it
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// No request should be made since CWE-79 is already in registry
		t.Errorf("Unexpected request to %s - CWE-79 should already be in registry", r.URL.Path)
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{
			Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	registry := NewRegistry()
	_ = registry.Register(&CWE{ID: 79, Name: "XSS", Description: "desc", CWEType: "weakness"})

	fetcher := NewTreeFetcher(client, registry, 5)
	err := fetcher.FetchWithAncestors(context.Background(), 79)
	if err != nil {
		t.Errorf("FetchWithAncestors() error = %v, want nil (should skip already-fetched)", err)
	}
}

func TestTreeFetcher_FetchWithDescendants_AddsRelationships(t *testing.T) {
	// Verify that FetchWithDescendants adds ParentOf relationships to the existing CWE entry
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/cwe/weakness/74":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/74/children":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ParentOf","cweId":79,"viewId":1000}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/children":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)

	err := fetcher.FetchWithDescendants(context.Background(), 74)
	if err != nil {
		t.Fatalf("FetchWithDescendants() error = %v", err)
	}

	// Verify relationships were added to CWE-74
	cwe74, ok := registry.Get(74)
	if !ok {
		t.Fatal("CWE-74 should be in registry")
	}

	foundParentOf := false
	for _, rel := range cwe74.Relationships {
		if rel.Nature == RelationshipParentOf && rel.CWEID == 79 {
			foundParentOf = true
		}
	}
	if !foundParentOf {
		t.Error("CWE-74 should have ParentOf relationship to CWE-79")
	}
}

func TestTreeFetcher_FetchWithAncestors_AddsRelationships(t *testing.T) {
	// Verify that FetchWithAncestors adds ChildOf relationships
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/weakness/74":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/74/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)

	err := fetcher.FetchWithAncestors(context.Background(), 79)
	if err != nil {
		t.Fatalf("FetchWithAncestors() error = %v", err)
	}

	// Verify relationships were added to CWE-79
	cwe79, ok := registry.Get(79)
	if !ok {
		t.Fatal("CWE-79 should be in registry")
	}

	foundChildOf := false
	for _, rel := range cwe79.Relationships {
		if rel.Nature == RelationshipChildOf && rel.CWEID == 74 {
			foundChildOf = true
		}
	}
	if !foundChildOf {
		t.Error("CWE-79 should have ChildOf relationship to CWE-74")
	}
}

func TestTreeFetcher_FetchWithAncestors_ParentsFetchError(t *testing.T) {
	// When GetParents fails, FetchWithAncestors should return an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/parents":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)

	err := fetcher.FetchWithAncestors(context.Background(), 79)
	if err == nil {
		t.Error("FetchWithAncestors() should return error when GetParents fails")
	}
}

// ==================== MultipleFetcher_FetchMultipleToRegistry with nil values ====================

func TestMultipleFetcher_FetchMultipleToRegistry_SkipsNilCWEs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := struct {
			Data json.RawMessage `json:"Data"`
		}{
			Data: json.RawMessage(`{"79":{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"},"89":null}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	fetcher := NewMultipleFetcher(client)
	registry := NewRegistry()

	err := fetcher.FetchMultipleToRegistry(context.Background(), []int{79, 89}, registry)
	if err != nil {
		t.Errorf("FetchMultipleToRegistry() error = %v", err)
	}
	// Only CWE-79 should be registered (CWE-89 is nil)
	if registry.Size() != 1 {
		t.Errorf("FetchMultipleToRegistry() registry size = %v, want 1", registry.Size())
	}
}

// ==================== FetchWithRelations with viewID ====================

func TestBasicFetcher_FetchWithRelations_WithViewID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Verify the view parameter is passed through
		if r.URL.Path == "/cwe/79/parents" || r.URL.Path == "/cwe/79/children" {
			viewID := r.URL.Query().Get("view")
			if viewID != "1000" {
				t.Errorf("Expected view=1000 query param, got view=%s", viewID)
			}
		}

		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/children":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	fetcher := NewBasicFetcher(client)

	cwe, err := fetcher.FetchWithRelations(context.Background(), 79, 1000)
	if err != nil {
		t.Fatalf("FetchWithRelations() error = %v", err)
	}
	if cwe.ID != 79 {
		t.Errorf("FetchWithRelations() ID = %v, want 79", cwe.ID)
	}
}

// ==================== DataFetcher interface compliance ====================

func TestBasicFetcher_ImplementsDataFetcher(t *testing.T) {
	// Verify that BasicFetcher implements the DataFetcher interface
	var _ DataFetcher = (*BasicFetcher)(nil)
}

// ==================== FetchFullTree with error on descendants ====================

func TestTreeFetcher_FetchFullTree_ErrorOnDescendants(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		callCount++

		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(fmt.Sprintf(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`)),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/children":
			// Return error on children fetch
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)

	err := fetcher.FetchFullTree(context.Background(), 79)
	// FetchFullTree calls FetchWithAncestors first (which succeeds),
	// then FetchWithDescendants (which fails on children fetch)
	if err == nil {
		t.Error("FetchFullTree() should return error when FetchWithDescendants fails")
	}

	// But CWE-79 should still be in registry from the ancestors phase
	if !registry.Contains(79) {
		t.Error("CWE-79 should be in registry from ancestors phase")
	}
}

// ==================== FetchWithAncestors recursive depth test ====================

func TestTreeFetcher_FetchWithAncestors_DeepRecursion(t *testing.T) {
	// Test that the recursive ancestor fetching works correctly through multiple levels
	requestLog := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		requestLog = append(requestLog, r.URL.Path)

		switch r.URL.Path {
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ChildOf","cweId":74,"viewId":1000}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/weakness/74":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/74/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ChildOf","cweId":664,"viewId":1000}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/weakness/664":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":664,"name":"Resource Control","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/664/parents":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)

	err := fetcher.FetchWithAncestors(context.Background(), 79)
	if err != nil {
		t.Fatalf("FetchWithAncestors() error = %v", err)
	}

	// All three CWEs should be in the registry
	if !registry.Contains(79) {
		t.Error("CWE-79 should be in registry")
	}
	if !registry.Contains(74) {
		t.Error("CWE-74 should be in registry")
	}
	if !registry.Contains(664) {
		t.Error("CWE-664 should be in registry")
	}
	if registry.Size() != 3 {
		t.Errorf("Registry size = %d, want 3", registry.Size())
	}
}

// ==================== FetchWithDescendants recursive depth test ====================

func TestTreeFetcher_FetchWithDescendants_DeepRecursion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/cwe/weakness/664":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":664,"name":"Resource Control","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/664/children":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ParentOf","cweId":74,"viewId":1000}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/weakness/74":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":74,"name":"Injection","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/74/children":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"nature":"ParentOf","cweId":79,"viewId":1000}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/weakness/79":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[{"id":79,"name":"XSS","description":"desc","cwe_type":"weakness"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "/cwe/79/children":
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			resp := struct {
				Data json.RawMessage `json:"Data"`
			}{
				Data: json.RawMessage(`[]`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := newTestAPIClient(server)
	registry := NewRegistry()
	fetcher := NewTreeFetcher(client, registry, 5)

	err := fetcher.FetchWithDescendants(context.Background(), 664)
	if err != nil {
		t.Fatalf("FetchWithDescendants() error = %v", err)
	}

	// All three CWEs should be in the registry
	if !registry.Contains(664) {
		t.Error("CWE-664 should be in registry")
	}
	if !registry.Contains(74) {
		t.Error("CWE-74 should be in registry")
	}
	if !registry.Contains(79) {
		t.Error("CWE-79 should be in registry")
	}
	if registry.Size() != 3 {
		t.Errorf("Registry size = %d, want 3", registry.Size())
	}
}
