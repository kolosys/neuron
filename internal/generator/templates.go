package generator

import (
	"strings"
	"text/template"
)

// Template names
const (
	TmplModels    = "models"
	TmplClient    = "client"
	TmplPaths     = "paths"
	TmplSecurity  = "security"
	TmplResponses = "responses"
	TmplHeaders   = "headers"
)

// TemplateData holds data for template execution
type TemplateData struct {
	PackageName         string
	ModelsPackage       string
	ModelsPackageImport string
	ClientPackage       string
	BaseURL             string
	Groups              map[string]*SchemaGroup
	Routes              []RouteDefinition
	Security            []SecurityConfig
	Headers             []HeaderDefinition
	Responses           []ResponseDefinition
	Imports             []string
	GeneratedAt         string
	Version             string
}

// Templates holds all code generation templates
type Templates struct {
	models    *template.Template
	client    *template.Template
	paths     *template.Template
	security  *template.Template
	responses *template.Template
	headers   *template.Template
}

// NewTemplates creates and parses all templates
func NewTemplates() (*Templates, error) {
	funcs := template.FuncMap{
		"lower":            strings.ToLower,
		"upper":            strings.ToUpper,
		"title":            strings.Title,
		"hasPrefix":        strings.HasPrefix,
		"hasSuffix":        strings.HasSuffix,
		"trimPrefix":       strings.TrimPrefix,
		"trimSuffix":       strings.TrimSuffix,
		"replace":          strings.ReplaceAll,
		"join":             strings.Join,
		"quote":            func(s string) string { return `"` + s + `"` },
		"backtick":         func(s string) string { return "`" + s + "`" },
		"add":              func(a, b int) int { return a + b },
		"sub":              func(a, b int) int { return a - b },
		"isLast":           func(i, length int) bool { return i == length-1 },
		"formatType":       formatGoType,
		"formatClientType": formatClientGoType,
		"jsonTag":          formatJSONTag,
		"validTag":         formatValidationTag,
	}

	t := &Templates{}
	var err error

	t.models, err = template.New(TmplModels).Funcs(funcs).Parse(modelsTemplate)
	if err != nil {
		return nil, err
	}

	t.client, err = template.New(TmplClient).Funcs(funcs).Parse(clientTemplate)
	if err != nil {
		return nil, err
	}

	t.paths, err = template.New(TmplPaths).Funcs(funcs).Parse(pathsTemplate)
	if err != nil {
		return nil, err
	}

	t.security, err = template.New(TmplSecurity).Funcs(funcs).Parse(securityTemplate)
	if err != nil {
		return nil, err
	}

	t.responses, err = template.New(TmplResponses).Funcs(funcs).Parse(responsesTemplate)
	if err != nil {
		return nil, err
	}

	t.headers, err = template.New(TmplHeaders).Funcs(funcs).Parse(headersTemplate)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Models returns the models template
func (t *Templates) Models() *template.Template { return t.models }

// Client returns the client template
func (t *Templates) Client() *template.Template { return t.client }

// Paths returns the paths template
func (t *Templates) Paths() *template.Template { return t.paths }

// Security returns the security template
func (t *Templates) Security() *template.Template { return t.security }

// Responses returns the responses template
func (t *Templates) Responses() *template.Template { return t.responses }

// Headers returns the headers template
func (t *Templates) Headers() *template.Template { return t.headers }

func formatGoType(t *GoType) string {
	if t == nil {
		return "any"
	}
	return t.String()
}

// formatClientGoType formats a type for use in client code, adding models. prefix for custom types
func formatClientGoType(t *GoType, modelsPackage string) string {
	if t == nil {
		return "any"
	}

	typeStr := t.String()

	// Check if it's a primitive type that doesn't need a prefix
	primitiveTypes := map[string]bool{
		"string": true, "int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"float32": true, "float64": true, "bool": true, "byte": true, "rune": true,
		"any": true, "error": true,
	}

	// Handle pointer types
	baseType := strings.TrimPrefix(typeStr, "*")

	// Handle slice types
	if after, ok := strings.CutPrefix(baseType, "[]"); ok {
		baseType = after
	}

	// Handle map types - for simplicity, keep as-is if already has package
	if strings.HasPrefix(baseType, "map[") {
		return typeStr
	}

	// If already has a package prefix or is a primitive type, return as-is
	if strings.Contains(baseType, ".") || primitiveTypes[baseType] {
		return typeStr
	}

	// For time.Time and other stdlib types, the package is already set
	if t.Package != "" {
		return typeStr
	}

	// Add models package prefix for custom types
	if strings.HasPrefix(typeStr, "*") {
		return "*" + modelsPackage + "." + strings.TrimPrefix(typeStr, "*")
	}
	if strings.HasPrefix(typeStr, "[]") {
		elemType := strings.TrimPrefix(typeStr, "[]")
		if strings.HasPrefix(elemType, "*") {
			return "[]*" + modelsPackage + "." + strings.TrimPrefix(elemType, "*")
		}
		return "[]" + modelsPackage + "." + elemType
	}
	return modelsPackage + "." + typeStr
}

func formatJSONTag(field GoField) string {
	tag := field.JSONName
	if field.OmitEmpty {
		tag += ",omitempty"
	}
	return tag
}

func formatValidationTag(field GoField) string {
	if field.Validation == "" {
		return ""
	}
	return ` validate:"` + field.Validation + `"`
}

const modelsTemplate = `// Code generated by neuron-cli. DO NOT EDIT.
// Generated at: {{.GeneratedAt}}

package {{.PackageName}}
{{if .Imports}}
import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)
{{end}}
{{range $name, $group := .Groups}}
{{range $struct := .Structs}}
{{if .Description}}// {{.Name}} {{.Description}}
{{else}}// {{.Name}}
{{end -}}
{{if .IsEnum}}type {{.Name}} {{(index .EnumValues 0).Type}}

const (
{{- range .EnumValues}}
	{{.Name}} {{$struct.Name}} = {{if eq .Type "string"}}"{{.Value}}"{{else}}{{.Value}}{{end}}
{{- end}}
)
{{else}}type {{.Name}} struct {
{{- range .Embeds}}
	{{.}}
{{- end}}
{{- range .Fields}}
	{{.Name}} {{formatType .Type}} ` + "`" + `json:"{{jsonTag .}}"{{validTag .}}` + "`" + `{{if .Description}} // {{.Description}}{{end}}
{{- end}}
}
{{end}}
{{end}}
{{end}}`

const clientTemplate = `// Code generated by neuron-cli. DO NOT EDIT.
// Generated at: {{.GeneratedAt}}

package {{.ClientPackage}}

import (
	"context"
	"net/http"
	"time"

	"github.com/kolosys/neuron"
)

// RESTClient wraps a neuron.Client with generated route methods
type RESTClient struct {
	*neuron.Client
	baseURL string
}

// ClientOptions holds client configuration
type ClientOptions struct {
	BaseURL    string
	UserAgent  string
	Headers    http.Header
	Timeout    time.Duration
	MaxRetries int
}

// NewRESTClient creates a new REST client
func NewRESTClient(opts ClientOptions) *RESTClient {
	if opts.UserAgent == "" {
		opts.UserAgent = "neuron-client/1.0"
	}
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 3
	}
	{{if .BaseURL}}if opts.BaseURL == "" {
		opts.BaseURL = "{{.BaseURL}}"
	}
	{{end}}
	client := neuron.NewClient(neuron.ClientOptions{
		BaseURL:    opts.BaseURL,
		UserAgent:  opts.UserAgent,
		Headers:    opts.Headers,
		Timeout:    opts.Timeout,
		MaxRetries: opts.MaxRetries,
	})

	return &RESTClient{
		Client:  client,
		baseURL: opts.BaseURL,
	}
}

// BaseURL returns the client's base URL
func (c *RESTClient) BaseURL() string {
	return c.baseURL
}

// WithContext returns a context-aware request builder
func (c *RESTClient) WithContext(ctx context.Context) *RequestBuilder {
	return &RequestBuilder{
		client: c,
		ctx:    ctx,
	}
}

// RequestBuilder helps build requests with context
type RequestBuilder struct {
	client  *RESTClient
	ctx     context.Context
	headers http.Header
	query   map[string]any
}

// WithHeader adds a header to the request
func (b *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(http.Header)
	}
	b.headers.Set(key, value)
	return b
}

// WithQuery adds a query parameter
func (b *RequestBuilder) WithQuery(key string, value any) *RequestBuilder {
	if b.query == nil {
		b.query = make(map[string]any)
	}
	b.query[key] = value
	return b
}

// Options returns RequestOptions for the builder
func (b *RequestBuilder) Options() *neuron.RequestOptions {
	return &neuron.RequestOptions{
		Context: b.ctx,
		Headers: b.headers,
		Query:   b.query,
	}
}
`

const pathsTemplate = `// Code generated by neuron-cli. DO NOT EDIT.
// Generated at: {{.GeneratedAt}}

package {{.ClientPackage}}

import (
	"context"
	"fmt"
	"strings"

	"github.com/kolosys/neuron"
	"{{.ModelsPackageImport}}"
)

// Route definitions
var (
{{- range .Routes}}
	{{.Name}}Route = neuron.NewRoute[{{if .HasRequestBody}}{{if eq .RequestType "any"}}any{{else}}{{$.ModelsPackage}}.{{.RequestType}}{{end}}{{else}}neuron.EmptyRequest{{end}}, {{if hasPrefix .ResponseType "[]"}}[]{{$.ModelsPackage}}.{{trimPrefix .ResponseType "[]"}}{{else if eq .ResponseType "EmptyResponse"}}neuron.EmptyResponse{{else if eq .ResponseType "any"}}any{{else}}{{$.ModelsPackage}}.{{.ResponseType}}{{end}}](
		neuron.Method{{.Method}},
		"{{.Path}}",
	)
{{- end}}
)
{{range .Routes}}
{{if .Description}}// {{.Name}} {{.Description}}
{{else if .Summary}}// {{.Name}} {{.Summary}}
{{else}}// {{.Name}} executes {{.Method}} {{.Path}}
{{end -}}
{{if .Deprecated}}// Deprecated: This operation is deprecated.
{{end -}}
func (c *RESTClient) {{.Name}}(ctx context.Context{{range .PathParams}}, {{lower .GoName}} {{formatClientType .Type $.ModelsPackage}}{{end}}{{if .HasRequestBody}}, body {{if eq .RequestType "any"}}any{{else}}{{$.ModelsPackage}}.{{.RequestType}}{{end}}{{end}}, opts ...*neuron.RequestOptions) (*neuron.Response[{{if hasPrefix .ResponseType "[]"}}[]{{$.ModelsPackage}}.{{trimPrefix .ResponseType "[]"}}{{else if eq .ResponseType "EmptyResponse"}}neuron.EmptyResponse{{else if eq .ResponseType "any"}}any{{else}}{{$.ModelsPackage}}.{{.ResponseType}}{{end}}], error) {
	path := "{{.Path}}"
{{- range .PathParams}}
	path = strings.Replace(path, "{{"{"}}{{.Name}}{{"}"}}", fmt.Sprintf("%v", {{lower .GoName}}), 1)
{{- end}}

	route := neuron.NewRoute[{{if .HasRequestBody}}{{if eq .RequestType "any"}}any{{else}}{{$.ModelsPackage}}.{{.RequestType}}{{end}}{{else}}neuron.EmptyRequest{{end}}, {{if hasPrefix .ResponseType "[]"}}[]{{$.ModelsPackage}}.{{trimPrefix .ResponseType "[]"}}{{else if eq .ResponseType "EmptyResponse"}}neuron.EmptyResponse{{else if eq .ResponseType "any"}}any{{else}}{{$.ModelsPackage}}.{{.ResponseType}}{{end}}](
		neuron.Method{{.Method}},
		path,
	)

	var options *neuron.RequestOptions
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	} else {
		options = &neuron.RequestOptions{}
	}
	options.Context = ctx

	{{if .HasRequestBody}}return neuron.Execute(c.Client, route, body, options){{else}}return neuron.Execute(c.Client, route, neuron.EmptyRequest{}, options){{end}}
}
{{end}}`

const securityTemplate = `// Code generated by neuron-cli. DO NOT EDIT.
// Generated at: {{.GeneratedAt}}

package {{.ClientPackage}}

import (
	"net/http"
)

// Security scheme types
const (
	SecurityTypeHTTP   = "http"
	SecurityTypeAPIKey = "apiKey"
	SecurityTypeOAuth2 = "oauth2"
	SecurityTypeOpenID = "openIdConnect"
)

// SecurityScheme defines a security scheme
type SecurityScheme struct {
	Type         string
	Scheme       string
	BearerFormat string
	In           string
	Name         string
}
{{range .Security}}
// {{.Name}}Security represents the {{.Name}} security scheme
var {{.Name}}Security = SecurityScheme{
	Type:         "{{.Type}}",
	{{- if .Scheme}}
	Scheme:       "{{.Scheme}}",
	{{- end}}
	{{- if .BearerFormat}}
	BearerFormat: "{{.BearerFormat}}",
	{{- end}}
	{{- if .In}}
	In:           "{{.In}}",
	{{- end}}
	{{- if .HeaderName}}
	Name:         "{{.HeaderName}}",
	{{- end}}
}
{{end}}
// ApplyBearerToken adds a bearer token to the request headers
func ApplyBearerToken(token string) func(*http.Request) error {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// ApplyAPIKey adds an API key to the request
func ApplyAPIKey(name, value, location string) func(*http.Request) error {
	return func(req *http.Request) error {
		switch location {
		case "header":
			req.Header.Set(name, value)
		case "query":
			q := req.URL.Query()
			q.Set(name, value)
			req.URL.RawQuery = q.Encode()
		case "cookie":
			req.AddCookie(&http.Cookie{Name: name, Value: value})
		}
		return nil
	}
}

// ApplyBasicAuth adds basic authentication to the request
func ApplyBasicAuth(username, password string) func(*http.Request) error {
	return func(req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	}
}
`

const responsesTemplate = `// Code generated by neuron-cli. DO NOT EDIT.
// Generated at: {{.GeneratedAt}}

package {{.ClientPackage}}

import (
	"{{.ModelsPackageImport}}"
)

// Response type aliases for common responses
type (
	// EmptyResponse represents a response with no body
	EmptyResponse = struct{}
)
{{range .Responses}}
// {{.Name}}Response wraps the {{.Name}} response
type {{.Name}}Response struct {
	{{- if .Type}}
	Data {{$.ModelsPackage}}.{{.Type}}
	{{- end}}
	{{- range .Headers}}
	{{.GoName}} {{formatType .Type}} ` + "`" + `header:"{{.Name}}"` + "`" + `{{if .Description}} // {{.Description}}{{end}}
	{{- end}}
}
{{end}}`

const headersTemplate = `// Code generated by neuron-cli. DO NOT EDIT.
// Generated at: {{.GeneratedAt}}

package {{.ClientPackage}}

// Standard header names
const (
	HeaderContentType     = "Content-Type"
	HeaderAccept          = "Accept"
	HeaderAuthorization   = "Authorization"
	HeaderUserAgent       = "User-Agent"
	HeaderXRequestID      = "X-Request-ID"
	HeaderXCorrelationID  = "X-Correlation-ID"
)

// Content types
const (
	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
	ContentTypeForm = "application/x-www-form-urlencoded"
)
{{if .Headers}}
// API-specific headers
const (
{{- range .Headers}}
	Header{{.GoName}} = "{{.Name}}"
{{- end}}
)

// HeaderValues holds typed header values
type HeaderValues struct {
{{- range .Headers}}
	{{.GoName}} {{formatType .Type}} {{if .Description}}// {{.Description}}{{end}}
{{- end}}
}
{{end}}`
