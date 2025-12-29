package mock

import (
	"bytes"
	"io"
	"sync"
	"sync/atomic"
)

// MockValidator is a mock implementation of neuron.Validator for testing.
type MockValidator struct {
	calls       []ValidateCall
	validateErr atomic.Value // *errHolder
	recordCalls atomic.Bool
	mu          sync.RWMutex
}

// ValidateCall records a call to Validate.
type ValidateCall struct {
	Data        []byte
	ContentType string
	Result      error
}

// MockValidatorOptions configures mock validator behavior.
type MockValidatorOptions struct {
	RecordCalls bool
}

// NewMockValidator creates a new mock validator.
func NewMockValidator(opts *MockValidatorOptions) *MockValidator {
	mv := &MockValidator{
		calls: make([]ValidateCall, 0),
	}

	if opts == nil {
		mv.recordCalls.Store(true)
		return mv
	}

	mv.recordCalls.Store(opts.RecordCalls)
	return mv
}

// Validate validates the provided data.
func (mv *MockValidator) Validate(data []byte, contentType string) error {
	err := mv.getAndConsumeError()

	if mv.recordCalls.Load() {
		mv.mu.Lock()
		mv.calls = append(mv.calls, ValidateCall{
			Data:        data,
			ContentType: contentType,
			Result:      err,
		})
		mv.mu.Unlock()
	}

	return err
}

// getAndConsumeError retrieves and clears an injected error (one-shot).
func (mv *MockValidator) getAndConsumeError() error {
	val := mv.validateErr.Load()
	if val == nil {
		return nil
	}

	holder := val.(*errHolder)
	if holder == nil || holder.err == nil {
		return nil
	}

	err := holder.err
	if holder.oneShot {
		mv.validateErr.Store((*errHolder)(nil))
	}
	return err
}

// InjectValidationError injects an error for the next Validate call.
func (mv *MockValidator) InjectValidationError(err error) {
	if err == nil {
		mv.validateErr.Store((*errHolder)(nil))
	} else {
		mv.validateErr.Store(&errHolder{err: err, oneShot: true})
	}
}

// ClearInjectedErrors clears all injected errors.
func (mv *MockValidator) ClearInjectedErrors() {
	mv.validateErr.Store((*errHolder)(nil))
}

// ValidateCalls returns a copy of all recorded Validate calls.
func (mv *MockValidator) ValidateCalls() []ValidateCall {
	mv.mu.RLock()
	defer mv.mu.RUnlock()

	calls := make([]ValidateCall, len(mv.calls))
	copy(calls, mv.calls)
	return calls
}

// ClearRecorded clears all recorded calls.
func (mv *MockValidator) ClearRecorded() {
	mv.mu.Lock()
	mv.calls = mv.calls[:0]
	mv.mu.Unlock()
}

// Reset resets all state.
func (mv *MockValidator) Reset() {
	mv.mu.Lock()
	mv.calls = mv.calls[:0]
	mv.mu.Unlock()
	mv.validateErr.Store((*errHolder)(nil))
}

// MockRequestIDGenerator is a mock implementation of neuron.RequestIDGenerator for testing.
type MockRequestIDGenerator struct {
	ids       []string
	idIndex   atomic.Int32
	callCount atomic.Int64
}

// MockRequestIDGeneratorOptions configures mock request ID generator behavior.
type MockRequestIDGeneratorOptions struct {
	IDs []string
}

// NewMockRequestIDGenerator creates a new mock request ID generator.
func NewMockRequestIDGenerator(opts *MockRequestIDGeneratorOptions) *MockRequestIDGenerator {
	rig := &MockRequestIDGenerator{
		ids: make([]string, 0),
	}

	if opts != nil && len(opts.IDs) > 0 {
		rig.ids = make([]string, len(opts.IDs))
		copy(rig.ids, opts.IDs)
	}

	return rig
}

// Generate returns the next request ID in the sequence.
func (rig *MockRequestIDGenerator) Generate() string {
	rig.callCount.Add(1)

	if len(rig.ids) == 0 {
		return ""
	}

	idx := rig.idIndex.Load()
	if int(idx) >= len(rig.ids) {
		return rig.ids[len(rig.ids)-1]
	}

	id := rig.ids[idx]
	rig.idIndex.Store(idx + 1)
	return id
}

// CallCount returns the number of times Generate was called.
func (rig *MockRequestIDGenerator) CallCount() int64 {
	return rig.callCount.Load()
}

// SetIDs sets the ID sequence.
func (rig *MockRequestIDGenerator) SetIDs(ids []string) {
	rig.ids = make([]string, len(ids))
	copy(rig.ids, ids)
	rig.idIndex.Store(0)
	rig.callCount.Store(0)
}

// Reset resets the ID generator to the first ID.
func (rig *MockRequestIDGenerator) Reset() {
	rig.idIndex.Store(0)
	rig.callCount.Store(0)
}

// CurrentIndex returns the current index in the ID sequence.
func (rig *MockRequestIDGenerator) CurrentIndex() int {
	return int(rig.idIndex.Load())
}

// MockBodyProvider is a mock implementation of neuron.BodyProvider for testing.
type MockBodyProvider struct {
	contentType string
	body        []byte
	bodyErr     atomic.Value // *errHolder
	callCount   atomic.Int64
}

// MockBodyProviderOptions configures mock body provider behavior.
type MockBodyProviderOptions struct {
	ContentType string
	Body        []byte
}

// NewMockBodyProvider creates a new mock body provider.
func NewMockBodyProvider(opts *MockBodyProviderOptions) *MockBodyProvider {
	mbp := &MockBodyProvider{
		contentType: "application/octet-stream",
		body:        []byte{},
	}

	if opts == nil {
		return mbp
	}

	if opts.ContentType != "" {
		mbp.contentType = opts.ContentType
	}
	if opts.Body != nil {
		mbp.body = opts.Body
	}

	return mbp
}

// ContentType returns the content type of the body.
func (mbp *MockBodyProvider) ContentType() string {
	return mbp.contentType
}

// Body returns the body as an io.Reader.
func (mbp *MockBodyProvider) Body() (io.Reader, error) {
	mbp.callCount.Add(1)

	// Check for injected error
	if err := mbp.getAndConsumeError(); err != nil {
		return nil, err
	}

	return bytes.NewReader(mbp.body), nil
}

// getAndConsumeError retrieves and clears an injected error (one-shot).
func (mbp *MockBodyProvider) getAndConsumeError() error {
	val := mbp.bodyErr.Load()
	if val == nil {
		return nil
	}

	holder := val.(*errHolder)
	if holder == nil || holder.err == nil {
		return nil
	}

	err := holder.err
	if holder.oneShot {
		mbp.bodyErr.Store((*errHolder)(nil))
	}
	return err
}

// SetBody sets the body content.
func (mbp *MockBodyProvider) SetBody(body []byte) {
	mbp.body = body
}

// SetContentType sets the content type.
func (mbp *MockBodyProvider) SetContentType(contentType string) {
	mbp.contentType = contentType
}

// InjectBodyError injects an error for the next Body call.
func (mbp *MockBodyProvider) InjectBodyError(err error) {
	if err == nil {
		mbp.bodyErr.Store((*errHolder)(nil))
	} else {
		mbp.bodyErr.Store(&errHolder{err: err, oneShot: true})
	}
}

// ClearInjectedErrors clears all injected errors.
func (mbp *MockBodyProvider) ClearInjectedErrors() {
	mbp.bodyErr.Store((*errHolder)(nil))
}

// CallCount returns the number of times Body was called.
func (mbp *MockBodyProvider) CallCount() int64 {
	return mbp.callCount.Load()
}

// Reset resets the mock to initial state.
func (mbp *MockBodyProvider) Reset() {
	mbp.callCount.Store(0)
	mbp.bodyErr.Store((*errHolder)(nil))
}
