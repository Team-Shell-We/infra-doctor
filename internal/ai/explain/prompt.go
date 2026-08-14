package explain

import (
	"strings"
	"text/template"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

// systemPrompt : topic과 무관하게 고정

/*
* 모델의 역할을 정하고 답변을 JSON 스키마로 강제하는 시스템 프롬프트. 이 프롬프트는
* internal/ai/explain/result.go가 기대하는 JSON 형태로 강제함.
* 사용자는 프롬프트를 전혀 쓰지 않고,
* CLI가 모든 topic에 항상 같은 지시를 보낸다.
* 이 섹션은 이제 Go 코드가
* 결정론적으로 계산한다(status.go 참고), 모델은 관여하지 않는다.
 */

const systemPrompt = `You are Infra Doctor, a CLI tool that explains infrastructure concepts to backend developers strictly in the context of their own scanned project.

Base every claim ONLY on the scanned project facts you are given in the user message. Do not invent technologies, files, or configuration that were not listed there.

Respond with a single JSON object and nothing else (no markdown, no code fences, no text before or after the JSON) matching exactly this schema:

{
  "current_project": ["short bullet describing a relevant piece of the detected stack", "..."],
  "build_flow": ["step 1", "step 2", "..."],
  "why_topic": ["one concrete reason this technology matters for THIS project", "..."]
}

Rules:
- current_project: 2-5 bullets, only things actually present in the scanned facts.
- build_flow: 3-6 short steps showing how this topic fits into this project's build/deploy flow.
- why_topic: 2-4 bullets, specific to this project's actual stack, not generic textbook reasons.
- The user message includes an "already-verified status" list marking each fact as "present" or "absent". For any step or reason that depends on a fact marked "absent", phrase it as what adopting it WOULD look like (e.g. "You would deploy...", "Adding a Kubernetes manifest would...") — never state it as something already happening. Only describe a step as already happening when its fact is marked "present".
- Every string must be plain text, no markdown formatting.
- Do not include any field other than current_project, build_flow, and why_topic — a separate, already-verified file-status checklist is shown to the user by other means, so do not attempt to reproduce or guess one.`

const userPromptTemplate = `Explain {{.DisplayName}} ({{.Topic}}) in the context of my project.

Scanned project facts:
{{.Summary}}
{{- if .StatusFacts}}

Already-verified status for this topic (for context only, do not repeat these as a separate field):
{{.StatusFacts}}
{{- end}}

Respond entirely in {{.Language}}.`

type userPromptData struct {
	Topic       string
	DisplayName string
	Summary     string
	StatusFacts string
	Language    string
}

// languageNames : i18n 언어 코드를 모델에게 지시할 사람이 읽는 이름으로
// 매핑. 고정 문구 번역(internal/i18n)과는 별개 메커니즘 — AI 응답은
// 정적 번역이 아니라 프롬프트 지시로 언어를 맞춘다.
var languageNames = map[string]string{
	"ko": "Korean",
	"en": "English",
}

// BuildRequest는 topic, 스캔 요약, 결정론적으로 계산된 상태 사실(status.go),
// 언어 설정으로 completion 요청 전체를 조립한다. 순수 함수라 네트워크
// 호출 없이 프롬프트 생성을 단위 테스트할 수 있다.
func BuildRequest(topic string, summary ai.Summary, status []StatusItem, lang string) (ai.CompletionRequest, error) {

	tmpl, err := template.New("explain-user").Parse(userPromptTemplate)
	if err != nil {
		return ai.CompletionRequest{}, err
	}

	summaryText := summary.String()
	if summaryText == "" {
		summaryText = "(no relevant technologies were detected by the scanner)"
	}

	languageName, ok := languageNames[lang]
	if !ok {
		languageName = languageNames["en"]
	}

	data := userPromptData{
		Topic:       topic,
		DisplayName: DisplayName(topic),
		Summary:     summaryText,
		StatusFacts: formatStatusFacts(status),
		Language:    languageName,
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

func formatStatusFacts(status []StatusItem) string {

	var lines []string

	for _, item := range status {

		state := "absent"
		if item.Present {
			state = "present"
		}

		lines = append(lines, "- "+item.Label+": "+state)
	}

	return strings.Join(lines, "\n")
}
