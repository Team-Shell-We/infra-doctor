package recommend

import "testing"

func TestParseValid(t *testing.T) {

	result, err := Parse(`{"reasons": ["Single API server", "Low infrastructure complexity"]}`)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Reasons) != 2 {
		t.Errorf("Reasons = %v, want 2 items", result.Reasons)
	}
}

func TestParseEmptyReasons(t *testing.T) {

	if _, err := Parse(`{"reasons": []}`); err == nil {
		t.Fatal("expected an error for empty reasons, got nil")
	}
}

func TestParseMalformedJSON(t *testing.T) {

	if _, err := Parse("not json"); err == nil {
		t.Fatal("expected an error for malformed JSON, got nil")
	}
}
