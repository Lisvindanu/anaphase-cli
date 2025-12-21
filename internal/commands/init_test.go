package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidProjectName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid lowercase", "myproject", true},
		{"valid with hyphen", "my-project", true},
		{"valid with underscore", "my_project", true},
		{"valid with number", "my-project-123", true},
		{"invalid empty", "", false},
		{"invalid starts with number", "123project", false},
		{"invalid starts with hyphen", "-project", false},
		{"invalid special char", "my@project", false},
		{"invalid space", "my project", false},
		{"invalid dot", "my.project", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidProjectName(tt.input)
			if got != tt.want {
				t.Errorf("isValidProjectName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsLetter(t *testing.T) {
	tests := []struct {
		name  string
		input rune
		want  bool
	}{
		{"lowercase a", 'a', true},
		{"lowercase z", 'z', true},
		{"uppercase A", 'A', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', false},
		{"digit 9", '9', false},
		{"hyphen", '-', false},
		{"underscore", '_', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLetter(tt.input)
			if got != tt.want {
				t.Errorf("isLetter(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		name  string
		input rune
		want  bool
	}{
		{"digit 0", '0', true},
		{"digit 5", '5', true},
		{"digit 9", '9', true},
		{"letter a", 'a', false},
		{"letter Z", 'Z', false},
		{"hyphen", '-', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDigit(tt.input)
			if got != tt.want {
				t.Errorf("isDigit(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestInitCommand(t *testing.T) {
	// Test that init command exists and has correct name
	if initCmd.Use != "init <project_name>" {
		t.Errorf("Expected command use 'init <project_name>', got '%s'", initCmd.Use)
	}

	// Test flags are registered
	flags := []string{"module", "db", "cache", "event-bus", "no-docker"}
	for _, flag := range flags {
		if initCmd.Flags().Lookup(flag) == nil {
			t.Errorf("Flag '%s' should be registered", flag)
		}
	}
}

func TestInitCommandDefaults(t *testing.T) {
	// Test default values
	dbFlag := initCmd.Flags().Lookup("db")
	if dbFlag.DefValue != "postgres" {
		t.Errorf("Expected default db 'postgres', got '%s'", dbFlag.DefValue)
	}

	moduleFlag := initCmd.Flags().Lookup("module")
	if moduleFlag.DefValue != "" {
		t.Errorf("Expected default module '', got '%s'", moduleFlag.DefValue)
	}
}

// Integration test - actually create a project
func TestInitCommandIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "anaphase-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Run init command
	rootCmd.SetArgs([]string{"init", "test-project", "--module", "github.com/test/test-project"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Init command failed: %v", err)
	}

	// Verify project structure
	projectDir := filepath.Join(tmpDir, "test-project")

	// Check key files exist
	files := []string{
		"go.mod",
		"Makefile",
		"README.md",
		".env.example",
		".gitignore",
		"Dockerfile",
		"docker-compose.yml",
		"cmd/api/main.go",
		"internal/config/config.go",
		"internal/server/server.go",
	}

	for _, file := range files {
		path := filepath.Join(projectDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file '%s' does not exist", file)
		}
	}

	// Check go.mod contains correct module name
	content, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	if !contains(string(content), "github.com/test/test-project") {
		t.Error("go.mod does not contain correct module name")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
