package neuron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
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

// AddResponseTimeout creates a response timeout middleware that checks if a request took too long.
func AddResponseTimeout(timeout time.Duration) ResponseHook {
	return func(resp *http.Response) error {
		start, ok := resp.Request.Context().Value(requestStartKey).(time.Time)
		if !ok {
			return nil
		}

		if time.Since(start) > timeout {
			return ClientError{
				Type:    ErrorTypeTimeout,
				Message: "response took longer than specified timeout",
				Route:   resp.Request.URL.Path,
			}
		}

		return nil
	}
}

// AddRateLimitHandler creates a response hook that parses rate limit headers
// and updates the provided RateLimitUpdater with the current rate limit state.
func AddRateLimitHandler(updater RateLimitUpdater) ResponseHook {
	return func(resp *http.Response) error {
		if updater == nil {
			return nil
		}

		info := ParseRateLimitHeaders(resp.Header)
		if info != nil {
			return updater.UpdateFromHeaders(
				resp.Request.Method,
				resp.Request.URL.Path,
				info,
			)
		}
		return nil
	}
}

// AddRateLimitRetry creates a response hook that returns an error on 429 responses.
// The error can be used to trigger automatic retry logic at the client level.
// For automatic 429 handling, use ClientOptions.AutoHandleRateLimit instead.
func AddRateLimitRetry() ResponseHook {
	return func(resp *http.Response) error {
		if resp.StatusCode == 429 {
			info := ParseRateLimitHeaders(resp.Header)
			retryAfter := time.Second // Default

			if info != nil && info.RetryAfter > 0 {
				retryAfter = info.RetryAfter
			}

			return ClientError{
				Type:       ErrorTypeRateLimit,
				StatusCode: 429,
				Message:    fmt.Sprintf("rate limited, retry after %v", retryAfter),
				Route:      resp.Request.URL.Path,
			}
		}
		return nil
	}
}

// AddRateLimitLogging creates a response hook that logs rate limit information.
// The logFn receives the method, path, and parsed rate limit info.
func AddRateLimitLogging(logFn func(method, path string, info *RateLimitInfo)) ResponseHook {
	return func(resp *http.Response) error {
		info := ParseRateLimitHeaders(resp.Header)
		if info != nil {
			logFn(resp.Request.Method, resp.Request.URL.Path, info)
		}
		return nil
	}
}
