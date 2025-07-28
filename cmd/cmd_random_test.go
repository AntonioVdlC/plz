package cmd

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRandomCommand(t *testing.T) {
	tests := []struct {
		name     string
		flags    map[string]string
		validate func(string) bool
		wantErr  bool
	}{
		{
			name:  "random string default length",
			flags: map[string]string{"type": "string"},
			validate: func(output string) bool {
				// Should have format: "Random string (16 chars): <string>"
				parts := strings.Split(output, ": ")
				if len(parts) != 2 {
					return false
				}
				if !strings.HasPrefix(parts[0], "Random string (16 chars)") {
					return false
				}
				return len(parts[1]) == 16
			},
			wantErr: false,
		},
		{
			name:  "random string custom length",
			flags: map[string]string{"type": "string", "length": "10"},
			validate: func(output string) bool {
				parts := strings.Split(output, ": ")
				if len(parts) != 2 {
					return false
				}
				if !strings.HasPrefix(parts[0], "Random string (10 chars)") {
					return false
				}
				return len(parts[1]) == 10
			},
			wantErr: false,
		},
		{
			name:  "random number default range",
			flags: map[string]string{"type": "number"},
			validate: func(output string) bool {
				// Should have format: "Random number (0-100): <number>"
				parts := strings.Split(output, ": ")
				if len(parts) != 2 {
					return false
				}
				if !strings.HasPrefix(parts[0], "Random number (0-100)") {
					return false
				}
				num, err := strconv.Atoi(parts[1])
				return err == nil && num >= 0 && num <= 100
			},
			wantErr: false,
		},
		{
			name:  "random number custom range",
			flags: map[string]string{"type": "number", "min": "50", "max": "60"},
			validate: func(output string) bool {
				parts := strings.Split(output, ": ")
				if len(parts) != 2 {
					return false
				}
				if !strings.HasPrefix(parts[0], "Random number (50-60)") {
					return false
				}
				num, err := strconv.Atoi(parts[1])
				return err == nil && num >= 50 && num <= 60
			},
			wantErr: false,
		},
		{
			name:  "random uuid",
			flags: map[string]string{"type": "uuid"},
			validate: func(output string) bool {
				// Should have format: "Random UUID: <uuid>"
				parts := strings.Split(output, ": ")
				if len(parts) != 2 {
					return false
				}
				if parts[0] != "Random UUID" {
					return false
				}
				// Check UUID format: 8-4-4-4-12
				uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
				return uuidRegex.MatchString(parts[1])
			},
			wantErr: false,
		},
		{
			name:     "invalid random type",
			flags:    map[string]string{"type": "invalid"},
			validate: func(output string) bool { return true },
			wantErr:  true,
		},
		{
			name:     "invalid number range",
			flags:    map[string]string{"type": "number", "min": "100", "max": "50"},
			validate: func(output string) bool { return true },
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to default values
			randomType = "string"
			randomLength = 16
			randomMin = 0
			randomMax = 100

			// Create a new command instance for testing
			cmd := &cobra.Command{
				Use:  "random",
				RunE: runRandom,
			}

			// Add flags
			cmd.Flags().StringVarP(&randomType, "type", "t", "string", "Type: string, number, uuid")
			cmd.Flags().IntVarP(&randomLength, "length", "l", 16, "Length for string generation")
			cmd.Flags().Int64Var(&randomMin, "min", 0, "Minimum value for number generation")
			cmd.Flags().Int64Var(&randomMax, "max", 100, "Maximum value for number generation")

			// Set flags
			for flag, value := range tt.flags {
				cmd.Flags().Set(flag, value)
			}

			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run command
			err := cmd.RunE(cmd, []string{})

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
				if !tt.validate(output) {
					t.Errorf("Output validation failed for: %q", output)
				}
			}
		})
	}
}

func TestRandomFunctions(t *testing.T) {
	t.Run("generateRandomString", func(t *testing.T) {
		result, err := generateRandomString(10)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(result) != 10 {
			t.Errorf("Expected length 10, got %d", len(result))
		}
		// Check that it contains valid characters
		validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		for _, char := range result {
			if !strings.ContainsRune(validChars, char) {
				t.Errorf("Invalid character found: %c", char)
			}
		}
	})

	t.Run("generateRandomNumber", func(t *testing.T) {
		result, err := generateRandomNumber(10, 20)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result < 10 || result >= 20 {
			t.Errorf("Result %d not in range [10, 20)", result)
		}
	})

	t.Run("generateUUID", func(t *testing.T) {
		result, err := generateUUID()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Check UUID format
		uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
		if !uuidRegex.MatchString(result) {
			t.Errorf("Invalid UUID format: %s", result)
		}
	})
}