# neuron API

Complete API documentation for the neuron package.

**Import Path:** `github.com/kolosys/neuron`

## Package Documentation

## Types

### AuthConfig

AuthConfig configures authentication middleware

#### Example Usage

```go
// Create a new AuthConfig
authconfig := AuthConfig{
    Type: AuthType{},
    Token: "example",
    Username: "example",
    Password: "example",
    HeaderName: "example",
    HeaderValue: "example",
}
```

#### Type Definition

```go
type AuthConfig struct {
    Type AuthType
    Token string
    Username string
    Password string
    HeaderName string
    HeaderValue string
}
```

### Fields

| Field       | Type       | Description |
| ----------- | ---------- | ----------- |
| Type        | `AuthType` |             |
| Token       | `string`   |             |
| Username    | `string`   |             |
| Password    | `string`   |             |
| HeaderName  | `string`   |             |
| HeaderValue | `string`   |             |

### AuthProvider

AuthProvider interface for authentication

#### Example Usage

```go
// Example implementation of AuthProvider
type MyAuthProvider struct {
    // Add your fields here
}

func (m MyAuthProvider) GetToken(param1 context.Context) string {
    // Implement your logic here
    return
}

func (m MyAuthProvider) GetAuthHeader(param1 string) string {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type AuthProvider interface {
    GetToken(ctx context.Context) (string, error)
    GetAuthHeader(token string) string
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### AuthType

AuthType represents the type of authentication

#### Example Usage

```go
// Example usage of AuthType
var value AuthType
// Initialize with appropriate value
```

#### Type Definition

```go
type AuthType int
```

### BackoffStrategy

BackoffStrategy defines how to handle backoff

#### Example Usage

```go
// Example usage of BackoffStrategy
var value BackoffStrategy
// Initialize with appropriate value
```

#### Type Definition

```go
type BackoffStrategy int
```

### BodyProvider

BodyProvider allows custom body serialization

#### Example Usage

```go
// Example implementation of BodyProvider
type MyBodyProvider struct {
    // Add your fields here
}

func (m MyBodyProvider) ContentType() string {
    // Implement your logic here
    return
}

func (m MyBodyProvider) Body() io.Reader {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type BodyProvider interface {
    ContentType() string
    Body() (io.Reader, error)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Cache

Cache interface for caching middleware

#### Example Usage

```go
// Example implementation of Cache
type MyCache struct {
    // Add your fields here
}

func (m MyCache) Get(param1 string) *CacheEntry {
    // Implement your logic here
    return
}

func (m MyCache) Set(param1 string, param2 CacheEntry)  {
    // Implement your logic here
    return
}

func (m MyCache) Delete(param1 string)  {
    // Implement your logic here
    return
}

func (m MyCache) Clear()  {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Cache interface {
    Get(key string) (*CacheEntry, bool)
    Set(key string, entry CacheEntry)
    Delete(key string)
    Clear()
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### CacheEntry

CacheEntry represents a cached response

#### Example Usage

```go
// Create a new CacheEntry
cacheentry := CacheEntry{
    Data: [],
    Headers: /* value */,
    StatusCode: 42,
    Timestamp: /* value */,
    TTL: /* value */,
}
```

#### Type Definition

```go
type CacheEntry struct {
    Data []byte
    Headers http.Header
    StatusCode int
    Timestamp time.Time
    TTL time.Duration
}
```

### Fields

| Field      | Type            | Description |
| ---------- | --------------- | ----------- |
| Data       | `[]byte`        |             |
| Headers    | `http.Header`   |             |
| StatusCode | `int`           |             |
| Timestamp  | `time.Time`     |             |
| TTL        | `time.Duration` |             |

### CircuitBreaker

CircuitBreaker interface for middleware integration

#### Example Usage

```go
// Example implementation of CircuitBreaker
type MyCircuitBreaker struct {
    // Add your fields here
}

func (m MyCircuitBreaker) AllowRequest() bool {
    // Implement your logic here
    return
}

func (m MyCircuitBreaker) RecordSuccess()  {
    // Implement your logic here
    return
}

func (m MyCircuitBreaker) RecordFailure()  {
    // Implement your logic here
    return
}

func (m MyCircuitBreaker) GetState() CircuitBreakerState {
    // Implement your logic here
    return
}

func (m MyCircuitBreaker) Close() error {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type CircuitBreaker interface {
    AllowRequest() bool
    RecordSuccess()
    RecordFailure()
    GetState() CircuitBreakerState
    Close() error
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### CircuitBreakerConfig

CircuitBreakerConfig defines circuit breaker behavior for HTTP clients

#### Example Usage

```go
// Create a new CircuitBreakerConfig
circuitbreakerconfig := CircuitBreakerConfig{
    Enabled: true,
    PerRouteCircuitBreakers: true,
    FailureThreshold: 42,
    RecoveryTimeout: /* value */,
    HalfOpenMaxRequests: 42,
    SuccessThreshold: 42,
}
```

#### Type Definition

```go
type CircuitBreakerConfig struct {
    Enabled bool
    PerRouteCircuitBreakers bool
    FailureThreshold int
    RecoveryTimeout time.Duration
    HalfOpenMaxRequests int
    SuccessThreshold int
}
```

### Fields

| Field                   | Type            | Description                                                   |
| ----------------------- | --------------- | ------------------------------------------------------------- |
| Enabled                 | `bool`          | Enable circuit breaker functionality                          |
| PerRouteCircuitBreakers | `bool`          | Per-route circuit breakers (vs single global circuit breaker) |
| FailureThreshold        | `int`           | Circuit breaker options                                       |
| RecoveryTimeout         | `time.Duration` |                                                               |
| HalfOpenMaxRequests     | `int`           |                                                               |
| SuccessThreshold        | `int`           |                                                               |

### Constructor Functions

### DefaultCircuitBreakerConfig

DefaultCircuitBreakerConfig returns sensible defaults for HTTP clients

```go
func DefaultCircuitBreakerConfig() CircuitBreakerConfig
```

**Parameters:**
None

**Returns:**

- CircuitBreakerConfig

### CircuitBreakerState

CircuitBreakerState represents the state of a circuit breaker

#### Example Usage

```go
// Example usage of CircuitBreakerState
var value CircuitBreakerState
// Initialize with appropriate value
```

#### Type Definition

```go
type CircuitBreakerState string
```

### Client

Client provides a type-safe HTTP client with rate limiting, queuing, and circuit breaking

#### Example Usage

```go
// Create a new Client
client := Client{
    Options: ClientOptions{},
    Metrics: RequestMetrics{},
}
```

#### Type Definition

```go
type Client struct {
    Options ClientOptions
    Metrics RequestMetrics
}
```

### Fields

| Field   | Type             | Description      |
| ------- | ---------------- | ---------------- |
| Options | `ClientOptions`  |                  |
| Metrics | `RequestMetrics` | Metrics tracking |

### Constructor Functions

### NewClient

NewClient creates a new type-safe HTTP client

```go
func NewClient(options ClientOptions) *Client
```

**Parameters:**

- `options` (ClientOptions)

**Returns:**

- \*Client

## Methods

### GetMetrics

GetMetrics returns current client metrics

```go
func (*Client) GetMetrics() RequestMetrics
```

**Parameters:**
None

**Returns:**

- RequestMetrics

### Shutdown

Shutdown gracefully shuts down the client

```go
func (*Client) Shutdown(timeout time.Duration) error
```

**Parameters:**

- `timeout` (time.Duration)

**Returns:**

- error

### ClientError

Error types for type-safe error handling

#### Example Usage

```go
// Create a new ClientError
clienterror := ClientError{
    Type: ErrorType{},
    Message: "example",
    StatusCode: 42,
    Route: "example",
    Method: "example",
    URL: "example",
    Attempt: 42,
    Timestamp: /* value */,
    Cause: error{},
    Context: RequestContext{},
}
```

#### Type Definition

```go
type ClientError struct {
    Type ErrorType
    Message string
    StatusCode int
    Route string
    Method string
    URL string
    Attempt int
    Timestamp time.Time
    Cause error
    Context RequestContext
}
```

### Fields

| Field      | Type             | Description          |
| ---------- | ---------------- | -------------------- |
| Type       | `ErrorType`      |                      |
| Message    | `string`         |                      |
| StatusCode | `int`            |                      |
| Route      | `string`         |                      |
| Method     | `string`         | HTTP method          |
| URL        | `string`         | Full URL             |
| Attempt    | `int`            | Retry attempt number |
| Timestamp  | `time.Time`      | When error occurred  |
| Cause      | `error`          |                      |
| Context    | `RequestContext` |                      |

## Methods

### Error

```go
func (ClientError) Error() string
```

**Parameters:**
None

**Returns:**

- string

### Unwrap

```go
func (ClientError) Unwrap() error
```

**Parameters:**
None

**Returns:**

- error

### WithContext

WithContext adds request context information to the error

```go
func (ClientError) WithContext(req *http.Request, attempt int) ClientError
```

**Parameters:**

- `req` (\*http.Request)
- `attempt` (int)

**Returns:**

- ClientError

### ClientOptions

ClientOptions configures the HTTP client behavior

#### Example Usage

```go
// Create a new ClientOptions
clientoptions := ClientOptions{
    BaseURL: "example",
    UserAgent: "example",
    Headers: /* value */,
    Timeout: /* value */,
    GlobalRateLimit: RateLimiter{},
    PerRouteRateLimit: true,
    RateLimitConfig: RateLimitConfig{},
    CircuitBreakerConfig: CircuitBreakerConfig{},
    MaxRetries: 42,
    RetryDelay: /* value */,
    RetryMultiplier: 3.14,
    QueueTimeout: /* value */,
    MaxQueueSize: 42,
    RequestMiddleware: [],
    ResponseMiddleware: [],
    HTTPClient: &/* value */{},
    SweepInterval: /* value */,
    SweepEnabled: true,
}
```

#### Type Definition

```go
type ClientOptions struct {
    BaseURL string
    UserAgent string
    Headers http.Header
    Timeout time.Duration
    GlobalRateLimit RateLimiter
    PerRouteRateLimit bool
    RateLimitConfig RateLimitConfig
    CircuitBreakerConfig CircuitBreakerConfig
    MaxRetries int
    RetryDelay time.Duration
    RetryMultiplier float64
    QueueTimeout time.Duration
    MaxQueueSize int
    RequestMiddleware []RequestMiddleware
    ResponseMiddleware []ResponseMiddleware
    HTTPClient *http.Client
    SweepInterval time.Duration
    SweepEnabled bool
}
```

### Fields

| Field                | Type                   | Description                   |
| -------------------- | ---------------------- | ----------------------------- |
| BaseURL              | `string`               | Base configuration            |
| UserAgent            | `string`               |                               |
| Headers              | `http.Header`          |                               |
| Timeout              | `time.Duration`        |                               |
| GlobalRateLimit      | `RateLimiter`          | Rate limiting                 |
| PerRouteRateLimit    | `bool`                 |                               |
| RateLimitConfig      | `RateLimitConfig`      |                               |
| CircuitBreakerConfig | `CircuitBreakerConfig` | Circuit breaker configuration |
| MaxRetries           | `int`                  | Request handling              |
| RetryDelay           | `time.Duration`        |                               |
| RetryMultiplier      | `float64`              |                               |
| QueueTimeout         | `time.Duration`        | Queue management              |
| MaxQueueSize         | `int`                  |                               |
| RequestMiddleware    | `[]RequestMiddleware`  | Middleware                    |
| ResponseMiddleware   | `[]ResponseMiddleware` |                               |
| HTTPClient           | `*http.Client`         | HTTP client                   |
| SweepInterval        | `time.Duration`        | Sweeping configuration        |
| SweepEnabled         | `bool`                 |                               |

### CompressionConfig

CompressionConfig configures compression middleware

#### Example Usage

```go
// Create a new CompressionConfig
compressionconfig := CompressionConfig{
    Type: CompressionType{},
    Level: 42,
    MinSize: 42,
    ContentTypes: [],
}
```

#### Type Definition

```go
type CompressionConfig struct {
    Type CompressionType
    Level int
    MinSize int
    ContentTypes []string
}
```

### Fields

| Field        | Type              | Description |
| ------------ | ----------------- | ----------- |
| Type         | `CompressionType` |             |
| Level        | `int`             |             |
| MinSize      | `int`             |             |
| ContentTypes | `[]string`        |             |

### Constructor Functions

### DefaultCompressionConfig

DefaultCompressionConfig returns a default compression configuration

```go
func DefaultCompressionConfig() CompressionConfig
```

**Parameters:**
None

**Returns:**

- CompressionConfig

### CompressionType

CompressionType represents the type of compression

#### Example Usage

```go
// Example usage of CompressionType
var value CompressionType
// Initialize with appropriate value
```

#### Type Definition

```go
type CompressionType int
```

### Deserializable

Deserializable represents types that can be deserialized from responses

#### Example Usage

```go
// Example implementation of Deserializable
type MyDeserializable struct {
    // Add your fields here
}

func (m MyDeserializable) UnmarshalJSON(param1 []byte) error {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Deserializable interface {
    UnmarshalJSON(data []byte) error
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### EmptyRequest

EmptyRequest represents requests with no body

#### Example Usage

```go
// Create a new EmptyRequest
emptyrequest := EmptyRequest{

}
```

#### Type Definition

```go
type EmptyRequest struct {
}
```

### EmptyResponse

EmptyResponse represents responses with no body

#### Example Usage

```go
// Create a new EmptyResponse
emptyresponse := EmptyResponse{

}
```

#### Type Definition

```go
type EmptyResponse struct {
}
```

### ErrorType

_No documentation available_

#### Example Usage

```go
// Example usage of ErrorType
var value ErrorType
// Initialize with appropriate value
```

#### Type Definition

```go
type ErrorType int
```

### HTTPMethod

HTTPMethod represents supported HTTP methods

#### Example Usage

```go
// Example usage of HTTPMethod
var value HTTPMethod
// Initialize with appropriate value
```

#### Type Definition

```go
type HTTPMethod string
```

### HeaderConfig

HeaderConfig configures header middleware

#### Example Usage

```go
// Create a new HeaderConfig
headerconfig := HeaderConfig{
    UserAgent: "example",
    CustomHeaders: map[],
    RemoveHeaders: [],
    AddHeaders: map[],
}
```

#### Type Definition

```go
type HeaderConfig struct {
    UserAgent string
    CustomHeaders map[string]string
    RemoveHeaders []string
    AddHeaders map[string]string
}
```

### Fields

| Field         | Type                | Description |
| ------------- | ------------------- | ----------- |
| UserAgent     | `string`            |             |
| CustomHeaders | `map[string]string` |             |
| RemoveHeaders | `[]string`          |             |
| AddHeaders    | `map[string]string` |             |

### Constructor Functions

### DefaultHeaderConfig

DefaultHeaderConfig returns a default header configuration

```go
func DefaultHeaderConfig() HeaderConfig
```

**Parameters:**
None

**Returns:**

- HeaderConfig

### InMemoryCache

InMemoryCache provides a simple in-memory cache

#### Example Usage

```go
// Create a new InMemoryCache
inmemorycache := InMemoryCache{

}
```

#### Type Definition

```go
type InMemoryCache struct {
}
```

### Constructor Functions

### NewInMemoryCache

```go
func NewInMemoryCache() *InMemoryCache
```

**Parameters:**
None

**Returns:**

- \*InMemoryCache

## Methods

### Clear

```go
func (*InMemoryCache) Clear()
```

**Parameters:**
None

**Returns:**
None

### Delete

```go
func (*InMemoryCache) Delete(key string)
```

**Parameters:**

- `key` (string)

**Returns:**
None

### Get

```go
func (*InMemoryCache) Get(key string) (*CacheEntry, bool)
```

**Parameters:**

- `key` (string)

**Returns:**

- \*CacheEntry
- bool

### Set

```go
func (*InMemoryCache) Set(key string, entry CacheEntry)
```

**Parameters:**

- `key` (string)
- `entry` (CacheEntry)

**Returns:**
None

### JSONValidator

JSONValidator provides JSON schema validation

#### Example Usage

```go
// Create a new JSONValidator
jsonvalidator := JSONValidator{

}
```

#### Type Definition

```go
type JSONValidator struct {
}
```

## Methods

### Validate

```go
func (*JSONValidator) Validate(data []byte, contentType string) error
```

**Parameters:**

- `data` ([]byte)
- `contentType` (string)

**Returns:**

- error

### LogEntry

LogEntry represents a log entry

#### Example Usage

```go
// Create a new LogEntry
logentry := LogEntry{
    Method: "example",
    URL: "example",
    StatusCode: 42,
    Headers: /* value */,
    Duration: /* value */,
    Timestamp: /* value */,
    Body: [],
    Error: error{},
}
```

#### Type Definition

```go
type LogEntry struct {
    Method string
    URL string
    StatusCode int
    Headers http.Header
    Duration time.Duration
    Timestamp time.Time
    Body []byte
    Error error
}
```

### Fields

| Field      | Type            | Description |
| ---------- | --------------- | ----------- |
| Method     | `string`        |             |
| URL        | `string`        |             |
| StatusCode | `int`           |             |
| Headers    | `http.Header`   |             |
| Duration   | `time.Duration` |             |
| Timestamp  | `time.Time`     |             |
| Body       | `[]byte`        |             |
| Error      | `error`         |             |

### LogLevel

LogLevel represents the logging level

#### Example Usage

```go
// Example usage of LogLevel
var value LogLevel
// Initialize with appropriate value
```

#### Type Definition

```go
type LogLevel int
```

### Logger

Interfaces for extensibility Logger interface for logging middleware

#### Example Usage

```go
// Example implementation of Logger
type MyLogger struct {
    // Add your fields here
}

func (m MyLogger) LogRequest(param1 LogEntry)  {
    // Implement your logic here
    return
}

func (m MyLogger) LogResponse(param1 LogEntry)  {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Logger interface {
    LogRequest(entry LogEntry)
    LogResponse(entry LogEntry)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### LoggingConfig

LoggingConfig configures the logging middleware

#### Example Usage

```go
// Create a new LoggingConfig
loggingconfig := LoggingConfig{
    Level: LogLevel{},
    IncludeBody: true,
    MaxBodySize: 42,
    Logger: &/* value */{},
}
```

#### Type Definition

```go
type LoggingConfig struct {
    Level LogLevel
    IncludeBody bool
    MaxBodySize int
    Logger *log.Logger
}
```

### Fields

| Field       | Type          | Description |
| ----------- | ------------- | ----------- |
| Level       | `LogLevel`    |             |
| IncludeBody | `bool`        |             |
| MaxBodySize | `int`         |             |
| Logger      | `*log.Logger` |             |

### Constructor Functions

### DefaultLoggingConfig

DefaultLoggingConfig returns a default logging configuration

```go
func DefaultLoggingConfig() LoggingConfig
```

**Parameters:**
None

**Returns:**

- LoggingConfig

### LoggingContextKey

_No documentation available_

#### Example Usage

```go
// Example usage of LoggingContextKey
var value LoggingContextKey
// Initialize with appropriate value
```

#### Type Definition

```go
type LoggingContextKey string
```

### MetricsCollector

MetricsCollector collects and stores metrics

#### Example Usage

```go
// Create a new MetricsCollector
metricscollector := MetricsCollector{
    RequestCount: 42,
    ResponseCount: 42,
    ErrorCount: 42,
    TotalDuration: /* value */,
    MinDuration: /* value */,
    MaxDuration: /* value */,
    StatusCodeCounts: map[],
    StartTime: /* value */,
}
```

#### Type Definition

```go
type MetricsCollector struct {
    RequestCount int64
    ResponseCount int64
    ErrorCount int64
    TotalDuration time.Duration
    MinDuration time.Duration
    MaxDuration time.Duration
    StatusCodeCounts map[int]int64
    StartTime time.Time
}
```

### Fields

| Field            | Type            | Description                       |
| ---------------- | --------------- | --------------------------------- |
| RequestCount     | `int64`         | Request metrics                   |
| ResponseCount    | `int64`         |                                   |
| ErrorCount       | `int64`         |                                   |
| TotalDuration    | `time.Duration` | Duration metrics                  |
| MinDuration      | `time.Duration` |                                   |
| MaxDuration      | `time.Duration` |                                   |
| StatusCodeCounts | `map[int]int64` | Status code metrics               |
| StartTime        | `time.Time`     | Start time for uptime calculation |

### Constructor Functions

### NewMetricsCollector

NewMetricsCollector creates a new metrics collector

```go
func NewMetricsCollector() *MetricsCollector
```

**Parameters:**
None

**Returns:**

- \*MetricsCollector

## Methods

### GetMetrics

GetMetrics returns current metrics

```go
func (*Client) GetMetrics() RequestMetrics
```

**Parameters:**
None

**Returns:**

- RequestMetrics

### RecordRequest

RecordRequest records a request metric

```go
func (*MetricsCollector) RecordRequest()
```

**Parameters:**
None

**Returns:**
None

### RecordResponse

RecordResponse records a response metric

```go
func (*MetricsCollector) RecordResponse(statusCode int, duration time.Duration)
```

**Parameters:**

- `statusCode` (int)
- `duration` (time.Duration)

**Returns:**
None

### MetricsSnapshot

MetricsSnapshot represents a snapshot of metrics at a point in time

#### Example Usage

```go
// Create a new MetricsSnapshot
metricssnapshot := MetricsSnapshot{
    RequestCount: 42,
    ResponseCount: 42,
    ErrorCount: 42,
    AverageDuration: /* value */,
    MinDuration: /* value */,
    MaxDuration: /* value */,
    StatusCodeCounts: map[],
    Uptime: /* value */,
}
```

#### Type Definition

```go
type MetricsSnapshot struct {
    RequestCount int64
    ResponseCount int64
    ErrorCount int64
    AverageDuration time.Duration
    MinDuration time.Duration
    MaxDuration time.Duration
    StatusCodeCounts map[int]int64
    Uptime time.Duration
}
```

### Fields

| Field            | Type            | Description |
| ---------------- | --------------- | ----------- |
| RequestCount     | `int64`         |             |
| ResponseCount    | `int64`         |             |
| ErrorCount       | `int64`         |             |
| AverageDuration  | `time.Duration` |             |
| MinDuration      | `time.Duration` |             |
| MaxDuration      | `time.Duration` |             |
| StatusCodeCounts | `map[int]int64` |             |
| Uptime           | `time.Duration` |             |

## Methods

### ErrorRate

ErrorRate returns the error rate as a percentage

```go
func (*MetricsSnapshot) ErrorRate() float64
```

**Parameters:**
None

**Returns:**

- float64

### RequestsPerSecond

RequestsPerSecond returns the average requests per second

```go
func (*MetricsSnapshot) RequestsPerSecond() float64
```

**Parameters:**
None

**Returns:**

- float64

### MiddlewareChain

MiddlewareChain manages a chain of middleware functions

#### Example Usage

```go
// Create a new MiddlewareChain
middlewarechain := MiddlewareChain{

}
```

#### Type Definition

```go
type MiddlewareChain struct {
}
```

### Constructor Functions

### NewMiddlewareChain

NewMiddlewareChain creates a new middleware chain

```go
func NewMiddlewareChain() *MiddlewareChain
```

**Parameters:**
None

**Returns:**

- \*MiddlewareChain

## Methods

### AddRequestMiddleware

AddRequestMiddleware adds a request middleware to the chain

```go
func (*MiddlewareChain) AddRequestMiddleware(middleware RequestMiddleware) *MiddlewareChain
```

**Parameters:**

- `middleware` (RequestMiddleware)

**Returns:**

- \*MiddlewareChain

### AddResponseMiddleware

AddResponseMiddleware adds a response middleware to the chain

```go
func (*MiddlewareChain) AddResponseMiddleware(middleware ResponseMiddleware) *MiddlewareChain
```

**Parameters:**

- `middleware` (ResponseMiddleware)

**Returns:**

- \*MiddlewareChain

### ApplyRequestMiddleware

ApplyRequestMiddleware applies all request middleware in order

```go
func (*MiddlewareChain) ApplyRequestMiddleware(req *http.Request) error
```

**Parameters:**

- `req` (\*http.Request)

**Returns:**

- error

### ApplyResponseMiddleware

ApplyResponseMiddleware applies all response middleware in order

```go
func (*MiddlewareChain) ApplyResponseMiddleware(resp *http.Response) error
```

**Parameters:**

- `resp` (\*http.Response)

**Returns:**

- error

### QueuedRequest

QueuedRequest represents a request waiting in queue

#### Example Usage

```go
// Create a new QueuedRequest
queuedrequest := QueuedRequest{
    Request: &/* value */{},
    ResponseCh: /* value */,
    Context: /* value */,
    Retries: 42,
    EnqueueTime: /* value */,
}
```

#### Type Definition

```go
type QueuedRequest struct {
    Request *http.Request
    ResponseCh chan *QueuedResponse
    Context context.Context
    Retries int
    EnqueueTime time.Time
}
```

### Fields

| Field       | Type                   | Description |
| ----------- | ---------------------- | ----------- |
| Request     | `*http.Request`        |             |
| ResponseCh  | `chan *QueuedResponse` |             |
| Context     | `context.Context`      |             |
| Retries     | `int`                  |             |
| EnqueueTime | `time.Time`            |             |

### QueuedResponse

QueuedResponse represents the result of a queued request

#### Example Usage

```go
// Create a new QueuedResponse
queuedresponse := QueuedResponse{
    Response: &/* value */{},
    Error: error{},
}
```

#### Type Definition

```go
type QueuedResponse struct {
    Response *http.Response
    Error error
}
```

### Fields

| Field    | Type             | Description |
| -------- | ---------------- | ----------- |
| Response | `*http.Response` |             |
| Error    | `error`          |             |

### RateLimitConfig

RateLimitConfig defines rate limiting behavior

#### Example Usage

```go
// Create a new RateLimitConfig
ratelimitconfig := RateLimitConfig{
    GlobalRequestsPerSecond: 42,
    GlobalBurstSize: 42,
    RouteRequestsPerSecond: 42,
    RouteBurstSize: 42,
    RespectDiscordHeaders: true,
    BackoffStrategy: BackoffStrategy{},
    QueueOnRateLimit: true,
    RateLimitTimeout: /* value */,
}
```

#### Type Definition

```go
type RateLimitConfig struct {
    GlobalRequestsPerSecond int
    GlobalBurstSize int
    RouteRequestsPerSecond int
    RouteBurstSize int
    RespectDiscordHeaders bool
    BackoffStrategy BackoffStrategy
    QueueOnRateLimit bool
    RateLimitTimeout time.Duration
}
```

### Fields

| Field                   | Type              | Description                   |
| ----------------------- | ----------------- | ----------------------------- |
| GlobalRequestsPerSecond | `int`             | Global rate limiting          |
| GlobalBurstSize         | `int`             |                               |
| RouteRequestsPerSecond  | `int`             | Per-route rate limiting       |
| RouteBurstSize          | `int`             |                               |
| RespectDiscordHeaders   | `bool`            | Rate limit detection          |
| BackoffStrategy         | `BackoffStrategy` |                               |
| QueueOnRateLimit        | `bool`            | Queue behavior on rate limits |
| RateLimitTimeout        | `time.Duration`   |                               |

### RateLimitInfo

RateLimitInfo contains information about current rate limit status

#### Example Usage

```go
// Create a new RateLimitInfo
ratelimitinfo := RateLimitInfo{
    RouteID: "example",
    Bucket: "example",
    Limit: 42,
    Remaining: 42,
    ResetAfter: /* value */,
    RetryAfter: /* value */,
    Global: true,
}
```

#### Type Definition

```go
type RateLimitInfo struct {
    RouteID string
    Bucket string
    Limit int
    Remaining int
    ResetAfter time.Duration
    RetryAfter time.Duration
    Global bool
}
```

### Fields

| Field      | Type            | Description |
| ---------- | --------------- | ----------- |
| RouteID    | `string`        |             |
| Bucket     | `string`        |             |
| Limit      | `int`           |             |
| Remaining  | `int`           |             |
| ResetAfter | `time.Duration` |             |
| RetryAfter | `time.Duration` |             |
| Global     | `bool`          |             |

### RateLimitInfoProvider

RateLimitInfoProvider provides rate limit information

#### Example Usage

```go
// Example implementation of RateLimitInfoProvider
type MyRateLimitInfoProvider struct {
    // Add your fields here
}

func (m MyRateLimitInfoProvider) GetRateLimitInfo(param1 string) *RateLimitInfo {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type RateLimitInfoProvider interface {
    GetRateLimitInfo(path string) *RateLimitInfo
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### RateLimiter

Middleware interfaces for external library integration RateLimiter interface for middleware integration

#### Example Usage

```go
// Example implementation of RateLimiter
type MyRateLimiter struct {
    // Add your fields here
}

func (m MyRateLimiter) WaitN(param1 context.Context, param2 int) error {
    // Implement your logic here
    return
}

func (m MyRateLimiter) AllowN(param1 time.Time, param2 int) bool {
    // Implement your logic here
    return
}

func (m MyRateLimiter) Tokens() float64 {
    // Implement your logic here
    return
}

func (m MyRateLimiter) Burst() int {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type RateLimiter interface {
    WaitN(ctx context.Context, n int) error
    AllowN(now time.Time, n int) bool
    Tokens() float64
    Burst() int
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### RequestContext

RequestContext provides metadata about the current request

#### Example Usage

```go
// Create a new RequestContext
requestcontext := RequestContext{
    RouteID: "example",
    Attempt: 42,
    QueueTime: /* value */,
    StartTime: /* value */,
    Metadata: map[],
}
```

#### Type Definition

```go
type RequestContext struct {
    RouteID string
    Attempt int
    QueueTime time.Duration
    StartTime time.Time
    Metadata map[string]any
}
```

### Fields

| Field     | Type             | Description |
| --------- | ---------------- | ----------- |
| RouteID   | `string`         |             |
| Attempt   | `int`            |             |
| QueueTime | `time.Duration`  |             |
| StartTime | `time.Time`      |             |
| Metadata  | `map[string]any` |             |

### RequestIDConfig

RequestIDConfig configures request ID generation

#### Example Usage

```go
// Create a new RequestIDConfig
requestidconfig := RequestIDConfig{
    Generator: RequestIDGenerator{},
    HeaderName: "example",
    ContextKey: RequestIDContextKey{},
    Propagate: true,
}
```

#### Type Definition

```go
type RequestIDConfig struct {
    Generator RequestIDGenerator
    HeaderName string
    ContextKey RequestIDContextKey
    Propagate bool
}
```

### Fields

| Field      | Type                  | Description |
| ---------- | --------------------- | ----------- |
| Generator  | `RequestIDGenerator`  |             |
| HeaderName | `string`              |             |
| ContextKey | `RequestIDContextKey` |             |
| Propagate  | `bool`                |             |

### Constructor Functions

### DefaultRequestIDConfig

DefaultRequestIDConfig returns a default request ID configuration

```go
func DefaultRequestIDConfig() RequestIDConfig
```

**Parameters:**
None

**Returns:**

- RequestIDConfig

### RequestIDContextKey

_No documentation available_

#### Example Usage

```go
// Example usage of RequestIDContextKey
var value RequestIDContextKey
// Initialize with appropriate value
```

#### Type Definition

```go
type RequestIDContextKey string
```

### RequestIDGenerator

RequestIDGenerator generates unique request IDs

#### Example Usage

```go
// Example implementation of RequestIDGenerator
type MyRequestIDGenerator struct {
    // Add your fields here
}

func (m MyRequestIDGenerator) Generate() string {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type RequestIDGenerator interface {
    Generate() string
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### RequestMetrics

RequestMetrics provides insights into request performance

#### Example Usage

```go
// Create a new RequestMetrics
requestmetrics := RequestMetrics{
    TotalRequests: 42,
    SuccessfulRequests: 42,
    FailedRequests: 42,
    QueuedRequests: 42,
    AverageQueueTime: /* value */,
    AverageResponseTime: /* value */,
    RateLimitHits: 42,
}
```

#### Type Definition

```go
type RequestMetrics struct {
    TotalRequests int64
    SuccessfulRequests int64
    FailedRequests int64
    QueuedRequests int64
    AverageQueueTime time.Duration
    AverageResponseTime time.Duration
    RateLimitHits int64
}
```

### Fields

| Field               | Type            | Description |
| ------------------- | --------------- | ----------- |
| TotalRequests       | `int64`         |             |
| SuccessfulRequests  | `int64`         |             |
| FailedRequests      | `int64`         |             |
| QueuedRequests      | `int64`         |             |
| AverageQueueTime    | `time.Duration` |             |
| AverageResponseTime | `time.Duration` |             |
| RateLimitHits       | `int64`         |             |

### RequestMiddleware

RequestMiddleware processes requests before they are sent

#### Example Usage

```go
// Example usage of RequestMiddleware
var value RequestMiddleware
// Initialize with appropriate value
```

#### Type Definition

```go
type RequestMiddleware func(req *http.Request) error
```

### Constructor Functions

### AddAPIKeyAuth

AddAPIKeyAuth creates an API key authentication middleware

```go
func AddAPIKeyAuth(apiKey, headerName string) RequestMiddleware
```

**Parameters:**

- `apiKey` (string)
- `headerName` (string)

**Returns:**

- RequestMiddleware

### AddAPIKeyHeaderAuth

AddAPIKeyHeaderAuth creates an API key authentication middleware with common header names

```go
func AddAPIKeyHeaderAuth(apiKey string) RequestMiddleware
```

**Parameters:**

- `apiKey` (string)

**Returns:**

- RequestMiddleware

### AddAccept

AddAccept creates an accept header middleware

```go
func AddAccept(accept string) RequestMiddleware
```

**Parameters:**

- `accept` (string)

**Returns:**

- RequestMiddleware

### AddAcceptAll

AddAcceptAll creates an accept all middleware

```go
func AddAcceptAll() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddAcceptJSON

AddAcceptJSON creates an accept JSON middleware

```go
func AddAcceptJSON() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddAcceptXML

AddAcceptXML creates an accept XML middleware

```go
func AddAcceptXML() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddAdaptiveTimeout

AddAdaptiveTimeout creates an adaptive timeout middleware

```go
func AddAdaptiveTimeout(baseTimeout time.Duration, multiplier float64) RequestMiddleware
```

**Parameters:**

- `baseTimeout` (time.Duration)
- `multiplier` (float64)

**Returns:**

- RequestMiddleware

### AddAuth

AddAuth creates an authentication middleware based on configuration

```go
func AddAuth(config AuthConfig) RequestMiddleware
```

**Parameters:**

- `config` (AuthConfig)

**Returns:**

- RequestMiddleware

### AddAuthFromContext

AddAuthFromContext creates an authentication middleware that gets credentials from context

```go
func AddAuthFromContext(headerName string, contextKey string) RequestMiddleware
```

**Parameters:**

- `headerName` (string)
- `contextKey` (string)

**Returns:**

- RequestMiddleware

### AddAuthentication

AddAuthentication adds authentication headers

```go
func AddAuthentication(authProvider AuthProvider) RequestMiddleware
```

**Parameters:**

- `authProvider` (AuthProvider)

**Returns:**

- RequestMiddleware

### AddAutoCompression

AddAutoCompression creates an automatic compression middleware

```go
func AddAutoCompression() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddBasicAuth

AddBasicAuth creates a Basic authentication middleware

```go
func AddBasicAuth(username, password string) RequestMiddleware
```

**Parameters:**

- `username` (string)
- `password` (string)

**Returns:**

- RequestMiddleware

### AddBearerAuth

AddBearerAuth creates a Bearer token authentication middleware

```go
func AddBearerAuth(token string) RequestMiddleware
```

**Parameters:**

- `token` (string)

**Returns:**

- RequestMiddleware

### AddCORSHeaders

AddCORSHeaders creates a CORS headers middleware

```go
func AddCORSHeaders() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddCacheControl

AddCacheControl creates a cache control middleware

```go
func AddCacheControl(cacheControl string) RequestMiddleware
```

**Parameters:**

- `cacheControl` (string)

**Returns:**

- RequestMiddleware

### AddCircuitBreaker

AddCircuitBreaker implements circuit breaker pattern

```go
func AddCircuitBreaker(circuitBreaker CircuitBreaker) RequestMiddleware
```

**Parameters:**

- `circuitBreaker` (CircuitBreaker)

**Returns:**

- RequestMiddleware

### AddCompression

AddCompression creates a request compression middleware

```go
func AddCompression(config CompressionConfig) RequestMiddleware
```

**Parameters:**

- `config` (CompressionConfig)

**Returns:**

- RequestMiddleware

### AddCompressionLevel

AddCompressionLevel creates a compression level middleware

```go
func AddCompressionLevel(level int) RequestMiddleware
```

**Parameters:**

- `level` (int)

**Returns:**

- RequestMiddleware

### AddConditionalAuth

AddConditionalAuth creates a conditional authentication middleware

```go
func AddConditionalAuth(condition func(*http.Request) bool, authMiddleware RequestMiddleware) RequestMiddleware
```

**Parameters:**

- `condition` (func(\*http.Request) bool)
- `authMiddleware` (RequestMiddleware)

**Returns:**

- RequestMiddleware

### AddConditionalRequest

AddConditionalRequest creates a conditional request middleware

```go
func AddConditionalRequest(ifModifiedSince, ifNoneMatch string) RequestMiddleware
```

**Parameters:**

- `ifModifiedSince` (string)
- `ifNoneMatch` (string)

**Returns:**

- RequestMiddleware

### AddConditionalTimeout

AddConditionalTimeout creates a conditional timeout middleware

```go
func AddConditionalTimeout(condition func(*http.Request) bool, timeout time.Duration) RequestMiddleware
```

**Parameters:**

- `condition` (func(\*http.Request) bool)
- `timeout` (time.Duration)

**Returns:**

- RequestMiddleware

### AddContentType

AddContentType creates a content type middleware

```go
func AddContentType(contentType string) RequestMiddleware
```

**Parameters:**

- `contentType` (string)

**Returns:**

- RequestMiddleware

### AddCorrelationID

AddCorrelationID creates a correlation ID middleware for distributed tracing

```go
func AddCorrelationID(config RequestIDConfig) RequestMiddleware
```

**Parameters:**

- `config` (RequestIDConfig)

**Returns:**

- RequestMiddleware

### AddCustomAuth

AddCustomAuth creates a custom authentication middleware

```go
func AddCustomAuth(headerName, headerValue string) RequestMiddleware
```

**Parameters:**

- `headerName` (string)
- `headerValue` (string)

**Returns:**

- RequestMiddleware

### AddCustomHeader

AddCustomHeader creates a custom header middleware

```go
func AddCustomHeader(headers map[string]string) RequestMiddleware
```

**Parameters:**

- `headers` (map[string]string)

**Returns:**

- RequestMiddleware

### AddDeadline

AddDeadline creates a deadline middleware

```go
func AddDeadline(deadline time.Time) RequestMiddleware
```

**Parameters:**

- `deadline` (time.Time)

**Returns:**

- RequestMiddleware

### AddDebugLogging

AddDebugLogging creates a debug logging middleware with detailed information

```go
func AddDebugLogging(config LoggingConfig) RequestMiddleware
```

**Parameters:**

- `config` (LoggingConfig)

**Returns:**

- RequestMiddleware

### AddDigestAuth

AddDigestAuth creates a Digest authentication middleware (simplified)

```go
func AddDigestAuth(username, password, realm string) RequestMiddleware
```

**Parameters:**

- `username` (string)
- `password` (string)
- `realm` (string)

**Returns:**

- RequestMiddleware

### AddDynamicAuth

AddDynamicAuth creates a dynamic authentication middleware

```go
func AddDynamicAuth(authFunc func(*http.Request) error) RequestMiddleware
```

**Parameters:**

- `authFunc` (func(\*http.Request) error)

**Returns:**

- RequestMiddleware

### AddETag

AddETag creates an ETag middleware

```go
func AddETag(etag string) RequestMiddleware
```

**Parameters:**

- `etag` (string)

**Returns:**

- RequestMiddleware

### AddFormContentType

AddFormContentType creates a form content type middleware

```go
func AddFormContentType() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddGlobalTimeout

AddGlobalTimeout creates a global timeout middleware

```go
func AddGlobalTimeout(timeout time.Duration) RequestMiddleware
```

**Parameters:**

- `timeout` (time.Duration)

**Returns:**

- RequestMiddleware

### AddHeader

AddHeader creates a header middleware based on configuration

```go
func AddHeader(config HeaderConfig) RequestMiddleware
```

**Parameters:**

- `config` (HeaderConfig)

**Returns:**

- RequestMiddleware

### AddHeaderFilter

AddHeaderFilter creates a header filter middleware

```go
func AddHeaderFilter(allowedHeaders []string) RequestMiddleware
```

**Parameters:**

- `allowedHeaders` ([]string)

**Returns:**

- RequestMiddleware

### AddHeaderTransformation

AddHeaderTransformation creates a header transformation middleware

```go
func AddHeaderTransformation(transformFunc func(string, string) (string, string)) RequestMiddleware
```

**Parameters:**

- `transformFunc` (func(string, string) (string, string))

**Returns:**

- RequestMiddleware

### AddJSONContentType

AddJSONContentType creates a JSON content type middleware

```go
func AddJSONContentType() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddJWTTokenAuth

AddJWTTokenAuth creates a JWT token authentication middleware

```go
func AddJWTTokenAuth(token string) RequestMiddleware
```

**Parameters:**

- `token` (string)

**Returns:**

- RequestMiddleware

### AddLogging

AddLogging creates a logging middleware

```go
func AddLogging(config LoggingConfig) RequestMiddleware
```

**Parameters:**

- `config` (LoggingConfig)

**Returns:**

- RequestMiddleware

### AddMetrics

AddMetrics creates a metrics collection middleware

```go
func AddMetrics(collector *MetricsCollector) RequestMiddleware
```

**Parameters:**

- `collector` (\*MetricsCollector)

**Returns:**

- RequestMiddleware

### AddMultiAuth

AddMultiAuth creates a middleware that tries multiple authentication methods

```go
func AddMultiAuth(authMethods ...RequestMiddleware) RequestMiddleware
```

**Parameters:**

- `authMethods` (...RequestMiddleware)

**Returns:**

- RequestMiddleware

### AddMultipartContentType

AddMultipartContentType creates a multipart content type middleware

```go
func AddMultipartContentType() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddNoCache

AddNoCache creates a no-cache middleware

```go
func AddNoCache() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddOAuth2Auth

AddOAuth2Auth creates an OAuth2 authentication middleware

```go
func AddOAuth2Auth(accessToken string) RequestMiddleware
```

**Parameters:**

- `accessToken` (string)

**Returns:**

- RequestMiddleware

### AddPerRequestTimeout

AddPerRequestTimeout creates a per-request timeout middleware

```go
func AddPerRequestTimeout(timeout time.Duration) RequestMiddleware
```

**Parameters:**

- `timeout` (time.Duration)

**Returns:**

- RequestMiddleware

### AddRateLimit

AddRateLimit adds custom rate limiting headers

```go
func AddRateLimit(rateLimitInfo RateLimitInfoProvider) RequestMiddleware
```

**Parameters:**

- `rateLimitInfo` (RateLimitInfoProvider)

**Returns:**

- RequestMiddleware

### AddRequestID

AddRequestID creates a request ID middleware

```go
func AddRequestID(config RequestIDConfig) RequestMiddleware
```

**Parameters:**

- `config` (RequestIDConfig)

**Returns:**

- RequestMiddleware

### AddRetry

AddRetry implements retry logic at the middleware level

```go
func AddRetry(maxRetries int, retryCondition RetryCondition) RequestMiddleware
```

**Parameters:**

- `maxRetries` (int)
- `retryCondition` (RetryCondition)

**Returns:**

- RequestMiddleware

### AddRotatingAuth

AddRotatingAuth creates a rotating authentication middleware

```go
func AddRotatingAuth(tokens []string, headerName string) RequestMiddleware
```

**Parameters:**

- `tokens` ([]string)
- `headerName` (string)

**Returns:**

- RequestMiddleware

### AddSecurityHeaders

AddSecurityHeaders creates a security headers middleware

```go
func AddSecurityHeaders() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### AddStructuredLogging

AddStructuredLogging creates a structured logging middleware

```go
func AddStructuredLogging(config LoggingConfig) RequestMiddleware
```

**Parameters:**

- `config` (LoggingConfig)

**Returns:**

- RequestMiddleware

### AddTimeout

AddTimeout creates a timeout middleware

```go
func AddTimeout(config TimeoutConfig) RequestMiddleware
```

**Parameters:**

- `config` (TimeoutConfig)

**Returns:**

- RequestMiddleware

### AddTimeoutChain

AddTimeoutChain creates a chain of timeout middlewares

```go
func AddTimeoutChain(timeouts ...time.Duration) RequestMiddleware
```

**Parameters:**

- `timeouts` (...time.Duration)

**Returns:**

- RequestMiddleware

### AddTimeoutFromContext

AddTimeoutFromContext creates a timeout middleware that gets timeout from context

```go
func AddTimeoutFromContext(contextKey string, defaultTimeout time.Duration) RequestMiddleware
```

**Parameters:**

- `contextKey` (string)
- `defaultTimeout` (time.Duration)

**Returns:**

- RequestMiddleware

### AddTracing

AddTracing creates a tracing middleware that propagates request ID

```go
func AddTracing(config RequestIDConfig) RequestMiddleware
```

**Parameters:**

- `config` (RequestIDConfig)

**Returns:**

- RequestMiddleware

### AddUserAgent

AddUserAgent creates a user agent middleware

```go
func AddUserAgent(userAgent string) RequestMiddleware
```

**Parameters:**

- `userAgent` (string)

**Returns:**

- RequestMiddleware

### AddValidation

AddValidation validates request payloads

```go
func AddValidation(validator Validator) RequestMiddleware
```

**Parameters:**

- `validator` (Validator)

**Returns:**

- RequestMiddleware

### AddXMLContentType

AddXMLContentType creates an XML content type middleware

```go
func AddXMLContentType() RequestMiddleware
```

**Parameters:**
None

**Returns:**

- RequestMiddleware

### RequestOptions

RequestOptions contains configuration for individual requests

#### Example Usage

```go
// Create a new RequestOptions
requestoptions := RequestOptions{
    Headers: /* value */,
    Query: map[],
    Timeout: &/* value */{},
    Context: /* value */,
    Retries: &42{},
    RateLimitID: "example",
}
```

#### Type Definition

```go
type RequestOptions struct {
    Headers http.Header
    Query map[string]any
    Timeout *time.Duration
    Context context.Context
    Retries *int
    RateLimitID string
}
```

### Fields

| Field       | Type              | Description                 |
| ----------- | ----------------- | --------------------------- |
| Headers     | `http.Header`     |                             |
| Query       | `map[string]any`  |                             |
| Timeout     | `*time.Duration`  |                             |
| Context     | `context.Context` |                             |
| Retries     | `*int`            |                             |
| RateLimitID | `string`          | Custom rate limit bucket ID |

### RequestQueue

RequestQueue manages queued requests for a specific route/bucket

#### Example Usage

```go
// Create a new RequestQueue
requestqueue := RequestQueue{
    Queue: [],
    RateLimiter: RateLimiter{},
    Processing: true,
    LastUsed: /* value */,
}
```

#### Type Definition

```go
type RequestQueue struct {
    Queue []QueuedRequest
    RateLimiter RateLimiter
    Processing bool
    LastUsed time.Time
}
```

### Fields

| Field       | Type              | Description |
| ----------- | ----------------- | ----------- |
| Queue       | `[]QueuedRequest` |             |
| RateLimiter | `RateLimiter`     |             |
| Processing  | `bool`            |             |
| LastUsed    | `time.Time`       |             |

### Response

Response wraps HTTP response data with type safety

#### Example Usage

```go
// Create a new Response
response := Response{
    Data: T{},
    StatusCode: 42,
    Headers: /* value */,
    Raw: &/* value */{},
}
```

#### Type Definition

```go
type Response struct {
    Data T
    StatusCode int
    Headers http.Header
    Raw *http.Response
}
```

### Fields

| Field      | Type             | Description |
| ---------- | ---------------- | ----------- |
| Data       | `T`              |             |
| StatusCode | `int`            |             |
| Headers    | `http.Header`    |             |
| Raw        | `*http.Response` |             |

### Constructor Functions

### Execute

Execute performs a type-safe HTTP request

```go
func Execute(client *Client, route *ast.IndexListExpr, request TRequest, options *RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**

- `client` (\*Client)
- `route` (\*ast.IndexListExpr)
- `request` (TRequest)
- `options` (\*RequestOptions)

**Returns:**

- \*\*ast.IndexExpr
- error

### ResponseMiddleware

ResponseMiddleware processes responses after they are received

#### Example Usage

```go
// Example usage of ResponseMiddleware
var value ResponseMiddleware
// Initialize with appropriate value
```

#### Type Definition

```go
type ResponseMiddleware func(resp *http.Response) error
```

### Constructor Functions

### AddResponseCache

AddResponseCache implements response caching

```go
func AddResponseCache(cache Cache) ResponseMiddleware
```

**Parameters:**

- `cache` (Cache)

**Returns:**

- ResponseMiddleware

### AddResponseCircuitBreaker

AddResponseCircuitBreaker handles circuit breaker state updates

```go
func AddResponseCircuitBreaker() ResponseMiddleware
```

**Parameters:**
None

**Returns:**

- ResponseMiddleware

### AddResponseCompression

AddResponseCompression creates a response compression middleware

```go
func AddResponseCompression(config CompressionConfig) ResponseMiddleware
```

**Parameters:**

- `config` (CompressionConfig)

**Returns:**

- ResponseMiddleware

### AddResponseDebug

AddResponseDebug creates a debug response logging middleware

```go
func AddResponseDebug(config LoggingConfig) ResponseMiddleware
```

**Parameters:**

- `config` (LoggingConfig)

**Returns:**

- ResponseMiddleware

### AddResponseDecompression

AddResponseDecompression creates a response decompression middleware

```go
func AddResponseDecompression() ResponseMiddleware
```

**Parameters:**
None

**Returns:**

- ResponseMiddleware

### AddResponseErrorLogging

AddResponseErrorLogging creates an error logging middleware

```go
func AddResponseErrorLogging(config LoggingConfig) ResponseMiddleware
```

**Parameters:**

- `config` (LoggingConfig)

**Returns:**

- ResponseMiddleware

### AddResponseLogging

AddResponseLogging creates a response logging middleware

```go
func AddResponseLogging(config LoggingConfig) ResponseMiddleware
```

**Parameters:**

- `config` (LoggingConfig)

**Returns:**

- ResponseMiddleware

### AddResponseMetrics

AddResponseMetrics creates a response metrics collection middleware

```go
func AddResponseMetrics(collector *MetricsCollector) ResponseMiddleware
```

**Parameters:**

- `collector` (\*MetricsCollector)

**Returns:**

- ResponseMiddleware

### AddResponseRequestID

AddResponseRequestID creates a response middleware that logs request ID

```go
func AddResponseRequestID(config RequestIDConfig) ResponseMiddleware
```

**Parameters:**

- `config` (RequestIDConfig)

**Returns:**

- ResponseMiddleware

### AddResponseStructured

AddResponseStructured creates a structured response logging middleware

```go
func AddResponseStructured(config LoggingConfig) ResponseMiddleware
```

**Parameters:**

- `config` (LoggingConfig)

**Returns:**

- ResponseMiddleware

### AddResponseTimeout

AddResponseTimeout creates a response timeout middleware

```go
func AddResponseTimeout(timeout time.Duration) ResponseMiddleware
```

**Parameters:**

- `timeout` (time.Duration)

**Returns:**

- ResponseMiddleware

### RetryCondition

RetryCondition determines if a request should be retried

#### Example Usage

```go
// Example usage of RetryCondition
var value RetryCondition
// Initialize with appropriate value
```

#### Type Definition

```go
type RetryCondition func(resp *http.Response, err error) bool
```

### Route

Route represents a type-safe route definition

#### Example Usage

```go
// Create a new Route
route := Route{
    Method: HTTPMethod{},
    Path: "example",
}
```

#### Type Definition

```go
type Route struct {
    Method HTTPMethod
    Path string
}
```

### Fields

| Field  | Type         | Description |
| ------ | ------------ | ----------- |
| Method | `HTTPMethod` |             |
| Path   | `string`     |             |

### Constructor Functions

### NewRoute

NewRoute creates a new type-safe route

```go
func NewRoute(method HTTPMethod, path string) Route
```

**Parameters:**

- `method` (HTTPMethod)
- `path` (string)

**Returns:**

- \*ast.IndexListExpr

### SequentialGenerator

SequentialGenerator generates sequential request IDs

#### Example Usage

```go
// Create a new SequentialGenerator
sequentialgenerator := SequentialGenerator{

}
```

#### Type Definition

```go
type SequentialGenerator struct {
}
```

## Methods

### Generate

Generate creates a sequential request ID

```go
func (*SequentialGenerator) Generate() string
```

**Parameters:**
None

**Returns:**

- string

### Serializable

Serializable represents types that can be serialized for requests

#### Example Usage

```go
// Example implementation of Serializable
type MySerializable struct {
    // Add your fields here
}

func (m MySerializable) MarshalJSON() []byte {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Serializable interface {
    MarshalJSON() ([]byte, error)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### SimpleLogger

Default implementations SimpleLogger provides a basic logger implementation

#### Example Usage

```go
// Create a new SimpleLogger
simplelogger := SimpleLogger{

}
```

#### Type Definition

```go
type SimpleLogger struct {
}
```

## Methods

### LogRequest

```go
func (*SimpleLogger) LogRequest(entry LogEntry)
```

**Parameters:**

- `entry` (LogEntry)

**Returns:**
None

### LogResponse

```go
func (*SimpleLogger) LogResponse(entry LogEntry)
```

**Parameters:**

- `entry` (LogEntry)

**Returns:**
None

### StaticAuthProvider

StaticAuthProvider provides static token authentication

#### Example Usage

```go
// Create a new StaticAuthProvider
staticauthprovider := StaticAuthProvider{
    Token: "example",
    Prefix: "example",
}
```

#### Type Definition

```go
type StaticAuthProvider struct {
    Token string
    Prefix string
}
```

### Fields

| Field  | Type     | Description |
| ------ | -------- | ----------- |
| Token  | `string` |             |
| Prefix | `string` |             |

## Methods

### GetAuthHeader

```go
func (*StaticAuthProvider) GetAuthHeader(token string) string
```

**Parameters:**

- `token` (string)

**Returns:**

- string

### GetToken

```go
func (*StaticAuthProvider) GetToken(ctx context.Context) (string, error)
```

**Parameters:**

- `ctx` (context.Context)

**Returns:**

- string
- error

### TimeoutConfig

TimeoutConfig configures timeout middleware

#### Example Usage

```go
// Create a new TimeoutConfig
timeoutconfig := TimeoutConfig{
    Timeout: /* value */,
    PerRequest: true,
    GlobalTimeout: /* value */,
    RequestTimeout: /* value */,
}
```

#### Type Definition

```go
type TimeoutConfig struct {
    Timeout time.Duration
    PerRequest bool
    GlobalTimeout time.Duration
    RequestTimeout time.Duration
}
```

### Fields

| Field          | Type            | Description |
| -------------- | --------------- | ----------- |
| Timeout        | `time.Duration` |             |
| PerRequest     | `bool`          |             |
| GlobalTimeout  | `time.Duration` |             |
| RequestTimeout | `time.Duration` |             |

### Constructor Functions

### DefaultTimeoutConfig

DefaultTimeoutConfig returns a default timeout configuration

```go
func DefaultTimeoutConfig() TimeoutConfig
```

**Parameters:**
None

**Returns:**

- TimeoutConfig

### TimestampGenerator

TimestampGenerator generates timestamp-based request IDs

#### Example Usage

```go
// Create a new TimestampGenerator
timestampgenerator := TimestampGenerator{

}
```

#### Type Definition

```go
type TimestampGenerator struct {
}
```

## Methods

### Generate

Generate creates a timestamp-based request ID

```go
func (*SequentialGenerator) Generate() string
```

**Parameters:**
None

**Returns:**

- string

### UUIDGenerator

UUIDGenerator generates UUID-style request IDs

#### Example Usage

```go
// Create a new UUIDGenerator
uuidgenerator := UUIDGenerator{

}
```

#### Type Definition

```go
type UUIDGenerator struct {
}
```

## Methods

### Generate

Generate creates a new UUID-style request ID

```go
func (*SequentialGenerator) Generate() string
```

**Parameters:**
None

**Returns:**

- string

### Validator

Validator interface for request validation

#### Example Usage

```go
// Example implementation of Validator
type MyValidator struct {
    // Add your fields here
}

func (m MyValidator) Validate(param1 []byte, param2 string) error {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Validator interface {
    Validate(data []byte, contentType string) error
}
```

## Methods

| Method | Description |
| ------ | ----------- |

## Functions

### AddAutoMetrics

AddAutoMetrics creates a simple metrics middleware that doesn't require a collector

```go
func AddAutoMetrics() (RequestMiddleware, ResponseMiddleware, func() MetricsSnapshot)
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `RequestMiddleware` | |
| `ResponseMiddleware` | |
| `func() MetricsSnapshot` | |

**Example:**

```go
// Example usage of AddAutoMetrics
result := AddAutoMetrics(/* parameters */)
```

### GetRequestID

GetRequestID extracts the request ID from context

```go
func GetRequestID(ctx context.Context) (string, bool)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |
| `bool` | |

**Example:**

```go
// Example usage of GetRequestID
result := GetRequestID(/* parameters */)
```

### RequestIDFromRequest

RequestIDFromRequest extracts the request ID from an HTTP request

```go
func RequestIDFromRequest(req *http.Request) (string, bool)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `req` | `*http.Request` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |
| `bool` | |

**Example:**

```go
// Example usage of RequestIDFromRequest
result := RequestIDFromRequest(/* parameters */)
```

### RequestIDFromResponse

RequestIDFromResponse extracts the request ID from an HTTP response

```go
func RequestIDFromResponse(resp *http.Response) (string, bool)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `resp` | `*http.Response` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |
| `bool` | |

**Example:**

```go
// Example usage of RequestIDFromResponse
result := RequestIDFromResponse(/* parameters */)
```

### WithRequestID

WithRequestID adds a request ID to the context

```go
func WithRequestID(ctx context.Context, requestID string) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `requestID` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithRequestID
result := WithRequestID(/* parameters */)
```

## External Links

- [Package Overview](../packages/neuron.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/neuron)
- [Source Code](https://github.com/kolosys/neuron/tree/dev/github.com/kolosys/neuron)
