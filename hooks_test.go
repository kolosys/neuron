package neuron_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/kolosys/neuron"
)

func TestClient_Hooks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Hook") != "true" {
			t.Error("request hook header missing")
		}
		w.Header().Set("X-Response-Hook", "true")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	reqHooked := false
	respHooked := false

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		RequestHooks: []RequestHook{
			func(req *http.Request) error {
				req.Header.Set("X-Request-Hook", "true")
				reqHooked = true
				return nil
			},
		},
		ResponseHooks: []ResponseHook{
			func(resp *http.Response) error {
				if resp.Header.Get("X-Response-Hook") != "true" {
					return fmt.Errorf("response hook header missing")
				}
				respHooked = true
				return nil
			},
		},
	})

	_, err := client.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !reqHooked {
		t.Error("request hook not called")
	}
	if !respHooked {
		t.Error("response hook not called")
	}
}

func TestAddResponseTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		ResponseHooks: []ResponseHook{
			AddResponseTimeout(10 * time.Millisecond),
		},
	})

	_, err := client.Get("/")
	if err == nil {
		t.Fatal("expected timeout error from hook, got nil")
	}

	ce, ok := err.(ClientError)
	if !ok {
		t.Fatalf("expected ClientError, got %T", err)
	}
	if ce.Type != ErrorTypeTimeout {
		t.Errorf("expected ErrorTypeTimeout, got %v", ce.Type)
	}
}

type invalidJSONProvider struct{}

func (p invalidJSONProvider) ContentType() string { return "application/json" }
func (p invalidJSONProvider) Body() (io.Reader, error) {
	return strings.NewReader("{invalid json"), nil
}

func TestAddValidation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		RequestHooks: []RequestHook{
			AddValidation(&JSONValidator{}),
		},
	})

	// Invalid JSON using a provider
	_, err := client.Post("/", &RequestOptions{
		Body: invalidJSONProvider{},
	})

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
