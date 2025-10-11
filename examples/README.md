# Neuron Examples

This directory contains working examples demonstrating how to use Neuron with various adapters and configurations.

## Examples

### Middleware Features Example

**Location:** `middleware_features/`

Demonstrates the built-in middleware features available in Neuron, including logging, metrics, authentication, compression, and more.

**Features:**

- Request ID generation and propagation
- Request/response logging
- Authentication (Bearer, API Key, Basic)
- Header management
- Timeout configuration
- Compression support

**Usage:**

```bash
cd examples/middleware_features
go run main.go
```

### Ion Adapter Example

**Location:** `ion_adapter/`

Demonstrates how to use the Neuron HTTP client with the Ion adapter for rate limiting and circuit breaking.

**Features:**

- Rate limiting with token bucket algorithm
- Circuit breaker pattern implementation
- Custom middleware integration
- Error handling and recovery

**Usage:**

```bash
cd examples/ion_adapter
go run main.go
```

**Code Example:**

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

    // Create client with adapter
    client := neuron.NewClient(neuron.ClientOptions{
        BaseURL: "https://api.example.com",
        Adapter: ionAdapter,
    })
    defer client.Shutdown(5 * time.Second)

    // Use the client...
}
```

## Running Examples

### Prerequisites

- Go 1.21 or later
- Neuron library installed
- Required adapter modules

### Installation

```bash
# Install neuron
go get github.com/kolosys/neuron

# Install ion adapter
go get github.com/kolosys/neuron/adapter/ion

# Install ion library
go get github.com/kolosys/ion
```

### Running Examples

```bash
# Run middleware features example
cd examples/middleware_features
go run main.go

# Run ion adapter example
cd examples/ion_adapter
go run main.go
```

## Creating Your Own Examples

To create a new example:

1. Create a new directory under `examples/`
2. Add a `main.go` file with your example
3. Add a `go.mod` file with dependencies
4. Add a `README.md` with documentation
5. Test your example thoroughly

### Example Structure

```
examples/
├── README.md                   # This file
├── middleware_features/        # Middleware features example
│   ├── main.go                 # Example code
│   ├── go.mod                  # Dependencies
│   └── README.md               # Example documentation
├── ion_adapter/                # Ion adapter example
│   ├── main.go                 # Example code
│   ├── go.mod                  # Dependencies
│   └── README.md               # Example documentation
└── custom_adapter/             # Custom adapter example (coming soon)
    ├── main.go
    ├── go.mod
    └── README.md
```

### Example Template

```go
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
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Hello from Neuron!"}`))
    }))
    defer server.Close()

    // Create client
    client := neuron.NewClient(neuron.ClientOptions{
        BaseURL: server.URL,
        // Add your adapter here
    })
    defer client.Shutdown(5 * time.Second)

    // Define route
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
}
```

## Contributing Examples

To contribute a new example:

1. Create a new directory under `examples/`
2. Add your example code
3. Add proper documentation
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
