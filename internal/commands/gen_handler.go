package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/spf13/cobra"
)

var genHandlerCmd = &cobra.Command{
	Use:   "handler <domain_name>",
	Short: "Generate HTTP handlers for a domain",
	Long: `Generate HTTP handlers, DTOs, and tests for a domain entity.

Example:
  anaphase gen handler customer
  anaphase gen handler order --auth --validate`,
	Args: cobra.ExactArgs(1),
	RunE: runGenHandler,
}

var (
	handlerProtocol string
	handlerAuth     bool
	handlerValidate bool
)

func init() {
	genHandlerCmd.Flags().StringVar(&handlerProtocol, "protocol", "http", "Protocol type: http|grpc")
	genHandlerCmd.Flags().BoolVar(&handlerAuth, "auth", false, "Include JWT auth middleware")
	genHandlerCmd.Flags().BoolVar(&handlerValidate, "validate", false, "Include request validation")
	genCmd.AddCommand(genHandlerCmd)
}

func runGenHandler(cmd *cobra.Command, args []string) error {
	domainName := args[0]

	fmt.Printf("🔨 Generating %s handlers for domain: %s\n\n", handlerProtocol, domainName)

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create handler generator
	gen := generator.NewHandlerGenerator(domainName, &generator.HandlerConfig{
		Protocol: handlerProtocol,
		Auth:     handlerAuth,
		Validate: handlerValidate,
		Logger:   logger,
	})

	// Generate files
	ctx := context.Background()
	files, err := gen.Generate(ctx)
	if err != nil {
		return fmt.Errorf("generate handlers: %w", err)
	}

	// Report results
	fmt.Println("✅ Generated files:")
	for _, file := range files {
		fmt.Printf("  ✓ %s\n", file)
	}

	fmt.Println("\n🎉 Handler generation complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review generated handlers")
	fmt.Println("  2. Run: go build ./...")
	fmt.Println("  3. Wire handlers in main.go")

	return nil
}
