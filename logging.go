package neuron

import (
	"context"
	"log"
	"net/http"
	"strings"
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

// safeHeaderNames is a whitelist of header names that are safe to log
var safeHeaderNames = map[string]bool{
	"content-type":    true,
	"content-length":  true,
	"accept":          true,
	"accept-encoding": true,
	"accept-language": true,
	"user-agent":      true,
	"referer":         true,
	"origin":          true,
	"cache-control":   true,
	"connection":      true,
	"date":            true,
	"server":          true,
	"x-forwarded-for": true,
	"x-real-ip":       true,
}

// sensitiveHeaderNames is a set of header names that contain sensitive data
var sensitiveHeaderNames = map[string]bool{
	"authorization":   true,
	"cookie":          true,
	"set-cookie":      true,
	"x-api-key":       true,
	"x-auth-token":    true,
	"x-access-token":  true,
	"x-refresh-token": true,
	"api-key":         true,
	"auth-token":      true,
}

// safeHeadersForLogging creates a copy containing only whitelisted safe headers
// This ensures no sensitive data can flow to logging calls
func safeHeadersForLogging(headers http.Header) http.Header {
	safe := make(http.Header)
	for key, values := range headers {
		keyLower := strings.ToLower(key)
		// Only include explicitly whitelisted safe headers
		if safeHeaderNames[keyLower] {
			// Create a new slice to avoid sharing the underlying array
			safe[key] = append([]string(nil), values...)
		}
	}
	return safe
}

// sanitizeHeadersForLogging redacts sensitive header values before logging
// This function ensures no sensitive data flows to logging calls
func sanitizeHeadersForLogging(headers http.Header) http.Header {
	sanitized := make(http.Header)
	for key, values := range headers {
		keyLower := strings.ToLower(key)
		if sensitiveHeaderNames[keyLower] {
			// Redact sensitive headers
			sanitized[key] = []string{"<redacted>"}
		} else {
			// Copy non-sensitive headers as-is
			sanitized[key] = append([]string(nil), values...)
		}
	}
	return sanitized
}

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
		// Build safe headers map - only include whitelisted headers to prevent sensitive data flow
		safeHeaders := safeHeadersForLogging(req.Header)
		// Sanitize headers before logging to ensure no sensitive data flows to logs
		sanitizedHeaders := sanitizeHeadersForLogging(safeHeaders)
		if len(sanitizedHeaders) > 0 {
			config.Logger.Printf("[DEBUG] Headers: %v", sanitizedHeaders)
		}

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
		// Build safe headers map - only include whitelisted headers to prevent sensitive data flow
		safeHeaders := safeHeadersForLogging(resp.Header)
		// Sanitize headers before logging to ensure no sensitive data flows to logs
		sanitizedHeaders := sanitizeHeadersForLogging(safeHeaders)
		if len(sanitizedHeaders) > 0 {
			config.Logger.Printf("[DEBUG] Headers: %v", sanitizedHeaders)
		}

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
