package neuron_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kolosys/neuron"
)

// MockAdapter is a simple adapter for testing
type MockAdapter struct {
	name string
}

func (m *MockAdapter) Name() string {
	return m.name
}

func (m *MockAdapter) WrapHTTPClient(client *http.Client) *http.Client {
	// Return the client as-is for testing
	return client
}

func (m *MockAdapter) CreateRequestMiddleware() []neuron.RequestMiddleware {
	return []neuron.RequestMiddleware{
		func(req *http.Request) error {
			fmt.Printf("[MOCK] Request to %s\n", req.URL.Path)
			return nil
		},
	}
}

func (m *MockAdapter) CreateResponseMiddleware() []neuron.ResponseMiddleware {
	return []neuron.ResponseMiddleware{
		func(resp *http.Response) error {
			fmt.Printf("[MOCK] Response from %s (status: %d)\n", resp.Request.URL.Path, resp.StatusCode)
			return nil
		},
	}
}

func (m *MockAdapter) Shutdown(ctx context.Context) error {
	fmt.Printf("[MOCK] Adapter %s shutting down\n", m.name)
	return nil
}

func ExampleClient_withAdapter() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello from Adapter!"}`))
	}))
	defer server.Close()

	// Create a mock adapter
	adapter := &MockAdapter{name: "test-adapter"}

	// Create client with adapter
	client := neuron.NewClient(neuron.ClientOptions{
		BaseURL: server.URL,
		Adapter: adapter,
	})
	defer client.Shutdown(5 * time.Second)

	// Define a route
	route := neuron.NewRoute[neuron.EmptyRequest, map[string]string](
		neuron.MethodGET, "/test",
	)

	// Execute request
	resp, err := neuron.Execute(client, route, neuron.EmptyRequest{}, &neuron.RequestOptions{
		Context: context.Background(),
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Data: %+v\n", resp.Data)
	// Output:
	// [MOCK] Request to /test
	// [MOCK] Response from /test (status: 200)
	// Status: 200
	// Data: map[message:Hello from Adapter!]
}

func ExampleClient_withIonAdapter() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello from Ion Adapter!"}`))
	}))
	defer server.Close()

	// Create ion adapter (this would be imported from the adapter module)
	// ionAdapter := ion.NewAdapter().
	// 	WithRateLimiter(ion.NewRateLimiter(10, 5)).
	// 	WithCircuitBreaker(ion.NewCircuitBreaker("test-service"))

	// For now, use mock adapter to demonstrate the pattern
	adapter := &MockAdapter{name: "ion-adapter"}

	// Create client with adapter
	client := neuron.NewClient(neuron.ClientOptions{
		BaseURL: server.URL,
		Adapter: adapter,
	})
	defer client.Shutdown(5 * time.Second)

	// Define a route
	route := neuron.NewRoute[neuron.EmptyRequest, map[string]string](
		neuron.MethodGET, "/test",
	)

	// Execute request
	resp, err := neuron.Execute(client, route, neuron.EmptyRequest{}, &neuron.RequestOptions{
		Context: context.Background(),
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Data: %+v\n", resp.Data)
	// Output:
	// [MOCK] Request to /test
	// [MOCK] Response from /test (status: 200)
	// Status: 200
	// Data: map[message:Hello from Ion Adapter!]
}
