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

	assertValidCode(t, strings.TrimSpace(stdout.String()), 4, defaultLetters, defaultDigits)
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

	assertValidCode(t, strings.TrimSpace(stdout.String()), 5, defaultLetters, defaultDigits)
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

	assertValidCode(t, strings.TrimSpace(stdout.String()), 4, "AbCdEfGhIj", defaultDigits)
}

func TestRootCommandSupportsCustomDigits(t *testing.T) {
	var stdout bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--digits", "01"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute root command: %v", err)
	}

	assertValidCode(t, strings.TrimSpace(stdout.String()), 4, defaultLetters, "01")
}

func TestRootCommandOutputsVersion(t *testing.T) {
	var stdout bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute root command: %v", err)
	}

	if strings.TrimSpace(stdout.String()) != currentVersion {
		t.Fatalf("expected version %q, got %q", currentVersion, strings.TrimSpace(stdout.String()))
	}
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

func TestRootCommandRejectsDigitsWithNonDigits(t *testing.T) {
	var stderr bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--digits", "12a3"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for digits containing non-digits")
	}

	if !strings.Contains(err.Error(), "only digits 0-9 are allowed") {
		t.Fatalf("expected digit character validation error, got %v", err)
	}
}

func TestRootCommandRejectsRepeatedDigits(t *testing.T) {
	var stderr bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--digits", "2234"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for repeated digits")
	}

	if !strings.Contains(err.Error(), "digits must not repeat") {
		t.Fatalf("expected repeated digit validation error, got %v", err)
	}
}

func TestRootCommandRejectsEmptyDigits(t *testing.T) {
	var stderr bytes.Buffer
	cmd := newRootCmd()

	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--digits", ""})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for empty digits")
	}

	if !strings.Contains(err.Error(), "must contain 1 to 10 digits") {
		t.Fatalf("expected digit length validation error, got %v", err)
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
	if !strings.Contains(helpText, "Hyphen-separated Code Generator") {
		t.Fatalf("expected help to mention the tool name, got %q", helpText)
	}

	if !strings.Contains(helpText, "generates hyphen-separated codes with configurable") {
		t.Fatalf("expected help to mention the updated summary, got %q", helpText)
	}

	if !strings.Contains(helpText, "--group-size") {
		t.Fatalf("expected help to mention --group-size, got %q", helpText)
	}

	if !strings.Contains(helpText, "allowed values: 4 or 5") {
		t.Fatalf("expected help to mention allowed group sizes, got %q", helpText)
	}

	if !strings.Contains(helpText, "--letters") {
		t.Fatalf("expected help to mention --letters, got %q", helpText)
	}

	if !strings.Contains(helpText, "at least 8 unique letters") {
		t.Fatalf("expected help to mention letters constraints, got %q", helpText)
	}

	if !strings.Contains(helpText, "--digits") {
		t.Fatalf("expected help to mention --digits, got %q", helpText)
	}

	if !strings.Contains(helpText, "length 1-10") {
		t.Fatalf("expected help to mention digit constraints, got %q", helpText)
	}

	if !strings.Contains(helpText, "--version") {
		t.Fatalf("expected help to mention --version, got %q", helpText)
	}

	if !strings.Contains(helpText, "Examples:") {
		t.Fatalf("expected help to include examples, got %q", helpText)
	}

	if strings.Contains(helpText, "--toggle") {
		t.Fatalf("did not expect help to mention --toggle, got %q", helpText)
	}
}
