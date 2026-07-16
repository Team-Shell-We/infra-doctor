package analyzer

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeGradle(buildFile string) (*project.FrameworkInfo, error) {

	content, err := os.ReadFile(buildFile)
	if err != nil {
		return nil, err
	}

	text := string(content)

	info := &project.FrameworkInfo{
		BuildTool: project.BuildToolInfo{
			Type: "Gradle",
			File: filepath.Base(buildFile),
			Path: buildFile,
		},
		SpringBoot: project.SpringBootInfo{
			Enabled: false,
		},
		Java: project.JavaInfo{},
	}

	// Spring Boot Version
	springRegex := regexp.MustCompile(`org\.springframework\.boot['"]?\s*version\s*['"]([0-9.]+)['"]`)
	if match := springRegex.FindStringSubmatch(text); len(match) == 2 {
		info.SpringBoot.Enabled = true
		info.SpringBoot.Version = match[1]
	}

	// Java Version
	javaRegex := regexp.MustCompile(`JavaLanguageVersion\.of\((\d+)\)`)
	if match := javaRegex.FindStringSubmatch(text); len(match) == 2 {
		info.Java.Version = match[1]
	}

	return info, nil
}