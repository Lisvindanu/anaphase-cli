package generator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// HandlerConfig holds configuration for handler generation
type HandlerConfig struct {
	Protocol string
	Auth     bool
	Validate bool
	Logger   *slog.Logger
}

// HandlerGenerator generates HTTP/gRPC handlers
type HandlerGenerator struct {
	domainName string
	config     *HandlerConfig
}

// NewHandlerGenerator creates a new handler generator
func NewHandlerGenerator(domainName string, config *HandlerConfig) *HandlerGenerator {
	return &HandlerGenerator{
		domainName: domainName,
		config:     config,
	}
}

// Generate creates handler files
func (g *HandlerGenerator) Generate(ctx context.Context) ([]string, error) {
	var generatedFiles []string

	// Create output directory
	outputDir := filepath.Join("internal", "adapter", "handler", g.config.Protocol)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	// Generate DTO file
	dtoFile, err := g.generateDTO(outputDir)
	if err != nil {
		return nil, fmt.Errorf("generate DTO: %w", err)
	}
	generatedFiles = append(generatedFiles, dtoFile)

	// Generate handler file
	handlerFile, err := g.generateHandler(outputDir)
	if err != nil {
		return nil, fmt.Errorf("generate handler: %w", err)
	}
	generatedFiles = append(generatedFiles, handlerFile)

	// Generate test file
	testFile, err := g.generateHandlerTest(outputDir)
	if err != nil {
		return nil, fmt.Errorf("generate handler test: %w", err)
	}
	generatedFiles = append(generatedFiles, testFile)

	return generatedFiles, nil
}

func (g *HandlerGenerator) generateDTO(outputDir string) (string, error) {
	filename := filepath.Join(outputDir, g.domainName+"_dto.go")

	var b strings.Builder

	// Package
	b.WriteString(fmt.Sprintf("package %s\n\n", g.config.Protocol))

	// Request DTOs
	entityName := toPascalCase(g.domainName)
	b.WriteString(fmt.Sprintf("// Create%sRequest represents HTTP request to create %s\n", entityName, g.domainName))
	b.WriteString(fmt.Sprintf("type Create%sRequest struct {\n", entityName))
	b.WriteString("\t// TODO: Add fields based on domain entity\n")
	b.WriteString("}\n\n")

	b.WriteString(fmt.Sprintf("// Update%sRequest represents HTTP request to update %s\n", entityName, g.domainName))
	b.WriteString(fmt.Sprintf("type Update%sRequest struct {\n", entityName))
	b.WriteString("\t// TODO: Add fields based on domain entity\n")
	b.WriteString("}\n\n")

	// Response DTOs
	b.WriteString(fmt.Sprintf("// %sResponse represents HTTP response with %s data\n", entityName, g.domainName))
	b.WriteString(fmt.Sprintf("type %sResponse struct {\n", entityName))
	b.WriteString("\tID        string `json:\"id\"`\n")
	b.WriteString("\tCreatedAt string `json:\"created_at\"`\n")
	b.WriteString("\tUpdatedAt string `json:\"updated_at\"`\n")
	b.WriteString("\t// TODO: Add fields based on domain entity\n")
	b.WriteString("}\n\n")

	// Error response
	b.WriteString("// ErrorResponse represents standard error response\n")
	b.WriteString("type ErrorResponse struct {\n")
	b.WriteString("\tError   string            `json:\"error\"`\n")
	b.WriteString("\tMessage string            `json:\"message\"`\n")
	b.WriteString("\tDetails map[string]string `json:\"details,omitempty\"`\n")
	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *HandlerGenerator) generateHandler(outputDir string) (string, error) {
	filename := filepath.Join(outputDir, g.domainName+"_handler.go")

	var b strings.Builder

	// Package
	b.WriteString(fmt.Sprintf("package %s\n\n", g.config.Protocol))

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"log/slog\"\n")
	b.WriteString("\t\"net/http\"\n\n")
	b.WriteString("\t\"github.com/go-chi/chi/v5\"\n")
	b.WriteString("\t\"github.com/google/uuid\"\n\n")
	b.WriteString("\t\"github.com/lisvindanu/anaphase-cli/internal/core/port\"\n")
	b.WriteString(")\n\n")

	// Handler struct
	entityName := toPascalCase(g.domainName)
	b.WriteString(fmt.Sprintf("// %sHandler handles HTTP requests for %s domain\n", entityName, g.domainName))
	b.WriteString(fmt.Sprintf("type %sHandler struct {\n", entityName))
	b.WriteString(fmt.Sprintf("\tservice port.%sService\n", entityName))
	b.WriteString("\tlogger  *slog.Logger\n")
	b.WriteString("}\n\n")

	// Constructor
	b.WriteString(fmt.Sprintf("// New%sHandler creates a new %s handler\n", entityName, g.domainName))
	b.WriteString(fmt.Sprintf("func New%sHandler(service port.%sService, logger *slog.Logger) *%sHandler {\n", entityName, entityName, entityName))
	b.WriteString(fmt.Sprintf("\treturn &%sHandler{\n", entityName))
	b.WriteString("\t\tservice: service,\n")
	b.WriteString("\t\tlogger:  logger,\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n\n")

	// RegisterRoutes
	b.WriteString("// RegisterRoutes registers all routes for this handler\n")
	b.WriteString(fmt.Sprintf("func (h *%sHandler) RegisterRoutes(r chi.Router) {\n", entityName))
	b.WriteString(fmt.Sprintf("\tr.Route(\"/%s\", func(r chi.Router) {\n", strings.ToLower(g.domainName)+"s"))
	b.WriteString("\t\tr.Post(\"/\", h.Create)\n")
	b.WriteString("\t\tr.Get(\"/{id}\", h.GetByID)\n")
	b.WriteString("\t\tr.Put(\"/{id}\", h.Update)\n")
	b.WriteString("\t\tr.Delete(\"/{id}\", h.Delete)\n")
	b.WriteString("\t})\n")
	b.WriteString("}\n\n")

	// Create handler
	b.WriteString(fmt.Sprintf("// Create creates a new %s\n", g.domainName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) Create(w http.ResponseWriter, r *http.Request) {\n", entityName))
	b.WriteString("\tctx := r.Context()\n\n")
	b.WriteString(fmt.Sprintf("\tvar req Create%sRequest\n", entityName))
	b.WriteString("\tif err := json.NewDecoder(r.Body).Decode(&req); err != nil {\n")
	b.WriteString("\t\th.respondError(w, http.StatusBadRequest, \"invalid request body\", err)\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\t// TODO: Call service to create entity\n")
	b.WriteString("\t_ = ctx\n")
	b.WriteString("\t_ = req\n\n")
	b.WriteString(fmt.Sprintf("\th.respondJSON(w, http.StatusCreated, %sResponse{})\n", entityName))
	b.WriteString("}\n\n")

	// GetByID handler
	b.WriteString(fmt.Sprintf("// GetByID retrieves %s by ID\n", g.domainName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) GetByID(w http.ResponseWriter, r *http.Request) {\n", entityName))
	b.WriteString("\tctx := r.Context()\n")
	b.WriteString("\tid := chi.URLParam(r, \"id\")\n\n")
	b.WriteString("\tuuid, err := uuid.Parse(id)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\th.respondError(w, http.StatusBadRequest, \"invalid ID\", err)\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\t// TODO: Call service to get entity\n")
	b.WriteString("\t_ = ctx\n")
	b.WriteString("\t_ = uuid\n\n")
	b.WriteString(fmt.Sprintf("\th.respondJSON(w, http.StatusOK, %sResponse{})\n", entityName))
	b.WriteString("}\n\n")

	// Update handler
	b.WriteString(fmt.Sprintf("// Update updates an existing %s\n", g.domainName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) Update(w http.ResponseWriter, r *http.Request) {\n", entityName))
	b.WriteString("\tctx := r.Context()\n")
	b.WriteString("\tid := chi.URLParam(r, \"id\")\n\n")
	b.WriteString("\tuuid, err := uuid.Parse(id)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\th.respondError(w, http.StatusBadRequest, \"invalid ID\", err)\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")
	b.WriteString(fmt.Sprintf("\tvar req Update%sRequest\n", entityName))
	b.WriteString("\tif err := json.NewDecoder(r.Body).Decode(&req); err != nil {\n")
	b.WriteString("\t\th.respondError(w, http.StatusBadRequest, \"invalid request body\", err)\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\t// TODO: Call service to update entity\n")
	b.WriteString("\t_ = ctx\n")
	b.WriteString("\t_ = uuid\n")
	b.WriteString("\t_ = req\n\n")
	b.WriteString(fmt.Sprintf("\th.respondJSON(w, http.StatusOK, %sResponse{})\n", entityName))
	b.WriteString("}\n\n")

	// Delete handler
	b.WriteString(fmt.Sprintf("// Delete removes a %s\n", g.domainName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) Delete(w http.ResponseWriter, r *http.Request) {\n", entityName))
	b.WriteString("\tctx := r.Context()\n")
	b.WriteString("\tid := chi.URLParam(r, \"id\")\n\n")
	b.WriteString("\tuuid, err := uuid.Parse(id)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\th.respondError(w, http.StatusBadRequest, \"invalid ID\", err)\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\t// TODO: Call service to delete entity\n")
	b.WriteString("\t_ = ctx\n")
	b.WriteString("\t_ = uuid\n\n")
	b.WriteString("\tw.WriteHeader(http.StatusNoContent)\n")
	b.WriteString("}\n\n")

	// Helper methods
	b.WriteString("// respondJSON sends a JSON response\n")
	b.WriteString(fmt.Sprintf("func (h *%sHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {\n", entityName))
	b.WriteString("\tw.Header().Set(\"Content-Type\", \"application/json\")\n")
	b.WriteString("\tw.WriteHeader(status)\n")
	b.WriteString("\tif err := json.NewEncoder(w).Encode(data); err != nil {\n")
	b.WriteString("\t\th.logger.Error(\"failed to encode response\", \"error\", err)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n\n")

	b.WriteString("// respondError sends an error response\n")
	b.WriteString(fmt.Sprintf("func (h *%sHandler) respondError(w http.ResponseWriter, status int, message string, err error) {\n", entityName))
	b.WriteString("\th.logger.Error(message, \"error\", err)\n")
	b.WriteString("\th.respondJSON(w, status, ErrorResponse{\n")
	b.WriteString("\t\tError:   http.StatusText(status),\n")
	b.WriteString("\t\tMessage: message,\n")
	b.WriteString("\t})\n")
	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *HandlerGenerator) generateHandlerTest(outputDir string) (string, error) {
	filename := filepath.Join(outputDir, g.domainName+"_handler_test.go")

	var b strings.Builder

	// Package
	b.WriteString(fmt.Sprintf("package %s_test\n\n", g.config.Protocol))

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"testing\"\n\n")
	b.WriteString("\t\"github.com/stretchr/testify/assert\"\n")
	b.WriteString(")\n\n")

	// Test placeholder
	entityName := toPascalCase(g.domainName)
	b.WriteString(fmt.Sprintf("func Test%sHandler_Create(t *testing.T) {\n", entityName))
	b.WriteString("\t// TODO: Implement handler tests\n")
	b.WriteString("\tassert.True(t, true)\n")
	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

// toPascalCase converts string to PascalCase
func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
