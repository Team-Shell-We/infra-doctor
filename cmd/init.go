package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/fs"
	"github.com/Team-Shell-We/infra-doctor/internal/utils"
	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize Infra Doctor",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf(
					"failed to get current directory: %w",
					err,
				)
			}

			projectDir, err := findSpringBootProject(currentDir)
			if err != nil {
				return err
			}

			fmt.Println("Spring Boot project:", projectDir)

			result, err := utils.Initialize(projectDir)
			if err != nil {
				return err
			}

			for _, createdPath := range result.CreatedPaths {
				fmt.Println("Created:", createdPath)
			}

			return nil
		},
	}
}

func findSpringBootProject(currentDir string) (string, error) {
	// 현재 위치가 Spring Boot 프로젝트인지 먼저 확인
	if isSpringBootProject(currentDir) {
		return currentDir, nil
	}

	// 예제 디렉터리 내부를 재귀적으로 탐색
	examplesDir := filepath.Join(currentDir, "examples")

	if _, err := os.Stat(examplesDir); err != nil {
		return "", fmt.Errorf(
			"examples directory not found: %s",
			examplesDir,
		)
	}

	var foundDir string

	err := filepath.WalkDir(
		examplesDir,
		func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if !entry.IsDir() {
				return nil
			}

			// 불필요한 디렉터리는 탐색하지 않음
			switch entry.Name() {
			case ".git", ".gradle", ".idea", "build", "node_modules":
				if path != examplesDir {
					return filepath.SkipDir
				}
			}

			if isSpringBootProject(path) {
				foundDir = path
				return fs.SkipAll
			}

			return nil
		},
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to search examples directory: %w",
			err,
		)
	}

	if foundDir == "" {
		return "", fmt.Errorf(
			"Spring Boot project not found under %s; "+
				"build.gradle, build.gradle.kts, or pom.xml is required",
			examplesDir,
		)
	}

	return foundDir, nil
}

func isSpringBootProject(projectDir string) bool {
	buildFiles := []string{
		"build.gradle",
		"build.gradle.kts",
		"pom.xml",
	}

	for _, fileName := range buildFiles {
		filePath := filepath.Join(projectDir, fileName)

		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		buildContent := strings.ToLower(string(content))

		if strings.Contains(buildContent, "org.springframework.boot") ||
			strings.Contains(buildContent, "spring-boot") {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(initCommand())
}