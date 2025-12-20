package neuron

import (
	"net/http"
	"sync"
	"time"
)

// MetricsCollector collects and stores metrics
type MetricsCollector struct {
	mu sync.RWMutex

	// Request metrics
	RequestCount  int64
	ResponseCount int64
	ErrorCount    int64

	// Duration metrics
	TotalDuration time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration

	// Status code metrics
	StatusCodeCounts map[int]int64

	// Start time for uptime calculation
	StartTime time.Time
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		StatusCodeCounts: make(map[int]int64),
		StartTime:        time.Now(),
	}
}

// RecordRequest records a request metric
func (m *MetricsCollector) RecordRequest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestCount++
}

// RecordResponse records a response metric
func (m *MetricsCollector) RecordResponse(statusCode int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ResponseCount++
	m.StatusCodeCounts[statusCode]++

	// Update duration metrics
	if m.ResponseCount == 1 {
		m.MinDuration = duration
		m.MaxDuration = duration
	} else {
		if duration < m.MinDuration {
			m.MinDuration = duration
		}
		if duration > m.MaxDuration {
			m.MaxDuration = duration
		}
	}

	m.TotalDuration += duration

	// Count errors (4xx, 5xx)
	if statusCode >= 400 {
		m.ErrorCount++
	}
}

// GetMetrics returns current metrics
func (m *MetricsCollector) GetMetrics() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	avgDuration := time.Duration(0)
	if m.ResponseCount > 0 {
		avgDuration = m.TotalDuration / time.Duration(m.ResponseCount)
	}

	return MetricsSnapshot{
		RequestCount:     m.RequestCount,
		ResponseCount:    m.ResponseCount,
		ErrorCount:       m.ErrorCount,
		AverageDuration:  avgDuration,
		MinDuration:      m.MinDuration,
		MaxDuration:      m.MaxDuration,
		StatusCodeCounts: copyStatusCodeCounts(m.StatusCodeCounts),
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

// copyStatusCodeCounts creates a copy of the status code counts map
func copyStatusCodeCounts(src map[int]int64) map[int]int64 {
	dst := make(map[int]int64)
	for k, v := range src {
		dst[k] = v
	}
	return dst
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
		start, ok := resp.Request.Context().Value("request_start").(time.Time)
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
		start, ok := resp.Request.Context().Value("request_start").(time.Time)
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
