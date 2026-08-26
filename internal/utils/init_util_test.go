package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeCreatesConfigAndGitignore(t *testing.T) {

	root := t.TempDir()

	result, err := Initialize(root)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if len(result.CreatedPaths) != 3 {
		t.Errorf("CreatedPaths = %v, want 3 entries", result.CreatedPaths)
	}

	config, err := os.ReadFile(filepath.Join(root, ".infra-doctor", "config.yaml"))
	if err != nil {
		t.Fatalf("config.yaml not created: %v", err)
	}
	if string(config) != defaultConfig {
		t.Errorf("config.yaml content = %q, want %q", config, defaultConfig)
	}

	gitignore, err := os.ReadFile(filepath.Join(root, ".infra-doctor", ".gitignore"))
	if err != nil {
		t.Fatalf(".gitignore not created: %v", err)
	}
	if string(gitignore) != defaultGitignore {
		t.Errorf(".gitignore content = %q, want %q", gitignore, defaultGitignore)
	}
}

func TestInitializeReturnsErrAlreadyInitializedWhenDirExists(t *testing.T) {

	root := t.TempDir()

	if _, err := Initialize(root); err != nil {
		t.Fatalf("first Initialize failed: %v", err)
	}

	_, err := Initialize(root)
	if err != ErrAlreadyInitialized {
		t.Errorf("second Initialize err = %v, want %v", err, ErrAlreadyInitialized)
	}
}

func TestInitializeReturnsErrorForEmptyProjectDir(t *testing.T) {

	if _, err := Initialize(""); err == nil {
		t.Fatal("Initialize(\"\"): want error, got nil")
	}
}
