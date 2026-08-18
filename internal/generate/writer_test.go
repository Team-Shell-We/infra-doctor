package generate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriterProtectsAndOverwritesExistingFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Dockerfile")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	files := []File{{Path: "Dockerfile", Content: []byte("new")}}

	protected, err := (Writer{}).Write(root, files, WriteOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(protected.Skipped) != 1 {
		t.Fatalf("Skipped = %v, want Dockerfile", protected.Skipped)
	}

	overwritten, err := (Writer{}).Write(root, files, WriteOptions{Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(overwritten.Overwritten) != 1 {
		t.Fatalf("Overwritten = %v, want Dockerfile", overwritten.Overwritten)
	}
	content, err := os.ReadFile(path)
	if err != nil || string(content) != "new" {
		t.Fatalf("content = %q, err = %v", content, err)
	}
}
