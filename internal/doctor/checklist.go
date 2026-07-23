package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type Check struct {
	Name     string
	Category Category
	Passed   bool
}

// Checklist reflects the current pass/fail state of the fixed set of
// infrastructure checks shown in the "Infrastructure Check" section.
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
