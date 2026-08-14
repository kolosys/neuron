package neuron

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kolosys/ion/circuit"
)

// HTTPMethod represents supported HTTP methods
type HTTPMethod string

const (
	MethodGET     HTTPMethod = "GET"
	MethodPOST    HTTPMethod = "POST"
	MethodPUT     HTTPMethod = "PUT"
	MethodPATCH   HTTPMethod = "PATCH"
	MethodDELETE  HTTPMethod = "DELETE"
	MethodHEAD    HTTPMethod = "HEAD"
	MethodOPTIONS HTTPMethod = "OPTIONS"
)

// Route represents a type-safe route definition
type Route[TRequest, TResponse any] struct {
	Method HTTPMethod
	Path   string
}

// NewRoute creates a new type-safe route
func NewRoute[TRequest, TResponse any](method HTTPMethod, path string) Route[TRequest, TResponse] {
	return Route[TRequest, TResponse]{
		Method: method,
		Path:   path,
	}
}

// RequestOptions contains configuration for individual requests
type RequestOptions struct {
	Headers http.Header
	Query   map[string]any
	Timeout *time.Duration
	Context context.Context
	Retries *int
	Body    any // Request body (JSON, form data, io.Reader, BodyProvider)

	// Per-request hooks (runs after client-level hooks)
	RequestHooks  []RequestHook
	ResponseHooks []ResponseHook

	// IdempotencyKey is used for request deduplication
	// If set, concurrent requests with the same key will share the response
	IdempotencyKey string

	// DisableDedupe skips deduplication for this request
	DisableDedupe bool
}

// Response wraps HTTP response data with type safety
type Response[T any] struct {
	Data       T
	StatusCode int
	Headers    http.Header
	Raw        *http.Response
	body       []byte // Cached body for helper methods
}

// IsSuccess returns true if the status code is in the 2xx range
func (r *Response[T]) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if the status code is 4xx or 5xx
func (r *Response[T]) IsError() bool {
	return r.StatusCode >= 400
}

// IsClientError returns true if the status code is in the 4xx range
func (r *Response[T]) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the status code is in the 5xx range
func (r *Response[T]) IsServerError() bool {
	return r.StatusCode >= 500
}

// JSON unmarshals the response body as JSON
func (r *Response[T]) JSON(target any) error {
	if r.Raw == nil || r.Raw.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	data, err := r.Bytes()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// String returns the response body as a string
func (r *Response[T]) String() (string, error) {
	data, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Bytes returns the response body as bytes
func (r *Response[T]) Bytes() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}

	if r.Raw == nil || r.Raw.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	data, err := io.ReadAll(r.Raw.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	r.body = data
	return data, nil
}

// ClientOptions configures the HTTP client behavior
type ClientOptions struct {
	// Base configuration
	BaseURL   string
	UserAgent string
	Headers   http.Header
	Timeout   time.Duration

	// Request handling
	MaxRetries      int
	RetryDelay      time.Duration
	RetryMultiplier float64

	// Hooks
	RequestHooks  []RequestHook
	ResponseHooks []ResponseHook

	// HTTP client
	HTTPClient *http.Client

	// Circuit, if set, wraps HTTPClient.Do in circuit.Execute.
	// An open circuit fails fast as ClientError with Type ErrorTypeCircuit.
	Circuit circuit.CircuitBreaker

	// Resilience configuration
	AdaptiveTimeout bool

	// Rate limiting
	RateLimiter         RateLimiter
	RateLimitUpdater    RateLimitUpdater
	AutoHandleRateLimit bool

	// Deduplication
	Deduplicator *Deduplicator
}

// RequestHook processes requests before they are sent
type RequestHook func(req *http.Request) error

// ResponseHook processes responses after they are received
type ResponseHook func(resp *http.Response) error

// Error types for type-safe error handling
type ClientError struct {
	Type       ErrorType
	Message    string
	StatusCode int
	Route      string
	Method     string    // HTTP method
	URL        string    // Full URL
	Attempt    int       // Retry attempt number
	Timestamp  time.Time // When error occurred
	Cause      error
}

func (e ClientError) Error() string {
	if e.Method != "" && e.URL != "" {
		if e.Attempt > 0 {
			return fmt.Sprintf("%s %s (attempt %d): %s", e.Method, e.URL, e.Attempt, e.Message)
		}
		return fmt.Sprintf("%s %s: %s", e.Method, e.URL, e.Message)
	}
	return e.Message
}

func (e ClientError) Unwrap() error {
	return e.Cause
}

// WithContext adds request context information to the error
func (e ClientError) WithContext(req *http.Request, attempt int) ClientError {
	e.Method = req.Method
	e.URL = req.URL.String()
	e.Attempt = attempt
	e.Timestamp = time.Now()
	return e
}

type ErrorType int

const (
	ErrorTypeRequest ErrorType = iota
	ErrorTypeTimeout
	ErrorTypeResponse
	ErrorTypeNetwork
	ErrorTypeAuth
	ErrorTypeRateLimit
	ErrorTypeCircuit
)

// Serializable represents types that can be serialized for requests
type Serializable interface {
	MarshalJSON() ([]byte, error)
}

// Deserializable represents types that can be deserialized from responses
type Deserializable interface {
	UnmarshalJSON(data []byte) error
}

// BodyProvider allows custom body serialization
type BodyProvider interface {
	ContentType() string
	Body() (io.Reader, error)
}

// RequestMetrics provides a snapshot of insights into request performance
type RequestMetrics struct {
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	AverageResponseTime time.Duration
}

// EmptyRequest represents requests with no body
type EmptyRequest struct{}

// EmptyResponse represents responses with no body
type EmptyResponse struct{}
