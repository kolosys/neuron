package neuron

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// DedupeConfig configures request deduplication behavior
type DedupeConfig struct {
	// Enabled enables request deduplication
	Enabled bool

	// WindowSize is the time window for deduplication
	// Requests with the same key within this window will be deduplicated
	WindowSize time.Duration

	// MaxSize is the maximum number of entries in the dedupe cache
	MaxSize int

	// KeyGenerator generates a unique key for a request
	// If nil, defaults to method + URL
	KeyGenerator func(req *http.Request) string
}

// DefaultDedupeConfig returns sensible defaults for deduplication
func DefaultDedupeConfig() DedupeConfig {
	return DedupeConfig{
		Enabled:    true,
		WindowSize: 5 * time.Second,
		MaxSize:    10000,
	}
}

// Deduplicator prevents duplicate concurrent requests
type Deduplicator struct {
	config  DedupeConfig
	pending sync.Map

	// Metrics
	deduped   atomic.Int64
	inflight  atomic.Int64
	cleanupMu sync.Mutex
	lastClean time.Time
}

// pendingRequest tracks an in-flight request
type pendingRequest struct {
	done   chan struct{}
	result *dedupeResult
	mu     sync.Mutex
	expiry time.Time
}

// dedupeResult stores the result of a deduplicated request
type dedupeResult struct {
	body       []byte
	statusCode int
	header     http.Header
	err        error
}

// NewDeduplicator creates a new request deduplicator
func NewDeduplicator(config DedupeConfig) *Deduplicator {
	if config.WindowSize == 0 {
		config.WindowSize = 5 * time.Second
	}
	if config.MaxSize == 0 {
		config.MaxSize = 10000
	}

	return &Deduplicator{
		config:    config,
		lastClean: time.Now(),
	}
}

// Dedupe executes a request with deduplication
// If a request with the same key is already in flight, wait for its result
func (d *Deduplicator) Dedupe(ctx context.Context, key string, fn func() (*http.Response, error)) (*http.Response, error) {
	if !d.config.Enabled || key == "" {
		return fn()
	}

	d.maybeCleanup() // Cleanup stale entries periodically

	now := time.Now()
	pending := &pendingRequest{
		done:   make(chan struct{}),
		expiry: now.Add(d.config.WindowSize),
	}

	actual, loaded := d.pending.LoadOrStore(key, pending)
	if loaded {
		existingPending := actual.(*pendingRequest)
		d.deduped.Add(1)

		select {
		case <-existingPending.done:
			existingPending.mu.Lock()
			result := existingPending.result
			existingPending.mu.Unlock()

			if result == nil {
				return fn()
			}

			if result.err != nil {
				return nil, result.err
			}

			return d.cloneResponse(result), nil

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	d.inflight.Add(1)
	defer d.inflight.Add(-1)

	resp, err := fn()

	pending.mu.Lock()
	if err != nil {
		pending.result = &dedupeResult{err: err}
	} else {
		var body []byte
		var readErr error

		if resp.Body != nil {
			body, readErr = io.ReadAll(resp.Body)
			resp.Body.Close()
		}

		if readErr != nil {
			pending.result = &dedupeResult{err: readErr}
		} else {
			var header http.Header
			if resp.Header != nil {
				header = resp.Header.Clone()
			}

			pending.result = &dedupeResult{
				body:       body,
				statusCode: resp.StatusCode,
				header:     header,
			}

			resp.Body = io.NopCloser(bytes.NewReader(body))
		}
	}
	pending.mu.Unlock()

	close(pending.done) // Signal waiting requests

	go func() {
		time.Sleep(d.config.WindowSize)
		d.pending.Delete(key)
	}()

	if pending.result.err != nil {
		return nil, pending.result.err
	}

	return resp, nil
}

// cloneResponse creates a new response from cached data
func (d *Deduplicator) cloneResponse(result *dedupeResult) *http.Response {
	return &http.Response{
		StatusCode: result.statusCode,
		Header:     result.header.Clone(),
		Body:       io.NopCloser(bytes.NewReader(result.body)),
	}
}

// maybeCleanup removes expired entries if needed
func (d *Deduplicator) maybeCleanup() {
	now := time.Now()
	d.cleanupMu.Lock()
	if now.Sub(d.lastClean) < d.config.WindowSize {
		d.cleanupMu.Unlock()
		return
	}
	d.lastClean = now
	d.cleanupMu.Unlock()

	count := 0
	d.pending.Range(func(key, value any) bool {
		count++
		pending := value.(*pendingRequest)
		if now.After(pending.expiry) {
			d.pending.Delete(key)
		}
		return true
	})
}

// Stats returns deduplication statistics
func (d *Deduplicator) Stats() DedupeStats {
	return DedupeStats{
		Deduped:  d.deduped.Load(),
		Inflight: d.inflight.Load(),
	}
}

// DedupeStats contains deduplication statistics
type DedupeStats struct {
	// Deduped is the number of requests that were deduplicated
	Deduped int64

	// Inflight is the current number of in-flight requests
	Inflight int64
}

// GenerateDedupeKey creates a default deduplication key from a request
func GenerateDedupeKey(req *http.Request) string {
	return req.Method + " " + req.URL.String()
}
