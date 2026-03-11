/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	digits := defaultDigits
	groupSize := defaultGroupSize
	letters := defaultLetters

	cmd := &cobra.Command{
		Use:   "hsc",
		Short: "Generate hyphen-separated codes with configurable letters, digits, and group size",
		Long: `Hyphen-separated Code Generator generates hyphen-separated codes with configurable
letters, digits, and group size.

Each group always contains exactly 2 letters and the remaining characters are digits.
The first character of the first group is always a letter, and each generated letter
is used at most once per code.`,
		Example: `  hsc
  hsc --group-size 5
  hsc --letters AbCdEfGhIj
  hsc --digits 0123456789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRootCommand(cmd, groupSize, letters, digits)
		},
	}

	cmd.Flags().StringVar(&digits, "digits", defaultDigits, "candidate digits for generated code (digits only, no duplicates, length 1-10)")
	cmd.Flags().IntVar(&groupSize, "group-size", defaultGroupSize, "characters per group (allowed values: 4 or 5)")
	cmd.Flags().StringVar(&letters, "letters", defaultLetters, "candidate letters for generated code (letters only, case-insensitive deduplication, at least 8 unique letters)")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func runRootCommand(cmd *cobra.Command, groupSize int, letters string, digits string) error {
	generator, err := NewCodeGenerator(nil, groupSize, letters, digits)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), generator.Generate())
	return err
}
