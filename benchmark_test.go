package neuron

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// BenchmarkExecute benchmarks the main Execute function
func BenchmarkExecute(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success","id":123}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Execute(client, route, EmptyRequest{}, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkExecuteWithMiddleware benchmarks Execute with middleware
func BenchmarkExecuteWithMiddleware(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	logConfig := DefaultLoggingConfig()
	client := NewClient(ClientOptions{
		BaseURL: server.URL,
		RequestMiddleware: []RequestMiddleware{
			LoggingMiddleware(logConfig),
			func(req *http.Request) error {
				req.Header.Set("X-Custom-Header", "value")
				return nil
			},
		},
		ResponseMiddleware: []ResponseMiddleware{
			LoggingResponseMiddleware(logConfig),
		},
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Execute(client, route, EmptyRequest{}, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkExecuteWithMetrics benchmarks Execute with metrics collection
func BenchmarkExecuteWithMetrics(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	collector := NewMetricsCollector()
	client := NewClient(ClientOptions{
		BaseURL: server.URL,
		RequestMiddleware: []RequestMiddleware{
			MetricsMiddleware(collector),
		},
		ResponseMiddleware: []ResponseMiddleware{
			MetricsResponseMiddleware(collector),
		},
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Execute(client, route, EmptyRequest{}, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkRetryWithBackoff benchmarks retry logic with backoff
func BenchmarkRetryWithBackoff(b *testing.B) {
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount%3 != 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL:    server.URL,
		MaxRetries: 3,
		RetryDelay: 1 * time.Millisecond,
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Execute(client, route, EmptyRequest{}, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkMiddlewareChain benchmarks middleware chain execution
func BenchmarkMiddlewareChain(b *testing.B) {
	chain := NewMiddlewareChain()

	// Add 5 middleware
	for i := 0; i < 5; i++ {
		chain.AddRequestMiddleware(func(req *http.Request) error {
			req.Header.Set(fmt.Sprintf("X-Middleware-%d", i), "value")
			return nil
		})
	}

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := chain.ApplyRequestMiddleware(req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkMetricsCollector benchmarks metrics collection
func BenchmarkMetricsCollector(b *testing.B) {
	collector := NewMetricsCollector()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		collector.RecordRequest()
		collector.RecordResponse(200, time.Duration(i)*time.Millisecond)
	}
}

// BenchmarkJitterCalculation benchmarks jitter calculation
func BenchmarkJitterCalculation(b *testing.B) {
	client := NewClient(ClientOptions{
		BaseURL:    "https://api.example.com",
		RetryDelay: 100 * time.Millisecond,
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = client.calculateBackoffDelayWithJitter(i % 5)
	}
}

// BenchmarkAuthMiddleware benchmarks authentication middleware
func BenchmarkAuthMiddleware(b *testing.B) {
	middleware := BearerAuthMiddleware("test-token")
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := middleware(req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkConcurrentRequests benchmarks concurrent request handling
func BenchmarkConcurrentRequests(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Execute(client, route, EmptyRequest{}, nil)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})
}

// BenchmarkQuerySerialization benchmarks query parameter serialization
func BenchmarkQuerySerialization(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := NewClient(ClientOptions{
		BaseURL: server.URL,
	})

	route := NewRoute[EmptyRequest, TestResponse](MethodGET, "/test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Execute(client, route, EmptyRequest{}, &RequestOptions{
			Query: map[string]any{
				"string": "test",
				"int":    42,
				"bool":   true,
				"float":  3.14,
			},
		})
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkErrorCreation benchmarks enhanced error creation
func BenchmarkErrorCreation(b *testing.B) {
	req := httptest.NewRequest("GET", "https://api.example.com/test", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := ClientError{
			Type:    ErrorTypeNetwork,
			Message: "test error",
			Route:   "/test",
		}.WithContext(req, i)
		_ = err.Error()
	}
}
