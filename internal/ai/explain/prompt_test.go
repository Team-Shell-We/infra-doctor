package explain

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

func TestBuildRequestIncludesTopicAndSummary(t *testing.T) {

	summary := ai.Summary{Framework: "Spring Boot 3.5.7 (Gradle), Java 17"}
	status := []StatusItem{{Label: "Dockerfile", Present: true}, {Label: "Docker Compose", Present: false}}

	req, err := BuildRequest("docker", summary, status)
	if err != nil {
		t.Fatalf("BuildRequest failed: %v", err)
	}

	if !req.JSONMode {
		t.Error("expected JSONMode to be true")
	}

	if !strings.Contains(req.UserPrompt, "Docker") {
		t.Errorf("UserPrompt should mention the display name Docker, got: %s", req.UserPrompt)
	}

	if !strings.Contains(req.UserPrompt, "docker") {
		t.Errorf("UserPrompt should mention the raw topic docker, got: %s", req.UserPrompt)
	}

	if !strings.Contains(req.UserPrompt, "Spring Boot 3.5.7") {
		t.Errorf("UserPrompt should embed the scanned summary, got: %s", req.UserPrompt)
	}

	if !strings.Contains(req.UserPrompt, "Dockerfile: present") || !strings.Contains(req.UserPrompt, "Docker Compose: absent") {
		t.Errorf("UserPrompt should embed the deterministic status facts, got: %s", req.UserPrompt)
	}

	// OpenAI JSON 모드는 프롬프트 안에 "json"이라는 단어가 있어야 요청이 통과됨
	if !strings.Contains(strings.ToLower(req.SystemPrompt+req.UserPrompt), "json") {
		t.Error("prompt must mention 'json' somewhere for OpenAI JSON mode to be accepted")
	}

	// 시스템 프롬프트의 스키마 필드명은 result.go의 json 태그와 정확히 일치해야 모델 응답을 파싱 가능
	for _, field := range []string{"current_project", "build_flow", "why_topic"} {
		if !strings.Contains(req.SystemPrompt, field) {
			t.Errorf("SystemPrompt is missing schema field %q — result.go's json tags must match", field)
		}
	}

	// current_status는 절대 모델에게 요청하지 않음! — status.go에서 결정론적으로 판단
	if strings.Contains(req.SystemPrompt, "current_status") {
		t.Error("SystemPrompt must not ask the model for current_status — that section is computed in Go, not by the AI")
	}
}

func TestBuildRequestEmptySummaryFallsBack(t *testing.T) {

	req, err := BuildRequest("redis", ai.Summary{}, nil)
	if err != nil {
		t.Fatalf("BuildRequest failed: %v", err)
	}

	if !strings.Contains(req.UserPrompt, "no relevant technologies") {
		t.Errorf("expected a fallback message for an empty summary, got: %s", req.UserPrompt)
	}
}

func TestBuildRequestNoStatusOmitsSection(t *testing.T) {

	req, err := BuildRequest("redis", ai.Summary{Framework: "Spring Boot"}, nil)
	if err != nil {
		t.Fatalf("BuildRequest failed: %v", err)
	}

	if strings.Contains(req.UserPrompt, "Already-verified status") {
		t.Errorf("UserPrompt should omit the status section entirely when there is no status data, got: %s", req.UserPrompt)
	}
}
