package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWorkflowRootFindsDirectoryContainingWorkflows(t *testing.T) {

	root := t.TempDir()
	workflows := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(workflows, 0o755); err != nil {
		t.Fatal(err)
	}

	nested := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	found, ok := findWorkflowRoot(nested)
	if !ok {
		t.Fatal("findWorkflowRoot: want ok=true, got false")
	}

	wantRoot, _ := filepath.EvalSymlinks(root)
	gotRoot, _ := filepath.EvalSymlinks(found)
	if gotRoot != wantRoot {
		t.Errorf("findWorkflowRoot found = %q, want %q", found, root)
	}
}

func TestFindWorkflowRootReturnsFalseWhenNoWorkflowsExist(t *testing.T) {

	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	if _, ok := findWorkflowRoot(nested); ok {
		t.Error("findWorkflowRoot with no .github/workflows anywhere: want ok=false, got true")
	}
}
