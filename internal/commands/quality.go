package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/spf13/cobra"
)

var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "Code quality tools",
	Long: `Code quality tools for linting, formatting, and validating generated code.

Available subcommands:
  lint      - Run linters on generated code
  format    - Format code using gofmt/goimports
  validate  - Validate code structure and build`,
}

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Run linters on code",
	Long: `Run Go linters on the codebase.

This command will:
  1. Try to use golangci-lint if available
  2. Fall back to go vet if golangci-lint is not installed
  3. Report any issues found

Example:
  anaphase quality lint
  anaphase quality lint ./internal/core
  anaphase quality lint --fix`,
	RunE: runLint,
}

var formatCmd = &cobra.Command{
	Use:   "format [path]",
	Short: "Format code",
	Long: `Format Go code using gofmt and goimports.

This command will:
  1. Format code using gofmt
  2. Organize imports using goimports if available
  3. Show which files were formatted

Example:
  anaphase quality format
  anaphase quality format ./internal/core`,
	RunE: runFormat,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate code structure",
	Long: `Validate that generated code compiles and has correct structure.

This command will:
  1. Check that all Go files compile
  2. Run go build to ensure no compilation errors
  3. Verify package structure

Example:
  anaphase quality validate`,
	RunE: runValidate,
}

var (
	lintFix     bool
	formatWrite bool
)

func init() {
	rootCmd.AddCommand(qualityCmd)
	qualityCmd.AddCommand(lintCmd)
	qualityCmd.AddCommand(formatCmd)
	qualityCmd.AddCommand(validateCmd)

	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "Automatically fix issues when possible")
	formatCmd.Flags().BoolVarP(&formatWrite, "write", "w", true, "Write result to source file instead of stdout")
}

func runLint(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderTitle("Code Linting"))

	// Determine path
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	ui.PrintInfo(fmt.Sprintf("Linting: %s", path))
	fmt.Println()

	// Try golangci-lint first
	if isCommandAvailable("golangci-lint") {
		fmt.Println("📋 Running golangci-lint...")
		cmdArgs := []string{"run", path}
		if lintFix {
			cmdArgs = append(cmdArgs, "--fix")
		}

		output, err := runCommand("golangci-lint", cmdArgs...)
		if err != nil {
			if output != "" {
				fmt.Println(output)
			}
			ui.PrintWarning("Linting found issues")
			fmt.Println()
			ui.PrintInfo("Run with --fix to automatically fix some issues")
			return nil // Don't return error, just show issues
		}

		ui.PrintSuccess("No linting issues found!")
	} else {
		// Fall back to go vet
		fmt.Println("📋 Running go vet...")
		ui.PrintWarning("golangci-lint not found, using go vet instead")
		ui.PrintInfo("Install golangci-lint for better linting: https://golangci-lint.run/usage/install/")
		fmt.Println()

		output, err := runCommand("go", "vet", path)
		if err != nil {
			if output != "" {
				fmt.Println(output)
			}
			ui.PrintWarning("go vet found issues")
			return nil
		}

		ui.PrintSuccess("No issues found by go vet!")
	}

	fmt.Println()
	return nil
}

func runFormat(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderTitle("Code Formatting"))

	// Determine path
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	ui.PrintInfo(fmt.Sprintf("Formatting: %s", path))
	fmt.Println()

	// Run gofmt
	fmt.Println("📝 Running gofmt...")
	gofmtArgs := []string{"-s"}
	if formatWrite {
		gofmtArgs = append(gofmtArgs, "-w")
	} else {
		gofmtArgs = append(gofmtArgs, "-d")
	}

	// Find all .go files
	goFiles, err := findGoFiles(path)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to find Go files: %v", err))
		return err
	}

	if len(goFiles) == 0 {
		ui.PrintWarning("No Go files found")
		return nil
	}

	formattedCount := 0
	for _, file := range goFiles {
		output, err := runCommand("gofmt", append(gofmtArgs, file)...)
		if err != nil {
			ui.PrintError(fmt.Sprintf("Failed to format %s: %v", file, err))
			continue
		}

		// If output is not empty, file was formatted
		if strings.TrimSpace(output) != "" {
			formattedCount++
			if !formatWrite {
				fmt.Println(output)
			}
		}
	}

	if formatWrite {
		if formattedCount > 0 {
			ui.PrintSuccess(fmt.Sprintf("Formatted %d file(s)", formattedCount))
		} else {
			ui.PrintSuccess("All files already formatted")
		}
	}

	// Try goimports if available
	if isCommandAvailable("goimports") {
		fmt.Println("\n📦 Running goimports...")
		goimportsArgs := []string{}
		if formatWrite {
			goimportsArgs = append(goimportsArgs, "-w")
		} else {
			goimportsArgs = append(goimportsArgs, "-d")
		}

		importsCount := 0
		for _, file := range goFiles {
			output, err := runCommand("goimports", append(goimportsArgs, file)...)
			if err != nil {
				continue
			}

			if strings.TrimSpace(output) != "" {
				importsCount++
				if !formatWrite {
					fmt.Println(output)
				}
			}
		}

		if formatWrite {
			if importsCount > 0 {
				ui.PrintSuccess(fmt.Sprintf("Organized imports in %d file(s)", importsCount))
			} else {
				ui.PrintSuccess("All imports already organized")
			}
		}
	} else {
		ui.PrintInfo("\nInstall goimports for automatic import organization:")
		fmt.Println("  go install golang.org/x/tools/cmd/goimports@latest")
	}

	fmt.Println()
	return nil
}

func runValidate(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.RenderTitle("Code Validation"))
	fmt.Println()

	// Step 1: Check syntax
	fmt.Println("📋 Step 1/3: Checking syntax...")
	output, err := runCommand("go", "fmt", "./...")
	if err != nil {
		ui.PrintError("Syntax errors found")
		if output != "" {
			fmt.Println(output)
		}
		return err
	}
	ui.PrintSuccess("Syntax OK")

	// Step 2: Run go vet
	fmt.Println("\n📋 Step 2/3: Running go vet...")
	output, err = runCommand("go", "vet", "./...")
	if err != nil {
		ui.PrintWarning("go vet found potential issues")
		if output != "" {
			fmt.Println(output)
		}
		// Don't fail on vet warnings
	} else {
		ui.PrintSuccess("No vet issues found")
	}

	// Step 3: Try to build
	fmt.Println("\n📋 Step 3/3: Building code...")
	output, err = runCommand("go", "build", "./...")
	if err != nil {
		ui.PrintError("Build failed")
		if output != "" {
			fmt.Println(output)
		}
		return err
	}
	ui.PrintSuccess("Build successful")

	fmt.Println()
	ui.PrintSuccess("Validation complete! Code is ready to use.")
	fmt.Println()

	return nil
}

// Helper functions

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func findGoFiles(root string) ([]string, error) {
	var goFiles []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor and hidden directories
		if info.IsDir() {
			name := info.Name()
			if name == "vendor" || name == ".git" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only include .go files (not _test.go for now)
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			goFiles = append(goFiles, path)
		}

		return nil
	})

	return goFiles, err
}
