package neuron_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/kolosys/neuron"
)

func TestMetricsCollector_Concurrency(t *testing.T) {
	collector := NewMetricsCollector()
	const iterations = 1000
	const goroutines = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				collector.RecordRequest()
				collector.RecordResponse(200, 10*time.Millisecond)
				collector.RecordResponse(500, 20*time.Millisecond)
			}
		}()
	}

	wg.Wait()

	metrics := collector.GetMetrics()

	expectedRequests := int64(goroutines * iterations)
	expectedResponses := int64(goroutines * iterations * 2)
	expectedErrors := int64(goroutines * iterations)

	if metrics.RequestCount != expectedRequests {
		t.Errorf("expected %d requests, got %d", expectedRequests, metrics.RequestCount)
	}
	if metrics.ResponseCount != expectedResponses {
		t.Errorf("expected %d responses, got %d", expectedResponses, metrics.ResponseCount)
	}
	if metrics.ErrorCount != expectedErrors {
		t.Errorf("expected %d errors, got %d", expectedErrors, metrics.ErrorCount)
	}

	if metrics.StatusCodeCounts[200] != int64(goroutines*iterations) {
		t.Errorf("expected %d OKs, got %d", goroutines*iterations, metrics.StatusCodeCounts[200])
	}
	if metrics.StatusCodeCounts[500] != int64(goroutines*iterations) {
		t.Errorf("expected %d Errors, got %d", goroutines*iterations, metrics.StatusCodeCounts[500])
	}

	if metrics.MinDuration != 10*time.Millisecond {
		t.Errorf("expected min duration 10ms, got %v", metrics.MinDuration)
	}
	if metrics.MaxDuration != 20*time.Millisecond {
		t.Errorf("expected max duration 20ms, got %v", metrics.MaxDuration)
	}
}
