package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

var rootCmd = &cobra.Command{
	Use:   "anaphase",
	Short: "Anaphase - AI-Powered Golang Microservice Generator",
	Long: `Anaphase CLI is an intelligent code scaffolding tool that generates
production-ready Golang microservices with Domain-Driven Design architecture.

Features:
  - AI-powered domain generation from natural language
  - Automatic dependency injection with AST manipulation
  - Clean Architecture enforcement
  - Complete test generation
  - OpenAPI/Swagger documentation`,
	Version: version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")
}

// exitWithError prints an error message and exits with status 1
func exitWithError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+msg+"\n", args...)
	os.Exit(1)
}
