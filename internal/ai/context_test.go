package ai

import (
	"strings"
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func TestBuildSummary(t *testing.T) {

	info := &project.Info{}
	info.Framework.SpringBoot.Enabled = true
	info.Framework.SpringBoot.Version = "3.5.7"
	info.Framework.BuildTool.Type = "Gradle"
	info.Framework.Java.Version = "17"
	info.Dependencies.Security.Enabled = true
	info.Dependencies.JPA.Enabled = true
	info.Database.Primary.Type = "PostgreSQL"
	info.Database.Redis = &project.RedisInfo{Enabled: true}
	info.Infrastructure.Docker.Enabled = true
	info.Infrastructure.Docker.Compose = []project.ComposeInfo{{File: "docker-compose.yml"}}
	info.Github.Workflows = []project.WorkflowInfo{{Name: "CI/CD to EC2"}}
	info.Profiles = []project.ProfileInfo{{Name: "dev"}, {Name: "prod"}}

	summary := BuildSummary(info)

	if summary.Framework != "Spring Boot 3.5.7 (Gradle), Java 17" {
		t.Errorf("unexpected Framework: %q", summary.Framework)
	}

	wantDeps := []string{"Spring Security", "Spring Data JPA"}
	if !equalStrings(summary.Dependencies, wantDeps) {
		t.Errorf("Dependencies = %v, want %v", summary.Dependencies, wantDeps)
	}

	wantDatabase := []string{"PostgreSQL", "Redis"}
	if !equalStrings(summary.Database, wantDatabase) {
		t.Errorf("Database = %v, want %v", summary.Database, wantDatabase)
	}

	wantInfra := []string{"Docker", "Docker Compose"}
	if !equalStrings(summary.Infrastructure, wantInfra) {
		t.Errorf("Infrastructure = %v, want %v", summary.Infrastructure, wantInfra)
	}

	if len(summary.CICD) != 1 || summary.CICD[0] != "CI/CD to EC2" {
		t.Errorf("CICD = %v, want [CI/CD to EC2]", summary.CICD)
	}

	wantProfiles := []string{"dev", "prod"}
	if !equalStrings(summary.Profiles, wantProfiles) {
		t.Errorf("Profiles = %v, want %v", summary.Profiles, wantProfiles)
	}

	text := summary.String()
	for _, fragment := range []string{"Spring Boot 3.5.7", "PostgreSQL, Redis", "Docker, Docker Compose"} {
		if !strings.Contains(text, fragment) {
			t.Errorf("summary text %q does not contain %q", text, fragment)
		}
	}
}

func TestBuildSummaryEmptyProject(t *testing.T) {

	summary := BuildSummary(&project.Info{})

	if summary.Framework != "" {
		t.Errorf("expected empty Framework, got %q", summary.Framework)
	}

	if len(summary.Dependencies) != 0 || len(summary.Database) != 0 ||
		len(summary.Infrastructure) != 0 || len(summary.CICD) != 0 || len(summary.Profiles) != 0 {
		t.Errorf("expected all empty slices for a zero-value project.Info, got %+v", summary)
	}

	if summary.String() != "" {
		t.Errorf("expected empty summary text, got %q", summary.String())
	}
}

func equalStrings(got, want []string) bool {

	if len(got) != len(want) {
		return false
	}

	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}

	return true
}
