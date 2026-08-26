package visualize

import (
	"fmt"
	"strings"
)

func RenderMermaid(diagram Diagram) string {
	var builder strings.Builder

	builder.WriteString("flowchart TD\n\n")

	writeNodes(
		&builder,
		"Entry Point",
		diagram.Nodes,
		Client,
	)

	writeNodes(
		&builder,
		"Runtime",
		diagram.Nodes,
		Proxy,
		Application,
	)

	writeNodes(
		&builder,
		"Data Layer",
		diagram.Nodes,
		Database,
		Cache,
	)

	writeNodes(
		&builder,
		"Deployment",
		diagram.Nodes,
		Pipeline,
		Orchestrator,
		Container,
	)

	builder.WriteString("    %% Request Flow\n")

	for _, edge := range diagram.Edges {
		if isDeploymentEdge(edge) {
			continue
		}

		writeEdge(&builder, edge)
	}

	builder.WriteString("\n")
	builder.WriteString("    %% Deployment Flow\n")

	for _, edge := range diagram.Edges {
		if !isDeploymentEdge(edge) {
			continue
		}

		writeDeploymentEdge(&builder, edge)
	}

	return builder.String()
}

func writeNodes(
	builder *strings.Builder,
	section string,
	nodes []Node,
	kinds ...NodeKind,
) {
	filtered := filterNodes(nodes, kinds...)

	if len(filtered) == 0 {
		return
	}

	fmt.Fprintf(
		builder,
		"    %%%% %s\n",
		section,
	)

	for _, node := range filtered {
		fmt.Fprintf(
			builder,
			"    %s\n",
			mermaidNode(node),
		)
	}

	builder.WriteString("\n")
}

func filterNodes(
	nodes []Node,
	kinds ...NodeKind,
) []Node {
	allowed := make(map[NodeKind]bool)

	for _, kind := range kinds {
		allowed[kind] = true
	}

	var result []Node

	for _, node := range nodes {
		if allowed[node.Kind] {
			result = append(result, node)
		}
	}

	return result
}

func mermaidNode(node Node) string {
	label := escapeMermaid(node.Label)

	if node.Description != "" {
		label += "<br/>" +
			escapeMermaid(node.Description)
	}

	switch node.Kind {
	case Client:
		return fmt.Sprintf(
			"%s([\"%s\"])",
			node.ID,
			label,
		)

	case Database, Cache:
		return fmt.Sprintf(
			"%s[(\"%s\")]",
			node.ID,
			label,
		)

	default:
		return fmt.Sprintf(
			"%s[\"%s\"]",
			node.ID,
			label,
		)
	}
}

func writeEdge(
	builder *strings.Builder,
	edge Edge,
) {
	if edge.Label == "" {
		fmt.Fprintf(
			builder,
			"    %s --> %s\n",
			edge.From,
			edge.To,
		)

		return
	}

	fmt.Fprintf(
		builder,
		"    %s -->|%s| %s\n",
		edge.From,
		escapeMermaid(edge.Label),
		edge.To,
	)
}

func writeDeploymentEdge(
	builder *strings.Builder,
	edge Edge,
) {
	if edge.Label == "" {
		fmt.Fprintf(
			builder,
			"    %s -.-> %s\n",
			edge.From,
			edge.To,
		)

		return
	}

	fmt.Fprintf(
		builder,
		"    %s -.->|%s| %s\n",
		edge.From,
		escapeMermaid(edge.Label),
		edge.To,
	)
}

func isDeploymentEdge(edge Edge) bool {
	switch edge.Label {
	case "runs", "deploys", "build/deploy":
		return true
	default:
		return false
	}
}

func escapeMermaid(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"\"", "&quot;",
		"<", "&lt;",
		">", "&gt;",
		"\n", " ",
	)

	return replacer.Replace(value)
}
