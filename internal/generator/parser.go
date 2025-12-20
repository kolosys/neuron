package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// OpenAPI represents the root OpenAPI v3 specification
type OpenAPI struct {
	OpenAPI      string                `json:"openapi"`
	Info         Info                  `json:"info"`
	Servers      []Server              `json:"servers,omitempty"`
	Paths        map[string]*PathItem  `json:"paths"`
	Components   *Components           `json:"components,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Tags         []Tag                 `json:"tags,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty"`
}

// Info provides metadata about the API
type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version"`
}

// Contact provides contact information for the API
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License provides license information for the API
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Server represents a server for the API
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

// ServerVariable represents a variable for a server URL template
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// Components holds reusable objects for the API
type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty"`
	Responses       map[string]*Response       `json:"responses,omitempty"`
	Parameters      map[string]*Parameter      `json:"parameters,omitempty"`
	Examples        map[string]*Example        `json:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]*Header         `json:"headers,omitempty"`
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty"`
	Links           map[string]*Link           `json:"links,omitempty"`
	Callbacks       map[string]*Callback       `json:"callbacks,omitempty"`
}

// PathItem represents a path item object
type PathItem struct {
	Ref         string       `json:"$ref,omitempty"`
	Summary     string       `json:"summary,omitempty"`
	Description string       `json:"description,omitempty"`
	Get         *Operation   `json:"get,omitempty"`
	Put         *Operation   `json:"put,omitempty"`
	Post        *Operation   `json:"post,omitempty"`
	Delete      *Operation   `json:"delete,omitempty"`
	Options     *Operation   `json:"options,omitempty"`
	Head        *Operation   `json:"head,omitempty"`
	Patch       *Operation   `json:"patch,omitempty"`
	Trace       *Operation   `json:"trace,omitempty"`
	Servers     []Server     `json:"servers,omitempty"`
	Parameters  []*Parameter `json:"parameters,omitempty"`
}

// Operation represents an API operation
type Operation struct {
	Tags         []string              `json:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"`
	Description  string                `json:"description,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty"`
	Parameters   []*Parameter          `json:"parameters,omitempty"`
	RequestBody  *RequestBody          `json:"requestBody,omitempty"`
	Responses    map[string]*Response  `json:"responses,omitempty"`
	Callbacks    map[string]*Callback  `json:"callbacks,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty"`
	Security     []SecurityRequirement `json:"security,omitempty"`
	Servers      []Server              `json:"servers,omitempty"`
}

// SchemaType handles OpenAPI 3.1.0 type field which can be a string or array of strings
type SchemaType struct {
	Types    []string
	Nullable bool
}

// UnmarshalJSON implements custom unmarshaling for SchemaType
func (st *SchemaType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		st.Types = []string{str}
		return nil
	}

	// Try to unmarshal as an array of strings
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		for _, t := range arr {
			if t == "null" {
				st.Nullable = true
			} else {
				st.Types = append(st.Types, t)
			}
		}
		return nil
	}

	return fmt.Errorf("type must be a string or array of strings")
}

// String returns the primary type (first non-null type)
func (st SchemaType) String() string {
	if len(st.Types) > 0 {
		return st.Types[0]
	}
	return ""
}

// IsNullable returns whether the type includes null
func (st SchemaType) IsNullable() bool {
	return st.Nullable
}

// NewSchemaType creates a SchemaType from a string
func NewSchemaType(t string) SchemaType {
	return SchemaType{Types: []string{t}}
}

// NewNullableSchemaType creates a nullable SchemaType from a string
func NewNullableSchemaType(t string) SchemaType {
	return SchemaType{Types: []string{t}, Nullable: true}
}

// AdditionalProperties handles the additionalProperties field which can be boolean or Schema
type AdditionalProperties struct {
	Allowed bool
	Schema  *Schema
}

// UnmarshalJSON implements custom unmarshaling for AdditionalProperties
func (ap *AdditionalProperties) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a boolean first
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		ap.Allowed = b
		ap.Schema = nil
		return nil
	}

	// Try to unmarshal as a Schema
	var schema Schema
	if err := json.Unmarshal(data, &schema); err == nil {
		ap.Allowed = true
		ap.Schema = &schema
		return nil
	}

	return fmt.Errorf("additionalProperties must be a boolean or schema object")
}

// GetSchema returns the schema if present, or nil
func (ap *AdditionalProperties) GetSchema() *Schema {
	return ap.Schema
}

// IsAllowed returns true if additional properties are allowed
func (ap *AdditionalProperties) IsAllowed() bool {
	return ap.Allowed
}

// Schema represents a JSON Schema object
type Schema struct {
	Ref                  string                `json:"$ref,omitempty"`
	Type                 SchemaType            `json:"type,omitempty"`
	Format               string                `json:"format,omitempty"`
	Title                string                `json:"title,omitempty"`
	Description          string                `json:"description,omitempty"`
	Default              any                   `json:"default,omitempty"`
	Enum                 []any                 `json:"enum,omitempty"`
	Const                any                   `json:"const,omitempty"`
	MultipleOf           *float64              `json:"multipleOf,omitempty"`
	Maximum              *float64              `json:"maximum,omitempty"`
	ExclusiveMaximum     *float64              `json:"exclusiveMaximum,omitempty"`
	Minimum              *float64              `json:"minimum,omitempty"`
	ExclusiveMinimum     *float64              `json:"exclusiveMinimum,omitempty"`
	MaxLength            *int                  `json:"maxLength,omitempty"`
	MinLength            *int                  `json:"minLength,omitempty"`
	Pattern              string                `json:"pattern,omitempty"`
	MaxItems             *int                  `json:"maxItems,omitempty"`
	MinItems             *int                  `json:"minItems,omitempty"`
	UniqueItems          bool                  `json:"uniqueItems,omitempty"`
	MaxContains          *int                  `json:"maxContains,omitempty"`
	MinContains          *int                  `json:"minContains,omitempty"`
	MaxProperties        *int                  `json:"maxProperties,omitempty"`
	MinProperties        *int                  `json:"minProperties,omitempty"`
	Required             []string              `json:"required,omitempty"`
	DependentRequired    map[string][]string   `json:"dependentRequired,omitempty"`
	AllOf                []*Schema             `json:"allOf,omitempty"`
	AnyOf                []*Schema             `json:"anyOf,omitempty"`
	OneOf                []*Schema             `json:"oneOf,omitempty"`
	Not                  *Schema               `json:"not,omitempty"`
	If                   *Schema               `json:"if,omitempty"`
	Then                 *Schema               `json:"then,omitempty"`
	Else                 *Schema               `json:"else,omitempty"`
	Items                *Schema               `json:"items,omitempty"`
	PrefixItems          []*Schema             `json:"prefixItems,omitempty"`
	Contains             *Schema               `json:"contains,omitempty"`
	Properties           map[string]*Schema    `json:"properties,omitempty"`
	PatternProperties    map[string]*Schema    `json:"patternProperties,omitempty"`
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
	PropertyNames        *Schema               `json:"propertyNames,omitempty"`
	Nullable             bool                  `json:"nullable,omitempty"`
	Discriminator        *Discriminator        `json:"discriminator,omitempty"`
	ReadOnly             bool                  `json:"readOnly,omitempty"`
	WriteOnly            bool                  `json:"writeOnly,omitempty"`
	XML                  *XML                  `json:"xml,omitempty"`
	ExternalDocs         *ExternalDocs         `json:"externalDocs,omitempty"`
	Example              any                   `json:"example,omitempty"`
	Deprecated           bool                  `json:"deprecated,omitempty"`

	// Internal fields for processing
	Name       string `json:"-"`
	GoType     string `json:"-"`
	IsResolved bool   `json:"-"`
}

// Discriminator is used for polymorphism
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// XML provides additional XML-specific information
type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

// Parameter describes a single operation parameter
type Parameter struct {
	Ref             string                `json:"$ref,omitempty"`
	Name            string                `json:"name,omitempty"`
	In              string                `json:"in,omitempty"`
	Description     string                `json:"description,omitempty"`
	Required        bool                  `json:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty"`
	Explode         *bool                 `json:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty"`
	Schema          *Schema               `json:"schema,omitempty"`
	Example         any                   `json:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
}

// RequestBody describes a single request body
type RequestBody struct {
	Ref         string                `json:"$ref,omitempty"`
	Description string                `json:"description,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
	Required    bool                  `json:"required,omitempty"`
}

// Response describes a single response from an API Operation
type Response struct {
	Ref         string                `json:"$ref,omitempty"`
	Description string                `json:"description,omitempty"`
	Headers     map[string]*Header    `json:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty"`
}

// Header represents an HTTP header
type Header struct {
	Ref             string                `json:"$ref,omitempty"`
	Description     string                `json:"description,omitempty"`
	Required        bool                  `json:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty"`
	Explode         *bool                 `json:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty"`
	Schema          *Schema               `json:"schema,omitempty"`
	Example         any                   `json:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
}

// MediaType provides schema and examples for a media type
type MediaType struct {
	Schema   *Schema              `json:"schema,omitempty"`
	Example  any                  `json:"example,omitempty"`
	Examples map[string]*Example  `json:"examples,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty"`
}

// Encoding describes encoding for a media type property
type Encoding struct {
	ContentType   string             `json:"contentType,omitempty"`
	Headers       map[string]*Header `json:"headers,omitempty"`
	Style         string             `json:"style,omitempty"`
	Explode       *bool              `json:"explode,omitempty"`
	AllowReserved bool               `json:"allowReserved,omitempty"`
}

// Example provides an example value
type Example struct {
	Ref           string `json:"$ref,omitempty"`
	Summary       string `json:"summary,omitempty"`
	Description   string `json:"description,omitempty"`
	Value         any    `json:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty"`
}

// Link describes a possible design-time link for a response
type Link struct {
	Ref          string         `json:"$ref,omitempty"`
	OperationRef string         `json:"operationRef,omitempty"`
	OperationID  string         `json:"operationId,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	RequestBody  any            `json:"requestBody,omitempty"`
	Description  string         `json:"description,omitempty"`
	Server       *Server        `json:"server,omitempty"`
}

// Callback is a map of expressions to path items
type Callback map[string]*PathItem

// SecurityScheme defines a security scheme
type SecurityScheme struct {
	Ref              string      `json:"$ref,omitempty"`
	Type             string      `json:"type,omitempty"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty"`
}

// OAuthFlows allows configuration of OAuth flows
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// OAuthFlow configuration details for a supported OAuth flow
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// SecurityRequirement lists the required security schemes
type SecurityRequirement map[string][]string

// Tag adds metadata to a single tag
type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// ExternalDocs allows referencing external documentation
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

// Parser handles parsing OpenAPI specifications
type Parser struct {
	spec *OpenAPI
}

// NewParser creates a new OpenAPI parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses an OpenAPI specification from a file
func (p *Parser) ParseFile(path string) (*OpenAPI, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return p.Parse(data)
}

// ParseURL fetches and parses an OpenAPI specification from a URL
func (p *Parser) ParseURL(ctx context.Context, url string) (*OpenAPI, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "neuron-cli/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spec from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return p.Parse(data)
}

// Parse parses an OpenAPI specification from JSON data
func (p *Parser) Parse(data []byte) (*OpenAPI, error) {
	var spec OpenAPI
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	if !strings.HasPrefix(spec.OpenAPI, "3.") {
		return nil, fmt.Errorf("unsupported OpenAPI version: %s (only 3.x is supported)", spec.OpenAPI)
	}

	p.spec = &spec

	if err := p.resolveReferences(); err != nil {
		return nil, fmt.Errorf("failed to resolve references: %w", err)
	}

	return &spec, nil
}

// Spec returns the parsed specification
func (p *Parser) Spec() *OpenAPI {
	return p.spec
}

// resolveReferences resolves $ref references in the specification
func (p *Parser) resolveReferences() error {
	if p.spec.Components == nil || p.spec.Components.Schemas == nil {
		return nil
	}

	for name, schema := range p.spec.Components.Schemas {
		schema.Name = name
		if err := p.resolveSchemaRefs(schema); err != nil {
			return fmt.Errorf("failed to resolve schema %s: %w", name, err)
		}
	}

	return nil
}

// resolveSchemaRefs resolves references within a schema
func (p *Parser) resolveSchemaRefs(schema *Schema) error {
	if schema == nil || schema.IsResolved {
		return nil
	}

	schema.IsResolved = true

	if schema.Ref != "" {
		resolved, err := p.resolveRef(schema.Ref)
		if err != nil {
			return err
		}
		if resolved != nil {
			*schema = *resolved
			schema.IsResolved = true
		}
	}

	for _, s := range schema.AllOf {
		if err := p.resolveSchemaRefs(s); err != nil {
			return err
		}
	}

	for _, s := range schema.AnyOf {
		if err := p.resolveSchemaRefs(s); err != nil {
			return err
		}
	}

	for _, s := range schema.OneOf {
		if err := p.resolveSchemaRefs(s); err != nil {
			return err
		}
	}

	if schema.Items != nil {
		if err := p.resolveSchemaRefs(schema.Items); err != nil {
			return err
		}
	}

	for _, prop := range schema.Properties {
		if err := p.resolveSchemaRefs(prop); err != nil {
			return err
		}
	}

	if schema.AdditionalProperties != nil && schema.AdditionalProperties.GetSchema() != nil {
		if err := p.resolveSchemaRefs(schema.AdditionalProperties.GetSchema()); err != nil {
			return err
		}
	}

	return nil
}

// resolveRef resolves a $ref string to its schema
func (p *Parser) resolveRef(ref string) (*Schema, error) {
	if !strings.HasPrefix(ref, "#/components/schemas/") {
		return nil, nil
	}

	name := strings.TrimPrefix(ref, "#/components/schemas/")
	if p.spec.Components == nil || p.spec.Components.Schemas == nil {
		return nil, fmt.Errorf("schema not found: %s", name)
	}

	schema, ok := p.spec.Components.Schemas[name]
	if !ok {
		return nil, fmt.Errorf("schema not found: %s", name)
	}

	return schema, nil
}

// GetRefName extracts the schema name from a $ref string
func GetRefName(ref string) string {
	if strings.HasPrefix(ref, "#/components/schemas/") {
		return strings.TrimPrefix(ref, "#/components/schemas/")
	}
	if strings.HasPrefix(ref, "#/components/responses/") {
		return strings.TrimPrefix(ref, "#/components/responses/")
	}
	if strings.HasPrefix(ref, "#/components/parameters/") {
		return strings.TrimPrefix(ref, "#/components/parameters/")
	}
	if strings.HasPrefix(ref, "#/components/requestBodies/") {
		return strings.TrimPrefix(ref, "#/components/requestBodies/")
	}
	if strings.HasPrefix(ref, "#/components/headers/") {
		return strings.TrimPrefix(ref, "#/components/headers/")
	}
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}
