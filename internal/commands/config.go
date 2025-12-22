package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lisvindanu/anaphase-cli/internal/ai"
	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Anaphase configuration",
	Long: `Manage Anaphase configuration including AI providers, cache settings, and more.

Available subcommands:
  list            - Show current configuration
  set-provider    - Set default AI provider
  check           - Health check all providers
  show-providers  - List available providers`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show current configuration",
	Long:  "Display the current Anaphase configuration including AI providers, cache settings, and more.",
	RunE:  runConfigList,
}

var configSetProviderCmd = &cobra.Command{
	Use:   "set-provider <provider>",
	Short: "Set default AI provider",
	Long: `Set the default AI provider for code generation.

Available providers: gemini, groq, openai, claude, ollama

Example:
  anaphase config set-provider groq
  anaphase config set-provider gemini`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigSetProvider,
}

var configCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Health check all providers",
	Long:  "Check the health and availability of all configured AI providers.",
	RunE:  runConfigCheck,
}

var configShowProvidersCmd = &cobra.Command{
	Use:   "show-providers",
	Short: "List available providers",
	Long:  "Show all available AI providers and their current status.",
	RunE:  runConfigShowProviders,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetProviderCmd)
	configCmd.AddCommand(configCheckCmd)
	configCmd.AddCommand(configShowProvidersCmd)
}

func runConfigList(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderTitle("Anaphase Configuration"))

	cfg, err := ai.LoadConfig()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to load config: %v", err))
		return err
	}

	// AI Configuration
	fmt.Println(ui.InfoStyle.Render("\n🤖 AI Configuration:"))
	fmt.Printf("  Primary Provider: %s\n", ui.SuccessStyle.Render(cfg.AI.PrimaryProvider))
	fmt.Printf("  Fallback Providers: %v\n", cfg.AI.FallbackProviders)
	fmt.Println()

	// Provider Details
	fmt.Println(ui.InfoStyle.Render("📡 Configured Providers:"))

	// Gemini
	if cfg.AI.Providers.Gemini.Enabled || cfg.AI.Providers.Gemini.APIKey != "" {
		fmt.Printf("  %s Gemini\n", ui.CheckmarkStyle.Render())
		fmt.Printf("    Model: %s\n", cfg.AI.Providers.Gemini.Model)
		fmt.Printf("    Timeout: %s\n", cfg.AI.Providers.Gemini.Timeout)
		fmt.Printf("    Max Retries: %d\n", cfg.AI.Providers.Gemini.MaxRetries)
	}

	// Groq
	if cfg.AI.Providers.Groq.Enabled || cfg.AI.Providers.Groq.APIKey != "" {
		fmt.Printf("  %s Groq\n", ui.CheckmarkStyle.Render())
		fmt.Printf("    Model: %s\n", cfg.AI.Providers.Groq.Model)
		fmt.Printf("    Timeout: %s\n", cfg.AI.Providers.Groq.Timeout)
		fmt.Printf("    Max Retries: %d\n", cfg.AI.Providers.Groq.MaxRetries)
	}

	// OpenAI
	if cfg.AI.Providers.OpenAI.Enabled || cfg.AI.Providers.OpenAI.APIKey != "" {
		fmt.Printf("  %s OpenAI\n", ui.CheckmarkStyle.Render())
		fmt.Printf("    Model: %s\n", cfg.AI.Providers.OpenAI.Model)
	}

	// Cache Configuration
	fmt.Println(ui.InfoStyle.Render("\n💾 Cache Configuration:"))
	fmt.Printf("  Enabled: %v\n", cfg.Cache.Enabled)
	fmt.Printf("  Directory: %s\n", cfg.Cache.Directory)
	fmt.Printf("  TTL: %s\n", cfg.Cache.TTL)
	fmt.Println()

	// Generator Configuration
	fmt.Println(ui.InfoStyle.Render("⚙️  Generator Settings:"))
	fmt.Printf("  Go Version: %s\n", cfg.Generator.GoVersion)
	fmt.Printf("  Code Style: %s\n", cfg.Generator.CodeStyle)

	return nil
}

func runConfigSetProvider(cmd *cobra.Command, args []string) error {
	provider := args[0]

	// Validate provider
	validProviders := []string{"gemini", "groq", "openai", "claude", "ollama"}
	valid := false
	for _, p := range validProviders {
		if p == provider {
			valid = true
			break
		}
	}

	if !valid {
		ui.PrintError(fmt.Sprintf("Invalid provider: %s", provider))
		fmt.Println("\nValid providers:", validProviders)
		return fmt.Errorf("invalid provider")
	}

	cfg, err := ai.LoadConfig()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to load config: %v", err))
		return err
	}

	// Update primary provider
	cfg.AI.PrimaryProvider = provider

	ui.PrintSuccess(fmt.Sprintf("Default provider set to: %s", provider))
	ui.PrintInfo("Note: Config changes are temporary. To persist, edit ~/.anaphase/config.yaml")

	return nil
}

func runConfigCheck(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderTitle("Provider Health Check"))

	cfg, err := ai.LoadConfig()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to load config: %v", err))
		return err
	}

	// Create orchestrator
	logger := cmd.Context().Value("logger")
	if logger == nil {
		logger = os.Stdout
	}

	orchestrator, err := ai.NewOrchestrator(cfg, nil)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to create orchestrator: %v", err))
		return err
	}

	// Check all providers
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := orchestrator.ValidateProviders(ctx)

	fmt.Println()
	for provider, err := range results {
		if err == nil {
			fmt.Printf("  %s %s - %s\n",
				ui.CheckmarkStyle.Render(),
				provider,
				ui.SuccessStyle.Render("Healthy"))
		} else {
			fmt.Printf("  %s %s - %s\n",
				ui.CrossStyle.Render(),
				provider,
				ui.ErrorStyle.Render(err.Error()))
		}
	}

	fmt.Println()
	return nil
}

func runConfigShowProviders(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderTitle("Available AI Providers"))

	providers := []struct {
		name        string
		description string
		models      []string
		free        bool
	}{
		{
			name:        "Gemini",
			description: "Google's AI model - Reliable, free tier available",
			models:      []string{"gemini-2.0-flash-exp", "gemini-pro"},
			free:        true,
		},
		{
			name:        "Groq",
			description: "Extremely fast inference - Free during preview",
			models:      []string{"llama-3.3-70b-versatile", "mixtral-8x7b-32768"},
			free:        true,
		},
		{
			name:        "OpenAI",
			description: "GPT models - Paid service",
			models:      []string{"gpt-4o", "gpt-4o-mini"},
			free:        false,
		},
		{
			name:        "Claude",
			description: "Anthropic's AI - Paid service",
			models:      []string{"claude-3-5-sonnet-20241022"},
			free:        false,
		},
		{
			name:        "Ollama",
			description: "Local AI - Run models on your machine",
			models:      []string{"qwen2.5-coder:7b", "codellama"},
			free:        true,
		},
	}

	fmt.Println()
	for _, p := range providers {
		fmt.Printf("\n%s %s\n", ui.InfoStyle.Render("📡"), ui.SuccessStyle.Render(p.name))
		fmt.Printf("  %s\n", p.description)
		if p.free {
			fmt.Printf("  💰 Free: %s\n", ui.SuccessStyle.Render("Yes"))
		} else {
			fmt.Printf("  💰 Free: %s\n", ui.ErrorStyle.Render("No (Paid)"))
		}
		fmt.Printf("  🎯 Models: %v\n", p.models)
	}

	fmt.Println()
	fmt.Println(ui.RenderSubtle("To set a provider as default:"))
	fmt.Println(ui.RenderSubtle("  anaphase config set-provider <provider>"))
	fmt.Println()

	return nil
}
