package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/spf13/cobra"
)

var genRepositoryCmd = &cobra.Command{
	Use:   "repository <domain_name>",
	Short: "Generate repository implementation for a domain",
	Long: `Generate database-specific repository implementation.

Example:
  anaphase gen repository customer --db postgres
  anaphase gen repository order --db postgres --cache`,
	Args: cobra.ExactArgs(1),
	RunE: runGenRepository,
}

var (
	repositoryDB    string
	repositoryCache bool
)

func init() {
	genRepositoryCmd.Flags().StringVar(&repositoryDB, "db", "postgres", "Database type: postgres|mysql|mongodb")
	genRepositoryCmd.Flags().BoolVar(&repositoryCache, "cache", false, "Include Redis caching layer")
	genCmd.AddCommand(genRepositoryCmd)
}

func runGenRepository(cmd *cobra.Command, args []string) error {
	domainName := args[0]

	fmt.Printf("💾 Generating %s repository for domain: %s\n\n", repositoryDB, domainName)

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create repository generator
	gen := generator.NewRepositoryGenerator(domainName, &generator.RepositoryConfig{
		Database: repositoryDB,
		Cache:    repositoryCache,
		Logger:   logger,
	})

	// Generate files
	ctx := context.Background()
	files, err := gen.Generate(ctx)
	if err != nil {
		return fmt.Errorf("generate repository: %w", err)
	}

	// Report results
	fmt.Println("✅ Generated files:")
	for _, file := range files {
		fmt.Printf("  ✓ %s\n", file)
	}

	fmt.Println("\n🎉 Repository generation complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review generated repository")
	fmt.Println("  2. Set up database connection")
	fmt.Println("  3. Run: go build ./...")

	return nil
}
