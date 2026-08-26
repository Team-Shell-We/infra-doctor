package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

func Analyze(info *project.Info) *Result {

	result := &Result{
		Score: 100,
	}

	rules := []Rule{
		DeploymentRule{},
		LocalDevelopmentRule{},
		ProductionRule{},
	}

	for _, rule := range rules {

		diagnoses := rule.Check(info)

		for _, diagnosis := range diagnoses {

			result.Score += diagnosis.ScoreImpact
			result.Diagnoses = append(result.Diagnoses, diagnosis)
		}
	}

	if result.Score < 0 {
		result.Score = 0
	}

	result.Checks = Checklist(info)

	return result
}
