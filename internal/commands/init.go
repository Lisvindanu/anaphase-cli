package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/lisvindanu/anaphase-cli/internal/setup"
	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	initModule   string
	initDB       string
	initCache    bool
	initEventBus string
	initNoDocker bool
)

var initCmd = &cobra.Command{
	Use:   "init <project_name>",
	Short: "Initialize a new microservice project",
	Long: `Initialize a new microservice project with Domain-Driven Design architecture.

This command creates a complete project structure with:
  - Clean Architecture folder layout
  - Go module initialization
  - Docker Compose for local development
  - Makefile for common tasks
  - Configuration management
  - Basic server setup

Example:
  anaphase init my-shop
  anaphase init my-shop --module github.com/mycompany/my-shop
  anaphase init my-shop --db postgres --cache --event-bus nats`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags
	initCmd.Flags().StringVar(&initModule, "module", "", "Go module name (default: github.com/user/<project_name>)")
	initCmd.Flags().StringVar(&initDB, "db", "postgres", "Database type: postgres|mysql|mongodb")
	initCmd.Flags().BoolVar(&initCache, "cache", false, "Include Redis in docker-compose")
	initCmd.Flags().StringVar(&initEventBus, "event-bus", "", "Message broker: nats|kafka|rabbitmq")
	initCmd.Flags().BoolVar(&initNoDocker, "no-docker", false, "Skip docker-compose.yml generation")
}

func runInit(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Validate project name
	if !isValidProjectName(projectName) {
		return fmt.Errorf("invalid project name '%s': must contain only letters, numbers, hyphens, and underscores", projectName)
	}

	// Set default module name if not provided
	if initModule == "" {
		username := os.Getenv("USER")
		if username == "" {
			username = "user"
		}
		initModule = fmt.Sprintf("github.com/%s/%s", username, projectName)
	}

	// Validate database type
	validDBs := map[string]bool{"postgres": true, "mysql": true, "mongodb": true}
	if !validDBs[initDB] {
		return fmt.Errorf("invalid database type '%s': must be postgres, mysql, or mongodb", initDB)
	}

	// Validate event bus if provided
	if initEventBus != "" {
		validEventBus := map[string]bool{"nats": true, "kafka": true, "rabbitmq": true}
		if !validEventBus[initEventBus] {
			return fmt.Errorf("invalid event bus '%s': must be nats, kafka, or rabbitmq", initEventBus)
		}
	}

	// Create project configuration
	config := &generator.ProjectConfig{
		Name:       projectName,
		Module:     initModule,
		Database:   initDB,
		Cache:      initCache,
		EventBus:   initEventBus,
		SkipDocker: initNoDocker,
		OutputDir:  projectName,
	}

	// Check if directory already exists
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	// Create project generator
	gen := generator.NewProjectGenerator(config)

	// Generate project
	fmt.Printf("Creating new project '%s'...\n", projectName)
	fmt.Printf("  Module: %s\n", initModule)
	fmt.Printf("  Database: %s\n", initDB)
	if initCache {
		fmt.Printf("  Cache: Redis\n")
	}
	if initEventBus != "" {
		fmt.Printf("  Event Bus: %s\n", initEventBus)
	}
	fmt.Println()

	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Printf("\n✅ Project '%s' created successfully!\n\n", projectName)

	// Auto-setup: Ensure Anaphase config exists
	fmt.Println("🔧 Setting up Anaphase configuration...")
	if err := setup.CreateAnaphaseConfig(); err != nil {
		fmt.Printf("Warning: Could not create Anaphase config: %v\n", err)
	}

	// Auto-setup: Ensure project configs
	fmt.Println("📝 Setting up project configuration files...")
	// Change to project directory temporarily
	currentDir, _ := os.Getwd()
	os.Chdir(projectName)

	if err := setup.EnsureProjectConfig(); err != nil {
		fmt.Printf("Warning: Could not create project configs: %v\n", err)
	}

	if err := setup.EnsureGitignore(); err != nil {
		fmt.Printf("Warning: Could not update .gitignore: %v\n", err)
	}

	// Create .env.example file
	ui.PrintInfo("📝 Creating .env.example...")

	// Generate DATABASE_URL based on selected database type
	var dbURL string
	switch initDB {
	case "postgres":
		dbURL = "postgresql://username:password@localhost:5432/" + projectName + "?sslmode=disable"
	case "mysql":
		dbURL = "mysql://username:password@localhost:3306/" + projectName + "?parseTime=true"
	case "sqlite":
		dbURL = "sqlite://./data/" + projectName + ".db"
	case "mongodb":
		dbURL = "mongodb://localhost:27017/" + projectName
	default:
		dbURL = "postgresql://username:password@localhost:5432/" + projectName + "?sslmode=disable"
	}

	envExample := `# Database Configuration
DATABASE_URL=` + dbURL + `

# Server Configuration
PORT=8080
ENV=development

# JWT Configuration (if using auth)
JWT_SECRET=your-secret-key-change-this

# Redis Configuration (if using cache)
REDIS_URL=redis://localhost:6379

# Logging
LOG_LEVEL=info
`
	if err := os.WriteFile(".env.example", []byte(envExample), 0644); err != nil {
		ui.PrintWarning(fmt.Sprintf("Warning: Could not create .env.example: %v", err))
	} else {
		ui.PrintSuccess("✅ Created .env.example")
		ui.PrintInfo("💡 Copy .env.example to .env and update with your credentials")
	}

	// Run go mod tidy to download dependencies
	fmt.Println()
	ui.PrintInfo("📦 Installing dependencies (go mod tidy)...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr

	if err := tidyCmd.Run(); err != nil {
		ui.PrintWarning(fmt.Sprintf("Warning: go mod tidy failed: %v", err))
		ui.PrintInfo("You can run 'go mod tidy' manually later")
	} else {
		ui.PrintSuccess("✅ Dependencies installed successfully!")
	}

	// Return to original directory
	os.Chdir(currentDir)

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  make run")
	fmt.Println()
	fmt.Println("To generate a domain:")
	fmt.Printf("  anaphase gen domain \"your domain description\"\n")
	fmt.Println()

	return nil
}

func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}

	// Check first character is letter
	if !isLetter(rune(name[0])) {
		return false
	}

	// Check all characters are valid
	for _, ch := range name {
		if !isLetter(ch) && !isDigit(ch) && ch != '-' && ch != '_' {
			return false
		}
	}

	return true
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
