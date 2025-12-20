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

// Client provides a type-safe HTTP client with rate limiting, queuing, and circuit breaking
type Client struct {
	Config ClientOptions
	mu     sync.RWMutex

	// Request queues per route (for hook integration)
	queues map[string]*RequestQueue

	// Metrics tracking
	Metrics RequestMetrics

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

	if options.QueueTimeout <= 0 {
		options.QueueTimeout = 30 * time.Second
	}

	if options.MaxQueueSize <= 0 {
		options.MaxQueueSize = 1000
	}

	if options.SweepInterval <= 0 {
		options.SweepInterval = 5 * time.Minute
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		Config: options,
		queues: make(map[string]*RequestQueue),
		ctx:    ctx,
		cancel: cancel,
	}

	// Start background workers
	if options.SweepEnabled {
		client.startSweeper()
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
		clientErr := ClientError{
			Type:      ErrorTypeRequest,
			Message:   fmt.Sprintf("failed to build request: %v", err),
			Route:     route.Path,
			Cause:     err,
			Timestamp: time.Now(),
		}
		if req != nil {
			clientErr = clientErr.WithContext(req, 0)
		}
		return nil, &clientErr
	}

	// Merge client-level and request-level hooks
	requestHooks := append(client.Config.RequestHooks, options.RequestHooks...)
	responseHooks := append(client.Config.ResponseHooks, options.ResponseHooks...)

	// Apply request hooks (client-level + per-request)
	for _, hook := range requestHooks {
		if err := hook(req); err != nil {
			clientErr := ClientError{
				Type:    ErrorTypeRequest,
				Message: fmt.Sprintf("request hook failed: %v", err),
				Route:   route.Path,
				Cause:   err,
			}.WithContext(req, 0)
			return nil, &clientErr
		}
	}

	// Execute request (resilience handled by hooks)
	resp, err := client.executeRequest(req, options)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Apply response hooks (client-level + per-request)
	for _, hook := range responseHooks {
		if err := hook(resp); err != nil {
			clientErr := ClientError{
				Type:    ErrorTypeResponse,
				Message: fmt.Sprintf("response hook failed: %v", err),
				Route:   route.Path,
				Cause:   err,
			}.WithContext(req, 0)
			return nil, &clientErr
		}
	}

	// Read response body
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		clientErr := ClientError{
			Type:    ErrorTypeResponse,
			Message: fmt.Sprintf("failed to read response body: %v", err),
			Route:   route.Path,
			Cause:   err,
		}.WithContext(req, 0)
		return nil, &clientErr
	}

	// Replace body for potential reuse
	resp.Body = io.NopCloser(bytes.NewReader(bodyData))

	// Parse response
	var response TResponse
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if len(bodyData) > 0 {
			if err := client.parseResponseBody(bodyData, &response); err != nil {
				clientErr := ClientError{
					Type:    ErrorTypeResponse,
					Message: fmt.Sprintf("failed to parse response: %v", err),
					Route:   route.Path,
					Cause:   err,
				}.WithContext(req, 0)
				return nil, &clientErr
			}
		}

		client.mu.Lock()
		client.Metrics.SuccessfulRequests++
		client.mu.Unlock()
	} else {
		client.mu.Lock()
		client.Metrics.FailedRequests++
		client.mu.Unlock()

		clientErr := ClientError{
			Type:       ErrorTypeResponse,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			StatusCode: resp.StatusCode,
			Route:      route.Path,
		}.WithContext(req, 0)
		return nil, &clientErr
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
			clientErr := ClientError{
				Type:    ErrorTypeTimeout,
				Message: "request cancelled or timed out",
				Cause:   req.Context().Err(),
			}.WithContext(req, attempt)
			return nil, &clientErr
		default:
		}

		resp, err := c.Config.HTTPClient.Do(req)
		if err == nil {
			c.mu.Lock()
			c.Metrics.TotalRequests++
			c.mu.Unlock()

			return resp, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if !c.shouldRetry(err, attempt, maxRetries) {
			break
		}

		// Calculate backoff delay with jitter
		delay := c.calculateBackoffDelayWithJitter(attempt)

		// Sleep with context awareness
		select {
		case <-time.After(delay):
			// Continue to next retry
		case <-req.Context().Done():
			clientErr := ClientError{
				Type:    ErrorTypeTimeout,
				Message: "request cancelled during retry backoff",
				Cause:   req.Context().Err(),
			}.WithContext(req, attempt)
			return nil, &clientErr
		}
	}

	c.mu.Lock()
	c.Metrics.TotalRequests++
	c.Metrics.FailedRequests++
	c.mu.Unlock()

	clientErr := ClientError{
		Type:    ErrorTypeNetwork,
		Message: fmt.Sprintf("request failed after %d attempts: %v", maxRetries+1, lastErr),
		Cause:   lastErr,
	}.WithContext(req, maxRetries)
	return nil, &clientErr
}

// Helper methods

// Helper methods for hook integration

func (c *Client) parseResponseBody(data []byte, target any) error {
	if len(data) == 0 {
		return nil
	}

	// Try custom deserializer first
	if deserializer, ok := target.(Deserializable); ok {
		return deserializer.UnmarshalJSON(data)
	}

	// Fall back to standard JSON
	return json.Unmarshal(data, target)
}

func (c *Client) shouldRetry(err error, attempt, maxRetries int) bool {
	if attempt >= maxRetries {
		return false
	}

	// Check for retryable errors
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "temporary failure") {
		return true
	}

	return false
}

func (c *Client) calculateBackoffDelay(attempt int) time.Duration {
	// Exponential backoff
	delay := float64(c.Config.RetryDelay) * (c.Config.RetryMultiplier * float64(attempt))
	return time.Duration(delay)
}

// calculateBackoffDelayWithJitter calculates backoff delay with jitter to prevent thundering herd
func (c *Client) calculateBackoffDelayWithJitter(attempt int) time.Duration {
	baseDelay := c.calculateBackoffDelay(attempt)

	// Add ±25% jitter to prevent thundering herd
	// jitter = baseDelay * (0.75 + rand(0, 0.5))
	jitterFactor := 0.75 + rand.Float64()*0.5

	return time.Duration(float64(baseDelay) * jitterFactor)
}

func (c *Client) startSweeper() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.Config.SweepInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.sweepInactiveQueues()
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *Client) sweepInactiveQueues() {
	cutoff := time.Now().Add(-c.Config.SweepInterval * 2)

	toRemove := make([]string, 0)
	for routeID, queue := range c.queues {
		if queue.LastUsed.Before(cutoff) && len(queue.Queue) == 0 && !queue.Processing {
			toRemove = append(toRemove, routeID)
		}
	}

	for _, routeID := range toRemove {
		delete(c.queues, routeID)
	}
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

// GetMetrics returns current client metrics
func (c *Client) GetMetrics() RequestMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Metrics
}

// New creates a new Neuron client with default options
func New() *Client {
	return NewWithOptions(ClientOptions{})
}

// NewWithOptions creates a new Neuron client with custom options
func NewWithOptions(options ClientOptions) *Client {
	return NewClient(options)
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
// This is the base method that Get/Post/Put/etc call internally
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

	// Build request body if provided
	var requestBody any
	if options.Body != nil && method != MethodGET && method != MethodHEAD {
		requestBody = options.Body
	}

	// Create route
	route := NewRoute[any, any](method, path)

	// Execute request (hook merging happens in Execute)
	resp, err := Execute(c, route, requestBody, options)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DoWithType executes a request and unmarshals the response to the specified type
func DoWithType[T any](client *Client, method HTTPMethod, path string, opts ...*RequestOptions) (*Response[T], error) {
	resp, err := client.Do(method, path, opts...)
	if err != nil {
		return nil, err
	}

	var target T
	if err := resp.JSON(&target); err != nil {
		return nil, err
	}

	return &Response[T]{
		Data:       target,
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Raw:        resp.Raw,
		body:       resp.body,
	}, nil
}

// Helper functions

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

// serializeBody serializes the request body
func serializeBody(data any) io.Reader {
	if data == nil {
		return nil
	}

	// Handle BodyProvider
	if provider, ok := data.(BodyProvider); ok {
		body, err := provider.Body()
		if err != nil {
			return bytes.NewReader([]byte{})
		}
		return body
	}

	// Handle io.Reader
	if reader, ok := data.(io.Reader); ok {
		return reader
	}

	// Default to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return bytes.NewReader([]byte{})
	}
	return bytes.NewReader(jsonData)
}

func isSnowflake(s string) bool {
	if len(s) < 17 || len(s) > 19 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
