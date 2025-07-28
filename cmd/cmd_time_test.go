package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestTimeCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		flags      map[string]string
		checkFunc  func(string) bool
		wantErr    bool
	}{
		{
			name:  "current time default",
			args:  []string{},
			flags: map[string]string{},
			checkFunc: func(output string) bool {
				lines := strings.Split(output, "\n")
				return len(lines) >= 4 &&
					strings.Contains(lines[0], "Formatted time (UTC):") &&
					strings.Contains(lines[1], "Unix timestamp:") &&
					strings.Contains(lines[2], "ISO 8601:") &&
					strings.Contains(lines[3], "Weekday:")
			},
			wantErr: false,
		},
		{
			name:  "unix timestamp input",
			args:  []string{"1609459200"},
			flags: map[string]string{},
			checkFunc: func(output string) bool {
				lines := strings.Split(output, "\n")
				return len(lines) >= 4 &&
					strings.Contains(lines[0], "Formatted time (UTC): 2021-01-01 00:00:00") &&
					strings.Contains(lines[1], "Unix timestamp: 1609459200") &&
					strings.Contains(lines[2], "ISO 8601: 2021-01-01T00:00:00Z") &&
					strings.Contains(lines[3], "Weekday: Friday")
			},
			wantErr: false,
		},
		{
			name:  "date string input",
			args:  []string{"2021-01-01"},
			flags: map[string]string{},
			checkFunc: func(output string) bool {
				lines := strings.Split(output, "\n")
				return len(lines) >= 4 &&
					strings.Contains(lines[0], "Formatted time (UTC): 2021-01-01 00:00:00")
			},
			wantErr: false,
		},
		{
			name:  "datetime string input",
			args:  []string{"2021-01-01 12:30:45"},
			flags: map[string]string{},
			checkFunc: func(output string) bool {
				lines := strings.Split(output, "\n")
				return len(lines) >= 4 &&
					strings.Contains(lines[0], "Formatted time (UTC): 2021-01-01 12:30:45")
			},
			wantErr: false,
		},
		{
			name:  "timestamp output",
			args:  []string{"2021-01-01"},
			flags: map[string]string{"timestamp": "true"},
			checkFunc: func(output string) bool {
				return strings.TrimSpace(output) == "Unix timestamp: 1609459200"
			},
			wantErr: false,
		},
		{
			name:  "custom format",
			args:  []string{"2021-01-01"},
			flags: map[string]string{"format": "2006/01/02"},
			checkFunc: func(output string) bool {
				lines := strings.Split(output, "\n")
				return len(lines) >= 1 &&
					strings.Contains(lines[0], "Formatted time (UTC): 2021/01/01")
			},
			wantErr: false,
		},
		{
			name:     "invalid date format",
			args:     []string{"invalid-date"},
			flags:    map[string]string{},
			checkFunc: func(output string) bool { return true },
			wantErr:  true,
		},
		{
			name:     "invalid timezone",
			args:     []string{"2021-01-01"},
			flags:    map[string]string{"timezone": "Invalid/Zone"},
			checkFunc: func(output string) bool { return true },
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to default values
			timeFormat = "2006-01-02 15:04:05"
			toTimestamp = false
			timeTimezone = "UTC"

			// Create a new command instance for testing
			cmd := &cobra.Command{
				Use:  "time [timestamp]",
				Args: cobra.MaximumNArgs(1),
				RunE: runTime,
			}

			// Add flags
			cmd.Flags().StringVarP(&timeFormat, "format", "f", "2006-01-02 15:04:05", "Output format for time")
			cmd.Flags().BoolVarP(&toTimestamp, "timestamp", "t", false, "Convert to Unix timestamp instead")
			cmd.Flags().StringVarP(&timeTimezone, "timezone", "z", "UTC", "Timezone (UTC, Local, or IANA name)")

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
				if !tt.checkFunc(output) {
					t.Errorf("Output validation failed for: %q", output)
				}
			}
		})
	}
}

func TestParseTimezone(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{
			name:     "UTC timezone",
			timezone: "UTC",
			wantErr:  false,
		},
		{
			name:     "local timezone",
			timezone: "local",
			wantErr:  false,
		},
		{
			name:     "case insensitive UTC",
			timezone: "utc",
			wantErr:  false,
		},
		{
			name:     "case insensitive local",
			timezone: "LOCAL",
			wantErr:  false,
		},
		{
			name:     "valid IANA timezone",
			timezone: "America/New_York",
			wantErr:  false,
		},
		{
			name:     "invalid timezone",
			timezone: "Invalid/Zone",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := parseTimezone(tt.timezone)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if location == nil {
					t.Error("Expected valid location, got nil")
				}
			}
		})
	}
}

func TestTimeCommandArgValidation(t *testing.T) {
	// Test with too many arguments
	cmd := &cobra.Command{
		Use:  "time [timestamp]",
		Args: cobra.MaximumNArgs(1),
		RunE: runTime,
	}

	err := cmd.Args(cmd, []string{"arg1", "arg2"})
	if err == nil {
		t.Error("Expected error when too many arguments provided")
	}
}