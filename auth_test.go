package neuron_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kolosys/neuron"
)

func TestAuthHooks(t *testing.T) {
	tests := []struct {
		name       string
		hook       RequestHook
		wantHeader string
		wantValue  string
	}{
		{
			name:       "Bearer Auth",
			hook:       AddBearerAuth("token123"),
			wantHeader: "Authorization",
			wantValue:  "Bearer token123",
		},
		{
			name:       "API Key Auth",
			hook:       AddAPIKeyAuth("key123", "X-API-KEY"),
			wantHeader: "X-Api-Key",
			wantValue:  "key123",
		},
		{
			name:       "Basic Auth",
			hook:       AddBasicAuth("user", "pass"),
			wantHeader: "Authorization",
			wantValue:  "Basic dXNlcjpwYXNz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get(tt.wantHeader) != tt.wantValue {
					t.Errorf("expected header %s=%s, got %s", tt.wantHeader, tt.wantValue, r.Header.Get(tt.wantHeader))
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			client := NewClient(ClientOptions{
				BaseURL:      ts.URL,
				RequestHooks: []RequestHook{tt.hook},
			})

			_, err := client.Get("/")
			if err != nil {
				t.Fatalf("Get failed: %v", err)
			}
		})
	}
}
