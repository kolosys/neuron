# mock API

Complete API documentation for the mock package.

**Import Path:** `github.com/kolosys/neuron/mock`

## Package Documentation

Package mock provides testing utilities and mocks for the neuron HTTP client library.

This package offers comprehensive mocks for all neuron interfaces and the HTTP
transport layer, enabling isolated unit testing without making real HTTP requests.

# Overview

The mock package provides mocks for:
  - HTTP transport layer (MockRoundTripper)
  - Rate limiting (MockRateLimiter, MockRateLimitHandler)
  - Authentication (MockAuthProvider)
  - Caching (MockCache)
  - Request validation (MockValidator)
  - Request ID generation (MockRequestIDGenerator)
  - Body serialization (MockBodyProvider)
  - Utilities for recording, error injection, and assertions

# Quick Start

The simplest way to test code using neuron is to use MockRoundTripper:

	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{
		StatusCode: 200,
		Body:       []byte(`{"status":"ok"}`),
	})

	client := neuron.NewClient(neuron.ClientOptions{
		BaseURL:    "http://api.example.com",
		HTTPClient: &http.Client{Transport: rt},
	})

	// Your code using client here...

	// Verify the request
	mock.AssertRequestCount(t, rt, 1)
	mock.AssertRequestMethod(t, rt, 0, "GET")

# Thread Safety

All mocks in this package are designed to be thread-safe. They use atomic
operations for counters and sync.RWMutex for protecting shared state. This
allows for concurrent testing scenarios without data races.

Run tests with the -race flag to verify thread safety:

	go test -race ./...

# Error Injection

All mocks support error injection for testing error handling paths. By default,
injected errors are consumed after first use (one-shot). This allows testing
retry logic:

	rt := mock.NewMockRoundTripper(nil)
	rt.InjectError(io.ErrUnexpectedEOF)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	resp, err := rt.RoundTrip(req)
	// First call returns error
	if err != io.ErrUnexpectedEOF { t.Fatal() }

	resp, err = rt.RoundTrip(req)
	// Second call succeeds (error was cleared)
	if err != nil || resp.StatusCode != 200 { t.Fatal() }

# Rate Limiting

Test rate limiting behavior without actual delays:

	rl := mock.NewMockRateLimiter(nil)
	rl.SetAllow(false)  // Block all requests

	client := neuron.NewClient(neuron.ClientOptions{
		RateLimiter: rl,
	})

	mock.AssertWaitCalled(t, rl, "GET", "/api/data")

Or test rate limit header parsing:

	rlh := mock.NewMockRateLimitHandler(nil)

	info := &neuron.RateLimitInfo{
		Limit:     100,
		Remaining: 10,
	}
	rlh.UpdateFromHeaders("GET", "/api/users", info)

	mock.AssertRateLimitUpdated(t, rlh, "GET", "/api/users")

# Authentication

Test token rotation and auth provider behavior:

	ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
		InitialToken: "token1",
		Tokens:       []string{"token2", "token3"},
	})

	token1, _ := ap.GetToken(context.Background())
	ap.RotateToken()
	token2, _ := ap.GetToken(context.Background())

	mock.AssertTokenValue(t, ap, 0, "token1")
	mock.AssertTokenValue(t, ap, 1, "token2")
	mock.AssertTokenRotated(t, ap)

# Caching

Test cache behavior and hit/miss patterns:

	cache := mock.NewMockCache(nil)
	entry := neuron.CacheEntry{Data: []byte("cached"), StatusCode: 200}

	cache.Set("key1", entry)
	cache.Get("key1")  // hit
	cache.Get("key2")  // miss

	mock.AssertCacheHits(t, cache, 1)
	mock.AssertCacheMisses(t, cache, 1)
	mock.AssertCacheHitRate(t, cache, 0.5, 0.01)

# Comprehensive Example

Here's a complete example combining multiple mocks:

	func TestCompleteWorkflow(t *testing.T) {
		// Setup mocks
		rt := mock.NewMockRoundTripper(nil)
		rl := mock.NewMockRateLimiter(nil)
		ap := mock.NewMockAuthProvider(&mock.MockAuthProviderOptions{
			InitialToken: "secret-token",
		})
		cache := mock.NewMockCache(nil)

		// Configure responses
		rt.QueueResponse(mock.ResponseConfig{
			StatusCode: 200,
			Body:       []byte(`[{"id":1,"name":"Alice"}]`),
			Headers:    http.Header{"Content-Type": []string{"application/json"}},
		})

		// Create client with mocks
		client := neuron.NewClient(neuron.ClientOptions{
			BaseURL:    "http://api.example.com",
			HTTPClient: &http.Client{Transport: rt},
			RateLimiter: rl,
		})

		// Add auth hook
		client.AddRequestHook(func(req *http.Request) error {
			token, _ := ap.GetToken(req.Context())
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		})

		// Test your code
		// ...

		// Verify interactions
		mock.AssertRequestCount(t, rt, 1)
		mock.AssertRequestHeader(t, rt, 0, "Authorization", "Bearer secret-token")
		mock.AssertTokenCalled(t, ap, 1)
		mock.AssertAllowCalled(t, rl, "GET", "/api/users")
	}

# Recording and Statistics

All mocks can record their operations for inspection:

	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		RecordRequests: true,
	})

	// Make requests...
	requests := rt.Requests()
	for i, req := range requests {
		fmt.Printf("Request %d: %s %s\n", i, req.Method, req.URL)
	}

Cache mocks provide detailed statistics:

	cache := mock.NewMockCache(nil)
	// ... cache operations ...
	fmt.Printf("Hits: %d, Misses: %d, Rate: %.2f%%\n",
		cache.Hits(), cache.Misses(), cache.HitRate()*100)

# Assertion Helpers

The package provides 40+ assertion functions for common test patterns:

	// Request assertions
	mock.AssertRequestCount(t, rt, 1)
	mock.AssertRequestMethod(t, rt, 0, "POST")
	mock.AssertRequestURL(t, rt, 0, "http://example.com/api/users")
	mock.AssertRequestHeader(t, rt, 0, "Content-Type", "application/json")
	mock.AssertRequestBody(t, rt, 0, expectedBody)
	mock.AssertRequestBodyJSON(t, rt, 0, expectedJSON)

	// Rate limit assertions
	mock.AssertAllowCalled(t, rl, "GET", "/api/data")
	mock.AssertWaitCalled(t, rl, "POST", "/api/users")
	mock.AssertRateLimitUpdated(t, rlh, "GET", "/api/data")
	mock.AssertRateLimitExhausted(t, rlh)

	// Auth assertions
	mock.AssertTokenCalled(t, ap, 2)
	mock.AssertTokenValue(t, ap, 0, "expected-token")
	mock.AssertTokenRotated(t, ap)

	// Cache assertions
	mock.AssertCacheHits(t, cache, 5)
	mock.AssertCacheMisses(t, cache, 2)
	mock.AssertCacheSize(t, cache, 3)
	mock.AssertCacheContains(t, cache, "key1")

	// Wait helpers
	if !mock.WaitForRequest(t, rt, 100*time.Millisecond) {
		t.Fatal("request never arrived")
	}

# Best Practices

 1. Always enable recording during development - helps debug test failures
 2. Use specific assertions rather than just checking request count
 3. Test error paths with error injection
 4. Disable recording in performance-critical tests
 5. Use WaitFor helpers for async scenarios
 6. Call Reset() or ClearRecorded() between test cases
 7. Run tests with -race flag to catch concurrency issues

# Latency Simulation

Simulate network latency for timeout testing:

	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		Latency: 50*time.Millisecond,
	})

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	start := time.Now()
	rt.RoundTrip(req)
	// Will take at least 50ms

# Callback Recording

Record callback invocations for testing hooks:

	recorder := mock.NewCallbackRecorder()
	client.OnError(recorder.Record)

	// ... trigger errors ...

	mock.AssertCallCount(t, recorder, 2)

# Response Matching

Use pattern matchers for complex multi-endpoint scenarios:

	rt := mock.NewMockRoundTripper(nil)
	rt.AddMatcher(func(req *http.Request) (*http.Response, bool) {
		if req.URL.Path == "/api/users" {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`[]`))),
				Header:     make(http.Header),
			}, true
		}
		return nil, false
	})

# Performance Testing

Disable recording for performance-critical tests:

	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		RecordRequests: false,
	})

	// Allocations are minimal without recording


## Types

### AllowCall
AllowCall records a call to Allow.

#### Example Usage

```go
// Create a new AllowCall
allowcall := AllowCall{
    Method: "example",
    Endpoint: "example",
    Time: /* value */,
    Result: true,
}
```

#### Type Definition

```go
type AllowCall struct {
    Method string
    Endpoint string
    Time time.Time
    Result bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Method | `string` |  |
| Endpoint | `string` |  |
| Time | `time.Time` |  |
| Result | `bool` |  |

### CacheOperation
CacheOperation records a cache operation for analysis.

#### Example Usage

```go
// Create a new CacheOperation
cacheoperation := CacheOperation{
    Op: "example",
    Key: "example",
    Time: 42,
    Hit: true,
    Entry: any{},
}
```

#### Type Definition

```go
type CacheOperation struct {
    Op string
    Key string
    Time int64
    Hit bool
    Entry any
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Op | `string` | "get", "set", "delete", "clear" |
| Key | `string` |  |
| Time | `int64` | nanoseconds since epoch |
| Hit | `bool` | for Get operations |
| Entry | `any` | the entry involved |

### CallbackRecorder
CallbackRecorder records the number of times a callback has been invoked. It's designed to be called as a callback function and provides assertion helpers.

#### Example Usage

```go
// Create a new CallbackRecorder
callbackrecorder := CallbackRecorder{

}
```

#### Type Definition

```go
type CallbackRecorder struct {
}
```

### Constructor Functions

### NewCallbackRecorder

NewCallbackRecorder creates a new callback recorder.

```go
func NewCallbackRecorder() *CallbackRecorder
```

**Parameters:**
  None

**Returns:**
- *CallbackRecorder

## Methods

### AssertCallCount

AssertCallCount fails the test if the call count doesn't match the expected value.

```go
func (*CallbackRecorder) AssertCallCount(t *testing.T, expected int64)
```

**Parameters:**
- `t` (*testing.T)
- `expected` (int64)

**Returns:**
  None

### AssertCalled

AssertCalled fails the test if the callback was never called.

```go
func (*CallbackRecorder) AssertCalled(t *testing.T)
```

**Parameters:**
- `t` (*testing.T)

**Returns:**
  None

### AssertNotCalled

AssertNotCalled fails the test if the callback was called at least once.

```go
func (*CallbackRecorder) AssertNotCalled(t *testing.T)
```

**Parameters:**
- `t` (*testing.T)

**Returns:**
  None

### Calls

Calls returns the current call count.

```go
func (*CallbackRecorder) Calls() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Record

Record increments the call counter. This is designed to be used as a callback.

```go
func (*CallbackRecorder) Record()
```

**Parameters:**
  None

**Returns:**
  None

### Reset

Reset resets the call counter to zero.

```go
func (*MockCache) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### CommonErrors
CommonErrors provides sentinel errors for common test scenarios.

#### Example Usage

```go
// Create a new CommonErrors
commonerrors := CommonErrors{
    Timeout: error{},
    NetworkError: error{},
    DNSError: error{},
    TLSError: error{},
    ContextCanceled: error{},
    RateLimit: error{},
    AuthFailed: error{},
}
```

#### Type Definition

```go
type CommonErrors struct {
    Timeout error
    NetworkError error
    DNSError error
    TLSError error
    ContextCanceled error
    RateLimit error
    AuthFailed error
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Timeout | `error` |  |
| NetworkError | `error` |  |
| DNSError | `error` |  |
| TLSError | `error` |  |
| ContextCanceled | `error` |  |
| RateLimit | `error` |  |
| AuthFailed | `error` |  |

### Constructor Functions

### DefaultErrors

DefaultErrors returns a CommonErrors struct with standard test errors.

```go
func DefaultErrors() *CommonErrors
```

**Parameters:**
  None

**Returns:**
- *CommonErrors

### ErrorSequence
ErrorSequence provides a sequence of errors for testing retry behavior. Each call to Next() returns the next error in the sequence.

#### Example Usage

```go
// Create a new ErrorSequence
errorsequence := ErrorSequence{

}
```

#### Type Definition

```go
type ErrorSequence struct {
}
```

### Constructor Functions

### NewErrorSequence

NewErrorSequence creates a new error sequence from the provided errors.

```go
func NewErrorSequence(errs ...error) *ErrorSequence
```

**Parameters:**
- `errs` (...error)

**Returns:**
- *ErrorSequence

## Methods

### Error

Error implements the error interface.

```go
func (*ErrorSequence) Error() string
```

**Parameters:**
  None

**Returns:**
- string

### Next

Next returns the next error in the sequence. If the end is reached, it returns the last error.

```go
func (*ErrorSequence) Next() error
```

**Parameters:**
  None

**Returns:**
- error

### Reset

Reset resets the sequence to the first error.

```go
func (*MockCache) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### GetHeaderCall
GetHeaderCall records a call to GetAuthHeader.

#### Example Usage

```go
// Create a new GetHeaderCall
getheadercall := GetHeaderCall{
    Token: "example",
    Result: "example",
    Time: /* value */,
}
```

#### Type Definition

```go
type GetHeaderCall struct {
    Token string
    Result string
    Time time.Time
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Token | `string` |  |
| Result | `string` |  |
| Time | `time.Time` |  |

### GetTokenCall
GetTokenCall records a call to GetToken.

#### Example Usage

```go
// Create a new GetTokenCall
gettokencall := GetTokenCall{
    Time: /* value */,
    Result: "example",
    Err: error{},
}
```

#### Type Definition

```go
type GetTokenCall struct {
    Time time.Time
    Result string
    Err error
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Time | `time.Time` |  |
| Result | `string` |  |
| Err | `error` |  |

### MockAuthProvider
MockAuthProvider is a mock implementation of neuron.AuthProvider for testing.

#### Example Usage

```go
// Create a new MockAuthProvider
mockauthprovider := MockAuthProvider{

}
```

#### Type Definition

```go
type MockAuthProvider struct {
}
```

### Constructor Functions

### NewMockAuthProvider

NewMockAuthProvider creates a new mock auth provider.

```go
func NewMockAuthProvider(opts *MockAuthProviderOptions) *MockAuthProvider
```

**Parameters:**
- `opts` (*MockAuthProviderOptions)

**Returns:**
- *MockAuthProvider

## Methods

### ClearInjectedErrors

ClearInjectedErrors clears all injected errors.

```go
func (*MockRateLimiter) ClearInjectedErrors()
```

**Parameters:**
  None

**Returns:**
  None

### ClearRecorded

ClearRecorded clears all recorded calls.

```go
func (*MockValidator) ClearRecorded()
```

**Parameters:**
  None

**Returns:**
  None

### CurrentTokenIndex

CurrentTokenIndex returns the index of the current token.

```go
func (*MockAuthProvider) CurrentTokenIndex() int
```

**Parameters:**
  None

**Returns:**
- int

### GetAuthHeader

GetAuthHeader returns the formatted authentication header value.

```go
func (*MockAuthProvider) GetAuthHeader(token string) string
```

**Parameters:**
- `token` (string)

**Returns:**
- string

### GetHeaderCalls

GetHeaderCalls returns a copy of all recorded GetAuthHeader calls.

```go
func (*MockAuthProvider) GetHeaderCalls() []GetHeaderCall
```

**Parameters:**
  None

**Returns:**
- []GetHeaderCall

### GetToken

GetToken returns the current token, or the next token in the sequence if configured.

```go
func (*MockAuthProvider) GetToken(ctx context.Context) (string, error)
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- string
- error

### GetTokenCalls

GetTokenCalls returns a copy of all recorded GetToken calls.

```go
func (*MockAuthProvider) GetTokenCalls() []GetTokenCall
```

**Parameters:**
  None

**Returns:**
- []GetTokenCall

### InjectTokenError

InjectTokenError injects an error for the next GetToken call (one-shot by default).

```go
func (*MockAuthProvider) InjectTokenError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### Reset

Reset resets the token sequence to the first token.

```go
func (*MockAuthProvider) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### RotateToken

RotateToken advances to the next token in the sequence.

```go
func (*MockAuthProvider) RotateToken()
```

**Parameters:**
  None

**Returns:**
  None

### SetHeaderFormat

SetHeaderFormat sets the format for GetAuthHeader. Use "{}" as a placeholder for the token.

```go
func (*MockAuthProvider) SetHeaderFormat(format string)
```

**Parameters:**
- `format` (string)

**Returns:**
  None

### SetTokens

SetTokens sets the token sequence for rotation testing.

```go
func (*MockAuthProvider) SetTokens(tokens []string)
```

**Parameters:**
- `tokens` ([]string)

**Returns:**
  None

### MockAuthProviderOptions
MockAuthProviderOptions configures mock auth provider behavior.

#### Example Usage

```go
// Create a new MockAuthProviderOptions
mockauthprovideroptions := MockAuthProviderOptions{
    InitialToken: "example",
    RecordCalls: true,
    Tokens: [],
    HeaderFormat: "example",
}
```

#### Type Definition

```go
type MockAuthProviderOptions struct {
    InitialToken string
    RecordCalls bool
    Tokens []string
    HeaderFormat string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| InitialToken | `string` |  |
| RecordCalls | `bool` |  |
| Tokens | `[]string` |  |
| HeaderFormat | `string` | e.g., "Bearer {}" or "X-API-Key {}" |

### MockBodyProvider
MockBodyProvider is a mock implementation of neuron.BodyProvider for testing.

#### Example Usage

```go
// Create a new MockBodyProvider
mockbodyprovider := MockBodyProvider{

}
```

#### Type Definition

```go
type MockBodyProvider struct {
}
```

### Constructor Functions

### NewMockBodyProvider

NewMockBodyProvider creates a new mock body provider.

```go
func NewMockBodyProvider(opts *MockBodyProviderOptions) *MockBodyProvider
```

**Parameters:**
- `opts` (*MockBodyProviderOptions)

**Returns:**
- *MockBodyProvider

## Methods

### Body

Body returns the body as an io.Reader.

```go
func (*MockBodyProvider) Body() (io.Reader, error)
```

**Parameters:**
  None

**Returns:**
- io.Reader
- error

### CallCount

CallCount returns the number of times Body was called.

```go
func (*MockBodyProvider) CallCount() int64
```

**Parameters:**
  None

**Returns:**
- int64

### ClearInjectedErrors

ClearInjectedErrors clears all injected errors.

```go
func (*MockBodyProvider) ClearInjectedErrors()
```

**Parameters:**
  None

**Returns:**
  None

### ContentType

ContentType returns the content type of the body.

```go
func (*MockBodyProvider) ContentType() string
```

**Parameters:**
  None

**Returns:**
- string

### InjectBodyError

InjectBodyError injects an error for the next Body call.

```go
func (*MockBodyProvider) InjectBodyError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### Reset

Reset resets the mock to initial state.

```go
func (*MockBodyProvider) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### SetBody

SetBody sets the body content.

```go
func (*MockBodyProvider) SetBody(body []byte)
```

**Parameters:**
- `body` ([]byte)

**Returns:**
  None

### SetContentType

SetContentType sets the content type.

```go
func (*MockBodyProvider) SetContentType(contentType string)
```

**Parameters:**
- `contentType` (string)

**Returns:**
  None

### MockBodyProviderOptions
MockBodyProviderOptions configures mock body provider behavior.

#### Example Usage

```go
// Create a new MockBodyProviderOptions
mockbodyprovideroptions := MockBodyProviderOptions{
    ContentType: "example",
    Body: [],
}
```

#### Type Definition

```go
type MockBodyProviderOptions struct {
    ContentType string
    Body []byte
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ContentType | `string` |  |
| Body | `[]byte` |  |

### MockCache
MockCache is a mock implementation of neuron.Cache for testing.

#### Example Usage

```go
// Create a new MockCache
mockcache := MockCache{

}
```

#### Type Definition

```go
type MockCache struct {
}
```

### Constructor Functions

### NewMockCache

NewMockCache creates a new mock cache.

```go
func NewMockCache(opts *MockCacheOptions) *MockCache
```

**Parameters:**
- `opts` (*MockCacheOptions)

**Returns:**
- *MockCache

## Methods

### Clear

Clear removes all entries from the cache.

```go
func (*MockCache) Clear()
```

**Parameters:**
  None

**Returns:**
  None

### ClearRecorded

ClearRecorded clears all recorded operations and statistics.

```go
func (*MockAuthProvider) ClearRecorded()
```

**Parameters:**
  None

**Returns:**
  None

### Clears

Clears returns the number of cache clears.

```go
func (*MockCache) Clears() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Contains

Contains checks if a key exists in the cache.

```go
func (*MockCache) Contains(key string) bool
```

**Parameters:**
- `key` (string)

**Returns:**
- bool

### Delete

Delete removes an entry from the cache.

```go
func (*MockCache) Delete(key string)
```

**Parameters:**
- `key` (string)

**Returns:**
  None

### Deletes

Deletes returns the number of cache deletes.

```go
func (*MockCache) Deletes() int64
```

**Parameters:**
  None

**Returns:**
- int64

### EnableRecording

EnableRecording enables or disables operation recording.

```go
func (*MockCache) EnableRecording(enabled bool)
```

**Parameters:**
- `enabled` (bool)

**Returns:**
  None

### Get

Get retrieves an entry from the cache.

```go
func (*MockCache) Get(key string) (*neuron.CacheEntry, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- *neuron.CacheEntry
- bool

### GetEntry

GetEntry returns a copy of an entry without updating hit counts.

```go
func (*MockCache) GetEntry(key string) (*neuron.CacheEntry, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- *neuron.CacheEntry
- bool

### HitRate

HitRate returns the ratio of hits to total requests.

```go
func (*MockCache) HitRate() float64
```

**Parameters:**
  None

**Returns:**
- float64

### Hits

Hits returns the number of cache hits.

```go
func (*MockCache) Hits() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Keys

Keys returns a copy of all keys in the cache.

```go
func (*MockCache) Keys() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Misses

Misses returns the number of cache misses.

```go
func (*MockCache) Misses() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Operations

Operations returns a copy of all recorded operations.

```go
func (*MockCache) Operations() []CacheOperation
```

**Parameters:**
  None

**Returns:**
- []CacheOperation

### Reset

Reset clears all cache state, statistics, and recorded operations.

```go
func (*MockAuthProvider) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### Set

Set stores an entry in the cache.

```go
func (*MockCache) Set(key string, entry neuron.CacheEntry)
```

**Parameters:**
- `key` (string)
- `entry` (neuron.CacheEntry)

**Returns:**
  None

### Sets

Sets returns the number of cache sets.

```go
func (*MockCache) Sets() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Size

Size returns the number of entries in the cache.

```go
func (*MockCache) Size() int
```

**Parameters:**
  None

**Returns:**
- int

### MockCacheOptions
MockCacheOptions configures mock cache behavior.

#### Example Usage

```go
// Create a new MockCacheOptions
mockcacheoptions := MockCacheOptions{
    RecordOperations: true,
    InitialData: map[],
}
```

#### Type Definition

```go
type MockCacheOptions struct {
    RecordOperations bool
    InitialData map[string]neuron.CacheEntry
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| RecordOperations | `bool` |  |
| InitialData | `map[string]neuron.CacheEntry` |  |

### MockRateLimitHandler
MockRateLimitHandler is a mock implementation of neuron.RateLimitHandler. It embeds MockRateLimiter and adds UpdateFromHeaders support.

#### Example Usage

```go
// Create a new MockRateLimitHandler
mockratelimithandler := MockRateLimitHandler{

}
```

#### Type Definition

```go
type MockRateLimitHandler struct {
    *MockRateLimiter
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| **MockRateLimiter | `*MockRateLimiter` |  |

### Constructor Functions

### NewMockRateLimitHandler

NewMockRateLimitHandler creates a new mock rate limit handler.

```go
func NewMockRateLimitHandler(opts *MockRateLimiterOptions) *MockRateLimitHandler
```

**Parameters:**
- `opts` (*MockRateLimiterOptions)

**Returns:**
- *MockRateLimitHandler

## Methods

### ClearRecordedUpdates

ClearRecordedUpdates clears all recorded updates.

```go
func (*MockRateLimitHandler) ClearRecordedUpdates()
```

**Parameters:**
  None

**Returns:**
  None

### LastUpdate

LastUpdate returns the most recent update, or nil if none.

```go
func (*MockRateLimitHandler) LastUpdate() *RateLimitUpdate
```

**Parameters:**
  None

**Returns:**
- *RateLimitUpdate

### Reset

Reset resets all state including updates.

```go
func (*MockAuthProvider) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### UpdateFromHeaders

UpdateFromHeaders records a call to UpdateFromHeaders.

```go
func (*MockRateLimitHandler) UpdateFromHeaders(method, endpoint string, info *neuron.RateLimitInfo) error
```

**Parameters:**
- `method` (string)
- `endpoint` (string)
- `info` (*neuron.RateLimitInfo)

**Returns:**
- error

### Updates

Updates returns a copy of all recorded UpdateFromHeaders calls.

```go
func (*MockRateLimitHandler) Updates() []RateLimitUpdate
```

**Parameters:**
  None

**Returns:**
- []RateLimitUpdate

### UpdatesForEndpoint

UpdatesForEndpoint returns updates for a specific endpoint.

```go
func (*MockRateLimitHandler) UpdatesForEndpoint(method, endpoint string) []RateLimitUpdate
```

**Parameters:**
- `method` (string)
- `endpoint` (string)

**Returns:**
- []RateLimitUpdate

### WasExhausted

WasExhausted returns true if any recorded update indicated exhaustion.

```go
func (*MockRateLimitHandler) WasExhausted() bool
```

**Parameters:**
  None

**Returns:**
- bool

### MockRateLimiter
MockRateLimiter is a mock implementation of neuron.RateLimiter for testing.

#### Example Usage

```go
// Create a new MockRateLimiter
mockratelimiter := MockRateLimiter{

}
```

#### Type Definition

```go
type MockRateLimiter struct {
}
```

### Constructor Functions

### NewMockRateLimiter

NewMockRateLimiter creates a new mock rate limiter.

```go
func NewMockRateLimiter(opts *MockRateLimiterOptions) *MockRateLimiter
```

**Parameters:**
- `opts` (*MockRateLimiterOptions)

**Returns:**
- *MockRateLimiter

## Methods

### Allow

Allow checks if a request is allowed without blocking.

```go
func (*MockRateLimiter) Allow(ctx context.Context, method, endpoint string) bool
```

**Parameters:**
- `ctx` (context.Context)
- `method` (string)
- `endpoint` (string)

**Returns:**
- bool

### AllowCalls

AllowCalls returns a copy of all recorded Allow calls.

```go
func (*MockRateLimiter) AllowCalls() []AllowCall
```

**Parameters:**
  None

**Returns:**
- []AllowCall

### ClearInjectedErrors

ClearInjectedErrors clears all injected errors.

```go
func (*MockBodyProvider) ClearInjectedErrors()
```

**Parameters:**
  None

**Returns:**
  None

### ClearRecorded

ClearRecorded clears all recorded calls.

```go
func (*MockRateLimiter) ClearRecorded()
```

**Parameters:**
  None

**Returns:**
  None

### InjectWaitError

InjectWaitError injects an error for the next Wait call (one-shot by default).

```go
func (*MockRateLimiter) InjectWaitError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### Reset

Reset resets all state including configuration and recorded calls.

```go
func (*MockRateLimitHandler) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### SetAllow

SetAllow sets the global allow state.

```go
func (*MockRateLimiter) SetAllow(allow bool)
```

**Parameters:**
- `allow` (bool)

**Returns:**
  None

### SetAllowForEndpoint

SetAllowForEndpoint sets allow state for a specific endpoint.

```go
func (*MockRateLimiter) SetAllowForEndpoint(method, endpoint string, allow bool)
```

**Parameters:**
- `method` (string)
- `endpoint` (string)
- `allow` (bool)

**Returns:**
  None

### SetWaitDuration

SetWaitDuration sets the duration that Wait should sleep.

```go
func (*MockRateLimiter) SetWaitDuration(duration time.Duration)
```

**Parameters:**
- `duration` (time.Duration)

**Returns:**
  None

### String

String returns a human-readable description of the rate limit state.

```go
func (*MockRateLimiter) String() string
```

**Parameters:**
  None

**Returns:**
- string

### Wait

Wait blocks until the request is allowed or context is cancelled.

```go
func (*MockRateLimiter) Wait(ctx context.Context, method, endpoint string) error
```

**Parameters:**
- `ctx` (context.Context)
- `method` (string)
- `endpoint` (string)

**Returns:**
- error

### WaitCalls

WaitCalls returns a copy of all recorded Wait calls.

```go
func (*MockRateLimiter) WaitCalls() []WaitCall
```

**Parameters:**
  None

**Returns:**
- []WaitCall

### MockRateLimiterOptions
MockRateLimiterOptions configures mock rate limiter behavior.

#### Example Usage

```go
// Create a new MockRateLimiterOptions
mockratelimiteroptions := MockRateLimiterOptions{
    RecordCalls: true,
    DefaultAllow: true,
}
```

#### Type Definition

```go
type MockRateLimiterOptions struct {
    RecordCalls bool
    DefaultAllow bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| RecordCalls | `bool` |  |
| DefaultAllow | `bool` |  |

### MockRequestIDGenerator
MockRequestIDGenerator is a mock implementation of neuron.RequestIDGenerator for testing.

#### Example Usage

```go
// Create a new MockRequestIDGenerator
mockrequestidgenerator := MockRequestIDGenerator{

}
```

#### Type Definition

```go
type MockRequestIDGenerator struct {
}
```

### Constructor Functions

### NewMockRequestIDGenerator

NewMockRequestIDGenerator creates a new mock request ID generator.

```go
func NewMockRequestIDGenerator(opts *MockRequestIDGeneratorOptions) *MockRequestIDGenerator
```

**Parameters:**
- `opts` (*MockRequestIDGeneratorOptions)

**Returns:**
- *MockRequestIDGenerator

## Methods

### CallCount

CallCount returns the number of times Generate was called.

```go
func (*MockBodyProvider) CallCount() int64
```

**Parameters:**
  None

**Returns:**
- int64

### CurrentIndex

CurrentIndex returns the current index in the ID sequence.

```go
func (*MockRequestIDGenerator) CurrentIndex() int
```

**Parameters:**
  None

**Returns:**
- int

### Generate

Generate returns the next request ID in the sequence.

```go
func (*MockRequestIDGenerator) Generate() string
```

**Parameters:**
  None

**Returns:**
- string

### Reset

Reset resets the ID generator to the first ID.

```go
func (*MockAuthProvider) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### SetIDs

SetIDs sets the ID sequence.

```go
func (*MockRequestIDGenerator) SetIDs(ids []string)
```

**Parameters:**
- `ids` ([]string)

**Returns:**
  None

### MockRequestIDGeneratorOptions
MockRequestIDGeneratorOptions configures mock request ID generator behavior.

#### Example Usage

```go
// Create a new MockRequestIDGeneratorOptions
mockrequestidgeneratoroptions := MockRequestIDGeneratorOptions{
    IDs: [],
}
```

#### Type Definition

```go
type MockRequestIDGeneratorOptions struct {
    IDs []string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| IDs | `[]string` |  |

### MockRoundTripper
MockRoundTripper is a mock implementation of http.RoundTripper for testing. It allows queuing responses, recording requests, and injecting errors.

#### Example Usage

```go
// Create a new MockRoundTripper
mockroundtripper := MockRoundTripper{

}
```

#### Type Definition

```go
type MockRoundTripper struct {
}
```

### Constructor Functions

### NewMockRoundTripper

NewMockRoundTripper creates a new mock round tripper.

```go
func NewMockRoundTripper(opts *MockRoundTripperOptions) *MockRoundTripper
```

**Parameters:**
- `opts` (*MockRoundTripperOptions)

**Returns:**
- *MockRoundTripper

## Methods

### AddMatcher

AddMatcher adds a response matcher for pattern-based matching.

```go
func (*MockRoundTripper) AddMatcher(matcher ResponseMatcher)
```

**Parameters:**
- `matcher` (ResponseMatcher)

**Returns:**
  None

### ClearInjectedErrors

ClearInjectedErrors clears all injected errors.

```go
func (*MockBodyProvider) ClearInjectedErrors()
```

**Parameters:**
  None

**Returns:**
  None

### ClearRecorded

ClearRecorded clears all recorded requests and responses.

```go
func (*MockRoundTripper) ClearRecorded()
```

**Parameters:**
  None

**Returns:**
  None

### EnableRecording

EnableRecording enables request recording.

```go
func (*MockRoundTripper) EnableRecording(enabled bool)
```

**Parameters:**
- `enabled` (bool)

**Returns:**
  None

### InjectError

InjectError injects an error for the next request (one-shot by default).

```go
func (*MockRoundTripper) InjectError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### InjectErrorSequence

InjectErrorSequence injects a sequence of errors for successive requests.

```go
func (*MockRoundTripper) InjectErrorSequence(errs []error)
```

**Parameters:**
- `errs` ([]error)

**Returns:**
  None

### QueueResponse

QueueResponse adds a response to the response queue.

```go
func (*MockRoundTripper) QueueResponse(config ResponseConfig)
```

**Parameters:**
- `config` (ResponseConfig)

**Returns:**
  None

### RequestCount

RequestCount returns the total number of requests made.

```go
func (*MockRoundTripper) RequestCount() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Requests

Requests returns a copy of all recorded requests.

```go
func (*MockRoundTripper) Requests() []RecordedRequest
```

**Parameters:**
  None

**Returns:**
- []RecordedRequest

### Reset

Reset resets all state including requests, responses, and injected errors.

```go
func (*ErrorSequence) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### RoundTrip

RoundTrip implements http.RoundTripper by returning queued responses or executing matchers.

```go
func (*MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error)
```

**Parameters:**
- `req` (*http.Request)

**Returns:**
- *http.Response
- error

### SetLatency

SetLatency sets the simulated network latency.

```go
func (*MockRoundTripper) SetLatency(latency time.Duration)
```

**Parameters:**
- `latency` (time.Duration)

**Returns:**
  None

### MockRoundTripperOptions
MockRoundTripperOptions configures MockRoundTripper behavior.

#### Example Usage

```go
// Create a new MockRoundTripperOptions
mockroundtripperoptions := MockRoundTripperOptions{
    RecordRequests: true,
    Latency: /* value */,
    BufferSize: 42,
}
```

#### Type Definition

```go
type MockRoundTripperOptions struct {
    RecordRequests bool
    Latency time.Duration
    BufferSize int
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| RecordRequests | `bool` |  |
| Latency | `time.Duration` |  |
| BufferSize | `int` |  |

### MockValidator
MockValidator is a mock implementation of neuron.Validator for testing.

#### Example Usage

```go
// Create a new MockValidator
mockvalidator := MockValidator{

}
```

#### Type Definition

```go
type MockValidator struct {
}
```

### Constructor Functions

### NewMockValidator

NewMockValidator creates a new mock validator.

```go
func NewMockValidator(opts *MockValidatorOptions) *MockValidator
```

**Parameters:**
- `opts` (*MockValidatorOptions)

**Returns:**
- *MockValidator

## Methods

### ClearInjectedErrors

ClearInjectedErrors clears all injected errors.

```go
func (*MockBodyProvider) ClearInjectedErrors()
```

**Parameters:**
  None

**Returns:**
  None

### ClearRecorded

ClearRecorded clears all recorded calls.

```go
func (*MockCache) ClearRecorded()
```

**Parameters:**
  None

**Returns:**
  None

### InjectValidationError

InjectValidationError injects an error for the next Validate call.

```go
func (*MockValidator) InjectValidationError(err error)
```

**Parameters:**
- `err` (error)

**Returns:**
  None

### Reset

Reset resets all state.

```go
func (*MockCache) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### Validate

Validate validates the provided data.

```go
func (*MockValidator) Validate(data []byte, contentType string) error
```

**Parameters:**
- `data` ([]byte)
- `contentType` (string)

**Returns:**
- error

### ValidateCalls

ValidateCalls returns a copy of all recorded Validate calls.

```go
func (*MockValidator) ValidateCalls() []ValidateCall
```

**Parameters:**
  None

**Returns:**
- []ValidateCall

### MockValidatorOptions
MockValidatorOptions configures mock validator behavior.

#### Example Usage

```go
// Create a new MockValidatorOptions
mockvalidatoroptions := MockValidatorOptions{
    RecordCalls: true,
}
```

#### Type Definition

```go
type MockValidatorOptions struct {
    RecordCalls bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| RecordCalls | `bool` |  |

### RateLimitUpdate
RateLimitUpdate records a call to UpdateFromHeaders.

#### Example Usage

```go
// Create a new RateLimitUpdate
ratelimitupdate := RateLimitUpdate{
    Method: "example",
    Endpoint: "example",
    Info: &/* value */{},
    Time: /* value */,
}
```

#### Type Definition

```go
type RateLimitUpdate struct {
    Method string
    Endpoint string
    Info *neuron.RateLimitInfo
    Time time.Time
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Method | `string` |  |
| Endpoint | `string` |  |
| Info | `*neuron.RateLimitInfo` |  |
| Time | `time.Time` |  |

### RecordedRequest
RecordedRequest represents a recorded HTTP request.

#### Example Usage

```go
// Create a new RecordedRequest
recordedrequest := RecordedRequest{
    Method: "example",
    URL: &/* value */{},
    Headers: /* value */,
    Body: [],
    Time: /* value */,
}
```

#### Type Definition

```go
type RecordedRequest struct {
    Method string
    URL *url.URL
    Headers http.Header
    Body []byte
    Time time.Time
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Method | `string` |  |
| URL | `*url.URL` |  |
| Headers | `http.Header` |  |
| Body | `[]byte` |  |
| Time | `time.Time` |  |

### Constructor Functions

### AssertRequestsForMethod

AssertRequestsForMethod returns requests matching the given HTTP method.

```go
func AssertRequestsForMethod(t *testing.T, rt *MockRoundTripper, method string) []RecordedRequest
```

**Parameters:**
- `t` (*testing.T)
- `rt` (*MockRoundTripper)
- `method` (string)

**Returns:**
- []RecordedRequest

### RequestCapture
RequestCapture captures the essential details of an HTTP request for analysis.

#### Example Usage

```go
// Create a new RequestCapture
requestcapture := RequestCapture{
    Method: "example",
    URL: "example",
    Headers: map[],
    Body: [],
    Time: 42,
}
```

#### Type Definition

```go
type RequestCapture struct {
    Method string
    URL string
    Headers map[string][]string
    Body []byte
    Time int64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Method | `string` |  |
| URL | `string` |  |
| Headers | `map[string][]string` |  |
| Body | `[]byte` |  |
| Time | `int64` | nanoseconds since epoch |

### ResponseCapture
ResponseCapture captures the essential details of an HTTP response for analysis.

#### Example Usage

```go
// Create a new ResponseCapture
responsecapture := ResponseCapture{
    StatusCode: 42,
    Headers: map[],
    Body: [],
    Time: 42,
}
```

#### Type Definition

```go
type ResponseCapture struct {
    StatusCode int
    Headers map[string][]string
    Body []byte
    Time int64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| StatusCode | `int` |  |
| Headers | `map[string][]string` |  |
| Body | `[]byte` |  |
| Time | `int64` | nanoseconds since epoch |

### ResponseConfig
ResponseConfig represents a configured response.

#### Example Usage

```go
// Create a new ResponseConfig
responseconfig := ResponseConfig{
    StatusCode: 42,
    Body: [],
    Headers: /* value */,
    Err: error{},
    Delay: /* value */,
}
```

#### Type Definition

```go
type ResponseConfig struct {
    StatusCode int
    Body []byte
    Headers http.Header
    Err error
    Delay time.Duration
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| StatusCode | `int` |  |
| Body | `[]byte` |  |
| Headers | `http.Header` |  |
| Err | `error` |  |
| Delay | `time.Duration` |  |

### ResponseMatcher
ResponseMatcher is a function that matches a request and returns a response.

#### Example Usage

```go
// Example usage of ResponseMatcher
var value ResponseMatcher
// Initialize with appropriate value
```

#### Type Definition

```go
type ResponseMatcher func(req *http.Request) (*http.Response, bool)
```

### ValidateCall
ValidateCall records a call to Validate.

#### Example Usage

```go
// Create a new ValidateCall
validatecall := ValidateCall{
    Data: [],
    ContentType: "example",
    Result: error{},
}
```

#### Type Definition

```go
type ValidateCall struct {
    Data []byte
    ContentType string
    Result error
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Data | `[]byte` |  |
| ContentType | `string` |  |
| Result | `error` |  |

### WaitCall
WaitCall records a call to Wait.

#### Example Usage

```go
// Create a new WaitCall
waitcall := WaitCall{
    Method: "example",
    Endpoint: "example",
    Time: /* value */,
    Waited: /* value */,
    Err: error{},
}
```

#### Type Definition

```go
type WaitCall struct {
    Method string
    Endpoint string
    Time time.Time
    Waited time.Duration
    Err error
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Method | `string` |  |
| Endpoint | `string` |  |
| Time | `time.Time` |  |
| Waited | `time.Duration` |  |
| Err | `error` |  |

## Functions

### AssertAllowCalled
AssertAllowCalled fails the test if Allow was not called for the given method and endpoint.

```go
func AssertAllowCalled(t *testing.T, rl *MockRateLimiter, method, endpoint string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rl` | `*MockRateLimiter` | |
| `method` | `string` | |
| `endpoint` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertAllowCalled
result := AssertAllowCalled(/* parameters */)
```

### AssertAllowCount
AssertAllowCount fails the test if the Allow call count doesn't match expected.

```go
func AssertAllowCount(t *testing.T, rl *MockRateLimiter, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rl` | `*MockRateLimiter` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertAllowCount
result := AssertAllowCount(/* parameters */)
```

### AssertAllowNotCalled
AssertAllowNotCalled fails the test if Allow was called for the given method and endpoint.

```go
func AssertAllowNotCalled(t *testing.T, rl *MockRateLimiter, method, endpoint string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rl` | `*MockRateLimiter` | |
| `method` | `string` | |
| `endpoint` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertAllowNotCalled
result := AssertAllowNotCalled(/* parameters */)
```

### AssertCacheContains
AssertCacheContains fails the test if the key doesn't exist in the cache.

```go
func AssertCacheContains(t *testing.T, cache *MockCache, key string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `cache` | `*MockCache` | |
| `key` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertCacheContains
result := AssertCacheContains(/* parameters */)
```

### AssertCacheEmpty
AssertCacheEmpty fails the test if the cache is not empty.

```go
func AssertCacheEmpty(t *testing.T, cache *MockCache)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `cache` | `*MockCache` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertCacheEmpty
result := AssertCacheEmpty(/* parameters */)
```

### AssertCacheHitRate
AssertCacheHitRate fails the test if the hit rate doesn't match the expected value (within tolerance).

```go
func AssertCacheHitRate(t *testing.T, cache *MockCache, expected float64, tolerance float64)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `cache` | `*MockCache` | |
| `expected` | `float64` | |
| `tolerance` | `float64` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertCacheHitRate
result := AssertCacheHitRate(/* parameters */)
```

### AssertCacheHits
AssertCacheHits fails the test if the hit count doesn't match expected.

```go
func AssertCacheHits(t *testing.T, cache *MockCache, expected int64)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `cache` | `*MockCache` | |
| `expected` | `int64` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertCacheHits
result := AssertCacheHits(/* parameters */)
```

### AssertCacheMisses
AssertCacheMisses fails the test if the miss count doesn't match expected.

```go
func AssertCacheMisses(t *testing.T, cache *MockCache, expected int64)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `cache` | `*MockCache` | |
| `expected` | `int64` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertCacheMisses
result := AssertCacheMisses(/* parameters */)
```

### AssertCacheNotContains
AssertCacheNotContains fails the test if the key exists in the cache.

```go
func AssertCacheNotContains(t *testing.T, cache *MockCache, key string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `cache` | `*MockCache` | |
| `key` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertCacheNotContains
result := AssertCacheNotContains(/* parameters */)
```

### AssertCacheSize
AssertCacheSize fails the test if the cache size doesn't match expected.

```go
func AssertCacheSize(t *testing.T, cache *MockCache, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `cache` | `*MockCache` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertCacheSize
result := AssertCacheSize(/* parameters */)
```

### AssertHeaderGenerated
AssertHeaderGenerated fails the test if GetAuthHeader was not called for the token.

```go
func AssertHeaderGenerated(t *testing.T, ap *MockAuthProvider, token string, expectedHeader string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `ap` | `*MockAuthProvider` | |
| `token` | `string` | |
| `expectedHeader` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertHeaderGenerated
result := AssertHeaderGenerated(/* parameters */)
```

### AssertLastRequestMethod
AssertLastRequestMethod fails the test if the last request's method doesn't match.

```go
func AssertLastRequestMethod(t *testing.T, rt *MockRoundTripper, method string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `method` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertLastRequestMethod
result := AssertLastRequestMethod(/* parameters */)
```

### AssertLastRequestURL
AssertLastRequestURL fails the test if the last request's URL doesn't match.

```go
func AssertLastRequestURL(t *testing.T, rt *MockRoundTripper, expectedURL string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `expectedURL` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertLastRequestURL
result := AssertLastRequestURL(/* parameters */)
```

### AssertNoRequests
AssertNoRequests fails the test if any requests were made.

```go
func AssertNoRequests(t *testing.T, rt *MockRoundTripper)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertNoRequests
result := AssertNoRequests(/* parameters */)
```

### AssertQueryParam
AssertQueryParam fails the test if a query parameter doesn't match the expected value.

```go
func AssertQueryParam(t *testing.T, rt *MockRoundTripper, index int, key, expectedValue string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `key` | `string` | |
| `expectedValue` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertQueryParam
result := AssertQueryParam(/* parameters */)
```

### AssertRateLimitExhausted
AssertRateLimitExhausted fails the test if a rate limit exhaustion was not recorded.

```go
func AssertRateLimitExhausted(t *testing.T, rlh *MockRateLimitHandler)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rlh` | `*MockRateLimitHandler` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRateLimitExhausted
result := AssertRateLimitExhausted(/* parameters */)
```

### AssertRateLimitUpdateCount
AssertRateLimitUpdateCount fails the test if the update count doesn't match expected.

```go
func AssertRateLimitUpdateCount(t *testing.T, rlh *MockRateLimitHandler, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rlh` | `*MockRateLimitHandler` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRateLimitUpdateCount
result := AssertRateLimitUpdateCount(/* parameters */)
```

### AssertRateLimitUpdated
AssertRateLimitUpdated fails the test if UpdateFromHeaders was not called for the endpoint.

```go
func AssertRateLimitUpdated(t *testing.T, rlh *MockRateLimitHandler, method, endpoint string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rlh` | `*MockRateLimitHandler` | |
| `method` | `string` | |
| `endpoint` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRateLimitUpdated
result := AssertRateLimitUpdated(/* parameters */)
```

### AssertRequestBody
AssertRequestBody fails the test if the body doesn't match the expected bytes.

```go
func AssertRequestBody(t *testing.T, rt *MockRoundTripper, index int, expectedBody []byte)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `expectedBody` | `[]byte` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestBody
result := AssertRequestBody(/* parameters */)
```

### AssertRequestBodyJSON
AssertRequestBodyJSON fails the test if the body can't be unmarshaled as JSON matching the expected value.

```go
func AssertRequestBodyJSON(t *testing.T, rt *MockRoundTripper, index int, expected any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `expected` | `any` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestBodyJSON
result := AssertRequestBodyJSON(/* parameters */)
```

### AssertRequestBodyString
AssertRequestBodyString fails the test if the body doesn't match the expected string.

```go
func AssertRequestBodyString(t *testing.T, rt *MockRoundTripper, index int, expectedBody string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `expectedBody` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestBodyString
result := AssertRequestBodyString(/* parameters */)
```

### AssertRequestCount
AssertRequestCount fails the test if the request count doesn't match expected.

```go
func AssertRequestCount(t *testing.T, rt *MockRoundTripper, expected int64)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `expected` | `int64` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestCount
result := AssertRequestCount(/* parameters */)
```

### AssertRequestCountBetween
AssertRequestCountBetween fails the test if the request count is not within the range.

```go
func AssertRequestCountBetween(t *testing.T, rt *MockRoundTripper, min, max int64)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `min` | `int64` | |
| `max` | `int64` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestCountBetween
result := AssertRequestCountBetween(/* parameters */)
```

### AssertRequestHeader
AssertRequestHeader fails the test if a header value doesn't match.

```go
func AssertRequestHeader(t *testing.T, rt *MockRoundTripper, index int, key, expectedValue string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `key` | `string` | |
| `expectedValue` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestHeader
result := AssertRequestHeader(/* parameters */)
```

### AssertRequestMethod
AssertRequestMethod fails the test if the method of the i-th request doesn't match.

```go
func AssertRequestMethod(t *testing.T, rt *MockRoundTripper, index int, method string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `method` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestMethod
result := AssertRequestMethod(/* parameters */)
```

### AssertRequestURL
AssertRequestURL fails the test if the URL of the i-th request doesn't match.

```go
func AssertRequestURL(t *testing.T, rt *MockRoundTripper, index int, expectedURL string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `expectedURL` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestURL
result := AssertRequestURL(/* parameters */)
```

### AssertRequestURLPath
AssertRequestURLPath fails the test if the path of the i-th request doesn't match.

```go
func AssertRequestURLPath(t *testing.T, rt *MockRoundTripper, index int, expectedPath string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `index` | `int` | |
| `expectedPath` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestURLPath
result := AssertRequestURLPath(/* parameters */)
```

### AssertRequestsContainPath
AssertRequestsContainPath fails the test if none of the requests match the given path.

```go
func AssertRequestsContainPath(t *testing.T, rt *MockRoundTripper, path string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `path` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertRequestsContainPath
result := AssertRequestsContainPath(/* parameters */)
```

### AssertTokenCalled
AssertTokenCalled fails the test if GetToken was not called exactly count times.

```go
func AssertTokenCalled(t *testing.T, ap *MockAuthProvider, count int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `ap` | `*MockAuthProvider` | |
| `count` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertTokenCalled
result := AssertTokenCalled(/* parameters */)
```

### AssertTokenRotated
AssertTokenRotated fails the test if the token was not rotated (multiple distinct tokens used).

```go
func AssertTokenRotated(t *testing.T, ap *MockAuthProvider)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `ap` | `*MockAuthProvider` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertTokenRotated
result := AssertTokenRotated(/* parameters */)
```

### AssertTokenValue
AssertTokenValue fails the test if the i-th token call didn't return the expected token.

```go
func AssertTokenValue(t *testing.T, ap *MockAuthProvider, index int, expected string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `ap` | `*MockAuthProvider` | |
| `index` | `int` | |
| `expected` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertTokenValue
result := AssertTokenValue(/* parameters */)
```

### AssertWaitCalled
AssertWaitCalled fails the test if Wait was not called for the given method and endpoint.

```go
func AssertWaitCalled(t *testing.T, rl *MockRateLimiter, method, endpoint string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rl` | `*MockRateLimiter` | |
| `method` | `string` | |
| `endpoint` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertWaitCalled
result := AssertWaitCalled(/* parameters */)
```

### AssertWaitCount
AssertWaitCount fails the test if the Wait call count doesn't match expected.

```go
func AssertWaitCount(t *testing.T, rl *MockRateLimiter, expected int)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rl` | `*MockRateLimiter` | |
| `expected` | `int` | |

**Returns:**
None

**Example:**

```go
// Example usage of AssertWaitCount
result := AssertWaitCount(/* parameters */)
```

### IsContextCanceledError
IsContextCanceledError checks if an error is a context cancellation error.

```go
func IsContextCanceledError(err error) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of IsContextCanceledError
result := IsContextCanceledError(/* parameters */)
```

### IsNetworkError
IsNetworkError checks if an error is a network-related error.

```go
func IsNetworkError(err error) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of IsNetworkError
result := IsNetworkError(/* parameters */)
```

### IsRateLimitError
IsRateLimitError checks if an error message contains "rate limit".

```go
func IsRateLimitError(err error) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of IsRateLimitError
result := IsRateLimitError(/* parameters */)
```

### IsTimeoutError
IsTimeoutError checks if an error is a timeout error.

```go
func IsTimeoutError(err error) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of IsTimeoutError
result := IsTimeoutError(/* parameters */)
```

### WaitForRequest
WaitForRequest polls until a request is received or timeout expires.

```go
func WaitForRequest(t *testing.T, rt *MockRoundTripper, timeout time.Duration) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `timeout` | `time.Duration` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of WaitForRequest
result := WaitForRequest(/* parameters */)
```

### WaitForRequestCount
WaitForRequestCount polls until the request count reaches expected or timeout expires.

```go
func WaitForRequestCount(t *testing.T, rt *MockRoundTripper, expected int64, timeout time.Duration) bool
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `t` | `*testing.T` | |
| `rt` | `*MockRoundTripper` | |
| `expected` | `int64` | |
| `timeout` | `time.Duration` | |

**Returns:**
| Type | Description |
|------|-------------|
| `bool` | |

**Example:**

```go
// Example usage of WaitForRequestCount
result := WaitForRequestCount(/* parameters */)
```

## External Links

- [Package Overview](../packages/mock.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/neuron/mock)
- [Source Code](https://github.com/kolosys/neuron/tree/dev/mock)
