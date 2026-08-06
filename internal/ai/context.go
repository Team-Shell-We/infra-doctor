package ai

import (
	"fmt"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

// Summary is a compact, prompt-ready view of a scanned project. AI commands
// send this instead of raw source code: it is cheaper, grounds the model in
// facts the CLI already verified deterministically, and keeps every AI
// command answering from the same known-accurate picture of the project.
type Summary struct {
	Framework      string
	Dependencies   []string
	Database       []string
	Infrastructure []string
	CICD           []string
	Profiles       []string
}

// BuildSummary turns a scanned project.Info into a Summary. It is a pure
// function of its input so prompt construction can be unit tested without
// re-running the analyzer.
func BuildSummary(info *project.Info) Summary {

	summary := Summary{}

	if info.Framework.SpringBoot.Enabled {
		summary.Framework = fmt.Sprintf("Spring Boot %s", info.Framework.SpringBoot.Version)
	}

	if info.Framework.BuildTool.Type != "" {
		summary.Framework = strings.TrimSpace(summary.Framework + " (" + info.Framework.BuildTool.Type + ")")
	}

	if info.Framework.Java.Version != "" {
		summary.Framework = strings.TrimSpace(fmt.Sprintf("%s, Java %s", summary.Framework, info.Framework.Java.Version))
	}

	deps := map[string]bool{
		"Spring Security":      info.Dependencies.Security.Enabled,
		"Spring Data JPA":      info.Dependencies.JPA.Enabled,
		"Kafka":                info.Dependencies.Kafka.Enabled,
		"AWS SDK":              info.Dependencies.AWS.Enabled,
		"Lombok":               info.Dependencies.Lombok.Enabled,
		"Spring Boot Actuator": info.Dependencies.Actuator.Enabled,
		"SpringDoc OpenAPI":    info.Dependencies.OpenAPI.Enabled,
	}
	summary.Dependencies = enabledNames(deps)

	if info.Database.Primary.Type != "" && info.Database.Primary.Type != "Unknown" {
		summary.Database = append(summary.Database, info.Database.Primary.Type)
	}

	if info.Database.Redis != nil {
		summary.Database = append(summary.Database, "Redis")
	}

	infra := map[string]bool{
		"Docker":         info.Infrastructure.Docker.Enabled,
		"Docker Compose": len(info.Infrastructure.Docker.Compose) > 0,
		"Kubernetes":     info.Infrastructure.Kubernetes.Enabled,
		"Nginx":          info.Infrastructure.Nginx.Enabled,
		"Health Check":   info.Infrastructure.HealthCheck.Enabled,
		"Monitoring":     info.Infrastructure.Monitoring.Enabled,
		"Log Rotation":   info.Infrastructure.LogRotation.Enabled,
		"DB Backup":      info.Infrastructure.Backup.Enabled,
	}
	summary.Infrastructure = enabledNames(infra)

	for _, workflow := range info.Github.Workflows {
		summary.CICD = append(summary.CICD, workflow.Name)
	}

	for _, profile := range info.Profiles {
		summary.Profiles = append(summary.Profiles, profile.Name)
	}

	return summary
}

// String renders the summary as a compact bullet list suitable for
// embedding directly into a prompt.
func (s Summary) String() string {

	var b strings.Builder

	writeLine := func(label string, values []string) {
		if len(values) == 0 {
			return
		}
		fmt.Fprintf(&b, "%s: %s\n", label, strings.Join(values, ", "))
	}

	if s.Framework != "" {
		fmt.Fprintf(&b, "Framework: %s\n", s.Framework)
	}

	writeLine("Dependencies", s.Dependencies)
	writeLine("Database", s.Database)
	writeLine("Infrastructure", s.Infrastructure)
	writeLine("CI/CD", s.CICD)
	writeLine("Profiles", s.Profiles)

	return strings.TrimSpace(b.String())
}

func enabledNames(flags map[string]bool) []string {

	// Fixed, deterministic order so prompt text (and therefore model
	// output) doesn't vary run-to-run because of Go's random map order.
	order := []string{
		"Spring Security", "Spring Data JPA", "Kafka", "AWS SDK", "Lombok",
		"Spring Boot Actuator", "SpringDoc OpenAPI",
		"Docker", "Docker Compose", "Kubernetes", "Nginx",
		"Health Check", "Monitoring", "Log Rotation", "DB Backup",
	}

	var names []string

	for _, name := range order {
		if enabled, ok := flags[name]; ok && enabled {
			names = append(names, name)
		}
	}

	return names
}
