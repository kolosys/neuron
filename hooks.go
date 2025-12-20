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
)

// HookChain manages a chain of hook functions
// This is an optional utility for advanced use cases; most users will use slices directly
type HookChain struct {
	requestHooks  []RequestHook
	responseHooks []ResponseHook
}

// NewHookChain creates a new hook chain
func NewHookChain() *HookChain {
	return &HookChain{
		requestHooks:  make([]RequestHook, 0),
		responseHooks: make([]ResponseHook, 0),
	}
}

// AddRequestHook adds a request hook to the chain
func (hc *HookChain) AddRequestHook(hook RequestHook) *HookChain {
	hc.requestHooks = append(hc.requestHooks, hook)
	return hc
}

// AddResponseHook adds a response hook to the chain
func (hc *HookChain) AddResponseHook(hook ResponseHook) *HookChain {
	hc.responseHooks = append(hc.responseHooks, hook)
	return hc
}

// ApplyRequestHooks applies all request hooks in order
func (hc *HookChain) ApplyRequestHooks(req *http.Request) error {
	for _, hook := range hc.requestHooks {
		if err := hook(req); err != nil {
			return err
		}
	}
	return nil
}

// ApplyResponseHooks applies all response hooks in order
func (hc *HookChain) ApplyResponseHooks(resp *http.Response) error {
	for _, hook := range hc.responseHooks {
		if err := hook(resp); err != nil {
			return err
		}
	}
	return nil
}

// AddRetry implements retry logic at the hook level
// Note: The client already has built-in retry support via ClientOptions.MaxRetries
func AddRetry(maxRetries int, retryCondition RetryCondition) RequestHook {
	return func(req *http.Request) error {
		ctx := context.WithValue(req.Context(), maxRetriesKey, maxRetries)
		ctx = context.WithValue(ctx, retryConditionKey, retryCondition)
		*req = *req.WithContext(ctx)
		return nil
	}
}

// RetryCondition determines if a request should be retried
type RetryCondition func(resp *http.Response, err error) bool

// AddResponseCache implements response caching
func AddResponseCache(cache Cache) ResponseHook {
	return func(resp *http.Response) error {
		if resp.Request.Method == "GET" && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			cacheKey := resp.Request.URL.String()
			cacheEntry := CacheEntry{
				Data:       body,
				Headers:    resp.Header,
				StatusCode: resp.StatusCode,
				Timestamp:  time.Now(),
			}

			cache.Set(cacheKey, cacheEntry)
			resp.Body = io.NopCloser(bytes.NewReader(body))
		}
		return nil
	}
}

// Cache interface for caching hooks
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

// InMemoryCache provides a simple in-memory cache implementation
type InMemoryCache struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
}

// NewInMemoryCache creates a new in-memory cache
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

// AddValidation validates request payloads
func AddValidation(validator Validator) RequestHook {
	return func(req *http.Request) error {
		if req.Body != nil && req.ContentLength > 0 {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return fmt.Errorf("failed to read request body: %w", err)
			}

			if err := validator.Validate(body, req.Header.Get("Content-Type")); err != nil {
				return fmt.Errorf("request validation failed: %w", err)
			}

			req.Body = io.NopCloser(bytes.NewReader(body))
		}
		return nil
	}
}

// Validator interface for request validation
type Validator interface {
	Validate(data []byte, contentType string) error
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
