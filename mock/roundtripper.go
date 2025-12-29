package mock

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// MockRoundTripper is a mock implementation of http.RoundTripper for testing.
// It allows queuing responses, recording requests, and injecting errors.
type MockRoundTripper struct {
	responses    []ResponseConfig
	requests     []RecordedRequest
	matchers     []ResponseMatcher
	responseErr  atomic.Value // *errHolder
	latency      time.Duration
	recordReqs   atomic.Bool
	requestCount atomic.Int64
	mu           sync.RWMutex
}

// ResponseConfig represents a configured response.
type ResponseConfig struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Err        error
	Delay      time.Duration
}

// RecordedRequest represents a recorded HTTP request.
type RecordedRequest struct {
	Method  string
	URL     *url.URL
	Headers http.Header
	Body    []byte
	Time    time.Time
}

// ResponseMatcher is a function that matches a request and returns a response.
type ResponseMatcher func(req *http.Request) (*http.Response, bool)

// MockRoundTripperOptions configures MockRoundTripper behavior.
type MockRoundTripperOptions struct {
	RecordRequests bool
	Latency        time.Duration
	BufferSize     int
}

// NewMockRoundTripper creates a new mock round tripper.
func NewMockRoundTripper(opts *MockRoundTripperOptions) *MockRoundTripper {
	rt := &MockRoundTripper{
		responses: make([]ResponseConfig, 0),
		requests:  make([]RecordedRequest, 0),
		matchers:  make([]ResponseMatcher, 0),
	}

	if opts == nil {
		rt.recordReqs.Store(true)
		return rt
	}

	rt.recordReqs.Store(opts.RecordRequests)
	rt.latency = opts.Latency

	if opts.BufferSize > 0 {
		rt.responses = make([]ResponseConfig, 0, opts.BufferSize)
		rt.requests = make([]RecordedRequest, 0, opts.BufferSize)
	}

	return rt
}

// RoundTrip implements http.RoundTripper by returning queued responses or executing matchers.
func (rt *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.requestCount.Add(1)

	// Record the request if enabled
	if rt.recordReqs.Load() {
		rt.recordRequest(req)
	}

	// Check for injected error
	if err := rt.getAndConsumeError(); err != nil {
		return nil, err
	}

	// Try matchers first
	rt.mu.RLock()
	for _, matcher := range rt.matchers {
		if resp, matched := matcher(req); matched {
			rt.mu.RUnlock()

			if resp.Request == nil {
				resp.Request = req
			}
			return rt.applyLatency(resp), nil
		}
	}
	rt.mu.RUnlock()

	// Fall back to queued responses
	config := rt.dequeueResponse()
	if config == nil {
		return nil, io.EOF
	}

	if config.Err != nil {
		return nil, config.Err
	}

	resp := &http.Response{
		Status:        http.StatusText(config.StatusCode),
		StatusCode:    config.StatusCode,
		Header:        config.Headers,
		Body:          io.NopCloser(bytes.NewReader(config.Body)),
		ContentLength: int64(len(config.Body)),
		Request:       req,
	}

	return rt.applyLatency(resp), nil
}

// recordRequest captures request details.
func (rt *MockRoundTripper) recordRequest(req *http.Request) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	rt.mu.Lock()
	rt.requests = append(rt.requests, RecordedRequest{
		Method:  req.Method,
		URL:     req.URL,
		Headers: req.Header.Clone(),
		Body:    body,
		Time:    time.Now(),
	})
	rt.mu.Unlock()
}

// dequeueResponse removes and returns the next queued response.
func (rt *MockRoundTripper) dequeueResponse() *ResponseConfig {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if len(rt.responses) == 0 {
		return nil
	}

	config := rt.responses[0]
	rt.responses = rt.responses[1:]
	return &config
}

// getAndConsumeError retrieves and clears an injected error (one-shot).
func (rt *MockRoundTripper) getAndConsumeError() error {
	val := rt.responseErr.Load()
	if val == nil {
		return nil
	}

	holder := val.(*errHolder)
	if holder == nil || holder.err == nil {
		return nil
	}

	err := holder.err
	if holder.oneShot {
		rt.responseErr.Store((*errHolder)(nil))
	}
	return err
}

// applyLatency applies configured latency to a response.
func (rt *MockRoundTripper) applyLatency(resp *http.Response) *http.Response {
	delay := rt.latency
	if resp == nil {
		return resp
	}

	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-resp.Request.Context().Done():
			return resp
		}
	}
	return resp
}

// QueueResponse adds a response to the response queue.
func (rt *MockRoundTripper) QueueResponse(config ResponseConfig) {
	rt.mu.Lock()
	rt.responses = append(rt.responses, config)
	rt.mu.Unlock()
}

// AddMatcher adds a response matcher for pattern-based matching.
func (rt *MockRoundTripper) AddMatcher(matcher ResponseMatcher) {
	rt.mu.Lock()
	rt.matchers = append(rt.matchers, matcher)
	rt.mu.Unlock()
}

// InjectError injects an error for the next request (one-shot by default).
func (rt *MockRoundTripper) InjectError(err error) {
	if err == nil {
		rt.responseErr.Store((*errHolder)(nil))
	} else {
		rt.responseErr.Store(&errHolder{err: err, oneShot: true})
	}
}

// InjectErrorSequence injects a sequence of errors for successive requests.
func (rt *MockRoundTripper) InjectErrorSequence(errs []error) {
	seq := NewErrorSequence(errs...)
	rt.responseErr.Store(&errHolder{
		err:     seq,
		oneShot: false,
	})
}

// ClearInjectedErrors clears all injected errors.
func (rt *MockRoundTripper) ClearInjectedErrors() {
	rt.responseErr.Store((*errHolder)(nil))
}

// Requests returns a copy of all recorded requests.
func (rt *MockRoundTripper) Requests() []RecordedRequest {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	reqs := make([]RecordedRequest, len(rt.requests))
	copy(reqs, rt.requests)
	return reqs
}

// RequestCount returns the total number of requests made.
func (rt *MockRoundTripper) RequestCount() int64 {
	return rt.requestCount.Load()
}

// ClearRecorded clears all recorded requests and responses.
func (rt *MockRoundTripper) ClearRecorded() {
	rt.mu.Lock()
	rt.requests = rt.requests[:0]
	rt.responses = rt.responses[:0]
	rt.mu.Unlock()
	rt.requestCount.Store(0)
}

// Reset resets all state including requests, responses, and injected errors.
func (rt *MockRoundTripper) Reset() {
	rt.mu.Lock()
	rt.requests = rt.requests[:0]
	rt.responses = rt.responses[:0]
	rt.matchers = rt.matchers[:0]
	rt.mu.Unlock()
	rt.responseErr.Store((*errHolder)(nil))
	rt.requestCount.Store(0)
}

// SetLatency sets the simulated network latency.
func (rt *MockRoundTripper) SetLatency(latency time.Duration) {
	rt.latency = latency
}

// EnableRecording enables request recording.
func (rt *MockRoundTripper) EnableRecording(enabled bool) {
	rt.recordReqs.Store(enabled)
}
