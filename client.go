package neuron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// Client provides a type-safe HTTP client with rate limiting and circuit breaking
type Client struct {
	Config ClientOptions

	// Metrics tracking
	Metrics *MetricsCollector

	// Shutdown management
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewClient creates a new type-safe HTTP client
func NewClient(options ClientOptions) *Client {
	// Set defaults
	if options.HTTPClient == nil {
		options.HTTPClient = &http.Client{
			Timeout: options.Timeout,
		}
	}

	if options.Headers == nil {
		options.Headers = make(http.Header)
	}

	if options.UserAgent != "" {
		options.Headers.Set("User-Agent", options.UserAgent)
	}

	if options.MaxRetries <= 0 {
		options.MaxRetries = 3
	}

	if options.RetryDelay <= 0 {
		options.RetryDelay = time.Second
	}

	if options.RetryMultiplier <= 0 {
		options.RetryMultiplier = 2.0
	}

	c, cancel := context.WithCancel(context.Background())

	client := &Client{
		Config:  options,
		Metrics: NewMetricsCollector(),
		ctx:     c,
		cancel:  cancel,
	}

	return client
}

// Execute performs a type-safe HTTP request
func Execute[TRequest any, TResponse any](
	client *Client,
	route Route[TRequest, TResponse],
	request TRequest,
	options *RequestOptions,
) (*Response[TResponse], error) {
	if options == nil {
		options = &RequestOptions{}
	}

	if options.Context == nil {
		options.Context = context.Background()
	}

	// Build the request
	req, err := client.buildRequest(route.Method, route.Path, request, options)
	if err != nil {
		ce := ClientError{
			Type:      ErrorTypeRequest,
			Message:   fmt.Sprintf("failed to build request: %v", err),
			Route:     route.Path,
			Cause:     err,
			Timestamp: time.Now(),
		}
		if req != nil {
			ce = ce.WithContext(req, 0)
		}
		return nil, ce
	}

	// Apply adaptive timeout if enabled
	reqCtx := req.Context()
	if client.Config.AdaptiveTimeout {
		timeout := calculateAdaptiveTimeout(req.Method, client.Config.Timeout)
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(reqCtx, timeout)
		defer cancel()
	}

	// Store start time for metrics/logging
	startTime := time.Now()
	reqCtx = context.WithValue(reqCtx, requestStartKey, startTime)
	req = req.WithContext(reqCtx)

	// Merge client-level and request-level hooks
	requestHooks := append(client.Config.RequestHooks, options.RequestHooks...)
	responseHooks := append(client.Config.ResponseHooks, options.ResponseHooks...)

	// Apply request hooks
	for _, hook := range requestHooks {
		if err := hook(req); err != nil {
			if ce, ok := err.(ClientError); ok {
				return nil, ce
			}
			return nil, ClientError{
				Type:    ErrorTypeRequest,
				Message: fmt.Sprintf("request hook failed: %v", err),
				Route:   route.Path,
				Cause:   err,
			}.WithContext(req, 0)
		}
	}

	// Execute request
	resp, err := client.executeRequest(req, options)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Record metrics after request completes
	duration := time.Since(startTime)
	client.Metrics.RecordResponse(resp.StatusCode, duration)

	// Apply response hooks
	for _, hook := range responseHooks {
		if err := hook(resp); err != nil {
			if ce, ok := err.(ClientError); ok {
				return nil, ce
			}
			return nil, ClientError{
				Type:    ErrorTypeResponse,
				Message: fmt.Sprintf("response hook failed: %v", err),
				Route:   route.Path,
				Cause:   err,
			}.WithContext(req, 0)
		}
	}

	// Read response body once
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ClientError{
			Type:    ErrorTypeResponse,
			Message: fmt.Sprintf("failed to read response body: %v", err),
			Route:   route.Path,
			Cause:   err,
		}.WithContext(req, 0)
	}

	// Parse response
	var response TResponse
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if len(bodyData) > 0 {
			if err := client.parseResponseBody(bodyData, &response); err != nil {
				return nil, ClientError{
					Type:    ErrorTypeResponse,
					Message: fmt.Sprintf("failed to parse response: %v", err),
					Route:   route.Path,
					Cause:   err,
				}.WithContext(req, 0)
			}
		}
	} else {
		return nil, ClientError{
			Type:       ErrorTypeResponse,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			StatusCode: resp.StatusCode,
			Route:      route.Path,
		}.WithContext(req, 0)
	}

	return &Response[TResponse]{
		Data:       response,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Raw:        resp,
		body:       bodyData,
	}, nil
}

// buildRequest constructs an HTTP request from route and data
func (c *Client) buildRequest(method HTTPMethod, path string, data any, options *RequestOptions) (*http.Request, error) {
	// Build URL
	fullURL := c.Config.BaseURL + path

	// Add query parameters
	if len(options.Query) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for key, value := range options.Query {
			if str := serializeQueryParam(value); str != "" {
				q.Add(key, str)
			}
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	// Serialize body
	var body io.Reader
	if data != nil && method != MethodGET && method != MethodHEAD {
		body = serializeBody(data)
	}

	// Create request
	req, err := http.NewRequestWithContext(options.Context, string(method), fullURL, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, values := range c.Config.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if options.Headers != nil {
		for key, values := range options.Headers {
			for _, value := range values {
				req.Header.Set(key, value) // Override defaults
			}
		}
	}

	// Set content type for body requests
	if body != nil && req.Header.Get("Content-Type") == "" {
		if provider, ok := data.(BodyProvider); ok {
			req.Header.Set("Content-Type", provider.ContentType())
		} else if _, ok := data.(io.Reader); !ok {
			// Only set JSON for non-Reader types
			req.Header.Set("Content-Type", "application/json")
		}
	}

	return req, nil
}

// executeRequest executes a request (resilience handled by middleware)
func (c *Client) executeRequest(req *http.Request, options *RequestOptions) (*http.Response, error) {
	// Execute request with retries
	return c.executeWithRetries(req, options)
}

// Simple request execution - resilience handled by middleware

// executeWithRetries executes a request with retry logic
func (c *Client) executeWithRetries(req *http.Request, options *RequestOptions) (*http.Response, error) {
	maxRetries := c.Config.MaxRetries
	if options.Retries != nil {
		maxRetries = *options.Retries
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check context cancellation before attempting request
		select {
		case <-req.Context().Done():
			return nil, ClientError{
				Type:    ErrorTypeTimeout,
				Message: "request cancelled or timed out",
				Cause:   req.Context().Err(),
			}.WithContext(req, attempt)
		default:
		}

		resp, err := c.Config.HTTPClient.Do(req)
		if err == nil {
			if resp.StatusCode >= 500 && attempt < maxRetries {
				// Retry on server errors
				resp.Body.Close()
				goto retry
			}
			c.Metrics.RequestCount.Add(1)
			return resp, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if !c.shouldRetry(err, attempt, maxRetries) {
			break
		}

	retry:
		// Calculate backoff delay with jitter
		delay := c.calculateBackoffDelayWithJitter(attempt)

		// Sleep with context awareness
		select {
		case <-time.After(delay):
			// Continue to next retry
		case <-req.Context().Done():
			return nil, ClientError{
				Type:    ErrorTypeTimeout,
				Message: "request cancelled during retry backoff",
				Cause:   req.Context().Err(),
			}.WithContext(req, attempt)
		}
	}

	c.Metrics.RequestCount.Add(1)
	c.Metrics.ErrorCount.Add(1)

	return nil, ClientError{
		Type:    ErrorTypeNetwork,
		Message: fmt.Sprintf("request failed after %d attempts: %v", maxRetries+1, lastErr),
		Cause:   lastErr,
	}.WithContext(req, maxRetries)
}

// GetMetrics returns current client metrics
func (c *Client) GetMetrics() MetricsSnapshot {
	return c.Metrics.GetMetrics()
}

// Get executes a GET request and returns the response
func (c *Client) Get(path string, opts ...*RequestOptions) (*Response[any], error) {
	return c.Do(MethodGET, path, opts...)
}

// Post executes a POST request and returns the response
func (c *Client) Post(path string, opts ...*RequestOptions) (*Response[any], error) {
	return c.Do(MethodPOST, path, opts...)
}

// Put executes a PUT request and returns the response
func (c *Client) Put(path string, opts ...*RequestOptions) (*Response[any], error) {
	return c.Do(MethodPUT, path, opts...)
}

// Patch executes a PATCH request and returns the response
func (c *Client) Patch(path string, opts ...*RequestOptions) (*Response[any], error) {
	return c.Do(MethodPATCH, path, opts...)
}

// Delete executes a DELETE request and returns the response
func (c *Client) Delete(path string, opts ...*RequestOptions) (*Response[any], error) {
	return c.Do(MethodDELETE, path, opts...)
}

// Head executes a HEAD request and returns the response
func (c *Client) Head(path string, opts ...*RequestOptions) (*Response[any], error) {
	return c.Do(MethodHEAD, path, opts...)
}

// Options executes an OPTIONS request and returns the response
func (c *Client) Options(path string, opts ...*RequestOptions) (*Response[any], error) {
	return c.Do(MethodOPTIONS, path, opts...)
}

// Do executes an HTTP request with the specified method and returns the response
func (c *Client) Do(method HTTPMethod, path string, opts ...*RequestOptions) (*Response[any], error) {
	var options *RequestOptions
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	} else {
		options = &RequestOptions{}
	}

	if options.Context == nil {
		options.Context = context.Background()
	}

	// Create route
	route := NewRoute[any, any](method, path)

	// Execute request (hook merging and body serialization happen in Execute/buildRequest)
	return Execute(c, route, options.Body, options)
}

// DoWithType executes a request and unmarshals the response to the specified type
func DoWithType[T any](client *Client, method HTTPMethod, path string, opts ...*RequestOptions) (*Response[T], error) {
	var options *RequestOptions
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	} else {
		options = &RequestOptions{}
	}

	route := NewRoute[any, T](method, path)
	return Execute(client, route, options.Body, options)
}

// serializeBody serializes the request body
func serializeBody(data any) io.Reader {
	if data == nil {
		return nil
	}

	if provider, ok := data.(BodyProvider); ok {
		body, err := provider.Body()
		if err != nil {
			return bytes.NewReader(nil)
		}
		return body
	}

	if reader, ok := data.(io.Reader); ok {
		return reader
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		buf.Reset()
		bufferPool.Put(buf)
		return bytes.NewReader(nil)
	}

	return &pooledBufferReader{
		buf: buf,
	}
}

type pooledBufferReader struct {
	buf *bytes.Buffer
	off int
}

func (r *pooledBufferReader) Read(p []byte) (n int, err error) {
	if r.buf == nil {
		return 0, io.EOF
	}
	n, err = r.buf.Read(p)
	if err == io.EOF {
		r.Close()
	}
	return n, err
}

func (r *pooledBufferReader) Close() error {
	if r.buf != nil {
		r.buf.Reset()
		bufferPool.Put(r.buf)
		r.buf = nil
	}
	return nil
}

func (c *Client) parseResponseBody(data []byte, target any) error {
	if len(data) == 0 {
		return nil
	}
	if deserializer, ok := target.(Deserializable); ok {
		return deserializer.UnmarshalJSON(data)
	}
	return json.Unmarshal(data, target)
}

func (c *Client) shouldRetry(err error, attempt, maxRetries int) bool {
	if attempt >= maxRetries {
		return false
	}

	if err != nil {
		msg := err.Error()
		return strings.Contains(msg, "connection refused") ||
			strings.Contains(msg, "timeout") ||
			strings.Contains(msg, "temporary failure")
	}

	return false
}

func (c *Client) calculateBackoffDelay(attempt int) time.Duration {
	delay := float64(c.Config.RetryDelay) * (c.Config.RetryMultiplier * float64(attempt))
	return time.Duration(delay)
}

func (c *Client) calculateBackoffDelayWithJitter(attempt int) time.Duration {
	baseDelay := c.calculateBackoffDelay(attempt)
	jitterFactor := 0.75 + rand.Float64()*0.5
	return time.Duration(float64(baseDelay) * jitterFactor)
}

// Shutdown gracefully shuts down the client
func (c *Client) Shutdown(timeout time.Duration) error {
	c.cancel()
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("shutdown timeout after %v", timeout)
	}
}

func serializeQueryParam(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func calculateAdaptiveTimeout(method string, baseTimeout time.Duration) time.Duration {
	if baseTimeout <= 0 {
		baseTimeout = 30 * time.Second
	}

	switch method {
	case string(MethodGET):
		return time.Duration(float64(baseTimeout) * 0.8)
	case string(MethodPOST), string(MethodPUT), string(MethodPATCH):
		return time.Duration(float64(baseTimeout) * 1.2)
	case string(MethodDELETE):
		return time.Duration(float64(baseTimeout) * 1.1)
	default:
		return baseTimeout
	}
}
