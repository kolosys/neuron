package mock

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// AssertRequestCount fails the test if the request count doesn't match expected.
func AssertRequestCount(t *testing.T, rt *MockRoundTripper, expected int64) {
	t.Helper()
	if actual := rt.RequestCount(); actual != expected {
		t.Errorf("request count: expected %d, got %d", expected, actual)
	}
}

// AssertNoRequests fails the test if any requests were made.
func AssertNoRequests(t *testing.T, rt *MockRoundTripper) {
	t.Helper()
	if rt.RequestCount() > 0 {
		t.Errorf("expected no requests, but %d were made", rt.RequestCount())
	}
}

// AssertRequestMethod fails the test if the method of the i-th request doesn't match.
func AssertRequestMethod(t *testing.T, rt *MockRoundTripper, index int, method string) {
	t.Helper()
	requests := rt.Requests()

	if index >= len(requests) {
		t.Fatalf("request index out of bounds: %d >= %d", index, len(requests))
	}

	if requests[index].Method != method {
		t.Errorf("request %d method: expected %s, got %s", index, method, requests[index].Method)
	}
}

// AssertRequestURL fails the test if the URL of the i-th request doesn't match.
func AssertRequestURL(t *testing.T, rt *MockRoundTripper, index int, expectedURL string) {
	t.Helper()
	requests := rt.Requests()

	if index >= len(requests) {
		t.Fatalf("request index out of bounds: %d >= %d", index, len(requests))
	}

	actualURL := requests[index].URL.String()
	if actualURL != expectedURL {
		t.Errorf("request %d URL: expected %s, got %s", index, expectedURL, actualURL)
	}
}

// AssertRequestURLPath fails the test if the path of the i-th request doesn't match.
func AssertRequestURLPath(t *testing.T, rt *MockRoundTripper, index int, expectedPath string) {
	t.Helper()
	requests := rt.Requests()

	if index >= len(requests) {
		t.Fatalf("request index out of bounds: %d >= %d", index, len(requests))
	}

	actualPath := requests[index].URL.Path
	if actualPath != expectedPath {
		t.Errorf("request %d path: expected %s, got %s", index, expectedPath, actualPath)
	}
}

// AssertRequestHeader fails the test if a header value doesn't match.
func AssertRequestHeader(t *testing.T, rt *MockRoundTripper, index int, key, expectedValue string) {
	t.Helper()
	requests := rt.Requests()

	if index >= len(requests) {
		t.Fatalf("request index out of bounds: %d >= %d", index, len(requests))
	}

	actualValue := requests[index].Headers.Get(key)
	if actualValue != expectedValue {
		t.Errorf("request %d header %s: expected %s, got %s", index, key, expectedValue, actualValue)
	}
}

// AssertRequestBody fails the test if the body doesn't match the expected bytes.
func AssertRequestBody(t *testing.T, rt *MockRoundTripper, index int, expectedBody []byte) {
	t.Helper()
	requests := rt.Requests()

	if index >= len(requests) {
		t.Fatalf("request index out of bounds: %d >= %d", index, len(requests))
	}

	if string(requests[index].Body) != string(expectedBody) {
		t.Errorf("request %d body: expected %s, got %s", index, expectedBody, requests[index].Body)
	}
}

// AssertRequestBodyString fails the test if the body doesn't match the expected string.
func AssertRequestBodyString(t *testing.T, rt *MockRoundTripper, index int, expectedBody string) {
	t.Helper()
	AssertRequestBody(t, rt, index, []byte(expectedBody))
}

// AssertRequestBodyJSON fails the test if the body can't be unmarshaled as JSON matching the expected value.
func AssertRequestBodyJSON(t *testing.T, rt *MockRoundTripper, index int, expected any) {
	t.Helper()
	requests := rt.Requests()

	if index >= len(requests) {
		t.Fatalf("request index out of bounds: %d >= %d", index, len(requests))
	}

	var actual any
	if err := json.Unmarshal(requests[index].Body, &actual); err != nil {
		t.Fatalf("request %d body not valid JSON: %v", index, err)
	}

	actualJSON, _ := json.Marshal(actual)
	expectedJSON, _ := json.Marshal(expected)

	if string(actualJSON) != string(expectedJSON) {
		t.Errorf("request %d body JSON: expected %s, got %s", index, expectedJSON, actualJSON)
	}
}

// AssertAllowCalled fails the test if Allow was not called for the given method and endpoint.
func AssertAllowCalled(t *testing.T, rl *MockRateLimiter, method, endpoint string) {
	t.Helper()
	calls := rl.AllowCalls()

	for _, call := range calls {
		if call.Method == method && call.Endpoint == endpoint {
			return
		}
	}

	t.Errorf("Allow not called for %s %s", method, endpoint)
}

// AssertAllowNotCalled fails the test if Allow was called for the given method and endpoint.
func AssertAllowNotCalled(t *testing.T, rl *MockRateLimiter, method, endpoint string) {
	t.Helper()
	calls := rl.AllowCalls()

	for _, call := range calls {
		if call.Method == method && call.Endpoint == endpoint {
			t.Errorf("Allow was called for %s %s but should not have been", method, endpoint)
			return
		}
	}
}

// AssertAllowCount fails the test if the Allow call count doesn't match expected.
func AssertAllowCount(t *testing.T, rl *MockRateLimiter, expected int) {
	t.Helper()
	actual := len(rl.AllowCalls())
	if actual != expected {
		t.Errorf("Allow call count: expected %d, got %d", expected, actual)
	}
}

// AssertWaitCalled fails the test if Wait was not called for the given method and endpoint.
func AssertWaitCalled(t *testing.T, rl *MockRateLimiter, method, endpoint string) {
	t.Helper()
	calls := rl.WaitCalls()

	for _, call := range calls {
		if call.Method == method && call.Endpoint == endpoint {
			return
		}
	}

	t.Errorf("Wait not called for %s %s", method, endpoint)
}

// AssertWaitCount fails the test if the Wait call count doesn't match expected.
func AssertWaitCount(t *testing.T, rl *MockRateLimiter, expected int) {
	t.Helper()
	actual := len(rl.WaitCalls())
	if actual != expected {
		t.Errorf("Wait call count: expected %d, got %d", expected, actual)
	}
}

// AssertRateLimitUpdated fails the test if UpdateFromHeaders was not called for the endpoint.
func AssertRateLimitUpdated(t *testing.T, rlh *MockRateLimitHandler, method, endpoint string) {
	t.Helper()
	updates := rlh.UpdatesForEndpoint(method, endpoint)

	if len(updates) == 0 {
		t.Errorf("UpdateFromHeaders not called for %s %s", method, endpoint)
	}
}

// AssertRateLimitUpdateCount fails the test if the update count doesn't match expected.
func AssertRateLimitUpdateCount(t *testing.T, rlh *MockRateLimitHandler, expected int) {
	t.Helper()
	actual := len(rlh.Updates())
	if actual != expected {
		t.Errorf("rate limit update count: expected %d, got %d", expected, actual)
	}
}

// AssertRateLimitExhausted fails the test if a rate limit exhaustion was not recorded.
func AssertRateLimitExhausted(t *testing.T, rlh *MockRateLimitHandler) {
	t.Helper()
	if !rlh.WasExhausted() {
		t.Error("expected rate limit to have been exhausted")
	}
}

// AssertTokenCalled fails the test if GetToken was not called exactly count times.
func AssertTokenCalled(t *testing.T, ap *MockAuthProvider, count int) {
	t.Helper()
	actual := len(ap.GetTokenCalls())
	if actual != count {
		t.Errorf("GetToken call count: expected %d, got %d", count, actual)
	}
}

// AssertTokenValue fails the test if the i-th token call didn't return the expected token.
func AssertTokenValue(t *testing.T, ap *MockAuthProvider, index int, expected string) {
	t.Helper()
	calls := ap.GetTokenCalls()

	if index >= len(calls) {
		t.Fatalf("token call index out of bounds: %d >= %d", index, len(calls))
	}

	if calls[index].Result != expected {
		t.Errorf("token call %d: expected %s, got %s", index, expected, calls[index].Result)
	}
}

// AssertHeaderGenerated fails the test if GetAuthHeader was not called for the token.
func AssertHeaderGenerated(t *testing.T, ap *MockAuthProvider, token string, expectedHeader string) {
	t.Helper()
	calls := ap.GetHeaderCalls()

	for _, call := range calls {
		if call.Token == token {
			if call.Result != expectedHeader {
				t.Errorf("auth header for token %s: expected %s, got %s", token, expectedHeader, call.Result)
			}
			return
		}
	}

	t.Errorf("GetAuthHeader not called for token %s", token)
}

// AssertTokenRotated fails the test if the token was not rotated (multiple distinct tokens used).
func AssertTokenRotated(t *testing.T, ap *MockAuthProvider) {
	t.Helper()
	calls := ap.GetTokenCalls()

	if len(calls) < 2 {
		t.Errorf("expected at least 2 token calls for rotation, got %d", len(calls))
		return
	}

	// Check if any tokens are different
	firstToken := calls[0].Result
	for _, call := range calls[1:] {
		if call.Result != firstToken {
			return
		}
	}

	t.Error("token was not rotated (all tokens are the same)")
}

// AssertCacheHits fails the test if the hit count doesn't match expected.
func AssertCacheHits(t *testing.T, cache *MockCache, expected int64) {
	t.Helper()
	if actual := cache.Hits(); actual != expected {
		t.Errorf("cache hits: expected %d, got %d", expected, actual)
	}
}

// AssertCacheMisses fails the test if the miss count doesn't match expected.
func AssertCacheMisses(t *testing.T, cache *MockCache, expected int64) {
	t.Helper()
	if actual := cache.Misses(); actual != expected {
		t.Errorf("cache misses: expected %d, got %d", expected, actual)
	}
}

// AssertCacheSize fails the test if the cache size doesn't match expected.
func AssertCacheSize(t *testing.T, cache *MockCache, expected int) {
	t.Helper()
	if actual := cache.Size(); actual != expected {
		t.Errorf("cache size: expected %d, got %d", expected, actual)
	}
}

// AssertCacheContains fails the test if the key doesn't exist in the cache.
func AssertCacheContains(t *testing.T, cache *MockCache, key string) {
	t.Helper()
	if !cache.Contains(key) {
		t.Errorf("expected cache to contain key %s", key)
	}
}

// AssertCacheNotContains fails the test if the key exists in the cache.
func AssertCacheNotContains(t *testing.T, cache *MockCache, key string) {
	t.Helper()
	if cache.Contains(key) {
		t.Errorf("expected cache to not contain key %s", key)
	}
}

// AssertCacheEmpty fails the test if the cache is not empty.
func AssertCacheEmpty(t *testing.T, cache *MockCache) {
	t.Helper()
	if cache.Size() != 0 {
		t.Errorf("expected empty cache, but size is %d", cache.Size())
	}
}

// WaitForRequest polls until a request is received or timeout expires.
func WaitForRequest(t *testing.T, rt *MockRoundTripper, timeout time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if rt.RequestCount() > 0 {
			return true
		}
		time.Sleep(time.Millisecond)
	}

	t.Logf("timeout waiting for request")
	return false
}

// WaitForRequestCount polls until the request count reaches expected or timeout expires.
func WaitForRequestCount(t *testing.T, rt *MockRoundTripper, expected int64, timeout time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if rt.RequestCount() >= expected {
			return true
		}
		time.Sleep(time.Millisecond)
	}

	t.Logf("timeout waiting for %d requests (got %d)", expected, rt.RequestCount())
	return false
}

// AssertLastRequestMethod fails the test if the last request's method doesn't match.
func AssertLastRequestMethod(t *testing.T, rt *MockRoundTripper, method string) {
	t.Helper()
	requests := rt.Requests()

	if len(requests) == 0 {
		t.Fatal("no requests recorded")
	}

	lastIndex := len(requests) - 1
	AssertRequestMethod(t, rt, lastIndex, method)
}

// AssertLastRequestURL fails the test if the last request's URL doesn't match.
func AssertLastRequestURL(t *testing.T, rt *MockRoundTripper, expectedURL string) {
	t.Helper()
	requests := rt.Requests()

	if len(requests) == 0 {
		t.Fatal("no requests recorded")
	}

	lastIndex := len(requests) - 1
	AssertRequestURL(t, rt, lastIndex, expectedURL)
}

// AssertRequestsContainPath fails the test if none of the requests match the given path.
func AssertRequestsContainPath(t *testing.T, rt *MockRoundTripper, path string) {
	t.Helper()
	requests := rt.Requests()

	for _, req := range requests {
		if req.URL.Path == path {
			return
		}
	}

	t.Errorf("no requests found with path %s", path)
}

// AssertRequestCountBetween fails the test if the request count is not within the range.
func AssertRequestCountBetween(t *testing.T, rt *MockRoundTripper, min, max int64) {
	t.Helper()
	actual := rt.RequestCount()

	if actual < min || actual > max {
		t.Errorf("request count: expected between %d and %d, got %d", min, max, actual)
	}
}

// Helper function to format assertion error messages consistently
func formatAssertionMessage(testName, expected, actual any) string {
	return fmt.Sprintf("%s: expected %v, got %v", testName, expected, actual)
}

// AssertQueryParam fails the test if a query parameter doesn't match the expected value.
func AssertQueryParam(t *testing.T, rt *MockRoundTripper, index int, key, expectedValue string) {
	t.Helper()
	requests := rt.Requests()

	if index >= len(requests) {
		t.Fatalf("request index out of bounds: %d >= %d", index, len(requests))
	}

	actualValue := requests[index].URL.Query().Get(key)
	if actualValue != expectedValue {
		t.Errorf("request %d query param %s: expected %s, got %s", index, key, expectedValue, actualValue)
	}
}

// AssertRequestsForMethod returns requests matching the given HTTP method.
func AssertRequestsForMethod(t *testing.T, rt *MockRoundTripper, method string) []RecordedRequest {
	t.Helper()
	requests := rt.Requests()

	var filtered []RecordedRequest
	for _, req := range requests {
		if req.Method == method {
			filtered = append(filtered, req)
		}
	}

	return filtered
}

// AssertCacheHitRate fails the test if the hit rate doesn't match the expected value (within tolerance).
func AssertCacheHitRate(t *testing.T, cache *MockCache, expected float64, tolerance float64) {
	t.Helper()
	actual := cache.HitRate()

	if actual < expected-tolerance || actual > expected+tolerance {
		t.Errorf("cache hit rate: expected ~%.3f, got %.3f", expected, actual)
	}
}
