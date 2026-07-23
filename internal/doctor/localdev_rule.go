package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type LocalDevelopmentRule struct{}

func (r LocalDevelopmentRule) Check(info *project.Info) []Diagnosis {

	var diagnoses []Diagnosis

	registry, err := LoadRules()
	if err != nil {
		return diagnoses
	}

	hasDependency := (info.Database.Primary.Type != "" && info.Database.Primary.Type != "Unknown") ||
		info.Database.Redis != nil

	if hasDependency && len(info.Infrastructure.Docker.Compose) == 0 {

		if rule, err := registry.LocalDevRule("no_compose_with_dependencies"); err == nil {
			diagnoses = append(diagnoses, rule)
		}
	}

	hasDevProfile := false

	for _, profile := range info.Profiles {
		if profile.Name == "dev" {
			hasDevProfile = true
			break
		}
	}

	if !hasDevProfile {

		if rule, err := registry.LocalDevRule("no_dev_profile"); err == nil {
			diagnoses = append(diagnoses, rule)
		}
	}

	return diagnoses
}
