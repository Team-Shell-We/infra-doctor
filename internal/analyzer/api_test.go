package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeAPICountsControllersAndEndpoints(t *testing.T) {

	root := t.TempDir()

	src := `package com.example.demo;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/users")
public class UserController {

    @GetMapping
    public String list() { return ""; }

    @GetMapping("/{id}")
    public String get() { return ""; }

    @PostMapping
    public String create() { return ""; }
}
`
	if err := os.WriteFile(filepath.Join(root, "UserController.java"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	info, err := AnalyzeAPI(root)
	if err != nil {
		t.Fatalf("AnalyzeAPI failed: %v", err)
	}

	if info.ControllerCount != 1 {
		t.Errorf("ControllerCount = %d, want 1", info.ControllerCount)
	}
	if info.EndpointCount != 3 {
		t.Errorf("EndpointCount = %d, want 3 (GetMapping x2 + PostMapping x1)", info.EndpointCount)
	}
}

func TestAnalyzeAPIIgnoresNonJavaFiles(t *testing.T) {

	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "notes.md"), []byte("@RestController @GetMapping"), 0644); err != nil {
		t.Fatal(err)
	}

	info, err := AnalyzeAPI(root)
	if err != nil {
		t.Fatalf("AnalyzeAPI failed: %v", err)
	}

	if info.ControllerCount != 0 || info.EndpointCount != 0 {
		t.Errorf("expected zero counts for non-.java files, got %+v", info)
	}
}
