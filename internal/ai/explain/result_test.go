package explain

import "testing"

func TestParseValid(t *testing.T) {

	raw := `{
		"current_project": ["Spring Boot", "PostgreSQL"],
		"build_flow": ["Source", "Build", "Image"],
		"why_topic": ["Consistent runtime environment"]
	}`

	result, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.CurrentProject) != 2 {
		t.Errorf("CurrentProject = %v, want 2 items", result.CurrentProject)
	}
}

func TestParseIgnoresUnexpectedFields(t *testing.T) {

	// If the model disobeys the "don't include current_status" instruction
	// and sends it anyway, Parse must not choke on the unknown field.
	raw := `{
		"current_project": ["Spring Boot"],
		"build_flow": ["Source", "Build"],
		"why_topic": ["Consistent environments"],
		"current_status": [{"label": "Dockerfile", "present": true}]
	}`

	if _, err := Parse(raw); err != nil {
		t.Fatalf("Parse should ignore unexpected fields, got error: %v", err)
	}
}

func TestParseMissingField(t *testing.T) {

	raw := `{
		"current_project": ["Spring Boot"],
		"build_flow": ["Source", "Build"],
		"why_topic": []
	}`

	_, err := Parse(raw)
	if err == nil {
		t.Fatal("expected an error for missing/empty required fields, got nil")
	}
}

func TestParseMalformedJSON(t *testing.T) {

	_, err := Parse("not json at all")
	if err == nil {
		t.Fatal("expected an error for malformed JSON, got nil")
	}
}

func TestParseWhitespaceTolerant(t *testing.T) {

	raw := `

	{"current_project": ["a"], "build_flow": ["b"], "why_topic": ["c"]}

	`

	if _, err := Parse(raw); err != nil {
		t.Fatalf("Parse should tolerate surrounding whitespace, got error: %v", err)
	}
}
