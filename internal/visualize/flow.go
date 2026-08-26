package visualize

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func BuildDeploymentFlow(root string, info project.Info) (Diagram, error) {
	workflowText, err := workflowContents(
		root,
		info.Github.Workflows,
	)
	if err != nil {
		return Diagram{}, err
	}

	diagram := Diagram{
		Title: "Deployment Flow",
	}

	// 노드를 추가하면서 직전 노드와 연결한다.
	addStep := func(id, label string, kind NodeKind) {
		diagram.Nodes = append(
			diagram.Nodes,
			Node{
				ID:    id,
				Label: label,
				Kind:  kind,
			},
		)

		if len(diagram.Nodes) <= 1 {
			return
		}

		previous := diagram.Nodes[len(diagram.Nodes)-2]

		diagram.Edges = append(
			diagram.Edges,
			Edge{
				From: previous.ID,
				To:   id,
			},
		)
	}

	addStep(
		"developer",
		"Developer",
		Client,
	)
	addStep(
		"git-push",
		"Git Push",
		Pipeline,
	)

	if len(info.Github.Workflows) > 0 {
		addStep(
			"github-actions",
			"GitHub Actions",
			Pipeline,
		)
	}

	if build := detectedBuild(root, info); build != "" {
		addStep(
			"build",
			build+" Build",
			Pipeline,
		)
	}

	hasDocker := info.Infrastructure.Docker.Enabled ||
		fileExists(filepath.Join(root, "Dockerfile"))

	if hasDocker {
		addStep(
			"docker-image",
			"Docker Image",
			Container,
		)
	}

	if registry := detectedRegistry(workflowText); registry != "" {
		addStep(
			"registry",
			registry,
			Container,
		)
	}

	target := detectedDeployTarget(
		root,
		workflowText,
		info.Infrastructure.Kubernetes.Enabled,
	)

	if target != "" {
		addStep(
			"deploy-target",
			target,
			Orchestrator,
		)
	}

	// Kubernetes가 직접 컨테이너를 관리하므로 일반 Docker Container는 생략한다.
	if hasDocker && target != "Kubernetes" {
		addStep(
			"docker-container",
			"Docker Container",
			Container,
		)
	}

	addStep(
		"running-service",
		"Running Service",
		Application,
	)

	return diagram, nil
}

// workflowContents는 발견된 GitHub Actions 워크플로 파일을 모두 읽는다.
func workflowContents(
	root string,
	workflows []project.WorkflowInfo,
) (string, error) {
	var contents strings.Builder

	for _, workflow := range workflows {
		path := workflow.Path

		// Path가 없으면 GitHub Actions의 기본 위치를 사용한다.
		if path == "" {
			path = filepath.Join(
				root,
				".github",
				"workflows",
				workflow.File,
			)
		} else if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}

		contents.Write(data)
		contents.WriteByte('\n')
	}

	return strings.ToLower(contents.String()), nil
}

// detectedBuild : analyzer가 감지한 빌드 도구를 우선 쓰고, 없으면 대표 파일로 판단한다.
func detectedBuild(root string, info project.Info) string {
	if build := strings.TrimSpace(
		info.Framework.BuildTool.Type,
	); build != "" {
		return build
	}

	candidates := []struct {
		name  string
		label string
	}{
		{
			name:  "build.gradle",
			label: "Gradle",
		},
		{
			name:  "build.gradle.kts",
			label: "Gradle",
		},
		{
			name:  "pom.xml",
			label: "Maven",
		},
	}

	for _, candidate := range candidates {
		if fileExists(
			filepath.Join(root, candidate.name),
		) {
			return candidate.label
		}
	}

	return ""
}

func detectedRegistry(text string) string {
	switch {
	case strings.Contains(text, "amazon-ecr"),
		strings.Contains(text, ".dkr.ecr."):
		return "Amazon ECR"

	case strings.Contains(text, "ghcr.io"):
		return "GitHub Container Registry"

	case strings.Contains(text, "docker/login-action"),
		strings.Contains(text, "docker push"):
		return "Docker Registry"

	default:
		return ""
	}
}

func detectedDeployTarget(
	root string,
	text string,
	kubernetes bool,
) string {
	switch {
	case kubernetes,
		strings.Contains(text, "kubectl"),
		strings.Contains(text, "helm "),
		fileExists(filepath.Join(root, "Chart.yaml")):
		return "Kubernetes"

	case strings.Contains(text, "docker pull"),
		strings.Contains(text, "docker compose pull"):
		if strings.Contains(text, "ec2") ||
			strings.Contains(text, "appleboy/ssh-action") ||
			strings.Contains(text, "ssh ") {
			return "EC2 Pull"
		}
		return "Container Pull"

	case strings.Contains(text, "ec2"),
		strings.Contains(text, "appleboy/ssh-action"),
		strings.Contains(text, "ssh "):
		return "EC2"

	default:
		return ""
	}
}

func fileExists(path string) bool {
	stat, err := os.Stat(path)

	return err == nil && !stat.IsDir()
}
