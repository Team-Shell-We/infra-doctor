package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

// withFakeHome : os.UserHomeDir()가 임시 디렉터리를 가리키게 해
// ai.LoadOrDefault()가 실제 ~/.infra-doctor/config.json 대신 테스트용
// 파일을 읽게 만든다(HOME은 Unix, USERPROFILE은 Windows에서 쓰임)
func withFakeHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	return home
}

func TestCurrentLangDefaultsToEnglishWhenNotConfigured(t *testing.T) {

	withFakeHome(t)

	if got := currentLang(); got != i18n.English {
		t.Errorf("currentLang() = %q, want %q when no config exists", got, i18n.English)
	}
}

func TestCurrentLangReadsSavedLanguage(t *testing.T) {

	home := withFakeHome(t)

	configDir := filepath.Join(home, ".infra-doctor")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(`{"language":"ko"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	if got := currentLang(); got != i18n.Korean {
		t.Errorf("currentLang() = %q, want %q", got, i18n.Korean)
	}
}
