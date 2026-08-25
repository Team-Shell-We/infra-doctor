package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type Check struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Category Category `json:"category"`
	Passed   bool     `json:"passed"`
}

// Checklist : "Infrastructure Check" 섹션에 표시되는 고정된 인프라 체크
// 항목들의 현재 통과/실패 상태를 반환
func Checklist(info *project.Info) []Check {

	return []Check{
		{ID: "docker", Name: "Docker", Category: Infrastructure, Passed: info.Infrastructure.Docker.Enabled},
		{ID: "docker_compose", Name: "Docker Compose", Category: Infrastructure, Passed: len(info.Infrastructure.Docker.Compose) > 0},
		{ID: "health_check", Name: "Health Check", Category: Infrastructure, Passed: info.Infrastructure.HealthCheck.Enabled},
		{ID: "reverse_proxy", Name: "Reverse Proxy", Category: Infrastructure, Passed: info.Infrastructure.Nginx.Enabled},
		{ID: "monitoring", Name: "Monitoring", Category: Monitoring, Passed: info.Infrastructure.Monitoring.Enabled},
		{ID: "log_rotation", Name: "Log Rotation", Category: Infrastructure, Passed: info.Infrastructure.LogRotation.Enabled},
		{ID: "db_backup", Name: "DB Backup", Category: Database, Passed: info.Infrastructure.Backup.Enabled},
	}
}
