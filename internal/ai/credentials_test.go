package ai

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSaveToAndLoadFromRoundTrip(t *testing.T) {

	path := filepath.Join(t.TempDir(), "nested", "config.json")

	creds := Credentials{Provider: "openai", APIKey: "sk-test", Login: true}

	if err := SaveTo(path, creds); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	loaded, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}

	if *loaded != creds {
		t.Errorf("round trip mismatch: got %+v, want %+v", *loaded, creds)
	}
}

func TestSaveToFilePermissions(t *testing.T) {

	if runtime.GOOS == "windows" {
		t.Skip("POSIX file permission bits are not meaningful on Windows")
	}

	// 디렉터리가 미리 존재하면 안 됨 : os.MkdirAll은 이미 있는 디렉터리엔
	// 아무 일도 안 하므로, 요청한 권한이 실제로 적용됐는지 검증 못한 채
	// 테스트가 그냥 통과해버림(t.TempDir()은 SaveTo 실행 전에 이미 존재).
	path := filepath.Join(t.TempDir(), "nested", "config.json")

	if err := SaveTo(path, Credentials{Provider: "openai", APIKey: "sk-test", Login: true}); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file permissions 0600, got %o", perm)
	}

	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("Stat dir failed: %v", err)
	}

	if perm := dirInfo.Mode().Perm(); perm != 0700 {
		t.Errorf("expected directory permissions 0700, got %o", perm)
	}
}

func TestLoadFromMissingFile(t *testing.T) {

	path := filepath.Join(t.TempDir(), "does-not-exist.json")

	_, err := LoadFrom(path)
	if !errors.Is(err, ErrNotLoggedIn) {
		t.Errorf("expected ErrNotLoggedIn, got %v", err)
	}
}

func TestLoadFromLoggedOutCredentials(t *testing.T) {

	path := filepath.Join(t.TempDir(), "config.json")

	if err := SaveTo(path, Credentials{Provider: "openai", APIKey: "sk-test", Login: false}); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	_, err := LoadFrom(path)
	if !errors.Is(err, ErrNotLoggedIn) {
		t.Errorf("expected ErrNotLoggedIn for login:false, got %v", err)
	}
}

func TestPath(t *testing.T) {

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // os.UserHomeDir() reads this on Windows

	path, err := Path()
	if err != nil {
		t.Fatalf("Path failed: %v", err)
	}

	want := filepath.Join(home, ".infra-doctor", "config.json")
	if path != want {
		t.Errorf("Path() = %q, want %q", path, want)
	}
}
