package neuron_test

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kolosys/neuron"
)

func TestNewDeduplicator(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		d := neuron.NewDeduplicator(neuron.DedupeConfig{})
		if d == nil {
			t.Fatal("NewDeduplicator returned nil")
		}
	})

	t.Run("custom config", func(t *testing.T) {
		config := neuron.DedupeConfig{
			Enabled:    true,
			WindowSize: 10 * time.Second,
			MaxSize:    5000,
		}
		d := neuron.NewDeduplicator(config)
		if d == nil {
			t.Fatal("NewDeduplicator returned nil")
		}
	})
}

func TestDeduplicator_Dedupe(t *testing.T) {
	t.Run("disabled deduplication passes through", func(t *testing.T) {
		d := neuron.NewDeduplicator(neuron.DedupeConfig{Enabled: false})

		var callCount atomic.Int32
		fn := func() (*http.Response, error) {
			callCount.Add(1)
			return &http.Response{StatusCode: 200}, nil
		}

		_, err := d.Dedupe(context.Background(), "key", fn)
		if err != nil {
			t.Errorf("Dedupe() error = %v", err)
		}

		if callCount.Load() != 1 {
			t.Errorf("expected 1 call, got %d", callCount.Load())
		}
	})

	t.Run("empty key passes through", func(t *testing.T) {
		d := neuron.NewDeduplicator(neuron.DedupeConfig{Enabled: true})

		var callCount atomic.Int32
		fn := func() (*http.Response, error) {
			callCount.Add(1)
			return &http.Response{StatusCode: 200}, nil
		}

		_, err := d.Dedupe(context.Background(), "", fn)
		if err != nil {
			t.Errorf("Dedupe() error = %v", err)
		}

		if callCount.Load() != 1 {
			t.Errorf("expected 1 call, got %d", callCount.Load())
		}
	})

	t.Run("concurrent requests are deduplicated", func(t *testing.T) {
		d := neuron.NewDeduplicator(neuron.DedupeConfig{
			Enabled:    true,
			WindowSize: 5 * time.Second,
		})

		var callCount atomic.Int32
		fn := func() (*http.Response, error) {
			callCount.Add(1)
			time.Sleep(100 * time.Millisecond) // Simulate network delay
			return &http.Response{StatusCode: 200}, nil
		}

		var wg sync.WaitGroup
		results := make([]*http.Response, 5)
		errs := make([]error, 5)

		for i := range 5 {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx], errs[idx] = d.Dedupe(context.Background(), "same-key", fn)
			}(i)
		}

		wg.Wait()

		// Only one actual request should be made
		if callCount.Load() != 1 {
			t.Errorf("expected 1 call, got %d", callCount.Load())
		}

		// All requests should succeed
		for i, err := range errs {
			if err != nil {
				t.Errorf("request %d error = %v", i, err)
			}
		}

		for i, resp := range results {
			if resp == nil || resp.StatusCode != 200 {
				t.Errorf("request %d got invalid response", i)
			}
		}
	})

	t.Run("different keys are not deduplicated", func(t *testing.T) {
		d := neuron.NewDeduplicator(neuron.DedupeConfig{
			Enabled:    true,
			WindowSize: 5 * time.Second,
		})

		var callCount atomic.Int32
		fn := func() (*http.Response, error) {
			callCount.Add(1)
			return &http.Response{StatusCode: 200}, nil
		}

		var wg sync.WaitGroup
		for i := range 5 {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				key := "key-" + string(rune('a'+idx))
				d.Dedupe(context.Background(), key, fn)
			}(i)
		}

		wg.Wait()

		if callCount.Load() != 5 {
			t.Errorf("expected 5 calls, got %d", callCount.Load())
		}
	})

	t.Run("error propagates to all waiters", func(t *testing.T) {
		d := neuron.NewDeduplicator(neuron.DedupeConfig{
			Enabled:    true,
			WindowSize: 5 * time.Second,
		})

		expectedErr := errors.New("network error")
		fn := func() (*http.Response, error) {
			time.Sleep(50 * time.Millisecond)
			return nil, expectedErr
		}

		var wg sync.WaitGroup
		errs := make([]error, 3)

		for i := range 3 {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				_, errs[idx] = d.Dedupe(context.Background(), "error-key", fn)
			}(i)
		}

		wg.Wait()

		for i, err := range errs {
			if err == nil || err.Error() != expectedErr.Error() {
				t.Errorf("request %d: expected error %q, got %v", i, expectedErr, err)
			}
		}
	})

	t.Run("context cancellation for waiting request", func(t *testing.T) {
		d := neuron.NewDeduplicator(neuron.DedupeConfig{
			Enabled:    true,
			WindowSize: 5 * time.Second,
		})

		fn := func() (*http.Response, error) {
			time.Sleep(500 * time.Millisecond)
			return &http.Response{StatusCode: 200}, nil
		}

		var wg sync.WaitGroup
		var waitingErr error

		// Start the first request (will execute the function)
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Dedupe(context.Background(), "timeout-key", fn)
		}()

		// Give the first request time to start
		time.Sleep(20 * time.Millisecond)

		// Start a second request with a short timeout (will wait for first request)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_, waitingErr = d.Dedupe(ctx, "timeout-key", fn)
		}()

		wg.Wait()

		if !errors.Is(waitingErr, context.DeadlineExceeded) {
			t.Errorf("expected context.DeadlineExceeded, got %v", waitingErr)
		}
	})
}

func TestDeduplicator_Stats(t *testing.T) {
	d := neuron.NewDeduplicator(neuron.DedupeConfig{
		Enabled:    true,
		WindowSize: 5 * time.Second,
	})

	fn := func() (*http.Response, error) {
		time.Sleep(50 * time.Millisecond)
		return &http.Response{StatusCode: 200}, nil
	}

	var wg sync.WaitGroup
	for range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Dedupe(context.Background(), "stats-key", fn)
		}()
	}

	wg.Wait()

	stats := d.Stats()
	if stats.Deduped < 2 {
		t.Errorf("expected at least 2 deduped requests, got %d", stats.Deduped)
	}
}

func TestGenerateDedupeKey(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		url      string
		expected string
	}{
		{
			name:     "GET request",
			method:   "GET",
			url:      "https://api.example.com/users",
			expected: "GET https://api.example.com/users",
		},
		{
			name:     "POST request with query params",
			method:   "POST",
			url:      "https://api.example.com/users?page=1",
			expected: "POST https://api.example.com/users?page=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			key := neuron.GenerateDedupeKey(req)
			if key != tt.expected {
				t.Errorf("GenerateDedupeKey() = %q, want %q", key, tt.expected)
			}
		})
	}
}

func TestDefaultDedupeConfig(t *testing.T) {
	config := neuron.DefaultDedupeConfig()

	if !config.Enabled {
		t.Error("default config should be enabled")
	}
	if config.WindowSize != 5*time.Second {
		t.Errorf("WindowSize = %v, want 5s", config.WindowSize)
	}
	if config.MaxSize != 10000 {
		t.Errorf("MaxSize = %d, want 10000", config.MaxSize)
	}
}

func BenchmarkDeduplicator_NoDedupe(b *testing.B) {
	d := neuron.NewDeduplicator(neuron.DedupeConfig{
		Enabled:    true,
		WindowSize: 5 * time.Second,
	})

	fn := func() (*http.Response, error) {
		return &http.Response{StatusCode: 200}, nil
	}

	for i := 0; b.Loop(); i++ {
		key := "unique-key-" + string(rune(i%1000))
		d.Dedupe(context.Background(), key, fn)
	}
}

func BenchmarkDeduplicator_WithDedupe(b *testing.B) {
	d := neuron.NewDeduplicator(neuron.DedupeConfig{
		Enabled:    true,
		WindowSize: 5 * time.Second,
	})

	fn := func() (*http.Response, error) {
		time.Sleep(time.Millisecond)
		return &http.Response{StatusCode: 200}, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			d.Dedupe(context.Background(), "same-key", fn)
		}
	})
}
