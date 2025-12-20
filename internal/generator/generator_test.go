package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kolosys/neuron/internal/generator"
)

func TestGenerator_Generate(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	opts := generator.NewOptions(
		generator.WithInputPath("testdata/sample_spec.json"),
		generator.WithOutputPath(tmpDir),
		generator.WithPackageName("testapi"),
		generator.WithModelsPackage("models"),
		generator.WithClientPackage("client"),
		generator.WithOmitEmpty(true),
	)

	gen := generator.New(opts)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("failed to generate: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tmpDir, "models", "user.go"),
		filepath.Join(tmpDir, "models", "order.go"),
		filepath.Join(tmpDir, "client", "rest_client.go"),
		filepath.Join(tmpDir, "client", "rest_paths.go"),
		filepath.Join(tmpDir, "client", "security.go"),
		filepath.Join(tmpDir, "client", "responses.go"),
		filepath.Join(tmpDir, "client", "headers.go"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestGenerator_Generate_DirectoryStructure(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	opts := generator.NewOptions(
		generator.WithInputPath("testdata/sample_spec.json"),
		generator.WithOutputPath(tmpDir),
	)

	gen := generator.New(opts)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("failed to generate: %v", err)
	}

	dirs := []string{
		filepath.Join(tmpDir, "models"),
		filepath.Join(tmpDir, "client"),
	}

	for _, d := range dirs {
		info, err := os.Stat(d)
		if os.IsNotExist(err) {
			t.Errorf("expected directory %s to exist", d)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s should be a directory", d)
		}
	}
}

func TestGenerator_Generate_InvalidInput(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	opts := generator.NewOptions(
		generator.WithInputPath("testdata/nonexistent.json"),
		generator.WithOutputPath(tmpDir),
	)

	gen := generator.New(opts)
	err := gen.Generate()
	if err == nil {
		t.Error("expected error for non-existent input file")
	}
}

func TestGenerator_Generate_WithBaseURL(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	opts := generator.NewOptions(
		generator.WithInputPath("testdata/sample_spec.json"),
		generator.WithOutputPath(tmpDir),
		generator.WithBaseURL("https://api.example.com"),
	)

	gen := generator.New(opts)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("failed to generate: %v", err)
	}

	clientPath := filepath.Join(tmpDir, "client", "rest_client.go")
	content, err := os.ReadFile(clientPath)
	if err != nil {
		t.Fatalf("failed to read client file: %v", err)
	}

	if !contains(string(content), "https://api.example.com") {
		t.Error("generated client should contain base URL")
	}
}

func TestGenerator_Generate_ModelContent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	opts := generator.NewOptions(
		generator.WithInputPath("testdata/sample_spec.json"),
		generator.WithOutputPath(tmpDir),
		generator.WithModelsPackage("models"),
	)

	gen := generator.New(opts)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("failed to generate: %v", err)
	}

	userPath := filepath.Join(tmpDir, "models", "user.go")
	content, err := os.ReadFile(userPath)
	if err != nil {
		t.Fatalf("failed to read user models file: %v", err)
	}

	contentStr := string(content)

	if !contains(contentStr, "package models") {
		t.Error("models file should have correct package")
	}

	if !contains(contentStr, "type User struct") {
		t.Error("models file should contain User struct")
	}

	if !contains(contentStr, `json:"id"`) {
		t.Error("models file should contain JSON tags")
	}
}

func TestGenerator_Generate_SecurityContent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	opts := generator.NewOptions(
		generator.WithInputPath("testdata/sample_spec.json"),
		generator.WithOutputPath(tmpDir),
	)

	gen := generator.New(opts)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("failed to generate: %v", err)
	}

	secPath := filepath.Join(tmpDir, "client", "security.go")
	content, err := os.ReadFile(secPath)
	if err != nil {
		t.Fatalf("failed to read security file: %v", err)
	}

	contentStr := string(content)

	if !contains(contentStr, "ApplyBearerToken") {
		t.Error("security file should contain ApplyBearerToken")
	}

	if !contains(contentStr, "ApplyAPIKey") {
		t.Error("security file should contain ApplyAPIKey")
	}

	if !contains(contentStr, "ApplyBasicAuth") {
		t.Error("security file should contain ApplyBasicAuth")
	}
}

func TestOptions_Defaults(t *testing.T) {
	t.Parallel()

	opts := generator.DefaultOptions()

	if opts.PackageName != "generated" {
		t.Errorf("default PackageName = %v, want generated", opts.PackageName)
	}

	if opts.ModelsPackage != "models" {
		t.Errorf("default ModelsPackage = %v, want models", opts.ModelsPackage)
	}

	if opts.ClientPackage != "client" {
		t.Errorf("default ClientPackage = %v, want client", opts.ClientPackage)
	}

	if opts.NamingConvention != generator.NamingPascalCase {
		t.Errorf("default NamingConvention = %v, want PascalCase", opts.NamingConvention)
	}

	if !opts.OmitEmpty {
		t.Error("default OmitEmpty should be true")
	}

	if opts.ValidationTags {
		t.Error("default ValidationTags should be false")
	}

	if !opts.GenerateHelpers {
		t.Error("default GenerateHelpers should be true")
	}
}

func TestOptions_FunctionalOptions(t *testing.T) {
	t.Parallel()

	opts := generator.NewOptions(
		generator.WithInputPath("/path/to/spec.json"),
		generator.WithOutputPath("/path/to/output"),
		generator.WithPackageName("myapi"),
		generator.WithModelsPackage("types"),
		generator.WithClientPackage("http"),
		generator.WithNamingConvention(generator.NamingCamelCase),
		generator.WithOmitEmpty(false),
		generator.WithValidationTags(true),
		generator.WithBaseURL("https://api.test.com"),
		generator.WithMethodPrefix("Do"),
		generator.WithGenerateHelpers(false),
	)

	if opts.InputPath != "/path/to/spec.json" {
		t.Errorf("InputPath = %v, want /path/to/spec.json", opts.InputPath)
	}

	if opts.OutputPath != "/path/to/output" {
		t.Errorf("OutputPath = %v, want /path/to/output", opts.OutputPath)
	}

	if opts.PackageName != "myapi" {
		t.Errorf("PackageName = %v, want myapi", opts.PackageName)
	}

	if opts.ModelsPackage != "types" {
		t.Errorf("ModelsPackage = %v, want types", opts.ModelsPackage)
	}

	if opts.ClientPackage != "http" {
		t.Errorf("ClientPackage = %v, want http", opts.ClientPackage)
	}

	if opts.NamingConvention != generator.NamingCamelCase {
		t.Errorf("NamingConvention = %v, want camelCase", opts.NamingConvention)
	}

	if opts.OmitEmpty {
		t.Error("OmitEmpty should be false")
	}

	if !opts.ValidationTags {
		t.Error("ValidationTags should be true")
	}

	if opts.BaseURL != "https://api.test.com" {
		t.Errorf("BaseURL = %v, want https://api.test.com", opts.BaseURL)
	}

	if opts.MethodPrefix != "Do" {
		t.Errorf("MethodPrefix = %v, want Do", opts.MethodPrefix)
	}

	if opts.GenerateHelpers {
		t.Error("GenerateHelpers should be false")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
