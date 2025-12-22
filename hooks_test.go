package neuron_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/kolosys/neuron"
)

func TestClient_Hooks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Hook") != "true" {
			t.Error("request hook header missing")
		}
		w.Header().Set("X-Response-Hook", "true")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	reqHooked := false
	respHooked := false

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		RequestHooks: []RequestHook{
			func(req *http.Request) error {
				req.Header.Set("X-Request-Hook", "true")
				reqHooked = true
				return nil
			},
		},
		ResponseHooks: []ResponseHook{
			func(resp *http.Response) error {
				if resp.Header.Get("X-Response-Hook") != "true" {
					return fmt.Errorf("response hook header missing")
				}
				respHooked = true
				return nil
			},
		},
	})

	_, err := client.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !reqHooked {
		t.Error("request hook not called")
	}
	if !respHooked {
		t.Error("response hook not called")
	}
}

func TestAddResponseTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		ResponseHooks: []ResponseHook{
			AddResponseTimeout(10 * time.Millisecond),
		},
	})

	_, err := client.Get("/")
	if err == nil {
		t.Fatal("expected timeout error from hook, got nil")
	}

	ce, ok := err.(ClientError)
	if !ok {
		t.Fatalf("expected ClientError, got %T", err)
	}
	if ce.Type != ErrorTypeTimeout {
		t.Errorf("expected ErrorTypeTimeout, got %v", ce.Type)
	}
}

type invalidJSONProvider struct{}

func (p invalidJSONProvider) ContentType() string { return "application/json" }
func (p invalidJSONProvider) Body() (io.Reader, error) {
	return strings.NewReader("{invalid json"), nil
}

func TestAddValidation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		RequestHooks: []RequestHook{
			AddValidation(&JSONValidator{}),
		},
	})

	// Invalid JSON using a provider
	_, err := client.Post("/", &RequestOptions{
		Body: invalidJSONProvider{},
	})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestAddRateLimitHandler(t *testing.T) {
	updater := &testRateLimitUpdater{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "50")
		w.Header().Set("X-RateLimit-Bucket", "test-bucket")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		ResponseHooks: []ResponseHook{
			AddRateLimitHandler(updater),
		},
	})

	_, err := client.Get("/api/test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(updater.calls) != 1 {
		t.Fatalf("expected 1 update call, got %d", len(updater.calls))
	}

	call := updater.calls[0]
	if call.method != "GET" {
		t.Errorf("method = %q, want GET", call.method)
	}
	if call.info.Limit != 100 {
		t.Errorf("Limit = %d, want 100", call.info.Limit)
	}
	if call.info.Remaining != 50 {
		t.Errorf("Remaining = %d, want 50", call.info.Remaining)
	}
	if call.info.Bucket != "test-bucket" {
		t.Errorf("Bucket = %q, want test-bucket", call.info.Bucket)
	}
}

func TestAddRateLimitHandler_NilUpdater(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "100")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		ResponseHooks: []ResponseHook{
			AddRateLimitHandler(nil),
		},
	})

	_, err := client.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestAddRateLimitRetry(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL:    ts.URL,
		MaxRetries: 0, // Disable retries so hook error propagates
		ResponseHooks: []ResponseHook{
			AddRateLimitRetry(),
		},
	})

	_, err := client.Get("/")
	if err == nil {
		t.Fatal("expected rate limit error, got nil")
	}

	ce, ok := err.(ClientError)
	if !ok {
		t.Fatalf("expected ClientError, got %T", err)
	}
	if ce.Type != ErrorTypeRateLimit {
		t.Errorf("expected ErrorTypeRateLimit, got %v", ce.Type)
	}
	if ce.StatusCode != 429 {
		t.Errorf("StatusCode = %d, want 429", ce.StatusCode)
	}
}

func TestAddRateLimitRetry_NoRateLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		ResponseHooks: []ResponseHook{
			AddRateLimitRetry(),
		},
	})

	_, err := client.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestAddRateLimitLogging(t *testing.T) {
	var logCalls []struct {
		method string
		path   string
		info   *RateLimitInfo
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "25")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		ResponseHooks: []ResponseHook{
			AddRateLimitLogging(func(method, path string, info *RateLimitInfo) {
				logCalls = append(logCalls, struct {
					method string
					path   string
					info   *RateLimitInfo
				}{method, path, info})
			}),
		},
	})

	_, err := client.Get("/api/resource")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(logCalls) != 1 {
		t.Fatalf("expected 1 log call, got %d", len(logCalls))
	}

	if logCalls[0].method != "GET" {
		t.Errorf("method = %q, want GET", logCalls[0].method)
	}
	if logCalls[0].info.Remaining != 25 {
		t.Errorf("Remaining = %d, want 25", logCalls[0].info.Remaining)
	}
}

type testRateLimitUpdater struct {
	calls []struct {
		method   string
		endpoint string
		info     *RateLimitInfo
	}
}

func (u *testRateLimitUpdater) UpdateFromHeaders(method, endpoint string, info *RateLimitInfo) error {
	u.calls = append(u.calls, struct {
		method   string
		endpoint string
		info     *RateLimitInfo
	}{method, endpoint, info})
	return nil
}
