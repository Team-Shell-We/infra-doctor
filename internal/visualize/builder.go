package visualize

import (
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/project"
)

func Build(info project.Info) Diagram {
	diagram := Diagram{
		Title: "Infrastructure Overview",
	}

	addNode(&diagram, Node{
		ID:    "client",
		Label: "Client",
		Kind:  Client,
	})

	entry := "client"

	if info.Infrastructure.Nginx.Enabled {
		addNode(&diagram, Node{
			ID:          "nginx",
			Label:       "Nginx",
			Description: "Reverse Proxy",
			Kind:        Proxy,
		})

		addEdge(&diagram, Edge{
			From: "client",
			To:   "nginx",
		})

		entry = "nginx"
	}

	appLabel := "Application"

	if info.Framework.SpringBoot.Enabled {
		appLabel = "Spring Boot"
	}

	addNode(&diagram, Node{
		ID:          "application",
		Label:       appLabel,
		Description: "Business Logic",
		Kind:        Application,
	})

	addEdge(&diagram, Edge{
		From: entry,
		To:   "application",
	})

	database := strings.TrimSpace(
		info.Database.Primary.Type,
	)

	if database != "" && database != "Unknown" {
		databaseID := slug(database)

		addNode(&diagram, Node{
			ID:          databaseID,
			Label:       database,
			Description: "Persistent Storage",
			Kind:        Database,
		})

		addEdge(&diagram, Edge{
			From: "application",
			To:   databaseID,
		})
	}

	if info.Database.Redis != nil &&
		info.Database.Redis.Enabled {

		addNode(&diagram, Node{
			ID:          "redis",
			Label:       "Redis",
			Description: "Cache",
			Kind:        Cache,
		})

		addEdge(&diagram, Edge{
			From: "application",
			To:   "redis",
		})
	}

	if info.Infrastructure.Docker.Enabled {
		addNode(&diagram, Node{
			ID:    "container",
			Label: "Docker Container",
			Kind:  Container,
		})

		addEdge(&diagram, Edge{
			From:  "container",
			To:    "application",
			Label: "runs",
		})
	}

	if info.Infrastructure.Kubernetes.Enabled {
		addNode(&diagram, Node{
			ID:    "kubernetes",
			Label: "Kubernetes",
			Kind:  Orchestrator,
		})

		target := "application"

		if info.Infrastructure.Docker.Enabled {
			target = "container"
		}

		addEdge(&diagram, Edge{
			From:  "kubernetes",
			To:    target,
			Label: "deploys",
		})
	}

	if len(info.Github.Workflows) > 0 {
		addNode(&diagram, Node{
			ID:          "github-actions",
			Label:       "GitHub Actions",
			Description: "CI/CD",
			Kind:        Pipeline,
		})

		target := "application"

		if info.Infrastructure.Kubernetes.Enabled {
			target = "kubernetes"
		} else if info.Infrastructure.Docker.Enabled {
			target = "container"
		}

		addEdge(&diagram, Edge{
			From:  "github-actions",
			To:    target,
			Label: "build/deploy",
		})
	}

	return diagram
}

func FilterFlow(diagram Diagram) Diagram {
	allowedKinds := map[NodeKind]bool{
		Client:      true,
		Proxy:       true,
		Application: true,
		Database:    true,
		Cache:       true,
	}

	result := Diagram{
		Title: "Request Flow",
	}

	includedNodeIDs := make(map[string]bool)

	for _, node := range diagram.Nodes {
		if allowedKinds[node.Kind] {
			result.Nodes = append(
				result.Nodes,
				node,
			)

			includedNodeIDs[node.ID] = true
		}
	}

	for _, edge := range diagram.Edges {
		fromIncluded := includedNodeIDs[edge.From]
		toIncluded := includedNodeIDs[edge.To]

		if fromIncluded && toIncluded {
			result.Edges = append(
				result.Edges,
				edge,
			)
		}
	}

	return result
}

func addNode(diagram *Diagram, node Node) {
	for _, existing := range diagram.Nodes {
		if existing.ID == node.ID {
			return
		}
	}

	diagram.Nodes = append(
		diagram.Nodes,
		node,
	)
}

func addEdge(diagram *Diagram, edge Edge) {
	diagram.Edges = append(
		diagram.Edges,
		edge,
	)
}

func slug(value string) string {
	value = strings.ToLower(value)

	var builder strings.Builder

	for _, character := range value {
		isLetter := character >= 'a' &&
			character <= 'z'

		isNumber := character >= '0' &&
			character <= '9'

		if isLetter || isNumber {
			builder.WriteRune(character)
			continue
		}

		if builder.Len() > 0 {
			builder.WriteByte('-')
		}
	}

	return strings.Trim(
		builder.String(),
		"-",
	)
}
