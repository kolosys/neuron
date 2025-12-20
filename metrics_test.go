package neuron

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()

	// Test request recording
	collector.RecordRequest()
	collector.RecordRequest()
	collector.RecordRequest()

	// Test response recording
	collector.RecordResponse(200, 100*time.Millisecond)
	collector.RecordResponse(201, 200*time.Millisecond)
	collector.RecordResponse(400, 150*time.Millisecond)

	metrics := collector.GetMetrics()

	if metrics.RequestCount != 3 {
		t.Errorf("expected RequestCount=3, got %d", metrics.RequestCount)
	}

	if metrics.ResponseCount != 3 {
		t.Errorf("expected ResponseCount=3, got %d", metrics.ResponseCount)
	}

	if metrics.ErrorCount != 1 { // Only 400 is an error
		t.Errorf("expected ErrorCount=1, got %d", metrics.ErrorCount)
	}

	// Check status code counts
	if metrics.StatusCodeCounts[200] != 1 {
		t.Errorf("expected status 200 count=1, got %d", metrics.StatusCodeCounts[200])
	}
}

func TestAddMetrics(t *testing.T) {
	collector := NewMetricsCollector()
	middleware := AddMetrics(collector)

	req := httptest.NewRequest("GET", "/test", nil)
	err := middleware(req)
	if err != nil {
		t.Errorf("metrics middleware failed: %v", err)
	}

	metrics := collector.GetMetrics()
	if metrics.RequestCount != 1 {
		t.Errorf("expected RequestCount=1, got %d", metrics.RequestCount)
	}
}

func TestAddMetricsResponse(t *testing.T) {
	collector := NewMetricsCollector()
	middleware := AddResponseMetrics(collector)

	req := httptest.NewRequest("GET", "/test", nil)
	// Add start time to context
	startTime := time.Now().Add(-100 * time.Millisecond)
	ctx := context.WithValue(req.Context(), "metrics_start", startTime)
	req = req.WithContext(ctx)

	resp := &http.Response{
		StatusCode: 200,
		Request:    req,
	}

	err := middleware(resp)
	if err != nil {
		t.Errorf("metrics response middleware failed: %v", err)
	}

	metrics := collector.GetMetrics()
	if metrics.ResponseCount != 1 {
		t.Errorf("expected ResponseCount=1, got %d", metrics.ResponseCount)
	}

	if metrics.StatusCodeCounts[200] != 1 {
		t.Errorf("expected status 200 count=1, got %d", metrics.StatusCodeCounts[200])
	}
}

func TestAddAutoMetrics(t *testing.T) {
	reqMw, respMw, getMetrics := AddAutoMetrics()

	req := httptest.NewRequest("GET", "/test", nil)
	err := reqMw(req)
	if err != nil {
		t.Errorf("auto metrics request middleware failed: %v", err)
	}

	// Test response middleware
	resp := &http.Response{
		StatusCode: 200,
		Request:    req,
	}

	err = respMw(resp)
	if err != nil {
		t.Errorf("auto metrics response middleware failed: %v", err)
	}

	metrics := getMetrics()
	if metrics.RequestCount != 1 {
		t.Errorf("expected RequestCount=1, got %d", metrics.RequestCount)
	}

	if metrics.ResponseCount != 1 {
		t.Errorf("expected ResponseCount=1, got %d", metrics.ResponseCount)
	}
}

func TestMetricsWithDifferentStatusCodes(t *testing.T) {
	collector := NewMetricsCollector()

	statusCodes := []int{200, 201, 400, 404, 500, 503}

	for _, code := range statusCodes {
		collector.RecordResponse(code, 100*time.Millisecond)
	}

	metrics := collector.GetMetrics()

	// Verify error count (4xx, 5xx)
	expectedErrors := int64(4) // 400, 404, 500, 503
	if metrics.ErrorCount != expectedErrors {
		t.Errorf("expected ErrorCount=%d, got %d", expectedErrors, metrics.ErrorCount)
	}

	// Verify response count
	if metrics.ResponseCount != int64(len(statusCodes)) {
		t.Errorf("expected ResponseCount=%d, got %d", len(statusCodes), metrics.ResponseCount)
	}
}

func TestMetricsDurationTracking(t *testing.T) {
	collector := NewMetricsCollector()

	durations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		150 * time.Millisecond,
	}

	for _, d := range durations {
		collector.RecordResponse(200, d)
	}

	metrics := collector.GetMetrics()

	if metrics.MinDuration != 100*time.Millisecond {
		t.Errorf("expected MinDuration=100ms, got %v", metrics.MinDuration)
	}

	if metrics.MaxDuration != 200*time.Millisecond {
		t.Errorf("expected MaxDuration=200ms, got %v", metrics.MaxDuration)
	}

	expectedAvg := 150 * time.Millisecond
	if metrics.AverageDuration != expectedAvg {
		t.Errorf("expected AverageDuration=%v, got %v", expectedAvg, metrics.AverageDuration)
	}
}
