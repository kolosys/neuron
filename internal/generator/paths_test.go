package generator_test

import (
	"testing"

	"github.com/kolosys/neuron/internal/generator"
)

func TestPathProcessor_Process(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	schemaProcessor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)
	_, err = schemaProcessor.Process()
	if err != nil {
		t.Fatalf("failed to process schemas: %v", err)
	}

	pathProcessor := generator.NewPathProcessor(opts, schemaProcessor)
	err = pathProcessor.Process(spec)
	if err != nil {
		t.Fatalf("failed to process paths: %v", err)
	}

	routes := pathProcessor.GetRoutes()
	if len(routes) == 0 {
		t.Error("expected at least one route")
	}

	routeMap := make(map[string]generator.RouteDefinition)
	for _, r := range routes {
		routeMap[r.Name] = r
	}

	listUsers := routeMap["ListUsers"]
	if listUsers.Name == "" {
		t.Error("expected ListUsers route")
	} else {
		if listUsers.Method != generator.MethodGET {
			t.Errorf("ListUsers method = %v, want GET", listUsers.Method)
		}
		if listUsers.Path != "/users" {
			t.Errorf("ListUsers path = %v, want /users", listUsers.Path)
		}
		if len(listUsers.QueryParams) != 2 {
			t.Errorf("ListUsers query params = %d, want 2", len(listUsers.QueryParams))
		}
	}

	getUser := routeMap["GetUser"]
	if getUser.Name == "" {
		t.Error("expected GetUser route")
	} else {
		if len(getUser.PathParams) != 1 {
			t.Errorf("GetUser path params = %d, want 1", len(getUser.PathParams))
		}
		if getUser.PathParams[0].Name != "id" {
			t.Errorf("GetUser path param name = %v, want id", getUser.PathParams[0].Name)
		}
	}

	createUser := routeMap["CreateUser"]
	if createUser.Name == "" {
		t.Error("expected CreateUser route")
	} else {
		if !createUser.HasRequestBody {
			t.Error("CreateUser should have request body")
		}
		if createUser.RequestBodyType != "UserCreateRequest" {
			t.Errorf("CreateUser request body type = %v, want UserCreateRequest", createUser.RequestBodyType)
		}
	}

	deleteUser := routeMap["DeleteUser"]
	if deleteUser.Name == "" {
		t.Error("expected DeleteUser route")
	} else {
		if !deleteUser.Deprecated {
			t.Error("DeleteUser should be deprecated")
		}
	}
}

func TestPathProcessor_Security(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	schemaProcessor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)
	pathProcessor := generator.NewPathProcessor(opts, schemaProcessor)
	err = pathProcessor.Process(spec)
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	security := pathProcessor.GetSecurity()
	if len(security) != 2 {
		t.Errorf("security schemes = %d, want 2", len(security))
	}

	secMap := make(map[string]generator.SecurityConfig)
	for _, s := range security {
		secMap[s.Name] = s
	}

	bearer := secMap["BearerAuth"]
	if bearer.Name == "" {
		t.Error("expected BearerAuth security scheme")
	} else {
		if bearer.Type != "http" {
			t.Errorf("BearerAuth type = %v, want http", bearer.Type)
		}
		if bearer.Scheme != "bearer" {
			t.Errorf("BearerAuth scheme = %v, want bearer", bearer.Scheme)
		}
	}

	apiKey := secMap["ApiKey"]
	if apiKey.Name == "" {
		t.Error("expected ApiKey security scheme")
	} else {
		if apiKey.Type != "apiKey" {
			t.Errorf("ApiKey type = %v, want apiKey", apiKey.Type)
		}
		if apiKey.In != "header" {
			t.Errorf("ApiKey in = %v, want header", apiKey.In)
		}
	}
}

func TestPathProcessor_Headers(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	schemaProcessor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)
	pathProcessor := generator.NewPathProcessor(opts, schemaProcessor)
	err = pathProcessor.Process(spec)
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	headers := pathProcessor.GetHeaders()
	if len(headers) != 2 {
		t.Errorf("headers = %d, want 2", len(headers))
	}

	headerMap := make(map[string]generator.HeaderDefinition)
	for _, h := range headers {
		headerMap[h.Name] = h
	}

	reqID := headerMap["X-Request-ID"]
	if reqID.Name == "" {
		t.Error("expected X-Request-ID header")
	} else {
		if reqID.Type.Name != "string" {
			t.Errorf("X-Request-ID type = %v, want string", reqID.Type.Name)
		}
	}

	rateLimit := headerMap["X-RateLimit-Remaining"]
	if rateLimit.Name == "" {
		t.Error("expected X-RateLimit-Remaining header")
	} else {
		if rateLimit.Type.Name != "int32" {
			t.Errorf("X-RateLimit-Remaining type = %v, want int32", rateLimit.Type.Name)
		}
	}
}

func TestPathProcessor_Responses(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	schemaProcessor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)
	pathProcessor := generator.NewPathProcessor(opts, schemaProcessor)
	err = pathProcessor.Process(spec)
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	responses := pathProcessor.GetResponses()
	if len(responses) != 2 {
		t.Errorf("responses = %d, want 2", len(responses))
	}

	respMap := make(map[string]generator.ResponseDefinition)
	for _, r := range responses {
		respMap[r.Name] = r
	}

	notFound := respMap["NotFound"]
	if notFound.Name == "" {
		t.Error("expected NotFound response")
	}

	unauthorized := respMap["Unauthorized"]
	if unauthorized.Name == "" {
		t.Error("expected Unauthorized response")
	}
}

func TestPathProcessor_MethodPrefix(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	opts.MethodPrefix = "API"
	schemaProcessor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)
	pathProcessor := generator.NewPathProcessor(opts, schemaProcessor)
	err = pathProcessor.Process(spec)
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	routes := pathProcessor.GetRoutes()
	for _, r := range routes {
		if len(r.Name) < 3 || r.Name[:3] != "API" {
			t.Errorf("route %q should have API prefix", r.Name)
		}
	}
}

func TestRouteDefinition_ResponseTypes(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	schemaProcessor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)
	pathProcessor := generator.NewPathProcessor(opts, schemaProcessor)
	err = pathProcessor.Process(spec)
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	routes := pathProcessor.GetRoutes()
	routeMap := make(map[string]generator.RouteDefinition)
	for _, r := range routes {
		routeMap[r.Name] = r
	}

	listUsers := routeMap["ListUsers"]
	if listUsers.ResponseType != "[]User" {
		t.Errorf("ListUsers response type = %v, want []User", listUsers.ResponseType)
	}

	getUser := routeMap["GetUser"]
	if getUser.ResponseType != "User" {
		t.Errorf("GetUser response type = %v, want User", getUser.ResponseType)
	}

	deleteUser := routeMap["DeleteUser"]
	if deleteUser.ResponseType != "EmptyResponse" {
		t.Errorf("DeleteUser response type = %v, want EmptyResponse", deleteUser.ResponseType)
	}
}
