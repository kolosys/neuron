package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kolosys/neuron/internal/generator"
)

const (
	version = "0.1.0"
	usage   = `neuron-cli - OpenAPI v3 Code Generator for Neuron

Usage:
  neuron-cli generate [options]
  neuron-cli version
  neuron-cli help

Commands:
  generate    Generate Go client code from an OpenAPI specification
  version     Print the version
  help        Show this help message

Generate Options:
  -i, --input        Path to OpenAPI v3 JSON specification (required if --url not set)
  -u, --url          URL to fetch OpenAPI v3 JSON specification from
  -o, --output       Output directory for generated code (default: ./generated)
  -p, --package      Base package name (default: generated)
  --models-pkg       Package name for models (default: models)
  --client-pkg       Package name for client (default: client)
  --naming           Naming convention: PascalCase, camelCase, snake_case (default: PascalCase)
  --omit-empty       Add omitempty to optional JSON fields (default: true)
  --validation       Add validation tags to struct fields (default: false)
  --base-url         Default base URL for generated client
  --method-prefix    Prefix for generated client methods
  --no-helpers       Disable helper function generation

Examples:
  neuron-cli generate -i openapi.json -o ./api
  neuron-cli generate --input spec.json --output ./client --package myapi
  neuron-cli generate -i api.json -o ./gen --base-url https://api.example.com
  neuron-cli generate --url https://example.com/openapi.json -o ./api
`
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		runGenerate(os.Args[2:])
	case "version":
		fmt.Printf("neuron-cli version %s\n", version)
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		fmt.Print(usage)
		os.Exit(1)
	}
}

func runGenerate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)

	var (
		input         string
		inputURL      string
		output        string
		packageName   string
		modelsPackage string
		clientPackage string
		naming        string
		omitEmpty     bool
		validation    bool
		baseURL       string
		methodPrefix  string
		noHelpers     bool
	)

	fs.StringVar(&input, "i", "", "Path to OpenAPI v3 JSON specification")
	fs.StringVar(&input, "input", "", "Path to OpenAPI v3 JSON specification")
	fs.StringVar(&inputURL, "u", "", "URL to fetch OpenAPI v3 JSON specification from")
	fs.StringVar(&inputURL, "url", "", "URL to fetch OpenAPI v3 JSON specification from")
	fs.StringVar(&output, "o", "./generated", "Output directory for generated code")
	fs.StringVar(&output, "output", "./generated", "Output directory for generated code")
	fs.StringVar(&packageName, "p", "generated", "Base package name")
	fs.StringVar(&packageName, "package", "generated", "Base package name")
	fs.StringVar(&modelsPackage, "models-pkg", "models", "Package name for models")
	fs.StringVar(&clientPackage, "client-pkg", "client", "Package name for client")
	fs.StringVar(&naming, "naming", "PascalCase", "Naming convention: PascalCase, camelCase, snake_case")
	fs.BoolVar(&omitEmpty, "omit-empty", true, "Add omitempty to optional JSON fields")
	fs.BoolVar(&validation, "validation", false, "Add validation tags to struct fields")
	fs.StringVar(&baseURL, "base-url", "", "Default base URL for generated client")
	fs.StringVar(&methodPrefix, "method-prefix", "", "Prefix for generated client methods")
	fs.BoolVar(&noHelpers, "no-helpers", false, "Disable helper function generation")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if input == "" && inputURL == "" {
		fmt.Fprintln(os.Stderr, "Error: input is required (-i/--input for file or -u/--url for URL)")
		fs.PrintDefaults()
		os.Exit(1)
	}

	if input != "" && inputURL != "" {
		fmt.Fprintln(os.Stderr, "Error: cannot specify both --input and --url")
		os.Exit(1)
	}

	if input != "" {
		if _, err := os.Stat(input); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: input file does not exist: %s\n", input)
			os.Exit(1)
		}
	}

	namingConvention := parseNamingConvention(naming)

	opts := generator.NewOptions(
		generator.WithInputPath(input),
		generator.WithInputURL(inputURL),
		generator.WithOutputPath(output),
		generator.WithPackageName(packageName),
		generator.WithModelsPackage(modelsPackage),
		generator.WithClientPackage(clientPackage),
		generator.WithNamingConvention(namingConvention),
		generator.WithOmitEmpty(omitEmpty),
		generator.WithValidationTags(validation),
		generator.WithBaseURL(baseURL),
		generator.WithMethodPrefix(methodPrefix),
		generator.WithGenerateHelpers(!noHelpers),
	)

	gen := generator.New(opts)

	source := input
	if inputURL != "" {
		source = inputURL
	}
	fmt.Printf("Generating code from %s to %s...\n", source, output)

	if err := gen.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Code generation complete!")
}

func parseNamingConvention(s string) generator.NamingConvention {
	switch s {
	case "camelCase":
		return generator.NamingCamelCase
	case "snake_case":
		return generator.NamingSnakeCase
	default:
		return generator.NamingPascalCase
	}
}
