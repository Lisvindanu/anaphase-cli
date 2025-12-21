package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lisvindanuu/anaphase-cli/internal/ai"
	"github.com/lisvindanuu/anaphase-cli/internal/generator"
	"github.com/spf13/cobra"
)

var (
	genDomainOutput string
)

var genDomainCmd = &cobra.Command{
	Use:   "domain <description>",
	Short: "Generate domain entities and business logic using AI",
	Long: `Generate domain entities, value objects, repositories, and services using AI.

This command uses AI to analyze your business requirement and generate:
  - Domain entities with validation
  - Value objects (immutable types)
  - Repository interfaces
  - Service interfaces
  - Business logic methods

Example:
  anaphase gen domain "Cart with Items. User can add, remove, update quantity"
  anaphase gen domain "Order has ID, Total, Status. Can be cancelled if pending"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runGenDomain,
}

func init() {
	genCmd.AddCommand(genDomainCmd)

	genDomainCmd.Flags().StringVar(&genDomainOutput, "output", "internal/core", "Output directory for generated files")
}

func runGenDomain(cmd *cobra.Command, args []string) error {
	// Combine all args as description
	description := ""
	for i, arg := range args {
		if i > 0 {
			description += " "
		}
		description += arg
	}

	fmt.Println("🤖 Generating domain using AI...")
	fmt.Printf("📝 Description: %s\n\n", description)

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load AI configuration
	fmt.Println("⚙️  Loading configuration...")
	cfg, err := ai.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Create orchestrator
	orchestrator, err := ai.NewOrchestrator(cfg, logger)
	if err != nil {
		return fmt.Errorf("create orchestrator: %w", err)
	}

	// Generate domain spec using AI
	fmt.Println("🧠 Analyzing with AI...")
	ctx := context.Background()
	spec, err := ai.GenerateDomain(ctx, orchestrator, description)
	if err != nil {
		return fmt.Errorf("generate domain: %w", err)
	}

	fmt.Printf("✅ AI Analysis Complete!\n\n")
	fmt.Printf("Domain: %s\n", spec.DomainName)
	fmt.Printf("Entities: %d\n", len(spec.Entities))
	fmt.Printf("Value Objects: %d\n", len(spec.ValueObjects))
	fmt.Printf("Repository: %s\n", spec.RepositoryInterface.Name)
	fmt.Printf("Service: %s\n\n", spec.ServiceInterface.Name)

	// Generate code files
	fmt.Println("📂 Generating files...")
	domainGen := generator.NewDomainGenerator(spec, genDomainOutput)

	files, err := domainGen.Generate()
	if err != nil {
		return fmt.Errorf("generate files: %w", err)
	}

	// Show generated files
	fmt.Println("\n✅ Generated files:")
	for _, file := range files {
		fmt.Printf("  ✓ %s\n", file)
	}

	fmt.Println("\n🎉 Domain generation complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review generated files")
	fmt.Println("  2. Run: go build ./...")
	fmt.Println("  3. Generate handler: anaphase gen handler", spec.DomainName)

	return nil
}
