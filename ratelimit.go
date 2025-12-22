package neuron

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RateLimitInfo contains parsed rate limit information from response headers
type RateLimitInfo struct {
	// Limit is the maximum number of requests allowed in the current window
	Limit int

	// Remaining is the number of requests remaining in the current window
	Remaining int

	// Reset is the absolute time when the rate limit resets
	Reset time.Time

	// ResetAfter is the duration until the rate limit resets
	ResetAfter time.Duration

	// Bucket is the rate limit bucket identifier (Discord-style)
	Bucket string

	// Global indicates if this is a global rate limit
	Global bool

	// RetryAfter is the duration to wait before retrying (from Retry-After header)
	RetryAfter time.Duration

	// Scope indicates the rate limit scope (user, global, shared)
	Scope string
}

// IsExhausted returns true if the rate limit has been exhausted
func (r *RateLimitInfo) IsExhausted() bool {
	return r.Remaining <= 0
}

// WaitDuration returns how long to wait before the next request
func (r *RateLimitInfo) WaitDuration() time.Duration {
	if r.RetryAfter > 0 {
		return r.RetryAfter
	}
	if r.ResetAfter > 0 {
		return r.ResetAfter
	}
	if !r.Reset.IsZero() {
		return time.Until(r.Reset)
	}
	return 0
}

// RateLimiter is the interface for rate limiting requests
type RateLimiter interface {
	// Allow checks if a request is allowed without blocking
	// Returns true if the request can proceed immediately
	Allow(ctx context.Context, method, endpoint string) bool

	// Wait blocks until the request is allowed or context is cancelled
	Wait(ctx context.Context, method, endpoint string) error
}

// RateLimitUpdater receives rate limit info from response headers
type RateLimitUpdater interface {
	// UpdateFromHeaders updates rate limit state from response headers
	UpdateFromHeaders(method, endpoint string, info *RateLimitInfo) error
}

// RateLimitHandler combines limiter and updater for full integration
type RateLimitHandler interface {
	RateLimiter
	RateLimitUpdater
}

// ParseRateLimitHeaders extracts rate limit info from HTTP response headers
// Supports standard rate limit headers and Discord-style bucket headers
func ParseRateLimitHeaders(headers http.Header) *RateLimitInfo {
	info := &RateLimitInfo{}
	hasRateLimitInfo := false

	// Parse X-RateLimit-Limit
	if limit := headers.Get("X-RateLimit-Limit"); limit != "" {
		if v, err := strconv.Atoi(limit); err == nil {
			info.Limit = v
			hasRateLimitInfo = true
		}
	}

	// Parse X-RateLimit-Remaining
	if remaining := headers.Get("X-RateLimit-Remaining"); remaining != "" {
		if v, err := strconv.Atoi(remaining); err == nil {
			info.Remaining = v
			hasRateLimitInfo = true
		}
	}

	// Parse X-RateLimit-Reset (Unix timestamp, can be float for Discord)
	if reset := headers.Get("X-RateLimit-Reset"); reset != "" {
		if v, err := strconv.ParseFloat(reset, 64); err == nil {
			sec := int64(v)
			nsec := int64((v - float64(sec)) * 1e9)
			info.Reset = time.Unix(sec, nsec)
			hasRateLimitInfo = true
		}
	}

	// Parse X-RateLimit-Reset-After (seconds, can be float)
	if resetAfter := headers.Get("X-RateLimit-Reset-After"); resetAfter != "" {
		if v, err := strconv.ParseFloat(resetAfter, 64); err == nil {
			info.ResetAfter = time.Duration(v * float64(time.Second))
			hasRateLimitInfo = true
		}
	}

	// Parse X-RateLimit-Bucket (Discord-style)
	if bucket := headers.Get("X-RateLimit-Bucket"); bucket != "" {
		info.Bucket = bucket
		hasRateLimitInfo = true
	}

	// Parse X-RateLimit-Global
	if global := headers.Get("X-RateLimit-Global"); global != "" {
		info.Global = strings.EqualFold(global, "true")
		if info.Global {
			hasRateLimitInfo = true
		}
	}

	// Parse X-RateLimit-Scope
	if scope := headers.Get("X-RateLimit-Scope"); scope != "" {
		info.Scope = scope
		hasRateLimitInfo = true
	}

	// Parse Retry-After header (seconds or HTTP date)
	if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
		hasRateLimitInfo = true
		// Try parsing as seconds first
		if seconds, err := strconv.ParseFloat(retryAfter, 64); err == nil {
			info.RetryAfter = time.Duration(seconds * float64(time.Second))
		} else {
			// Try parsing as HTTP date (RFC 7231)
			if t, err := http.ParseTime(retryAfter); err == nil {
				info.RetryAfter = max(time.Until(t), 0)
			}
		}
	}

	if !hasRateLimitInfo {
		return nil
	}

	return info
}

// NoopRateLimiter is a rate limiter that allows all requests
type NoopRateLimiter struct{}

// Allow always returns true
func (n *NoopRateLimiter) Allow(ctx context.Context, method, endpoint string) bool {
	return true
}

// Wait always returns immediately
func (n *NoopRateLimiter) Wait(ctx context.Context, method, endpoint string) error {
	return nil
}

// NoopRateLimitUpdater is an updater that does nothing
type NoopRateLimitUpdater struct{}

// UpdateFromHeaders does nothing
func (n *NoopRateLimitUpdater) UpdateFromHeaders(method, endpoint string, info *RateLimitInfo) error {
	return nil
}
