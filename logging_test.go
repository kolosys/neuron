package neuron

import (
	"net/http"
	"strings"
	"testing"
)

func TestSafeHeadersForLogging(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		expected map[string][]string
	}{
		{
			name: "sensitive headers are excluded",
			headers: http.Header{
				"Authorization": []string{"Bearer secret-token"},
				"Cookie":        []string{"session=abc123"},
				"X-Api-Key":     []string{"api-key-value"},
				"Content-Type":  []string{"application/json"},
				"User-Agent":    []string{"test-agent"},
			},
			expected: map[string][]string{
				"Content-Type": []string{"application/json"},
				"User-Agent":   []string{"test-agent"},
			},
		},
		{
			name: "case insensitive matching",
			headers: http.Header{
				"AUTHORIZATION": []string{"Bearer token"},
				"authorization": []string{"Bearer token2"},
				"Content-Type":  []string{"application/json"},
			},
			expected: map[string][]string{
				"Content-Type": []string{"application/json"},
			},
		},
		{
			name: "no sensitive headers",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
			expected: map[string][]string{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
		},
		{
			name: "only whitelisted headers included",
			headers: http.Header{
				"Authorization":   []string{"Bearer token"},
				"Cookie":          []string{"session=123"},
				"Set-Cookie":      []string{"session=456"},
				"X-Api-Key":       []string{"key"},
				"X-Auth-Token":    []string{"token"},
				"X-Access-Token":  []string{"access"},
				"X-Refresh-Token": []string{"refresh"},
				"Api-Key":         []string{"api"},
				"Auth-Token":      []string{"auth"},
				"Content-Type":    []string{"application/json"},
				"User-Agent":      []string{"test-agent"},
				"Accept":          []string{"application/json"},
			},
			expected: map[string][]string{
				"Content-Type": []string{"application/json"},
				"User-Agent":   []string{"test-agent"},
				"Accept":       []string{"application/json"},
			},
		},
		{
			name:     "only sensitive headers",
			headers:  http.Header{"Authorization": []string{"Bearer token"}},
			expected: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := safeHeadersForLogging(tt.headers)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d headers, got %d", len(tt.expected), len(result))
			}

			for key, expectedValues := range tt.expected {
				actualValues, ok := result[key]
				if !ok {
					t.Errorf("expected header %q not found", key)
					continue
				}

				if len(actualValues) != len(expectedValues) {
					t.Errorf("header %q: expected %d values, got %d", key, len(expectedValues), len(actualValues))
					continue
				}

				for i, expected := range expectedValues {
					if actualValues[i] != expected {
						t.Errorf("header %q[%d]: expected %q, got %q", key, i, expected, actualValues[i])
					}
				}
			}

			// Verify only safe headers are present
			for key := range result {
				keyLower := strings.ToLower(key)
				if !safeHeaderNames[keyLower] {
					t.Errorf("unsafe header %q should not be in result", key)
				}
			}

			// Verify original headers are not modified
			for key, originalValues := range tt.headers {
				if len(originalValues) == 0 || originalValues[0] == "" {
					t.Errorf("original header %q was modified", key)
				}
			}
		})
	}
}

func TestSanitizeHeadersForLogging(t *testing.T) {
	tests := []struct {
		name     string
		headers  http.Header
		expected map[string][]string
	}{
		{
			name: "sensitive headers are redacted",
			headers: http.Header{
				"Authorization": []string{"Bearer secret-token"},
				"Cookie":        []string{"session=abc123"},
				"X-Api-Key":     []string{"api-key-value"},
				"Content-Type":  []string{"application/json"},
				"User-Agent":    []string{"test-agent"},
			},
			expected: map[string][]string{
				"Authorization": []string{"<redacted>"},
				"Cookie":        []string{"<redacted>"},
				"X-Api-Key":     []string{"<redacted>"},
				"Content-Type":  []string{"application/json"},
				"User-Agent":    []string{"test-agent"},
			},
		},
		{
			name: "case insensitive matching",
			headers: http.Header{
				"AUTHORIZATION": []string{"Bearer token"},
				"authorization": []string{"Bearer token2"},
				"Content-Type":  []string{"application/json"},
			},
			expected: map[string][]string{
				"AUTHORIZATION": []string{"<redacted>"},
				"authorization": []string{"<redacted>"},
				"Content-Type":  []string{"application/json"},
			},
		},
		{
			name: "no sensitive headers",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
			expected: map[string][]string{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
		},
		{
			name: "all sensitive headers redacted",
			headers: http.Header{
				"Authorization":   []string{"Bearer token"},
				"Cookie":          []string{"session=123"},
				"Set-Cookie":      []string{"session=456"},
				"X-Api-Key":       []string{"key"},
				"X-Auth-Token":    []string{"token"},
				"X-Access-Token":  []string{"access"},
				"X-Refresh-Token": []string{"refresh"},
				"Api-Key":         []string{"api"},
				"Auth-Token":      []string{"auth"},
				"Content-Type":    []string{"application/json"},
			},
			expected: map[string][]string{
				"Authorization":   []string{"<redacted>"},
				"Cookie":          []string{"<redacted>"},
				"Set-Cookie":      []string{"<redacted>"},
				"X-Api-Key":       []string{"<redacted>"},
				"X-Auth-Token":    []string{"<redacted>"},
				"X-Access-Token":  []string{"<redacted>"},
				"X-Refresh-Token": []string{"<redacted>"},
				"Api-Key":         []string{"<redacted>"},
				"Auth-Token":      []string{"<redacted>"},
				"Content-Type":    []string{"application/json"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeHeadersForLogging(tt.headers)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d headers, got %d", len(tt.expected), len(result))
			}

			for key, expectedValues := range tt.expected {
				actualValues, ok := result[key]
				if !ok {
					t.Errorf("expected header %q not found", key)
					continue
				}

				if len(actualValues) != len(expectedValues) {
					t.Errorf("header %q: expected %d values, got %d", key, len(expectedValues), len(actualValues))
					continue
				}

				for i, expected := range expectedValues {
					if actualValues[i] != expected {
						t.Errorf("header %q[%d]: expected %q, got %q", key, i, expected, actualValues[i])
					}
				}
			}

			// Verify original headers are not modified
			for key, originalValues := range tt.headers {
				keyLower := strings.ToLower(key)
				if sensitiveHeaderNames[keyLower] {
					// Sensitive headers should still have original values
					if len(originalValues) == 0 || originalValues[0] == "<redacted>" {
						t.Errorf("original header %q was modified", key)
					}
				}
			}
		})
	}
}
