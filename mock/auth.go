package mock

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// MockAuthProvider is a mock implementation of neuron.AuthProvider for testing.
type MockAuthProvider struct {
	tokens       []string
	tokenIndex   atomic.Int32
	getCalls     []GetTokenCall
	headerCalls  []GetHeaderCall
	tokenErr     atomic.Value // *errHolder
	headerFormat string
	recordCalls  atomic.Bool
	mu           sync.RWMutex
}

// GetTokenCall records a call to GetToken.
type GetTokenCall struct {
	Time   time.Time
	Result string
	Err    error
}

// GetHeaderCall records a call to GetAuthHeader.
type GetHeaderCall struct {
	Token  string
	Result string
	Time   time.Time
}

// MockAuthProviderOptions configures mock auth provider behavior.
type MockAuthProviderOptions struct {
	InitialToken string
	RecordCalls  bool
	Tokens       []string
	HeaderFormat string // e.g., "Bearer {}" or "X-API-Key {}"
}

// NewMockAuthProvider creates a new mock auth provider.
func NewMockAuthProvider(opts *MockAuthProviderOptions) *MockAuthProvider {
	ap := &MockAuthProvider{
		tokens:       make([]string, 0),
		getCalls:     make([]GetTokenCall, 0),
		headerCalls:  make([]GetHeaderCall, 0),
		headerFormat: "Bearer {}",
	}

	if opts == nil {
		ap.recordCalls.Store(true)
		return ap
	}

	ap.recordCalls.Store(opts.RecordCalls)

	if opts.InitialToken != "" {
		ap.tokens = append(ap.tokens, opts.InitialToken)
	}

	if len(opts.Tokens) > 0 {
		ap.tokens = append(ap.tokens, opts.Tokens...)
	}

	if opts.HeaderFormat != "" {
		ap.headerFormat = opts.HeaderFormat
	}

	return ap
}

// GetToken returns the current token, or the next token in the sequence if configured.
func (ap *MockAuthProvider) GetToken(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Check for injected error
	if err := ap.getAndConsumeError(); err != nil {
		if ap.recordCalls.Load() {
			ap.mu.Lock()
			ap.getCalls = append(ap.getCalls, GetTokenCall{
				Time:   time.Now(),
				Result: "",
				Err:    err,
			})
			ap.mu.Unlock()
		}
		return "", err
	}

	token := ap.getCurrentToken()

	if ap.recordCalls.Load() {
		ap.mu.Lock()
		ap.getCalls = append(ap.getCalls, GetTokenCall{
			Time:   time.Now(),
			Result: token,
			Err:    nil,
		})
		ap.mu.Unlock()
	}

	return token, nil
}

// getCurrentToken returns the current or next token in the sequence.
func (ap *MockAuthProvider) getCurrentToken() string {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	if len(ap.tokens) == 0 {
		return ""
	}

	idx := ap.tokenIndex.Load()
	if int(idx) >= len(ap.tokens) {
		return ap.tokens[len(ap.tokens)-1]
	}

	return ap.tokens[idx]
}

// GetAuthHeader returns the formatted authentication header value.
func (ap *MockAuthProvider) GetAuthHeader(token string) string {
	header := formatHeader(ap.headerFormat, token)

	if ap.recordCalls.Load() {
		ap.mu.Lock()
		ap.headerCalls = append(ap.headerCalls, GetHeaderCall{
			Token:  token,
			Result: header,
			Time:   time.Now(),
		})
		ap.mu.Unlock()
	}

	return header
}

// formatHeader substitutes {} with the token in the format string.
func formatHeader(format, token string) string {
	for i := 0; i < len(format)-1; i++ {
		if format[i] == '{' && format[i+1] == '}' {
			return format[:i] + token + format[i+2:]
		}
	}
	return format + " " + token
}

// getAndConsumeError retrieves and clears an injected error (one-shot).
func (ap *MockAuthProvider) getAndConsumeError() error {
	val := ap.tokenErr.Load()
	if val == nil {
		return nil
	}

	holder := val.(*errHolder)
	if holder == nil || holder.err == nil {
		return nil
	}

	err := holder.err
	if holder.oneShot {
		ap.tokenErr.Store((*errHolder)(nil))
	}
	return err
}

// InjectTokenError injects an error for the next GetToken call (one-shot by default).
func (ap *MockAuthProvider) InjectTokenError(err error) {
	if err == nil {
		ap.tokenErr.Store((*errHolder)(nil))
	} else {
		ap.tokenErr.Store(&errHolder{err: err, oneShot: true})
	}
}

// ClearInjectedErrors clears all injected errors.
func (ap *MockAuthProvider) ClearInjectedErrors() {
	ap.tokenErr.Store((*errHolder)(nil))
}

// SetTokens sets the token sequence for rotation testing.
func (ap *MockAuthProvider) SetTokens(tokens []string) {
	ap.mu.Lock()
	ap.tokens = make([]string, len(tokens))
	copy(ap.tokens, tokens)
	ap.mu.Unlock()
	ap.tokenIndex.Store(0)
}

// RotateToken advances to the next token in the sequence.
func (ap *MockAuthProvider) RotateToken() {
	ap.mu.RLock()
	tokenCount := len(ap.tokens)
	ap.mu.RUnlock()

	if tokenCount == 0 {
		return
	}

	idx := ap.tokenIndex.Load()
	if int(idx)+1 < tokenCount {
		ap.tokenIndex.Store(idx + 1)
	}
}

// CurrentTokenIndex returns the index of the current token.
func (ap *MockAuthProvider) CurrentTokenIndex() int {
	ap.mu.RLock()
	tokenCount := len(ap.tokens)
	ap.mu.RUnlock()

	idx := ap.tokenIndex.Load()
	if tokenCount == 0 {
		return -1
	}
	if int(idx) >= tokenCount {
		return tokenCount - 1
	}
	return int(idx)
}

// Reset resets the token sequence to the first token.
func (ap *MockAuthProvider) Reset() {
	ap.tokenIndex.Store(0)
	ap.mu.Lock()
	ap.getCalls = ap.getCalls[:0]
	ap.headerCalls = ap.headerCalls[:0]
	ap.mu.Unlock()
	ap.tokenErr.Store((*errHolder)(nil))
}

// GetTokenCalls returns a copy of all recorded GetToken calls.
func (ap *MockAuthProvider) GetTokenCalls() []GetTokenCall {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	calls := make([]GetTokenCall, len(ap.getCalls))
	copy(calls, ap.getCalls)
	return calls
}

// GetHeaderCalls returns a copy of all recorded GetAuthHeader calls.
func (ap *MockAuthProvider) GetHeaderCalls() []GetHeaderCall {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	calls := make([]GetHeaderCall, len(ap.headerCalls))
	copy(calls, ap.headerCalls)
	return calls
}

// ClearRecorded clears all recorded calls.
func (ap *MockAuthProvider) ClearRecorded() {
	ap.mu.Lock()
	ap.getCalls = ap.getCalls[:0]
	ap.headerCalls = ap.headerCalls[:0]
	ap.mu.Unlock()
}

// SetHeaderFormat sets the format for GetAuthHeader.
// Use "{}" as a placeholder for the token.
func (ap *MockAuthProvider) SetHeaderFormat(format string) {
	ap.headerFormat = format
}
