package analyzer

import (
	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeProject(root string) (*project.Info, error) {

	info := &project.Info{}

	buildFile, err := FindBuildFile(root)
	if err != nil {
		return nil, err
	}

	switch buildFile.Tool {

	case Gradle:

		framework, database, err := AnalyzeGradle(buildFile.Path)
		if err != nil {
			return nil, err
		}

		info.Framework = *framework
		info.Database = *database

		profiles, err := FindProfiles(root)
		if err != nil {
			return nil, err
		}
		info.Profiles = profiles

		infrastructure, err := AnalyzeInfrastructure(root)
		if err != nil {
			return nil, err
		}
		info.Infrastructure = *infrastructure

		github, err := AnalyzeGitHub(root)
		if err != nil {
			return nil, err
		}
		info.Github = *github

	case Maven:
		// TODO: AnalyzeMaven(buildFile.Path)
	}

	return info, nil
}
