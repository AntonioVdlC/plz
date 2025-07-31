package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestHashCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flags    map[string]string
		expected string
		wantErr  bool
	}{
		{
			name:     "sha256 string",
			args:     []string{"hello"},
			flags:    map[string]string{"type": "sha256"},
			expected: "SHA256: 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			wantErr:  false,
		},
		{
			name:     "md5 string",
			args:     []string{"hello"},
			flags:    map[string]string{"type": "md5"},
			expected: "MD5: 5d41402abc4b2a76b9719d911017c592",
			wantErr:  false,
		},
		{
			name:     "sha1 string",
			args:     []string{"hello"},
			flags:    map[string]string{"type": "sha1"},
			expected: "SHA1: aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
			wantErr:  false,
		},
		{
			name:     "unsupported hash type",
			args:     []string{"hello"},
			flags:    map[string]string{"type": "invalid"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to default values
			hashType = "sha256"
			hashFromFile = false

			// Create a new command instance for testing
			cmd := &cobra.Command{
				Use:  "hash [string|file]",
				Args: cobra.ExactArgs(1),
				RunE: runHash,
			}

			// Add flags
			cmd.Flags().StringVarP(&hashType, "type", "t", "sha256", "Hash type: md5, sha1, sha256")
			cmd.Flags().BoolVarP(&hashFromFile, "file", "f", false, "Hash file contents instead of string")

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

func TestHashFromFile(t *testing.T) {
	// Reset flags to default values
	hashType = "sha256"
	hashFromFile = false

	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use:  "hash [string|file]",
		Args: cobra.ExactArgs(1),
		RunE: runHash,
	}

	// Add flags
	cmd.Flags().StringVarP(&hashType, "type", "t", "sha256", "Hash type: md5, sha1, sha256")
	cmd.Flags().BoolVarP(&hashFromFile, "file", "f", false, "Hash file contents instead of string")

	// Set flags for file reading
	cmd.Flags().Set("file", "true")

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command with test file
	err := cmd.RunE(cmd, []string{"testdata/test.txt"})

	// Restore stdout and get output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "SHA256 (testdata/test.txt): a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestHashFileNotFound(t *testing.T) {
	// Reset flags to default values
	hashType = "sha256"
	hashFromFile = false

	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use:  "hash [string|file]",
		Args: cobra.ExactArgs(1),
		RunE: runHash,
	}

	// Add flags
	cmd.Flags().StringVarP(&hashType, "type", "t", "sha256", "Hash type: md5, sha1, sha256")
	cmd.Flags().BoolVarP(&hashFromFile, "file", "f", false, "Hash file contents instead of string")

	// Set flags for file reading
	cmd.Flags().Set("file", "true")

	// Run command with non-existent file
	err := cmd.RunE(cmd, []string{"nonexistent.txt"})

	if err == nil {
		t.Error("Expected error when file doesn't exist")
	}
}
