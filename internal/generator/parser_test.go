package generator_test

import (
	"testing"

	"github.com/kolosys/neuron/internal/generator"
)

func TestParser_ParseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		wantErr     bool
		wantVersion string
		wantSchemas int
		wantPaths   int
	}{
		{
			name:        "valid sample spec",
			path:        "testdata/sample_spec.json",
			wantErr:     false,
			wantVersion: "3.1.0",
			wantSchemas: 10,
			wantPaths:   3,
		},
		{
			name:    "non-existent file",
			path:    "testdata/nonexistent.json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := generator.NewParser()
			spec, err := p.ParseFile(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if spec.OpenAPI != tt.wantVersion {
				t.Errorf("OpenAPI version = %v, want %v", spec.OpenAPI, tt.wantVersion)
			}

			if spec.Components == nil || len(spec.Components.Schemas) != tt.wantSchemas {
				got := 0
				if spec.Components != nil {
					got = len(spec.Components.Schemas)
				}
				t.Errorf("schemas count = %d, want %d", got, tt.wantSchemas)
			}

			if len(spec.Paths) != tt.wantPaths {
				t.Errorf("paths count = %d, want %d", len(spec.Paths), tt.wantPaths)
			}
		})
	}
}

func TestParser_Parse_InvalidJSON(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	_, err := p.Parse([]byte("invalid json"))

	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestParser_Parse_UnsupportedVersion(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	_, err := p.Parse([]byte(`{"openapi": "2.0.0"}`))

	if err == nil {
		t.Error("expected error for unsupported version, got nil")
	}
}

func TestGetRefName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ref  string
		want string
	}{
		{"#/components/schemas/User", "User"},
		{"#/components/schemas/OrderItem", "OrderItem"},
		{"#/components/responses/NotFound", "NotFound"},
		{"#/components/parameters/PageLimit", "PageLimit"},
		{"#/components/headers/X-Request-ID", "X-Request-ID"},
		{"#/components/requestBodies/UserInput", "UserInput"},
		{"SomethingElse", "SomethingElse"},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			t.Parallel()

			got := generator.GetRefName(tt.ref)
			if got != tt.want {
				t.Errorf("GetRefName(%q) = %q, want %q", tt.ref, got, tt.want)
			}
		})
	}
}

func TestParser_SpecInfo(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if spec.Info.Title != "Sample API" {
		t.Errorf("Title = %q, want %q", spec.Info.Title, "Sample API")
	}

	if spec.Info.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", spec.Info.Version, "1.0.0")
	}

	if len(spec.Servers) != 1 {
		t.Errorf("Servers count = %d, want 1", len(spec.Servers))
	}

	if spec.Servers[0].URL != "https://api.example.com/v1" {
		t.Errorf("Server URL = %q, want %q", spec.Servers[0].URL, "https://api.example.com/v1")
	}
}

func TestParser_SecuritySchemes(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if spec.Components == nil {
		t.Fatal("Components is nil")
	}

	if len(spec.Components.SecuritySchemes) != 2 {
		t.Errorf("SecuritySchemes count = %d, want 2", len(spec.Components.SecuritySchemes))
	}

	bearer := spec.Components.SecuritySchemes["bearerAuth"]
	if bearer == nil {
		t.Fatal("bearerAuth scheme not found")
	}

	if bearer.Type != "http" {
		t.Errorf("bearerAuth.Type = %q, want %q", bearer.Type, "http")
	}

	if bearer.Scheme != "bearer" {
		t.Errorf("bearerAuth.Scheme = %q, want %q", bearer.Scheme, "bearer")
	}
}
