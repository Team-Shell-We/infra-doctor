package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type DeploymentRule struct{}

func (r DeploymentRule) Check(info *project.Info) []Diagnosis {

	var diagnoses []Diagnosis

	registry, err := LoadRules()
	if err != nil {
		return diagnoses
	}

	hasDocker := info.Infrastructure.Docker.Enabled
	hasGitHubActions := len(info.Github.Workflows) > 0

	switch {

	case !hasDocker && !hasGitHubActions:

		if rule, err := registry.DeploymentRule("no_deployment"); err == nil {
			diagnoses = append(diagnoses, rule)
		}

	case !hasDocker && hasGitHubActions:

		if rule, err := registry.DeploymentRule("no_docker"); err == nil {
			diagnoses = append(diagnoses, rule)
		}

	case hasDocker && !hasGitHubActions:

		if rule, err := registry.DeploymentRule("no_github_actions"); err == nil {
			diagnoses = append(diagnoses, rule)
		}

	case hasDocker && hasGitHubActions:
		// Deployment configuration looks good.
	}

	return diagnoses
}