package generator

import (
	"regexp"
	"sort"
	"strings"
)

// HTTPMethod represents an HTTP method
type HTTPMethod string

const (
	MethodGET     HTTPMethod = "GET"
	MethodPOST    HTTPMethod = "POST"
	MethodPUT     HTTPMethod = "PUT"
	MethodPATCH   HTTPMethod = "PATCH"
	MethodDELETE  HTTPMethod = "DELETE"
	MethodHEAD    HTTPMethod = "HEAD"
	MethodOPTIONS HTTPMethod = "OPTIONS"
)

// RouteDefinition represents a generated route
type RouteDefinition struct {
	Name            string
	Method          HTTPMethod
	Path            string
	OperationID     string
	Summary         string
	Description     string
	Tags            []string
	RequestType     string
	ResponseType    string
	PathParams      []RouteParam
	QueryParams     []RouteParam
	HeaderParams    []RouteParam
	HasRequestBody  bool
	RequestBodyType string
	Security        []SecurityRequirement
	Deprecated      bool
}

// RouteParam represents a route parameter
type RouteParam struct {
	Name        string
	GoName      string
	Type        *GoType
	In          string
	Required    bool
	Description string
	Style       string
}

// SecurityConfig represents a security scheme configuration
type SecurityConfig struct {
	Name         string
	Type         string
	Scheme       string
	BearerFormat string
	In           string
	HeaderName   string
	Description  string
}

// HeaderDefinition represents a reusable header
type HeaderDefinition struct {
	Name        string
	GoName      string
	Type        *GoType
	Required    bool
	Description string
}

// ResponseDefinition represents a reusable response
type ResponseDefinition struct {
	Name        string
	Description string
	Type        string
	Headers     []HeaderDefinition
}

// PathProcessor handles path processing and route generation
type PathProcessor struct {
	opts            Options
	schemaProcessor *SchemaProcessor
	routes          []RouteDefinition
	security        []SecurityConfig
	headers         []HeaderDefinition
	responses       []ResponseDefinition
}

// NewPathProcessor creates a new path processor
func NewPathProcessor(opts Options, schemaProcessor *SchemaProcessor) *PathProcessor {
	return &PathProcessor{
		opts:            opts,
		schemaProcessor: schemaProcessor,
	}
}

// Process processes all paths and generates route definitions
func (p *PathProcessor) Process(spec *OpenAPI) error {
	for path, item := range spec.Paths {
		if err := p.processPathItem(path, item); err != nil {
			return err
		}
	}

	if spec.Components != nil {
		p.processSecuritySchemes(spec.Components.SecuritySchemes)
		p.processHeaders(spec.Components.Headers)
		p.processResponses(spec.Components.Responses)
	}

	sort.Slice(p.routes, func(i, j int) bool {
		if p.routes[i].Path != p.routes[j].Path {
			return p.routes[i].Path < p.routes[j].Path
		}
		return methodOrder(p.routes[i].Method) < methodOrder(p.routes[j].Method)
	})

	return nil
}

// GetRoutes returns all processed routes
func (p *PathProcessor) GetRoutes() []RouteDefinition {
	return p.routes
}

// GetSecurity returns all security configurations
func (p *PathProcessor) GetSecurity() []SecurityConfig {
	return p.security
}

// GetHeaders returns all header definitions
func (p *PathProcessor) GetHeaders() []HeaderDefinition {
	return p.headers
}

// GetResponses returns all response definitions
func (p *PathProcessor) GetResponses() []ResponseDefinition {
	return p.responses
}

// processPathItem processes a single path item
func (p *PathProcessor) processPathItem(path string, item *PathItem) error {
	pathParams := p.extractPathParams(item.Parameters)

	if item.Get != nil {
		p.processOperation(path, MethodGET, item.Get, pathParams)
	}
	if item.Post != nil {
		p.processOperation(path, MethodPOST, item.Post, pathParams)
	}
	if item.Put != nil {
		p.processOperation(path, MethodPUT, item.Put, pathParams)
	}
	if item.Patch != nil {
		p.processOperation(path, MethodPATCH, item.Patch, pathParams)
	}
	if item.Delete != nil {
		p.processOperation(path, MethodDELETE, item.Delete, pathParams)
	}
	if item.Head != nil {
		p.processOperation(path, MethodHEAD, item.Head, pathParams)
	}
	if item.Options != nil {
		p.processOperation(path, MethodOPTIONS, item.Options, pathParams)
	}

	return nil
}

// processOperation processes a single operation
func (p *PathProcessor) processOperation(path string, method HTTPMethod, op *Operation, pathParams []RouteParam) {
	route := RouteDefinition{
		Method:      method,
		Path:        path,
		OperationID: op.OperationID,
		Summary:     op.Summary,
		Description: op.Description,
		Tags:        op.Tags,
		PathParams:  pathParams,
		Security:    op.Security,
		Deprecated:  op.Deprecated,
	}

	route.Name = p.generateRouteName(method, path, op.OperationID)

	for _, param := range op.Parameters {
		rp := p.processParameter(param)
		switch param.In {
		case "path":
			exists := false
			for _, pp := range route.PathParams {
				if pp.Name == rp.Name {
					exists = true
					break
				}
			}
			if !exists {
				route.PathParams = append(route.PathParams, rp)
			}
		case "query":
			route.QueryParams = append(route.QueryParams, rp)
		case "header":
			route.HeaderParams = append(route.HeaderParams, rp)
		}
	}

	if op.RequestBody != nil && op.RequestBody.Content != nil {
		route.HasRequestBody = true
		route.RequestBodyType = p.extractRequestType(op.RequestBody)
		route.RequestType = route.RequestBodyType
	} else {
		route.RequestType = "EmptyRequest"
	}

	route.ResponseType = p.extractResponseType(op.Responses)

	p.routes = append(p.routes, route)
}

// extractPathParams extracts path parameters from parameter list
func (p *PathProcessor) extractPathParams(params []*Parameter) []RouteParam {
	var pathParams []RouteParam
	for _, param := range params {
		if param.In == "path" {
			pathParams = append(pathParams, p.processParameter(param))
		}
	}
	return pathParams
}

// processParameter processes a single parameter
func (p *PathProcessor) processParameter(param *Parameter) RouteParam {
	goType := &GoType{Name: "string"}
	if param.Schema != nil {
		goType = p.schemaProcessor.schemaToGoType(param.Schema)
	}

	return RouteParam{
		Name:        param.Name,
		GoName:      toPascalCase(param.Name),
		Type:        goType,
		In:          param.In,
		Required:    param.Required,
		Description: param.Description,
		Style:       param.Style,
	}
}

// extractRequestType extracts the request body type
func (p *PathProcessor) extractRequestType(reqBody *RequestBody) string {
	if reqBody == nil || reqBody.Content == nil {
		return "EmptyRequest"
	}

	for contentType, media := range reqBody.Content {
		if strings.Contains(contentType, "json") && media.Schema != nil {
			if media.Schema.Ref != "" {
				return toPascalCase(GetRefName(media.Schema.Ref))
			}
			return "any"
		}
	}

	return "any"
}

// extractResponseType extracts the primary response type
func (p *PathProcessor) extractResponseType(responses map[string]*Response) string {
	successCodes := []string{"200", "201", "202", "204"}

	for _, code := range successCodes {
		if resp, ok := responses[code]; ok {
			return p.getResponseSchemaType(resp)
		}
	}

	for code, resp := range responses {
		if strings.HasPrefix(code, "2") {
			return p.getResponseSchemaType(resp)
		}
	}

	return "EmptyResponse"
}

// getResponseSchemaType extracts the type from a response
func (p *PathProcessor) getResponseSchemaType(resp *Response) string {
	if resp == nil {
		return "EmptyResponse"
	}

	if resp.Ref != "" {
		return toPascalCase(GetRefName(resp.Ref))
	}

	if resp.Content == nil {
		return "EmptyResponse"
	}

	for contentType, media := range resp.Content {
		if strings.Contains(contentType, "json") && media.Schema != nil {
			if media.Schema.Ref != "" {
				return toPascalCase(GetRefName(media.Schema.Ref))
			}
			if media.Schema.Type.String() == "array" && media.Schema.Items != nil {
				if media.Schema.Items.Ref != "" {
					return "[]" + toPascalCase(GetRefName(media.Schema.Items.Ref))
				}
			}
			return "any"
		}
	}

	return "EmptyResponse"
}

// generateRouteName generates a name for a route
func (p *PathProcessor) generateRouteName(method HTTPMethod, path, operationID string) string {
	if operationID != "" {
		name := toPascalCase(operationID)
		if p.opts.MethodPrefix != "" {
			return p.opts.MethodPrefix + name
		}
		return name
	}

	cleanPath := cleanPathForName(path)
	name := toPascalCase(string(method)) + toPascalCase(cleanPath)

	if p.opts.MethodPrefix != "" {
		return p.opts.MethodPrefix + name
	}
	return name
}

// processSecuritySchemes processes security schemes
func (p *PathProcessor) processSecuritySchemes(schemes map[string]*SecurityScheme) {
	for name, scheme := range schemes {
		if scheme == nil {
			continue
		}

		config := SecurityConfig{
			Name:        toPascalCase(name),
			Type:        scheme.Type,
			Scheme:      scheme.Scheme,
			Description: scheme.Description,
		}

		if scheme.Type == "http" {
			config.BearerFormat = scheme.BearerFormat
		}
		if scheme.Type == "apiKey" {
			config.In = scheme.In
			config.HeaderName = scheme.Name
		}

		p.security = append(p.security, config)
	}

	sort.Slice(p.security, func(i, j int) bool {
		return p.security[i].Name < p.security[j].Name
	})
}

// processHeaders processes reusable headers
func (p *PathProcessor) processHeaders(headers map[string]*Header) {
	for name, header := range headers {
		if header == nil {
			continue
		}

		goType := &GoType{Name: "string"}
		if header.Schema != nil {
			goType = p.schemaProcessor.schemaToGoType(header.Schema)
		}

		p.headers = append(p.headers, HeaderDefinition{
			Name:        name,
			GoName:      toPascalCase(name),
			Type:        goType,
			Required:    header.Required,
			Description: header.Description,
		})
	}

	sort.Slice(p.headers, func(i, j int) bool {
		return p.headers[i].Name < p.headers[j].Name
	})
}

// processResponses processes reusable responses
func (p *PathProcessor) processResponses(responses map[string]*Response) {
	for name, resp := range responses {
		if resp == nil {
			continue
		}

		def := ResponseDefinition{
			Name:        toPascalCase(name),
			Description: resp.Description,
		}

		if resp.Content != nil {
			for contentType, media := range resp.Content {
				if strings.Contains(contentType, "json") && media.Schema != nil {
					if media.Schema.Ref != "" {
						def.Type = toPascalCase(GetRefName(media.Schema.Ref))
					}
					break
				}
			}
		}

		for hName, header := range resp.Headers {
			if header == nil {
				continue
			}
			goType := &GoType{Name: "string"}
			if header.Schema != nil {
				goType = p.schemaProcessor.schemaToGoType(header.Schema)
			}
			def.Headers = append(def.Headers, HeaderDefinition{
				Name:        hName,
				GoName:      toPascalCase(hName),
				Type:        goType,
				Required:    header.Required,
				Description: header.Description,
			})
		}

		p.responses = append(p.responses, def)
	}

	sort.Slice(p.responses, func(i, j int) bool {
		return p.responses[i].Name < p.responses[j].Name
	})
}

// cleanPathForName cleans a path for use in a name
func cleanPathForName(path string) string {
	path = regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(path, "")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	path = strings.Trim(path, "_")
	path = regexp.MustCompile(`_+`).ReplaceAllString(path, "_")
	return path
}

// methodOrder returns the sort order for HTTP methods
func methodOrder(m HTTPMethod) int {
	switch m {
	case MethodGET:
		return 0
	case MethodPOST:
		return 1
	case MethodPUT:
		return 2
	case MethodPATCH:
		return 3
	case MethodDELETE:
		return 4
	case MethodHEAD:
		return 5
	case MethodOPTIONS:
		return 6
	default:
		return 7
	}
}
