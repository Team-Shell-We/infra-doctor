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

	path := filepath.Join(t.TempDir(), "config.json")

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
