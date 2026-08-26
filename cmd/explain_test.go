package cmd

import "testing"

func TestExplainArgsAcceptsValidTopic(t *testing.T) {

	if err := explainArgs(explainCmd, []string{"docker"}); err != nil {
		t.Errorf("explainArgs([docker]) = %v, want nil", err)
	}
}

func TestExplainArgsAcceptsValidTopicWithPath(t *testing.T) {

	if err := explainArgs(explainCmd, []string{"docker", "./some/path"}); err != nil {
		t.Errorf("explainArgs([docker, path]) = %v, want nil", err)
	}
}

func TestExplainArgsRejectsUnknownTopic(t *testing.T) {

	if err := explainArgs(explainCmd, []string{"not-a-real-topic"}); err == nil {
		t.Error("explainArgs([not-a-real-topic]) = nil, want error")
	}
}

func TestExplainArgsRejectsWrongArgCount(t *testing.T) {

	cases := [][]string{
		{},
		{"docker", "path", "extra"},
	}

	for _, args := range cases {
		if err := explainArgs(explainCmd, args); err == nil {
			t.Errorf("explainArgs(%v) = nil, want error", args)
		}
	}
}
