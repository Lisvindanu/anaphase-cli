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
	config     *WireConfig
	domains    []string
	dbType     string // postgres, mysql, or sqlite
	dbDriver   string // pgxpool, sql.DB, etc
	moduleName string // detected from go.mod
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

	// Detect module name from go.mod
	if err := g.detectModuleName(); err != nil {
		return nil, fmt.Errorf("detect module name: %w", err)
	}

	// Detect database type from .env
	if err := g.detectDatabaseType(); err != nil {
		g.config.Logger.Warn("failed to detect database type, defaulting to postgres", "error", err)
		g.dbType = "postgres"
		g.dbDriver = "pgxpool"
	}

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

// detectModuleName reads go.mod and extracts module name
func (g *WireGenerator) detectModuleName() error {
	goModFile := "go.mod"
	data, err := os.ReadFile(goModFile)
	if err != nil {
		return fmt.Errorf("read go.mod: %w", err)
	}

	// Parse module line
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			g.moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module"))
			g.config.Logger.Info("detected module name", "module", g.moduleName)
			return nil
		}
	}

	return fmt.Errorf("module name not found in go.mod")
}

// detectDatabaseType reads .env file and detects database type
func (g *WireGenerator) detectDatabaseType() error {
	envFile := ".env"
	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("read .env: %w", err)
	}

	envContent := string(data)

	// Parse DATABASE_URL
	if strings.Contains(envContent, "DATABASE_URL=") {
		for _, line := range strings.Split(envContent, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "DATABASE_URL=") {
				dbURL := strings.TrimPrefix(line, "DATABASE_URL=")

				if strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://") {
					g.dbType = "postgres"
					g.dbDriver = "pgxpool"
				} else if strings.HasPrefix(dbURL, "mysql://") {
					g.dbType = "mysql"
					g.dbDriver = "sql.DB"
				} else if strings.HasPrefix(dbURL, "sqlite://") || strings.Contains(dbURL, ".db") {
					g.dbType = "sqlite"
					g.dbDriver = "sql.DB"
				}

				g.config.Logger.Info("detected database type", "type", g.dbType, "driver", g.dbDriver)
				return nil
			}
		}
	}

	return fmt.Errorf("DATABASE_URL not found in .env")
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

	// Context only needed for PostgreSQL
	if g.dbType == "postgres" {
		b.WriteString("\t\"context\"\n")
	}

	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"log/slog\"\n")
	b.WriteString("\t\"os\"\n")

	// Strings package needed for MySQL/SQLite DSN parsing
	if g.dbType == "mysql" || g.dbType == "sqlite" {
		b.WriteString("\t\"strings\"\n")
	}

	b.WriteString("\n\t\"github.com/go-chi/chi/v5\"\n")

	
	// Database driver import based on detected type
	switch g.dbType {
	case "postgres":
		b.WriteString("\t\"github.com/jackc/pgx/v5/pgxpool\"\n\n")
	case "mysql", "sqlite":
		b.WriteString("\t\"database/sql\"\n")
		if g.dbType == "mysql" {
			b.WriteString("\t_ \"github.com/go-sql-driver/mysql\"\n\n")
		} else {
			b.WriteString("\t_ \"github.com/mattn/go-sqlite3\"\n\n")
		}
	}

	// Import handlers and repositories for each domain
	if len(g.domains) > 0 {
		b.WriteString(fmt.Sprintf("\thandlerhttp \"%s/internal/adapter/handler/http\"\n", g.moduleName))
		// Repository import based on database type
		switch g.dbType {
		case "postgres":
			b.WriteString(fmt.Sprintf("\t\"%s/internal/adapter/repository/postgres\"\n", g.moduleName))
		case "mysql":
			b.WriteString(fmt.Sprintf("\t\"%s/internal/adapter/repository/mysql\"\n", g.moduleName))
		case "sqlite":
			b.WriteString(fmt.Sprintf("\t\"%s/internal/adapter/repository/sqlite\"\n", g.moduleName))
		}
	}

	b.WriteString(")\n\n")

	// App struct
	b.WriteString("// App holds all application dependencies\n")
	b.WriteString("type App struct {\n")
	b.WriteString("\tlogger *slog.Logger\n")

	// Database field type based on driver
	switch g.dbType {
	case "postgres":
		b.WriteString("\tdb     *pgxpool.Pool\n\n")
	case "mysql", "sqlite":
		b.WriteString("\tdb     *sql.DB\n\n")
	}

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

	// Default database URL based on type
	switch g.dbType {
	case "postgres":
		b.WriteString("\t\tdbURL = \"postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable\"\n")
	case "mysql":
		b.WriteString("\t\tdbURL = \"mysql://root:password@localhost:3306/anaphase?parseTime=true\"\n")
	case "sqlite":
		b.WriteString("\t\tdbURL = \"sqlite://./anaphase.db\"\n")
	}

	b.WriteString("\t}\n\n")

	// Database connection code based on type
	switch g.dbType {
	case "postgres":
		b.WriteString("\tdb, err := pgxpool.New(context.Background(), dbURL)\n")
		b.WriteString("\tif err != nil {\n")
		b.WriteString("\t\treturn nil, fmt.Errorf(\"connect to database: %w\", err)\n")
		b.WriteString("\t}\n\n")
		b.WriteString("\t// Ping database\n")
		b.WriteString("\tif err := db.Ping(context.Background()); err != nil {\n")
		b.WriteString("\t\treturn nil, fmt.Errorf(\"ping database: %w\", err)\n")
		b.WriteString("\t}\n\n")
	case "mysql":
		b.WriteString("\t// Parse MySQL DSN from URL (mysql://user:pass@host:port/db?params)\n")
		b.WriteString("\t// Convert to MySQL driver format: user:pass@tcp(host:port)/db?params\n")
		b.WriteString("\tdsn := strings.TrimPrefix(dbURL, \"mysql://\")\n")
		b.WriteString("\t// Replace @host:port with @tcp(host:port)\n")
		b.WriteString("\tif idx := strings.Index(dsn, \"@\"); idx != -1 {\n")
		b.WriteString("\t\trest := dsn[idx+1:]\n")
		b.WriteString("\t\tif dbIdx := strings.Index(rest, \"/\"); dbIdx != -1 {\n")
		b.WriteString("\t\t\thost := rest[:dbIdx]\n")
		b.WriteString("\t\t\tdsn = dsn[:idx+1] + \"tcp(\" + host + \")\" + rest[dbIdx:]\n")
		b.WriteString("\t\t}\n")
		b.WriteString("\t}\n\n")
		b.WriteString("\tdb, err := sql.Open(\"mysql\", dsn)\n")
		b.WriteString("\tif err != nil {\n")
		b.WriteString("\t\treturn nil, fmt.Errorf(\"open database: %w\", err)\n")
		b.WriteString("\t}\n\n")
		b.WriteString("\t// Ping database\n")
		b.WriteString("\tif err := db.Ping(); err != nil {\n")
		b.WriteString("\t\treturn nil, fmt.Errorf(\"ping database: %w\", err)\n")
		b.WriteString("\t}\n\n")
	case "sqlite":
		b.WriteString("\t// Parse SQLite path from URL\n")
		b.WriteString("\tdbPath := strings.TrimPrefix(dbURL, \"sqlite://\")\n")
		b.WriteString("\tdb, err := sql.Open(\"sqlite3\", dbPath)\n")
		b.WriteString("\tif err != nil {\n")
		b.WriteString("\t\treturn nil, fmt.Errorf(\"open database: %w\", err)\n")
		b.WriteString("\t}\n\n")
		b.WriteString("\t// Ping database\n")
		b.WriteString("\tif err := db.Ping(); err != nil {\n")
		b.WriteString("\t\treturn nil, fmt.Errorf(\"ping database: %w\", err)\n")
		b.WriteString("\t}\n\n")
	}

	b.WriteString("\tlogger.Info(\"database connected\")\n\n")

	// Initialize repositories and handlers for each domain
	for _, domain := range g.domains {
		entityName := toPascalCase(domain)

		b.WriteString(fmt.Sprintf("\t// Initialize %s dependencies\n", domain))

		// Repository instantiation based on database type
		switch g.dbType {
		case "postgres":
			b.WriteString(fmt.Sprintf("\t%sRepo := postgres.New%sRepository(db)\n", domain, entityName))
		case "mysql":
			b.WriteString(fmt.Sprintf("\t%sRepo := mysql.New%sRepository(db)\n", domain, entityName))
		case "sqlite":
			b.WriteString(fmt.Sprintf("\t%sRepo := sqlite.New%sRepository(db)\n", domain, entityName))
		}

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
