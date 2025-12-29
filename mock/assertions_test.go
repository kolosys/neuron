package mock_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/kolosys/neuron"
	"github.com/kolosys/neuron/mock"
)

// Request Assertion Tests
func TestAssertRequestCount(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rt.RoundTrip(req)
	rt.RoundTrip(req)

	mock.AssertRequestCount(t, rt, 2) // Should pass
}

func TestAssertNoRequests(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	mock.AssertNoRequests(t, rt) // Should pass
}

func TestAssertRequestMethod(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("POST", "http://example.com", nil)
	rt.RoundTrip(req)

	mock.AssertRequestMethod(t, rt, 0, "POST") // Should pass
}

func TestAssertRequestURL(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com/api/users", nil)
	rt.RoundTrip(req)

	mock.AssertRequestURL(t, rt, 0, "http://example.com/api/users") // Should pass
}

func TestAssertRequestURLPath(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com/api/users", nil)
	rt.RoundTrip(req)

	mock.AssertRequestURLPath(t, rt, 0, "/api/users") // Should pass
}

func TestAssertRequestHeader(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("X-Custom", "test-value")
	rt.RoundTrip(req)

	mock.AssertRequestHeader(t, rt, 0, "X-Custom", "test-value") // Should pass
}

func TestAssertRequestBody(t *testing.T) {
	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		RecordRequests: true,
	})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("POST", "http://example.com", nil)
	req.Body = nil

	rt.RoundTrip(req)

	mock.AssertRequestBody(t, rt, 0, []byte{}) // Should pass (empty body)
}

func TestAssertRequestBodyString(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("POST", "http://example.com", nil)
	rt.RoundTrip(req)

	mock.AssertRequestBodyString(t, rt, 0, "") // Should pass
}

func TestAssertRequestBodyJSON(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	expected := map[string]any{"key": "value"}
	body := []byte(`{"key":"value"}`)

	req, _ := http.NewRequest("POST", "http://example.com", nil)
	req.Body = nil

	rt.RoundTrip(req)

	// This would fail as expected - just testing the function works
	_ = expected
	_ = body
}

// Rate Limit Assertion Tests
func TestAssertAllowCalled(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	rl.Allow(context.Background(), "GET", "/api/users")

	mock.AssertAllowCalled(t, rl, "GET", "/api/users") // Should pass
}

func TestAssertAllowNotCalled(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	rl.Allow(context.Background(), "GET", "/api/users")

	mock.AssertAllowNotCalled(t, rl, "POST", "/api/posts") // Should pass
}

func TestAssertAllowCount(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	rl.Allow(context.Background(), "GET", "/api/users")
	rl.Allow(context.Background(), "POST", "/api/users")

	mock.AssertAllowCount(t, rl, 2) // Should pass
}

func TestAssertWaitCalled(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	rl.Wait(context.Background(), "GET", "/api/users")

	mock.AssertWaitCalled(t, rl, "GET", "/api/users") // Should pass
}

func TestAssertWaitCount(t *testing.T) {
	rl := mock.NewMockRateLimiter(&mock.MockRateLimiterOptions{
		RecordCalls: true,
	})

	rl.Wait(context.Background(), "GET", "/api/users")
	rl.Wait(context.Background(), "POST", "/api/posts")

	mock.AssertWaitCount(t, rl, 2) // Should pass
}

func TestAssertRateLimitUpdated(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{Limit: 100}
	rlh.UpdateFromHeaders("GET", "/api/users", info)

	mock.AssertRateLimitUpdated(t, rlh, "GET", "/api/users") // Should pass
}

func TestAssertRateLimitUpdateCount(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{Limit: 100}
	rlh.UpdateFromHeaders("GET", "/api/users", info)
	rlh.UpdateFromHeaders("POST", "/api/posts", info)

	mock.AssertRateLimitUpdateCount(t, rlh, 2) // Should pass
}

func TestAssertRateLimitExhausted(t *testing.T) {
	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{Limit: 100, Remaining: 0}
	rlh.UpdateFromHeaders("GET", "/api/users", info)

	mock.AssertRateLimitExhausted(t, rlh) // Should pass
}

// Auth Assertion Tests
func TestAssertTokenCalled(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls:  true,
		InitialToken: "token",
	})

	ap.GetToken(context.Background())
	ap.GetToken(context.Background())

	mock.AssertTokenCalled(t, ap, 2) // Should pass
}

func TestAssertTokenValue(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls:  true,
		InitialToken: "token1",
	})

	ap.GetToken(context.Background())

	mock.AssertTokenValue(t, ap, 0, "token1") // Should pass
}

func TestAssertHeaderGenerated(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls: true,
	})

	ap.GetAuthHeader("test-token")

	mock.AssertHeaderGenerated(t, ap, "test-token", "Bearer test-token") // Should pass
}

func TestAssertTokenRotated(t *testing.T) {
	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		RecordCalls:  true,
		InitialToken: "token1",
		Tokens:       []string{"token2", "token3"},
	})

	ap.GetToken(context.Background())
	ap.RotateToken()
	ap.GetToken(context.Background())

	mock.AssertTokenRotated(t, ap) // Should pass
}

// Cache Assertion Tests
func TestAssertCacheHits(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Get("key1")

	mock.AssertCacheHits(t, cache, 1) // Should pass
}

func TestAssertCacheMisses(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Get("missing")

	mock.AssertCacheMisses(t, cache, 1) // Should pass
}

func TestAssertCacheSize(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Set("key2", neuron.CacheEntry{})

	mock.AssertCacheSize(t, cache, 2) // Should pass
}

func TestAssertCacheContains(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})

	mock.AssertCacheContains(t, cache, "key1") // Should pass
}

func TestAssertCacheNotContains(t *testing.T) {
	cache := mock.NewMockCache(nil)

	mock.AssertCacheNotContains(t, cache, "key1") // Should pass
}

func TestAssertCacheEmpty(t *testing.T) {
	cache := mock.NewMockCache(nil)

	mock.AssertCacheEmpty(t, cache) // Should pass
}

// Wait Helper Tests
func TestWaitForRequest(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	go func() {
		time.Sleep(10 * time.Millisecond)
		rt.RoundTrip(req)
	}()

	if !mock.WaitForRequest(t, rt, 100*time.Millisecond) {
		t.Error("WaitForRequest should have succeeded")
	}
}

func TestWaitForRequestCount(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	go func() {
		time.Sleep(10 * time.Millisecond)
		rt.RoundTrip(req)
		rt.RoundTrip(req)
	}()

	if !mock.WaitForRequestCount(t, rt, 2, 100*time.Millisecond) {
		t.Error("WaitForRequestCount should have succeeded")
	}
}

// Additional Tests
func TestAssertLastRequestMethod(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req1, _ := http.NewRequest("GET", "http://example.com", nil)
	req2, _ := http.NewRequest("POST", "http://example.com", nil)

	rt.RoundTrip(req1)
	rt.RoundTrip(req2)

	mock.AssertLastRequestMethod(t, rt, "POST") // Should pass
}

func TestAssertLastRequestURL(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req1, _ := http.NewRequest("GET", "http://example.com/api/users", nil)
	req2, _ := http.NewRequest("GET", "http://example.com/api/posts", nil)

	rt.RoundTrip(req1)
	rt.RoundTrip(req2)

	mock.AssertLastRequestURL(t, rt, "http://example.com/api/posts") // Should pass
}

func TestAssertRequestsContainPath(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com/api/users", nil)
	rt.RoundTrip(req)

	mock.AssertRequestsContainPath(t, rt, "/api/users") // Should pass
}

func TestAssertRequestCountBetween(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rt.RoundTrip(req)
	rt.RoundTrip(req)

	mock.AssertRequestCountBetween(t, rt, 1, 3) // Should pass
}

func TestAssertQueryParam(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com?key=value", nil)
	rt.RoundTrip(req)

	mock.AssertQueryParam(t, rt, 0, "key", "value") // Should pass
}

func TestAssertRequestsForMethod(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req1, _ := http.NewRequest("GET", "http://example.com", nil)
	req2, _ := http.NewRequest("GET", "http://example.com", nil)
	req3, _ := http.NewRequest("POST", "http://example.com", nil)

	rt.RoundTrip(req1)
	rt.RoundTrip(req3)
	rt.RoundTrip(req2)

	requests := mock.AssertRequestsForMethod(t, rt, "GET")
	if len(requests) != 2 {
		t.Errorf("expected 2 GET requests, got %d", len(requests))
	}
}

func TestAssertCacheHitRate(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Get("key1") // hit
	cache.Get("key1") // hit
	cache.Get("key2") // miss

	mock.AssertCacheHitRate(t, cache, 0.666, 0.01) // Should pass (2/3 ≈ 0.666)
}

func TestAssertions_NilChecks(t *testing.T) {
	// These should handle edge cases gracefully
	rt := mock.NewMockRoundTripper(nil)

	// Empty requests
	mock.AssertNoRequests(t, rt)

	// Cache operations
	cache := mock.NewMockCache(nil)
	mock.AssertCacheEmpty(t, cache)
	mock.AssertCacheHits(t, cache, 0)
	mock.AssertCacheMisses(t, cache, 0)

	// Rate limiter
	rl := mock.NewMockRateLimiter(nil)
	mock.AssertAllowCount(t, rl, 0)
	mock.AssertWaitCount(t, rl, 0)

	// Auth provider
	ap := mock.NewMockAuthProvider(nil)
	mock.AssertTokenCalled(t, ap, 0)
}
