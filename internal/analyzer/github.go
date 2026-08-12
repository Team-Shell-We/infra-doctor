package analyzer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeGitHub(root string) (*project.GithubInfo, error) {

	info := &project.GithubInfo{}

	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			if shouldSkipDir(fileInfo.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if !strings.Contains(rel, ".github") ||
			!strings.Contains(rel, "workflows") {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}

		workflow, err := ParseWorkflow(path, fileInfo.Name())
		if err != nil {
			return err
		}

		info.Workflows = append(info.Workflows, *workflow)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}
