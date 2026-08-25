package cmd

import (
	"os"
	"path/filepath"
	"testing"
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
