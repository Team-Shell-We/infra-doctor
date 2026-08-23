package recommend

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

// Decision: 배포 전략 추천 결과. Recommended/KubernetesFit은 project.Info만으로 결정

type Decision struct {
	Recommended   string // "Docker Compose" | "Kubernetes"
	KubernetesFit bool

	// Reasons: AI 설명의 근거가 되는 짧은 사실 라벨
	Reasons []string
}

// complexityThreshold: 아직 Kubernetes를 쓰지 않는 프로젝트에서, 이 개수 이상의
// 복잡도 신호(complexitySignals)가 감지되면 Kubernetes를 추천
const complexityThreshold = 4

// replicaScaleThreshold: 이미 Kubernetes를 쓰는 프로젝트에서, replicas가 이
// 값 이상이면 실제로 수평 확장 중이라고 보고 Reasons에 근거 추가
const replicaScaleThreshold = 2

// Decide: 스캔된 프로젝트의 배포 전략을 결정
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
// 신호와 같은 방식(boolean 존재 여부)이 아니라 개수 기반이라 별도 기준값이 필요
const (
	endpointCountThreshold = 20
	moduleCountThreshold   = 1
)

// complexitySignals: 인프라 복잡도 신호 목록
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
