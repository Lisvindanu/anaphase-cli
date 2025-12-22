package commands

import (
	"fmt"

	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	genMiddlewareOutput string
)

var genMiddlewareCmd = &cobra.Command{
	Use:   "middleware",
	Short: "Generate HTTP middleware",
	Long: `Generate common HTTP middleware for your application.

Available middleware types:
  --type auth       - JWT authentication middleware
  --type ratelimit  - Rate limiting middleware
  --type logging    - Structured logging middleware
  --type cors       - CORS configuration middleware

Example:
  anaphase gen middleware --type auth
  anaphase gen middleware --type ratelimit --output internal/middleware
  anaphase gen middleware --type logging
  anaphase gen middleware --type cors`,
	RunE: runGenMiddleware,
}

var middlewareType string

func init() {
	genCmd.AddCommand(genMiddlewareCmd)

	genMiddlewareCmd.Flags().StringVar(&middlewareType, "type", "", "Middleware type (auth, ratelimit, logging, cors)")
	genMiddlewareCmd.Flags().StringVar(&genMiddlewareOutput, "output", "internal/middleware", "Output directory for generated middleware")

	genMiddlewareCmd.MarkFlagRequired("type")
}

func runGenMiddleware(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderTitle("Middleware Generator"))

	// Validate middleware type
	validTypes := map[string]generator.MiddlewareType{
		"auth":      generator.MiddlewareAuth,
		"ratelimit": generator.MiddlewareRateLimit,
		"logging":   generator.MiddlewareLogging,
		"cors":      generator.MiddlewareCORS,
	}

	mwType, valid := validTypes[middlewareType]
	if !valid {
		ui.PrintError(fmt.Sprintf("Invalid middleware type: %s", middlewareType))
		fmt.Println("\nValid types: auth, ratelimit, logging, cors")
		return fmt.Errorf("invalid middleware type")
	}

	ui.PrintInfo(fmt.Sprintf("Type: %s", middlewareType))
	ui.PrintInfo(fmt.Sprintf("Output: %s", genMiddlewareOutput))
	fmt.Println()

	// Generate middleware
	fmt.Println("📦 Generating middleware...")
	gen := generator.NewMiddlewareGenerator(mwType, genMiddlewareOutput)
	files, err := gen.Generate()

	if err != nil {
		ui.PrintError(fmt.Sprintf("Generation failed: %v", err))
		return fmt.Errorf("generate middleware: %w", err)
	}

	// Show generated files
	fmt.Println(ui.SuccessStyle.Render("\nGenerated Files:"))
	for _, file := range files {
		fmt.Println(ui.RenderListItem(file, true))
	}

	fmt.Println()
	ui.PrintSuccess("Middleware generation complete!")

	// Show usage instructions based on type
	fmt.Println(ui.RenderSubtle("\nUsage Instructions:"))
	switch mwType {
	case generator.MiddlewareAuth:
		fmt.Println("  1. Set your JWT secret key (environment variable or config)")
		fmt.Println("  2. Import: import \"yourproject/internal/middleware\"")
		fmt.Println("  3. Use in your HTTP router:")
		fmt.Println("     config := middleware.AuthConfig{")
		fmt.Println("         SecretKey: os.Getenv(\"JWT_SECRET\"),")
		fmt.Println("         SkipPaths: []string{\"/health\", \"/login\"},")
		fmt.Println("     }")
		fmt.Println("     router.Use(middleware.AuthMiddleware(config))")

	case generator.MiddlewareRateLimit:
		fmt.Println("  1. Import: import \"yourproject/internal/middleware\"")
		fmt.Println("  2. Configure rate limits:")
		fmt.Println("     config := middleware.RateLimitConfig{")
		fmt.Println("         Rate:     100,              // 100 requests")
		fmt.Println("         Interval: time.Minute,      // per minute")
		fmt.Println("         MaxBurst: 120,              // allow burst of 120")
		fmt.Println("     }")
		fmt.Println("     router.Use(middleware.RateLimitMiddleware(config))")

	case generator.MiddlewareLogging:
		fmt.Println("  1. Import: import \"yourproject/internal/middleware\"")
		fmt.Println("  2. Set up logger and apply middleware:")
		fmt.Println("     logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))")
		fmt.Println("     config := middleware.LoggingConfig{")
		fmt.Println("         Logger:    logger,")
		fmt.Println("         SkipPaths: []string{\"/health\"},")
		fmt.Println("     }")
		fmt.Println("     router.Use(middleware.RequestIDMiddleware())")
		fmt.Println("     router.Use(middleware.LoggingMiddleware(config))")

	case generator.MiddlewareCORS:
		fmt.Println("  1. Import: import \"yourproject/internal/middleware\"")
		fmt.Println("  2. Configure CORS (Development):")
		fmt.Println("     router.Use(middleware.CORSMiddleware(middleware.DefaultCORSConfig()))")
		fmt.Println()
		fmt.Println("  3. Configure CORS (Production):")
		fmt.Println("     config := middleware.ProductionCORSConfig([]string{")
		fmt.Println("         \"https://example.com\",")
		fmt.Println("         \"https://app.example.com\",")
		fmt.Println("     })")
		fmt.Println("     router.Use(middleware.CORSMiddleware(config))")
	}

	fmt.Println()
	fmt.Println(ui.RenderSubtle("Next Steps:"))
	fmt.Println("  1. Review generated middleware code")
	fmt.Println("  2. Customize configuration as needed")
	fmt.Println("  3. Integrate into your HTTP router/server")
	fmt.Println("  4. Run: go build ./...")
	fmt.Println()

	return nil
}
