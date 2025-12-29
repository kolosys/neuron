package mock

import (
	"sync"
	"sync/atomic"
	"testing"
)

// CallbackRecorder records the number of times a callback has been invoked.
// It's designed to be called as a callback function and provides assertion helpers.
type CallbackRecorder struct {
	calls atomic.Int64
	mu    sync.Mutex
}

// NewCallbackRecorder creates a new callback recorder.
func NewCallbackRecorder() *CallbackRecorder {
	return &CallbackRecorder{}
}

// Record increments the call counter. This is designed to be used as a callback.
func (cr *CallbackRecorder) Record() {
	cr.calls.Add(1)
}

// Calls returns the current call count.
func (cr *CallbackRecorder) Calls() int64 {
	return cr.calls.Load()
}

// Reset resets the call counter to zero.
func (cr *CallbackRecorder) Reset() {
	cr.calls.Store(0)
}

// AssertCalled fails the test if the callback was never called.
func (cr *CallbackRecorder) AssertCalled(t *testing.T) {
	t.Helper()
	if cr.calls.Load() == 0 {
		t.Error("callback was not called")
	}
}

// AssertNotCalled fails the test if the callback was called at least once.
func (cr *CallbackRecorder) AssertNotCalled(t *testing.T) {
	t.Helper()
	if cr.calls.Load() > 0 {
		t.Error("callback was called but should not have been")
	}
}

// AssertCallCount fails the test if the call count doesn't match the expected value.
func (cr *CallbackRecorder) AssertCallCount(t *testing.T, expected int64) {
	t.Helper()
	if actual := cr.calls.Load(); actual != expected {
		t.Errorf("callback call count: expected %d, got %d", expected, actual)
	}
}

// CacheOperation records a cache operation for analysis.
type CacheOperation struct {
	Op    string // "get", "set", "delete", "clear"
	Key   string
	Time  int64 // nanoseconds since epoch
	Hit   bool  // for Get operations
	Entry any   // the entry involved
}

// RequestCapture captures the essential details of an HTTP request for analysis.
type RequestCapture struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    []byte
	Time    int64 // nanoseconds since epoch
}

// ResponseCapture captures the essential details of an HTTP response for analysis.
type ResponseCapture struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
	Time       int64 // nanoseconds since epoch
}
