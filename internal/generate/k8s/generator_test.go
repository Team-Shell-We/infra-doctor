package k8s

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

func TestPlanIncludesResourceLimitsAndSpringProfile(t *testing.T) {
	ctx := generate.Context{
		ServiceName: "application", Namespace: "default", Replicas: 1,
		DockerImage: "application:latest", ApplicationPort: 8080,
		Lang: "en",
	}

	files, err := (Generator{}).Plan(ctx)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	var deployment string
	for _, f := range files {
		if f.Path == "k8s/deployment.yml" {
			deployment = string(f.Content)
		}
	}
	if deployment == "" {
		t.Fatal("expected k8s/deployment.yml in the planned files")
	}

	// resources 없이 스케줄링하면 노드 자원을 고갈시킬 수 있음
	for _, want := range []string{"resources:", "requests:", "limits:"} {
		if !strings.Contains(deployment, want) {
			t.Errorf("expected %q in the Deployment spec", want)
		}
	}
	if !strings.Contains(deployment, "SPRING_PROFILES_ACTIVE") {
		t.Error("expected SPRING_PROFILES_ACTIVE to match Dockerfile/compose")
	}
}
