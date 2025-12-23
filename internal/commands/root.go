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
	version = "0.4.3"
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
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running menu: %v\n", err)
		os.Exit(1)
	}

	// Clear screen after menu exits
	fmt.Print("\033[2J\033[H")

	// Get the selected command
	if menuModel, ok := finalModel.(ui.MenuModel); ok {
		choice := menuModel.GetChoice()
		selectedItem := menuModel.GetSelectedItem()

		if choice != "" && selectedItem != nil {
			// Parse command parts
			cmdParts := ui.FormatCommand(choice)

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
			fmt.Printf("\n%s\n", ui.RenderInfo(fmt.Sprintf("Running: anaphase %s", fullCmd)))

			// Find and execute the subcommand
			subCmd, _, err := cmd.Root().Find(cmdParts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\n%s Error finding command: %v\n\n", ui.RenderError(""), err)
				return
			}

			// Set args: only the user inputs
			// The subcommand is already found, we just need to pass the inputs
			args := inputs

			// Handle special cases for commands with flags
			if choice == "init" && len(inputs) >= 2 {
				// For init: inputs[0] = project name, inputs[1] = database type
				projectName := inputs[0]
				dbType := strings.ToLower(strings.TrimSpace(inputs[1]))

				// Default to postgres if empty
				if dbType == "" {
					dbType = "postgres"
				}

				// Validate database type
				validDBs := map[string]bool{"postgres": true, "mysql": true, "sqlite": true, "mongodb": true}
				if !validDBs[dbType] {
					fmt.Fprintf(os.Stderr, "\n%s Invalid database type '%s'. Valid types: postgres, mysql, sqlite, mongodb\n\n", ui.RenderError(""), dbType)
					return
				}

				// Set args: project name only, database as flag
				args = []string{projectName}
				subCmd.Flags().Set("db", dbType)
			}

			// Ensure command output goes to stdout/stderr properly
			subCmd.SetOut(os.Stdout)
			subCmd.SetErr(os.Stderr)

			// Call RunE directly instead of Execute to avoid triggering root command
			if subCmd.RunE != nil {
				if err := subCmd.RunE(subCmd, args); err != nil {
					fmt.Fprintf(os.Stderr, "\n%s Command failed: %v\n\n", ui.RenderError(""), err)
					os.Exit(1)
				}
			} else if subCmd.Run != nil {
				subCmd.Run(subCmd, args)
			} else {
				// Command has subcommands - show help
				if subCmd.HasSubCommands() {
					fmt.Println()
					ui.PrintInfo(fmt.Sprintf("'%s' has subcommands. Available subcommands:", choice))
					fmt.Println()
					for _, subcmd := range subCmd.Commands() {
						if !subcmd.Hidden {
							fmt.Printf("  • %s - %s\n", subcmd.Name(), subcmd.Short)
						}
					}
					fmt.Println()
					ui.PrintInfo(fmt.Sprintf("Run: anaphase %s <subcommand>", choice))
					fmt.Println()
				} else {
					fmt.Fprintf(os.Stderr, "\n%s Command has no run function\n\n", ui.RenderError(""))
				}
				return
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
