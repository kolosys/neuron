package mock

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kolosys/neuron"
)

// MockRateLimiter is a mock implementation of neuron.RateLimiter for testing.
type MockRateLimiter struct {
	allowCalls   []AllowCall
	waitCalls    []WaitCall
	allowState   map[endpointKey]bool
	globalAllow  atomic.Bool
	waitDuration time.Duration
	waitErr      atomic.Value // *errHolder
	recordCalls  atomic.Bool
	mu           sync.RWMutex
}

// endpointKey uniquely identifies an endpoint.
type endpointKey struct {
	method   string
	endpoint string
}

// AllowCall records a call to Allow.
type AllowCall struct {
	Method   string
	Endpoint string
	Time     time.Time
	Result   bool
}

// WaitCall records a call to Wait.
type WaitCall struct {
	Method   string
	Endpoint string
	Time     time.Time
	Waited   time.Duration
	Err      error
}

// MockRateLimiterOptions configures mock rate limiter behavior.
type MockRateLimiterOptions struct {
	RecordCalls  bool
	DefaultAllow bool
}

// NewMockRateLimiter creates a new mock rate limiter.
func NewMockRateLimiter(opts *MockRateLimiterOptions) *MockRateLimiter {
	rl := &MockRateLimiter{
		allowCalls:   make([]AllowCall, 0),
		waitCalls:    make([]WaitCall, 0),
		allowState:   make(map[endpointKey]bool),
		waitDuration: 0,
	}

	if opts == nil {
		rl.recordCalls.Store(true)
		rl.globalAllow.Store(true)
		return rl
	}

	rl.recordCalls.Store(opts.RecordCalls)
	rl.globalAllow.Store(opts.DefaultAllow)

	return rl
}

// Allow checks if a request is allowed without blocking.
func (rl *MockRateLimiter) Allow(ctx context.Context, method, endpoint string) bool {
	result := rl.computeAllow(method, endpoint)

	if rl.recordCalls.Load() {
		rl.mu.Lock()
		rl.allowCalls = append(rl.allowCalls, AllowCall{
			Method:   method,
			Endpoint: endpoint,
			Time:     time.Now(),
			Result:   result,
		})
		rl.mu.Unlock()
	}

	return result
}

// Wait blocks until the request is allowed or context is cancelled.
func (rl *MockRateLimiter) Wait(ctx context.Context, method, endpoint string) error {
	start := time.Now()

	select {
	case <-time.After(rl.waitDuration):
	case <-ctx.Done():
		return ctx.Err()
	}

	err := rl.getAndConsumeWaitError()

	if rl.recordCalls.Load() {
		rl.mu.Lock()
		rl.waitCalls = append(rl.waitCalls, WaitCall{
			Method:   method,
			Endpoint: endpoint,
			Time:     time.Now(),
			Waited:   time.Since(start),
			Err:      err,
		})
		rl.mu.Unlock()
	}

	return err
}

// computeAllow determines if a request is allowed based on configuration.
func (rl *MockRateLimiter) computeAllow(method, endpoint string) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	key := endpointKey{method: method, endpoint: endpoint}
	if allow, ok := rl.allowState[key]; ok {
		return allow
	}

	return rl.globalAllow.Load()
}

// getAndConsumeWaitError retrieves and clears a wait error (one-shot).
func (rl *MockRateLimiter) getAndConsumeWaitError() error {
	val := rl.waitErr.Load()
	if val == nil {
		return nil
	}

	holder := val.(*errHolder)
	if holder == nil || holder.err == nil {
		return nil
	}

	err := holder.err
	if holder.oneShot {
		rl.waitErr.Store((*errHolder)(nil))
	}
	return err
}

// SetAllow sets the global allow state.
func (rl *MockRateLimiter) SetAllow(allow bool) {
	rl.globalAllow.Store(allow)
}

// SetAllowForEndpoint sets allow state for a specific endpoint.
func (rl *MockRateLimiter) SetAllowForEndpoint(method, endpoint string, allow bool) {
	rl.mu.Lock()
	key := endpointKey{method: method, endpoint: endpoint}
	rl.allowState[key] = allow
	rl.mu.Unlock()
}

// SetWaitDuration sets the duration that Wait should sleep.
func (rl *MockRateLimiter) SetWaitDuration(duration time.Duration) {
	rl.waitDuration = duration
}

// InjectWaitError injects an error for the next Wait call (one-shot by default).
func (rl *MockRateLimiter) InjectWaitError(err error) {
	if err == nil {
		rl.waitErr.Store((*errHolder)(nil))
	} else {
		rl.waitErr.Store(&errHolder{err: err, oneShot: true})
	}
}

// ClearInjectedErrors clears all injected errors.
func (rl *MockRateLimiter) ClearInjectedErrors() {
	rl.waitErr.Store((*errHolder)(nil))
}

// AllowCalls returns a copy of all recorded Allow calls.
func (rl *MockRateLimiter) AllowCalls() []AllowCall {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	calls := make([]AllowCall, len(rl.allowCalls))
	copy(calls, rl.allowCalls)
	return calls
}

// WaitCalls returns a copy of all recorded Wait calls.
func (rl *MockRateLimiter) WaitCalls() []WaitCall {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	calls := make([]WaitCall, len(rl.waitCalls))
	copy(calls, rl.waitCalls)
	return calls
}

// ClearRecorded clears all recorded calls.
func (rl *MockRateLimiter) ClearRecorded() {
	rl.mu.Lock()
	rl.allowCalls = rl.allowCalls[:0]
	rl.waitCalls = rl.waitCalls[:0]
	rl.mu.Unlock()
}

// Reset resets all state including configuration and recorded calls.
func (rl *MockRateLimiter) Reset() {
	rl.mu.Lock()
	rl.allowCalls = rl.allowCalls[:0]
	rl.waitCalls = rl.waitCalls[:0]
	rl.allowState = make(map[endpointKey]bool)
	rl.mu.Unlock()
	rl.globalAllow.Store(true)
	rl.waitErr.Store((*errHolder)(nil))
}

// MockRateLimitHandler is a mock implementation of neuron.RateLimitHandler.
// It embeds MockRateLimiter and adds UpdateFromHeaders support.
type MockRateLimitHandler struct {
	*MockRateLimiter
	updates []RateLimitUpdate
	mu      sync.RWMutex
}

// RateLimitUpdate records a call to UpdateFromHeaders.
type RateLimitUpdate struct {
	Method   string
	Endpoint string
	Info     *neuron.RateLimitInfo
	Time     time.Time
}

// NewMockRateLimitHandler creates a new mock rate limit handler.
func NewMockRateLimitHandler(opts *MockRateLimiterOptions) *MockRateLimitHandler {
	return &MockRateLimitHandler{
		MockRateLimiter: NewMockRateLimiter(opts),
		updates:         make([]RateLimitUpdate, 0),
	}
}

// UpdateFromHeaders records a call to UpdateFromHeaders.
func (rlh *MockRateLimitHandler) UpdateFromHeaders(method, endpoint string, info *neuron.RateLimitInfo) error {
	rlh.mu.Lock()
	rlh.updates = append(rlh.updates, RateLimitUpdate{
		Method:   method,
		Endpoint: endpoint,
		Info:     info,
		Time:     time.Now(),
	})
	rlh.mu.Unlock()

	return nil
}

// Updates returns a copy of all recorded UpdateFromHeaders calls.
func (rlh *MockRateLimitHandler) Updates() []RateLimitUpdate {
	rlh.mu.RLock()
	defer rlh.mu.RUnlock()

	updates := make([]RateLimitUpdate, len(rlh.updates))
	copy(updates, rlh.updates)
	return updates
}

// ClearRecordedUpdates clears all recorded updates.
func (rlh *MockRateLimitHandler) ClearRecordedUpdates() {
	rlh.mu.Lock()
	rlh.updates = rlh.updates[:0]
	rlh.mu.Unlock()
}

// Reset resets all state including updates.
func (rlh *MockRateLimitHandler) Reset() {
	rlh.MockRateLimiter.Reset()
	rlh.mu.Lock()
	rlh.updates = rlh.updates[:0]
	rlh.mu.Unlock()
}

// WasExhausted returns true if any recorded update indicated exhaustion.
func (rlh *MockRateLimitHandler) WasExhausted() bool {
	rlh.mu.RLock()
	defer rlh.mu.RUnlock()

	for _, update := range rlh.updates {
		if update.Info != nil && update.Info.IsExhausted() {
			return true
		}
	}
	return false
}

// LastUpdate returns the most recent update, or nil if none.
func (rlh *MockRateLimitHandler) LastUpdate() *RateLimitUpdate {
	rlh.mu.RLock()
	defer rlh.mu.RUnlock()

	if len(rlh.updates) == 0 {
		return nil
	}

	return &rlh.updates[len(rlh.updates)-1]
}

// UpdatesForEndpoint returns updates for a specific endpoint.
func (rlh *MockRateLimitHandler) UpdatesForEndpoint(method, endpoint string) []RateLimitUpdate {
	rlh.mu.RLock()
	defer rlh.mu.RUnlock()

	var filtered []RateLimitUpdate
	for _, update := range rlh.updates {
		if update.Method == method && update.Endpoint == endpoint {
			filtered = append(filtered, update)
		}
	}
	return filtered
}

// String returns a human-readable description of the rate limit state.
func (rl *MockRateLimiter) String() string {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	global := rl.globalAllow.Load()
	return fmt.Sprintf("MockRateLimiter{global=%v, endpoints=%d}", global, len(rl.allowState))
}
