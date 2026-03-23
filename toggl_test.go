package toggl

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("empty token returns error", func(t *testing.T) {
		_, err := NewClient("")
		if err == nil {
			t.Fatal("expected error for empty token, got nil")
		}
	})

	t.Run("valid token returns client", func(t *testing.T) {
		c, err := NewClient("test-token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})
}

func TestWithBaseURL(t *testing.T) {
	c, err := NewClient("test-token", WithBaseURL("https://example.com"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.baseURL != "https://example.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://example.com")
	}
}

func TestWithHTTPClient(t *testing.T) {
	wantTimeout := 42 * time.Second
	custom := &http.Client{Timeout: wantTimeout}
	c, err := NewClient("test-token", WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient != custom {
		t.Error("httpClient was not replaced by WithHTTPClient")
	}
	if c.httpClient.Timeout != wantTimeout {
		t.Errorf("Timeout = %v, want %v; WithHTTPClient should preserve the custom client's timeout", c.httpClient.Timeout, wantTimeout)
	}
}

func TestWithHTTPClient_noTimeout(t *testing.T) {
	// A custom client with no timeout should receive defaultTimeout.
	custom := &http.Client{}
	c, err := NewClient("test-token", WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient.Timeout != defaultTimeout {
		t.Errorf("Timeout = %v, want defaultTimeout %v", c.httpClient.Timeout, defaultTimeout)
	}
}

func TestWithHTTPClient_nil(t *testing.T) {
	// A nil value must not replace the default client (would panic on timeout assignment).
	c, err := NewClient("test-token", WithHTTPClient(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient == nil {
		t.Error("httpClient is nil after WithHTTPClient(nil); default client should be retained")
	}
}

func TestWithTimeout(t *testing.T) {
	want := 5 * time.Second
	c, err := NewClient("test-token", WithTimeout(want))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient.Timeout != want {
		t.Errorf("Timeout = %v, want %v", c.httpClient.Timeout, want)
	}
}

func TestWithTimeout_zero(t *testing.T) {
	// WithTimeout(0) must disable the timeout, not fall back to defaultTimeout.
	c, err := NewClient("test-token", WithTimeout(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient.Timeout != 0 {
		t.Errorf("Timeout = %v, want 0; explicit WithTimeout(0) should disable the timeout", c.httpClient.Timeout)
	}
}

func TestWithTimeout_AfterWithHTTPClient(t *testing.T) {
	// WithTimeout must apply even when WithHTTPClient appears before it.
	want := 7 * time.Second
	custom := &http.Client{}
	c, err := NewClient("test-token", WithHTTPClient(custom), WithTimeout(want))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient.Timeout != want {
		t.Errorf("Timeout = %v, want %v (WithTimeout after WithHTTPClient)", c.httpClient.Timeout, want)
	}
}

func TestWithTimeout_BeforeWithHTTPClient(t *testing.T) {
	// WithTimeout must apply even when it appears before WithHTTPClient.
	want := 7 * time.Second
	custom := &http.Client{}
	c, err := NewClient("test-token", WithTimeout(want), WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.httpClient.Timeout != want {
		t.Errorf("Timeout = %v, want %v (WithTimeout before WithHTTPClient)", c.httpClient.Timeout, want)
	}
}
