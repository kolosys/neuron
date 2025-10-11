package neuron

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	config := DefaultLoggingConfig()
	middleware := LoggingMiddleware(config)
	req := httptest.NewRequest("GET", "/test/endpoint", nil)

	err := middleware(req)
	if err != nil {
		t.Errorf("logging middleware failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "GET") || !strings.Contains(output, "/test/endpoint") {
		t.Errorf("log output missing request info: %s", output)
	}
}

func TestLoggingResponseMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	config := DefaultLoggingConfig()
	middleware := LoggingResponseMiddleware(config)
	req := httptest.NewRequest("GET", "/test", nil)
	resp := &http.Response{
		StatusCode: 200,
		Request:    req,
	}

	err := middleware(resp)
	if err != nil {
		t.Errorf("logging response middleware failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "200") {
		t.Errorf("log output missing status code: %s", output)
	}
}

func TestDebugLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	config := DefaultLoggingConfig()
	config.Level = LogLevelDebug
	config.IncludeBody = true
	middleware := DebugLoggingMiddleware(config)
	req := httptest.NewRequest("POST", "/test", strings.NewReader("test body"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")

	err := middleware(req)
	if err != nil {
		t.Errorf("debug logging middleware failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "POST") {
		t.Errorf("log output missing method: %s", output)
	}
	if !strings.Contains(output, "Content-Type") {
		t.Errorf("log output missing headers: %s", output)
	}
}

func TestErrorLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	config := DefaultLoggingConfig()
	middleware := ErrorLoggingMiddleware(config)
	req := httptest.NewRequest("GET", "/test", nil)

	tests := []struct {
		name       string
		statusCode int
		shouldLog  bool
	}{
		{"success response", 200, false},
		{"client error", 400, true},
		{"server error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Request:    req,
				Status:     http.StatusText(tt.statusCode),
			}

			err := middleware(resp)
			if err != nil {
				t.Errorf("error logging middleware failed: %v", err)
			}

			output := buf.String()
			hasLog := len(output) > 0

			if tt.shouldLog && !hasLog {
				t.Error("expected error to be logged")
			}
			if !tt.shouldLog && hasLog {
				t.Error("expected no error logging")
			}
		})
	}
}

func TestStructuredLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	config := DefaultLoggingConfig()
	middleware := StructuredLoggingMiddleware(config)
	req := httptest.NewRequest("GET", "/api/users/123", nil)
	req.Header.Set("X-Request-ID", "req-123")

	err := middleware(req)
	if err != nil {
		t.Errorf("structured logging middleware failed: %v", err)
	}

	output := buf.String()
	// Check for structured fields
	if !strings.Contains(output, "method") {
		t.Errorf("log missing method field: %s", output)
	}
	if !strings.Contains(output, "GET") {
		t.Errorf("log missing method value: %s", output)
	}
	if !strings.Contains(output, "/api/users/123") {
		t.Errorf("log missing path: %s", output)
	}
}
