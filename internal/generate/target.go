package generate

import "fmt"

type Target string

// 사용자가 입력한 값이 Target 타입에 정의된 값인지 확인하는 함수입니다.
const (
	TargetDocker       Target = "docker"
	TargetCompose      Target = "compose"
	TargetNginx        Target = "nginx"
	TargetCI           Target = "ci"
	TargetK8s          Target = "k8s"
	TargetArchitecture Target = "architecture"
)

func Targets() []Target {
	return []Target{
		TargetDocker,
		TargetCompose,
		TargetNginx,
		TargetCI,
		TargetK8s,
		TargetArchitecture,
	}
}

func ParseTarget(value string) (Target, error) {
	target := Target(value)

	for _, candidate := range Targets() {
		if target == candidate {
			return target, nil
		}
	}

	return "", fmt.Errorf(
		"unknown generate target %q",
		value,
	)
}
