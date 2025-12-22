package commands

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/lisvindanu/anaphase-cli/internal/ai"
	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	genDomainOutput      string
	genDomainProvider    string
	genDomainInteractive bool
)

var genDomainCmd = &cobra.Command{
	Use:   "domain [description]",
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
  anaphase gen domain "Order has ID, Total, Status. Can be cancelled if pending"
  anaphase gen domain "User with email" --provider groq
  anaphase gen domain "Product catalog" --provider gemini
  anaphase gen domain --interactive`,
	RunE: runGenDomain,
}

func init() {
	genCmd.AddCommand(genDomainCmd)

	genDomainCmd.Flags().StringVar(&genDomainOutput, "output", "internal/core", "Output directory for generated files")
	genDomainCmd.Flags().StringVar(&genDomainProvider, "provider", "", "AI provider to use (gemini, groq, openai, claude)")
	genDomainCmd.Flags().BoolVarP(&genDomainInteractive, "interactive", "i", false, "Run in interactive mode")
}

// promptInput prompts the user for input with a message
func promptInput(message string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", message, ui.SubtleStyle.Render(defaultValue))
	} else {
		fmt.Printf("%s: ", message)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}

	return input
}

// promptChoice prompts the user to choose from options
func promptChoice(message string, options []string, defaultIdx int) string {
	fmt.Println(message)
	for i, opt := range options {
		if i == defaultIdx {
			fmt.Printf("  %d) %s %s\n", i+1, opt, ui.SubtleStyle.Render("(default)"))
		} else {
			fmt.Printf("  %d) %s\n", i+1, opt)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter choice [%d]: ", defaultIdx+1)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return options[defaultIdx]
	}

	// Parse choice
	var choice int
	fmt.Sscanf(input, "%d", &choice)

	if choice < 1 || choice > len(options) {
		return options[defaultIdx]
	}

	return options[choice-1]
}

func runGenDomain(cmd *cobra.Command, args []string) error {
	var description string
	var provider string
	var output string

	// Interactive mode
	if genDomainInteractive {
		fmt.Println(ui.RenderTitle("Interactive Domain Generation"))
		fmt.Println()

		// Prompt for description
		description = promptInput("Enter domain description", "")
		for description == "" {
			ui.PrintError("Description cannot be empty")
			description = promptInput("Enter domain description", "")
		}
		fmt.Println()

		// Prompt for AI provider
		providers := []string{"gemini", "groq", "openai", "claude"}
		provider = promptChoice("Select AI provider:", providers, 0)
		fmt.Println()

		// Prompt for output directory
		output = promptInput("Output directory", "internal/core")
		fmt.Println()

	} else {
		// Non-interactive mode - validate args
		if len(args) == 0 {
			ui.PrintError("Description is required. Use --interactive for guided mode.")
			return fmt.Errorf("description required")
		}

		// Combine all args as description
		description = ""
		for i, arg := range args {
			if i > 0 {
				description += " "
			}
			description += arg
		}

		provider = genDomainProvider
		output = genDomainOutput
	}

	fmt.Println(ui.RenderTitle("AI-Powered Domain Generation"))
	ui.PrintInfo(fmt.Sprintf("Description: %s", description))
	fmt.Println()

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load AI configuration
	fmt.Println("⚙️  Step 1/3: Loading configuration...")
	cfg, err := ai.LoadConfig()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to load config: %v", err))
		return fmt.Errorf("load config: %w", err)
	}

	// Override provider if specified
	if provider != "" {
		cfg.AI.PrimaryProvider = provider
		ui.PrintInfo(fmt.Sprintf("Using provider: %s", provider))
	} else {
		ui.PrintInfo(fmt.Sprintf("Using provider: %s", cfg.AI.PrimaryProvider))
	}

	// Create orchestrator
	orchestrator, err := ai.NewOrchestrator(cfg, logger)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to create orchestrator: %v", err))
		return fmt.Errorf("create orchestrator: %w", err)
	}

	// Generate domain spec using AI
	fmt.Println("\n🧠 Step 2/3: Analyzing with AI...")
	ctx := context.Background()
	spec, err := ai.GenerateDomain(ctx, orchestrator, description)

	if err != nil {
		ui.PrintError(fmt.Sprintf("AI generation failed: %v", err))
		return fmt.Errorf("generate domain: %w", err)
	}

	ui.PrintSuccess("AI Analysis Complete!")
	fmt.Println()

	fmt.Println(ui.InfoStyle.Render("Generated Specification:"))
	fmt.Printf("  📦 Domain: %s\n", spec.DomainName)
	fmt.Printf("  📄 Entities: %d\n", len(spec.Entities))
	fmt.Printf("  📄 Value Objects: %d\n", len(spec.ValueObjects))
	fmt.Printf("  ⚙️  Repository: %s\n", spec.RepositoryInterface.Name)
	fmt.Printf("  ⚙️  Service: %s\n", spec.ServiceInterface.Name)
	fmt.Println()

	// Generate code files
	fmt.Println("📂 Step 3/3: Generating code files...")
	domainGen := generator.NewDomainGenerator(spec, output)
	files, err := domainGen.Generate()

	if err != nil {
		ui.PrintError(fmt.Sprintf("Code generation failed: %v", err))
		return fmt.Errorf("generate files: %w", err)
	}

	// Show generated files
	fmt.Println(ui.SuccessStyle.Render("\nGenerated Files:"))
	for _, file := range files {
		fmt.Println(ui.RenderListItem(file, true))
	}

	fmt.Println()
	ui.PrintSuccess("Domain generation complete! 🚀")

	fmt.Println(ui.RenderSubtle("\nNext Steps:"))
	fmt.Println("  1. Review generated files")
	fmt.Println("  2. Run: go build ./...")
	fmt.Println("  3. Generate handler: anaphase gen handler " + spec.DomainName)

	return nil
}
