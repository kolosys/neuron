package neuron

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddBearerAuth(t *testing.T) {
	middleware := AddBearerAuth("my-secret-token")
	req := httptest.NewRequest("GET", "/test", nil)

	err := middleware(req)
	if err != nil {
		t.Errorf("bearer auth middleware failed: %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader != "Bearer my-secret-token" {
		t.Errorf("expected 'Bearer my-secret-token', got '%s'", authHeader)
	}
}

func TestAddAPIKeyAuth(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string
		headerName     string
		expectedHeader string
	}{
		{
			name:           "default X-API-Key",
			apiKey:         "test-key-123",
			headerName:     "X-API-Key",
			expectedHeader: "X-API-Key",
		},
		{
			name:           "custom header",
			apiKey:         "custom-key-456",
			headerName:     "X-Custom-Key",
			expectedHeader: "X-Custom-Key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := AddAPIKeyAuth(tt.apiKey, tt.headerName)
			req := httptest.NewRequest("GET", "/test", nil)

			err := middleware(req)
			if err != nil {
				t.Errorf("API key auth middleware failed: %v", err)
			}

			value := req.Header.Get(tt.expectedHeader)
			if value != tt.apiKey {
				t.Errorf("expected '%s', got '%s'", tt.apiKey, value)
			}
		})
	}
}

func TestAddBasicAuth(t *testing.T) {
	username := "testuser"
	password := "testpass"
	middleware := AddBasicAuth(username, password)
	req := httptest.NewRequest("GET", "/test", nil)

	err := middleware(req)
	if err != nil {
		t.Errorf("basic auth middleware failed: %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Basic ") {
		t.Errorf("expected 'Basic' prefix, got '%s'", authHeader)
	}

	// Decode and verify credentials
	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Errorf("failed to decode auth header: %v", err)
	}

	expected := username + ":" + password
	if string(decoded) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(decoded))
	}
}

func TestAddAuth(t *testing.T) {
	config := AuthConfig{
		Type:  AuthTypeBearer,
		Token: "test-token",
	}

	middleware := AddAuth(config)
	req := httptest.NewRequest("GET", "/test", nil)

	err := middleware(req)
	if err != nil {
		t.Errorf("auth middleware failed: %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader != "Bearer test-token" {
		t.Errorf("expected 'Bearer test-token', got '%s'", authHeader)
	}
}

func TestAddDynamicAuth(t *testing.T) {
	callCount := 0
	authFunc := func(req *http.Request) error {
		callCount++
		req.Header.Set("Authorization", "Bearer dynamic-token")
		return nil
	}

	middleware := AddDynamicAuth(authFunc)

	// First request
	req1 := httptest.NewRequest("GET", "/test", nil)
	err := middleware(req1)
	if err != nil {
		t.Errorf("dynamic auth middleware failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected authFunc to be called once, got %d", callCount)
	}

	authHeader1 := req1.Header.Get("Authorization")
	if authHeader1 != "Bearer dynamic-token" {
		t.Errorf("expected 'Bearer dynamic-token', got '%s'", authHeader1)
	}
}

func TestAddCustomAuth(t *testing.T) {
	headerName := "X-Custom-Auth"
	headerValue := "custom-value-123"

	middleware := AddCustomAuth(headerName, headerValue)
	req := httptest.NewRequest("GET", "/test", nil)

	err := middleware(req)
	if err != nil {
		t.Errorf("custom auth middleware failed: %v", err)
	}

	value := req.Header.Get(headerName)
	if value != headerValue {
		t.Errorf("expected '%s', got '%s'", headerValue, value)
	}
}

func TestAuthConfigWithDifferentTypes(t *testing.T) {
	tests := []struct {
		name           string
		config         AuthConfig
		expectedHeader string
		expectedValue  string
	}{
		{
			name: "Bearer token",
			config: AuthConfig{
				Type:  AuthTypeBearer,
				Token: "bearer-token",
			},
			expectedHeader: "Authorization",
			expectedValue:  "Bearer bearer-token",
		},
		{
			name: "API Key",
			config: AuthConfig{
				Type:       AuthTypeAPIKey,
				Token:      "api-key-123",
				HeaderName: "X-API-Key",
			},
			expectedHeader: "X-API-Key",
			expectedValue:  "api-key-123",
		},
		{
			name: "Basic Auth",
			config: AuthConfig{
				Type:     AuthTypeBasic,
				Username: "user",
				Password: "pass",
			},
			expectedHeader: "Authorization",
			expectedValue:  "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass")),
		},
		{
			name: "Custom Auth",
			config: AuthConfig{
				Type:        AuthTypeCustom,
				HeaderName:  "X-Custom-Token",
				HeaderValue: "custom-value",
			},
			expectedHeader: "X-Custom-Token",
			expectedValue:  "custom-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := AddAuth(tt.config)
			req := httptest.NewRequest("GET", "/test", nil)

			err := middleware(req)
			if err != nil {
				t.Errorf("auth middleware failed: %v", err)
			}

			value := req.Header.Get(tt.expectedHeader)
			if value != tt.expectedValue {
				t.Errorf("expected '%s', got '%s'", tt.expectedValue, value)
			}
		})
	}
}
