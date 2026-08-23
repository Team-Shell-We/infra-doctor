package recommend

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

// Decision: 배포 전략 추천 결과. Recommended/KubernetesFit은 project.Info만으로
// 결정되며 AI는 절대 관여하지 않는다(설명만 담당, prompt.go 참고).
// explain/status.go에서 얻은 교훈과 동일 — AI에게 "사실 판단"을 맡기면
// 재현 불가능하고 근거 없는 답이 나올 수 있다.
type Decision struct {
	Recommended   string // "Docker Compose" | "Kubernetes"
	KubernetesFit bool

	// Reasons: AI 설명의 근거가 되는 짧은 사실 라벨. 사용자에게 직접
	// 렌더링하지 않고 AI가 문장으로 풀어쓰는 재료로만 쓰인다.
	Reasons []string
}

// complexityThreshold: 아직 Kubernetes를 쓰지 않는 프로젝트에서, 이 개수 이상의
// 복잡도 신호(complexitySignals)가 감지되면 Kubernetes를 추천한다. 신호가
// 4→6개로 늘어난 만큼(엔드포인트 수, 멀티모듈 구조 추가) 기존 "4개 중 3개"
// 비율에 맞춰 3에서 4로 올렸다.
const complexityThreshold = 4

// replicaScaleThreshold: 이미 Kubernetes를 쓰는 프로젝트에서, replicas가 이
// 값 이상이면 실제로 수평 확장 중이라고 보고 Reasons에 근거를 추가한다.
const replicaScaleThreshold = 2

// Decide는 스캔된 프로젝트의 배포 전략을 결정한다.
func Decide(info *project.Info) Decision {

	if info.Infrastructure.Kubernetes.Enabled {

		reasons := []string{"Kubernetes manifests already present in this project"}

		if info.Infrastructure.Kubernetes.Replicas >= replicaScaleThreshold {
			reasons = append(reasons, fmt.Sprintf("Already scaled to %d replicas", info.Infrastructure.Kubernetes.Replicas))
		}

		return Decision{
			Recommended:   "Kubernetes",
			KubernetesFit: true,
			Reasons:       reasons,
		}
	}

	signals := complexitySignals(info)

	if len(signals) >= complexityThreshold {
		return Decision{
			Recommended:   "Kubernetes",
			KubernetesFit: true,
			Reasons:       signals,
		}
	}

	return Decision{
		Recommended:   "Docker Compose",
		KubernetesFit: false,
		Reasons:       []string{"Single API server", "Low infrastructure complexity"},
	}
}

// endpointCountThreshold, moduleCountThreshold: complexitySignals의 나머지
// 신호와 같은 방식(boolean 존재 여부)이 아니라 개수 기반이라 별도 기준값이
// 필요하다. 둘 다 실측 데이터가 아니라 경험적 추정치.
const (
	endpointCountThreshold = 20
	moduleCountThreshold   = 1
)

// complexitySignals: 인프라 복잡도 신호 목록. 이 도구는 Spring Boot 서비스
// 하나만 스캔하므로 마이크로서비스 개수 같은 건 알 수 없다 — 여기 나열된
// 항목들이 "Compose 단일 호스트로는 부족한 규모"를 가늠하는 가장 가까운 대리 지표다.
func complexitySignals(info *project.Info) []string {

	var signals []string

	if info.Dependencies.Kafka.Enabled {
		signals = append(signals, "Message queue (Kafka) in use")
	}

	if info.Database.Primary.Type != "" && info.Database.Primary.Type != "Unknown" {
		signals = append(signals, "Relational database in use")
	}

	if info.Database.Redis != nil {
		signals = append(signals, "Redis cache in use")
	}

	if len(info.Github.Workflows) > 1 {
		signals = append(signals, "Multiple CI/CD workflows configured")
	}

	if info.API.EndpointCount > endpointCountThreshold {
		signals = append(signals, fmt.Sprintf("Large number of API endpoints (%d)", info.API.EndpointCount))
	}

	if info.Framework.Modules.Count > moduleCountThreshold {
		signals = append(signals, fmt.Sprintf("Multi-module Gradle project (%d modules)", info.Framework.Modules.Count))
	}

	return signals
}

// NextSteps는 결정 결과와 현재 프로젝트에 부족한 부분을 기반으로 다음 명령어를
// 제안한다. AI가 아니라 코드가 결정 — 명령어를 잘못 지어내면(hallucination)
// 사용자에게 그대로 오해를 줄 수 있기 때문이다.
func NextSteps(info *project.Info, decision Decision) []string {

	var steps []string

	if decision.Recommended == "Docker Compose" && len(info.Infrastructure.Docker.Compose) == 0 {
		steps = append(steps, "infra-doctor generate docker")
	}

	if decision.Recommended == "Kubernetes" && !info.Infrastructure.Kubernetes.Enabled {
		steps = append(steps, "infra-doctor generate k8s")
	}

	if len(info.Github.Workflows) == 0 {
		steps = append(steps, "infra-doctor generate ci")
	}

	if len(steps) == 0 {
		steps = append(steps, "infra-doctor doctor")
	}

	return steps
}
