package generator

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Generator handles OpenAPI to Go code generation
type Generator struct {
	opts      Options
	parser    *Parser
	templates *Templates
}

// New creates a new Generator with the given options
func New(opts Options) *Generator {
	return &Generator{
		opts:   opts,
		parser: NewParser(),
	}
}

// Generate parses the OpenAPI spec and generates Go code
func (g *Generator) Generate() error {
	return g.GenerateWithContext(context.Background())
}

// GenerateWithContext parses the OpenAPI spec and generates Go code with context support
func (g *Generator) GenerateWithContext(ctx context.Context) error {
	var spec *OpenAPI
	var err error

	if g.opts.InputURL != "" {
		spec, err = g.parser.ParseURL(ctx, g.opts.InputURL)
		if err != nil {
			return fmt.Errorf("failed to parse OpenAPI spec from URL: %w", err)
		}
	} else {
		spec, err = g.parser.ParseFile(g.opts.InputPath)
		if err != nil {
			return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
		}
	}

	g.templates, err = NewTemplates()
	if err != nil {
		return fmt.Errorf("failed to initialize templates: %w", err)
	}

	if err := g.createOutputDirectories(); err != nil {
		return fmt.Errorf("failed to create output directories: %w", err)
	}

	schemas := make(map[string]*Schema)
	if spec.Components != nil && spec.Components.Schemas != nil {
		schemas = spec.Components.Schemas
	}

	schemaProcessor := NewSchemaProcessor(g.opts, schemas)
	groups, err := schemaProcessor.Process()
	if err != nil {
		return fmt.Errorf("failed to process schemas: %w", err)
	}

	pathProcessor := NewPathProcessor(g.opts, schemaProcessor)
	if err := pathProcessor.Process(spec); err != nil {
		return fmt.Errorf("failed to process paths: %w", err)
	}

	if err := g.generateModels(groups); err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	if err := g.generateClient(pathProcessor); err != nil {
		return fmt.Errorf("failed to generate client: %w", err)
	}

	return nil
}

// createOutputDirectories creates the output directory structure
func (g *Generator) createOutputDirectories() error {
	dirs := []string{
		g.opts.OutputPath,
		filepath.Join(g.opts.OutputPath, g.opts.ModelsPackage),
		filepath.Join(g.opts.OutputPath, g.opts.ClientPackage),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateModels generates model files for each group
func (g *Generator) generateModels(groups map[string]*SchemaGroup) error {
	for groupName, group := range groups {
		if len(group.Structs) == 0 {
			continue
		}

		imports := g.collectModelImports(group)

		data := TemplateData{
			PackageName: g.opts.ModelsPackage,
			Groups:      map[string]*SchemaGroup{groupName: group},
			Imports:     imports,
			GeneratedAt: time.Now().Format(time.RFC3339),
		}

		var buf bytes.Buffer
		if err := g.templates.Models().Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to execute models template for %s: %w", groupName, err)
		}

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			formatted = buf.Bytes()
		}

		filename := filepath.Join(g.opts.OutputPath, g.opts.ModelsPackage, groupName+".go")
		if err := os.WriteFile(filename, formatted, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}

		fmt.Printf("  Generated: %s\n", filename)
	}

	return nil
}

// generateClient generates client files
func (g *Generator) generateClient(pathProcessor *PathProcessor) error {
	routes := pathProcessor.GetRoutes()
	security := pathProcessor.GetSecurity()
	headers := pathProcessor.GetHeaders()
	responses := pathProcessor.GetResponses()

	generatedAt := time.Now().Format(time.RFC3339)

	if err := g.generateClientFile(generatedAt); err != nil {
		return err
	}

	if err := g.generatePathsFile(routes, generatedAt); err != nil {
		return err
	}

	if err := g.generateSecurityFile(security, generatedAt); err != nil {
		return err
	}

	if err := g.generateResponsesFile(responses, generatedAt); err != nil {
		return err
	}

	if err := g.generateHeadersFile(headers, generatedAt); err != nil {
		return err
	}

	return nil
}

// generateClientFile generates rest_client.go
func (g *Generator) generateClientFile(generatedAt string) error {
	data := TemplateData{
		PackageName:   g.opts.PackageName,
		ClientPackage: g.opts.ClientPackage,
		ModelsPackage: g.opts.ModelsPackage,
		BaseURL:       g.opts.BaseURL,
		GeneratedAt:   generatedAt,
	}

	var buf bytes.Buffer
	if err := g.templates.Client().Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute client template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		formatted = buf.Bytes()
	}

	filename := filepath.Join(g.opts.OutputPath, g.opts.ClientPackage, "rest_client.go")
	if err := os.WriteFile(filename, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	fmt.Printf("  Generated: %s\n", filename)
	return nil
}

// generatePathsFile generates rest_paths.go
func (g *Generator) generatePathsFile(routes []RouteDefinition, generatedAt string) error {
	data := TemplateData{
		PackageName:         g.opts.PackageName,
		ClientPackage:       g.opts.ClientPackage,
		ModelsPackage:       g.opts.ModelsPackage,
		ModelsPackageImport: g.formatModelsPackageImport(),
		Routes:              routes,
		GeneratedAt:         generatedAt,
	}

	var buf bytes.Buffer
	if err := g.templates.Paths().Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute paths template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		formatted = buf.Bytes()
	}

	filename := filepath.Join(g.opts.OutputPath, g.opts.ClientPackage, "rest_paths.go")
	if err := os.WriteFile(filename, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	fmt.Printf("  Generated: %s\n", filename)
	return nil
}

// generateSecurityFile generates security.go
func (g *Generator) generateSecurityFile(security []SecurityConfig, generatedAt string) error {
	data := TemplateData{
		PackageName:   g.opts.PackageName,
		ClientPackage: g.opts.ClientPackage,
		Security:      security,
		GeneratedAt:   generatedAt,
	}

	var buf bytes.Buffer
	if err := g.templates.Security().Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute security template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		formatted = buf.Bytes()
	}

	filename := filepath.Join(g.opts.OutputPath, g.opts.ClientPackage, "security.go")
	if err := os.WriteFile(filename, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	fmt.Printf("  Generated: %s\n", filename)
	return nil
}

// generateResponsesFile generates responses.go
func (g *Generator) generateResponsesFile(responses []ResponseDefinition, generatedAt string) error {
	data := TemplateData{
		PackageName:         g.opts.PackageName,
		ClientPackage:       g.opts.ClientPackage,
		ModelsPackage:       g.opts.ModelsPackage,
		ModelsPackageImport: g.formatModelsPackageImport(),
		Responses:           responses,
		GeneratedAt:         generatedAt,
	}

	var buf bytes.Buffer
	if err := g.templates.Responses().Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute responses template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		formatted = buf.Bytes()
	}

	filename := filepath.Join(g.opts.OutputPath, g.opts.ClientPackage, "responses.go")
	if err := os.WriteFile(filename, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	fmt.Printf("  Generated: %s\n", filename)
	return nil
}

// generateHeadersFile generates headers.go
func (g *Generator) generateHeadersFile(headers []HeaderDefinition, generatedAt string) error {
	data := TemplateData{
		PackageName:   g.opts.PackageName,
		ClientPackage: g.opts.ClientPackage,
		Headers:       headers,
		GeneratedAt:   generatedAt,
	}

	var buf bytes.Buffer
	if err := g.templates.Headers().Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute headers template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		formatted = buf.Bytes()
	}

	filename := filepath.Join(g.opts.OutputPath, g.opts.ClientPackage, "headers.go")
	if err := os.WriteFile(filename, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	fmt.Printf("  Generated: %s\n", filename)
	return nil
}

// collectModelImports collects required imports for a schema group
func (g *Generator) collectModelImports(group *SchemaGroup) []string {
	importsMap := make(map[string]bool)

	for _, s := range group.Structs {
		for _, f := range s.Fields {
			if f.Type != nil && f.Type.Package != "" {
				importsMap[f.Type.Package] = true
			}
			if f.Type != nil && f.Type.ElementType != nil && f.Type.ElementType.Package != "" {
				importsMap[f.Type.ElementType.Package] = true
			}
		}
	}

	var imports []string
	for imp := range importsMap {
		imports = append(imports, imp)
	}
	sort.Strings(imports)

	return imports
}

// formatModelsPackageImport returns the full import path for the models package
func (g *Generator) formatModelsPackageImport() string {
	if strings.Contains(g.opts.PackageName, "/") {
		return g.opts.PackageName + "/" + g.opts.ModelsPackage
	}
	return g.opts.ModelsPackage
}
