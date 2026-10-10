package recommend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
	"github.com/Team-Shell-We/infra-doctor/internal/ai/openai"
	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

// TestRecommendEvaluation : 실제 OpenAI 호출이 필요해 기본 `go test ./...`에서는 건너뛴다.
// Decide()의 결정 자체는 decision_test.go가 이미 결정론적으로 검증하므로, 여기서는 그 결정을
// LLM이 그대로 따르는지(재결정하지 않는지)와, 탐지되지 않은 기술을 언급하지 않는지를 대표
// 시나리오 5개로 점검하고 docs/evaluation.md에 실제 응답을 기록한다.
//
//	INFRA_DOCTOR_EVAL=1 go test ./internal/ai/recommend/... -run TestRecommendEvaluation -v
//
// 사전에 `infra-doctor login`으로 OpenAI API Key가 등록돼 있어야 한다.
func TestRecommendEvaluation(t *testing.T) {

	if os.Getenv("INFRA_DOCTOR_EVAL") != "1" {
		t.Skip("set INFRA_DOCTOR_EVAL=1 to run this (requires a live OpenAI call via `infra-doctor login`)")
	}

	creds, err := ai.Load()
	if err != nil {
		t.Fatalf("not logged in: %v", err)
	}

	client := openai.New(creds.APIKey)

	var results []evalResult

	for _, c := range evalCases() {

		summary := ai.BuildSummary(&c.Info)
		decision := Decide(&c.Info)

		if decision.Recommended != c.ExpectRecommended {
			t.Errorf("%s: fixture itself is wrong — Decide() = %q, want %q", c.Name, decision.Recommended, c.ExpectRecommended)
			continue
		}

		req, err := BuildRequest(summary, decision, "en")
		if err != nil {
			t.Fatalf("%s: BuildRequest failed: %v", c.Name, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		resp, err := client.Complete(ctx, req)
		cancel()
		if err != nil {
			t.Fatalf("%s: OpenAI call failed: %v", c.Name, err)
		}

		parsed, err := Parse(resp.Content)
		if err != nil {
			t.Errorf("%s: response failed schema validation: %v", c.Name, err)
			continue
		}

		grounded, flagged := groundingCheck(summary, parsed.Reasons)
		if !grounded {
			t.Errorf("%s: AI reasons mention undetected technology: %v", c.Name, flagged)
		}

		results = append(results, evalResult{
			Case:     c,
			Summary:  summary,
			Decision: decision,
			Reasons:  parsed.Reasons,
			Grounded: grounded,
			Flagged:  flagged,
		})
	}

	if err := writeEvaluationReport(results); err != nil {
		t.Errorf("failed to write docs/evaluation.md: %v", err)
	}
}

type evalCase struct {
	Name              string
	Info              project.Info
	ExpectRecommended string // Decide()가 결정론적으로 내놓아야 하는 값 — 픽스처 자체의 자가 점검용
}

func evalCases() []evalCase {
	return []evalCase{
		{
			Name: "최소 구성 (인프라 없음)",
			Info: project.Info{
				Framework: project.FrameworkInfo{
					SpringBoot: project.SpringBootInfo{Enabled: true, Version: "3.2.0"},
					BuildTool:  project.BuildToolInfo{Type: "Gradle"},
					Java:       project.JavaInfo{Version: "17"},
				},
			},
			ExpectRecommended: "Docker Compose",
		},
		{
			Name: "DB+Redis 사용, 낮은 복잡도",
			Info: project.Info{
				Framework: project.FrameworkInfo{
					SpringBoot: project.SpringBootInfo{Enabled: true, Version: "3.2.0"},
					BuildTool:  project.BuildToolInfo{Type: "Gradle"},
					Java:       project.JavaInfo{Version: "17"},
				},
				Database: project.DatabaseInfo{
					Primary: project.Database{Type: "PostgreSQL"},
					Redis:   &project.RedisInfo{Enabled: true},
				},
				Infrastructure: project.InfrastructureInfo{
					Docker: project.DockerInfo{Enabled: true},
				},
			},
			ExpectRecommended: "Docker Compose",
		},
		{
			Name: "고복잡도 (Kafka+DB+Redis+멀티모듈+API 25개)",
			Info: project.Info{
				Framework: project.FrameworkInfo{
					SpringBoot: project.SpringBootInfo{Enabled: true, Version: "3.2.0"},
					BuildTool:  project.BuildToolInfo{Type: "Gradle"},
					Java:       project.JavaInfo{Version: "17"},
					Modules:    project.ModuleInfo{Count: 4},
				},
				Dependencies: project.DependencyInfo{
					Kafka: project.KafkaInfo{Enabled: true},
				},
				Database: project.DatabaseInfo{
					Primary: project.Database{Type: "MySQL"},
					Redis:   &project.RedisInfo{Enabled: true},
				},
				API: project.APIInfo{EndpointCount: 25},
			},
			ExpectRecommended: "Kubernetes",
		},
		{
			Name: "이미 Kubernetes 사용 중, replicas 낮음",
			Info: project.Info{
				Framework: project.FrameworkInfo{
					SpringBoot: project.SpringBootInfo{Enabled: true, Version: "3.2.0"},
					BuildTool:  project.BuildToolInfo{Type: "Gradle"},
				},
				Infrastructure: project.InfrastructureInfo{
					Kubernetes: project.KubernetesInfo{Enabled: true, Replicas: 1},
				},
			},
			ExpectRecommended: "Kubernetes",
		},
		{
			Name: "이미 Kubernetes 사용 중, replicas 높음(스케일 중)",
			Info: project.Info{
				Framework: project.FrameworkInfo{
					SpringBoot: project.SpringBootInfo{Enabled: true, Version: "3.2.0"},
					BuildTool:  project.BuildToolInfo{Type: "Gradle"},
				},
				Infrastructure: project.InfrastructureInfo{
					Kubernetes: project.KubernetesInfo{Enabled: true, Replicas: 5},
				},
			},
			ExpectRecommended: "Kubernetes",
		},
	}
}

type evalResult struct {
	Case     evalCase
	Summary  ai.Summary
	Decision Decision
	Reasons  []string
	Grounded bool
	Flagged  []string
}

// watchedKeywords : AI 응답에 등장하면 "탐지된 사실에 있는지" 대조해보는 기술명 목록.
// Docker/Kubernetes는 추천 자체의 주제라 항상 언급되므로 대상에서 뺐다 — 진짜 위험한 건
// 감지되지 않은 '다른' 기술을 있는 것처럼 서술하는 경우다.
var watchedKeywords = []string{
	"Nginx", "Terraform", "Redis", "PostgreSQL", "MySQL", "MariaDB",
	"Kafka", "AWS", "Lombok", "Prometheus", "Grafana", "Jenkins", "CircleCI",
}

// groundingCheck : reasons에 등장하는 watchedKeywords 중, summary(실제 탐지된 사실)에
// 없는 게 있으면 hallucination 후보로 플래그한다.
func groundingCheck(summary ai.Summary, reasons []string) (ok bool, flagged []string) {

	allowed := append([]string{}, summary.Dependencies...)
	allowed = append(allowed, summary.Database...)
	allowed = append(allowed, summary.Infrastructure...)
	allowed = append(allowed, summary.CICD...)

	text := strings.Join(reasons, " ")

	for _, keyword := range watchedKeywords {

		if !strings.Contains(text, keyword) {
			continue
		}

		found := false
		for _, name := range allowed {
			if strings.Contains(name, keyword) {
				found = true
				break
			}
		}

		if !found {
			flagged = append(flagged, keyword)
		}
	}

	return len(flagged) == 0, flagged
}

func writeEvaluationReport(results []evalResult) error {

	var b strings.Builder

	b.WriteString("# recommend Evaluation Report\n\n")
	b.WriteString("`INFRA_DOCTOR_EVAL=1 go test ./internal/ai/recommend/... -run TestRecommendEvaluation -v`로 재생성할 수 있다. ")
	b.WriteString("배포 전략 결정 자체는 `decision_test.go`가 결정론적으로 검증하므로, 이 리포트는 두 가지만 본다 — ")
	b.WriteString("① LLM이 그 결정을 그대로 설명하는지(재결정하지 않는지), ② 탐지되지 않은 기술을 언급하지 않는지.\n\n")

	for _, r := range results {

		fmt.Fprintf(&b, "## %s\n\n", r.Case.Name)
		fmt.Fprintf(&b, "- **감지된 사실**: %s\n", summaryOneLine(r.Summary))
		fmt.Fprintf(&b, "- **결정(코드)**: %s\n", r.Decision.Recommended)

		if r.Grounded {
			b.WriteString("- **근거 점검**: ✅ 탐지된 사실 범위 안에서만 설명함\n\n")
		} else {
			fmt.Fprintf(&b, "- **근거 점검**: ❌ 미탐지 기술 언급: %v\n\n", r.Flagged)
		}

		b.WriteString("**AI 설명(실제 응답)**:\n\n")
		for _, reason := range r.Reasons {
			fmt.Fprintf(&b, "- %s\n", reason)
		}
		b.WriteString("\n")
	}

	path, err := filepath.Abs(filepath.Join("..", "..", "..", "docs", "evaluation.md"))
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func summaryOneLine(s ai.Summary) string {

	var parts []string

	if s.Framework != "" {
		parts = append(parts, s.Framework)
	}

	parts = append(parts, s.Dependencies...)
	parts = append(parts, s.Database...)
	parts = append(parts, s.Infrastructure...)
	parts = append(parts, s.CICD...)

	if len(parts) == 0 {
		return "(없음)"
	}

	return strings.Join(parts, ", ")
}
