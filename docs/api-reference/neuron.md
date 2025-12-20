# neuron API

Complete API documentation for the neuron package.

**Import Path:** `github.com/kolosys/neuron`

## Package Documentation



## Types

### AuthConfig
AuthConfig configures authentication hooks

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| Type | `AuthType` |  |
| Token | `string` |  |
| Username | `string` |  |
| Password | `string` |  |
| HeaderName | `string` |  |
| HeaderValue | `string` |  |

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
Cache interface for caching hooks

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| Data | `[]byte` |  |
| Headers | `http.Header` |  |
| StatusCode | `int` |  |
| Timestamp | `time.Time` |  |
| TTL | `time.Duration` |  |

### Client
Client provides a type-safe HTTP client with rate limiting and circuit breaking

#### Example Usage

```go
// Create a new Client
client := Client{
    Config: ClientOptions{},
    Metrics: &MetricsCollector{}{},
}
```

#### Type Definition

```go
type Client struct {
    Config ClientOptions
    Metrics *MetricsCollector
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Config | `ClientOptions` |  |
| Metrics | `*MetricsCollector` | Metrics tracking |

### Constructor Functions

### NewClient

NewClient creates a new type-safe HTTP client

```go
func NewClient(options ClientOptions) *Client
```

**Parameters:**
- `options` (ClientOptions)

**Returns:**
- *Client

## Methods

### Delete

Delete executes a DELETE request and returns the response

```go
func (*InMemoryCache) Delete(key string)
```

**Parameters:**
- `key` (string)

**Returns:**
  None

### Do

Do executes an HTTP request with the specified method and returns the response

```go
func (*Client) Do(method HTTPMethod, path string, opts ...*RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `method` (HTTPMethod)
- `path` (string)
- `opts` (...*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

### Get

Get executes a GET request and returns the response

```go
func (*InMemoryCache) Get(key string) (*CacheEntry, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- *CacheEntry
- bool

### GetMetrics

GetMetrics returns current client metrics

```go
func (*Client) GetMetrics() MetricsSnapshot
```

**Parameters:**
  None

**Returns:**
- MetricsSnapshot

### Head

Head executes a HEAD request and returns the response

```go
func (*Client) Head(path string, opts ...*RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `path` (string)
- `opts` (...*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

### Options

Options executes an OPTIONS request and returns the response

```go
func (*Client) Options(path string, opts ...*RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `path` (string)
- `opts` (...*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

### Patch

Patch executes a PATCH request and returns the response

```go
func (*Client) Patch(path string, opts ...*RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `path` (string)
- `opts` (...*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

### Post

Post executes a POST request and returns the response

```go
func (*Client) Post(path string, opts ...*RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `path` (string)
- `opts` (...*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

### Put

Put executes a PUT request and returns the response

```go
func (*Client) Put(path string, opts ...*RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `path` (string)
- `opts` (...*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

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
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Type | `ErrorType` |  |
| Message | `string` |  |
| StatusCode | `int` |  |
| Route | `string` |  |
| Method | `string` | HTTP method |
| URL | `string` | Full URL |
| Attempt | `int` | Retry attempt number |
| Timestamp | `time.Time` | When error occurred |
| Cause | `error` |  |

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
- `req` (*http.Request)
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
    MaxRetries: 42,
    RetryDelay: /* value */,
    RetryMultiplier: 3.14,
    RequestHooks: [],
    ResponseHooks: [],
    HTTPClient: &/* value */{},
    AdaptiveTimeout: true,
}
```

#### Type Definition

```go
type ClientOptions struct {
    BaseURL string
    UserAgent string
    Headers http.Header
    Timeout time.Duration
    MaxRetries int
    RetryDelay time.Duration
    RetryMultiplier float64
    RequestHooks []RequestHook
    ResponseHooks []ResponseHook
    HTTPClient *http.Client
    AdaptiveTimeout bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| BaseURL | `string` | Base configuration |
| UserAgent | `string` |  |
| Headers | `http.Header` |  |
| Timeout | `time.Duration` |  |
| MaxRetries | `int` | Request handling |
| RetryDelay | `time.Duration` |  |
| RetryMultiplier | `float64` |  |
| RequestHooks | `[]RequestHook` | Hooks |
| ResponseHooks | `[]ResponseHook` |  |
| HTTPClient | `*http.Client` | HTTP client |
| AdaptiveTimeout | `bool` | Resilience configuration |

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| Type | `CompressionType` |  |
| Level | `int` |  |
| MinSize | `int` |  |
| ContentTypes | `[]string` |  |

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
HeaderConfig configures header hooks

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| UserAgent | `string` |  |
| CustomHeaders | `map[string]string` |  |
| RemoveHeaders | `[]string` |  |
| AddHeaders | `map[string]string` |  |

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

### HookChain
HookChain manages a chain of hook functions This is an optional utility for advanced use cases; most users will use slices directly

#### Example Usage

```go
// Create a new HookChain
hookchain := HookChain{

}
```

#### Type Definition

```go
type HookChain struct {
}
```

### Constructor Functions

### NewHookChain

NewHookChain creates a new hook chain

```go
func NewHookChain() *HookChain
```

**Parameters:**
  None

**Returns:**
- *HookChain

## Methods

### AddRequestHook

AddRequestHook adds a request hook to the chain

```go
func (*HookChain) AddRequestHook(hook RequestHook) *HookChain
```

**Parameters:**
- `hook` (RequestHook)

**Returns:**
- *HookChain

### AddResponseHook

AddResponseHook adds a response hook to the chain

```go
func (*HookChain) AddResponseHook(hook ResponseHook) *HookChain
```

**Parameters:**
- `hook` (ResponseHook)

**Returns:**
- *HookChain

### ApplyRequestHooks

ApplyRequestHooks applies all request hooks in order

```go
func (*HookChain) ApplyRequestHooks(req *http.Request) error
```

**Parameters:**
- `req` (*http.Request)

**Returns:**
- error

### ApplyResponseHooks

ApplyResponseHooks applies all response hooks in order

```go
func (*HookChain) ApplyResponseHooks(resp *http.Response) error
```

**Parameters:**
- `resp` (*http.Response)

**Returns:**
- error

### InMemoryCache
InMemoryCache provides a simple in-memory cache implementation

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

NewInMemoryCache creates a new in-memory cache

```go
func NewInMemoryCache() *InMemoryCache
```

**Parameters:**
  None

**Returns:**
- *InMemoryCache

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
- *CacheEntry
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

### LoggingConfig
LoggingConfig configures the logging hooks

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| Level | `LogLevel` |  |
| IncludeBody | `bool` |  |
| MaxBodySize | `int` |  |
| Logger | `*log.Logger` |  |

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

### MetricsCollector
MetricsCollector collects and stores metrics

#### Example Usage

```go
// Create a new MetricsCollector
metricscollector := MetricsCollector{
    RequestCount: /* value */,
    ResponseCount: /* value */,
    ErrorCount: /* value */,
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
    RequestCount atomic.Int64
    ResponseCount atomic.Int64
    ErrorCount atomic.Int64
    TotalDuration atomic.Int64
    MinDuration atomic.Int64
    MaxDuration atomic.Int64
    StatusCodeCounts map[int]*atomic.Int64
    StartTime time.Time
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| RequestCount | `atomic.Int64` | Request metrics |
| ResponseCount | `atomic.Int64` |  |
| ErrorCount | `atomic.Int64` |  |
| TotalDuration | `atomic.Int64` | Duration metrics |
| MinDuration | `atomic.Int64` | nanoseconds |
| MaxDuration | `atomic.Int64` | nanoseconds |
| StatusCodeCounts | `map[int]*atomic.Int64` | Status code metrics |
| StartTime | `time.Time` | Start time for uptime calculation |

### Constructor Functions

### NewMetricsCollector

NewMetricsCollector creates a new metrics collector

```go
func NewMetricsCollector() *MetricsCollector
```

**Parameters:**
  None

**Returns:**
- *MetricsCollector

## Methods

### GetMetrics

GetMetrics returns current metrics

```go
func (*Client) GetMetrics() MetricsSnapshot
```

**Parameters:**
  None

**Returns:**
- MetricsSnapshot

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| RequestCount | `int64` |  |
| ResponseCount | `int64` |  |
| ErrorCount | `int64` |  |
| AverageDuration | `time.Duration` |  |
| MinDuration | `time.Duration` |  |
| MaxDuration | `time.Duration` |  |
| StatusCodeCounts | `map[int]int64` |  |
| Uptime | `time.Duration` |  |

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

### RequestHook
RequestHook processes requests before they are sent

#### Example Usage

```go
// Example usage of RequestHook
var value RequestHook
// Initialize with appropriate value
```

#### Type Definition

```go
type RequestHook func(req *http.Request) error
```

### Constructor Functions

### AddAPIKeyAuth

AddAPIKeyAuth creates an API key authentication hook

```go
func AddAPIKeyAuth(apiKey, headerName string) RequestHook
```

**Parameters:**
- `apiKey` (string)
- `headerName` (string)

**Returns:**
- RequestHook

### AddAPIKeyHeaderAuth

AddAPIKeyHeaderAuth creates an API key authentication middleware with common header names

```go
func AddAPIKeyHeaderAuth(apiKey string) RequestHook
```

**Parameters:**
- `apiKey` (string)

**Returns:**
- RequestHook

### AddAccept

AddAccept creates an accept header middleware

```go
func AddAccept(accept string) RequestHook
```

**Parameters:**
- `accept` (string)

**Returns:**
- RequestHook

### AddAuth

AddAuth creates an authentication hook based on configuration

```go
func AddAuth(config AuthConfig) RequestHook
```

**Parameters:**
- `config` (AuthConfig)

**Returns:**
- RequestHook

### AddAuthHeader

AddAuthHeader creates a generic header-based authentication hook

```go
func AddAuthHeader(header, value string) RequestHook
```

**Parameters:**
- `header` (string)
- `value` (string)

**Returns:**
- RequestHook

### AddAuthentication

AddAuthentication adds authentication headers using an AuthProvider

```go
func AddAuthentication(authProvider AuthProvider) RequestHook
```

**Parameters:**
- `authProvider` (AuthProvider)

**Returns:**
- RequestHook

### AddAutoCompression

AddAutoCompression creates an automatic compression middleware

```go
func AddAutoCompression() RequestHook
```

**Parameters:**
  None

**Returns:**
- RequestHook

### AddBasicAuth

AddBasicAuth creates a Basic authentication hook

```go
func AddBasicAuth(username, password string) RequestHook
```

**Parameters:**
- `username` (string)
- `password` (string)

**Returns:**
- RequestHook

### AddBearerAuth

AddBearerAuth creates a Bearer token authentication hook

```go
func AddBearerAuth(token string) RequestHook
```

**Parameters:**
- `token` (string)

**Returns:**
- RequestHook

### AddCompression

AddCompression creates a request compression middleware

```go
func AddCompression(config CompressionConfig) RequestHook
```

**Parameters:**
- `config` (CompressionConfig)

**Returns:**
- RequestHook

### AddCompressionLevel

AddCompressionLevel creates a compression level middleware

```go
func AddCompressionLevel(level int) RequestHook
```

**Parameters:**
- `level` (int)

**Returns:**
- RequestHook

### AddContentType

AddContentType creates a content type middleware

```go
func AddContentType(contentType string) RequestHook
```

**Parameters:**
- `contentType` (string)

**Returns:**
- RequestHook

### AddCorrelationID

AddCorrelationID creates a correlation ID middleware for distributed tracing

```go
func AddCorrelationID(config RequestIDConfig) RequestHook
```

**Parameters:**
- `config` (RequestIDConfig)

**Returns:**
- RequestHook

### AddDebugLogging

AddDebugLogging creates a debug logging middleware with detailed information

```go
func AddDebugLogging(config LoggingConfig) RequestHook
```

**Parameters:**
- `config` (LoggingConfig)

**Returns:**
- RequestHook

### AddHeader

AddHeader creates a header hook based on configuration

```go
func AddHeader(config HeaderConfig) RequestHook
```

**Parameters:**
- `config` (HeaderConfig)

**Returns:**
- RequestHook

### AddHeaderSet

AddHeaderSet creates a middleware that sets a single header

```go
func AddHeaderSet(key, value string) RequestHook
```

**Parameters:**
- `key` (string)
- `value` (string)

**Returns:**
- RequestHook

### AddLogging

AddLogging creates a logging hook

```go
func AddLogging(config LoggingConfig) RequestHook
```

**Parameters:**
- `config` (LoggingConfig)

**Returns:**
- RequestHook

### AddMetrics

AddMetrics creates a metrics collection middleware

```go
func AddMetrics(collector *MetricsCollector) RequestHook
```

**Parameters:**
- `collector` (*MetricsCollector)

**Returns:**
- RequestHook

### AddNoCache

AddNoCache creates a no-cache middleware

```go
func AddNoCache() RequestHook
```

**Parameters:**
  None

**Returns:**
- RequestHook

### AddRequestID

AddRequestID creates a request ID middleware

```go
func AddRequestID(config RequestIDConfig) RequestHook
```

**Parameters:**
- `config` (RequestIDConfig)

**Returns:**
- RequestHook

### AddSecurityHeaders

AddSecurityHeaders creates a security headers middleware

```go
func AddSecurityHeaders() RequestHook
```

**Parameters:**
  None

**Returns:**
- RequestHook

### AddStructuredLogging

AddStructuredLogging creates a structured logging middleware

```go
func AddStructuredLogging(config LoggingConfig) RequestHook
```

**Parameters:**
- `config` (LoggingConfig)

**Returns:**
- RequestHook

### AddTracing

AddTracing creates a tracing middleware that propagates request ID

```go
func AddTracing(config RequestIDConfig) RequestHook
```

**Parameters:**
- `config` (RequestIDConfig)

**Returns:**
- RequestHook

### AddUserAgent

AddUserAgent creates a user agent middleware

```go
func AddUserAgent(userAgent string) RequestHook
```

**Parameters:**
- `userAgent` (string)

**Returns:**
- RequestHook

### AddValidation

AddValidation validates request payloads

```go
func AddValidation(validator Validator) RequestHook
```

**Parameters:**
- `validator` (Validator)

**Returns:**
- RequestHook

### RequestIDConfig
RequestIDConfig configures request ID generation

#### Example Usage

```go
// Create a new RequestIDConfig
requestidconfig := RequestIDConfig{
    Generator: RequestIDGenerator{},
    HeaderName: "example",
    ContextKey: any{},
    Propagate: true,
}
```

#### Type Definition

```go
type RequestIDConfig struct {
    Generator RequestIDGenerator
    HeaderName string
    ContextKey any
    Propagate bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Generator | `RequestIDGenerator` |  |
| HeaderName | `string` |  |
| ContextKey | `any` |  |
| Propagate | `bool` |  |

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
RequestMetrics provides a snapshot of insights into request performance

#### Example Usage

```go
// Create a new RequestMetrics
requestmetrics := RequestMetrics{
    TotalRequests: 42,
    SuccessfulRequests: 42,
    FailedRequests: 42,
    AverageResponseTime: /* value */,
}
```

#### Type Definition

```go
type RequestMetrics struct {
    TotalRequests int64
    SuccessfulRequests int64
    FailedRequests int64
    AverageResponseTime time.Duration
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| TotalRequests | `int64` |  |
| SuccessfulRequests | `int64` |  |
| FailedRequests | `int64` |  |
| AverageResponseTime | `time.Duration` |  |

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
    Body: any{},
    RequestHooks: [],
    ResponseHooks: [],
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
    Body any
    RequestHooks []RequestHook
    ResponseHooks []ResponseHook
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Headers | `http.Header` |  |
| Query | `map[string]any` |  |
| Timeout | `*time.Duration` |  |
| Context | `context.Context` |  |
| Retries | `*int` |  |
| Body | `any` | Request body (JSON, form data, io.Reader, BodyProvider) |
| RequestHooks | `[]RequestHook` | Per-request hooks (runs after client-level hooks) |
| ResponseHooks | `[]ResponseHook` |  |

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| Data | `T` |  |
| StatusCode | `int` |  |
| Headers | `http.Header` |  |
| Raw | `*http.Response` |  |

### Constructor Functions

### DoWithType

DoWithType executes a request and unmarshals the response to the specified type

```go
func DoWithType(client *Client, method HTTPMethod, path string, opts ...*RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `client` (*Client)
- `method` (HTTPMethod)
- `path` (string)
- `opts` (...*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

### Execute

Execute performs a type-safe HTTP request

```go
func Execute(client *Client, route *ast.IndexListExpr, request TRequest, options *RequestOptions) (**ast.IndexExpr, error)
```

**Parameters:**
- `client` (*Client)
- `route` (*ast.IndexListExpr)
- `request` (TRequest)
- `options` (*RequestOptions)

**Returns:**
- **ast.IndexExpr
- error

## Methods

### Bytes

Bytes returns the response body as bytes

```go
func (**ast.IndexExpr) Bytes() ([]byte, error)
```

**Parameters:**
  None

**Returns:**
- []byte
- error

### IsClientError

IsClientError returns true if the status code is in the 4xx range

```go
func (**ast.IndexExpr) IsClientError() bool
```

**Parameters:**
  None

**Returns:**
- bool

### IsError

IsError returns true if the status code is 4xx or 5xx

```go
func (**ast.IndexExpr) IsError() bool
```

**Parameters:**
  None

**Returns:**
- bool

### IsServerError

IsServerError returns true if the status code is in the 5xx range

```go
func (**ast.IndexExpr) IsServerError() bool
```

**Parameters:**
  None

**Returns:**
- bool

### IsSuccess

IsSuccess returns true if the status code is in the 2xx range

```go
func (**ast.IndexExpr) IsSuccess() bool
```

**Parameters:**
  None

**Returns:**
- bool

### JSON

JSON unmarshals the response body as JSON

```go
func (**ast.IndexExpr) JSON(target any) error
```

**Parameters:**
- `target` (any)

**Returns:**
- error

### String

String returns the response body as a string

```go
func (**ast.IndexExpr) String() (string, error)
```

**Parameters:**
  None

**Returns:**
- string
- error

### ResponseHook
ResponseHook processes responses after they are received

#### Example Usage

```go
// Example usage of ResponseHook
var value ResponseHook
// Initialize with appropriate value
```

#### Type Definition

```go
type ResponseHook func(resp *http.Response) error
```

### Constructor Functions

### AddResponseCache

AddResponseCache implements response caching

```go
func AddResponseCache(cache Cache) ResponseHook
```

**Parameters:**
- `cache` (Cache)

**Returns:**
- ResponseHook

### AddResponseCompression

AddResponseCompression creates a response compression middleware

```go
func AddResponseCompression(config CompressionConfig) ResponseHook
```

**Parameters:**
- `config` (CompressionConfig)

**Returns:**
- ResponseHook

### AddResponseDebug

AddResponseDebug creates a debug response logging middleware

```go
func AddResponseDebug(config LoggingConfig) ResponseHook
```

**Parameters:**
- `config` (LoggingConfig)

**Returns:**
- ResponseHook

### AddResponseDecompression

AddResponseDecompression creates a response decompression middleware

```go
func AddResponseDecompression() ResponseHook
```

**Parameters:**
  None

**Returns:**
- ResponseHook

### AddResponseErrorLogging

AddResponseErrorLogging creates an error logging middleware

```go
func AddResponseErrorLogging(config LoggingConfig) ResponseHook
```

**Parameters:**
- `config` (LoggingConfig)

**Returns:**
- ResponseHook

### AddResponseLogging

AddResponseLogging creates a response logging middleware

```go
func AddResponseLogging(config LoggingConfig) ResponseHook
```

**Parameters:**
- `config` (LoggingConfig)

**Returns:**
- ResponseHook

### AddResponseMetrics

AddResponseMetrics creates a response metrics collection middleware

```go
func AddResponseMetrics(collector *MetricsCollector) ResponseHook
```

**Parameters:**
- `collector` (*MetricsCollector)

**Returns:**
- ResponseHook

### AddResponseRequestID

AddResponseRequestID creates a response middleware that logs request ID

```go
func AddResponseRequestID(config RequestIDConfig) ResponseHook
```

**Parameters:**
- `config` (RequestIDConfig)

**Returns:**
- ResponseHook

### AddResponseStructured

AddResponseStructured creates a structured response logging middleware

```go
func AddResponseStructured(config LoggingConfig) ResponseHook
```

**Parameters:**
- `config` (LoggingConfig)

**Returns:**
- ResponseHook

### AddResponseTimeout

AddResponseTimeout creates a response timeout middleware that checks if a request took too long.

```go
func AddResponseTimeout(timeout time.Duration) ResponseHook
```

**Parameters:**
- `timeout` (time.Duration)

**Returns:**
- ResponseHook

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| Method | `HTTPMethod` |  |
| Path | `string` |  |

### Constructor Functions

### NewRoute

NewRoute creates a new type-safe route

```go
func NewRoute(method HTTPMethod, path string) *ast.IndexListExpr
```

**Parameters:**
- `method` (HTTPMethod)
- `path` (string)

**Returns:**
- *ast.IndexListExpr

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

| Field | Type | Description |
| ----- | ---- | ----------- |
| Token | `string` |  |
| Prefix | `string` |  |

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
func AddAutoMetrics() (RequestHook, ResponseHook, func() MetricsSnapshot)
```

**Parameters:**
None

**Returns:**
| Type | Description |
|------|-------------|
| `RequestHook` | |
| `ResponseHook` | |
| `func() MetricsSnapshot` | |

**Example:**

```go
// Example usage of AddAutoMetrics
result := AddAutoMetrics(/* parameters */)
```

### GetRequestID
GetRequestID extracts the request ID from context

```go
func GetRequestID(c context.Context) (string, bool)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `c` | `context.Context` | |

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
func WithRequestID(c context.Context, requestID string) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `c` | `context.Context` | |
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
