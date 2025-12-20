package neuron

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsCollector collects and stores metrics
type MetricsCollector struct {
	mu sync.RWMutex

	// Request metrics
	RequestCount  atomic.Int64
	ResponseCount atomic.Int64
	ErrorCount    atomic.Int64

	// Duration metrics
	TotalDuration atomic.Int64 // nanoseconds
	MinDuration   atomic.Int64 // nanoseconds
	MaxDuration   atomic.Int64 // nanoseconds

	// Status code metrics
	StatusCodeCounts map[int]*atomic.Int64

	// Start time for uptime calculation
	StartTime time.Time
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		StatusCodeCounts: make(map[int]*atomic.Int64),
		StartTime:        time.Now(),
	}
}

// RecordRequest records a request metric
func (m *MetricsCollector) RecordRequest() {
	m.RequestCount.Add(1)
}

// RecordResponse records a response metric
func (m *MetricsCollector) RecordResponse(statusCode int, duration time.Duration) {
	m.ResponseCount.Add(1)

	// Update status code counts
	m.mu.RLock()
	counter, ok := m.StatusCodeCounts[statusCode]
	m.mu.RUnlock()

	if !ok {
		m.mu.Lock()
		// Double check after acquiring write lock
		if counter, ok = m.StatusCodeCounts[statusCode]; !ok {
			counter = &atomic.Int64{}
			m.StatusCodeCounts[statusCode] = counter
		}
		m.mu.Unlock()
	}
	counter.Add(1)

	// Update duration metrics
	durNs := duration.Nanoseconds()
	m.TotalDuration.Add(durNs)

	// Atomic update for MinDuration
	for {
		oldMin := m.MinDuration.Load()
		if oldMin != 0 && durNs >= oldMin {
			break
		}
		if m.MinDuration.CompareAndSwap(oldMin, durNs) {
			break
		}
	}

	// Atomic update for MaxDuration
	for {
		oldMax := m.MaxDuration.Load()
		if durNs <= oldMax {
			break
		}
		if m.MaxDuration.CompareAndSwap(oldMax, durNs) {
			break
		}
	}

	// Count errors (4xx, 5xx)
	if statusCode >= 400 {
		m.ErrorCount.Add(1)
	}
}

// GetMetrics returns current metrics
func (m *MetricsCollector) GetMetrics() MetricsSnapshot {
	respCount := m.ResponseCount.Load()
	totalDur := m.TotalDuration.Load()

	avgDuration := time.Duration(0)
	if respCount > 0 {
		avgDuration = time.Duration(totalDur / respCount)
	}

	m.mu.RLock()
	counts := make(map[int]int64, len(m.StatusCodeCounts))
	for k, v := range m.StatusCodeCounts {
		counts[k] = v.Load()
	}
	m.mu.RUnlock()

	return MetricsSnapshot{
		RequestCount:     m.RequestCount.Load(),
		ResponseCount:    respCount,
		ErrorCount:       m.ErrorCount.Load(),
		AverageDuration:  avgDuration,
		MinDuration:      time.Duration(m.MinDuration.Load()),
		MaxDuration:      time.Duration(m.MaxDuration.Load()),
		StatusCodeCounts: counts,
		Uptime:           time.Since(m.StartTime),
	}
}

// MetricsSnapshot represents a snapshot of metrics at a point in time
type MetricsSnapshot struct {
	RequestCount     int64
	ResponseCount    int64
	ErrorCount       int64
	AverageDuration  time.Duration
	MinDuration      time.Duration
	MaxDuration      time.Duration
	StatusCodeCounts map[int]int64
	Uptime           time.Duration
}

// ErrorRate returns the error rate as a percentage
func (m *MetricsSnapshot) ErrorRate() float64 {
	if m.ResponseCount == 0 {
		return 0
	}
	return float64(m.ErrorCount) / float64(m.ResponseCount) * 100
}

// RequestsPerSecond returns the average requests per second
func (m *MetricsSnapshot) RequestsPerSecond() float64 {
	if m.Uptime.Seconds() == 0 {
		return 0
	}
	return float64(m.ResponseCount) / m.Uptime.Seconds()
}

// AddMetrics creates a metrics collection middleware
func AddMetrics(collector *MetricsCollector) RequestHook {
	return func(req *http.Request) error {
		collector.RecordRequest()
		return nil
	}
}

// AddResponseMetrics creates a response metrics collection middleware
func AddResponseMetrics(collector *MetricsCollector) ResponseHook {
	return func(resp *http.Response) error {
		// Get duration from context if available
		start, ok := resp.Request.Context().Value(requestStartKey).(time.Time)
		if !ok {
			start = time.Now()
		}

		duration := time.Since(start)
		collector.RecordResponse(resp.StatusCode, duration)

		return nil
	}
}

// AddAutoMetrics creates a simple metrics middleware that doesn't require a collector
func AddAutoMetrics() (RequestHook, ResponseHook, func() MetricsSnapshot) {
	collector := NewMetricsCollector()

	requestMiddleware := func(req *http.Request) error {
		collector.RecordRequest()
		return nil
	}

	responseMiddleware := func(resp *http.Response) error {
		start, ok := resp.Request.Context().Value(requestStartKey).(time.Time)
		if !ok {
			start = time.Now()
		}

		duration := time.Since(start)
		collector.RecordResponse(resp.StatusCode, duration)

		return nil
	}

	getMetrics := func() MetricsSnapshot {
		return collector.GetMetrics()
	}

	return requestMiddleware, responseMiddleware, getMetrics
}
