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
	Entity       *ComponentInfo
	ServicePort  *ComponentInfo
	RepoPort     *ComponentInfo
	RepoAdapter  *ComponentInfo
	Handler      *ComponentInfo
	Dependencies []string
}

// ComponentInfo holds information about a component
type ComponentInfo struct {
	Name      string
	Type      string // "entity", "interface", "struct"
	Path      string
	Package   string
	IsPort    bool // Is this a port/interface?
	IsAdapter bool // Is this an adapter/implementation?
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

// scanProject scans the project structure flexibly
func (g *DiagramGenerator) scanProject() error {
	// Scan all components in internal/
	components, err := g.scanAllComponents()
	if err != nil {
		return err
	}

	// Group by domain
	domainMap := make(map[string]*DomainInfo)

	for _, comp := range components {
		domainName := comp.Name

		// Initialize domain if not exists
		if _, exists := domainMap[domainName]; !exists {
			domainMap[domainName] = &DomainInfo{
				Name: domainName,
			}
		}

		domain := domainMap[domainName]

		// Categorize component based on type and package
		switch {
		case comp.Type == "entity" || strings.Contains(comp.Package, "entity"):
			domain.Entity = &comp
		case comp.IsPort && (strings.Contains(comp.Path, "service") || strings.HasSuffix(comp.Name, "Service")):
			domain.ServicePort = &comp
		case comp.IsPort && (strings.Contains(comp.Path, "repository") || strings.HasSuffix(comp.Name, "Repository")):
			domain.RepoPort = &comp
		case comp.IsAdapter && (strings.Contains(comp.Path, "repository") || strings.Contains(comp.Package, "postgres") || strings.Contains(comp.Package, "mysql") || strings.Contains(comp.Package, "mongo")):
			domain.RepoAdapter = &comp
		case strings.Contains(comp.Package, "handler") || strings.Contains(comp.Path, "handler"):
			domain.Handler = &comp
		}
	}

	// Convert map to slice, filter out empty domains
	for _, domain := range domainMap {
		// Only include domains that have at least one meaningful component
		if domain.Entity != nil || domain.Handler != nil || domain.RepoAdapter != nil ||
			(domain.ServicePort != nil && domain.RepoPort != nil) {
			g.domains = append(g.domains, *domain)
		}
	}

	return nil
}

// scanAllComponents scans all Go files in internal/ and detects components
func (g *DiagramGenerator) scanAllComponents() ([]ComponentInfo, error) {
	var components []ComponentInfo
	internalDir := "internal"

	if _, err := os.Stat(internalDir); os.IsNotExist(err) {
		return components, nil
	}

	fset := token.NewFileSet()

	err := filepath.WalkDir(internalDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip test files, directories, and non-Go files
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Skip internal tool packages (not part of generated architecture)
		skipPaths := []string{
			"/generator/", "/commands/", "/ui/", "/config/",
			"/provider/", "/cache/", "/orchestrator/", "/utils/",
		}
		for _, skipPath := range skipPaths {
			if strings.Contains(path, skipPath) {
				return nil
			}
		}

		// Only process files in core/ and adapter/ directories
		if !strings.Contains(path, "/core/") && !strings.Contains(path, "/adapter/") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files with parse errors
		}

		packageName := file.Name.Name

		ast.Inspect(file, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			typeName := typeSpec.Name.Name

			// Skip DTOs, requests, responses, configs, and other non-domain types
			skipSuffixes := []string{"Request", "Response", "DTO", "Config", "Error", "Result", "Option"}
			for _, suffix := range skipSuffixes {
				if strings.HasSuffix(typeName, suffix) {
					return true
				}
			}

			var comp ComponentInfo
			comp.Path = path
			comp.Package = packageName

			// Detect type: interface or struct
			switch typeSpec.Type.(type) {
			case *ast.InterfaceType:
				comp.Type = "interface"
				comp.IsPort = true
				comp.Name = extractDomainName(typeName)
			case *ast.StructType:
				comp.Type = "struct"
				comp.IsAdapter = !strings.Contains(path, "/entity/") && !strings.Contains(path, "/valueobject/")
				comp.Name = extractDomainName(typeName)

				// Entities are structs in entity package
				if strings.Contains(path, "/entity/") {
					comp.Type = "entity"
					comp.IsAdapter = false
				}

				// Skip value objects - they're not domain components we track
				if strings.Contains(path, "/valueobject/") {
					return true
				}
			}

			// Only add if we have a valid name
			if comp.Name != "" {
				components = append(components, comp)
			}

			return true
		})

		return nil
	})

	return components, err
}

// extractDomainName extracts domain name from type name
// e.g., "CustomerRepository" -> "customer", "OrderService" -> "order"
func extractDomainName(typeName string) string {
	name := typeName

	// Remove common suffixes
	suffixes := []string{"Repository", "Service", "Handler", "Entity", "Port", "Adapter", "Repo"}
	for _, suffix := range suffixes {
		name = strings.TrimSuffix(name, suffix)
	}

	// Convert to lowercase
	return strings.ToLower(name)
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

// generateFullDiagram generates full architecture diagram with ports and adapters
func (g *DiagramGenerator) generateFullDiagram() string {
	var b strings.Builder

	b.WriteString("graph TB\n")

	// Presentation Layer
	b.WriteString("    subgraph \"Presentation Layer\"\n")
	for _, domain := range g.domains {
		if domain.Handler != nil {
			handlerID := domain.Name + "Handler"
			b.WriteString(fmt.Sprintf("        %s[\"%s Handler\"]\n", handlerID, toPascalCase(domain.Name)))
		}
	}
	b.WriteString("    end\n\n")

	// Application Layer (Service Ports/Interfaces)
	hasAppLayer := false
	for _, domain := range g.domains {
		if domain.ServicePort != nil {
			hasAppLayer = true
			break
		}
	}
	if hasAppLayer {
		b.WriteString("    subgraph \"Application Layer\"\n")
		for _, domain := range g.domains {
			if domain.ServicePort != nil {
				serviceID := domain.Name + "Service"
				b.WriteString(fmt.Sprintf("        %s[\"%s Service<br/>(Port/Interface)\"]\n", serviceID, toPascalCase(domain.Name)))
			}
		}
		b.WriteString("    end\n\n")
	}

	// Domain Layer (Entities and Repository Ports)
	b.WriteString("    subgraph \"Domain Layer\"\n")
	for _, domain := range g.domains {
		if domain.Entity != nil {
			entityID := domain.Name + "Entity"
			b.WriteString(fmt.Sprintf("        %s[\"%s Entity\"]\n", entityID, toPascalCase(domain.Name)))
		}
		if domain.RepoPort != nil {
			repoPortID := domain.Name + "RepoPort"
			b.WriteString(fmt.Sprintf("        %s[\"%s Repository<br/>(Port/Interface)\"]\n", repoPortID, toPascalCase(domain.Name)))
		}
	}
	b.WriteString("    end\n\n")

	// Infrastructure Layer (Repository Adapters)
	b.WriteString("    subgraph \"Infrastructure Layer\"\n")
	for _, domain := range g.domains {
		if domain.RepoAdapter != nil {
			repoID := domain.Name + "Repo"
			adapterName := toPascalCase(domain.RepoAdapter.Package) + " " + toPascalCase(domain.Name) + " Repo"
			b.WriteString(fmt.Sprintf("        %s[\"%s<br/>(Adapter)\"]\n", repoID, adapterName))
		}
	}
	b.WriteString("    DB[(Database)]\n")
	b.WriteString("    end\n\n")

	// Add connections (following Clean Architecture dependency rule)
	for _, domain := range g.domains {
		handlerID := domain.Name + "Handler"
		serviceID := domain.Name + "Service"
		entityID := domain.Name + "Entity"
		repoPortID := domain.Name + "RepoPort"
		repoID := domain.Name + "Repo"

		// Handler -> Service Port
		if domain.Handler != nil && domain.ServicePort != nil {
			b.WriteString(fmt.Sprintf("    %s -->|calls| %s\n", handlerID, serviceID))
		}

		// Service Port -> Entity
		if domain.ServicePort != nil && domain.Entity != nil {
			b.WriteString(fmt.Sprintf("    %s -->|uses| %s\n", serviceID, entityID))
		}

		// Service Port -> Repository Port
		if domain.ServicePort != nil && domain.RepoPort != nil {
			b.WriteString(fmt.Sprintf("    %s -->|depends on| %s\n", serviceID, repoPortID))
		}

		// Handler -> Repository Port (if no service)
		if domain.Handler != nil && domain.ServicePort == nil && domain.RepoPort != nil {
			b.WriteString(fmt.Sprintf("    %s -->|uses| %s\n", handlerID, repoPortID))
		}

		// Repository Adapter implements Repository Port
		if domain.RepoAdapter != nil && domain.RepoPort != nil {
			b.WriteString(fmt.Sprintf("    %s -.->|implements| %s\n", repoID, repoPortID))
		}

		// Repository Adapter -> Database
		if domain.RepoAdapter != nil {
			b.WriteString(fmt.Sprintf("    %s --> DB\n", repoID))
		}
	}

	// Styling
	b.WriteString("\n    classDef handler fill:#667eea,stroke:#764ba2,color:#fff\n")
	b.WriteString("    classDef service fill:#10b981,stroke:#059669,color:#fff\n")
	b.WriteString("    classDef entity fill:#f59e0b,stroke:#d97706,color:#fff\n")
	b.WriteString("    classDef port fill:#ec4899,stroke:#db2777,color:#fff,stroke-dasharray: 5 5\n")
	b.WriteString("    classDef repo fill:#3b82f6,stroke:#2563eb,color:#fff\n")

	for _, domain := range g.domains {
		if domain.Handler != nil {
			b.WriteString(fmt.Sprintf("    class %sHandler handler\n", domain.Name))
		}
		if domain.ServicePort != nil {
			b.WriteString(fmt.Sprintf("    class %sService port\n", domain.Name))
		}
		if domain.Entity != nil {
			b.WriteString(fmt.Sprintf("    class %sEntity entity\n", domain.Name))
		}
		if domain.RepoPort != nil {
			b.WriteString(fmt.Sprintf("    class %sRepoPort port\n", domain.Name))
		}
		if domain.RepoAdapter != nil {
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

		if domain.Entity != nil {
			b.WriteString(fmt.Sprintf("    %s --> %sE[Entity]\n", domainID, domainID))
		}
		if domain.ServicePort != nil {
			b.WriteString(fmt.Sprintf("    %s --> %sSP[Service Port]\n", domainID, domainID))
		}
		if domain.RepoPort != nil {
			b.WriteString(fmt.Sprintf("    %s --> %sRP[Repo Port]\n", domainID, domainID))
		}
		if domain.RepoAdapter != nil {
			b.WriteString(fmt.Sprintf("    %s --> %sR[Repo Adapter]\n", domainID, domainID))
		}
		if domain.Handler != nil {
			b.WriteString(fmt.Sprintf("    %s --> %sH[Handler]\n", domainID, domainID))
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

	// List domains with detailed component info
	if len(g.domains) > 0 {
		b.WriteString("Discovered Domains:\n")
		for _, domain := range g.domains {
			components := []string{}
			if domain.Entity != nil {
				components = append(components, "Entity")
			}
			if domain.ServicePort != nil {
				components = append(components, "Service Port")
			}
			if domain.RepoPort != nil {
				components = append(components, "Repository Port")
			}
			if domain.RepoAdapter != nil {
				adapterType := toPascalCase(domain.RepoAdapter.Package)
				components = append(components, adapterType+" Repository")
			}
			if domain.Handler != nil {
				components = append(components, "Handler")
			}

			b.WriteString(fmt.Sprintf("  • %s: [%s]\n", toPascalCase(domain.Name), strings.Join(components, ", ")))
		}

		// Legend
		b.WriteString("\nLegend:\n")
		b.WriteString("  Port = Interface (defines contract)\n")
		b.WriteString("  Adapter = Implementation (concrete)\n")
	}

	return b.String()
}
