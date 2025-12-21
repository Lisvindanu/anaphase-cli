package commands

import (
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code components",
	Long: `Generate various code components using AI or templates.

Available subcommands:
  domain      - Generate domain entities and business logic
  handler     - Generate HTTP/gRPC handlers
  repository  - Generate database repositories
  test        - Generate unit tests
  docs        - Generate documentation`,
}

func init() {
	rootCmd.AddCommand(genCmd)
}
