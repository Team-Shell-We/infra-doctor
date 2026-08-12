package doctor

import (
	"regexp"
	"strings"
	"testing"
)

var validRuleID = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

var validCategories = map[Category]bool{
	Infrastructure: true,
	CICD:           true,
	Database:       true,
	Monitoring:     true,
	Security:       true,
}

var validLevels = map[Level]bool{
	Info:     true,
	Warning:  true,
	Critical: true,
}

// TestRuleRegistrySchema : internal/doctor/rules/*.yaml에 정의된 모든 룰이 doctor 패키지가 기대하는 스키마를 따르는지 검증
// 기여자의 룰 PR이 통과해야 하는 자동 검증
func TestRuleRegistrySchema(t *testing.T) {

	registry, err := LoadRules()
	if err != nil {
		t.Fatalf("LoadRules() failed: %v", err)
	}

	groups := map[string]map[string]Diagnosis{
		"deployment": registry.Deployment,
		"production": registry.Production,
		"localdev":   registry.LocalDev,
	}

	seenIDs := make(map[string]string)

	for group, rules := range groups {

		if len(rules) == 0 {
			t.Errorf("rule group %q has no rules loaded", group)
		}

		for id, diagnosis := range rules {

			validateRule(t, group, id, diagnosis)

			if firstGroup, ok := seenIDs[id]; ok {
				t.Errorf("rule id %q is used in both %q and %q; ids should be unique across all rule files", id, firstGroup, group)
			}

			seenIDs[id] = group
		}
	}
}

func validateRule(t *testing.T, group, id string, d Diagnosis) {

	t.Helper()

	if !validRuleID.MatchString(id) {
		t.Errorf("[%s/%s] id must be snake_case (match %s)", group, id, validRuleID.String())
	}

	if !validCategories[d.Category] {
		t.Errorf("[%s/%s] category %q is not a recognized Category", group, id, d.Category)
	}

	if !validLevels[d.Level] {
		t.Errorf("[%s/%s] level %q is not a recognized Level", group, id, d.Level)
	}

	if d.ScoreImpact >= 0 {
		t.Errorf("[%s/%s] score must be negative, got %d", group, id, d.ScoreImpact)
	}

	requiredText := map[string]string{
		"title":   d.Title,
		"message": d.Message,
		"reason":  d.Reason,
		"fix":     d.Fix,
	}

	for field, value := range requiredText {
		if strings.TrimSpace(value) == "" {
			t.Errorf("[%s/%s] %s must not be empty", group, id, field)
		}
	}
}
