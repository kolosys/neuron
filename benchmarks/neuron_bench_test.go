package neuron_bench_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/kolosys/neuron"
)

type benchRequest struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

type benchResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func BenchmarkExecute(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req benchRequest
		json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(benchResponse{ID: req.ID, Message: "ok"})
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
	})

	route := NewRoute[benchRequest, benchResponse](MethodPOST, "/")
	reqData := benchRequest{ID: 1, Data: "benchmark"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Execute(client, route, reqData, nil)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMetricsRecord(b *testing.B) {
	collector := NewMetricsCollector()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			collector.RecordRequest()
			collector.RecordResponse(http.StatusOK, 100*time.Microsecond)
		}
	})
}

// BenchmarkClient_Post_Serialization benchmarks the client.Post path which includes body serialization.
func BenchmarkClient_Post_Serialization(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ClientOptions{
		BaseURL: ts.URL,
	})

	reqData := benchRequest{ID: 1, Data: "benchmark serialized body data test"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.Post("/", &RequestOptions{Body: reqData})
		}
	})
}
