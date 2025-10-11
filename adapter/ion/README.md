# Neuron Ion Adapter

The Neuron Ion Adapter provides integration with the `ion` library for advanced concurrency and resilience features in the Neuron HTTP client.

## Overview

The Ion adapter brings the power of the `ion` library to Neuron, providing:

- **Rate Limiting** - Token bucket and sliding window rate limiters
- **Circuit Breaking** - Circuit breaker pattern implementation
- **Worker Pool Integration** - Worker pool concurrency control (coming soon)
- **Semaphore Integration** - Weighted fair semaphores (coming soon)

## Installation

```bash
go get github.com/kolosys/neuron/adapter/ion
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/kolosys/neuron"
    "github.com/kolosys/neuron/adapter/ion"
)

func main() {
    // Create ion adapter
    ionAdapter := ion.NewAdapter().
        WithRateLimiter(ion.NewRateLimiter(10, 5)).
        WithCircuitBreaker(ion.NewCircuitBreaker("my-service"))

    // Create client with ion adapter
    client := neuron.NewClient(neuron.ClientOptions{
        BaseURL: "https://api.example.com",
        Adapter: ionAdapter,
    })
    defer client.Shutdown(5 * time.Second)

    // Define a route
    route := neuron.NewRoute[neuron.EmptyRequest, map[string]string](
        neuron.MethodGET, "/users",
    )

    // Execute request with ion rate limiting and circuit breaking
    resp, err := neuron.Execute(client, route, neuron.EmptyRequest{}, &neuron.RequestOptions{
        Context: context.Background(),
    })

    if err != nil {
        panic(err)
    }

    fmt.Printf("Response: %+v\n", resp.Data)
}
```

## Features

### Rate Limiting

The Ion adapter supports rate limiting using ion's rate limiter implementations:

```go
// Create rate limiter with 10 requests per second and burst of 5
rateLimiter := ion.NewRateLimiter(10, 5)

// Create adapter with rate limiter
adapter := ion.NewAdapter().
    WithRateLimiter(rateLimiter)
```

#### Rate Limiter Options

- **RequestsPerSecond** - Maximum requests per second
- **BurstSize** - Maximum burst size for token bucket
- **QueueOnRateLimit** - Whether to queue requests when rate limited
- **Timeout** - Timeout for rate limit operations

### Circuit Breaking

The Ion adapter supports circuit breaking using ion's circuit breaker implementation:

```go
// Create circuit breaker with custom options
circuitBreaker := ion.NewCircuitBreakerWithOptions("my-service",
    circuit.WithFailureThreshold(5),
    circuit.WithRecoveryTimeout(30*time.Second),
    circuit.WithHalfOpenMaxRequests(3),
    circuit.WithHalfOpenSuccessThreshold(2),
)

// Create adapter with circuit breaker
adapter := ion.NewAdapter().
    WithCircuitBreaker(circuitBreaker)
```

#### Circuit Breaker Options

- **FailureThreshold** - Number of failures before opening circuit
- **RecoveryTimeout** - Time to wait before attempting recovery
- **HalfOpenMaxRequests** - Maximum requests in half-open state
- **SuccessThreshold** - Successes needed to close circuit
- **FailurePredicate** - Custom function to determine failures

### Configuration

The Ion adapter supports configuration through the adapter interface:

```go
// Configure rate limiting
adapter.ConfigureRateLimiting(adapter.RateLimitConfig{
    RequestsPerSecond: 20,
    BurstSize:         10,
    QueueOnRateLimit:  true,
    Timeout:           30,
})

// Configure circuit breaker
adapter.ConfigureCircuitBreaker(adapter.CircuitBreakerConfig{
    FailureThreshold:       3,
    RecoveryTimeout:        15,
    HalfOpenMaxRequests:    2,
    SuccessThreshold:       1,
    FailurePredicate:       func(err error) bool { return err != nil },
    PerRouteCircuitBreakers: false,
})
```

## Advanced Usage

### Custom Rate Limiter

```go
import (
    "github.com/kolosys/ion/ratelimit"
)

// Create custom rate limiter
customRateLimiter := ratelimit.NewTokenBucket(
    ratelimit.PerSecond(20), // 20 requests per second
    10,                      // burst size of 10
)

adapter := ion.NewAdapter().
    WithRateLimiter(customRateLimiter)
```

### Custom Circuit Breaker

```go
import (
    "github.com/kolosys/ion/circuit"
)

// Create custom circuit breaker
customCircuitBreaker := circuit.New("my-service",
    circuit.WithFailureThreshold(3),
    circuit.WithRecoveryTimeout(15*time.Second),
    circuit.WithHalfOpenMaxRequests(2),
    circuit.WithHalfOpenSuccessThreshold(1),
)

adapter := ion.NewAdapter().
    WithCircuitBreaker(customCircuitBreaker)
```

### Middleware Integration

The Ion adapter automatically creates middleware for:

- **Rate Limiting** - Request middleware that enforces rate limits
- **Circuit Breaking** - Request/response middleware for circuit breaker state
- **Logging** - Response middleware for circuit breaker state logging

## Examples

See the `examples/` directory for complete working examples:

- **Basic Usage** - Simple rate limiting and circuit breaking
- **Custom Configuration** - Advanced configuration options
- **Middleware Integration** - Custom middleware examples

## Dependencies

The Ion adapter depends on:

- `github.com/kolosys/ion` - The ion library for concurrency primitives
- `github.com/kolosys/neuron` - The core Neuron library
- `github.com/kolosys/neuron/adapter` - The adapter interface

## Contributing

To contribute to the Ion adapter:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
