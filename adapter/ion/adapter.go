package ion

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kolosys/ion/circuit"
	"github.com/kolosys/ion/ratelimit"
	"github.com/kolosys/neuron/adapter"
)

// IonAdapter implements the adapter interface for ion integration
type IonAdapter struct {
	*adapter.BaseAdapter

	// Ion components
	rateLimiter    ratelimit.Limiter
	circuitBreaker circuit.CircuitBreaker

	// Configuration
	rateLimitConfig      adapter.RateLimitConfig
	circuitBreakerConfig adapter.CircuitBreakerConfig
}

// NewAdapter creates a new ion adapter
func NewAdapter() *IonAdapter {
	return &IonAdapter{
		BaseAdapter: adapter.NewBaseAdapter("ion"),
	}
}

// WithRateLimiter configures the rate limiter
func (a *IonAdapter) WithRateLimiter(limiter ratelimit.Limiter) *IonAdapter {
	a.rateLimiter = limiter
	return a
}

// WithCircuitBreaker configures the circuit breaker
func (a *IonAdapter) WithCircuitBreaker(circuitBreaker circuit.CircuitBreaker) *IonAdapter {
	a.circuitBreaker = circuitBreaker
	return a
}

// ConfigureRateLimiting configures rate limiting for the adapter
func (a *IonAdapter) ConfigureRateLimiting(config adapter.RateLimitConfig) *IonAdapter {
	a.rateLimitConfig = config

	// Create ion rate limiter if not already set
	if a.rateLimiter == nil {
		rate := ratelimit.PerSecond(config.RequestsPerSecond)
		a.rateLimiter = ratelimit.NewTokenBucket(rate, config.BurstSize)
	}

	return a
}

// ConfigureCircuitBreaker configures circuit breaker for the adapter
func (a *IonAdapter) ConfigureCircuitBreaker(config adapter.CircuitBreakerConfig) *IonAdapter {
	a.circuitBreakerConfig = config

	// Create ion circuit breaker if not already set
	if a.circuitBreaker == nil {
		options := []circuit.Option{
			circuit.WithFailureThreshold(int64(config.FailureThreshold)),
			circuit.WithRecoveryTimeout(time.Duration(config.RecoveryTimeout) * time.Second),
			circuit.WithHalfOpenMaxRequests(int64(config.HalfOpenMaxRequests)),
			circuit.WithHalfOpenSuccessThreshold(int64(config.SuccessThreshold)),
		}

		if config.FailurePredicate != nil {
			options = append(options, circuit.WithFailurePredicate(config.FailurePredicate))
		}

		a.circuitBreaker = circuit.New("neuron-ion-adapter", options...)
	}

	return a
}

// WrapHTTPClient wraps an HTTP client with ion functionality
func (a *IonAdapter) WrapHTTPClient(client *http.Client) *http.Client {
	if a.circuitBreaker == nil {
		return client
	}

	// Create a circuit breaker protected HTTP client
	return &http.Client{
		Transport: &circuitBreakerTransport{
			client:         client,
			circuitBreaker: a.circuitBreaker,
		},
	}
}

// CreateRequestMiddleware creates request middleware for the adapter
func (a *IonAdapter) CreateRequestMiddleware() []adapter.RequestMiddleware {
	middleware := []adapter.RequestMiddleware{}

	// Add rate limiting middleware
	if a.rateLimiter != nil {
		middleware = append(middleware, a.createRateLimitMiddleware())
	}

	// Add circuit breaker middleware (for logging/state checking)
	if a.circuitBreaker != nil {
		middleware = append(middleware, a.createCircuitBreakerMiddleware())
	}

	return middleware
}

// CreateResponseMiddleware creates response middleware for the adapter
func (a *IonAdapter) CreateResponseMiddleware() []adapter.ResponseMiddleware {
	middleware := []adapter.ResponseMiddleware{}

	// Add circuit breaker response middleware (for logging)
	if a.circuitBreaker != nil {
		middleware = append(middleware, a.createCircuitBreakerResponseMiddleware())
	}

	return middleware
}

// Shutdown gracefully shuts down the adapter
func (a *IonAdapter) Shutdown(ctx context.Context) error {
	if a.circuitBreaker != nil {
		return a.circuitBreaker.Close()
	}
	return nil
}

// Helper methods

func (a *IonAdapter) createRateLimitMiddleware() adapter.RequestMiddleware {
	return func(req *http.Request) error {
		return a.rateLimiter.WaitN(req.Context(), 1)
	}
}

func (a *IonAdapter) createCircuitBreakerMiddleware() adapter.RequestMiddleware {
	return func(req *http.Request) error {
		// Log circuit breaker state
		state := a.circuitBreaker.State()
		fmt.Printf("[ION] Circuit breaker state for %s: %v\n", req.URL.Path, state)
		return nil
	}
}

func (a *IonAdapter) createCircuitBreakerResponseMiddleware() adapter.ResponseMiddleware {
	return func(resp *http.Response) error {
		// Log circuit breaker response
		state := a.circuitBreaker.State()
		fmt.Printf("[ION] Circuit breaker response for %s (status: %d, state: %v)\n",
			resp.Request.URL.Path, resp.StatusCode, state)
		return nil
	}
}

// circuitBreakerTransport implements http.RoundTripper with circuit breaker protection
type circuitBreakerTransport struct {
	client         *http.Client
	circuitBreaker circuit.CircuitBreaker
}

func (t *circuitBreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Use ion's circuit breaker Execute method
	result, err := t.circuitBreaker.Execute(req.Context(), func(ctx context.Context) (any, error) {
		// Update request context
		reqWithCtx := req.WithContext(ctx)

		// Execute the actual HTTP request
		resp, err := t.client.Do(reqWithCtx)
		if err != nil {
			return nil, err
		}

		// Check for HTTP error status codes and convert them to errors for circuit breaker
		if resp.StatusCode >= 500 {
			return resp, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}

		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*http.Response), nil
}

// Factory functions for creating ion components

// NewRateLimiter creates a new ion rate limiter
func NewRateLimiter(requestsPerSecond, burst int) ratelimit.Limiter {
	rate := ratelimit.PerSecond(requestsPerSecond)
	return ratelimit.NewTokenBucket(rate, burst)
}

// NewCircuitBreaker creates a new ion circuit breaker with default options
func NewCircuitBreaker(name string) circuit.CircuitBreaker {
	return circuit.New(name,
		circuit.WithFailureThreshold(5),
		circuit.WithRecoveryTimeout(30*time.Second),
		circuit.WithHalfOpenMaxRequests(3),
		circuit.WithHalfOpenSuccessThreshold(2),
	)
}

// NewCircuitBreakerWithOptions creates a new ion circuit breaker with custom options
func NewCircuitBreakerWithOptions(name string, options ...circuit.Option) circuit.CircuitBreaker {
	return circuit.New(name, options...)
}
