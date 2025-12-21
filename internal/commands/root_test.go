package commands

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Test that root command can be created
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	// Test command name
	if rootCmd.Use != "anaphase" {
		t.Errorf("Expected command name 'anaphase', got '%s'", rootCmd.Use)
	}

	// Test version
	if rootCmd.Version != version {
		t.Errorf("Expected version '%s', got '%s'", version, rootCmd.Version)
	}

	// Test that flags are registered
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("verbose flag should be registered")
	}

	debugFlag := rootCmd.PersistentFlags().Lookup("debug")
	if debugFlag == nil {
		t.Error("debug flag should be registered")
	}
}
