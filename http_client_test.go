package cwe

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		client := NewHTTPClient("https://example.com/api")
		defer client.Close()

		if client.baseURL != "https://example.com/api" {
			t.Errorf("expected baseURL 'https://example.com/api', got %q", client.baseURL)
		}
		if client.userAgent != DefaultUserAgent {
			t.Errorf("expected userAgent %q, got %q", DefaultUserAgent, client.userAgent)
		}
		if client.maxRetries != 0 {
			t.Errorf("expected maxRetries 0, got %d", client.maxRetries)
		}
		if client.retryDelay != 1*time.Second {
			t.Errorf("expected retryDelay 1s, got %v", client.retryDelay)
		}
		if client.rateLimiter != nil {
			t.Error("expected nil rateLimiter by default")
		}
	})

	t.Run("WithRetry", func(t *testing.T) {
		client := NewHTTPClient("https://example.com", WithRetry(3, 5*time.Second))
		defer client.Close()

		if client.maxRetries != 3 {
			t.Errorf("expected maxRetries 3, got %d", client.maxRetries)
		}
		if client.retryDelay != 5*time.Second {
			t.Errorf("expected retryDelay 5s, got %v", client.retryDelay)
		}
	})

	t.Run("WithHTTPRateLimiter", func(t *testing.T) {
		client := NewHTTPClient("https://example.com", WithHTTPRateLimiter(1.0, 5))
		defer client.Close()

		if client.rateLimiter == nil {
			t.Fatal("expected rateLimiter to be set")
		}
	})

	t.Run("WithHTTPTimeout", func(t *testing.T) {
		client := NewHTTPClient("https://example.com", WithHTTPTimeout(10*time.Second))
		defer client.Close()

		if client.client.Timeout != 10*time.Second {
			t.Errorf("expected timeout 10s, got %v", client.client.Timeout)
		}
	})

	t.Run("WithUserAgent", func(t *testing.T) {
		client := NewHTTPClient("https://example.com", WithUserAgent("test-agent/1.0"))
		defer client.Close()

		if client.userAgent != "test-agent/1.0" {
			t.Errorf("expected userAgent 'test-agent/1.0', got %q", client.userAgent)
		}
	})

	t.Run("WithHTTPClient", func(t *testing.T) {
		customHTTP := &http.Client{Timeout: 99 * time.Second}
		client := NewHTTPClient("https://example.com", WithHTTPClient(customHTTP))
		defer client.Close()

		if client.client != customHTTP {
			t.Error("expected custom http.Client to be set")
		}
	})

	t.Run("multiple options", func(t *testing.T) {
		client := NewHTTPClient("https://example.com",
			WithRetry(2, 3*time.Second),
			WithHTTPRateLimiter(0.5, 2),
			WithUserAgent("multi-opt"),
		)
		defer client.Close()

		if client.maxRetries != 2 {
			t.Errorf("expected maxRetries 2, got %d", client.maxRetries)
		}
		if client.retryDelay != 3*time.Second {
			t.Errorf("expected retryDelay 3s, got %v", client.retryDelay)
		}
		if client.rateLimiter == nil {
			t.Error("expected rateLimiter to be set")
		}
		if client.userAgent != "multi-opt" {
			t.Errorf("expected userAgent 'multi-opt', got %q", client.userAgent)
		}
	})
}

func TestHTTPClient_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := map[string]interface{}{"key": "value"}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.Header.Get("User-Agent") == "" {
				t.Error("expected User-Agent header to be set")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expected)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		var result map[string]interface{}
		if err := client.Get(context.Background(), "/test", &result); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["key"] != "value" {
			t.Errorf("expected result[key]=value, got %v", result["key"])
		}
	})

	t.Run("404 error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		var result map[string]interface{}
		err := client.Get(context.Background(), "/missing", &result)
		if err == nil {
			t.Fatal("expected error for 404")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		} else if apiErr.StatusCode != 404 {
			t.Errorf("expected status 404, got %d", apiErr.StatusCode)
		}
	})

	t.Run("500 error with retry", func(t *testing.T) {
		callCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if callCount < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL, WithRetry(3, 10*time.Millisecond))
		defer client.Close()

		var result map[string]string
		if err := client.Get(context.Background(), "/retry", &result); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got %d", callCount)
		}
	})

	t.Run("500 error exhausts retries", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL, WithRetry(2, 10*time.Millisecond))
		defer client.Close()

		var result map[string]string
		err := client.Get(context.Background(), "/fail", &result)
		if err == nil {
			t.Fatal("expected error when retries exhausted")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(500 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL, WithRetry(3, 50*time.Millisecond))
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		var result map[string]string
		err := client.Get(ctx, "/slow", &result)
		if err == nil {
			t.Fatal("expected error from cancelled context")
		}
	})

	t.Run("nil result", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"key": "value"}`))
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		if err := client.Get(context.Background(), "/test", nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`not valid json`))
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		var result map[string]interface{}
		err := client.Get(context.Background(), "/badjson", &result)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}

		var parseErr *ParseError
		if !errors.As(err, &parseErr) {
			t.Errorf("expected ParseError, got %T: %v", err, err)
		}
	})
}

func TestHTTPClient_GetRaw(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedBody := []byte("raw response body")
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(expectedBody)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		data, err := client.GetRaw(context.Background(), "/raw")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "raw response body" {
			t.Errorf("expected 'raw response body', got %q", string(data))
		}
	})

	t.Run("error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		_, err := client.GetRaw(context.Background(), "/bad")
		if err == nil {
			t.Fatal("expected error for 400 response")
		}
	})
}

func TestHTTPClient_Post(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
			}
			var body map[string]string
			json.NewDecoder(r.Body).Decode(&body)
			if body["name"] != "test" {
				t.Errorf("expected body name=test, got %q", body["name"])
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		var result map[string]string
		if err := client.Post(context.Background(), "/submit", map[string]string{"name": "test"}, &result); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["result"] != "ok" {
			t.Errorf("expected result=ok, got %q", result["result"])
		}
	})

	t.Run("serialization error", func(t *testing.T) {
		client := NewHTTPClient("https://example.com")
		defer client.Close()

		// Channels cannot be serialized to JSON
		badBody := map[string]chan int{"ch": make(chan int)}
		err := client.Post(context.Background(), "/test", badBody, nil)
		if err == nil {
			t.Fatal("expected error for unserializable body")
		}
	})

	t.Run("server error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		err := client.Post(context.Background(), "/fail", map[string]string{"k": "v"}, nil)
		if err == nil {
			t.Fatal("expected error for 500 response")
		}
	})

	t.Run("nil body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		var result map[string]string
		if err := client.Post(context.Background(), "/test", nil, &result); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil result", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		if err := client.Post(context.Background(), "/test", map[string]string{"k": "v"}, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestHTTPClient_PostForm(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Errorf("expected Content-Type application/x-www-form-urlencoded, got %q", r.Header.Get("Content-Type"))
			}
			r.ParseForm()
			if r.FormValue("key") != "value" {
				t.Errorf("expected form key=value, got %q", r.FormValue("key"))
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		form := url.Values{}
		form.Set("key", "value")

		var result map[string]string
		if err := client.PostForm(context.Background(), "/form", form, &result); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["result"] != "ok" {
			t.Errorf("expected result=ok, got %q", result["result"])
		}
	})

	t.Run("error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		form := url.Values{}
		form.Set("key", "value")

		err := client.PostForm(context.Background(), "/form", form, nil)
		if err == nil {
			t.Fatal("expected error for 400 response")
		}
	})

	t.Run("server error with retry", func(t *testing.T) {
		callCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if callCount < 2 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL, WithRetry(2, 10*time.Millisecond))
		defer client.Close()

		form := url.Values{}
		form.Set("key", "value")

		var result map[string]string
		if err := client.PostForm(context.Background(), "/form", form, &result); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("context cancellation on retry", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL, WithRetry(5, 100*time.Millisecond))
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		form := url.Values{}
		form.Set("key", "value")

		err := client.PostForm(ctx, "/form", form, nil)
		if err == nil {
			t.Fatal("expected error from cancelled context")
		}
	})

	t.Run("nil result", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
		}))
		defer srv.Close()

		client := NewHTTPClient(srv.URL)
		defer client.Close()

		form := url.Values{}
		form.Set("key", "value")

		if err := client.PostForm(context.Background(), "/form", form, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rate limiter error", func(t *testing.T) {
		client := NewHTTPClient("https://example.com", WithHTTPRateLimiter(0.0001, 0))
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		form := url.Values{}
		form.Set("key", "value")

		err := client.PostForm(ctx, "/form", form, nil)
		if err == nil {
			t.Fatal("expected error from rate limiter")
		}
	})
}

func TestHTTPClient_GettersSetters(t *testing.T) {
	client := NewHTTPClient("https://original.com")

	t.Run("GetBaseURL", func(t *testing.T) {
		if got := client.GetBaseURL(); got != "https://original.com" {
			t.Errorf("expected https://original.com, got %q", got)
		}
	})

	t.Run("SetBaseURL", func(t *testing.T) {
		client.SetBaseURL("https://new.com")
		if got := client.GetBaseURL(); got != "https://new.com" {
			t.Errorf("expected https://new.com, got %q", got)
		}
	})

	t.Run("GetMaxRetries", func(t *testing.T) {
		if got := client.GetMaxRetries(); got != 0 {
			t.Errorf("expected 0, got %d", got)
		}
	})

	t.Run("SetMaxRetries", func(t *testing.T) {
		client.SetMaxRetries(5)
		if got := client.GetMaxRetries(); got != 5 {
			t.Errorf("expected 5, got %d", got)
		}
	})

	t.Run("GetRetryDelay", func(t *testing.T) {
		if got := client.GetRetryDelay(); got != 1*time.Second {
			t.Errorf("expected 1s, got %v", got)
		}
	})

	t.Run("SetRetryDelay", func(t *testing.T) {
		client.SetRetryDelay(5 * time.Second)
		if got := client.GetRetryDelay(); got != 5*time.Second {
			t.Errorf("expected 5s, got %v", got)
		}
	})

	t.Run("GetRateLimiter", func(t *testing.T) {
		if got := client.GetRateLimiter(); got != nil {
			t.Error("expected nil rate limiter")
		}
	})

	t.Run("SetRateLimiter", func(t *testing.T) {
		limiter := NewRateLimiter(1.0, 5)
		client.SetRateLimiter(limiter)
		if got := client.GetRateLimiter(); got != limiter {
			t.Error("expected set rate limiter")
		}
	})

	t.Run("GetHTTPClient", func(t *testing.T) {
		if got := client.GetHTTPClient(); got == nil {
			t.Error("expected non-nil http.Client")
		}
	})

	t.Run("SetHTTPClient", func(t *testing.T) {
		customClient := &http.Client{Timeout: 99 * time.Second}
		client.SetHTTPClient(customClient)
		if got := client.GetHTTPClient(); got != customClient {
			t.Error("expected custom http.Client")
		}
	})

	client.Close()
}

func TestHTTPClient_Close(t *testing.T) {
	client := NewHTTPClient("https://example.com")
	// Close should not panic
	client.Close()
	// Calling Close again should also not panic
	client.Close()
}

func TestHTTPClient_RateLimiterWithGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, WithHTTPRateLimiter(100.0, 10))
	defer client.Close()

	var result map[string]string
	if err := client.Get(context.Background(), "/test", &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
