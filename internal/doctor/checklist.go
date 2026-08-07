package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type Check struct {
	Name     string
	Category Category
	Passed   bool
}

// Checklist : "Infrastructure Check" 섹션에 표시되는 고정된 인프라 체크
// 항목들의 현재 통과/실패 상태를 반환
func Checklist(info *project.Info) []Check {

	return []Check{
		{Name: "Docker", Category: Infrastructure, Passed: info.Infrastructure.Docker.Enabled},
		{Name: "Docker Compose", Category: Infrastructure, Passed: len(info.Infrastructure.Docker.Compose) > 0},
		{Name: "Health Check", Category: Infrastructure, Passed: info.Infrastructure.HealthCheck.Enabled},
		{Name: "Reverse Proxy", Category: Infrastructure, Passed: info.Infrastructure.Nginx.Enabled},
		{Name: "Monitoring", Category: Monitoring, Passed: info.Infrastructure.Monitoring.Enabled},
		{Name: "Log Rotation", Category: Infrastructure, Passed: info.Infrastructure.LogRotation.Enabled},
		{Name: "DB Backup", Category: Database, Passed: info.Infrastructure.Backup.Enabled},
	}
}
