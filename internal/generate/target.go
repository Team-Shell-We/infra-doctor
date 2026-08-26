package generate

import "fmt"

type Target string

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

// ParseTarget : 사용자가 입력한 값이 Target 타입에 정의된 값인지 확인
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
