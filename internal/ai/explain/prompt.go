package explain

import (
	"strings"
	"text/template"

	"github.com/Team-Shell-We/infra-doctor/internal/ai"
)

// systemPrompt is topic-independent: it fixes the model's role and forces
// its answer into the exact JSON shape internal/ai/explain/result.go
// expects. This is what keeps `explain` from being a free-form chatbot —
// the user never writes any of this, the CLI always sends the same
// instructions for every topic.
//
// Deliberately NOT part of the schema: a "current status" / "which files
// exist" field. Asking the model to freely list "2-5 relevant files" was
// tried and reliably produced plausible-but-fabricated filenames (e.g.
// "Dockerfile.dev") that were never actually scanned for — a direct
// violation of "only state scanned facts". That section is now computed
// deterministically in Go (see status.go) and never touches the model.
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
- Every string must be plain text, no markdown formatting.
- Do not include any field other than current_project, build_flow, and why_topic — a separate, already-verified file-status checklist is shown to the user by other means, so do not attempt to reproduce or guess one.`

const userPromptTemplate = `Explain {{.DisplayName}} ({{.Topic}}) in the context of my project.

Scanned project facts:
{{.Summary}}
{{- if .StatusFacts}}

Already-verified status for this topic (for context only, do not repeat these as a separate field):
{{.StatusFacts}}
{{- end}}`

type userPromptData struct {
	Topic       string
	DisplayName string
	Summary     string
	StatusFacts string
}

// BuildRequest assembles the full completion request for a given topic,
// scanned project summary, and deterministically-computed status facts
// (see status.go). It's a pure function so prompt construction can be unit
// tested without any network call.
func BuildRequest(topic string, summary ai.Summary, status []StatusItem) (ai.CompletionRequest, error) {

	tmpl, err := template.New("explain-user").Parse(userPromptTemplate)
	if err != nil {
		return ai.CompletionRequest{}, err
	}

	summaryText := summary.String()
	if summaryText == "" {
		summaryText = "(no relevant technologies were detected by the scanner)"
	}

	data := userPromptData{
		Topic:       topic,
		DisplayName: DisplayName(topic),
		Summary:     summaryText,
		StatusFacts: formatStatusFacts(status),
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
