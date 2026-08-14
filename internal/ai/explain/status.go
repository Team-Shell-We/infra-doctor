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
