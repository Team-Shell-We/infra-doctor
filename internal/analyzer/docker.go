package analyzer

import (
	"os"
	"path/filepath"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeDocker(root string) (*project.DockerInfo, error) {

	info := &project.DockerInfo{}

	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		switch fileInfo.Name() {

		case "Dockerfile":

			// Docker를 사용한다고 판단
			info.Enabled = true

			info.Dockerfiles = append(info.Dockerfiles, project.DockerfileInfo{
				File: fileInfo.Name(),
				Path: path,
			})

		case "docker-compose.yml", "docker-compose.yaml":

			info.Compose = append(info.Compose, project.ComposeInfo{
				File: fileInfo.Name(),
				Path: path,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}