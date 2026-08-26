package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type BuildTool string

const (
	BuildToolUnknown BuildTool = "unknown"
	BuildToolMaven   BuildTool = "maven"
	BuildToolGradle  BuildTool = "gradle" // Gradle or Gradle Kotlin DSL
)

type DetectionResult struct {
	Detected  bool
	BuildTool BuildTool
	BuildFile string
}

func DetectSpringBoot(dir string) (*DetectionResult, error) {

	if dir == "" {
		return nil, errors.New("project directory is empty")
	}

	info, err := os.Stat(dir)
	if err != nil {
		return nil, errors.New("project directory does not exist")
	}

	if !info.IsDir() {
		return &DetectionResult{}, fmt.Errorf(
			"project directory %q is not a directory",
			dir,
		)
	}

	candidateFiles := []struct {
		FileName  string
		BuildTool BuildTool
		Markers   []string
	}{
		{
			FileName:  "build.gradle",
			BuildTool: BuildToolGradle,
			Markers: []string{
				"org.springframework.boot",
				"spring-boot-starter",
			},
		},
		{
			FileName:  "build.gradle.kts",
			BuildTool: BuildToolGradle,
			Markers: []string{
				"org.springframework.boot",
				"spring-boot-starter",
			},
		},
		{
			FileName:  "pom.xml",
			BuildTool: BuildToolMaven,
			Markers: []string{
				"org.springframework.boot",
				"spring-boot-starter",
				"spring-boot-maven-plugin",
			},
		},
	}

	for _, candidate := range candidateFiles {
		buildFilePath := filepath.Join(dir, candidate.FileName)

		content, err := os.ReadFile(buildFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf(
				"cannot read build file %q: %w",
				buildFilePath,
				err,
			)
		}

		if containsAnyMarker(string(content), candidate.Markers) {
			return &DetectionResult{
				Detected:  true,
				BuildTool: candidate.BuildTool,
				BuildFile: buildFilePath,
			}, nil
		}
	}

	return &DetectionResult{
		Detected:  false,
		BuildTool: BuildToolUnknown,
	}, nil
}

func containsAnyMarker(content string, markers []string) bool {
	lowerContent := strings.ToLower(content)

	for _, marker := range markers {
		if strings.Contains(lowerContent, strings.ToLower(marker)) {
			return true
		}
	}

	return false
}
