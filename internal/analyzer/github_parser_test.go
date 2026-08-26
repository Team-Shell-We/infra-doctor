package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseWorkflowOnAsScalar(t *testing.T) {

	root := t.TempDir()
	path := filepath.Join(root, "ci.yml")

	content := "name: CI\non: push\njobs:\n  build:\n    runs-on: ubuntu-latest\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	workflow, err := ParseWorkflow(path, "ci.yml")
	if err != nil {
		t.Fatalf("ParseWorkflow failed: %v", err)
	}

	if len(workflow.Triggers) != 1 || workflow.Triggers[0].Event != "push" {
		t.Errorf("Triggers = %+v, want [{Event: push}]", workflow.Triggers)
	}
}

func TestParseWorkflowOnAsSequence(t *testing.T) {

	root := t.TempDir()
	path := filepath.Join(root, "ci.yml")

	content := "name: CI\non: [push, pull_request]\njobs:\n  build:\n    runs-on: ubuntu-latest\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	workflow, err := ParseWorkflow(path, "ci.yml")
	if err != nil {
		t.Fatalf("ParseWorkflow failed: %v", err)
	}

	if len(workflow.Triggers) != 2 {
		t.Fatalf("Triggers = %+v, want 2 entries", workflow.Triggers)
	}
	if workflow.Triggers[0].Event != "push" || workflow.Triggers[1].Event != "pull_request" {
		t.Errorf("Triggers = %+v, want push then pull_request", workflow.Triggers)
	}
}

func TestParseWorkflowOnAsMappingWithBranches(t *testing.T) {

	root := t.TempDir()
	path := filepath.Join(root, "ci.yml")

	content := "name: CI\non:\n  push:\n    branches: [main, develop]\n  pull_request: {}\njobs:\n  build:\n    runs-on: ubuntu-latest\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	workflow, err := ParseWorkflow(path, "ci.yml")
	if err != nil {
		t.Fatalf("ParseWorkflow failed: %v", err)
	}

	if len(workflow.Triggers) != 2 {
		t.Fatalf("Triggers = %+v, want 2 entries", workflow.Triggers)
	}
	if workflow.Triggers[0].Event != "push" || len(workflow.Triggers[0].Branches) != 2 {
		t.Errorf("push trigger = %+v, want branches [main develop]", workflow.Triggers[0])
	}
	if workflow.Triggers[1].Event != "pull_request" {
		t.Errorf("second trigger = %+v, want pull_request", workflow.Triggers[1])
	}
}
