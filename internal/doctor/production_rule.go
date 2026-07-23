package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type ProductionRule struct{}

func (r ProductionRule) Check(info *project.Info) []Diagnosis {

	var diagnoses []Diagnosis

	return diagnoses
}
