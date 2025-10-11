package neuron

import (
	"context"
	"net/http"
)

// SimpleAdapter is a basic adapter interface for external integrations
type SimpleAdapter interface {
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

// AdapterBridge converts between external adapter types and neuron types
type AdapterBridge struct {
	adapter interface {
		Name() string
		WrapHTTPClient(client *http.Client) *http.Client
		CreateRequestMiddleware() []RequestMiddleware
		CreateResponseMiddleware() []ResponseMiddleware
		Shutdown(ctx context.Context) error
	}
}

// NewAdapterBridge creates a bridge between external adapter and neuron
func NewAdapterBridge(adapter SimpleAdapter) *AdapterBridge {
	return &AdapterBridge{adapter: adapter}
}

// Name returns the adapter name
func (b *AdapterBridge) Name() string {
	return b.adapter.Name()
}

// WrapHTTPClient delegates to the adapter
func (b *AdapterBridge) WrapHTTPClient(client *http.Client) *http.Client {
	return b.adapter.WrapHTTPClient(client)
}

// CreateRequestMiddleware delegates to the adapter
func (b *AdapterBridge) CreateRequestMiddleware() []RequestMiddleware {
	return b.adapter.CreateRequestMiddleware()
}

// CreateResponseMiddleware delegates to the adapter
func (b *AdapterBridge) CreateResponseMiddleware() []ResponseMiddleware {
	return b.adapter.CreateResponseMiddleware()
}

// Shutdown delegates to the adapter
func (b *AdapterBridge) Shutdown(ctx context.Context) error {
	return b.adapter.Shutdown(ctx)
}
