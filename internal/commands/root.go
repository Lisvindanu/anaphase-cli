package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lisvindanu/anaphase-cli/internal/ui"
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
	Run: func(cmd *cobra.Command, args []string) {
		// Show interactive menu when no subcommand is provided
		showInteractiveMenu(cmd)
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// showInteractiveMenu shows an interactive TUI menu
func showInteractiveMenu(cmd *cobra.Command) {
	m := ui.NewMenuModel()
	// Don't use alternate screen - let menu and output share the same buffer
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running menu: %v\n", err)
		os.Exit(1)
	}

	// Get the selected command
	if menuModel, ok := finalModel.(ui.MenuModel); ok {
		choice := menuModel.GetChoice()
		if choice != "" {
			// Clear the screen first to remove the menu
			fmt.Print("\033[2J\033[H")

			// Parse and execute the selected command
			cmdParts := ui.FormatCommand(choice)

			// Show info about the selected command
			fmt.Printf("%s Running: anaphase %s\n\n", ui.RenderInfo("ℹ"), choice)

			// Find and execute the subcommand
			subCmd, _, err := cmd.Root().Find(cmdParts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding command: %v\n", err)
				return
			}

			// Ensure command output goes to stdout/stderr properly
			subCmd.SetOut(os.Stdout)
			subCmd.SetErr(os.Stderr)

			// Set args and execute
			subCmd.SetArgs(cmdParts[1:]) // Skip the first part which is the command name
			if err := subCmd.Execute(); err != nil {
				os.Exit(1)
			}
		}
	}
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
