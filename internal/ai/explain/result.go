package explain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Result is the fixed 3-section shape every `explain` answer must fill in
// (see prompt.go's system prompt, which instructs the model to return
// exactly this JSON shape). Field names/tags here and the field names named
// in the prompt template must stay in sync.
//
// "Current Status" is deliberately NOT part of this struct — it is never
// AI-generated, see status.go for why.
type Result struct {
	CurrentProject []string `json:"current_project"`
	BuildFlow      []string `json:"build_flow"`
	WhyTopic       []string `json:"why_topic"`
}

// Parse decodes and validates a model completion. It deliberately requires
// every section to be non-empty — a model reply missing a section is
// treated as a failure to render rather than silently shown as an empty box.
func Parse(raw string) (*Result, error) {

	var result Result

	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &result); err != nil {
		return nil, fmt.Errorf("explain: failed to parse AI response as JSON: %w", err)
	}

	var missing []string

	if len(result.CurrentProject) == 0 {
		missing = append(missing, "current_project")
	}

	if len(result.BuildFlow) == 0 {
		missing = append(missing, "build_flow")
	}

	if len(result.WhyTopic) == 0 {
		missing = append(missing, "why_topic")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("explain: AI response is missing required field(s): %s", strings.Join(missing, ", "))
	}

	return &result, nil
}
