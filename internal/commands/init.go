package commands

import (
	"fmt"
	"os"

	"github.com/lisvindanuu/anaphase-cli/internal/generator"
	"github.com/spf13/cobra"
)

var (
	initModule    string
	initDB        string
	initCache     bool
	initEventBus  string
	initNoDocker  bool
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
		Name:         projectName,
		Module:       initModule,
		Database:     initDB,
		Cache:        initCache,
		EventBus:     initEventBus,
		SkipDocker:   initNoDocker,
		OutputDir:    projectName,
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
