# Quick Start

This guide will help you get started with neuron quickly with a basic example.

## Basic Usage

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/neuron"
)

func main() {
    // Basic usage example
    fmt.Println("Welcome to neuron!")
    
    // TODO: Add your code here
}
```

## Common Use Cases

### Using neuron

**Import Path:** `github.com/kolosys/neuron`



```go
package main

import (
    "fmt"
    "github.com/kolosys/neuron"
)

func main() {
    // Example usage of neuron
    fmt.Println("Using neuron package")
}
```

#### Available Types
- **AuthConfig** - AuthConfig configures authentication middleware
- **AuthProvider** - AuthProvider interface for authentication
- **AuthType** - AuthType represents the type of authentication
- **BackoffStrategy** - BackoffStrategy defines how to handle backoff
- **BodyProvider** - BodyProvider allows custom body serialization
- **Cache** - Cache interface for caching middleware
- **CacheEntry** - CacheEntry represents a cached response
- **CircuitBreaker** - CircuitBreaker interface for middleware integration
- **CircuitBreakerConfig** - CircuitBreakerConfig defines circuit breaker behavior for HTTP clients
- **CircuitBreakerState** - CircuitBreakerState represents the state of a circuit breaker
- **Client** - Client provides a type-safe HTTP client with rate limiting, queuing, and circuit breaking
- **ClientError** - Error types for type-safe error handling
- **ClientOptions** - ClientOptions configures the HTTP client behavior
- **CompressionConfig** - CompressionConfig configures compression middleware
- **CompressionType** - CompressionType represents the type of compression
- **Deserializable** - Deserializable represents types that can be deserialized from responses
- **EmptyRequest** - EmptyRequest represents requests with no body
- **EmptyResponse** - EmptyResponse represents responses with no body
- **ErrorType** - 
- **HTTPMethod** - HTTPMethod represents supported HTTP methods
- **HeaderConfig** - HeaderConfig configures header middleware
- **InMemoryCache** - InMemoryCache provides a simple in-memory cache
- **JSONValidator** - JSONValidator provides JSON schema validation
- **LogEntry** - LogEntry represents a log entry
- **LogLevel** - LogLevel represents the logging level
- **Logger** - Interfaces for extensibility Logger interface for logging middleware
- **LoggingConfig** - LoggingConfig configures the logging middleware
- **LoggingContextKey** - 
- **MetricsCollector** - MetricsCollector collects and stores metrics
- **MetricsSnapshot** - MetricsSnapshot represents a snapshot of metrics at a point in time
- **MiddlewareChain** - MiddlewareChain manages a chain of middleware functions
- **QueuedRequest** - QueuedRequest represents a request waiting in queue
- **QueuedResponse** - QueuedResponse represents the result of a queued request
- **RateLimitConfig** - RateLimitConfig defines rate limiting behavior
- **RateLimitInfo** - RateLimitInfo contains information about current rate limit status
- **RateLimitInfoProvider** - RateLimitInfoProvider provides rate limit information
- **RateLimiter** - Middleware interfaces for external library integration RateLimiter interface for middleware integration
- **RequestContext** - RequestContext provides metadata about the current request
- **RequestIDConfig** - RequestIDConfig configures request ID generation
- **RequestIDContextKey** - 
- **RequestIDGenerator** - RequestIDGenerator generates unique request IDs
- **RequestMetrics** - RequestMetrics provides insights into request performance
- **RequestMiddleware** - RequestMiddleware processes requests before they are sent
- **RequestOptions** - RequestOptions contains configuration for individual requests
- **RequestQueue** - RequestQueue manages queued requests for a specific route/bucket
- **Response** - Response wraps HTTP response data with type safety
- **ResponseMiddleware** - ResponseMiddleware processes responses after they are received
- **RetryCondition** - RetryCondition determines if a request should be retried
- **Route** - Route represents a type-safe route definition
- **SequentialGenerator** - SequentialGenerator generates sequential request IDs
- **Serializable** - Serializable represents types that can be serialized for requests
- **SimpleLogger** - Default implementations SimpleLogger provides a basic logger implementation
- **StaticAuthProvider** - StaticAuthProvider provides static token authentication
- **TimeoutConfig** - TimeoutConfig configures timeout middleware
- **TimestampGenerator** - TimestampGenerator generates timestamp-based request IDs
- **UUIDGenerator** - UUIDGenerator generates UUID-style request IDs
- **Validator** - Validator interface for request validation

#### Available Functions
- **AddAutoMetrics** - AddAutoMetrics creates a simple metrics middleware that doesn't require a collector
- **GetRequestID** - GetRequestID extracts the request ID from context
- **RequestIDFromRequest** - RequestIDFromRequest extracts the request ID from an HTTP request
- **RequestIDFromResponse** - RequestIDFromResponse extracts the request ID from an HTTP response
- **WithRequestID** - WithRequestID adds a request ID to the context

For detailed API documentation, see the [neuron API Reference](../api-reference/neuron.md).

## Step-by-Step Tutorial

### Step 1: Import the Package

First, import the necessary packages in your Go file:

```go
import (
    "fmt"
    "github.com/kolosys/neuron"
)
```

### Step 2: Initialize

Set up the basic configuration:

```go
func main() {
    // Initialize your application
    fmt.Println("Initializing neuron...")
}
```

### Step 3: Use the Library

Implement your specific use case:

```go
func main() {
    // Your implementation here
}
```

## Running Your Code

To run your Go program:

```bash
go run main.go
```

To build an executable:

```bash
go build -o myapp
./myapp
```

## Configuration Options

neuron can be configured to suit your needs. Check the [Core Concepts](../core-concepts/) section for detailed information about configuration options.

## Error Handling

Always handle errors appropriately:

```go
result, err := someFunction()
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

## Best Practices

- Always handle errors returned by library functions
- Check the API documentation for detailed parameter information
- Use meaningful variable and function names
- Add comments to document your code

## Complete Example

Here's a complete working example:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/neuron"
)

func main() {
    fmt.Println("Starting neuron application...")
    
    // Add your implementation here
    
    fmt.Println("Application completed successfully!")
}
```

## Next Steps

Now that you've seen the basics, explore:

- **[Core Concepts](../core-concepts/)** - Understanding the library architecture
- **[API Reference](../api-reference/)** - Complete API documentation
- **[Examples](../examples/README.md)** - More detailed examples
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced patterns

## Getting Help

If you run into issues:

1. Check the [API Reference](../api-reference/)
2. Browse the [Examples](../examples/README.md)
3. Visit the [GitHub Issues](https://github.com/kolosys/neuron/issues) page

