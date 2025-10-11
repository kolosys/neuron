# Neuron Ion Adapter Example

This example demonstrates how to use the Neuron HTTP client with the Ion adapter for rate limiting and circuit breaking.

## Overview

This example shows:

- How to create and configure an Ion adapter
- How to integrate rate limiting with token bucket algorithm
- How to implement circuit breaker pattern
- How to handle errors and recovery
- How to use custom middleware

## Features Demonstrated

- **Rate Limiting** - Token bucket rate limiter with configurable rate and burst
- **Circuit Breaking** - Circuit breaker with configurable thresholds
- **Middleware Integration** - Custom request/response middleware
- **Error Handling** - Proper error handling and recovery
- **Shutdown** - Graceful shutdown of the client and adapter

## Code Structure

```go
// MockIonAdapter implements the SimpleAdapter interface for testing
type MockIonAdapter struct {
    name string
}

// Methods implement the SimpleAdapter interface:
// - Name() string
// - WrapHTTPClient(client *http.Client) *http.Client
// - CreateRequestMiddleware() []neuron.RequestMiddleware
// - CreateResponseMiddleware() []neuron.ResponseMiddleware
// - Shutdown(ctx context.Context) error
```

## Usage

### Prerequisites

- Go 1.21 or later
- Neuron library
- Ion adapter module

### Installation

```bash
# Install dependencies
go mod tidy
```

### Running the Example

```bash
# Run the example
go run main.go
```

### Expected Output

```
[ION] Rate limiting request to /test
[ION] Circuit breaker response for /test (status: 200)
Status: 200
Data: map[message:Hello from Modular Ion Adapter!]
[ION] Adapter ion-adapter shutting down
```

## Configuration

The example uses a mock adapter that simulates the behavior of the real Ion adapter. In a real implementation, you would use:

```go
import (
    "github.com/kolosys/neuron/adapter/ion"
)

// Create real ion adapter
ionAdapter := ion.NewAdapter().
    WithRateLimiter(ion.NewRateLimiter(10, 5)).
    WithCircuitBreaker(ion.NewCircuitBreaker("my-service"))
```

## Customization

### Rate Limiting

```go
// Custom rate limiter configuration
rateLimiter := ion.NewRateLimiter(20, 10) // 20 req/s, burst 10
adapter := ion.NewAdapter().WithRateLimiter(rateLimiter)
```

### Circuit Breaking

```go
// Custom circuit breaker configuration
circuitBreaker := ion.NewCircuitBreakerWithOptions("my-service",
    circuit.WithFailureThreshold(3),
    circuit.WithRecoveryTimeout(15*time.Second),
)
adapter := ion.NewAdapter().WithCircuitBreaker(circuitBreaker)
```

### Custom Middleware

```go
// Add custom request middleware
func (m *MockIonAdapter) CreateRequestMiddleware() []neuron.RequestMiddleware {
    return []neuron.RequestMiddleware{
        func(req *http.Request) error {
            // Your custom logic here
            fmt.Printf("[CUSTOM] Processing request to %s\n", req.URL.Path)
            return nil
        },
    }
}
```

## Error Handling

The example includes proper error handling:

```go
resp, err := neuron.Execute(client, route, neuron.EmptyRequest{}, &neuron.RequestOptions{
    Context: context.Background(),
})

if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
```

## Shutdown

The example demonstrates graceful shutdown:

```go
defer client.Shutdown(5 * time.Second)
```

This ensures that:

- All pending requests are completed
- The adapter is properly shut down
- Resources are cleaned up

## Real-World Usage

In a real application, you would:

1. Replace the mock adapter with the real Ion adapter
2. Configure appropriate rate limits for your use case
3. Set up circuit breaker thresholds based on your error rates
4. Add proper logging and monitoring
5. Handle errors appropriately

## Troubleshooting

### Common Issues

1. **Import Errors** - Make sure all dependencies are installed
2. **Type Errors** - Ensure you're using the correct types from the adapter
3. **Shutdown Issues** - Make sure to call `client.Shutdown()` properly

### Debug Mode

To enable debug logging, you can modify the middleware:

```go
func (m *MockIonAdapter) CreateRequestMiddleware() []neuron.RequestMiddleware {
    return []neuron.RequestMiddleware{
        func(req *http.Request) error {
            fmt.Printf("[DEBUG] Request: %s %s\n", req.Method, req.URL.Path)
            return nil
        },
    }
}
```

## Contributing

To contribute to this example:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
