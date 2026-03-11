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
	groupSize := defaultGroupSize

	cmd := &cobra.Command{
		Use:   "hsc",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRootCommand(cmd, groupSize)
		},
	}

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.Flags().IntVar(&groupSize, "group-size", defaultGroupSize, "number of characters in each group (allowed values: 4 or 5)")

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

func runRootCommand(cmd *cobra.Command, groupSize int) error {
	generator, err := NewCodeGenerator(nil, groupSize)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), generator.Generate())
	return err
}
