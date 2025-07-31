package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestJSONCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flags    map[string]string
		expected string
		wantErr  bool
	}{
		{
			name:     "pretty print json string",
			args:     []string{`{"name":"test","value":123}`},
			flags:    map[string]string{},
			expected: "Pretty printed JSON:\n{\n  \"name\": \"test\",\n  \"value\": 123\n}",
			wantErr:  false,
		},
		{
			name:     "minify json string",
			args:     []string{`{"name": "test", "value": 123}`},
			flags:    map[string]string{"minify": "true"},
			expected: "Minified JSON:\n{\"name\":\"test\",\"value\":123}",
			wantErr:  false,
		},
		{
			name:     "validate json string",
			args:     []string{`{"name":"test","value":123}`},
			flags:    map[string]string{"validate": "true"},
			expected: "✓ Valid JSON",
			wantErr:  false,
		},
		{
			name:     "invalid json string",
			args:     []string{`{"name":"test","value":}`},
			flags:    map[string]string{},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "validate invalid json",
			args:     []string{`{"invalid":`},
			flags:    map[string]string{"validate": "true"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to default values
			jsonMinify = false
			jsonValidate = false
			jsonFromFile = false

			// Create a new command instance for testing
			cmd := &cobra.Command{
				Use:  "json [json-string|file]",
				Args: cobra.ExactArgs(1),
				RunE: runJSON,
			}

			// Add flags
			cmd.Flags().BoolVarP(&jsonMinify, "minify", "m", false, "Minify JSON instead of pretty printing")
			cmd.Flags().BoolVarP(&jsonValidate, "validate", "v", false, "Only validate JSON, don't output")
			cmd.Flags().BoolVarP(&jsonFromFile, "file", "f", false, "Read JSON from file instead of string")

			// Set flags
			for flag, value := range tt.flags {
				cmd.Flags().Set(flag, value)
			}

			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run command
			err := cmd.RunE(cmd, tt.args)

			// Restore stdout and get output
			w.Close()
			os.Stdout = old
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := strings.TrimSpace(buf.String())

			// Check results
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if output != tt.expected {
					t.Errorf("Expected output %q, got %q", tt.expected, output)
				}
			}
		})
	}
}

func TestJSONFromFile(t *testing.T) {
	// Reset flags to default values
	jsonMinify = false
	jsonValidate = false
	jsonFromFile = false

	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use:  "json [json-string|file]",
		Args: cobra.ExactArgs(1),
		RunE: runJSON,
	}

	// Add flags
	cmd.Flags().BoolVarP(&jsonMinify, "minify", "m", false, "Minify JSON instead of pretty printing")
	cmd.Flags().BoolVarP(&jsonValidate, "validate", "v", false, "Only validate JSON, don't output")
	cmd.Flags().BoolVarP(&jsonFromFile, "file", "f", false, "Read JSON from file instead of string")

	// Set flags for file reading
	cmd.Flags().Set("file", "true")

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command with test file
	err := cmd.RunE(cmd, []string{"testdata/test.json"})

	// Restore stdout and get output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedOutput := "Pretty printed JSON from testdata/test.json:\n{\n  \"name\": \"test\",\n  \"nested\": {\n    \"key\": \"value\"\n  },\n  \"value\": 123\n}"
	if output != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, output)
	}
}

func TestJSONFileNotFound(t *testing.T) {
	// Reset flags to default values
	jsonMinify = false
	jsonValidate = false
	jsonFromFile = false

	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use:  "json [json-string|file]",
		Args: cobra.ExactArgs(1),
		RunE: runJSON,
	}

	// Add flags
	cmd.Flags().BoolVarP(&jsonMinify, "minify", "m", false, "Minify JSON instead of pretty printing")
	cmd.Flags().BoolVarP(&jsonValidate, "validate", "v", false, "Only validate JSON, don't output")
	cmd.Flags().BoolVarP(&jsonFromFile, "file", "f", false, "Read JSON from file instead of string")

	// Set flags for file reading
	cmd.Flags().Set("file", "true")

	// Run command with non-existent file
	err := cmd.RunE(cmd, []string{"nonexistent.json"})

	if err == nil {
		t.Error("Expected error when file doesn't exist")
	}
}

func TestJSONValidateFromFile(t *testing.T) {
	// Reset flags to default values
	jsonMinify = false
	jsonValidate = false
	jsonFromFile = false

	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use:  "json [json-string|file]",
		Args: cobra.ExactArgs(1),
		RunE: runJSON,
	}

	// Add flags
	cmd.Flags().BoolVarP(&jsonMinify, "minify", "m", false, "Minify JSON instead of pretty printing")
	cmd.Flags().BoolVarP(&jsonValidate, "validate", "v", false, "Only validate JSON, don't output")
	cmd.Flags().BoolVarP(&jsonFromFile, "file", "f", false, "Read JSON from file instead of string")

	// Set flags for file reading and validation
	cmd.Flags().Set("file", "true")
	cmd.Flags().Set("validate", "true")

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command with test file
	err := cmd.RunE(cmd, []string{"testdata/test.json"})

	// Restore stdout and get output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "✓ Valid JSON file: testdata/test.json"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}
