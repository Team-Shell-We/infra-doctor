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

	// OpenAI's JSON mode requires the literal word "json" to appear
	// somewhere in the prompt or the API rejects the request outright.
	if !strings.Contains(strings.ToLower(req.SystemPrompt+req.UserPrompt), "json") {
		t.Error("prompt must mention 'json' somewhere for OpenAI JSON mode to be accepted")
	}

	// The schema field names in the system prompt must match result.go's
	// json tags exactly, or the model's replies won't parse.
	for _, field := range []string{"current_project", "build_flow", "why_topic"} {
		if !strings.Contains(req.SystemPrompt, field) {
			t.Errorf("SystemPrompt is missing schema field %q — result.go's json tags must match", field)
		}
	}

	// current_status must never be requested from the model — it's
	// computed deterministically in status.go instead (see BuildStatus's
	// doc comment for why: the model was observed inventing filenames).
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
