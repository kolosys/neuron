# generator API

Complete API documentation for the generator package.

**Import Path:** `github.com/kolosys/neuron/internal/generator`

## Package Documentation



## Constants

**TmplModels, TmplClient, TmplPaths, TmplSecurity, TmplResponses, TmplHeaders**

Template names


```go
const TmplModels = "models"
const TmplClient = "client"
const TmplPaths = "paths"
const TmplSecurity = "security"
const TmplResponses = "responses"
const TmplHeaders = "headers"
```

## Types

### AdditionalProperties
AdditionalProperties handles the additionalProperties field which can be boolean or Schema

#### Example Usage

```go
// Create a new AdditionalProperties
additionalproperties := AdditionalProperties{
    Allowed: true,
    Schema: &Schema{}{},
}
```

#### Type Definition

```go
type AdditionalProperties struct {
    Allowed bool
    Schema *Schema
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Allowed | `bool` |  |
| Schema | `*Schema` |  |

## Methods

### GetSchema

GetSchema returns the schema if present, or nil

```go
func (*AdditionalProperties) GetSchema() *Schema
```

**Parameters:**
  None

**Returns:**
- *Schema

### IsAllowed

IsAllowed returns true if additional properties are allowed

```go
func (*AdditionalProperties) IsAllowed() bool
```

**Parameters:**
  None

**Returns:**
- bool

### UnmarshalJSON

UnmarshalJSON implements custom unmarshaling for AdditionalProperties

```go
func (*AdditionalProperties) UnmarshalJSON(data []byte) error
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- error

### Callback
Callback is a map of expressions to path items

#### Example Usage

```go
// Example usage of Callback
var value Callback
// Initialize with appropriate value
```

#### Type Definition

```go
type Callback map[string]*PathItem
```

### Components
Components holds reusable objects for the API

#### Example Usage

```go
// Create a new Components
components := Components{
    Schemas: map[],
    Responses: map[],
    Parameters: map[],
    Examples: map[],
    RequestBodies: map[],
    Headers: map[],
    SecuritySchemes: map[],
    Links: map[],
    Callbacks: map[],
}
```

#### Type Definition

```go
type Components struct {
    Schemas map[string]*Schema `json:"schemas,omitempty"`
    Responses map[string]*Response `json:"responses,omitempty"`
    Parameters map[string]*Parameter `json:"parameters,omitempty"`
    Examples map[string]*Example `json:"examples,omitempty"`
    RequestBodies map[string]*RequestBody `json:"requestBodies,omitempty"`
    Headers map[string]*Header `json:"headers,omitempty"`
    SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty"`
    Links map[string]*Link `json:"links,omitempty"`
    Callbacks map[string]*Callback `json:"callbacks,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Schemas | `map[string]*Schema` |  |
| Responses | `map[string]*Response` |  |
| Parameters | `map[string]*Parameter` |  |
| Examples | `map[string]*Example` |  |
| RequestBodies | `map[string]*RequestBody` |  |
| Headers | `map[string]*Header` |  |
| SecuritySchemes | `map[string]*SecurityScheme` |  |
| Links | `map[string]*Link` |  |
| Callbacks | `map[string]*Callback` |  |

### Contact
Contact provides contact information for the API

#### Example Usage

```go
// Create a new Contact
contact := Contact{
    Name: "example",
    URL: "example",
    Email: "example",
}
```

#### Type Definition

```go
type Contact struct {
    Name string `json:"name,omitempty"`
    URL string `json:"url,omitempty"`
    Email string `json:"email,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| URL | `string` |  |
| Email | `string` |  |

### Discriminator
Discriminator is used for polymorphism

#### Example Usage

```go
// Create a new Discriminator
discriminator := Discriminator{
    PropertyName: "example",
    Mapping: map[],
}
```

#### Type Definition

```go
type Discriminator struct {
    PropertyName string `json:"propertyName"`
    Mapping map[string]string `json:"mapping,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| PropertyName | `string` |  |
| Mapping | `map[string]string` |  |

### Encoding
Encoding describes encoding for a media type property

#### Example Usage

```go
// Create a new Encoding
encoding := Encoding{
    ContentType: "example",
    Headers: map[],
    Style: "example",
    Explode: &true{},
    AllowReserved: true,
}
```

#### Type Definition

```go
type Encoding struct {
    ContentType string `json:"contentType,omitempty"`
    Headers map[string]*Header `json:"headers,omitempty"`
    Style string `json:"style,omitempty"`
    Explode *bool `json:"explode,omitempty"`
    AllowReserved bool `json:"allowReserved,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ContentType | `string` |  |
| Headers | `map[string]*Header` |  |
| Style | `string` |  |
| Explode | `*bool` |  |
| AllowReserved | `bool` |  |

### EnumValue
EnumValue represents an enum constant

#### Example Usage

```go
// Create a new EnumValue
enumvalue := EnumValue{
    Name: "example",
    Value: any{},
    Type: "example",
    Description: "example",
}
```

#### Type Definition

```go
type EnumValue struct {
    Name string
    Value any
    Type string
    Description string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Value | `any` |  |
| Type | `string` |  |
| Description | `string` |  |

### Example
Example provides an example value

#### Example Usage

```go
// Create a new Example
example := Example{
    Ref: "example",
    Summary: "example",
    Description: "example",
    Value: any{},
    ExternalValue: "example",
}
```

#### Type Definition

```go
type Example struct {
    Ref string `json:"$ref,omitempty"`
    Summary string `json:"summary,omitempty"`
    Description string `json:"description,omitempty"`
    Value any `json:"value,omitempty"`
    ExternalValue string `json:"externalValue,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Summary | `string` |  |
| Description | `string` |  |
| Value | `any` |  |
| ExternalValue | `string` |  |

### ExternalDocs
ExternalDocs allows referencing external documentation

#### Example Usage

```go
// Create a new ExternalDocs
externaldocs := ExternalDocs{
    Description: "example",
    URL: "example",
}
```

#### Type Definition

```go
type ExternalDocs struct {
    Description string `json:"description,omitempty"`
    URL string `json:"url"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Description | `string` |  |
| URL | `string` |  |

### Generator
Generator handles OpenAPI to Go code generation

#### Example Usage

```go
// Create a new Generator
generator := Generator{

}
```

#### Type Definition

```go
type Generator struct {
}
```

### Constructor Functions

### New

New creates a new Generator with the given options

```go
func New(opts Options) *Generator
```

**Parameters:**
- `opts` (Options)

**Returns:**
- *Generator

## Methods

### Generate

Generate parses the OpenAPI spec and generates Go code

```go
func (*Generator) Generate() error
```

**Parameters:**
  None

**Returns:**
- error

### GenerateWithContext

GenerateWithContext parses the OpenAPI spec and generates Go code with context support

```go
func (*Generator) GenerateWithContext(ctx context.Context) error
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- error

### GoField
GoField represents a struct field

#### Example Usage

```go
// Create a new GoField
gofield := GoField{
    Name: "example",
    Type: &GoType{}{},
    JSONName: "example",
    Description: "example",
    Required: true,
    OmitEmpty: true,
    Deprecated: true,
    ReadOnly: true,
    WriteOnly: true,
    Validation: "example",
}
```

#### Type Definition

```go
type GoField struct {
    Name string
    Type *GoType
    JSONName string
    Description string
    Required bool
    OmitEmpty bool
    Deprecated bool
    ReadOnly bool
    WriteOnly bool
    Validation string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Type | `*GoType` |  |
| JSONName | `string` |  |
| Description | `string` |  |
| Required | `bool` |  |
| OmitEmpty | `bool` |  |
| Deprecated | `bool` |  |
| ReadOnly | `bool` |  |
| WriteOnly | `bool` |  |
| Validation | `string` |  |

### GoStruct
GoStruct represents a Go struct type

#### Example Usage

```go
// Create a new GoStruct
gostruct := GoStruct{
    Name: "example",
    Description: "example",
    Fields: [],
    Embeds: [],
    IsEnum: true,
    EnumValues: [],
    Group: "example",
}
```

#### Type Definition

```go
type GoStruct struct {
    Name string
    Description string
    Fields []GoField
    Embeds []string
    IsEnum bool
    EnumValues []EnumValue
    Group string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Description | `string` |  |
| Fields | `[]GoField` |  |
| Embeds | `[]string` |  |
| IsEnum | `bool` |  |
| EnumValues | `[]EnumValue` |  |
| Group | `string` |  |

### GoType
GoType represents a normalized Go type

#### Example Usage

```go
// Create a new GoType
gotype := GoType{
    Name: "example",
    Package: "example",
    IsPointer: true,
    IsSlice: true,
    IsMap: true,
    KeyType: &GoType{}{},
    ElementType: &GoType{}{},
}
```

#### Type Definition

```go
type GoType struct {
    Name string
    Package string
    IsPointer bool
    IsSlice bool
    IsMap bool
    KeyType *GoType
    ElementType *GoType
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Package | `string` |  |
| IsPointer | `bool` |  |
| IsSlice | `bool` |  |
| IsMap | `bool` |  |
| KeyType | `*GoType` |  |
| ElementType | `*GoType` |  |

## Methods

### String

String returns the Go type as a string

```go
func (SchemaType) String() string
```

**Parameters:**
  None

**Returns:**
- string

### HTTPMethod
HTTPMethod represents an HTTP method

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

### Header
Header represents an HTTP header

#### Example Usage

```go
// Create a new Header
header := Header{
    Ref: "example",
    Description: "example",
    Required: true,
    Deprecated: true,
    AllowEmptyValue: true,
    Style: "example",
    Explode: &true{},
    AllowReserved: true,
    Schema: &Schema{}{},
    Example: any{},
    Examples: map[],
    Content: map[],
}
```

#### Type Definition

```go
type Header struct {
    Ref string `json:"$ref,omitempty"`
    Description string `json:"description,omitempty"`
    Required bool `json:"required,omitempty"`
    Deprecated bool `json:"deprecated,omitempty"`
    AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`
    Style string `json:"style,omitempty"`
    Explode *bool `json:"explode,omitempty"`
    AllowReserved bool `json:"allowReserved,omitempty"`
    Schema *Schema `json:"schema,omitempty"`
    Example any `json:"example,omitempty"`
    Examples map[string]*Example `json:"examples,omitempty"`
    Content map[string]*MediaType `json:"content,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Description | `string` |  |
| Required | `bool` |  |
| Deprecated | `bool` |  |
| AllowEmptyValue | `bool` |  |
| Style | `string` |  |
| Explode | `*bool` |  |
| AllowReserved | `bool` |  |
| Schema | `*Schema` |  |
| Example | `any` |  |
| Examples | `map[string]*Example` |  |
| Content | `map[string]*MediaType` |  |

### HeaderDefinition
HeaderDefinition represents a reusable header

#### Example Usage

```go
// Create a new HeaderDefinition
headerdefinition := HeaderDefinition{
    Name: "example",
    GoName: "example",
    Type: &GoType{}{},
    Required: true,
    Description: "example",
}
```

#### Type Definition

```go
type HeaderDefinition struct {
    Name string
    GoName string
    Type *GoType
    Required bool
    Description string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| GoName | `string` |  |
| Type | `*GoType` |  |
| Required | `bool` |  |
| Description | `string` |  |

### Info
Info provides metadata about the API

#### Example Usage

```go
// Create a new Info
info := Info{
    Title: "example",
    Description: "example",
    TermsOfService: "example",
    Contact: &Contact{}{},
    License: &License{}{},
    Version: "example",
}
```

#### Type Definition

```go
type Info struct {
    Title string `json:"title"`
    Description string `json:"description,omitempty"`
    TermsOfService string `json:"termsOfService,omitempty"`
    Contact *Contact `json:"contact,omitempty"`
    License *License `json:"license,omitempty"`
    Version string `json:"version"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Title | `string` |  |
| Description | `string` |  |
| TermsOfService | `string` |  |
| Contact | `*Contact` |  |
| License | `*License` |  |
| Version | `string` |  |

### License
License provides license information for the API

#### Example Usage

```go
// Create a new License
license := License{
    Name: "example",
    URL: "example",
}
```

#### Type Definition

```go
type License struct {
    Name string `json:"name"`
    URL string `json:"url,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| URL | `string` |  |

### Link
Link describes a possible design-time link for a response

#### Example Usage

```go
// Create a new Link
link := Link{
    Ref: "example",
    OperationRef: "example",
    OperationID: "example",
    Parameters: map[],
    RequestBody: any{},
    Description: "example",
    Server: &Server{}{},
}
```

#### Type Definition

```go
type Link struct {
    Ref string `json:"$ref,omitempty"`
    OperationRef string `json:"operationRef,omitempty"`
    OperationID string `json:"operationId,omitempty"`
    Parameters map[string]any `json:"parameters,omitempty"`
    RequestBody any `json:"requestBody,omitempty"`
    Description string `json:"description,omitempty"`
    Server *Server `json:"server,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| OperationRef | `string` |  |
| OperationID | `string` |  |
| Parameters | `map[string]any` |  |
| RequestBody | `any` |  |
| Description | `string` |  |
| Server | `*Server` |  |

### MediaType
MediaType provides schema and examples for a media type

#### Example Usage

```go
// Create a new MediaType
mediatype := MediaType{
    Schema: &Schema{}{},
    Example: any{},
    Examples: map[],
    Encoding: map[],
}
```

#### Type Definition

```go
type MediaType struct {
    Schema *Schema `json:"schema,omitempty"`
    Example any `json:"example,omitempty"`
    Examples map[string]*Example `json:"examples,omitempty"`
    Encoding map[string]*Encoding `json:"encoding,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Schema | `*Schema` |  |
| Example | `any` |  |
| Examples | `map[string]*Example` |  |
| Encoding | `map[string]*Encoding` |  |

### NamingConvention
NamingConvention defines the naming style for generated code

#### Example Usage

```go
// Example usage of NamingConvention
var value NamingConvention
// Initialize with appropriate value
```

#### Type Definition

```go
type NamingConvention string
```

### OAuthFlow
OAuthFlow configuration details for a supported OAuth flow

#### Example Usage

```go
// Create a new OAuthFlow
oauthflow := OAuthFlow{
    AuthorizationURL: "example",
    TokenURL: "example",
    RefreshURL: "example",
    Scopes: map[],
}
```

#### Type Definition

```go
type OAuthFlow struct {
    AuthorizationURL string `json:"authorizationUrl,omitempty"`
    TokenURL string `json:"tokenUrl,omitempty"`
    RefreshURL string `json:"refreshUrl,omitempty"`
    Scopes map[string]string `json:"scopes,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| AuthorizationURL | `string` |  |
| TokenURL | `string` |  |
| RefreshURL | `string` |  |
| Scopes | `map[string]string` |  |

### OAuthFlows
OAuthFlows allows configuration of OAuth flows

#### Example Usage

```go
// Create a new OAuthFlows
oauthflows := OAuthFlows{
    Implicit: &OAuthFlow{}{},
    Password: &OAuthFlow{}{},
    ClientCredentials: &OAuthFlow{}{},
    AuthorizationCode: &OAuthFlow{}{},
}
```

#### Type Definition

```go
type OAuthFlows struct {
    Implicit *OAuthFlow `json:"implicit,omitempty"`
    Password *OAuthFlow `json:"password,omitempty"`
    ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
    AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Implicit | `*OAuthFlow` |  |
| Password | `*OAuthFlow` |  |
| ClientCredentials | `*OAuthFlow` |  |
| AuthorizationCode | `*OAuthFlow` |  |

### OpenAPI
OpenAPI represents the root OpenAPI v3 specification

#### Example Usage

```go
// Create a new OpenAPI
openapi := OpenAPI{
    OpenAPI: "example",
    Info: Info{},
    Servers: [],
    Paths: map[],
    Components: &Components{}{},
    Security: [],
    Tags: [],
    ExternalDocs: &ExternalDocs{}{},
}
```

#### Type Definition

```go
type OpenAPI struct {
    OpenAPI string `json:"openapi"`
    Info Info `json:"info"`
    Servers []Server `json:"servers,omitempty"`
    Paths map[string]*PathItem `json:"paths"`
    Components *Components `json:"components,omitempty"`
    Security []SecurityRequirement `json:"security,omitempty"`
    Tags []Tag `json:"tags,omitempty"`
    ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| OpenAPI | `string` |  |
| Info | `Info` |  |
| Servers | `[]Server` |  |
| Paths | `map[string]*PathItem` |  |
| Components | `*Components` |  |
| Security | `[]SecurityRequirement` |  |
| Tags | `[]Tag` |  |
| ExternalDocs | `*ExternalDocs` |  |

### Operation
Operation represents an API operation

#### Example Usage

```go
// Create a new Operation
operation := Operation{
    Tags: [],
    Summary: "example",
    Description: "example",
    ExternalDocs: &ExternalDocs{}{},
    OperationID: "example",
    Parameters: [],
    RequestBody: &RequestBody{}{},
    Responses: map[],
    Callbacks: map[],
    Deprecated: true,
    Security: [],
    Servers: [],
}
```

#### Type Definition

```go
type Operation struct {
    Tags []string `json:"tags,omitempty"`
    Summary string `json:"summary,omitempty"`
    Description string `json:"description,omitempty"`
    ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
    OperationID string `json:"operationId,omitempty"`
    Parameters []*Parameter `json:"parameters,omitempty"`
    RequestBody *RequestBody `json:"requestBody,omitempty"`
    Responses map[string]*Response `json:"responses,omitempty"`
    Callbacks map[string]*Callback `json:"callbacks,omitempty"`
    Deprecated bool `json:"deprecated,omitempty"`
    Security []SecurityRequirement `json:"security,omitempty"`
    Servers []Server `json:"servers,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Tags | `[]string` |  |
| Summary | `string` |  |
| Description | `string` |  |
| ExternalDocs | `*ExternalDocs` |  |
| OperationID | `string` |  |
| Parameters | `[]*Parameter` |  |
| RequestBody | `*RequestBody` |  |
| Responses | `map[string]*Response` |  |
| Callbacks | `map[string]*Callback` |  |
| Deprecated | `bool` |  |
| Security | `[]SecurityRequirement` |  |
| Servers | `[]Server` |  |

### Option
Option is a functional option for configuring the generator

#### Example Usage

```go
// Example usage of Option
var value Option
// Initialize with appropriate value
```

#### Type Definition

```go
type Option func(*Options)
```

### Constructor Functions

### WithBaseURL

WithBaseURL sets the default base URL for the client

```go
func WithBaseURL(url string) Option
```

**Parameters:**
- `url` (string)

**Returns:**
- Option

### WithClientPackage

WithClientPackage sets the client package name

```go
func WithClientPackage(name string) Option
```

**Parameters:**
- `name` (string)

**Returns:**
- Option

### WithGenerateHelpers

WithGenerateHelpers enables or disables helper function generation

```go
func WithGenerateHelpers(enabled bool) Option
```

**Parameters:**
- `enabled` (bool)

**Returns:**
- Option

### WithInputPath

WithInputPath sets the input OpenAPI specification path

```go
func WithInputPath(path string) Option
```

**Parameters:**
- `path` (string)

**Returns:**
- Option

### WithInputURL

WithInputURL sets the URL to fetch the OpenAPI specification from

```go
func WithInputURL(url string) Option
```

**Parameters:**
- `url` (string)

**Returns:**
- Option

### WithMethodPrefix

WithMethodPrefix sets the prefix for generated client methods

```go
func WithMethodPrefix(prefix string) Option
```

**Parameters:**
- `prefix` (string)

**Returns:**
- Option

### WithModelsPackage

WithModelsPackage sets the models package name

```go
func WithModelsPackage(name string) Option
```

**Parameters:**
- `name` (string)

**Returns:**
- Option

### WithNamingConvention

WithNamingConvention sets the naming convention for generated types

```go
func WithNamingConvention(convention NamingConvention) Option
```

**Parameters:**
- `convention` (NamingConvention)

**Returns:**
- Option

### WithOmitEmpty

WithOmitEmpty enables or disables omitempty in JSON tags

```go
func WithOmitEmpty(enabled bool) Option
```

**Parameters:**
- `enabled` (bool)

**Returns:**
- Option

### WithOutputPath

WithOutputPath sets the output directory for generated code

```go
func WithOutputPath(path string) Option
```

**Parameters:**
- `path` (string)

**Returns:**
- Option

### WithPackageName

WithPackageName sets the base package name

```go
func WithPackageName(name string) Option
```

**Parameters:**
- `name` (string)

**Returns:**
- Option

### WithValidationTags

WithValidationTags enables or disables validation tags

```go
func WithValidationTags(enabled bool) Option
```

**Parameters:**
- `enabled` (bool)

**Returns:**
- Option

### Options
Options configures the code generator behavior

#### Example Usage

```go
// Create a new Options
options := Options{
    InputPath: "example",
    InputURL: "example",
    OutputPath: "example",
    PackageName: "example",
    ModelsPackage: "example",
    ClientPackage: "example",
    NamingConvention: NamingConvention{},
    OmitEmpty: true,
    ValidationTags: true,
    BaseURL: "example",
    MethodPrefix: "example",
    GenerateHelpers: true,
}
```

#### Type Definition

```go
type Options struct {
    InputPath string
    InputURL string
    OutputPath string
    PackageName string
    ModelsPackage string
    ClientPackage string
    NamingConvention NamingConvention
    OmitEmpty bool
    ValidationTags bool
    BaseURL string
    MethodPrefix string
    GenerateHelpers bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| InputPath | `string` | InputPath is the path to the OpenAPI JSON specification |
| InputURL | `string` | InputURL is the URL to fetch the OpenAPI JSON specification from |
| OutputPath | `string` | OutputPath is the directory where generated code will be written |
| PackageName | `string` | PackageName is the base package name for generated code |
| ModelsPackage | `string` | ModelsPackage is the package name for model types (default: "models") |
| ClientPackage | `string` | ClientPackage is the package name for client code (default: "client") |
| NamingConvention | `NamingConvention` | NamingConvention controls the naming style for generated types |
| OmitEmpty | `bool` | OmitEmpty adds omitempty to JSON tags for optional fields |
| ValidationTags | `bool` | ValidationTags adds validation tags to struct fields |
| BaseURL | `string` | BaseURL is the default base URL for the generated client |
| MethodPrefix | `string` | MethodPrefix is an optional prefix for generated client methods |
| GenerateHelpers | `bool` | GenerateHelpers generates helper functions for common operations |

### Constructor Functions

### DefaultOptions

DefaultOptions returns options with sensible defaults

```go
func DefaultOptions() Options
```

**Parameters:**
  None

**Returns:**
- Options

### NewOptions

NewOptions creates a new Options with the given functional options

```go
func NewOptions(opts ...Option) Options
```

**Parameters:**
- `opts` (...Option)

**Returns:**
- Options

## Methods

### Apply

Apply applies functional options to the Options struct

```go
func (*Options) Apply(opts ...Option)
```

**Parameters:**
- `opts` (...Option)

**Returns:**
  None

### Parameter
Parameter describes a single operation parameter

#### Example Usage

```go
// Create a new Parameter
parameter := Parameter{
    Ref: "example",
    Name: "example",
    In: "example",
    Description: "example",
    Required: true,
    Deprecated: true,
    AllowEmptyValue: true,
    Style: "example",
    Explode: &true{},
    AllowReserved: true,
    Schema: &Schema{}{},
    Example: any{},
    Examples: map[],
    Content: map[],
}
```

#### Type Definition

```go
type Parameter struct {
    Ref string `json:"$ref,omitempty"`
    Name string `json:"name,omitempty"`
    In string `json:"in,omitempty"`
    Description string `json:"description,omitempty"`
    Required bool `json:"required,omitempty"`
    Deprecated bool `json:"deprecated,omitempty"`
    AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`
    Style string `json:"style,omitempty"`
    Explode *bool `json:"explode,omitempty"`
    AllowReserved bool `json:"allowReserved,omitempty"`
    Schema *Schema `json:"schema,omitempty"`
    Example any `json:"example,omitempty"`
    Examples map[string]*Example `json:"examples,omitempty"`
    Content map[string]*MediaType `json:"content,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Name | `string` |  |
| In | `string` |  |
| Description | `string` |  |
| Required | `bool` |  |
| Deprecated | `bool` |  |
| AllowEmptyValue | `bool` |  |
| Style | `string` |  |
| Explode | `*bool` |  |
| AllowReserved | `bool` |  |
| Schema | `*Schema` |  |
| Example | `any` |  |
| Examples | `map[string]*Example` |  |
| Content | `map[string]*MediaType` |  |

### Parser
Parser handles parsing OpenAPI specifications

#### Example Usage

```go
// Create a new Parser
parser := Parser{

}
```

#### Type Definition

```go
type Parser struct {
}
```

### Constructor Functions

### NewParser

NewParser creates a new OpenAPI parser

```go
func NewParser() *Parser
```

**Parameters:**
  None

**Returns:**
- *Parser

## Methods

### Parse

Parse parses an OpenAPI specification from JSON data

```go
func (*Parser) Parse(data []byte) (*OpenAPI, error)
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- *OpenAPI
- error

### ParseFile

ParseFile parses an OpenAPI specification from a file

```go
func (*Parser) ParseFile(path string) (*OpenAPI, error)
```

**Parameters:**
- `path` (string)

**Returns:**
- *OpenAPI
- error

### ParseURL

ParseURL fetches and parses an OpenAPI specification from a URL

```go
func (*Parser) ParseURL(ctx context.Context, url string) (*OpenAPI, error)
```

**Parameters:**
- `ctx` (context.Context)
- `url` (string)

**Returns:**
- *OpenAPI
- error

### Spec

Spec returns the parsed specification

```go
func (*Parser) Spec() *OpenAPI
```

**Parameters:**
  None

**Returns:**
- *OpenAPI

### PathItem
PathItem represents a path item object

#### Example Usage

```go
// Create a new PathItem
pathitem := PathItem{
    Ref: "example",
    Summary: "example",
    Description: "example",
    Get: &Operation{}{},
    Put: &Operation{}{},
    Post: &Operation{}{},
    Delete: &Operation{}{},
    Options: &Operation{}{},
    Head: &Operation{}{},
    Patch: &Operation{}{},
    Trace: &Operation{}{},
    Servers: [],
    Parameters: [],
}
```

#### Type Definition

```go
type PathItem struct {
    Ref string `json:"$ref,omitempty"`
    Summary string `json:"summary,omitempty"`
    Description string `json:"description,omitempty"`
    Get *Operation `json:"get,omitempty"`
    Put *Operation `json:"put,omitempty"`
    Post *Operation `json:"post,omitempty"`
    Delete *Operation `json:"delete,omitempty"`
    Options *Operation `json:"options,omitempty"`
    Head *Operation `json:"head,omitempty"`
    Patch *Operation `json:"patch,omitempty"`
    Trace *Operation `json:"trace,omitempty"`
    Servers []Server `json:"servers,omitempty"`
    Parameters []*Parameter `json:"parameters,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Summary | `string` |  |
| Description | `string` |  |
| Get | `*Operation` |  |
| Put | `*Operation` |  |
| Post | `*Operation` |  |
| Delete | `*Operation` |  |
| Options | `*Operation` |  |
| Head | `*Operation` |  |
| Patch | `*Operation` |  |
| Trace | `*Operation` |  |
| Servers | `[]Server` |  |
| Parameters | `[]*Parameter` |  |

### PathProcessor
PathProcessor handles path processing and route generation

#### Example Usage

```go
// Create a new PathProcessor
pathprocessor := PathProcessor{

}
```

#### Type Definition

```go
type PathProcessor struct {
}
```

### Constructor Functions

### NewPathProcessor

NewPathProcessor creates a new path processor

```go
func NewPathProcessor(opts Options, schemaProcessor *SchemaProcessor) *PathProcessor
```

**Parameters:**
- `opts` (Options)
- `schemaProcessor` (*SchemaProcessor)

**Returns:**
- *PathProcessor

## Methods

### GetHeaders

GetHeaders returns all header definitions

```go
func (*PathProcessor) GetHeaders() []HeaderDefinition
```

**Parameters:**
  None

**Returns:**
- []HeaderDefinition

### GetResponses

GetResponses returns all response definitions

```go
func (*PathProcessor) GetResponses() []ResponseDefinition
```

**Parameters:**
  None

**Returns:**
- []ResponseDefinition

### GetRoutes

GetRoutes returns all processed routes

```go
func (*PathProcessor) GetRoutes() []RouteDefinition
```

**Parameters:**
  None

**Returns:**
- []RouteDefinition

### GetSecurity

GetSecurity returns all security configurations

```go
func (*PathProcessor) GetSecurity() []SecurityConfig
```

**Parameters:**
  None

**Returns:**
- []SecurityConfig

### Process

Process processes all paths and generates route definitions

```go
func (*PathProcessor) Process(spec *OpenAPI) error
```

**Parameters:**
- `spec` (*OpenAPI)

**Returns:**
- error

### RequestBody
RequestBody describes a single request body

#### Example Usage

```go
// Create a new RequestBody
requestbody := RequestBody{
    Ref: "example",
    Description: "example",
    Content: map[],
    Required: true,
}
```

#### Type Definition

```go
type RequestBody struct {
    Ref string `json:"$ref,omitempty"`
    Description string `json:"description,omitempty"`
    Content map[string]*MediaType `json:"content,omitempty"`
    Required bool `json:"required,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Description | `string` |  |
| Content | `map[string]*MediaType` |  |
| Required | `bool` |  |

### Response
Response describes a single response from an API Operation

#### Example Usage

```go
// Create a new Response
response := Response{
    Ref: "example",
    Description: "example",
    Headers: map[],
    Content: map[],
    Links: map[],
}
```

#### Type Definition

```go
type Response struct {
    Ref string `json:"$ref,omitempty"`
    Description string `json:"description,omitempty"`
    Headers map[string]*Header `json:"headers,omitempty"`
    Content map[string]*MediaType `json:"content,omitempty"`
    Links map[string]*Link `json:"links,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Description | `string` |  |
| Headers | `map[string]*Header` |  |
| Content | `map[string]*MediaType` |  |
| Links | `map[string]*Link` |  |

### ResponseDefinition
ResponseDefinition represents a reusable response

#### Example Usage

```go
// Create a new ResponseDefinition
responsedefinition := ResponseDefinition{
    Name: "example",
    Description: "example",
    Type: "example",
    Headers: [],
}
```

#### Type Definition

```go
type ResponseDefinition struct {
    Name string
    Description string
    Type string
    Headers []HeaderDefinition
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Description | `string` |  |
| Type | `string` |  |
| Headers | `[]HeaderDefinition` |  |

### RouteDefinition
RouteDefinition represents a generated route

#### Example Usage

```go
// Create a new RouteDefinition
routedefinition := RouteDefinition{
    Name: "example",
    Method: HTTPMethod{},
    Path: "example",
    OperationID: "example",
    Summary: "example",
    Description: "example",
    Tags: [],
    RequestType: "example",
    ResponseType: "example",
    PathParams: [],
    QueryParams: [],
    HeaderParams: [],
    HasRequestBody: true,
    RequestBodyType: "example",
    Security: [],
    Deprecated: true,
}
```

#### Type Definition

```go
type RouteDefinition struct {
    Name string
    Method HTTPMethod
    Path string
    OperationID string
    Summary string
    Description string
    Tags []string
    RequestType string
    ResponseType string
    PathParams []RouteParam
    QueryParams []RouteParam
    HeaderParams []RouteParam
    HasRequestBody bool
    RequestBodyType string
    Security []SecurityRequirement
    Deprecated bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Method | `HTTPMethod` |  |
| Path | `string` |  |
| OperationID | `string` |  |
| Summary | `string` |  |
| Description | `string` |  |
| Tags | `[]string` |  |
| RequestType | `string` |  |
| ResponseType | `string` |  |
| PathParams | `[]RouteParam` |  |
| QueryParams | `[]RouteParam` |  |
| HeaderParams | `[]RouteParam` |  |
| HasRequestBody | `bool` |  |
| RequestBodyType | `string` |  |
| Security | `[]SecurityRequirement` |  |
| Deprecated | `bool` |  |

### RouteParam
RouteParam represents a route parameter

#### Example Usage

```go
// Create a new RouteParam
routeparam := RouteParam{
    Name: "example",
    GoName: "example",
    Type: &GoType{}{},
    In: "example",
    Required: true,
    Description: "example",
    Style: "example",
}
```

#### Type Definition

```go
type RouteParam struct {
    Name string
    GoName string
    Type *GoType
    In string
    Required bool
    Description string
    Style string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| GoName | `string` |  |
| Type | `*GoType` |  |
| In | `string` |  |
| Required | `bool` |  |
| Description | `string` |  |
| Style | `string` |  |

### Schema
Schema represents a JSON Schema object

#### Example Usage

```go
// Create a new Schema
schema := Schema{
    Ref: "example",
    Type: SchemaType{},
    Format: "example",
    Title: "example",
    Description: "example",
    Default: any{},
    Enum: [],
    Const: any{},
    MultipleOf: &3.14{},
    Maximum: &3.14{},
    ExclusiveMaximum: &3.14{},
    Minimum: &3.14{},
    ExclusiveMinimum: &3.14{},
    MaxLength: &42{},
    MinLength: &42{},
    Pattern: "example",
    MaxItems: &42{},
    MinItems: &42{},
    UniqueItems: true,
    MaxContains: &42{},
    MinContains: &42{},
    MaxProperties: &42{},
    MinProperties: &42{},
    Required: [],
    DependentRequired: map[],
    AllOf: [],
    AnyOf: [],
    OneOf: [],
    Not: &Schema{}{},
    If: &Schema{}{},
    Then: &Schema{}{},
    Else: &Schema{}{},
    Items: &Schema{}{},
    PrefixItems: [],
    Contains: &Schema{}{},
    Properties: map[],
    PatternProperties: map[],
    AdditionalProperties: &AdditionalProperties{}{},
    PropertyNames: &Schema{}{},
    Nullable: true,
    Discriminator: &Discriminator{}{},
    ReadOnly: true,
    WriteOnly: true,
    XML: &XML{}{},
    ExternalDocs: &ExternalDocs{}{},
    Example: any{},
    Deprecated: true,
    Name: "example",
    GoType: "example",
    IsResolved: true,
}
```

#### Type Definition

```go
type Schema struct {
    Ref string `json:"$ref,omitempty"`
    Type SchemaType `json:"type,omitempty"`
    Format string `json:"format,omitempty"`
    Title string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    Default any `json:"default,omitempty"`
    Enum []any `json:"enum,omitempty"`
    Const any `json:"const,omitempty"`
    MultipleOf *float64 `json:"multipleOf,omitempty"`
    Maximum *float64 `json:"maximum,omitempty"`
    ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
    Minimum *float64 `json:"minimum,omitempty"`
    ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
    MaxLength *int `json:"maxLength,omitempty"`
    MinLength *int `json:"minLength,omitempty"`
    Pattern string `json:"pattern,omitempty"`
    MaxItems *int `json:"maxItems,omitempty"`
    MinItems *int `json:"minItems,omitempty"`
    UniqueItems bool `json:"uniqueItems,omitempty"`
    MaxContains *int `json:"maxContains,omitempty"`
    MinContains *int `json:"minContains,omitempty"`
    MaxProperties *int `json:"maxProperties,omitempty"`
    MinProperties *int `json:"minProperties,omitempty"`
    Required []string `json:"required,omitempty"`
    DependentRequired map[string][]string `json:"dependentRequired,omitempty"`
    AllOf []*Schema `json:"allOf,omitempty"`
    AnyOf []*Schema `json:"anyOf,omitempty"`
    OneOf []*Schema `json:"oneOf,omitempty"`
    Not *Schema `json:"not,omitempty"`
    If *Schema `json:"if,omitempty"`
    Then *Schema `json:"then,omitempty"`
    Else *Schema `json:"else,omitempty"`
    Items *Schema `json:"items,omitempty"`
    PrefixItems []*Schema `json:"prefixItems,omitempty"`
    Contains *Schema `json:"contains,omitempty"`
    Properties map[string]*Schema `json:"properties,omitempty"`
    PatternProperties map[string]*Schema `json:"patternProperties,omitempty"`
    AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
    PropertyNames *Schema `json:"propertyNames,omitempty"`
    Nullable bool `json:"nullable,omitempty"`
    Discriminator *Discriminator `json:"discriminator,omitempty"`
    ReadOnly bool `json:"readOnly,omitempty"`
    WriteOnly bool `json:"writeOnly,omitempty"`
    XML *XML `json:"xml,omitempty"`
    ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
    Example any `json:"example,omitempty"`
    Deprecated bool `json:"deprecated,omitempty"`
    Name string `json:"-"`
    GoType string `json:"-"`
    IsResolved bool `json:"-"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Type | `SchemaType` |  |
| Format | `string` |  |
| Title | `string` |  |
| Description | `string` |  |
| Default | `any` |  |
| Enum | `[]any` |  |
| Const | `any` |  |
| MultipleOf | `*float64` |  |
| Maximum | `*float64` |  |
| ExclusiveMaximum | `*float64` |  |
| Minimum | `*float64` |  |
| ExclusiveMinimum | `*float64` |  |
| MaxLength | `*int` |  |
| MinLength | `*int` |  |
| Pattern | `string` |  |
| MaxItems | `*int` |  |
| MinItems | `*int` |  |
| UniqueItems | `bool` |  |
| MaxContains | `*int` |  |
| MinContains | `*int` |  |
| MaxProperties | `*int` |  |
| MinProperties | `*int` |  |
| Required | `[]string` |  |
| DependentRequired | `map[string][]string` |  |
| AllOf | `[]*Schema` |  |
| AnyOf | `[]*Schema` |  |
| OneOf | `[]*Schema` |  |
| Not | `*Schema` |  |
| If | `*Schema` |  |
| Then | `*Schema` |  |
| Else | `*Schema` |  |
| Items | `*Schema` |  |
| PrefixItems | `[]*Schema` |  |
| Contains | `*Schema` |  |
| Properties | `map[string]*Schema` |  |
| PatternProperties | `map[string]*Schema` |  |
| AdditionalProperties | `*AdditionalProperties` |  |
| PropertyNames | `*Schema` |  |
| Nullable | `bool` |  |
| Discriminator | `*Discriminator` |  |
| ReadOnly | `bool` |  |
| WriteOnly | `bool` |  |
| XML | `*XML` |  |
| ExternalDocs | `*ExternalDocs` |  |
| Example | `any` |  |
| Deprecated | `bool` |  |
| Name | `string` | Internal fields for processing |
| GoType | `string` |  |
| IsResolved | `bool` |  |

### SchemaGroup
SchemaGroup represents a group of related schemas

#### Example Usage

```go
// Create a new SchemaGroup
schemagroup := SchemaGroup{
    Name: "example",
    Structs: [],
}
```

#### Type Definition

```go
type SchemaGroup struct {
    Name string
    Structs []GoStruct
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Structs | `[]GoStruct` |  |

### SchemaProcessor
SchemaProcessor handles schema normalization and grouping

#### Example Usage

```go
// Create a new SchemaProcessor
schemaprocessor := SchemaProcessor{

}
```

#### Type Definition

```go
type SchemaProcessor struct {
}
```

### Constructor Functions

### NewSchemaProcessor

NewSchemaProcessor creates a new schema processor

```go
func NewSchemaProcessor(opts Options, schemas map[string]*Schema) *SchemaProcessor
```

**Parameters:**
- `opts` (Options)
- `schemas` (map[string]*Schema)

**Returns:**
- *SchemaProcessor

## Methods

### GetProcessedSchemas

GetProcessedSchemas returns all processed schemas

```go
func (*SchemaProcessor) GetProcessedSchemas() map[string]*GoStruct
```

**Parameters:**
  None

**Returns:**
- map[string]*GoStruct

### Process

Process normalizes and groups all schemas

```go
func (*PathProcessor) Process(spec *OpenAPI) error
```

**Parameters:**
- `spec` (*OpenAPI)

**Returns:**
- error

### SchemaType
SchemaType handles OpenAPI 3.1.0 type field which can be a string or array of strings

#### Example Usage

```go
// Create a new SchemaType
schematype := SchemaType{
    Types: [],
    Nullable: true,
}
```

#### Type Definition

```go
type SchemaType struct {
    Types []string
    Nullable bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Types | `[]string` |  |
| Nullable | `bool` |  |

### Constructor Functions

### NewNullableSchemaType

NewNullableSchemaType creates a nullable SchemaType from a string

```go
func NewNullableSchemaType(t string) SchemaType
```

**Parameters:**
- `t` (string)

**Returns:**
- SchemaType

### NewSchemaType

NewSchemaType creates a SchemaType from a string

```go
func NewSchemaType(t string) SchemaType
```

**Parameters:**
- `t` (string)

**Returns:**
- SchemaType

## Methods

### IsNullable

IsNullable returns whether the type includes null

```go
func (SchemaType) IsNullable() bool
```

**Parameters:**
  None

**Returns:**
- bool

### String

String returns the primary type (first non-null type)

```go
func (*GoType) String() string
```

**Parameters:**
  None

**Returns:**
- string

### UnmarshalJSON

UnmarshalJSON implements custom unmarshaling for SchemaType

```go
func (*AdditionalProperties) UnmarshalJSON(data []byte) error
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- error

### SecurityConfig
SecurityConfig represents a security scheme configuration

#### Example Usage

```go
// Create a new SecurityConfig
securityconfig := SecurityConfig{
    Name: "example",
    Type: "example",
    Scheme: "example",
    BearerFormat: "example",
    In: "example",
    HeaderName: "example",
    Description: "example",
}
```

#### Type Definition

```go
type SecurityConfig struct {
    Name string
    Type string
    Scheme string
    BearerFormat string
    In string
    HeaderName string
    Description string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Type | `string` |  |
| Scheme | `string` |  |
| BearerFormat | `string` |  |
| In | `string` |  |
| HeaderName | `string` |  |
| Description | `string` |  |

### SecurityRequirement
SecurityRequirement lists the required security schemes

#### Example Usage

```go
// Example usage of SecurityRequirement
var value SecurityRequirement
// Initialize with appropriate value
```

#### Type Definition

```go
type SecurityRequirement map[string][]string
```

### SecurityScheme
SecurityScheme defines a security scheme

#### Example Usage

```go
// Create a new SecurityScheme
securityscheme := SecurityScheme{
    Ref: "example",
    Type: "example",
    Description: "example",
    Name: "example",
    In: "example",
    Scheme: "example",
    BearerFormat: "example",
    Flows: &OAuthFlows{}{},
    OpenIDConnectURL: "example",
}
```

#### Type Definition

```go
type SecurityScheme struct {
    Ref string `json:"$ref,omitempty"`
    Type string `json:"type,omitempty"`
    Description string `json:"description,omitempty"`
    Name string `json:"name,omitempty"`
    In string `json:"in,omitempty"`
    Scheme string `json:"scheme,omitempty"`
    BearerFormat string `json:"bearerFormat,omitempty"`
    Flows *OAuthFlows `json:"flows,omitempty"`
    OpenIDConnectURL string `json:"openIdConnectUrl,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Ref | `string` |  |
| Type | `string` |  |
| Description | `string` |  |
| Name | `string` |  |
| In | `string` |  |
| Scheme | `string` |  |
| BearerFormat | `string` |  |
| Flows | `*OAuthFlows` |  |
| OpenIDConnectURL | `string` |  |

### Server
Server represents a server for the API

#### Example Usage

```go
// Create a new Server
server := Server{
    URL: "example",
    Description: "example",
    Variables: map[],
}
```

#### Type Definition

```go
type Server struct {
    URL string `json:"url"`
    Description string `json:"description,omitempty"`
    Variables map[string]ServerVariable `json:"variables,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| URL | `string` |  |
| Description | `string` |  |
| Variables | `map[string]ServerVariable` |  |

### ServerVariable
ServerVariable represents a variable for a server URL template

#### Example Usage

```go
// Create a new ServerVariable
servervariable := ServerVariable{
    Enum: [],
    Default: "example",
    Description: "example",
}
```

#### Type Definition

```go
type ServerVariable struct {
    Enum []string `json:"enum,omitempty"`
    Default string `json:"default"`
    Description string `json:"description,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Enum | `[]string` |  |
| Default | `string` |  |
| Description | `string` |  |

### Tag
Tag adds metadata to a single tag

#### Example Usage

```go
// Create a new Tag
tag := Tag{
    Name: "example",
    Description: "example",
    ExternalDocs: &ExternalDocs{}{},
}
```

#### Type Definition

```go
type Tag struct {
    Name string `json:"name"`
    Description string `json:"description,omitempty"`
    ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Description | `string` |  |
| ExternalDocs | `*ExternalDocs` |  |

### TemplateData
TemplateData holds data for template execution

#### Example Usage

```go
// Create a new TemplateData
templatedata := TemplateData{
    PackageName: "example",
    ModelsPackage: "example",
    ModelsPackageImport: "example",
    ClientPackage: "example",
    BaseURL: "example",
    Groups: map[],
    Routes: [],
    Security: [],
    Headers: [],
    Responses: [],
    Imports: [],
    GeneratedAt: "example",
    Version: "example",
}
```

#### Type Definition

```go
type TemplateData struct {
    PackageName string
    ModelsPackage string
    ModelsPackageImport string
    ClientPackage string
    BaseURL string
    Groups map[string]*SchemaGroup
    Routes []RouteDefinition
    Security []SecurityConfig
    Headers []HeaderDefinition
    Responses []ResponseDefinition
    Imports []string
    GeneratedAt string
    Version string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| PackageName | `string` |  |
| ModelsPackage | `string` |  |
| ModelsPackageImport | `string` |  |
| ClientPackage | `string` |  |
| BaseURL | `string` |  |
| Groups | `map[string]*SchemaGroup` |  |
| Routes | `[]RouteDefinition` |  |
| Security | `[]SecurityConfig` |  |
| Headers | `[]HeaderDefinition` |  |
| Responses | `[]ResponseDefinition` |  |
| Imports | `[]string` |  |
| GeneratedAt | `string` |  |
| Version | `string` |  |

### Templates
Templates holds all code generation templates

#### Example Usage

```go
// Create a new Templates
templates := Templates{

}
```

#### Type Definition

```go
type Templates struct {
}
```

### Constructor Functions

### NewTemplates

NewTemplates creates and parses all templates

```go
func NewTemplates() (*Templates, error)
```

**Parameters:**
  None

**Returns:**
- *Templates
- error

## Methods

### Client

Client returns the client template

```go
func (*Templates) Client() *template.Template
```

**Parameters:**
  None

**Returns:**
- *template.Template

### Headers

Headers returns the headers template

```go
func (*Templates) Headers() *template.Template
```

**Parameters:**
  None

**Returns:**
- *template.Template

### Models

Models returns the models template

```go
func (*Templates) Models() *template.Template
```

**Parameters:**
  None

**Returns:**
- *template.Template

### Paths

Paths returns the paths template

```go
func (*Templates) Paths() *template.Template
```

**Parameters:**
  None

**Returns:**
- *template.Template

### Responses

Responses returns the responses template

```go
func (*Templates) Responses() *template.Template
```

**Parameters:**
  None

**Returns:**
- *template.Template

### Security

Security returns the security template

```go
func (*Templates) Security() *template.Template
```

**Parameters:**
  None

**Returns:**
- *template.Template

### XML
XML provides additional XML-specific information

#### Example Usage

```go
// Create a new XML
xml := XML{
    Name: "example",
    Namespace: "example",
    Prefix: "example",
    Attribute: true,
    Wrapped: true,
}
```

#### Type Definition

```go
type XML struct {
    Name string `json:"name,omitempty"`
    Namespace string `json:"namespace,omitempty"`
    Prefix string `json:"prefix,omitempty"`
    Attribute bool `json:"attribute,omitempty"`
    Wrapped bool `json:"wrapped,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Namespace | `string` |  |
| Prefix | `string` |  |
| Attribute | `bool` |  |
| Wrapped | `bool` |  |

## Functions

### GetRefName
GetRefName extracts the schema name from a $ref string

```go
func GetRefName(ref string) string
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ref` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `string` | |

**Example:**

```go
// Example usage of GetRefName
result := GetRefName(/* parameters */)
```

## External Links

- [Package Overview](../packages/generator.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/neuron/internal/generator)
- [Source Code](https://github.com/kolosys/neuron/tree/dev/internal/generator)
