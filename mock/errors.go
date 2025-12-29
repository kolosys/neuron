package mock

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
)

// CommonErrors provides sentinel errors for common test scenarios.
type CommonErrors struct {
	Timeout         error
	NetworkError    error
	DNSError        error
	TLSError        error
	ContextCanceled error
	RateLimit       error
	AuthFailed      error
}

// DefaultErrors returns a CommonErrors struct with standard test errors.
func DefaultErrors() *CommonErrors {
	return &CommonErrors{
		Timeout:         errors.New("i/o timeout"),
		NetworkError:    errors.New("network connection refused"),
		DNSError:        errors.New("no such host"),
		TLSError:        errors.New("tls: handshake failure"),
		ContextCanceled: errors.New("context canceled"),
		RateLimit:       errors.New("rate limit exceeded"),
		AuthFailed:      errors.New("authentication failed"),
	}
}

// ErrorSequence provides a sequence of errors for testing retry behavior.
// Each call to Next() returns the next error in the sequence.
type ErrorSequence struct {
	errors []error
	index  atomic.Int32
}

// NewErrorSequence creates a new error sequence from the provided errors.
func NewErrorSequence(errs ...error) *ErrorSequence {
	return &ErrorSequence{
		errors: errs,
	}
}

// Next returns the next error in the sequence.
// If the end is reached, it returns the last error.
func (es *ErrorSequence) Next() error {
	if len(es.errors) == 0 {
		return nil
	}

	idx := es.index.Load()
	if int(idx) >= len(es.errors) {
		return es.errors[len(es.errors)-1]
	}

	err := es.errors[idx]
	es.index.Store(idx + 1)
	return err
}

// Error implements the error interface.
func (es *ErrorSequence) Error() string {
	if len(es.errors) == 0 {
		return "error sequence is empty"
	}

	idx := es.index.Load()
	if int(idx) >= len(es.errors) {
		return es.errors[len(es.errors)-1].Error()
	}

	return es.errors[idx].Error()
}

// Reset resets the sequence to the first error.
func (es *ErrorSequence) Reset() {
	es.index.Store(0)
}

// errHolder is a wrapper for nil-safe atomic storage of errors.
type errHolder struct {
	err     error
	oneShot bool
}

// consumeError retrieves and optionally consumes an error from atomic storage.
// If oneShot is true, the error is cleared after retrieval.
func consumeError(val any) error {
	if val == nil {
		return nil
	}

	holder := val.(*errHolder)
	if holder == nil || holder.err == nil {
		return nil
	}

	return holder.err
}

// IsTimeoutError checks if an error is a timeout error.
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return false
}

// IsNetworkError checks if an error is a network-related error.
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	return errors.As(err, &netErr)
}

// IsContextCanceledError checks if an error is a context cancellation error.
func IsContextCanceledError(err error) bool {
	return errors.Is(err, context.Canceled)
}

// IsRateLimitError checks if an error message contains "rate limit".
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, context.Canceled) || contains(err.Error(), "rate limit")
}

// contains checks if a string contains a substring (helper for error checking).
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
