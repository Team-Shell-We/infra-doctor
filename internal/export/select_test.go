package export

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

func sampleFiles() []generate.File {
	return []generate.File{
		{Path: "report.md"},
		{Path: "architecture.md"},
		{Path: "architecture.mmd"},
		{Path: "erd.md"},
		{Path: "erd.mmd"},
		{Path: "deployment-flow.md"},
		{Path: "recommendations.md"},
		{Path: "docker/Dockerfile"},
		{Path: "docker/docker-compose.yml"},
		{Path: "kubernetes/deployment.yaml"},
		{Path: "github/ci.yml"},
	}
}

func TestCategorizeGroupsFilesInFixedOrder(t *testing.T) {

	categories := Categorize(sampleFiles())

	wantOrder := []string{"report", "architecture", "erd", "flow", "recommendations", "docker", "kubernetes", "github"}
	if len(categories) != len(wantOrder) {
		t.Fatalf("len(categories) = %d, want %d", len(categories), len(wantOrder))
	}
	for i, id := range wantOrder {
		if categories[i].ID != id {
			t.Errorf("categories[%d].ID = %q, want %q", i, categories[i].ID, id)
		}
	}

	for _, c := range categories {
		if c.ID == "docker" && len(c.Files) != 2 {
			t.Errorf("docker category has %d files, want 2", len(c.Files))
		}
	}
}

func TestCategorizeOmitsEmptyCategories(t *testing.T) {

	files := []generate.File{{Path: "report.md"}}

	categories := Categorize(files)

	if len(categories) != 1 || categories[0].ID != "report" {
		t.Errorf("expected only the report category, got %+v", categories)
	}
}

func TestPromptSelectionParsesCommaSeparatedNumbers(t *testing.T) {

	categories := Categorize(sampleFiles())
	var out strings.Builder

	selected, err := PromptSelection(strings.NewReader("1,3\n"), &out, categories, "en")
	if err != nil {
		t.Fatalf("PromptSelection failed: %v", err)
	}

	// 1=report(1개), 3=erd(2개) -> 총 3개
	if len(selected) != 3 {
		t.Errorf("len(selected) = %d, want 3 (report + erd.md + erd.mmd)", len(selected))
	}
}

func TestPromptSelectionAll(t *testing.T) {

	categories := Categorize(sampleFiles())
	var out strings.Builder

	selected, err := PromptSelection(strings.NewReader("all\n"), &out, categories, "en")
	if err != nil {
		t.Fatalf("PromptSelection failed: %v", err)
	}

	if len(selected) != len(sampleFiles()) {
		t.Errorf("len(selected) = %d, want %d", len(selected), len(sampleFiles()))
	}
}

func TestPromptSelectionRejectsOutOfRangeNumber(t *testing.T) {

	categories := Categorize(sampleFiles())
	var out strings.Builder

	_, err := PromptSelection(strings.NewReader("9\n"), &out, categories, "en")
	if err == nil {
		t.Error("expected an error for an out-of-range selection, got nil")
	}
}

func TestPromptSelectionRejectsEmptyInput(t *testing.T) {

	categories := Categorize(sampleFiles())
	var out strings.Builder

	_, err := PromptSelection(strings.NewReader("\n"), &out, categories, "en")
	if err == nil {
		t.Error("expected an error for empty input, got nil")
	}
}
