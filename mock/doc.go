// Package mock provides testing utilities and mocks for the neuron HTTP client library.
//
// This package offers comprehensive mocks for all neuron interfaces and the HTTP
// transport layer, enabling isolated unit testing without making real HTTP requests.
//
// # Overview
//
// The mock package provides mocks for:
//   - HTTP transport layer (MockRoundTripper)
//   - Rate limiting (MockRateLimiter, MockRateLimitHandler)
//   - Authentication (MockAuthProvider)
//   - Caching (MockCache)
//   - Request validation (MockValidator)
//   - Request ID generation (MockRequestIDGenerator)
//   - Body serialization (MockBodyProvider)
//   - Utilities for recording, error injection, and assertions
//
// # Quick Start
//
// The simplest way to test code using neuron is to use MockRoundTripper:
//
//	rt := mock.NewMockRoundTripper(nil)
//	rt.QueueResponse(mock.ResponseConfig{
//		StatusCode: 200,
//		Body:       []byte(`{"status":"ok"}`),
//	})
//
//	client := neuron.NewClient(neuron.ClientOptions{
//		BaseURL:    "http://api.example.com",
//		HTTPClient: &http.Client{Transport: rt},
//	})
//
//	// Your code using client here...
//
//	// Verify the request
//	mock.AssertRequestCount(t, rt, 1)
//	mock.AssertRequestMethod(t, rt, 0, "GET")
//
// # Thread Safety
//
// All mocks in this package are designed to be thread-safe. They use atomic
// operations for counters and sync.RWMutex for protecting shared state. This
// allows for concurrent testing scenarios without data races.
//
// Run tests with the -race flag to verify thread safety:
//
//	go test -race ./...
//
// # Error Injection
//
// All mocks support error injection for testing error handling paths. By default,
// injected errors are consumed after first use (one-shot). This allows testing
// retry logic:
//
//	rt := mock.NewMockRoundTripper(nil)
//	rt.InjectError(io.ErrUnexpectedEOF)
//	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
//
//	req, _ := http.NewRequest("GET", "http://example.com", nil)
//	resp, err := rt.RoundTrip(req)
//	// First call returns error
//	if err != io.ErrUnexpectedEOF { t.Fatal() }
//
//	resp, err = rt.RoundTrip(req)
//	// Second call succeeds (error was cleared)
//	if err != nil || resp.StatusCode != 200 { t.Fatal() }
//
// # Rate Limiting
//
// Test rate limiting behavior without actual delays:
//
//	rl := mock.NewMockRateLimiter(nil)
//	rl.SetAllow(false)  // Block all requests
//
//	client := neuron.NewClient(neuron.ClientOptions{
//		RateLimiter: rl,
//	})
//
//	mock.AssertWaitCalled(t, rl, "GET", "/api/data")
//
// Or test rate limit header parsing:
//
//	rlh := mock.NewMockRateLimitHandler(nil)
//
//	info := &neuron.RateLimitInfo{
//		Limit:     100,
//		Remaining: 10,
//	}
//	rlh.UpdateFromHeaders("GET", "/api/users", info)
//
//	mock.AssertRateLimitUpdated(t, rlh, "GET", "/api/users")
//
// # Authentication
//
// Test token rotation and auth provider behavior:
//
//	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
//		InitialToken: "token1",
//		Tokens:       []string{"token2", "token3"},
//	})
//
//	token1, _ := ap.GetToken(context.Background())
//	ap.RotateToken()
//	token2, _ := ap.GetToken(context.Background())
//
//	mock.AssertTokenValue(t, ap, 0, "token1")
//	mock.AssertTokenValue(t, ap, 1, "token2")
//	mock.AssertTokenRotated(t, ap)
//
// # Caching
//
// Test cache behavior and hit/miss patterns:
//
//	cache := mock.NewMockCache(nil)
//	entry := neuron.CacheEntry{Data: []byte("cached"), StatusCode: 200}
//
//	cache.Set("key1", entry)
//	cache.Get("key1")  // hit
//	cache.Get("key2")  // miss
//
//	mock.AssertCacheHits(t, cache, 1)
//	mock.AssertCacheMisses(t, cache, 1)
//	mock.AssertCacheHitRate(t, cache, 0.5, 0.01)
//
// # Comprehensive Example
//
// Here's a complete example combining multiple mocks:
//
//	func TestCompleteWorkflow(t *testing.T) {
//		// Setup mocks
//		rt := mock.NewMockRoundTripper(nil)
//		rl := mock.NewMockRateLimiter(nil)
//		ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
//			InitialToken: "secret-token",
//		})
//		cache := mock.NewMockCache(nil)
//
//		// Configure responses
//		rt.QueueResponse(mock.ResponseConfig{
//			StatusCode: 200,
//			Body:       []byte(`[{"id":1,"name":"Alice"}]`),
//			Headers:    http.Header{"Content-Type": []string{"application/json"}},
//		})
//
//		// Create client with mocks
//		client := neuron.NewClient(neuron.ClientOptions{
//			BaseURL:    "http://api.example.com",
//			HTTPClient: &http.Client{Transport: rt},
//			RateLimiter: rl,
//		})
//
//		// Add auth hook
//		client.AddRequestHook(func(req *http.Request) error {
//			token, _ := ap.GetToken(req.Context())
//			req.Header.Set("Authorization", "Bearer "+token)
//			return nil
//		})
//
//		// Test your code
//		// ...
//
//		// Verify interactions
//		mock.AssertRequestCount(t, rt, 1)
//		mock.AssertRequestHeader(t, rt, 0, "Authorization", "Bearer secret-token")
//		mock.AssertTokenCalled(t, ap, 1)
//		mock.AssertAllowCalled(t, rl, "GET", "/api/users")
//	}
//
// # Recording and Statistics
//
// All mocks can record their operations for inspection:
//
//	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
//		RecordRequests: true,
//	})
//
//	// Make requests...
//	requests := rt.Requests()
//	for i, req := range requests {
//		fmt.Printf("Request %d: %s %s\n", i, req.Method, req.URL)
//	}
//
// Cache mocks provide detailed statistics:
//
//	cache := mock.NewMockCache(nil)
//	// ... cache operations ...
//	fmt.Printf("Hits: %d, Misses: %d, Rate: %.2f%%\n",
//		cache.Hits(), cache.Misses(), cache.HitRate()*100)
//
// # Assertion Helpers
//
// The package provides 40+ assertion functions for common test patterns:
//
//	// Request assertions
//	mock.AssertRequestCount(t, rt, 1)
//	mock.AssertRequestMethod(t, rt, 0, "POST")
//	mock.AssertRequestURL(t, rt, 0, "http://example.com/api/users")
//	mock.AssertRequestHeader(t, rt, 0, "Content-Type", "application/json")
//	mock.AssertRequestBody(t, rt, 0, expectedBody)
//	mock.AssertRequestBodyJSON(t, rt, 0, expectedJSON)
//
//	// Rate limit assertions
//	mock.AssertAllowCalled(t, rl, "GET", "/api/data")
//	mock.AssertWaitCalled(t, rl, "POST", "/api/users")
//	mock.AssertRateLimitUpdated(t, rlh, "GET", "/api/data")
//	mock.AssertRateLimitExhausted(t, rlh)
//
//	// Auth assertions
//	mock.AssertTokenCalled(t, ap, 2)
//	mock.AssertTokenValue(t, ap, 0, "expected-token")
//	mock.AssertTokenRotated(t, ap)
//
//	// Cache assertions
//	mock.AssertCacheHits(t, cache, 5)
//	mock.AssertCacheMisses(t, cache, 2)
//	mock.AssertCacheSize(t, cache, 3)
//	mock.AssertCacheContains(t, cache, "key1")
//
//	// Wait helpers
//	if !mock.WaitForRequest(t, rt, 100*time.Millisecond) {
//		t.Fatal("request never arrived")
//	}
//
// # Best Practices
//
//  1. Always enable recording during development - helps debug test failures
//  2. Use specific assertions rather than just checking request count
//  3. Test error paths with error injection
//  4. Disable recording in performance-critical tests
//  5. Use WaitFor helpers for async scenarios
//  6. Call Reset() or ClearRecorded() between test cases
//  7. Run tests with -race flag to catch concurrency issues
//
// # Latency Simulation
//
// Simulate network latency for timeout testing:
//
//	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
//		Latency: 50*time.Millisecond,
//	})
//
//	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
//
//	start := time.Now()
//	rt.RoundTrip(req)
//	// Will take at least 50ms
//
// # Callback Recording
//
// Record callback invocations for testing hooks:
//
//	recorder := mock.NewCallbackRecorder()
//	client.OnError(recorder.Record)
//
//	// ... trigger errors ...
//
//	mock.AssertCallCount(t, recorder, 2)
//
// # Response Matching
//
// Use pattern matchers for complex multi-endpoint scenarios:
//
//	rt := mock.NewMockRoundTripper(nil)
//	rt.AddMatcher(func(req *http.Request) (*http.Response, bool) {
//		if req.URL.Path == "/api/users" {
//			return &http.Response{
//				StatusCode: 200,
//				Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
//				Header:     make(http.Header),
//			}, true
//		}
//		return nil, false
//	})
//
// # Performance Testing
//
// Disable recording for performance-critical tests:
//
//	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
//		RecordRequests: false,
//	})
//
//	// Allocations are minimal without recording
package mock
