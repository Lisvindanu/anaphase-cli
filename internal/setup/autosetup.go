package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lisvindanu/anaphase-cli/internal/ui"
)

// EnsureGolangciLintConfig checks if .golangci.yml exists, creates it if not
func EnsureGolangciLintConfig() error {
	configPath := ".golangci.yml"

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // Config already exists
	}

	ui.PrintInfo("📝 Creating .golangci.yml config...")

	config := `version: 2

run:
  timeout: 5m
  tests: false
  skip-dirs:
    - vendor
    - .git
  skip-files:
    - ".*_test.go$"

linters:
  enable:
    - errcheck      # Check for unchecked errors
    - govet         # Reports suspicious constructs
    - staticcheck   # Static analysis checks
    - ineffassign   # Detects ineffectual assignments
    - unused        # Checks for unused code

issues:
  max-issues-per-linter: 0
  max-same-issues: 0

  exclude-rules:
    # Exclude errors in test files
    - path: _test\.go
      linters:
        - errcheck

    # Exclude errors in generated files
    - path: wire_gen\.go
      linters:
        - errcheck
        - staticcheck

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true
  sort-results: true
`

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("create config: %w", err)
	}

	ui.PrintSuccess("✅ Created .golangci.yml")
	return nil
}

// EnsureGolangciLint checks if golangci-lint is installed, offers to install if not
func EnsureGolangciLint() (bool, error) {
	// Check if golangci-lint is available
	if isCommandAvailable("golangci-lint") {
		return true, nil
	}

	ui.PrintWarning("⚠️  golangci-lint not installed")
	fmt.Println()

	// Offer to install
	fmt.Println(ui.RenderInfo("ℹ") + " Would you like to install it now? (recommended)")
	fmt.Print("  [Y/n]: ")

	var response string
	fmt.Scanln(&response)

	if response == "" || response == "y" || response == "Y" {
		return installGolangciLint()
	}

	ui.PrintInfo("Falling back to 'go vet' instead")
	return false, nil
}

// installGolangciLint installs golangci-lint using go install
func installGolangciLint() (bool, error) {
	ui.PrintInfo("📦 Installing golangci-lint...")
	fmt.Println()

	cmd := exec.Command("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to install: %v", err))
		fmt.Println()
		ui.PrintInfo("You can install manually:")
		fmt.Println("  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
		return false, err
	}

	ui.PrintSuccess("✅ golangci-lint installed successfully!")
	fmt.Println()
	return true, nil
}

// EnsureGoimports checks if goimports is installed, offers to install if not
func EnsureGoimports() error {
	if isCommandAvailable("goimports") {
		return nil
	}

	ui.PrintInfo("📦 Installing goimports for import organization...")

	cmd := exec.Command("go", "install", "golang.org/x/tools/cmd/goimports@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ui.PrintWarning("Failed to install goimports (optional)")
		return err
	}

	ui.PrintSuccess("✅ goimports installed")
	return nil
}

// EnsureProjectConfig ensures all necessary config files exist
func EnsureProjectConfig() error {
	// Check if this is a Go project
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return nil // Not a Go project, skip
	}

	// Create .golangci.yml if needed
	if err := EnsureGolangciLintConfig(); err != nil {
		return err
	}

	return nil
}

// EnsureGitignore adds common entries to .gitignore if they don't exist
func EnsureGitignore() error {
	gitignorePath := ".gitignore"

	// Entries to add
	entries := []string{
		"# Binaries",
		"bin/",
		"*.exe",
		"",
		"# IDE",
		".vscode/",
		".idea/",
		"*.swp",
		"*.swo",
		"",
		"# OS",
		".DS_Store",
		"Thumbs.db",
		"",
		"# Anaphase",
		".anaphase/",
	}

	// Read existing gitignore if it exists
	var existing string
	if data, err := os.ReadFile(gitignorePath); err == nil {
		existing = string(data)
	}

	// Check if entries already exist
	needsUpdate := false
	for _, entry := range entries {
		if entry != "" && !contains(existing, entry) {
			needsUpdate = true
			break
		}
	}

	if !needsUpdate {
		return nil
	}

	ui.PrintInfo("📝 Updating .gitignore...")

	// Append entries
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if existing != "" && existing[len(existing)-1] != '\n' {
		f.WriteString("\n")
	}

	for _, entry := range entries {
		f.WriteString(entry + "\n")
	}

	ui.PrintSuccess("✅ Updated .gitignore")
	return nil
}

// Helper functions

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		len(s) > len(substr)+1 && containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CreateAnaphaseConfig creates ~/.anaphase/config.yaml if it doesn't exist
func CreateAnaphaseConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".anaphase")
	configFile := filepath.Join(configDir, "config.yaml")

	// Check if config exists
	if _, err := os.Stat(configFile); err == nil {
		return nil
	}

	// Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	ui.PrintInfo("📝 Creating Anaphase config...")

	config := `# Anaphase Configuration
# Auto-generated configuration file

ai:
  primary:
    type: gemini
    # Get your API key from: https://makersuite.google.com/app/apikey
    apiKey: ""
    model: gemini-2.0-flash-exp
    timeout: 30s
    retries: 3

cache:
  enabled: true
  ttl: 24h
  dir: ~/.anaphase/cache

# To configure your API key, run:
#   anaphase config set-provider gemini
# Or set environment variable:
#   export GEMINI_API_KEY="your-key-here"
`

	if err := os.WriteFile(configFile, []byte(config), 0644); err != nil {
		return err
	}

	ui.PrintSuccess("✅ Created config at: " + configFile)
	fmt.Println()
	ui.PrintInfo("Next steps:")
	fmt.Println("  1. Get API key: https://makersuite.google.com/app/apikey")
	fmt.Println("  2. Set it: export GEMINI_API_KEY=\"your-key\"")
	fmt.Println("  3. Or run: anaphase config set-provider gemini")
	fmt.Println()

	return nil
}
