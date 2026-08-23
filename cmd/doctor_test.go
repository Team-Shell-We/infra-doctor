package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/doctor"
)

func TestWriteDoctorJSONProducesValidResult(t *testing.T) {

	result := &doctor.Result{
		Score: 42,
		Diagnoses: []doctor.Diagnosis{
			{Category: doctor.CICD, Level: doctor.Warning, ScoreImpact: -15, Title: "t", Message: "m", Reason: "r", Fix: "f"},
		},
		Checks: []doctor.Check{
			{ID: "docker", Name: "Docker", Category: doctor.Infrastructure, Passed: true},
		},
	}

	var buf bytes.Buffer
	if err := writeDoctorJSON(&buf, result); err != nil {
		t.Fatalf("writeDoctorJSON failed: %v", err)
	}

	var decoded doctor.Result
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}

	if decoded.Score != 42 {
		t.Errorf("Score = %d, want 42", decoded.Score)
	}
	if len(decoded.Diagnoses) != 1 || decoded.Diagnoses[0].Fix != "f" {
		t.Errorf("Diagnoses = %+v", decoded.Diagnoses)
	}
	if len(decoded.Checks) != 1 || decoded.Checks[0].ID != "docker" {
		t.Errorf("Checks = %+v", decoded.Checks)
	}
}

func TestDoctorShouldFail(t *testing.T) {

	cases := []struct {
		name         string
		score        int
		failUnder    int
		failUnderSet bool
		want         bool
	}{
		{"below threshold", 15, 70, true, true},
		{"at threshold", 70, 70, true, false},
		{"above threshold", 90, 70, true, false},
		{"flag not set, score would fail", 0, 70, false, false},
		{"flag not set, default zero threshold", 15, 0, false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := doctorShouldFail(c.score, c.failUnder, c.failUnderSet); got != c.want {
				t.Errorf("doctorShouldFail(%d, %d, %v) = %v, want %v", c.score, c.failUnder, c.failUnderSet, got, c.want)
			}
		})
	}
}
