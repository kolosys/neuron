package neuron

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

type RequestIDContextKey string

const (
	request_id     RequestIDContextKey = "request_id"
	correlation_id RequestIDContextKey = "correlation_id"
)

// RequestIDGenerator generates unique request IDs
type RequestIDGenerator interface {
	Generate() string
}

// UUIDGenerator generates UUID-style request IDs
type UUIDGenerator struct{}

// Generate creates a new UUID-style request ID
func (g *UUIDGenerator) Generate() string {
	// Simple UUID v4 implementation
	b := make([]byte, 16)
	rand.Read(b)

	// Set version (4) and variant bits
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// TimestampGenerator generates timestamp-based request IDs
type TimestampGenerator struct{}

// Generate creates a timestamp-based request ID
func (g *TimestampGenerator) Generate() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// SequentialGenerator generates sequential request IDs
type SequentialGenerator struct {
	counter int64
}

// Generate creates a sequential request ID
func (g *SequentialGenerator) Generate() string {
	g.counter++
	return fmt.Sprintf("req_%d", g.counter)
}

// RequestIDConfig configures request ID generation
type RequestIDConfig struct {
	Generator  RequestIDGenerator
	HeaderName string
	ContextKey RequestIDContextKey
	Propagate  bool
}

// DefaultRequestIDConfig returns a default request ID configuration
func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		Generator:  &UUIDGenerator{},
		HeaderName: "X-Request-ID",
		ContextKey: request_id,
		Propagate:  true,
	}
}

// AddRequestID creates a request ID middleware
func AddRequestID(config RequestIDConfig) RequestMiddleware {
	return func(req *http.Request) error {
		// Generate request ID
		requestID := config.Generator.Generate()

		// Add to context
		ctx := context.WithValue(req.Context(), config.ContextKey, requestID)
		*req = *req.WithContext(ctx)

		// Add to headers if propagation is enabled
		if config.Propagate {
			req.Header.Set(config.HeaderName, requestID)
		}

		return nil
	}
}

// AddResponseRequestID creates a response middleware that logs request ID
func AddResponseRequestID(config RequestIDConfig) ResponseMiddleware {
	return func(resp *http.Response) error {
		// Get request ID from context
		requestID, ok := resp.Request.Context().Value(config.ContextKey).(string)
		if ok {
			// Add request ID to response headers for tracing
			resp.Header.Set(config.HeaderName, requestID)
		}

		return nil
	}
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value("request_id").(string)
	return requestID, ok
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, request_id, requestID)
}

// RequestIDFromRequest extracts the request ID from an HTTP request
func RequestIDFromRequest(req *http.Request) (string, bool) {
	return GetRequestID(req.Context())
}

// RequestIDFromResponse extracts the request ID from an HTTP response
func RequestIDFromResponse(resp *http.Response) (string, bool) {
	return GetRequestID(resp.Request.Context())
}

// AddTracing creates a tracing middleware that propagates request ID
func AddTracing(config RequestIDConfig) RequestMiddleware {
	return func(req *http.Request) error {
		// Check if request ID already exists in headers
		existingID := req.Header.Get(config.HeaderName)
		if existingID != "" {
			// Use existing request ID
			ctx := context.WithValue(req.Context(), config.ContextKey, existingID)
			*req = *req.WithContext(ctx)
			return nil
		}

		// Generate new request ID
		requestID := config.Generator.Generate()

		// Add to context and headers
		ctx := context.WithValue(req.Context(), config.ContextKey, requestID)
		*req = *req.WithContext(ctx)
		req.Header.Set(config.HeaderName, requestID)

		return nil
	}
}

// AddCorrelationID creates a correlation ID middleware for distributed tracing
func AddCorrelationID(config RequestIDConfig) RequestMiddleware {
	return func(req *http.Request) error {
		// Check for existing correlation ID
		correlationID := req.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			// Generate new correlation ID
			correlationID = config.Generator.Generate()
			req.Header.Set("X-Correlation-ID", correlationID)
		}

		// Add to context
		ctx := context.WithValue(req.Context(), correlation_id, correlationID)
		*req = *req.WithContext(ctx)

		return nil
	}
}
