package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/spf13/cobra"
)

var wireCmd = &cobra.Command{
	Use:   "wire",
	Short: "Auto-wire dependencies and generate main.go",
	Long: `Automatically detect all domains and wire dependencies.

This command will:
- Scan for existing domains (entities, repositories, handlers)
- Generate dependency injection code
- Wire everything to main.go

Example:
  anaphase wire
  anaphase wire --output cmd/api`,
	RunE: runWire,
}

var (
	wireOutput string
)

func init() {
	wireCmd.Flags().StringVar(&wireOutput, "output", "cmd/api", "Output directory for main.go")
	rootCmd.AddCommand(wireCmd)
}

func runWire(cmd *cobra.Command, args []string) error {
	fmt.Println("⚡ Auto-wiring dependencies...\n")

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create wire generator
	gen := generator.NewWireGenerator(&generator.WireConfig{
		OutputDir: wireOutput,
		Logger:    logger,
	})

	// Generate wiring code
	ctx := context.Background()
	files, err := gen.Generate(ctx)
	if err != nil {
		return fmt.Errorf("generate wiring: %w", err)
	}

	// Report results
	fmt.Println("✅ Generated files:")
	for _, file := range files {
		fmt.Printf("  ✓ %s\n", file)
	}

	fmt.Println("\n🎉 Auto-wiring complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review generated main.go")
	fmt.Println("  2. Set up .env file with DB credentials")
	fmt.Println("  3. Run: go run cmd/api/main.go")

	return nil
}
