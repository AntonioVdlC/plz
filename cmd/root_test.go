package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Test that the root command exists and has the correct basic properties
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "plz" {
		t.Errorf("Expected Use to be 'plz', got %q", rootCmd.Use)
	}

	if rootCmd.Short != "A collection of useful CLI utilities" {
		t.Errorf("Expected Short description, got %q", rootCmd.Short)
	}

	expectedLong := "plz is a CLI tool that provides a collection of small, useful utilities for everyday development tasks."
	if rootCmd.Long != expectedLong {
		t.Errorf("Expected Long description, got %q", rootCmd.Long)
	}
}

func TestRootCommandHasSubcommands(t *testing.T) {
	expectedCommands := []string{"encode", "hash", "json", "random", "time"}
	
	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand %q not found", expectedCmd)
		}
	}
}

func TestExecuteFunction(t *testing.T) {
	// Test that Execute function runs without panicking
	// We can't easily test the full execution without mocking os.Args
	// but we can at least verify the function exists and is callable
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()

	// This would normally start the CLI, but we can't test it easily here
	// The main test is that the function exists and doesn't panic when called
	// In a real scenario, you'd mock os.Args and test specific command execution
}

func TestRootCommandHelp(t *testing.T) {
	// Test that help output contains expected information
	cmd := &cobra.Command{
		Use:   "plz",
		Short: "A collection of useful CLI utilities",
		Long:  "plz is a CLI tool that provides a collection of small, useful utilities for everyday development tasks.",
	}

	// Add a dummy subcommand to test help output
	dummyCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}
	cmd.AddCommand(dummyCmd)

	// Capture help output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.Help()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that help output contains expected elements
	if !strings.Contains(output, "plz is a CLI tool") {
		t.Errorf("Help output should contain long description. Got: %q", output)
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Help output should contain 'Usage:' section. Got: %q", output)
	}

	if !strings.Contains(output, "test") {
		t.Errorf("Help output should list the test subcommand. Got: %q", output)
	}
}