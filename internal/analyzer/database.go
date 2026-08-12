package analyzer

import (
	"os"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeDatabase(buildFile string) (*project.DatabaseInfo, error) {

	content, err := os.ReadFile(buildFile)
	if err != nil {
		return nil, err
	}

	text := strings.ToLower(string(content))

	info := &project.DatabaseInfo{}

	switch {
		case strings.Contains(text, "postgresql"):
			info.Primary.Type = "PostgreSQL"

		case strings.Contains(text, "mysql"):
			info.Primary.Type = "MySQL"

		case strings.Contains(text, "mariadb"):
			info.Primary.Type = "MariaDB"

		default:
			info.Primary.Type = "Unknown"
	}

	if strings.Contains(text, "spring-boot-starter-data-redis") ||
		strings.Contains(text, "spring-data-redis") {

		info.Redis = &project.RedisInfo{
			Enabled: true,
		}
	}

	return info, nil
}