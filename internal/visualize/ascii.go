package visualize

import (
	"strings"

	"github.com/Team-Shell-We/infra-doctor/internal/ui"
)

const (
	diagramWidth = 58
	nodeWidth    = 16
)

func RenderASCII(diagram Diagram) string {
	// 배포 흐름은 노드가 이미 실행 순서대로 저장되어 있다.
	if diagram.Title == "Deployment Flow" {
		return renderDeploymentASCII(diagram)
	}

	return renderArchitectureASCII(diagram)
}

func renderDeploymentASCII(diagram Diagram) string {
	lines := []string{
		ui.Center(diagram.Title, diagramWidth),
		strings.Repeat("-", diagramWidth),
		"",
	}

	for index, node := range diagram.Nodes {
		lines = append(
			lines,
			nodeBox(node, 2)...,
		)

		if index < len(diagram.Nodes)-1 {
			lines = append(
				lines,
				"        |",
				"        v",
			)
		}
	}

	lines = append(lines, "")

	return renderASCIIFrame(lines)
}

func renderArchitectureASCII(diagram Diagram) string {
	nodes := make(
		map[string]Node,
		len(diagram.Nodes),
	)

	for _, node := range diagram.Nodes {
		nodes[node.ID] = node
	}

	lines := []string{
		ui.Center(diagram.Title, diagramWidth),
		strings.Repeat("-", diagramWidth),
		"",
	}

	client, hasClient := firstNode(
		diagram.Nodes,
		Client,
	)
	if hasClient {
		lines = append(
			lines,
			"  "+client.Label,
		)
	}

	proxy, hasProxy := firstNode(
		diagram.Nodes,
		Proxy,
	)
	application, hasApplication := firstNode(
		diagram.Nodes,
		Application,
	)

	if hasProxy {
		lines = append(
			lines,
			"        |",
			"        v",
		)

		lines = append(
			lines,
			nodeBox(proxy, 2)...,
		)
	}

	if hasApplication {
		if hasClient || hasProxy {
			lines = append(
				lines,
				"        |",
				"        v",
			)
		}

		lines = append(
			lines,
			nodeBox(application, 2)...,
		)
	}

	stores := nodesOfKinds(
		diagram.Nodes,
		Database,
		Cache,
	)

	if len(stores) > 0 {
		lines = append(
			lines,
			renderStores(stores)...,
		)
	}

	lines = append(lines, "")

	return renderASCIIFrame(lines)
}

func renderASCIIFrame(lines []string) string {
	var result strings.Builder

	result.WriteString(
		"+" +
			strings.Repeat("-", diagramWidth) +
			"+\n",
	)

	for _, line := range lines {
		result.WriteString(
			"|" +
				ui.PadRight(line, diagramWidth) +
				"|\n",
		)
	}

	result.WriteString(
		"+" +
			strings.Repeat("-", diagramWidth) +
			"+\n",
	)

	return result.String()
}

func nodeBox(node Node, indent int) []string {
	prefix := strings.Repeat(" ", indent)

	border := prefix +
		"+" +
		strings.Repeat("-", nodeWidth) +
		"+"

	content := prefix +
		"|" +
		ui.Center(node.Label, nodeWidth) +
		"|"

	if node.Description != "" {
		content += "  <- " + node.Description
	}

	return []string{
		border,
		content,
		border,
	}
}

func renderStores(stores []Node) []string {
	if len(stores) == 1 {
		store := stores[0]

		return []string{
			"        |",
			"        v",
			"  " + store.Label,
			"  " + store.Description,
		}
	}

	left := stores[0]
	right := stores[1]

	lines := []string{
		"        |",
		"        +-------------------+",
		"        |                   |",
		"        v                   v",
		columns(left.Label, right.Label),
		columns(left.Description, right.Description),
	}

	for _, store := range stores[2:] {
		lines = append(
			lines,
			"",
			"  "+store.Label+"  <- "+store.Description,
		)
	}

	return lines
}

func columns(left, right string) string {
	return "  " +
		ui.PadRight(left, 20) +
		"  " +
		right
}

func firstNode(
	nodes []Node,
	kind NodeKind,
) (Node, bool) {
	for _, node := range nodes {
		if node.Kind == kind {
			return node, true
		}
	}

	return Node{}, false
}

func nodesOfKinds(
	nodes []Node,
	kinds ...NodeKind,
) []Node {
	allowed := make(
		map[NodeKind]bool,
		len(kinds),
	)

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
