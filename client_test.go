package neuron_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/kolosys/neuron"
)

type testRequest struct {
	Name string `json:"name"`
}

type testResponse struct {
	Message string `json:"message"`
}

func TestClient_Execute(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("expected path /test, got %s", r.URL.Path)
		}

		var req testRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{Message: "hello " + req.Name})
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
	})

	route := NewRoute[testRequest, testResponse](MethodPOST, "/test")
	resp, err := Execute(client, route, testRequest{Name: "world"}, nil)

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !resp.IsSuccess() {
		t.Errorf("expected success response, got %d", resp.StatusCode)
	}

	if resp.Data.Message != "hello world" {
		t.Errorf("expected message 'hello world', got '%s'", resp.Data.Message)
	}
}

func TestClient_Retries(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL:    ts.URL,
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	})

	// Do a simple GET
	_, err := client.Get("/")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestClient_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
		Timeout: 10 * time.Millisecond,
	})

	_, err := client.Get("/")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	ce, ok := err.(ClientError)
	if !ok {
		t.Fatalf("expected ClientError, got %T", err)
	}

	if ce.Type != ErrorTypeTimeout && ce.Type != ErrorTypeNetwork {
		// http.Client returns a generic error for timeout which might be mapped to Network if not handled specifically in Do
		// But in our executeWithRetries we check context.
		t.Errorf("expected timeout or network error type, got %v", ce.Type)
	}
}

func TestClient_AdaptiveTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL:         ts.URL,
		Timeout:         100 * time.Millisecond,
		AdaptiveTimeout: true,
	})

	// GET should have 0.8 * 100ms = 80ms
	resp, err := client.Get("/")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected OK, got %d", resp.StatusCode)
	}
}

func TestClient_DoWithType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testResponse{Message: "typed"})
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{BaseURL: ts.URL})
	resp, err := DoWithType[testResponse](client, MethodGET, "/")
	if err != nil {
		t.Fatalf("DoWithType failed: %v", err)
	}

	if resp.Data.Message != "typed" {
		t.Errorf("expected 'typed', got '%s'", resp.Data.Message)
	}
}
