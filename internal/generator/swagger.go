package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// SwaggerConfig holds configuration for Swagger generation
type SwaggerConfig struct {
	Domain      string
	Version     string
	Title       string
	Description string
	Host        string
	BasePath    string
}

// SwaggerGenerator generates Swagger annotations
type SwaggerGenerator struct {
	config     *SwaggerConfig
	entityInfo *SwaggerEntityInfo
}

// SwaggerEntityInfo holds entity information for Swagger
type SwaggerEntityInfo struct {
	EntityName           string
	EntityNameLower      string
	EntityNameLowerPlural string
	Fields               []SwaggerFieldInfo
}

// SwaggerFieldInfo holds field information for Swagger
type SwaggerFieldInfo struct {
	Name        string
	JSONName    string
	Type        string
	SwaggerType string
	Required    bool
	Example     string
}

// NewSwaggerGenerator creates a new Swagger generator
func NewSwaggerGenerator(config *SwaggerConfig) *SwaggerGenerator {
	return &SwaggerGenerator{
		config: config,
	}
}

// Generate generates Swagger annotations
func (g *SwaggerGenerator) Generate() error {
	// Scan entity
	if err := g.scanEntity(); err != nil {
		return fmt.Errorf("scan entity: %w", err)
	}

	// Add annotations to handler
	if err := g.addHandlerAnnotations(); err != nil {
		return fmt.Errorf("add handler annotations: %w", err)
	}

	// Generate main.go annotations
	if err := g.addMainAnnotations(); err != nil {
		return fmt.Errorf("add main annotations: %w", err)
	}

	return nil
}

// scanEntity scans the entity file
func (g *SwaggerGenerator) scanEntity() error {
	entityFile := filepath.Join("internal", "core", "entity", g.config.Domain+".go")

	if _, err := os.Stat(entityFile); os.IsNotExist(err) {
		return fmt.Errorf("entity file not found: %s", entityFile)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, entityFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse entity file: %w", err)
	}

	info := &SwaggerEntityInfo{
		EntityName:            toPascalCase(g.config.Domain),
		EntityNameLower:       strings.ToLower(g.config.Domain),
		EntityNameLowerPlural: strings.ToLower(g.config.Domain) + "s",
		Fields:                []SwaggerFieldInfo{},
	}

	// Find entity struct
	ast.Inspect(file, func(n ast.Node) bool {
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

			// Get JSON tag
			jsonName := strings.ToLower(fieldName)
			if field.Tag != nil {
				tag := field.Tag.Value
				if strings.Contains(tag, "json:") {
					parts := strings.Split(tag, "json:\"")
					if len(parts) > 1 {
						jsonName = strings.Split(parts[1], "\"")[0]
					}
				}
			}

			fieldInfo := SwaggerFieldInfo{
				Name:        fieldName,
				JSONName:    jsonName,
				Type:        fieldType,
				SwaggerType: getSwaggerType(fieldType),
				Required:    isRequiredField(fieldName),
				Example:     getSwaggerExample(fieldType, fieldName),
			}

			info.Fields = append(info.Fields, fieldInfo)
		}

		return true
	})

	g.entityInfo = info
	return nil
}

// addHandlerAnnotations adds Swagger annotations to handler
func (g *SwaggerGenerator) addHandlerAnnotations() error {
	handlerFile := filepath.Join("internal", "adapter", "handler", "http", g.config.Domain+"_handler.go")

	content, err := os.ReadFile(handlerFile)
	if err != nil {
		return fmt.Errorf("read handler file: %w", err)
	}

	// Insert annotations before each handler method
	updated := string(content)

	// Add Create annotations
	createAnnotation := g.generateCreateAnnotation()
	updated = strings.Replace(updated,
		"func (h *"+g.entityInfo.EntityName+"Handler) Create(",
		createAnnotation+"func (h *"+g.entityInfo.EntityName+"Handler) Create(",
		1)

	// Add GetByID annotations
	getAnnotation := g.generateGetAnnotation()
	updated = strings.Replace(updated,
		"func (h *"+g.entityInfo.EntityName+"Handler) GetByID(",
		getAnnotation+"func (h *"+g.entityInfo.EntityName+"Handler) GetByID(",
		1)

	// Add Update annotations
	updateAnnotation := g.generateUpdateAnnotation()
	updated = strings.Replace(updated,
		"func (h *"+g.entityInfo.EntityName+"Handler) Update(",
		updateAnnotation+"func (h *"+g.entityInfo.EntityName+"Handler) Update(",
		1)

	// Add Delete annotations
	deleteAnnotation := g.generateDeleteAnnotation()
	updated = strings.Replace(updated,
		"func (h *"+g.entityInfo.EntityName+"Handler) Delete(",
		deleteAnnotation+"func (h *"+g.entityInfo.EntityName+"Handler) Delete(",
		1)

	// Write back
	if err := os.WriteFile(handlerFile, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write handler file: %w", err)
	}

	return nil
}

// addMainAnnotations adds Swagger annotations to main.go
func (g *SwaggerGenerator) addMainAnnotations() error {
	mainFile := filepath.Join("cmd", "api", "main.go")

	content, err := os.ReadFile(mainFile)
	if err != nil {
		return fmt.Errorf("read main file: %w", err)
	}

	// Check if annotations already exist
	if strings.Contains(string(content), "@title") {
		return nil // Already annotated
	}

	// Generate main annotations
	mainAnnotation := fmt.Sprintf(`// @title %s
// @version %s
// @description %s
// @host %s
// @BasePath %s
// @schemes http https

`, g.config.Title, g.config.Version, g.config.Description, g.config.Host, g.config.BasePath)

	// Insert before package declaration
	updated := strings.Replace(string(content), "package main", mainAnnotation+"package main", 1)

	// Add swagger import if not present
	if !strings.Contains(updated, "github.com/swaggo/http-swagger") {
		// Add import
		updated = strings.Replace(updated,
			`"github.com/go-chi/chi/v5/middleware"`,
			`"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/lisvindanuu/anaphase-cli/docs" // swagger docs`,
			1)

		// Add swagger route
		updated = strings.Replace(updated,
			`// Health check`,
			`// Swagger documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Health check`,
			1)
	}

	// Write back
	if err := os.WriteFile(mainFile, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write main file: %w", err)
	}

	return nil
}

// generateCreateAnnotation generates Create endpoint annotation
func (g *SwaggerGenerator) generateCreateAnnotation() string {
	return fmt.Sprintf(`
// Create godoc
// @Summary Create a new %s
// @Description Create a new %s with the provided data
// @Tags %s
// @Accept json
// @Produce json
// @Param %s body Create%sRequest true "Create %s"
// @Success 201 {object} %sResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /%s [post]
`,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLowerPlural,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityName,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityName,
		g.entityInfo.EntityNameLowerPlural,
	)
}

// generateGetAnnotation generates GetByID endpoint annotation
func (g *SwaggerGenerator) generateGetAnnotation() string {
	return fmt.Sprintf(`
// GetByID godoc
// @Summary Get %s by ID
// @Description Get a %s by its ID
// @Tags %s
// @Accept json
// @Produce json
// @Param id path string true "%s ID"
// @Success 200 {object} %sResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /%s/{id} [get]
`,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLowerPlural,
		g.entityInfo.EntityName,
		g.entityInfo.EntityName,
		g.entityInfo.EntityNameLowerPlural,
	)
}

// generateUpdateAnnotation generates Update endpoint annotation
func (g *SwaggerGenerator) generateUpdateAnnotation() string {
	return fmt.Sprintf(`
// Update godoc
// @Summary Update %s
// @Description Update a %s by its ID
// @Tags %s
// @Accept json
// @Produce json
// @Param id path string true "%s ID"
// @Param %s body Update%sRequest true "Update %s"
// @Success 200 {object} %sResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /%s/{id} [put]
`,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLowerPlural,
		g.entityInfo.EntityName,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityName,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityName,
		g.entityInfo.EntityNameLowerPlural,
	)
}

// generateDeleteAnnotation generates Delete endpoint annotation
func (g *SwaggerGenerator) generateDeleteAnnotation() string {
	return fmt.Sprintf(`
// Delete godoc
// @Summary Delete %s
// @Description Delete a %s by its ID
// @Tags %s
// @Accept json
// @Produce json
// @Param id path string true "%s ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /%s/{id} [delete]
`,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLower,
		g.entityInfo.EntityNameLowerPlural,
		g.entityInfo.EntityName,
		g.entityInfo.EntityNameLowerPlural,
	)
}

// Helper functions

func getSwaggerType(goType string) string {
	switch {
	case strings.Contains(goType, "string"):
		return "string"
	case strings.Contains(goType, "int"):
		return "integer"
	case strings.Contains(goType, "float"):
		return "number"
	case strings.Contains(goType, "bool"):
		return "boolean"
	case strings.Contains(goType, "time.Time"):
		return "string"
	case strings.Contains(goType, "uuid.UUID"):
		return "string"
	default:
		return "string"
	}
}

func getSwaggerExample(fieldType, fieldName string) string {
	switch {
	case strings.Contains(fieldType, "Email"):
		return "user@example.com"
	case strings.Contains(fieldType, "Money"):
		return "99.99"
	case strings.Contains(fieldType, "Phone"):
		return "+1234567890"
	case fieldType == "string":
		return "example-" + strings.ToLower(fieldName)
	case strings.Contains(fieldType, "int"):
		return "42"
	case strings.Contains(fieldType, "float"):
		return "99.99"
	case strings.Contains(fieldType, "bool"):
		return "true"
	case strings.Contains(fieldType, "uuid.UUID"):
		return "550e8400-e29b-41d4-a716-446655440000"
	case strings.Contains(fieldType, "time.Time"):
		return "2024-01-15T10:30:00Z"
	default:
		return ""
	}
}

func isRequiredField(fieldName string) bool {
	nonRequired := []string{"ID", "CreatedAt", "UpdatedAt"}
	for _, nr := range nonRequired {
		if fieldName == nr {
			return false
		}
	}
	return true
}
