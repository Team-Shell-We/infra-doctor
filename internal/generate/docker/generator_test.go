package docker

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

func TestPlanRejectsNonJavaRuntime(t *testing.T) {
	ctx := generate.Context{Runtime: "node"}

	if _, err := (Generator{}).Plan(ctx); err == nil {
		t.Fatal("expected an error for a non-java runtime")
	}
}

func TestPlanIncludesHealthCheckAndNonRootUser(t *testing.T) {
	ctx := generate.Context{
		Runtime:         "java",
		RuntimeVersion:  "21",
		BuildCommand:    "./gradlew clean bootJar",
		ArtifactPath:    "build/libs/*.jar",
		ApplicationPort: 8080,
		HealthPath:      "/actuator/health",
		Lang:            "en",
	}

	files, err := (Generator{}).Plan(ctx)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	var dockerfile string
	for _, f := range files {
		if f.Path == "Dockerfile" {
			dockerfile = string(f.Content)
		}
	}
	if dockerfile == "" {
		t.Fatal("expected a Dockerfile in the planned files")
	}

	// doctor의 no_health_check 룰은 "HEALTHCHECK" 문자열로 감지
	if !strings.Contains(dockerfile, "HEALTHCHECK") {
		t.Error("Dockerfile must contain a HEALTHCHECK instruction")
	}
	if !strings.Contains(dockerfile, "USER app") {
		t.Error("Dockerfile must switch to a non-root user before ENTRYPOINT")
	}
	if !strings.Contains(dockerfile, "eclipse-temurin:21-jre") {
		t.Error("runtime stage should use the jre image")
	}
}

func TestPlanRespectsLanguage(t *testing.T) {
	base := generate.Context{
		Runtime: "java", RuntimeVersion: "21",
		BuildCommand: "./gradlew clean bootJar", ArtifactPath: "build/libs/*.jar",
		ApplicationPort: 8080, HealthPath: "/",
	}

	en, ko := base, base
	en.Lang, ko.Lang = "en", "ko"

	enFiles, err := (Generator{}).Plan(en)
	if err != nil {
		t.Fatalf("Plan(en) failed: %v", err)
	}
	koFiles, err := (Generator{}).Plan(ko)
	if err != nil {
		t.Fatalf("Plan(ko) failed: %v", err)
	}

	if !strings.Contains(string(enFiles[0].Content), "Next steps:") {
		t.Error("lang=en should contain the English next-steps banner")
	}
	if !strings.Contains(string(koFiles[0].Content), "다음 할 일:") {
		t.Error("lang=ko should contain the Korean next-steps banner")
	}
}
