package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
	"github.com/spf13/cobra"
)

func TestLoadGenerateConfigReadsOverwrite(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, ".infra-doctor")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := []byte("output:\n  directory: generated\n  overwrite: true\n")
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), content, 0o644); err != nil {
		t.Fatal(err)
	}

	loaded, err := loadGenerateConfig(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.OutputDir != "generated" || !loaded.Overwrite {
		t.Fatalf("loaded config = %+v", loaded)
	}
}

func TestLoadGenerateConfigReturnsZeroValueWhenNoConfigFile(t *testing.T) {
	root := t.TempDir()

	loaded, err := loadGenerateConfig(root, "")
	if err != nil {
		t.Fatalf("loadGenerateConfig with no config file: want nil error, got %v", err)
	}
	if loaded != (loadedGenerateConfig{}) {
		t.Errorf("loadGenerateConfig with no config file = %+v, want zero value", loaded)
	}
}

func TestLoadGenerateConfigReturnsErrorWhenExplicitPathMissing(t *testing.T) {
	root := t.TempDir()

	_, err := loadGenerateConfig(root, filepath.Join(root, "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("loadGenerateConfig with missing explicit path: want error, got nil")
	}
}

func TestLoadGenerateConfigReturnsErrorOnInvalidYAML(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, ".infra-doctor")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("not: [valid: yaml"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadGenerateConfig(root, ""); err == nil {
		t.Fatal("loadGenerateConfig with invalid YAML: want error, got nil")
	}
}

func TestPrintGenerateResultWritesWarningsAndPlannedPaths(t *testing.T) {
	result := generate.Result{
		DryRun:   true,
		Warnings: []string{"port not detected"},
		Planned:  []string{"Dockerfile", "docker-compose.yml"},
		Skipped:  []string{"docker-compose.yml"},
	}

	command := &cobra.Command{}
	var buf bytes.Buffer
	command.SetOut(&buf)

	if err := printGenerateResult(command, result); err != nil {
		t.Fatalf("printGenerateResult failed: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "warning: port not detected") {
		t.Errorf("output missing warning line, got %q", out)
	}
	if !strings.Contains(out, "planned: Dockerfile") {
		t.Errorf("output missing planned Dockerfile line, got %q", out)
	}
	if !strings.Contains(out, "skipped: docker-compose.yml") {
		t.Errorf("output missing skipped docker-compose.yml line, got %q", out)
	}
}
