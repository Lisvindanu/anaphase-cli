package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running menu: %v\n", err)
		os.Exit(1)
	}

	// Get the selected command
	if menuModel, ok := finalModel.(ui.MenuModel); ok {
		choice := menuModel.GetChoice()
		selectedItem := menuModel.GetSelectedItem()

		if choice != "" && selectedItem != nil {
			// Parse command parts
			cmdParts := ui.FormatCommand(choice)
			originalCmdLen := len(cmdParts)

			// Collect inputs if needed
			var inputs []string
			if selectedItem.NeedsInput() {
				fmt.Println()
				fmt.Println(ui.RenderTitle(selectedItem.Title()))
				fmt.Println()

				scanner := bufio.NewScanner(os.Stdin)
				for _, prompt := range selectedItem.InputPrompts() {
					fmt.Printf("%s %s: ", ui.RenderInfo(""), prompt)

					if !scanner.Scan() {
						fmt.Println(ui.RenderWarning("Input cancelled"))
						return
					}

					input := strings.TrimSpace(scanner.Text())
					if input == "" {
						fmt.Println(ui.RenderWarning("Input cannot be empty"))
						return
					}

					inputs = append(inputs, input)
				}
			}

			// Build full command for display
			fullCmd := choice
			if len(inputs) > 0 {
				fullCmd = choice + " " + strings.Join(inputs, " ")
			}

			// Show info about the selected command
			fmt.Printf("\n%s\n\n", ui.RenderInfo(fmt.Sprintf("Running: anaphase %s", fullCmd)))

			// Find and execute the subcommand
			subCmd, _, err := cmd.Root().Find(cmdParts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding command: %v\n", err)
				return
			}

			// Ensure command output goes to stdout/stderr properly
			subCmd.SetOut(os.Stdout)
			subCmd.SetErr(os.Stderr)

			// Set args: everything after the original command parts
			args := inputs
			if originalCmdLen > 1 {
				// For nested commands like "gen domain", skip the parent command
				args = append(cmdParts[1:], inputs...)
			}

			subCmd.SetArgs(args)
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
