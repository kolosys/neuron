# Neuron Adapter System

The Neuron Adapter System provides a modular way to integrate external libraries with the Neuron HTTP client. This system allows you to add advanced features like rate limiting, circuit breaking, and concurrency control without coupling them directly to the core Neuron library.

## Overview

The adapter system is designed with the following principles:

- **Zero Dependencies** - The core adapter interface has no external dependencies
- **Modular Design** - Each adapter is a separate, independent module
- **Type Safety** - Compile-time type checking for all adapters
- **Clean Architecture** - Clear separation of concerns
- **No Import Cycles** - Clean dependency structure

## Architecture

```
adapter/
├── interface.go            # Core adapter interface
├── go.mod                   # Adapter module
├── ion/                     # Ion-specific adapter
│   ├── adapter.go          # Ion adapter implementation
│   ├── go.mod              # Ion adapter module
│   └── go.sum              # Ion adapter dependencies
└── README.md               # This file
```

## Core Interface

The `Adapter` interface defines the contract that all adapters must implement:

```go
type Adapter interface {
    // Name returns the adapter name
    Name() string

    // WrapHTTPClient wraps an HTTP client with adapter functionality
    WrapHTTPClient(client *http.Client) *http.Client

    // CreateRequestMiddleware creates request middleware for the adapter
    CreateRequestMiddleware() []RequestMiddleware

    // CreateResponseMiddleware creates response middleware for the adapter
    CreateResponseMiddleware() []ResponseMiddleware

    // Shutdown gracefully shuts down the adapter
    Shutdown(ctx context.Context) error
}
```

## Available Adapters

### Ion Adapter

The Ion adapter provides integration with the `ion` library for advanced concurrency and resilience features:

- **Rate Limiting** - Token bucket and sliding window rate limiters
- **Circuit Breaking** - Circuit breaker pattern implementation
- **Worker Pool Integration** - Worker pool concurrency control (coming soon)
- **Semaphore Integration** - Weighted fair semaphores (coming soon)

#### Usage

```go
import (
    "github.com/kolosys/neuron"
    "github.com/kolosys/neuron/adapter/ion"
)

// Create ion adapter
ionAdapter := ion.NewAdapter().
    WithRateLimiter(ion.NewRateLimiter(10, 5)).
    WithCircuitBreaker(ion.NewCircuitBreaker("my-service"))

// Use with neuron client
client := neuron.NewClient(neuron.ClientOptions{
    BaseURL: "https://api.example.com",
    Adapter: ionAdapter,
})
```

## Creating Custom Adapters

To create a custom adapter, implement the `Adapter` interface:

```go
package myadapter

import (
    "context"
    "net/http"
    "github.com/kolosys/neuron/adapter"
)

type MyAdapter struct {
    *adapter.BaseAdapter
    // Your custom fields
}

func NewAdapter() *MyAdapter {
    return &MyAdapter{
        BaseAdapter: adapter.NewBaseAdapter("my-adapter"),
    }
}

func (a *MyAdapter) WrapHTTPClient(client *http.Client) *http.Client {
    // Add custom HTTP client modifications
    return client
}

func (a *MyAdapter) CreateRequestMiddleware() []adapter.RequestMiddleware {
    return []adapter.RequestMiddleware{
        func(req *http.Request) error {
            // Your custom request processing
            return nil
        },
    }
}

func (a *MyAdapter) CreateResponseMiddleware() []adapter.ResponseMiddleware {
    return []adapter.ResponseMiddleware{
        func(resp *http.Response) error {
            // Your custom response processing
            return nil
        },
    }
}

func (a *MyAdapter) Shutdown(ctx context.Context) error {
    // Your custom shutdown logic
    return nil
}
```

## Configuration Types

The adapter system provides configuration types for common patterns:

### RateLimitConfig

```go
type RateLimitConfig struct {
    RequestsPerSecond int
    BurstSize         int
    QueueOnRateLimit  bool
    Timeout           int // seconds
}
```

### CircuitBreakerConfig

```go
type CircuitBreakerConfig struct {
    FailureThreshold        int
    RecoveryTimeout         int // seconds
    HalfOpenMaxRequests     int
    SuccessThreshold        int
    FailurePredicate        func(error) bool
    PerRouteCircuitBreakers bool
}
```

## Benefits

- **Modularity** - Each adapter is a separate module
- **Extensibility** - Easy to add new adapters
- **Type Safety** - Compile-time type checking
- **Clean Architecture** - Clear separation of concerns
- **No Import Cycles** - Clean dependency structure
- **Zero Dependencies** - Core interface has no external dependencies

## Contributing

To contribute a new adapter:

1. Create a new directory under `adapter/`
2. Implement the `Adapter` interface
3. Add appropriate tests
4. Update this README with usage examples
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
