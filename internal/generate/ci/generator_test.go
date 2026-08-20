package ci

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/generate"
)

func TestPlanIncludesCacheAndPermissions(t *testing.T) {
	ctx := generate.Context{
		BuildTool: "gradle", BuildCommand: "./gradlew clean bootJar",
		RuntimeVersion: "21", Lang: "en",
	}

	files, err := (Generator{}).Plan(ctx)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	content := string(files[0].Content)

	// actions/setup-java의 내장 캐싱을 안 켜면 매 실행마다 의존성을
	// 처음부터 다시 받음
	if !strings.Contains(content, "cache: gradle") {
		t.Error("expected setup-java's built-in gradle cache to be enabled")
	}
	// GITHUB_TOKEN 기본 권한은 필요 이상으로 넓은 경우가 많음
	if !strings.Contains(content, "permissions:") || !strings.Contains(content, "contents: read") {
		t.Error("expected an explicit minimal permissions block")
	}
}

func TestPlanUsesBuildToolAsCacheKey(t *testing.T) {
	ctx := generate.Context{
		BuildTool: "maven", BuildCommand: "./mvnw clean package -DskipTests",
		RuntimeVersion: "21", Lang: "en",
	}

	files, err := (Generator{}).Plan(ctx)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	content := string(files[0].Content)

	if !strings.Contains(content, "cache: maven") {
		t.Error("expected the cache key to match the detected build tool (maven)")
	}
}

func TestPlanRejectsEmptyBuildCommand(t *testing.T) {
	if _, err := (Generator{}).Plan(generate.Context{}); err == nil {
		t.Fatal("expected an error when there's no build command")
	}
}
