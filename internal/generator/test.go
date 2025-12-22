package generator

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Note: embed path is relative to this file
// We'll load templates from files directly instead of embedding
// var testTemplates embed.FS

// TestConfig holds configuration for test generation
type TestConfig struct {
	Domain   string
	TestType string
}

// TestGenerator generates tests
type TestGenerator struct {
	config     *TestConfig
	entityInfo *EntityInfo
}

// EntityInfo holds information about the entity
type EntityInfo struct {
	EntityName            string
	EntityNameLower       string
	EntityNameLowerPlural string
	Module                string
	Fields                []FieldInfo
	Methods               []MethodInfo
	HasUniqueFields       bool
}

// FieldInfo holds field information
type FieldInfo struct {
	Name             string
	Type             string
	TestValue        string
	EmptyValue       string
	UpdatedTestValue string
	IsUnique         bool
}

// MethodInfo holds method information
type MethodInfo struct {
	Name string
}

// NewTestGenerator creates a new test generator
func NewTestGenerator(config *TestConfig) *TestGenerator {
	return &TestGenerator{
		config: config,
	}
}

// ScanDomain scans the domain to extract entity information
func (g *TestGenerator) ScanDomain() error {
	entityFile := filepath.Join("internal", "core", "entity", g.config.Domain+".go")

	if _, err := os.Stat(entityFile); os.IsNotExist(err) {
		return fmt.Errorf("entity file not found: %s", entityFile)
	}

	// Parse the entity file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, entityFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse entity file: %w", err)
	}

	// Extract entity information
	info := &EntityInfo{
		EntityName:            toPascalCase(g.config.Domain),
		EntityNameLower:       strings.ToLower(g.config.Domain),
		EntityNameLowerPlural: strings.ToLower(g.config.Domain) + "s",
		Module:                getModuleName(),
		Fields:                []FieldInfo{},
		Methods:               []MethodInfo{},
	}

	// Find the entity struct
	ast.Inspect(file, func(n ast.Node) bool {
		// Find struct type
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != info.EntityName {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Extract fields
		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			fieldName := field.Names[0].Name
			fieldType := formatType(field.Type)

			fieldInfo := FieldInfo{
				Name:             fieldName,
				Type:             fieldType,
				TestValue:        getTestValue(fieldType, fieldName),
				EmptyValue:       getEmptyValue(fieldType),
				UpdatedTestValue: getUpdatedTestValue(fieldType, fieldName),
				IsUnique:         isUniqueField(fieldName),
			}

			info.Fields = append(info.Fields, fieldInfo)

			if fieldInfo.IsUnique {
				info.HasUniqueFields = true
			}
		}

		return true
	})

	// Find methods
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil {
			return true
		}

		info.Methods = append(info.Methods, MethodInfo{
			Name: funcDecl.Name.Name,
		})

		return true
	})

	g.entityInfo = info
	return nil
}

// GenerateAllTests generates all types of tests
func (g *TestGenerator) GenerateAllTests(ctx context.Context) ([]string, error) {
	var files []string

	// Generate entity tests
	entityFile, err := g.generateEntityTests()
	if err != nil {
		return nil, fmt.Errorf("generate entity tests: %w", err)
	}
	files = append(files, entityFile)

	// Generate repository tests
	repoFile, err := g.generateRepositoryTests()
	if err != nil {
		return nil, fmt.Errorf("generate repository tests: %w", err)
	}
	files = append(files, repoFile)

	// Generate handler tests
	handlerFile, err := g.generateHandlerTests()
	if err != nil {
		return nil, fmt.Errorf("generate handler tests: %w", err)
	}
	files = append(files, handlerFile)

	return files, nil
}

// GenerateUnitTests generates unit tests only
func (g *TestGenerator) GenerateUnitTests(ctx context.Context) ([]string, error) {
	var files []string

	// Entity tests are unit tests
	entityFile, err := g.generateEntityTests()
	if err != nil {
		return nil, err
	}
	files = append(files, entityFile)

	// Handler unit tests
	handlerFile, err := g.generateHandlerTests()
	if err != nil {
		return nil, err
	}
	files = append(files, handlerFile)

	return files, nil
}

// GenerateIntegrationTests generates integration tests only
func (g *TestGenerator) GenerateIntegrationTests(ctx context.Context) ([]string, error) {
	var files []string

	// Repository tests are integration tests
	repoFile, err := g.generateRepositoryTests()
	if err != nil {
		return nil, err
	}
	files = append(files, repoFile)

	return files, nil
}

// generateEntityTests generates entity tests
func (g *TestGenerator) generateEntityTests() (string, error) {
	tmplFile := filepath.Join("internal", "templates", "test_entity.go.tmpl")
	tmplContent, err := os.ReadFile(tmplFile)
	if err != nil {
		return "", fmt.Errorf("read template: %w", err)
	}

	tmpl, err := template.New("test_entity").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	outputFile := filepath.Join("internal", "core", "entity", g.config.Domain+"_test.go")

	file, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, g.entityInfo); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return outputFile, nil
}

// generateRepositoryTests generates repository tests
func (g *TestGenerator) generateRepositoryTests() (string, error) {
	tmplFile := filepath.Join("internal", "templates", "test_repository.go.tmpl")
	tmplContent, err := os.ReadFile(tmplFile)
	if err != nil {
		return "", fmt.Errorf("read template: %w", err)
	}

	tmpl, err := template.New("test_repository").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	outputFile := filepath.Join("internal", "adapter", "repository", "postgres", g.config.Domain+"_repo_test.go")

	file, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, g.entityInfo); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return outputFile, nil
}

// generateHandlerTests generates handler tests
func (g *TestGenerator) generateHandlerTests() (string, error) {
	tmplFile := filepath.Join("internal", "templates", "test_handler.go.tmpl")
	tmplContent, err := os.ReadFile(tmplFile)
	if err != nil {
		return "", fmt.Errorf("read template: %w", err)
	}

	tmpl, err := template.New("test_handler").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	outputFile := filepath.Join("internal", "adapter", "handler", "http", g.config.Domain+"_handler_test.go")

	file, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, g.entityInfo); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return outputFile, nil
}

// Helper functions

func formatType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + formatType(t.X)
	case *ast.SelectorExpr:
		return formatType(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + formatType(t.Elt)
	default:
		return "unknown"
	}
}

func getTestValue(fieldType, fieldName string) string {
	switch {
	case strings.Contains(fieldType, "Email"):
		return `valueobject.NewEmail("test@example.com")`
	case strings.Contains(fieldType, "Money"):
		return `valueobject.NewMoney(99.99, "USD")`
	case strings.Contains(fieldType, "Phone"):
		return `valueobject.NewPhone("+1234567890")`
	case strings.Contains(fieldType, "Address"):
		return `valueobject.NewAddress("123 Main St", "City", "State", "12345", "Country")`
	case fieldType == "string":
		return `"test-` + strings.ToLower(fieldName) + `"`
	case fieldType == "int", fieldType == "int64":
		return "42"
	case fieldType == "float64":
		return "99.99"
	case fieldType == "bool":
		return "true"
	case fieldType == "uuid.UUID":
		return "uuid.New()"
	default:
		return "nil"
	}
}

func getEmptyValue(fieldType string) string {
	switch {
	case strings.HasPrefix(fieldType, "*"):
		return "nil"
	case fieldType == "string":
		return `""`
	case fieldType == "int", fieldType == "int64", fieldType == "float64":
		return "0"
	case fieldType == "bool":
		return "false"
	default:
		return "nil"
	}
}

func getUpdatedTestValue(fieldType, fieldName string) string {
	switch {
	case strings.Contains(fieldType, "Email"):
		return `valueobject.NewEmail("updated@example.com")`
	case strings.Contains(fieldType, "Money"):
		return `valueobject.NewMoney(199.99, "USD")`
	case fieldType == "string":
		return `"updated-` + strings.ToLower(fieldName) + `"`
	case fieldType == "int", fieldType == "int64":
		return "84"
	case fieldType == "float64":
		return "199.99"
	case fieldType == "bool":
		return "false"
	default:
		return getTestValue(fieldType, fieldName)
	}
}

func isUniqueField(fieldName string) bool {
	uniqueFields := []string{"Email", "SKU", "Username", "Code"}
	for _, unique := range uniqueFields {
		if fieldName == unique {
			return true
		}
	}
	return false
}

func getModuleName() string {
	// Read go.mod to get module name
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "myapp"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	return "myapp"
}
