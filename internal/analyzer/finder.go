package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
)

type BuildTool string

const (
	Gradle BuildTool = "Gradle"
	Maven  BuildTool = "Maven"
)

type BuildFile struct {
	Tool BuildTool
	Path string
}

// FindBuildFile: root에서 build.gradle(.kts) 또는 pom.xml을 찾는다
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
