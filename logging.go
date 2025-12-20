package neuron

import (
	"context"
	"log"
	"net/http"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// LoggingConfig configures the logging hooks
type LoggingConfig struct {
	Level       LogLevel
	IncludeBody bool
	MaxBodySize int
	Logger      *log.Logger
}

// DefaultLoggingConfig returns a default logging configuration
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:       LogLevelInfo,
		IncludeBody: false,
		MaxBodySize: 1024, // 1KB max body size for logging
		Logger:      log.Default(),
	}
}

// AddLogging creates a logging hook
func AddLogging(config LoggingConfig) RequestHook {
	return func(req *http.Request) error {
		start := time.Now()

		// Store start time in context for response logging
		c := context.WithValue(req.Context(), requestStartKey, start)
		*req = *req.WithContext(c)

		// Log request
		if config.Level <= LogLevelInfo {
			config.Logger.Printf("[REQUEST] %s %s", req.Method, req.URL.String())
		}

		return nil
	}
}

// AddResponseLogging creates a response logging middleware
func AddResponseLogging(config LoggingConfig) ResponseHook {
	return func(resp *http.Response) error {
		// Get start time from context
		start, ok := resp.Request.Context().Value(requestStartKey).(time.Time)
		if !ok {
			start = time.Now()
		}

		duration := time.Since(start)

		// Log response
		config.Logger.Printf("[RESPONSE] %s %s - %d (%v)",
			resp.Request.Method,
			resp.Request.URL.String(),
			resp.StatusCode,
			duration,
		)

		return nil
	}
}

// AddDebugLogging creates a debug logging middleware with detailed information
func AddDebugLogging(config LoggingConfig) RequestHook {
	return func(req *http.Request) error {
		start := time.Now()

		// Store start time in context
		c := context.WithValue(req.Context(), requestStartKey, start)
		*req = *req.WithContext(c)

		// Log detailed request information
		config.Logger.Printf("[DEBUG] Request: %s %s", req.Method, req.URL.String())
		config.Logger.Printf("[DEBUG] Headers: %v", req.Header)

		if config.IncludeBody && req.Body != nil {
			// Note: In a real implementation, you'd need to read and restore the body
			config.Logger.Printf("[DEBUG] Body: [body content would be logged here]")
		}

		return nil
	}
}

// AddResponseDebug creates a debug response logging middleware
func AddResponseDebug(config LoggingConfig) ResponseHook {
	return func(resp *http.Response) error {
		// Get start time from context
		start, ok := resp.Request.Context().Value(requestStartKey).(time.Time)
		if !ok {
			start = time.Now()
		}

		duration := time.Since(start)

		// Log detailed response information
		config.Logger.Printf("[DEBUG] Response: %d %s", resp.StatusCode, resp.Status)
		config.Logger.Printf("[DEBUG] Duration: %v", duration)
		config.Logger.Printf("[DEBUG] Headers: %v", resp.Header)

		return nil
	}
}

// AddResponseErrorLogging creates an error logging middleware
func AddResponseErrorLogging(config LoggingConfig) ResponseHook {
	return func(resp *http.Response) error {
		// Only log errors (4xx, 5xx status codes)
		if resp.StatusCode >= 400 {
			config.Logger.Printf("[ERROR] %s %s - %d %s",
				resp.Request.Method,
				resp.Request.URL.String(),
				resp.StatusCode,
				resp.Status,
			)
		}

		return nil
	}
}

// AddStructuredLogging creates a structured logging middleware
func AddStructuredLogging(config LoggingConfig) RequestHook {
	return func(req *http.Request) error {
		start := time.Now()

		// Store start time in context
		c := context.WithValue(req.Context(), requestStartKey, start)
		*req = *req.WithContext(c)

		// Structured log entry
		config.Logger.Printf(`{"level":"info","type":"request","method":"%s","url":"%s","timestamp":"%s"}`,
			req.Method,
			req.URL.String(),
			start.Format(time.RFC3339),
		)

		return nil
	}
}

// AddResponseStructured creates a structured response logging middleware
func AddResponseStructured(config LoggingConfig) ResponseHook {
	return func(resp *http.Response) error {
		// Get start time from context
		start, ok := resp.Request.Context().Value(requestStartKey).(time.Time)
		if !ok {
			start = time.Now()
		}

		duration := time.Since(start)

		// Structured log entry
		config.Logger.Printf(`{"level":"info","type":"response","method":"%s","url":"%s","status":%d,"duration":"%v","timestamp":"%s"}`,
			resp.Request.Method,
			resp.Request.URL.String(),
			resp.StatusCode,
			duration,
			time.Now().Format(time.RFC3339),
		)

		return nil
	}
}
