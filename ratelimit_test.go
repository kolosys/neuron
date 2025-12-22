package neuron_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/kolosys/neuron"
)

func TestParseRateLimitHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected *neuron.RateLimitInfo
	}{
		{
			name:     "no headers returns nil",
			headers:  map[string]string{},
			expected: nil,
		},
		{
			name: "standard headers",
			headers: map[string]string{
				"X-RateLimit-Limit":     "100",
				"X-RateLimit-Remaining": "99",
			},
			expected: &neuron.RateLimitInfo{
				Limit:     100,
				Remaining: 99,
			},
		},
		{
			name: "discord style headers",
			headers: map[string]string{
				"X-RateLimit-Limit":       "10",
				"X-RateLimit-Remaining":   "5",
				"X-RateLimit-Reset":       strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10),
				"X-RateLimit-Reset-After": "60.5",
				"X-RateLimit-Bucket":      "abc123",
				"X-RateLimit-Global":      "false",
				"X-RateLimit-Scope":       "user",
			},
			expected: &neuron.RateLimitInfo{
				Limit:      10,
				Remaining:  5,
				ResetAfter: 60500 * time.Millisecond,
				Bucket:     "abc123",
				Global:     false,
				Scope:      "user",
			},
		},
		{
			name: "global rate limit",
			headers: map[string]string{
				"X-RateLimit-Global": "true",
				"Retry-After":        "30",
			},
			expected: &neuron.RateLimitInfo{
				Global:     true,
				RetryAfter: 30 * time.Second,
			},
		},
		{
			name: "retry-after as seconds float",
			headers: map[string]string{
				"X-RateLimit-Remaining": "0",
				"Retry-After":           "1.5",
			},
			expected: &neuron.RateLimitInfo{
				Remaining:  0,
				RetryAfter: 1500 * time.Millisecond,
			},
		},
		{
			name: "reset timestamp with fractional seconds",
			headers: map[string]string{
				"X-RateLimit-Reset": "1704067200.5",
			},
			expected: &neuron.RateLimitInfo{
				Reset: time.Unix(1704067200, 500000000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			for k, v := range tt.headers {
				headers.Set(k, v)
			}

			result := neuron.ParseRateLimitHeaders(headers)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Limit != tt.expected.Limit {
				t.Errorf("Limit: expected %d, got %d", tt.expected.Limit, result.Limit)
			}
			if result.Remaining != tt.expected.Remaining {
				t.Errorf("Remaining: expected %d, got %d", tt.expected.Remaining, result.Remaining)
			}
			if result.Bucket != tt.expected.Bucket {
				t.Errorf("Bucket: expected %s, got %s", tt.expected.Bucket, result.Bucket)
			}
			if result.Global != tt.expected.Global {
				t.Errorf("Global: expected %v, got %v", tt.expected.Global, result.Global)
			}
			if result.Scope != tt.expected.Scope {
				t.Errorf("Scope: expected %s, got %s", tt.expected.Scope, result.Scope)
			}
			if result.ResetAfter != tt.expected.ResetAfter {
				t.Errorf("ResetAfter: expected %v, got %v", tt.expected.ResetAfter, result.ResetAfter)
			}
			if result.RetryAfter != tt.expected.RetryAfter {
				t.Errorf("RetryAfter: expected %v, got %v", tt.expected.RetryAfter, result.RetryAfter)
			}
			if !tt.expected.Reset.IsZero() && !result.Reset.Equal(tt.expected.Reset) {
				t.Errorf("Reset: expected %v, got %v", tt.expected.Reset, result.Reset)
			}
		})
	}
}

func TestRateLimitInfo_IsExhausted(t *testing.T) {
	tests := []struct {
		name      string
		remaining int
		expected  bool
	}{
		{"zero remaining", 0, true},
		{"negative remaining", -1, true},
		{"some remaining", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &neuron.RateLimitInfo{Remaining: tt.remaining}
			if got := info.IsExhausted(); got != tt.expected {
				t.Errorf("IsExhausted() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRateLimitInfo_WaitDuration(t *testing.T) {
	tests := []struct {
		name     string
		info     neuron.RateLimitInfo
		expected time.Duration
	}{
		{
			name:     "prefer RetryAfter",
			info:     neuron.RateLimitInfo{RetryAfter: 5 * time.Second, ResetAfter: 10 * time.Second},
			expected: 5 * time.Second,
		},
		{
			name:     "use ResetAfter when no RetryAfter",
			info:     neuron.RateLimitInfo{ResetAfter: 10 * time.Second},
			expected: 10 * time.Second,
		},
		{
			name:     "use Reset time when no durations",
			info:     neuron.RateLimitInfo{Reset: time.Now().Add(3 * time.Second)},
			expected: 3 * time.Second,
		},
		{
			name:     "zero when nothing set",
			info:     neuron.RateLimitInfo{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.WaitDuration()
			// Allow some tolerance for time-based calculations
			if tt.expected == 0 {
				if got != 0 {
					t.Errorf("WaitDuration() = %v, want 0", got)
				}
			} else {
				diff := got - tt.expected
				if diff < 0 {
					diff = -diff
				}
				if diff > 100*time.Millisecond {
					t.Errorf("WaitDuration() = %v, want ~%v", got, tt.expected)
				}
			}
		})
	}
}

func TestNoopRateLimiter(t *testing.T) {
	limiter := &neuron.NoopRateLimiter{}

	t.Run("Allow always returns true", func(t *testing.T) {
		if !limiter.Allow(context.Background(), "GET", "/test") {
			t.Error("Allow() should return true")
		}
	})

	t.Run("Wait returns immediately", func(t *testing.T) {
		start := time.Now()
		err := limiter.Wait(context.Background(), "GET", "/test")
		if err != nil {
			t.Errorf("Wait() error = %v", err)
		}
		if time.Since(start) > time.Millisecond {
			t.Error("Wait() should return immediately")
		}
	})
}

func TestNoopRateLimitUpdater(t *testing.T) {
	updater := &neuron.NoopRateLimitUpdater{}

	err := updater.UpdateFromHeaders("GET", "/test", &neuron.RateLimitInfo{
		Limit:     100,
		Remaining: 50,
	})

	if err != nil {
		t.Errorf("UpdateFromHeaders() error = %v", err)
	}
}
