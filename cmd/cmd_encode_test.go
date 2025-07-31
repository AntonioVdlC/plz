package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestEncodeCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flags    map[string]string
		expected string
		wantErr  bool
	}{
		{
			name:     "base64 encode",
			args:     []string{"hello"},
			flags:    map[string]string{"type": "base64"},
			expected: "Encoded (BASE64): aGVsbG8=",
			wantErr:  false,
		},
		{
			name:     "base64 decode",
			args:     []string{"aGVsbG8="},
			flags:    map[string]string{"type": "base64", "decode": "true"},
			expected: "Decoded (BASE64): hello",
			wantErr:  false,
		},
		{
			name:     "url encode",
			args:     []string{"hello world"},
			flags:    map[string]string{"type": "url"},
			expected: "Encoded (URL): hello+world",
			wantErr:  false,
		},
		{
			name:     "url decode",
			args:     []string{"hello+world"},
			flags:    map[string]string{"type": "url", "decode": "true"},
			expected: "Decoded (URL): hello world",
			wantErr:  false,
		},
		{
			name:     "invalid base64 decode",
			args:     []string{"invalid!!!"},
			flags:    map[string]string{"type": "base64", "decode": "true"},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "unsupported encoding type",
			args:     []string{"hello"},
			flags:    map[string]string{"type": "invalid"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to default values
			encodeType = "base64"
			shouldDecode = false

			// Create a new command instance for testing
			cmd := &cobra.Command{
				Use:  "encode [string]",
				Args: cobra.ExactArgs(1),
				RunE: runEncode,
			}

			// Add flags
			cmd.Flags().StringVarP(&encodeType, "type", "t", "base64", "Encoding type: base64, url")
			cmd.Flags().BoolVarP(&shouldDecode, "decode", "d", false, "Decode instead of encode")

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

func TestEncodeValidation(t *testing.T) {
	// Test with no arguments
	cmd := &cobra.Command{
		Use:  "encode [string]",
		Args: cobra.ExactArgs(1),
		RunE: runEncode,
	}

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}

	// Test with too many arguments
	err = cmd.Args(cmd, []string{"arg1", "arg2"})
	if err == nil {
		t.Error("Expected error when too many arguments provided")
	}
}
