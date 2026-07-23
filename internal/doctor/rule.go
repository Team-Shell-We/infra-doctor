package doctor

import "github.com/Team-Shell-We/infra-doctor/internal/project"

type Rule interface {
	Check(info *project.Info) []Diagnosis
}
