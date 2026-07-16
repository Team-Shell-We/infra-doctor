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
		info.Type = "PostgreSQL"

	case strings.Contains(text, "mysql"):
		info.Type = "MySQL"

	case strings.Contains(text, "mariadb"):
		info.Type = "MariaDB"

	default:
		info.Type = "Unknown"
	}

	if strings.Contains(text, "spring-boot-starter-data-redis") ||
		strings.Contains(text, "spring-data-redis") {
		info.Redis = true
	}

	return info, nil
}