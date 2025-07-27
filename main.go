package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plz",
	Short: "A collection of useful CLI utilities",
	Long:  `plz is a CLI tool that provides a collection of small, useful utilities for everyday development tasks.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}