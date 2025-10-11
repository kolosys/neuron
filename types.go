package neuron

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
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

	// Type hints for compile-time type safety
	RequestType  func() TRequest
	ResponseType func() TResponse
}

// NewRoute creates a new type-safe route
func NewRoute[TRequest, TResponse any](method HTTPMethod, path string) Route[TRequest, TResponse] {
	return Route[TRequest, TResponse]{
		Method:       method,
		Path:         path,
		RequestType:  func() TRequest { var zero TRequest; return zero },
		ResponseType: func() TResponse { var zero TResponse; return zero },
	}
}

// RequestOptions contains configuration for individual requests
type RequestOptions struct {
	Headers     http.Header
	Query       map[string]any
	Timeout     *time.Duration
	Context     context.Context
	Retries     *int
	RateLimitID string // Custom rate limit bucket ID
}

// Response wraps HTTP response data with type safety
type Response[T any] struct {
	Data       T
	StatusCode int
	Headers    http.Header
	Raw        *http.Response
}

// ClientOptions configures the HTTP client behavior
type ClientOptions struct {
	// Base configuration
	BaseURL   string
	UserAgent string
	Headers   http.Header
	Timeout   time.Duration

	// Rate limiting
	GlobalRateLimit   RateLimiter
	PerRouteRateLimit bool
	RateLimitConfig   RateLimitConfig

	// Circuit breaker configuration
	CircuitBreakerConfig CircuitBreakerConfig

	// Request handling
	MaxRetries      int
	RetryDelay      time.Duration
	RetryMultiplier float64

	// Queue management
	QueueTimeout time.Duration
	MaxQueueSize int

	// Middleware
	RequestMiddleware  []RequestMiddleware
	ResponseMiddleware []ResponseMiddleware

	// Adapter for external integrations
	Adapter SimpleAdapter

	// HTTP client
	HTTPClient *http.Client

	// Sweeping configuration
	SweepInterval time.Duration
	SweepEnabled  bool
}

// RateLimitConfig defines rate limiting behavior
type RateLimitConfig struct {
	// Global rate limiting
	GlobalRequestsPerSecond int
	GlobalBurstSize         int

	// Per-route rate limiting
	RouteRequestsPerSecond int
	RouteBurstSize         int

	// Rate limit detection
	RespectDiscordHeaders bool
	BackoffStrategy       BackoffStrategy

	// Queue behavior on rate limits
	QueueOnRateLimit bool
	RateLimitTimeout time.Duration
}

// BackoffStrategy defines how to handle backoff
type BackoffStrategy int

const (
	BackoffExponential BackoffStrategy = iota
	BackoffLinear
	BackoffFixed
)

// RequestMiddleware processes requests before they are sent
type RequestMiddleware func(req *http.Request) error

// ResponseMiddleware processes responses after they are received
type ResponseMiddleware func(resp *http.Response) error

// Adapter defines the interface for external library integrations
type Adapter interface {
	// Name returns the adapter name
	Name() string

	// WrapHTTPClient wraps an HTTP client with adapter functionality
	WrapHTTPClient(client *http.Client) *http.Client

	// CreateRequestMiddleware creates request middleware for the adapter
	CreateRequestMiddleware() []RequestMiddleware

	// CreateResponseMiddleware creates response middleware for the adapter
	CreateResponseMiddleware() []ResponseMiddleware

	// Shutdown gracefully shuts down the adapter
	Shutdown(ctx context.Context) error
}

// Adapter configuration types are defined in adapter modules

// RequestQueue manages queued requests for a specific route/bucket
type RequestQueue struct {
	Queue       []QueuedRequest
	RateLimiter RateLimiter
	Processing  bool
	LastUsed    time.Time
}

// QueuedRequest represents a request waiting in queue
type QueuedRequest struct {
	Request     *http.Request
	ResponseCh  chan *QueuedResponse
	Context     context.Context
	Retries     int
	EnqueueTime time.Time
}

// QueuedResponse represents the result of a queued request
type QueuedResponse struct {
	Response *http.Response
	Error    error
}

// RequestContext provides metadata about the current request
type RequestContext struct {
	RouteID   string
	Attempt   int
	QueueTime time.Duration
	StartTime time.Time
	Metadata  map[string]any
}

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
	Context    RequestContext
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
	ErrorTypeRateLimit
	ErrorTypeTimeout
	ErrorTypeQueue
	ErrorTypeResponse
	ErrorTypeNetwork
	ErrorTypeAuth
	ErrorTypeCircuitBreaker
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

// RateLimitInfo contains information about current rate limit status
type RateLimitInfo struct {
	RouteID    string
	Bucket     string
	Limit      int
	Remaining  int
	ResetAfter time.Duration
	RetryAfter time.Duration
	Global     bool
}

// RequestMetrics provides insights into request performance
type RequestMetrics struct {
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	QueuedRequests      int64
	AverageQueueTime    time.Duration
	AverageResponseTime time.Duration
	RateLimitHits       int64
}

// EmptyRequest represents requests with no body
type EmptyRequest struct{}

// EmptyResponse represents responses with no body
type EmptyResponse struct{}

// CircuitBreakerConfig defines circuit breaker behavior for HTTP clients
type CircuitBreakerConfig struct {
	// Enable circuit breaker functionality
	Enabled bool

	// Per-route circuit breakers (vs single global circuit breaker)
	PerRouteCircuitBreakers bool

	// Circuit breaker options
	FailureThreshold    int
	RecoveryTimeout     time.Duration
	HalfOpenMaxRequests int
	SuccessThreshold    int
}

// DefaultCircuitBreakerConfig returns sensible defaults for HTTP clients
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Enabled:                 true,
		PerRouteCircuitBreakers: true,
		FailureThreshold:        5,
		RecoveryTimeout:         30 * time.Second,
		HalfOpenMaxRequests:     3,
		SuccessThreshold:        2,
	}
}

// Middleware interfaces for external library integration

// RateLimiter interface for middleware integration
type RateLimiter interface {
	WaitN(ctx context.Context, n int) error
	AllowN(now time.Time, n int) bool
	Tokens() float64
	Burst() int
}

// CircuitBreaker interface for middleware integration
type CircuitBreaker interface {
	AllowRequest() bool
	RecordSuccess()
	RecordFailure()
	GetState() CircuitBreakerState
	Close() error
}

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState string

const (
	CircuitBreakerClosed   CircuitBreakerState = "closed"
	CircuitBreakerOpen     CircuitBreakerState = "open"
	CircuitBreakerHalfOpen CircuitBreakerState = "half-open"
)
