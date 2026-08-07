package recommend

import "github.com/Team-Shell-We/infra-doctor/internal/project"

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
// 복잡도 신호(complexitySignals)가 감지되면 Kubernetes를 추천한다.
const complexityThreshold = 3

// Decide는 스캔된 프로젝트의 배포 전략을 결정한다.
func Decide(info *project.Info) Decision {

	if info.Infrastructure.Kubernetes.Enabled {
		return Decision{
			Recommended:   "Kubernetes",
			KubernetesFit: true,
			Reasons:       []string{"Kubernetes manifests already present in this project"},
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
