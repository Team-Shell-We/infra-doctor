package analyzer

import (
	"os"
	"path/filepath"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func AnalyzeInfrastructure(root string) (*project.InfrastructureInfo, error) {

	info := &project.InfrastructureInfo{}

	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// --------------------------------------------------------
		// Directory
		// --------------------------------------------------------

		if fileInfo.IsDir() {

			// TODO: Kubernetes 디렉터리 탐지
			// ex) k8s/, helm/

			return nil
		}

		// --------------------------------------------------------
		// Infrastructure Files
		// --------------------------------------------------------

		switch fileInfo.Name() {

		// ----------------------------
		// Docker
		// ----------------------------

		case "Dockerfile":

			info.Docker.Enabled = true

			info.Docker.Dockerfiles = append(
				info.Docker.Dockerfiles,
				project.DockerfileInfo{
					File: fileInfo.Name(),
					Path: path,
				},
			)

		// ----------------------------
		// Docker Compose
		// ----------------------------

		case "docker-compose.yml",
			"docker-compose.yaml":

			info.Docker.Compose = append(
				info.Docker.Compose,
				project.ComposeInfo{
					File: fileInfo.Name(),
					Path: path,
				},
			)

		// ----------------------------
		// Kubernetes
		// ----------------------------

		// TODO: 다음 단계에서 구현
		// case "Chart.yaml":
		// case "deployment.yaml":
		// case "service.yaml":
		// case "ingress.yaml":
		// case "kustomization.yaml":
		//     info.Kubernetes.Enabled = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}