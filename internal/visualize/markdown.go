package visualize

import (
	"fmt"
	"strings"
)

func RenderMarkdown(d Diagram) string {
	var b strings.Builder

	fmt.Fprintf(
		&b,
		"# %s\n\n",
		d.Title,
	)

	b.WriteString("## Architecture Diagram\n\n")
	b.WriteString("```text\n")
	b.WriteString(RenderASCII(d))
	b.WriteString("```\n\n")

	b.WriteString("## Components\n\n")
	b.WriteString("| Component | Type | Description |\n")
	b.WriteString("|---|---|---|\n")

	for _, node := range d.Nodes {
		description := node.Description
		if description == "" {
			description = "-"
		}

		fmt.Fprintf(
			&b,
			"| %s | %s | %s |\n",
			markdownCell(node.Label),
			markdownCell(string(node.Kind)),
			markdownCell(description),
		)
	}

	return b.String()
}

func nodeMap(nodes []Node) map[string]Node {
	result := make(map[string]Node, len(nodes))

	for _, node := range nodes {
		result[node.ID] = node
	}

	return result
}

func nodeLabel(
	nodes map[string]Node,
	id string,
) string {
	node, found := nodes[id]
	if !found {
		return id
	}

	return node.Label
}

func markdownCell(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")

	return value
}