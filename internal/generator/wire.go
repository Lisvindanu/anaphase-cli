package generator

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// WireConfig holds configuration for wire generation
type WireConfig struct {
	OutputDir string
	Logger    *slog.Logger
}

// WireGenerator generates dependency wiring code
type WireGenerator struct {
	config  *WireConfig
	domains []string
}

// NewWireGenerator creates a new wire generator
func NewWireGenerator(config *WireConfig) *WireGenerator {
	return &WireGenerator{
		config:  config,
		domains: []string{},
	}
}

// Generate creates wiring files
func (g *WireGenerator) Generate(ctx context.Context) ([]string, error) {
	var generatedFiles []string

	// Scan for existing domains
	if err := g.scanDomains(); err != nil {
		return nil, fmt.Errorf("scan domains: %w", err)
	}

	g.config.Logger.Info("discovered domains", "count", len(g.domains), "domains", g.domains)

	// Generate main.go
	if err := os.MkdirAll(g.config.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	mainFile, err := g.generateMain()
	if err != nil {
		return nil, fmt.Errorf("generate main: %w", err)
	}
	generatedFiles = append(generatedFiles, mainFile)

	// Generate wire.go (dependency injection)
	wireFile, err := g.generateWire()
	if err != nil {
		return nil, fmt.Errorf("generate wire: %w", err)
	}
	generatedFiles = append(generatedFiles, wireFile)

	return generatedFiles, nil
}

// scanDomains discovers existing domains by scanning entity files
func (g *WireGenerator) scanDomains() error {
	entityDir := "internal/core/entity"

	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		g.config.Logger.Warn("entity directory not found", "path", entityDir)
		return nil
	}

	fset := token.NewFileSet()

	err := filepath.WalkDir(entityDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Parse the file
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			g.config.Logger.Warn("failed to parse file", "path", path, "error", err)
			return nil
		}

		// Look for struct declarations that might be entities
		ast.Inspect(file, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			// Check if it's a struct
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				domainName := strings.ToLower(typeSpec.Name.Name)

				// Avoid duplicates
				found := false
				for _, existing := range g.domains {
					if existing == domainName {
						found = true
						break
					}
				}

				if !found {
					g.domains = append(g.domains, domainName)
				}
			}

			return true
		})

		return nil
	})

	return err
}

func (g *WireGenerator) generateMain() (string, error) {
	filename := filepath.Join(g.config.OutputDir, "main.go")

	var b strings.Builder

	b.WriteString("package main\n\n")

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n")
	b.WriteString("\t\"log/slog\"\n")
	b.WriteString("\t\"net/http\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"os/signal\"\n")
	b.WriteString("\t\"syscall\"\n")
	b.WriteString("\t\"time\"\n\n")
	b.WriteString("\t\"github.com/go-chi/chi/v5\"\n")
	b.WriteString("\t\"github.com/go-chi/chi/v5/middleware\"\n")
	b.WriteString(")\n\n")

	// Main function
	b.WriteString("func main() {\n")
	b.WriteString("\t// Setup logger\n")
	b.WriteString("\tlogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{\n")
	b.WriteString("\t\tLevel: slog.LevelInfo,\n")
	b.WriteString("\t}))\n")
	b.WriteString("\tslog.SetDefault(logger)\n\n")

	b.WriteString("\t// Create context with cancellation\n")
	b.WriteString("\tctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)\n")
	b.WriteString("\tdefer cancel()\n\n")

	b.WriteString("\t// Initialize dependencies\n")
	b.WriteString("\tapp, err := InitializeApp(logger)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\tlogger.Error(\"failed to initialize app\", \"error\", err)\n")
	b.WriteString("\t\tos.Exit(1)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tdefer app.Cleanup()\n\n")

	b.WriteString("\t// Setup router\n")
	b.WriteString("\tr := chi.NewRouter()\n")
	b.WriteString("\tr.Use(middleware.Logger)\n")
	b.WriteString("\tr.Use(middleware.Recoverer)\n")
	b.WriteString("\tr.Use(middleware.RequestID)\n")
	b.WriteString("\tr.Use(middleware.Timeout(60 * time.Second))\n\n")

	b.WriteString("\t// Health check\n")
	b.WriteString("\tr.Get(\"/health\", func(w http.ResponseWriter, r *http.Request) {\n")
	b.WriteString("\t\tw.WriteHeader(http.StatusOK)\n")
	b.WriteString("\t\tw.Write([]byte(\"OK\"))\n")
	b.WriteString("\t})\n\n")

	b.WriteString("\t// API routes\n")
	b.WriteString("\tr.Route(\"/api/v1\", func(r chi.Router) {\n")
	b.WriteString("\t\tapp.RegisterRoutes(r)\n")
	b.WriteString("\t})\n\n")

	b.WriteString("\t// Start server\n")
	b.WriteString("\tport := os.Getenv(\"PORT\")\n")
	b.WriteString("\tif port == \"\" {\n")
	b.WriteString("\t\tport = \"8080\"\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tsrv := &http.Server{\n")
	b.WriteString("\t\tAddr:    \":\" + port,\n")
	b.WriteString("\t\tHandler: r,\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\t// Start server in goroutine\n")
	b.WriteString("\tgo func() {\n")
	b.WriteString("\t\tlogger.Info(\"starting server\", \"port\", port)\n")
	b.WriteString("\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {\n")
	b.WriteString("\t\t\tlogger.Error(\"server error\", \"error\", err)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}()\n\n")

	b.WriteString("\t// Wait for interrupt signal\n")
	b.WriteString("\t<-ctx.Done()\n")
	b.WriteString("\tlogger.Info(\"shutting down gracefully...\")\n\n")

	b.WriteString("\t// Graceful shutdown\n")
	b.WriteString("\tshutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)\n")
	b.WriteString("\tdefer shutdownCancel()\n\n")

	b.WriteString("\tif err := srv.Shutdown(shutdownCtx); err != nil {\n")
	b.WriteString("\t\tlogger.Error(\"shutdown error\", \"error\", err)\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tlogger.Info(\"server stopped\")\n")
	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *WireGenerator) generateWire() (string, error) {
	filename := filepath.Join(g.config.OutputDir, "wire.go")

	var b strings.Builder

	b.WriteString("package main\n\n")

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"log/slog\"\n")
	b.WriteString("\t\"os\"\n\n")
	b.WriteString("\t\"github.com/go-chi/chi/v5\"\n")
	b.WriteString("\t\"github.com/jackc/pgx/v5/pgxpool\"\n\n")

	// Import handlers and repositories for each domain
	if len(g.domains) > 0 {
		b.WriteString("\thandlerhttp \"github.com/lisvindanuu/anaphase-cli/internal/adapter/handler/http\"\n")
		b.WriteString("\t\"github.com/lisvindanuu/anaphase-cli/internal/adapter/repository/postgres\"\n")
	}

	b.WriteString(")\n\n")

	// App struct
	b.WriteString("// App holds all application dependencies\n")
	b.WriteString("type App struct {\n")
	b.WriteString("\tlogger *slog.Logger\n")
	b.WriteString("\tdb     *pgxpool.Pool\n\n")

	// Handlers for each domain
	for _, domain := range g.domains {
		entityName := toPascalCase(domain)
		b.WriteString(fmt.Sprintf("\t%sHandler *handlerhttp.%sHandler\n", domain, entityName))
	}

	b.WriteString("}\n\n")

	// InitializeApp function
	b.WriteString("// InitializeApp initializes all application dependencies\n")
	b.WriteString("func InitializeApp(logger *slog.Logger) (*App, error) {\n")
	b.WriteString("\t// Database connection\n")
	b.WriteString("\tdbURL := os.Getenv(\"DATABASE_URL\")\n")
	b.WriteString("\tif dbURL == \"\" {\n")
	b.WriteString("\t\tdbURL = \"postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable\"\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tdb, err := pgxpool.New(context.Background(), dbURL)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn nil, fmt.Errorf(\"connect to database: %w\", err)\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\t// Ping database\n")
	b.WriteString("\tif err := db.Ping(context.Background()); err != nil {\n")
	b.WriteString("\t\treturn nil, fmt.Errorf(\"ping database: %w\", err)\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tlogger.Info(\"database connected\")\n\n")

	// Initialize repositories and handlers for each domain
	for _, domain := range g.domains {
		entityName := toPascalCase(domain)

		b.WriteString(fmt.Sprintf("\t// Initialize %s dependencies\n", domain))
		b.WriteString(fmt.Sprintf("\t%sRepo := postgres.New%sRepository(db)\n", domain, entityName))
		b.WriteString(fmt.Sprintf("\t// TODO: Create %s service implementation\n", domain))
		b.WriteString(fmt.Sprintf("\t// %sService := service.New%sService(%sRepo)\n", domain, entityName, domain))
		b.WriteString(fmt.Sprintf("\t%sHandler := handlerhttp.New%sHandler(nil, logger) // Pass service when implemented\n\n", domain, entityName))
	}

	b.WriteString("\treturn &App{\n")
	b.WriteString("\t\tlogger: logger,\n")
	b.WriteString("\t\tdb:     db,\n")

	for _, domain := range g.domains {
		b.WriteString(fmt.Sprintf("\t\t%sHandler: %sHandler,\n", domain, domain))
	}

	b.WriteString("\t}, nil\n")
	b.WriteString("}\n\n")

	// RegisterRoutes method
	b.WriteString("// RegisterRoutes registers all HTTP routes\n")
	b.WriteString("func (a *App) RegisterRoutes(r chi.Router) {\n")

	for _, domain := range g.domains {
		b.WriteString(fmt.Sprintf("\ta.%sHandler.RegisterRoutes(r)\n", domain))
	}

	b.WriteString("}\n\n")

	// Cleanup method
	b.WriteString("// Cleanup cleans up application resources\n")
	b.WriteString("func (a *App) Cleanup() {\n")
	b.WriteString("\tif a.db != nil {\n")
	b.WriteString("\t\ta.db.Close()\n")
	b.WriteString("\t\ta.logger.Info(\"database connection closed\")\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}
