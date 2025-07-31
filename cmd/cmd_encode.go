package cmd

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
	Use:   "encode [string]",
	Short: "Encode/decode strings with various formats",
	Long:  `Encode or decode strings using base64, URL encoding, or other formats.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEncode,
}

var (
	encodeType   string
	shouldDecode bool
)

func init() {
	encodeCmd.Flags().StringVarP(&encodeType, "type", "t", "base64", "Encoding type: base64, url")
	encodeCmd.Flags().BoolVarP(&shouldDecode, "decode", "d", false, "Decode instead of encode")
	rootCmd.AddCommand(encodeCmd)
}

func runEncode(cmd *cobra.Command, args []string) error {
	input := args[0]
	var result string
	var err error

	switch strings.ToLower(encodeType) {
	case "base64":
		if shouldDecode {
			decoded, err := base64.StdEncoding.DecodeString(input)
			if err != nil {
				return fmt.Errorf("failed to decode base64: %w", err)
			}
			result = string(decoded)
		} else {
			result = base64.StdEncoding.EncodeToString([]byte(input))
		}
	case "url":
		if shouldDecode {
			result, err = url.QueryUnescape(input)
			if err != nil {
				return fmt.Errorf("failed to decode URL: %w", err)
			}
		} else {
			result = url.QueryEscape(input)
		}
	default:
		return fmt.Errorf("unsupported encoding type: %s (supported: base64, url)", encodeType)
	}

	operation := "Encoded"
	if shouldDecode {
		operation = "Decoded"
	}

	fmt.Printf("%s (%s): %s\n", operation, strings.ToUpper(encodeType), result)
	return nil
}
