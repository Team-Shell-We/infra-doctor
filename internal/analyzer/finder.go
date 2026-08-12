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

// FindBuildFile : root에서 build.gradle 또는 pom.xml을 찾음(하위 디렉터리
// 재귀 탐색 안 함 — 재귀 탐색하면 examples/ 등 하위의 다른 프로젝트를
// 잘못 집어올 수 있어서 제거)
func FindBuildFile(root string) (*BuildFile, error) {

	candidates := []struct {
		FileName string
		Tool     BuildTool
	}{
		{"build.gradle", Gradle},
		{"build.gradle.kts", Gradle},
		{"pom.xml", Maven},
	}

	for _, candidate := range candidates {

		path := filepath.Join(root, candidate.FileName)

		if _, err := os.Stat(path); err == nil {
			return &BuildFile{Tool: candidate.Tool, Path: path}, nil
		}
	}

	return nil, fmt.Errorf("no build file found")
}
