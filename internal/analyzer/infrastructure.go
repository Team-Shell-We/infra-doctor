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

		case "Chart.yaml",
			"deployment.yaml",
			"deployment.yml",
			"service.yaml",
			"service.yml",
			"ingress.yaml",
			"ingress.yml",
			"kustomization.yaml":

			info.Kubernetes.Enabled = true

			info.Kubernetes.Files = append(
				info.Kubernetes.Files,
				project.KubernetesFileInfo{
					File: fileInfo.Name(),
					Path: path,
				},
			)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}
