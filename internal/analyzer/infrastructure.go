package analyzer

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

var replicasRegex = regexp.MustCompile(`(?m)^\s*replicas:\s*(\d+)`)

func AnalyzeInfrastructure(root string) (*project.InfrastructureInfo, error) {

	info := &project.InfrastructureInfo{}

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

		switch fileInfo.Name() {

		case "Dockerfile":

			info.Docker.Enabled = true

			info.Docker.Dockerfiles = append(
				info.Docker.Dockerfiles,
				project.DockerfileInfo{
					File: fileInfo.Name(),
					Path: path,
				},
			)

			if data, readErr := os.ReadFile(path); readErr == nil &&
				strings.Contains(string(data), "HEALTHCHECK") {

				info.HealthCheck.Enabled = true
			}

		case "docker-compose.yml",
			"docker-compose.yaml":

			info.Docker.Compose = append(
				info.Docker.Compose,
				project.ComposeInfo{
					File: fileInfo.Name(),
					Path: path,
				},
			)

			if data, readErr := os.ReadFile(path); readErr == nil &&
				strings.Contains(strings.ToLower(string(data)), "healthcheck:") {

				info.HealthCheck.Enabled = true
			}

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

			if data, readErr := os.ReadFile(path); readErr == nil {
				if match := replicasRegex.FindSubmatch(data); match != nil {
					if n, convErr := strconv.Atoi(string(match[1])); convErr == nil && n > info.Kubernetes.Replicas {
						info.Kubernetes.Replicas = n
					}
				}
			}

		case "nginx.conf":

			info.Nginx.Enabled = true

		case "prometheus.yml",
			"prometheus.yaml",
			"grafana.ini":

			info.Monitoring.Enabled = true

		case "logrotate.conf":

			info.LogRotation.Enabled = true

		case "backup.sh":

			info.Backup.Enabled = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}
