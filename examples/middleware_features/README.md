# Neuron Middleware Features Example

This example demonstrates the built-in middleware features available in Neuron.

## Features Demonstrated

- **Request ID Generation** - Automatic unique request ID generation and propagation
- **Logging** - Request/response logging with configurable levels
- **Authentication** - Bearer token authentication
- **Headers** - User agent and custom header management
- **Timeout** - Request timeout configuration
- **Compression** - Automatic compression support

## Usage

```bash
cd examples/middleware_features
go run main.go
```

## Middleware Used

### Request Middleware

```go
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
    neuron.TimeoutMiddleware(neuron.TimeoutConfig{Timeout: 30 * time.Second}),

    // Compression
    neuron.AutoCompressionMiddleware(),
}
```

### Response Middleware

```go
ResponseMiddleware: []neuron.ResponseMiddleware{
    // Response logging
    neuron.LoggingResponseMiddleware(neuron.DefaultLoggingConfig()),
}
```

## Available Middleware

### Logging

- `LoggingMiddleware()` - Basic request logging
- `LoggingResponseMiddleware()` - Response logging
- `DebugLoggingMiddleware()` - Debug-level logging
- `ErrorLoggingMiddleware()` - Error-only logging
- `StructuredLoggingMiddleware()` - JSON structured logging

### Metrics

- `MetricsMiddleware()` - Request metrics collection
- `MetricsResponseMiddleware()` - Response metrics collection
- `AutoMetricsMiddleware()` - Auto-create metrics collector

### Request ID

- `RequestIDMiddleware()` - Request ID generation
- `TracingMiddleware()` - Request tracing
- `CorrelationIDMiddleware()` - Distributed tracing

### Authentication

- `BearerAuthMiddleware()` - Bearer token authentication
- `APIKeyAuthMiddleware()` - API key authentication
- `BasicAuthMiddleware()` - Basic authentication
- `CustomAuthMiddleware()` - Custom authentication

### Compression

- `CompressionMiddleware()` - Request compression
- `AutoCompressionMiddleware()` - Auto accept-encoding

### Timeout

- `TimeoutMiddleware()` - Request timeout
- `PerRequestTimeoutMiddleware()` - Per-request timeout
- `GlobalTimeoutMiddleware()` - Global timeout

### Headers

- `UserAgentMiddleware()` - User agent header
- `CustomHeaderMiddleware()` - Custom headers
- `ContentTypeMiddleware()` - Content-Type header
- `AcceptMiddleware()` - Accept header

## Expected Output

```
[REQUEST] GET http://localhost:xxxxx/test
[RESPONSE] GET http://localhost:xxxxx/test - 200 (5ms)

--- Request 1 ---
Status: 200
Data: map[message:Hello from Neuron with Middleware Features! request_id:xxx-xxx-xxx]

--- Request 2 ---
Status: 200
Data: map[message:Hello from Neuron with Middleware Features! request_id:xxx-xxx-xxx]

--- Request 3 ---
Status: 200
Data: map[message:Hello from Neuron with Middleware Features! request_id:xxx-xxx-xxx]

--- Metrics ---
Request count: 3
Average duration: ~100ms
Success rate: 100%
```

## Customization

You can customize each middleware by providing configuration:

```go
// Custom logging configuration
neuron.LoggingMiddleware(neuron.LoggingConfig{
    Level:       neuron.LogLevelDebug,
    IncludeBody: true,
    MaxBodySize: 2048,
    Logger:      log.Default(),
})

// Custom timeout configuration
neuron.TimeoutMiddleware(neuron.TimeoutConfig{
    Timeout:        30 * time.Second,
    PerRequest:     true,
    GlobalTimeout:  60 * time.Second,
    RequestTimeout: 30 * time.Second,
})

// Custom request ID configuration
neuron.RequestIDMiddleware(neuron.RequestIDConfig{
    Generator:  &neuron.UUIDGenerator{},
    HeaderName: "X-Request-ID",
    Propagate:  true,
})
```

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
