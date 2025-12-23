package generator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// RepositoryConfig holds configuration for repository generation
type RepositoryConfig struct {
	Database string
	Cache    bool
	Logger   *slog.Logger
}

// RepositoryGenerator generates repository implementations
type RepositoryGenerator struct {
	domainName string
	config     *RepositoryConfig
	moduleName string
}

// NewRepositoryGenerator creates a new repository generator
func NewRepositoryGenerator(domainName string, config *RepositoryConfig) *RepositoryGenerator {
	return &RepositoryGenerator{
		domainName: domainName,
		config:     config,
	}
}

// Generate creates repository files
func (g *RepositoryGenerator) Generate(ctx context.Context) ([]string, error) {
	var generatedFiles []string

	// Detect module name from go.mod
	if err := g.detectModuleName(); err != nil {
		return nil, fmt.Errorf("detect module name: %w", err)
	}

	// Detect database type from .env if not specified
	if g.config.Database == "" || g.config.Database == "postgres" {
		if err := g.detectDatabaseType(); err != nil {
			g.config.Logger.Warn("failed to detect database type, using postgres", "error", err)
			g.config.Database = "postgres"
		}
	}

	// Create output directory
	outputDir := filepath.Join("internal", "adapter", "repository", g.config.Database)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	// Generate repository implementation
	repoFile, err := g.generateRepository(outputDir)
	if err != nil {
		return nil, fmt.Errorf("generate repository: %w", err)
	}
	generatedFiles = append(generatedFiles, repoFile)

	// Generate SQL queries
	if g.config.Database == "postgres" || g.config.Database == "mysql" {
		sqlFile, err := g.generateSQL(outputDir)
		if err != nil {
			return nil, fmt.Errorf("generate SQL: %w", err)
		}
		generatedFiles = append(generatedFiles, sqlFile)
	}

	// Generate test file
	testFile, err := g.generateRepositoryTest(outputDir)
	if err != nil {
		return nil, fmt.Errorf("generate repository test: %w", err)
	}
	generatedFiles = append(generatedFiles, testFile)

	return generatedFiles, nil
}

func (g *RepositoryGenerator) generateRepository(outputDir string) (string, error) {
	filename := filepath.Join(outputDir, g.domainName+"_repo.go")

	var b strings.Builder

	// Package
	b.WriteString(fmt.Sprintf("package %s\n\n", g.config.Database))

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"context\"\n")
	b.WriteString("\t\"fmt\"\n\n")
	b.WriteString("\t\"github.com/google/uuid\"\n")

	switch g.config.Database {
	case "postgres":
		b.WriteString("\t\"github.com/jackc/pgx/v5\"\n")
		b.WriteString("\t\"github.com/jackc/pgx/v5/pgxpool\"\n")
	case "mysql":
		b.WriteString("\t\"database/sql\"\n")
		b.WriteString("\t_ \"github.com/go-sql-driver/mysql\"\n")
	case "mongodb":
		b.WriteString("\t\"go.mongodb.org/mongo-driver/mongo\"\n")
		b.WriteString("\t\"go.mongodb.org/mongo-driver/bson\"\n")
	}

	b.WriteString(fmt.Sprintf("\n\t\"%s/internal/core/entity\"\n", g.moduleName))
	b.WriteString(fmt.Sprintf("\t\"%s/internal/core/port\"\n", g.moduleName))
	b.WriteString(")\n\n")

	// Repository struct
	entityName := toPascalCase(g.domainName)
	structName := strings.ToLower(g.domainName) + "Repository"

	b.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	switch g.config.Database {
	case "postgres":
		b.WriteString("\tdb *pgxpool.Pool\n")
	case "mysql":
		b.WriteString("\tdb *sql.DB\n")
	case "mongodb":
		b.WriteString("\tcollection *mongo.Collection\n")
	}
	b.WriteString("}\n\n")

	// Constructor
	b.WriteString(fmt.Sprintf("// New%sRepository creates a new %s repository\n", entityName, g.domainName))
	switch g.config.Database {
	case "postgres":
		b.WriteString(fmt.Sprintf("func New%sRepository(db *pgxpool.Pool) port.%sRepository {\n", entityName, entityName))
	case "mysql":
		b.WriteString(fmt.Sprintf("func New%sRepository(db *sql.DB) port.%sRepository {\n", entityName, entityName))
	case "mongodb":
		b.WriteString(fmt.Sprintf("func New%sRepository(collection *mongo.Collection) port.%sRepository {\n", entityName, entityName))
	}
	b.WriteString(fmt.Sprintf("\treturn &%s{\n", structName))
	switch g.config.Database {
	case "postgres", "mysql":
		b.WriteString("\t\tdb: db,\n")
	case "mongodb":
		b.WriteString("\t\tcollection: collection,\n")
	}
	b.WriteString("\t}\n")
	b.WriteString("}\n\n")

	// Save method
	b.WriteString(fmt.Sprintf("// Save saves a %s to the repository\n", g.domainName))
	b.WriteString(fmt.Sprintf("func (r *%s) Save(ctx context.Context, entity *entity.%s) error {\n", structName, entityName))

	switch g.config.Database {
	case "postgres":
		g.generatePostgresSave(&b, entityName)
	case "mysql":
		g.generateMySQLSave(&b, entityName)
	case "mongodb":
		g.generateMongoDBSave(&b, entityName)
	}

	b.WriteString("}\n\n")

	// FindByID method
	b.WriteString(fmt.Sprintf("// FindByID retrieves a %s by ID\n", g.domainName))
	b.WriteString(fmt.Sprintf("func (r *%s) FindByID(ctx context.Context, id uuid.UUID) (*entity.%s, error) {\n", structName, entityName))

	switch g.config.Database {
	case "postgres":
		g.generatePostgresFindByID(&b, entityName)
	case "mysql":
		g.generateMySQLFindByID(&b, entityName)
	case "mongodb":
		g.generateMongoDBFindByID(&b, entityName)
	}

	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *RepositoryGenerator) generatePostgresSave(b *strings.Builder, entityName string) {
	tableName := toSnakeCase(entityName) + "s"

	b.WriteString("\tquery := `\n")
	fmt.Fprintf(b, "\t\tINSERT INTO %s (id, created_at, updated_at)\n", tableName)
	b.WriteString("\t\tVALUES ($1, $2, $3)\n")
	b.WriteString("\t\tON CONFLICT (id) DO UPDATE\n")
	b.WriteString("\t\tSET updated_at = $3\n")
	b.WriteString("\t`\n\n")

	b.WriteString("\t_, err := r.db.Exec(ctx, query,\n")
	b.WriteString("\t\tentity.ID,\n")
	b.WriteString("\t\tentity.CreatedAt,\n")
	b.WriteString("\t\tentity.UpdatedAt,\n")
	b.WriteString("\t)\n\n")

	b.WriteString("\tif err != nil {\n")
	fmt.Fprintf(b, "\t\treturn fmt.Errorf(\"save %s: %%w\", err)\n", g.domainName)
	b.WriteString("\t}\n\n")

	b.WriteString("\treturn nil\n")
}

func (g *RepositoryGenerator) generatePostgresFindByID(b *strings.Builder, entityName string) {
	tableName := toSnakeCase(entityName) + "s"

	fmt.Fprintf(b, "\tvar result entity.%s\n\n", entityName)

	b.WriteString("\tquery := `\n")
	b.WriteString("\t\tSELECT id, created_at, updated_at\n")
	fmt.Fprintf(b, "\t\tFROM %s\n", tableName)
	b.WriteString("\t\tWHERE id = $1\n")
	b.WriteString("\t`\n\n")

	b.WriteString("\terr := r.db.QueryRow(ctx, query, id).Scan(\n")
	b.WriteString("\t\t&result.ID,\n")
	b.WriteString("\t\t&result.CreatedAt,\n")
	b.WriteString("\t\t&result.UpdatedAt,\n")
	b.WriteString("\t)\n\n")

	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\tif err == pgx.ErrNoRows {\n")
	fmt.Fprintf(b, "\t\t\treturn nil, fmt.Errorf(\"%s not found\")\n", g.domainName)
	b.WriteString("\t\t}\n")
	fmt.Fprintf(b, "\t\treturn nil, fmt.Errorf(\"find %s: %%w\", err)\n", g.domainName)
	b.WriteString("\t}\n\n")

	b.WriteString("\treturn &result, nil\n")
}

func (g *RepositoryGenerator) generateMySQLSave(b *strings.Builder, entityName string) {
	b.WriteString("\t// TODO: Implement MySQL save\n")
	b.WriteString("\treturn fmt.Errorf(\"not implemented\")\n")
}

func (g *RepositoryGenerator) generateMySQLFindByID(b *strings.Builder, entityName string) {
	b.WriteString("\t// TODO: Implement MySQL find by ID\n")
	b.WriteString("\treturn nil, fmt.Errorf(\"not implemented\")\n")
}

func (g *RepositoryGenerator) generateMongoDBSave(b *strings.Builder, entityName string) {
	b.WriteString("\t// TODO: Implement MongoDB save\n")
	b.WriteString("\treturn fmt.Errorf(\"not implemented\")\n")
}

func (g *RepositoryGenerator) generateMongoDBFindByID(b *strings.Builder, entityName string) {
	b.WriteString("\t// TODO: Implement MongoDB find by ID\n")
	b.WriteString("\treturn nil, fmt.Errorf(\"not implemented\")\n")
}

func (g *RepositoryGenerator) generateSQL(outputDir string) (string, error) {
	filename := filepath.Join(outputDir, "schema.sql")

	var b strings.Builder

	entityName := toPascalCase(g.domainName)
	tableName := toSnakeCase(entityName) + "s"

	b.WriteString(fmt.Sprintf("-- Schema for %s table\n\n", tableName))

	switch g.config.Database {
	case "postgres":
		b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))
		b.WriteString("\tid UUID PRIMARY KEY,\n")
		b.WriteString("\tcreated_at TIMESTAMP NOT NULL DEFAULT NOW(),\n")
		b.WriteString("\tupdated_at TIMESTAMP NOT NULL DEFAULT NOW()\n")
		b.WriteString("\t-- TODO: Add domain-specific columns\n")
		b.WriteString(");\n\n")

		b.WriteString(fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_created_at ON %s(created_at);\n", tableName, tableName))
	case "mysql":
		b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))
		b.WriteString("\tid CHAR(36) PRIMARY KEY,\n")
		b.WriteString("\tcreated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n")
		b.WriteString("\tupdated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP\n")
		b.WriteString("\t-- TODO: Add domain-specific columns\n")
		b.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n\n")

		b.WriteString(fmt.Sprintf("CREATE INDEX idx_%s_created_at ON %s(created_at);\n", tableName, tableName))
	}

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *RepositoryGenerator) generateRepositoryTest(outputDir string) (string, error) {
	filename := filepath.Join(outputDir, g.domainName+"_repo_test.go")

	var b strings.Builder

	// Package
	b.WriteString(fmt.Sprintf("package %s_test\n\n", g.config.Database))

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"testing\"\n\n")
	b.WriteString("\t\"github.com/stretchr/testify/assert\"\n")
	b.WriteString(")\n\n")

	// Test placeholder
	entityName := toPascalCase(g.domainName)
	b.WriteString(fmt.Sprintf("func Test%sRepository_Save(t *testing.T) {\n", entityName))
	b.WriteString("\t// TODO: Implement repository tests\n")
	b.WriteString("\tassert.True(t, true)\n")
	b.WriteString("}\n")

	// Write file
	if err := os.WriteFile(filename, []byte(b.String()), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

// detectModuleName reads go.mod and extracts module name  
func (g *RepositoryGenerator) detectModuleName() error {
	goModFile := "go.mod"
	data, err := os.ReadFile(goModFile)
	if err != nil {
		return fmt.Errorf("read go.mod: %w", err)
	}

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
func (g *RepositoryGenerator) detectDatabaseType() error {
	envFile := ".env"
	data, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("read .env: %w", err)
	}

	envContent := string(data)

	if strings.Contains(envContent, "DATABASE_URL=") {
		for _, line := range strings.Split(envContent, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "DATABASE_URL=") {
				dbURL := strings.TrimPrefix(line, "DATABASE_URL=")

				if strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://") {
					g.config.Database = "postgres"
				} else if strings.HasPrefix(dbURL, "mysql://") {
					g.config.Database = "mysql"
				} else if strings.HasPrefix(dbURL, "sqlite://") || strings.Contains(dbURL, ".db") {
					g.config.Database = "sqlite"
				} else if strings.HasPrefix(dbURL, "mongodb://") {
					g.config.Database = "mongodb"
				}

				g.config.Logger.Info("detected database type", "type", g.config.Database)
				return nil
			}
		}
	}

	return fmt.Errorf("DATABASE_URL not found in .env")
}
