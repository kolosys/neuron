package neuron_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kolosys/ion/circuit"
	. "github.com/kolosys/neuron"
	"github.com/kolosys/neuron/mock"
)

func TestStream_200LeavesBodyReadable(t *testing.T) {
	const payload = "chunk-one\nchunk-two\n"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(payload))
	}))
	t.Cleanup(ts.Close)

	client := NewClient(ClientOptions{BaseURL: ts.URL})
	resp, err := Stream[any](client, MethodPOST, "/v1/stream", nil, nil)
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if string(got) != payload {
		t.Fatalf("body: got %q want %q", got, payload)
	}
}

func TestStream_500DoesNotRetry(t *testing.T) {
	var posts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			posts.Add(1)
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("upstream failed"))
	}))
	t.Cleanup(ts.Close)

	client := NewClient(ClientOptions{
		BaseURL:    ts.URL,
		MaxRetries: 3,
		RetryDelay: time.Millisecond,
	})
	resp, err := Stream[any](client, MethodPOST, "/v1/stream", map[string]string{"n": "1"}, nil)
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status: got %d", resp.StatusCode)
	}
	if n := posts.Load(); n != 1 {
		t.Fatalf("POSTs: got %d want 1", n)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if string(got) != "upstream failed" {
		t.Fatalf("body: got %q", got)
	}
}

func TestStream_TransportErrorRetriesOnce(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.InjectError(errors.New("connection refused"))
	rt.QueueResponse(mock.ResponseConfig{StatusCode: http.StatusOK, Body: []byte("ok")})

	client := NewClient(ClientOptions{
		BaseURL:    "http://example.com",
		HTTPClient: &http.Client{Transport: rt},
		MaxRetries: 1,
		RetryDelay: time.Millisecond,
	})
	resp, err := Stream[any](client, MethodPOST, "/v1/stream", nil, nil)
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d", resp.StatusCode)
	}
	mock.AssertRequestCount(t, rt, 2)
}

func TestStream_CircuitOpenFailsBeforeDo(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: http.StatusOK, Body: []byte("should-not-see")})

	cb := circuit.New("neuron-stream",
		circuit.WithFailureThreshold(1),
		circuit.WithRecovery(circuit.RecoveryManual),
	)
	cb.RecordFailure()

	client := NewClient(ClientOptions{
		BaseURL:    "http://example.com",
		HTTPClient: &http.Client{Transport: rt},
		Circuit:    cb,
	})
	resp, err := Stream[any](client, MethodPOST, "/v1/stream", nil, nil)
	if resp != nil {
		resp.Body.Close()
		t.Fatal("expected no response")
	}
	if err == nil {
		t.Fatal("expected circuit-open error")
	}

	var ce ClientError
	if !errors.As(err, &ce) {
		t.Fatalf("error type: %T (%v)", err, err)
	}
	if ce.Type != ErrorTypeCircuit {
		t.Fatalf("ErrorType: got %v want ErrorTypeCircuit", ce.Type)
	}

	var ionErr *circuit.CircuitError
	if !errors.As(err, &ionErr) || !ionErr.IsCircuitOpen() {
		t.Fatalf("cause: want open circuit.CircuitError, got %v", err)
	}
	mock.AssertNoRequests(t, rt)
}
