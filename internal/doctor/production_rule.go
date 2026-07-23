package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type ProductionRule struct{}

func (r ProductionRule) Check(info *project.Info) []Diagnosis {

	var diagnoses []Diagnosis

	registry, err := LoadRules()
	if err != nil {
		return diagnoses
	}

	checks := []struct {
		id      string
		missing bool
	}{
		{"no_health_check", !info.Infrastructure.HealthCheck.Enabled},
		{"no_reverse_proxy", !info.Infrastructure.Nginx.Enabled},
		{"no_monitoring", !info.Infrastructure.Monitoring.Enabled},
		{"no_log_rotation", !info.Infrastructure.LogRotation.Enabled},
		{"no_backup", !info.Infrastructure.Backup.Enabled},
	}

	for _, c := range checks {

		if !c.missing {
			continue
		}

		if rule, err := registry.ProductionRule(c.id); err == nil {
			diagnoses = append(diagnoses, rule)
		}
	}

	return diagnoses
}
