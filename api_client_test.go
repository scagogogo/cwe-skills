package cwe

import (
	"net/http"
	"testing"
	"time"
)

func TestNewAPIClient(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		client := NewAPIClient()
		defer client.Close()

		if client.baseURL != DefaultBaseURL {
			t.Errorf("expected baseURL %q, got %q", DefaultBaseURL, client.baseURL)
		}
		if client.httpClient == nil {
			t.Error("expected non-nil httpClient")
		}
	})

	t.Run("WithAPIBaseURL", func(t *testing.T) {
		client := NewAPIClient(WithAPIBaseURL("https://custom.api.com"))
		defer client.Close()

		if client.baseURL != "https://custom.api.com" {
			t.Errorf("expected baseURL 'https://custom.api.com', got %q", client.baseURL)
		}
		// Ensure HTTP client's baseURL is also updated
		if client.httpClient.GetBaseURL() != "https://custom.api.com" {
			t.Errorf("expected httpClient baseURL 'https://custom.api.com', got %q", client.httpClient.GetBaseURL())
		}
	})

	t.Run("WithAPITimeout", func(t *testing.T) {
		client := NewAPIClient(WithAPITimeout(60 * time.Second))
		defer client.Close()

		if client.httpClient.client.Timeout != 60*time.Second {
			t.Errorf("expected timeout 60s, got %v", client.httpClient.client.Timeout)
		}
	})

	t.Run("WithAPIRateLimit", func(t *testing.T) {
		client := NewAPIClient(WithAPIRateLimit(1.0, 5))
		defer client.Close()

		limiter := client.GetRateLimiter()
		if limiter == nil {
			t.Fatal("expected rate limiter to be set")
		}
	})

	t.Run("WithAPIRetry", func(t *testing.T) {
		client := NewAPIClient(WithAPIRetry(3, 5*time.Second))
		defer client.Close()

		if client.httpClient.GetMaxRetries() != 3 {
			t.Errorf("expected maxRetries 3, got %d", client.httpClient.GetMaxRetries())
		}
		if client.httpClient.GetRetryDelay() != 5*time.Second {
			t.Errorf("expected retryDelay 5s, got %v", client.httpClient.GetRetryDelay())
		}
	})

	t.Run("WithAPIHTTPClient", func(t *testing.T) {
		client := NewAPIClient(WithAPIHTTPClient(
			WithUserAgent("custom-ua"),
			WithRetry(2, 3*time.Second),
		))
		defer client.Close()

		if client.httpClient.userAgent != "custom-ua" {
			t.Errorf("expected userAgent 'custom-ua', got %q", client.httpClient.userAgent)
		}
		if client.httpClient.GetMaxRetries() != 2 {
			t.Errorf("expected maxRetries 2, got %d", client.httpClient.GetMaxRetries())
		}
	})

	t.Run("multiple options", func(t *testing.T) {
		client := NewAPIClient(
			WithAPIBaseURL("https://test.com"),
			WithAPITimeout(45*time.Second),
			WithAPIRetry(2, 2*time.Second),
		)
		defer client.Close()

		if client.baseURL != "https://test.com" {
			t.Errorf("expected baseURL 'https://test.com', got %q", client.baseURL)
		}
	})
}

func TestAPIClient_GetBaseURL(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	if got := client.GetBaseURL(); got != DefaultBaseURL {
		t.Errorf("expected %q, got %q", DefaultBaseURL, got)
	}
}

func TestAPIClient_SetBaseURL(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	client.SetBaseURL("https://new-base.com")
	if got := client.GetBaseURL(); got != "https://new-base.com" {
		t.Errorf("expected 'https://new-base.com', got %q", got)
	}
	// Verify HTTP client is also updated
	if got := client.httpClient.GetBaseURL(); got != "https://new-base.com" {
		t.Errorf("expected httpClient baseURL 'https://new-base.com', got %q", got)
	}
}

func TestAPIClient_GetHTTPClient(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	hc := client.GetHTTPClient()
	if hc == nil {
		t.Error("expected non-nil HTTPClient")
	}
}

func TestAPIClient_SetHTTPClient(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	newHTTPClient := NewHTTPClient("https://custom.com")
	client.SetHTTPClient(newHTTPClient)

	if client.GetHTTPClient() != newHTTPClient {
		t.Error("expected HTTPClient to be updated")
	}
	newHTTPClient.Close()
}

func TestAPIClient_GetRateLimiter(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	limiter := client.GetRateLimiter()
	if limiter == nil {
		t.Error("expected non-nil rate limiter (default)")
	}
}

func TestAPIClient_SetRateLimiter(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	newLimiter := NewRateLimiter(5.0, 10)
	client.SetRateLimiter(newLimiter)

	if client.GetRateLimiter() != newLimiter {
		t.Error("expected rate limiter to be updated")
	}
}

func TestAPIClient_Close(t *testing.T) {
	client := NewAPIClient()
	// Close should not panic
	client.Close()
	// Double close should not panic
	client.Close()
}

func TestAPIClient_HTTPClientSync(t *testing.T) {
	t.Run("baseURL sync on creation", func(t *testing.T) {
		client := NewAPIClient(WithAPIBaseURL("https://sync-test.com"))
		defer client.Close()

		if client.baseURL != client.httpClient.baseURL {
			t.Errorf("baseURL mismatch: api=%q http=%q",
				client.baseURL, client.httpClient.baseURL)
		}
	})

	t.Run("baseURL sync on set", func(t *testing.T) {
		client := NewAPIClient()
		defer client.Close()

		client.SetBaseURL("https://new-sync.com")
		if client.baseURL != client.httpClient.baseURL {
			t.Errorf("baseURL mismatch after SetBaseURL: api=%q http=%q",
				client.baseURL, client.httpClient.baseURL)
		}
	})

	t.Run("custom HTTPClient override", func(t *testing.T) {
		client := NewAPIClient()
		defer client.Close()

		customHTTP := NewHTTPClient("https://override.com", WithHTTPTimeout(99*time.Second))
		client.SetHTTPClient(customHTTP)

		hc := client.GetHTTPClient()
		if hc.GetBaseURL() != "https://override.com" {
			t.Errorf("expected override.com, got %q", hc.GetBaseURL())
		}
		customHTTP.Close()
	})
}

func TestAPIClient_DefaultHTTPClientFields(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	hc := client.GetHTTPClient()

	// The default APIClient should have a rate limiter
	if hc.GetRateLimiter() == nil {
		t.Error("expected default rate limiter on APIClient's HTTPClient")
	}

	// The default base URL should be synced
	if hc.GetBaseURL() != DefaultBaseURL {
		t.Errorf("expected httpClient baseURL %q, got %q", DefaultBaseURL, hc.GetBaseURL())
	}

	// Verify the underlying http.Client
	if hc.GetHTTPClient() == nil {
		t.Error("expected non-nil underlying http.Client")
	}

	// Verify the underlying http.Client has a timeout
	if hc.GetHTTPClient().Timeout != DefaultTimeout {
		t.Errorf("expected timeout %v, got %v", DefaultTimeout, hc.GetHTTPClient().Timeout)
	}
}

func TestAPIClient_SetHTTPClientThenGetHTTPClient(t *testing.T) {
	client := NewAPIClient()
	defer client.Close()

	customHTTP := &http.Client{Timeout: 5 * time.Second}
	newHC := NewHTTPClient("https://custom.com", WithHTTPClient(customHTTP))
	client.SetHTTPClient(newHC)

	retrievedHC := client.GetHTTPClient()
	if retrievedHC != newHC {
		t.Error("expected GetHTTPClient to return the same HTTPClient that was set")
	}
	retrievedHC.Close()
}
