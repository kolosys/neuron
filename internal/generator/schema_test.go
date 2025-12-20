package generator_test

import (
	"testing"

	"github.com/kolosys/neuron/internal/generator"
)

func TestSchemaProcessor_Process(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	processor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)

	groups, err := processor.Process()
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	if len(groups) == 0 {
		t.Error("expected at least one group")
	}

	userGroup := groups["user"]
	if userGroup == nil {
		t.Fatal("expected 'user' group")
	}

	foundUser := false
	for _, s := range userGroup.Structs {
		if s.Name == "User" {
			foundUser = true
			if s.Description != "A user in the system" {
				t.Errorf("User description = %q, want %q", s.Description, "A user in the system")
			}
		}
	}

	if !foundUser {
		t.Error("expected to find 'User' struct in user group")
	}

	orderGroup := groups["order"]
	if orderGroup == nil {
		t.Fatal("expected 'order' group")
	}
}

func TestSchemaProcessor_EnumProcessing(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	processor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)

	_, err = processor.Process()
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	processed := processor.GetProcessedSchemas()

	userRole := processed["UserRole"]
	if userRole == nil {
		t.Fatal("expected UserRole schema")
	}

	if !userRole.IsEnum {
		t.Error("UserRole should be an enum")
	}

	if len(userRole.EnumValues) != 3 {
		t.Errorf("UserRole enum values count = %d, want 3", len(userRole.EnumValues))
	}
}

func TestGoType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		typ  generator.GoType
		want string
	}{
		{
			name: "simple string",
			typ:  generator.GoType{Name: "string"},
			want: "string",
		},
		{
			name: "pointer string",
			typ:  generator.GoType{Name: "string", IsPointer: true},
			want: "*string",
		},
		{
			name: "slice of strings",
			typ:  generator.GoType{Name: "string", IsSlice: true, ElementType: &generator.GoType{Name: "string"}},
			want: "[]string",
		},
		{
			name: "map string to any",
			typ:  generator.GoType{Name: "any", IsMap: true, KeyType: &generator.GoType{Name: "string"}, ElementType: &generator.GoType{Name: "any"}},
			want: "map[string]any",
		},
		{
			name: "time.Time with package",
			typ:  generator.GoType{Name: "Time", Package: "time"},
			want: "time.Time",
		},
		{
			name: "pointer to time.Time",
			typ:  generator.GoType{Name: "Time", Package: "time", IsPointer: true},
			want: "*time.Time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.typ.String()
			if got != tt.want {
				t.Errorf("GoType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSchemaProcessor_FieldSorting(t *testing.T) {
	t.Parallel()

	p := generator.NewParser()
	spec, err := p.ParseFile("testdata/sample_spec.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	opts := generator.DefaultOptions()
	processor := generator.NewSchemaProcessor(opts, spec.Components.Schemas)

	_, err = processor.Process()
	if err != nil {
		t.Fatalf("failed to process: %v", err)
	}

	processed := processor.GetProcessedSchemas()

	user := processed["User"]
	if user == nil {
		t.Fatal("expected User schema")
	}

	if len(user.Fields) == 0 {
		t.Fatal("expected User to have fields")
	}

	if user.Fields[0].Name != "ID" {
		t.Errorf("first field should be ID, got %q", user.Fields[0].Name)
	}
}

func TestSchemaProcessor_NamingConventions(t *testing.T) {
	t.Parallel()

	schemas := map[string]*generator.Schema{
		"user_profile": {
			Type: generator.NewSchemaType("object"),
			Properties: map[string]*generator.Schema{
				"first_name": {Type: generator.NewSchemaType("string")},
				"last_name":  {Type: generator.NewSchemaType("string")},
			},
		},
	}

	tests := []struct {
		name       string
		convention generator.NamingConvention
		wantStruct string
	}{
		{
			name:       "PascalCase",
			convention: generator.NamingPascalCase,
			wantStruct: "UserProfile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := generator.DefaultOptions()
			opts.NamingConvention = tt.convention

			processor := generator.NewSchemaProcessor(opts, schemas)
			_, err := processor.Process()
			if err != nil {
				t.Fatalf("failed to process: %v", err)
			}

			processed := processor.GetProcessedSchemas()
			found := false
			for name := range processed {
				if processed[name].Name == tt.wantStruct {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected struct named %q", tt.wantStruct)
			}
		})
	}
}

func TestSchemaProcessor_OmitEmpty(t *testing.T) {
	t.Parallel()

	schemas := map[string]*generator.Schema{
		"TestObject": {
			Type:     generator.NewSchemaType("object"),
			Required: []string{"id"},
			Properties: map[string]*generator.Schema{
				"id":    {Type: generator.NewSchemaType("integer"), Format: "int64"},
				"name":  {Type: generator.NewSchemaType("string")},
				"email": {Type: generator.NewSchemaType("string")},
			},
		},
	}

	t.Run("with omit empty", func(t *testing.T) {
		t.Parallel()

		opts := generator.DefaultOptions()
		opts.OmitEmpty = true

		processor := generator.NewSchemaProcessor(opts, schemas)
		_, err := processor.Process()
		if err != nil {
			t.Fatalf("failed to process: %v", err)
		}

		processed := processor.GetProcessedSchemas()
		testObj := processed["TestObject"]

		for _, field := range testObj.Fields {
			if field.Name == "ID" {
				if field.OmitEmpty {
					t.Error("ID field should not have omitempty (it's required)")
				}
			} else {
				if !field.OmitEmpty {
					t.Errorf("%s field should have omitempty (it's optional)", field.Name)
				}
			}
		}
	})

	t.Run("without omit empty", func(t *testing.T) {
		t.Parallel()

		opts := generator.DefaultOptions()
		opts.OmitEmpty = false

		processor := generator.NewSchemaProcessor(opts, schemas)
		_, err := processor.Process()
		if err != nil {
			t.Fatalf("failed to process: %v", err)
		}

		processed := processor.GetProcessedSchemas()
		testObj := processed["TestObject"]

		for _, field := range testObj.Fields {
			if field.OmitEmpty {
				t.Errorf("%s field should not have omitempty", field.Name)
			}
		}
	})
}
