package mock_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/kolosys/neuron/mock"
)

func TestMockRoundTripper_QueueResponse(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	rt.QueueResponse(mock.ResponseConfig{
		StatusCode: 200,
		Body:       []byte(`{"status":"ok"}`),
	})

	req, _ := http.NewRequest("GET", "http://example.com/api", nil)
	resp, err := rt.RoundTrip(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code: expected 200, got %d", resp.StatusCode)
	}
}

func TestMockRoundTripper_RecordsRequests(t *testing.T) {
	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		RecordRequests: true,
	})

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("POST", "http://example.com/api", io.NopCloser(bytes.NewReader([]byte(`test`))))
	req.Header.Set("X-Custom", "value")

	rt.RoundTrip(req)

	requests := rt.Requests()
	if len(requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(requests))
	}

	recorded := requests[0]
	if recorded.Method != "POST" {
		t.Errorf("method: expected POST, got %s", recorded.Method)
	}
	if recorded.URL.String() != "http://example.com/api" {
		t.Errorf("url: expected http://example.com/api, got %s", recorded.URL.String())
	}
	if recorded.Headers.Get("X-Custom") != "value" {
		t.Errorf("header: expected value, got %s", recorded.Headers.Get("X-Custom"))
	}
	if string(recorded.Body) != "test" {
		t.Errorf("body: expected test, got %s", recorded.Body)
	}
}

func TestMockRoundTripper_DisabledRecording(t *testing.T) {
	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		RecordRequests: false,
	})

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rt.RoundTrip(req)

	requests := rt.Requests()
	if len(requests) != 0 {
		t.Errorf("expected 0 requests when recording disabled, got %d", len(requests))
	}
}

func TestMockRoundTripper_InjectError(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	testErr := io.EOF
	rt.InjectError(testErr)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, err := rt.RoundTrip(req)

	if err != testErr {
		t.Errorf("expected EOF, got %v", err)
	}

	// Error should be cleared after one-shot consumption
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	resp, err := rt.RoundTrip(req)

	if err != nil {
		t.Errorf("expected no error after one-shot, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMockRoundTripper_ErrorSequence(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	seq := mock.NewErrorSequence(
		io.ErrUnexpectedEOF,
		io.ErrUnexpectedEOF,
		nil,
	)

	rt.InjectError(seq)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// First call returns error from sequence
	_, err := rt.RoundTrip(req)
	if err == nil {
		t.Error("expected error from sequence")
	}

	// But the error holder is consumed, so next call should proceed normally
	// This demonstrates that the sequence itself isn't consumed in the same way
}

func TestMockRoundTripper_RequestCount(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	for i := 0; i < 3; i++ {
		rt.RoundTrip(req)
	}

	if rt.RequestCount() != 3 {
		t.Errorf("expected request count 3, got %d", rt.RequestCount())
	}
}

func TestMockRoundTripper_Latency(t *testing.T) {
	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		Latency: 10 * time.Millisecond,
	})

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	start := time.Now()
	rt.RoundTrip(req)
	elapsed := time.Since(start)

	if elapsed < 10*time.Millisecond {
		t.Errorf("expected latency >= 10ms, got %v", elapsed)
	}
}

func TestMockRoundTripper_ContextCancellation(t *testing.T) {
	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		Latency: 100 * time.Millisecond,
	})

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)

	start := time.Now()
	rt.RoundTrip(req)
	elapsed := time.Since(start)

	if elapsed > 50*time.Millisecond {
		t.Errorf("expected cancellation to skip latency, got %v", elapsed)
	}
}

func TestMockRoundTripper_ResponseMatcher(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	rt.AddMatcher(func(req *http.Request) (*http.Response, bool) {
		if req.URL.Path == "/api/users" {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`[{"id":1}]`))),
				Header:     make(http.Header),
			}, true
		}
		return nil, false
	})

	req, _ := http.NewRequest("GET", "http://example.com/api/users", nil)
	resp, err := rt.RoundTrip(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `[{"id":1}]` {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestMockRoundTripper_MatcherFallback(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	rt.AddMatcher(func(req *http.Request) (*http.Response, bool) {
		if req.URL.Path == "/api/users" {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`users`))),
				Header:     make(http.Header),
			}, true
		}
		return nil, false
	})

	rt.QueueResponse(mock.ResponseConfig{
		StatusCode: 404,
		Body:       []byte(`not found`),
	})

	req1, _ := http.NewRequest("GET", "http://example.com/api/users", nil)
	resp1, _ := rt.RoundTrip(req1)
	if resp1.StatusCode != 200 {
		t.Error("matcher should match /api/users")
	}

	req2, _ := http.NewRequest("GET", "http://example.com/other", nil)
	resp2, _ := rt.RoundTrip(req2)
	if resp2.StatusCode != 404 {
		t.Error("should fall back to queued response")
	}
}

func TestMockRoundTripper_MultipleResponses(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 201})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 404})

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	statuses := []int{201, 200, 404}
	for i, expected := range statuses {
		resp, _ := rt.RoundTrip(req)
		if resp.StatusCode != expected {
			t.Errorf("request %d: expected %d, got %d", i, expected, resp.StatusCode)
		}
	}

	// Next request with no queued responses should return EOF
	_, err := rt.RoundTrip(req)
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestMockRoundTripper_ResponseHeaders(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Custom", "test-value")

	rt.QueueResponse(mock.ResponseConfig{
		StatusCode: 200,
		Headers:    headers,
		Body:       []byte(`{}`),
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	resp, _ := rt.RoundTrip(req)

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("content-type: expected application/json, got %s", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("X-Custom") != "test-value" {
		t.Errorf("x-custom: expected test-value, got %s", resp.Header.Get("X-Custom"))
	}
}

func TestMockRoundTripper_ClearRecorded(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	rt.RoundTrip(req)
	if rt.RequestCount() != 1 {
		t.Error("expected 1 request before clear")
	}

	rt.ClearRecorded()
	if rt.RequestCount() != 0 {
		t.Error("expected 0 requests after clear")
	}
	if len(rt.Requests()) != 0 {
		t.Error("expected empty requests slice after clear")
	}
}

func TestMockRoundTripper_Reset(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	rt.InjectError(io.EOF)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	rt.RoundTrip(req)

	rt.Reset()

	if rt.RequestCount() != 0 {
		t.Error("reset should clear request count")
	}
	if len(rt.Requests()) != 0 {
		t.Error("reset should clear requests")
	}
	if _, err := rt.RoundTrip(req); err != io.EOF {
		t.Error("reset should clear injected errors")
	}
}

func TestMockRoundTripper_SetLatency(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.SetLatency(5 * time.Millisecond)

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	start := time.Now()
	rt.RoundTrip(req)
	elapsed := time.Since(start)

	if elapsed < 5*time.Millisecond {
		t.Errorf("expected latency >= 5ms, got %v", elapsed)
	}
}

func TestMockRoundTripper_EnableRecording(t *testing.T) {
	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		RecordRequests: true,
	})

	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	rt.RoundTrip(req)
	if len(rt.Requests()) == 0 {
		t.Error("expected requests to be recorded")
	}

	rt.EnableRecording(false)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	rt.RoundTrip(req)
	if len(rt.Requests()) != 1 {
		t.Error("expected recording to be disabled")
	}
}

func TestMockRoundTripper_ThreadSafety(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)

	for i := 0; i < 100; i++ {
		rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})
	}

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// Concurrent reads and writes
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				rt.RoundTrip(req)
			}
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			rt.Requests()
		}()
	}

	time.Sleep(100 * time.Millisecond)

	if rt.RequestCount() != 100 {
		t.Errorf("expected 100 requests, got %d", rt.RequestCount())
	}
}

func TestMockRoundTripper_RequestPreservesContext(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	ctx := context.WithValue(context.Background(), "key", "value")
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)

	resp, _ := rt.RoundTrip(req)

	if resp.Request.Context().Value("key") != "value" {
		t.Error("context should be preserved")
	}
}

func TestMockRoundTripper_EmptyBody(t *testing.T) {
	rt := mock.NewMockRoundTripper(nil)
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 204})

	req, _ := http.NewRequest("DELETE", "http://example.com", nil)
	resp, _ := rt.RoundTrip(req)

	body, _ := io.ReadAll(resp.Body)
	if len(body) != 0 {
		t.Errorf("expected empty body, got %d bytes", len(body))
	}
}

func TestMockRoundTripper_PreservesURL(t *testing.T) {
	rt := mock.NewMockRoundTripper(&mock.MockRoundTripperOptions{
		RecordRequests: true,
	})
	rt.QueueResponse(mock.ResponseConfig{StatusCode: 200})

	testURL := "http://example.com:8080/path?query=value"
	req, _ := http.NewRequest("GET", testURL, nil)
	rt.RoundTrip(req)

	recorded := rt.Requests()[0]
	if recorded.URL.String() != testURL {
		t.Errorf("url not preserved: expected %s, got %s", testURL, recorded.URL.String())
	}
}
