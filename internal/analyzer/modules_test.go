package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeGradleModulesCountsIncludes(t *testing.T) {

	root := t.TempDir()

	content := "rootProject.name = 'demo'\ninclude 'api', 'worker', 'common'\n"
	if err := os.WriteFile(filepath.Join(root, "settings.gradle"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	count, err := AnalyzeGradleModules(root)
	if err != nil {
		t.Fatalf("AnalyzeGradleModules failed: %v", err)
	}

	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestAnalyzeGradleModulesNoSettingsFile(t *testing.T) {

	root := t.TempDir()

	count, err := AnalyzeGradleModules(root)
	if err != nil {
		t.Fatalf("AnalyzeGradleModules failed: %v", err)
	}

	if count != 0 {
		t.Errorf("count = %d, want 0 when there's no settings.gradle", count)
	}
}

func TestAnalyzeGradleModulesMultipleIncludeLines(t *testing.T) {

	root := t.TempDir()

	content := "include(\":api\")\ninclude(\":worker\")\n"
	if err := os.WriteFile(filepath.Join(root, "settings.gradle.kts"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	count, err := AnalyzeGradleModules(root)
	if err != nil {
		t.Fatalf("AnalyzeGradleModules failed: %v", err)
	}

	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}
