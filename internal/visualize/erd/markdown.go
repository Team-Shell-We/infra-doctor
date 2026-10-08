package erd

import (
	"fmt"
	"strings"
)

func RenderMarkdown(d Diagram) string {

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", d.Title)

	for _, entity := range d.Entities {
		fmt.Fprintf(&b, "## %s\n\n", entity.Name)
		b.WriteString("| Column | Type | Key |\n")
		b.WriteString("|---|---|---|\n")
		for _, column := range entity.Columns {
			fmt.Fprintf(&b, "| %s | %s | %s |\n",
				markdownCell(column.Name),
				markdownCell(typeOrDefault(column.Type)),
				keyLabel(column),
			)
		}
		b.WriteString("\n")
	}

	if len(d.Relations) > 0 {
		b.WriteString("## Relations\n\n")
		for _, relation := range d.Relations {
			fmt.Fprintf(&b, "- `%s` %s `%s`", relation.From, relation.Cardinality, relation.To)
			if relation.Label != "" {
				fmt.Fprintf(&b, " (%s)", relation.Label)
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func markdownCell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
