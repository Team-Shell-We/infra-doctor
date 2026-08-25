package export

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
	"github.com/Team-Shell-We/infra-doctor/internal/generate"
	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func TestBuildFilesIncludesNginx(t *testing.T) {

	root := t.TempDir()
	info := project.Info{Framework: project.FrameworkInfo{BuildTool: project.BuildToolInfo{Type: "Gradle"}}}

	files, _, err := buildFiles(root, info, &doctor.Result{}, "en")
	if err != nil {
		t.Fatalf("buildFiles failed: %v", err)
	}

	found := false
	for _, f := range files {
		if filepath.ToSlash(f.Path) == "docker/nginx.conf" {
			found = true
		}
	}
	if !found {
		t.Error("expected docker/nginx.conf among the exported files")
	}
}

func TestBuildFilesPropagatesContextWarnings(t *testing.T) {

	root := t.TempDir()
	// Java 버전 미감지 상태로 두면 BuildContext가 경고를 반환
	info := project.Info{Framework: project.FrameworkInfo{BuildTool: project.BuildToolInfo{Type: "Gradle"}}}

	_, warnings, err := buildFiles(root, info, &doctor.Result{}, "en")
	if err != nil {
		t.Fatalf("buildFiles failed: %v", err)
	}

	if len(warnings) == 0 {
		t.Error("expected a warning about undetected Java version")
	}
}

func TestPrintResultPrintsWarnings(t *testing.T) {

	var buf strings.Builder

	err := printResult(&buf, generate.Result{Warnings: []string{"something"}})
	if err != nil {
		t.Fatalf("printResult failed: %v", err)
	}

	if !strings.Contains(buf.String(), "warning: something") {
		t.Errorf("expected warning line in output, got %q", buf.String())
	}
}
