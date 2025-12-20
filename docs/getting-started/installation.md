# Installation

This guide will help you install and set up neuron in your Go project.

## Prerequisites

Before installing neuron, ensure you have:

- **Go ** or later installed
- A Go module initialized in your project (run `go mod init` if needed)
- Access to the GitHub repository (for private repositories)

## Installation Steps

### Step 1: Install the Package

Use `go get` to install neuron:

```bash
go get github.com/kolosys/neuron
```

This will download the package and add it to your `go.mod` file.

### Step 2: Import in Your Code

Import the package in your Go source files:

```go
import "github.com/kolosys/neuron"
```

### Multiple Packages

neuron includes several packages. Import the ones you need:

```go
// 
import "github.com/kolosys/neuron"
```

### Step 3: Verify Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "github.com/kolosys/neuron"
)

func main() {
    fmt.Println("neuron installed successfully!")
}
```

Run the test:

```bash
go run main.go
```

## Updating the Package

To update to the latest version:

```bash
go get -u github.com/kolosys/neuron
```

To update to a specific version:

```bash
go get github.com/kolosys/neuron@v1.2.3
```

## Installing a Specific Version

To install a specific version of the package:

```bash
go get github.com/kolosys/neuron@v1.0.0
```

Check available versions on the [GitHub releases page](https://github.com/kolosys/neuron/releases).

## Development Setup

If you want to contribute or modify the library:

### Clone the Repository

```bash
git clone https://github.com/kolosys/neuron.git
cd neuron
```

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
go test ./...
```

## Troubleshooting

### Module Not Found

If you encounter a "module not found" error:

1. Ensure your `GOPATH` is set correctly
2. Check that you have network access to GitHub
3. Try running `go clean -modcache` and reinstall

### Private Repository Access

For private repositories, configure Git to use SSH or a personal access token:

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Or set up GOPRIVATE:

```bash
export GOPRIVATE=github.com/kolosys/neuron
```

## Next Steps

Now that you have neuron installed, check out the [Quick Start Guide](quick-start.md) to learn how to use it.

## Additional Resources

- [Quick Start Guide](quick-start.md)
- [API Reference](../api-reference/)
- [Examples](../examples/README.md)

