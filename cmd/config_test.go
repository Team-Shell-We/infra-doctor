package cmd

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = original

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}

	return string(out)
}

func TestPrintConfigShowsUnconfiguredProviderWhenNotLoggedIn(t *testing.T) {

	out := captureStdout(t, func() { printConfig(ai.Credentials{}) })

	if !strings.Contains(out, "Not configured") {
		t.Errorf("printConfig with zero-value creds: want \"Not configured\" provider, got %q", out)
	}

	if !strings.Contains(out, "ASCII + Mermaid") {
		t.Errorf("printConfig: want default output format \"ASCII + Mermaid\", got %q", out)
	}
}

func TestPrintConfigShowsProviderAndLanguage(t *testing.T) {

	out := captureStdout(t, func() {
		printConfig(ai.Credentials{Provider: "openai", Language: "ko"})
	})

	if !strings.Contains(out, "openai") {
		t.Errorf("printConfig: want provider %q in output, got %q", "openai", out)
	}

	if !strings.Contains(out, "ko") {
		t.Errorf("printConfig: want language %q in output, got %q", "ko", out)
	}
}
