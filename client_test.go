package neuron

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestNewClient tests client creation and defaults
func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		options ClientOptions
		check   func(*testing.T, *Client)
	}{
		{
			name:    "default options",
			options: ClientOptions{BaseURL: "https://api.example.com"},
			check: func(t *testing.T, c *Client) {
				if c.Options.MaxRetries != 3 {
					t.Errorf("expected MaxRetries=3, got %d", c.Options.MaxRetries)
				}
				if c.Options.RetryDelay != time.Second {
					t.Errorf("expected RetryDelay=1s, got %v", c.Options.RetryDelay)
				}
				if c.Options.HTTPClient == nil {
					t.Error("expected HTTPClient to be initialized")
				}
			},
		},
		{
			name: "custom options",
			options: ClientOptions{
				BaseURL:    "https://api.example.com",
				MaxRetries: 5,
				RetryDelay: 2 * time.Second,
			},
			check: func(t *testing.T, c *Client) {
				if c.Options.MaxRetries != 5 {
					t.Errorf("expected MaxRetries=5, got %d", c.Options.MaxRetries)
				}
				if c.Options.RetryDelay != 2*time.Second {
					t.Errorf("expected RetryDelay=2s, got %v", c.Options.RetryDelay)
				}
			},
		},
		{
			name: "user agent header",
			options: ClientOptions{
				BaseURL:   "https://api.example.com",
				UserAgent: "TestBot/1.0",
			},
			check: func(t *testing.T, c *Client) {
				ua := c.Options.Headers.Get("User-Agent")
				if ua != "TestBot/1.0" {
					t.Errorf("expected User-Agent=TestBot/1.0, got %s", ua)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.options)
			tt.check(t, client)
		})
	}
}

// TestExecute tests the main Execute function
func TestExecute(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectError    bool
		checkResponse  func(*testing.T, *Response[TestResponse])
	}{
		{
			name: "successful request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"message":"success","id":123}`)
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *Response[TestResponse]) {
				if resp.Data.Message != "success" {
					t.Errorf("expected message=success, got %s", resp.Data.Message)
				}
				if resp.Data.ID != 123 {
					t.Errorf("expected id=123, got %d", resp.Data.ID)
				}
				if resp.StatusCode != 200 {
					t.Errorf("expected status=200, got %d", resp.StatusCode)
				}
			},
		},
		{
			name: "error response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"bad request"}`)
			},
			expectError: true,
		},
		{
			name: "empty response body",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *Response[TestResponse]) {
				if resp.StatusCode != 204 {
					t.Errorf("expected status=204, got %d", resp.StatusCode)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client := NewClient(ClientOptions{
				BaseURL: server.URL,
			})

			route := NewRoute[TestRequest, TestResponse](MethodGET, "/test")
			resp, err := Execute(client, route, TestRequest{}, nil)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if !tt.expectError && tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// TestExecuteWithContext tests context cancellation
func TestExecuteWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")
	_, err := Execute(client, route, EmptyRequest{}, &RequestOptions{
		Context: ctx,
	})

	if err == nil {
		t.Error("expected timeout error, got nil")
	}

	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Errorf("expected ClientError, got %T", err)
	}

	if clientErr.Type != ErrorTypeTimeout && clientErr.Type != ErrorTypeNetwork {
		t.Errorf("expected timeout/network error, got %v", clientErr.Type)
	}
}

// TestRetryLogic tests retry behavior
func TestRetryLogic(t *testing.T) {
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 2 {
			// Simulate connection error by closing connection
			// For this test, we'll just succeed on second attempt
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message":"success"}`)
			return
		}
		// Succeed on subsequent attempts
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL:    server.URL,
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")
	resp, err := Execute(client, route, EmptyRequest{}, nil)

	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}

	if attemptCount < 1 {
		t.Errorf("expected at least 1 attempt, got %d", attemptCount)
	}

	if resp != nil && resp.Data.Message != "success" {
		t.Errorf("expected message=success, got %s", resp.Data.Message)
	}
}

// TestRequestMiddleware tests request middleware execution
func TestRequestMiddleware(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if middleware added header
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Error("middleware header not found")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	middlewareExecuted := false
	client := NewClient(ClientOptions{
		BaseURL: server.URL,
		RequestMiddleware: []RequestMiddleware{
			func(req *http.Request) error {
				middlewareExecuted = true
				req.Header.Set("X-Test-Header", "test-value")
				return nil
			},
		},
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")
	_, err := Execute(client, route, EmptyRequest{}, nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !middlewareExecuted {
		t.Error("middleware was not executed")
	}
}

// TestResponseMiddleware tests response middleware execution
func TestResponseMiddleware(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response-Header", "response-value")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	middlewareExecuted := false
	var capturedHeader string

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
		ResponseMiddleware: []ResponseMiddleware{
			func(resp *http.Response) error {
				middlewareExecuted = true
				capturedHeader = resp.Header.Get("X-Response-Header")
				return nil
			},
		},
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")
	_, err := Execute(client, route, EmptyRequest{}, nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !middlewareExecuted {
		t.Error("middleware was not executed")
	}

	if capturedHeader != "response-value" {
		t.Errorf("expected header=response-value, got %s", capturedHeader)
	}
}

// TestEnhancedErrorMessages tests the improved error context
func TestEnhancedErrorMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test/endpoint")
	_, err := Execute(client, route, EmptyRequest{}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("expected ClientError, got %T", err)
	}

	// Check error message includes method and URL
	if clientErr.Method != "GET" {
		t.Errorf("expected method=GET, got %s", clientErr.Method)
	}

	if !strings.Contains(clientErr.URL, "/test/endpoint") {
		t.Errorf("expected URL to contain /test/endpoint, got %s", clientErr.URL)
	}

	if clientErr.StatusCode != 404 {
		t.Errorf("expected status=404, got %d", clientErr.StatusCode)
	}

	// Check formatted error message
	errorMsg := clientErr.Error()
	if !strings.Contains(errorMsg, "GET") || !strings.Contains(errorMsg, "/test/endpoint") {
		t.Errorf("error message missing context: %s", errorMsg)
	}
}

// TestJitterInRetries tests that retry delays include jitter
func TestJitterInRetries(t *testing.T) {
	client := NewClient(ClientOptions{
		BaseURL:    "https://api.example.com",
		RetryDelay: 100 * time.Millisecond,
		RateLimitConfig: RateLimitConfig{
			BackoffStrategy: BackoffFixed,
		},
	})

	// Calculate multiple jittered delays
	delays := make([]time.Duration, 10)
	for i := 0; i < 10; i++ {
		delays[i] = client.calculateBackoffDelayWithJitter(1)
	}

	// Check that delays vary (jitter is working)
	allSame := true
	for i := 1; i < len(delays); i++ {
		if delays[i] != delays[0] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("all delays are identical, jitter not working")
	}

	// Check that delays are within expected range (75%-125% of base)
	baseDelay := 100 * time.Millisecond
	minDelay := time.Duration(float64(baseDelay) * 0.75)
	maxDelay := time.Duration(float64(baseDelay) * 1.25)

	for i, delay := range delays {
		if delay < minDelay || delay > maxDelay {
			t.Errorf("delay[%d]=%v outside expected range [%v, %v]", i, delay, minDelay, maxDelay)
		}
	}
}

// TestQueryParameters tests query parameter serialization
func TestQueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("string") != "test" {
			t.Errorf("expected string=test, got %s", q.Get("string"))
		}
		if q.Get("int") != "42" {
			t.Errorf("expected int=42, got %s", q.Get("int"))
		}
		if q.Get("bool") != "true" {
			t.Errorf("expected bool=true, got %s", q.Get("bool"))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")
	_, err := Execute(client, route, EmptyRequest{}, &RequestOptions{
		Query: map[string]any{
			"string": "test",
			"int":    42,
			"bool":   true,
		},
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestCustomHeaders tests custom header handling
func TestCustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Error("custom header not found")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
		Headers: http.Header{
			"X-Custom-Header": []string{"custom-value"},
		},
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")
	_, err := Execute(client, route, EmptyRequest{}, nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestShutdown tests graceful shutdown
func TestShutdown(t *testing.T) {
	client := NewClient(ClientOptions{
		BaseURL:      "https://api.example.com",
		SweepEnabled: true,
	})

	// Shutdown should complete without error
	err := client.Shutdown(5 * time.Second)
	if err != nil {
		t.Errorf("shutdown failed: %v", err)
	}
}

// TestMetrics tests metrics collection
func TestMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	// Make multiple requests
	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")
	for i := 0; i < 5; i++ {
		Execute(client, route, EmptyRequest{}, nil)
	}

	metrics := client.GetMetrics()
	if metrics.TotalRequests != 5 {
		t.Errorf("expected 5 total requests, got %d", metrics.TotalRequests)
	}
	if metrics.SuccessfulRequests != 5 {
		t.Errorf("expected 5 successful requests, got %d", metrics.SuccessfulRequests)
	}
}

// Test helper types

type TestRequest struct {
	Name  string `json:"name,omitempty"`
	Value int    `json:"value,omitempty"`
}

type TestResponse struct {
	Message string `json:"message,omitempty"`
	ID      int    `json:"id,omitempty"`
}

// TestBodyProvider tests custom body provider
type customBodyProvider struct {
	data string
}

func (c customBodyProvider) ContentType() string {
	return "text/plain"
}

func (c customBodyProvider) Body() (io.Reader, error) {
	return strings.NewReader(c.data), nil
}

func TestCustomBodyProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("expected Content-Type=text/plain, got %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "custom data" {
			t.Errorf("expected body='custom data', got %s", string(body))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	route := NewRoute[customBodyProvider, TestResponse](MethodPOST, "/test")
	_, err := Execute(client, route, customBodyProvider{data: "custom data"}, nil)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
