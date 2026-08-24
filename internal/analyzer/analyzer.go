package analyzer

import (
	"fmt"
	"os"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeProject(root string) (*project.Info, error) {

	// root가 아예 존재하지 않는 경우와, root는 있지만 Spring Boot
	// 프로젝트가 아닌 경우(FindBuildFile의 "no build file found")를
	// 구분해서 알려준다 — 안 그러면 예: `recommend k8s`처럼 인자를
	// 실수로 경로로 착각하고 입력했을 때 원인 파악이 어렵다.
	stat, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", root)
		}
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}

	info := &project.Info{}

	buildFile, err := FindBuildFile(root)
	if err != nil {
		return nil, err
	}

	switch buildFile.Tool {

	case Gradle:

		framework, dependency, database, err := AnalyzeGradle(buildFile.Path)
		if err != nil {
			return nil, err
		}

		info.Framework = *framework
		info.Dependencies = *dependency
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
		framework, dependency, database, err:=AnalyzeMaven(buildFile.Path)
		
		if err!=nil{
			return nil, err
		}

		info.Framework=*framework
		info.Dependencies=*dependency
		info.Database=*database
	}

	return info, nil
}
