package analyzer

import (
	"fmt"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeProject(root string) (*project.Info, error) {

	info := &project.Info{}

	// Repository에서 build.gradle 또는 pom.xml을 찾음
	buildFile, err := FindBuildFile(root)
	if err != nil {
		return nil, err
	}

	// 어떤 빌드 도구를 사용하는지 출력
	fmt.Printf("✓ %s detected\n", buildFile.Tool)

	// 찾은 build 파일 경로를 출력
	fmt.Printf("✓ Build File : %s\n", buildFile.Path)

	return info, nil
}