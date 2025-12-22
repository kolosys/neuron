package generator

import (
	"strings"
)

// NamingConvention defines the naming style for generated code
type NamingConvention string

const (
	NamingPascalCase NamingConvention = "PascalCase"
	NamingCamelCase  NamingConvention = "camelCase"
	NamingSnakeCase  NamingConvention = "snake_case"
)

// NameMapping represents a pattern-to-replacement mapping for type names
type NameMapping struct {
	Pattern     string
	Replacement string
	IsGlob      bool
}

// Options configures the code generator behavior
type Options struct {
	// InputPath is the path to the OpenAPI JSON specification
	InputPath string

	// InputURL is the URL to fetch the OpenAPI JSON specification from
	InputURL string

	// OutputPath is the directory where generated code will be written
	OutputPath string

	// PackageName is the base package name for generated code
	PackageName string

	// ModelsPackage is the package name for model types (default: "models")
	ModelsPackage string

	// ClientPackage is the package name for client code (default: "client")
	ClientPackage string

	// NamingConvention controls the naming style for generated types
	NamingConvention NamingConvention

	// OmitEmpty adds omitempty to JSON tags for optional fields
	OmitEmpty bool

	// ValidationTags adds validation tags to struct fields
	ValidationTags bool

	// BaseURL is the default base URL for the generated client
	BaseURL string

	// MethodPrefix is an optional prefix for generated client methods
	MethodPrefix string

	// GenerateHelpers generates helper functions for common operations
	GenerateHelpers bool

	// NameMappings contains pattern-based name remapping rules
	NameMappings []NameMapping

	// ModelsOnly generates only model files, skipping client code generation
	ModelsOnly bool
}

// Option is a functional option for configuring the generator
type Option func(*Options)

// DefaultOptions returns options with sensible defaults
func DefaultOptions() Options {
	return Options{
		PackageName:      "generated",
		ModelsPackage:    "models",
		ClientPackage:    "client",
		NamingConvention: NamingPascalCase,
		OmitEmpty:        true,
		ValidationTags:   false,
		GenerateHelpers:  true,
	}
}

// WithInputPath sets the input OpenAPI specification path
func WithInputPath(path string) Option {
	return func(o *Options) {
		o.InputPath = path
	}
}

// WithInputURL sets the URL to fetch the OpenAPI specification from
func WithInputURL(url string) Option {
	return func(o *Options) {
		o.InputURL = url
	}
}

// WithOutputPath sets the output directory for generated code
func WithOutputPath(path string) Option {
	return func(o *Options) {
		o.OutputPath = path
	}
}

// WithPackageName sets the base package name
func WithPackageName(name string) Option {
	return func(o *Options) {
		o.PackageName = name
	}
}

// WithModelsPackage sets the models package name
func WithModelsPackage(name string) Option {
	return func(o *Options) {
		o.ModelsPackage = name
	}
}

// WithClientPackage sets the client package name
func WithClientPackage(name string) Option {
	return func(o *Options) {
		o.ClientPackage = name
	}
}

// WithNamingConvention sets the naming convention for generated types
func WithNamingConvention(convention NamingConvention) Option {
	return func(o *Options) {
		o.NamingConvention = convention
	}
}

// WithOmitEmpty enables or disables omitempty in JSON tags
func WithOmitEmpty(enabled bool) Option {
	return func(o *Options) {
		o.OmitEmpty = enabled
	}
}

// WithValidationTags enables or disables validation tags
func WithValidationTags(enabled bool) Option {
	return func(o *Options) {
		o.ValidationTags = enabled
	}
}

// WithBaseURL sets the default base URL for the client
func WithBaseURL(url string) Option {
	return func(o *Options) {
		o.BaseURL = url
	}
}

// WithMethodPrefix sets the prefix for generated client methods
func WithMethodPrefix(prefix string) Option {
	return func(o *Options) {
		o.MethodPrefix = prefix
	}
}

// WithGenerateHelpers enables or disables helper function generation
func WithGenerateHelpers(enabled bool) Option {
	return func(o *Options) {
		o.GenerateHelpers = enabled
	}
}

// WithNameMappings sets name remapping patterns
func WithNameMappings(mappings []NameMapping) Option {
	return func(o *Options) {
		o.NameMappings = mappings
	}
}

// WithModelsOnly sets whether to generate only models (skip client code)
func WithModelsOnly(enabled bool) Option {
	return func(o *Options) {
		o.ModelsOnly = enabled
	}
}

// ParseNameMappings parses comma-separated mapping patterns
// Format: "pattern:replacement" or "pattern" (removes matched suffix)
// Supports glob patterns: "*Response", "*Request", etc.
func ParseNameMappings(s string) ([]NameMapping, error) {
	if s == "" {
		return nil, nil
	}

	var mappings []NameMapping
	parts := strings.Split(s, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		var pattern, replacement string
		if idx := strings.Index(part, ":"); idx >= 0 {
			pattern = strings.TrimSpace(part[:idx])
			replacement = strings.TrimSpace(part[idx+1:])
		} else {
			pattern = part
			replacement = ""
		}

		isGlob := strings.Contains(pattern, "*") || strings.Contains(pattern, "?")
		mappings = append(mappings, NameMapping{
			Pattern:     pattern,
			Replacement: replacement,
			IsGlob:      isGlob,
		})
	}

	return mappings, nil
}

// Apply applies functional options to the Options struct
func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// NewOptions creates a new Options with the given functional options
func NewOptions(opts ...Option) Options {
	o := DefaultOptions()
	o.Apply(opts...)
	return o
}
