package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 회귀 테스트: `recommend k8s`처럼 존재하지 않는 경로를 인자로 잘못
// 넘겼을 때, "no build file found" 같은 모호한 메시지 대신 경로 자체가
// 없다는 걸 명확히 알려줘야 함
func TestAnalyzeProjectNonexistentPath(t *testing.T) {

	root := filepath.Join(t.TempDir(), "does-not-exist")

	_, err := AnalyzeProject(root)
	if err == nil {
		t.Fatal("expected an error for a nonexistent path, got nil")
	}

	wantSubstring := "path does not exist"
	if got := err.Error(); !strings.Contains(got, wantSubstring) {
		t.Errorf("error = %q, want it to contain %q", got, wantSubstring)
	}
}

func TestAnalyzeProjectPathIsAFile(t *testing.T) {

	path := filepath.Join(t.TempDir(), "not-a-dir.txt")

	if err := os.WriteFile(path, []byte("hi"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := AnalyzeProject(path)
	if err == nil {
		t.Fatal("expected an error when path is a file, not a directory")
	}

	wantSubstring := "is not a directory"
	if got := err.Error(); !strings.Contains(got, wantSubstring) {
		t.Errorf("error = %q, want it to contain %q", got, wantSubstring)
	}
}

// 경로는 존재하지만 build.gradle/pom.xml이 없는 경우엔 기존 그대로
// FindBuildFile의 "no build file found" 에러가 나야 함
func TestAnalyzeProjectValidDirWithoutBuildFile(t *testing.T) {

	root := t.TempDir()

	_, err := AnalyzeProject(root)
	if err == nil {
		t.Fatal("expected an error for a directory with no build file")
	}

	wantSubstring := "no build file found"
	if got := err.Error(); !strings.Contains(got, wantSubstring) {
		t.Errorf("error = %q, want it to contain %q", got, wantSubstring)
	}
}
