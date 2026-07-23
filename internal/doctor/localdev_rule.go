package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type LocalDevelopmentRule struct{}

func (r LocalDevelopmentRule) Check(info *project.Info) []Diagnosis {

	var diagnoses []Diagnosis

	return diagnoses
}
