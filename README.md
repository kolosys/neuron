# Neuron 🕸️

![GoVersion](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![stdlib + optional Ion](https://img.shields.io/badge/deps-stdlib%20%2B%20optional%20Ion-blue.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/kolosys/neuron.svg)](https://pkg.go.dev/github.com/kolosys/neuron)
[![Go Report Card](https://goreportcard.com/badge/github.com/kolosys/neuron)](https://goreportcard.com/report/github.com/kolosys/neuron)

Neuron is a type-safe HTTP client for Go. The core is the standard library; [Ion](https://github.com/kolosys/ion) is an optional circuit breaker you compose in when you want fail-fast around `HTTPClient.Do`. Built by [Kolosys](https://github.com/kolosys).

## Install

```bash
go get github.com/kolosys/neuron
```

## Execute

Typed request and response. Neuron serializes the body, runs retries (including 5xx), and materializes `TResponse`.

```go
client := neuron.NewClient(neuron.ClientOptions{
    BaseURL: "https://api.example.com",
})

route := neuron.NewRoute[CreateUserReq, User](neuron.MethodPOST, "/users")
resp, err := neuron.Execute(client, route, CreateUserReq{Name: "ada"}, nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Data.ID)
```

Convenience methods (`Get`, `Post`, `Do`, `DoWithType`) go through the same path.

## Stream

`Stream` returns `*http.Response` with **Body still open**. You `Close` it. Use this for SSE and other long-lived LLM responses — Neuron does not parse SSE.

Retries happen **only** on transport errors before any response exists. Once headers or a body are received, including 5xx, there is no second POST: replaying an LLM stream can duplicate tool calls.

`Stream` does not apply `AdaptiveTimeout`. Set `http.Client.Timeout` to something like 10 minutes, or put a deadline on the request context.

```go
client := neuron.NewClient(neuron.ClientOptions{
    BaseURL: "https://api.example.com",
    HTTPClient: &http.Client{Timeout: 10 * time.Minute},
})

resp, err := neuron.Stream(client, neuron.MethodPOST, "/v1/messages", body, &neuron.RequestOptions{
    Context: ctx,
})
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
io.Copy(os.Stdout, resp.Body)
```

## Optional Ion circuit

Leave `Circuit` nil for a plain client. Set it to wrap every `HTTPClient.Do` in `circuit.Execute`. An open circuit fails fast as `ClientError` with `Type: ErrorTypeCircuit`; the Ion `*circuit.CircuitError` is the cause.

```go
cb := circuit.New("provider",
    circuit.WithFailureThreshold(5),
    circuit.WithRecovery(circuit.RecoveryManual),
)

client := neuron.NewClient(neuron.ClientOptions{
    BaseURL: "https://api.example.com",
    Circuit: cb,
})
```

## Hooks

Request hooks run after the request is built and before `Do`. Response hooks run after a response is received. Attach them on the client or per request.

```go
client := neuron.NewClient(neuron.ClientOptions{
    BaseURL: "https://api.example.com",
    RequestHooks: []neuron.RequestHook{
        neuron.AddBearerAuth(os.Getenv("API_TOKEN")),
    },
})

_, err := client.Get("/me", &neuron.RequestOptions{
    RequestHooks: []neuron.RequestHook{
        func(req *http.Request) error {
            req.Header.Set("X-Request-ID", id)
            return nil
        },
    },
})
```

Do not attach body-consuming response hooks to `Stream`.

## Testing with mock/

`neuron/mock` provides a `RoundTripper`, rate limiter, auth, and cache fakes plus assertions.

```go
rt := mock.NewMockRoundTripper(nil)
rt.QueueResponse(mock.ResponseConfig{
    StatusCode: 200,
    Body:       []byte(`{"ok":true}`),
})

client := neuron.NewClient(neuron.ClientOptions{
    BaseURL:    "http://api.example.com",
    HTTPClient: &http.Client{Transport: rt},
})

_, err := client.Get("/health")
mock.AssertRequestCount(t, rt, 1)
mock.AssertRequestMethod(t, rt, 0, "GET")
```

## License

MIT
