package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/lisvindanuu/anaphase-cli/pkg/fileutil"
)

//go:embed templates/*
var templatesFS embed.FS

// ProjectConfig holds project generation configuration
type ProjectConfig struct {
	Name       string
	Module     string
	Database   string
	Cache      bool
	EventBus   string
	SkipDocker bool
	OutputDir  string
}

// ProjectGenerator generates new project structure
type ProjectGenerator struct {
	config    *ProjectConfig
	templates *template.Template
}

// NewProjectGenerator creates a new project generator
func NewProjectGenerator(config *ProjectConfig) *ProjectGenerator {
	return &ProjectGenerator{
		config: config,
	}
}

// Generate creates the complete project structure
func (g *ProjectGenerator) Generate() error {
	// Load templates
	if err := g.loadTemplates(); err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	// Create directory structure
	if err := g.createDirectoryStructure(); err != nil {
		return fmt.Errorf("create directory structure: %w", err)
	}

	// Generate files
	if err := g.generateFiles(); err != nil {
		return fmt.Errorf("generate files: %w", err)
	}

	return nil
}

func (g *ProjectGenerator) loadTemplates() error {
	// Parse all template files
	tmpl, err := template.New("").ParseFS(templatesFS, "templates/*.tmpl")
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}

	g.templates = tmpl
	return nil
}

func (g *ProjectGenerator) createDirectoryStructure() error {
	dirs := []string{
		g.config.OutputDir,
		filepath.Join(g.config.OutputDir, "cmd", "api"),
		filepath.Join(g.config.OutputDir, "internal", "config"),
		filepath.Join(g.config.OutputDir, "internal", "core", "entity"),
		filepath.Join(g.config.OutputDir, "internal", "core", "port"),
		filepath.Join(g.config.OutputDir, "internal", "core", "service"),
		filepath.Join(g.config.OutputDir, "internal", "core", "valueobject"),
		filepath.Join(g.config.OutputDir, "internal", "adapter", "handler", "http"),
		filepath.Join(g.config.OutputDir, "internal", "adapter", "repository", "postgres"),
		filepath.Join(g.config.OutputDir, "internal", "adapter", "integration"),
		filepath.Join(g.config.OutputDir, "internal", "server"),
		filepath.Join(g.config.OutputDir, "migrations"),
		filepath.Join(g.config.OutputDir, "docs"),
		filepath.Join(g.config.OutputDir, "pkg", "logger"),
		filepath.Join(g.config.OutputDir, "pkg", "validator"),
		filepath.Join(g.config.OutputDir, "pkg", "errors"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (g *ProjectGenerator) generateFiles() error {
	files := []struct {
		template string
		output   string
	}{
		{"main.go.tmpl", "cmd/api/main.go"},
		{"config.go.tmpl", "internal/config/config.go"},
		{"env.go.tmpl", "internal/config/env.go"},
		{"server.go.tmpl", "internal/server/server.go"},
		{"router.go.tmpl", "internal/server/router.go"},
		{"middleware.go.tmpl", "internal/server/middleware.go"},
		{"go.mod.tmpl", "go.mod"},
		{"Makefile.tmpl", "Makefile"},
		{"env.example.tmpl", ".env.example"},
		{"gitignore.tmpl", ".gitignore"},
		{"README.md.tmpl", "README.md"},
	}

	// Add docker-compose if not skipped
	if !g.config.SkipDocker {
		files = append(files, struct {
			template string
			output   string
		}{"docker-compose.yml.tmpl", "docker-compose.yml"})
	}

	// Add Dockerfile
	files = append(files, struct {
		template string
		output   string
	}{"Dockerfile.tmpl", "Dockerfile"})

	for _, f := range files {
		outputPath := filepath.Join(g.config.OutputDir, f.output)

		if err := g.generateFile(f.template, outputPath); err != nil {
			return fmt.Errorf("generate %s: %w", f.output, err)
		}
	}

	return nil
}

func (g *ProjectGenerator) generateFile(templateName, outputPath string) error {
	// Create output file
	if err := fileutil.EnsureDir(filepath.Dir(outputPath)); err != nil {
		return fmt.Errorf("ensure directory: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	// Execute template
	tmpl := g.templates.Lookup(templateName)
	if tmpl == nil {
		return fmt.Errorf("template not found: %s", templateName)
	}

	if err := tmpl.Execute(f, g.config); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
