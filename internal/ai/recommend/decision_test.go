package recommend

import (
	"reflect"
	"strings"
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

// 회귀 테스트: 이미 Kubernetes를 쓰는 프로젝트가 replicas도 높으면, 그냥
// "manifest가 있다"보다 더 구체적인 근거(실제로 스케일 중이라는 사실)를 Reasons에 추가해야 함
func TestDecideExistingKubernetesWithHighReplicasAddsReason(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Kubernetes.Enabled = true
	info.Infrastructure.Kubernetes.Replicas = 5

	decision := Decide(info)

	found := false
	for _, reason := range decision.Reasons {
		if reason == "Already scaled to 5 replicas" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a replica-scale reason, got %+v", decision.Reasons)
	}
}

func TestDecideExistingKubernetesWithLowReplicasOmitsReason(t *testing.T) {

	info := &project.Info{}
	info.Infrastructure.Kubernetes.Enabled = true
	info.Infrastructure.Kubernetes.Replicas = 1

	decision := Decide(info)

	if len(decision.Reasons) != 1 {
		t.Errorf("expected no replica-scale reason for a single replica, got %+v", decision.Reasons)
	}
}

// 회귀 테스트: API 엔드포인트 수/멀티모듈 구조도 복잡도 신호로 잡혀야 함
func TestDecideCountsEndpointAndModuleSignals(t *testing.T) {

	info := &project.Info{}
	info.API.EndpointCount = endpointCountThreshold + 1
	info.Framework.Modules.Count = moduleCountThreshold + 1
	info.Dependencies.Kafka.Enabled = true
	info.Database.Redis = &project.RedisInfo{Enabled: true}

	decision := Decide(info)

	if decision.Recommended != "Kubernetes" {
		t.Errorf("expected Kubernetes once endpoint+module signals push past the threshold, got %q (reasons: %v)", decision.Recommended, decision.Reasons)
	}
}

// 회귀 테스트: 멀티모듈 신호 라벨이 빌드 도구를 하드코딩하면 안 됨
func TestDecideMultiModuleSignalDoesNotHardcodeBuildTool(t *testing.T) {

	info := &project.Info{}
	info.Framework.BuildTool.Type = "Maven"
	info.Framework.Modules.Count = moduleCountThreshold + 1

	signals := complexitySignals(info)

	for _, s := range signals {
		if strings.Contains(s, "Gradle") {
			t.Errorf("multi-module signal must not name a specific build tool, got %q", s)
		}
	}
}

func TestDecideLowEndpointAndModuleCountsDoNotSignal(t *testing.T) {

	info := &project.Info{}
	info.API.EndpointCount = 1
	info.Framework.Modules.Count = 1

	signals := complexitySignals(info)

	for _, s := range signals {
		if s == "Large number of API endpoints (1)" || s == "Multi-module project (1 modules)" {
			t.Errorf("did not expect a signal from a small single-module project, got %v", signals)
		}
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

// 같은 입력이면 항상 같은 결정이 나와야 함
func TestDecideIsDeterministic(t *testing.T) {

	info := &project.Info{}
	info.Dependencies.Kafka.Enabled = true

	first := Decide(info)
	second := Decide(info)

	if !reflect.DeepEqual(first, second) {
		t.Errorf("Decide should be deterministic, got %+v then %+v", first, second)
	}
}
