package neuron

import (
	"context"
	"net/http"
	"time"
)

// TimeoutConfig configures timeout middleware
type TimeoutConfig struct {
	Timeout        time.Duration
	PerRequest     bool
	GlobalTimeout  time.Duration
	RequestTimeout time.Duration
}

// DefaultTimeoutConfig returns a default timeout configuration
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout:        30 * time.Second,
		PerRequest:     true,
		GlobalTimeout:  60 * time.Second,
		RequestTimeout: 30 * time.Second,
	}
}

// AddTimeout creates a timeout middleware
func AddTimeout(config TimeoutConfig) RequestMiddleware {
	return func(req *http.Request) error {
		// Create timeout context
		ctx, cancel := context.WithTimeout(req.Context(), config.Timeout)
		defer cancel()

		// Update request context
		*req = *req.WithContext(ctx)

		return nil
	}
}

// AddPerRequestTimeout creates a per-request timeout middleware
func AddPerRequestTimeout(timeout time.Duration) RequestMiddleware {
	return func(req *http.Request) error {
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		*req = *req.WithContext(ctx)

		return nil
	}
}

// AddGlobalTimeout creates a global timeout middleware
func AddGlobalTimeout(timeout time.Duration) RequestMiddleware {
	return func(req *http.Request) error {
		// Check if request already has a timeout
		_, hasTimeout := req.Context().Deadline()
		if hasTimeout {
			return nil
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		*req = *req.WithContext(ctx)

		return nil
	}
}

// AddAdaptiveTimeout creates an adaptive timeout middleware
func AddAdaptiveTimeout(baseTimeout time.Duration, multiplier float64) RequestMiddleware {
	return func(req *http.Request) error {
		// Calculate adaptive timeout based on request characteristics
		timeout := baseTimeout

		// Adjust timeout based on request method
		switch req.Method {
		case "GET":
			timeout = time.Duration(float64(timeout) * 0.8) // Shorter for GET
		case "POST", "PUT", "PATCH":
			timeout = time.Duration(float64(timeout) * 1.2) // Longer for write operations
		case "DELETE":
			timeout = time.Duration(float64(timeout) * 1.1) // Slightly longer for DELETE
		}

		// Apply multiplier
		timeout = time.Duration(float64(timeout) * multiplier)

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		*req = *req.WithContext(ctx)

		return nil
	}
}

// AddConditionalTimeout creates a conditional timeout middleware
func AddConditionalTimeout(condition func(*http.Request) bool, timeout time.Duration) RequestMiddleware {
	return func(req *http.Request) error {
		if condition(req) {
			ctx, cancel := context.WithTimeout(req.Context(), timeout)
			defer cancel()

			*req = *req.WithContext(ctx)
		}

		return nil
	}
}

// AddTimeoutFromContext creates a timeout middleware that gets timeout from context
func AddTimeoutFromContext(contextKey string, defaultTimeout time.Duration) RequestMiddleware {
	return func(req *http.Request) error {
		timeout := defaultTimeout

		if ctxTimeout, ok := req.Context().Value(contextKey).(time.Duration); ok {
			timeout = ctxTimeout
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		*req = *req.WithContext(ctx)

		return nil
	}
}

// AddResponseTimeout creates a response timeout middleware
func AddResponseTimeout(timeout time.Duration) ResponseMiddleware {
	return func(resp *http.Response) error {
		// Check if response took too long
		start, ok := resp.Request.Context().Value("request_start").(time.Time)
		if !ok {
			return nil
		}

		duration := time.Since(start)
		if duration > timeout {
			// Log or handle timeout
			return nil
		}

		return nil
	}
}

// AddDeadline creates a deadline middleware
func AddDeadline(deadline time.Time) RequestMiddleware {
	return func(req *http.Request) error {
		ctx, cancel := context.WithDeadline(req.Context(), deadline)
		defer cancel()

		*req = *req.WithContext(ctx)

		return nil
	}
}

// AddTimeoutChain creates a chain of timeout middlewares
func AddTimeoutChain(timeouts ...time.Duration) RequestMiddleware {
	return func(req *http.Request) error {
		ctx := req.Context()

		for _, timeout := range timeouts {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		*req = *req.WithContext(ctx)

		return nil
	}
}
