package analyzer

import (
	"os"
	// regexp 패키지 : 정규표현식을 사용한 문자열 검색, 매칭, 대체 기능 제공
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
		BuildTool:  "Gradle",
		SpringBoot: true,
	}

	// Spring Boot Version
	springRegex := regexp.MustCompile(`org\.springframework\.boot['"]?\s*version\s*['"]([0-9.]+)['"]`)
	if match := springRegex.FindStringSubmatch(text); len(match) == 2 {
		info.SpringVersion = match[1]
	}

	// Java Version
	javaRegex := regexp.MustCompile(`JavaLanguageVersion\.of\((\d+)\)`)
	if match := javaRegex.FindStringSubmatch(text); len(match) == 2 {
		info.JavaVersion = match[1]
	}

	return info, nil
}