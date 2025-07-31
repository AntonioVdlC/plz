package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash [string|file]",
	Short: "Generate hash values for strings or files",
	Long:  `Generate MD5, SHA1, or SHA256 hash values for input strings or files.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runHash,
}

var (
	hashType     string
	hashFromFile bool
)

func init() {
	hashCmd.Flags().StringVarP(&hashType, "type", "t", "sha256", "Hash type: md5, sha1, sha256")
	hashCmd.Flags().BoolVarP(&hashFromFile, "file", "f", false, "Hash file contents instead of string")
	rootCmd.AddCommand(hashCmd)
}

func runHash(cmd *cobra.Command, args []string) error {
	input := args[0]
	var data []byte
	var err error

	if hashFromFile {
		data, err = os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", input, err)
		}
	} else {
		data = []byte(input)
	}

	var hash string
	switch strings.ToLower(hashType) {
	case "md5":
		hash = fmt.Sprintf("%x", md5.Sum(data))
	case "sha1":
		hash = fmt.Sprintf("%x", sha1.Sum(data))
	case "sha256":
		hash = fmt.Sprintf("%x", sha256.Sum256(data))
	default:
		return fmt.Errorf("unsupported hash type: %s (supported: md5, sha1, sha256)", hashType)
	}

	if hashFromFile {
		fmt.Printf("%s (%s): %s\n", strings.ToUpper(hashType), input, hash)
	} else {
		fmt.Printf("%s: %s\n", strings.ToUpper(hashType), hash)
	}

	return nil
}
