package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type DeploymentRule struct{}

func (r DeploymentRule) Check(info *project.Info) []Diagnosis {

	var diagnoses []Diagnosis

	hasDocker := info.Infrastructure.Docker.Enabled
	hasGitHubActions := len(info.Github.Workflows) > 0

	switch {

	// Docker ❌ GitHub Actions ❌
	case !hasDocker && !hasGitHubActions:

		diagnoses = append(diagnoses, Diagnosis{
			Category: Infrastructure,
			Level:    Critical,

			ScoreImpact: -40,

			Title:   "Deployment environment is not prepared",
			Message: "Dockerfile and GitHub Actions workflow were not found.",

			Reason: "The project does not appear to have a standardized deployment process. Manual deployment increases the risk of configuration drift and deployment errors.",

			Fix: "Create a Dockerfile and configure a GitHub Actions workflow.",
		})

	// Docker ❌ GitHub Actions ✅
	case !hasDocker && hasGitHubActions:

		diagnoses = append(diagnoses, Diagnosis{
			Category: Infrastructure,
			Level:    Warning,

			ScoreImpact: -20,

			Title:   "Docker is not configured",
			Message: "GitHub Actions was detected, but no Dockerfile was found.",

			Reason: "CI/CD is configured, but deployment still depends on the server environment. Containerization provides consistent builds and deployments.",

			Fix: "Create a Dockerfile and update the workflow to build and deploy the Docker image.",
		})

	// Docker ✅ GitHub Actions ❌
	case hasDocker && !hasGitHubActions:

		diagnoses = append(diagnoses, Diagnosis{
			Category: CICD,
			Level:    Warning,

			ScoreImpact: -15,

			Title:   "CI/CD pipeline is not configured",
			Message: "Docker is configured, but no GitHub Actions workflow was found.",

			Reason: "Containerization is complete, but deployment still appears to be manual. CI/CD helps automate testing and deployment.",

			Fix: "Configure GitHub Actions to automate build and deployment.",
		})

	// Docker ✅ GitHub Actions ✅
	case hasDocker && hasGitHubActions:
		// No deployment issues detected.
	}

	return diagnoses
}
