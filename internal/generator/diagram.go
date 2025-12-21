package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// DiagramConfig holds configuration for diagram generation
type DiagramConfig struct {
	Format string // mermaid, ascii, both
	Type   string // all, domain, layers, dependencies
}

// DiagramGenerator generates architecture diagrams
type DiagramGenerator struct {
	config  *DiagramConfig
	domains []DomainInfo
}

// DomainInfo holds information about a domain
type DomainInfo struct {
	Name         string
	HasEntity    bool
	HasRepo      bool
	HasHandler   bool
	HasService   bool
	Dependencies []string
}

// NewDiagramGenerator creates a new diagram generator
func NewDiagramGenerator(config *DiagramConfig) *DiagramGenerator {
	return &DiagramGenerator{
		config:  config,
		domains: []DomainInfo{},
	}
}

// Generate generates the architecture diagram
func (g *DiagramGenerator) Generate() (string, error) {
	// Scan project
	if err := g.scanProject(); err != nil {
		return "", fmt.Errorf("scan project: %w", err)
	}

	// Generate diagram based on format
	var output strings.Builder

	switch g.config.Format {
	case "mermaid":
		output.WriteString(g.generateMermaid())
	case "ascii":
		output.WriteString(g.generateASCII())
	case "both":
		output.WriteString("# Mermaid Diagram\n\n")
		output.WriteString(g.generateMermaid())
		output.WriteString("\n\n# ASCII Diagram\n\n")
		output.WriteString(g.generateASCII())
	default:
		output.WriteString(g.generateMermaid())
	}

	return output.String(), nil
}

// scanProject scans the project structure
func (g *DiagramGenerator) scanProject() error {
	// Scan entities
	entities, err := g.scanEntities()
	if err != nil {
		return err
	}

	// Build domain info
	for _, entity := range entities {
		domain := DomainInfo{
			Name:      entity,
			HasEntity: true,
		}

		// Check for repository
		repoFile := filepath.Join("internal", "adapter", "repository", "postgres", entity+"_repo.go")
		if _, err := os.Stat(repoFile); err == nil {
			domain.HasRepo = true
		}

		// Check for handler
		handlerFile := filepath.Join("internal", "adapter", "handler", "http", entity+"_handler.go")
		if _, err := os.Stat(handlerFile); err == nil {
			domain.HasHandler = true
		}

		// Check for service
		serviceFile := filepath.Join("internal", "core", "service", entity+"_service.go")
		if _, err := os.Stat(serviceFile); err == nil {
			domain.HasService = true
		}

		g.domains = append(g.domains, domain)
	}

	return nil
}

// scanEntities scans for entities
func (g *DiagramGenerator) scanEntities() ([]string, error) {
	entityDir := "internal/core/entity"
	var entities []string

	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return entities, nil
	}

	fset := token.NewFileSet()

	err := filepath.WalkDir(entityDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		ast.Inspect(file, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				entityName := strings.ToLower(typeSpec.Name.Name)
				entities = append(entities, entityName)
			}

			return true
		})

		return nil
	})

	return entities, err
}

// generateMermaid generates a Mermaid diagram
func (g *DiagramGenerator) generateMermaid() string {
	var b strings.Builder

	b.WriteString("```mermaid\n")

	switch g.config.Type {
	case "layers":
		b.WriteString(g.generateLayersDiagram())
	case "domain":
		b.WriteString(g.generateDomainDiagram())
	case "dependencies":
		b.WriteString(g.generateDependencyDiagram())
	default: // "all"
		b.WriteString(g.generateFullDiagram())
	}

	b.WriteString("```\n")

	return b.String()
}

// generateFullDiagram generates full architecture diagram
func (g *DiagramGenerator) generateFullDiagram() string {
	var b strings.Builder

	b.WriteString("graph TB\n")
	b.WriteString("    subgraph \"Presentation Layer\"\n")

	for _, domain := range g.domains {
		if domain.HasHandler {
			handlerID := domain.Name + "Handler"
			b.WriteString(fmt.Sprintf("        %s[\"%s Handler\"]\n", handlerID, toPascalCase(domain.Name)))
		}
	}

	b.WriteString("    end\n\n")

	b.WriteString("    subgraph \"Application Layer\"\n")
	for _, domain := range g.domains {
		if domain.HasService {
			serviceID := domain.Name + "Service"
			b.WriteString(fmt.Sprintf("        %s[\"%s Service\"]\n", serviceID, toPascalCase(domain.Name)))
		}
	}
	b.WriteString("    end\n\n")

	b.WriteString("    subgraph \"Domain Layer\"\n")
	for _, domain := range g.domains {
		if domain.HasEntity {
			entityID := domain.Name + "Entity"
			b.WriteString(fmt.Sprintf("        %s[\"%s Entity\"]\n", entityID, toPascalCase(domain.Name)))
		}
	}
	b.WriteString("    end\n\n")

	b.WriteString("    subgraph \"Infrastructure Layer\"\n")
	for _, domain := range g.domains {
		if domain.HasRepo {
			repoID := domain.Name + "Repo"
			b.WriteString(fmt.Sprintf("        %s[\"%s Repository\"]\n", repoID, toPascalCase(domain.Name)))
		}
	}
	b.WriteString("    DB[(Database)]\n")
	b.WriteString("    end\n\n")

	// Add connections
	for _, domain := range g.domains {
		if domain.HasHandler && domain.HasService {
			b.WriteString(fmt.Sprintf("    %sHandler --> %sService\n", domain.Name, domain.Name))
		}
		if domain.HasService && domain.HasEntity {
			b.WriteString(fmt.Sprintf("    %sService --> %sEntity\n", domain.Name, domain.Name))
		}
		if domain.HasService && domain.HasRepo {
			b.WriteString(fmt.Sprintf("    %sService --> %sRepo\n", domain.Name, domain.Name))
		}
		if domain.HasRepo {
			b.WriteString(fmt.Sprintf("    %sRepo --> DB\n", domain.Name))
		}
	}

	// Styling
	b.WriteString("\n    classDef handler fill:#667eea,stroke:#764ba2,color:#fff\n")
	b.WriteString("    classDef service fill:#10b981,stroke:#059669,color:#fff\n")
	b.WriteString("    classDef entity fill:#f59e0b,stroke:#d97706,color:#fff\n")
	b.WriteString("    classDef repo fill:#3b82f6,stroke:#2563eb,color:#fff\n")

	for _, domain := range g.domains {
		if domain.HasHandler {
			b.WriteString(fmt.Sprintf("    class %sHandler handler\n", domain.Name))
		}
		if domain.HasService {
			b.WriteString(fmt.Sprintf("    class %sService service\n", domain.Name))
		}
		if domain.HasEntity {
			b.WriteString(fmt.Sprintf("    class %sEntity entity\n", domain.Name))
		}
		if domain.HasRepo {
			b.WriteString(fmt.Sprintf("    class %sRepo repo\n", domain.Name))
		}
	}

	return b.String()
}

// generateLayersDiagram generates layers diagram
func (g *DiagramGenerator) generateLayersDiagram() string {
	var b strings.Builder

	b.WriteString("graph TB\n")
	b.WriteString("    A[Presentation Layer<br/>HTTP Handlers] --> B[Application Layer<br/>Services]\n")
	b.WriteString("    B --> C[Domain Layer<br/>Entities & Ports]\n")
	b.WriteString("    B --> D[Infrastructure Layer<br/>Repositories]\n")
	b.WriteString("    D --> E[(Database)]\n")
	b.WriteString("\n    classDef layer fill:#667eea,stroke:#764ba2,color:#fff\n")
	b.WriteString("    class A,B,C,D layer\n")

	return b.String()
}

// generateDomainDiagram generates domain diagram
func (g *DiagramGenerator) generateDomainDiagram() string {
	var b strings.Builder

	b.WriteString("graph LR\n")

	for _, domain := range g.domains {
		domainID := domain.Name
		b.WriteString(fmt.Sprintf("    %s[\"%s Domain\"]\n", domainID, toPascalCase(domain.Name)))

		if domain.HasEntity {
			b.WriteString(fmt.Sprintf("    %s --> %sE[Entity]\n", domainID, domainID))
		}
		if domain.HasRepo {
			b.WriteString(fmt.Sprintf("    %s --> %sR[Repository]\n", domainID, domainID))
		}
		if domain.HasHandler {
			b.WriteString(fmt.Sprintf("    %s --> %sH[Handler]\n", domainID, domainID))
		}
		if domain.HasService {
			b.WriteString(fmt.Sprintf("    %s --> %sS[Service]\n", domainID, domainID))
		}
	}

	return b.String()
}

// generateDependencyDiagram generates dependency diagram
func (g *DiagramGenerator) generateDependencyDiagram() string {
	var b strings.Builder

	b.WriteString("graph LR\n")
	b.WriteString("    CLI[Anaphase CLI] --> GEN[Generators]\n")
	b.WriteString("    GEN --> AI[AI Provider<br/>Google Gemini]\n")
	b.WriteString("    GEN --> TMPL[Templates]\n")
	b.WriteString("    CLI --> UI[Bubble Tea UI]\n")
	b.WriteString("    CLI --> CMD[Commands]\n")
	b.WriteString("    CMD --> GEN\n")

	return b.String()
}

// generateASCII generates ASCII art diagram
func (g *DiagramGenerator) generateASCII() string {
	var b strings.Builder

	b.WriteString("┌─────────────────────────────────────────────┐\n")
	b.WriteString("│        CLEAN ARCHITECTURE LAYERS            │\n")
	b.WriteString("├─────────────────────────────────────────────┤\n")
	b.WriteString("│                                             │\n")
	b.WriteString("│  ┌───────────────────────────────────────┐  │\n")
	b.WriteString("│  │   Presentation (HTTP Handlers)        │  │\n")
	b.WriteString("│  └───────────────────────────────────────┘  │\n")
	b.WriteString("│                    ▼                        │\n")
	b.WriteString("│  ┌───────────────────────────────────────┐  │\n")
	b.WriteString("│  │   Application (Services)              │  │\n")
	b.WriteString("│  └───────────────────────────────────────┘  │\n")
	b.WriteString("│                    ▼                        │\n")
	b.WriteString("│  ┌───────────────────────────────────────┐  │\n")
	b.WriteString("│  │   Domain (Entities & Ports)           │  │\n")
	b.WriteString("│  └───────────────────────────────────────┘  │\n")
	b.WriteString("│                    ▼                        │\n")
	b.WriteString("│  ┌───────────────────────────────────────┐  │\n")
	b.WriteString("│  │   Infrastructure (Repositories)       │  │\n")
	b.WriteString("│  └───────────────────────────────────────┘  │\n")
	b.WriteString("│                    ▼                        │\n")
	b.WriteString("│             ┌──────────┐                    │\n")
	b.WriteString("│             │ Database │                    │\n")
	b.WriteString("│             └──────────┘                    │\n")
	b.WriteString("│                                             │\n")
	b.WriteString("└─────────────────────────────────────────────┘\n\n")

	// List domains
	if len(g.domains) > 0 {
		b.WriteString("Discovered Domains:\n")
		for _, domain := range g.domains {
			components := []string{}
			if domain.HasEntity {
				components = append(components, "Entity")
			}
			if domain.HasRepo {
				components = append(components, "Repository")
			}
			if domain.HasHandler {
				components = append(components, "Handler")
			}
			if domain.HasService {
				components = append(components, "Service")
			}

			b.WriteString(fmt.Sprintf("  • %s: [%s]\n", toPascalCase(domain.Name), strings.Join(components, ", ")))
		}
	}

	return b.String()
}
