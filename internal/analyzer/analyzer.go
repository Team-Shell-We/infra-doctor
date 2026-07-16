package analyzer

import (
	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeProject(root string) (*project.Info, error) {

	info := &project.Info{}

	// Repository에서 build.gradle 또는 pom.xml을 찾음
	buildFile, err := FindBuildFile(root)
	if err != nil {
		return nil, err
	}

	switch buildFile.Tool {

	case Gradle:

		framework, err := AnalyzeGradle(buildFile.Path)
		if err != nil {
			return nil, err
		}

		info.Framework = *framework
		
		profiles, err := FindProfiles(root)
		if err != nil {
			return nil, err
		}

		info.Profiles = profiles

	case Maven:
		// TODO: AnalyzeMaven(buildFile.Path)
	}

	return info, nil
}