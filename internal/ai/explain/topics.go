package explain

// Topics는 `infra-doctor explain <topic>`이 허용하는 고정된 주제 목록
// cobra의 ValidArgs로도 쓰여, 잘못된 topic은 AI 호출 전에 거부
var Topics = []string{
	"compose",
	"container",
	"docker",
	"github-actions",
	"image",
	"k8s",
	"nginx",
	"postgres",
	"rds",
	"redis",
}

var displayNames = map[string]string{
	"compose":        "Docker Compose",
	"container":      "Container",
	"docker":         "Docker",
	"github-actions": "GitHub Actions",
	"image":          "Docker Image",
	"k8s":            "Kubernetes",
	"nginx":          "Nginx",
	"postgres":       "PostgreSQL",
	"rds":            "Amazon RDS",
	"redis":          "Redis",
}

// DisplayName은 topic의 사람이 읽기 좋은 이름을 반환한다(예: "k8s" ->
// "Kubernetes"). 알 수 없는 topic이면 원래 문자열을 그대로 반환한다.
func DisplayName(topic string) string {

	if name, ok := displayNames[topic]; ok {
		return name
	}

	return topic
}
