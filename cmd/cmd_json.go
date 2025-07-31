package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json [json-string|file]",
	Short: "Process JSON data",
	Long:  `Pretty print, minify, or validate JSON data from string or file.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runJSON,
}

var (
	jsonMinify   bool
	jsonValidate bool
	jsonFromFile bool
)

func init() {
	jsonCmd.Flags().BoolVarP(&jsonMinify, "minify", "m", false, "Minify JSON instead of pretty printing")
	jsonCmd.Flags().BoolVarP(&jsonValidate, "validate", "v", false, "Only validate JSON, don't output")
	jsonCmd.Flags().BoolVarP(&jsonFromFile, "file", "f", false, "Read JSON from file instead of string")
	rootCmd.AddCommand(jsonCmd)
}

func runJSON(cmd *cobra.Command, args []string) error {
	input := args[0]
	var jsonData []byte
	var err error

	if jsonFromFile {
		jsonData, err = os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", input, err)
		}
	} else {
		jsonData = []byte(input)
	}

	var data interface{}
	if unmarshalErr := json.Unmarshal(jsonData, &data); unmarshalErr != nil {
		return fmt.Errorf("invalid JSON: %w", unmarshalErr)
	}

	if jsonValidate {
		if jsonFromFile {
			fmt.Printf("✓ Valid JSON file: %s\n", input)
		} else {
			fmt.Println("✓ Valid JSON")
		}
		return nil
	}

	var output []byte
	if jsonMinify {
		output, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	} else {
		output, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	operation := "Pretty printed"
	if jsonMinify {
		operation = "Minified"
	}

	if jsonFromFile {
		fmt.Printf("%s JSON from %s:\n", operation, input)
	} else {
		fmt.Printf("%s JSON:\n", operation)
	}
	fmt.Println(string(output))

	return nil
}
