package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommandOutputsGeneratedCode(t *testing.T) {
	var stdout bytes.Buffer

	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)
	rootCmd.SetArgs([]string{})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute root command: %v", err)
	}

	assertValidCode(t, strings.TrimSpace(stdout.String()))
}
