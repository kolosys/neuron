package generator

import (
	"sort"
	"strings"
	"unicode"
)

// GoType represents a normalized Go type
type GoType struct {
	Name        string
	Package     string
	IsPointer   bool
	IsSlice     bool
	IsMap       bool
	KeyType     *GoType
	ElementType *GoType
}

// String returns the Go type as a string
func (t *GoType) String() string {
	var sb strings.Builder

	if t.IsSlice {
		sb.WriteString("[]")
	}
	if t.IsMap {
		sb.WriteString("map[")
		if t.KeyType != nil {
			sb.WriteString(t.KeyType.String())
		} else {
			sb.WriteString("string")
		}
		sb.WriteString("]")
	}
	if t.IsPointer {
		sb.WriteString("*")
	}
	if t.Package != "" {
		sb.WriteString(t.Package)
		sb.WriteString(".")
	}
	if t.ElementType != nil && (t.IsSlice || t.IsMap) {
		sb.WriteString(t.ElementType.String())
	} else {
		sb.WriteString(t.Name)
	}

	return sb.String()
}

// GoField represents a struct field
type GoField struct {
	Name        string
	Type        *GoType
	JSONName    string
	Description string
	Required    bool
	OmitEmpty   bool
	Deprecated  bool
	ReadOnly    bool
	WriteOnly   bool
	Validation  string
}

// GoStruct represents a Go struct type
type GoStruct struct {
	Name        string
	Description string
	Fields      []GoField
	Embeds      []string
	IsEnum      bool
	EnumValues  []EnumValue
	Group       string
}

// EnumValue represents an enum constant
type EnumValue struct {
	Name        string
	Value       any
	Type        string
	Description string
}

// SchemaGroup represents a group of related schemas
type SchemaGroup struct {
	Name    string
	Structs []GoStruct
}

// SchemaProcessor handles schema normalization and grouping
type SchemaProcessor struct {
	opts      Options
	schemas   map[string]*Schema
	processed map[string]*GoStruct
	groups    map[string]*SchemaGroup
}

// NewSchemaProcessor creates a new schema processor
func NewSchemaProcessor(opts Options, schemas map[string]*Schema) *SchemaProcessor {
	return &SchemaProcessor{
		opts:      opts,
		schemas:   schemas,
		processed: make(map[string]*GoStruct),
		groups:    make(map[string]*SchemaGroup),
	}
}

// Process normalizes and groups all schemas
func (p *SchemaProcessor) Process() (map[string]*SchemaGroup, error) {
	for name, schema := range p.schemas {
		if _, err := p.processSchema(name, schema); err != nil {
			return nil, err
		}
	}

	for _, gs := range p.processed {
		group := p.extractDomain(gs.Name)
		gs.Group = group

		if _, ok := p.groups[group]; !ok {
			p.groups[group] = &SchemaGroup{Name: group}
		}
		p.groups[group].Structs = append(p.groups[group].Structs, *gs)
	}

	for _, g := range p.groups {
		sort.Slice(g.Structs, func(i, j int) bool {
			return g.Structs[i].Name < g.Structs[j].Name
		})
	}

	return p.groups, nil
}

// GetProcessedSchemas returns all processed schemas
func (p *SchemaProcessor) GetProcessedSchemas() map[string]*GoStruct {
	return p.processed
}

// processSchema converts an OpenAPI schema to a Go struct
func (p *SchemaProcessor) processSchema(name string, schema *Schema) (*GoStruct, error) {
	if gs, ok := p.processed[name]; ok {
		return gs, nil
	}

	goName := p.toGoName(name)
	gs := &GoStruct{
		Name:        goName,
		Description: schema.Description,
	}

	// Handle traditional enum arrays
	if len(schema.Enum) > 0 {
		gs.IsEnum = true
		goType := p.schemaToGoType(schema)
		for _, v := range schema.Enum {
			enumName := p.enumValueName(goName, v)
			gs.EnumValues = append(gs.EnumValues, EnumValue{
				Name:  enumName,
				Value: v,
				Type:  goType.Name,
			})
		}
		p.processed[name] = gs
		return gs, nil
	}

	// Handle oneOf with const values (OpenAPI 3.1.0 style enums)
	if len(schema.OneOf) > 0 && p.isOneOfEnum(schema.OneOf) {
		gs.IsEnum = true
		// Determine the type from schema type or infer from first const value
		enumType := p.determineEnumType(schema)
		for _, s := range schema.OneOf {
			if s.Const != nil {
				enumName := p.enumValueName(goName, s.Const)
				gs.EnumValues = append(gs.EnumValues, EnumValue{
					Name:        enumName,
					Value:       s.Const,
					Type:        enumType,
					Description: s.Description,
				})
			}
		}
		p.processed[name] = gs
		return gs, nil
	}

	if len(schema.AllOf) > 0 {
		for _, s := range schema.AllOf {
			if s.Ref != "" {
				refName := GetRefName(s.Ref)
				gs.Embeds = append(gs.Embeds, p.toGoName(refName))
			} else if s.Properties != nil {
				fields := p.processProperties(s.Properties, s.Required)
				gs.Fields = append(gs.Fields, fields...)
			}
		}
	}

	if schema.Properties != nil {
		fields := p.processProperties(schema.Properties, schema.Required)
		gs.Fields = append(gs.Fields, fields...)
	}

	p.processed[name] = gs
	return gs, nil
}

// processProperties converts schema properties to Go fields
func (p *SchemaProcessor) processProperties(props map[string]*Schema, required []string) []GoField {
	requiredMap := make(map[string]bool)
	for _, r := range required {
		requiredMap[r] = true
	}

	var fields []GoField
	for name, prop := range props {
		field := GoField{
			Name:        p.toGoName(name),
			Type:        p.schemaToGoType(prop),
			JSONName:    name,
			Description: prop.Description,
			Required:    requiredMap[name],
			OmitEmpty:   p.opts.OmitEmpty && !requiredMap[name],
			Deprecated:  prop.Deprecated,
			ReadOnly:    prop.ReadOnly,
			WriteOnly:   prop.WriteOnly,
		}

		if p.opts.ValidationTags {
			field.Validation = p.buildValidation(prop, requiredMap[name])
		}

		fields = append(fields, field)
	}

	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Name == "ID" {
			return true
		}
		if fields[j].Name == "ID" {
			return false
		}
		return fields[i].Name < fields[j].Name
	})

	return fields
}

// schemaToGoType converts an OpenAPI schema to a Go type
func (p *SchemaProcessor) schemaToGoType(schema *Schema) *GoType {
	if schema == nil {
		return &GoType{Name: "any"}
	}

	if schema.Ref != "" {
		refName := GetRefName(schema.Ref)
		return &GoType{
			Name:      p.toGoName(refName),
			IsPointer: schema.Nullable || schema.Type.IsNullable(),
		}
	}

	// Handle OpenAPI 3.1.0 nullable types (type: ["string", "null"])
	isNullable := schema.Nullable || schema.Type.IsNullable()
	typeName := schema.Type.String()

	switch typeName {
	case "string":
		return p.stringType(schema, isNullable)
	case "integer":
		return p.integerType(schema, isNullable)
	case "number":
		return p.numberType(schema, isNullable)
	case "boolean":
		return &GoType{Name: "bool", IsPointer: isNullable}
	case "array":
		return p.arrayType(schema)
	case "object":
		return p.objectType(schema)
	default:
		if len(schema.AnyOf) > 0 || len(schema.OneOf) > 0 {
			return &GoType{Name: "any"}
		}
		return &GoType{Name: "any"}
	}
}

// stringType returns the Go type for a string schema
func (p *SchemaProcessor) stringType(schema *Schema, isNullable bool) *GoType {
	switch schema.Format {
	case "date-time", "date":
		return &GoType{Name: "Time", Package: "time", IsPointer: isNullable}
	case "uuid":
		return &GoType{Name: "string", IsPointer: isNullable}
	case "uri", "url":
		return &GoType{Name: "string", IsPointer: isNullable}
	case "email":
		return &GoType{Name: "string", IsPointer: isNullable}
	case "byte":
		return &GoType{Name: "byte", IsSlice: true}
	case "binary":
		return &GoType{Name: "byte", IsSlice: true}
	default:
		return &GoType{Name: "string", IsPointer: isNullable}
	}
}

// integerType returns the Go type for an integer schema
func (p *SchemaProcessor) integerType(schema *Schema, isNullable bool) *GoType {
	switch schema.Format {
	case "int32":
		return &GoType{Name: "int32", IsPointer: isNullable}
	case "int64":
		return &GoType{Name: "int64", IsPointer: isNullable}
	default:
		return &GoType{Name: "int64", IsPointer: isNullable}
	}
}

// numberType returns the Go type for a number schema
func (p *SchemaProcessor) numberType(schema *Schema, isNullable bool) *GoType {
	switch schema.Format {
	case "float":
		return &GoType{Name: "float32", IsPointer: isNullable}
	case "double":
		return &GoType{Name: "float64", IsPointer: isNullable}
	default:
		return &GoType{Name: "float64", IsPointer: isNullable}
	}
}

// arrayType returns the Go type for an array schema
func (p *SchemaProcessor) arrayType(schema *Schema) *GoType {
	elemType := p.schemaToGoType(schema.Items)
	return &GoType{
		Name:        elemType.Name,
		Package:     elemType.Package,
		IsSlice:     true,
		ElementType: elemType,
	}
}

// objectType returns the Go type for an object schema
func (p *SchemaProcessor) objectType(schema *Schema) *GoType {
	if schema.AdditionalProperties != nil {
		if innerSchema := schema.AdditionalProperties.GetSchema(); innerSchema != nil {
			valueType := p.schemaToGoType(innerSchema)
			return &GoType{
				IsMap:       true,
				KeyType:     &GoType{Name: "string"},
				ElementType: valueType,
			}
		}
		if schema.AdditionalProperties.IsAllowed() {
			return &GoType{Name: "any", IsMap: true, KeyType: &GoType{Name: "string"}}
		}
	}

	if schema.Properties == nil {
		return &GoType{Name: "any", IsMap: true, KeyType: &GoType{Name: "string"}}
	}

	return &GoType{Name: "any"}
}

// extractDomain extracts the domain/group from a schema name
func (p *SchemaProcessor) extractDomain(name string) string {
	words := splitCamelCase(name)
	if len(words) == 0 {
		return "common"
	}

	domain := strings.ToLower(words[0])

	commonSuffixes := []string{
		"request", "response", "input", "output", "dto",
		"params", "config", "options", "settings", "result",
		"list", "item", "detail", "summary", "info",
		"create", "update", "delete", "patch",
	}

	for _, suffix := range commonSuffixes {
		if strings.EqualFold(words[len(words)-1], suffix) && len(words) > 1 {
			remaining := words[:len(words)-1]
			domain = strings.ToLower(remaining[0])
			break
		}
	}

	if len(domain) <= 2 {
		if len(words) > 1 {
			domain = strings.ToLower(words[0] + words[1])
		} else {
			domain = "common"
		}
	}

	return domain
}

// toGoName converts a name to a valid Go identifier
func (p *SchemaProcessor) toGoName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, " ", "_")

	switch p.opts.NamingConvention {
	case NamingCamelCase:
		return toCamelCase(name)
	case NamingSnakeCase:
		return toSnakeCase(name)
	default:
		return toPascalCase(name)
	}
}

// enumValueName creates a Go name for an enum value
func (p *SchemaProcessor) enumValueName(typeName string, value any) string {
	var strVal string
	switch v := value.(type) {
	case string:
		strVal = v
	case float64:
		// JSON numbers are unmarshaled as float64
		if v == float64(int64(v)) {
			strVal = itoa(int(v))
		} else {
			strVal = ftoa(v)
		}
	case int:
		strVal = itoa(v)
	case int64:
		strVal = itoa(int(v))
	case bool:
		if v {
			strVal = "True"
		} else {
			strVal = "False"
		}
	default:
		strVal = "Unknown"
	}

	strVal = strings.ReplaceAll(strVal, "-", "_")
	strVal = strings.ReplaceAll(strVal, " ", "_")
	strVal = strings.ReplaceAll(strVal, ".", "_")

	return typeName + toPascalCase(strVal)
}

// isOneOfEnum checks if a oneOf array represents an enum (all items have const values)
func (p *SchemaProcessor) isOneOfEnum(oneOf []*Schema) bool {
	if len(oneOf) == 0 {
		return false
	}
	for _, s := range oneOf {
		if s.Const == nil {
			return false
		}
	}
	return true
}

// determineEnumType determines the Go type for an enum based on schema or const values
func (p *SchemaProcessor) determineEnumType(schema *Schema) string {
	// First check if the schema has an explicit type
	typeName := schema.Type.String()
	switch typeName {
	case "string":
		return "string"
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	}

	// If no explicit type, infer from the first const value
	if len(schema.OneOf) > 0 && schema.OneOf[0].Const != nil {
		switch schema.OneOf[0].Const.(type) {
		case string:
			return "string"
		case float64:
			return "int"
		case int, int64:
			return "int"
		case bool:
			return "bool"
		}
	}

	// Default to string
	return "string"
}

// buildValidation creates validation tags for a field
func (p *SchemaProcessor) buildValidation(schema *Schema, required bool) string {
	var parts []string

	if required {
		parts = append(parts, "required")
	}

	if schema.MinLength != nil {
		parts = append(parts, "min="+itoa(*schema.MinLength))
	}
	if schema.MaxLength != nil {
		parts = append(parts, "max="+itoa(*schema.MaxLength))
	}
	if schema.Minimum != nil {
		parts = append(parts, "gte="+ftoa(*schema.Minimum))
	}
	if schema.Maximum != nil {
		parts = append(parts, "lte="+ftoa(*schema.Maximum))
	}
	if schema.Pattern != "" {
		parts = append(parts, "regexp="+schema.Pattern)
	}
	if schema.Format == "email" {
		parts = append(parts, "email")
	}
	if schema.Format == "uri" || schema.Format == "url" {
		parts = append(parts, "url")
	}

	return strings.Join(parts, ",")
}

// splitCamelCase splits a camelCase or PascalCase string into words
func splitCamelCase(s string) []string {
	var words []string
	var current strings.Builder

	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// toPascalCase converts a string to PascalCase
func toPascalCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	words := strings.Fields(s)

	for i, word := range words {
		if len(word) > 0 {
			word = handleAcronyms(word)
			if word == strings.ToUpper(word) {
				words[i] = word
			} else {
				words[i] = strings.ToUpper(string(word[0])) + word[1:]
			}
		}
	}

	return strings.Join(words, "")
}

// toCamelCase converts a string to camelCase
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}

	i := 0
	for ; i < len(pascal); i++ {
		if i > 0 && (i+1 >= len(pascal) || unicode.IsLower(rune(pascal[i+1]))) {
			break
		}
	}
	if i > 1 {
		i--
	}

	return strings.ToLower(pascal[:i]) + pascal[i:]
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// handleAcronyms handles common acronyms
func handleAcronyms(word string) string {
	upper := strings.ToUpper(word)
	acronyms := map[string]bool{
		"ID": true, "URL": true, "URI": true, "API": true,
		"HTTP": true, "HTTPS": true, "JSON": true, "XML": true,
		"HTML": true, "CSS": true, "SQL": true, "UUID": true,
		"IP": true, "TCP": true, "UDP": true, "DNS": true,
		"TLS": true, "SSL": true, "SSH": true, "JWT": true,
	}

	if acronyms[upper] {
		return upper
	}
	return word
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var neg bool
	if i < 0 {
		neg = true
		i = -i
	}
	var digits []byte
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}

func ftoa(f float64) string {
	return itoa(int(f))
}
