package nextstep

import "github.com/Team-Shell-We/infra-doctor/internal/project"

// Fallback : 제안할 게 없을 때 보여줄 기본 명령어
const Fallback = "infra-doctor doctor"

// Suggest : project.Info에서 부족한 인프라 요소를 보고 generate 명령어를
// 제안한다. wantsKubernetes는 호출자가 이미 판단한 배포 전략(recommend의
// Decision.Recommended == "Kubernetes")을 반영 — doctor는 이 판단을 하지
// 않으므로 항상 false로 호출한다.
func Suggest(info *project.Info, wantsKubernetes bool) []string {

	var steps []string

	if !info.Infrastructure.Docker.Enabled {
		steps = append(steps, "infra-doctor generate docker")
	}

	if wantsKubernetes {
		if !info.Infrastructure.Kubernetes.Enabled {
			steps = append(steps, "infra-doctor generate k8s")
		}
	} else if len(info.Infrastructure.Docker.Compose) == 0 {
		steps = append(steps, "infra-doctor generate compose")
	}

	if len(info.Github.Workflows) == 0 {
		steps = append(steps, "infra-doctor generate ci")
	}

	if len(steps) == 0 {
		steps = append(steps, Fallback)
	}

	return steps
}
