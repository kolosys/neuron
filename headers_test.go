package neuron_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kolosys/neuron"
)

func TestHeaderHooks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "TestAgent" {
			t.Errorf("expected User-Agent TestAgent, got %s", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Custom") != "Value" {
			t.Errorf("expected X-Custom Value, got %s", r.Header.Get("X-Custom"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		RequestHooks: []RequestHook{
			AddUserAgent("TestAgent"),
			AddHeaderSet("X-Custom", "Value"),
		},
	})

	_, err := client.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestSecurityHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Content-Type-Options") != "nosniff" {
			t.Error("nosniff header missing")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		RequestHooks: []RequestHook{
			AddSecurityHeaders(),
		},
	})

	_, err := client.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}
