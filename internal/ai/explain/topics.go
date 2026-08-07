package explain

// Topics is the fixed set of subjects `infra-doctor explain <topic>`
// accepts, matching the feature spec exactly. It doubles as cobra's
// ValidArgs list so an invalid topic is rejected before any AI call.
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

// DisplayName returns the human-readable name for a topic (e.g. "k8s" ->
// "Kubernetes"), falling back to the raw topic if it's somehow unknown.
func DisplayName(topic string) string {

	if name, ok := displayNames[topic]; ok {
		return name
	}

	return topic
}
