package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommandOutputsGeneratedCode(t *testing.T) {
	var stdout bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute root command: %v", err)
	}

	assertValidCode(t, strings.TrimSpace(stdout.String()), 4, defaultLetters)
}

func TestRootCommandSupportsGroupSizeFive(t *testing.T) {
	var stdout bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--group-size", "5"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute root command: %v", err)
	}

	assertValidCode(t, strings.TrimSpace(stdout.String()), 5, defaultLetters)
}

func TestRootCommandSupportsCustomLetters(t *testing.T) {
	var stdout bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--letters", "AbCdEfGhIj"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute root command: %v", err)
	}

	assertValidCode(t, strings.TrimSpace(stdout.String()), 4, "AbCdEfGhIj")
}

func TestRootCommandRejectsInvalidGroupSize(t *testing.T) {
	var stderr bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--group-size", "3"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid group size")
	}

	if !strings.Contains(err.Error(), "allowed values are 4 or 5") {
		t.Fatalf("expected group size validation error, got %v", err)
	}
}

func TestRootCommandRejectsLettersWithNonLetters(t *testing.T) {
	var stderr bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--letters", "abcd1234"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for letters containing non-letters")
	}

	if !strings.Contains(err.Error(), "only letters are allowed") {
		t.Fatalf("expected letters validation error, got %v", err)
	}
}

func TestRootCommandRejectsLettersThatAreTooShortAfterDeduplication(t *testing.T) {
	var stderr bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--letters", "AaBbCcDd"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for too few unique letters")
	}

	if !strings.Contains(err.Error(), "need at least 8 unique letters") {
		t.Fatalf("expected unique letters validation error, got %v", err)
	}
}

func TestRootCommandHelpMentionsFlagConstraints(t *testing.T) {
	var stdout bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute help: %v", err)
	}

	helpText := stdout.String()
	if !strings.Contains(helpText, "--group-size") {
		t.Fatalf("expected help to mention --group-size, got %q", helpText)
	}

	if !strings.Contains(helpText, "allowed values: 4 or 5") {
		t.Fatalf("expected help to mention allowed group sizes, got %q", helpText)
	}

	if !strings.Contains(helpText, "--letters") {
		t.Fatalf("expected help to mention --letters, got %q", helpText)
	}

	if !strings.Contains(helpText, "need at least 8 unique letters") {
		t.Fatalf("expected help to mention letters constraints, got %q", helpText)
	}
}
