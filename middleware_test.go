package neuron

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestMiddlewareChain tests the middleware chain functionality
func TestMiddlewareChain(t *testing.T) {
	chain := NewMiddlewareChain()
	executionOrder := []string{}

	// Add request middleware
	chain.AddRequestMiddleware(func(req *http.Request) error {
		executionOrder = append(executionOrder, "request1")
		return nil
	})
	chain.AddRequestMiddleware(func(req *http.Request) error {
		executionOrder = append(executionOrder, "request2")
		return nil
	})

	// Add response middleware
	chain.AddResponseMiddleware(func(resp *http.Response) error {
		executionOrder = append(executionOrder, "response1")
		return nil
	})
	chain.AddResponseMiddleware(func(resp *http.Response) error {
		executionOrder = append(executionOrder, "response2")
		return nil
	})

	// Create test request and response
	req := httptest.NewRequest("GET", "/test", nil)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}

	// Apply middleware
	err := chain.ApplyRequestMiddleware(req)
	if err != nil {
		t.Errorf("request middleware failed: %v", err)
	}

	err = chain.ApplyResponseMiddleware(resp)
	if err != nil {
		t.Errorf("response middleware failed: %v", err)
	}

	// Check execution order
	expectedOrder := []string{"request1", "request2", "response1", "response2"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d middleware executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("middleware[%d]: expected %s, got %s", i, expected, executionOrder[i])
		}
	}
}

// TestAddAuthentication tests authentication middleware
func TestAddAuthentication(t *testing.T) {
	authProvider := &StaticAuthProvider{
		Token:  "test-token",
		Prefix: "Bearer",
	}

	middleware := AddAuthentication(authProvider)
	req := httptest.NewRequest("GET", "/test", nil)

	err := middleware(req)
	if err != nil {
		t.Errorf("authentication middleware failed: %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader != "Bearer test-token" {
		t.Errorf("expected Authorization='Bearer test-token', got %s", authHeader)
	}
}

// TestAddRetry tests retry middleware
func TestAddRetry(t *testing.T) {
	middleware := AddRetry(3, func(resp *http.Response, err error) bool {
		return err != nil || resp.StatusCode >= 500
	})

	req := httptest.NewRequest("GET", "/test", nil)
	err := middleware(req)
	if err != nil {
		t.Errorf("retry middleware failed: %v", err)
	}

	// Check context values
	maxRetries := req.Context().Value(maxRetriesKey)
	if maxRetries != 3 {
		t.Errorf("expected maxRetries=3, got %v", maxRetries)
	}

	retryCondition := req.Context().Value(retryConditionKey)
	if retryCondition == nil {
		t.Error("retry condition not set in context")
	}
}

// TestAddRateLimit tests rate limit middleware
func TestAddRateLimit(t *testing.T) {
	provider := &testRateLimitInfoProvider{
		info: &RateLimitInfo{
			Bucket:    "test-bucket",
			Limit:     100,
			Remaining: 50,
		},
	}

	middleware := AddRateLimit(provider)
	req := httptest.NewRequest("GET", "/test", nil)

	err := middleware(req)
	if err != nil {
		t.Errorf("rate limit middleware failed: %v", err)
	}

	if req.Header.Get("X-RateLimit-Bucket") != "test-bucket" {
		t.Error("rate limit bucket header not set")
	}
	if req.Header.Get("X-RateLimit-Limit") != "100" {
		t.Error("rate limit limit header not set")
	}
	if req.Header.Get("X-RateLimit-Remaining") != "50" {
		t.Error("rate limit remaining header not set")
	}
}

// TestAddCache tests cache middleware
func TestAddCache(t *testing.T) {
	cache := NewInMemoryCache()
	middleware := AddResponseCache(cache)

	// Create test response
	req := httptest.NewRequest("GET", "/test", nil)
	responseBody := []byte(`{"message":"test"}`)
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Request:    req,
	}

	err := middleware(resp)
	if err != nil {
		t.Errorf("cache middleware failed: %v", err)
	}

	// Check if response was cached
	cacheKey := "/test"
	entry, found := cache.Get(cacheKey)
	if !found {
		t.Error("response not cached")
	}

	if string(entry.Data) != string(responseBody) {
		t.Errorf("cached data mismatch: expected %s, got %s", responseBody, entry.Data)
	}
}

// TestAddValidation tests validation middleware
func TestAddValidation(t *testing.T) {
	validator := &JSONValidator{}
	middleware := AddValidation(validator)

	tests := []struct {
		name        string
		body        string
		contentType string
		expectError bool
	}{
		{
			name:        "valid JSON",
			body:        `{"message":"test"}`,
			contentType: "application/json",
			expectError: false,
		},
		{
			name:        "invalid JSON",
			body:        `{invalid json}`,
			contentType: "application/json",
			expectError: true,
		},
		{
			name:        "non-JSON content",
			body:        "plain text",
			contentType: "text/plain",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", tt.contentType)
			req.ContentLength = int64(len(tt.body))

			err := middleware(req)
			if tt.expectError && err == nil {
				t.Error("expected validation error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

// TestAddCircuitBreaker tests circuit breaker middleware
func TestAddCircuitBreaker(t *testing.T) {
	cb := &testCircuitBreaker{allowRequest: true}
	middleware := AddCircuitBreaker(cb)

	req := httptest.NewRequest("GET", "/test", nil)
	err := middleware(req)
	if err != nil {
		t.Errorf("circuit breaker middleware failed: %v", err)
	}

	// Test when circuit breaker is open
	cb.allowRequest = false
	err = middleware(req)
	if err == nil {
		t.Error("expected circuit breaker error, got nil")
	}
}

// TestAddCircuitBreakerResponse tests circuit breaker response handling
func TestAddCircuitBreakerResponse(t *testing.T) {
	cb := &testCircuitBreaker{}
	middleware := AddResponseCircuitBreaker()

	// Test success recording
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), circuitBreakerKey, cb)
	req = req.WithContext(ctx)

	resp := &http.Response{
		StatusCode: 200,
		Request:    req,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}

	err := middleware(resp)
	if err != nil {
		t.Errorf("response middleware failed: %v", err)
	}

	if !cb.successRecorded {
		t.Error("success not recorded")
	}

	// Test failure recording
	cb = &testCircuitBreaker{}
	ctx = context.WithValue(req.Context(), circuitBreakerKey, cb)
	req = req.WithContext(ctx)

	resp.StatusCode = 500
	resp.Request = req

	err = middleware(resp)
	if err != nil {
		t.Errorf("response middleware failed: %v", err)
	}

	if !cb.failureRecorded {
		t.Error("failure not recorded")
	}
}

// TestInMemoryCache tests the in-memory cache implementation
func TestInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache()

	// Test Set and Get
	entry := CacheEntry{
		Data:       []byte("test data"),
		StatusCode: 200,
		Timestamp:  time.Now(),
	}

	cache.Set("key1", entry)
	retrieved, found := cache.Get("key1")
	if !found {
		t.Error("expected to find cached entry")
	}

	if string(retrieved.Data) != "test data" {
		t.Errorf("data mismatch: expected 'test data', got %s", retrieved.Data)
	}

	// Test TTL expiration
	entry.TTL = 10 * time.Millisecond
	entry.Timestamp = time.Now().Add(-20 * time.Millisecond)
	cache.Set("key2", entry)

	_, found = cache.Get("key2")
	if found {
		t.Error("expected expired entry to not be found")
	}

	// Test Delete
	cache.Delete("key1")
	_, found = cache.Get("key1")
	if found {
		t.Error("expected deleted entry to not be found")
	}

	// Test Clear
	cache.Set("key3", entry)
	cache.Set("key4", entry)
	cache.Clear()

	_, found = cache.Get("key3")
	if found {
		t.Error("expected cache to be cleared")
	}
}

// Test helper types

type testRateLimitInfoProvider struct {
	info *RateLimitInfo
}

func (p *testRateLimitInfoProvider) GetRateLimitInfo(path string) *RateLimitInfo {
	return p.info
}

type testCircuitBreaker struct {
	allowRequest    bool
	successRecorded bool
	failureRecorded bool
}

func (cb *testCircuitBreaker) AllowRequest() bool {
	return cb.allowRequest
}

func (cb *testCircuitBreaker) RecordSuccess() {
	cb.successRecorded = true
}

func (cb *testCircuitBreaker) RecordFailure() {
	cb.failureRecorded = true
}

func (cb *testCircuitBreaker) GetState() CircuitBreakerState {
	if cb.allowRequest {
		return CircuitBreakerClosed
	}
	return CircuitBreakerOpen
}

func (cb *testCircuitBreaker) Close() error {
	return nil
}
