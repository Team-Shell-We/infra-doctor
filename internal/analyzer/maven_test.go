package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeMavenDetectsSpringBootFromParent(t *testing.T) {

	root := t.TempDir()
	pomFile := filepath.Join(root, "pom.xml")

	content := `<project>
	<modelVersion>4.0.0</modelVersion>
	<parent>
		<groupId>org.springframework.boot</groupId>
		<artifactId>spring-boot-starter-parent</artifactId>
		<version>3.2.0</version>
	</parent>
	<artifactId>demo</artifactId>
	<properties>
		<java.version>17</java.version>
	</properties>
</project>
`
	if err := os.WriteFile(pomFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	framework, _, _, err := AnalyzeMaven(pomFile)
	if err != nil {
		t.Fatalf("AnalyzeMaven failed: %v", err)
	}

	if !framework.SpringBoot.Enabled || framework.SpringBoot.Version != "3.2.0" {
		t.Errorf("SpringBoot = %+v, want Enabled=true Version=3.2.0", framework.SpringBoot)
	}
	if framework.Java.Version != "17" {
		t.Errorf("Java.Version = %q, want 17", framework.Java.Version)
	}
}

// dependencyManagement로 spring-boot-dependencies BOM만 import하고
// 부모/직접 버전이 없는 멀티모듈 루트 pom 패턴
func TestAnalyzeMavenDetectsSpringBootFromBOMImport(t *testing.T) {

	root := t.TempDir()
	pomFile := filepath.Join(root, "pom.xml")

	content := `<project>
	<modelVersion>4.0.0</modelVersion>
	<artifactId>demo</artifactId>
	<dependencyManagement>
		<dependencies>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-dependencies</artifactId>
				<version>3.2.0</version>
				<type>pom</type>
				<scope>import</scope>
			</dependency>
		</dependencies>
	</dependencyManagement>
	<dependencies>
		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-web</artifactId>
		</dependency>
	</dependencies>
</project>
`
	if err := os.WriteFile(pomFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	framework, _, _, err := AnalyzeMaven(pomFile)
	if err != nil {
		t.Fatalf("AnalyzeMaven failed: %v", err)
	}

	if !framework.SpringBoot.Enabled || framework.SpringBoot.Version != "3.2.0" {
		t.Errorf("SpringBoot = %+v, want Enabled=true Version=3.2.0 (from BOM)", framework.SpringBoot)
	}
}

// 자식 모듈의 <dependency>가 BOM이 관리하는 버전을 쓰느라 자체 <version>이 없는 경우
func TestAnalyzeMavenResolvesVersionFromDependencyManagement(t *testing.T) {

	root := t.TempDir()
	pomFile := filepath.Join(root, "pom.xml")

	content := `<project>
	<modelVersion>4.0.0</modelVersion>
	<artifactId>demo</artifactId>
	<dependencyManagement>
		<dependencies>
			<dependency>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-starter-web</artifactId>
				<version>3.2.0</version>
			</dependency>
		</dependencies>
	</dependencyManagement>
	<dependencies>
		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-web</artifactId>
		</dependency>
	</dependencies>
</project>
`
	if err := os.WriteFile(pomFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	framework, _, _, err := AnalyzeMaven(pomFile)
	if err != nil {
		t.Fatalf("AnalyzeMaven failed: %v", err)
	}

	if framework.SpringBoot.Version != "3.2.0" {
		t.Errorf("SpringBoot.Version = %q, want 3.2.0 resolved from dependencyManagement", framework.SpringBoot.Version)
	}
}

func TestAnalyzeMavenModulesCountsModules(t *testing.T) {

	root := t.TempDir()
	pomFile := filepath.Join(root, "pom.xml")

	content := `<project>
	<modelVersion>4.0.0</modelVersion>
	<artifactId>demo-parent</artifactId>
	<packaging>pom</packaging>
	<modules>
		<module>api</module>
		<module>worker</module>
		<module>common</module>
	</modules>
</project>
`
	if err := os.WriteFile(pomFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	count, err := AnalyzeMavenModules(pomFile)
	if err != nil {
		t.Fatalf("AnalyzeMavenModules failed: %v", err)
	}

	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

// 회귀 테스트: AnalyzeProject의 Maven 분기가 인프라/CI/프로파일도 채우는지 확인
func TestAnalyzeProjectMavenIncludesInfrastructureAndCI(t *testing.T) {

	root := t.TempDir()

	pom := `<project>
	<modelVersion>4.0.0</modelVersion>
	<artifactId>demo</artifactId>
	<parent>
		<groupId>org.springframework.boot</groupId>
		<artifactId>spring-boot-starter-parent</artifactId>
		<version>3.2.0</version>
	</parent>
</project>
`
	if err := os.WriteFile(filepath.Join(root, "pom.xml"), []byte(pom), 0644); err != nil {
		t.Fatal(err)
	}

	dockerfile := "FROM eclipse-temurin:17-jre\nHEALTHCHECK CMD curl -f http://localhost:8080/actuator/health || exit 1\n"
	if err := os.WriteFile(filepath.Join(root, "Dockerfile"), []byte(dockerfile), 0644); err != nil {
		t.Fatal(err)
	}

	workflowDir := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatal(err)
	}
	workflow := "name: CI\non:\n  push:\n    branches: [main]\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n"
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte(workflow), 0644); err != nil {
		t.Fatal(err)
	}

	info, err := AnalyzeProject(root)
	if err != nil {
		t.Fatalf("AnalyzeProject failed: %v", err)
	}

	if !info.Infrastructure.Docker.Enabled {
		t.Error("expected Docker to be detected for a Maven project")
	}
	if len(info.Github.Workflows) != 1 {
		t.Errorf("expected 1 GitHub workflow to be detected, got %d", len(info.Github.Workflows))
	}
}
