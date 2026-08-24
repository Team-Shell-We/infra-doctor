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

// AnalyzeMavenModules: pom.xml의 <modules><module> 개수를 셈
func AnalyzeMavenModules(buildFile string) (int, error) {
	pom, err := parseMavenProject(buildFile)
	if err != nil {
		return 0, fmt.Errorf("parse Maven project %q: %w", buildFile, err)
	}

	return len(pom.Modules), nil
}

// parseMavenProject: pom.xml을 Maven 모델로 변환
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

// extractSpringBoot: parent, dependencyManagement(BOM), dependency,
// plugin 순서로 Spring Boot 사용 여부와 버전 확인
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

	// spring-boot-dependencies를 <dependencyManagement>로 import하는
	// 프로젝트는 부모나 직접 의존성 없이도 Spring Boot 프로젝트임
	for _, item := range pom.DependencyManagement.Dependencies {
		groupID := resolveMavenValue(item.GroupID, pom.Properties)
		artifactID := resolveMavenValue(item.ArtifactID, pom.Properties)

		if groupID == "org.springframework.boot" &&
			artifactID == "spring-boot-dependencies" {
			result.SpringBoot.Enabled = true

			if result.SpringBoot.Version == "" {
				result.SpringBoot.Version = resolveMavenValue(item.Version, pom.Properties)
			}
		}
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
			version := resolveMavenValue(item.Version, pom.Properties)
			if version == "" {
				// <version>이 없으면 BOM이 관리하는 버전
				version = managedMavenVersion(pom, groupID, item.ArtifactID)
			}
			result.SpringBoot.Version = version
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

// extractMavenDependencies: 프로젝트에서 사용하는 주요 의존성 분석
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

// extractMavenDatabase: 주 데이터베이스와 Redis 사용 여부를 분석
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

// findMavenJavaVersion: properties와 maven-compiler-plugin에서
// Java 버전을 찾음
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

// managedMavenVersion: <dependencyManagement>에서 groupID/artifactID가 일치하는 항목의 버전을 찾음 
// <version> 없이 BOM으로만 버전이 관리되는 의존성의 버전을 찾을 때 사용
func managedMavenVersion(pom *mavenProject, groupID, artifactID string) string {
	artifactID = resolveMavenValue(artifactID, pom.Properties)

	for _, item := range pom.DependencyManagement.Dependencies {
		if resolveMavenValue(item.GroupID, pom.Properties) == groupID &&
			resolveMavenValue(item.ArtifactID, pom.Properties) == artifactID {
			return resolveMavenValue(item.Version, pom.Properties)
		}
	}

	return ""
}

// resolveMavenValue: ${property.name} 형식의 값을 실제 값으로 변환
func resolveMavenValue(
	value string,
	properties mavenProperties,
) string {
	value = strings.TrimSpace(value)

	// 순환 참조에 빠지지 않도록 최대 해석 횟수를 제한
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

// normalizeJavaVersion: 1.8 형식을 8로 변환
func normalizeJavaVersion(value string) string {
	value = strings.TrimSpace(value)
	return strings.TrimPrefix(value, "1.")
}
