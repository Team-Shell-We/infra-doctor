package analyzer

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeProject(root string) (*project.Info, error) {

	info := &project.Info{}

	if HasGradle(root) {
		fmt.Println("✓ Gradle detected")
	}

	return info, nil
} 