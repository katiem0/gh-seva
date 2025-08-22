package secrets

import (
	"testing"
)

func TestNewCmdSecrets(t *testing.T) {
	cmd := NewCmdSecrets()

	if cmd == nil {
		t.Fatal("NewCmdSecrets() returned nil")
	}

	// Test basic properties
	if cmd.Use != "secrets <command> [flags]" {
		t.Errorf("Expected Use to be 'secrets <command> [flags]', got %s", cmd.Use)
	}

	// Test that subcommands are added
	subcommands := cmd.Commands()

	// Verify we have export and create subcommands
	exportFound := false
	createFound := false

	for _, subcmd := range subcommands {
		if subcmd.Name() == "export" {
			exportFound = true
		}
		if subcmd.Name() == "create" {
			createFound = true
		}
	}

	if !exportFound {
		t.Error("Export subcommand not found")
	}

	if !createFound {
		t.Error("Create subcommand not found")
	}

	// Test command short description
	if cmd.Short == "" {
		t.Error("Command short description should not be empty")
	}
}
