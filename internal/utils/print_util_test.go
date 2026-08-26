package utils

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

// captureStdout : f 실행 중 os.Stdout에 쓰인 내용을 문자열로 반환.
// PrintAbout 등이 io.Writer 대신 fmt.Printf로 표준출력에 직접 쓰기 때문에 필요
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

func TestPrintAboutIncludesTaglineForBothLanguages(t *testing.T) {

	for _, lang := range []string{i18n.English, i18n.Korean} {
		out := captureStdout(t, func() { PrintAbout(lang) })

		if !strings.Contains(out, i18n.Get(lang, "about.tagline")) {
			t.Errorf("PrintAbout(%q) output missing tagline, got %q", lang, out)
		}
	}
}

func TestPrintDonateInfoIncludesKofiLink(t *testing.T) {

	out := captureStdout(t, func() { PrintDonateInfo(i18n.English) })

	if !strings.Contains(out, "https://ko-fi.com/shellwe") {
		t.Errorf("PrintDonateInfo output missing Ko-fi link, got %q", out)
	}
}

func TestPrintHelpIncludesCoreCommandsForBothLanguages(t *testing.T) {

	for _, lang := range []string{i18n.English, i18n.Korean} {
		out := captureStdout(t, func() { PrintHelp(lang) })

		if !strings.Contains(out, i18n.Get(lang, "help.coreCommands")) {
			t.Errorf("PrintHelp(%q) output missing core commands label, got %q", lang, out)
		}

		if !strings.Contains(out, "infra-doctor <command> [arguments]") {
			t.Errorf("PrintHelp(%q) output missing usage line, got %q", lang, out)
		}
	}
}
