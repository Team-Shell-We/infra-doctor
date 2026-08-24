package analyzer

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

// AnalyzeMaven은 pom.xml을 분석해 프레임워크, 주요 의존성,
// 데이터베이스 정보를 반환한다.
func AnalyzeMaven(buildFile string) (
	*project.FrameworkInfo,
	*project.DependencyInfo,
	*project.DatabaseInfo,
	error,
) {
	pom, err := parseMavenProject(buildFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"parse Maven project %q: %w",
			buildFile,
			err,
		)
	}

	framework := &project.FrameworkInfo{
		BuildTool: project.BuildToolInfo{
			Type: "Maven",
			File: filepath.Base(buildFile),
			Path: buildFile,
		},
	}

	dependencies := &project.DependencyInfo{}

	database := &project.DatabaseInfo{
		Primary: project.Database{
			Type: "Unknown",
		},
	}

	extractMavenFramework(pom, framework)
	extractMavenDependencies(pom, dependencies)
	extractMavenDatabase(pom, database)

	return framework, dependencies, database, nil
}

// parseMavenProject는 pom.xml을 Maven 모델로 변환한다.
func parseMavenProject(path string) (*mavenProject, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open pom.xml: %w", err)
	}
	defer file.Close()

	var pom mavenProject

	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&pom); err != nil {
		return nil, fmt.Errorf("decode pom.xml: %w", err)
	}

	if strings.TrimSpace(pom.ArtifactID) == "" {
		return nil, fmt.Errorf("missing project artifactId")
	}

	return &pom, nil
}

// extractMavenFramework는 Spring Boot와 Java 버전을 분석한다.
func extractMavenFramework(
	pom *mavenProject,
	result *project.FrameworkInfo,
) {
	extractSpringBoot(pom, result)
	result.Java.Version = findMavenJavaVersion(pom)
}

// extractSpringBoot는 parent, dependency, plugin 순서로
// Spring Boot 사용 여부와 버전을 확인한다.
func extractSpringBoot(
	pom *mavenProject,
	result *project.FrameworkInfo,
) {
	parentGroupID := resolveMavenValue(
		pom.Parent.GroupID,
		pom.Properties,
	)
	parentArtifactID := resolveMavenValue(
		pom.Parent.ArtifactID,
		pom.Properties,
	)

	if parentGroupID == "org.springframework.boot" &&
		parentArtifactID == "spring-boot-starter-parent" {
		result.SpringBoot.Enabled = true
		result.SpringBoot.Version = resolveMavenValue(
			pom.Parent.Version,
			pom.Properties,
		)
	}

	for _, item := range pom.Dependencies {
		groupID := resolveMavenValue(
			item.GroupID,
			pom.Properties,
		)

		if groupID != "org.springframework.boot" {
			continue
		}

		result.SpringBoot.Enabled = true

		if result.SpringBoot.Version == "" {
			result.SpringBoot.Version = resolveMavenValue(
				item.Version,
				pom.Properties,
			)
		}
	}

	for _, plugin := range pom.Plugins {
		groupID := resolveMavenValue(
			plugin.GroupID,
			pom.Properties,
		)
		artifactID := resolveMavenValue(
			plugin.ArtifactID,
			pom.Properties,
		)

		if groupID == "org.springframework.boot" &&
			artifactID == "spring-boot-maven-plugin" {
			result.SpringBoot.Enabled = true

			if result.SpringBoot.Version == "" {
				result.SpringBoot.Version = resolveMavenValue(
					plugin.Version,
					pom.Properties,
				)
			}
		}
	}
}

// extractMavenDependencies는 프로젝트에서 사용하는 주요 의존성을 분석한다.
func extractMavenDependencies(
	pom *mavenProject,
	result *project.DependencyInfo,
) {
	for _, item := range pom.Dependencies {
		groupID := strings.ToLower(resolveMavenValue(
			item.GroupID,
			pom.Properties,
		))
		artifactID := strings.ToLower(resolveMavenValue(
			item.ArtifactID,
			pom.Properties,
		))

		switch {
		case artifactID == "spring-boot-starter-security":
			result.Security.Enabled = true

		case artifactID == "spring-boot-starter-data-jpa":
			result.JPA.Enabled = true

		case artifactID == "spring-kafka":
			result.Kafka.Enabled = true

		case strings.HasPrefix(groupID, "software.amazon.awssdk"):
			result.AWS.Enabled = true

		case strings.Contains(
			artifactID,
			"spring-cloud-starter-aws",
		):
			result.AWS.Enabled = true

		case artifactID == "lombok":
			result.Lombok.Enabled = true

		case artifactID == "spring-boot-starter-actuator":
			result.Actuator.Enabled = true

		case strings.Contains(artifactID, "springdoc-openapi"):
			result.OpenAPI.Enabled = true
		}
	}
}

// extractMavenDatabase는 주 데이터베이스와 Redis 사용 여부를 분석한다.
func extractMavenDatabase(
	pom *mavenProject,
	result *project.DatabaseInfo,
) {
	for _, item := range pom.Dependencies {
		groupID := strings.ToLower(resolveMavenValue(
			item.GroupID,
			pom.Properties,
		))
		artifactID := strings.ToLower(resolveMavenValue(
			item.ArtifactID,
			pom.Properties,
		))

		switch {
		case groupID == "org.postgresql" &&
			artifactID == "postgresql":
			result.Primary.Type = "PostgreSQL"

		case groupID == "com.mysql" &&
			artifactID == "mysql-connector-j":
			result.Primary.Type = "MySQL"

		case groupID == "mysql" &&
			artifactID == "mysql-connector-java":
			result.Primary.Type = "MySQL"

		case groupID == "org.mariadb.jdbc" &&
			artifactID == "mariadb-java-client":
			result.Primary.Type = "MariaDB"
		}

		if isRedisDependency(groupID, artifactID) {
			result.Redis = &project.RedisInfo{
				Enabled: true,
			}
		}
	}
}

func isRedisDependency(groupID, artifactID string) bool {
	switch {
	case groupID == "org.springframework.boot" &&
		artifactID == "spring-boot-starter-data-redis":
		return true

	case groupID == "org.springframework.data" &&
		artifactID == "spring-data-redis":
		return true

	case groupID == "redis.clients" &&
		artifactID == "jedis":
		return true

	case groupID == "io.lettuce" &&
		artifactID == "lettuce-core":
		return true

	default:
		return false
	}
}

// findMavenJavaVersion은 properties와 maven-compiler-plugin에서
// Java 버전을 찾는다.
func findMavenJavaVersion(pom *mavenProject) string {
	propertyNames := []string{
		"java.version",
		"maven.compiler.release",
		"maven.compiler.target",
		"maven.compiler.source",
	}

	for _, name := range propertyNames {
		value := resolveMavenValue(
			pom.Properties[name],
			pom.Properties,
		)

		if value != "" {
			return normalizeJavaVersion(value)
		}
	}

	for _, plugin := range pom.Plugins {
		artifactID := resolveMavenValue(
			plugin.ArtifactID,
			pom.Properties,
		)

		if artifactID != "maven-compiler-plugin" {
			continue
		}

		candidates := []string{
			plugin.Configuration.Release,
			plugin.Configuration.Target,
			plugin.Configuration.Source,
		}

		for _, candidate := range candidates {
			value := resolveMavenValue(
				candidate,
				pom.Properties,
			)

			if value != "" {
				return normalizeJavaVersion(value)
			}
		}
	}

	return ""
}

// resolveMavenValue는 ${property.name} 형식의 값을 실제 값으로 변환한다.
func resolveMavenValue(
	value string,
	properties mavenProperties,
) string {
	value = strings.TrimSpace(value)

	// 순환 참조에 빠지지 않도록 최대 해석 횟수를 제한한다.
	for range 10 {
		if !strings.HasPrefix(value, "${") ||
			!strings.HasSuffix(value, "}") {
			break
		}

		propertyName := strings.TrimSuffix(
			strings.TrimPrefix(value, "${"),
			"}",
		)

		resolved, exists := properties[propertyName]
		if !exists {
			break
		}

		resolved = strings.TrimSpace(resolved)
		if resolved == value {
			break
		}

		value = resolved
	}

	return value
}

// normalizeJavaVersion은 1.8 형식을 8로 변환한다.
func normalizeJavaVersion(value string) string {
	value = strings.TrimSpace(value)
	return strings.TrimPrefix(value, "1.")
}