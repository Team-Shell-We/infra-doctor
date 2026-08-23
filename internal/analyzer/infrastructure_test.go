package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

// 실사용 시나리오 회귀 테스트: node_modules 안에 우연히 Dockerfile이 있어도
// (예: 서드파티 패키지가 자기 Dockerfile을 포함) 그건 무시하고, 프로젝트
// 루트의 진짜 Dockerfile만 감지해야 함
func TestAnalyzeInfrastructureSkipsExcludedDirs(t *testing.T) {

	root := t.TempDir()

	mustWriteFile(t, filepath.Join(root, "Dockerfile"), "FROM eclipse-temurin:17")
	mustWriteFile(t, filepath.Join(root, "node_modules", "some-pkg", "Dockerfile"), "FROM node:20")
	mustWriteFile(t, filepath.Join(root, ".git", "hooks", "Dockerfile"), "not real")

	info, err := AnalyzeInfrastructure(root)
	if err != nil {
		t.Fatalf("AnalyzeInfrastructure failed: %v", err)
	}

	if len(info.Docker.Dockerfiles) != 1 {
		t.Fatalf("expected exactly 1 Dockerfile found, got %d: %+v", len(info.Docker.Dockerfiles), info.Docker.Dockerfiles)
	}

	if info.Docker.Dockerfiles[0].Path != filepath.Join(root, "Dockerfile") {
		t.Errorf("expected the root Dockerfile, got %q", info.Docker.Dockerfiles[0].Path)
	}
}

// 회귀 테스트: recommend가 "이미 replicas를 늘려서 쓰고 있다"는 근거를
// 대려면 manifest에서 replicas 값을 실제로 읽어와야 함
func TestAnalyzeInfrastructureReadsKubernetesReplicas(t *testing.T) {

	root := t.TempDir()

	mustWriteFile(t, filepath.Join(root, "k8s", "deployment.yaml"), "apiVersion: apps/v1\nkind: Deployment\nspec:\n  replicas: 5\n")

	info, err := AnalyzeInfrastructure(root)
	if err != nil {
		t.Fatalf("AnalyzeInfrastructure failed: %v", err)
	}

	if info.Kubernetes.Replicas != 5 {
		t.Errorf("Replicas = %d, want 5", info.Kubernetes.Replicas)
	}
}

func TestAnalyzeInfrastructureNoReplicasFieldDefaultsToZero(t *testing.T) {

	root := t.TempDir()

	mustWriteFile(t, filepath.Join(root, "k8s", "service.yaml"), "apiVersion: v1\nkind: Service\n")

	info, err := AnalyzeInfrastructure(root)
	if err != nil {
		t.Fatalf("AnalyzeInfrastructure failed: %v", err)
	}

	if info.Kubernetes.Replicas != 0 {
		t.Errorf("Replicas = %d, want 0 when no manifest declares it", info.Kubernetes.Replicas)
	}
}

func mustWriteFile(t *testing.T, path, content string) {

	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create dir for %q: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %q: %v", path, err)
	}
}
