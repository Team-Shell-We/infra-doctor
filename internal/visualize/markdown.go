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

func markdownCell(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")

	return value
}
