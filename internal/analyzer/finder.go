package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
)

// BuildTool은 빌드 도구 종류를 나타낸다.
type BuildTool string

const (
	Gradle BuildTool = "Gradle"
	Maven  BuildTool = "Maven"
)

// BuildFile : 찾은 빌드 파일 정보 저장 구조체
type BuildFile struct {
	Tool BuildTool
	Path string
}

// Repository를 재귀 탐색하여 build.gradle 또는 pom.xml을 찾음
func FindBuildFile(root string) (*BuildFile, error) {

	var result *BuildFile

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// 폴더는 건너뜀
		if info.IsDir() {
			return nil
		}

		switch info.Name() {

		case "build.gradle", "build.gradle.kts":
			result = &BuildFile{
				Tool: Gradle,
				Path: path,
			}
			return filepath.SkipAll

		case "pom.xml":
			result = &BuildFile{
				Tool: Maven,
				Path: path,
			}
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("no build file found")
	}

	return result, nil
}