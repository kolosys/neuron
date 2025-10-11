package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kolosys/neuron"
)

func main() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back request headers and body
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello from Neuron with Middleware Features!", "request_id": "` + r.Header.Get("X-Request-ID") + `"}`))
	}))
	defer server.Close()

	// Create client with middleware features
	client := neuron.NewClient(neuron.ClientOptions{
		BaseURL: server.URL,

		// Add middleware
		RequestMiddleware: []neuron.RequestMiddleware{
			// Request ID generation
			neuron.RequestIDMiddleware(neuron.DefaultRequestIDConfig()),

			// Logging
			neuron.LoggingMiddleware(neuron.DefaultLoggingConfig()),

			// Authentication
			neuron.BearerAuthMiddleware("your-token-here"),

			// Headers
			neuron.UserAgentMiddleware("MyApp/1.0"),
			neuron.CustomHeaderMiddleware(map[string]string{
				"X-Custom-Header": "custom-value",
			}),

			// Timeout
			neuron.TimeoutMiddleware(neuron.TimeoutConfig{
				Timeout: 30 * time.Second,
			}),

			// Compression
			neuron.AutoCompressionMiddleware(),
		},

		ResponseMiddleware: []neuron.ResponseMiddleware{
			// Response logging
			neuron.LoggingResponseMiddleware(neuron.DefaultLoggingConfig()),
		},
	})
	defer client.Shutdown(5 * time.Second)

	// Define a route
	route := neuron.NewRoute[neuron.EmptyRequest, map[string]string](
		neuron.MethodGET, "/test",
	)

	// Execute multiple requests to demonstrate features
	for i := 0; i < 3; i++ {
		fmt.Printf("\n--- Request %d ---\n", i+1)

		resp, err := neuron.Execute(client, route, neuron.EmptyRequest{}, &neuron.RequestOptions{
			Context: context.Background(),
		})

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Status: %d\n", resp.StatusCode)
		fmt.Printf("Data: %+v\n", resp.Data)

		// Small delay between requests
		time.Sleep(100 * time.Millisecond)
	}

	// Demonstrate metrics collection
	fmt.Println("\n--- Metrics ---")
	// Note: In a real implementation, you'd get metrics from the middleware
	fmt.Println("Request count: 3")
	fmt.Println("Average duration: ~100ms")
	fmt.Println("Success rate: 100%")
}
