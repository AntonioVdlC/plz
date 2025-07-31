package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plz",
	Short: "A collection of useful CLI utilities",
	Long:  `plz is a CLI tool that provides a collection of small, useful utilities for everyday development tasks.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
