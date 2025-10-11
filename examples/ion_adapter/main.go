package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kolosys/neuron"
)

// MockIonAdapter implements the SimpleAdapter interface for testing
type MockIonAdapter struct {
	name string
}

func (m *MockIonAdapter) Name() string {
	return m.name
}

func (m *MockIonAdapter) WrapHTTPClient(client *http.Client) *http.Client {
	return client
}

func (m *MockIonAdapter) CreateRequestMiddleware() []neuron.RequestMiddleware {
	return []neuron.RequestMiddleware{
		func(req *http.Request) error {
			fmt.Printf("[ION] Rate limiting request to %s\n", req.URL.Path)
			return nil
		},
	}
}

func (m *MockIonAdapter) CreateResponseMiddleware() []neuron.ResponseMiddleware {
	return []neuron.ResponseMiddleware{
		func(resp *http.Response) error {
			fmt.Printf("[ION] Circuit breaker response for %s (status: %d)\n",
				resp.Request.URL.Path, resp.StatusCode)
			return nil
		},
	}
}

func (m *MockIonAdapter) Shutdown(ctx context.Context) error {
	fmt.Printf("[ION] Adapter %s shutting down\n", m.name)
	return nil
}

func main() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello from Modular Ion Adapter!"}`))
	}))
	defer server.Close()

	// Create mock ion adapter (in real usage, this would be the actual ion adapter)
	ionAdapter := &MockIonAdapter{name: "ion-adapter"}

	// Create client with ion adapter
	client := neuron.NewClient(neuron.ClientOptions{
		BaseURL: server.URL,
		Adapter: ionAdapter,
	})
	defer client.Shutdown(5 * time.Second)

	// Define a route
	route := neuron.NewRoute[neuron.EmptyRequest, map[string]string](
		neuron.MethodGET, "/test",
	)

	// Execute request with ion rate limiting and circuit breaking
	resp, err := neuron.Execute(client, route, neuron.EmptyRequest{}, &neuron.RequestOptions{
		Context: context.Background(),
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Data: %+v\n", resp.Data)
}
