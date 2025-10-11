package adapter

import (
	"context"
	"net/http"
)

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

// RequestMiddleware processes requests before they are sent
type RequestMiddleware func(req *http.Request) error

// ResponseMiddleware processes responses after they are received
type ResponseMiddleware func(resp *http.Response) error

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int
	BurstSize         int
	QueueOnRateLimit  bool
	Timeout           int // seconds
}

// CircuitBreakerConfig defines circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold        int
	RecoveryTimeout         int // seconds
	HalfOpenMaxRequests     int
	SuccessThreshold        int
	FailurePredicate        func(error) bool
	PerRouteCircuitBreakers bool
}

// BaseAdapter provides common functionality for all adapters
type BaseAdapter struct {
	name string
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(name string) *BaseAdapter {
	return &BaseAdapter{name: name}
}

// Name returns the adapter name
func (b *BaseAdapter) Name() string {
	return b.name
}

// Default implementations that can be overridden
func (b *BaseAdapter) WrapHTTPClient(client *http.Client) *http.Client {
	return client
}

func (b *BaseAdapter) CreateRequestMiddleware() []RequestMiddleware {
	return nil
}

func (b *BaseAdapter) CreateResponseMiddleware() []ResponseMiddleware {
	return nil
}

func (b *BaseAdapter) Shutdown(ctx context.Context) error {
	return nil
}
