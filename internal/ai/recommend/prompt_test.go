package recommend

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

func TestBuildRequestIncludesDecisionAndSummary(t *testing.T) {

	summary := ai.Summary{Framework: "Spring Boot 3.5.7"}
	decision := Decision{Recommended: "Docker Compose", Reasons: []string{"Single API server"}}

	req, err := BuildRequest(summary, decision, "en")
	if err != nil {
		t.Fatalf("BuildRequest failed: %v", err)
	}

	if !req.JSONMode {
		t.Error("expected JSONMode to be true")
	}

	if !strings.Contains(req.UserPrompt, "Docker Compose") {
		t.Errorf("UserPrompt should mention the recommendation, got: %s", req.UserPrompt)
	}

	if !strings.Contains(req.UserPrompt, "Spring Boot 3.5.7") {
		t.Errorf("UserPrompt should embed the scanned summary, got: %s", req.UserPrompt)
	}

	if !strings.Contains(req.UserPrompt, "Single API server") {
		t.Errorf("UserPrompt should embed the decision's reason labels, got: %s", req.UserPrompt)
	}

	// OpenAI JSON 모드는 프롬프트 안에 "json"이라는 단어가 있어야 요청이 통과된다.
	if !strings.Contains(strings.ToLower(req.SystemPrompt+req.UserPrompt), "json") {
		t.Error("prompt must mention 'json' somewhere for OpenAI JSON mode to be accepted")
	}

	if !strings.Contains(req.SystemPrompt, "reasons") {
		t.Error("SystemPrompt is missing schema field \"reasons\" — result.go's json tag must match")
	}
}

func TestBuildRequestNeverAsksModelToRedecide(t *testing.T) {

	req, err := BuildRequest(ai.Summary{}, Decision{Recommended: "Kubernetes"}, "en")
	if err != nil {
		t.Fatalf("BuildRequest failed: %v", err)
	}

	if strings.Contains(req.SystemPrompt, "recommended_strategy") || strings.Contains(req.SystemPrompt, `"recommended"`) {
		t.Error("SystemPrompt must not ask the model for a recommendation field — that's decided in decision.go")
	}
}

func TestBuildRequestLanguageDirective(t *testing.T) {

	reqKo, err := BuildRequest(ai.Summary{}, Decision{Recommended: "Docker Compose"}, "ko")
	if err != nil {
		t.Fatalf("BuildRequest failed: %v", err)
	}

	if !strings.Contains(reqKo.UserPrompt, "Respond entirely in Korean.") {
		t.Errorf("expected a Korean language directive, got: %s", reqKo.UserPrompt)
	}

	reqUnknown, err := BuildRequest(ai.Summary{}, Decision{Recommended: "Docker Compose"}, "fr")
	if err != nil {
		t.Fatalf("BuildRequest failed: %v", err)
	}

	if !strings.Contains(reqUnknown.UserPrompt, "Respond entirely in English.") {
		t.Errorf("expected an unknown language code to fall back to English, got: %s", reqUnknown.UserPrompt)
	}
}
