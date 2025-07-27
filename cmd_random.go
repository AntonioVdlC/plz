package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Generate random values",
	Long:  `Generate random strings, numbers, or UUIDs.`,
	RunE:  runRandom,
}

var (
	randomType   string
	randomLength int
	randomMin    int64
	randomMax    int64
)

func init() {
	randomCmd.Flags().StringVarP(&randomType, "type", "t", "string", "Type: string, number, uuid")
	randomCmd.Flags().IntVarP(&randomLength, "length", "l", 16, "Length for string generation")
	randomCmd.Flags().Int64Var(&randomMin, "min", 0, "Minimum value for number generation")
	randomCmd.Flags().Int64Var(&randomMax, "max", 100, "Maximum value for number generation")
	rootCmd.AddCommand(randomCmd)
}

func runRandom(cmd *cobra.Command, args []string) error {
	switch strings.ToLower(randomType) {
	case "string":
		result, err := generateRandomString(randomLength)
		if err != nil {
			return fmt.Errorf("failed to generate random string: %w", err)
		}
		fmt.Printf("Random string (%d chars): %s\n", randomLength, result)
	case "number":
		if randomMax <= randomMin {
			return fmt.Errorf("max value must be greater than min value")
		}
		result, err := generateRandomNumber(randomMin, randomMax)
		if err != nil {
			return fmt.Errorf("failed to generate random number: %w", err)
		}
		fmt.Printf("Random number (%d-%d): %d\n", randomMin, randomMax, result)
	case "uuid":
		result, err := generateUUID()
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}
		fmt.Printf("Random UUID: %s\n", result)
	default:
		return fmt.Errorf("unsupported random type: %s (supported: string, number, uuid)", randomType)
	}

	return nil
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

func generateRandomNumber(min, max int64) (int64, error) {
	diff := max - min
	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, err
	}
	return min + n.Int64(), nil
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}