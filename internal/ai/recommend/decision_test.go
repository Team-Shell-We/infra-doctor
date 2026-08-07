package recommend

import (
	"reflect"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func TestDecideSimpleProjectRecommendsCompose(t *testing.T) {

	decision := Decide(&project.Info{})

	if decision.Recommended != "Docker Compose" {
		t.Errorf("Recommended = %q, want %q", decision.Recommended, "Docker Compose")
	}

	if decision.KubernetesFit {
		t.Error("KubernetesFit should be false for a simple project")
	}
}

func TestDecideExistingKubernetesIsKept(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Kubernetes.Enabled = true

	decision := Decide(info)

	if decision.Recommended != "Kubernetes" || !decision.KubernetesFit {
		t.Errorf("expected Kubernetes to be recommended when manifests already exist, got %+v", decision)
	}
}

func TestDecideHighComplexityRecommendsKubernetes(t *testing.T) {

	info := &project.Info{}
	info.Dependencies.Kafka.Enabled = true
	info.Database.Primary.Type = "PostgreSQL"
	info.Database.Redis = &project.RedisInfo{Enabled: true}
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}, {Name: "deploy"}}

	decision := Decide(info)

	if decision.Recommended != "Kubernetes" {
		t.Errorf("expected Kubernetes for a high-complexity project, got %q", decision.Recommended)
	}
}

// 같은 입력이면 항상 같은 결정이 나와야 한다 — AI 호출 없이 순수 함수로만
// 판단하는 게 이 설계의 핵심이므로, 결정론성 자체를 테스트로 못박아 둔다.
func TestDecideIsDeterministic(t *testing.T) {

	info := &project.Info{}
	info.Dependencies.Kafka.Enabled = true

	first := Decide(info)
	second := Decide(info)

	if !reflect.DeepEqual(first, second) {
		t.Errorf("Decide should be deterministic, got %+v then %+v", first, second)
	}
}

func TestNextStepsSuggestsDockerWhenMissing(t *testing.T) {

	info := &project.Info{}
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}}

	steps := NextSteps(info, Decision{Recommended: "Docker Compose"})

	if len(steps) != 1 || steps[0] != "infra-doctor generate docker" {
		t.Errorf("NextSteps = %v, want [infra-doctor generate docker]", steps)
	}
}

func TestNextStepsSuggestsCIWhenMissing(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Compose = []project.ComposeInfo{{File: "docker-compose.yml"}}

	steps := NextSteps(info, Decision{Recommended: "Docker Compose"})

	if len(steps) != 1 || steps[0] != "infra-doctor generate ci" {
		t.Errorf("NextSteps = %v, want [infra-doctor generate ci]", steps)
	}
}

func TestNextStepsFallbackWhenNothingMissing(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Docker.Compose = []project.ComposeInfo{{File: "docker-compose.yml"}}
	info.Github.Workflows = []project.WorkflowInfo{{Name: "ci"}}

	steps := NextSteps(info, Decision{Recommended: "Docker Compose"})

	if len(steps) != 1 || steps[0] != "infra-doctor doctor" {
		t.Errorf("NextSteps = %v, want [infra-doctor doctor]", steps)
	}
}
