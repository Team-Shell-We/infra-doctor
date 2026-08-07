package recommend

import (
	"strings"
	"text/template"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

// systemPrompt: 배포 전략 자체는 이미 코드가 결정했음을 명시하고, AI에게는
// 그 결정을 뒤집지 않고 근거만 자연어로 풀어쓰도록 강제한다.
const systemPrompt = `You are Infra Doctor, a CLI tool that explains a deployment strategy recommendation to backend developers in the context of their own scanned project.

A deployment strategy has ALREADY been decided by deterministic rules, not by you. You are only explaining WHY it fits this project — never propose a different strategy, never contradict the given recommendation.

Base every claim ONLY on the scanned project facts and the given reason labels. Do not invent technologies, files, or configuration that were not listed.

Respond with a single JSON object and nothing else (no markdown, no code fences, no text before or after the JSON) matching exactly this schema:

{
  "reasons": ["one concrete, specific explanation grounded in this project's actual stack", "..."]
}

Rules:
- reasons: 2-4 bullets, each a full sentence expanding on the given reason labels — not generic textbook advice.
- Every string must be plain text, no markdown formatting.
- Do not include any field other than reasons.`

const userPromptTemplate = `A deployment strategy of {{.Recommended}} has been recommended for my project. Explain why, in the context of my project.

Scanned project facts:
{{.Summary}}

Reasons this strategy was chosen (already decided, explain these — do not change the recommendation):
{{.ReasonLabels}}`

type userPromptData struct {
	Recommended  string
	Summary      string
	ReasonLabels string
}

// BuildRequest는 이미 결정된 Decision과 스캔 요약으로 completion 요청을
// 조립하는 순수 함수 — 네트워크 호출 없이 단위 테스트 가능하다.
func BuildRequest(summary ai.Summary, decision Decision) (ai.CompletionRequest, error) {

	tmpl, err := template.New("recommend-user").Parse(userPromptTemplate)
	if err != nil {
		return ai.CompletionRequest{}, err
	}

	summaryText := summary.String()
	if summaryText == "" {
		summaryText = "(no relevant technologies were detected by the scanner)"
	}

	var reasonLines []string
	for _, reason := range decision.Reasons {
		reasonLines = append(reasonLines, "- "+reason)
	}

	data := userPromptData{
		Recommended:  decision.Recommended,
		Summary:      summaryText,
		ReasonLabels: strings.Join(reasonLines, "\n"),
	}

	var b strings.Builder

	if err := tmpl.Execute(&b, data); err != nil {
		return ai.CompletionRequest{}, err
	}

	return ai.CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   b.String(),
		JSONMode:     true,
	}, nil
}
