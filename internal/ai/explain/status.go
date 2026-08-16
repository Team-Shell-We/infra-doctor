package explain

import "github.com/Team-Shell-We/infra-doctor/internal/project"

// StatusItem : "Current Status" 섹션의 한 줄, 예: {Label: "Dockerfile", Present: true}.
type StatusItem struct {
	Label   string
	Present bool
}

func BuildStatus(topic string, info *project.Info) []StatusItem {

	switch topic {

	case "docker":
		return []StatusItem{
			{"Dockerfile", info.Infrastructure.Docker.Enabled},
			{"Docker Compose", len(info.Infrastructure.Docker.Compose) > 0},
			{"Health Check", info.Infrastructure.HealthCheck.Enabled},
		}

	case "compose":
		return []StatusItem{
			{"docker-compose.yml", len(info.Infrastructure.Docker.Compose) > 0},
			{"Health Check", info.Infrastructure.HealthCheck.Enabled},
		}

	case "container", "image":
		return []StatusItem{
			{"Dockerfile", info.Infrastructure.Docker.Enabled},
		}

	case "github-actions":
		return []StatusItem{
			{"GitHub Actions workflow", len(info.Github.Workflows) > 0},
		}

	case "k8s":
		return []StatusItem{
			{"Kubernetes manifests", info.Infrastructure.Kubernetes.Enabled},
		}

	case "nginx":
		return []StatusItem{
			{"Nginx configuration", info.Infrastructure.Nginx.Enabled},
		}

	case "postgres":
		return []StatusItem{
			{"PostgreSQL", info.Database.Primary.Type == "PostgreSQL"},
		}

	case "rds":
		return []StatusItem{
			{"AWS SDK dependency", info.Dependencies.AWS.Enabled},
			{"Relational database (PostgreSQL/MySQL)", isRelationalDatabase(info.Database.Primary.Type)},
		}

	case "redis":
		return []StatusItem{
			{"Redis", info.Database.Redis != nil},
		}
	}

	return nil
}

func isRelationalDatabase(dbType string) bool {
	return dbType == "PostgreSQL" || dbType == "MySQL" || dbType == "MariaDB"
}

// TopicPresent : topic이 이 프로젝트에 이미 도입돼 있는지. 각 topic의 첫
// StatusItem을 그 topic의 핵심 신호로 취급한다(docker→Dockerfile,
// k8s→Kubernetes manifests 등, BuildStatus의 항목 순서와 일치). 아직
// 도입 안 된 topic을 explain할 때 cmd/explain.go가 안내 배너를 보여줄지
// 판단하는 데 쓴다.
func TopicPresent(status []StatusItem) bool {
	return len(status) > 0 && status[0].Present
}
