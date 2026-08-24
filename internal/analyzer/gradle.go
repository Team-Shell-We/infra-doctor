package analyzer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeGradle(buildFile string) (
	*project.FrameworkInfo,
	*project.DependencyInfo,
	*project.DatabaseInfo,
	error,
) {

	content, err := os.ReadFile(buildFile)
	if err != nil {
		return nil, nil, nil, err
	}

	text := string(content)
	lowerText := strings.ToLower(text)

	framework := &project.FrameworkInfo{
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

	dependency := &project.DependencyInfo{}

	database := &project.DatabaseInfo{}

	// --------------------------------------------------------
	// Spring Boot
	// --------------------------------------------------------

	// ['"]?\)? : Kotlin DSL은 plugin id 뒤에 닫는 괄호가 온다
	// (id("org.springframework.boot") version "x" vs Groovy의 id 'x' version 'y')
	springRegex := regexp.MustCompile(`org\.springframework\.boot['"]?\)?\s*version\s*['"]([0-9.]+)['"]`)
	if match := springRegex.FindStringSubmatch(text); len(match) == 2 {
		framework.SpringBoot.Enabled = true
		framework.SpringBoot.Version = match[1]
	}

	// --------------------------------------------------------
	// Java
	// --------------------------------------------------------

	javaToolchainRegex := regexp.MustCompile(`JavaLanguageVersion\.of\((\d+)\)`)
	javaSourceCompatRegex := regexp.MustCompile(`sourceCompatibility\s*=\s*['"]?(\d+)`)

	if match := javaToolchainRegex.FindStringSubmatch(text); len(match) == 2 {
		framework.Java.Version = match[1]
	} else if match := javaSourceCompatRegex.FindStringSubmatch(text); len(match) == 2 {
		framework.Java.Version = match[1]
	}

	// --------------------------------------------------------
	// Dependencies
	// --------------------------------------------------------

	if strings.Contains(text, "spring-boot-starter-security") {
		dependency.Security.Enabled = true
	}

	if strings.Contains(text, "spring-boot-starter-data-jpa") {
		dependency.JPA.Enabled = true
	}

	if strings.Contains(text, "spring-kafka") {
		dependency.Kafka.Enabled = true
	}

	if strings.Contains(text, "software.amazon.awssdk") ||
		strings.Contains(text, "spring-cloud-starter-aws") {
		dependency.AWS.Enabled = true
	}

	if strings.Contains(lowerText, "lombok") {
		dependency.Lombok.Enabled = true
	}

	if strings.Contains(text, "spring-boot-starter-actuator") {
		dependency.Actuator.Enabled = true
	}

	if strings.Contains(text, "springdoc-openapi") {
		dependency.OpenAPI.Enabled = true
	}

	// --------------------------------------------------------
	// Database
	// --------------------------------------------------------

	switch {
	case strings.Contains(lowerText, "postgresql"):
		database.Primary.Type = "PostgreSQL"

	case strings.Contains(lowerText, "mysql"):
		database.Primary.Type = "MySQL"

	case strings.Contains(lowerText, "mariadb"):
		database.Primary.Type = "MariaDB"

	default:
		database.Primary.Type = "Unknown"
	}

	if strings.Contains(lowerText, "spring-boot-starter-data-redis") ||
		strings.Contains(lowerText, "spring-data-redis") {

		database.Redis = &project.RedisInfo{
			Enabled: true,
		}
	}

	return framework, dependency, database, nil
}
