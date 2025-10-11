package neuron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type contextKey string

const (
	maxRetriesKey     contextKey = "max_retries"
	retryConditionKey contextKey = "retry_condition"
	circuitBreakerKey contextKey = "circuit_breaker"
)

// MiddlewareChain manages a chain of middleware functions
type MiddlewareChain struct {
	requestMiddleware  []RequestMiddleware
	responseMiddleware []ResponseMiddleware
}

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		requestMiddleware:  make([]RequestMiddleware, 0),
		responseMiddleware: make([]ResponseMiddleware, 0),
	}
}

// AddRequestMiddleware adds a request middleware to the chain
func (mc *MiddlewareChain) AddRequestMiddleware(middleware RequestMiddleware) *MiddlewareChain {
	mc.requestMiddleware = append(mc.requestMiddleware, middleware)
	return mc
}

// AddResponseMiddleware adds a response middleware to the chain
func (mc *MiddlewareChain) AddResponseMiddleware(middleware ResponseMiddleware) *MiddlewareChain {
	mc.responseMiddleware = append(mc.responseMiddleware, middleware)
	return mc
}

// ApplyRequestMiddleware applies all request middleware in order
func (mc *MiddlewareChain) ApplyRequestMiddleware(req *http.Request) error {
	for _, middleware := range mc.requestMiddleware {
		if err := middleware(req); err != nil {
			return err
		}
	}
	return nil
}

// ApplyResponseMiddleware applies all response middleware in order
func (mc *MiddlewareChain) ApplyResponseMiddleware(resp *http.Response) error {
	for _, middleware := range mc.responseMiddleware {
		if err := middleware(resp); err != nil {
			return err
		}
	}
	return nil
}

// AuthenticationMiddleware adds authentication headers
func AuthenticationMiddleware(authProvider AuthProvider) RequestMiddleware {
	return func(req *http.Request) error {
		token, err := authProvider.GetToken(req.Context())
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if token != "" {
			authHeader := authProvider.GetAuthHeader(token)
			req.Header.Set("Authorization", authHeader)
		}

		return nil
	}
}

// RetryMiddleware implements retry logic at the middleware level
func RetryMiddleware(maxRetries int, retryCondition RetryCondition) RequestMiddleware {
	return func(req *http.Request) error {
		// Store retry configuration in context
		ctx := context.WithValue(req.Context(), maxRetriesKey, maxRetries)
		ctx = context.WithValue(ctx, retryConditionKey, retryCondition)
		*req = *req.WithContext(ctx)
		return nil
	}
}

// RateLimitMiddleware adds custom rate limiting headers
func RateLimitMiddleware(rateLimitInfo RateLimitInfoProvider) RequestMiddleware {
	return func(req *http.Request) error {
		info := rateLimitInfo.GetRateLimitInfo(req.URL.Path)
		if info != nil {
			req.Header.Set("X-RateLimit-Bucket", info.Bucket)
			req.Header.Set("X-RateLimit-Limit", fmt.Sprintf("%d", info.Limit))
			req.Header.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", info.Remaining))
		}
		return nil
	}
}

// CacheMiddleware implements response caching
func CacheMiddleware(cache Cache) ResponseMiddleware {
	return func(resp *http.Response) error {
		// Only cache successful GET requests
		if resp.Request.Method == "GET" && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			// Cache the response
			cacheKey := resp.Request.URL.String()
			cacheEntry := CacheEntry{
				Data:       body,
				Headers:    resp.Header,
				StatusCode: resp.StatusCode,
				Timestamp:  time.Now(),
			}

			cache.Set(cacheKey, cacheEntry)

			// Replace response body
			resp.Body = io.NopCloser(bytes.NewReader(body))
		}

		return nil
	}
}

// ValidationMiddleware validates request payloads
func ValidationMiddleware(validator Validator) RequestMiddleware {
	return func(req *http.Request) error {
		// Only validate requests with bodies
		if req.Body != nil && req.ContentLength > 0 {
			// Read body
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return fmt.Errorf("failed to read request body: %w", err)
			}

			// Validate
			if err := validator.Validate(body, req.Header.Get("Content-Type")); err != nil {
				return fmt.Errorf("request validation failed: %w", err)
			}

			// Replace body
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		return nil
	}
}

// CircuitBreakerMiddleware implements circuit breaker pattern
func CircuitBreakerMiddleware(circuitBreaker CircuitBreaker) RequestMiddleware {
	return func(req *http.Request) error {
		if !circuitBreaker.AllowRequest() {
			return fmt.Errorf("circuit breaker is open")
		}

		// Store circuit breaker in context for response handling
		ctx := context.WithValue(req.Context(), circuitBreakerKey, circuitBreaker)
		*req = *req.WithContext(ctx)

		return nil
	}
}

// CircuitBreakerResponseMiddleware handles circuit breaker state updates
func CircuitBreakerResponseMiddleware() ResponseMiddleware {
	return func(resp *http.Response) error {
		circuitBreaker, ok := resp.Request.Context().Value(circuitBreakerKey).(CircuitBreaker)
		if !ok {
			return nil
		}

		// Update circuit breaker based on response
		if resp.StatusCode >= 500 {
			circuitBreaker.RecordFailure()
		} else {
			circuitBreaker.RecordSuccess()
		}

		return nil
	}
}

// Interfaces for extensibility

// Logger interface for logging middleware
type Logger interface {
	LogRequest(entry LogEntry)
	LogResponse(entry LogEntry)
}

// LogEntry represents a log entry
type LogEntry struct {
	Method     string
	URL        string
	StatusCode int
	Headers    http.Header
	Duration   time.Duration
	Timestamp  time.Time
	Body       []byte
	Error      error
}

// AuthProvider interface for authentication
type AuthProvider interface {
	GetToken(ctx context.Context) (string, error)
	GetAuthHeader(token string) string
}

// RetryCondition determines if a request should be retried
type RetryCondition func(resp *http.Response, err error) bool

// RateLimitInfoProvider provides rate limit information
type RateLimitInfoProvider interface {
	GetRateLimitInfo(path string) *RateLimitInfo
}

// Cache interface for caching middleware
type Cache interface {
	Get(key string) (*CacheEntry, bool)
	Set(key string, entry CacheEntry)
	Delete(key string)
	Clear()
}

// CacheEntry represents a cached response
type CacheEntry struct {
	Data       []byte
	Headers    http.Header
	StatusCode int
	Timestamp  time.Time
	TTL        time.Duration
}

// Validator interface for request validation
type Validator interface {
	Validate(data []byte, contentType string) error
}

// Default implementations

// SimpleLogger provides a basic logger implementation
type SimpleLogger struct{}

func (l *SimpleLogger) LogRequest(entry LogEntry) {
	fmt.Printf("[REQUEST] %s %s\n", entry.Method, entry.URL)
}

func (l *SimpleLogger) LogResponse(entry LogEntry) {
	fmt.Printf("[RESPONSE] %s %s %d (%v)\n", entry.Method, entry.URL, entry.StatusCode, entry.Duration)
}

// StaticAuthProvider provides static token authentication
type StaticAuthProvider struct {
	Token  string
	Prefix string
}

func (a *StaticAuthProvider) GetToken(ctx context.Context) (string, error) {
	return a.Token, nil
}

func (a *StaticAuthProvider) GetAuthHeader(token string) string {
	if a.Prefix != "" {
		return a.Prefix + " " + token
	}
	return "Bearer " + token
}

// InMemoryCache provides a simple in-memory cache
type InMemoryCache struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]CacheEntry),
	}
}

func (c *InMemoryCache) Get(key string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Check TTL
	if entry.TTL > 0 && time.Since(entry.Timestamp) > entry.TTL {
		delete(c.data, key)
		return nil, false
	}

	return &entry, true
}

func (c *InMemoryCache) Set(key string, entry CacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = entry
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]CacheEntry)
}

// JSONValidator provides JSON schema validation
type JSONValidator struct{}

func (v *JSONValidator) Validate(data []byte, contentType string) error {
	if contentType == "application/json" {
		var obj any
		return json.Unmarshal(data, &obj)
	}
	return nil
}
