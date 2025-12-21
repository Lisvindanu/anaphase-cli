package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lisvindanuu/anaphase-cli/internal/ai"
	"github.com/lisvindanuu/anaphase-cli/pkg/fileutil"
)

// DomainGenerator generates domain code files from AI spec
type DomainGenerator struct {
	spec      *ai.DomainSpec
	outputDir string
}

// NewDomainGenerator creates a new domain generator
func NewDomainGenerator(spec *ai.DomainSpec, outputDir string) *DomainGenerator {
	return &DomainGenerator{
		spec:      spec,
		outputDir: outputDir,
	}
}

// Generate creates all domain files
func (g *DomainGenerator) Generate() ([]string, error) {
	var generatedFiles []string

	// Create directory structure
	if err := g.createDirectories(); err != nil {
		return nil, fmt.Errorf("create directories: %w", err)
	}

	// Generate entities
	for _, entity := range g.spec.Entities {
		file, err := g.generateEntity(entity)
		if err != nil {
			return nil, fmt.Errorf("generate entity %s: %w", entity.Name, err)
		}
		generatedFiles = append(generatedFiles, file)
	}

	// Generate value objects
	for _, vo := range g.spec.ValueObjects {
		file, err := g.generateValueObject(vo)
		if err != nil {
			return nil, fmt.Errorf("generate value object %s: %w", vo.Name, err)
		}
		generatedFiles = append(generatedFiles, file)
	}

	// Generate repository interface
	file, err := g.generateRepository()
	if err != nil {
		return nil, fmt.Errorf("generate repository: %w", err)
	}
	generatedFiles = append(generatedFiles, file)

	// Generate service interface
	file, err = g.generateService()
	if err != nil {
		return nil, fmt.Errorf("generate service: %w", err)
	}
	generatedFiles = append(generatedFiles, file)

	return generatedFiles, nil
}

func (g *DomainGenerator) createDirectories() error {
	dirs := []string{
		filepath.Join(g.outputDir, "entity"),
		filepath.Join(g.outputDir, "valueobject"),
		filepath.Join(g.outputDir, "port"),
		filepath.Join(g.outputDir, "service"),
	}

	for _, dir := range dirs {
		if err := fileutil.EnsureDir(dir); err != nil {
			return err
		}
	}

	return nil
}

func (g *DomainGenerator) generateEntity(entity ai.EntitySpec) (string, error) {
	filename := filepath.Join(g.outputDir, "entity", toSnakeCase(entity.Name)+".go")

	var b strings.Builder

	// Detect required imports
	needsValueObject := false
	needsTime := false
	for _, field := range entity.Fields {
		// Check if field type is a value object (PascalCase, not a standard Go type)
		if isValueObjectType(field.Type) {
			needsValueObject = true
		}
		// Check if we need time package
		if field.Type == "time.Time" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" {
			needsTime = true
		}
	}

	// Package and imports
	b.WriteString("package entity\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"errors\"\n")
	if needsTime {
		b.WriteString("\t\"time\"\n")
	}
	b.WriteString("\n\t\"github.com/google/uuid\"\n")
	if needsValueObject {
		b.WriteString("\t\"github.com/lisvindanuu/anaphase-cli/internal/core/valueobject\"\n")
	}
	b.WriteString(")\n\n")

	// Error definitions
	b.WriteString("// Common errors\n")
	b.WriteString("var (\n")
	b.WriteString(fmt.Sprintf("\tErr%sNotFound = errors.New(\"%s not found\")\n", entity.Name, strings.ToLower(entity.Name)))
	b.WriteString(fmt.Sprintf("\tErrInvalid%s = errors.New(\"invalid %s\")\n", entity.Name, strings.ToLower(entity.Name)))
	b.WriteString(")\n\n")

	// Struct definition
	if entity.IsAggregateRoot {
		b.WriteString(fmt.Sprintf("// %s is an aggregate root\n", entity.Name))
	} else {
		b.WriteString(fmt.Sprintf("// %s represents a %s entity\n", entity.Name, strings.ToLower(entity.Name)))
	}
	b.WriteString(fmt.Sprintf("type %s struct {\n", entity.Name))

	for _, field := range entity.Fields {
		fieldType := field.Type
		// Add valueobject package prefix if it's a value object and not already qualified
		if isValueObjectType(field.Type) && !strings.HasPrefix(field.Type, "valueobject.") {
			fieldType = "valueobject." + field.Type
		}

		if field.Description != "" {
			b.WriteString(fmt.Sprintf("\t%s %s // %s\n", field.Name, fieldType, field.Description))
		} else {
			b.WriteString(fmt.Sprintf("\t%s %s\n", field.Name, fieldType))
		}
	}

	b.WriteString("}\n\n")

	// Constructor
	b.WriteString(fmt.Sprintf("// New%s creates a new %s\n", entity.Name, strings.ToLower(entity.Name)))
	b.WriteString(fmt.Sprintf("func New%s() *%s {\n", entity.Name, entity.Name))
	b.WriteString(fmt.Sprintf("\treturn &%s{\n", entity.Name))
	b.WriteString("\t\tID: uuid.New(),\n")

	// Only add CreatedAt if it exists in fields
	hasCreatedAt := false
	hasUpdatedAt := false
	for _, field := range entity.Fields {
		if field.Name == "CreatedAt" {
			hasCreatedAt = true
		}
		if field.Name == "UpdatedAt" {
			hasUpdatedAt = true
		}
	}

	if hasCreatedAt {
		b.WriteString("\t\tCreatedAt: time.Now(),\n")
	}
	if hasUpdatedAt {
		b.WriteString("\t\tUpdatedAt: time.Now(),\n")
	}

	b.WriteString("\t}\n")
	b.WriteString("}\n\n")

	// Validate method
	b.WriteString(fmt.Sprintf("// Validate validates the %s\n", strings.ToLower(entity.Name)))
	b.WriteString(fmt.Sprintf("func (e *%s) Validate() error {\n", entity.Name))
	b.WriteString("\tif e.ID == uuid.Nil {\n")
	b.WriteString(fmt.Sprintf("\t\treturn ErrInvalid%s\n", entity.Name))
	b.WriteString("\t}\n")

	for _, field := range entity.Fields {
		if field.Validation != "" {
			b.WriteString(fmt.Sprintf("\t// Validate %s: %s\n", field.Name, field.Validation))
		}
	}

	b.WriteString("\treturn nil\n")
	b.WriteString("}\n")

	// Methods
	for _, method := range entity.Methods {
		b.WriteString("\n")
		if method.Description != "" {
			b.WriteString(fmt.Sprintf("// %s %s\n", method.Name, method.Description))
		}
		b.WriteString(method.Signature + " {\n")
		b.WriteString("\t// TODO: Implement business logic\n")

		// Add return statement only if method has return type (doesn't end with just ")")
		trimmedSig := strings.TrimSpace(method.Signature)
		if !strings.HasSuffix(trimmedSig, ")") {
			b.WriteString("\treturn nil\n")
		}

		b.WriteString("}\n")
	}

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *DomainGenerator) generateValueObject(vo ai.ValueObjectSpec) (string, error) {
	filename := filepath.Join(g.outputDir, "valueobject", toSnakeCase(vo.Name)+".go")

	var b strings.Builder

	// Package
	b.WriteString("package valueobject\n\n")

	// Struct
	b.WriteString(fmt.Sprintf("// %s is a value object\n", vo.Name))
	b.WriteString(fmt.Sprintf("type %s struct {\n", vo.Name))

	for _, field := range vo.Fields {
		b.WriteString(fmt.Sprintf("\t%s %s", field.Name, field.Type))
		if field.Description != "" {
			b.WriteString(fmt.Sprintf(" // %s", field.Description))
		}
		b.WriteString("\n")
	}

	b.WriteString("}\n\n")

	// Constructor
	b.WriteString(fmt.Sprintf("// New%s creates a new %s\n", vo.Name, vo.Name))
	b.WriteString(fmt.Sprintf("func New%s(", vo.Name))

	// Constructor parameters
	for i, field := range vo.Fields {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%s %s", strings.ToLower(field.Name), field.Type))
	}

	b.WriteString(fmt.Sprintf(") (%s, error) {\n", vo.Name))
	b.WriteString(fmt.Sprintf("\tv := %s{\n", vo.Name))

	for _, field := range vo.Fields {
		b.WriteString(fmt.Sprintf("\t\t%s: %s,\n", field.Name, strings.ToLower(field.Name)))
	}

	b.WriteString("\t}\n\n")
	b.WriteString("\tif err := v.Validate(); err != nil {\n")
	b.WriteString(fmt.Sprintf("\t\treturn %s{}, err\n", vo.Name))
	b.WriteString("\t}\n\n")
	b.WriteString("\treturn v, nil\n")
	b.WriteString("}\n\n")

	// Validate
	b.WriteString(fmt.Sprintf("// Validate validates the %s\n", vo.Name))
	b.WriteString(fmt.Sprintf("func (v %s) Validate() error {\n", vo.Name))
	if vo.Validation != "" {
		b.WriteString(fmt.Sprintf("\t// %s\n", vo.Validation))
	}
	b.WriteString("\t// TODO: Add validation logic\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *DomainGenerator) generateRepository() (string, error) {
	filename := filepath.Join(g.outputDir, "port", toSnakeCase(g.spec.RepositoryInterface.Name)+".go")

	var b strings.Builder

	// Package and imports
	b.WriteString("package port\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n\n")
	b.WriteString("\t\"github.com/google/uuid\"\n")
	b.WriteString("\t\"github.com/lisvindanuu/anaphase-cli/internal/core/entity\"\n")
	b.WriteString("\t\"github.com/lisvindanuu/anaphase-cli/internal/core/valueobject\"\n")
	b.WriteString(")\n\n")

	// Interface
	b.WriteString(fmt.Sprintf("// %s defines the contract for %s persistence\n",
		g.spec.RepositoryInterface.Name, strings.ToLower(g.spec.DomainName)))
	b.WriteString(fmt.Sprintf("type %s interface {\n", g.spec.RepositoryInterface.Name))

	for _, method := range g.spec.RepositoryInterface.Methods {
		if method.Description != "" {
			b.WriteString(fmt.Sprintf("\t// %s %s\n", method.Name, method.Description))
		}
		b.WriteString(fmt.Sprintf("\t%s\n\n", method.Signature))
	}

	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *DomainGenerator) generateService() (string, error) {
	filename := filepath.Join(g.outputDir, "port", toSnakeCase(g.spec.ServiceInterface.Name)+".go")

	var b strings.Builder

	// Package and imports
	b.WriteString("package port\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n\n")
	b.WriteString("\t\"github.com/google/uuid\"\n")
	b.WriteString("\t\"github.com/lisvindanuu/anaphase-cli/internal/core/entity\"\n")
	b.WriteString("\t\"github.com/lisvindanuu/anaphase-cli/internal/core/valueobject\"\n")
	b.WriteString(")\n\n")

	// Interface
	b.WriteString(fmt.Sprintf("// %s defines the contract for %s business logic\n",
		g.spec.ServiceInterface.Name, strings.ToLower(g.spec.DomainName)))
	b.WriteString(fmt.Sprintf("type %s interface {\n", g.spec.ServiceInterface.Name))

	for _, method := range g.spec.ServiceInterface.Methods {
		if method.Description != "" {
			b.WriteString(fmt.Sprintf("\t// %s %s\n", method.Name, method.Description))
		}
		b.WriteString(fmt.Sprintf("\t%s\n\n", method.Signature))
	}

	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// isValueObjectType checks if a type is a custom value object (not a standard Go type)
func isValueObjectType(typeName string) bool {
	// Check if already qualified with valueobject package
	if strings.HasPrefix(typeName, "valueobject.") {
		return true
	}

	// Standard Go types that are NOT value objects
	standardTypes := map[string]bool{
		"string":     true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"float32":    true,
		"float64":    true,
		"bool":       true,
		"byte":       true,
		"rune":       true,
		"time.Time":  true,
		"uuid.UUID":  true,
		"error":      true,
	}

	// Check if it starts with map, []slice, or chan (built-in types)
	if strings.HasPrefix(typeName, "map[") ||
	   strings.HasPrefix(typeName, "[]") ||
	   strings.HasPrefix(typeName, "chan ") ||
	   strings.HasPrefix(typeName, "*") {
		return false
	}

	// If not a standard type and PascalCase, it's likely a value object
	return !standardTypes[typeName] && len(typeName) > 0 && typeName[0] >= 'A' && typeName[0] <= 'Z'
}
