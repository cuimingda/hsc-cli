package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMakeInstallTarget(t *testing.T) {
	output := runMakeDryRun(t, "install")

	if !strings.Contains(output, "go install ./cmd/hsc") {
		t.Fatalf("expected install target to run %q, got %q", "go install ./cmd/hsc", output)
	}
}

func TestMakeTestTarget(t *testing.T) {
	output := runMakeDryRun(t, "test")

	if !strings.Contains(output, "go test ./...") {
		t.Fatalf("expected test target to run %q, got %q", "go test ./...", output)
	}
}

func runMakeDryRun(t *testing.T, target string) string {
	t.Helper()

	cmd := exec.Command("make", "-n", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run make -n %s: %v\n%s", target, err, output)
	}

	return string(output)
}
