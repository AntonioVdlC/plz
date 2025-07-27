package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time [timestamp]",
	Short: "Convert and format timestamps",
	Long:  `Convert between Unix timestamps and human-readable dates, or get current time.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runTime,
}

var (
	timeFormat   string
	toTimestamp  bool
	timeTimezone string
)

func init() {
	timeCmd.Flags().StringVarP(&timeFormat, "format", "f", "2006-01-02 15:04:05", "Output format for time")
	timeCmd.Flags().BoolVarP(&toTimestamp, "timestamp", "t", false, "Convert to Unix timestamp instead")
	timeCmd.Flags().StringVarP(&timeTimezone, "timezone", "z", "UTC", "Timezone (UTC, Local, or IANA name)")
	rootCmd.AddCommand(timeCmd)
}

func runTime(cmd *cobra.Command, args []string) error {
	var targetTime time.Time
	var err error

	location, err := parseTimezone(timeTimezone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}

	if len(args) == 0 {
		targetTime = time.Now().In(location)
	} else {
		input := args[0]
		
		if timestamp, err := strconv.ParseInt(input, 10, 64); err == nil {
			targetTime = time.Unix(timestamp, 0).In(location)
		} else {
			targetTime, err = time.Parse("2006-01-02 15:04:05", input)
			if err != nil {
				targetTime, err = time.Parse("2006-01-02", input)
				if err != nil {
					return fmt.Errorf("failed to parse time: %w", err)
				}
			}
			targetTime = targetTime.In(location)
		}
	}

	if toTimestamp {
		fmt.Printf("Unix timestamp: %d\n", targetTime.Unix())
	} else {
		fmt.Printf("Formatted time (%s): %s\n", timeTimezone, targetTime.Format(timeFormat))
		fmt.Printf("Unix timestamp: %d\n", targetTime.Unix())
		fmt.Printf("ISO 8601: %s\n", targetTime.Format(time.RFC3339))
		fmt.Printf("Weekday: %s\n", targetTime.Weekday())
	}

	return nil
}

func parseTimezone(tz string) (*time.Location, error) {
	switch strings.ToLower(tz) {
	case "utc":
		return time.UTC, nil
	case "local":
		return time.Local, nil
	default:
		return time.LoadLocation(tz)
	}
}