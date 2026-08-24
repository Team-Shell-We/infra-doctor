package nextstep

import (
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func TestSuggestDockerWhenMissing(t *testing.T) {

	info := &project.Info{}
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}}
	info.Infrastructure.Docker.Compose = []project.ComposeInfo{{File: "docker-compose.yml"}}

	steps := Suggest(info, false)

	if len(steps) != 1 || steps[0] != "infra-doctor generate docker" {
		t.Errorf("Suggest = %v, want [infra-doctor generate docker]", steps)
	}
}

func TestSuggestComposeWhenMissingAndNotKubernetes(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Enabled = true
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}}

	steps := Suggest(info, false)

	if len(steps) != 1 || steps[0] != "infra-doctor generate compose" {
		t.Errorf("Suggest = %v, want [infra-doctor generate compose]", steps)
	}
}

func TestSuggestK8sWhenMissingAndWantsKubernetes(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Enabled = true
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}}

	steps := Suggest(info, true)

	if len(steps) != 1 || steps[0] != "infra-doctor generate k8s" {
		t.Errorf("Suggest = %v, want [infra-doctor generate k8s]", steps)
	}
}

// 회귀 테스트: wantsKubernetes=true면 compose 부재는 체크하지 않는다(반대로
// wantsKubernetes=false면 k8s 부재도 체크하지 않는다) — 두 전략은 배타적.
func TestSuggestDoesNotMixStrategies(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Enabled = true
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}}

	k8sSteps := Suggest(info, true)
	for _, s := range k8sSteps {
		if s == "infra-doctor generate compose" {
			t.Errorf("wantsKubernetes=true should not suggest compose, got %v", k8sSteps)
		}
	}

	composeSteps := Suggest(info, false)
	for _, s := range composeSteps {
		if s == "infra-doctor generate k8s" {
			t.Errorf("wantsKubernetes=false should not suggest k8s, got %v", composeSteps)
		}
	}
}

func TestSuggestCIWhenMissing(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Enabled = true
	info.Infrastructure.Docker.Compose = []project.ComposeInfo{{File: "docker-compose.yml"}}

	steps := Suggest(info, false)

	if len(steps) != 1 || steps[0] != "infra-doctor generate ci" {
		t.Errorf("Suggest = %v, want [infra-doctor generate ci]", steps)
	}
}

func TestSuggestFallbackWhenNothingMissing(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Enabled = true
	info.Infrastructure.Docker.Compose = []project.ComposeInfo{{File: "docker-compose.yml"}}
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}}

	steps := Suggest(info, false)

	if len(steps) != 1 || steps[0] != Fallback {
		t.Errorf("Suggest = %v, want [%s]", steps, Fallback)
	}
}

func TestSuggestMultipleGapsAllListed(t *testing.T) {

	info := &project.Info{}

	steps := Suggest(info, false)

	want := []string{"infra-doctor generate docker", "infra-doctor generate compose", "infra-doctor generate ci"}
	if len(steps) != len(want) {
		t.Fatalf("Suggest = %v, want %v", steps, want)
	}
	for i := range want {
		if steps[i] != want[i] {
			t.Errorf("step %d = %q, want %q", i, steps[i], want[i])
		}
	}
}
